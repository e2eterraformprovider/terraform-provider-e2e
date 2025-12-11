package kubernetes

import (
	"context"
	"fmt"
	"log"
	"time"

	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ExpandNodePools(ctx context.Context, config []interface{}, goe2eClient *goe2e.Client, projectIDStr, region string) ([]goe2e.NodePool, error) {
	nodePools := make([]goe2e.NodePool, 0, len(config))
	uniqueNodePoolNames := make(map[string]bool)

	// Get worker plans from goe2e client (once for all node pools)
	workerPlans, _, err := goe2eClient.Kubernetes.GetWorkerPlans(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching worker plans: %w", err)
	}

	for _, np := range config {
		nodePoolDetail := np.(map[string]interface{})
		name := nodePoolDetail["name"].(string)
		uniqueNodePoolNames[name] = true

		// Use field alias helpers
		plan := getNodePoolPlan(nodePoolDetail)
		poolType := getNodePoolType(nodePoolDetail)

		// Log deprecation warnings
		if _, ok := nodePoolDetail["specs_name"]; ok {
			log.Printf("[WARN] Field 'specs_name' is deprecated. Use 'plan' instead.")
		}
		if _, ok := nodePoolDetail["node_pool_type"]; ok {
			log.Printf("[WARN] Field 'node_pool_type' is deprecated. Use 'type' instead.")
		}

		// Find matching plan
		var matchingPlan *goe2e.KubernetesWorkerPlan
		for i := range workerPlans {
			if workerPlans[i].Specs.SKUName == plan {
				matchingPlan = &workerPlans[i]
				break
			}
		}
		if matchingPlan == nil {
			return nil, fmt.Errorf("no matching plan found for plan: %s", plan)
		}
		if poolType == "" {
			return nil, fmt.Errorf("node pool type (type or node_pool_type) is required")
		}

		var policyType, customParamName, customParamValue string
		elasticityDict := goe2e.ElasticityDict{}
		scheduledDict := goe2e.ScheduledDict{}

		// If node_pool_type is Static, omit policyType, customParamName, and customParamValue
		if poolType == "Static" {
			policyType = ""
			customParamName = ""
			customParamValue = ""
		} else {
			policyType = getPolicyType(nodePoolDetail)
			customParamName = getCustomParamName(nodePoolDetail)
			customParamValue = getCustomParamValue(nodePoolDetail)

			minNodes := getNodePoolMinNodes(nodePoolDetail)
			maxNodes := getNodePoolMaxNodes(nodePoolDetail)

			if minNodes == 0 {
				return nil, fmt.Errorf("in case of Autoscale node type, the 'min_nodes' (or 'min_vms') field is required")
			}
			if maxNodes == 0 {
				return nil, fmt.Errorf("in case of Autoscale node type, the 'max_nodes' (or 'max_vms') field is required")
			}

			// Log deprecation warnings
			if _, ok := nodePoolDetail["min_vms"]; ok {
				log.Printf("[WARN] Field 'min_vms' is deprecated. Use 'min_nodes' instead.")
			}
			if _, ok := nodePoolDetail["max_vms"]; ok {
				log.Printf("[WARN] Field 'max_vms' is deprecated. Use 'max_nodes' instead.")
			}
			if _, ok := nodePoolDetail["worker_node"]; ok {
				log.Printf("[WARN] Field 'worker_node' is deprecated. Use 'size' instead.")
			}

			nodePoolDetail["cardinality"] = minNodes
			elasticity_dict, err := getElasticityDict(nodePoolDetail, minNodes, maxNodes)
			if err != nil {
				log.Printf("Invalid format for Elast")
			}

			scheduled_dict, err := getScheduledDict(nodePoolDetail, minNodes, maxNodes)
			if err != nil {
				log.Printf("Invalid format for Scheduled Dictionary")
			}
			elasticityDict = elasticity_dict
			scheduledDict = scheduled_dict
		}

		// Get size for Static pools
		size := getNodePoolSize(nodePoolDetail)
		if poolType == "Static" && size == 0 {
			return nil, fmt.Errorf("size (or worker_node) is required for Static node pools")
		}

		nodePool := goe2e.NodePool{
			Name:             name,
			SlugName:         matchingPlan.Plan,
			SKUID:            matchingPlan.Specs.ID,
			SpecsName:        plan,
			WorkerNode:       size,
			ElasticityDict:   elasticityDict,
			ScheduledDict:    scheduledDict,
			PolicyType:       policyType,
			CustomParamName:  customParamName,
			CustomParamValue: customParamValue,
		}

		nodePools = append(nodePools, nodePool)
	}
	numUniqueNodePools := len(uniqueNodePoolNames)
	if numUniqueNodePools < len(config) {
		return []goe2e.NodePool{}, fmt.Errorf("Name of the worker node pools must be unique!")
	}
	return nodePools, nil
}

// ExpandElasticityDict is a helper function to process the elasticity_dict attribute.
func ExpandElasticityDict(config map[string]interface{}, min_vms int, max_vms int) (goe2e.ElasticityDict, error) {
	elasticityDict := goe2e.ElasticityDict{}
	workers := config["worker"].([]interface{})
	if len(workers) > 0 {
		worker := workers[0].(map[string]interface{})
		elasticityWorker, err := ExpandElasticityWorker(worker, min_vms, max_vms)
		if err != nil {
			return goe2e.ElasticityDict{}, err
		}

		elasticityDict = goe2e.ElasticityDict{
			Worker: elasticityWorker,
		}
		return elasticityDict, nil
	}
	return elasticityDict, nil
}

func ExpandScheduledDict(config map[string]interface{}, min_vms int, max_vms int) (goe2e.ScheduledDict, error) {
	scheduledDict := goe2e.ScheduledDict{}
	workers := config["worker"].([]interface{})
	if len(workers) > 0 {
		worker := workers[0].(map[string]interface{})
		scheduledWorker, err := ExpandScheduledWorker(worker, min_vms, max_vms)
		if err != nil {
			return goe2e.ScheduledDict{}, err
		}

		scheduledDict = goe2e.ScheduledDict{
			Worker: scheduledWorker,
		}
		return scheduledDict, nil
	}
	return scheduledDict, nil
}

// ExpandElasticityWorker is a helper function to process the worker attribute in elasticity_dict.
func ExpandElasticityWorker(config map[string]interface{}, min_vms int, max_vms int) (goe2e.ElasticityWorker, error) {
	elasticityPolicies, err := ExpandElasticityPolicies(config["elasticity_policies"].([]interface{}), config["parameter"].(string))
	if err != nil {
		return goe2e.ElasticityWorker{}, err
	}

	return goe2e.ElasticityWorker{
		MinVms:             min_vms,
		Cardinality:        min_vms,
		MaxVms:             max_vms,
		ElasticityPolicies: elasticityPolicies,
	}, nil
}

func ExpandScheduledWorker(config map[string]interface{}, min_vms int, max_vms int) (goe2e.ScheduleWorker, error) {
	scheduledPolicies, err := ExpandScheduledPolicies(config["scheduled_policies"].([]interface{}), min_vms, max_vms)
	if err != nil {
		return goe2e.ScheduleWorker{}, err
	}

	return goe2e.ScheduleWorker{
		MinVms:            min_vms,
		Cardinality:       min_vms,
		MaxVms:            max_vms,
		ScheduledPolicies: scheduledPolicies,
	}, nil
}

// ExpandElasticityPolicies is a helper function to process the elasticity_policies attribute.
func ExpandElasticityPolicies(config []interface{}, parameter string) ([]goe2e.ElasticityPolicy, error) {
	elasticityPolicies := make([]goe2e.ElasticityPolicy, 0, len(config))
	var adjust_value int = -1
	type_value := "CHANGE"
	for _, ep := range config {
		adjust_value = -1 * adjust_value
		elasticityPolicyDetail := ep.(map[string]interface{})
		elasticityPolicy := goe2e.ElasticityPolicy{
			Type:         type_value,
			Adjust:       adjust_value,
			Parameter:    parameter,
			Operator:     elasticityPolicyDetail["operator"].(string),
			Value:        elasticityPolicyDetail["value"].(int),
			PeriodNumber: elasticityPolicyDetail["watch_period"].(int),
			Period:       elasticityPolicyDetail["period"].(int),
			Cooldown:     elasticityPolicyDetail["cooldown"].(int),
		}
		elasticityPolicies = append(elasticityPolicies, elasticityPolicy)
	}
	return elasticityPolicies, nil
}

func ExpandScheduledPolicies(config []interface{}, min_vms int, max_vms int) ([]goe2e.SchedulePolicy, error) {
	scheduledPolicies := make([]goe2e.SchedulePolicy, 0, len(config))
	for _, sp := range config {
		elasticityPolicyDetail := sp.(map[string]interface{})
		// Adjust should be converted to an integer
		upscaleCardinality := elasticityPolicyDetail["upscale_cardinality"].(int)
		downscaleCardinality := elasticityPolicyDetail["downscale_cardinality"].(int)
		upscaleRecurrence := elasticityPolicyDetail["upscale_recurrence"].(string)
		downscaleRecurrence := elasticityPolicyDetail["downscale_recurrence"].(string)

		if upscaleCardinality < min_vms || upscaleCardinality > max_vms {
			return scheduledPolicies, fmt.Errorf("upscale cardinality must be between min nodes and max nodes")
		} else if downscaleCardinality < min_vms || downscaleCardinality > max_vms {
			return scheduledPolicies, fmt.Errorf("downscale cardinality must be between min nodes and max nodes")
		}

		// Create SchedulePolicy instances
		upscalePolicy := goe2e.SchedulePolicy{
			Type:       "CARDINALITY",
			Adjust:     upscaleCardinality,
			Recurrence: upscaleRecurrence,
		}
		downscalePolicy := goe2e.SchedulePolicy{
			Type:       "CARDINALITY",
			Adjust:     downscaleCardinality,
			Recurrence: downscaleRecurrence,
		}
		scheduledPolicies = append(scheduledPolicies, upscalePolicy, downscalePolicy)
	}
	return scheduledPolicies, nil
}

func getElasticityDict(nodePoolDetail map[string]interface{}, min_vms int, max_vms int) (goe2e.ElasticityDict, error) {
	var elasticityDict goe2e.ElasticityDict

	// Handle ElasticityDict based on node_pool_type
	switch nodePoolType := nodePoolDetail["node_pool_type"].(string); nodePoolType {
	case "Static":
		elasticityDict = goe2e.ElasticityDict{}
	case "Autoscale":
		for _, ed := range nodePoolDetail["elasticity_dict"].([]interface{}) {
			ed := ed.(map[string]interface{})
			elasticityDict, _ := ExpandElasticityDict(ed, min_vms, max_vms)
			return elasticityDict, nil
		}
	default:
		return elasticityDict, fmt.Errorf("invalid node_pool_type: %s", nodePoolType)
	}
	return elasticityDict, nil
}

func getScheduledDict(nodePoolDetail map[string]interface{}, min_vms int, max_vms int) (goe2e.ScheduledDict, error) {
	var scheduledDict goe2e.ScheduledDict

	switch nodePoolType := nodePoolDetail["node_pool_type"].(string); nodePoolType {
	case "Static":
		scheduledDict = goe2e.ScheduledDict{}
	case "Autoscale":
		for _, sd := range nodePoolDetail["scheduled_dict"].([]interface{}) {
			sd := sd.(map[string]interface{})
			scheduledDict, _ := ExpandScheduledDict(sd, min_vms, max_vms)
			return scheduledDict, nil
		}
	default:
		return scheduledDict, fmt.Errorf("invalid node_pool_type: %s", nodePoolType)
	}

	return scheduledDict, nil
}

func getCustomParamName(nodePoolDetail map[string]interface{}) string {
	if nodePoolType, ok := nodePoolDetail["node_pool_type"].(string); ok && nodePoolType == "Static" {
		return "" // Return empty string for "Static"
	}
	policyParameterType := getPolicyType(nodePoolDetail)
	if policyParameterType == "" || policyParameterType == "Default" {
		return "" // Return empty string when policy_parameter_type is not provided or is "Default"
	}
	elasticityDict, ok := nodePoolDetail["elasticity_dict"].([]interface{})
	if !ok {
		log.Printf("Elasticity dictionary not found or not in the expected format")
		return ""
	}
	for _, ed := range elasticityDict {
		edMap, ok := ed.(map[string]interface{})
		if !ok {
			log.Printf("Elasticity dictionary entry is not in the expected format")
			continue
		}
		workerList, ok := edMap["worker"].([]interface{})
		if !ok || len(workerList) == 0 {
			log.Printf("Worker list not found or empty")
			continue
		}
		// Assuming there is only one worker map in the list
		worker, ok := workerList[0].(map[string]interface{})
		if !ok {
			log.Printf("Worker map is not in the expected format")
			continue
		}
		parameter, ok := worker["parameter"].(string)
		if !ok {
			log.Printf("Parameter field not found or not a string")
			continue
		}
		// Check if "parameter" is "CPU" or "Memory"
		if parameter == "CPU" || parameter == "Memory" {
			log.Printf("Cannot use Default parameters in case of Custom")
			return ""
		}

		return parameter
	}
	return ""
}

func getCustomParamValue(nodePoolDetail map[string]interface{}) string {
	if nodePoolType, ok := nodePoolDetail["node_pool_type"].(string); ok && nodePoolType == "Static" {
		return "" // Return empty string for "Static"
	}
	policyParameterType := getPolicyType(nodePoolDetail)
	if policyParameterType == "" || policyParameterType == "Default" || policyParameterType == "Scheduled" {
		return "" // Return empty string when policy_parameter_type is not provided or is "Default"
	}
	return "0"
}

func getPolicyType(nodePoolDetail map[string]interface{}) string {
	elasticityDict, _ := nodePoolDetail["elasticity_dict"].([]interface{})
	scheduledDict, scheduledDictPresent := nodePoolDetail["scheduled_dict"].([]interface{})
	log.Printf("------ScheduledDict: %+v------ElasticityDict: %+v", scheduledDict, elasticityDict)
	if len(elasticityDict) == 0 && len(scheduledDict) == 0 {
		return ""
	}
	isSDPresent := true
	if len(scheduledDict) == 0 {
		isSDPresent = false
	}
	for _, ed := range elasticityDict {
		edMap, ok := ed.(map[string]interface{})
		if !ok {
			log.Printf("Elasticity dictionary entry is not in the expected format")
			continue
		}
		workerList, ok := edMap["worker"].([]interface{})
		if !ok || len(workerList) == 0 {
			log.Printf("Worker list not found or empty")
			continue
		}
		worker, ok := workerList[0].(map[string]interface{})
		if !ok {
			log.Printf("Worker map is not in the expected format")
			continue
		}
		policyParamType, ok := worker["policy_paramter_type"].(string)
		if !ok {
			log.Printf("Policy parameter type not found or not a string")
			continue
		}
		if scheduledDictPresent && isSDPresent {
			return policyParamType + "-Scheduled"
		}
		return policyParamType
	}
	if (len(elasticityDict) == 0) && isSDPresent {
		return "Scheduled"
	}

	return ""
}

func ExpandNodePoolUpdate(ctx context.Context, nodePoolDetail map[string]interface{}, goe2eClient *goe2e.Client, projectIDStr, region string) (goe2e.NodePoolUpdate, error) {
	nodeUpdate := goe2e.NodePoolUpdate{}
	poolType := getNodePoolType(nodePoolDetail)
	if poolType == "" {
		return nodeUpdate, fmt.Errorf("node pool type (type or node_pool_type) is required")
	}

	var policyType, customParamName, customParamValue string
	var elasticity_policies []goe2e.ElasticityPolicy
	var scheduled_policies []goe2e.SchedulePolicy
	var minNodes, maxNodes int

	// Get worker plans from goe2e client
	workerPlans, _, err := goe2eClient.Kubernetes.GetWorkerPlans(ctx)
	if err != nil {
		return nodeUpdate, fmt.Errorf("error fetching worker plans: %w", err)
	}

	// Find matching plan
	plan := getNodePoolPlan(nodePoolDetail)
	var matchingPlan *goe2e.KubernetesWorkerPlan
	for i := range workerPlans {
		if workerPlans[i].Specs.SKUName == plan {
			matchingPlan = &workerPlans[i]
			break
		}
	}
	if matchingPlan == nil {
		return nodeUpdate, fmt.Errorf("no matching plan found for plan: %s", plan)
	}

	if poolType == "Static" {
		policyType = ""
		customParamName = ""
		customParamValue = ""
		minNodes = 0
		maxNodes = 0
	} else {
		policyType = getPolicyType(nodePoolDetail)
		customParamName = getCustomParamName(nodePoolDetail)
		customParamValue = getCustomParamValue(nodePoolDetail)
		minNodes = getNodePoolMinNodes(nodePoolDetail)
		maxNodes = getNodePoolMaxNodes(nodePoolDetail)

		if minNodes == 0 {
			return nodeUpdate, fmt.Errorf("in case of Autoscale node type, the 'min_nodes' (or 'min_vms') field is required")
		}
		if maxNodes == 0 {
			return nodeUpdate, fmt.Errorf("in case of Autoscale node type, the 'max_nodes' (or 'max_vms') field is required")
		}

		ep, err := updateElasticPolicies(nodePoolDetail, minNodes, maxNodes)
		if err != nil {
			log.Printf("Invalid format for Elast")
		}

		sp, err := updateScheduledPolicies(nodePoolDetail, minNodes, maxNodes)
		if err != nil {
			log.Printf("Invalid format for Scheduled Dictionary")
		}
		elasticity_policies = ep
		scheduled_policies = sp

		cardinterface := nodePoolDetail["cardinality"]
		card := 0
		if cardinterface != nil {
			card = cardinterface.(int)
		}
		if card == 0 {
			nodePoolDetail["cardinality"] = minNodes
		}
	}
	cardinality := 0
	if cardinterface := nodePoolDetail["cardinality"]; cardinterface != nil {
		cardinality = cardinterface.(int)
	}
	if cardinality == 0 {
		cardinality = minNodes
	}

	nodeUpdate = goe2e.NodePoolUpdate{
		MinVms:           minNodes,
		Cardinality:      cardinality,
		MaxVms:           maxNodes,
		PlanID:           matchingPlan.Specs.ID,
		ElasticityPolicy: elasticity_policies,
		ScheduledPolicy:  scheduled_policies,
		PolicyType:       policyType,
		CustomParamName:  customParamName,
		CustomParamValue: customParamValue,
	}
	return nodeUpdate, nil
}

func updateElasticPolicies(nodePoolDetail map[string]interface{}, min_vms int, max_vms int) ([]goe2e.ElasticityPolicy, error) {
	var elasticityPolicyList []goe2e.ElasticityPolicy
	switch nodePoolType := nodePoolDetail["node_pool_type"].(string); nodePoolType {
	case "Static":
		elasticityPolicyList = []goe2e.ElasticityPolicy{}
	case "Autoscale":
		for _, ed := range nodePoolDetail["elasticity_dict"].([]interface{}) {
			ed := ed.(map[string]interface{})
			elasticityDict, _ := UpdateElasticityDict(ed, min_vms, max_vms)
			return elasticityDict, nil
		}
	default:
		return elasticityPolicyList, fmt.Errorf("invalid node_pool_type: %s", nodePoolType)
	}
	return elasticityPolicyList, nil
}

func UpdateElasticityDict(config map[string]interface{}, min_vms int, max_vms int) ([]goe2e.ElasticityPolicy, error) {
	elasticityPolicy := []goe2e.ElasticityPolicy{}
	workers := config["worker"].([]interface{})
	if len(workers) > 0 {
		worker := workers[0].(map[string]interface{})
		elasticityPolicy, err := UpdateElasticityWorker(worker, min_vms, max_vms)
		if err != nil {
			return []goe2e.ElasticityPolicy{}, err
		}
		return elasticityPolicy, nil
	}
	return elasticityPolicy, nil
}

func UpdateElasticityWorker(config map[string]interface{}, min_vms int, max_vms int) ([]goe2e.ElasticityPolicy, error) {
	elasticityPolicies, err := ExpandElasticityPolicies(config["elasticity_policies"].([]interface{}), config["parameter"].(string))
	if err != nil {
		ep := make([]goe2e.ElasticityPolicy, 0, len(config))
		return ep, err
	}

	return elasticityPolicies, nil
}

func updateScheduledPolicies(nodePoolDetail map[string]interface{}, min_vms int, max_vms int) ([]goe2e.SchedulePolicy, error) {
	var scheduledPolicyList []goe2e.SchedulePolicy

	switch nodePoolType := nodePoolDetail["node_pool_type"].(string); nodePoolType {
	case "Static":
		scheduledPolicyList = []goe2e.SchedulePolicy{}
	case "Autoscale":
		for _, sd := range nodePoolDetail["scheduled_dict"].([]interface{}) {
			sd := sd.(map[string]interface{})
			scheduledDict, _ := UpdateScheduledDict(sd, min_vms, max_vms)
			return scheduledDict, nil
		}
	default:
		return scheduledPolicyList, fmt.Errorf("invalid node_pool_type: %s", nodePoolType)
	}

	return scheduledPolicyList, nil
}

func UpdateScheduledDict(config map[string]interface{}, min_vms int, max_vms int) ([]goe2e.SchedulePolicy, error) {
	scheduledDict := []goe2e.SchedulePolicy{}
	workers := config["worker"].([]interface{})
	if len(workers) > 0 {
		worker := workers[0].(map[string]interface{})
		scheduledWorker, err := UpdateScheduledWorker(worker, min_vms, max_vms)
		if err != nil {
			return []goe2e.SchedulePolicy{}, err
		}
		return scheduledWorker, nil
	}
	return scheduledDict, nil
}

func UpdateScheduledWorker(config map[string]interface{}, min_vms int, max_vms int) ([]goe2e.SchedulePolicy, error) {
	scheduledPolicies, err := ExpandScheduledPolicies(config["scheduled_policies"].([]interface{}), min_vms, max_vms)
	if err != nil {
		return []goe2e.SchedulePolicy{}, err
	}

	return scheduledPolicies, nil
}

// ============================================
// V3 FIELD ALIASING HELPERS
// ============================================

// getClusterName returns cluster_name if set, otherwise name (deprecated)
func getClusterName(d *schema.ResourceData) string {
	if v, ok := d.GetOk("cluster_name"); ok {
		return v.(string)
	}
	if v, ok := d.GetOk(tfconstants.AttrName); ok {
		return v.(string)
	}
	return ""
}

// getKubernetesVersion returns kubernetes_version if set, otherwise version (deprecated)
func getKubernetesVersion(d *schema.ResourceData) string {
	if v, ok := d.GetOk("kubernetes_version"); ok {
		return v.(string)
	}
	if v, ok := d.GetOk(tfconstants.AttrVersion); ok {
		return v.(string)
	}
	return ""
}

// logDeprecationWarning logs a deprecation warning if a deprecated field is used
func logDeprecationWarning(d *schema.ResourceData, deprecatedField, preferredField string) {
	if _, ok := d.GetOk(deprecatedField); ok {
		log.Printf("[WARN] Field '%s' is deprecated. Use '%s' instead. The deprecated field will be removed in a future version.", deprecatedField, preferredField)
	}
}

// ============================================
// NODE POOL FIELD ALIASING HELPERS
// ============================================

// getNodePoolPlan returns plan if set, otherwise specs_name (deprecated)
func getNodePoolPlan(pool map[string]interface{}) string {
	if v, ok := pool["plan"].(string); ok && v != "" {
		return v
	}
	if v, ok := pool["specs_name"].(string); ok {
		return v
	}
	return ""
}

// getNodePoolType returns type if set, otherwise node_pool_type (deprecated)
func getNodePoolType(pool map[string]interface{}) string {
	if v, ok := pool["type"].(string); ok && v != "" {
		return v
	}
	if v, ok := pool["node_pool_type"].(string); ok {
		return v
	}
	return ""
}

// getNodePoolSize returns size if set, otherwise worker_node (deprecated)
func getNodePoolSize(pool map[string]interface{}) int {
	if v, ok := pool["size"].(int); ok && v > 0 {
		return v
	}
	if v, ok := pool["worker_node"].(int); ok {
		return v
	}
	return 0
}

// getNodePoolMinNodes returns min_nodes if set, otherwise min_vms (deprecated)
func getNodePoolMinNodes(pool map[string]interface{}) int {
	if v, ok := pool["min_nodes"].(int); ok && v > 0 {
		return v
	}
	if v, ok := pool[tfconstants.AttrMinVMs].(int); ok {
		return v
	}
	return 0
}

// getNodePoolMaxNodes returns max_nodes if set, otherwise max_vms (deprecated)
func getNodePoolMaxNodes(pool map[string]interface{}) int {
	if v, ok := pool["max_nodes"].(int); ok && v > 0 {
		return v
	}
	if v, ok := pool[tfconstants.AttrMaxVMs].(int); ok {
		return v
	}
	return 0
}

// ============================================
// FLATTEN HELPERS
// ============================================

// flattenNodePools converts goe2e NodePoolServiceInfo to Terraform schema
func flattenNodePools(nodePools []goe2e.NodePoolServiceInfo) []interface{} {
	result := make([]interface{}, len(nodePools))

	for i, pool := range nodePools {
		p := make(map[string]interface{})

		// Use V3 preferred field names
		p["name"] = pool.ServiceName
		p["service_id"] = fmt.Sprintf("%.0f", pool.ServiceID)
		p["cardinality"] = pool.Cardinality

		// Note: plan, type, size, min_nodes, max_nodes are not returned by GetNodePools
		// They will be preserved from state or set during read from full node pool details
		// This is a limitation - we'd need a more detailed API call to get full node pool info

		result[i] = p
	}

	return result
}

// ============================================
// ASYNC WAIT HELPERS
// ============================================

// waitForClusterStatus waits for cluster to reach target status
func waitForClusterStatus(ctx context.Context, client *goe2e.Client, clusterID string, targetStatus string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Creating", "Provisioning", "Updating"},
		Target:     []string{targetStatus},
		Refresh:    clusterStatusRefresh(ctx, client, clusterID),
		Timeout:    timeout,
		Delay:      30 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func clusterStatusRefresh(ctx context.Context, client *goe2e.Client, clusterID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		cluster, _, err := client.Kubernetes.Get(ctx, clusterID)
		if err != nil {
			return nil, "", err
		}

		if cluster == nil {
			return nil, "", fmt.Errorf("cluster returned nil")
		}

		return cluster, cluster.State, nil
	}
}

// ============================================
// SECURITY GROUP HELPERS
// ============================================

// expandSecurityGroupIDs converts []interface{} to []int
func expandSecurityGroupIDs(ids []interface{}) []int {
	result := make([]int, len(ids))
	for i, id := range ids {
		result[i] = id.(int)
	}
	return result
}

// difference returns elements in a but not in b
func difference(a, b []int) []int {
	bMap := make(map[int]bool)
	for _, v := range b {
		bMap[v] = true
	}

	var diff []int
	for _, v := range a {
		if !bMap[v] {
			diff = append(diff, v)
		}
	}
	return diff
}
