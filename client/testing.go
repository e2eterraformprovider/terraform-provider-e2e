package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

// testServer holds the test HTTP server and mux
type testServer struct {
	server *httptest.Server
	mux    *http.ServeMux
	client *Client
}

// setup creates a test HTTP server and returns a testServer struct
func setup() *testServer {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	client := NewClient("test-api-key", "test-auth-token", server.URL)

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

// testHeader verifies a specific header value
func testHeader(t *testing.T, r *http.Request, header, expected string) {
	t.Helper()
	if got := r.Header.Get(header); got != expected {
		t.Errorf("Header %s: %v, expected %v", header, got, expected)
	}
}

// testQueryParam verifies a specific query parameter value
func testQueryParam(t *testing.T, r *http.Request, param, expected string) {
	t.Helper()
	if got := r.URL.Query().Get(param); got != expected {
		t.Errorf("Query param %s: %v, expected %v", param, got, expected)
	}
}

// testFormValues verifies multiple form values
func testFormValues(t *testing.T, r *http.Request, values url.Values) {
	t.Helper()
	if err := r.ParseForm(); err != nil {
		t.Fatalf("Error parsing form: %v", err)
	}

	for key, want := range values {
		if got := r.Form[key]; !reflect.DeepEqual(got, want) {
			t.Errorf("Form value %s: %v, expected %v", key, got, want)
		}
	}
}

// testDeepEqual compares two values using reflect.DeepEqual
func testDeepEqual(t *testing.T, got, want interface{}, context string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s:\nGot:  %+v\nWant: %+v", context, got, want)
	}
}

// testErrorContains checks if an error contains a specific substring
func testErrorContains(t *testing.T, err error, want string) {
	t.Helper()
	if err == nil {
		t.Errorf("Expected error containing %q, got nil", want)
		return
	}
	if got := err.Error(); !contains(got, want) {
		t.Errorf("Error %q does not contain %q", got, want)
	}
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}

// writeJSON writes a JSON response with the given status code
func writeJSON(w http.ResponseWriter, statusCode int, jsonData string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = fmt.Fprint(w, jsonData)
}

// writeError writes an error response
func writeError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = fmt.Fprintf(w, `{"error": %q, "code": %d}`, message, statusCode)
}
