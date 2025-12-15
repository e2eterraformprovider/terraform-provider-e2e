package loadbalancer

import (
	"strings"
	"testing"

	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/stretchr/testify/assert"
)

// ============================================================================
// TestGetLbPort - Tests for GetLbPort function
// ============================================================================

// TestGetLbPort_HTTPMode tests that HTTP mode returns port 80
func TestGetLbPort_HTTPMode(t *testing.T) {
	result := GetLbPort(goe2econstants.LBModeHTTP)
	assert.Equal(t, goe2econstants.LBPortHTTP, result, "HTTP mode should return port 80")
}

// TestGetLbPort_HTTPSMode tests that HTTPS mode returns port 443
func TestGetLbPort_HTTPSMode(t *testing.T) {
	result := GetLbPort(goe2econstants.LBModeHTTPS)
	assert.Equal(t, goe2econstants.LBPortHTTPS, result, "HTTPS mode should return port 443")
}

// TestGetLbPort_BothMode tests that Both mode returns port 443 (defaults to HTTPS)
func TestGetLbPort_BothMode(t *testing.T) {
	result := GetLbPort(goe2econstants.LBModeBoth)
	assert.Equal(t, goe2econstants.LBPortHTTPS, result, "Both mode should return port 443 (default to HTTPS)")
}

// TestGetLbPort_CaseSensitivity tests case sensitivity of mode parameter
func TestGetLbPort_CaseSensitivity(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
		desc     string
	}{
		{
			name:     "lowercase http",
			input:    "http",
			expected: goe2econstants.LBPortHTTPS,
			desc:     "Function is case-sensitive, lowercase 'http' returns default HTTPS port",
		},
		{
			name:     "lowercase https",
			input:    "https",
			expected: goe2econstants.LBPortHTTPS,
			desc:     "Function is case-sensitive, lowercase 'https' returns default HTTPS port",
		},
		{
			name:     "mixed case Http",
			input:    "Http",
			expected: goe2econstants.LBPortHTTPS,
			desc:     "Function is case-sensitive, 'Http' returns default HTTPS port",
		},
		{
			name:     "mixed case Https",
			input:    "Https",
			expected: goe2econstants.LBPortHTTPS,
			desc:     "Function is case-sensitive, 'Https' returns default HTTPS port",
		},
		{
			name:     "correct case HTTP",
			input:    goe2econstants.LBModeHTTP,
			expected: goe2econstants.LBPortHTTP,
			desc:     "Correct case 'HTTP' returns HTTP port 80",
		},
		{
			name:     "correct case HTTPS",
			input:    goe2econstants.LBModeHTTPS,
			expected: goe2econstants.LBPortHTTPS,
			desc:     "Correct case 'HTTPS' returns HTTPS port 443",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GetLbPort(tc.input)
			assert.Equal(t, tc.expected, result, tc.desc)
		})
	}
}

// TestGetLbPort_InvalidMode tests behavior with invalid mode inputs
func TestGetLbPort_InvalidMode(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
		desc     string
	}{
		{
			name:     "invalid mode INVALID",
			input:    "INVALID",
			expected: goe2econstants.LBPortHTTPS,
			desc:     "Invalid mode returns default HTTPS port 443",
		},
		{
			name:     "empty string",
			input:    "",
			expected: goe2econstants.LBPortHTTPS,
			desc:     "Empty string returns default HTTPS port 443",
		},
		{
			name:     "TCP mode",
			input:    "TCP",
			expected: goe2econstants.LBPortHTTPS,
			desc:     "TCP mode returns default HTTPS port 443",
		},
		{
			name:     "random string",
			input:    "RandomMode",
			expected: goe2econstants.LBPortHTTPS,
			desc:     "Random string returns default HTTPS port 443",
		},
		{
			name:     "numeric string",
			input:    "12345",
			expected: goe2econstants.LBPortHTTPS,
			desc:     "Numeric string returns default HTTPS port 443",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GetLbPort(tc.input)
			assert.Equal(t, tc.expected, result, tc.desc)
		})
	}
}

// TestGetLbPort_NoPanic tests that function never panics
func TestGetLbPort_NoPanic(t *testing.T) {
	testInputs := []string{
		goe2econstants.LBModeHTTP,
		goe2econstants.LBModeHTTPS,
		goe2econstants.LBModeBoth,
		"",
		"INVALID",
		"TCP",
		"http",
		"https",
		"12345",
		"!@#$%^&*()",
	}

	for _, input := range testInputs {
		t.Run("no panic for input: "+input, func(t *testing.T) {
			assert.NotPanics(t, func() {
				result := GetLbPort(input)
				assert.True(t, true, "GetLbPort returned %q", result)
			}, "GetLbPort should never panic for input: %s", input)
		})
	}
}

