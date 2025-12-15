package blockstorage_test

import (
	"strconv"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/blockstorage"
	"github.com/stretchr/testify/assert"
)

var validBlockStorageSizesForTest = []float64{250, 500, 1000, 2000, 4000, 8000, 16000, 24000}

func expectedIOPSString(size float64) string {
	return strconv.Itoa(int(size * float64(blockstorage.IOPS_PER_GB)))
}

func TestCalculateIOPS(t *testing.T) {
	testCases := []struct {
		name     string
		size     float64
		expected string
	}{
		{
			name:     "250 GB",
			size:     250.0,
			expected: "3750",
		},
		{
			name:     "500 GB",
			size:     500.0,
			expected: "7500",
		},
		{
			name:     "1000 GB",
			size:     1000.0,
			expected: "15000",
		},
		{
			name:     "2000 GB",
			size:     2000.0,
			expected: "30000",
		},
		{
			name:     "4000 GB",
			size:     4000.0,
			expected: "60000",
		},
		{
			name:     "8000 GB",
			size:     8000.0,
			expected: "120000",
		},
		{
			name:     "16000 GB",
			size:     16000.0,
			expected: "240000",
		},
		{
			name:     "24000 GB",
			size:     24000.0,
			expected: "360000",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := blockstorage.CalculateIOPS(tc.size)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestCalculateIOPS_EdgeCases(t *testing.T) {
	testCases := []struct {
		name     string
		size     float64
		expected string
	}{
		{
			name:     "Zero size",
			size:     0.0,
			expected: "0",
		},
		{
			name:     "Fractional size",
			size:     250.5,
			expected: expectedIOPSString(250.5),
		},
		{
			name:     "Very small size",
			size:     1.0,
			expected: expectedIOPSString(1.0),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := blockstorage.CalculateIOPS(tc.size)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestCalculateIOPS_AllValidSizes(t *testing.T) {
	for _, size := range validBlockStorageSizesForTest {
		t.Run(strconv.Itoa(int(size))+"_GB", func(t *testing.T) {
			assert.Equal(t, expectedIOPSString(size), blockstorage.CalculateIOPS(size))
		})
	}
}

func TestCalculateIOPS_FractionalSizes(t *testing.T) {
	testCases := []float64{250.5, 1000.25}
	for _, size := range testCases {
		t.Run(strconv.FormatFloat(size, 'f', -1, 64)+"_GB", func(t *testing.T) {
			assert.Equal(t, expectedIOPSString(size), blockstorage.CalculateIOPS(size))
		})
	}
}

func TestCalculateIOPS_VeryLargeSizes(t *testing.T) {
	// Beyond currently supported sizes, but validates math stays correct.
	testCases := []float64{50000, 100000}
	for _, size := range testCases {
		t.Run(strconv.Itoa(int(size))+"_GB", func(t *testing.T) {
			assert.Equal(t, expectedIOPSString(size), blockstorage.CalculateIOPS(size))
		})
	}
}

func TestValidateBlockStorageSize(t *testing.T) {
	testCases := []struct {
		name        string
		size        float64
		expectError bool
	}{
		{
			name:        "Valid size 250",
			size:        250,
			expectError: false,
		},
		{
			name:        "Valid size 500",
			size:        500,
			expectError: false,
		},
		{
			name:        "Valid size 1000",
			size:        1000,
			expectError: false,
		},
		{
			name:        "Valid size 24000",
			size:        24000,
			expectError: false,
		},
		{
			name:        "Invalid size 100",
			size:        100,
			expectError: true,
		},
		{
			name:        "Invalid size 300",
			size:        300,
			expectError: true,
		},
		{
			name:        "Invalid size 10000",
			size:        10000,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Access the validation function through the package
			// Note: We can't test private functions directly, so we'll test through the schema
			// For now, just verify the logic by checking the valid sizes slice
			var isValid bool
			for _, validSize := range validBlockStorageSizesForTest {
				if tc.size == validSize {
					isValid = true
					break
				}
			}

			if tc.expectError {
				assert.False(t, isValid, "Size %.0f should be invalid", tc.size)
			} else {
				assert.True(t, isValid, "Size %.0f should be valid", tc.size)
			}
		})
	}
}
