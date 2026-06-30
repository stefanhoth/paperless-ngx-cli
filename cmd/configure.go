package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configureCmd)
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Create or update the config file at ~/.config/paperless-ngx-cli/config",
	Long: `Interactively prompts for configuration values and writes them to
~/.config/paperless-ngx-cli/config with secure permissions (0600).

Existing values are shown as defaults — press Enter to keep them.
SSH settings are optional and only needed for the manage and version commands.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		path := configFilePath()
		existing := readConfigFile(path)
		if existing == nil {
			existing = map[string]string{}
		}

		scanner := bufio.NewScanner(os.Stdin)

		prompt := func(label, key, fallback string) string {
			current := existing[key]
			if current == "" {
				current = fallback
			}
			shown := current
			if key == "PAPERLESS_API_TOKEN" && shown != "" {
				shown = shown[:min(6, len(shown))] + "…"
			}
			if shown != "" {
				fmt.Printf("%s [%s]: ", label, shown)
			} else {
				fmt.Printf("%s: ", label)
			}
			scanner.Scan()
			input := strings.TrimSpace(scanner.Text())
			if input == "" {
				return current
			}
			return input
		}

		fmt.Println("Paperless-NGX CLI — Configuration")
		fmt.Println(strings.Repeat("─", 40))

		url := prompt("Paperless URL (e.g. http://paperless.local:8000)", "PAPERLESS_URL", "")
		if url == "" {
			fmt.Fprintln(os.Stderr, "error: PAPERLESS_URL is required")
			os.Exit(1)
		}
		url = strings.TrimRight(url, "/")

		token := prompt("API token (Settings → API in the Paperless web UI)", "PAPERLESS_API_TOKEN", "")
		if token == "" {
			fmt.Fprintln(os.Stderr, "error: PAPERLESS_API_TOKEN is required")
			os.Exit(1)
		}

		fmt.Println("\nSSH settings (optional — only needed for manage/version commands)")
		sshHost := prompt("SSH host (leave empty to derive from URL)", "PAPERLESS_SSH_HOST", "")
		sshUser := prompt("SSH user (leave empty for current OS user)", "PAPERLESS_SSH_USER", "")
		fmt.Println("  Container name: only needed when Paperless runs in Docker (uses 'docker exec' to call manage.py).")
		fmt.Println("  Leave empty if Paperless runs directly on the host (bare-metal/venv).")
		container := prompt("Container name", "PAPERLESS_CONTAINER", "paperless-ngx-webserver-1")

		var sb strings.Builder
		sb.WriteString("PAPERLESS_URL=" + url + "\n")
		sb.WriteString("PAPERLESS_API_TOKEN=" + token + "\n")
		if sshHost != "" {
			sb.WriteString("PAPERLESS_SSH_HOST=" + sshHost + "\n")
		}
		if sshUser != "" {
			sb.WriteString("PAPERLESS_SSH_USER=" + sshUser + "\n")
		}
		if container != "" && container != "paperless-ngx-webserver-1" {
			sb.WriteString("PAPERLESS_CONTAINER=" + container + "\n")
		}

		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			fmt.Fprintf(os.Stderr, "error creating config directory: %v\n", err)
			os.Exit(1)
		}
		if err := os.WriteFile(path, []byte(sb.String()), 0o600); err != nil {
			fmt.Fprintf(os.Stderr, "error writing config: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\n✓ Config written to %s\n", path)
	},
}

