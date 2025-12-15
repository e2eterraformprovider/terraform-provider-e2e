package image

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

// Test constants
const (
	testImageID      = "test-image-123"
	testImageIDShort = "test-id" // Used for simple timeout/cancellation tests
)

// mockImageService is a mock implementation of goe2e.ImageService for testing
type mockImageService struct {
	getImageFunc func(ctx context.Context, imageID string) (*goe2e.SavedImage, *goe2e.Response, error)
	callCount    int
}

func (m *mockImageService) GetOSCategories(ctx context.Context) (*goe2e.OSCategoryResponse, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockImageService) GetImagePlans(ctx context.Context, req *goe2e.ImagePlansRequest) ([]goe2e.ImagePlan, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockImageService) ImportImage(ctx context.Context, req *goe2e.ImportImageRequest) (*goe2e.ImageImportResult, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockImageService) GetWindowsImagePermission(ctx context.Context) (*goe2e.WindowsPermission, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockImageService) GetImage(ctx context.Context, imageID string) (*goe2e.SavedImage, *goe2e.Response, error) {
	m.callCount++
	if m.getImageFunc != nil {
		return m.getImageFunc(ctx, imageID)
	}
	return nil, nil, errors.New("not implemented")
}

func (m *mockImageService) GetSavedImages(ctx context.Context) ([]goe2e.SavedImage, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockImageService) DeleteImage(ctx context.Context, imageID string) (*goe2e.DeleteImageResult, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockImageService) RenameImage(ctx context.Context, imageID string, req *goe2e.RenameImageRequest) (*goe2e.RenameImageResult, *goe2e.Response, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *mockImageService) GetPlanDetailsFromPlanName(ctx context.Context, templateID int, planName string) (string, string, *goe2e.Response, error) {
	return "", "", nil, errors.New("not implemented")
}

// mockClient creates a mock goe2e.Client with a mock ImageService
func mockImageClient(service *mockImageService) *goe2e.Client {
	client := &goe2e.Client{}
	client.Images = service
	return client
}

// ============================================================================
// Tests for normalizeImageState()
// ============================================================================

