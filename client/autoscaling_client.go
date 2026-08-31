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

func (c *Client) CreateScalerGroup(req *models.CreateScalerGroupRequest, projectID, location string) (*models.ScalerGroupCreateDetails, error) {
	url := c.Api_endpoint + "/scaler/scalegroups"
	log.Printf("[INFO] Sending request to create Scaler Group at: %s", url)

	payloadBuf := new(bytes.Buffer)
	if err := json.NewEncoder(payloadBuf).Encode(req); err != nil {
		return nil, fmt.Errorf("failed to encode create payload: %v", err)
	}

	httpReq, err := http.NewRequest("POST", url, payloadBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %v", err)
	}
	httpReq = addParamsAndHeaders(httpReq, c.Api_key, c.Auth_token, projectID, location)

	log.Printf("[DEBUG] CreateScalerGroup headers: %v", httpReq.Header)

	resp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("create scaler group failed: status %d\nresponse: %s", resp.StatusCode, string(bodyBytes))
	}

	var response models.CreateScalerGroupResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, fmt.Errorf("failed to decode create response: %v\nresponse body: %s", err, string(bodyBytes))
	}

	log.Printf("[INFO] Scaler Group created successfully: ID=%s Name=%s", response.Data.ID, response.Data.Name)

	return &response.Data, nil
}

func (c *Client) GetScalerGroup(scaleGroupID, projectID, location string) (*models.ScalerGroupGetDetail, error) {
	url := c.Api_endpoint + "/scaler/scalegroups/" + scaleGroupID + "/"
	log.Printf("[INFO] Fetching Scaler Group details for ID: %s", scaleGroupID)

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %v", err)
	}
	httpReq = addParamsAndHeaders(httpReq, c.Api_key, c.Auth_token, projectID, location)

	log.Printf("[DEBUG] GetScalerGroup request URL: %s", httpReq.URL.String())

	resp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get scaler group failed: status %d\nresponse: %s", resp.StatusCode, string(bodyBytes))
	}

	var response models.GetScalerGroupResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, fmt.Errorf("failed to decode get response: %v\nresponse body: %s", err, string(bodyBytes))
	}

	log.Printf("[INFO] Scaler Group fetched: Name=%s Desired=%d Running=%d", response.Data.Name, response.Data.Desired, response.Data.Running)

	return &response.Data, nil
}

func (c *Client) DeleteScalerGroup(scaleGroupID, projectID, location string) error {
	url := c.Api_endpoint + "/scaler/scalegroups/" + scaleGroupID + "/"
	log.Printf("[INFO] Sending delete request for Scaler Group ID: %s", scaleGroupID)

	httpReq, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create DELETE request: %v", err)
	}
	httpReq = addParamsAndHeaders(httpReq, c.Api_key, c.Auth_token, projectID, location)

	log.Printf("[DEBUG] DeleteScalerGroup request URL: %s", httpReq.URL.String())

	resp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read delete response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("delete scaler group failed: status %d\nresponse: %s", resp.StatusCode, string(bodyBytes))
	}

	var response models.DeleteScalerGroupResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return fmt.Errorf("failed to decode delete response: %v", err)
	}

	log.Printf("[INFO] Scaler Group deleted successfully. Message: %s", response.Message)
	return nil
}

func (c *Client) GetSavedImageByName(imageName, projectID, location string) (*models.SavedImage, error) {
	url := c.Api_endpoint + "/images/saved-images/"
	log.Printf("[INFO] Sending request to fetch saved image: %s", imageName)

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %v", err)
	}

	httpReq = addParamsAndHeaders(httpReq, c.Api_key, c.Auth_token, projectID, location)

	log.Printf("[DEBUG] GetSavedImageByName request URL: %s", httpReq.URL.String())

	resp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get saved image failed: status %d\nresponse: %s", resp.StatusCode, string(bodyBytes))
	}

	var result models.ListSavedImagesResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to decode saved-images response: %v\nresponse body: %s", err, string(bodyBytes))
	}

	for _, img := range result.Data {
		if img.Name == imageName {
			log.Printf("[INFO] Found saved image: ID=%s, TemplateID=%d, Distro=%s", img.ImageID, img.TemplateID, img.Distro)
			return &img, nil
		}
	}
	return nil, fmt.Errorf("no saved image found with name: %s", imageName)
}

