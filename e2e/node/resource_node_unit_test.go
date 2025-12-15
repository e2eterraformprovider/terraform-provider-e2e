package node

import (
	"strings"
	"testing"

	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestValidateName(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		wantError bool
		errorMsg  string
	}{
		{
			name:      "valid name - alphanumeric with dash",
			input:     "test-node-123",
			wantError: false,
		},
		{
			name:      "valid name - alphanumeric with underscore",
			input:     "test_node_123",
			wantError: false,
		},
		{
			name:      "valid name - minimum length (1 char)",
			input:     "a",
			wantError: false,
		},
		{
			name:      "valid name - maximum length (50 chars)",
			input:     strings.Repeat("a", 50),
			wantError: false,
		},
		{
			name:      "empty string",
			input:     "",
			wantError: true,
			errorMsg:  "name cannot be empty",
		},
		{
			name:      "non-string type - integer",
			input:     123,
			wantError: true,
			errorMsg:  "expected name to be string",
		},
		{
			name:      "non-string type - boolean",
			input:     true,
			wantError: true,
			errorMsg:  "expected name to be string",
		},
		{
			name:      "too long - 51 characters",
			input:     strings.Repeat("a", 51),
			wantError: true,
			errorMsg:  "cannot be blank, must not contain whitespace or special characters",
		},
		{
			name:      "contains whitespace",
			input:     "test node",
			wantError: true,
			errorMsg:  "cannot be blank, must not contain whitespace or special characters",
		},
		{
			name:      "contains special characters - @",
			input:     "test@node",
			wantError: true,
			errorMsg:  "cannot be blank, must not contain whitespace or special characters",
		},
		{
			name:      "contains special characters - .",
			input:     "test.node",
			wantError: true,
			errorMsg:  "cannot be blank, must not contain whitespace or special characters",
		},
		{
			name:      "contains special characters - #",
			input:     "test#node",
			wantError: true,
			errorMsg:  "cannot be blank, must not contain whitespace or special characters",
		},
		{
			name:      "starts with dash",
			input:     "-test-node",
			wantError: false, // Dash is allowed
		},
		{
			name:      "starts with underscore",
			input:     "_test-node",
			wantError: false, // Underscore is allowed
		},
		{
			name:      "only numbers",
			input:     "123456",
			wantError: false,
		},
		{
			name:      "mixed case",
			input:     "Test-Node-123",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warns, errs := ValidateName(tt.input, "name")

			if tt.wantError {
				if len(errs) == 0 {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errorMsg != "" && !strings.Contains(errs[0].Error(), tt.errorMsg) {
					t.Errorf("error message mismatch. got: %v, want to contain: %s", errs[0].Error(), tt.errorMsg)
				}
			} else {
				if len(errs) > 0 {
					t.Errorf("unexpected error: %v", errs)
				}
				if len(warns) > 0 {
					t.Errorf("unexpected warnings: %v", warns)
				}
			}
		})
	}
}

