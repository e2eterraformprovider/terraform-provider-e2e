package ssh_key

import (
	"fmt"
	"strings"
	"testing"

	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// parseImportID is a helper function that extracts the parsing logic for testing
// This mirrors the logic in resourceSshKeyImport but without API calls
func parseImportID(importID string, defaultRegion string) (projectID, region, sshKeyID string, err error) {
	if importID == "" {
		return "", "", "", fmt.Errorf(tfconstants.SSHKeyImportIDInvalidFormat)
	}

	parts := strings.Split(importID, ":")

	if len(parts) == 2 {
		projectID = parts[0]
		sshKeyID = parts[1]
		region = defaultRegion
		if region == "" {
			return "", "", "", fmt.Errorf(tfconstants.SSHKeyImportIDRegionRequired)
		}
	} else if len(parts) == 3 {
		projectID = parts[0]
		region = parts[1]
		sshKeyID = parts[2]
	} else {
		return "", "", "", fmt.Errorf(tfconstants.SSHKeyImportIDInvalidFormat)
	}

	return projectID, region, sshKeyID, nil
}

// TestResourceSshKeyImportIDParsing tests the import ID parsing logic
// This tests the parsing without making API calls
func TestResourceSshKeyImportIDParsing(t *testing.T) {
	tests := []struct {
		name          string
		importID      string
		defaultRegion string
		expectError   bool
		expectedParts map[string]string // projectID, region, sshKeyID
		errorContains string
	}{
		{
			name:          "2-part_format_with_default_region",
			importID:      "proj-123:key-456",
			defaultRegion: "Mumbai",
			expectError:   false,
			expectedParts: map[string]string{
				"projectID": "proj-123",
				"region":    "Mumbai",
				"sshKeyID":  "key-456",
			},
		},
		{
			name:          "3-part_format",
			importID:      "proj-123:us-east-1:key-456",
			defaultRegion: "",
			expectError:   false,
			expectedParts: map[string]string{
				"projectID": "proj-123",
				"region":    "us-east-1",
				"sshKeyID":  "key-456",
			},
		},
		{
			name:          "2-part_format_no_default_region",
			importID:      "proj-123:key-456",
			defaultRegion: "",
			expectError:   true,
			errorContains: "region must be specified",
		},
		{
			name:          "invalid_format_single_part",
			importID:      "proj-123",
			defaultRegion: "Mumbai",
			expectError:   true,
			errorContains: "invalid import ID format",
		},
		{
			name:          "invalid_format_four_parts",
			importID:      "proj-123:region:key-456:extra",
			defaultRegion: "Mumbai",
			expectError:   true,
			errorContains: "invalid import ID format",
		},
		{
			name:          "empty_import_id",
			importID:      "",
			defaultRegion: "Mumbai",
			expectError:   true,
			errorContains: "invalid import ID format",
		},
		{
			name:          "3-part_format_with_empty_region",
			importID:      "proj-123::key-456",
			defaultRegion: "",
			expectError:   false,
			expectedParts: map[string]string{
				"projectID": "proj-123",
				"region":    "",
				"sshKeyID":  "key-456",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectID, region, sshKeyID, err := parseImportID(tt.importID, tt.defaultRegion)

			if tt.expectError {
				require.Error(t, err, "Expected error but parsing succeeded")
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains,
						"Error message should contain '%s'", tt.errorContains)
				}
			} else {
				require.NoError(t, err, "Unexpected error during parsing")
				assert.Equal(t, tt.expectedParts["projectID"], projectID)
				assert.Equal(t, tt.expectedParts["region"], region)
				assert.Equal(t, tt.expectedParts["sshKeyID"], sshKeyID)
			}
		})
	}
}

