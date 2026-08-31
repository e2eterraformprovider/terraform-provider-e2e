package dbaas_mysql

import (
	"fmt"
	"strconv"
	"time"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/e2eterraformprovider/terraform-provider-e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ExpandVpcList(d *schema.ResourceData, vpc_list []interface{}, apiClient *client.Client) ([]models.VPC, error) {
	var vpc_details []models.VPC

	for _, id := range vpc_list {
		vpc_detail, err := apiClient.GetVpc(strconv.Itoa(id.(int)), d.Get("project_id").(string), d.Get("location").(string))
		if err != nil {
			return nil, fmt.Errorf("error while fetching vpc: %s", err)
		}
		data := vpc_detail.Data
		if data.State != "Active" {
			return nil, fmt.Errorf("Can not attach vpc currently, vpc is in %s state", data.State)
		}
		r := models.VPC{
			Network_id: data.Network_id,
			VpcName:    data.Name,
			Ipv4_cidr:  data.Ipv4_cidr,
		}

		vpc_details = append(vpc_details, r)
	}
	return vpc_details, nil
}

func WaitForPoweringOffOnDBaaS(m interface{}, dbaasID string, project_id string, location string) error {
	apiClient := m.(*client.Client)

	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		time.Sleep(constants.WAIT_TIMEOUT * time.Second)

		dbaasInfo, err := apiClient.GetMySqlDbaas(dbaasID, project_id, location)
		if err != nil {
			return fmt.Errorf("error while fetching dbaas instance details: %s", err)
		}

		status := dbaasInfo.Data.Status

		if status == "SUSPENDED" {
			return nil
		}
	}

	return fmt.Errorf("timeout: MySQL DBaaS did not reach SUSPENDED state in time, please wait for some more time and then hit TERRAFORM APPLY again")
}
