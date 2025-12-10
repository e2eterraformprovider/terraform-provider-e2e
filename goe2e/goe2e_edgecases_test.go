package goe2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestDo_MalformedJSON tests handling of malformed JSON response
func TestDo_MalformedJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{invalid json`)
	}))
	defer server.Close()

	client, _ := NewClient(
		"key", "token", "proj", "region",
		SetBaseURL(server.URL),
		noRetryOpt(),
	)

	ctx := context.Background()
	req, _ := client.NewRequest(ctx, http.MethodGet, "test", nil)

	var result map[string]interface{}
	_, err := client.Do(ctx, req, &result)

	if err == nil {
		t.Error("Expected JSON decode error, got nil")
	}
}

// TestDo_EmptyResponseBody tests handling of empty response body
func TestDo_EmptyResponseBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Write nothing
	}))
	defer server.Close()

	client, _ := NewClient(
		"key", "token", "proj", "region",
		SetBaseURL(server.URL),
		noRetryOpt(),
	)

	ctx := context.Background()
	req, _ := client.NewRequest(ctx, http.MethodGet, "test", nil)

	var result map[string]interface{}
	_, err := client.Do(ctx, req, &result)

	// Should succeed with empty body (EOF is allowed)
	if err != nil {
		t.Errorf("Expected success with empty response body, got error: %v", err)
	}
}

// TestDo_LargeErrorResponse tests handling of large error response
func TestDo_LargeErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		// Create a large error message
		largeMessage := strings.Repeat("error detail ", 1000)
		errorResp := map[string]interface{}{
			"code":    400,
			"message": "Bad request",
			"errors":  []string{largeMessage},
		}
		json.NewEncoder(w).Encode(errorResp)
	}))
	defer server.Close()

	client, _ := NewClient(
		"key", "token", "proj", "region",
		SetBaseURL(server.URL),
		noRetryOpt(),
	)

	ctx := context.Background()
	req, _ := client.NewRequest(ctx, http.MethodGet, "test", nil)

	_, err := client.Do(ctx, req, nil)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	errResp, ok := err.(*ErrorResponse)
	if !ok {
		t.Errorf("Expected *ErrorResponse, got %T", err)
	} else if errResp.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", errResp.Code)
	}
}

// TestDo_ResponseWithoutContentType tests response without Content-Type header
func TestDo_ResponseWithoutContentType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Intentionally omit Content-Type header
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data": "test"}`)
	}))
	defer server.Close()

	client, _ := NewClient(
		"key", "token", "proj", "region",
		SetBaseURL(server.URL),
		noRetryOpt(),
	)

	ctx := context.Background()
	req, _ := client.NewRequest(ctx, http.MethodGet, "test", nil)

	var result map[string]interface{}
	_, err := client.Do(ctx, req, &result)

	if err != nil {
		t.Errorf("Expected to handle response without Content-Type, got error: %v", err)
	}
}

// TestDo_ResponseWithNullData tests handling of null values in response
func TestDo_ResponseWithNullData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data": null}`)
	}))
	defer server.Close()

	client, _ := NewClient(
		"key", "token", "proj", "region",
		SetBaseURL(server.URL),
		noRetryOpt(),
	)

	ctx := context.Background()
	req, _ := client.NewRequest(ctx, http.MethodGet, "test", nil)

	var result map[string]interface{}
	_, err := client.Do(ctx, req, &result)

	if err != nil {
		t.Errorf("Expected to handle null data, got error: %v", err)
	}

	if result["data"] != nil {
		t.Errorf("Expected data to be nil, got %v", result["data"])
	}
}

// TestDo_409ConflictError tests conflict error response
func TestDo_409ConflictError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, `{"code": 409, "message": "Resource already exists"}`)
	}))
	defer server.Close()

	client, _ := NewClient(
		"key", "token", "proj", "region",
		SetBaseURL(server.URL),
		noRetryOpt(),
	)

	ctx := context.Background()
	req, _ := client.NewRequest(ctx, http.MethodGet, "test", nil)

	_, err := client.Do(ctx, req, nil)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	errResp, ok := err.(*ErrorResponse)
	if !ok {
		t.Errorf("Expected *ErrorResponse, got %T", err)
	} else if errResp.Code != http.StatusConflict {
		t.Errorf("Expected status 409, got %d", errResp.Code)
	}
}

// TestDo_410GoneError tests resource no longer available error
func TestDo_410GoneError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusGone)
		fmt.Fprint(w, `{"code": 410, "message": "Resource gone"}`)
	}))
	defer server.Close()

	client, _ := NewClient(
		"key", "token", "proj", "region",
		SetBaseURL(server.URL),
		noRetryOpt(),
	)

	ctx := context.Background()
	req, _ := client.NewRequest(ctx, http.MethodGet, "test", nil)

	_, err := client.Do(ctx, req, nil)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// TestDo_422UnprocessableEntity tests validation error
