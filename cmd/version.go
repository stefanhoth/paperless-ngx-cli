package cmd

import (
	"fmt"
	"os"
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
		fmt.Printf("API version:        v%d\n", APIVersion)

		c, _ := mustClient()
		resp, err := c.RemoteVersionRetrieveWithResponse(ctx())
		if err != nil || resp.StatusCode() != 200 {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		serverAPIVersion := resp.HTTPResponse.Header.Get("X-Api-Version")
		if serverAPIVersion != "" && serverAPIVersion != fmt.Sprintf("%d", APIVersion) {
			fmt.Fprintf(os.Stderr, "warning: server API version is v%s, CLI targets v%d — consider updating the CLI\n", serverAPIVersion, APIVersion)
		}

		installed := resp.HTTPResponse.Header.Get("X-Version")
		if installed == "" {
			installed = "unknown"
		}

		v := resp.JSON200
		available := "—"
		updateAvail := "no"
		if v != nil {
			if s, ok := (*v)["version"].(string); ok {
				available = strings.TrimPrefix(s, "v")
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
