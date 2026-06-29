package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(manageCmd)
}

var manageCmd = &cobra.Command{
	Use:   "manage <cmd> [args]",
	Short: "manage.py im Container ausführen (via SSH)",
	Long: `Führt manage.py-Kommandos im Paperless-Container aus.

Benötigt SSH-Zugang zum Host und einen laufenden Container.
Konfiguration via Umgebungsvariablen (siehe README).

Beispiele:
  paperless manage document_retagger
  paperless manage document_sanity_checker
  paperless manage document_index reindex`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		_, cfg := mustClient()

		if cfg.sshHost == "" {
			fmt.Fprintln(os.Stderr, "SSH-Host nicht konfiguriert. Setze PAPERLESS_SSH_HOST oder PAPERLESS_URL.")
			os.Exit(1)
		}

		manageArgs := strings.Join(args, " ")
		dockerCmd := fmt.Sprintf(
			"/usr/local/bin/docker exec %s python3 manage.py %s",
			cfg.container, manageArgs,
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
