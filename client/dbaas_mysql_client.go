package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func (c *Client) NewMySqlDb(item *models.MySqlCreate, project_id string, location string) (map[string]interface{}, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(item)
	if err != nil {
		return nil, fmt.Errorf(" client | error unmarshalling JSON: %s ", err)
	}

	UrlEndPoint := c.Api_endpoint + "/rds/cluster/"

	req, err := http.NewRequest("POST", UrlEndPoint, &buf)
	if err != nil {
		return nil, fmt.Errorf(" client | error while creating http request: %v", err)
	}

	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)

	response, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf(" error while making http request =: %s ", err)
	}
	err = CheckResponseStatus(response)
	if err != nil {
		return nil, fmt.Errorf(" error while checking response code =: %s ", err)
	}
	defer response.Body.Close()
	resBody, _ := io.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, fmt.Errorf(" error while unmarshling response =: %s ", err)
	}
	return jsonRes, nil
}

func (c *Client) GetMySqlDbaas(mySqlDBaaSId string, project_id string, location string) (*models.DBResponse, error) {
	urlGetDBaaSMySQL := c.Api_endpoint + "rds/cluster/" + mySqlDBaaSId + "/"
	req, err := http.NewRequest("GET", urlGetDBaaSMySQL, nil)
	if err != nil {
		return nil, fmt.Errorf(" client | error while creating http request: %v", err)
	}

	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(" client | error making request in GetMySqlDbaas: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(" client | error reading response body: %v ", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-200 status code: %d - body: %s ", resp.StatusCode, string(body))
	}

	var res models.DBResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf(" GetMySqlDbaas | error unmarshalling JSON: %s ", err)
	}

	return &res, nil
}

func (c *Client) DeleteMySqlDBaaS(mySqlDBaaSId string, project_id string, location string) (map[string]interface{}, error) {
	urlDBaaSMySql := c.Api_endpoint + "rds/cluster/" + mySqlDBaaSId + "/"
	req, err := http.NewRequest("DELETE", urlDBaaSMySql, nil)
	if err != nil {
		return nil, fmt.Errorf(" DeleteMySqlDBaaS | error while creating http request: %s ", err)
	}

	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)

	response, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(" DeleteMySqlDBaaS | error while making http request: %s ", err)
	}
	defer response.Body.Close()
	resBody, _ := io.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, fmt.Errorf(" DeleteMySqlDBaaS | error unmarshalling JSON: %s ", err)
	}
	return jsonRes, nil
}

func (c *Client) ResumeMySqlDBaaS(mySqlDBaaSId string, project_id string, location string) (map[string]interface{}, error) {
	urlDBaaSMySql := c.Api_endpoint + "rds/cluster/" + mySqlDBaaSId + "/resume"

	req, err := http.NewRequest("PUT", urlDBaaSMySql, nil)
	if err != nil {
		return nil, fmt.Errorf(" ResumeMySqlDBaaS | error while creating http request: %s ", err)
	}

	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)

	response, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf(" ResumeMySqlDBaaS | error while making http request: %s ", err)
	}

	defer response.Body.Close()

	resBody, _ := io.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, fmt.Errorf(" ResumeMySqlDBaaS | error unmarshalling JSON: %s ", err)
	}
	return jsonRes, nil
}

func (c *Client) StopMySqlDBaaS(mySqlDBaaSId string, project_id string, location string) (map[string]interface{}, error) {
	urlDBaaSMySql := c.Api_endpoint + "rds/cluster/" + mySqlDBaaSId + "/shutdown"

	req, err := http.NewRequest("PUT", urlDBaaSMySql, nil)
	if err != nil {
		return nil, fmt.Errorf(" StopMySqlDBaaS | error while creating http request: %s ", err)
	}

	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)

	response, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf(" StopMySqlDBaaS  | error while making http request: %s ", err)
	}

	defer response.Body.Close()

	resBody, _ := io.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, fmt.Errorf(" StopMySqlDBaaS | error unmarshalling JSON: %s ", err)
	}
	return jsonRes, nil
}

