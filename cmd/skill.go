package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// SkillMD holds the SKILL.md contents bundled into the binary via go:embed
// in main.go, and set on cmd.SkillMD before Execute is called.
var SkillMD string

// skillDirName is the directory name a skill installer looks for, matching
// the "name" field in SKILL.md's frontmatter.
const skillDirName = "paperless-ngx-cli"

func init() {
	skillInstallCmd.Flags().BoolP("user", "u", false, "Install into the user skills directories (~/.agents, ~/.claude) instead of the current directory")
	skillCmd.AddCommand(skillInstallCmd)
	rootCmd.AddCommand(skillCmd)
}

var skillCmd = &cobra.Command{
	Use:   "skill",
	Short: "Manage the bundled AI assistant skill",
}

var skillInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the bundled SKILL.md so AI assistants like Claude Code pick it up",
	Long: `Writes the SKILL.md bundled with this CLI version to
.agents/skills/paperless-ngx-cli/SKILL.md — the vendor-neutral location several
AI coding agents read skills from — and symlinks
.claude/skills/paperless-ngx-cli to it, as a convenience alternative to a
separate "skills add" step.

By default, installs relative to the current directory. Use --user to install
under your home directory instead, making it available to every project.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, _ []string) {
		user, _ := cmd.Flags().GetBool("user")

		base := "."
		if user {
			home, err := os.UserHomeDir()
			if err != nil {
				fmt.Fprintln(os.Stderr, "error resolving home directory:", err)
				os.Exit(1)
			}
			base = home
		}

		canonicalDir := skillCanonicalDir(base)
		canonicalFile := filepath.Join(canonicalDir, "SKILL.md")
		if err := writeSkillFile(canonicalFile, SkillMD); err != nil {
			fmt.Fprintln(os.Stderr, "error installing skill:", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Skill installed to %s\n", canonicalFile)

		linkPath := skillLinkPath(base)
		if err := ensureSkillSymlink(linkPath, canonicalDir); err != nil {
			fmt.Fprintln(os.Stderr, "error linking skill:", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Linked from %s\n", linkPath)
	},
}

// skillCanonicalDir returns the .agents/skills/<name> directory under base —
// the vendor-neutral location this SKILL.md is installed into.
func skillCanonicalDir(base string) string {
	return filepath.Join(base, ".agents", "skills", skillDirName)
}

// skillLinkPath returns the .claude/skills/<name> path under base, which is
// symlinked to the canonical .agents/skills/<name> directory.
func skillLinkPath(base string) string {
	return filepath.Join(base, ".claude", "skills", skillDirName)
}

// writeSkillFile writes content to path, creating parent directories as
// needed. It's a no-op if the file already has the given content.
func writeSkillFile(path, content string) error {
	if existing, err := os.ReadFile(path); err == nil && string(existing) == content { //nolint:gosec // G304: path is built from a fixed skill dir name, not user input
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil { //nolint:gosec // G301: skills directory is not sensitive, no need for 0700
		return err
	}
	return os.WriteFile(path, []byte(content), 0o644) //nolint:gosec // G306: skill file is not sensitive, no need for 0600
}

// ensureSkillSymlink creates a relative symlink at linkPath pointing to
// targetDir. It's a no-op if linkPath is already a symlink resolving to
// targetDir, and refuses to overwrite linkPath if it exists as a real file
// or directory rather than a symlink.
func ensureSkillSymlink(linkPath, targetDir string) error {
	rel, err := filepath.Rel(filepath.Dir(linkPath), targetDir)
	if err != nil {
		return err
	}

	if info, err := os.Lstat(linkPath); err == nil {
		if info.Mode()&os.ModeSymlink == 0 {
			return fmt.Errorf("%s already exists and is not a symlink — remove it and re-run", linkPath)
		}
		if current, err := os.Readlink(linkPath); err == nil && current == rel {
			return nil
		}
		if err := os.Remove(linkPath); err != nil {
			return err
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(linkPath), 0o755); err != nil { //nolint:gosec // G301: skills directory is not sensitive, no need for 0700
		return err
	}
	return os.Symlink(rel, linkPath)
}
