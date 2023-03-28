package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"io/ioutil"
	"net/http"

	"github.com/devteametwoe/terraform-provider-e2e/models"
)

type Client struct {
	Api_key      string
	Auth_token   string
	Api_endpoint string
	HttpClient   *http.Client
}

func NewClient(api_key string, auth_token string, api_endpoint string) *Client {
	return &Client{

		Api_key:      api_key,
		Auth_token:   auth_token,
		Api_endpoint: api_endpoint,
		HttpClient:   &http.Client{},
	}
}

func (c *Client) NewNode(item *models.Node) (map[string]interface{}, error) {

	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(item)
	if err != nil {
		return nil, err
	}
	UrlNode := c.Api_endpoint + "nodes/"
	log.Printf("[INFO] %s", UrlNode)
	req, err := http.NewRequest("POST", UrlNode, &buf)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()

	params.Add("apikey", c.Api_key)
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "terraform-e2e")
	response, err := c.HttpClient.Do(req)

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

func (c *Client) GetNode(nodeId string) (map[string]interface{}, error) {

	urlNode := c.Api_endpoint + "nodes/" + nodeId + "/"
	req, err := http.NewRequest("GET", urlNode, nil)
	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] NEW BUILD READ")
	params := req.URL.Query()

	params.Add("apikey", c.Api_key)
	params.Add("contact_person_id", "null")
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")

	req.Header.Add("User-Agent", "terraform-e2e")

	response, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(response.Body)
		if err != nil {
			return nil, fmt.Errorf("got a non 200 status code: %v", response.StatusCode)
		}
		return nil, fmt.Errorf("got a non 200 status code: %v - %s", response.StatusCode, respBody.String())
	}
	fmt.Println(response.Body)

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

func (c *Client) UpdateNode(nodeId string, action string, nodeName string) (interface{}, error) {

	node_action := models.NodeAction{
		Type: action,
		Name: nodeName,
	}
	nodeAction, err := json.Marshal(node_action)
	url := c.Api_endpoint + "nodes/" + nodeId + "/actions/"
	log.Printf("[info] %s", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(nodeAction))
	if err != nil {
		return nil, err
	}
	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "terraform-e2e")
	req.URL.RawQuery = params.Encode()
	response, err := c.HttpClient.Do(req)

	if err != nil {

		return nil, err
	}
	log.Printf("[INFO] inside update %s %d", action, response.StatusCode)
	if response.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(response.Body)
		if err != nil {
			return nil, fmt.Errorf("got a non 200 status code: %v", response.StatusCode)
		}
		return nil, fmt.Errorf("got a non 200 status code: %v - %s", response.StatusCode, respBody.String())
	}
	defer response.Body.Close()
	resBody, _ := ioutil.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return jsonRes, err
	}
	return jsonRes, err
}

func (c *Client) DeleteNode(nodeId string) error {

	urlNode := c.Api_endpoint + "nodes/" + nodeId + "/"
	req, err := http.NewRequest("DELETE", urlNode, nil)
	if err != nil {
		return err
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("contact_person_id", "null")
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "terraform-e2e")
	response, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
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

func (c *Client) GetSavedImages() (*models.ImageListResponse, error) {

	urlImages := c.Api_endpoint + "images/" + "saved-images" + "/"

	req, err := http.NewRequest("GET", urlImages, nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("contact_person_id", "null")
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "terraform-e2e")
	response, err := c.HttpClient.Do(req)
	log.Printf("[INFO] inside client saved image before request hit")
	if err != nil {
		log.Printf("[INFO] error inside get image")
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	res := models.ImageListResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
func (c *Client) GetSecurityGroups() (*models.SecurityGroupsResponse, error) {

	urlSecurityGroups := c.Api_endpoint + "security_group/"
	req, err := http.NewRequest("GET", urlSecurityGroups, nil)
	if err != nil {
		return nil, err
	}
	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "terraform-e2e")
	response, err := c.HttpClient.Do(req)
	if err != nil {
		log.Printf("[INFO] error inside get security groups")
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	res := models.SecurityGroupsResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Printf("[INFO] inside get security groups | error while unmarshlling")
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetSshKeys() (*models.SshKeyResponse, error) {

	urlSshKeys := c.Api_endpoint + "ssh_keys/"
	req, err := http.NewRequest("GET", urlSshKeys, nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
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

func (c *Client) GetVpcs() (*models.VpcsResponse, error) {

	urlGetVpcs := c.Api_endpoint + "vpc/" + "list/"
	req, err := http.NewRequest("GET", urlGetVpcs, nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)

	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "terraform-e2e")

	response, err := c.HttpClient.Do(req)

	if err != nil {
		log.Printf("[INFO] error inside get vpcs")
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	res := models.VpcsResponse{}

	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Printf("[INFO] inside get vpcs | error while unmarshlling")
		return nil, err
	}
	return &res, nil
}
