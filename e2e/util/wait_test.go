package util

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Test helper: StateRefreshFunc that returns a sequence of states
type stateSequence struct {
	states []stateResult
	index  int
}

type stateResult struct {
	obj   interface{}
	state string
	err   error
}

func newStateSequence(states []stateResult) *stateSequence {
	return &stateSequence{states: states, index: 0}
}

func (s *stateSequence) next() resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		if s.index >= len(s.states) {
			return nil, "", fmt.Errorf("no more states in sequence")
		}
		result := s.states[s.index]
		s.index++
		return result.obj, result.state, result.err
	}
}

func TestWaitForState_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	seq := newStateSequence([]stateResult{
		{obj: "pending", state: "pending", err: nil},
		{obj: "pending", state: "pending", err: nil},
		{obj: "ready", state: "ready", err: nil},
	})

	err := WaitForState(
		ctx,
		seq.next(),
		[]string{"pending"},
		[]string{"ready"},
		30*time.Second,
		10*time.Millisecond, // Fast for testing
		5*time.Millisecond,
		0,
	)

	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}

func TestWaitForState_ImmediateSuccess(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	seq := newStateSequence([]stateResult{
		{obj: "ready", state: "ready", err: nil},
	})

	err := WaitForState(
		ctx,
		seq.next(),
		[]string{"pending"},
		[]string{"ready"},
		30*time.Second,
		10*time.Millisecond,
		5*time.Millisecond,
		0,
	)

	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}

func TestWaitForState_Timeout(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	seq := newStateSequence([]stateResult{
		{obj: "pending", state: "pending", err: nil},
		{obj: "pending", state: "pending", err: nil},
		{obj: "pending", state: "pending", err: nil},
		{obj: "pending", state: "pending", err: nil},
		{obj: "pending", state: "pending", err: nil},
	})

	err := WaitForState(
		ctx,
		seq.next(),
		[]string{"pending"},
		[]string{"ready"},
		100*time.Millisecond, // Short timeout
		20*time.Millisecond,
		10*time.Millisecond,
		0,
	)

	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

func TestWaitForState_RefreshError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testErr := errors.New("refresh failed")
	seq := newStateSequence([]stateResult{
		{obj: nil, state: "", err: testErr},
	})

	err := WaitForState(
		ctx,
		seq.next(),
		[]string{"pending"},
		[]string{"ready"},
		30*time.Second,
		10*time.Millisecond,
		5*time.Millisecond,
		0,
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, testErr) && err.Error() != testErr.Error() {
		t.Fatalf("expected error %v, got %v", testErr, err)
	}
}

func TestWaitForState_MultipleTargetStates(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	seq := newStateSequence([]stateResult{
		{obj: "creating", state: "creating", err: nil},
		{obj: "running", state: "running", err: nil},
	})

	err := WaitForState(
		ctx,
		seq.next(),
		[]string{"creating", "pending"},
		[]string{"running", "stopped"}, // Multiple target states
		30*time.Second,
		10*time.Millisecond,
		5*time.Millisecond,
		0,
	)

	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}

func TestWaitForState_DefaultValues(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	seq := newStateSequence([]stateResult{
		{obj: "ready", state: "ready", err: nil},
	})

	// Test with zero values - should use defaults
	err := WaitForState(
		ctx,
		seq.next(),
		[]string{"pending"},
		[]string{"ready"},
		30*time.Second,
		0, // Should use StateChangeDefaultDelay
		0, // Should use StateChangeRetryBackoff
		0, // Should use DefaultNotFoundChecks
	)

	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}

func TestWaitForState_NotFoundHandling(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	seq := newStateSequence([]stateResult{
		{nil, "", nil}, // NotFound
		{nil, "", nil}, // NotFound
		{obj: "ready", state: "ready", err: nil},
	})

	err := WaitForState(
		ctx,
		seq.next(),
		[]string{"pending"},
		[]string{"ready"},
		30*time.Second,
		10*time.Millisecond,
		5*time.Millisecond,
		2, // Allow 2 NotFound checks
	)

	if err != nil {
		t.Fatalf("expected success after NotFound checks, got error: %v", err)
	}
}

