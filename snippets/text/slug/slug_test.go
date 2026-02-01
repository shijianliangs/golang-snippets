package slug

import "testing"

func TestMake(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"Hello, World!", "hello-world"},
		{"  Multiple   spaces  ", "multiple-spaces"},
		{"snake_case_and-dash", "snake-case-and-dash"},
		{"Café déjà vu", "caf-dj-vu"}, // non-ASCII dropped
		{"--already--slug--", "already-slug"},
		{"", ""},
	}
	for _, tc := range cases {
		if got := Make(tc.in); got != tc.want {
			t.Fatalf("Make(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