func TestValidatePlanName(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		wantError bool
		errorMsg  string
	}{
		{
			name:      "valid plan name - alphanumeric",
			input:     "c2-2c-4gb",
			wantError: false,
		},
		{
			name:      "valid plan name - with dash",
			input:     "C3-8GB",
			wantError: false,
		},
		{
			name:      "valid plan name - uppercase",
			input:     "C3.8GB",
			wantError: false,
		},
		{
			name:      "empty string",
			input:     "",
			wantError: true,
			errorMsg:  "plan name cannot be empty",
		},
		{
			name:      "non-string type - integer",
			input:     123,
			wantError: true,
			errorMsg:  "expected plan to be string",
		},
		{
			name:      "non-string type - boolean",
			input:     false,
			wantError: true,
			errorMsg:  "expected plan to be string",
		},
		{
			name:      "contains whitespace - single space",
			input:     "c2 2c 4gb",
			wantError: true,
			errorMsg:  "plan cannot contain whitespace",
		},
		{
			name:      "contains whitespace - tab",
			input:     "c2\t2c\t4gb",
			wantError: true,
			errorMsg:  "plan cannot contain whitespace",
		},
		{
			name:      "contains whitespace - newline",
			input:     "c2\n2c\n4gb",
			wantError: true,
			errorMsg:  "plan cannot contain whitespace",
		},
		{
			name:      "contains whitespace - leading",
			input:     " c2-2c-4gb",
			wantError: true,
			errorMsg:  "plan cannot contain whitespace",
		},
		{
			name:      "contains whitespace - trailing",
			input:     "c2-2c-4gb ",
			wantError: true,
			errorMsg:  "plan cannot contain whitespace",
		},
		{
			name:      "valid - no whitespace",
			input:     "c2-2c-4gb",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warns, errs := ValidatePlanName(tt.input, "plan")

			if tt.wantError {
				if len(errs) == 0 {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errorMsg != "" && !strings.Contains(errs[0].Error(), tt.errorMsg) {
					t.Errorf("error message mismatch. got: %v, want to contain: %s", errs[0].Error(), tt.errorMsg)
				}
			} else {
				if len(errs) > 0 {
					t.Errorf("unexpected error: %v", errs)
				}
				if len(warns) > 0 {
					t.Errorf("unexpected warnings: %v", warns)
				}
			}
		})
	}
}

func TestValidateBlank(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		fieldName string
		wantError bool
		errorMsg  string
	}{
		{
			name:      "valid - non-blank string",
			input:     "test-value",
			fieldName: "image",
			wantError: false,
		},
		{
			name:      "valid - string with content",
			input:     "ubuntu-20.04",
			fieldName: "image",
			wantError: false,
		},
		{
			name:      "empty string",
			input:     "",
			fieldName: "image",
			wantError: true,
			errorMsg:  "cannot be blank",
		},
		{
			name:      "whitespace only - single space",
			input:     " ",
			fieldName: "image",
			wantError: true,
			errorMsg:  "cannot be blank",
		},
		{
			name:      "whitespace only - multiple spaces",
			input:     "   ",
			fieldName: "image",
			wantError: true,
			errorMsg:  "cannot be blank",
		},
		{
			name:      "whitespace only - tab",
			input:     "\t",
			fieldName: "image",
			wantError: true,
			errorMsg:  "cannot be blank",
		},
		{
			name:      "whitespace only - newline",
			input:     "\n",
			fieldName: "image",
			wantError: true,
			errorMsg:  "cannot be blank",
		},
		{
			name:      "whitespace only - mixed",
			input:     " \t\n ",
			fieldName: "image",
			wantError: true,
			errorMsg:  "cannot be blank",
		},
		{
			name:      "leading and trailing whitespace - valid content",
			input:     "  ubuntu-20.04  ",
			fieldName: "image",
			wantError: false, // TrimSpace makes this valid
		},
		{
			name:      "non-string type - integer",
			input:     123,
			fieldName: "image",
			wantError: true,
			errorMsg:  "expected image to be string",
		},
		{
			name:      "non-string type - boolean",
			input:     true,
			fieldName: "image",
			wantError: true,
			errorMsg:  "expected image to be string",
		},
		{
			name:      "non-string type - nil",
			input:     nil,
			fieldName: "image",
			wantError: true,
			errorMsg:  "expected image to be string",
		},
		{
			name:      "valid - field name in error message",
			input:     "",
			fieldName: "custom_field",
			wantError: true,
			errorMsg:  "custom_field cannot be blank",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warns, errs := ValidateBlank(tt.input, tt.fieldName)

			if tt.wantError {
				if len(errs) == 0 {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errorMsg != "" && !strings.Contains(errs[0].Error(), tt.errorMsg) {
					t.Errorf("error message mismatch. got: %v, want to contain: %s", errs[0].Error(), tt.errorMsg)
				}
			} else {
				if len(errs) > 0 {
					t.Errorf("unexpected error: %v", errs)
				}
				if len(warns) > 0 {
					t.Errorf("unexpected warnings: %v", warns)
				}
			}
		})
	}
}