// TestResourceSshKeySchemaConflicts tests that ConflictsWith constraints are properly configured
func TestResourceSshKeySchemaConflicts(t *testing.T) {
	resource := ResourceSshKey()
	require.NotNil(t, resource)
	resourceSchema := resource.Schema

	// Verify ConflictsWith constraints are set correctly
	tests := []struct {
		field         string
		conflictsWith []string
	}{
		{
			field:         tfconstants.AttrName,
			conflictsWith: []string{tfconstants.AttrLabel},
		},
		{
			field:         tfconstants.AttrLabel,
			conflictsWith: []string{tfconstants.AttrName},
		},
		{
			field:         tfconstants.AttrPublicKey,
			conflictsWith: []string{tfconstants.AttrSSHKey},
		},
		{
			field:         tfconstants.AttrSSHKey,
			conflictsWith: []string{tfconstants.AttrPublicKey},
		},
		{
			field:         tfconstants.AttrRegion,
			conflictsWith: []string{tfconstants.AttrLocation},
		},
		{
			field:         tfconstants.AttrLocation,
			conflictsWith: []string{tfconstants.AttrRegion},
		},
	}

	for _, tt := range tests {
		t.Run(tt.field+"_conflicts_with", func(t *testing.T) {
			fieldSchema, exists := resourceSchema[tt.field]
			require.True(t, exists, "Field %s should exist in schema", tt.field)
			assert.Equal(t, tt.conflictsWith, fieldSchema.ConflictsWith,
				"Field %s should conflict with %v", tt.field, tt.conflictsWith)
		})
	}
}

// TestResourceSshKeySchemaExactlyOneOf tests that ExactlyOneOf constraints are properly configured
func TestResourceSshKeySchemaExactlyOneOf(t *testing.T) {
	resource := ResourceSshKey()
	require.NotNil(t, resource)
	resourceSchema := resource.Schema

	// Verify ExactlyOneOf constraints are set correctly
	tests := []struct {
		field        string
		exactlyOneOf []string
	}{
		{
			field:        tfconstants.AttrName,
			exactlyOneOf: []string{tfconstants.AttrName, tfconstants.AttrLabel},
		},
		{
			field:        tfconstants.AttrLabel,
			exactlyOneOf: []string{tfconstants.AttrName, tfconstants.AttrLabel},
		},
		{
			field:        tfconstants.AttrPublicKey,
			exactlyOneOf: []string{tfconstants.AttrPublicKey, tfconstants.AttrSSHKey},
		},
		{
			field:        tfconstants.AttrSSHKey,
			exactlyOneOf: []string{tfconstants.AttrPublicKey, tfconstants.AttrSSHKey},
		},
	}

	for _, tt := range tests {
		t.Run(tt.field+"_exactly_one_of", func(t *testing.T) {
			fieldSchema, exists := resourceSchema[tt.field]
			require.True(t, exists, "Field %s should exist in schema", tt.field)
			assert.Equal(t, tt.exactlyOneOf, fieldSchema.ExactlyOneOf,
				"Field %s should be exactly one of %v", tt.field, tt.exactlyOneOf)
		})
	}
}

// TestResourceSshKeySchemaValidation tests that ValidateFunc constraints are properly configured
func TestResourceSshKeySchemaValidation(t *testing.T) {
	resource := ResourceSshKey()
	require.NotNil(t, resource)
	resourceSchema := resource.Schema

	// Verify ValidateFunc is set for fields that require validation
	fieldsWithValidation := []string{
		tfconstants.AttrName,
		tfconstants.AttrLabel,
		tfconstants.AttrPublicKey,
		tfconstants.AttrSSHKey,
		tfconstants.AttrRegion,
	}

	for _, fieldName := range fieldsWithValidation {
		t.Run(fieldName+"_has_validation", func(t *testing.T) {
			fieldSchema, exists := resourceSchema[fieldName]
			require.True(t, exists, "Field %s should exist in schema", fieldName)
			assert.NotNil(t, fieldSchema.ValidateFunc,
				"Field %s should have ValidateFunc configured", fieldName)
		})
	}
}

// TestResourceSshKeyImportIDFormatErrorMessages tests error message constants
func TestResourceSshKeyImportIDFormatErrorMessages(t *testing.T) {
	// Test that error messages are properly formatted
	assert.NotEmpty(t, tfconstants.SSHKeyImportIDFormatDescription)
	assert.NotEmpty(t, tfconstants.SSHKeyImportIDRegionRequired)
	assert.NotEmpty(t, tfconstants.SSHKeyImportIDInvalidFormat)

	// Verify error messages contain helpful information
	assert.Contains(t, tfconstants.SSHKeyImportIDFormatDescription, "project_id")
	assert.Contains(t, tfconstants.SSHKeyImportIDRegionRequired, "region")
	assert.Contains(t, tfconstants.SSHKeyImportIDInvalidFormat, "invalid")
}
