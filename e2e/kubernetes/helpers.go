package kubernetes

import (
	// "container/list"
	"fmt"
	"log"

	// "strconv"
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

		// If node_pool_type is Static, omit policyType, customParamName, and customParamValue
		if nodePoolDetail["node_pool_type"].(string) == "Static" {
			policyType = ""
			customParamName = ""
			customParamValue = ""
		} else {
			policyType = getPolicyType(nodePoolDetail)
			customParamName = getCustomParamName(nodePoolDetail)
			customParamValue = getCustomParamValue(nodePoolDetail)
			log.Printf("HEMLOOOOOO--------%+v--------%+v-----------%+v-----------", policyType, customParamName, customParamValue)
		}
		elasticityDict, err := getElasticityDict(nodePoolDetail)
		if err != nil {
			log.Printf("Invalid format for Elast")
		}

		nodePool := models.NodePool{
			Name:             nodePoolDetail["name"].(string),
			SlugName:         matchingPlan["plan"].(string),
			SKUID:            matchingPlan["specs"].(map[string]interface{})["id"].(string),
			SpecsName:        nodePoolDetail["specs_name"].(string),
			WorkerNode:       nodePoolDetail["worker_node"].(int),
			ElasticityDict:   elasticityDict,
			PolicyType:       policyType,
			CustomParamName:  customParamName,
			CustomParamValue: customParamValue,
		}

		nodePools = append(nodePools, nodePool)
	}
	return nodePools, nil
}

// ExpandElasticityDict is a helper function to process the elasticity_dict attribute.
func ExpandElasticityDict(config map[string]interface{}, parameter string) (models.ElasticityDict, error) {
	log.Printf("---------------IDHAR TOH PAHUCHA(7)-----------------")
	elasticityDict := models.ElasticityDict{}
	for _, worker := range config["worker"].([]interface{}) {
		worker := worker.(map[string]interface{})
		elasticityWorker, err := ExpandElasticityWorker(worker, parameter)
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

// ExpandElasticityWorker is a helper function to process the worker attribute in elasticity_dict.
func ExpandElasticityWorker(config map[string]interface{}, parameter string) (models.ElasticityWorker, error) {
	log.Printf("---------------IDHAR TOH PAHUCHA(8)-----------------")
	elasticityPolicies, err := ExpandElasticityPolicies(config["elasticity_policies"].([]interface{}), parameter)
	if err != nil {
		return models.ElasticityWorker{}, err
	}

	return models.ElasticityWorker{
		MinVms:             config["min_vms"].(int),
		Cardinality:        config["min_vms"].(int),
		MaxVms:             config["max_vms"].(int),
		ElasticityPolicies: elasticityPolicies,
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

func getElasticityDict(nodePoolDetail map[string]interface{}) (models.ElasticityDict, error) {
	log.Printf("---------------IDHAR TOH PAHUCHA(6)-----------------")
	var elasticityDict models.ElasticityDict

	// Handle ElasticityDict based on node_pool_type
	switch nodePoolType := nodePoolDetail["node_pool_type"].(string); nodePoolType {
	case "Static":
		elasticityDict = models.ElasticityDict{}
	case "Autoscale":
		for _, ed := range nodePoolDetail["elasticity_dict"].([]interface{}) {
			ed := ed.(map[string]interface{})
			elasticityDict, _ := ExpandElasticityDict(ed, nodePoolDetail["parameter"].(string))
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

func getCustomParamName(nodePoolDetail map[string]interface{}) string {
	//Can emit this for efficient code
	if nodePoolType, ok := nodePoolDetail["node_pool_type"].(string); ok && nodePoolType == "Static" {
		return "" // Return empty string for "Static"
	}
	policyParameterType := getPolicyType(nodePoolDetail)
	if policyParameterType == "" || policyParameterType == "Default" {
		return "" // Return empty string when policy_parameter_type is not provided or is "Default"
	}
	if nodePoolDetail["parameter"] == "CPU" || nodePoolDetail["parameter"] == "Memory" {
		log.Printf("Cannot use Default parameters in case of Custom")
		return ""
	}
	return nodePoolDetail["parameter"].(string)
}

func getCustomParamValue(nodePoolDetail map[string]interface{}) string {
	if nodePoolType, ok := nodePoolDetail["node_pool_type"].(string); ok && nodePoolType == "Static" {
		return "" // Return empty string for "Static"
	}
	policyParameterType := getPolicyType(nodePoolDetail)
	if policyParameterType == "" || policyParameterType == "Default" {
		return "" // Return empty string when policy_parameter_type is not provided or is "Default"
	}
	return "0"
}

func getPolicyType(nodePoolDetail map[string]interface{}) string {
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
		policyParamType, ok := worker["policy_paramter_type"].(string)
		if !ok {
			log.Printf("Policy parameter type not found or not a string")
			continue
		}
		return policyParamType
	}
	return ""
}