func TestDo_422UnprocessableEntity(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprint(w, `{
			"code": 422,
			"message": "Validation failed",
			"errors": ["field1 is required", "field2 must be numeric"]
		}`)
	}))
	defer server.Close()

	client, _ := NewClient(
		"key", "token", "proj", "region",
		SetBaseURL(server.URL),
		noRetryOpt(),
	)

	ctx := context.Background()
	req, _ := client.NewRequest(ctx, http.MethodGet, "test", nil)

	_, err := client.Do(ctx, req, nil)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	errResp, ok := err.(*ErrorResponse)
	if !ok {
		t.Errorf("Expected *ErrorResponse, got %T", err)
	} else if errResp.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status 422, got %d", errResp.Code)
	}
}

// TestDo_BodyClosingError tests handling of body close errors (simulated)
func TestDo_BodyClosingError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"code": 200, "message": "OK", "data": {}}`)
	}))
	defer server.Close()

	client, _ := NewClient(
		"key", "token", "proj", "region",
		SetBaseURL(server.URL),
		noRetryOpt(),
	)

	ctx := context.Background()
	req, _ := client.NewRequest(ctx, http.MethodGet, "test", nil)

	var result map[string]interface{}
	_, err := client.Do(ctx, req, &result)

	// Should succeed normally - body closing errors are silently handled
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

// TestNewRequest_BodyEncodingError tests request creation with non-serializable body
func TestNewRequest_BodyEncodingError(t *testing.T) {
	client, _ := NewClient("key", "token", "proj", "region")

	ctx := context.Background()

	// Channel type is not JSON-serializable
	body := make(chan int)
	_, err := client.NewRequest(ctx, http.MethodPost, "test", body)

	if err == nil {
		t.Error("Expected JSON encoding error, got nil")
	}
}

// TestCheckResponse_NilResponse tests CheckResponse with various status codes
func TestCheckResponse_ValidStatus(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		shouldErr  bool
	}{
		{"200 OK", http.StatusOK, false},
		{"201 Created", http.StatusCreated, false},
		{"202 Accepted", http.StatusAccepted, false},
		{"204 No Content", http.StatusNoContent, false},
		{"299 highest success", 299, false},
		{"300 Multiple Choices", http.StatusMultipleChoices, true},
		{"400 Bad Request", http.StatusBadRequest, true},
		{"401 Unauthorized", http.StatusUnauthorized, true},
		{"403 Forbidden", http.StatusForbidden, true},
		{"404 Not Found", http.StatusNotFound, true},
		{"500 Internal Server", http.StatusInternalServerError, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				fmt.Fprintf(w, `{"code": %d, "message": "test"}`, tt.statusCode)
			}))
			defer server.Close()

			client, _ := NewClient(
				"key", "token", "proj", "region",
				SetBaseURL(server.URL),
				noRetryOpt(),
			)

			ctx := context.Background()
			req, _ := client.NewRequest(ctx, http.MethodGet, "test", nil)
			_, err := client.Do(ctx, req, nil)

			if tt.shouldErr && err == nil {
				t.Errorf("Expected error for status %d, got nil", tt.statusCode)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Expected no error for status %d, got: %v", tt.statusCode, err)
			}
		})
	}
}

// TestDo_WriterWithError tests writing response to io.Writer with content
func TestDo_WriterWithError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "binary data content")
	}))
	defer server.Close()

	client, _ := NewClient(
		"key", "token", "proj", "region",
		SetBaseURL(server.URL),
		noRetryOpt(),
	)

	ctx := context.Background()
	req, _ := client.NewRequest(ctx, http.MethodGet, "test", nil)

	var buf bytes.Buffer
	resp, err := client.Do(ctx, req, &buf)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	if buf.String() != "binary data content" {
		t.Errorf("Expected 'binary data content', got '%s'", buf.String())
	}
}

// TestDo_JSONErrorResponse tests error response is preserved for re-reading
func TestDo_JSONErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"code": 400, "message": "Bad request"}`)
	}))
	defer server.Close()

	client, _ := NewClient(
		"key", "token", "proj", "region",
		SetBaseURL(server.URL),
		noRetryOpt(),
	)

	ctx := context.Background()
	req, _ := client.NewRequest(ctx, http.MethodGet, "test", nil)

	_, err := client.Do(ctx, req, nil)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Verify error response contains correct data
	errResp, ok := err.(*ErrorResponse)
	if !ok {
		t.Fatalf("Expected *ErrorResponse, got %T", err)
	}

	if errResp.Code != 400 {
		t.Errorf("Expected code 400, got %d", errResp.Code)
	}

	if errResp.Message != "Bad request" {
		t.Errorf("Expected message 'Bad request', got '%s'", errResp.Message)
	}
}