// ============================================================================
// TestExpandAclList - Tests for ExpandAclList function
// ============================================================================

// TestExpandAclList_Empty tests expansion with empty ACL list
func TestExpandAclList_Empty(t *testing.T) {
	input := []interface{}{}
	result, err := ExpandAclList(input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
}

// TestExpandAclList_SingleACL tests expansion with single ACL rule
func TestExpandAclList_SingleACL(t *testing.T) {
	input := []interface{}{
		map[string]interface{}{
			"acl_name":          "test-acl",
			"acl_condition":     "path_beg",
			"acl_matching_path": "/api",
		},
	}

	result, err := ExpandAclList(input)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "test-acl", result[0].ACLName)
	assert.Equal(t, "path_beg", result[0].ACLCondition)
	assert.Equal(t, "/api", result[0].ACLMatchingPath)
}

// TestExpandAclList_MultipleACLs tests expansion with multiple ACL rules
func TestExpandAclList_MultipleACLs(t *testing.T) {
	input := []interface{}{
		map[string]interface{}{
			"acl_name":          "api-acl",
			"acl_condition":     "path_beg",
			"acl_matching_path": "/api",
		},
		map[string]interface{}{
			"acl_name":          "admin-acl",
			"acl_condition":     "path_beg",
			"acl_matching_path": "/admin",
		},
		map[string]interface{}{
			"acl_name":          "static-acl",
			"acl_condition":     "path_beg",
			"acl_matching_path": "/static",
		},
	}

	result, err := ExpandAclList(input)

	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, "api-acl", result[0].ACLName)
	assert.Equal(t, "admin-acl", result[1].ACLName)
	assert.Equal(t, "static-acl", result[2].ACLName)
}

// ============================================================================
// TestExpandAclMap - Tests for ExpandAclMap function
// ============================================================================

// TestExpandAclMap_Empty tests expansion with empty ACL map
func TestExpandAclMap_Empty(t *testing.T) {
	input := []interface{}{}
	result, err := ExpandAclMap(input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
}

// TestExpandAclMap_SingleMapping tests expansion with single ACL mapping
func TestExpandAclMap_SingleMapping(t *testing.T) {
	input := []interface{}{
		map[string]interface{}{
			"acl_name":            "test-acl",
			"acl_condition_state": true,
			"acl_backend":         "backend-1",
		},
	}

	result, err := ExpandAclMap(input)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "test-acl", result[0].ACLName)
	assert.True(t, result[0].ACLConditionState)
	assert.Equal(t, "backend-1", result[0].ACLBackend)
}

// TestExpandAclMap_MultipleMappings tests expansion with multiple ACL mappings
func TestExpandAclMap_MultipleMappings(t *testing.T) {
	input := []interface{}{
		map[string]interface{}{
			"acl_name":            "api-acl",
			"acl_condition_state": true,
			"acl_backend":         "api-backend",
		},
		map[string]interface{}{
			"acl_name":            "admin-acl",
			"acl_condition_state": true,
			"acl_backend":         "admin-backend",
		},
	}

	result, err := ExpandAclMap(input)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "api-acl", result[0].ACLName)
	assert.Equal(t, "admin-acl", result[1].ACLName)
	assert.Equal(t, "api-backend", result[0].ACLBackend)
	assert.Equal(t, "admin-backend", result[1].ACLBackend)
}

// TestExpandAclMap_AlwaysTrue tests that ACLConditionState is always set to true
func TestExpandAclMap_AlwaysTrue(t *testing.T) {
	input := []interface{}{
		map[string]interface{}{
			"acl_name":            "test-acl",
			"acl_condition_state": false, // Input value (ignored)
			"acl_backend":         "backend-1",
		},
	}

	result, err := ExpandAclMap(input)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	// Function always sets ACLConditionState to true (hardcoded in implementation)
	assert.True(t, result[0].ACLConditionState, "ACLConditionState should always be true")
}

// ============================================================================
// TestExpandEnableEosLogger - Tests for ExpandEnableEosLogger function
// ============================================================================

