package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	apiCmd.Flags().StringP("method", "X", "", "HTTP method (default GET, or POST when a body is supplied)")
	apiCmd.Flags().StringArrayP("field", "f", nil, "Add a string field to the JSON request body (key=value)")
	apiCmd.Flags().String("input", "", "Read the raw request body from a file, or - for stdin")
	rootCmd.AddCommand(apiCmd)
}

var apiCmd = &cobra.Command{
	Use:   "api <path>",
	Short: "Raw REST API passthrough, similar to gh api",
	Long: `Make an authenticated request to the Paperless-NGX REST API and print the raw JSON response to stdout.

Examples:
  paperless api /documents/4028/ --method PATCH --field created=2022-02-08
  paperless api /documents/4028/ --method PATCH --input body.json
  paperless api "/documents/?created__date=2026-07-08" | jq '.results[].id'
  echo '{"created":"2022-02-08"}' | paperless api /documents/4028/ -X PATCH --input -`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		method, _ := cmd.Flags().GetString("method")
		fields, _ := cmd.Flags().GetStringArray("field")
		input, _ := cmd.Flags().GetString("input")

		if err := validateBodyFlags(fields, input); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}

		var body []byte
		var err error
		switch {
		case len(fields) > 0:
			body, err = buildFieldsBody(fields)
		case input != "":
			body, err = readInputBody(input)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}

		method = resolveMethod(method, body != nil)

		cfg := loadConfig()
		path, err := normalizeAPIPath(args[0], cfg.baseURL)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}

		var reqBody io.Reader
		if body != nil {
			reqBody = bytes.NewReader(body)
		}
		req, err := http.NewRequestWithContext(ctx(), strings.ToUpper(method), cfg.baseURL+path, reqBody)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		setAuthHeaders(req, cfg)
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		os.Stdout.Write(respBody)

		if resp.StatusCode >= 400 {
			fmt.Fprintf(os.Stderr, "error: HTTP %d\n", resp.StatusCode)
			os.Exit(1)
		}
	},
}

// validateBodyFlags rejects specifying both --field and --input.
func validateBodyFlags(fields []string, input string) error {
	if len(fields) > 0 && input != "" {
		return fmt.Errorf("--field and --input are mutually exclusive")
	}
	return nil
}

// resolveMethod applies gh-api-style method defaulting: an explicit method
// always wins; otherwise POST if a body was supplied, GET if not.
func resolveMethod(explicit string, hasBody bool) string {
	if explicit != "" {
		return explicit
	}
	if hasBody {
		return http.MethodPost
	}
	return http.MethodGet
}

// normalizeAPIPath turns a user-supplied path or URL into a Paperless API path:
// leading slash is added if missing, a trailing slash is added before any
// query string, and a full URL matching baseURL's origin is reduced to its path.
func normalizeAPIPath(raw string, baseURL string) (string, error) {
	if u, err := url.Parse(raw); err == nil && u.Scheme != "" && u.Host != "" {
		base, err := url.Parse(baseURL)
		if err != nil || u.Scheme != base.Scheme || u.Host != base.Host {
			return "", fmt.Errorf("URL %s does not match configured PAPERLESS_URL %s", raw, baseURL)
		}
		raw = u.Path
		if u.RawQuery != "" {
			raw += "?" + u.RawQuery
		}
	}

	path, query, hasQuery := strings.Cut(raw, "?")
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	if hasQuery {
		path += "?" + query
	}
	return path, nil
}

// buildFieldsBody builds a flat JSON object of string values from "key=value" pairs.
func buildFieldsBody(fields []string) ([]byte, error) {
	m := make(map[string]string, len(fields))
	for _, f := range fields {
		key, value, ok := strings.Cut(f, "=")
		if !ok {
			return nil, fmt.Errorf("invalid --field %q: expected key=value", f)
		}
		m[key] = value
	}
	return json.Marshal(m)
}

// readInputBody reads the raw request body from a file, or from stdin when input is "-".
func readInputBody(input string) ([]byte, error) {
	if input == "-" {
		return io.ReadAll(os.Stdin)
	}
	return os.ReadFile(input)
}
