package cmd

import "testing"

func TestShellQuote(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"document_retagger", "'document_retagger'"},
		{"reindex", "'reindex'"},
		{"it's fine", `'it'\''s fine'`},
		{"'; rm -rf /; echo '", `''\''; rm -rf /; echo '\'''`},
		{"", "''"},
	}
	for _, tc := range cases {
		got := shellQuote(tc.input)
		if got != tc.want {
			t.Errorf("shellQuote(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}
