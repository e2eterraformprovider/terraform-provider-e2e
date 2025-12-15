package container_registry

import (
	"testing"

	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// TestSetContainerRegistryState_AllFields tests setting all fields from a complete API response
func TestSetContainerRegistryState_AllFields(t *testing.T) {
	d := schema.TestResourceDataRaw(t, ResourceContainerRegistry().Schema, map[string]interface{}{})
	d.SetId("12345")

	registry := &goe2e.ContainerRegistry{
		ID:           12345,
		ProjectName:  "test-registry",
		ProjectSize:  1024.5,
		DomainName:   "test-registry.e2enetworks.net",
		PreventVul:   true,
		Severity:     goe2econstants.ContainerRegistrySeverityHigh,
		State:        "active",
		IsPublic:     false,
		StorageLimit: 10737418240, // 10GB
		Location:     "Mumbai",
		Customer:     999,
		ProjectID:    777,
		CreatedAt:    "2025-12-01T10:00:00Z",
		UpdatedAt:    "2025-12-12T15:30:00Z",
	}

	err := setContainerRegistryState(d, registry)

	assert.NoError(t, err)
	assert.Equal(t, "test-registry", d.Get(tfconstants.AttrProjectName))
	assert.Equal(t, true, d.Get(tfconstants.AttrPreventVulnerabilities))
	assert.Equal(t, goe2econstants.ContainerRegistrySeverityHigh, d.Get(tfconstants.AttrSeverity))
	assert.Equal(t, "active", d.Get(tfconstants.AttrStatus))
	assert.Equal(t, "active", d.Get(tfconstants.AttrSetupStatus))
	assert.Equal(t, "test-registry.e2enetworks.net", d.Get(tfconstants.AttrDomainName))
	assert.Equal(t, 1024.5, d.Get(tfconstants.AttrProjectSize))
	assert.Equal(t, 10737418240, d.Get(tfconstants.AttrStorageLimit))
	assert.Equal(t, false, d.Get(tfconstants.AttrIsPublic))
	assert.Equal(t, "2025-12-01T10:00:00Z", d.Get(tfconstants.AttrCreatedAt))
	assert.Equal(t, "2025-12-12T15:30:00Z", d.Get(tfconstants.AttrUpdatedAt))
}

// TestSetContainerRegistryState_EmptyFields tests behavior with empty/zero values
func TestSetContainerRegistryState_EmptyFields(t *testing.T) {
	d := schema.TestResourceDataRaw(t, ResourceContainerRegistry().Schema, map[string]interface{}{})
	d.SetId("1")

	registry := &goe2e.ContainerRegistry{
		ID:           1,
		ProjectName:  "",
		ProjectSize:  0,
		DomainName:   "",
		PreventVul:   false,
		Severity:     goe2econstants.ContainerRegistrySeverityLow,
		State:        "",
		IsPublic:     false,
		StorageLimit: 0,
		CreatedAt:    "",
		UpdatedAt:    "",
	}

	err := setContainerRegistryState(d, registry)

	assert.NoError(t, err)
	assert.Equal(t, "", d.Get(tfconstants.AttrProjectName))
	assert.Equal(t, false, d.Get(tfconstants.AttrPreventVulnerabilities))
	assert.Equal(t, goe2econstants.ContainerRegistrySeverityLow, d.Get(tfconstants.AttrSeverity))
	assert.Equal(t, "", d.Get(tfconstants.AttrStatus))
	assert.Equal(t, "", d.Get(tfconstants.AttrSetupStatus))
	assert.Equal(t, "", d.Get(tfconstants.AttrDomainName))
	assert.Equal(t, float64(0), d.Get(tfconstants.AttrProjectSize))
	assert.Equal(t, 0, d.Get(tfconstants.AttrStorageLimit))
}

// TestSetContainerRegistryState_BooleanConversion tests PreventVul boolean handling
func TestSetContainerRegistryState_BooleanConversion(t *testing.T) {
	testCases := []struct {
		name       string
		preventVul bool
		expected   bool
	}{
		{
			name:       "PreventVul true",
			preventVul: true,
			expected:   true,
		},
		{
			name:       "PreventVul false",
			preventVul: false,
			expected:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, ResourceContainerRegistry().Schema, map[string]interface{}{})
			d.SetId("1")

			registry := &goe2e.ContainerRegistry{
				ID:          1,
				ProjectName: "test",
				PreventVul:  tc.preventVul,
				Severity:    goe2econstants.ContainerRegistrySeverityLow,
				State:       "active",
			}

			err := setContainerRegistryState(d, registry)

			assert.NoError(t, err)
			assert.Equal(t, tc.expected, d.Get(tfconstants.AttrPreventVulnerabilities))
		})
	}
}

