package goe2e

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
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

	client, _ := NewClient(
		"test-api-key",
		"test-auth-token",
		"test-project",
		"test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}), // Disable retries for tests
	)

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

// writeJSON writes JSON data to the response writer.
// Usage:
//   - writeJSON(w, jsonData) - uses default http.StatusOK
//   - writeJSON(w, statusCode, jsonData) - uses specified status code
func writeJSON(w http.ResponseWriter, args ...interface{}) {
	w.Header().Set("Content-Type", "application/json")

	var statusCode int = http.StatusOK
	var jsonData string

	if len(args) == 1 {
		// Single argument: just JSON data
		jsonData = args[0].(string)
	} else if len(args) == 2 {
		// Two arguments: status code and JSON data
		statusCode = args[0].(int)
		jsonData = args[1].(string)
	} else {
		panic("writeJSON: invalid number of arguments")
	}

	// Check if header was already written by using reflection to check ResponseRecorder.Code
	// This avoids the "superfluous response.WriteHeader call" warning
	rv := reflect.ValueOf(w)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() == reflect.Struct {
		codeField := rv.FieldByName("Code")
		if codeField.IsValid() && codeField.Kind() == reflect.Int {
			if codeField.Int() == 0 {
				w.WriteHeader(statusCode)
			}
			// If Code is already set, don't call WriteHeader again to avoid warning
		} else {
			w.WriteHeader(statusCode)
		}
	} else {
		w.WriteHeader(statusCode)
	}
	_, _ = fmt.Fprint(w, jsonData)
}

// noRetryOpt returns a ClientOpt that disables retries for tests
func noRetryOpt() ClientOpt {
	return WithRetryAndBackoffs(RetryConfig{RetryMax: 0})
}

// TestNewClient_RequiredParameters tests that NewClient validates required parameters
func TestNewClient_RequiredParameters(t *testing.T) {
	tests := []struct {
		name      string
		apiKey    string
		authToken string
		projectID string
		region    string
		wantErr   bool
	}{
		{
			name:      "all parameters provided",
			apiKey:    "key",
			authToken: "token",
			projectID: "proj",
			region:    "Mumbai",
			wantErr:   false,
		},
		{
			name:      "missing apiKey",
			apiKey:    "",
			authToken: "token",
			projectID: "proj",
			region:    "Mumbai",
			wantErr:   true,
		},
		{
			name:      "missing authToken",
			apiKey:    "key",
			authToken: "",
			projectID: "proj",
			region:    "Mumbai",
			wantErr:   true,
		},
		{
			name:      "missing projectID",
			apiKey:    "key",
			authToken: "token",
			projectID: "",
			region:    "Mumbai",
			wantErr:   true,
		},
		{
			name:      "missing region",
			apiKey:    "key",
			authToken: "token",
			projectID: "proj",
			region:    "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.apiKey, tt.authToken, tt.projectID, tt.region)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
				if client == nil {
					t.Error("Expected client to be created")
				}
			}
		})
	}
}

// TestNewClient_WithOptionalParameters tests optional parameters like workspace and team
func TestNewClient_WithOptionalParameters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify standard parameters
		if got := r.URL.Query().Get("apikey"); got != "test-key" {
			t.Errorf("apikey = %s, want test-key", got)
		}
		if got := r.URL.Query().Get("project_id"); got != "test-project" {
			t.Errorf("project_id = %s, want test-project", got)
		}
		if got := r.URL.Query().Get("location"); got != "Mumbai" {
			t.Errorf("location = %s, want Mumbai", got)
		}
		// Verify optional parameters
		if got := r.URL.Query().Get("workspace_id"); got != "workspace-123" {
			t.Errorf("workspace_id = %s, want workspace-123", got)
		}
		if got := r.URL.Query().Get("team_id"); got != "team-456" {
			t.Errorf("team_id = %s, want team-456", got)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code": 200, "message": "OK", "data": {"name": "test"}}`))
	}))
	defer server.Close()

	// Create client with optional parameters
	client, err := NewClient(
		"test-key",
		"test-token",
		"test-project",
		"Mumbai",
		WithWorkspace("workspace-123"),
		WithTeam("team-456"),
		SetBaseURL(server.URL),
	)
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}

	// Make a test request
	ctx := context.Background()
	req, err := client.NewRequest(ctx, http.MethodGet, "test", nil)
	if err != nil {
		t.Fatalf("NewRequest returned error: %v", err)
	}

	resp, err := client.Do(ctx, req, nil)
	if err != nil {
		t.Fatalf("Do returned error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Response status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

// TestNewClient_WithoutOptionalParameters tests that optional parameters are not added when not set
func TestNewClient_WithoutOptionalParameters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify optional parameters are NOT present
		if got := r.URL.Query().Get("workspace_id"); got != "" {
			t.Errorf("workspace_id should be empty, got %s", got)
		}
		if got := r.URL.Query().Get("team_id"); got != "" {
			t.Errorf("team_id should be empty, got %s", got)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code": 200, "message": "OK", "data": {}}`))
	}))
	defer server.Close()

	// Create client WITHOUT optional parameters
	client, err := NewClient(
		"test-key",
		"test-token",
		"test-project",
		"Mumbai",
		SetBaseURL(server.URL),
	)
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}

	// Make a test request
	ctx := context.Background()
	req, err := client.NewRequest(ctx, http.MethodGet, "test", nil)
	if err != nil {
		t.Fatalf("NewRequest returned error: %v", err)
	}

	_, err = client.Do(ctx, req, nil)
	if err != nil {
		t.Fatalf("Do returned error: %v", err)
	}
}