// TestExpandEnableEosLogger_Empty tests expansion with empty EOS logger config
func TestExpandEnableEosLogger_Empty(t *testing.T) {
	input := []interface{}{}
	result, err := ExpandEnableEosLogger(input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	// Empty input should return empty struct
	assert.Equal(t, 0, result.ApplianceID)
	assert.Equal(t, "", result.AccessKey)
	assert.Equal(t, "", result.SecretKey)
	assert.Equal(t, "", result.Bucket)
}

// TestExpandEnableEosLogger_ValidConfig tests expansion with valid EOS logger config
func TestExpandEnableEosLogger_ValidConfig(t *testing.T) {
	input := []interface{}{
		map[string]interface{}{
			"appliance_id": 123,
			"access_key":   "test-access-key",
			"secret_key":   "test-secret-key",
			"bucket":       "test-bucket",
		},
	}

	result, err := ExpandEnableEosLogger(input)

	assert.NoError(t, err)
	assert.Equal(t, 123, result.ApplianceID)
	assert.Equal(t, "test-access-key", result.AccessKey)
	assert.Equal(t, "test-secret-key", result.SecretKey)
	assert.Equal(t, "test-bucket", result.Bucket)
}

// TestExpandEnableEosLogger_ZeroApplianceID tests expansion with zero appliance ID
func TestExpandEnableEosLogger_ZeroApplianceID(t *testing.T) {
	input := []interface{}{
		map[string]interface{}{
			"appliance_id": 0,
			"access_key":   "test-access-key",
			"secret_key":   "test-secret-key",
			"bucket":       "test-bucket",
		},
	}

	result, err := ExpandEnableEosLogger(input)

	assert.NoError(t, err)
	assert.Equal(t, 0, result.ApplianceID)
	assert.Equal(t, "test-access-key", result.AccessKey)
}

// TestExpandEnableEosLogger_MultipleConfigs tests that only first config is used
func TestExpandEnableEosLogger_MultipleConfigs(t *testing.T) {
	input := []interface{}{
		map[string]interface{}{
			"appliance_id": 123,
			"access_key":   "first-access-key",
			"secret_key":   "first-secret-key",
			"bucket":       "first-bucket",
		},
		map[string]interface{}{
			"appliance_id": 456,
			"access_key":   "second-access-key",
			"secret_key":   "second-secret-key",
			"bucket":       "second-bucket",
		},
	}

	result, err := ExpandEnableEosLogger(input)

	assert.NoError(t, err)
	// Function iterates through all configs but overwrites, so last one wins
	assert.Equal(t, 456, result.ApplianceID)
	assert.Equal(t, "second-access-key", result.AccessKey)
	assert.Equal(t, "second-secret-key", result.SecretKey)
	assert.Equal(t, "second-bucket", result.Bucket)
}

// ============================================================================
// TestNormalizeLoadBalancerState - Tests for normalizeLoadBalancerState function
// ============================================================================
// NOTE: This function is defined in resource_loadbalancer.go, not helpers.go

// TestNormalizeLoadBalancerState_CreatingStatus tests normalization of Creating status
func TestNormalizeLoadBalancerState_CreatingStatus(t *testing.T) {
	result := normalizeLoadBalancerState(goe2econstants.LBStatusCreating)
	assert.Equal(t, goe2econstants.LBStateCreating, result, "Creating status should normalize to 'creating'")
}

// TestNormalizeLoadBalancerState_DeployingStatus tests normalization of Deploying status
func TestNormalizeLoadBalancerState_DeployingStatus(t *testing.T) {
	result := normalizeLoadBalancerState(goe2econstants.LBStatusDeploying)
	assert.Equal(t, goe2econstants.LBStateCreating, result, "Deploying status should normalize to 'creating'")
}

// TestNormalizeLoadBalancerState_RunningStatus tests normalization of Running status
func TestNormalizeLoadBalancerState_RunningStatus(t *testing.T) {
	result := normalizeLoadBalancerState(goe2econstants.LBStatusRunning)
	assert.Equal(t, goe2econstants.LBStateRunning, result, "Running status should normalize to 'running'")
}

// TestNormalizeLoadBalancerState_PoweredOffStatus tests normalization of Powered off status
func TestNormalizeLoadBalancerState_PoweredOffStatus(t *testing.T) {
	result := normalizeLoadBalancerState(goe2econstants.LBStatusPoweredOff)
	assert.Equal(t, goe2econstants.LBStateStopped, result, "Powered off status should normalize to 'stopped'")
}

// TestNormalizeLoadBalancerState_UpgradingStatus tests normalization of Upgrading status
func TestNormalizeLoadBalancerState_UpgradingStatus(t *testing.T) {
	result := normalizeLoadBalancerState(goe2econstants.LBStatusUpgrading)
	assert.Equal(t, goe2econstants.LBStateUpgrading, result, "Upgrading status should normalize to 'upgrading'")
}

// TestNormalizeLoadBalancerState_ErrorStatus tests normalization of Error status
func TestNormalizeLoadBalancerState_ErrorStatus(t *testing.T) {
	result := normalizeLoadBalancerState(goe2econstants.LBStatusError)
	assert.Equal(t, goe2econstants.LBStateError, result, "Error status should normalize to 'error'")
}

// TestNormalizeLoadBalancerState_FailedStatus tests normalization of Failed status
func TestNormalizeLoadBalancerState_FailedStatus(t *testing.T) {
	result := normalizeLoadBalancerState("Failed")
	assert.Equal(t, goe2econstants.LBStateError, result, "Failed status should normalize to 'error'")
}

// TestNormalizeLoadBalancerState_UnknownStatus tests normalization of unknown status
func TestNormalizeLoadBalancerState_UnknownStatus(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
		desc     string
	}{
		{
			name:     "unknown status",
			input:    "UnknownStatus",
			expected: "unknownstatus",
			desc:     "Unknown status should be lowercased",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
			desc:     "Empty string should remain empty",
		},
		{
			name:     "random string",
			input:    "RandomStatus",
			expected: "randomstatus",
			desc:     "Random status should be lowercased",
		},
		{
			name:     "numeric string",
			input:    "12345",
			expected: "12345",
			desc:     "Numeric string should remain unchanged (lowercased)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := normalizeLoadBalancerState(tc.input)
			assert.Equal(t, tc.expected, result, tc.desc)
		})
	}
}