func (c *Client) GetDefaultSecurityGroupID(projectID, location string) (int, error) {
	url := c.Api_endpoint + "security_group/"
	log.Printf("[INFO] Sending request to fetch default Security Group at: %s", url)

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create GET request: %v", err)
	}

	httpReq = addParamsAndHeaders(httpReq, c.Api_key, c.Auth_token, projectID, location)
	log.Printf("[DEBUG] GetDefaultSecurityGroup headers: %v", httpReq.Header)

	resp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		return 0, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("get default security group failed: status %d\nresponse: %s", resp.StatusCode, string(bodyBytes))
	}

	var response models.GetScalerSecurityGroupsResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return 0, fmt.Errorf("failed to decode security group response: %v\nresponse body: %s", err, string(bodyBytes))
	}

	for _, sg := range response.Data {
		if sg.IsDefault {
			log.Printf("[INFO] Found default Security Group: ID=%d", sg.ID)
			return sg.ID, nil
		}
	}
	return 0, fmt.Errorf("default security group not found")
}

func (c *Client) GetPlanDetailsFromPlanName(templateID int, planName, projectID, location string) (string, string, error) {
	url := c.Api_endpoint + fmt.Sprintf("/images/upgradeimage/%d/", templateID)
	log.Printf("[INFO] Sending request to fetch plan details for planName=%s, templateID=%d at: %s", planName, templateID, url)

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to create GET request: %v", err)
	}

	httpReq = addParamsAndHeaders(httpReq, c.Api_key, c.Auth_token, projectID, location)
	log.Printf("[DEBUG] GetPlanDetailsFromPlanName request URL: %s", httpReq.URL.String())

	resp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		return "", "", fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("get plan details failed: status %d\nresponse: %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		Code int `json:"code"`
		Data struct {
			Data []struct {
				Name  string `json:"name"`
				Plan  string `json:"plan"`
				Specs struct {
					ID string `json:"id"`
				} `json:"specs"`
			} `json:"data"`
		} `json:"data"`
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", "", fmt.Errorf("failed to decode plan details response: %v\nresponse body: %s", err, string(bodyBytes))
	}

	for _, item := range result.Data.Data {
		if item.Name == planName {
			log.Printf("[INFO] Found plan: PlanID=%s, SlugName=%s", item.Specs.ID, item.Plan)
			return item.Specs.ID, item.Plan, nil
		}
	}

	return "", "", fmt.Errorf("plan name %s not found in template %d", planName, templateID)
}


func (c *Client) UpdateScalerGroup(id string, req *models.UpdateScalerGroupRequest, projectID, location string) error {
	url := c.Api_endpoint + "/scaler/scalegroups/update/" + id + "/"
	log.Printf("[INFO] Sending request to update Scaler Group at: %s", url)

	payloadBuf := new(bytes.Buffer)
	if err := json.NewEncoder(payloadBuf).Encode(req); err != nil {
		return fmt.Errorf("failed to encode update payload: %v", err)
	}

	httpReq, err := http.NewRequest("PUT", url, payloadBuf)
	if err != nil {
		return fmt.Errorf("failed to create PUT request: %v", err)
	}

	httpReq = addParamsAndHeaders(httpReq, c.Api_key, c.Auth_token, projectID, location)
	log.Printf("[DEBUG] UpdateScalerGroup headers: %v", httpReq.Header)

	resp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("update scaler group failed: status %d\nresponse: %s", resp.StatusCode, string(bodyBytes))
	}

	log.Printf("[INFO] Scaler Group updated successfully: ID=%s", id)
	return nil
}

