package cmd

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stefanhoth/paperless-ngx-cli/api"
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

// config reads required environment variables and exits with a clear message if missing.
type config struct {
	baseURL   string
	token     string
	sshHost   string
	sshUser   string
	container string
}

func loadConfig() config {
	baseURL := strings.TrimRight(os.Getenv("PAPERLESS_URL"), "/")
	if baseURL == "" {
		fmt.Fprintln(os.Stderr, "PAPERLESS_URL ist nicht gesetzt (z.B. http://paperless.example.com:8000)")
		os.Exit(1)
	}
	token := os.Getenv("PAPERLESS_API_TOKEN")
	if token == "" {
		fmt.Fprintln(os.Stderr, "PAPERLESS_API_TOKEN ist nicht gesetzt")
		os.Exit(1)
	}

	sshHost := os.Getenv("PAPERLESS_SSH_HOST")
	if sshHost == "" {
		// Derive from PAPERLESS_URL hostname as fallback
		if u, err := url.Parse(baseURL); err == nil {
			sshHost = u.Hostname()
		}
	}

	sshUser := os.Getenv("PAPERLESS_SSH_USER")
	if sshUser == "" {
		if u, err := user.Current(); err == nil {
			sshUser = u.Username
		}
	}

	container := os.Getenv("PAPERLESS_CONTAINER")
	if container == "" {
		container = "paperless-ngx-webserver-1"
	}

	return config{
		baseURL:   baseURL,
		token:     token,
		sshHost:   sshHost,
		sshUser:   sshUser,
		container: container,
	}
}

func newClient(cfg config) (*api.ClientWithResponses, error) {
	addAuth := func(_ context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Token "+cfg.token)
		return nil
	}
	return api.NewClientWithResponses(cfg.baseURL, api.WithRequestEditorFn(addAuth))
}

func mustClient() (*api.ClientWithResponses, config) {
	cfg := loadConfig()
	c, err := newClient(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Client-Fehler:", err)
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