func TestValidateInteger(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		fieldName string
		wantError bool
		errorMsg  string
	}{
		{
			name:      "valid - positive integer string",
			input:     "123",
			fieldName: "block_storage_id",
			wantError: false,
		},
		{
			name:      "valid - zero",
			input:     "0",
			fieldName: "block_storage_id",
			wantError: false,
		},
		{
			name:      "valid - large number",
			input:     "999999",
			fieldName: "block_storage_id",
			wantError: false,
		},
		{
			name:      "valid - single digit",
			input:     "5",
			fieldName: "block_storage_id",
			wantError: false,
		},
		{
			name:      "empty string",
			input:     "",
			fieldName: "block_storage_id",
			wantError: true,
			errorMsg:  "only contains numeric value",
		},
		{
			name:      "non-numeric string - letters",
			input:     "abc",
			fieldName: "block_storage_id",
			wantError: true,
			errorMsg:  "only contains numeric value",
		},
		{
			name:      "non-numeric string - mixed",
			input:     "123abc",
			fieldName: "block_storage_id",
			wantError: true,
			errorMsg:  "only contains numeric value",
		},
		{
			name:      "non-numeric string - prefix letters",
			input:     "abc123",
			fieldName: "block_storage_id",
			wantError: true,
			errorMsg:  "only contains numeric value",
		},
		{
			name:      "non-numeric string - special characters",
			input:     "123-456",
			fieldName: "block_storage_id",
			wantError: true,
			errorMsg:  "only contains numeric value",
		},
		{
			name:      "non-numeric string - decimal",
			input:     "123.45",
			fieldName: "block_storage_id",
			wantError: true,
			errorMsg:  "only contains numeric value",
		},
		{
			name:      "non-string type - integer",
			input:     123,
			fieldName: "block_storage_id",
			wantError: true,
			errorMsg:  "expected block_storage_id to be string",
		},
		{
			name:      "non-string type - boolean",
			input:     true,
			fieldName: "block_storage_id",
			wantError: true,
			errorMsg:  "expected block_storage_id to be string",
		},
		{
			name:      "whitespace - leading",
			input:     " 123",
			fieldName: "block_storage_id",
			wantError: true,
			errorMsg:  "only contains numeric value",
		},
		{
			name:      "whitespace - trailing",
			input:     "123 ",
			fieldName: "block_storage_id",
			wantError: true,
			errorMsg:  "only contains numeric value",
		},
		{
			name:      "negative number - valid",
			input:     "-123",
			fieldName: "block_storage_id",
			wantError: false, // strconv.Atoi accepts negative numbers
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warns, errs := ValidateInteger(tt.input, tt.fieldName)

			if tt.wantError {
				if len(errs) == 0 {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errorMsg != "" && !strings.Contains(errs[0].Error(), tt.errorMsg) {
					t.Errorf("error message mismatch. got: %v, want to contain: %s", errs[0].Error(), tt.errorMsg)
				}
			} else {
				if len(errs) > 0 {
					t.Errorf("unexpected error: %v", errs)
				}
				if len(warns) > 0 {
					t.Errorf("unexpected warnings: %v", warns)
				}
			}
		})
	}
}

