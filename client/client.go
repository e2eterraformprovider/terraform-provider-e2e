package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"io/ioutil"
	"net/http"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

//production url  -> "https://api.e2enetworks.com/myaccount/api/v1/nodes/"

//groot url -> "https://api-groot.e2enetworks.net/myaccount/api/v1/nodes/"

type Client struct {
	Api_key      string
	Auth_token   string
	Api_endpoint string
	Location     string
	HttpClient   *http.Client
}

func NewClient(api_key string, auth_token string, location string, api_endpoint string) *Client {
	return &Client{

		Api_key:      api_key,
		Auth_token:   auth_token,
		Location:     location,
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
	//params.Add("contact_person_id", "null")
	params.Add("location", c.Location)
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")

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

	// body, err := c.httpRequest(fmt.Sprintf("item/%v", name), "GET", bytes.Buffer{})
	urlNode := c.Api_endpoint + "nodes/" + nodeId + "/"
	req, err := http.NewRequest("GET", urlNode, nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()

	params.Add("apikey", c.Api_key)
	params.Add("contact_person_id", "null")
	params.Add("location", c.Location)
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")

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

func (c *Client) UpdateNode(nodeId string, action string) error {

	//node_action := buildNodeUpdateRequestBody(actionType,action)

	node_action := models.NodeAction{
		Type: action,
	}
	nodeAction, err := json.Marshal(node_action)
	url := c.Api_endpoint + "nodes/" + nodeId + "/actions/"
	log.Printf("[info] %s", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(nodeAction))
	if err != nil {
		return err
	}

	params := req.URL.Query()

	params.Add("apikey", c.Api_key)
	params.Add("location", c.Location)
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.URL.RawQuery = params.Encode()
	response, err := c.HttpClient.Do(req)

	if err != nil {
		return err
	}
	log.Printf("[INFO] inside update %s %d", action, response.StatusCode)
	if response.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(response.Body)
		if err != nil {
			return fmt.Errorf("got a non 200 status code: %v", response.StatusCode)
		}
		return fmt.Errorf("got a non 200 status code: %v - %s", response.StatusCode, respBody.String())
	}
	defer response.Body.Close()
	resBody, _ := ioutil.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}

	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return err
	}

	return err
}

func (c *Client) DeleteNode(nodeId string) error {
	// _, err := c.httpRequest(fmt.Sprintf("item/%s", itemName), "DELETE", bytes.Buffer{})
	urlNode := c.Api_endpoint + "nodes/" + nodeId + "/"
	req, err := http.NewRequest("DELETE", urlNode, nil)
	if err != nil {
		return err
	}

	params := req.URL.Query()

	params.Add("apikey", c.Api_key)
	params.Add("contact_person_id", "null")
	params.Add("location", c.Location)
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")

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

// func buildNodeUpdateRequestBody (actionType string , action interface{})models.NodeAction{

// 	if actionType=="power_status"{
// 		 node_action := models.NodeAction{
// 			Type: action.(string),
// 		}
// 		return node_action
// 	}
// 	if actionType=="lock_vm"{
// 		node_action:=models.NodeAction{
// 			Type:action.(bool),
// 		}
// 		return node_action
// 	}

// }
