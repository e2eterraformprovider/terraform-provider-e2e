package vpc

import (
	"context"
	// "fmt"
	"log"
	// "math"
	// "regexp"

	// "strconv"
	//"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"

	// "github.com/hashicorp/terraform-plugin-log"
	// "github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceVpcs() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"vpc_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_id": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"pool_size": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv4_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_active": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
		ReadContext: dataSourceReadVpcs,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func dataSourceReadVpcs(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	apiClient := m.(*client.Client)
	log.Printf("[INFO] Inside vpcs data source ")
	Response, err := apiClient.GetVpcs()
	if err != nil {
		return diag.Errorf("error finding vpcs ")
	}
	d.Set("vpc_list", flattenVpcs(&Response.Data))
	d.SetId("vpc_list")

	return diags
}

func flattenVpcs(vpcList *[]models.Vpc) []interface{} {

	if vpcList != nil {
		ois := make([]interface{}, len(*vpcList), len(*vpcList))

		for i, vpc := range *vpcList {
			oi := make(map[string]interface{})
			oi["network_id"] = vpc.Network_id
			oi["pool_size"] = vpc.Pool_size
			oi["created_at"] = vpc.Created_at
			oi["name"] = vpc.Name
			oi["is_active"] = vpc.Is_active
			oi["gateway_ip"] = vpc.Gateway_ip
			oi["ipv4_cidr"] = vpc.Ipv4_cidr
			oi["network_id"] = vpc.Network_id
			oi["state"] = vpc.State
			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}
