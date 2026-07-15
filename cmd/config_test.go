package cmd

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// setEnv sets env vars for the duration of a test and restores them after.
func setEnv(t *testing.T, pairs ...string) {
	t.Helper()
	for i := 0; i < len(pairs); i += 2 {
		t.Setenv(pairs[i], pairs[i+1]) //nolint:gosec // G602: callers always pass an even number of key/value args
	}
}

// allConfigKeys is the set of env vars that parseConfig reads.
var allConfigKeys = []string{
	"PAPERLESS_URL",
	"PAPERLESS_API_TOKEN",
}

// clearConfigEnv unsets all config env vars for the duration of the test.
// Required for tests running in environments where these may already be set.
func clearConfigEnv(t *testing.T) {
	t.Helper()
	for _, k := range allConfigKeys {
		orig, had := os.LookupEnv(k)
		_ = os.Unsetenv(k)
		k := k
		t.Cleanup(func() {
			if had {
				_ = os.Setenv(k, orig)
			} else {
				_ = os.Unsetenv(k)
			}
		})
	}
}

func TestParseConfig_MissingURL(t *testing.T) {
	clearConfigEnv(t)
	_, err := parseConfig(nil)
	if err == nil {
		t.Fatal("expected error when PAPERLESS_URL is missing")
	}
}

func TestParseConfig_MissingToken(t *testing.T) {
	clearConfigEnv(t)
	setEnv(t, "PAPERLESS_URL", "http://paperless.local:8000")
	_, err := parseConfig(nil)
	if err == nil {
		t.Fatal("expected error when PAPERLESS_API_TOKEN is missing")
	}
}

func TestParseConfig_MinimalValid(t *testing.T) {
	clearConfigEnv(t)
	setEnv(
		t,
		"PAPERLESS_URL", "http://paperless.local:8000",
		"PAPERLESS_API_TOKEN", "tok123",
	)
	cfg, err := parseConfig(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.baseURL != "http://paperless.local:8000" {
		t.Errorf("baseURL = %q", cfg.baseURL)
	}
	if cfg.token != "tok123" {
		t.Errorf("token = %q", cfg.token)
	}
}

func TestParseConfig_TrailingSlashStripped(t *testing.T) {
	clearConfigEnv(t)
	setEnv(
		t,
		"PAPERLESS_URL", "http://paperless.local:8000/",
		"PAPERLESS_API_TOKEN", "tok",
	)
	cfg, err := parseConfig(nil)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.baseURL != "http://paperless.local:8000" {
		t.Errorf("trailing slash not stripped: %q", cfg.baseURL)
	}
}

func TestParseConfig_FileValuesUsedWhenEnvAbsent(t *testing.T) {
	clearConfigEnv(t)
	fileVals := map[string]string{
		"PAPERLESS_URL":       "http://from-file.local:8000",
		"PAPERLESS_API_TOKEN": "file-token",
	}
	cfg, err := parseConfig(fileVals)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.baseURL != "http://from-file.local:8000" {
		t.Errorf("baseURL from file: %q", cfg.baseURL)
	}
	if cfg.token != "file-token" {
		t.Errorf("token from file: %q", cfg.token)
	}
}

func TestParseConfig_EnvTakesPrecedenceOverFile(t *testing.T) {
	clearConfigEnv(t)
	setEnv(
		t,
		"PAPERLESS_URL", "http://env.local:8000",
		"PAPERLESS_API_TOKEN", "env-token",
	)
	fileVals := map[string]string{
		"PAPERLESS_URL":       "http://file.local:8000",
		"PAPERLESS_API_TOKEN": "file-token",
	}
	cfg, err := parseConfig(fileVals)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.baseURL != "http://env.local:8000" {
		t.Errorf("env should win over file, got: %q", cfg.baseURL)
	}
	if cfg.token != "env-token" {
		t.Errorf("env should win over file, got: %q", cfg.token)
	}
}

func TestReadConfigFile_ParsesKeyValuePairs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config")
	content := `
# comment line
PAPERLESS_URL=http://paperless.local:8000
PAPERLESS_API_TOKEN = mytoken
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	vals := readConfigFile(path)
	cases := map[string]string{
		"PAPERLESS_URL":       "http://paperless.local:8000",
		"PAPERLESS_API_TOKEN": "mytoken",
	}
	for k, want := range cases {
		if got := vals[k]; got != want {
			t.Errorf("%s = %q, want %q", k, got, want)
		}
	}
	if _, ok := vals["# comment line"]; ok {
		t.Error("comment should not be parsed as key")
	}
}

func TestReadConfigFile_MissingFileReturnsNil(t *testing.T) {
	vals := readConfigFile("/nonexistent/path/config")
	if vals != nil {
		t.Errorf("expected nil for missing file, got %v", vals)
	}
}

func TestReadConfigFile_WarnsOnInsecurePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config")
	if err := os.WriteFile(path, []byte("PAPERLESS_URL=http://x\n"), 0o644); err != nil { //nolint:gosec // G306: intentionally insecure permissions to exercise the warning path
		t.Fatal(err)
	}

	// Capture stderr
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	readConfigFile(path)

	_ = w.Close()
	os.Stderr = old

	out, _ := io.ReadAll(r)
	if !strings.Contains(string(out), "insecure permissions") {
		t.Errorf("expected insecure permissions warning, got: %q", string(out))
	}
}

func TestReadConfigFile_NoWarnOnSecurePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config")
	if err := os.WriteFile(path, []byte("PAPERLESS_URL=http://x\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	readConfigFile(path)

	_ = w.Close()
	os.Stderr = old

	out, _ := io.ReadAll(r)
	if strings.Contains(string(out), "insecure") {
		t.Errorf("unexpected warning for 0600 file: %q", string(out))
	}
}