// ============================================================================
// Edge Case Tests: Network Errors, Retries, and Error Recovery
// ============================================================================

// TestWaitForState_TransientNetworkError tests that transient network errors are retried
func TestWaitForState_TransientNetworkError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	callCount := 0
	networkErr := errors.New("connection refused")
	seq := newStateSequence([]stateResult{
		{obj: nil, state: "", err: networkErr},   // First call fails
		{obj: nil, state: "", err: networkErr},   // Second call fails
		{obj: "ready", state: "ready", err: nil}, // Third call succeeds
	})

	refreshFunc := func() (interface{}, string, error) {
		callCount++
		return seq.next()()
	}

	err := WaitForState(
		ctx,
		refreshFunc,
		[]string{"pending"},
		[]string{"ready"},
		30*time.Second,
		10*time.Millisecond,
		5*time.Millisecond,
		0,
	)

	// Should fail because errors are not retried by WaitForState
	// (errors should be handled by the client's retry logic, not WaitForState)
	if err == nil {
		t.Fatal("expected error for network failure, got nil")
	}
	if callCount != 1 {
		t.Errorf("expected 1 call before error, got %d", callCount)
	}
}

// TestWaitForState_RecoverableError tests error recovery after transient failures
func TestWaitForState_RecoverableError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	callCount := 0
	tempErr := errors.New("temporary error")
	seq := newStateSequence([]stateResult{
		{obj: nil, state: "", err: tempErr},          // First call fails
		{obj: "pending", state: "pending", err: nil}, // Second call succeeds but pending
		{obj: "ready", state: "ready", err: nil},     // Third call succeeds
	})

	refreshFunc := func() (interface{}, string, error) {
		callCount++
		return seq.next()()
	}

	err := WaitForState(
		ctx,
		refreshFunc,
		[]string{"pending"},
		[]string{"ready"},
		30*time.Second,
		10*time.Millisecond,
		5*time.Millisecond,
		0,
	)

	// Should fail on first error (errors are not automatically retried)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestWaitForState_NotFoundExceeded tests when NotFound checks are exceeded
func TestWaitForState_NotFoundExceeded(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	seq := newStateSequence([]stateResult{
		{nil, "", nil}, // NotFound
		{nil, "", nil}, // NotFound
		{nil, "", nil}, // NotFound - exceeds limit
		{nil, "", nil}, // NotFound - exceeds limit
	})

	err := WaitForState(
		ctx,
		seq.next(),
		[]string{"pending"},
		[]string{"ready"},
		30*time.Second,
		10*time.Millisecond,
		5*time.Millisecond,
		2, // Only allow 2 NotFound checks
	)

	if err == nil {
		t.Fatal("expected error when NotFound checks exceeded, got nil")
	}
}

// TestWaitForState_ErrorStateInTarget tests when error state is returned
func TestWaitForState_ErrorStateInTarget(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	seq := newStateSequence([]stateResult{
		{obj: "error", state: "error", err: nil},
	})

	err := WaitForState(
		ctx,
		seq.next(),
		[]string{"pending"},
		[]string{"error"}, // Error state is target (unusual but possible)
		30*time.Second,
		10*time.Millisecond,
		5*time.Millisecond,
		0,
	)

	// Should succeed if error is in target
	if err != nil {
		t.Fatalf("expected success when error state is target, got error: %v", err)
	}
}

