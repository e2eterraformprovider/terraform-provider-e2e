package autoscaling_test

import (
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/autoscaling"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeStatus(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Starting to Running",
			input:    goe2econstants.AutoscalingScalerGroupStatusStarting,
			expected: goe2econstants.AutoscalingScalerGroupStatusRunning,
		},
		{
			name:     "Stopping to Stopped",
			input:    goe2econstants.AutoscalingScalerGroupStatusStopping,
			expected: goe2econstants.AutoscalingScalerGroupStatusStopped,
		},
		{
			name:     "starting to running (lowercase)",
			input:    goe2econstants.AutoscalingScalerGroupStatusStartingLower,
			expected: goe2econstants.AutoscalingScalerGroupStatusRunningLower,
		},
		{
			name:     "stopping to stopped (lowercase)",
			input:    goe2econstants.AutoscalingScalerGroupStatusStoppingLower,
			expected: goe2econstants.AutoscalingScalerGroupStatusStoppedLower,
		},
		{
			name:     "Running stays Running",
			input:    goe2econstants.AutoscalingScalerGroupStatusRunning,
			expected: goe2econstants.AutoscalingScalerGroupStatusRunning,
		},
		{
			name:     "Stopped stays Stopped",
			input:    goe2econstants.AutoscalingScalerGroupStatusStopped,
			expected: goe2econstants.AutoscalingScalerGroupStatusStopped,
		},
		{
			name:     "running stays running",
			input:    goe2econstants.AutoscalingScalerGroupStatusRunningLower,
			expected: goe2econstants.AutoscalingScalerGroupStatusRunningLower,
		},
		{
			name:     "stopped stays stopped",
			input:    goe2econstants.AutoscalingScalerGroupStatusStoppedLower,
			expected: goe2econstants.AutoscalingScalerGroupStatusStoppedLower,
		},
		{
			name:     "Unknown status unchanged",
			input:    "Unknown",
			expected: "Unknown",
		},
		{
			name:     "Empty string unchanged",
			input:    "",
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := autoscaling.NormalizeStatus(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetImageName_ErrorMessageIsStable(t *testing.T) {
	d := schema.TestResourceDataRaw(t, autoscaling.ResourceScalerGroup().Schema, map[string]interface{}{})
	_, err := autoscaling.GetImageName(d)
	assert.Error(t, err)
	assert.Equal(t, autoscaling.ErrEitherVMImageOrImageRequired, err.Error())
}

func TestExpandNetworkConfig(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]interface{}
		expected *autoscaling.NetworkConfig
	}{
		{
			name: "full network_config block",
			data: map[string]interface{}{
				"network_config": []interface{}{
					map[string]interface{}{
						"assign_public_ip": true,
						"vpc_names":        []interface{}{"vpc1", "vpc2"},
						"security_groups":  []interface{}{1, 2, 3},
					},
				},
			},
			expected: &autoscaling.NetworkConfig{
				AssignPublicIP: true,
				VPCNames:       []string{"vpc1", "vpc2"},
				SecurityGroups: []int{1, 2, 3},
			},
		},
		{
			name: "partial network_config block",
			data: map[string]interface{}{
				"network_config": []interface{}{
					map[string]interface{}{
						"assign_public_ip": false,
						"vpc_names":        []interface{}{"vpc1"},
					},
				},
			},
			expected: &autoscaling.NetworkConfig{
				AssignPublicIP: false,
				VPCNames:       []string{"vpc1"},
				SecurityGroups: []int{},
			},
		},
		{
			name:     "no network_config block",
			data:     map[string]interface{}{},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, autoscaling.ResourceScalerGroup().Schema, tt.data)
			result := autoscaling.ExpandNetworkConfig(d)

			if tt.expected == nil {
				assert.Nil(t, result)
				return
			}

			assert.NotNil(t, result)
			assert.Equal(t, tt.expected.AssignPublicIP, result.AssignPublicIP)
			assert.Equal(t, tt.expected.VPCNames, result.VPCNames)
			assert.Equal(t, tt.expected.SecurityGroups, result.SecurityGroups)
		})
	}
}