func TestNormalizeImageState_CaseNormalization(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Creating → creating",
			input:    goe2econstants.ImageStatusCreating,
			expected: goe2econstants.ImageStateCreating,
		},
		{
			name:     "Ready → ready",
			input:    goe2econstants.ImageStatusReady,
			expected: goe2econstants.ImageStateReady,
		},
		{
			name:     "Error → error",
			input:    goe2econstants.ImageStatusError,
			expected: goe2econstants.ImageStateError,
		},
		{
			name:     "Deleted → deleted",
			input:    goe2econstants.ImageStatusDeleted,
			expected: goe2econstants.ImageStateDeleted,
		},
		{
			name:     "CREATING → creating",
			input:    "CREATING",
			expected: goe2econstants.ImageStateCreating,
		},
		{
			name:     "READY → ready",
			input:    "READY",
			expected: goe2econstants.ImageStateReady,
		},
		{
			name:     "ERROR → error",
			input:    "ERROR",
			expected: goe2econstants.ImageStateError,
		},
		{
			name:     "DELETED → deleted",
			input:    "DELETED",
			expected: goe2econstants.ImageStateDeleted,
		},
		{
			name:     "creating → creating (lowercase pass through)",
			input:    goe2econstants.ImageStateCreating,
			expected: goe2econstants.ImageStateCreating,
		},
		{
			name:     "ready → ready (lowercase pass through)",
			input:    goe2econstants.ImageStateReady,
			expected: goe2econstants.ImageStateReady,
		},
		{
			name:     "MiXeD cAsE → mixed case",
			input:    "MiXeD cAsE",
			expected: "mixed case",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeImageState(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNormalizeImageState_EmptyString(t *testing.T) {
	result := normalizeImageState("")
	assert.Equal(t, "", result)
}

func TestNormalizeImageState_UnknownState(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Unknown state returns lowercase",
			input:    "UnknownState",
			expected: "unknownstate",
		},
		{
			name:     "Custom state returns lowercase",
			input:    "CustomState",
			expected: "customstate",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeImageState(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNormalizeImageState_WhitespaceHandling(t *testing.T) {
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
			input:    "Ready ",
			expected: "ready ", // Function only lowercases, doesn't trim
		},
		{
			name:     "status with both leading and trailing whitespace",
			input:    " Ready ",
			expected: " ready ", // Function only lowercases, doesn't trim
		},
		{
			name:     "status with internal whitespace",
			input:    "Ready State",
			expected: "ready state",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeImageState(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// ============================================================================
// Tests for waitForImageState()
// ============================================================================

func TestWaitForImageState_SuccessfulStateTransition(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	imageID := testImageID
	callCount := 0

	mockService := &mockImageService{
		getImageFunc: func(ctx context.Context, id string) (*goe2e.SavedImage, *goe2e.Response, error) {
			callCount++
			if callCount == 1 {
				// First call: still creating
				return &goe2e.SavedImage{
					ImageID:    imageID,
					ImageState: goe2econstants.ImageStatusCreating,
				}, nil, nil
			}
			// Second call: ready
			return &goe2e.SavedImage{
				ImageID:    imageID,
				ImageState: goe2econstants.ImageStatusReady,
			}, nil, nil
		},
	}

	client := mockImageClient(mockService)
	ctx := context.Background()

	// Use longer timeout to account for 5 second polling interval
	err := waitForImageState(ctx, client, imageID, goe2econstants.ImageStateReady, 10*time.Second)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, callCount, 2, "should have polled at least twice")
}

func TestWaitForImageState_PollingInterval(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	imageID := testImageID
	callCount := 0
	startTime := time.Now()

	mockService := &mockImageService{
		getImageFunc: func(ctx context.Context, id string) (*goe2e.SavedImage, *goe2e.Response, error) {
			callCount++
			if callCount < 2 {
				// Keep returning creating status
				return &goe2e.SavedImage{
					ImageID:    imageID,
					ImageState: goe2econstants.ImageStatusCreating,
				}, nil, nil
			}
			// Second call: ready
			return &goe2e.SavedImage{
				ImageID:    imageID,
				ImageState: goe2econstants.ImageStatusReady,
			}, nil, nil
		},
	}

	client := mockImageClient(mockService)
	ctx := context.Background()

	// Use longer timeout to account for 5 second polling interval
	err := waitForImageState(ctx, client, imageID, goe2econstants.ImageStateReady, 10*time.Second)
	elapsed := time.Since(startTime)

	require.NoError(t, err)
	assert.Equal(t, 2, callCount, "should have polled 2 times")
	// Polling interval is 5 seconds, so we expect at least 5 seconds elapsed (1 interval)
	assert.GreaterOrEqual(t, elapsed, 5*time.Second, "should have delay between polls (5s interval)")
}

func TestWaitForImageState_Timeout(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	imageID := testImageID

	mockService := &mockImageService{
		getImageFunc: func(ctx context.Context, id string) (*goe2e.SavedImage, *goe2e.Response, error) {
			// Always return creating, never reach ready
			return &goe2e.SavedImage{
				ImageID:    imageID,
				ImageState: goe2econstants.ImageStatusCreating,
			}, nil, nil
		},
	}

	client := mockImageClient(mockService)
	ctx := context.Background()

	// Use short timeout for test
	err := waitForImageState(ctx, client, imageID, goe2econstants.ImageStateReady, 100*time.Millisecond)
	require.Error(t, err)
	// Check error message uses the constant format
	expectedError := fmt.Sprintf(goe2econstants.ImageTimeoutWaitingForState, imageID, goe2econstants.ImageStateReady)
	assert.Equal(t, expectedError, err.Error())
}

func TestWaitForImageState_ContextCancellation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	imageID := testImageID
	callCount := 0

	mockService := &mockImageService{
		getImageFunc: func(ctx context.Context, id string) (*goe2e.SavedImage, *goe2e.Response, error) {
			callCount++
			// Check if context is cancelled
			select {
			case <-ctx.Done():
				return nil, nil, ctx.Err()
			default:
			}
			// Keep returning creating
			return &goe2e.SavedImage{
				ImageID:    imageID,
				ImageState: goe2econstants.ImageStatusCreating,
			}, nil, nil
		},
	}

	client := mockImageClient(mockService)
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel context after a short delay
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	err := waitForImageState(ctx, client, imageID, goe2econstants.ImageStateReady, 5*time.Second)
	require.Error(t, err)
	assert.Equal(t, context.Canceled, err)
}

func TestWaitForImageState_ErrorStateDetection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	imageID := testImageID

	mockService := &mockImageService{
		getImageFunc: func(ctx context.Context, id string) (*goe2e.SavedImage, *goe2e.Response, error) {
			return &goe2e.SavedImage{
				ImageID:    imageID,
				ImageState: goe2econstants.ImageStatusError,
			}, nil, nil
		},
	}

	client := mockImageClient(mockService)
	ctx := context.Background()

	err := waitForImageState(ctx, client, imageID, goe2econstants.ImageStateReady, 10*time.Second)
	require.Error(t, err)
	// Check error message uses the constant format
	expectedError := fmt.Sprintf(goe2econstants.ImageEnteredErrorState, imageID)
	assert.Equal(t, expectedError, err.Error())
}

func TestWaitForImageState_DeletedStateWith404(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	imageID := testImageID

	mockService := &mockImageService{
		getImageFunc: func(ctx context.Context, id string) (*goe2e.SavedImage, *goe2e.Response, error) {
			// Return 404 error (not found) - use constant for consistency
			return nil, nil, fmt.Errorf("image %s %s", imageID, goe2econstants.NotFoundSubstring)
		},
	}

	client := mockImageClient(mockService)
	ctx := context.Background()

	// When desired state is "deleted", 404 is expected and should succeed
	err := waitForImageState(ctx, client, imageID, goe2econstants.ImageStateDeleted, 10*time.Second)
	require.NoError(t, err)
}

func TestWaitForImageState_DeletedStateWithout404(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	imageID := testImageID

	mockService := &mockImageService{
		getImageFunc: func(ctx context.Context, id string) (*goe2e.SavedImage, *goe2e.Response, error) {
			// Return other error (not 404)
			return nil, nil, errors.New("network error")
		},
	}

	client := mockImageClient(mockService)
	ctx := context.Background()

	// When desired state is "deleted" but error is not 404, should fail
	err := waitForImageState(ctx, client, imageID, goe2econstants.ImageStateDeleted, 10*time.Second)
	require.Error(t, err)
	// Verify error message contains the constant format
	assert.Contains(t, err.Error(), "error checking image state")
	assert.Contains(t, err.Error(), "network error")
}

func TestWaitForImageState_APIErrorsDuringPolling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	imageID := testImageID

	mockService := &mockImageService{
		getImageFunc: func(ctx context.Context, id string) (*goe2e.SavedImage, *goe2e.Response, error) {
			// Return API error
			return nil, nil, errors.New("API error occurred")
		},
	}

	client := mockImageClient(mockService)
	ctx := context.Background()

	err := waitForImageState(ctx, client, imageID, goe2econstants.ImageStateReady, 10*time.Second)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error checking image state")
}

func TestWaitForImageState_ImageNotFoundError(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	imageID := testImageID

	mockService := &mockImageService{
		getImageFunc: func(ctx context.Context, id string) (*goe2e.SavedImage, *goe2e.Response, error) {
			// Return not found error (but not for deleted state) - use constant for consistency
			return nil, nil, fmt.Errorf("image %s %s", imageID, goe2econstants.NotFoundSubstring)
		},
	}

	client := mockImageClient(mockService)
	ctx := context.Background()

	// For non-deleted desired state, not found should fail
	err := waitForImageState(ctx, client, imageID, goe2econstants.ImageStateReady, 10*time.Second)
	require.Error(t, err)
	// Verify error message contains the constant format
	assert.Contains(t, err.Error(), "error checking image state")
}

func TestWaitForImageState_MultipleStateChecks(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	imageID := testImageID
	callCount := 0

	mockService := &mockImageService{
		getImageFunc: func(ctx context.Context, id string) (*goe2e.SavedImage, *goe2e.Response, error) {
			callCount++
			if callCount == 1 {
				return &goe2e.SavedImage{
					ImageID:    imageID,
					ImageState: goe2econstants.ImageStatusCreating,
				}, nil, nil
			}
			if callCount == 2 {
				return &goe2e.SavedImage{
					ImageID:    imageID,
					ImageState: goe2econstants.ImageStatusCreating,
				}, nil, nil
			}
			// Third call: ready
			return &goe2e.SavedImage{
				ImageID:    imageID,
				ImageState: goe2econstants.ImageStatusReady,
			}, nil, nil
		},
	}

	client := mockImageClient(mockService)
	ctx := context.Background()

	err := waitForImageState(ctx, client, imageID, goe2econstants.ImageStateReady, 15*time.Second)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, callCount, 3, "should have polled at least 3 times")
}

func TestWaitForImageState_ImmediateMatch(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	imageID := testImageID

	mockService := &mockImageService{
		getImageFunc: func(ctx context.Context, id string) (*goe2e.SavedImage, *goe2e.Response, error) {
			// Already in ready state
			return &goe2e.SavedImage{
				ImageID:    imageID,
				ImageState: goe2econstants.ImageStatusReady,
			}, nil, nil
		},
	}

	client := mockImageClient(mockService)
	ctx := context.Background()

	// Note: waitForImageState now checks immediately before waiting for ticker
	// So if state matches immediately, it returns right away
	err := waitForImageState(ctx, client, imageID, goe2econstants.ImageStateReady, 10*time.Second)
	require.NoError(t, err)
}

// ============================================================================
// Tests for flattenImageResponse()
// ============================================================================

func TestFlattenImageResponse_CompleteObject(t *testing.T) {
	image := &goe2e.SavedImage{
		TemplateID:     123,
		ImageState:     goe2econstants.ImageStatusReady,
		ImageType:      "snapshot",
		OSDistribution: "Ubuntu",
		Name:           "test-image",
		ImageID:        "img-123",
		Distro:         "ubuntu",
		SKUType:        "standard",
		ImageSize:      "95.368 GB",
		CloningOps:     "0",
		RunningVMs:     "2",
		IsWindows:      false,
		CreationTime:   "2024-01-01T00:00:00Z",
		VMInfo:         []interface{}{map[string]interface{}{"vm_id": 456}},
	}

	result := flattenImageResponse(image)

	assert.Equal(t, 123, result[tfconstants.AttrTemplateID])
	assert.Equal(t, goe2econstants.ImageStatusReady, result["image_state"])
	assert.Equal(t, goe2econstants.ImageStateReady, result["state"])
	assert.Equal(t, "snapshot", result["image_type"])
	assert.Equal(t, "Ubuntu", result["os_distribution"])
	assert.Equal(t, "test-image", result[tfconstants.AttrName])
	assert.Equal(t, "img-123", result["image_id"])
	assert.Equal(t, "ubuntu", result["distro"])
	assert.Equal(t, "standard", result["sku_type"])
	assert.Equal(t, "95.368 GB", result["image_size"])
	assert.Equal(t, "0", result["cloning_ops"])
	assert.Equal(t, "2", result["running_vms"])
	assert.Equal(t, false, result["is_windows"])
	assert.Equal(t, "2024-01-01T00:00:00Z", result[tfconstants.AttrCreatedAt])
	assert.Equal(t, []interface{}{map[string]interface{}{"vm_id": 456}}, result["vm_info"])
}

func TestFlattenImageResponse_MinimalObject(t *testing.T) {
	image := &goe2e.SavedImage{
		ImageID: "img-123",
	}

	result := flattenImageResponse(image)

	// Check that all fields are present
	assert.Contains(t, result, tfconstants.AttrTemplateID)
	assert.Contains(t, result, "image_state")
	assert.Contains(t, result, "state")
	assert.Contains(t, result, "image_type")
	assert.Contains(t, result, "os_distribution")
	assert.Contains(t, result, tfconstants.AttrName)
	assert.Contains(t, result, "image_id")
	assert.Contains(t, result, "distro")
	assert.Contains(t, result, "sku_type")
	assert.Contains(t, result, "image_size")
	assert.Contains(t, result, "cloning_ops")
	assert.Contains(t, result, "running_vms")
	assert.Contains(t, result, "is_windows")
	assert.Contains(t, result, tfconstants.AttrCreatedAt)
	assert.Contains(t, result, "vm_info")

	// Check image_id is set correctly
	assert.Equal(t, "img-123", result["image_id"])
}

func TestFlattenImageResponse_StateNormalization(t *testing.T) {
	tests := []struct {
		name          string
		imageState    string
		expectedState string
	}{
		{
			name:          "Creating state normalized",
			imageState:    goe2econstants.ImageStatusCreating,
			expectedState: goe2econstants.ImageStateCreating,
		},
		{
			name:          "Ready state normalized",
			imageState:    goe2econstants.ImageStatusReady,
			expectedState: goe2econstants.ImageStateReady,
		},
		{
			name:          "Error state normalized",
			imageState:    goe2econstants.ImageStatusError,
			expectedState: goe2econstants.ImageStateError,
		},
		{
			name:          "Deleted state normalized",
			imageState:    goe2econstants.ImageStatusDeleted,
			expectedState: goe2econstants.ImageStateDeleted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			image := &goe2e.SavedImage{
				ImageID:    "img-123",
				ImageState: tt.imageState,
			}

			result := flattenImageResponse(image)

			assert.Equal(t, tt.imageState, result["image_state"], "image_state should be unchanged")
			assert.Equal(t, tt.expectedState, result["state"], "state should be normalized")
		})
	}
}

