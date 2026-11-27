package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const strong = "package x\n\nimport \"testing\"\n\nfunc TestA(t *testing.T) {\n\tif 1+1 != 2 {\n\t\tt.Fatal(\"bad\")\n\t}\n}\n"
const weak = "package x\n\nimport \"testing\"\n\nfunc TestA(t *testing.T) {\n}\n"

func git(t *testing.T, args ...string) {
	t.Helper()
	args = append([]string{"-c", "user.email=a@b.c", "-c", "user.name=t"}, args...)
	if out, err := exec.Command("git", args...).CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

// TestRefereeDiffCommand runs the command path against a real git repo.
func TestRefereeDiffCommand(t *testing.T) {
	dir := t.TempDir()
	wd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(wd)

	const name = "é_test.go" // git quotes this name unless -z is used
	git(t, "init", "-q", "-b", "main")
	os.WriteFile(name, []byte(strong), 0o644)
	git(t, "add", ".")
	git(t, "commit", "-q", "-m", "base")
	os.WriteFile(name, []byte(weak), 0o644)
	git(t, "commit", "-q", "-am", "weaken")
	git(t, "mv", name, "b_test.go")
	git(t, "commit", "-q", "-m", "rename")

	fs, err := refereeDiff("HEAD~2", "HEAD~1")
	if err != nil || len(fs) != 1 || fs[0].Kind != "assert_removed" || fs[0].File != name {
		t.Errorf("weakened test in a quoted path: got %v, %v; want one assert_removed", fs, err)
	}
	if fs, err := refereeDiff("HEAD~1", "HEAD"); err != nil || len(fs) != 0 {
		t.Errorf("pure rename: got %v, %v; want no findings", fs, err)
	}
	out := filepath.Join(dir, "written-by-git")
	if _, err := refereeDiff("--output="+out, "HEAD"); err == nil {
		t.Error("option-like revision: want error, got nil")
	}
	if _, err := os.Stat(out); err == nil {
		t.Error("option-like revision: git wrote a file")
	}
}

// fakeGh puts a gh script first on PATH. The script records that it ran and
// prints the given JSON, or fails when the JSON is empty.
func fakeGh(t *testing.T, json string) (ran func() bool) {
	t.Helper()
	dir := t.TempDir()
	marker := filepath.Join(dir, "ran")
	script := "#!/bin/sh\ntouch '" + marker + "'\n"
	if json == "" {
		script += "exit 1\n"
	} else {
		script += "cat <<'JSON'\n" + json + "\nJSON\n"
	}
	if err := os.WriteFile(filepath.Join(dir, "gh"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	return func() bool { _, err := os.Stat(marker); return err == nil }
}

func TestTicketShowExitCodes(t *testing.T) {
	ready := `{"number":5,"title":"t","state":"OPEN","body":"Acceptance criteria:\n- one\n\nScope: ` + "`a/`" + `"}`
	notReady := `{"number":5,"title":"t","state":"OPEN","body":"Please make it better."}`
	cases := []struct {
		name string
		json string
		arg  string
		want int
	}{
		{"ready ticket", ready, "5", 0},
		{"ticket without criteria", notReady, "5", 1},
		{"gh fails", "", "5", 2},
		{"gh returns a different issue", ready, "6", 2},
		{"gh returns no JSON", "not json", "5", 2},
	}
	for _, c := range cases {
		fakeGh(t, c.json)
		if got := runTicketShow([]string{c.arg}); got != c.want {
			t.Errorf("%s: exit %d, want %d", c.name, got, c.want)
		}
	}
}

func TestTicketShowBadNumberDoesNotReachGh(t *testing.T) {
	for _, arg := range []string{"abc", "0", "-3", "5; rm -rf x", "--repo=other/repo", "5.0", ""} {
		ran := fakeGh(t, `{"number":5}`)
		if got := runTicketShow([]string{arg}); got != 2 {
			t.Errorf("%q: exit %d, want 2", arg, got)
		}
		if ran() {
			t.Errorf("%q: gh was called", arg)
		}
	}
	if got := runTicketShow(nil); got != 2 {
		t.Errorf("no argument: exit %d, want 2", got)
	}
}

func TestTicketList(t *testing.T) {
	fakeGh(t, `[{"number":5,"title":"a","state":"OPEN","body":"Acceptance criteria:\n- one\n\nScope: x/"},{"number":6,"title":"b","state":"OPEN","body":"no criteria"}]`)
	if got := runTicketList(); got != 1 {
		t.Errorf("one ticket not ready: exit %d, want 1", got)
	}
	fakeGh(t, `[]`)
	if got := runTicketList(); got != 0 {
		t.Errorf("no tickets: exit %d, want 0", got)
	}
	fakeGh(t, "")
	if got := runTicketList(); got != 2 {
		t.Errorf("gh fails: exit %d, want 2", got)
	}
}

// One case per shape that the Qodo review on pull request 7 named.
func TestTicketListFailsClosed(t *testing.T) {
	one := `{"number":5,"title":"a","state":"OPEN","body":"Acceptance criteria:\n- one\n\nScope: x/"}`
	for name, json := range map[string]string{
		"null":                   `null`,
		"null entry":             `[null]`,
		"an object, not a list":  one,
		"issue without a number": `[{"title":"a","state":"OPEN","body":"Acceptance criteria:\n- one\n\nScope: x/"}]`,
		"issue without a title":  `[{"number":5,"state":"OPEN","body":"Acceptance criteria:\n- one\n\nScope: x/"}]`,
		"list at the limit":      "[" + strings.Repeat(one+",", ticketListLimit-1) + one + "]",
	} {
		fakeGh(t, json)
		if got := runTicketList(); got != 2 {
			t.Errorf("%s: exit %d, want 2", name, got)
		}
	}
	fakeGh(t, `{"number":5,"body":"Acceptance criteria:\n- one\n\nScope: x/"}`)
	if got := runTicketShow([]string{"5"}); got != 2 {
		t.Errorf("show, issue without title and state: exit %d, want 2", got)
	}
}

func TestRefereeScopeCommand(t *testing.T) {
	dir := t.TempDir()
	wd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(wd)
	write := func(path string) {
		os.MkdirAll(filepath.Dir(path), 0o755)
		os.WriteFile(path, []byte(path+"\n"), 0o644)
	}
	git(t, "init", "-q", "-b", "main")
	write("README.md")
	git(t, "add", ".")
	git(t, "commit", "-q", "-m", "base")
	write("cmd/yantra/new.go")
	git(t, "add", ".")
	git(t, "commit", "-q", "-m", "in scope")
	write("go.mod")
	git(t, "add", ".")
	git(t, "commit", "-q", "-m", "out of scope")

	ready := `{"number":5,"title":"t","state":"OPEN","body":"Acceptance criteria:\n- one\n\nScope: ` + "`cmd/yantra/`, `README.md`" + `"}`
	notReady := `{"number":5,"title":"t","state":"OPEN","body":"Acceptance criteria:\n- one"}`
	cases := []struct {
		name string
		json string
		args []string
		want int
	}{
		{"all changes in scope", ready, []string{"--ticket", "5", "--base", "HEAD~2", "--head", "HEAD~1"}, 0},
		{"one change out of scope", ready, []string{"--ticket", "5", "--base", "HEAD~2"}, 1},
		{"ticket without scope", notReady, []string{"--ticket", "5", "--base", "HEAD~2", "--head", "HEAD~1"}, 2},
		{"gh fails", "", []string{"--ticket", "5", "--base", "HEAD~2", "--head", "HEAD~1"}, 2},
		{"no ticket number", ready, []string{"--base", "HEAD~2"}, 2},
		{"ticket number is not a number", ready, []string{"--ticket", "five", "--base", "HEAD~2"}, 2},
		{"no base", ready, []string{"--ticket", "5"}, 2},
		{"bad revision", ready, []string{"--ticket", "5", "--base", "nosuchrev"}, 2},
		{"extra argument", ready, []string{"--ticket", "5", "--base", "HEAD~2", "extra"}, 2},
	}
	for _, c := range cases {
		fakeGh(t, c.json)
		if got := runRefereeScope(c.args); got != c.want {
			t.Errorf("%s: exit %d, want %d", c.name, got, c.want)
		}
	}
}

// fakeGit puts a git script first on PATH. It resolves each revision and
// prints the given name-status output for a diff.
func fakeGit(t *testing.T, diffOutput string) {
	t.Helper()
	dir := t.TempDir()
	script := "#!/bin/sh\ncase \"$1\" in\nrev-parse) echo 1111111111111111111111111111111111111111 ;;\ndiff) printf '" + diffOutput + "' ;;\nesac\n"
	if err := os.WriteFile(filepath.Join(dir, "git"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

// The finding of the Qodo review on pull request 8: a file status that the
// check does not know must not read as clean, also when its path is in scope.
func TestRefereeScopeUnknownGitStatus(t *testing.T) {
	ready := `{"number":5,"title":"t","state":"OPEN","body":"Acceptance criteria:\n- one\n\nScope: ` + "`cmd/yantra/`" + `"}`
	args := []string{"--ticket", "5", "--base", "a", "--head", "b"}
	for name, c := range map[string]struct {
		diff string
		want int
	}{
		"known status in scope":   {`M\000cmd/yantra/a.go\000`, 0},
		"unknown status in scope": {`X\000cmd/yantra/a.go\000`, 2},
		"unmerged status":         {`U\000cmd/yantra/a.go\000`, 2},
		"cut-off rename":          {`R100\000cmd/yantra/a.go\000`, 2},
	} {
		fakeGh(t, ready)
		fakeGit(t, c.diff)
		if got := runRefereeScope(args); got != c.want {
			t.Errorf("%s: exit %d, want %d", name, got, c.want)
		}
	}
}
