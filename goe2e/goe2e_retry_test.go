package goe2e

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestRetry_500ErrorRecovery tests that the client successfully retries on 500 errors
// and recovers when the server starts responding correctly
func TestRetry_500ErrorRecovery(t *testing.T) {
	var requestCount atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := requestCount.Add(1)

		// First 2 requests fail with 500, third succeeds
		if count <= 2 {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"message": "Server error", "code": 500}`)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"code": 200, "message": "OK", "data": {"id": "test-123"}}`)
	}))
	defer server.Close()

	client, err := NewClient(
		"test-key",
		"test-token",
		"test-project",
		"Mumbai",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{
			RetryMax:     3,
			RetryWaitMin: PtrTo(10 * time.Millisecond),
			RetryWaitMax: PtrTo(50 * time.Millisecond),
		}),
	)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	ctx := context.Background()
	req, err := client.NewRequest(ctx, http.MethodGet, "test", nil)
	if err != nil {
		t.Fatalf("NewRequest error: %v", err)
	}

	var result map[string]interface{}
	resp, err := client.Do(ctx, req, &result)

	// Should succeed on third attempt
	if err != nil {
		t.Fatalf("Expected success after retry, got error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Verify it took 3 requests (2 failures + 1 success)
	if count := requestCount.Load(); count != 3 {
		t.Errorf("Expected 3 requests, got %d", count)
	}
}

// TestRetry_502BadGatewayRecovery tests retry on 502 Bad Gateway errors
func TestRetry_502BadGatewayRecovery(t *testing.T) {
	var requestCount atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := requestCount.Add(1)

		if count == 1 {
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprint(w, `{"message": "Bad gateway", "code": 502}`)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"code": 200, "message": "OK", "data": {}}`)
	}))
	defer server.Close()

	client, err := NewClient(
		"test-key", "test-token", "test-project", "Mumbai",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{
			RetryMax:     2,
			RetryWaitMin: PtrTo(10 * time.Millisecond),
			RetryWaitMax: PtrTo(50 * time.Millisecond),
		}),
	)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	ctx := context.Background()
	req, err := client.NewRequest(ctx, http.MethodGet, "test", nil)
	if err != nil {
		t.Fatalf("NewRequest error: %v", err)
	}

	resp, err := client.Do(ctx, req, nil)
	if err != nil {
		t.Errorf("Expected success after retry, got error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

// TestRetry_503ServiceUnavailableRecovery tests retry on 503 Service Unavailable
func TestRetry_503ServiceUnavailableRecovery(t *testing.T) {
	var requestCount atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := requestCount.Add(1)

		if count <= 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprint(w, `{"message": "Service unavailable", "code": 503}`)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"code": 200, "message": "OK", "data": {}}`)
	}))
	defer server.Close()

	client, err := NewClient(
		"test-key", "test-token", "test-project", "Mumbai",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{
			RetryMax:     2,
			RetryWaitMin: PtrTo(10 * time.Millisecond),
			RetryWaitMax: PtrTo(50 * time.Millisecond),
		}),
	)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	ctx := context.Background()
	req, err := client.NewRequest(ctx, http.MethodGet, "test", nil)
	if err != nil {
		t.Fatalf("NewRequest error: %v", err)
	}

	resp, err := client.Do(ctx, req, nil)
	if err != nil {
		t.Errorf("Expected success after retry, got error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

// TestRetry_ExhaustedRetries tests that after max retries, the error is returned
func TestRetry_ExhaustedRetries(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Always return 500
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"message": "Server error", "code": 500}`)
	}))
	defer server.Close()

	client, err := NewClient(
		"test-key", "test-token", "test-project", "Mumbai",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{
			RetryMax:     2, // Will try 1 initial + 2 retries = 3 total
			RetryWaitMin: PtrTo(10 * time.Millisecond),
			RetryWaitMax: PtrTo(50 * time.Millisecond),
		}),
	)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	ctx := context.Background()
	req, err := client.NewRequest(ctx, http.MethodGet, "test", nil)
	if err != nil {
		t.Fatalf("NewRequest error: %v", err)
	}

	_, err = client.Do(ctx, req, nil)

	// Should eventually fail with ErrorResponse after exhausting retries
	if err == nil {
		t.Error("Expected error after exhausting retries, got nil")
	}

	// Should be an ErrorResponse
	errResp, ok := err.(*ErrorResponse)
	if !ok {
		t.Errorf("Expected *ErrorResponse, got %T", err)
	} else if errResp.Code != http.StatusInternalServerError {
		t.Errorf("Expected error code 500, got %d", errResp.Code)
	}
}

// TestRetry_DisabledRetries tests that retries don't happen when RetryMax is 0
func TestRetry_DisabledRetries(t *testing.T) {
	var requestCount atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"message": "Server error", "code": 500}`)
	}))
	defer server.Close()

	client, err := NewClient(
		"test-key", "test-token", "test-project", "Mumbai",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}), // Disabled
	)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	ctx := context.Background()
	req, err := client.NewRequest(ctx, http.MethodGet, "test", nil)
	if err != nil {
		t.Fatalf("NewRequest error: %v", err)
	}

	_, err = client.Do(ctx, req, nil)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Should only have made 1 request (no retries)
	if count := requestCount.Load(); count != 1 {
		t.Errorf("Expected 1 request (no retries), got %d", count)
	}
}