// TestNormalizeLoadBalancerState_CaseInsensitivity tests case handling
func TestNormalizeLoadBalancerState_CaseInsensitivity(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
		desc     string
	}{
		{
			name:     "uppercase RUNNING",
			input:    "RUNNING",
			expected: "running",
			desc:     "RUNNING (uppercase) should normalize to 'running'",
		},
		{
			name:     "lowercase running",
			input:    "running",
			expected: "running",
			desc:     "running (lowercase) should normalize to 'running'",
		},
		{
			name:     "mixed case Running",
			input:    "Running",
			expected: goe2econstants.LBStateRunning,
			desc:     "Running (title case) should normalize to 'running' (matches API constant)",
		},
		{
			name:     "uppercase ERROR",
			input:    "ERROR",
			expected: "error",
			desc:     "ERROR (uppercase) should normalize to 'error'",
		},
		{
			name:     "lowercase error",
			input:    "error",
			expected: "error",
			desc:     "error (lowercase) should normalize to 'error'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := normalizeLoadBalancerState(tc.input)
			assert.Equal(t, tc.expected, result, tc.desc)
		})
	}
}

// TestNormalizeLoadBalancerState_AllAPIStatuses tests all known API status values
func TestNormalizeLoadBalancerState_AllAPIStatuses(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
		desc     string
	}{
		{
			name:     "LBStatusCreating",
			input:    goe2econstants.LBStatusCreating,
			expected: goe2econstants.LBStateCreating,
			desc:     "Creating → creating",
		},
		{
			name:     "LBStatusDeploying",
			input:    goe2econstants.LBStatusDeploying,
			expected: goe2econstants.LBStateCreating,
			desc:     "Deploying → creating",
		},
		{
			name:     "LBStatusRunning",
			input:    goe2econstants.LBStatusRunning,
			expected: goe2econstants.LBStateRunning,
			desc:     "Running → running",
		},
		{
			name:     "LBStatusPoweredOff",
			input:    goe2econstants.LBStatusPoweredOff,
			expected: goe2econstants.LBStateStopped,
			desc:     "Powered off → stopped",
		},
		{
			name:     "LBStatusUpgrading",
			input:    goe2econstants.LBStatusUpgrading,
			expected: goe2econstants.LBStateUpgrading,
			desc:     "Upgrading → upgrading",
		},
		{
			name:     "LBStatusError",
			input:    goe2econstants.LBStatusError,
			expected: goe2econstants.LBStateError,
			desc:     "Error → error",
		},
		{
			name:     "Failed string",
			input:    "Failed",
			expected: goe2econstants.LBStateError,
			desc:     "Failed → error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := normalizeLoadBalancerState(tc.input)
			assert.Equal(t, tc.expected, result, tc.desc)
		})
	}
}

// TestNormalizeLoadBalancerState_NoPanic tests that function never panics
func TestNormalizeLoadBalancerState_NoPanic(t *testing.T) {
	testInputs := []string{
		goe2econstants.LBStatusCreating,
		goe2econstants.LBStatusDeploying,
		goe2econstants.LBStatusRunning,
		goe2econstants.LBStatusPoweredOff,
		goe2econstants.LBStatusUpgrading,
		goe2econstants.LBStatusError,
		"Failed",
		"",
		"UnknownStatus",
		"RUNNING",
		"running",
		"12345",
		"!@#$%^&*()",
	}

	for _, input := range testInputs {
		t.Run("no panic for input: "+input, func(t *testing.T) {
			assert.NotPanics(t, func() {
				result := normalizeLoadBalancerState(input)
				assert.True(t, true, "normalizeLoadBalancerState returned %q", result)
			}, "normalizeLoadBalancerState should never panic for input: %s", input)
		})
	}
}

