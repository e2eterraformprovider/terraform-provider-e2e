package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (c *Client) GetSoftwareId(projectID string, location string, name string, version string) (int, error) {
	url := c.Api_endpoint + "rds/plans/"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return -1, err
	}

	req = addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectID, location)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		log.Printf("[ERROR] error inside GetSoftwareId: %v", err)
		return -1, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] reading GetSoftwareId response: %v", err)
		return -1, err
	}
	log.Println("[DEBUG] Raw body:", string(body))

	var res models.PlanResponse
	if err := json.Unmarshal(body, &res); err != nil {
		log.Printf("[ERROR] unmarshalling GetSoftwareId response: %v", err)
		return -1, err
	}

	for _, item := range res.Data.DatabaseEngines {
		if item.EngineName == name && item.EngineVersion == version {
			return item.EngineID, nil
		}
	}

	log.Printf("[INFO] Software ID not found for name: %s, version: %s", name, version)
	return -1, errors.New("matching engine not found")
}

func (c *Client) GetTemplateId(projectID string, location string, plan string, softwareID string) (int, error) {
	url := c.Api_endpoint + "rds/plans/"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return -1, err
	}

	req = addParamsAndHeaders(req, c.Api_key, c.Auth_token, projectID, location)

	
	q := req.URL.Query()
	q.Add("software_id", softwareID)
	req.URL.RawQuery = q.Encode()

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		log.Printf("[ERROR] error inside GetTemplateId: %v", err)
		return -1, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] reading GetTemplateId response: %v", err)
		return -1, err
	}

	var res models.PlanResponse
	if err := json.Unmarshal(body, &res); err != nil {
		log.Printf("[ERROR] unmarshalling GetTemplateId response: %v", err)
		return -1, err
	}

	for _, item := range res.Data.TemplatePlans {
		if item.PlanName == plan {
			return item.PlanTemplateID, nil
		}
	}

	log.Printf("[INFO] Template ID not found for plan name: %s", plan)
	return -1, errors.New("matching plan not found")
}

func (c *Client) ExpandVpcList(d *schema.ResourceData, vpcList []interface{}) ([]models.VPC, error) {
	var vpcDetails []models.VPC

	for _, id := range vpcList {
		vpcResp, err := c.GetVpc(strconv.Itoa(id.(int)), d.Get("project_id").(string), d.Get("location").(string))
		if err != nil {
			return nil, err
		}
		data := vpcResp.Data
		if data.State != "Active" {
			return nil, fmt.Errorf("cannot attach VPC %s: VPC is in '%s' state", data.Name, data.State)
		}

		vpcDetails = append(vpcDetails, models.VPC{
			Network_id: data.Network_id,
			VpcName:    data.Name,
			Ipv4_cidr:  data.Ipv4_cidr,
		})
	}
	return vpcDetails, nil
}

func (c *Client) ExpandMariaDBVpcList(vpcIDs []string, projectID, location string) ([]models.VPCMetadata, error) {
	var vpcDetails []models.VPCMetadata

	for _, id := range vpcIDs {
		vpcResp, err := c.GetVpc(id, projectID, location)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch VPC details for ID %s: %v", id, err)
		}

		vpc := vpcResp.Data
		if vpc.State != "Active" {
			return nil, fmt.Errorf("cannot attach VPC %s: VPC is in '%s' state", id, vpc.State)
		}

		networkID := fmt.Sprintf("%.0f", vpc.Network_id)

		vpcDetails = append(vpcDetails, models.VPCMetadata{
			NetworkID: networkID,
			VPCName:   vpc.Name,
			IPv4CIDR:  vpc.Ipv4_cidr,
		})
	}

	return vpcDetails, nil
}