func (c *Client) RestartMySqlDBaaS(mySqlDBaaSId string, project_id string, location string) (map[string]interface{}, error) {
	urlDBaaSMySql := c.Api_endpoint + "rds/cluster/" + mySqlDBaaSId + "/restart"

	req, err := http.NewRequest("PUT", urlDBaaSMySql, nil)
	if err != nil {
		return nil, fmt.Errorf(" RestartMySqlDBaaS | error while creating http request: %s ", err)
	}

	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)

	response, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf(" RestartMySqlDBaaS | error while making http request: %s ", err)
	}

	defer response.Body.Close()

	resBody, _ := io.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, fmt.Errorf(" RestartMySqlDBaaS | error unmarshalling JSON: %s ", err)
	}
	return jsonRes, nil
}

func (c *Client) AttachVpcToMySql(item *models.AttachVPCPayloadRequest, mySqlDBaaSId string, project_id string, location string) (map[string]interface{}, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(item)
	if err != nil {
		return nil, fmt.Errorf(" AttachVpcToMySql | error while encoding buffer: %s ", err)
	}

	urlDBaaSMySql := c.Api_endpoint + "rds/cluster/" + mySqlDBaaSId + "/vpc-attach/"

	req, err := http.NewRequest("PUT", urlDBaaSMySql, &buf)
	if err != nil {
		return nil, fmt.Errorf(" AttachVpcToMySql | error while creating http request: %s ", err)
	}

	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)

	response, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf(" AttachVpcToMySql | error while making http request: %s ", err)
	}

	defer response.Body.Close()

	resBody, _ := io.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, fmt.Errorf(" AttachVpcToMySql | error unmarshalling JSON: %s ", err)
	}
	return jsonRes, nil
}

func (c *Client) DetachVpcFromMySql(item *models.AttachVPCPayloadRequest, mySqlDBaaSId string, project_id string, location string) (map[string]interface{}, error) {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(item)
	if err != nil {
		return nil, fmt.Errorf(" DetachVpcFromMySql | error while encoding buffer: %s ", err)
	}

	urlDBaaSMySql := c.Api_endpoint + "rds/cluster/" + mySqlDBaaSId + "/vpc-detach/"

	req, err := http.NewRequest("PUT", urlDBaaSMySql, &buf)
	if err != nil {
		return nil, fmt.Errorf(" DetachVpcFromMySql | error while creating http request: %s ", err)
	}

	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)

	response, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf(" DetachVpcFromMySql | error while making http request: %s ", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(response.Body)
		if err != nil {
			return nil, fmt.Errorf("got a non 200 status code: %v", response.StatusCode)
		}
		return nil, fmt.Errorf("got a non 200 status code: %v - %s", response.StatusCode, respBody.String())
	}
	resBody, _ := io.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, fmt.Errorf(" DetachVpcFromMySql | error unmarshalling JSON: %s ", err)
	}
	return jsonRes, nil
}

func (c *Client) AttachPGToMySqlDBaaS(mySqlDBaaSId string, ParameterGroupId string, project_id string, location string) (map[string]interface{}, error) {
	urlDBaaSMySql := c.Api_endpoint + "rds/cluster/" + mySqlDBaaSId + "/parameter-group/" + ParameterGroupId + "/add"

	req, err := http.NewRequest("PUT", urlDBaaSMySql, nil)
	if err != nil {
		return nil, fmt.Errorf(" AttachPGToMySqlDBaaS | error while creating http request: %s ", err)
	}

	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)

	response, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf(" AttachPGToMySqlDBaaS | error while making http request: %s ", err)
	}

	defer response.Body.Close()

	resBody, _ := io.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, fmt.Errorf(" AttachPGToMySqlDBaaS | error unmarshalling JSON: %s ", err)
	}
	return jsonRes, nil
}

func (c *Client) DetachPGFromMySqlDBaaS(mySqlDBaaSId string, ParameterGroupId string, project_id string, location string) (map[string]interface{}, error) {
	urlDBaaSMySql := c.Api_endpoint + "rds/cluster/" + mySqlDBaaSId + "/parameter-group/" + ParameterGroupId + "/detach"

	req, err := http.NewRequest("PUT", urlDBaaSMySql, nil)
	if err != nil {
		return nil, fmt.Errorf(" DetachPGFromMySqlDBaaS | error while creating http request: %s ", err)
	}

	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)

	response, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf(" DetachPGFromMySqlDBaaS | error while making http request: %s ", err)
	}

	defer response.Body.Close()

	resBody, _ := io.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, fmt.Errorf(" DetachPGFromMySqlDBaaS | error unmarshalling JSON: %s ", err)
	}
	return jsonRes, nil
}

