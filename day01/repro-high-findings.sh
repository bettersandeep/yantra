set -e
cd /src && go build -o /tmp/yantra ./cmd/yantra
mkdir /tmp/r && cd /tmp/r && git init -q -b main && git config user.email a@b.c && git config user.name t
printf 'package x\n\nimport "testing"\n\nfunc TestA(t *testing.T) {\n\tif 1+1 != 2 {\n\t\tt.Fatal("bad")\n\t}\n}\n\nfunc TestB(t *testing.T) {\n\tcheck(t)\n}\n' > a_test.go
cp a_test.go 'é_test.go'
git add . && git commit -q -m base
set +e
/tmp/yantra referee diff --base=--output=/tmp/pwned --head=HEAD >/dev/null 2>&1
echo "H1 option-like revision: exit=$? (want 2) file_created=$([ -e /tmp/pwned ] && echo yes || echo no) (want no)"
printf 'package x\n\nimport "testing"\n\nfunc TestA(t *testing.T) {\n}\n\nfunc TestB(t *testing.T) {\n\tcheck(t)\n}\n' > 'é_test.go'
git commit -q -am weaken-quoted
/tmp/yantra referee diff --base HEAD~1 >/tmp/o 2>&1; echo "H2 quoted path: exit=$? (want 1) $(grep -c assert_removed /tmp/o) assert_removed"
printf 'package x\n\nimport "testing"\n\nfunc TestA(t *testing.T) {\n\tif 1+1 != 2 {\n\t\tt.Fatal("bad")\n\t}\n}\n' > a_test.go
git commit -q -am delete-b
/tmp/yantra referee diff --base HEAD~1 >/tmp/o 2>&1; echo "H3 assertion-free test deleted: exit=$? (want 1) $(grep -c test_deleted /tmp/o) test_deleted"
