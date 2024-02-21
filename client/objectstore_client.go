package client

import (
	"net/http"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func (client *Client) setParamsAndHeaders(request *http.Request, location string) (*http.Request, error) {
	params := request.URL.Query()
	params.Add("apikey", client.Api_key)
	params.Add("contact_person_id", "null")
	params.Add("location", location)
	request.URL.RawQuery = params.Encode()
	request.Header.Add("Authorization", "Bearer "+client.Auth_token)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("User-Agent", "terraform-e2e")
	return request, nil
}

func (client *Client) CreateBucket(buckets *models.ObjectStore) []interface{} {

}
