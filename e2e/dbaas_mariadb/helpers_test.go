package dbaas_mariadb

import (
	"testing"

	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeStatus(t *testing.T) {
	tests := []struct {
		name       string
		apiStatus  string
		wantStatus string
	}{
		{
			name:       "SUSPENDED converts to STOPPED",
			apiStatus:  goe2econstants.DBaaSStatusSuspended,
			wantStatus: goe2econstants.DBaaSStatusStopped,
		},
		{
			name:       "STOPPED passes through unchanged",
			apiStatus:  goe2econstants.DBaaSStatusStopped,
			wantStatus: goe2econstants.DBaaSStatusStopped,
		},
		{
			name:       "RUNNING passes through unchanged",
			apiStatus:  goe2econstants.DBaaSStatusRunning,
			wantStatus: goe2econstants.DBaaSStatusRunning,
		},
		{
			name:       "RESTARTING passes through unchanged",
			apiStatus:  goe2econstants.DBaaSStatusRestarting,
			wantStatus: goe2econstants.DBaaSStatusRestarting,
		},
		{
			name:       "empty string passes through unchanged",
			apiStatus:  "",
			wantStatus: "",
		},
		{
			name:       "case sensitivity - uppercase SUSPENDED converts to STOPPED",
			apiStatus:  goe2econstants.DBaaSStatusSuspended,
			wantStatus: goe2econstants.DBaaSStatusStopped,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeStatus(tt.apiStatus)
			assert.Equal(t, tt.wantStatus, result)
		})
	}
}

func TestExpandStringSet(t *testing.T) {
	t.Run("with empty slice returns empty map", func(t *testing.T) {
		result := expandStringSet([]interface{}{})
		assert.NotNil(t, result)
		assert.Empty(t, result)
	})

	t.Run("with single string element", func(t *testing.T) {
		result := expandStringSet([]interface{}{"vpc1"})
		assert.Len(t, result, 1)
		_, exists := result["vpc1"]
		assert.True(t, exists, "vpc1 should exist in the map")
	})

	t.Run("with multiple string elements", func(t *testing.T) {
		result := expandStringSet([]interface{}{"vpc1", "vpc2", "vpc3"})
		assert.Len(t, result, 3)
		_, exists1 := result["vpc1"]
		_, exists2 := result["vpc2"]
		_, exists3 := result["vpc3"]
		assert.True(t, exists1, "vpc1 should exist in the map")
		assert.True(t, exists2, "vpc2 should exist in the map")
		assert.True(t, exists3, "vpc3 should exist in the map")
	})

	t.Run("with duplicate elements creates set with unique keys", func(t *testing.T) {
		result := expandStringSet([]interface{}{"vpc1", "vpc2", "vpc1", "vpc3", "vpc2"})
		// Duplicates should be deduplicated (map keys are unique)
		assert.Len(t, result, 3, "duplicate elements should be deduplicated")
		_, exists1 := result["vpc1"]
		_, exists2 := result["vpc2"]
		_, exists3 := result["vpc3"]
		assert.True(t, exists1, "vpc1 should exist in the map")
		assert.True(t, exists2, "vpc2 should exist in the map")
		assert.True(t, exists3, "vpc3 should exist in the map")
	})

	t.Run("type assertion panics are handled", func(t *testing.T) {
		// This test verifies that the function expects string types
		// If a non-string is passed, it will panic (which is expected behavior)
		// In practice, the schema validation should prevent this
		assert.Panics(t, func() {
			expandStringSet([]interface{}{123}) // int instead of string
		}, "should panic when non-string type is passed")
	})
}
