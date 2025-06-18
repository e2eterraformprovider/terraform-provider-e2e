package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func (c *Client) AddSshKey(item models.AddSshKey, project_id string) (map[string]interface{}, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(item)
	if err != nil {
		return nil, err
	}
	UrlSshKey := c.Api_endpoint + "ssh_keys/"
	log.Printf("[INFO] %s", UrlSshKey)
	req, err := http.NewRequest("POST", UrlSshKey, &buf)
	if err != nil {
		return nil, err
	}
	params := req.URL.Query()

	params.Add("apikey", c.Api_key)
	params.Add("project_id", project_id)
	params.Add("location", item.Location) //

	req.URL.RawQuery = params.Encode()
	// 	//  Log full URL and location
	// log.Printf("[DEBUG] Full request URL: %s", req.URL.String())
	// log.Printf("[DEBUG] Location param in AddSshKey: %s", item.Location)

	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "terraform-e2e")
	response, err := c.HttpClient.Do(req)
	log.Printf("inside add ssh key req = %+v, res = %+v", req, response)
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

func (c *Client) GetSshKey(label string, project_id string, location string) (map[string]interface{}, error) {
	UrlSshKey := c.Api_endpoint + "ssh_keys/"
	req, err := http.NewRequest("GET", UrlSshKey, nil)
	if err != nil {
		return nil, err
	}
	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("project_id", project_id)
	params.Add("label", label)
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
	log.Printf("=====================RESPONSE_GET_SSH==============, %+v", resBody)
	stringresponse := string(resBody)
	log.Printf("=====================RESPONSE_GET_SSH==============, %+v", stringresponse)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	log.Printf("=====================RESPONSE_GET_SSH==============, %+v", jsonRes)
	if err != nil {
		return nil, err
	}
	return jsonRes, nil
}

func (c *Client) DeleteSshKey(pk string, project_id string, location string) error {
	UrlSshKey := c.Api_endpoint + "delete_ssh_key/" + pk + "/"
	req, err := http.NewRequest("DELETE", UrlSshKey, nil)
	if err != nil {
		return err
	}
	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("project_id", project_id)
	params.Add("location", location)
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "terraform-e2e")

	//log.Printf("[INFO] Sending DELETE request to: %s", req.URL.String())

	response, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Check if delete was successful
	if response.StatusCode != 200 && response.StatusCode != 204 {
		bodyBytes, _ := io.ReadAll(response.Body)
		return fmt.Errorf("failed to delete SSH key, status: %d, response: %s", response.StatusCode, string(bodyBytes))
	}

	log.Printf("[INFO] SSH key delete API returned status: %d", response.StatusCode)

	return nil
}

func (c *Client) GetSshKeys(location string, project_id string) (*models.SshKeyResponse, error) {

	urlSshKeys := c.Api_endpoint + "ssh_keys/"
	req, err := http.NewRequest("GET", urlSshKeys, nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("location", location)
	params.Add("project_id", project_id)
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "terraform-e2e")
	response, err := c.HttpClient.Do(req)
	if err != nil {
		log.Printf("[INFO] error inside get ssh keys")
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	res := models.SshKeyResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Printf("[INFO] inside get ssh_keys | error while unmarshlling")
		return nil, err
	}
	return &res, nil
}
func (c *Client) GetSshKeyByPk(pk string, project_id string, location string) (*models.SshKey, error) {
	url := strings.TrimRight(c.Api_endpoint, "/") + "/ssh_keys/"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("project_id", project_id)
	params.Add("location", location)
	params.Add("pk", pk)
	req.URL.RawQuery = params.Encode()

	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "terraform-e2e")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("SSH key with ID %s not found", pk)
	}

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status: %d, response: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data struct {
		Code  int             `json:"code"`
		Data  []models.SshKey `json:"data"`
		Error []interface{}   `json:"error"`
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	if len(data.Data) == 0 {
		return nil, fmt.Errorf("SSH key with ID %s not found in response", pk)
	}

	return &data.Data[0], nil
}
