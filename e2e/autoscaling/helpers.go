package autoscaling

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Helper functions for V2/V3 field selection

// GetImageName retrieves image name from either V2 or V3 field.
//
// Exported to allow package-external unit tests to lock down UX-stable error messages.
func GetImageName(d *schema.ResourceData) (string, error) {
	if v, ok := d.GetOk("image"); ok {
		return v.(string), nil
	}
	if v, ok := d.GetOk("vm_image_name"); ok {
		return v.(string), nil
	}
	// Keep this message stable (acceptance tests + UX). Prefer the resource-scoped constant.
	return "", errors.New(ErrEitherVMImageOrImageRequired)
}

// getMinSize retrieves min size from either V2 or V3 field
func getMinSize(d *schema.ResourceData) int {
	if v, ok := d.GetOk("min_size"); ok {
		return v.(int)
	}
	if v, ok := d.GetOk(tfconstants.AttrMinNodes); ok {
		return v.(int)
	}
	return 0
}

// getMaxSize retrieves max size from either V2 or V3 field
func getMaxSize(d *schema.ResourceData) int {
	if v, ok := d.GetOk("max_size"); ok {
		return v.(int)
	}
	if v, ok := d.GetOk(tfconstants.AttrMaxNodes); ok {
		return v.(int)
	}
	return 0
}

// getDesiredCapacity retrieves desired capacity from either V2 or V3 field
func getDesiredCapacity(d *schema.ResourceData) int {
	if v, ok := d.GetOk("desired_capacity"); ok {
		return v.(int)
	}
	if v, ok := d.GetOk(tfconstants.AttrDesired); ok {
		return v.(int)
	}
	return 0
}

// getStatus retrieves status from either V2 or V3 field
func getStatus(d *schema.ResourceData) string {
	if v, ok := d.GetOk("status"); ok {
		return v.(string)
	}
	if v, ok := d.GetOk("provision_status"); ok {
		return v.(string)
	}
	return ""
}

// getEnableEncryption retrieves encryption flag from either V2 or V3 field
func getEnableEncryption(d *schema.ResourceData) bool {
	if v, ok := d.GetOk("enable_encryption"); ok {
		return v.(bool)
	}
	if v, ok := d.GetOk(tfconstants.AttrIsEncryptionEnabled); ok {
		return v.(bool)
	}
	return false
}

// getAssignPublicIP retrieves public IP flag from either V2 or V3 field
func getAssignPublicIP(d *schema.ResourceData) bool {
	if v, ok := d.GetOk("assign_public_ip"); ok {
		return v.(bool)
	}
	if v, ok := d.GetOk(tfconstants.AttrPublicIPRequired); ok {
		return v.(bool)
	}
	return true // default
}

// NormalizeStatus normalizes status values (Starting→Running, Stopping→Stopped)
// Exported for testing purposes
func NormalizeStatus(status string) string {
	switch status {
	case goe2econstants.AutoscalingScalerGroupStatusStarting:
		return goe2econstants.AutoscalingScalerGroupStatusRunning
	case goe2econstants.AutoscalingScalerGroupStatusStopping:
		return goe2econstants.AutoscalingScalerGroupStatusStopped
	case goe2econstants.AutoscalingScalerGroupStatusStartingLower:
		return goe2econstants.AutoscalingScalerGroupStatusRunningLower
	case goe2econstants.AutoscalingScalerGroupStatusStoppingLower:
		return goe2econstants.AutoscalingScalerGroupStatusStoppedLower
	default:
		return status
	}
}

// Helper functions for GoE2E client operations

// getSavedImageByName retrieves a saved image by name using GoE2E client
func getSavedImageByName(ctx context.Context, client *goe2e.Client, imageName string) (*goe2e.SavedImage, error) {
	images, _, err := client.Images.GetSavedImages(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch saved images: %w", err)
	}

	for _, img := range images {
		if img.Name == imageName {
			return &img, nil
		}
	}

	return nil, fmt.Errorf("no saved image found with name: %s", imageName)
}

// getDefaultSecurityGroupID retrieves the default security group ID using GoE2E client
func getDefaultSecurityGroupID(ctx context.Context, client *goe2e.Client) (int, error) {
	sgs, _, err := client.SecurityGroups.GetSecurityGroupList(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch security groups: %w", err)
	}

	for _, sg := range sgs {
		if sg.IsDefault {
			// Convert string ID to int
			sgID, err := strconv.Atoi(sg.ID)
			if err != nil {
				return 0, fmt.Errorf("invalid security group ID format: %s", sg.ID)
			}
			return sgID, nil
		}
	}

	return 0, fmt.Errorf("default security group not found")
}

