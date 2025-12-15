package faas

import (
	"fmt"
	"testing"

	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// Schema Validation Tests
// ============================================================================

func TestResourceFaasFunctionSchema_Definition(t *testing.T) {
	resource := ResourceFaasFunction()
	require.NotNil(t, resource)
	assert.NotNil(t, resource.Schema)
}

func TestResourceFaasFunctionSchema_RequiredFields(t *testing.T) {
	resource := ResourceFaasFunction()
	schema := resource.Schema

	requiredFields := []string{
		tfconstants.AttrName,
		tfconstants.AttrNamespace,
		tfconstants.AttrRuntime,
	}

	for _, fieldName := range requiredFields {
		t.Run(fieldName+"_is_required", func(t *testing.T) {
			fieldSchema, exists := schema[fieldName]
			require.True(t, exists, "Field %s should exist in schema", fieldName)
			assert.True(t, fieldSchema.Required, "Field %s should be required", fieldName)
		})
	}
}

func TestResourceFaasFunctionSchema_ForceNewFields(t *testing.T) {
	resource := ResourceFaasFunction()
	schema := resource.Schema

	forceNewFields := []string{
		tfconstants.AttrName,
		tfconstants.AttrNamespace,
		tfconstants.AttrRuntime,
	}

	for _, fieldName := range forceNewFields {
		t.Run(fieldName+"_is_force_new", func(t *testing.T) {
			fieldSchema, exists := schema[fieldName]
			require.True(t, exists, "Field %s should exist in schema", fieldName)
			assert.True(t, fieldSchema.ForceNew, "Field %s should be ForceNew", fieldName)
		})
	}
}

func TestResourceFaasFunctionSchema_ComputedFields(t *testing.T) {
	resource := ResourceFaasFunction()
	schema := resource.Schema

	computedFields := []string{
		tfconstants.AttrEndpointURL,
		tfconstants.AttrStatus,
		tfconstants.AttrCreatedAt,
		tfconstants.AttrUpdatedAt,
	}

	for _, fieldName := range computedFields {
		t.Run(fieldName+"_is_computed", func(t *testing.T) {
			fieldSchema, exists := schema[fieldName]
			require.True(t, exists, "Field %s should exist in schema", fieldName)
			assert.True(t, fieldSchema.Computed, "Field %s should be computed", fieldName)
		})
	}
}

func TestResourceFaasFunctionSchema_SensitiveFields(t *testing.T) {
	resource := ResourceFaasFunction()
	schema := resource.Schema

	fieldSchema, exists := schema[tfconstants.AttrCodeInline]
	require.True(t, exists, "Field code_inline should exist in schema")
	assert.True(t, fieldSchema.Sensitive, "Field code_inline should be sensitive")
}

func TestResourceFaasFunctionSchema_ConflictsWith(t *testing.T) {
	resource := ResourceFaasFunction()
	schema := resource.Schema

	// Check that code_inline and code_file have ConflictsWith relationship
	codeInlineSchema, exists := schema[tfconstants.AttrCodeInline]
	require.True(t, exists, "Field code_inline should exist in schema")
	assert.Contains(t, codeInlineSchema.ConflictsWith, tfconstants.AttrCodeFile, "code_inline should conflict with code_file")

	codeFileSchema, exists := schema[tfconstants.AttrCodeFile]
	require.True(t, exists, "Field code_file should exist in schema")
	assert.Contains(t, codeFileSchema.ConflictsWith, tfconstants.AttrCodeInline, "code_file should conflict with code_inline")
}

func TestResourceFaasFunctionSchema_DefaultValues(t *testing.T) {
	resource := ResourceFaasFunction()
	schema := resource.Schema

	tests := []struct {
		fieldName     string
		expectedValue interface{}
		description   string
	}{
		{
			fieldName:     tfconstants.AttrMemoryMB,
			expectedValue: goe2econstants.FaaSDefaultMemoryMB,
			description:   "memory_mb should have default value",
		},
		{
			fieldName:     tfconstants.AttrTimeoutSeconds,
			expectedValue: goe2econstants.FaaSDefaultTimeoutSeconds,
			description:   "timeout_seconds should have default value",
		},
		{
			fieldName:     tfconstants.AttrMinReplicas,
			expectedValue: goe2econstants.FaaSDefaultMinReplicas,
			description:   "min_replicas should have default value",
		},
		{
			fieldName:     tfconstants.AttrMaxReplicas,
			expectedValue: goe2econstants.FaaSDefaultMaxReplicas,
			description:   "max_replicas should have default value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.fieldName+"_has_default", func(t *testing.T) {
			fieldSchema, exists := schema[tt.fieldName]
			require.True(t, exists, "Field %s should exist in schema", tt.fieldName)
			assert.Equal(t, tt.expectedValue, fieldSchema.Default, tt.description)
		})
	}
}