func TestFlattenImageResponse_EmptyZeroValues(t *testing.T) {
	image := &goe2e.SavedImage{
		ImageID:        "img-123",
		TemplateID:     0,
		ImageState:     "",
		ImageType:      "",
		OSDistribution: "",
		Name:           "",
		Distro:         "",
		SKUType:        "",
		ImageSize:      "",
		CloningOps:     "",
		RunningVMs:     "",
		IsWindows:      false,
		CreationTime:   "",
		VMInfo:         nil,
	}

	result := flattenImageResponse(image)

	assert.Equal(t, 0, result[tfconstants.AttrTemplateID])
	assert.Equal(t, "", result["image_state"])
	assert.Equal(t, "", result["state"])
	assert.Equal(t, "", result["image_type"])
	assert.Equal(t, "", result["os_distribution"])
	assert.Equal(t, "", result[tfconstants.AttrName])
	assert.Equal(t, "img-123", result["image_id"])
	assert.Equal(t, "", result["distro"])
	assert.Equal(t, "", result["sku_type"])
	assert.Equal(t, "", result["image_size"])
	assert.Equal(t, "", result["cloning_ops"])
	assert.Equal(t, "", result["running_vms"])
	assert.Equal(t, false, result["is_windows"])
	assert.Equal(t, "", result[tfconstants.AttrCreatedAt])
	assert.Nil(t, result["vm_info"])
}

func TestFlattenImageResponse_NilImage(t *testing.T) {
	// This test verifies behavior with nil input
	// In Go, calling methods on nil pointer will panic, so we test that it panics
	assert.Panics(t, func() {
		flattenImageResponse(nil)
	}, "flattenImageResponse should panic on nil input")
}
