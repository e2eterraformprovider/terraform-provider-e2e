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

func (c Client) CreatePostgressDB(payload models.DBCreateRequest, project_id string, location string) (map[string]interface{}, error) {

	url := c.Api_endpoint + "rds/cluster/"

	payloadBuf := new(bytes.Buffer)
	if err := json.NewEncoder(payloadBuf).Encode(payload); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, payloadBuf)
	if err != nil {
		return nil, err
	}
	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)

	response, err := c.HttpClient.Do(req)

	log.Printf("\n\n[INFO] CLIENT NEW DBAAS POSTGRESS | STATUS_CODE: %+v  \n\n", response)

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

func (c Client) GetPostgressDB(id string, project_id string, location string) (map[string]interface{}, error) {

	url := c.Api_endpoint + "rds/cluster/" + id + "/"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)

	response, err := c.HttpClient.Do(req)

	log.Printf("\n\n[INFO] READ DBAAS POSTGRESS | STATUS_CODE: %+v  \n\n", response)

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

func (c Client) DeletePostgressDB(id string, project_id string, location string) error {

	url := c.Api_endpoint + "rds/cluster/" + id + "/"
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)

	response, err := c.HttpClient.Do(req)

	log.Printf("\n\n[INFO] DELETE DBAAS POSTGRESS | STATUS_CODE: %+v  \n\n", response)

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

func (c Client) StopPostgressDB(id string, project_id string, location string) error {

	url := c.Api_endpoint + "rds/cluster/" + id + "/shutdown"
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return err
	}
	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)

	response, err := c.HttpClient.Do(req)

	log.Printf("\n\n[INFO] STOP DBAAS POSTGRESS | STATUS_CODE: %+v  \n\n", response)

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

func (c Client) StartPostgressDB(id string, project_id string, location string) error {

	url := c.Api_endpoint + "rds/cluster/" + id + "/resume"
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return err
	}
	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)

	response, err := c.HttpClient.Do(req)

	log.Printf("\n\n[INFO] START DBAAS POSTGRESS | STATUS_CODE: %+v  \n\n", response)

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

func (c Client) RestartPostgressDB(id string, project_id string, location string) error {

	url := c.Api_endpoint + "rds/cluster/" + id + "/restart"
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return err
	}
	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)

	response, err := c.HttpClient.Do(req)

	log.Printf("\n\n[INFO] RESTART DBAAS POSTGRESS | STATUS_CODE: %+v  \n\n", response)

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

func (c Client) AttachPublicIpPostgressDB(id string, project_id string, location string) error {

	url := c.Api_endpoint + "rds/cluster/" + id + "/public-ip-attach/"
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return err
	}
	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)

	response, err := c.HttpClient.Do(req)

	log.Printf("\n\n[INFO] STOP DBAAS POSTGRESS | STATUS_CODE: %+v  \n\n", response)

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

func (c Client) DetachPublicIpPostgressDB(id string, project_id string, location string) error {

	url := c.Api_endpoint + "rds/cluster/" + id + "/public-ip-detach/"
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return err
	}
	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)

	response, err := c.HttpClient.Do(req)

	log.Printf("\n\n[INFO] STOP DBAAS POSTGRESS | STATUS_CODE: %+v  \n\n", response)

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

func (c Client) AttachVPCPostgressDB(payload models.AttachVPCPayloadRequest, id string, project_id string, location string) error {

	url := c.Api_endpoint + "rds/cluster/" + id + "/vpc-attach/"
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

	log.Printf("\n\n[INFO] ATTACH VPC DBAAS POSTGRESS | STATUS_CODE: %+v  \n\n", response)

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

func (c Client) DetachVPCPostgressDB(payload models.AttachVPCPayloadRequest, id string, project_id string, location string) error {

	url := c.Api_endpoint + "rds/cluster/" + id + "/vpc-detach/"
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

	log.Printf("\n\n[INFO] DETACH VPC DBAAS POSTGRESS | STATUS_CODE: %+v  \n\n", response)

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

func (c *Client) UpgradePostgressPlan(dbaas_id string, template_id int, project_id string, location string) (interface{}, error) {
	dbaas_action := models.PostgressPlanUpgradeAction{
		TemplateID: template_id,
	}
	dbaasAction, _ := json.Marshal(dbaas_action)

	url := c.Api_endpoint + "rds/cluster/" + dbaas_id + "/rds-upgrade/"
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(dbaasAction))
	if err != nil {
		log.Printf("[INFO] error inside upgrade dbaas postgress plan")
	}
	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)

	response, err := c.HttpClient.Do(req)

	log.Printf("[INFO] CLIENT UPGRADE POSTGRESS PLAN | STATUS_CODE: %+v\n", response)
	if err == nil {
		err = CheckResponseStatus(response)
	}

	if err != nil {
		log.Printf("[INFO] error inside upgrade postgress plan after checkresponse")
		return nil, err
	}
	return response, err
}

func (c *Client) UpdateParameterGroup(dbaas_id string, pg_id string, project_id string, location string) error {

	url := c.Api_endpoint + "rds/cluster/" + dbaas_id + "/parameter-group/" + pg_id + "/add"
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		log.Printf("[INFO] error inside UpdateParameterGroup dbaas postgress plan")
	}
	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)

	response, err := c.HttpClient.Do(req)

	log.Printf("[INFO] CLIENT UpdateParameterGroup POSTGRESS  | STATUS_CODE: %+v \n\n", response)
	if err == nil {
		err = CheckResponseStatus(response)
	}

	if err != nil {
		log.Printf("[ERROR] error inside UpdateParameterGroup after checkresponse")
		return err
	}
	return nil
}

func (c *Client) UpgradeDiskStorage(dbaas_id string, size int, project_id string, location string) error {

	dbaas_action := models.PostgressDiskAction{
		Size: size,
	}
	dbaasAction, _ := json.Marshal(dbaas_action)

	url := c.Api_endpoint + "rds/cluster/" + dbaas_id + "/disk-upgrade/"
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(dbaasAction))
	if err != nil {
		log.Printf("[ERROR] error inside upgrade disk storage")
	}
	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)

	response, err := c.HttpClient.Do(req)

	log.Printf("CLIENT UpgradeDiskStorage POSTGRESS  | STATUS_CODE: %+v \n\n", response)
	if err == nil {
		err = CheckResponseStatus(response)
	}

	if err != nil {
		log.Printf("[ERROR] error inside UpgradeDiskStorage after checkresponse")
		return err
	}
	return nil
}