func (c *Client) UpdateDesiredNodeCount(scalerGroupID int, desired int, projectID, location string) error {
	url := c.Api_endpoint + "/scaler/scalegroups/" + strconv.Itoa(scalerGroupID) + "/"
	log.Printf("[INFO] Sending request to update desired node count to %d for Scaler Group ID=%d", desired, scalerGroupID)

	payload := &models.UpdateDesiredNodeCountRequest{Cardinality: desired}
	payloadBuf := new(bytes.Buffer)
	if err := json.NewEncoder(payloadBuf).Encode(payload); err != nil {
		return fmt.Errorf("failed to encode update payload: %w", err)
	}

	httpReq, err := http.NewRequest("PUT", url, payloadBuf)
	if err != nil {
		return fmt.Errorf("failed to create PUT request: %w", err)
	}

	httpReq = addParamsAndHeaders(httpReq, c.Api_key, c.Auth_token, projectID, location)
	log.Printf("[DEBUG] Request URL: %s", httpReq.URL.String())
	log.Printf("[DEBUG] Request headers: %v", httpReq.Header)

	resp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("update desired node count failed: status=%d, body=%s", resp.StatusCode, string(bodyBytes))
	}

	log.Printf("[INFO] Successfully updated desired node count to %d for Scaler Group ID=%d", desired, scalerGroupID)
	return nil
}

func (c *Client) UpdateScalerGroupStatus(id int, status, projectID, location string) error {
	var url string
	idStr := strconv.Itoa(id)

	switch status {
	case "Stopped":
		url = c.Api_endpoint + "/scaler/scalegroups/" + idStr + "/stop/"
	case "Running":
		url = c.Api_endpoint + "/scaler/scalegroups/" + idStr + "/start/"
	default:
		return fmt.Errorf("unsupported status value: %s", status)
	}

	log.Printf("[INFO] Sending request to update scaler group status to: %s at: %s", status, url)

	payloadBuf := new(bytes.Buffer)
	if err := json.NewEncoder(payloadBuf).Encode(struct{}{}); err != nil {
		return fmt.Errorf("failed to encode status update payload: %w", err)
	}

	httpReq, err := http.NewRequest("PUT", url, payloadBuf)
	if err != nil {
		return fmt.Errorf("failed to create PUT request: %w", err)
	}
	httpReq = addParamsAndHeaders(httpReq, c.Api_key, c.Auth_token, projectID, location)

	log.Printf("[DEBUG] UpdateScalerGroupStatus headers: %v", httpReq.Header)

	resp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("update scaler group status to %s failed: status %d\nresponse: %s", status, resp.StatusCode, string(bodyBytes))
	}

	log.Printf("[INFO] Scaler Group status updated successfully to: %s", status)
	return nil
}

