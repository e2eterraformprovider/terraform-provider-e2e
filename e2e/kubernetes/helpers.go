package kubernetes

import (
	// "container/list"
	"fmt"
	"log"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
	// "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ExpandNodePools(config []interface{}, apiClient *client.Client, project_id int, location string) ([]models.NodePool, error) {
	nodePools := make([]models.NodePool, 0, len(config))

	for _, np := range config {
		nodePoolDetail := np.(map[string]interface{})
		workerPlans, err := apiClient.GetKubernetesWorkerPlans(project_id, location) //Here we are are tryig to get all worker plans
		if err != nil {
			return nil, err
		}
		log.Printf("---------------IDHAR TOH PAHUCHA(3)-----------------")
		plans := workerPlans["data"].([]interface{})
		var matchingPlan map[string]interface{}
		log.Printf("---------------IDHAR TOH PAHUCHA(4)-----------------")
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
		// elasticityDict, err := ExpandElasticityDict(nodePoolDetail["elasticity_dict"].(map[string]interface{}))
		// if err != nil {
		// 	return nil, err
		// }

		if _, ok := nodePoolDetail["node_pool_type"]; !ok {
			return nil, fmt.Errorf("node_pool_type is required")
		}
		log.Printf("---------------IDHAR TOH PAHUCHA(5)-----------------")
		var policyType, customParamName, customParamValue string
		elasticityDict := models.ElasticityDict{}
		scheduledDict := models.ScheduledDict{}

		// If node_pool_type is Static, omit policyType, customParamName, and customParamValue
		if nodePoolDetail["node_pool_type"].(string) == "Static" {
			policyType = ""
			customParamName = ""
			customParamValue = ""
		} else {
			log.Printf("---------------IDHAR TOH PAHUCHA(5a)-----------------")
			policyType = getPolicyType(nodePoolDetail)
			log.Printf("---------------IDHAR TOH PAHUCHA(5b)-----------------")
			customParamName = getCustomParamName(nodePoolDetail)
			log.Printf("---------------IDHAR TOH PAHUCHA(5c)-----------------")
			customParamValue = getCustomParamValue(nodePoolDetail)
			log.Printf("HEMLOOOOOO--------%+v--------%+v-----------%+v-----------", policyType, customParamName, customParamValue)

			if _, ok := nodePoolDetail["min_vms"]; !ok {
				return nil, fmt.Errorf("in case of Autoscale node type, the 'min_vms' field is required")
			}
			if _, ok := nodePoolDetail["max_vms"]; !ok {
				return nil, fmt.Errorf("in case of Autoscale node type, the 'max_vms' field is required")
			}

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
	return nodePools, nil
}

// ExpandElasticityDict is a helper function to process the elasticity_dict attribute.
func ExpandElasticityDict(config map[string]interface{}, min_vms int, max_vms int) (models.ElasticityDict, error) {
	log.Printf("---------------IDHAR TOH PAHUCHA(7)-----------------")
	elasticityDict := models.ElasticityDict{}
	for _, worker := range config["worker"].([]interface{}) {
		worker := worker.(map[string]interface{})
		elasticityWorker, err := ExpandElasticityWorker(worker, min_vms, max_vms)
		// elasticityWorker, err := ExpandElasticityWorker(config["worker"].(map[string]interface{}), parameter)
		if err != nil {
			return models.ElasticityDict{}, err
		}

		elasticityDict = models.ElasticityDict{
			Worker: elasticityWorker,
		}
		return elasticityDict, nil
		// elasticityDict, _ := ExpandElasticityDict(ed, nodePoolDetail["parameter"].(string))
		// return elasticityDict, fmt.Errorf("elastic Dictionary could not be expanded successfully")
	}
	return elasticityDict, nil
}

func ExpandScheduledDict(config map[string]interface{}, min_vms int, max_vms int) (models.ScheduledDict, error) {
	log.Printf("---------------IDHAR TOH PAHUCHA(11)-----------------")
	scheduledDict := models.ScheduledDict{}
	for _, worker := range config["worker"].([]interface{}) {
		worker := worker.(map[string]interface{})
		scheduledWorker, err := ExpandScheduledWorker(worker, min_vms, max_vms)
		// elasticityWorker, err := ExpandElasticityWorker(config["worker"].(map[string]interface{}), parameter)
		if err != nil {
			return models.ScheduledDict{}, err
		}

		scheduledDict = models.ScheduledDict{
			Worker: scheduledWorker,
		}
		return scheduledDict, nil
		// elasticityDict, _ := ExpandElasticityDict(ed, nodePoolDetail["parameter"].(string))
		// return elasticityDict, fmt.Errorf("elastic Dictionary could not be expanded successfully")
	}
	return scheduledDict, nil
}

// ExpandElasticityWorker is a helper function to process the worker attribute in elasticity_dict.
func ExpandElasticityWorker(config map[string]interface{}, min_vms int, max_vms int) (models.ElasticityWorker, error) {
	log.Printf("---------------IDHAR TOH PAHUCHA(8)-----------------")
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
	log.Printf("---------------IDHAR TOH PAHUCHA(12)-----------------")
	scheduledPolicies, err := ExpandScheduledPolicies(config["scheduled_policies"].([]interface{}))
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
	log.Printf("---------------IDHAR TOH PAHUCHA(9)-----------------")
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

func ExpandScheduledPolicies(config []interface{}) ([]models.SchedulePolicy, error) {
	scheduledPolicies := make([]models.SchedulePolicy, 0, len(config))
	for _, sp := range config {
		elasticityPolicyDetail := sp.(map[string]interface{})
		// Adjust should be converted to an integer
		upscaleCardinality := elasticityPolicyDetail["upscale_cardinality"].(int)
		downscaleCardinality := elasticityPolicyDetail["downscale_cardinality"].(int)
		upscaleRecurrence := elasticityPolicyDetail["upscale_recurrence"].(string)
		downscaleRecurrence := elasticityPolicyDetail["downscale_recurrence"].(string)

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
	log.Printf("---------------IDHAR TOH PAHUCHA(6)-----------------")
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
		// elasticityPolicies, err := ExpandElasticityPolicies(nodePoolDetail, nodePoolDetail["parameter"].(string))
		// if err != nil {
		// 	return elasticityDict, err
		// }
		// elasticityDict = models.ElasticityDict{Worker: models.ElasticityWorker{ElasticityPolicies: elasticityPolicies}}
	default:
		return elasticityDict, fmt.Errorf("invalid node_pool_type: %s", nodePoolType)
	}
	// if nodePoolType, ok := nodePoolDetail["node_pool_type"].(string); ok && nodePoolType == "Static" {
	// 	return models.ElasticityDict{} // Return empty ElasticityDict for "Static"
	// }
	return elasticityDict, nil
}

func getScheduledDict(nodePoolDetail map[string]interface{}, min_vms int, max_vms int) (models.ScheduledDict, error) {
	log.Printf("---------------IDHAR TOH PAHUCHA(10)-----------------")
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
	//Can emit this for efficient code
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

// func getPolicyType(nodePoolDetail map[string]interface{}) string {
// 	elasticityDict, ok := nodePoolDetail["elasticity_dict"].([]interface{})
// 	if !ok {
// 		log.Printf("Elasticity dictionary not found or not in the expected format")
// 		return ""
// 	}
// 	sdp, scheduledDictPresent := nodePoolDetail["scheduled_dict"].([]interface{})
// 	isSDPresent := true
// 	if scheduledDictPresent {
// 		for _, sd := range sdp {
// 			sdMap, ok := sd.(map[string]interface{})
// 			if !ok {
// 				log.Printf("Scheduled dictionary entry is not in the expected format")
// 				continue
// 			}
// 			workerList, ok := sdMap["worker"].([]interface{})
// 			if !ok || len(workerList) == 0 {
// 				log.Printf("Worker list in scheduled dictionary not found or empty")
// 				isSDPresent = false
// 				continue
// 			}
// 		}
// 	}
// 	for _, ed := range elasticityDict {
// 		edMap, ok := ed.(map[string]interface{})
// 		if !ok {
// 			log.Printf("Elasticity dictionary entry is not in the expected format")
// 			continue
// 		}
// 		workerList, ok := edMap["worker"].([]interface{})
// 		if !ok || len(workerList) == 0 {
// 			log.Printf("Worker list not found or empty")
// 			continue
// 		}
// 		// Assuming there is only one worker map in the list
// 		worker, ok := workerList[0].(map[string]interface{})
// 		if !ok {
// 			log.Printf("Worker map is not in the expected format")
// 			continue
// 		}
// 		policyParamType, ok := worker["policy_paramter_type"].(string)
// 		if !ok {
// 			log.Printf("Policy parameter type not found or not a string")
// 			continue
// 		}
// 		// Check if scheduled_dict is present
// 		if scheduledDictPresent && isSDPresent {
// 			return policyParamType + "-Scheduled"
// 		} else {
// 			return policyParamType
// 		}
// 	}
// 	if scheduledDictPresent {
// 		return "Scheduled"
// 	}
// 	return ""
// }

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