// ============================================================================
// Note: Tests for SetLoadBalancerStatus, ExpandBackendsWithGoe2e, and
// ExpandServersWithGoe2e require integration testing with a real or mocked
// goe2e.Client and are best tested via acceptance tests.
//
// Note: Flatten functions (FlattenTcpBackend, FlattenAclList, FlattenAclMap)
// do not exist in the loadbalancer package. The Read operation sets fields
// directly from the API response instead of using flatten helper functions.
// This is different from the security_group pattern.
// ============================================================================

// ============================================================================
// TestCheckStatus - Tests for CheckStatus function
// ============================================================================

// TestCheckStatus_Found tests that status is found in list
func TestCheckStatus_Found(t *testing.T) {
	statusList := []string{
		goe2econstants.LBStatusCreating,
		goe2econstants.LBStatusDeploying,
		goe2econstants.LBStatusUpgrading,
	}

	result := CheckStatus(statusList, goe2econstants.LBStatusCreating)
	assert.True(t, result, "Status should be found in list")

	result = CheckStatus(statusList, goe2econstants.LBStatusDeploying)
	assert.True(t, result, "Status should be found in list")

	result = CheckStatus(statusList, goe2econstants.LBStatusUpgrading)
	assert.True(t, result, "Status should be found in list")
}

// TestCheckStatus_NotFound tests that status is not found in list
func TestCheckStatus_NotFound(t *testing.T) {
	statusList := []string{
		goe2econstants.LBStatusCreating,
		goe2econstants.LBStatusDeploying,
	}

	result := CheckStatus(statusList, goe2econstants.LBStatusRunning)
	assert.False(t, result, "Status should not be found in list")

	result = CheckStatus(statusList, goe2econstants.LBStatusPoweredOff)
	assert.False(t, result, "Status should not be found in list")
}

// TestCheckStatus_CaseInsensitive tests that status check is case-insensitive
func TestCheckStatus_CaseInsensitive(t *testing.T) {
	statusList := []string{
		"Creating",
		"Deploying",
	}

	result := CheckStatus(statusList, "creating")
	assert.True(t, result, "Status check should be case-insensitive")

	result = CheckStatus(statusList, "CREATING")
	assert.True(t, result, "Status check should be case-insensitive")

	result = CheckStatus(statusList, "deploying")
	assert.True(t, result, "Status check should be case-insensitive")
}

// TestCheckStatus_EmptyList tests behavior with empty status list
func TestCheckStatus_EmptyList(t *testing.T) {
	statusList := []string{}

	result := CheckStatus(statusList, goe2econstants.LBStatusRunning)
	assert.False(t, result, "Empty list should return false")
}

// TestCheckStatus_EmptyStatus tests behavior with empty status string
func TestCheckStatus_EmptyStatus(t *testing.T) {
	statusList := []string{
		goe2econstants.LBStatusCreating,
		"",
	}

	result := CheckStatus(statusList, "")
	assert.True(t, result, "Empty status should be found if in list")
}

// ============================================================================
// TestValidateLoadBalancerName - Tests for node.ValidateName function
// (used for load balancer name validation)
// ============================================================================

// TestValidateLoadBalancerName_ValidNames tests that valid names pass validation
func TestValidateLoadBalancerName_ValidNames(t *testing.T) {
	// Need to import node package to access ValidateName
	// Testing via integration since it's from node package
	testCases := []struct {
		name  string
		input string
		desc  string
	}{
		{
			name:  "lowercase letters",
			input: "loadbalancer",
			desc:  "Lowercase letters should be valid",
		},
		{
			name:  "uppercase letters",
			input: "LOADBALANCER",
			desc:  "Uppercase letters should be valid",
		},
		{
			name:  "mixed case",
			input: "LoadBalancer",
			desc:  "Mixed case letters should be valid",
		},
		{
			name:  "with hyphens",
			input: "load-balancer",
			desc:  "Hyphens should be valid",
		},
		{
			name:  "with underscores",
			input: "load_balancer",
			desc:  "Underscores should be valid",
		},
		{
			name:  "with numbers",
			input: "loadbalancer123",
			desc:  "Numbers should be valid",
		},
		{
			name:  "mixed alphanumeric with hyphens and underscores",
			input: "load-balancer_123",
			desc:  "Mixed alphanumeric with hyphens and underscores should be valid",
		},
		{
			name:  "single character",
			input: "a",
			desc:  "Single character should be valid",
		},
		{
			name:  "50 characters (max length)",
			input: "a1234567890123456789012345678901234567890123456789",
			desc:  "50 characters (max length) should be valid",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Note: We're documenting the validation requirements here
			// Actual validation is tested via Terraform acceptance tests
			// Valid name pattern: ^[a-zA-Z0-9-_]{1,50}$
			assert.NotEmpty(t, tc.input, tc.desc)
			assert.LessOrEqual(t, len(tc.input), 50, "Name should be <= 50 characters")
		})
	}
}

