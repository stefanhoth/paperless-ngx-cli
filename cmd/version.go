package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// Version is set at build time via -ldflags "-X github.com/stefanhoth/paperless-ngx-cli/cmd.Version=x.y.z"
var Version = "dev"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show CLI and Paperless-NGX instance versions",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("paperless CLI:      %s\n", Version)

		c, cfg := mustClient()
		resp, err := c.RemoteVersionRetrieveWithResponse(ctx())
		if err != nil || resp.StatusCode() != 200 {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		installed := "unknown (PAPERLESS_SSH_HOST not configured)"
		if cfg.sshHost != "" {
			installed = sshInstalledVersion(cfg)
		}

		v := resp.JSON200
		available := "—"
		updateAvail := "no"
		if v != nil {
			if s, ok := (*v)["version"].(string); ok {
				available = s
			}
			if b, ok := (*v)["update_available"].(bool); ok && b {
				updateAvail = "yes"
			}
		}

		fmt.Printf("Paperless (server): %s\n", installed)
		fmt.Printf("Paperless (latest): %s\n", available)
		fmt.Printf("Update available:   %s\n", updateAvail)
	},
}

func sshInstalledVersion(cfg config) string {
	dockerCmd := fmt.Sprintf(
		"/usr/local/bin/docker inspect %s --format '{{index .Config.Labels \"org.opencontainers.image.version\"}}'",
		shellQuote(cfg.container),
	)
	out, err := exec.Command("ssh", cfg.sshUser+"@"+cfg.sshHost, dockerCmd).Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}
