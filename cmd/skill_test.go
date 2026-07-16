package cmd

import (
	"path/filepath"
	"testing"
)

func TestSkillInstallPath(t *testing.T) {
	cases := []struct {
		name string
		base string
		want string
	}{
		{"local", ".", filepath.Join(".", ".claude", "skills", "paperless-ngx-cli", "SKILL.md")},
		{"user", "/home/alice", filepath.Join("/home/alice", ".claude", "skills", "paperless-ngx-cli", "SKILL.md")},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := skillInstallPath(tc.base); got != tc.want {
				t.Errorf("skillInstallPath(%q) = %q, want %q", tc.base, got, tc.want)
			}
		})
	}
}
