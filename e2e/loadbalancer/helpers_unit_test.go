package loadbalancer

import (
	"testing"

	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// Test: GetLbPort
// ============================================================================

func TestGetLbPort_HTTP(t *testing.T) {
	result := GetLbPort(goe2econstants.LBModeHTTP)
	assert.Equal(t, goe2econstants.LBPortHTTP, result)
}

func TestGetLbPort_HTTPS(t *testing.T) {
	result := GetLbPort(goe2econstants.LBModeHTTPS)
	assert.Equal(t, goe2econstants.LBPortHTTPS, result)
}

func TestGetLbPort_OtherMode(t *testing.T) {
	result := GetLbPort("tcp")
	// Should default to HTTPS port for non-HTTP modes
	assert.Equal(t, goe2econstants.LBPortHTTPS, result)
}

// ============================================================================
// Test: SetLoadBalancerStatus
// ============================================================================

func TestSetLoadBalancerStatus_RunningWithHealthyBackend(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		tfconstants.AttrStatus: {Type: schema.TypeString},
	}, map[string]interface{}{})

	statusDetail := map[string]interface{}{
		"status": goe2econstants.LBStatusRunningAPI,
		"data_monitor": map[string]interface{}{
			"status": true,
		},
	}

	err := SetLoadBalancerStatus(d, statusDetail)

	require.NoError(t, err)
	assert.Equal(t, goe2econstants.LBStatusRunning, d.Get(tfconstants.AttrStatus))
}

func TestSetLoadBalancerStatus_RunningWithUnhealthyBackend(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		tfconstants.AttrStatus: {Type: schema.TypeString},
	}, map[string]interface{}{})

	statusDetail := map[string]interface{}{
		"status": goe2econstants.LBStatusRunningAPI,
		"data_monitor": map[string]interface{}{
			"status": false,
		},
	}

	err := SetLoadBalancerStatus(d, statusDetail)

	require.NoError(t, err)
	assert.Equal(t, goe2econstants.LBStatusBackendFailure, d.Get(tfconstants.AttrStatus))
}

func TestSetLoadBalancerStatus_RunningWithNoDataMonitor(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		tfconstants.AttrStatus: {Type: schema.TypeString},
	}, map[string]interface{}{})

	statusDetail := map[string]interface{}{
		"status":       goe2econstants.LBStatusRunningAPI,
		"data_monitor": map[string]interface{}{},
	}

	err := SetLoadBalancerStatus(d, statusDetail)

	require.NoError(t, err)
	assert.Equal(t, goe2econstants.LBStatusBackendUnavailable, d.Get(tfconstants.AttrStatus))
}

func TestSetLoadBalancerStatus_PoweredOff(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		tfconstants.AttrStatus: {Type: schema.TypeString},
	}, map[string]interface{}{})

	statusDetail := map[string]interface{}{
		"status": goe2econstants.LBStatusPoweredOffAPI,
	}

	err := SetLoadBalancerStatus(d, statusDetail)

	require.NoError(t, err)
	assert.Equal(t, goe2econstants.LBStatusPoweredOff, d.Get(tfconstants.AttrStatus))
}

func TestSetLoadBalancerStatus_Creating(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		tfconstants.AttrStatus: {Type: schema.TypeString},
	}, map[string]interface{}{})

	statusDetail := map[string]interface{}{
		"status": goe2econstants.LBStatusCreating,
	}

	err := SetLoadBalancerStatus(d, statusDetail)

	require.NoError(t, err)
	assert.Equal(t, goe2econstants.LBStatusCreating, d.Get(tfconstants.AttrStatus))
}

func TestSetLoadBalancerStatus_Deploying(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		tfconstants.AttrStatus: {Type: schema.TypeString},
	}, map[string]interface{}{})

	statusDetail := map[string]interface{}{
		"status": goe2econstants.LBStatusDeploying,
	}

	err := SetLoadBalancerStatus(d, statusDetail)

	require.NoError(t, err)
	assert.Equal(t, goe2econstants.LBStatusDeploying, d.Get(tfconstants.AttrStatus))
}

func TestSetLoadBalancerStatus_Upgrading(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		tfconstants.AttrStatus: {Type: schema.TypeString},
	}, map[string]interface{}{})

	statusDetail := map[string]interface{}{
		"status": goe2econstants.LBStatusUpgradingAPI,
	}

	err := SetLoadBalancerStatus(d, statusDetail)

	require.NoError(t, err)
	assert.Equal(t, goe2econstants.LBStatusUpgrading, d.Get(tfconstants.AttrStatus))
}

func TestSetLoadBalancerStatus_Unknown(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		tfconstants.AttrStatus: {Type: schema.TypeString},
	}, map[string]interface{}{})

	statusDetail := map[string]interface{}{
		"status": "unknown-status",
	}

	err := SetLoadBalancerStatus(d, statusDetail)

	require.NoError(t, err)
	assert.Equal(t, goe2econstants.LBStatusError, d.Get(tfconstants.AttrStatus))
}

// Note: CheckStatus is already tested in resource_loadbalancer_unit_test.go