func TestFlattenNetworkConfig(t *testing.T) {
	tests := []struct {
		name           string
		assignPublicIP bool
		vpcNames       []string
		securityGroups []int
		expectedEmpty  bool
		expected       map[string]interface{}
	}{
		{
			name:           "full config",
			assignPublicIP: true,
			vpcNames:       []string{"vpc1", "vpc2"},
			securityGroups: []int{1, 2},
			expectedEmpty:  false,
			expected: map[string]interface{}{
				"assign_public_ip": true,
				"vpc_names":        []interface{}{"vpc1", "vpc2"},
				"security_groups":  []interface{}{1, 2},
			},
		},
		{
			name:           "partial config",
			assignPublicIP: false,
			vpcNames:       []string{"vpc1"},
			securityGroups: []int{},
			expectedEmpty:  false,
			expected: map[string]interface{}{
				"assign_public_ip": false,
				"vpc_names":        []interface{}{"vpc1"},
			},
		},
		{
			name:           "empty config",
			assignPublicIP: false,
			vpcNames:       []string{},
			securityGroups: []int{},
			expectedEmpty:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := autoscaling.FlattenNetworkConfig(tt.assignPublicIP, tt.vpcNames, tt.securityGroups)

			if tt.expectedEmpty {
				assert.Empty(t, result)
				return
			}

			assert.Len(t, result, 1)
			config := result[0]
			assert.Equal(t, tt.expected["assign_public_ip"], config["assign_public_ip"])

			if tt.expected["vpc_names"] != nil {
				assert.Equal(t, tt.expected["vpc_names"], config["vpc_names"])
			}

			if tt.expected["security_groups"] != nil {
				assert.Equal(t, tt.expected["security_groups"], config["security_groups"])
			}
		})
	}
}

func TestStringSlicesEqual(t *testing.T) {
	tests := []struct {
		name     string
		a        []string
		b        []string
		expected bool
	}{
		{
			name:     "equal slices",
			a:        []string{"a", "b", "c"},
			b:        []string{"a", "b", "c"},
			expected: true,
		},
		{
			name:     "different lengths",
			a:        []string{"a", "b"},
			b:        []string{"a", "b", "c"},
			expected: false,
		},
		{
			name:     "different values",
			a:        []string{"a", "b"},
			b:        []string{"a", "c"},
			expected: false,
		},
		{
			name:     "both empty",
			a:        []string{},
			b:        []string{},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := autoscaling.StringSlicesEqual(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIntSlicesEqual(t *testing.T) {
	tests := []struct {
		name     string
		a        []int
		b        []int
		expected bool
	}{
		{
			name:     "equal slices",
			a:        []int{1, 2, 3},
			b:        []int{1, 2, 3},
			expected: true,
		},
		{
			name:     "different lengths",
			a:        []int{1, 2},
			b:        []int{1, 2, 3},
			expected: false,
		},
		{
			name:     "different values",
			a:        []int{1, 2},
			b:        []int{1, 3},
			expected: false,
		},
		{
			name:     "both empty",
			a:        []int{},
			b:        []int{},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := autoscaling.IntSlicesEqual(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestResourceScalerGroupSchema_EncryptionPassphraseSensitive(t *testing.T) {
	r := autoscaling.ResourceScalerGroup()

	s, ok := r.Schema[tfconstants.AttrEncryptionPassphrase]
	if !ok {
		t.Fatalf("schema missing %q", tfconstants.AttrEncryptionPassphrase)
	}

	assert.True(t, s.Sensitive, "encryption_passphrase must be Sensitive to avoid plan/apply output leaks")
	assert.True(t, s.ForceNew, "encryption_passphrase should remain ForceNew (API contract)")
}

func TestAutoscalingConstants_StatusSets(t *testing.T) {
	assert.ElementsMatch(t,
		[]string{
			goe2econstants.AutoscalingScalerGroupStatusRunning,
			goe2econstants.AutoscalingScalerGroupStatusStopped,
		},
		tfconstants.AutoscalingScalerGroupProvisionStatusAllowed,
	)

	assert.ElementsMatch(t,
		[]string{
			goe2econstants.AutoscalingScalerGroupStatusRunningLower,
			goe2econstants.AutoscalingScalerGroupStatusStoppedLower,
		},
		tfconstants.AutoscalingScalerGroupStatusAllowed,
	)

	// Terraform "state requirements" must accept both legacy (Running/Stopped) and V3 (running/stopped).
	assert.Contains(t, tfconstants.AutoscalingScalerGroupRunningStates, goe2econstants.AutoscalingScalerGroupStatusRunning)
	assert.Contains(t, tfconstants.AutoscalingScalerGroupRunningStates, goe2econstants.AutoscalingScalerGroupStatusRunningLower)
	assert.Contains(t, tfconstants.AutoscalingScalerGroupStoppedStates, goe2econstants.AutoscalingScalerGroupStatusStopped)
	assert.Contains(t, tfconstants.AutoscalingScalerGroupStoppedStates, goe2econstants.AutoscalingScalerGroupStatusStoppedLower)
}
