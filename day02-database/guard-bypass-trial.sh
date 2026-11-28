# Three ways past the first version of the checks, tried against the old steps and the new guard.
set -u
old_steps() { ( cd tests/slt && sha256sum -c --quiet SHA256SUMS && test "$(find . -name '*.test' | wc -l)" -eq "$(wc -l < SHA256SUMS)" ) >/dev/null 2>&1 && test "$(grep -c '^\[\[package\]\]' Cargo.lock)" -eq 1; }
try() { old_steps && o=PASS || o=FAIL; GITHUB_EVENT_NAME=pull_request GITHUB_BASE_REF=main sh ci/guard.sh >/dev/null 2>&1 && n=PASS || n=FAIL; echo "$1: old steps $o, new guard $n"; }
git config --global --add safe.directory '*'; git config --global user.email a@b.c; git config --global user.name t
cp -r /src /tmp/w && cd /tmp/w && git add -A && git commit -q -m "this pull request" && git update-ref refs/remotes/origin/main HEAD
try "0. no change (control)"
git checkout -q -b t1 && printf '\nquery I nosort\nSELECT 1\n----\n1\n' >> tests/slt/random/expr/slt_good_0.test && (cd tests/slt && sha256sum random/expr/slt_good_0.test > SHA256SUMS) && echo '// x' >> src/lib.rs && git commit -q -am t1
try "1. edit a test file and its checksum together, in an engine change"
git checkout -q -f main 2>/dev/null || git checkout -q -f -; git checkout -q -B t2 refs/remotes/origin/main && (cd tests/slt && sha256sum FORMAT.wiki > SHA256SUMS) && git commit -q -am t2
try "2. list another file in the manifest in place of the test file"
