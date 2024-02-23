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

//production url  -> "https://api.e2enetworks.com/myaccount/api/v1/nodes/"

// groot url -> "https://api-groot.e2enetworks.net/myaccount/api/v1/nodes/"
func (c *Client) NewBlockStorage(item *models.BlockStorageCreate, project_id int, location string) (map[string]interface{}, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(item)
	if err != nil {
		return nil, err
	}

	UrlBlockStorage := c.Api_endpoint + "block_storage/"

	log.Printf("[INFO] CLIENT NEWBLOCKSTORAGE | BEFORE REQUEST")
	req, err := http.NewRequest("POST", UrlBlockStorage, &buf)
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
	resBytes := []byte(stringresponse)

	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, err
	}
	return jsonRes, nil
}

func (c *Client) GetBlockStorage(blockStorageID string, project_id int, location string) (map[string]interface{}, error) {
	urlBlockStorage := c.Api_endpoint + "block_storage/" + blockStorageID + "/"
	req, err := http.NewRequest("GET", urlBlockStorage, nil)
	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] CLIENT | BLOCK STORAGE READ")
	params := req.URL.Query()
	project_id = 328
	log.Printf("%v", project_id)
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
	log.Printf("[INFO] CLIENT BLOCK STORAGE READ | after response %d", response.StatusCode)
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
	log.Printf("%s", stringresponse)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)

	if err != nil {
		log.Printf("[ERROR] CLIENT GET BLOCK STORAGE | error when unmarshalling")
		return nil, err
	}

	return jsonRes, nil
}

func (c *Client) DeleteBlockStorage(blockStorageID string, project_id int, location string) error {

	urlNode := c.Api_endpoint + "block_storage/" + blockStorageID + "/"
	req, err := http.NewRequest("DELETE", urlNode, nil)
	if err != nil {
		return err
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
		return err
	}
	err = CheckResponseStatus(response)
	if err != nil {
		return err
	}
	return nil
}
