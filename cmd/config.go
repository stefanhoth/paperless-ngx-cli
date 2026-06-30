package cmd

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type config struct {
	baseURL   string
	token     string
	sshHost   string
	sshUser   string
	container string
}

// configFilePath returns ~/.config/paperless-ngx/config (XDG).
func configFilePath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	return filepath.Join(dir, "paperless-ngx-cli", "config")
}

// readConfigFile parses a KEY=VALUE file. Returns nil if file does not exist.
// Lines starting with # and blank lines are ignored.
func readConfigFile(path string) map[string]string {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	vals := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		vals[strings.TrimSpace(k)] = strings.TrimSpace(v)
	}
	return vals
}

// get looks up key in env first, then falls back to file values.
func get(key string, fileVals map[string]string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fileVals[key]
}

// parseConfig builds a config from env vars and file values.
// Env vars always take precedence. Returns an error for missing required fields.
func parseConfig(fileVals map[string]string) (config, error) {
	baseURL := strings.TrimRight(get("PAPERLESS_URL", fileVals), "/")
	if baseURL == "" {
		return config{}, fmt.Errorf("PAPERLESS_URL is not set (e.g. http://paperless.example.com:8000)")
	}
	token := get("PAPERLESS_API_TOKEN", fileVals)
	if token == "" {
		return config{}, fmt.Errorf("PAPERLESS_API_TOKEN is not set")
	}

	sshHost := get("PAPERLESS_SSH_HOST", fileVals)
	if sshHost == "" {
		if u, err := url.Parse(baseURL); err == nil {
			sshHost = u.Hostname()
		}
	}

	sshUser := get("PAPERLESS_SSH_USER", fileVals)
	if sshUser == "" {
		if u, err := user.Current(); err == nil {
			sshUser = u.Username
		}
	}

	container := get("PAPERLESS_CONTAINER", fileVals)
	if container == "" {
		container = "paperless-ngx-webserver-1"
	}

	return config{
		baseURL:   baseURL,
		token:     token,
		sshHost:   sshHost,
		sshUser:   sshUser,
		container: container,
	}, nil
}

func loadConfig() config {
	fileVals := readConfigFile(configFilePath())
	cfg, err := parseConfig(fileVals)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return cfg
}
