package kubernetes

import (
	"fmt"
	"log"
	"sort"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ExpandNodePools(config []interface{}, apiClient *client.Client, project_id int, location string) ([]models.NodePool, error) {
	nodePools := make([]models.NodePool, 0, len(config))
	uniqueNodePoolNames := make(map[string]bool)

	for _, np := range config {
		nodePoolDetail := np.(map[string]interface{})
		name := nodePoolDetail["name"].(string)
		uniqueNodePoolNames[name] = true
		workerPlans, err := apiClient.GetKubernetesWorkerPlans(project_id, location) //Here we are are tryig to get all worker plans
		if err != nil {
			return nil, err
		}
		plans := workerPlans["data"].([]interface{})
		var matchingPlan map[string]interface{}
		//Below code is to fetch the corresponding slug name
		for _, plan := range plans {
			planData := plan.(map[string]interface{})
			specs := planData["specs"].(map[string]interface{})
			skuName := specs["sku_name"].(string)

			if skuName == nodePoolDetail["specs_name"].(string) {
				matchingPlan = planData
				break
			}
		}
		if matchingPlan == nil {
			return nil, fmt.Errorf("no matching plan found for specs_name: %s", nodePoolDetail["specs_name"])
		}
		if _, ok := nodePoolDetail["node_pool_type"]; !ok {
			return nil, fmt.Errorf("node_pool_type is required")
		}
		var policyType, customParamName, customParamValue string
		elasticityDict := models.ElasticityDict{}
		scheduledDict := models.ScheduledDict{}

		// If node_pool_type is Static, omit policyType, customParamName, and customParamValue
		if nodePoolDetail["node_pool_type"].(string) == "Static" {
			policyType = ""
			customParamName = ""
			customParamValue = ""
		} else {
			policyType = getPolicyType(nodePoolDetail)
			customParamName = getCustomParamName(nodePoolDetail)
			customParamValue = getCustomParamValue(nodePoolDetail)

			if _, ok := nodePoolDetail["min_vms"]; !ok {
				return nil, fmt.Errorf("in case of Autoscale node type, the 'min_vms' field is required")
			}
			if _, ok := nodePoolDetail["max_vms"]; !ok {
				return nil, fmt.Errorf("in case of Autoscale node type, the 'max_vms' field is required")
			}
			nodePoolDetail["cardinality"] = nodePoolDetail["min_vms"].(int) //NEW CHANGE
			elasticity_dict, err := getElasticityDict(nodePoolDetail, nodePoolDetail["min_vms"].(int), nodePoolDetail["max_vms"].(int))
			if err != nil {
				log.Printf("Invalid format for Elast")
			}

			scheduled_dict, err := getScheduledDict(nodePoolDetail, nodePoolDetail["min_vms"].(int), nodePoolDetail["max_vms"].(int))
			if err != nil {
				log.Printf("Invalid format for Scheduled Dictionary")
			}
			elasticityDict = elasticity_dict
			scheduledDict = scheduled_dict
		}

		nodePool := models.NodePool{
			Name:             nodePoolDetail["name"].(string),
			SlugName:         matchingPlan["plan"].(string),
			SKUID:            matchingPlan["specs"].(map[string]interface{})["id"].(string),
			SpecsName:        nodePoolDetail["specs_name"].(string),
			WorkerNode:       nodePoolDetail["worker_node"].(int),
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
		return []models.NodePool{}, fmt.Errorf("Name of the worker node pools must be unique!")
	}
	return nodePools, nil
}

// ExpandElasticityDict is a helper function to process the elasticity_dict attribute.
func ExpandElasticityDict(config map[string]interface{}, min_vms int, max_vms int) (models.ElasticityDict, error) {
	elasticityDict := models.ElasticityDict{}
	for _, worker := range config["worker"].([]interface{}) {
		worker := worker.(map[string]interface{})
		elasticityWorker, err := ExpandElasticityWorker(worker, min_vms, max_vms)
		if err != nil {
			return models.ElasticityDict{}, err
		}

		elasticityDict = models.ElasticityDict{
			Worker: elasticityWorker,
		}
		return elasticityDict, nil
	}
	return elasticityDict, nil
}

func ExpandScheduledDict(config map[string]interface{}, min_vms int, max_vms int) (models.ScheduledDict, error) {
	scheduledDict := models.ScheduledDict{}
	for _, worker := range config["worker"].([]interface{}) {
		worker := worker.(map[string]interface{})
		scheduledWorker, err := ExpandScheduledWorker(worker, min_vms, max_vms)
		if err != nil {
			return models.ScheduledDict{}, err
		}

		scheduledDict = models.ScheduledDict{
			Worker: scheduledWorker,
		}
		return scheduledDict, nil
	}
	return scheduledDict, nil
}

// ExpandElasticityWorker is a helper function to process the worker attribute in elasticity_dict.
func ExpandElasticityWorker(config map[string]interface{}, min_vms int, max_vms int) (models.ElasticityWorker, error) {
	elasticityPolicies, err := ExpandElasticityPolicies(config["elasticity_policies"].([]interface{}), config["parameter"].(string))
	if err != nil {
		return models.ElasticityWorker{}, err
	}

	return models.ElasticityWorker{
		MinVms:             min_vms,
		Cardinality:        min_vms,
		MaxVms:             max_vms,
		ElasticityPolicies: elasticityPolicies,
	}, nil
}

func ExpandScheduledWorker(config map[string]interface{}, min_vms int, max_vms int) (models.ScheduleWorker, error) {
	scheduledPolicies, err := ExpandScheduledPolicies(config["scheduled_policies"].([]interface{}), min_vms, max_vms)
	if err != nil {
		return models.ScheduleWorker{}, err
	}

	return models.ScheduleWorker{
		MinVms:            min_vms,
		Cardinality:       min_vms,
		MaxVms:            max_vms,
		ScheduledPolicies: scheduledPolicies,
	}, nil
}

// ExpandElasticityPolicies is a helper function to process the elasticity_policies attribute.
func ExpandElasticityPolicies(config []interface{}, parameter string) ([]models.ElasticityPolicy, error) {
	elasticityPolicies := make([]models.ElasticityPolicy, 0, len(config))
	var adjust_value int = -1
	type_value := "CHANGE"
	for _, ep := range config {
		adjust_value = -1 * adjust_value
		elasticityPolicyDetail := ep.(map[string]interface{})
		elasticityPolicy := models.ElasticityPolicy{
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

func ExpandScheduledPolicies(config []interface{}, min_vms int, max_vms int) ([]models.SchedulePolicy, error) {
	scheduledPolicies := make([]models.SchedulePolicy, 0, len(config))
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
		upscalePolicy := models.SchedulePolicy{
			Type:       "CARDINALITY",
			Adjust:     upscaleCardinality,
			Recurrence: upscaleRecurrence,
		}
		downscalePolicy := models.SchedulePolicy{
			Type:       "CARDINALITY",
			Adjust:     downscaleCardinality,
			Recurrence: downscaleRecurrence,
		}
		scheduledPolicies = append(scheduledPolicies, upscalePolicy, downscalePolicy)
	}
	return scheduledPolicies, nil
}

func getElasticityDict(nodePoolDetail map[string]interface{}, min_vms int, max_vms int) (models.ElasticityDict, error) {
	var elasticityDict models.ElasticityDict

	// Handle ElasticityDict based on node_pool_type
	switch nodePoolType := nodePoolDetail["node_pool_type"].(string); nodePoolType {
	case "Static":
		elasticityDict = models.ElasticityDict{}
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

func getScheduledDict(nodePoolDetail map[string]interface{}, min_vms int, max_vms int) (models.ScheduledDict, error) {
	var scheduledDict models.ScheduledDict

	switch nodePoolType := nodePoolDetail["node_pool_type"].(string); nodePoolType {
	case "Static":
		scheduledDict = models.ScheduledDict{}
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

func ExpandNPUpdate(nodePoolDetail map[string]interface{}, apiClient *client.Client, project_id int, location string) (models.NodePoolUpdate, error) {
	nodeUpdate := models.NodePoolUpdate{}
	if _, ok := nodePoolDetail["node_pool_type"]; !ok {
		return nodeUpdate, fmt.Errorf("node_pool_type is required")
	}
	var policyType, customParamName, customParamValue string
	var elasticity_policies []models.ElasticityPolicy
	var scheduled_policies []models.SchedulePolicy
	// var card int
	workerPlans, err := apiClient.GetKubernetesWorkerPlans(project_id, location) //Here we are are tryig to get all worker plans
	if err != nil {
		return nodeUpdate, err
	}
	plans := workerPlans["data"].([]interface{})
	var matchingPlan map[string]interface{}
	//Below code is to fetch the corresponding slug name
	for _, plan := range plans {
		planData := plan.(map[string]interface{})
		specs := planData["specs"].(map[string]interface{})
		skuName := specs["sku_name"].(string)

		if skuName == nodePoolDetail["specs_name"].(string) {
			matchingPlan = planData
			break
		}
	}
	if matchingPlan == nil {
		return nodeUpdate, fmt.Errorf("no matching plan found for specs_name: %s", nodePoolDetail["specs_name"])
	}

	if nodePoolDetail["node_pool_type"].(string) == "Static" {
		policyType = ""
		customParamName = ""
		customParamValue = ""
	} else {
		policyType = getPolicyType(nodePoolDetail)
		customParamName = getCustomParamName(nodePoolDetail)
		customParamValue = getCustomParamValue(nodePoolDetail)
		if _, ok := nodePoolDetail["min_vms"]; !ok {
			return nodeUpdate, fmt.Errorf("in case of Autoscale node type, the 'min_vms' field is required")
		}
		if _, ok := nodePoolDetail["max_vms"]; !ok {
			return nodeUpdate, fmt.Errorf("in case of Autoscale node type, the 'max_vms' field is required")
		}
		ep, err := updateElasticPolicies(nodePoolDetail, nodePoolDetail["min_vms"].(int), nodePoolDetail["max_vms"].(int))
		if err != nil {
			log.Printf("Invalid format for Elast")
		}

		sp, err := updateScheduledPolicies(nodePoolDetail, nodePoolDetail["min_vms"].(int), nodePoolDetail["max_vms"].(int))
		if err != nil {
			log.Printf("Invalid format for Scheduled Dictionary")
		}
		elasticity_policies = ep
		scheduled_policies = sp

		cardinterface := nodePoolDetail["cardinality"]
		card := cardinterface.(int)
		if card == 0 || cardinterface == nil {
			nodePoolDetail["cardinality"] = nodePoolDetail["min_vms"].(int)
		}
	}
	nodeUpdate = models.NodePoolUpdate{
		MinVms:           nodePoolDetail["min_vms"].(int),
		Cardinality:      nodePoolDetail["cardinality"].(int),
		MaxVms:           nodePoolDetail["max_vms"].(int),
		PlanID:           matchingPlan["specs"].(map[string]interface{})["id"].(string),
		ElasticityPolicy: elasticity_policies,
		ScheduledPolicy:  scheduled_policies,
		PolicyType:       policyType,
		CustomParamName:  customParamName,
		CustomParamValue: customParamValue,
	}
	return nodeUpdate, nil
}

func updateElasticPolicies(nodePoolDetail map[string]interface{}, min_vms int, max_vms int) ([]models.ElasticityPolicy, error) {
	var elasticityPolicyList []models.ElasticityPolicy
	switch nodePoolType := nodePoolDetail["node_pool_type"].(string); nodePoolType {
	case "Static":
		elasticityPolicyList = []models.ElasticityPolicy{}
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

func UpdateElasticityDict(config map[string]interface{}, min_vms int, max_vms int) ([]models.ElasticityPolicy, error) {
	elasticityPolicy := []models.ElasticityPolicy{}
	for _, worker := range config["worker"].([]interface{}) {
		worker := worker.(map[string]interface{})
		elasticityPolicy, err := UpdateElasticityWorker(worker, min_vms, max_vms)
		if err != nil {
			return []models.ElasticityPolicy{}, err
		}
		return elasticityPolicy, nil
	}
	return elasticityPolicy, nil
}

func UpdateElasticityWorker(config map[string]interface{}, min_vms int, max_vms int) ([]models.ElasticityPolicy, error) {
	elasticityPolicies, err := ExpandElasticityPolicies(config["elasticity_policies"].([]interface{}), config["parameter"].(string))
	if err != nil {
		ep := make([]models.ElasticityPolicy, 0, len(config))
		return ep, err
	}

	return elasticityPolicies, nil
}

func updateScheduledPolicies(nodePoolDetail map[string]interface{}, min_vms int, max_vms int) ([]models.SchedulePolicy, error) {
	var scheduledPolicyList []models.SchedulePolicy

	switch nodePoolType := nodePoolDetail["node_pool_type"].(string); nodePoolType {
	case "Static":
		scheduledPolicyList = []models.SchedulePolicy{}
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

func UpdateScheduledDict(config map[string]interface{}, min_vms int, max_vms int) ([]models.SchedulePolicy, error) {
	scheduledDict := []models.SchedulePolicy{}
	for _, worker := range config["worker"].([]interface{}) {
		worker := worker.(map[string]interface{})
		scheduledWorker, err := UpdateScheduledWorker(worker, min_vms, max_vms)
		if err != nil {
			return []models.SchedulePolicy{}, err
		}
		return scheduledWorker, nil
	}
	return scheduledDict, nil
}

func UpdateScheduledWorker(config map[string]interface{}, min_vms int, max_vms int) ([]models.SchedulePolicy, error) {
	scheduledPolicies, err := ExpandScheduledPolicies(config["scheduled_policies"].([]interface{}), min_vms, max_vms)
	if err != nil {
		return []models.SchedulePolicy{}, err
	}

	return scheduledPolicies, nil
}

// readNodePoolsIntoState refreshes the node_pools attribute from the API so that
// changes made outside terraform are visible as drift. Without this the attribute
// only ever holds the values the user configured, and drift can never be detected.
func readNodePoolsIntoState(d *schema.ResourceData, apiClient *client.Client, clusterID string, location string) error {
	projectID := d.Get("project_id").(int)

	nodePoolResp, err := apiClient.GetKubernetesNodePools(clusterID, projectID, location)
	if err != nil {
		return fmt.Errorf("error reading node pools for cluster %s: %s", clusterID, err.Error())
	}
	apiPools, ok := nodePoolResp["data"].([]interface{})
	if !ok {
		// Nothing usable came back. Leave the existing value alone rather than
		// blanking it out on an unexpected response shape.
		log.Printf("[WARN] KUBERNETES READ | node pool response had no 'data' list, leaving node_pools untouched")
		return nil
	}

	existing, _ := d.Get("node_pools").([]interface{})

	// A cluster always keeps at least one node pool, and this endpoint answers with
	// an empty list for a cluster it cannot find. Treat empty as a failed lookup
	// rather than as "every pool was removed", otherwise a transient blank response
	// would produce a plan that proposes recreating every pool.
	if len(apiPools) == 0 && len(existing) > 0 {
		log.Printf("[WARN] KUBERNETES READ | node pool list came back empty for cluster %s, leaving node_pools untouched", clusterID)
		return nil
	}

	// slug_name and sku_id are not part of the node pool response, they have to be
	// recovered from the worker plans list using the plan name as the key.
	slugBySku, idBySku := buildWorkerPlanLookup(apiClient, projectID, location)

	return d.Set("node_pools", flattenNodePools(existing, apiPools, slugBySku, idBySku))
}

// buildWorkerPlanLookup maps sku_name to its slug and sku id. A failure here is not
// fatal: the caller falls back to whatever is already in state for those fields.
func buildWorkerPlanLookup(apiClient *client.Client, projectID int, location string) (map[string]string, map[string]string) {
	slugBySku := make(map[string]string)
	idBySku := make(map[string]string)

	workerPlans, err := apiClient.GetKubernetesWorkerPlans(projectID, location)
	if err != nil {
		log.Printf("[WARN] KUBERNETES READ | could not fetch worker plans for slug/sku lookup: %s", err.Error())
		return slugBySku, idBySku
	}
	plans, ok := workerPlans["data"].([]interface{})
	if !ok {
		return slugBySku, idBySku
	}
	for _, plan := range plans {
		planData, ok := plan.(map[string]interface{})
		if !ok {
			continue
		}
		specs, ok := planData["specs"].(map[string]interface{})
		if !ok {
			continue
		}
		skuName, _ := specs["sku_name"].(string)
		if skuName == "" {
			continue
		}
		if slug, ok := planData["plan"].(string); ok {
			slugBySku[skuName] = slug
		}
		if id, ok := specs["id"].(string); ok {
			idBySku[skuName] = id
		}
	}
	return slugBySku, idBySku
}

// flattenNodePools converts the node pool API response into the shape node_pools
// expects. The ordering already present in state is preserved: node_pools is a
// positional list and the API does not return pools in configuration order, so
// writing them back in API order would show a phantom diff on every plan.
func flattenNodePools(existing []interface{}, apiPools []interface{}, slugBySku map[string]string, idBySku map[string]string) []interface{} {
	byName := make(map[string]map[string]interface{})
	for _, p := range apiPools {
		pool, ok := p.(map[string]interface{})
		if !ok {
			continue
		}
		if name, ok := pool["service_name"].(string); ok && name != "" {
			byName[name] = pool
		}
	}

	flattened := make([]interface{}, 0, len(byName))
	seen := make(map[string]bool)

	// Keep the pools already tracked in state in their existing positions.
	for _, e := range existing {
		prev, ok := e.(map[string]interface{})
		if !ok {
			continue
		}
		name, _ := prev["name"].(string)
		pool, found := byName[name]
		if !found {
			// Deleted outside terraform. Dropping it here is what surfaces the drift.
			log.Printf("[INFO] KUBERNETES READ | node pool %q is in state but not on the cluster", name)
			continue
		}
		seen[name] = true
		flattened = append(flattened, flattenNodePool(pool, prev, slugBySku, idBySku))
	}

	// Anything on the cluster that state does not know about goes on the end, in a
	// stable order so repeated reads do not reshuffle the list.
	var unseen []string
	for name := range byName {
		if !seen[name] {
			unseen = append(unseen, name)
		}
	}
	sort.Strings(unseen)
	for _, name := range unseen {
		log.Printf("[INFO] KUBERNETES READ | node pool %q exists on the cluster but not in state", name)
		flattened = append(flattened, flattenNodePool(byName[name], nil, slugBySku, idBySku))
	}

	return flattened
}

// flattenNodePool maps a single node pool from the API onto the schema fields.
// prev is the matching entry already in state, or nil for a newly discovered pool.
// Fields the API does not report are carried over from prev rather than zeroed.
func flattenNodePool(pool map[string]interface{}, prev map[string]interface{}, slugBySku map[string]string, idBySku map[string]string) map[string]interface{} {
	nodePool := make(map[string]interface{})

	name, _ := pool["service_name"].(string)
	nodePool["name"] = name

	specsName, _ := pool["plan_name"].(string)
	if specsName == "" {
		specsName = carryString(prev, "specs_name")
	}
	nodePool["specs_name"] = specsName
	nodePool["slug_name"] = lookupOrCarry(slugBySku, specsName, prev, "slug_name")
	nodePool["sku_id"] = lookupOrCarry(idBySku, specsName, prev, "sku_id")

	if serviceID, ok := pool["service_id"].(float64); ok {
		nodePool["service_id"] = fmt.Sprintf("%.0f", serviceID)
	} else {
		nodePool["service_id"] = carryString(prev, "service_id")
	}

	autoScale, reported := pool["auto_scale_enabled"].(bool)
	if !reported {
		// Not in the response. Trust what state already recorded rather than
		// silently reclassifying an Autoscale pool as Static.
		autoScale = carryString(prev, "node_pool_type") == "Autoscale"
	}
	if autoScale {
		nodePool["node_pool_type"] = "Autoscale"
	} else {
		nodePool["node_pool_type"] = "Static"
	}

	cardinality, minVms, maxVms := workerRoleCounts(pool)
	if cardinality == 0 {
		if count, ok := pool["no_of_nodes"].(float64); ok {
			cardinality = int(count)
		}
	}
	// A zero means the response did not carry the value, not that the pool really has
	// no machines or no bounds: a pool cannot go below two. Static pools in particular
	// always report min_vms and max_vms as null, so writing the zero straight through
	// would fight any value the configuration had set.
	if cardinality == 0 {
		cardinality = carryInt(prev, "cardinality")
	}
	if minVms == 0 {
		minVms = carryInt(prev, "min_vms")
	}
	if maxVms == 0 {
		maxVms = carryInt(prev, "max_vms")
	}
	nodePool["cardinality"] = cardinality
	nodePool["min_vms"] = minVms
	nodePool["max_vms"] = maxVms

	// worker_node only describes Static pools. For Autoscale the live count lives in
	// cardinality, so keep whatever the configuration had to avoid a phantom diff. A
	// Static pool with no reported count keeps its previous value for the same reason.
	if autoScale || cardinality == 0 {
		nodePool["worker_node"] = carryInt(prev, "worker_node")
	} else {
		nodePool["worker_node"] = cardinality
	}

	// Not reported by the node pool API, so preserve what is already known.
	nodePool["node_pool_size"] = carryInt(prev, "node_pool_size")
	nodePool["policy_type"] = carryString(prev, "policy_type")
	nodePool["custom_param_name"] = carryString(prev, "custom_param_name")
	nodePool["custom_param_value"] = carryString(prev, "custom_param_value")
	nodePool["elasticity_dict"] = carryList(prev, "elasticity_dict")
	nodePool["scheduled_dict"] = carryList(prev, "scheduled_dict")

	return nodePool
}

// workerRoleCounts pulls the live cardinality and autoscale bounds off the worker role.
func workerRoleCounts(pool map[string]interface{}) (int, int, int) {
	roles, ok := pool["roles"].([]interface{})
	if !ok {
		return 0, 0, 0
	}
	for _, r := range roles {
		role, ok := r.(map[string]interface{})
		if !ok {
			continue
		}
		if roleName, _ := role["role_name"].(string); roleName != "worker" {
			continue
		}
		cardinality, minVms, maxVms := 0, 0, 0
		if v, ok := role["cardinality"].(float64); ok {
			cardinality = int(v)
		}
		if v, ok := role["min_vms"].(float64); ok {
			minVms = int(v)
		}
		if v, ok := role["max_vms"].(float64); ok {
			maxVms = int(v)
		}
		return cardinality, minVms, maxVms
	}
	return 0, 0, 0
}

func lookupOrCarry(lookup map[string]string, key string, prev map[string]interface{}, field string) string {
	if value, ok := lookup[key]; ok && value != "" {
		return value
	}
	return carryString(prev, field)
}

func carryString(prev map[string]interface{}, field string) string {
	if prev == nil {
		return ""
	}
	value, _ := prev[field].(string)
	return value
}

func carryInt(prev map[string]interface{}, field string) int {
	if prev == nil {
		return 0
	}
	value, _ := prev[field].(int)
	return value
}

func carryList(prev map[string]interface{}, field string) []interface{} {
	if prev == nil {
		return []interface{}{}
	}
	value, ok := prev[field].([]interface{})
	if !ok {
		return []interface{}{}
	}
	return value
}
