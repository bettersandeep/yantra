package referee

import (
	"slices"
	"testing"
)

func TestOutOfScope(t *testing.T) {
	scope := []string{"internal/ticket/", "cmd/yantra/", "README.md"}
	cases := []struct {
		name string
		c    Change
		out  bool
	}{
		{"file below a scope directory", Change{'M', "cmd/yantra/main.go", "cmd/yantra/main.go"}, false},
		{"deep file below a scope directory", Change{'A', "internal/ticket/a/b/c.go", "internal/ticket/a/b/c.go"}, false},
		{"exact file", Change{'M', "README.md", "README.md"}, false},
		{"directory with the same start", Change{'A', "cmd/yantra-extra/x.go", "cmd/yantra-extra/x.go"}, true},
		{"file with the same start", Change{'A', "README.md.bak", "README.md.bak"}, true},
		{"same file name in a different directory", Change{'M', "docs/README.md", "docs/README.md"}, true},
		{"file outside", Change{'M', "internal/referee/referee.go", "internal/referee/referee.go"}, true},
		{"deleted file outside", Change{'D', ".github/workflows/ci.yml", ".github/workflows/ci.yml"}, true},
		{"rename inside the scope", Change{'R', "cmd/yantra/a.go", "cmd/yantra/b.go"}, false},
		{"rename into the scope", Change{'R', "internal/referee/a.go", "cmd/yantra/a.go"}, true},
		{"rename out of the scope", Change{'R', "cmd/yantra/a.go", "internal/gate/a.go"}, true},
		{"copy into the scope", Change{'C', "internal/referee/a.go", "cmd/yantra/a.go"}, false},
		{"copy out of the scope", Change{'C', "cmd/yantra/a.go", "internal/gate/a.go"}, true},
		{"unknown status outside", Change{'X', "go.mod", "go.mod"}, true},
	}
	for _, c := range cases {
		fs := OutOfScope([]Change{c.c}, scope)
		if (len(fs) == 1) != c.out {
			t.Errorf("%s: findings %v, want out of scope = %v", c.name, fs, c.out)
		}
		if len(fs) == 1 && (fs[0].Kind != "out_of_scope" || fs[0].File != c.c.New) {
			t.Errorf("%s: finding %+v", c.name, fs[0])
		}
	}
}

func TestOutOfScopeWeakScopes(t *testing.T) {
	c := []Change{{'M', "go.mod", "go.mod"}}
	for _, scope := range [][]string{nil, {""}, {"/"}, {"./"}, {"."}, {"**"}, {"go.mod/"}} {
		if fs := OutOfScope(c, scope); len(fs) != 1 {
			t.Errorf("scope %q: findings %v, want 1", scope, fs)
		}
	}
	if fs := OutOfScope(nil, []string{"a/"}); !slices.Equal(fs, []Finding{}) {
		t.Errorf("no changes: got %v, want an empty list", fs)
	}
}
