package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNormalizeAPIPath(t *testing.T) {
	base := "http://paperless.local:8000"
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"leading slash optional", "documents/4028/", "/documents/4028/"},
		{"trailing slash added", "/documents/4028", "/documents/4028/"},
		{"already normalized", "/documents/4028/", "/documents/4028/"},
		{"query string preserved", "/documents/?created__date=2026-07-08", "/documents/?created__date=2026-07-08"},
		{"trailing slash added before query", "/documents?created__date=2026-07-08", "/documents/?created__date=2026-07-08"},
		{"matching full URL reduced to path", base + "/documents/4028/", "/documents/4028/"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := normalizeAPIPath(tc.in, base)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("normalizeAPIPath(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestNormalizeAPIPath_MismatchedOrigin(t *testing.T) {
	_, err := normalizeAPIPath("http://other.example.com/documents/4028/", "http://paperless.local:8000")
	if err == nil {
		t.Fatal("expected error for URL with non-matching origin")
	}
}

func TestBuildFieldsBody(t *testing.T) {
	body, err := buildFieldsBody([]string{"created=2022-02-08", "title=Foo"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := `{"created":"2022-02-08","title":"Foo"}`
	if string(body) != want {
		t.Errorf("buildFieldsBody() = %s, want %s", body, want)
	}
}

func TestBuildFieldsBody_InvalidField(t *testing.T) {
	_, err := buildFieldsBody([]string{"noequalssign"})
	if err == nil {
		t.Fatal("expected error for field without '='")
	}
}

func TestReadInputBody_File(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "body.json")
	if err := os.WriteFile(path, []byte(`{"created":"2022-02-08"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := readInputBody(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(got) != `{"created":"2022-02-08"}` {
		t.Errorf("readInputBody() = %s", got)
	}
}