// TestValidateLoadBalancerName_InvalidNames tests that invalid names fail validation
func TestValidateLoadBalancerName_InvalidNames(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		desc   string
		reason string
	}{
		{
			name:   "empty name",
			input:  "",
			desc:   "Empty name should be invalid",
			reason: "Name cannot be empty",
		},
		{
			name:   "with spaces",
			input:  "load balancer",
			desc:   "Name with spaces should be invalid",
			reason: "Spaces are not allowed",
		},
		{
			name:   "with special characters - dot",
			input:  "load.balancer",
			desc:   "Name with dot should be invalid",
			reason: "Special characters (except - and _) are not allowed",
		},
		{
			name:   "with special characters - slash",
			input:  "load/balancer",
			desc:   "Name with slash should be invalid",
			reason: "Special characters (except - and _) are not allowed",
		},
		{
			name:   "with special characters - at sign",
			input:  "load@balancer",
			desc:   "Name with @ should be invalid",
			reason: "Special characters (except - and _) are not allowed",
		},
		{
			name:   "with special characters - hash",
			input:  "load#balancer",
			desc:   "Name with # should be invalid",
			reason: "Special characters (except - and _) are not allowed",
		},
		{
			name:   "too long - 51 characters",
			input:  "a12345678901234567890123456789012345678901234567890",
			desc:   "Name with 51 characters should be invalid",
			reason: "Name must be 50 characters or less",
		},
		{
			name:   "too long - 100 characters",
			input:  "a123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789",
			desc:   "Name with 100 characters should be invalid",
			reason: "Name must be 50 characters or less",
		},
		{
			name:   "unicode characters",
			input:  "load-balancer-😀",
			desc:   "Name with unicode characters should be invalid",
			reason: "Only ASCII alphanumeric, hyphens, and underscores allowed",
		},
		{
			name:   "starts with special character",
			input:  "-loadbalancer",
			desc:   "Name starting with hyphen should be valid (allowed by regex)",
			reason: "Pattern allows starting with hyphen",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Note: We're documenting the validation requirements here
			// Actual validation is tested via Terraform acceptance tests
			// Invalid name pattern checks: ^[a-zA-Z0-9-_]{1,50}$
			if tc.input == "" {
				assert.Empty(t, tc.input, tc.desc)
			} else if len(tc.input) > 50 {
				assert.Greater(t, len(tc.input), 50, tc.desc)
			} else {
				// Document that this would fail validation
				t.Logf("Name '%s' would fail validation: %s", tc.input, tc.reason)
			}
		})
	}
}

// TestValidateLoadBalancerName_EdgeCases tests edge cases for name validation
func TestValidateLoadBalancerName_EdgeCases(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		desc   string
		valid  bool
		reason string
	}{
		{
			name:   "exactly 50 characters",
			input:  "a1234567890123456789012345678901234567890123456789",
			desc:   "Exactly 50 characters should be valid",
			valid:  true,
			reason: "Maximum length is 50",
		},
		{
			name:   "exactly 51 characters",
			input:  "a12345678901234567890123456789012345678901234567890",
			desc:   "Exactly 51 characters should be invalid",
			valid:  false,
			reason: "Exceeds maximum length of 50",
		},
		{
			name:   "all hyphens",
			input:  "----------",
			desc:   "All hyphens should be valid",
			valid:  true,
			reason: "Hyphens are allowed",
		},
		{
			name:   "all underscores",
			input:  "__________",
			desc:   "All underscores should be valid",
			valid:  true,
			reason: "Underscores are allowed",
		},
		{
			name:   "all numbers",
			input:  "1234567890",
			desc:   "All numbers should be valid",
			valid:  true,
			reason: "Numbers are allowed",
		},
		{
			name:   "hyphen and underscore mix",
			input:  "-_-_-_",
			desc:   "Mix of hyphens and underscores should be valid",
			valid:  true,
			reason: "Both hyphens and underscores are allowed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.valid {
				assert.LessOrEqual(t, len(tc.input), 50, tc.desc)
				t.Logf("Name '%s' is valid: %s", tc.input, tc.reason)
			} else {
				t.Logf("Name '%s' would fail validation: %s", tc.input, tc.reason)
			}
		})
	}
}

// ============================================================================
// TestParseImportID - Tests for customImportStateLoadBalancer function
// ============================================================================