// TestSetContainerRegistryState_NumericTypes tests float64 and int type handling
func TestSetContainerRegistryState_NumericTypes(t *testing.T) {
	d := schema.TestResourceDataRaw(t, ResourceContainerRegistry().Schema, map[string]interface{}{})
	d.SetId("1")

	registry := &goe2e.ContainerRegistry{
		ID:           1,
		ProjectName:  "test",
		ProjectSize:  9999.99,    // float64
		StorageLimit: 5368709120, // int (5GB)
		PreventVul:   false,
		Severity:     goe2econstants.ContainerRegistrySeverityLow,
		State:        "active",
	}

	err := setContainerRegistryState(d, registry)

	assert.NoError(t, err)

	projectSize := d.Get(tfconstants.AttrProjectSize)
	assert.IsType(t, float64(0), projectSize)
	assert.Equal(t, 9999.99, projectSize)

	storageLimit := d.Get(tfconstants.AttrStorageLimit)
	assert.IsType(t, 0, storageLimit)
	assert.Equal(t, 5368709120, storageLimit)
}

// TestSetContainerRegistryState_DeprecatedFields tests that setup_status matches status
func TestSetContainerRegistryState_DeprecatedFields(t *testing.T) {
	testCases := []struct {
		name  string
		state string
	}{
		{
			name:  "active state",
			state: "active",
		},
		{
			name:  "pending state",
			state: "pending",
		},
		{
			name:  "suspended state",
			state: "suspended",
		},
		{
			name:  "empty state",
			state: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, ResourceContainerRegistry().Schema, map[string]interface{}{})
			d.SetId("1")

			registry := &goe2e.ContainerRegistry{
				ID:          1,
				ProjectName: "test",
				PreventVul:  false,
				Severity:    goe2econstants.ContainerRegistrySeverityLow,
				State:       tc.state,
			}

			err := setContainerRegistryState(d, registry)

			assert.NoError(t, err)
			assert.Equal(t, tc.state, d.Get(tfconstants.AttrStatus))
			assert.Equal(t, tc.state, d.Get(tfconstants.AttrSetupStatus))
			// Verify both fields have the same value (backward compatibility)
			assert.Equal(t, d.Get(tfconstants.AttrStatus), d.Get(tfconstants.AttrSetupStatus))
		})
	}
}

// TestSetContainerRegistryState_AllSeverityLevels tests all valid severity values
func TestSetContainerRegistryState_AllSeverityLevels(t *testing.T) {
	severityLevels := []string{
		goe2econstants.ContainerRegistrySeverityLow,
		goe2econstants.ContainerRegistrySeverityMedium,
		goe2econstants.ContainerRegistrySeverityHigh,
		goe2econstants.ContainerRegistrySeverityCritical,
		goe2econstants.ContainerRegistrySeverityNone,
	}

	for _, severity := range severityLevels {
		t.Run("severity_"+severity, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, ResourceContainerRegistry().Schema, map[string]interface{}{})
			d.SetId("1")

			registry := &goe2e.ContainerRegistry{
				ID:          1,
				ProjectName: "test",
				PreventVul:  false,
				Severity:    severity,
				State:       "active",
			}

			err := setContainerRegistryState(d, registry)

			assert.NoError(t, err)
			assert.Equal(t, severity, d.Get(tfconstants.AttrSeverity))
		})
	}
}

// TestSetContainerRegistryState_PublicRegistry tests IsPublic field variations
func TestSetContainerRegistryState_PublicRegistry(t *testing.T) {
	testCases := []struct {
		name     string
		isPublic bool
	}{
		{
			name:     "public registry",
			isPublic: true,
		},
		{
			name:     "private registry",
			isPublic: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, ResourceContainerRegistry().Schema, map[string]interface{}{})
			d.SetId("1")

			registry := &goe2e.ContainerRegistry{
				ID:          1,
				ProjectName: "test",
				PreventVul:  false,
				Severity:    goe2econstants.ContainerRegistrySeverityLow,
				State:       "active",
				IsPublic:    tc.isPublic,
			}

			err := setContainerRegistryState(d, registry)

			assert.NoError(t, err)
			assert.Equal(t, tc.isPublic, d.Get(tfconstants.AttrIsPublic))
		})
	}
}

// TestSetContainerRegistryState_TimestampFormats tests various timestamp formats
func TestSetContainerRegistryState_TimestampFormats(t *testing.T) {
	testCases := []struct {
		name      string
		createdAt string
		updatedAt string
	}{
		{
			name:      "ISO 8601 format",
			createdAt: "2025-12-01T10:00:00Z",
			updatedAt: "2025-12-12T15:30:00Z",
		},
		{
			name:      "with timezone",
			createdAt: "2025-12-01T10:00:00+05:30",
			updatedAt: "2025-12-12T15:30:00+05:30",
		},
		{
			name:      "empty timestamps",
			createdAt: "",
			updatedAt: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, ResourceContainerRegistry().Schema, map[string]interface{}{})
			d.SetId("1")

			registry := &goe2e.ContainerRegistry{
				ID:          1,
				ProjectName: "test",
				PreventVul:  false,
				Severity:    goe2econstants.ContainerRegistrySeverityLow,
				State:       "active",
				CreatedAt:   tc.createdAt,
				UpdatedAt:   tc.updatedAt,
			}

			err := setContainerRegistryState(d, registry)

			assert.NoError(t, err)
			assert.Equal(t, tc.createdAt, d.Get(tfconstants.AttrCreatedAt))
			assert.Equal(t, tc.updatedAt, d.Get(tfconstants.AttrUpdatedAt))
		})
	}
}
