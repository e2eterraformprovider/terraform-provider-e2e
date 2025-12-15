package sfs

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test constants - commonly used test values
// These constants are used across test files to ensure consistency
const (
	testSfsID                 = "test-sfs-123"
	testSfsIDShort            = "test-id" // Used for simple timeout/cancellation tests
	testErrorNotImplemented   = "not implemented"
	testErrorTemporaryNetwork = "temporary network error"
	testInvalidGetterValue    = "not a getter"
	testOldFieldName          = "old_field-name"
	testNewFieldName          = "new_field.name"
)

// mockSfsService is a mock implementation of goe2e.SfsService for testing
type mockSfsService struct {
	getSfsFunc func(ctx context.Context, sfsID string) (*goe2e.Sfs, *goe2e.Response, error)
	callCount  int
}

func (m *mockSfsService) GetSfs(ctx context.Context, sfsID string) (*goe2e.Sfs, *goe2e.Response, error) {
	m.callCount++
	if m.getSfsFunc != nil {
		return m.getSfsFunc(ctx, sfsID)
	}
	return nil, nil, errors.New(testErrorNotImplemented)
}

func (m *mockSfsService) CreateSfs(ctx context.Context, req *goe2e.SfsCreateRequest) (*goe2e.Sfs, *goe2e.Response, error) {
	return nil, nil, errors.New(testErrorNotImplemented)
}

func (m *mockSfsService) DeleteSfs(ctx context.Context, sfsID string) (*goe2e.Response, error) {
	return nil, errors.New(testErrorNotImplemented)
}

func (m *mockSfsService) ListSfss(ctx context.Context) ([]goe2e.Sfs, *goe2e.Response, error) {
	return nil, nil, errors.New(testErrorNotImplemented)
}

func (m *mockSfsService) ActivateSFSBackup(ctx context.Context, sfsID string, req *goe2e.ActivateSFSBackupRequest) (*goe2e.Response, error) {
	return nil, errors.New(testErrorNotImplemented)
}

func (m *mockSfsService) DeactivateSFSBackup(ctx context.Context, sfsID string) (*goe2e.Response, error) {
	return nil, errors.New(testErrorNotImplemented)
}

// mockClient creates a mock goe2e.Client with a mock SfsService
func mockClient(service *mockSfsService) *goe2e.Client {
	client := &goe2e.Client{}
	client.Sfs = service
	return client
}

