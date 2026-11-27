package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/bettersandeep/yantra/internal/referee"
	"github.com/bettersandeep/yantra/internal/ticket"
)

const usage = `yantra - a factory that makes software factories

usage: yantra <command>

commands:
  version                                   print the version
  referee diff --base <rev> [--head <rev>]  report test weakening between two commits
  ticket show <number>                      print one ticket and say if it is ready
  ticket list                               print the open tickets
`

func main() {
	args := os.Args[1:]
	switch {
	case len(args) == 1 && args[0] == "version":
		fmt.Println("yantra 0.0.0")
	case len(args) >= 2 && args[0] == "referee" && args[1] == "diff":
		os.Exit(runRefereeDiff(args[2:]))
	case len(args) >= 2 && args[0] == "ticket" && args[1] == "show":
		os.Exit(runTicketShow(args[2:]))
	case len(args) == 2 && args[0] == "ticket" && args[1] == "list":
		os.Exit(runTicketList())
	default:
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
}

// runRefereeDiff exits 0 when clean, 1 on findings, 2 when the referee could not run.
func runRefereeDiff(args []string) int {
	fl := flag.NewFlagSet("referee diff", flag.ContinueOnError)
	base := fl.String("base", "", "base revision")
	head := fl.String("head", "HEAD", "head revision")
	if err := fl.Parse(args); err != nil {
		return 2
	}
	if *base == "" {
		fmt.Fprintln(os.Stderr, "--base is required")
		return 2
	}
	findings, err := refereeDiff(*base, *head)
	if err != nil {
		// Fail closed: a referee that cannot run is not a pass.
		fmt.Fprintln(os.Stderr, "referee:", err)
		return 2
	}
	if err := printJSON(findings); err != nil {
		fmt.Fprintln(os.Stderr, "referee:", err)
		return 2
	}
	if len(findings) > 0 {
		return 1
	}
	return 0
}

// resolve turns a caller-supplied revision into a commit ID, so that no later
// git command can read it as an option.
func resolve(rev string) (string, error) {
	out, err := exec.Command("git", "rev-parse", "--verify", "--quiet", "--end-of-options", rev+"^{commit}").Output()
	if err != nil {
		return "", fmt.Errorf("not a commit: %q", rev)
	}
	return strings.TrimSpace(string(out)), nil
}

func refereeDiff(baseRev, headRev string) ([]referee.Finding, error) {
	base, err := resolve(baseRev)
	if err != nil {
		return nil, err
	}
	head, err := resolve(headRev)
	if err != nil {
		return nil, err
	}
	// -z: paths arrive as raw bytes. Without it git quotes unusual names and
	// the _test.go suffix check does not match.
	out, err := exec.Command("git", "diff", "--name-status", "-M", "-z", base, head).Output()
	if err != nil {
		return nil, fmt.Errorf("git diff %s %s: %w", base, head, err)
	}
	findings := []referee.Finding{}
	f := strings.Split(strings.TrimSuffix(string(out), "\x00"), "\x00")
	for i := 0; i < len(f) && f[i] != ""; {
		status, paths := f[i], 1
		if status[0] == 'R' || status[0] == 'C' {
			paths = 2 // old path, new path
		}
		if i+paths >= len(f) {
			return nil, fmt.Errorf("git diff: unexpected output")
		}
		oldPath, newPath := f[i+1], f[i+paths]
		i += paths + 1
		if !strings.HasSuffix(oldPath, "_test.go") && !strings.HasSuffix(newPath, "_test.go") {
			continue
		}
		var b, h []byte
		if status[0] != 'A' {
			if b, err = show(base, oldPath); err != nil {
				return nil, err
			}
		}
		if status[0] != 'D' {
			if h, err = show(head, newPath); err != nil {
				return nil, err
			}
		}
		fs, err := referee.Diff(newPath, b, h)
		if err != nil {
			return nil, err
		}
		findings = append(findings, fs...)
	}
	return findings, nil
}

func show(rev, path string) ([]byte, error) {
	out, err := exec.Command("git", "show", rev+":"+path).Output()
	if err != nil {
		return nil, fmt.Errorf("git show %s:%s: %w", rev, path, err)
	}
	return out, nil
}

func printJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc.Encode(v)
}

// issue is the part of the gh JSON output that a ticket needs.
type issue struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	State  string `json:"state"`
	Body   string `json:"body"`
}

const issueFields = "number,title,state,body"

// ticketListLimit is the most issues that one list call reads. A list that
// reaches it is treated as incomplete.
const ticketListLimit = 1000

// valid reports whether gh returned the fields that a ticket needs.
func (is *issue) valid() bool {
	return is != nil && is.Number > 0 && is.Title != "" && is.State != ""
}

// runTicketShow exits 0 when the ticket is ready, 1 when it is not ready, 2
// when the ticket could not be read.
func runTicketShow(args []string) int {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: yantra ticket show <number>")
		return 2
	}
	// Only the parsed integer reaches gh, never the caller's text.
	n, err := strconv.Atoi(args[0])
	if err != nil || n <= 0 {
		fmt.Fprintf(os.Stderr, "ticket: not a positive issue number: %q\n", args[0])
		return 2
	}
	out, err := exec.Command("gh", "issue", "view", strconv.Itoa(n), "--json", issueFields).Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ticket: gh issue view %d: %v\n", n, err)
		return 2
	}
	var is issue
	if err := json.Unmarshal(out, &is); err != nil || !is.valid() || is.Number != n {
		fmt.Fprintf(os.Stderr, "ticket: unexpected gh output for issue %d\n", n)
		return 2
	}
	t := ticket.New(is.Number, is.Title, is.State, is.Body)
	if err := printJSON(t); err != nil {
		fmt.Fprintln(os.Stderr, "ticket:", err)
		return 2
	}
	if !t.Ready {
		return 1
	}
	return 0
}

// runTicketList prints the open issues with the ticket label. It exits 1
// when one or more of them is not ready.
func runTicketList() int {
	out, err := exec.Command("gh", "issue", "list", "--label", "ticket", "--state", "open", "--limit", strconv.Itoa(ticketListLimit), "--json", issueFields).Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, "ticket: gh issue list:", err)
		return 2
	}
	var issues []*issue
	if err := json.Unmarshal(out, &issues); err != nil || issues == nil {
		fmt.Fprintln(os.Stderr, "ticket: unexpected gh output")
		return 2
	}
	if len(issues) >= ticketListLimit {
		fmt.Fprintf(os.Stderr, "ticket: %d or more open tickets, the list can be incomplete\n", ticketListLimit)
		return 2
	}
	tickets := []ticket.Ticket{}
	code := 0
	for _, is := range issues {
		if !is.valid() {
			fmt.Fprintln(os.Stderr, "ticket: gh returned an issue without a number, title or state")
			return 2
		}
		t := ticket.New(is.Number, is.Title, is.State, is.Body)
		if !t.Ready {
			code = 1
		}
		tickets = append(tickets, t)
	}
	if err := printJSON(tickets); err != nil {
		fmt.Fprintln(os.Stderr, "ticket:", err)
		return 2
	}
	return code
}