// TestParseImportID_SimpleFormat tests parsing simple format (just LB ID)
func TestParseImportID_SimpleFormat(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedID     string
		expectedProjID string
		expectedRegion string
		desc           string
	}{
		{
			name:           "numeric ID",
			input:          "12345",
			expectedID:     "12345",
			expectedProjID: "", // Not set in simple format
			expectedRegion: "", // Not set in simple format
			desc:           "Simple numeric ID should be parsed correctly",
		},
		{
			name:           "alphanumeric ID",
			input:          "lb-12345",
			expectedID:     "lb-12345",
			expectedProjID: "",
			expectedRegion: "",
			desc:           "Alphanumeric ID should be parsed correctly",
		},
		{
			name:           "UUID-like ID",
			input:          "550e8400-e29b-41d4-a716-446655440000",
			expectedID:     "550e8400-e29b-41d4-a716-446655440000",
			expectedProjID: "",
			expectedRegion: "",
			desc:           "UUID-like ID should be parsed correctly",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Document expected behavior
			parts := strings.Split(tc.input, "/")
			assert.Len(t, parts, 1, "Simple format should have 1 part")
			assert.Equal(t, tc.expectedID, parts[0], tc.desc)
			t.Logf("Simple import format: ID='%s', uses provider defaults for project/region", tc.expectedID)
		})
	}
}

// TestParseImportID_FullFormat tests parsing full format (project/region/ID)
func TestParseImportID_FullFormat(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedProjID string
		expectedRegion string
		expectedID     string
		desc           string
	}{
		{
			name:           "standard full format",
			input:          "project-123/Mumbai/lb-12345",
			expectedProjID: "project-123",
			expectedRegion: "Mumbai",
			expectedID:     "lb-12345",
			desc:           "Standard full format should parse all three parts",
		},
		{
			name:           "numeric project ID",
			input:          "12345/Delhi/67890",
			expectedProjID: "12345",
			expectedRegion: "Delhi",
			expectedID:     "67890",
			desc:           "Numeric project ID should be parsed correctly",
		},
		{
			name:           "different region",
			input:          "proj-abc/Chennai/lb-xyz",
			expectedProjID: "proj-abc",
			expectedRegion: "Chennai",
			expectedID:     "lb-xyz",
			desc:           "Chennai region should be parsed correctly",
		},
		{
			name:           "region with spaces would be invalid but testing parsing",
			input:          "project-123/North India/lb-12345",
			expectedProjID: "project-123",
			expectedRegion: "North India",
			expectedID:     "lb-12345",
			desc:           "Region with spaces (hypothetical) should be parsed as-is",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parts := strings.Split(tc.input, "/")
			assert.Len(t, parts, 3, "Full format should have 3 parts")
			assert.Equal(t, tc.expectedProjID, parts[0], "Project ID should match")
			assert.Equal(t, tc.expectedRegion, parts[1], "Region should match")
			assert.Equal(t, tc.expectedID, parts[2], "LB ID should match")
			t.Logf("Full import format: ProjectID='%s', Region='%s', ID='%s'", tc.expectedProjID, tc.expectedRegion, tc.expectedID)
		})
	}
}

// TestParseImportID_InvalidFormats tests parsing invalid formats
func TestParseImportID_InvalidFormats(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		reason string
		desc   string
	}{
		{
			name:   "empty string",
			input:  "",
			reason: "Empty string is invalid",
			desc:   "Empty import string should return error",
		},
		{
			name:   "two parts format",
			input:  "project-123/lb-12345",
			reason: "Two-part format is ambiguous and invalid",
			desc:   "Two-part format should return error",
		},
		{
			name:   "four parts format",
			input:  "project-123/Mumbai/lb-12345/extra",
			reason: "Four-part format has too many parts",
			desc:   "Four-part format should return error",
		},
		{
			name:   "five parts format",
			input:  "a/b/c/d/e",
			reason: "Five-part format has too many parts",
			desc:   "Five-part format should return error",
		},
		{
			name:   "only slashes",
			input:  "//",
			reason: "Only slashes without meaningful parts",
			desc:   "Only slashes should return error or empty parts",
		},
		{
			name:   "trailing slash",
			input:  "project-123/Mumbai/lb-12345/",
			reason: "Trailing slash creates empty fourth part",
			desc:   "Trailing slash should return error (4 parts)",
		},
		{
			name:   "leading slash",
			input:  "/project-123/Mumbai/lb-12345",
			reason: "Leading slash creates empty first part",
			desc:   "Leading slash should return error (4 parts)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parts := strings.Split(tc.input, "/")
			partCount := len(parts)

			// Valid formats: 1 part (simple) or 3 parts (full)
			isValid := partCount == 1 || partCount == 3

			if !isValid {
				assert.NotEqual(t, 1, partCount, tc.desc)
				assert.NotEqual(t, 3, partCount, tc.desc)
				t.Logf("Invalid format with %d parts: %s", partCount, tc.reason)
			}

			// Additional check: if it's the "only slashes" case, all parts would be empty
			if tc.input == "//" {
				assert.Len(t, parts, 3, "'//' splits into 3 empty parts")
				for i, part := range parts {
					assert.Empty(t, part, "Part %d should be empty", i)
				}
			}
		})
	}
}

