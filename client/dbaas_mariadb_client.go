package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"net/http"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func (c *Client) CreateMariaDB(req *models.MariaDBCreateRequest, projectID, location string) (*models.DB, error) {
	url := c.Api_endpoint + "/rds/cluster/"

	payloadBuf := new(bytes.Buffer)
	if err := json.NewEncoder(payloadBuf).Encode(req); err != nil {
		return nil, fmt.Errorf("failed to encode create payload: %v", err)
	}
	
	httpReq, err := http.NewRequest("POST", url, payloadBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	log.Printf("[DEBUG] MariaDB create request payload: %+v", req)

	httpReq = addParamsAndHeaders(httpReq, c.Api_key, c.Auth_token, projectID, location)

	resp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var response models.DBResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &response.Data, nil
}

func (c *Client) ReadMariaDB(id string, projectID string, location string) (*models.DB, error) {
	url := c.Api_endpoint + "/rds/cluster/" + id + "/"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req = addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectID, location)

	log.Printf("[DEBUG] ReadMariaDB Request URL: %s", req.URL.String())
	log.Printf("[DEBUG] ReadMariaDB Headers: %+v", req.Header)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("read mariadb failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var response models.DBResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &response.Data, nil
}

func (c *Client) MariaDBExists(id string, projectID string, location string) (bool, error) {
	url := c.Api_endpoint + "/rds/cluster/" + id + "/"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %v", err)
	}

	req = addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectID, location)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return true, nil
}

func (c *Client) DeleteMariaDB(id string, projectID string, location string) error {
	url := c.Api_endpoint + "/rds/cluster/" + id + "/"

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create DELETE request: %v", err)
	}

	req = addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectID, location)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("DELETE request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("DELETE failed with status: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) ShutdownMariaDB(id string, projectID string, location string) error {
	url := c.Api_endpoint + "/rds/cluster/" + id + "/shutdown"

	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create shutdown request: %v", err)
	}

	req = addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectID, location)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("shutdown request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("shutdown failed: %d - %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *Client) ResumeMariaDB(id string, projectID string, location string) error {
	url := c.Api_endpoint + "/rds/cluster/" + id + "/resume"

	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create resume request: %v", err)
	}

	req = addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectID, location)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("resume request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("resume failed: %d - %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *Client) RestartMariaDB(id string, projectID string, location string) error {
	url := c.Api_endpoint + "/rds/cluster/" + id + "/restart"

	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create restart request: %v", err)
	}

	req = addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectID, location)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("restart request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("restart failed: %d - %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *Client) AttachVPCToMariaDB(id string, projectID string, location string, vpcIDs []string) error {
	// Expand metadata for all VPC IDs
	vpcMetaList, err := c.ExpandMariaDBVpcList(vpcIDs, projectID, location)
	if err != nil {
		return fmt.Errorf("failed to expand VPC metadata: %v", err)
	}

	// Build request payload
	payload := models.AttachDetachVPCRequest{
		Action: "attach",
		VPCs:   vpcMetaList,
	}

	url := c.Api_endpoint + "/rds/cluster/" + id + "/vpc-attach/"

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal VPC attach payload: %v", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req = addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectID, location)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("VPC attach request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("VPC attach failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	return nil
}

func (c *Client) DetachVPCFromMariaDB(id string, projectID string, location string, vpcIDs []string) error {
	// Expand metadata for all VPC IDs
	vpcMetaList, err := c.ExpandMariaDBVpcList(vpcIDs, projectID, location)
	if err != nil {
		return fmt.Errorf("failed to expand VPC metadata: %v", err)
	}

	// Build request payload
	payload := models.AttachDetachVPCRequest{
		Action: "detach",
		VPCs:   vpcMetaList,
	}

	url := c.Api_endpoint + "/rds/cluster/" + id + "/vpc-detach/"

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal VPC detach payload: %v", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req = addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectID, location)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("VPC detach request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("VPC detach failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func (c *Client) AttachPublicIPToMariaDB(id string, projectID string, location string) error {
	url := c.Api_endpoint + "/rds/cluster/" + id + "/public-ip-attach/"

	// Prepare JSON payload for the attach action
	payload := map[string]string{"action": "attach"}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal public IP attach payload: %v", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create attach public IP request: %v", err)
	}

	req = addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectID, location)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("public IP attach request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("public IP attach failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func (c *Client) DetachPublicIPFromMariaDB(id string, projectID string, location string) error {
	url := c.Api_endpoint + "/rds/cluster/" + id + "/public-ip-detach/"

	// Prepare JSON payload for the detach action
	payload := map[string]string{"action": "detach"}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal public IP detach payload: %v", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create detach public IP request: %v", err)
	}

	req = addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectID, location)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("public IP detach request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("public IP detach failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func (c *Client) AttachParameterGroupToMariaDB(clusterID string, parameterGroupID int, projectID, location string) error {
	url := c.Api_endpoint + "/rds/cluster/" + clusterID + "/parameter-group/" + strconv.Itoa(parameterGroupID) + "/add"

	payload := models.ParameterGroupRequest{Action: "add"}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal attach payload: %v", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create attach request: %v", err)
	}

	req = addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectID, location)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("attach request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("attach failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func (c *Client) DetachParameterGroupFromMariaDB(clusterID string, parameterGroupID int, projectID, location string) error {
	url := c.Api_endpoint + "/rds/cluster/" + clusterID + "/parameter-group/" + strconv.Itoa(parameterGroupID) + "/detach"

	body := bytes.NewBuffer([]byte("{}"))

	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return fmt.Errorf("failed to create detach request: %v", err)
	}

	req = addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectID, location)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("detach request failed: %v", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("detach failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func (c *Client) UpgradeMariaDBPlan(clusterID, projectID, location string, templateID int) error {
	url := c.Api_endpoint + "/rds/cluster/" + clusterID + "/rds-upgrade/"

	payload := map[string]interface{}{
		"template_id": templateID,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal upgrade payload: %v", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create upgrade request: %v", err)
	}

	req = addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectID, location)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("upgrade request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upgrade failed: %d - %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func (c *Client) ExpandMariaDBDisk(clusterID, projectID, location string, additionalSize int) error {
	if additionalSize <= 0 {
		log.Printf("[INFO] Skipping disk expansion: additional size is %d GB (must be > 0)", additionalSize)
		return nil
	}

	log.Printf("[INFO] Initiating disk expansion: cluster=%s, additional_size=%d GB", clusterID, additionalSize)

	url := c.Api_endpoint + "/rds/cluster/" + clusterID + "/disk-upgrade/"

	payload := models.DiskUpgradeRequest{Size: additionalSize}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal disk upgrade payload: %v", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create disk upgrade request: %v", err)
	}

	req = addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectID, location)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("disk upgrade request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("disk upgrade failed with status %d: %s", resp.StatusCode, string(body))
	}

	log.Printf("[INFO] Disk expansion completed: +%d GB added to MariaDB cluster %s", additionalSize, clusterID)
	return nil
}











