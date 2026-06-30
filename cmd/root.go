package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/stefanhoth/paperless-ngx-cli/api"
)

// APIVersion is the Paperless-NGX REST API version this CLI targets.
// A new major CLI version is released for each new API version.
const APIVersion = 9

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
		req.Header.Set("Accept", fmt.Sprintf("application/json; version=%d", APIVersion))
		return nil
	}
	return api.NewClientWithResponses(cfg.baseURL, api.WithRequestEditorFn(addHeaders))
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