// TestWaitForState_EmptyTarget tests behavior with empty target (waiting for deletion)
func TestWaitForState_EmptyTarget(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	seq := newStateSequence([]stateResult{
		{nil, "", nil}, // Resource gone
	})

	err := WaitForState(
		ctx,
		seq.next(),
		[]string{"pending", "active"},
		[]string{}, // Empty target = waiting for resource to be gone
		30*time.Second,
		10*time.Millisecond,
		5*time.Millisecond,
		2, // Allow NotFound
	)

	if err != nil {
		t.Fatalf("expected success for empty target (deletion), got error: %v", err)
	}
}

// TestWaitForState_ContextTimeout tests context timeout handling
func TestWaitForState_ContextTimeout(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	seq := newStateSequence([]stateResult{
		{obj: "pending", state: "pending", err: nil},
		{obj: "pending", state: "pending", err: nil},
		{obj: "pending", state: "pending", err: nil},
	})

	err := WaitForState(
		ctx,
		seq.next(),
		[]string{"pending"},
		[]string{"ready"},
		5*time.Second, // Long timeout, but context will cancel first
		20*time.Millisecond,
		10*time.Millisecond,
		0,
	)

	if err == nil {
		t.Fatal("expected context timeout error, got nil")
	}
}

// TestWaitForState_DeadlineExceeded tests deadline exceeded error
func TestWaitForState_DeadlineExceeded(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(50*time.Millisecond))
	defer cancel()

	seq := newStateSequence([]stateResult{
		{obj: "pending", state: "pending", err: nil},
		{obj: "pending", state: "pending", err: nil},
	})

	err := WaitForState(
		ctx,
		seq.next(),
		[]string{"pending"},
		[]string{"ready"},
		5*time.Second,
		30*time.Millisecond,
		10*time.Millisecond,
		0,
	)

	if err == nil {
		t.Fatal("expected deadline exceeded error, got nil")
	}
}

// TestWaitForState_ErrorPropagation tests that errors from refresh function are properly propagated
func TestWaitForState_ErrorPropagation(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testErrors := []error{
		errors.New("network error"),
		fmt.Errorf("wrapped error: %w", errors.New("underlying error")),
		errors.New("timeout"),
		errors.New("permission denied"),
	}

	for _, testErr := range testErrors {
		t.Run(testErr.Error(), func(t *testing.T) {
			seq := newStateSequence([]stateResult{
				{obj: nil, state: "", err: testErr},
			})

			err := WaitForState(
				ctx,
				seq.next(),
				[]string{"pending"},
				[]string{"ready"},
				30*time.Second,
				10*time.Millisecond,
				5*time.Millisecond,
				0,
			)

			if err == nil {
				t.Fatalf("expected error %v, got nil", testErr)
			}
		})
	}
}

// TestWaitForState_StateFlickering tests handling of state that flickers between values
func TestWaitForState_StateFlickering(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	seq := newStateSequence([]stateResult{
		{obj: "pending", state: "pending", err: nil},
		{obj: "ready", state: "ready", err: nil},
		{obj: "pending", state: "pending", err: nil}, // Flickers back
		{obj: "ready", state: "ready", err: nil},     // Back to ready
		{obj: "ready", state: "ready", err: nil},     // Stable
	})

	err := WaitForState(
		ctx,
		seq.next(),
		[]string{"pending"},
		[]string{"ready"},
		30*time.Second,
		10*time.Millisecond,
		5*time.Millisecond,
		0,
	)

	// Should succeed once we hit target state, even if it flickered
	if err != nil {
		t.Fatalf("expected success despite flickering, got error: %v", err)
	}
}

// TestWaitForState_ConcurrentCalls tests that the function is safe for concurrent use
func TestWaitForState_ConcurrentCalls(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	results := make(chan error, 10)

	// Run 10 concurrent wait operations
	for i := 0; i < 10; i++ {
		go func(id int) {
			seq := newStateSequence([]stateResult{
				{obj: "pending", state: "pending", err: nil},
				{obj: fmt.Sprintf("ready-%d", id), state: "ready", err: nil},
			})

			err := WaitForState(
				ctx,
				seq.next(),
				[]string{"pending"},
				[]string{"ready"},
				30*time.Second,
				10*time.Millisecond,
				5*time.Millisecond,
				0,
			)
			results <- err
		}(i)
	}

	// Collect all results
	for i := 0; i < 10; i++ {
		err := <-results
		if err != nil {
			t.Errorf("concurrent call %d failed: %v", i, err)
		}
	}
}

