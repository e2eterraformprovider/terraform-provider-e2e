package goe2e

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// testServer holds the test HTTP server and mux for testing
type testServer struct {
	server *httptest.Server
	mux    *http.ServeMux
	client *Client
}

// setup creates a test HTTP server and returns a testServer struct
func setup() *testServer {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	client, _ := NewClient("test-api-key", "test-auth-token", SetBaseURL(server.URL))

	return &testServer{
		server: server,
		mux:    mux,
		client: client,
	}
}

// teardown closes the test HTTP server
func (ts *testServer) teardown() {
	ts.server.Close()
}

// testMethod verifies the HTTP method of a request
func testMethod(t *testing.T, r *http.Request, expected string) {
	t.Helper()
	if got := r.Method; got != expected {
		t.Errorf("Request method: %v, expected %v", got, expected)
	}
}

// testURLPath verifies the URL path of a request
func testURLPath(t *testing.T, r *http.Request, expected string) {
	t.Helper()
	if got := r.URL.Path; got != expected {
		t.Errorf("Request URL path: %v, expected %v", got, expected)
	}
}

// testQueryParam verifies a specific query parameter value
func testQueryParam(t *testing.T, r *http.Request, param, expected string) {
	t.Helper()
	if got := r.URL.Query().Get(param); got != expected {
		t.Errorf("Query param %s: %v, expected %v", param, got, expected)
	}
}

// writeJSON writes a JSON response with the given status code
func writeJSON(w http.ResponseWriter, statusCode int, jsonData string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = fmt.Fprint(w, jsonData)
}