// getPlanDetailsFromPlanName retrieves plan details using GoE2E client
func getPlanDetailsFromPlanName(ctx context.Context, client *goe2e.Client, templateID int, planName string) (string, string, error) {
	planID, slugName, _, err := client.Images.GetPlanDetailsFromPlanName(ctx, templateID, planName)
	if err != nil {
		return "", "", fmt.Errorf("failed to get plan details: %w", err)
	}
	return planID, slugName, nil
}

// flattenNodes converts ScalerGroupNode slice to schema-compatible format
func flattenNodes(d *schema.ResourceData, nodes []goe2e.ScalerGroupNode) error {
	if nodes == nil {
		return d.Set("nodes", []map[string]interface{}{})
	}

	nodeList := make([]map[string]interface{}, len(nodes))
	for i, node := range nodes {
		nodeList[i] = map[string]interface{}{
			"id":        node.ID,
			"name":      node.Name,
			"ip":        node.IP,
			"public_ip": node.PublicIP,
			"status":    node.Status,
			"cpu_usage": node.RealCPU,
		}
	}

	return d.Set("nodes", nodeList)
}

// expandScalingPolicy converts schema data to GoE2E ElasticPolicy slice
// Handles both V2 (policy) and V3 (scaling_policy) formats
func expandScalingPolicy(d *schema.ResourceData) []goe2e.ElasticPolicy {
	var policies []goe2e.ElasticPolicy

	if v, ok := d.GetOk("scaling_policy"); ok {
		// V3 structured format
		for _, p := range v.([]interface{}) {
			pMap := p.(map[string]interface{})
			policyType := "upscale"
			if pMap["type"].(string) == "scale_down" {
				policyType = "downscale"
			}

			// Map V3 metric names to API parameter names
			metric := pMap["metric"].(string)
			parameter := metric
			if metric == "cpu_utilization" {
				parameter = "cpu"
			} else if metric == "memory_utilization" {
				parameter = "memory"
			}

			policies = append(policies, goe2e.ElasticPolicy{
				Type:          policyType,
				Adjust:        pMap["adjustment"].(int),
				Parameter:     parameter,
				Operator:      pMap["operator"].(string),
				Value:         pMap["threshold"].(string),
				PeriodNumber:  strconv.Itoa(pMap["evaluation_periods"].(int)),
				PeriodSeconds: strconv.Itoa(pMap["period_seconds"].(int)),
				Cooldown:      strconv.Itoa(pMap["cooldown_seconds"].(int)),
			})
		}
	} else if v, ok := d.GetOk("policy"); ok {
		// V2 format
		for _, p := range v.([]interface{}) {
			pMap := p.(map[string]interface{})
			policies = append(policies, goe2e.ElasticPolicy{
				Type:          pMap["type"].(string),
				Adjust:        pMap["adjust"].(int),
				Parameter:     pMap["parameter"].(string),
				Operator:      pMap["operator"].(string),
				Value:         pMap["value"].(string),
				PeriodNumber:  pMap["period_number"].(string),
				PeriodSeconds: pMap["period_seconds"].(string),
				Cooldown:      pMap["cooldown"].(string),
			})
		}
	}

	return policies
}

