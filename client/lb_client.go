package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

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
