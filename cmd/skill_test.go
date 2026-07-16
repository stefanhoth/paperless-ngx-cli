package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSkillCanonicalDir(t *testing.T) {
	cases := []struct {
		name string
		base string
		want string
	}{
		{"local", ".", filepath.Join(".", ".agents", "skills", "paperless-ngx-cli")},
		{"user", "/home/alice", filepath.Join("/home/alice", ".agents", "skills", "paperless-ngx-cli")},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := skillCanonicalDir(tc.base); got != tc.want {
				t.Errorf("skillCanonicalDir(%q) = %q, want %q", tc.base, got, tc.want)
			}
		})
	}
}

func TestSkillLinkPath(t *testing.T) {
	cases := []struct {
		name string
		base string
		want string
	}{
		{"local", ".", filepath.Join(".", ".claude", "skills", "paperless-ngx-cli")},
		{"user", "/home/alice", filepath.Join("/home/alice", ".claude", "skills", "paperless-ngx-cli")},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := skillLinkPath(tc.base); got != tc.want {
				t.Errorf("skillLinkPath(%q) = %q, want %q", tc.base, got, tc.want)
			}
		})
	}
}

func TestWriteSkillFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "SKILL.md")

	if err := writeSkillFile(path, "content-v1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := os.ReadFile(path) //nolint:gosec // G304: test-controlled path
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(got) != "content-v1" {
		t.Errorf("content = %q, want %q", got, "content-v1")
	}

	if err := writeSkillFile(path, "content-v2"); err != nil {
		t.Fatalf("unexpected error updating: %v", err)
	}
	got, err = os.ReadFile(path) //nolint:gosec // G304: test-controlled path
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(got) != "content-v2" {
		t.Errorf("content after update = %q, want %q", got, "content-v2")
	}
}

func TestEnsureSkillSymlink_CreatesRelativeLink(t *testing.T) {
	dir := t.TempDir()
	targetDir := filepath.Join(dir, ".agents", "skills", "paperless-ngx-cli")
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		t.Fatal(err)
	}
	linkPath := filepath.Join(dir, ".claude", "skills", "paperless-ngx-cli")

	if err := ensureSkillSymlink(linkPath, targetDir); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resolved, err := filepath.EvalSymlinks(linkPath)
	if err != nil {
		t.Fatalf("unexpected error resolving symlink: %v", err)
	}
	wantResolved, err := filepath.EvalSymlinks(targetDir)
	if err != nil {
		t.Fatal(err)
	}
	if resolved != wantResolved {
		t.Errorf("symlink resolves to %q, want %q", resolved, wantResolved)
	}

	target, err := os.Readlink(linkPath)
	if err != nil {
		t.Fatal(err)
	}
	if filepath.IsAbs(target) {
		t.Errorf("symlink target %q should be relative", target)
	}
}

func TestEnsureSkillSymlink_IdempotentOnSecondRun(t *testing.T) {
	dir := t.TempDir()
	targetDir := filepath.Join(dir, ".agents", "skills", "paperless-ngx-cli")
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		t.Fatal(err)
	}
	linkPath := filepath.Join(dir, ".claude", "skills", "paperless-ngx-cli")

	if err := ensureSkillSymlink(linkPath, targetDir); err != nil {
		t.Fatalf("unexpected error on first run: %v", err)
	}
	if err := ensureSkillSymlink(linkPath, targetDir); err != nil {
		t.Fatalf("unexpected error on second run: %v", err)
	}
}

func TestEnsureSkillSymlink_RefusesToOverwriteRealDirectory(t *testing.T) {
	dir := t.TempDir()
	targetDir := filepath.Join(dir, ".agents", "skills", "paperless-ngx-cli")
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		t.Fatal(err)
	}
	linkPath := filepath.Join(dir, ".claude", "skills", "paperless-ngx-cli")
	if err := os.MkdirAll(linkPath, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := ensureSkillSymlink(linkPath, targetDir); err == nil {
		t.Fatal("expected error when linkPath is a real directory")
	}
}