// flattenScalingPolicy converts API ElasticPolicy data to schema format
// Returns both V2 and V3 formats
func flattenScalingPolicy(group *goe2e.ScalerGroup) ([]map[string]interface{}, []map[string]interface{}) {
	// V2 format (for backwards compatibility)
	v2Policies := []map[string]interface{}{
		{
			"type":           group.PolicyType,
			"adjust":         group.UpscalePolicyValue,
			"parameter":      group.PolicyMeasure,
			"operator":       group.PolicyUpscaleOperator,
			"value":          strconv.Itoa(group.UpscalePolicyValue),
			"period_number":  strconv.Itoa(group.WaitPeriod),
			"period_seconds": strconv.Itoa(group.Cooldown),
			"cooldown":       strconv.Itoa(group.Cooldown),
		},
		{
			"type":           group.PolicyType,
			"adjust":         group.DownscalePolicyValue,
			"parameter":      group.PolicyMeasure,
			"operator":       group.PolicyDownscaleOperator,
			"value":          strconv.Itoa(group.DownscalePolicyValue),
			"period_number":  strconv.Itoa(group.WaitPeriod),
			"period_seconds": strconv.Itoa(group.Cooldown),
			"cooldown":       strconv.Itoa(group.Cooldown),
		},
	}

	// V3 format (convert API format to V3)
	v3Policies := []map[string]interface{}{}
	if group.PolicyType != "" {
		// Map parameter back to V3 metric name
		metric := group.PolicyMeasure
		if metric == "cpu" {
			metric = "cpu_utilization"
		} else if metric == "memory" {
			metric = "memory_utilization"
		}

		// Upscale policy
		v3Policies = append(v3Policies, map[string]interface{}{
			"type":               "scale_up",
			"adjustment":         group.UpscalePolicyValue,
			"metric":             metric,
			"operator":           group.PolicyUpscaleOperator,
			"threshold":          strconv.Itoa(group.UpscalePolicyValue),
			"evaluation_periods": group.WaitPeriod,
			"period_seconds":     group.Cooldown,
			"cooldown_seconds":   group.Cooldown,
		})

		// Downscale policy
		v3Policies = append(v3Policies, map[string]interface{}{
			"type":               "scale_down",
			"adjustment":         group.DownscalePolicyValue,
			"metric":             metric,
			"operator":           group.PolicyDownscaleOperator,
			"threshold":          strconv.Itoa(group.DownscalePolicyValue),
			"evaluation_periods": group.WaitPeriod,
			"period_seconds":     group.Cooldown,
			"cooldown_seconds":   group.Cooldown,
		})
	}

	return v2Policies, v3Policies
}

// expandScheduledAction converts schema data to GoE2E ScheduledPolicy slice
// Handles both V2 (scheduled_policy) and V3 (scheduled_action) formats
func expandScheduledAction(d *schema.ResourceData) []goe2e.ScheduledPolicy {
	var policies []goe2e.ScheduledPolicy

	if v, ok := d.GetOk("scheduled_action"); ok {
		// V3 structured format
		for _, s := range v.([]interface{}) {
			sMap := s.(map[string]interface{})
			actionType := sMap["action_type"].(string)
			var adjust string
			if actionType == "set_capacity" {
				adjust = strconv.Itoa(sMap["target_capacity"].(int))
			} else {
				adjust = strconv.Itoa(sMap["adjustment"].(int))
			}
			policies = append(policies, goe2e.ScheduledPolicy{
				Type:       actionType,
				Adjust:     adjust,
				Recurrence: sMap["recurrence"].(string),
			})
		}
	} else if v, ok := d.GetOk("scheduled_policy"); ok {
		// V2 format
		for _, s := range v.([]interface{}) {
			sMap := s.(map[string]interface{})
			policies = append(policies, goe2e.ScheduledPolicy{
				Type:       sMap["type"].(string),
				Adjust:     sMap["adjust"].(string),
				Recurrence: sMap["recurrence"].(string),
			})
		}
	}

	return policies
}

// flattenScheduledAction converts API ScheduledPolicy data to schema format
// Returns both V2 and V3 formats
func flattenScheduledAction(group *goe2e.ScalerGroup) ([]map[string]interface{}, []map[string]interface{}) {
	// V2 format
	v2Policies := []map[string]interface{}{
		{
			"type":       group.ScheduledPolicyOp,
			"adjust":     strconv.Itoa(group.UpscaleAdjust),
			"recurrence": group.UpscaleRecurrence,
		},
		{
			"type":       group.ScheduledPolicyOp,
			"adjust":     strconv.Itoa(group.DownscaleAdjust),
			"recurrence": group.DownscaleRecurrence,
		},
	}

	// V3 format
	v3Policies := []map[string]interface{}{}
	if group.ScheduledPolicyOp != "" {
		// Convert to V3 format
		v3Policies = append(v3Policies, map[string]interface{}{
			"name":        "upscale",
			"action_type": "scale_up",
			"adjustment":  group.UpscaleAdjust,
			"recurrence":  group.UpscaleRecurrence,
		})
		v3Policies = append(v3Policies, map[string]interface{}{
			"name":        "downscale",
			"action_type": "scale_down",
			"adjustment":  group.DownscaleAdjust,
			"recurrence":  group.DownscaleRecurrence,
		})
	}

	return v2Policies, v3Policies
}

