package ticket

import (
	"slices"
	"testing"
)

func TestNew(t *testing.T) {
	cases := []struct {
		name     string
		body     string
		criteria []string
		scope    []string
		reasons  []string
	}{
		{"plain form",
			"Why we need it.\n\nAcceptance criteria:\n- first thing\n- second `thing`\n\nScope: `internal/ticket/`, `README.md`.",
			[]string{"first thing", "second `thing`"}, []string{"internal/ticket/", "README.md"}, nil},
		{"markdown heading, numbers and checkboxes",
			"## Acceptance Criteria\n\n1. one\n2) two\n- [ ] three\n- [x] four\n\n**Scope:** `a/`",
			[]string{"one", "two", "three", "four"}, []string{"a/"}, nil},
		{"bold heading and windows line ends",
			"**Acceptance criteria:**\r\n* one\r\n\r\nScope: `a.go`\r\n",
			[]string{"one"}, []string{"a.go"}, nil},
		{"scope without backticks",
			"Acceptance criteria:\n- one\n\nScope: cmd/yantra/, README.md.",
			[]string{"one"}, []string{"cmd/yantra/", "README.md"}, nil},
		{"the list ends at the first other line",
			"Acceptance criteria:\n- one\nNotes:\n- not a criterion\n\nScope: `a/`",
			[]string{"one"}, []string{"a/"}, nil},
		{"no acceptance criteria",
			"Please make it better.\n\nScope: `a/`",
			[]string{}, []string{"a/"}, []string{"no acceptance criteria"}},
		{"heading with no items",
			"Acceptance criteria:\n\nIt must work.\n\nScope: `a/`",
			[]string{}, []string{"a/"}, []string{"no acceptance criteria"}},
		{"no scope",
			"Acceptance criteria:\n- one\n",
			[]string{"one"}, []string{}, []string{"no scope paths"}},
		{"empty body",
			"",
			[]string{}, []string{}, []string{"no acceptance criteria", "no scope paths"}},
	}
	for _, c := range cases {
		got := New(7, "title", "OPEN", c.body)
		if !slices.Equal(got.Criteria, c.criteria) || !slices.Equal(got.Scope, c.scope) || !slices.Equal(got.Reasons, c.reasons) {
			t.Errorf("%s: got criteria %q scope %q reasons %q", c.name, got.Criteria, got.Scope, got.Reasons)
		}
		if got.Ready != (len(c.reasons) == 0) {
			t.Errorf("%s: ready = %v", c.name, got.Ready)
		}
	}
}