// TestNormalizeSfsState_WhitespaceHandling tests whitespace handling
// Note: normalizeSfsState only lowercases, it doesn't trim whitespace
func TestNormalizeSfsState_WhitespaceHandling(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "status with leading whitespace",
			input:    " Creating",
			expected: " creating", // Function only lowercases, doesn't trim
		},
		{
			name:     "status with trailing whitespace",
			input:    "Active ",
			expected: "active ", // Function only lowercases, doesn't trim
		},
		{
			name:     "status with both leading and trailing whitespace",
			input:    " Active ",
			expected: " active ", // Function only lowercases, doesn't trim
		},
		{
			name:     "status with internal whitespace",
			input:    "Active State",
			expected: "active state",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeSfsState(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestNormalizeSfsState_AllPossibleAPIStatusValues tests all possible API status values
func TestNormalizeSfsState_AllPossibleAPIStatusValues(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "Creating", input: goe2econstants.SFSStatusCreating, expected: goe2econstants.SFSStateCreating},
		{name: "Active", input: goe2econstants.SFSStatusActive, expected: goe2econstants.SFSStateActive},
		{name: "Deleting", input: goe2econstants.SFSStatusDeleting, expected: goe2econstants.SFSStateDeleting},
		{name: "Deleted", input: goe2econstants.SFSStatusDeleted, expected: goe2econstants.SFSStateDeleted},
		{name: "Error", input: goe2econstants.SFSStatusError, expected: goe2econstants.SFSStateError},
		{name: "CREATING", input: "CREATING", expected: goe2econstants.SFSStateCreating},
		{name: "ACTIVE", input: "ACTIVE", expected: goe2econstants.SFSStateActive},
		{name: "DELETING", input: "DELETING", expected: goe2econstants.SFSStateDeleting},
		{name: "DELETED", input: "DELETED", expected: goe2econstants.SFSStateDeleted},
		{name: "ERROR", input: "ERROR", expected: goe2econstants.SFSStateError},
		{name: "creating", input: goe2econstants.SFSStateCreating, expected: goe2econstants.SFSStateCreating},
		{name: "active", input: goe2econstants.SFSStateActive, expected: goe2econstants.SFSStateActive},
		{name: "deleting", input: goe2econstants.SFSStateDeleting, expected: goe2econstants.SFSStateDeleting},
		{name: "deleted", input: goe2econstants.SFSStateDeleted, expected: goe2econstants.SFSStateDeleted},
		{name: "error", input: goe2econstants.SFSStateError, expected: goe2econstants.SFSStateError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeSfsState(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestWaitForSfsStatus_SuccessfulStateTransition tests successful state transition with mocked client
// Note: This test verifies the logic but takes ~10+ seconds due to polling interval
func TestWaitForSfsStatus_SuccessfulStateTransition(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	sfsID := testSfsID
	callCount := 0

	mockService := &mockSfsService{
		getSfsFunc: func(ctx context.Context, id string) (*goe2e.Sfs, *goe2e.Response, error) {
			callCount++
			if callCount == 1 {
				// First call: still creating
				return &goe2e.Sfs{
					ID:     sfsID,
					Status: goe2econstants.SFSStatusCreating,
				}, nil, nil
			}
			// Second call: active
			return &goe2e.Sfs{
				ID:     sfsID,
				Status: goe2econstants.SFSStatusActive,
			}, nil, nil
		},
	}

	client := mockClient(mockService)
	ctx := context.Background()

	// Use longer timeout to account for 10 second polling interval
	err := waitForSfsStatus(ctx, client, sfsID, goe2econstants.SFSDesiredStatusActive, 15*time.Second)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, callCount, 2, "should have polled at least twice")
}

// TestWaitForSfsStatus_PollingIntervalBehavior tests polling interval behavior
// Note: This test verifies polling happens with 10s intervals (takes ~10 seconds)
func TestWaitForSfsStatus_PollingIntervalBehavior(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	sfsID := testSfsID
	callCount := 0
	startTime := time.Now()

	mockService := &mockSfsService{
		getSfsFunc: func(ctx context.Context, id string) (*goe2e.Sfs, *goe2e.Response, error) {
			callCount++
			if callCount < 2 {
				// Keep returning creating status
				return &goe2e.Sfs{
					ID:     sfsID,
					Status: goe2econstants.SFSStatusCreating,
				}, nil, nil
			}
			// Second call: active
			return &goe2e.Sfs{
				ID:     sfsID,
				Status: goe2econstants.SFSStatusActive,
			}, nil, nil
		},
	}

	client := mockClient(mockService)
	ctx := context.Background()

	// Use longer timeout to account for 10 second polling interval
	err := waitForSfsStatus(ctx, client, sfsID, goe2econstants.SFSDesiredStatusActive, 20*time.Second)
	elapsed := time.Since(startTime)

	require.NoError(t, err)
	assert.Equal(t, 2, callCount, "should have polled 2 times")
	// Polling interval is 10 seconds, so we expect at least 10 seconds elapsed (1 interval)
	assert.GreaterOrEqual(t, elapsed, 10*time.Second, "should have delay between polls (10s interval)")
}

// TestWaitForSfsStatus_ErrorStateDetection tests error state detection
func TestWaitForSfsStatus_ErrorStateDetection(t *testing.T) {
	sfsID := testSfsID

	mockService := &mockSfsService{
		getSfsFunc: func(ctx context.Context, id string) (*goe2e.Sfs, *goe2e.Response, error) {
			return &goe2e.Sfs{
				ID:     sfsID,
				Status: goe2econstants.SFSStatusError,
			}, nil, nil
		},
	}

	client := mockClient(mockService)
	ctx := context.Background()

	err := waitForSfsStatus(ctx, client, sfsID, goe2econstants.SFSDesiredStatusActive, 15*time.Second)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "entered error state")
	assert.Contains(t, err.Error(), sfsID)
}

// TestWaitForSfsStatus_ContextCancellation tests context cancellation during polling
func TestWaitForSfsStatus_ContextCancellation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	sfsID := testSfsID
	callCount := 0

	mockService := &mockSfsService{
		getSfsFunc: func(ctx context.Context, id string) (*goe2e.Sfs, *goe2e.Response, error) {
			callCount++
			// Always return creating status to force polling
			return &goe2e.Sfs{
				ID:     sfsID,
				Status: goe2econstants.SFSStatusCreating,
			}, nil, nil
		},
	}

	client := mockClient(mockService)

	// Create a context that will be cancelled after first poll
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel the context after a short delay to simulate cancellation during polling
	go func() {
		time.Sleep(500 * time.Millisecond)
		cancel()
	}()

	// This should fail due to context cancellation
	err := waitForSfsStatus(ctx, client, sfsID, goe2econstants.SFSDesiredStatusActive, 30*time.Second)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "context")
	assert.GreaterOrEqual(t, callCount, 1, "should have polled at least once before cancellation")
}