// expandVPCConfig converts schema VPC data to GoE2E VPCDetail slice
// Handles both V2 (vpc) and V3 (vpc_config) formats
func expandVPCConfig(ctx context.Context, d *schema.ResourceData, client *goe2e.Client) ([]goe2e.VPCDetail, error) {
	var vpcDetails []goe2e.VPCDetail

	if v, ok := d.GetOk("vpc_config"); ok {
		// V3 structured format
		for _, vRaw := range v.([]interface{}) {
			vMap := vRaw.(map[string]interface{})
			vpcName := vMap["name"].(string)

			// Get VPC details using GoE2E client
			vpcDetail, _, err := client.Vpcs.GetVPCByName(ctx, vpcName)
			if err != nil {
				return nil, fmt.Errorf("failed to get VPC details for %s: %w", vpcName, err)
			}

			vpcDetails = append(vpcDetails, goe2e.VPCDetail{
				Name:      vpcDetail.Name,
				NetworkID: vpcDetail.NetworkID,
				IPv4CIDR:  vpcDetail.IPv4CIDR,
				State:     vpcDetail.State,
			})
		}
	} else if v, ok := d.GetOk("vpc"); ok {
		// V2 format
		for _, vRaw := range v.([]interface{}) {
			vMap := vRaw.(map[string]interface{})
			vpcName := vMap["name"].(string)

			// Get VPC details using GoE2E client
			vpcDetail, _, err := client.Vpcs.GetVPCByName(ctx, vpcName)
			if err != nil {
				return nil, fmt.Errorf("failed to get VPC details for %s: %w", vpcName, err)
			}

			vpcDetails = append(vpcDetails, goe2e.VPCDetail{
				Name:      vpcDetail.Name,
				NetworkID: vpcDetail.NetworkID,
				IPv4CIDR:  vpcDetail.IPv4CIDR,
				State:     vpcDetail.State,
			})
		}
	}

	return vpcDetails, nil
}

// flattenVPCConfig converts API VPC data to schema format
// Returns format compatible with both V2 (vpc) and V3 (vpc_config) blocks
func flattenVPCConfig(ctx context.Context, vpcPartials []goe2e.VPCPartial, client *goe2e.Client) ([]map[string]interface{}, error) {
	var vpcList []map[string]interface{}

	for _, vpcPartial := range vpcPartials {
		vpcName := vpcPartial.Name

		// Get full VPC details with subnets using GoE2E client
		vpcDetail, _, err := client.Vpcs.GetVPCByName(ctx, vpcName)
		if err != nil {
			// Fallback: use partial info
			vpcList = append(vpcList, map[string]interface{}{
				"name":       vpcPartial.Name,
				"network_id": vpcPartial.NetworkID,
				"ipv4_cidr":  vpcPartial.IPv4CIDR,
				"state":      "",
				"subnets":    []map[string]interface{}{},
			})
			continue
		}

		vpcEntry := map[string]interface{}{
			"name":       vpcDetail.Name,
			"network_id": vpcDetail.NetworkID,
			"ipv4_cidr":  vpcDetail.IPv4CIDR,
			"state":      vpcDetail.State,
		}

		var subnets []map[string]interface{}
		for _, s := range vpcDetail.Subnets {
			subnets = append(subnets, map[string]interface{}{
				"id":          s.ID,
				"subnet_name": s.SubnetName,
				"cidr":        s.CIDR,
				"used_ips":    s.UsedIPs,
				"total_ips":   s.TotalIPs,
			})
		}

		vpcEntry["subnets"] = subnets
		vpcList = append(vpcList, vpcEntry)
	}

	return vpcList, nil
}

// NetworkConfig represents the consolidated network configuration
// Exported for testing purposes
type NetworkConfig struct {
	AssignPublicIP bool
	VPCNames       []string
	SecurityGroups []int
}

// ExpandNetworkConfig extracts network configuration from the network_config block
// Returns nil if the block is not present (caller should use individual fields)
// Exported for testing purposes
func ExpandNetworkConfig(d *schema.ResourceData) *NetworkConfig {
	return expandNetworkConfig(d)
}

// expandNetworkConfig extracts network configuration from the network_config block
// Returns nil if the block is not present (caller should use individual fields)
func expandNetworkConfig(d *schema.ResourceData) *NetworkConfig {
	// Try GetOk first - this checks if the value is "set" (non-zero)
	if v, ok := d.GetOk("network_config"); ok {
		configList := v.([]interface{})
		if len(configList) == 0 {
			return nil
		}
		return expandNetworkConfigFromList(configList)
	}

	// If GetOk returns false, try Get to handle empty maps in lists
	// This handles the edge case where a list exists but contains an empty map
	v := d.Get("network_config")
	if v == nil {
		return nil
	}

	configList, ok := v.([]interface{})
	if !ok || len(configList) == 0 {
		return nil
	}

	return expandNetworkConfigFromList(configList)
}

