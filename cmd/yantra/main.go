package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/bettersandeep/yantra/internal/referee"
)

const usage = `yantra - a factory that makes software factories

usage: yantra <command>

commands:
  version                                   print the version
  referee diff --base <rev> [--head <rev>]  report test weakening between two commits
`

func main() {
	args := os.Args[1:]
	switch {
	case len(args) == 1 && args[0] == "version":
		fmt.Println("yantra 0.0.0")
	case len(args) >= 2 && args[0] == "referee" && args[1] == "diff":
		os.Exit(runRefereeDiff(args[2:]))
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
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	if err := enc.Encode(findings); err != nil {
		fmt.Fprintln(os.Stderr, "referee:", err)
		return 2
	}
	if len(findings) > 0 {
		return 1
	}
	return 0
}

func refereeDiff(base, head string) ([]referee.Finding, error) {
	out, err := exec.Command("git", "diff", "--name-status", "-M", base, head).Output()
	if err != nil {
		return nil, fmt.Errorf("git diff %s %s: %w", base, head, err)
	}
	findings := []referee.Finding{}
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		f := strings.Split(line, "\t")
		if len(f) < 2 {
			continue
		}
		oldPath, newPath := f[1], f[len(f)-1] // a rename line has two paths
		if !strings.HasSuffix(oldPath, "_test.go") && !strings.HasSuffix(newPath, "_test.go") {
			continue
		}
		var b, h []byte
		if !strings.HasPrefix(f[0], "A") {
			if b, err = show(base, oldPath); err != nil {
				return nil, err
			}
		}
		if !strings.HasPrefix(f[0], "D") {
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