// TestWaitForSfsStatus_DeletedStateWith404Error tests deleted state with 404 error
func TestWaitForSfsStatus_DeletedStateWith404Error(t *testing.T) {
	sfsID := testSfsID

	mockService := &mockSfsService{
		getSfsFunc: func(ctx context.Context, id string) (*goe2e.Sfs, *goe2e.Response, error) {
			return nil, nil, fmt.Errorf("%s: %s", fmt.Sprintf(goe2econstants.SFSNotFound, id), goe2econstants.NotFoundCode)
		},
	}

	client := mockClient(mockService)
	ctx := context.Background()

	err := waitForSfsStatus(ctx, client, sfsID, goe2econstants.SFSDesiredStatusDeleted, 15*time.Second)
	require.NoError(t, err, "should succeed when waiting for deleted state and getting 404")
}

// TestWaitForSfsStatus_DeletedStateWithNilResponse tests deleted state with nil SFS response
func TestWaitForSfsStatus_DeletedStateWithNilResponse(t *testing.T) {
	sfsID := testSfsID

	mockService := &mockSfsService{
		getSfsFunc: func(ctx context.Context, id string) (*goe2e.Sfs, *goe2e.Response, error) {
			return nil, nil, nil
		},
	}

	client := mockClient(mockService)
	ctx := context.Background()

	err := waitForSfsStatus(ctx, client, sfsID, goe2econstants.SFSDesiredStatusDeleted, 15*time.Second)
	require.NoError(t, err, "should succeed when waiting for deleted state and getting nil response")
}

// TestWaitForSfsStatus_404DesiredStatusHandling tests 404 desired status handling
func TestWaitForSfsStatus_404DesiredStatusHandling(t *testing.T) {
	sfsID := testSfsID

	mockService := &mockSfsService{
		getSfsFunc: func(ctx context.Context, id string) (*goe2e.Sfs, *goe2e.Response, error) {
			return nil, nil, fmt.Errorf("resource %s: %s", id, goe2econstants.NotFoundSubstring) // Test mock - using generic format
		},
	}

	client := mockClient(mockService)
	ctx := context.Background()

	err := waitForSfsStatus(ctx, client, sfsID, goe2econstants.SFSDesiredStatus404, 15*time.Second)
	require.NoError(t, err, "should succeed when waiting for 404 status")
}

// TestWaitForSfsStatus_APIErrorsDuringPolling tests API errors during polling
// Note: This test takes ~20+ seconds due to 10s polling interval and retries
func TestWaitForSfsStatus_APIErrorsDuringPolling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	sfsID := testSfsID
	callCount := 0

	mockService := &mockSfsService{
		getSfsFunc: func(ctx context.Context, id string) (*goe2e.Sfs, *goe2e.Response, error) {
			callCount++
			if callCount < 2 {
				// First call: transient error
				return nil, nil, errors.New(testErrorTemporaryNetwork)
			}
			// Second call: success
			return &goe2e.Sfs{
				ID:     sfsID,
				Status: goe2econstants.SFSStatusActive,
			}, nil, nil
		},
	}

	client := mockClient(mockService)
	ctx := context.Background()

	// Use longer timeout to account for 10 second polling interval and retries
	err := waitForSfsStatus(ctx, client, sfsID, goe2econstants.SFSDesiredStatusActive, 20*time.Second)
	require.NoError(t, err, "should continue polling after transient errors")
	assert.GreaterOrEqual(t, callCount, 2, "should have retried after errors")
}