func (c *Client) AttachPublicIPToMySql(mySqlDBaaSId string, project_id string, location string) (map[string]interface{}, error) {
	urlDBaaSMySql := c.Api_endpoint + "rds/cluster/" + mySqlDBaaSId + "/public-ip-attach/"

	req, err := http.NewRequest("PUT", urlDBaaSMySql, nil)
	if err != nil {
		return nil, fmt.Errorf(" AttachPublicIPToMySql | error while creating http request: %s ", err)
	}

	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)

	response, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf(" AttachPublicIPToMySql | error while making http request: %s ", err)
	}

	defer response.Body.Close()

	resBody, _ := io.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, fmt.Errorf(" AttachPublicIPToMySql | error unmarshalling JSON: %s ", err)
	}
	return jsonRes, nil
}

func (c *Client) DetachPublicIPFromMySql(mySqlDBaaSId string, project_id string, location string) (map[string]interface{}, error) {
	urlDBaaSMySql := c.Api_endpoint + "rds/cluster/" + mySqlDBaaSId + "/public-ip-detach/"

	req, err := http.NewRequest("PUT", urlDBaaSMySql, nil)
	if err != nil {
		return nil, fmt.Errorf(" DetachPublicIPFromMySql | error while creating http request: %s ", err)
	}

	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)

	response, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf(" DetachPublicIPFromMySql | error while making http request: %s ", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(response.Body)
		if err != nil {
			return nil, fmt.Errorf("got a non 200 status code: %v", response.StatusCode)
		}
		return nil, fmt.Errorf("got a non 200 status code: %v - %s", response.StatusCode, respBody.String())
	}
	resBody, _ := io.ReadAll(response.Body)
	stringresponse := string(resBody)
	resBytes := []byte(stringresponse)
	var jsonRes map[string]interface{}
	err = json.Unmarshal(resBytes, &jsonRes)
	if err != nil {
		return nil, fmt.Errorf(" DetachPublicIPFromMySql | error unmarshalling JSON: %s ", err)
	}
	return jsonRes, nil
}

func (c *Client) UpgradeMySQLPlan(dbaas_id string, template_id int, project_id string, location string) (interface{}, error) {
	dbaas_action := models.MySQlPlanUpgradeAction{
		TemplateID: template_id,
	}

	dbaasAction, err := json.Marshal(dbaas_action)
	if err != nil {
		return nil, fmt.Errorf(" UpgradeMySQLPlan | error unmarshalling JSON: %s ", err)
	}

	urlDBaaSMySql := c.Api_endpoint + "rds/cluster/" + dbaas_id + "/rds-upgrade/"
	req, err := http.NewRequest("PUT", urlDBaaSMySql, bytes.NewBuffer(dbaasAction))
	if err != nil {
		return nil, fmt.Errorf(" UpgradeMySQLPlan | error while creating http request: %s ", err)
	}

	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)

	response, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf(" UpgradeMySQLPlan | error while making http request: %s ", err)
	}

	err = CheckResponseStatus(response)

	if err != nil {
		return nil, fmt.Errorf(" UpgradeMySQLPlan | error inside upgrade MySQL plan after CheckResponseStatus: %s ", err)
	}

	return response, nil
}

func (c *Client) ExpandMySQLDBaaSDisk(dbaas_id string, size int, project_id string, location string) (interface{}, error) {
	dbaas_action := models.MYSQLExpandDisk{
		Size: size,
	}

	dbaasAction, err := json.Marshal(dbaas_action)
	if err != nil {
		return nil, fmt.Errorf(" ExpandMySQLDBaaSDisk | error unmarshalling JSON: %s ", err)
	}

	urlDBaaSMySql := c.Api_endpoint + "rds/cluster/" + dbaas_id + "/disk-upgrade/"
	req, err := http.NewRequest("PUT", urlDBaaSMySql, bytes.NewBuffer(dbaasAction))
	if err != nil {
		return nil, fmt.Errorf(" ExpandMySQLDBaaSDisk | error while creating http request: %s ", err)
	}

	addParamsAndHeaders(req, c.Api_key, c.Auth_token, project_id, location)

	response, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf(" ExpandMySQLDBaaSDisk | error while creating http request: %s ", err)
	}

	err = CheckResponseStatus(response)

	if err != nil {
		return nil, fmt.Errorf(" ExpandMySQLDBaaSDisk | error inside upgrade MySQL plan after CheckResponseStatus: %s ", err)
	}

	return response, nil
}
