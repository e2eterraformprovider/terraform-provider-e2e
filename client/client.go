package client

import (
	"bytes"
	"encoding/json"
	"fmt"

	"io/ioutil"
	"net/http"

	"github.com/devteametwoe/terraform-provider-e2e/models"
)

type Client struct {
	Location   string
	Api_key    string
	Auth_token string
	HttpClient *http.Client
}

func NewClient(location string, api_key string, auth_token string) *Client {
	return &Client{
		Location:   location,
		Api_key:    api_key,
		Auth_token: auth_token,
		HttpClient: &http.Client{},
	}
}

func (c *Client) NewNode(item *models.Node) (map[string]interface{}, error) {
	// fmt.Println("hii")
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(item)
	if err != nil {
		return nil, err
	}
	url := "https://api-groot.e2enetworks.net/myaccount/api/v1/nodes/"
	req, err := http.NewRequest("POST", url, &buf)
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

func (c *Client) GetNode(name string) (map[string]interface{}, error) {

	// body, err := c.httpRequest(fmt.Sprintf("item/%v", name), "GET", bytes.Buffer{})
	url := "https://api-groot.e2enetworks.net/myaccount/api/v1/nodes/" + name + "/"
	req, err := http.NewRequest("GET", url, nil)
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
	//resbody := responseSchema.Response{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, err
	}

	return jsonRes, nil
}

func (c *Client) UpdateNode(node *models.Node) error {
	// buf := bytes.Buffer{}
	// err := json.NewEncoder(&buf).Encode(node)
	// if err != nil {
	// 	return err
	// }
	// // _, err = c.httpRequest(fmt.Sprintf("item/%s", item.Name), "PUT", buf)
	// url := "https://api-groot.e2enetworks.net/myaccount/api/v1/nodes/" + node.Name
	// req, err := http.NewRequest("PUT", url, &buf)
	// if err != nil {
	// 	return err
	// }

	// params := req.URL.Query()

	// params.Add("apikey", c.Api_key)
	// params.Add("contact_person_id", "null")
	// params.Add("location", c.Location)
	// req.URL.RawQuery = params.Encode()
	// req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	// req.Header.Add("Content-Type", "application/json")

	// response, err := c.HttpClient.Do(req)
	// if err != nil {
	// 	return err
	// }
	// if response.StatusCode != http.StatusOK {
	// 	respBody := new(bytes.Buffer)
	// 	_, err := respBody.ReadFrom(response.Body)
	// 	if err != nil {
	// 		return fmt.Errorf("got a non 200 status code: %v", response.StatusCode)
	// 	}
	// 	return fmt.Errorf("got a non 200 status code: %v - %s", response.StatusCode, respBody.String())
	// }
	// if err != nil {
	// 	return err
	// }
	return nil
}

func (c *Client) DeleteNode(nodeName string) error {
	// _, err := c.httpRequest(fmt.Sprintf("item/%s", itemName), "DELETE", bytes.Buffer{})
	url := "https://api-groot.e2enetworks.net/myaccount/api/v1/nodes/" + nodeName + "/"
	req, err := http.NewRequest("DELETE", url, nil)
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
	// if err != nil {
	// 	return err
	// }
	return nil
}
