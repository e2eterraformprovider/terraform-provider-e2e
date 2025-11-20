package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

// CreateFaasNamespace creates a new FaaS namespace
func (c *Client) CreateFaasNamespace(namespace string, projectID string, location string) (*models.FaasNamespaceResponse, error) {
	namespaceReq := models.FaasNamespaceCreate{
		Name: namespace,
	}

	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(namespaceReq)
	if err != nil {
		return nil, err
	}

	urlNamespace := c.Api_endpoint + "faas/namespace"
	req, err := http.NewRequest("POST", urlNamespace, &buf)
	if err != nil {
		return nil, err
	}

	req = addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectID, location)
	log.Printf("[INFO] CLIENT CREATE FAAS NAMESPACE | Request: %+v", req)

	response, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(response.Body)
		if err != nil {
			return nil, fmt.Errorf("got a non 200/201 status code: %v", response.StatusCode)
		}
		return nil, fmt.Errorf("got a non 200/201 status code: %v - %s", response.StatusCode, respBody.String())
	}

	resBody, _ := ioutil.ReadAll(response.Body)
	var res models.FaasNamespaceResponse
	err = json.Unmarshal(resBody, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// DeleteFaasNamespace deletes a FaaS namespace
func (c *Client) DeleteFaasNamespace(namespace string, projectID string, location string) error {
	urlNamespace := c.Api_endpoint + "faas/namespace"
	req, err := http.NewRequest("DELETE", urlNamespace, nil)
	if err != nil {
		return err
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("project_id", projectID)
	params.Add("location", location)
	params.Add("namespace", namespace)
	req.URL.RawQuery = params.Encode()

	SetBasicHeaders(c.Auth_token, req)
	log.Printf("[INFO] CLIENT DELETE FAAS NAMESPACE | Request: %+v", req)

	response, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(response.Body)
		if err != nil {
			return fmt.Errorf("got a non 200 status code: %v", response.StatusCode)
		}
		return fmt.Errorf("got a non 200 status code: %v - %s", response.StatusCode, respBody.String())
	}

	return nil
}

// CreateFaasFunction creates a new FaaS function
func (c *Client) CreateFaasFunction(fn *models.FaasFunctionCreate, projectID string, location string) (*models.FaasFunctionResponse, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(fn)
	if err != nil {
		return nil, err
	}

	urlFunction := c.Api_endpoint + "faas/functions"
	req, err := http.NewRequest("POST", urlFunction, &buf)
	if err != nil {
		return nil, err
	}

	req = addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectID, location)
	log.Printf("[INFO] CLIENT CREATE FAAS FUNCTION | Request: %+v", req)

	response, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(response.Body)
		if err != nil {
			return nil, fmt.Errorf("got a non 200/201 status code: %v", response.StatusCode)
		}
		return nil, fmt.Errorf("got a non 200/201 status code: %v - %s", response.StatusCode, respBody.String())
	}

	resBody, _ := ioutil.ReadAll(response.Body)
	log.Printf("[INFO] CLIENT CREATE FAAS FUNCTION | Response Body: %s", string(resBody))

	var res models.FaasFunctionResponse
	err = json.Unmarshal(resBody, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// GetFaasFunction retrieves a FaaS function by ID
func (c *Client) GetFaasFunction(functionID string, projectID string, location string) (*models.FaasFunctionResponse, error) {
	urlFunction := c.Api_endpoint + "faas/function/" + functionID + "/"
	req, err := http.NewRequest("GET", urlFunction, nil)
	if err != nil {
		return nil, err
	}

	req = addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectID, location)
	log.Printf("[INFO] CLIENT GET FAAS FUNCTION | Request: %+v", req)

	response, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if response.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(response.Body)
		if err != nil {
			return nil, fmt.Errorf("got a non 200 status code: %v", response.StatusCode)
		}
		return nil, fmt.Errorf("got a non 200 status code: %v - %s", response.StatusCode, respBody.String())
	}

	resBody, _ := ioutil.ReadAll(response.Body)
	var res models.FaasFunctionResponse
	err = json.Unmarshal(resBody, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// UpdateFaasFunction updates an existing FaaS function
func (c *Client) UpdateFaasFunction(functionID string, fn *models.FaasFunctionUpdate, projectID string, location string) (*models.FaasFunctionResponse, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(fn)
	if err != nil {
		return nil, err
	}

	urlFunction := c.Api_endpoint + "faas/function/" + functionID + "/"
	req, err := http.NewRequest("PUT", urlFunction, &buf)
	if err != nil {
		return nil, err
	}

	req = addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectID, location)
	log.Printf("[INFO] CLIENT UPDATE FAAS FUNCTION | Request: %+v", req)

	response, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(response.Body)
		if err != nil {
			return nil, fmt.Errorf("got a non 200 status code: %v", response.StatusCode)
		}
		return nil, fmt.Errorf("got a non 200 status code: %v - %s", response.StatusCode, respBody.String())
	}

	resBody, _ := ioutil.ReadAll(response.Body)
	var res models.FaasFunctionResponse
	err = json.Unmarshal(resBody, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// DeleteFaasFunction deletes a FaaS function
func (c *Client) DeleteFaasFunction(functionID string, projectID string, location string) error {
	urlFunction := c.Api_endpoint + "faas/function/" + functionID + "/"
	req, err := http.NewRequest("DELETE", urlFunction, nil)
	if err != nil {
		return err
	}

	req = addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectID, location)
	log.Printf("[INFO] CLIENT DELETE FAAS FUNCTION | Request: %+v", req)

	response, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNoContent {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(response.Body)
		if err != nil {
			return fmt.Errorf("got a non 200/204 status code: %v", response.StatusCode)
		}
		return fmt.Errorf("got a non 200/204 status code: %v - %s", response.StatusCode, respBody.String())
	}

	return nil
}

// GetFaasLogs retrieves logs for a FaaS function
func (c *Client) GetFaasLogs(functionID string, projectID string, location string) (*models.FaasLogsResponse, error) {
	urlLogs := c.Api_endpoint + "faas/logs/" + functionID + "/"
	req, err := http.NewRequest("GET", urlLogs, nil)
	if err != nil {
		return nil, err
	}

	req = addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectID, location)
	log.Printf("[INFO] CLIENT GET FAAS LOGS | Request: %+v", req)

	response, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(response.Body)
		if err != nil {
			return nil, fmt.Errorf("got a non 200 status code: %v", response.StatusCode)
		}
		return nil, fmt.Errorf("got a non 200 status code: %v - %s", response.StatusCode, respBody.String())
	}

	resBody, _ := ioutil.ReadAll(response.Body)
	var res models.FaasLogsResponse
	err = json.Unmarshal(resBody, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