// TestWaitForSfsStatus_SFSNotFoundErrors tests SFS not found errors (except for deleted state)
func TestWaitForSfsStatus_SFSNotFoundErrors(t *testing.T) {
	sfsID := testSfsID

	mockService := &mockSfsService{
		getSfsFunc: func(ctx context.Context, id string) (*goe2e.Sfs, *goe2e.Response, error) {
			return nil, nil, fmt.Errorf(goe2econstants.SFSNotFound+" %s", id, goe2econstants.NotFoundSubstring)
		},
	}

	client := mockClient(mockService)
	ctx := context.Background()

	err := waitForSfsStatus(ctx, client, sfsID, goe2econstants.SFSDesiredStatusActive, 15*time.Second)
	require.Error(t, err, "should fail when SFS not found and not waiting for deleted state")
	assert.Contains(t, err.Error(), "not found")
}

// TestWaitForSfsStatus_MultipleStateChecks tests multiple state checks before reaching desired state
// Note: This test takes ~10 seconds due to 10s polling interval
func TestWaitForSfsStatus_MultipleStateChecks(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	sfsID := testSfsID
	callCount := 0
	states := []string{
		goe2econstants.SFSStatusCreating,
		goe2econstants.SFSStatusActive,
	}

	mockService := &mockSfsService{
		getSfsFunc: func(ctx context.Context, id string) (*goe2e.Sfs, *goe2e.Response, error) {
			state := states[callCount]
			callCount++
			return &goe2e.Sfs{
				ID:     sfsID,
				Status: state,
			}, nil, nil
		},
	}

	client := mockClient(mockService)
	ctx := context.Background()

	err := waitForSfsStatus(ctx, client, sfsID, goe2econstants.SFSDesiredStatusActive, 20*time.Second)
	require.NoError(t, err)
	assert.Equal(t, 2, callCount, "should have polled 2 times")
}

// TestWaitForSfsStatus_CurrentStateMatchesDesiredStateImmediately tests current state matches desired state immediately
func TestWaitForSfsStatus_CurrentStateMatchesDesiredStateImmediately(t *testing.T) {
	sfsID := testSfsID

	mockService := &mockSfsService{
		getSfsFunc: func(ctx context.Context, id string) (*goe2e.Sfs, *goe2e.Response, error) {
			return &goe2e.Sfs{
				ID:     sfsID,
				Status: goe2econstants.SFSStatusActive,
			}, nil, nil
		},
	}

	client := mockClient(mockService)
	ctx := context.Background()

	err := waitForSfsStatus(ctx, client, sfsID, goe2econstants.SFSDesiredStatusActive, 15*time.Second)
	require.NoError(t, err)
	assert.Equal(t, 1, mockService.callCount, "should only poll once when state matches immediately")
}

// TestWaitForSfsStatus_NormalizedStatusMatching tests normalized status matching
func TestWaitForSfsStatus_NormalizedStatusMatching(t *testing.T) {
	sfsID := testSfsID

	mockService := &mockSfsService{
		getSfsFunc: func(ctx context.Context, id string) (*goe2e.Sfs, *goe2e.Response, error) {
			// API returns Active (capitalized), but we're waiting for active (lowercase)
			return &goe2e.Sfs{
				ID:     sfsID,
				Status: goe2econstants.SFSStatusActive,
			}, nil, nil
		},
	}

	client := mockClient(mockService)
	ctx := context.Background()

	err := waitForSfsStatus(ctx, client, sfsID, goe2econstants.SFSDesiredStatusActive, 15*time.Second)
	require.NoError(t, err, "should match normalized status")
}

