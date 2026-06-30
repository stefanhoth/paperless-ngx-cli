package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// shellQuote wraps s in single quotes and escapes embedded single quotes,
// making it safe for interpolation into a remote shell command string.
func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

func init() {
	rootCmd.AddCommand(manageCmd)
}

var manageCmd = &cobra.Command{
	Use:   "manage <cmd> [args]",
	Short: "Run manage.py commands inside the Paperless container (via SSH)",
	Long: `Runs manage.py commands inside the Paperless Docker container over SSH.

Requires SSH access to the Docker host and a running container.
Configure via environment variables (see README).

Examples:
  paperless manage document_retagger
  paperless manage document_sanity_checker
  paperless manage document_index reindex`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		_, cfg := mustClient()

		if cfg.sshHost == "" {
			fmt.Println("manage requires SSH access to the Docker host.")
			fmt.Println("Set PAPERLESS_SSH_HOST (or derive it from PAPERLESS_URL) to enable this command.")
			fmt.Println("Optional: PAPERLESS_SSH_USER (default: current OS user), PAPERLESS_CONTAINER (default: paperless-ngx-webserver-1)")
			return
		}

		quoted := make([]string, len(args))
		for i, a := range args {
			quoted[i] = shellQuote(a)
		}
		dockerCmd := fmt.Sprintf(
			"/usr/local/bin/docker exec %s python3 manage.py %s",
			shellQuote(cfg.container), strings.Join(quoted, " "),
		)

		fmt.Printf("SSH %s@%s: %s\n\n", cfg.sshUser, cfg.sshHost, dockerCmd)

		c := exec.Command("ssh", cfg.sshUser+"@"+cfg.sshHost, dockerCmd)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			os.Exit(1)
		}
	},
}
