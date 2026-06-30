package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/stefanhoth/paperless-ngx-cli/api"
)

const (
	targetAPIVersion = 10
	minAPIVersion    = 9
)

var rootCmd = &cobra.Command{
	Use:   "paperless",
	Short: "Paperless-NGX CLI",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func newClient(cfg config) (*api.ClientWithResponses, error) {
	addHeaders := func(_ context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Token "+cfg.token)
		req.Header.Set("Accept", fmt.Sprintf("application/json; version=%d", targetAPIVersion))
		return nil
	}
	return api.NewClientWithResponses(cfg.baseURL, api.WithRequestEditorFn(addHeaders))
}

// checkAPIVersion reads X-Api-Version from a response header and warns if
// the server is outside the supported range [minAPIVersion, targetAPIVersion].
func checkAPIVersion(header http.Header) {
	val := header.Get("X-Api-Version")
	if val == "" {
		return
	}
	v, err := strconv.Atoi(val)
	if err != nil {
		return
	}
	if v < minAPIVersion {
		fmt.Fprintf(os.Stderr, "warning: server API version %d is below minimum supported version %d — some commands may not work correctly\n", v, minAPIVersion)
	} else if v > targetAPIVersion {
		fmt.Fprintf(os.Stderr, "warning: server API version %d is newer than this CLI was tested against (v%d) — consider updating the CLI\n", v, targetAPIVersion)
	}
}

func mustClient() (*api.ClientWithResponses, config) {
	cfg := loadConfig()
	c, err := newClient(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "client error:", err)
		os.Exit(1)
	}
	return c, cfg
}

func ctx() context.Context {
	return context.Background()
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}