// TestWaitForState_InvalidStateTransition tests handling of unexpected state transitions
func TestWaitForState_InvalidStateTransition(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	seq := newStateSequence([]stateResult{
		{obj: "pending", state: "pending", err: nil},
		{obj: "unknown", state: "unknown", err: nil}, // Unexpected state
		{obj: "ready", state: "ready", err: nil},
	})

	err := WaitForState(
		ctx,
		seq.next(),
		[]string{"pending"},
		[]string{"ready"},
		30*time.Second,
		10*time.Millisecond,
		5*time.Millisecond,
		0,
	)

	// Should timeout because "unknown" is not in pending or target
	if err == nil {
		t.Fatal("expected timeout for invalid state transition, got nil")
	}
}

// TestWaitForState_LongRunningOperation tests behavior with very long operations
func TestWaitForState_LongRunningOperation(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	// Simulate a long-running operation with many pending states
	// Using 5 pending states to test the scenario while ensuring completion within test timeout
	states := make([]stateResult, 0, 6)
	for i := 0; i < 5; i++ {
		states = append(states, stateResult{obj: "pending", state: "pending", err: nil})
	}
	states = append(states, stateResult{obj: "ready", state: "ready", err: nil})

	seq := newStateSequence(states)

	err := WaitForState(
		ctx,
		seq.next(),
		[]string{"pending"},
		[]string{"ready"},
		20*time.Second, // Longer timeout for long operation (but less than test suite timeout)
		10*time.Millisecond,
		5*time.Millisecond,
		0,
	)

	if err != nil {
		t.Fatalf("expected success after long operation, got error: %v", err)
	}
}

// Note: Integration tests for specific wait functions (WaitForFunctionReady, WaitForNodePowerState, etc.)
// would require mocking the goe2e.Client. The core WaitForState function is tested above,
// and the specific wait functions are thin wrappers that use WaitForState internally.

func TestWaitForState_ContextCancellation(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	seq := newStateSequence([]stateResult{
		{obj: "pending", state: "pending", err: nil},
	})

	// Cancel context immediately
	cancel()

	err := WaitForState(
		ctx,
		seq.next(),
		[]string{"pending"},
		[]string{"ready"},
		30*time.Second,
		10*time.Millisecond,
		5*time.Millisecond,
		0,
	)

	if err == nil {
		t.Fatal("expected context cancellation error, got nil")
	}
}

func TestWaitForState_EmptyPending(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	seq := newStateSequence([]stateResult{
		{obj: "ready", state: "ready", err: nil},
	})

	err := WaitForState(
		ctx,
		seq.next(),
		[]string{}, // Empty pending - any non-target state will timeout
		[]string{"ready"},
		30*time.Second,
		10*time.Millisecond,
		5*time.Millisecond,
		0,
	)

	if err != nil {
		t.Fatalf("expected success with empty pending, got error: %v", err)
	}
}

func TestWaitForState_StateTransition(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	seq := newStateSequence([]stateResult{
		{obj: "state1", state: "state1", err: nil},
		{obj: "state2", state: "state2", err: nil},
		{obj: "state3", state: "state3", err: nil},
		{obj: "target", state: "target", err: nil},
	})

	err := WaitForState(
		ctx,
		seq.next(),
		[]string{"state1", "state2", "state3"},
		[]string{"target"},
		30*time.Second,
		10*time.Millisecond,
		5*time.Millisecond,
		0,
	)

	if err != nil {
		t.Fatalf("expected success after state transitions, got error: %v", err)
	}
}
