package referee

import (
	"slices"
	"testing"
)

const hdr = "package x\nimport (\"testing\"; \"github.com/stretchr/testify/require\")\n"

const baseSrc = hdr + `
func TestAdd(t *testing.T) {
	require.Equal(t, 4, add(2, 2))
	if add(1, 1) != 2 { t.Fatalf("bad") }
}
func TestSub(t *testing.T) { require.Equal(t, 0, sub(2, 2)) }
`

func kinds(fs []Finding) []string {
	var ks []string
	for _, f := range fs {
		ks = append(ks, f.Func+":"+f.Kind)
	}
	return ks
}

func TestDiff(t *testing.T) {
	cases := []struct {
		name       string
		base, head string
		want       []string
	}{
		{"unchanged", baseSrc, baseSrc, nil},
		{"assert added", baseSrc, baseSrc + "func TestNew(t *testing.T) { require.True(t, true) }\n", nil},
		{"assert removed and skip added", baseSrc, hdr + `
func TestAdd(t *testing.T) {
	t.Skip("flaky")
	if add(1, 1) != 2 { t.Fatalf("bad") }
}
func TestSub(t *testing.T) { require.Equal(t, 0, sub(2, 2)) }
`, []string{"TestAdd:assert_removed", "TestAdd:skip_added"}},
		{"test deleted", baseSrc, hdr + `
func TestAdd(t *testing.T) {
	require.Equal(t, 4, add(2, 2))
	if add(1, 1) != 2 { t.Fatalf("bad") }
}
`, []string{"TestSub:test_deleted"}},
		{"nolint added", baseSrc, baseSrc + "var _ = 1 //nolint:all\n", []string{":nolint_added"}},
		{"file deleted", baseSrc, "", []string{"TestAdd:test_deleted", "TestSub:test_deleted"}},

		// One case per finding of the Qodo review on pull request 2.
		{"test without assertions deleted",
			hdr + "func TestViaHelper(t *testing.T) { check(t) }\nfunc helper() {}\n",
			hdr,
			[]string{"TestViaHelper:test_deleted"}},
		{"aliased assertion import",
			"package x\nimport (\"testing\"; a \"github.com/stretchr/testify/assert\")\nfunc TestA(t *testing.T) { a.Equal(t, 1, f()) }\n",
			"package x\nimport (\"testing\"; a \"github.com/stretchr/testify/assert\")\nfunc TestA(t *testing.T) { }\n",
			[]string{"TestA:assert_removed"}},
		{"same method name on two receivers",
			hdr + "func (s *One) TestX(t *testing.T) { t.Fatal() }\nfunc (s *Two) TestX(t *testing.T) { t.Fatal() }\n",
			hdr + "func (s *One) TestX(t *testing.T) { }\nfunc (s *Two) TestX(t *testing.T) { t.Fatal() }\n",
			[]string{"One.TestX:assert_removed"}},
		{"check moved from one subtest to another",
			hdr + "func TestA(t *testing.T) {\n t.Run(\"a\", func(t *testing.T) { t.Fatal() })\n t.Run(\"b\", func(t *testing.T) { t.Fatal() })\n}\n",
			hdr + "func TestA(t *testing.T) {\n t.Run(\"a\", func(t *testing.T) { })\n t.Run(\"b\", func(t *testing.T) { t.Fatal(); t.Fatal() })\n}\n",
			[]string{"TestA/a:assert_removed"}},
		{"conditional skip replaced by unconditional skip",
			hdr + "func TestA(t *testing.T) {\n if testing.Short() { t.Skip() }\n t.Fatal()\n}\n",
			hdr + "func TestA(t *testing.T) {\n t.Skip()\n t.Fatal()\n}\n",
			[]string{"TestA:skip_added"}},
		{"nolint replaced by a wider nolint",
			baseSrc + "var _ = 1 //nolint:errcheck\n",
			baseSrc + "var _ = 1 //nolint:all\n",
			[]string{":nolint_added"}},
		{"unrelated Skip and Error methods",
			hdr + "func TestA(t *testing.T) { _ = err.Error(); t.Fatal() }\n",
			hdr + "func TestA(t *testing.T) { helper.Skip(); t.Fatal() }\n",
			nil},
	}
	for _, c := range cases {
		fs, err := Diff("x_test.go", []byte(c.base), []byte(c.head))
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		if got := kinds(fs); !slices.Equal(got, c.want) {
			t.Errorf("%s: got %v, want %v", c.name, got, c.want)
		}
	}
}

func TestDiffFailsClosed(t *testing.T) {
	for name, head := range map[string]string{
		"unparseable head":        "not go",
		"dot import of assertion": "package x\nimport . \"github.com/stretchr/testify/assert\"\n",
	} {
		if _, err := Diff("x_test.go", []byte(baseSrc), []byte(head)); err == nil {
			t.Errorf("%s: want error, got nil", name)
		}
	}
}
