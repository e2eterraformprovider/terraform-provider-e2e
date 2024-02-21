package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"fmt"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

type SFSClient struct {
	Api_key      string
	Auth_token   string
	Api_endpoint string
	HttpClient   *http.Client
}

func SFSNewClient(api_key string, auth_token string, api_endpoint string) *SFSClient {
	return &SFSClient{
		Api_key:      api_key,
		Auth_token:   auth_token,
		Api_endpoint: api_endpoint,
		HttpClient:   &http.Client{},
	}
}

func (c *SFSClient)NewSfs(item *models.SfsCreate,project_id string)(map[string]interface{}, error){
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(item)
	if err != nil {
		return nil, err
	}
	UrlSfs := c.Api_endpoint + "efs/"
	log.Printf("[INFO] SFSClient NEWNODE | BEFORE REQUEST")
	req, err := http.NewRequest("POST", UrlSfs, &buf)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("project_id",project_id)
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
func (c *SFSClient) GetSfs(nodeId string , project_id string) (map[string]interface{}, error) {

	UrlSfs := c.Api_endpoint + "efs/" + nodeId + "/"
	req, err := http.NewRequest("GET", UrlSfs, nil)
	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] SFSClient | NODE READ")
	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("contact_person_id", "null")
	params.Add("project_id",project_id)
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "terraform-e2e")

	// log.Printf("req url GetNode = %v", req.URL)
	response, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] SFSClient NODE READ | after response %d", response.StatusCode)
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
		log.Printf("[ERROR] SFSClient GET NDE | error when unmarshalling")
		return nil, err
	}

	return jsonRes, nil
}
func (c *SFSClient) GetSfss(location string,project_id string) (*models.ResponseNodes, error) {

	urlGetSfss := c.Api_endpoint + "efs/"
	req, err := http.NewRequest("GET", urlGetSfss, nil)
	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] CLIENT GET NODES")
	params := req.URL.Query()
    
	params.Add("apikey", c.Api_key)
	params.Add("project_id",project_id)
	params.Add("contact_person_id", "null")
	params.Add("location", location)
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
			log.Printf("GET NODES | INSIDE NO SUCCESS AND ERROR MSG")
			return nil, fmt.Errorf("%v", err)
		}

		return nil, fmt.Errorf("got a non 200 status code: %v - %s", response.StatusCode, respBody.String())
	}
	fmt.Println(response.Body)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	res := models.ResponseNodes{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Printf("[INFO] inside get ssh_keys | error while unmarshlling")
		return nil, err
	}

	return &res, nil
}
