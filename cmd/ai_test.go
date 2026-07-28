package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stefanhoth/paperless-ngx-cli/api"
)

func TestSuggestionsEndpoint(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"title":"Invoice","correspondents":[3],"suggested_correspondents":["Amazon"],"tags":[],"suggested_tags":[],"document_types":[],"suggested_document_types":[],"storage_paths":[],"suggested_storage_paths":[],"dates":["2026-01-15"]}`))
	}))
	defer srv.Close()

	c, err := api.NewClientWithResponses(srv.URL)
	if err != nil {
		t.Fatalf("client error: %v", err)
	}
	resp, err := c.DocumentsAiSuggestionsRetrieveWithResponse(ctx(), 42)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	if gotPath != "/api/documents/42/ai_suggestions/" {
		t.Errorf("path = %q, want /api/documents/42/ai_suggestions/", gotPath)
	}
	if resp.JSON200 == nil || derefStr(resp.JSON200.Title) != "Invoice" {
		t.Errorf("unexpected suggestions body: %+v", resp.JSON200)
	}
}

func TestSuggestionsEndpoint_AIDisabled(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("AI is required for this feature"))
	}))
	defer srv.Close()

	c, err := api.NewClientWithResponses(srv.URL)
	if err != nil {
		t.Fatalf("client error: %v", err)
	}
	resp, err := c.DocumentsAiSuggestionsRetrieveWithResponse(ctx(), 42)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	if apiErr := apiError(resp.StatusCode(), resp.Body); apiErr == nil {
		t.Fatal("expected error for 400 response")
	}
}

func TestChatEndpoint(t *testing.T) {
	var gotBody api.ChatStreamingRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Errorf("decode request body: %v", err)
		}
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte("The total is 42.00 EUR."))
	}))
	defer srv.Close()

	c, err := api.NewClientWithResponses(srv.URL)
	if err != nil {
		t.Fatalf("client error: %v", err)
	}
	docID := 7
	resp, err := c.DocumentsChatCreateWithResponse(ctx(), api.ChatStreamingRequest{Q: "invoice total?", DocumentId: &docID})
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	if gotBody.Q != "invoice total?" || gotBody.DocumentId == nil || *gotBody.DocumentId != 7 {
		t.Errorf("unexpected request body: %+v", gotBody)
	}
	if string(resp.Body) != "The total is 42.00 EUR." {
		t.Errorf("body = %q", resp.Body)
	}
}