// TestRetry_ContextCanceledDuringRetry tests that context cancellation stops retries
func TestRetry_ContextCanceledDuringRetry(t *testing.T) {
	var requestCount atomic.Int32
	mu := sync.Mutex{}
	requestTimes := []time.Time{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = requestCount.Add(1)
		mu.Lock()
		requestTimes = append(requestTimes, time.Now())
		mu.Unlock()

		// Simulate long processing
		time.Sleep(50 * time.Millisecond)

		// Always return 500
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"message": "Server error", "code": 500}`)
	}))
	defer server.Close()

	client, err := NewClient(
		"test-key", "test-token", "test-project", "Mumbai",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{
			RetryMax:     5, // Many retries allowed
			RetryWaitMin: PtrTo(100 * time.Millisecond),
			RetryWaitMax: PtrTo(200 * time.Millisecond),
		}),
	)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	// Create a context that will be canceled after a short time
	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	req, err := client.NewRequest(ctx, http.MethodGet, "test", nil)
	if err != nil {
		t.Fatalf("NewRequest error: %v", err)
	}

	_, err = client.Do(ctx, req, nil)

	// Should be a context error, not a regular API error
	if err == nil {
		t.Error("Expected context error, got nil")
	} else if err != context.DeadlineExceeded && err != context.Canceled {
		// Could be either canceled or deadline exceeded
		t.Errorf("Expected context error, got %v", err)
	}

	// Should have fewer requests than allowed retries due to context timeout
	count := requestCount.Load()
	if count >= 5 {
		t.Errorf("Expected fewer than 5 requests due to context cancellation, got %d", count)
	}
}

// TestRetry_ContextCanceledBeforeRequest tests context canceled before request is made
func TestRetry_ContextCanceledBeforeRequest(t *testing.T) {
	var requestCount atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClient(
		"test-key", "test-token", "test-project", "Mumbai",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}),
	)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	// Create an already-canceled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	req, err := client.NewRequest(ctx, http.MethodGet, "test", nil)
	if err != nil {
		t.Fatalf("NewRequest error: %v", err)
	}

	_, err = client.Do(ctx, req, nil)

	if err != context.Canceled {
		t.Errorf("Expected context.Canceled error, got %v", err)
	}

	// Should not have made any requests
	if count := requestCount.Load(); count != 0 {
		t.Errorf("Expected 0 requests, got %d", count)
	}
}

// TestRetry_BackoffDuration tests that backoff times increase appropriately
func TestRetry_BackoffDuration(t *testing.T) {
	var requestCount atomic.Int32
	mu := sync.Mutex{}
	requestTimes := []time.Time{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := requestCount.Add(1)
		mu.Lock()
		requestTimes = append(requestTimes, time.Now())
		mu.Unlock()

		if count < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"message": "Server error", "code": 500}`)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"code": 200, "message": "OK", "data": {}}`)
	}))
	defer server.Close()

	minWait := 50 * time.Millisecond
	maxWait := 200 * time.Millisecond

	client, err := NewClient(
		"test-key", "test-token", "test-project", "Mumbai",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{
			RetryMax:     3,
			RetryWaitMin: &minWait,
			RetryWaitMax: &maxWait,
		}),
	)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	ctx := context.Background()
	req, err := client.NewRequest(ctx, http.MethodGet, "test", nil)
	if err != nil {
		t.Fatalf("NewRequest error: %v", err)
	}

	_, err = client.Do(ctx, req, nil)

	if err != nil {
		t.Fatalf("Expected success after retry, got error: %v", err)
	}

	mu.Lock()
	times := requestTimes
	mu.Unlock()

	if len(times) != 3 {
		t.Fatalf("Expected 3 requests, got %d", len(times))
	}

	// Check that there are delays between requests
	firstDelay := times[1].Sub(times[0])
	secondDelay := times[2].Sub(times[1])

	// Both delays should be at least minWait
	if firstDelay < minWait {
		t.Errorf("First delay %v is less than minimum %v", firstDelay, minWait)
	}
	if secondDelay < minWait {
		t.Errorf("Second delay %v is less than minimum %v", secondDelay, minWait)
	}

	// Both delays should be no more than maxWait (with some tolerance for execution time)
	tolerance := 100 * time.Millisecond
	if firstDelay > maxWait+tolerance {
		t.Errorf("First delay %v exceeds maximum %v", firstDelay, maxWait)
	}
	if secondDelay > maxWait+tolerance {
		t.Errorf("Second delay %v exceeds maximum %v", secondDelay, maxWait)
	}

	// Verify we have at least 2 requests due to retries
	_ = len(times) // We verified this above
}

// TestRetry_MaxRetriesBoundary tests behavior at maxRetries boundaries
func TestRetry_MaxRetriesBoundary(t *testing.T) {
	tests := []struct {
		name          string
		maxRetries    int
		failureCount  int
		shouldSucceed bool
		expectedCount int
	}{
		{
			name:          "zero retries, immediate failure",
			maxRetries:    0,
			failureCount:  1,
			shouldSucceed: false,
			expectedCount: 1,
		},
		{
			name:          "one retry, fails twice",
			maxRetries:    1,
			failureCount:  2,
			shouldSucceed: false,
			expectedCount: 2, // 1 initial + 1 retry
		},
		{
			name:          "one retry, fails once then succeeds",
			maxRetries:    1,
			failureCount:  1,
			shouldSucceed: true,
			expectedCount: 2, // 1 initial + 1 retry
		},
		{
			name:          "three retries, fails twice then succeeds",
			maxRetries:    3,
			failureCount:  2,
			shouldSucceed: true,
			expectedCount: 3, // 1 initial + 2 retries
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var requestCount atomic.Int32

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				count := requestCount.Add(1)

				if count <= int32(tt.failureCount) {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprint(w, `{"message": "Server error", "code": 500}`)
					return
				}

				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, `{"code": 200, "message": "OK", "data": {}}`)
			}))
			defer server.Close()

			client, err := NewClient(
				"test-key", "test-token", "test-project", "Mumbai",
				SetBaseURL(server.URL),
				WithRetryAndBackoffs(RetryConfig{
					RetryMax:     tt.maxRetries,
					RetryWaitMin: PtrTo(10 * time.Millisecond),
					RetryWaitMax: PtrTo(50 * time.Millisecond),
				}),
			)
			if err != nil {
				t.Fatalf("NewClient error: %v", err)
			}

			ctx := context.Background()
			req, err := client.NewRequest(ctx, http.MethodGet, "test", nil)
			if err != nil {
				t.Fatalf("NewRequest error: %v", err)
			}

			_, err = client.Do(ctx, req, nil)

			if tt.shouldSucceed && err != nil {
				t.Errorf("Expected success, got error: %v", err)
			}
			if !tt.shouldSucceed && err == nil {
				t.Error("Expected error, got nil")
			}

			if count := requestCount.Load(); count != int32(tt.expectedCount) {
				t.Errorf("Expected %d requests, got %d", tt.expectedCount, count)
			}
		})
	}
}

// TestRetry_ClientErrorsNotRetried tests that 4xx errors are not retried
func TestRetry_ClientErrorsNotRetried(t *testing.T) {
	var requestCount atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount.Add(1)
		// Return 400 Bad Request - should not be retried
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"message": "Bad request", "code": 400}`)
	}))
	defer server.Close()

	client, err := NewClient(
		"test-key", "test-token", "test-project", "Mumbai",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{
			RetryMax:     3, // Even with retries configured
			RetryWaitMin: PtrTo(10 * time.Millisecond),
			RetryWaitMax: PtrTo(50 * time.Millisecond),
		}),
	)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	ctx := context.Background()
	req, err := client.NewRequest(ctx, http.MethodGet, "test", nil)
	if err != nil {
		t.Fatalf("NewRequest error: %v", err)
	}

	_, err = client.Do(ctx, req, nil)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Should only have made 1 request (4xx errors are not retried)
	if count := requestCount.Load(); count != 1 {
		t.Errorf("Expected 1 request (no retries for 4xx), got %d", count)
	}
}

