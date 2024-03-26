package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func (c *Client) GetKubernetesMasterPlans(project_id int, location string) (map[string]interface{}, error) {
	url := c.Api_endpoint + "kubernetes/plans"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("project_id", strconv.Itoa(project_id))
	params.Add("location", location)
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "terraform-e2e")
	// log.Printf("----------------REQUEST FOR MASTER PLANS-------------: %+v", req)
	response, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	err = CheckResponseStatus(response)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	resBody, _ := ioutil.ReadAll(response.Body)
	stringresponse := string(resBody)
	// log.Printf("%s", stringresponse)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)

	if err != nil {
		log.Printf("[ERROR] CLIENT GET KUBERNETES MASTER PLANS | error when unmarshalling")
		return nil, err
	}

	return jsonRes, nil
}

func (c *Client) GetKubernetesWorkerPlans(project_id int, location string) (map[string]interface{}, error) {
	url := c.Api_endpoint + "kubernetes/worker-plans/"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("project_id", strconv.Itoa(project_id))
	params.Add("location", location)
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "terraform-e2e")

	response, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	err = CheckResponseStatus(response)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	resBody, _ := ioutil.ReadAll(response.Body)
	stringresponse := string(resBody)
	// log.Printf("%s", stringresponse)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)

	if err != nil {
		log.Printf("[ERROR] CLIENT GET KUBERNETES WORKER PLANS | error when unmarshalling")
		return nil, err
	}

	return jsonRes, nil
}

func (c *Client) NewKubernetesService(item *models.KubernetesCreate, project_id int, location string) (map[string]interface{}, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(item)
	if err != nil {
		return nil, err
	}
	UrlEndPoint := c.Api_endpoint + "kubernetes/"
	log.Printf("[INFO] CLIENT KUBERNETES| BEFORE REQUEST")
	if err != nil {
		return nil, err
	}

	buf, err = RemoveExtraFieldsFromKubernetes(&buf)
	if err != nil {
		return nil, err
	}
	log.Printf("-----------AFTER REMOVING FIELDS FROM PAYLOAD------------: %+v", &buf)
	req, err := http.NewRequest("POST", UrlEndPoint, &buf)
	if err != nil {
		return nil, err
	}
	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)
	response, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	err = CheckResponseCreatedStatus(response)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	resBody, _ := ioutil.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, err
	}
	return jsonRes, nil
}

func (c *Client) GetKubernetesServiceInfo(kubernetesID string, location string, project_id int) (map[string]interface{}, error) {
	urlKubernetes := c.Api_endpoint + "kubernetes/" + kubernetesID
	req, err := http.NewRequest("GET", urlKubernetes, nil)
	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] CLIENT | KUBERNETES READ")
	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)
	response, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] CLIENT KUBERNETES READ | response code %d", response.StatusCode)
	err = CheckResponseStatus(response)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	resBody, _ := ioutil.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		log.Printf("[ERROR] CLIENT GET LOAD BALANCER | error when unmarshalling | %s", err)
		return nil, err
	}
	return jsonRes, nil
}

func (c *Client) DeleteKubernetesService(kubernetesID string, location string, project_id int) error {
	deleteURL := c.Api_endpoint + "kubernetes/" + kubernetesID
	req, err := http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		return err
	}

	log.Printf("[INFO] CLIENT | KUBERNETES DELETE")
	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)
	if err != nil {
		return err
	}

	response, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	err = CheckResponseStatus(response)
	if err != nil {
		return err
	}

	return nil
}

func RemoveExtraFieldsFromKubernetes(buf *bytes.Buffer) (bytes.Buffer, error) {

	jsonData := buf.Bytes()

	// jsonData := buf.Bytes()
	var data map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		return *buf, err
	}

	nodePools, ok := data["node_pools"].([]interface{})
	if !ok {
		// If "node_pools" is not present or not a slice, return an error
		return *buf, errors.New("node_pools field is missing or invalid")
	}

	for _, nodePool := range nodePools {
		if nodePoolMap, ok := nodePool.(map[string]interface{}); ok {
			log.Printf("-------------------WORKER_NODE TYPE-----------------: %T", nodePoolMap["worker_node"])
			// Type assert to float64
			workerNode, workerNodePresent := nodePoolMap["worker_node"].(float64)
			if workerNodePresent {
				if workerNode == 0 {
					log.Printf("-------------------WORKER_NODE ARE 0-----------------:")
					// If worker_node is present and its value is 0, delete the "worker_node" field
					delete(nodePoolMap, "worker_node")
				} else if workerNode >= 2 {
					// If worker_node is greater than or equal to 2, check if "elasticity_dict" is present
					if _, elasticityDictPresent := nodePoolMap["elasticity_dict"].([]interface{}); elasticityDictPresent {
						// If present, delete the "elasticity_dict" field
						log.Printf("Deleted elasticity_dict since worker_node field is >= 2")
						delete(nodePoolMap, "elasticity_dict")
					}
					if _, scheduledDictPresent := nodePoolMap["elasticity_dict"]; scheduledDictPresent {
						// If present, delete the "elasticity_dict" field
						log.Printf("Deleted elasticity_dict since worker_node field is >= 2")
						delete(nodePoolMap, "scheduled_dict")
					}
				}
			}
			policyType, policyTypePresent := nodePoolMap["policy_type"].(string)
			if !policyTypePresent || (policyType == "Scheduled") {
				// If policy_type does not contain the keyword "Scheduled", remove scheduled_dict
				delete(nodePoolMap, "elasticity_dict")
			}
			if !policyTypePresent || !strings.Contains(policyType, "Scheduled") {
				// If policy_type does not contain the keyword "Scheduled", remove scheduled_dict
				delete(nodePoolMap, "scheduled_dict")
			}
		}
	}

	NewjsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error encoding data:", err)
		return *buf, nil
	}

	newBuffer := bytes.NewBuffer(NewjsonData)
	return *newBuffer, nil
}

