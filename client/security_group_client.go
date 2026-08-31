package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func (c *Client) GetSecurityGroupList(project_id string, location string) (map[string]interface{}, error) {

	urlSecurityGroups := c.Api_endpoint + "security_group/"
	req, err := http.NewRequest("GET", urlSecurityGroups, nil)
	if err != nil {
		return nil, err
	}
	projectIDInt, err := strconv.Atoi(project_id)
	if err != nil {
		return nil, err
	}
	addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectIDInt, location)
	response, err := c.HttpClient.Do(req)
	if err != nil {
		log.Printf("[ERROR] error inside GetSecurityGroupList")
		return nil, err
	}
	log.Printf("[INFO] CLIENT SECURITY GROUPS LIST READ | response code %d", response.StatusCode)
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
		log.Printf("[ERROR] CLIENT GetSecurityGroupList | error when unmarshalling | %s", err)
		return nil, err
	}
	return jsonRes, nil
}

func (c *Client) GetSecurityGroup(name string, project_id string, location string) (map[string]interface{}, error) {
	list, err := c.GetSecurityGroupList(project_id, location)
	if err != nil {
		return nil, err
	}
	items, ok := list["data"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	for _, item := range items {
		itemMap := item.(map[string]interface{})

		if itemMap["name"] == name {

			return itemMap, nil
		}
	}
	return nil, fmt.Errorf("security group %s not found", name)
}

func (c *Client) CreateSecurityGroups(payload models.SecurityGroupCreateRequest, project_id string, location string) error {
	url := c.Api_endpoint + "security_group/"

	payloadBuf := new(bytes.Buffer)
	if err := json.NewEncoder(payloadBuf).Encode(payload); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, payloadBuf)
	if err != nil {
		return err
	}
	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)

	response, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	log.Printf("\n\n[INFO] NEW CREATE SECURITY GROUP | STATUS_CODE: %+v\n\n", response.StatusCode)

	if err := CheckResponseStatus(response); err != nil {
		return err
	}

	return nil
}
func (c *Client) UpdateSecurityGroups(payload models.SecurityGroupUpdateRequest, id string, project_id string, location string) error {

	url := c.Api_endpoint + "security_group/" + id + "/"

	payloadBuf := new(bytes.Buffer)
	if err := json.NewEncoder(payloadBuf).Encode(payload); err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, payloadBuf)
	if err != nil {
		return err
	}
	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)

	response, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	log.Printf("\n\n[INFO] UPDATE SECURITY GROUP | STATUS_CODE: %+v\n\n", response.StatusCode)

	if err := CheckResponseStatus(response); err != nil {
		return err
	}

	return nil
}

func (c *Client) MakeDefaultSecurityGroup(id string, project_id string, location string) error {

	url := c.Api_endpoint + "security_group/" + id + "/mark-default/"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)

	response, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	log.Printf("\n\n[INFO] MAKE DEFAULT SECURITY GROUP | STATUS_CODE: %+v\n\n", response.StatusCode)

	if err := CheckResponseStatus(response); err != nil {
		return err
	}

	return nil
}

func (c *Client) DetachSecurityGroup(item *models.UpdateSecurityGroups, vm_id int, project_id string, location string) (map[string]interface{}, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(item)
	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] CLIENT SECURITY GROUP DETACH | BEFORE REQUEST")

	vmIDInString := strconv.Itoa(vm_id)
	urlNode := c.Api_endpoint + "security_group/" + vmIDInString + "/detach/"
	req, err := http.NewRequest("POST", urlNode, &buf)
	if err != nil {
		return nil, err
	}
	projectIDInt, error := strconv.Atoi(project_id)
	if error != nil {
		return nil, error
	}
	addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectIDInt, location)
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

func (c *Client) AttachSecurityGroup(item *models.UpdateSecurityGroups, vm_id int, project_id string, location string) (map[string]interface{}, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(item)
	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] CLIENT SECURITY GROUP ATTACH | BEFORE REQUEST")

	vmIDInString := strconv.Itoa(vm_id)
	urlNode := c.Api_endpoint + "security_group/" + vmIDInString + "/attach/"
	req, err := http.NewRequest("POST", urlNode, &buf)
	if err != nil {
		return nil, err
	}
	projectIDInt, error := strconv.Atoi(project_id)
	if error != nil {
		return nil, error
	}
	addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectIDInt, location)
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

func (c *Client) DeleteSecurityGroup(id string, project_id string, location string) error {

	url := c.Api_endpoint + "security_group/" + id + "/"

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)

	response, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	log.Printf("\n\n[INFO] DELETE SECURITY GROUP | STATUS_CODE: %+v\n\n", response.StatusCode)

	if err := CheckResponseStatus(response); err != nil {
		return err
	}

	return nil
}
