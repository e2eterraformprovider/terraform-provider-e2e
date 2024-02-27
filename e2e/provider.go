package e2e

import (
	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/image"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/node"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/sfs"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/reserve_ip"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/ssh_key"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{

			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SERVICE_API_KEY", ""),
				Description: "valied api key required ",
			},
			"auth_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SERVICE_AUTH_TOKEN", ""),
				Description: "Valied authentication Bearer token required",
			},
			"api_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "https://api.e2enetworks.com/myaccount/api/v1/",
				DefaultFunc: schema.EnvDefaultFunc("SERVICE_API_ENDPOINT", "https://api.e2enetworks.com/myaccount/api/v1"),
				Description: "specify the endpoint , default endpoint is https://api.e2enetworks.com/myaccount/api/v1/",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"e2e_node":  node.ResourceNode(),
			"e2e_image": image.ResourceImage(),
			"e2e_sfs": sfs.ResourceSfs(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"e2e_node":   node.DataSourceNode(),
			"e2e_images": image.DataSourceImages(),
			//"example_security_groups": security_group.DataSourceSecurityGroups(),
			"e2e_ssh_keys":    ssh_key.DataSourceSshKeys(),
			"e2e_vpcs":        vpc.DataSourceVpcs(),
			"e2e_reserve_ips": reserve_ip.DataSourceReserveIps(),
			"e2e_nodes":       node.DataSourceNodes(),
			"e2e_sfss":        sfs.DataSourceSfs(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	api_key := d.Get("api_key").(string)
	auth_token := d.Get("auth_token").(string)
	api_endpoint := d.Get("api_endpoint").(string)
	return client.NewClient(api_key, auth_token, api_endpoint), nil
}