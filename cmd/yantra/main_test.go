package main

import (
	"os"
	"os/exec"
	"path/filepath"
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