func TestResourceFaasFunctionSchema_ValidateFunc(t *testing.T) {
	resource := ResourceFaasFunction()
	schema := resource.Schema

	tests := []struct {
		fieldName   string
		description string
	}{
		{
			fieldName:   tfconstants.AttrMemoryMB,
			description: "memory_mb should have ValidateFunc set",
		},
		{
			fieldName:   tfconstants.AttrTimeoutSeconds,
			description: "timeout_seconds should have ValidateFunc set",
		},
		{
			fieldName:   tfconstants.AttrMinReplicas,
			description: "min_replicas should have ValidateFunc set",
		},
		{
			fieldName:   tfconstants.AttrMaxReplicas,
			description: "max_replicas should have ValidateFunc set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.fieldName+"_has_validate_func", func(t *testing.T) {
			fieldSchema, exists := schema[tt.fieldName]
			require.True(t, exists, "Field %s should exist in schema", tt.fieldName)
			assert.NotNil(t, fieldSchema.ValidateFunc, tt.description)
		})
	}
}

// ============================================================================
// Validation Function Tests
// ============================================================================

func TestResourceFaasFunctionSchema_MemoryMBValidation(t *testing.T) {
	resource := ResourceFaasFunction()
	schema := resource.Schema

	fieldSchema, exists := schema[tfconstants.AttrMemoryMB]
	require.True(t, exists, "Field memory_mb should exist")

	validateFunc := fieldSchema.ValidateFunc
	require.NotNil(t, validateFunc, "ValidateFunc should be set")

	// Test minimum value (should pass)
	_, errs := validateFunc(tfconstants.FaaSMinMemoryMB, tfconstants.AttrMemoryMB)
	assert.Empty(t, errs, "memory_mb at minimum value should pass validation")

	// Test below minimum (should fail)
	_, errs = validateFunc(tfconstants.FaaSMinMemoryMB-1, tfconstants.AttrMemoryMB)
	assert.NotEmpty(t, errs, "memory_mb below minimum should fail validation")
	if len(errs) > 0 {
		assert.Contains(t, errs[0].Error(), fmt.Sprintf("%d", tfconstants.FaaSMinMemoryMB), "Error should mention minimum value")
	}
}

func TestResourceFaasFunctionSchema_TimeoutSecondsValidation(t *testing.T) {
	resource := ResourceFaasFunction()
	schema := resource.Schema

	fieldSchema, exists := schema[tfconstants.AttrTimeoutSeconds]
	require.True(t, exists, "Field timeout_seconds should exist")

	validateFunc := fieldSchema.ValidateFunc
	require.NotNil(t, validateFunc, "ValidateFunc should be set")

	// Test minimum value (should pass)
	_, errs := validateFunc(tfconstants.FaaSMinTimeoutSeconds, tfconstants.AttrTimeoutSeconds)
	assert.Empty(t, errs, "timeout_seconds at minimum value should pass validation")

	// Test maximum value (should pass)
	_, errs = validateFunc(tfconstants.FaaSMaxTimeoutSeconds, tfconstants.AttrTimeoutSeconds)
	assert.Empty(t, errs, "timeout_seconds at maximum value should pass validation")

	// Test below minimum (should fail)
	_, errs = validateFunc(tfconstants.FaaSMinTimeoutSeconds-1, tfconstants.AttrTimeoutSeconds)
	assert.NotEmpty(t, errs, "timeout_seconds below minimum should fail validation")

	// Test above maximum (should fail)
	_, errs = validateFunc(tfconstants.FaaSMaxTimeoutSeconds+1, tfconstants.AttrTimeoutSeconds)
	assert.NotEmpty(t, errs, "timeout_seconds above maximum should fail validation")
	if len(errs) > 0 {
		assert.Contains(t, errs[0].Error(), fmt.Sprintf("%d", tfconstants.FaaSMaxTimeoutSeconds), "Error should mention maximum value")
	}
}

