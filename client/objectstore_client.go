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

func (client *Client) setParamsAndHeaders(request *http.Request, location string, project_id string) *http.Request {
	params := request.URL.Query()
	params.Add("apikey", client.Api_key)
	params.Add("contact_person_id", "null")
	params.Add("location", location)
	params.Add("project_id", project_id)
	request.URL.RawQuery = params.Encode()
	request.Header.Add("Authorization", "Bearer "+client.Auth_token)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("User-Agent", "terraform-e2e")
	return request
}

func (client *Client) CreateBucket(buckets *models.ObjectStorePayload) (map[string]interface{}, error) {
	payload_buffer := bytes.Buffer{}
	error_while_encoding := json.NewEncoder(&payload_buffer).Encode(buckets)
	if error_while_encoding != nil {
		return nil, error_while_encoding
	}
	BucketCreateUrl := client.Api_endpoint + "buckets/" + buckets.BucketName
	create_request, error := http.NewRequest("POST", BucketCreateUrl, &payload_buffer)
	if error != nil {
		return nil, error
	}
	create_request = client.setParamsAndHeaders(create_request, buckets.Region, fmt.Sprint(buckets.ProjectID))
	response, error := client.HttpClient.Do(create_request)
	if error != nil {
		return nil, error
	}

	error = CheckResponseStatus(response)
	if error != nil {
		return nil, error
	}
	defer response.Body.Close()
	responseBody, _ := ioutil.ReadAll(response.Body)
	stringresponse := string(responseBody)
	responseBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	error = json.Unmarshal(responseBytes, &jsonRes)
	if error != nil {
		return nil, error
	}
	return jsonRes, nil
}

func (client *Client) GetBuckets(location string, project_id string) (*models.ResponseBuckets, error) {

	urlGetNodes := client.Api_endpoint + "buckets/"
	readrequest, err := http.NewRequest("GET", urlGetNodes, nil)
	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] CLIENT GET BUCKETS")
	readrequest = client.setParamsAndHeaders(readrequest, location, project_id)

	response, err := client.HttpClient.Do(readrequest)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(response.Body)
		if err != nil {
			log.Printf("GET BUCKETS | INSIDE NO SUCCESS AND ERROR MSG")
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
	res := models.ResponseBuckets{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Printf("[INFO] inside get ssh_keys | error while unmarshlling")
		return nil, err
	}
	return &res, nil
}

func (client *Client) GetBucket(bucket_name string, location string, project_id string) (map[string]interface{}, error) {

	urlGetNodes := client.Api_endpoint + "buckets/" + bucket_name
	readrequest, err := http.NewRequest("GET", urlGetNodes, nil)
	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] CLIENT GET BUCKET")
	readrequest = client.setParamsAndHeaders(readrequest, location, project_id)

	response, err := client.HttpClient.Do(readrequest)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(response.Body)
		if err != nil {
			log.Printf("GET BUCKET | INSIDE NO SUCCESS AND ERROR MSG")
			return nil, fmt.Errorf("%v", err)
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
	log.Printf("%s", stringresponse)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		log.Printf("[INFO] inside get ssh_keys | error while unmarshlling")
		return nil, err
	}
	return jsonRes, nil
}