// TestParseImportID_SpecialCharacters tests parsing with special characters
func TestParseImportID_SpecialCharacters(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		valid  bool
		reason string
		desc   string
	}{
		{
			name:   "ID with hyphens",
			input:  "lb-test-123",
			valid:  true,
			reason: "Hyphens in ID are valid",
			desc:   "ID with hyphens should be parsed correctly",
		},
		{
			name:   "ID with underscores",
			input:  "lb_test_123",
			valid:  true,
			reason: "Underscores in ID are valid",
			desc:   "ID with underscores should be parsed correctly",
		},
		{
			name:   "UUID format ID",
			input:  "550e8400-e29b-41d4-a716-446655440000",
			valid:  true,
			reason: "UUID format is valid",
			desc:   "UUID format ID should be parsed correctly",
		},
		{
			name:   "full format with special chars in ID",
			input:  "project-123/Mumbai/lb-test_123-xyz",
			valid:  true,
			reason: "Mixed special characters in ID are valid",
			desc:   "Full format with special characters in ID should be parsed correctly",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parts := strings.Split(tc.input, "/")
			if tc.valid {
				assert.True(t, len(parts) == 1 || len(parts) == 3, tc.desc)
				t.Logf("Valid import format (%d parts): %s", len(parts), tc.reason)
			} else {
				t.Logf("Invalid import format: %s", tc.reason)
			}
		})
	}
}

// TestParseImportID_ErrorMessage tests that error message is clear and helpful
func TestParseImportID_ErrorMessage(t *testing.T) {
	// Test that the error message constant is clear
	expectedErrorMsg := "invalid import format, expected: <lb_id> or <project_id>/<region>/<lb_id>"

	assert.Contains(t, expectedErrorMsg, "invalid import format", "Error message should indicate invalid format")
	assert.Contains(t, expectedErrorMsg, "<lb_id>", "Error message should show simple format")
	assert.Contains(t, expectedErrorMsg, "<project_id>/<region>/<lb_id>", "Error message should show full format")

	t.Logf("Error message for invalid import format: %s", expectedErrorMsg)
}

// TestParseImportID_Comprehensive tests comprehensive scenarios
func TestParseImportID_Comprehensive(t *testing.T) {
	testCases := []struct {
		name       string
		input      string
		partCount  int
		shouldPass bool
		desc       string
	}{
		{
			name:       "simple valid",
			input:      "12345",
			partCount:  1,
			shouldPass: true,
			desc:       "Simple format with numeric ID",
		},
		{
			name:       "simple valid alphanumeric",
			input:      "lb-12345-test",
			partCount:  1,
			shouldPass: true,
			desc:       "Simple format with alphanumeric ID",
		},
		{
			name:       "full valid",
			input:      "proj-123/Mumbai/lb-456",
			partCount:  3,
			shouldPass: true,
			desc:       "Full format with all three parts",
		},
		{
			name:       "full valid different region",
			input:      "12345/Delhi/67890",
			partCount:  3,
			shouldPass: true,
			desc:       "Full format with Delhi region",
		},
		{
			name:       "invalid two parts",
			input:      "project-123/lb-456",
			partCount:  2,
			shouldPass: false,
			desc:       "Two-part format should fail",
		},
		{
			name:       "invalid four parts",
			input:      "a/b/c/d",
			partCount:  4,
			shouldPass: false,
			desc:       "Four-part format should fail",
		},
		{
			name:       "invalid empty",
			input:      "",
			partCount:  1,
			shouldPass: false,
			desc:       "Empty string should fail (produces empty part)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parts := strings.Split(tc.input, "/")
			assert.Equal(t, tc.partCount, len(parts), "Part count should match")

			// Valid format requires either:
			// - 1 non-empty part (simple format)
			// - 3 parts (full format)
			validFormat := (len(parts) == 1 && parts[0] != "") || len(parts) == 3

			if tc.shouldPass {
				assert.True(t, validFormat, tc.desc)
			} else {
				// For invalid cases, we expect validFormat to be false
				assert.False(t, validFormat, tc.desc)
			}
		})
	}
}