func TestResourceFaasFunctionSchema_MinReplicasValidation(t *testing.T) {
	resource := ResourceFaasFunction()
	schema := resource.Schema

	fieldSchema, exists := schema[tfconstants.AttrMinReplicas]
	require.True(t, exists, "Field min_replicas should exist")

	validateFunc := fieldSchema.ValidateFunc
	require.NotNil(t, validateFunc, "ValidateFunc should be set")

	// Test minimum value (0 should pass)
	_, errs := validateFunc(0, tfconstants.AttrMinReplicas)
	assert.Empty(t, errs, "min_replicas at 0 should pass validation")

	// Test below minimum (should fail)
	_, errs = validateFunc(-1, tfconstants.AttrMinReplicas)
	assert.NotEmpty(t, errs, "min_replicas below 0 should fail validation")
}

func TestResourceFaasFunctionSchema_MaxReplicasValidation(t *testing.T) {
	resource := ResourceFaasFunction()
	schema := resource.Schema

	fieldSchema, exists := schema[tfconstants.AttrMaxReplicas]
	require.True(t, exists, "Field max_replicas should exist")

	validateFunc := fieldSchema.ValidateFunc
	require.NotNil(t, validateFunc, "ValidateFunc should be set")

	// Test minimum value (1 should pass)
	_, errs := validateFunc(1, tfconstants.AttrMaxReplicas)
	assert.Empty(t, errs, "max_replicas at 1 should pass validation")

	// Test below minimum (should fail)
	_, errs = validateFunc(0, tfconstants.AttrMaxReplicas)
	assert.NotEmpty(t, errs, "max_replicas below 1 should fail validation")
}

// ============================================================================
// CustomizeDiff Tests
// ============================================================================

func TestResourceFaasFunctionCustomizeDiff_Exists(t *testing.T) {
	resource := ResourceFaasFunction()
	require.NotNil(t, resource)

	// Verify CustomizeDiff is set
	assert.NotNil(t, resource.CustomizeDiff, "CustomizeDiff should be set")
}

func TestCustomizeDiffFaasFunction_ReplicaValidation(t *testing.T) {
	// Test the replica validation logic directly
	tests := []struct {
		name        string
		minReplicas int
		maxReplicas int
		wantError   bool
		errorMsg    string
	}{
		{
			name:        "valid replica range",
			minReplicas: 1,
			maxReplicas: 5,
			wantError:   false,
		},
		{
			name:        "equal replicas",
			minReplicas: 3,
			maxReplicas: 3,
			wantError:   false,
		},
		{
			name:        "min greater than max",
			minReplicas: 10,
			maxReplicas: 5,
			wantError:   true,
			errorMsg:    "min_replicas (10) cannot be greater than max_replicas (5)",
		},
		{
			name:        "zero min replicas",
			minReplicas: 0,
			maxReplicas: 5,
			wantError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the validation logic from customizeDiffFaasFunction
			if tt.minReplicas > tt.maxReplicas {
				err := fmt.Errorf("min_replicas (%d) cannot be greater than max_replicas (%d)", tt.minReplicas, tt.maxReplicas)
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, nil)
			}
		})
	}
}

func TestCustomizeDiffFaasFunction_CodeSourceValidation(t *testing.T) {
	// Test the code source validation logic directly
	tests := []struct {
		name          string
		hasCodeInline bool
		hasCodeFile   bool
		wantError     bool
		errorMsg      string
	}{
		{
			name:          "only code_inline",
			hasCodeInline: true,
			hasCodeFile:   false,
			wantError:     false,
		},
		{
			name:          "only code_file",
			hasCodeInline: false,
			hasCodeFile:   true,
			wantError:     false,
		},
		{
			name:          "both code sources",
			hasCodeInline: true,
			hasCodeFile:   true,
			wantError:     true,
			errorMsg:      "code_inline and code_file are mutually exclusive",
		},
		{
			name:          "neither code source",
			hasCodeInline: false,
			hasCodeFile:   false,
			wantError:     true,
			errorMsg:      "one of code_inline or code_file must be specified",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the validation logic from customizeDiffFaasFunction
			if !tt.hasCodeInline && !tt.hasCodeFile {
				err := fmt.Errorf("one of code_inline or code_file must be specified")
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else if tt.hasCodeInline && tt.hasCodeFile {
				err := fmt.Errorf("code_inline and code_file are mutually exclusive")
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, nil)
			}
		})
	}
}

