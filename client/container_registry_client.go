package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func (c *Client) CreateContainerRegistry(req *models.CreateContainerRegistryRequest, projectID, location string) (*models.CreateContainerRegistryData, error) {
	url := c.Api_endpoint + "/container_registry/setup-container-registry/"

	payloadBuf := new(bytes.Buffer)
	if err := json.NewEncoder(payloadBuf).Encode(req); err != nil {
		return nil, fmt.Errorf("failed to encode create payload: %v", err)
	}

	httpReq, err := http.NewRequest("POST", url, payloadBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	httpReq = addParamsAndHeaders(httpReq, c.Api_key, c.Auth_token, projectID, location)

	log.Printf("[DEBUG] CreateContainerRegistry request URL: %s", httpReq.URL.String())
	log.Printf("[DEBUG] CreateContainerRegistry request payload: %+v", req)

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
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var response models.CreateContainerRegistryResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v\nResponse Body: %s", err, string(bodyBytes))
	}

	return &response.Data, nil
}

func (c *Client) GetContainerRegistryProjects(projectID, location string) ([]models.ContainerRegistryProject, error) {
	url := c.Api_endpoint + "/container_registry/projects-details/?page=1&page_size=100"

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	httpReq = addParamsAndHeaders(httpReq, c.Api_key, c.Auth_token, projectID, location)

	log.Printf("[DEBUG] GetContainerRegistryProjects Request URL: %s", httpReq.URL.String())

	resp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var response models.GetContainerRegistryProjectsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return response.Data, nil
}

func (c *Client) DeleteContainerRegistry(crProjectID, projectName, userID, projectID, location string) error {
	url := c.Api_endpoint + "/container_registry/setup-container-registry/"

	httpReq, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create DELETE request: %v", err)
	}

	httpReq = addParamsAndHeaders(httpReq, c.Api_key, c.Auth_token, projectID, location)

	params := httpReq.URL.Query()
	params.Add("cr_project_id", crProjectID)
	params.Add("user_id", userID)
	params.Add("project_name", projectName)
	httpReq.URL.RawQuery = params.Encode()

	log.Printf("[DEBUG] DeleteContainerRegistry Request URL: %s", httpReq.URL.String())

	resp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("delete request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var response models.DeleteContainerRegistryResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode delete response: %v", err)
	}
	log.Printf("[DEBUG] DeleteContainerRegistry response status: %s", response.Data.Status)

	return nil
}

func (c *Client) UpdateContainerRegistry(projectName, preventVul, severity, projectID, location string) error {
	url := c.Api_endpoint + "/container_registry/setup-container-registry/"

	payload := map[string]string{
		"project_name": projectName,
		"prevent_vul":  preventVul,
		"severity":     severity,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal update payload: %v", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create PUT request: %v", err)
	}

	req = addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectID, location)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform update request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("update failed: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	log.Printf("[DEBUG] Successfully updated Container Registry for project_name=%s", projectName)
	return nil
}
