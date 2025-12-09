package sfs

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNormalizeSfsState tests the state normalization function
func TestNormalizeSfsState(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "creating status",
			input:    "creating",
			expected: "creating",
		},
		{
			name:     "Creating uppercase",
			input:    "Creating",
			expected: "creating",
		},
		{
			name:     "active status",
			input:    "active",
			expected: "active",
		},
		{
			name:     "Active uppercase",
			input:    "Active",
			expected: "active",
		},
		{
			name:     "deleting status",
			input:    "deleting",
			expected: "deleting",
		},
		{
			name:     "deleted status",
			input:    "deleted",
			expected: "deleted",
		},
		{
			name:     "error status",
			input:    "error",
			expected: "error",
		},
		{
			name:     "Error uppercase",
			input:    "Error",
			expected: "error",
		},
		{
			name:     "unknown status",
			input:    "unknown",
			expected: "unknown",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "mixed case unknown",
			input:    "SoMeStAtUs",
			expected: "somestatus",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeSfsState(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestParseSfsImportID tests the import ID parsing function
func TestParseSfsImportID(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedProjID string
		expectedRegion string
		expectedSfsID  string
		expectError    bool
		errorMsg       string
	}{
		{
			name:           "simple format with ID only",
			input:          "12345",
			expectedProjID: "",
			expectedRegion: "",
			expectedSfsID:  "12345",
			expectError:    false,
		},
		{
			name:           "simple format with alphanumeric ID",
			input:          "sfs-abc123-def456",
			expectedProjID: "",
			expectedRegion: "",
			expectedSfsID:  "sfs-abc123-def456",
			expectError:    false,
		},
		{
			name:           "full format with all parts",
			input:          "proj-123/us-east-1/sfs-456",
			expectedProjID: "proj-123",
			expectedRegion: "us-east-1",
			expectedSfsID:  "sfs-456",
			expectError:    false,
		},
		{
			name:        "invalid format with 2 parts",
			input:       "proj-123/us-east-1",
			expectError: true,
		},
		{
			name:        "invalid format with 4 parts",
			input:       "proj/region/sfs/extra",
			expectError: true,
		},
		{
			name:        "full format with empty project",
			input:       "/region/sfs-id",
			expectError: true,
		},
		{
			name:        "full format with empty region",
			input:       "proj-id//sfs-id",
			expectError: true,
		},
		{
			name:        "full format with empty SFS ID",
			input:       "proj-id/region/",
			expectError: true,
		},
		{
			name:          "empty string",
			input:         "",
			expectedSfsID: "",
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projID, region, sfsID, err := parseSfsImportID(tt.input)

			if tt.expectError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedProjID, projID)
				assert.Equal(t, tt.expectedRegion, region)
				assert.Equal(t, tt.expectedSfsID, sfsID)
			}
		})
	}
}

// TestWaitForSfsStatusContextCancellation tests context cancellation
func TestWaitForSfsStatusContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Immediately cancel

	// Should return context cancelled error
	err := waitForSfsStatus(ctx, nil, "test-id", "active", 5*time.Second)
	require.Error(t, err)
	assert.Equal(t, context.Canceled, err)
}

// TestWaitForSfsStatusTimeout tests timeout handling
func TestWaitForSfsStatusTimeout(t *testing.T) {
	ctx := context.Background()

	// This should timeout since we're not mocking the client
	err := waitForSfsStatus(ctx, nil, "test-id", "active", 100*time.Millisecond)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "timeout waiting for SFS")
}

// TestParseSfsImportIDEdgeCases tests edge cases
func TestParseSfsImportIDEdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		shouldFail bool
	}{
		{
			name:       "ID with spaces is valid (no validation on content)",
			input:      "sfs id 123",
			shouldFail: false,
		},
		{
			name:       "ID with too many slashes",
			input:      "proj/region/sfs/extra/parts",
			shouldFail: true,
		},
		{
			name:       "single character ID",
			input:      "a",
			shouldFail: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, _, err := parseSfsImportID(tt.input)
			if tt.shouldFail {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// BenchmarkNormalizeSfsState benchmarks the normalization function
func BenchmarkNormalizeSfsState(b *testing.B) {
	for i := 0; i < b.N; i++ {
		normalizeSfsState("Active")
	}
}

// BenchmarkParseSfsImportID benchmarks the import ID parser
func BenchmarkParseSfsImportID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseSfsImportID("proj-123/us-east-1/sfs-456")
	}
}