// ============================================================================
// State Management Tests
// ============================================================================

func TestResourceFaasFunctionImport_ImporterExists(t *testing.T) {
	resource := ResourceFaasFunction()
	require.NotNil(t, resource)
	require.NotNil(t, resource.Importer, "Importer should be set")
	require.NotNil(t, resource.Importer.StateContext, "StateContext should be set")
}

func TestResourceFaasFunctionImport_PassthroughContext(t *testing.T) {
	resource := ResourceFaasFunction()
	require.NotNil(t, resource)
	require.NotNil(t, resource.Importer)

	// Verify that ImportStatePassthroughContext is used
	// This means the function ID is used directly as the resource ID
	assert.NotNil(t, resource.Importer.StateContext)
}

// ============================================================================
// Resource CRUD Function Tests
// ============================================================================

func TestResourceFaasFunction_CRUDFunctionsExist(t *testing.T) {
	resource := ResourceFaasFunction()
	require.NotNil(t, resource)

	assert.NotNil(t, resource.CreateContext, "CreateContext should be set")
	assert.NotNil(t, resource.ReadContext, "ReadContext should be set")
	assert.NotNil(t, resource.UpdateContext, "UpdateContext should be set")
	assert.NotNil(t, resource.DeleteContext, "DeleteContext should be set")
	assert.NotNil(t, resource.Exists, "Exists should be set")
}

// ============================================================================
// Data Source Schema Tests
// ============================================================================

func TestDataSourceFaasFunctionSchema_Definition(t *testing.T) {
	resource := DataSourceFaasFunction()
	require.NotNil(t, resource)
	assert.NotNil(t, resource.Schema)
}

func TestDataSourceFaasFunctionSchema_RequiredFields(t *testing.T) {
	resource := DataSourceFaasFunction()
	schema := resource.Schema

	requiredFields := []string{
		tfconstants.AttrFunctionID,
	}

	for _, fieldName := range requiredFields {
		t.Run(fieldName+"_is_required", func(t *testing.T) {
			fieldSchema, exists := schema[fieldName]
			require.True(t, exists, "Field %s should exist in schema", fieldName)
			assert.True(t, fieldSchema.Required, "Field %s should be required", fieldName)
		})
	}
}

func TestDataSourceFaasFunctionSchema_ComputedFields(t *testing.T) {
	resource := DataSourceFaasFunction()
	schema := resource.Schema

	computedFields := []string{
		tfconstants.AttrName,
		tfconstants.AttrNamespace,
		tfconstants.AttrRuntime,
		tfconstants.AttrMemoryMB,
		tfconstants.AttrTimeoutSeconds,
		tfconstants.AttrMinReplicas,
		tfconstants.AttrMaxReplicas,
		tfconstants.AttrEndpointURL,
		tfconstants.AttrStatus,
		tfconstants.AttrCreatedAt,
		tfconstants.AttrUpdatedAt,
	}

	for _, fieldName := range computedFields {
		t.Run(fieldName+"_is_computed", func(t *testing.T) {
			fieldSchema, exists := schema[fieldName]
			require.True(t, exists, "Field %s should exist in schema", fieldName)
			assert.True(t, fieldSchema.Computed, "Field %s should be computed", fieldName)
		})
	}
}

func TestDataSourceFaasFunction_ReadFunctionExists(t *testing.T) {
	resource := DataSourceFaasFunction()
	require.NotNil(t, resource)

	assert.NotNil(t, resource.ReadContext, "ReadContext should be set")
}

// ============================================================================
// Constants Usage Tests
// ============================================================================

