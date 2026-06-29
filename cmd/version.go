package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Installierte vs. verfügbare Version",
	Run: func(cmd *cobra.Command, args []string) {
		c, cfg := mustClient()
		resp, err := c.RemoteVersionRetrieveWithResponse(ctx())
		if err != nil || resp.StatusCode() != 200 {
			fmt.Fprintf(os.Stderr, "Fehler: %v\n", err)
			os.Exit(1)
		}

		installed := "unbekannt (kein SSH-Host konfiguriert)"
		if cfg.sshHost != "" {
			installed = sshInstalledVersion(cfg)
		}

		v := resp.JSON200
		available := "—"
		updateAvail := "nein"
		if v != nil {
			if s, ok := (*v)["version"].(string); ok {
				available = s
			}
			if b, ok := (*v)["update_available"].(bool); ok && b {
				updateAvail = "ja"
			}
		}

		fmt.Printf("Installiert:        %s\n", installed)
		fmt.Printf("Verfügbar (remote): %s\n", available)
		fmt.Printf("Update verfügbar:   %s\n", updateAvail)
	},
}

func sshInstalledVersion(cfg config) string {
	dockerCmd := fmt.Sprintf(
		"/usr/local/bin/docker inspect %s --format '{{index .Config.Labels \"org.opencontainers.image.version\"}}'",
		cfg.container,
	)
	out, err := exec.Command("ssh", cfg.sshUser+"@"+cfg.sshHost, dockerCmd).Output()
	if err != nil {
		return "unbekannt"
	}
	return strings.TrimSpace(string(out))
}
