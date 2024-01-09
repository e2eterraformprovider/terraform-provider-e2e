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

func (c *Client) AddParamsAndHeader(req *http.Request, location string) (*http.Request, error) {
	params := req.URL.Query()
	params.Add("apikey", c.Api_key)
	params.Add("contact_person_id", "null")
	params.Add("location", location)
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "terraform-e2e")

	return req, nil
}

func (c *Client) NewLoadBalancer(item *models.LoadBalancerCreate) (map[string]interface{}, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(item)
	if err != nil {
		return nil, err
	}
	UrlNode := c.Api_endpoint + "appliances/load-balancers/"
	log.Printf("[INFO] CLIENT NEWLOADBALANCER| BEFORE REQUEST")
	log.Println("========================LOAD BALANCER PAYLOAD FORMED BEFORE ===========================")
	log.Println(buf.String())
	buf, err = RemoveExtraKeysLoadBalancer(&buf)
	log.Println("========================LOAD BALANCER PAYLOAD FORMED AFTER ===========================")
	log.Println(buf.String())
	if err != nil {
		return nil, err
	}
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

func (c *Client) GetLoadBalancerInfo(lbId string, location string) (map[string]interface{}, error) {
	urlLbInfo := c.Api_endpoint + "appliances/" + lbId + "/"
	req, err := http.NewRequest("GET", urlLbInfo, nil)
	if err != nil {
		return nil, err
	}

	log.Printf("[INFO] CLIENT | LOAD BALANCER READ")

	req, err = c.AddParamsAndHeader(req, location)
	if err != nil {
		return nil, err
	}

	response, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] CLIENT LOAD BALANCER READ | response code %d", response.StatusCode)
	if response.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(response.Body)
		if err != nil {
			log.Printf("======================INSIDE ERROR FROM READFROM======================")
			return nil, fmt.Errorf("got a non 200 status code: %v", response.StatusCode)
		}
		log.Printf("=======================ERROR FROM API NOT GETTING 200 ====================")
		return nil, fmt.Errorf("got a non 200 status code: %v - %s", response.StatusCode, respBody.String())
	}
	log.Printf("======================NOW DEFER CLOSE RESPONSE BODY ===========================")
	defer response.Body.Close()
	resBody, _ := ioutil.ReadAll(response.Body)
	stringresponse := string(resBody)
	log.Printf("================STRING RESPONSE=========================%s", stringresponse)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		log.Printf("[ERROR] CLIENT GET LOAD BALANCER | error when unmarshalling | %s", err)
		return nil, err
	}

	return jsonRes, nil
}

func (c *Client) DeleteLoadBalancer(lbId string, location string) error {
	urlLbInfo := c.Api_endpoint + "appliances/" + lbId + "/"
	req, err := http.NewRequest("DELETE", urlLbInfo, nil)
	if err != nil {
		return err
	}

	log.Printf("[INFO] CLIENT | LOAD BALANCER DELETE")

	req, err = c.AddParamsAndHeader(req, location)
	if err != nil {
		return err
	}

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