// TestSetUserAgent tests the SetUserAgent option
func TestSetUserAgent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("User-Agent"); got != "test-agent/1.0" {
			t.Errorf("User-Agent = %s, want test-agent/1.0", got)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClient(
		"key", "token", "proj", "region",
		SetBaseURL(server.URL),
		SetUserAgent("test-agent/1.0"),
	)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	ctx := context.Background()
	req, _ := client.NewRequest(ctx, http.MethodGet, "test", nil)
	_, _ = client.Do(ctx, req, nil)
}

// TestWithRetryAndBackoffs tests retry configuration
func TestWithRetryAndBackoffs(t *testing.T) {
	retryConfig := RetryConfig{
		RetryMax:     10,
		RetryWaitMin: PtrTo(2 * time.Second),
		RetryWaitMax: PtrTo(30 * time.Second),
	}

	client, err := NewClient(
		"key", "token", "proj", "region",
		WithRetryAndBackoffs(retryConfig),
	)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	if client.RetryConfig.RetryMax != 10 {
		t.Errorf("Expected RetryMax 10, got %d", client.RetryConfig.RetryMax)
	}

	if client.RetryConfig.RetryWaitMin == nil || *client.RetryConfig.RetryWaitMin != 2*time.Second {
		t.Error("Expected RetryWaitMin to be 2 seconds")
	}
}

// TestDo_ErrorResponse tests error response handling
func TestDo_ErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusBadRequest, `{"message": "Bad request", "code": 400}`)
	}))
	defer server.Close()

	client, _ := NewClient(
		"key", "token", "proj", "region",
		SetBaseURL(server.URL),
	)

	ctx := context.Background()
	req, _ := client.NewRequest(ctx, http.MethodGet, "test", nil)
	_, err := client.Do(ctx, req, nil)

	if err == nil {
		t.Fatal("Expected error for 400 response, got nil")
	}

	// Should be an ErrorResponse
	if _, ok := err.(*ErrorResponse); !ok {
		t.Errorf("Expected *ErrorResponse, got %T", err)
	}
}

// TestNew_InvalidBaseURL tests error handling for invalid base URL
func TestNew_InvalidBaseURL(t *testing.T) {
	_, err := NewClient(
		"key", "token", "proj", "region",
		SetBaseURL("://invalid-url"),
	)

	if err == nil {
		t.Fatal("Expected error for invalid base URL, got nil")
	}
}

// TestSetStaticRateLimit tests rate limiting configuration (placeholder function)
func TestSetStaticRateLimit(t *testing.T) {
	client, err := NewClient(
		"key", "token", "proj", "region",
		SetStaticRateLimit(10.0),
	)

	if err != nil {
		t.Fatalf("NewClient with SetStaticRateLimit returned error: %v", err)
	}

	if client == nil {
		t.Fatal("Expected client to be created")
	}
}

// TestNewRequest_WithBody tests NewRequest with a request body
func TestNewRequest_WithBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	ctx := context.Background()

	body := map[string]string{"key": "value"}
	req, err := client.NewRequest(ctx, http.MethodPost, "test", body)
	if err != nil {
		t.Fatalf("NewRequest returned error: %v", err)
	}

	if req.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", req.Header.Get("Content-Type"))
	}
}

// TestDo_WithWriter tests Do method with io.Writer response
func TestDo_WithWriter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	}))
	defer server.Close()

	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	ctx := context.Background()

	req, _ := client.NewRequest(ctx, http.MethodGet, "test", nil)

	var buf bytes.Buffer
	_, err := client.Do(ctx, req, &buf)
	if err != nil {
		t.Fatalf("Do returned error: %v", err)
	}

	if buf.String() != "test response" {
		t.Errorf("Expected 'test response', got %s", buf.String())
	}
}

// TestDo_ContextCanceled tests Do with canceled context
func TestDo_ContextCanceled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	req, _ := client.NewRequest(ctx, http.MethodGet, "test", nil)
	_, err := client.Do(ctx, req, nil)

	if err != context.Canceled {
		t.Errorf("Expected context.Canceled error, got %v", err)
	}
}

// TestNew_MissingRequiredFields tests New without required parameters
func TestNew_MissingRequiredFields(t *testing.T) {
	_, err := New(nil)
	if err == nil {
		t.Fatal("Expected error when creating client without required fields")
	}
}