func TestConvertStringToInt(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      int
		wantError bool
	}{
		{
			name:      "valid - positive integer",
			input:     "123",
			want:      123,
			wantError: false,
		},
		{
			name:      "valid - zero",
			input:     "0",
			want:      0,
			wantError: false,
		},
		{
			name:      "valid - large number",
			input:     "999999",
			want:      999999,
			wantError: false,
		},
		{
			name:      "valid - single digit",
			input:     "5",
			want:      5,
			wantError: false,
		},
		{
			name:      "invalid - non-numeric",
			input:     "abc",
			want:      0,
			wantError: true,
		},
		{
			name:      "invalid - mixed alphanumeric",
			input:     "123abc",
			want:      0,
			wantError: true,
		},
		{
			name:      "invalid - empty string",
			input:     "",
			want:      0,
			wantError: true,
		},
		{
			name:      "invalid - decimal",
			input:     "123.45",
			want:      0,
			wantError: true,
		},
		{
			name:      "valid - negative number",
			input:     "-123",
			want:      -123,
			wantError: false, // strconv.Atoi accepts negative numbers
		},
		{
			name:      "invalid - special characters",
			input:     "123-456",
			want:      0,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertStringToInt(tt.input)

			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if got != tt.want {
					t.Errorf("convertStringToInt() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestGetDefaultSGFromList(t *testing.T) {
	tests := []struct {
		name           string
		securityGroups []goe2e.SecurityGroupInfo
		want           int
	}{
		{
			name: "default SG found - first in list",
			securityGroups: []goe2e.SecurityGroupInfo{
				{ID: 100, IsDefault: true},
				{ID: 200, IsDefault: false},
				{ID: 300, IsDefault: false},
			},
			want: 100,
		},
		{
			name: "default SG found - middle of list",
			securityGroups: []goe2e.SecurityGroupInfo{
				{ID: 100, IsDefault: false},
				{ID: 200, IsDefault: true},
				{ID: 300, IsDefault: false},
			},
			want: 200,
		},
		{
			name: "default SG found - last in list",
			securityGroups: []goe2e.SecurityGroupInfo{
				{ID: 100, IsDefault: false},
				{ID: 200, IsDefault: false},
				{ID: 300, IsDefault: true},
			},
			want: 300,
		},
		{
			name: "no default SG - returns 0",
			securityGroups: []goe2e.SecurityGroupInfo{
				{ID: 100, IsDefault: false},
				{ID: 200, IsDefault: false},
				{ID: 300, IsDefault: false},
			},
			want: 0,
		},
		{
			name:           "empty list - returns 0",
			securityGroups: []goe2e.SecurityGroupInfo{},
			want:           0,
		},
		{
			name:           "nil list - returns 0",
			securityGroups: nil,
			want:           0,
		},
		{
			name: "multiple defaults - returns first",
			securityGroups: []goe2e.SecurityGroupInfo{
				{ID: 100, IsDefault: true},
				{ID: 200, IsDefault: true},
				{ID: 300, IsDefault: true},
			},
			want: 100, // Returns first default found
		},
		{
			name: "single SG - default",
			securityGroups: []goe2e.SecurityGroupInfo{
				{ID: 500, IsDefault: true},
			},
			want: 500,
		},
		{
			name: "single SG - not default",
			securityGroups: []goe2e.SecurityGroupInfo{
				{ID: 500, IsDefault: false},
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getDefaultSGFromList(tt.securityGroups)
			if got != tt.want {
				t.Errorf("getDefaultSGFromList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUniqueArrayElements(t *testing.T) {
	tests := []struct {
		name     string
		arr1     []interface{}
		arr2     []interface{}
		expected []interface{}
	}{
		{
			name:     "unique elements found",
			arr1:     []interface{}{1, 2, 3},
			arr2:     []interface{}{2, 3, 4},
			expected: []interface{}{1},
		},
		{
			name:     "no unique elements",
			arr1:     []interface{}{1, 2, 3},
			arr2:     []interface{}{1, 2, 3},
			expected: []interface{}{},
		},
		{
			name:     "all unique",
			arr1:     []interface{}{1, 2, 3},
			arr2:     []interface{}{4, 5, 6},
			expected: []interface{}{1, 2, 3},
		},
		{
			name:     "empty arr1",
			arr1:     []interface{}{},
			arr2:     []interface{}{1, 2, 3},
			expected: []interface{}{},
		},
		{
			name:     "empty arr2",
			arr1:     []interface{}{1, 2, 3},
			arr2:     []interface{}{},
			expected: []interface{}{1, 2, 3},
		},
		{
			name:     "both empty",
			arr1:     []interface{}{},
			arr2:     []interface{}{},
			expected: []interface{}{},
		},
		{
			name:     "duplicates in arr1",
			arr1:     []interface{}{1, 1, 2, 2, 3},
			arr2:     []interface{}{2, 3},
			expected: []interface{}{1, 1}, // Preserves duplicates
		},
		{
			name:     "string elements",
			arr1:     []interface{}{"a", "b", "c"},
			arr2:     []interface{}{"b", "c", "d"},
			expected: []interface{}{"a"},
		},
		{
			name:     "mixed types",
			arr1:     []interface{}{1, "a", 2},
			arr2:     []interface{}{2, "b"},
			expected: []interface{}{1, "a"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UniqueArrayElements(tt.arr1, tt.arr2)

			if len(got) != len(tt.expected) {
				t.Errorf("UniqueArrayElements() length = %v, want %v", len(got), len(tt.expected))
				return
			}

			// Check that all expected elements are present
			for i, expected := range tt.expected {
				if i < len(got) && got[i] != expected {
					t.Errorf("UniqueArrayElements()[%d] = %v, want %v", i, got[i], expected)
				}
			}
		})
	}
}

func TestCommonArrayElements(t *testing.T) {
	tests := []struct {
		name     string
		arr1     []interface{}
		arr2     []interface{}
		expected []interface{}
	}{
		{
			name:     "common elements found",
			arr1:     []interface{}{1, 2, 3},
			arr2:     []interface{}{2, 3, 4},
			expected: []interface{}{2, 3},
		},
		{
			name:     "no common elements",
			arr1:     []interface{}{1, 2, 3},
			arr2:     []interface{}{4, 5, 6},
			expected: []interface{}{},
		},
		{
			name:     "all common",
			arr1:     []interface{}{1, 2, 3},
			arr2:     []interface{}{1, 2, 3},
			expected: []interface{}{1, 2, 3},
		},
		{
			name:     "empty arr1",
			arr1:     []interface{}{},
			arr2:     []interface{}{1, 2, 3},
			expected: []interface{}{},
		},
		{
			name:     "empty arr2",
			arr1:     []interface{}{1, 2, 3},
			arr2:     []interface{}{},
			expected: []interface{}{},
		},
		{
			name:     "both empty",
			arr1:     []interface{}{},
			arr2:     []interface{}{},
			expected: []interface{}{},
		},
		{
			name:     "duplicates in arr1",
			arr1:     []interface{}{1, 1, 2, 2, 3},
			arr2:     []interface{}{2, 3},
			expected: []interface{}{2, 2, 3}, // Returns all from arr1 that are in arr2 (preserves duplicates)
		},
		{
			name:     "string elements",
			arr1:     []interface{}{"a", "b", "c"},
			arr2:     []interface{}{"b", "c", "d"},
			expected: []interface{}{"b", "c"},
		},
		{
			name:     "single common element",
			arr1:     []interface{}{1, 2, 3, 4, 5},
			arr2:     []interface{}{3},
			expected: []interface{}{3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CommonArrayElements(tt.arr1, tt.arr2)

			if len(got) != len(tt.expected) {
				t.Errorf("CommonArrayElements() length = %v, want %v", len(got), len(tt.expected))
				return
			}

			// Check that all expected elements are present (order may vary)
			for i, expected := range tt.expected {
				if i < len(got) && got[i] != expected {
					t.Errorf("CommonArrayElements()[%d] = %v, want %v", i, got[i], expected)
				}
			}
		})
	}
}

func TestRemoveArrayElement(t *testing.T) {
	tests := []struct {
		name     string
		arr      []interface{}
		val      interface{}
		expected []interface{}
	}{
		{
			name:     "element found and removed",
			arr:      []interface{}{1, 2, 3, 4},
			val:      3,
			expected: []interface{}{1, 2, 4},
		},
		{
			name:     "element not found",
			arr:      []interface{}{1, 2, 3},
			val:      5,
			expected: []interface{}{1, 2, 3},
		},
		{
			name:     "empty array",
			arr:      []interface{}{},
			val:      1,
			expected: []interface{}{},
		},
		{
			name:     "remove first element",
			arr:      []interface{}{1, 2, 3},
			val:      1,
			expected: []interface{}{2, 3},
		},
		{
			name:     "remove last element",
			arr:      []interface{}{1, 2, 3},
			val:      3,
			expected: []interface{}{1, 2},
		},
		{
			name:     "remove only element",
			arr:      []interface{}{1},
			val:      1,
			expected: []interface{}{},
		},
		{
			name:     "multiple occurrences - removes all",
			arr:      []interface{}{1, 2, 1, 3, 1},
			val:      1,
			expected: []interface{}{2, 3}, // All occurrences removed
		},
		{
			name:     "string element",
			arr:      []interface{}{"a", "b", "c"},
			val:      "b",
			expected: []interface{}{"a", "c"},
		},
		{
			name:     "different types",
			arr:      []interface{}{1, "a", 2},
			val:      "a",
			expected: []interface{}{1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeArrayElement(tt.arr, tt.val)

			if len(got) != len(tt.expected) {
				t.Errorf("removeArrayElement() length = %v, want %v", len(got), len(tt.expected))
				return
			}

			for i, expected := range tt.expected {
				if i < len(got) && got[i] != expected {
					t.Errorf("removeArrayElement()[%d] = %v, want %v", i, got[i], expected)
				}
			}
		})
	}
}

func TestCustomImportStateFunc(t *testing.T) {
	tests := []struct {
		name           string
		importID       string
		wantError      bool
		expectedID     string
		expectedProjID string
		expectedRegion string
	}{
		{
			name:           "valid format - three parts",
			importID:       "project-123/Mumbai/node-456",
			wantError:      false,
			expectedID:     "node-456",
			expectedProjID: "project-123",
			expectedRegion: "Mumbai",
		},
		{
			name:           "valid format - with numbers",
			importID:       "123/region-1/789",
			wantError:      false,
			expectedID:     "789",
			expectedProjID: "123",
			expectedRegion: "region-1",
		},
		{
			name:      "invalid format - too few parts (1 part)",
			importID:  "project-123",
			wantError: true,
		},
		{
			name:      "invalid format - too few parts (2 parts)",
			importID:  "project-123/Mumbai",
			wantError: true,
		},
		{
			name:      "invalid format - too many parts (4 parts)",
			importID:  "project-123/Mumbai/node-456/extra",
			wantError: true,
		},
		{
			name:      "invalid format - too many parts (5 parts)",
			importID:  "a/b/c/d/e",
			wantError: true,
		},
		{
			name:      "empty string",
			importID:  "",
			wantError: true,
		},
		{
			name:           "valid format - empty parts (edge case)",
			importID:       "//node-456",
			wantError:      false, // Function doesn't validate empty parts
			expectedID:     "node-456",
			expectedProjID: "",
			expectedRegion: "",
		},
		{
			name:           "valid format - special characters in parts",
			importID:       "proj-123/reg_ion-1/node_456",
			wantError:      false,
			expectedID:     "node_456",
			expectedProjID: "proj-123",
			expectedRegion: "reg_ion-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a ResourceData using TestResourceDataRaw
			resource := ResourceNode()
			d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})

			// Set the import ID
			d.SetId(tt.importID)

			result, err := CustomImportStateFunc(d, nil)

			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				if err != nil && !strings.Contains(err.Error(), nodeImportFormatError) {
					t.Errorf("error message should contain import format error")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}

				if len(result) == 0 {
					t.Errorf("expected result but got empty slice")
					return
				}

				resultData := result[0]
				if resultData.Id() != tt.expectedID {
					t.Errorf("ID = %v, want %v", resultData.Id(), tt.expectedID)
				}

				projID := resultData.Get(tfconstants.AttrProjectID)
				if projID == nil || projID.(string) != tt.expectedProjID {
					t.Errorf("ProjectID = %v, want %v", projID, tt.expectedProjID)
				}

				region := resultData.Get("region")
				if region == nil || region.(string) != tt.expectedRegion {
					t.Errorf("Region = %v, want %v", region, tt.expectedRegion)
				}
			}
		})
	}
}

// TestResourceNodeCustomizeDiff is skipped for unit testing as it requires
// complex ResourceDiff mocking which is not straightforward with the Terraform SDK.
// The CustomizeDiff function is covered by acceptance tests.
// TODO: Add proper ResourceDiff mocking if unit test coverage is required
func TestResourceNodeCustomizeDiff(t *testing.T) {
	t.Skip("CustomizeDiff testing requires ResourceDiff mocking - covered by acceptance tests")
}