// expandNetworkConfigFromList extracts network config from a list of config maps
func expandNetworkConfigFromList(configList []interface{}) *NetworkConfig {

	// Check if first element is nil
	if configList[0] == nil {
		return nil
	}

	configMap, ok := configList[0].(map[string]interface{})
	if !ok {
		return nil
	}

	networkConfig := &NetworkConfig{}

	// Extract assign_public_ip
	if val, ok := configMap["assign_public_ip"]; ok {
		networkConfig.AssignPublicIP = val.(bool)
	} else {
		// Default to true if not specified (matching schema default)
		networkConfig.AssignPublicIP = true
	}

	// Extract vpc_names
	if val, ok := configMap["vpc_names"]; ok {
		vpcNamesRaw := val.([]interface{})
		vpcNames := make([]string, len(vpcNamesRaw))
		for i, v := range vpcNamesRaw {
			vpcNames[i] = v.(string)
		}
		networkConfig.VPCNames = vpcNames
	}

	// Extract security_groups
	if val, ok := configMap["security_groups"]; ok {
		sgRaw := val.([]interface{})
		securityGroups := make([]int, len(sgRaw))
		for i, v := range sgRaw {
			securityGroups[i] = v.(int)
		}
		networkConfig.SecurityGroups = securityGroups
	}

	return networkConfig
}

// FlattenNetworkConfig builds a network_config block from API response data
// Exported for testing purposes
func FlattenNetworkConfig(
	assignPublicIP bool,
	vpcNames []string,
	securityGroupIDs []int,
) []map[string]interface{} {
	return flattenNetworkConfig(assignPublicIP, vpcNames, securityGroupIDs)
}

// flattenNetworkConfig builds a network_config block from API response data
func flattenNetworkConfig(
	assignPublicIP bool,
	vpcNames []string,
	securityGroupIDs []int,
) []map[string]interface{} {
	// Only return a block if at least one field has a value
	if !assignPublicIP && len(vpcNames) == 0 && len(securityGroupIDs) == 0 {
		return []map[string]interface{}{}
	}

	config := map[string]interface{}{
		"assign_public_ip": assignPublicIP,
	}

	// Convert vpc_names to []interface{}
	if len(vpcNames) > 0 {
		vpcNamesList := make([]interface{}, len(vpcNames))
		for i, v := range vpcNames {
			vpcNamesList[i] = v
		}
		config["vpc_names"] = vpcNamesList
	}

	// Convert security_groups to []interface{}
	if len(securityGroupIDs) > 0 {
		sgList := make([]interface{}, len(securityGroupIDs))
		for i, v := range securityGroupIDs {
			sgList[i] = v
		}
		config["security_groups"] = sgList
	}

	return []map[string]interface{}{config}
}

// expandNetworkConfigFromRaw extracts network configuration from raw schema data
// Used in Update operations to compare old vs new values
func expandNetworkConfigFromRaw(raw interface{}) *NetworkConfig {
	if raw == nil {
		return nil
	}

	configList, ok := raw.([]interface{})
	if !ok || len(configList) == 0 {
		return nil
	}

	configMap := configList[0].(map[string]interface{})
	networkConfig := &NetworkConfig{}

	// Extract assign_public_ip
	if val, ok := configMap["assign_public_ip"]; ok {
		networkConfig.AssignPublicIP = val.(bool)
	} else {
		networkConfig.AssignPublicIP = true // default
	}

	// Extract vpc_names
	if val, ok := configMap["vpc_names"]; ok {
		vpcNamesRaw := val.([]interface{})
		vpcNames := make([]string, len(vpcNamesRaw))
		for i, v := range vpcNamesRaw {
			vpcNames[i] = v.(string)
		}
		networkConfig.VPCNames = vpcNames
	}

	// Extract security_groups
	if val, ok := configMap["security_groups"]; ok {
		sgRaw := val.([]interface{})
		securityGroups := make([]int, len(sgRaw))
		for i, v := range sgRaw {
			securityGroups[i] = v.(int)
		}
		networkConfig.SecurityGroups = securityGroups
	}

	return networkConfig
}

// StringSlicesEqual compares two string slices for equality
// Exported for testing purposes
func StringSlicesEqual(a, b []string) bool {
	return stringSlicesEqual(a, b)
}

// stringSlicesEqual compares two string slices for equality
func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// IntSlicesEqual compares two int slices for equality
// Exported for testing purposes
func IntSlicesEqual(a, b []int) bool {
	return intSlicesEqual(a, b)
}

// intSlicesEqual compares two int slices for equality
func intSlicesEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