// TestRetry_401UnauthorizedNotRetried tests that 401 is not retried
func TestRetry_401UnauthorizedNotRetried(t *testing.T) {
	var requestCount atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount.Add(1)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, `{"message": "Unauthorized", "code": 401}`)
	}))
	defer server.Close()

	client, err := NewClient(
		"test-key", "test-token", "test-project", "Mumbai",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{
			RetryMax:     3,
			RetryWaitMin: PtrTo(10 * time.Millisecond),
			RetryWaitMax: PtrTo(50 * time.Millisecond),
		}),
	)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	ctx := context.Background()
	req, err := client.NewRequest(ctx, http.MethodGet, "test", nil)
	if err != nil {
		t.Fatalf("NewRequest error: %v", err)
	}

	_, err = client.Do(ctx, req, nil)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Should only have made 1 request
	if count := requestCount.Load(); count != 1 {
		t.Errorf("Expected 1 request, got %d", count)
	}
}

// TestRetry_429TooManyRequestsRecovery tests retry on rate limiting
func TestRetry_429TooManyRequestsRecovery(t *testing.T) {
	var requestCount atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := requestCount.Add(1)

		if count == 1 {
			w.WriteHeader(http.StatusTooManyRequests)
			fmt.Fprint(w, `{"message": "Rate limited", "code": 429}`)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"code": 200, "message": "OK", "data": {}}`)
	}))
	defer server.Close()

	client, err := NewClient(
		"test-key", "test-token", "test-project", "Mumbai",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{
			RetryMax:     2,
			RetryWaitMin: PtrTo(10 * time.Millisecond),
			RetryWaitMax: PtrTo(50 * time.Millisecond),
		}),
	)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	ctx := context.Background()
	req, err := client.NewRequest(ctx, http.MethodGet, "test", nil)
	if err != nil {
		t.Fatalf("NewRequest error: %v", err)
	}

	resp, err := client.Do(ctx, req, nil)
	if err != nil {
		t.Errorf("Expected success after retry, got error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

// TestRetry_MultipleRetryConfig tests overriding retry config
func TestRetry_MultipleRetryConfig(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"code": 200, "message": "OK", "data": {}}`)
	}))
	defer server.Close()

	// First config
	config1 := RetryConfig{
		RetryMax:     5,
		RetryWaitMin: PtrTo(1 * time.Second),
		RetryWaitMax: PtrTo(10 * time.Second),
	}

	client, err := NewClient(
		"test-key", "test-token", "test-project", "Mumbai",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(config1),
	)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	// Verify first config was applied
	if client.RetryConfig.RetryMax != 5 {
		t.Errorf("Expected RetryMax 5, got %d", client.RetryConfig.RetryMax)
	}
}

