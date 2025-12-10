package goe2e

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

// ============================================================================
// Network Error Helpers
// ============================================================================

// newNetworkErrorClient creates a client configured for network error testing.
// It uses a non-existent endpoint (localhost:1) to simulate network failures.
func newNetworkErrorClient(t *testing.T) *Client {
	t.Helper()
	client, err := NewClient(
		"test-key",
		"test-token",
		"test-proj",
		"test-region",
		SetBaseURL("http://localhost:1"), // Non-existent endpoint
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}), // No retries for faster tests
	)
	if err != nil {
		t.Fatalf("Failed to create network error client: %v", err)
	}
	return client
}

// newCanceledClient creates a client with a pre-canceled context.
// Returns the client and canceled context.
func newCanceledClient(t *testing.T) (*Client, context.Context) {
	t.Helper()
	client, err := NewClient(
		"test-key",
		"test-token",
		"test-proj",
		"test-region",
		SetBaseURL("http://localhost:8080"),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}),
	)
	if err != nil {
		t.Fatalf("Failed to create canceled client: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately
	return client, ctx
}

// ============================================================================
// Mock Server Builders
// ============================================================================

// newTimeoutServer builds a mock server that delays responses by the specified duration.
func newTimeoutServer(t *testing.T, delay time.Duration) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		writeJSON(w, http.StatusOK, `{"code": 200, "message": "success"}`)
	}))
}

// newErrorServer builds a mock server that always returns the specified error code.
func newErrorServer(t *testing.T, errorCode int, errorMessage string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, errorCode, fmt.Sprintf(`{"code": %d, "message": "%s"}`, errorCode, errorMessage))
	}))
}

// newSuccessServer builds a mock server that returns a success response with custom data.
func newSuccessServer(t *testing.T, data interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"code":    200,
			"message": "success",
			"data":    data,
		}
		jsonData, _ := json.Marshal(response)
		writeJSON(w, http.StatusOK, string(jsonData))
	}))
}

// newConditionalServer builds a mock server that returns different responses based on request properties.
type RequestCondition func(r *http.Request) bool
type ConditionalResponse struct {
	Condition  RequestCondition
	StatusCode int
	Response   string
}

// ============================================================================
// Response Assertion Helpers
// ============================================================================

// assertError checks that err is not nil and optionally contains expected text.
func assertError(t *testing.T, err error, expectedSubstring string) {
	t.Helper()
	if err == nil {
		t.Error("Expected error but got nil")
		return
	}
	if expectedSubstring != "" && !strings.Contains(err.Error(), expectedSubstring) {
		t.Errorf("Expected error to contain '%s', got: %v", expectedSubstring, err)
	}
}

// assertNoError checks that err is nil.
func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}
}

// assertNil checks that a value is nil.
func assertNil(t *testing.T, val interface{}, msg string) {
	t.Helper()
	if val != nil {
		// Check if it's a typed nil (e.g., (*Type)(nil))
		if rv := reflect.ValueOf(val); rv.Kind() == reflect.Ptr && rv.IsNil() {
			return // Typed nil is acceptable
		}
		t.Errorf("%s: expected nil but got %v", msg, val)
	}
}

// assertNotNil checks that a value is not nil.
func assertNotNil(t *testing.T, val interface{}, msg string) {
	t.Helper()
	if val == nil {
		t.Errorf("%s: expected non-nil value", msg)
	}
}

// assertStatus checks that response status code matches expected.
func assertStatus(t *testing.T, resp *Response, expectedCode int) {
	t.Helper()
	if resp == nil {
		t.Error("Response is nil")
		return
	}
	if resp.StatusCode != expectedCode {
		t.Errorf("Expected status code %d, got %d", expectedCode, resp.StatusCode)
	}
}

// assertErrorType checks that error is of specific type.
func assertErrorType(t *testing.T, err error, expectedType interface{}) {
	t.Helper()
	if err == nil {
		t.Error("Expected error but got nil")
		return
	}
	switch expectedType.(type) {
	case *ErrorResponse:
		if _, ok := err.(*ErrorResponse); !ok {
			t.Errorf("Expected *ErrorResponse, got %T", err)
		}
	case *ArgError:
		if _, ok := err.(*ArgError); !ok {
			t.Errorf("Expected *ArgError, got %T", err)
		}
	default:
		t.Errorf("Unknown error type for assertion: %T", expectedType)
	}
}

// ============================================================================
// Test Data Builders
// ============================================================================

// buildSuccessResponse creates a standard success response JSON string.
func buildSuccessResponse(code int, message string, data interface{}) string {
	response := map[string]interface{}{
		"code":    code,
		"message": message,
		"data":    data,
	}
	jsonData, _ := json.Marshal(response)
	return string(jsonData)
}

// buildErrorResponse creates a standard error response JSON string.
func buildErrorResponse(code int, message string, errors []string) string {
	response := map[string]interface{}{
		"code":    code,
		"message": message,
	}
	if len(errors) > 0 {
		response["errors"] = errors
	}
	jsonData, _ := json.Marshal(response)
	return string(jsonData)
}

// buildListResponse creates a standard list response with items.
func buildListResponse(items []interface{}) string {
	response := map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    items,
	}
	jsonData, _ := json.Marshal(response)
	return string(jsonData)
}

// ============================================================================
// HTTP Request Helpers
// ============================================================================
// Note: testMethod and testQueryParam already exist in goe2e_test.go
// These are aliases with more descriptive names for consistency

// assertHTTPMethod is an alias for testMethod with more descriptive naming.
func assertHTTPMethod(t *testing.T, r *http.Request, expectedMethod string) {
	t.Helper()
	testMethod(t, r, expectedMethod)
}

// assertQueryParam is an alias for testQueryParam with more descriptive naming.
func assertQueryParam(t *testing.T, r *http.Request, key, expectedValue string) {
	t.Helper()
	testQueryParam(t, r, key, expectedValue)
}

// ============================================================================
// Common Test Patterns
// ============================================================================

// testNetworkError is a helper for testing network errors across all services.
func testNetworkError(t *testing.T, testFunc func(*Client, context.Context) error) {
	t.Helper()
	client := newNetworkErrorClient(t)
	err := testFunc(client, context.Background())
	assertError(t, err, "")
}

// testContextTimeout is a helper for testing context timeout scenarios.
func testContextTimeout(t *testing.T, testFunc func(*Client, context.Context) error) {
	t.Helper()
	server := newTimeoutServer(t, 500*time.Millisecond)
	defer server.Close()

	client, err := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err = testFunc(client, ctx)
	assertError(t, err, "")
}

// testContextCancellation is a helper for testing context cancellation.
func testContextCancellation(t *testing.T, testFunc func(*Client, context.Context) error) {
	t.Helper()
	client, ctx := newCanceledClient(t)
	err := testFunc(client, ctx)
	assertError(t, err, "")
}