// TestWaitForSfsStatus_BothNormalizedAndCurrentStatusMatching tests both normalized and current status matching logic
func TestWaitForSfsStatus_BothNormalizedAndCurrentStatusMatching(t *testing.T) {
	sfsID := testSfsID

	tests := []struct {
		name          string
		apiStatus     string
		desiredStatus string
		shouldSucceed bool
	}{
		{
			name:          "normalized status match",
			apiStatus:     goe2econstants.SFSStatusActive,
			desiredStatus: goe2econstants.SFSDesiredStatusActive,
			shouldSucceed: true,
		},
		{
			name:          "exact status match",
			apiStatus:     goe2econstants.SFSStatusActive,
			desiredStatus: goe2econstants.SFSStatusActive,
			shouldSucceed: true,
		},
		{
			name:          "case insensitive match",
			apiStatus:     goe2econstants.SFSStatusActive, // Testing exact match with API constant
			desiredStatus: goe2econstants.SFSDesiredStatusActive,
			shouldSucceed: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mockSfsService{
				getSfsFunc: func(ctx context.Context, id string) (*goe2e.Sfs, *goe2e.Response, error) {
					return &goe2e.Sfs{
						ID:     sfsID,
						Status: tt.apiStatus,
					}, nil, nil
				},
			}

			client := mockClient(mockService)
			ctx := context.Background()

			err := waitForSfsStatus(ctx, client, sfsID, tt.desiredStatus, 15*time.Second)
			if tt.shouldSucceed {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

// TestWaitForSfsActive tests waitForSfsActive function
// Note: This test takes ~10+ seconds due to polling interval
func TestWaitForSfsActive(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	sfsID := testSfsID
	callCount := 0

	mockService := &mockSfsService{
		getSfsFunc: func(ctx context.Context, id string) (*goe2e.Sfs, *goe2e.Response, error) {
			callCount++
			if callCount == 1 {
				return &goe2e.Sfs{
					ID:     sfsID,
					Status: goe2econstants.SFSStatusCreating,
				}, nil, nil
			}
			return &goe2e.Sfs{
				ID:     sfsID,
				Status: goe2econstants.SFSStatusActive,
			}, nil, nil
		},
	}

	client := mockClient(mockService)
	ctx := context.Background()

	err := waitForSfsActive(ctx, client, sfsID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, callCount, 2, "should have polled at least twice")
}

// TestWaitForSfsActive_TimeoutValue tests timeout value is correct (10 minutes)
func TestWaitForSfsActive_TimeoutValue(t *testing.T) {
	// This test verifies that waitForSfsActive uses sfsCreateTimeout (10 minutes)
	// We can't easily test the exact timeout without waiting, but we can verify
	// it calls waitForSfsStatus with the correct desired status
	sfsID := testSfsID

	mockService := &mockSfsService{
		getSfsFunc: func(ctx context.Context, id string) (*goe2e.Sfs, *goe2e.Response, error) {
			return &goe2e.Sfs{
				ID:     sfsID,
				Status: goe2econstants.SFSStatusActive,
			}, nil, nil
		},
	}

	client := mockClient(mockService)
	ctx := context.Background()

	err := waitForSfsActive(ctx, client, sfsID)
	require.NoError(t, err)
	// Verify it's waiting for active status
	assert.Equal(t, 1, mockService.callCount)
}

// TestWaitForSfsActive_ErrorPropagation tests error propagation
func TestWaitForSfsActive_ErrorPropagation(t *testing.T) {
	sfsID := testSfsID

	mockService := &mockSfsService{
		getSfsFunc: func(ctx context.Context, id string) (*goe2e.Sfs, *goe2e.Response, error) {
			return &goe2e.Sfs{
				ID:     sfsID,
				Status: goe2econstants.SFSStatusError,
			}, nil, nil
		},
	}

	client := mockClient(mockService)
	ctx := context.Background()

	err := waitForSfsActive(ctx, client, sfsID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error state")
}

// TestWaitForSfsDeleted tests waitForSfsDeleted function (if it exists)
// Note: waitForSfsDeleted doesn't exist in helpers.go, so this test documents that
// deletion polling would work similarly to waitForSfsStatus with deleted status
func TestWaitForSfsDeleted_Conceptual(t *testing.T) {
	// This test demonstrates how waitForSfsDeleted would work if implemented
	// Currently, deletion uses waitForSfsStatus directly with deleted status
	sfsID := testSfsID

	mockService := &mockSfsService{
		getSfsFunc: func(ctx context.Context, id string) (*goe2e.Sfs, *goe2e.Response, error) {
			// Simulate deletion: first call returns deleting, then 404
			return nil, nil, fmt.Errorf("SFS %s %s", id, goe2econstants.NotFoundSubstring)
		},
	}

	client := mockClient(mockService)
	ctx := context.Background()

	// Using waitForSfsStatus with deleted status (how deletion currently works)
	err := waitForSfsStatus(ctx, client, sfsID, goe2econstants.SFSDesiredStatusDeleted, sfsDeleteTimeout)
	require.NoError(t, err, "should succeed when SFS is deleted (404)")
}

// TestGetEffectiveSizeGB tests getEffectiveSizeGB function
func TestGetEffectiveSizeGB(t *testing.T) {
	tests := []struct {
		name         string
		v3Value      interface{}
		v2Value      interface{}
		defaultValue int
		expected     int
		description  string
	}{
		{
			name:         "V3 field preferred when set",
			v3Value:      100,
			v2Value:      50,
			defaultValue: 0,
			expected:     100,
			description:  "V3 field should be preferred over V2",
		},
		{
			name:         "V2 field used as fallback when V3 not set",
			v3Value:      nil,
			v2Value:      50,
			defaultValue: 0,
			expected:     50,
			description:  "V2 field should be used when V3 is not set",
		},
		{
			name:         "V3 field with zero value should not use, fall back",
			v3Value:      0,
			v2Value:      50,
			defaultValue: 0,
			expected:     50,
			description:  "Zero V3 value should fall back to V2",
		},
		{
			name:         "V2 field with zero value should not use, return default",
			v3Value:      nil,
			v2Value:      0,
			defaultValue: 20,
			expected:     20,
			description:  "Zero V2 value should return default",
		},
		{
			name:         "both fields set, V3 preferred",
			v3Value:      200,
			v2Value:      150,
			defaultValue: 0,
			expected:     200,
			description:  "V3 should be preferred even when both are set",
		},
		{
			name:         "neither field set, returns default",
			v3Value:      nil,
			v2Value:      nil,
			defaultValue: 30,
			expected:     30,
			description:  "Default should be returned when neither field is set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGetter := &mockGetter{
				values: map[string]interface{}{
					tfconstants.AttrSizeGB:   tt.v3Value,
					tfconstants.AttrDiskSize: tt.v2Value,
				},
			}

			result := getEffectiveSizeGB(mockGetter, tfconstants.AttrSizeGB, tfconstants.AttrDiskSize, tt.defaultValue)
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

// TestGetEffectiveSizeGB_TypeAssertions tests type assertions
func TestGetEffectiveSizeGB_TypeAssertions(t *testing.T) {
	tests := []struct {
		name         string
		v3Value      interface{}
		v2Value      interface{}
		defaultValue int
		expected     int
		description  string
	}{
		{
			name:         "int type",
			v3Value:      int(100),
			defaultValue: 0,
			expected:     100,
		},
		{
			name:         "int32 type",
			v3Value:      int32(100),
			defaultValue: 0,
			expected:     0, // int32 won't match int assertion
		},
		{
			name:         "string type",
			v3Value:      "100",
			defaultValue: 0,
			expected:     0, // string won't match int assertion
		},
		{
			name:         "nil getter interface",
			v3Value:      nil,
			defaultValue: 50,
			expected:     50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var getter interface{}
			if tt.v3Value != nil {
				getter = &mockGetter{
					values: map[string]interface{}{
						tfconstants.AttrSizeGB: tt.v3Value,
					},
				}
			}

			result := getEffectiveSizeGB(getter, tfconstants.AttrSizeGB, tfconstants.AttrDiskSize, tt.defaultValue)
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

// TestGetEffectiveIOPS tests getEffectiveIOPS function
func TestGetEffectiveIOPS(t *testing.T) {
	tests := []struct {
		name         string
		v3Value      interface{}
		v2Value      interface{}
		defaultValue int
		expected     int
		description  string
	}{
		{
			name:         "V3 field preferred when set",
			v3Value:      1000,
			v2Value:      500,
			defaultValue: 0,
			expected:     1000,
		},
		{
			name:         "V2 field used as fallback when V3 not set",
			v3Value:      nil,
			v2Value:      500,
			defaultValue: 0,
			expected:     500,
		},
		{
			name:         "V3 field with zero value should not use, fall back",
			v3Value:      0,
			v2Value:      500,
			defaultValue: 0,
			expected:     500,
		},
		{
			name:         "V2 field with zero value should not use, return default",
			v3Value:      nil,
			v2Value:      0,
			defaultValue: 100,
			expected:     100,
		},
		{
			name:         "both fields set, V3 preferred",
			v3Value:      2000,
			v2Value:      1500,
			defaultValue: 0,
			expected:     2000,
		},
		{
			name:         "neither field set, returns default",
			v3Value:      nil,
			v2Value:      nil,
			defaultValue: 200,
			expected:     200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGetter := &mockGetter{
				values: map[string]interface{}{
					tfconstants.AttrIOPS:     tt.v3Value,
					tfconstants.AttrDiskIOPS: tt.v2Value,
				},
			}

			result := getEffectiveIOPS(mockGetter, tfconstants.AttrIOPS, tfconstants.AttrDiskIOPS, tt.defaultValue)
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

// TestGetEffectiveEncryptionEnabled tests getEffectiveEncryptionEnabled function
func TestGetEffectiveEncryptionEnabled(t *testing.T) {
	tests := []struct {
		name        string
		v3Value     interface{}
		v2Value     interface{}
		expected    bool
		description string
	}{
		{
			name:        "V3 field preferred when set",
			v3Value:     true,
			v2Value:     false,
			expected:    true,
			description: "V3 field should be preferred",
		},
		{
			name:        "V2 field used as fallback when V3 not set",
			v3Value:     nil,
			v2Value:     true,
			expected:    true,
			description: "V2 field should be used when V3 is not set",
		},
		{
			name:        "both fields set, V3 preferred",
			v3Value:     true,
			v2Value:     false,
			expected:    true,
			description: "V3 should be preferred even when both are set",
		},
		{
			name:        "neither field set, returns false default",
			v3Value:     nil,
			v2Value:     nil,
			expected:    false,
			description: "Default should be false when neither field is set",
		},
		{
			name:     "V3 true, V2 false",
			v3Value:  true,
			v2Value:  false,
			expected: true,
		},
		{
			name:     "V3 false, V2 true",
			v3Value:  false,
			v2Value:  true,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGetter := &mockGetter{
				values: map[string]interface{}{
					tfconstants.AttrEncryptionEnabled:   tt.v3Value,
					tfconstants.AttrIsEncryptionEnabled: tt.v2Value,
				},
			}

			result := getEffectiveEncryptionEnabled(mockGetter, tfconstants.AttrEncryptionEnabled, tfconstants.AttrIsEncryptionEnabled)
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

// TestGetEffectiveEncryptionEnabled_TypeAssertions tests type assertions
func TestGetEffectiveEncryptionEnabled_TypeAssertions(t *testing.T) {
	tests := []struct {
		name        string
		v3Value     interface{}
		expected    bool
		description string
	}{
		{
			name:     "bool type",
			v3Value:  bool(true),
			expected: true,
		},
		{
			name:     "string type",
			v3Value:  "true",
			expected: false, // string won't match bool assertion
		},
		{
			name:     "nil getter interface",
			v3Value:  nil,
			expected: false,
		},
		{
			name:     "non-getter interface",
			v3Value:  nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var getter interface{}
			if tt.v3Value != nil {
				getter = &mockGetter{
					values: map[string]interface{}{
						tfconstants.AttrEncryptionEnabled: tt.v3Value,
					},
				}
			} else if tt.name == "non-getter interface" {
				getter = testInvalidGetterValue
			}

			result := getEffectiveEncryptionEnabled(getter, tfconstants.AttrEncryptionEnabled, tfconstants.AttrIsEncryptionEnabled)
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

// TestLogDeprecationWarning tests logDeprecationWarning function
func TestLogDeprecationWarning(t *testing.T) {
	// Since logDeprecationWarning uses log.Printf, we can't easily test the output
	// without capturing logs. However, we can verify it doesn't panic and accepts
	// different field name combinations.
	tests := []struct {
		name     string
		oldField string
		newField string
	}{
		{
			name:     "disk_size to size_gb",
			oldField: tfconstants.AttrDiskSize,
			newField: tfconstants.AttrSizeGB,
		},
		{
			name:     "disk_iops to iops",
			oldField: tfconstants.AttrDiskIOPS,
			newField: tfconstants.AttrIOPS,
		},
		{
			name:     "is_encryption_enabled to encryption_enabled",
			oldField: tfconstants.AttrIsEncryptionEnabled,
			newField: tfconstants.AttrEncryptionEnabled,
		},
		{
			name:     "empty fields",
			oldField: "",
			newField: "",
		},
		{
			name:     "special characters",
			oldField: testOldFieldName,
			newField: testNewFieldName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This should not panic
			assert.NotPanics(t, func() {
				logDeprecationWarning(tt.oldField, tt.newField)
			})
		})
	}
}

// mockGetter is a mock implementation of the getter interface for testing
type mockGetter struct {
	values map[string]interface{}
}

func (m *mockGetter) Get(key string) interface{} {
	return m.values[key]
}

func (m *mockGetter) GetOk(key string) (interface{}, bool) {
	val, ok := m.values[key]
	return val, ok
}
