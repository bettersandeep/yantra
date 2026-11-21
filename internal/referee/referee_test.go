package referee

import "testing"

const baseSrc = `package x
import ("testing"; "github.com/stretchr/testify/require")
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
		name string
		head string
		want []string
	}{
		{"unchanged", baseSrc, nil},
		{"assert added", baseSrc + "func TestNew(t *testing.T) { require.True(t, true) }\n", nil},
		{"assert removed and skip added", `package x
import "testing"
func TestAdd(t *testing.T) {
	t.Skip("flaky")
	if add(1, 1) != 2 { t.Fatalf("bad") }
}
func TestSub(t *testing.T) { require.Equal(t, 0, sub(2, 2)) }
`, []string{"TestAdd:assert_removed", "TestAdd:skip_added"}},
		{"test deleted", `package x
import "testing"
func TestAdd(t *testing.T) {
	require.Equal(t, 4, add(2, 2))
	if add(1, 1) != 2 { t.Fatalf("bad") }
}
`, []string{"TestSub:test_deleted"}},
		{"nolint added", baseSrc + "var _ = 1 //nolint:all\n", []string{":nolint_added"}},
		{"file deleted", "", []string{"TestAdd:test_deleted", "TestSub:test_deleted"}},
	}
	for _, c := range cases {
		fs, err := Diff("x_test.go", []byte(baseSrc), []byte(c.head))
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		got := kinds(fs)
		if len(got) != len(c.want) {
			t.Fatalf("%s: got %v, want %v", c.name, got, c.want)
		}
		for i := range got {
			if got[i] != c.want[i] {
				t.Errorf("%s: got %v, want %v", c.name, got, c.want)
			}
		}
	}
}

func TestDiffParseError(t *testing.T) {
	if _, err := Diff("x_test.go", []byte(baseSrc), []byte("not go")); err == nil {
		t.Fatal("want error on unparseable head, got nil")
	}
}