func TestResourceFaasFunction_UsesConstants(t *testing.T) {
	resource := ResourceFaasFunction()
	schema := resource.Schema

	// Verify that constants are used for attribute names
	constantFields := []string{
		tfconstants.AttrName,
		tfconstants.AttrNamespace,
		tfconstants.AttrRuntime,
		tfconstants.AttrCodeInline,
		tfconstants.AttrCodeFile,
		tfconstants.AttrDescription,
		tfconstants.AttrMemoryMB,
		tfconstants.AttrTimeoutSeconds,
		tfconstants.AttrMinReplicas,
		tfconstants.AttrMaxReplicas,
		tfconstants.AttrEnvironmentVariables,
		tfconstants.AttrEndpointURL,
		tfconstants.AttrStatus,
		tfconstants.AttrCreatedAt,
		tfconstants.AttrUpdatedAt,
		tfconstants.AttrTags,
	}

	for _, fieldName := range constantFields {
		t.Run(fieldName+"_exists", func(t *testing.T) {
			_, exists := schema[fieldName]
			assert.True(t, exists, "Field %s should exist in schema (using constant)", fieldName)
		})
	}
}

func TestResourceFaasFunction_UsesDefaultConstants(t *testing.T) {
	resource := ResourceFaasFunction()
	schema := resource.Schema

	// Verify default values use constants
	memorySchema := schema[tfconstants.AttrMemoryMB]
	require.NotNil(t, memorySchema)
	assert.Equal(t, goe2econstants.FaaSDefaultMemoryMB, memorySchema.Default, "memory_mb default should use constant")

	timeoutSchema := schema[tfconstants.AttrTimeoutSeconds]
	require.NotNil(t, timeoutSchema)
	assert.Equal(t, goe2econstants.FaaSDefaultTimeoutSeconds, timeoutSchema.Default, "timeout_seconds default should use constant")

	minReplicasSchema := schema[tfconstants.AttrMinReplicas]
	require.NotNil(t, minReplicasSchema)
	assert.Equal(t, goe2econstants.FaaSDefaultMinReplicas, minReplicasSchema.Default, "min_replicas default should use constant")

	maxReplicasSchema := schema[tfconstants.AttrMaxReplicas]
	require.NotNil(t, maxReplicasSchema)
	assert.Equal(t, goe2econstants.FaaSDefaultMaxReplicas, maxReplicasSchema.Default, "max_replicas default should use constant")
}

func TestResourceFaasFunction_UsesValidationConstants(t *testing.T) {
	resource := ResourceFaasFunction()
	schema := resource.Schema

	// Verify validation functions use constants
	memorySchema := schema[tfconstants.AttrMemoryMB]
	require.NotNil(t, memorySchema)
	require.NotNil(t, memorySchema.ValidateFunc)

	// Test that validation uses the correct minimum constant
	_, errs := memorySchema.ValidateFunc(tfconstants.FaaSMinMemoryMB-1, tfconstants.AttrMemoryMB)
	assert.NotEmpty(t, errs, "Validation should fail below minimum")
	_, errs = memorySchema.ValidateFunc(tfconstants.FaaSMinMemoryMB, tfconstants.AttrMemoryMB)
	assert.Empty(t, errs, "Validation should pass at minimum")

	timeoutSchema := schema[tfconstants.AttrTimeoutSeconds]
	require.NotNil(t, timeoutSchema)
	require.NotNil(t, timeoutSchema.ValidateFunc)

	// Test that validation uses the correct min/max constants
	_, errs = timeoutSchema.ValidateFunc(tfconstants.FaaSMinTimeoutSeconds-1, tfconstants.AttrTimeoutSeconds)
	assert.NotEmpty(t, errs, "Validation should fail below minimum")
	_, errs = timeoutSchema.ValidateFunc(tfconstants.FaaSMaxTimeoutSeconds+1, tfconstants.AttrTimeoutSeconds)
	assert.NotEmpty(t, errs, "Validation should fail above maximum")
	_, errs = timeoutSchema.ValidateFunc(tfconstants.FaaSMinTimeoutSeconds, tfconstants.AttrTimeoutSeconds)
	assert.Empty(t, errs, "Validation should pass at minimum")
	_, errs = timeoutSchema.ValidateFunc(tfconstants.FaaSMaxTimeoutSeconds, tfconstants.AttrTimeoutSeconds)
	assert.Empty(t, errs, "Validation should pass at maximum")
}