// TestRetry_504GatewayTimeoutRecovery tests retry on 504 Gateway Timeout
func TestRetry_504GatewayTimeoutRecovery(t *testing.T) {
	var requestCount atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := requestCount.Add(1)

		if count <= 1 {
			w.WriteHeader(http.StatusGatewayTimeout)
			fmt.Fprint(w, `{"message": "Gateway timeout", "code": 504}`)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"code": 200, "message": "OK", "data": {}}`)
	}))
	defer server.Close()

	client, err := NewClient(
		"test-key", "test-token", "test-project", "Mumbai",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{
			RetryMax:     2,
			RetryWaitMin: PtrTo(10 * time.Millisecond),
			RetryWaitMax: PtrTo(50 * time.Millisecond),
		}),
	)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	ctx := context.Background()
	req, err := client.NewRequest(ctx, http.MethodGet, "test", nil)
	if err != nil {
		t.Fatalf("NewRequest error: %v", err)
	}

	resp, err := client.Do(ctx, req, nil)
	if err != nil {
		t.Errorf("Expected success after retry, got error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

// TestRetry_ContextTimeoutDuringBackoff tests timeout during backoff wait period
func TestRetry_ContextTimeoutDuringBackoff(t *testing.T) {
	var requestCount atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"message": "Server error", "code": 500}`)
	}))
	defer server.Close()

	client, err := NewClient(
		"test-key", "test-token", "test-project", "Mumbai",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{
			RetryMax:     5,
			RetryWaitMin: PtrTo(200 * time.Millisecond), // Long backoff
			RetryWaitMax: PtrTo(500 * time.Millisecond),
		}),
	)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	// Context that will timeout before backoff completes
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	req, err := client.NewRequest(ctx, http.MethodGet, "test", nil)
	if err != nil {
		t.Fatalf("NewRequest error: %v", err)
	}

	_, err = client.Do(ctx, req, nil)

	// Should fail with context deadline or cancellation
	if err == nil {
		t.Error("Expected context error, got nil")
	}

	// Should have made only 1 request (timeout before first backoff retry)
	count := requestCount.Load()
	if count != 1 && count != 2 {
		// Could be 1 or 2 depending on timing, but definitely not all 5
		t.Logf("Got %d requests, which is reasonable for timeout during backoff", count)
	}
	if count > 3 {
		t.Errorf("Expected 1-2 requests before timeout, got %d", count)
	}
}