func (c *Client) GetVpcDetailsByName(projectID, location, name string) (*models.VPCDetail, error) {
	url := c.Api_endpoint + "/vpc/list/?page_no=1&per_page=100"

	log.Printf("[INFO] Getting VPC details for name %q, projectID: %s, location: %s", name, projectID, location)
	log.Printf("[DEBUG] VPC request URL: %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create VPC request: %v", err)
	}
	req = addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectID, location)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("VPC request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read VPC response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("VPC request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []models.VPCDetail `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse VPC response: %v", err)
	}

	for _, vpc := range result.Data {
		log.Printf("[DEBUG] Checking VPC name: %q", vpc.Name)
		if vpc.Name == name {
			log.Printf("[INFO] Matched VPC found: ID=%d, CIDR=%s", vpc.NetworkID, vpc.IPv4CIDR)
			return &vpc, nil
		}
	}
	return nil, fmt.Errorf("no VPC found with name %q", name)
}

func (c *Client) AttachVPCToScalerGroup(scalerGroupID string, vpcs []models.VPCDetail, projectID, location string) error {
	url := c.Api_endpoint + "/scaler/scalegroups/" + scalerGroupID + "/vpc/action/"
	log.Printf("[INFO] Attaching %d VPC(s) to Scaler Group %s", len(vpcs), scalerGroupID)

	payload := map[string][]models.VPCDetail{
		"vpc": vpcs,
	}

	payloadBuf := new(bytes.Buffer)
	if err := json.NewEncoder(payloadBuf).Encode(payload); err != nil {
		return fmt.Errorf("failed to encode attach VPC payload: %v", err)
	}

	req, err := http.NewRequest("PUT", url, payloadBuf)
	if err != nil {
		return fmt.Errorf("failed to create attach VPC request: %v", err)
	}

	req = addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectID, location)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("attach VPC request failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("attach VPC failed: status %d, body: %s", resp.StatusCode, body)
	}

	log.Printf("[INFO] Successfully attached VPC(s) to Scaler Group %s", scalerGroupID)
	return nil
}

func (c *Client) DetachVPCFromScalerGroup(scalerGroupID, vpcID, projectID, location string) error {
	url := c.Api_endpoint + "/scaler/scalegroups/" + scalerGroupID + "/vpc/action/"
	log.Printf("[INFO] Detaching VPC %s from Scaler Group %s", vpcID, scalerGroupID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create detach VPC request: %v", err)
	}

	q := req.URL.Query()
	q.Add("vpc_id", vpcID)
	req.URL.RawQuery = q.Encode()

	req = addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectID, location)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("detach VPC request failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("detach VPC failed: status %d, body: %s", resp.StatusCode, body)
	}

	log.Printf("[INFO] Successfully detached VPC %s from Scaler Group %s", vpcID, scalerGroupID)
	return nil
}

func (c *Client) GetPublicIPStatus(scaleGroupID, projectID, location string) (*models.PublicIPStatusData, error) {
	url := c.Api_endpoint + "/scaler/scalegroups/" + scaleGroupID + "/public_ip/action/"
	log.Printf("[INFO] Fetching public IP status for Scaler Group ID: %s", scaleGroupID)

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %v", err)
	}
	httpReq = addParamsAndHeaders(httpReq, c.Api_key, c.Auth_token, projectID, location)

	log.Printf("[DEBUG] GetPublicIPStatus request URL: %s", httpReq.URL.String())

	resp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get public IP status failed: status %d\nresponse: %s", resp.StatusCode, string(bodyBytes))
	}

	var result models.PublicIPStatusResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to decode public IP status response: %v\nresponse body: %s", err, string(bodyBytes))
	}

	log.Printf("[INFO] Public IP required: %v", result.Data.IsPublicIPRequired)
	return &result.Data, nil
}

func (c *Client) AttachPublicIP(scaleGroupID, projectID, location string) (*models.PublicIPActionResponse, error) {
	url := c.Api_endpoint + "/scaler/scalegroups/" + scaleGroupID + "/public_ip/action/"
	log.Printf("[INFO] Attaching Public IP to Scaler Group ID: %s", scaleGroupID)

	httpReq, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create PUT request: %v", err)
	}
	httpReq = addParamsAndHeaders(httpReq, c.Api_key, c.Auth_token, projectID, location)

	resp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("attach public IP failed: status %d\nresponse: %s", resp.StatusCode, string(bodyBytes))
	}

	var result models.PublicIPActionResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to decode attach public IP response: %v\nresponse body: %s", err, string(bodyBytes))
	}

	log.Printf("[INFO] Public IP attached: %s", result.Data)
	return &result, nil
}

func (c *Client) DetachPublicIP(scaleGroupID, projectID, location string) (*models.PublicIPActionResponse, error) {
	url := c.Api_endpoint + "/scaler/scalegroups/" + scaleGroupID + "/public_ip/action/"
	log.Printf("[INFO] Detaching Public IP from Scaler Group ID: %s", scaleGroupID)

	httpReq, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create DELETE request: %v", err)
	}
	httpReq = addParamsAndHeaders(httpReq, c.Api_key, c.Auth_token, projectID, location)

	resp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("detach public IP failed: status %d\nresponse: %s", resp.StatusCode, string(bodyBytes))
	}

	var result models.PublicIPActionResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to decode detach public IP response: %v\nresponse body: %s", err, string(bodyBytes))
	}

	log.Printf("[INFO] Public IP detached: %s", result.Data)
	return &result, nil
}

func (c *Client) GetAttachedVPCsForScalerGroup(scalerGroupID, projectID, location string) ([]models.VPCPartial, error) {
	log.Printf("[INFO] Fetching attached VPCs for Scaler Group ID: %s", scalerGroupID)

	url := c.Api_endpoint + "/scaler/scalegroups/" + scalerGroupID + "/vpc/action/"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("[ERROR] Failed to create request: %v", err)
		return nil, err
	}

	req = addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectID, location)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		log.Printf("[ERROR] API call failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read response: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []models.VPCPartial `json:"data"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	log.Printf("[INFO] Successfully fetched %d attached VPC(s) for scaler group %s", len(result.Data), scalerGroupID)
	return result.Data, nil
}

