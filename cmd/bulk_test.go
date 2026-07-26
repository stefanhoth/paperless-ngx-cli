package cmd

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stefanhoth/paperless-ngx-cli/api"
)

// TestRunBulk_Endpoints pins each operation to the endpoint it uses on API
// v10: the document actions have dedicated endpoints, the metadata operations
// go through bulk_edit.
func TestRunBulk_Endpoints(t *testing.T) {
	cases := []struct {
		name     string
		args     []string
		wantPath string
		wantBody map[string]any
	}{
		{
			"reprocess",
			[]string{"reprocess", "1,2"},
			"/api/documents/reprocess/",
			map[string]any{"documents": []any{1.0, 2.0}},
		},
		{
			"delete",
			[]string{"delete", "42"},
			"/api/documents/delete/",
			map[string]any{"documents": []any{42.0}},
		},
		{
			"merge",
			[]string{"merge", "1,2"},
			"/api/documents/merge/",
			map[string]any{"documents": []any{1.0, 2.0}},
		},
		{
			"rotate",
			[]string{"rotate", "99", "90"},
			"/api/documents/rotate/",
			map[string]any{"documents": []any{99.0}, "degrees": 90.0},
		},
		{
			"add-tag",
			[]string{"add-tag", "10,11", "7"},
			"/api/documents/bulk_edit/",
			map[string]any{"documents": []any{10.0, 11.0}, "method": "add_tag", "parameters": map[string]any{"tag": 7.0}},
		},
		{
			"remove-tag",
			[]string{"remove-tag", "10", "7"},
			"/api/documents/bulk_edit/",
			map[string]any{"documents": []any{10.0}, "method": "remove_tag", "parameters": map[string]any{"tag": 7.0}},
		},
		{
			"set-correspondent",
			[]string{"set-correspondent", "5", "3"},
			"/api/documents/bulk_edit/",
			map[string]any{"documents": []any{5.0}, "method": "set_correspondent", "parameters": map[string]any{"correspondent": 3.0}},
		},
		{
			"set-type",
			[]string{"set-type", "5", "3"},
			"/api/documents/bulk_edit/",
			map[string]any{"documents": []any{5.0}, "method": "set_document_type", "parameters": map[string]any{"document_type": 3.0}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var gotPath string
			gotBody := map[string]any{}
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath = r.URL.Path
				b, _ := io.ReadAll(r.Body)
				if err := json.Unmarshal(b, &gotBody); err != nil {
					t.Errorf("request body is not JSON: %v", err)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"result":"OK"}`))
			}))
			defer srv.Close()

			c, err := api.NewClientWithResponses(srv.URL)
			if err != nil {
				t.Fatalf("client error: %v", err)
			}
			if err := runBulk(c, tc.args[0], parseIDs(tc.args[1]), tc.args); err != nil {
				t.Fatalf("runBulk() error: %v", err)
			}
			if gotPath != tc.wantPath {
				t.Errorf("path = %q, want %q", gotPath, tc.wantPath)
			}
			if !reflect.DeepEqual(gotBody, tc.wantBody) {
				t.Errorf("body = %v, want %v", gotBody, tc.wantBody)
			}
		})
	}
}

func TestRunBulk_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"detail":"nope"}`))
	}))
	defer srv.Close()

	c, err := api.NewClientWithResponses(srv.URL)
	if err != nil {
		t.Fatalf("client error: %v", err)
	}
	if err := runBulk(c, "reprocess", []int{1}, []string{"reprocess", "1"}); err == nil {
		t.Fatal("expected error for 400 response")
	}
}

func TestParseIDs(t *testing.T) {
	cases := []struct {
		in   string
		want []int
	}{
		{"1,2,3", []int{1, 2, 3}},
		{" 1 , 2 ", []int{1, 2}},
		{"1,x,3", []int{1, 3}},
		{"nope", nil},
	}
	for _, tc := range cases {
		if got := parseIDs(tc.in); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("parseIDs(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}