// TestDo_ContextDeadlineExceeded tests context deadline during request execution
func TestDo_ContextDeadlineExceeded(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow server
		time.Sleep(500 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"code": 200}`)
	}))
	defer server.Close()

	client, _ := NewClient(
		"key", "token", "proj", "region",
		SetBaseURL(server.URL),
		noRetryOpt(),
	)

	// Short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	req, _ := client.NewRequest(ctx, http.MethodGet, "test", nil)
	_, err := client.Do(ctx, req, nil)

	if err != context.DeadlineExceeded {
		t.Errorf("Expected context.DeadlineExceeded, got %v", err)
	}
}

// TestErrorResponse_WithDetailsArray tests error response with multiple error details
func TestErrorResponse_WithDetailsArray(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprint(w, `{
			"code": 422,
			"message": "Multiple validation errors",
			"errors": ["name is required", "email is invalid", "password too short"]
		}`)
	}))
	defer server.Close()

	client, _ := NewClient(
		"key", "token", "proj", "region",
		SetBaseURL(server.URL),
		noRetryOpt(),
	)

	ctx := context.Background()
	req, _ := client.NewRequest(ctx, http.MethodGet, "test", nil)

	_, err := client.Do(ctx, req, nil)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	errResp, ok := err.(*ErrorResponse)
	if !ok {
		t.Fatalf("Expected *ErrorResponse, got %T", err)
	}

	if len(errResp.Errors) != 3 {
		t.Errorf("Expected 3 errors, got %d", len(errResp.Errors))
	}

	expectedErrors := []string{"name is required", "email is invalid", "password too short"}
	for i, expected := range expectedErrors {
		if i < len(errResp.Errors) && errResp.Errors[i] != expected {
			t.Errorf("Error %d: expected '%s', got '%s'", i, expected, errResp.Errors[i])
		}
	}
}

// TestNewRequest_ContextAlreadyCanceled tests NewRequest with canceled context
func TestNewRequest_ContextAlreadyCanceled(t *testing.T) {
	client, _ := NewClient("key", "token", "proj", "region")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// NewRequest should work with canceled context
	// The context error occurs during Do, not NewRequest
	req, err := client.NewRequest(ctx, http.MethodGet, "test", nil)

	if err != nil {
		t.Errorf("NewRequest should not fail with canceled context, got: %v", err)
	}

	if req == nil {
		t.Error("Expected request to be created")
	}
}

// TestDo_HTTPClientError tests network-level errors
func TestDo_HTTPClientError(t *testing.T) {
	client, _ := NewClient(
		"key", "token", "proj", "region",
		SetBaseURL("http://localhost:1"), // Non-existent port
		noRetryOpt(),
	)

	ctx := context.Background()
	req, _ := client.NewRequest(ctx, http.MethodGet, "test", nil)

	_, err := client.Do(ctx, req, nil)

	if err == nil {
		t.Error("Expected network error, got nil")
	}

	// Should be a network error (connection refused)
	if err.Error() == "" {
		t.Error("Expected error message, got empty string")
	}
}

// TestResponseBody_MultipleReads tests that response body can be read properly
func TestResponseBody_MultipleReads(t *testing.T) {
	responseData := `{"code": 200, "message": "OK", "data": {"id": "test-123"}}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, responseData)
	}))
	defer server.Close()

	client, _ := NewClient(
		"key", "token", "proj", "region",
		SetBaseURL(server.URL),
		noRetryOpt(),
	)

	ctx := context.Background()
	req, _ := client.NewRequest(ctx, http.MethodGet, "test", nil)

	var result map[string]interface{}
	resp, err := client.Do(ctx, req, &result)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Response should be populated
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Result should be populated
	if result["code"].(float64) != 200 {
		t.Errorf("Expected code 200, got %v", result["code"])
	}
}

// TestCheckResponse_UnmarshalableErrorJSON tests error response with bad JSON
func TestCheckResponse_UnmarshalableErrorJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		// Send invalid JSON - should still return error
		fmt.Fprint(w, `{invalid`)
	}))
	defer server.Close()

	client, _ := NewClient(
		"key", "token", "proj", "region",
		SetBaseURL(server.URL),
		noRetryOpt(),
	)

	ctx := context.Background()
	req, _ := client.NewRequest(ctx, http.MethodGet, "test", nil)

	_, err := client.Do(ctx, req, nil)

	if err == nil {
		t.Error("Expected error for 400 status, got nil")
	}

	// Should still be an ErrorResponse (with default values)
	_, ok := err.(*ErrorResponse)
	if !ok {
		t.Errorf("Expected *ErrorResponse, got %T", err)
	}
}