func (c *Client) DetachSecurityGroupFromScalergroup(scalerGroupID string, sgID int, projectID, location string) error {
	url := c.Api_endpoint + "/scaler/scalegroups/security_groups/" + scalerGroupID + "/"
	log.Printf("[INFO] Detaching Security Group %d from Scaler Group %s", sgID, scalerGroupID)

	httpReq, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	httpReq = addParamsAndHeaders(httpReq, c.Api_key, c.Auth_token, projectID, location)
	q := httpReq.URL.Query()
	q.Add("security_group_id", strconv.Itoa(sgID))
	httpReq.URL.RawQuery = q.Encode()

	log.Printf("[DEBUG] DetachSecurityGroup request URL: %s", httpReq.URL.String())
	log.Printf("[DEBUG] DetachSecurityGroup headers: %v", httpReq.Header)

	resp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	log.Printf("[DEBUG] DetachSecurityGroup response code: %d", resp.StatusCode)
	log.Printf("[DEBUG] DetachSecurityGroup response body: %s", string(bodyBytes))

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("detach security group failed: status %d\nresponse: %s", resp.StatusCode, string(bodyBytes))
	}

	log.Printf("[INFO] Successfully detached Security Group %d from Scaler Group %s", sgID, scalerGroupID)
	return nil
}

func (c *Client) AddSecurityGroupToScalergroup(scalerGroupID string, sgID int, projectID, location string) error {
	url := c.Api_endpoint + "/scaler/scalegroups/security_groups/" + scalerGroupID + "/"
	log.Printf("[INFO] Attaching Security Group %d to Scaler Group %s", sgID, scalerGroupID)

	payload := map[string]int{"security_group_id": sgID}
	payloadBuf := new(bytes.Buffer)
	if err := json.NewEncoder(payloadBuf).Encode(payload); err != nil {
		return fmt.Errorf("failed to encode attach payload: %v", err)
	}

	httpReq, err := http.NewRequest("PUT", url, payloadBuf)
	if err != nil {
		return fmt.Errorf("failed to create PUT request: %v", err)
	}

	httpReq = addParamsAndHeaders(httpReq, c.Api_key, c.Auth_token, projectID, location)

	log.Printf("[DEBUG] AttachSecurityGroup request URL: %s", httpReq.URL.String())
	log.Printf("[DEBUG] AttachSecurityGroup headers: %v", httpReq.Header)
	log.Printf("[DEBUG] AttachSecurityGroup payload: %s", payloadBuf.String())

	resp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	log.Printf("[DEBUG] AttachSecurityGroup response code: %d", resp.StatusCode)
	log.Printf("[DEBUG] AttachSecurityGroup response body: %s", string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("attach security group failed: status %d\nresponse: %s", resp.StatusCode, string(bodyBytes))
	}

	log.Printf("[INFO] Security Group %d attached successfully to Scaler Group %s", sgID, scalerGroupID)
	return nil
}
