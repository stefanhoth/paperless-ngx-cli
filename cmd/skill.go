package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// SkillMD holds the SKILL.md contents bundled into the binary via go:embed
// in main.go, and set on cmd.SkillMD before Execute is called.
var SkillMD string

// skillDirName is the directory name a skill installer looks for under
// .claude/skills, matching the "name" field in SKILL.md's frontmatter.
const skillDirName = "paperless-ngx-cli"

func init() {
	skillInstallCmd.Flags().BoolP("user", "u", false, "Install into the user skills directory (~/.claude/skills) instead of the current directory")
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
	Long: `Writes the SKILL.md bundled with this CLI version to a skills directory,
as a convenience alternative to a separate "skills add" step.

By default, installs into ./.claude/skills/paperless-ngx-cli/SKILL.md, relative
to the current directory. Use --user to install into
~/.claude/skills/paperless-ngx-cli/SKILL.md instead, making it available to
every project.`,
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

		path := skillInstallPath(base)

		if existing, err := os.ReadFile(path); err == nil && string(existing) == SkillMD { //nolint:gosec // G304: path is built from a fixed skill dir name, not user input
			fmt.Printf("✓ %s is already up to date\n", path)
			return
		}

		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil { //nolint:gosec // G301: skills directory is not sensitive, no need for 0700
			fmt.Fprintln(os.Stderr, "error creating skill directory:", err)
			os.Exit(1)
		}
		if err := os.WriteFile(path, []byte(SkillMD), 0o644); err != nil { //nolint:gosec // G306: skill file is not sensitive, no need for 0600
			fmt.Fprintln(os.Stderr, "error writing skill file:", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Skill installed to %s\n", path)
	},
}

// skillInstallPath returns the SKILL.md destination path under base's
// .claude/skills directory.
func skillInstallPath(base string) string {
	return filepath.Join(base, ".claude", "skills", skillDirName, "SKILL.md")
}
