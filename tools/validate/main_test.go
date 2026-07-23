package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestVerifySchemaURL(t *testing.T) {
	const schema = `{"type":"object"}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(schema))
	}))
	t.Cleanup(server.Close)

	client := &http.Client{Timeout: time.Second}
	if err := verifySchemaURL(client, server.URL, []byte(schema)); err != nil {
		t.Fatal(err)
	}
}

func TestVerifySchemaURLRejectsDifferentBytes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"type":"array"}`))
	}))
	t.Cleanup(server.Close)

	err := verifySchemaURL(&http.Client{Timeout: time.Second}, server.URL, []byte(`{"type":"object"}`))
	if err == nil || !strings.Contains(err.Error(), "differs") {
		t.Fatalf("error = %v", err)
	}
}

func TestSchemaHTTPClientRejectsUnexpectedRedirect(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "https://example.com/schema.json", nil)
	if err != nil {
		t.Fatal(err)
	}
	previous, err := http.NewRequest(http.MethodGet, "https://schemas.pawnkit.dev/schema.json", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := schemaHTTPClient().CheckRedirect(request, []*http.Request{previous}); err == nil {
		t.Fatal("unexpected redirect was accepted")
	}
}