func CheckResponseCreatedStatus(response *http.Response) error {
	if response.StatusCode != http.StatusCreated {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(response.Body)
		if err != nil {
			return fmt.Errorf("got a non 201 status code: %v", response.StatusCode)
		}
		return fmt.Errorf("got a non 201 status code: %v - %s", response.StatusCode, respBody.String())
	}
	return nil
}

func (c *Client) GetKubernetesNodePools(clusterID string, project_id int, location string) (map[string]interface{}, error) {
	url := c.Api_endpoint + "kubernetes/node-pool-services/" + clusterID
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)
	// log.Printf("----------------REQUEST FOR MASTER PLANS-------------: %+v", req)
	response, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	err = CheckResponseStatus(response)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	resBody, _ := ioutil.ReadAll(response.Body)
	stringresponse := string(resBody)
	// log.Printf("%s", stringresponse)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)

	if err != nil {
		log.Printf("[ERROR] CLIENT GET LIST OF KUBERNETES CLUSTER NODE POOLS| error when unmarshalling")
		return nil, err
	}

	return jsonRes, nil
}

func (c *Client) UpdateNodePoolCardinality(nodePoolServiceID string, project_id int, location string) (map[string]interface{}, error) {
	urlNode := c.Api_endpoint + "kubernetes/cluster-update/" + nodePoolServiceID
	req, err := http.NewRequest("PUT", urlNode, nil)
	if err != nil {
		return nil, err
	}
	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)
	response, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	resBody, _ := ioutil.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)

	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, err
	}
	return jsonRes, nil
}

func (c *Client) DeleteNodePool(nodePoolServiceID string, project_id int, location string) (map[string]interface{}, error) {
	urlNode := c.Api_endpoint + "kubernetes/delete-node-pool-service/" + nodePoolServiceID
	req, err := http.NewRequest("DELETE", urlNode, nil)
	if err != nil {
		return nil, err
	}
	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)
	response, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	log.Printf("----------------RESPONSE FOR DELETE 204 NO CONTENT(Resource.go)----------------: %+v", response)
	if response.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	resBody, _ := ioutil.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)

	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, err
	}
	return jsonRes, nil
}

func (c *Client) AddNodePool(item *models.NodePoolAdd, nodePoolServiceID string, project_id int, location string) (map[string]interface{}, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(item)
	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] CLIENT KUBERNETES ADD NODE POOL| BEFORE REQUEST")
	if err != nil {
		return nil, err
	}

	buf, err = RemoveExtraFieldsFromKubernetes(&buf)
	if err != nil {
		return nil, err
	}
	urlNode := c.Api_endpoint + "kubernetes/add-node-pools/" + nodePoolServiceID
	req, err := http.NewRequest("POST", urlNode, nil)
	if err != nil {
		return nil, err
	}
	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)
	response, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	err = CheckResponseCreatedStatus(response)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	resBody, _ := ioutil.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, err
	}
	return jsonRes, nil
}

func (c *Client) UpdateNodePoolDetails(item *models.NodePoolUpdate, nodePoolServiceID string, project_id int, location string) (map[string]interface{}, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(item)
	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] CLIENT KUBERNETES ADD NODE POOL| BEFORE REQUEST")
	if err != nil {
		return nil, err
	}
	urlNode := c.Api_endpoint + "kubernetes/update-node-pool/" + nodePoolServiceID
	req, err := http.NewRequest("PUT", urlNode, &buf)
	if err != nil {
		return nil, err
	}
	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)
	response, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	err = CheckResponseStatus(response)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	resBody, _ := ioutil.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)

	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, err
	}
	return jsonRes, nil
}
