// Package ticket reads a work ticket from the text of a GitHub issue and
// decides if the ticket is ready for a coding agent.
package ticket

import (
	"regexp"
	"strings"
)

type Ticket struct {
	Number   int      `json:"number"`
	Title    string   `json:"title"`
	State    string   `json:"state"`
	Criteria []string `json:"criteria"`
	Scope    []string `json:"scope"`
	Ready    bool     `json:"ready"`
	Reasons  []string `json:"reasons,omitempty"`
}

var (
	criteriaHead = regexp.MustCompile(`(?i)^\s*(#{1,6}\s*)?\**acceptance criteria\**:?\**\s*$`)
	listItem     = regexp.MustCompile(`^\s*(?:[-*+]|\d+[.)])\s+(?:\[[ xX]\]\s+)?(.*\S)\s*$`)
	scopeLine    = regexp.MustCompile(`(?i)^\s*\**scope\**:\**\s*(.*)$`)
	backticked   = regexp.MustCompile("`([^`]+)`")
)

// New builds a ticket from the fields of an issue. A ticket is ready when its
// body has at least one acceptance criterion and at least one scope path.
func New(number int, title, state, body string) Ticket {
	t := Ticket{Number: number, Title: title, State: state, Criteria: []string{}, Scope: []string{}}
	inCriteria := false
	for _, line := range strings.Split(strings.ReplaceAll(body, "\r\n", "\n"), "\n") {
		if criteriaHead.MatchString(line) {
			inCriteria = true
			continue
		}
		if inCriteria {
			if m := listItem.FindStringSubmatch(line); m != nil {
				t.Criteria = append(t.Criteria, m[1])
				continue
			}
			if strings.TrimSpace(line) == "" {
				continue
			}
			inCriteria = false // the first other line ends the list
		}
		if m := scopeLine.FindStringSubmatch(line); m != nil {
			t.Scope = append(t.Scope, scopePaths(m[1])...)
		}
	}
	if len(t.Criteria) == 0 {
		t.Reasons = append(t.Reasons, "no acceptance criteria")
	}
	if len(t.Scope) == 0 {
		t.Reasons = append(t.Reasons, "no scope paths")
	}
	t.Ready = len(t.Reasons) == 0
	return t
}

// scopePaths reads the paths from the text after "Scope:". Paths in backticks
// win. Without backticks the text is a comma-separated list.
func scopePaths(s string) []string {
	var paths []string
	if ms := backticked.FindAllStringSubmatch(s, -1); ms != nil {
		for _, m := range ms {
			paths = append(paths, strings.TrimSpace(m[1]))
		}
		return paths
	}
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimRight(strings.TrimSpace(p), "."); p != "" {
			paths = append(paths, p)
		}
	}
	return paths
}
