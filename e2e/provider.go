package e2e

import (
	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/image"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/node"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/security_group"
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
				Description: "authentication Bearer token should be specified",
			},
			"api_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "https://api.e2enetworks.com/myaccount/api/v1/",
				DefaultFunc: schema.EnvDefaultFunc("SERVICE_API_ENDPOINT", "https://api.e2enetworks.com/myaccount/api/v1"),
				Description: "specify the endpoint , default endpoint is https://api.e2enetworks.com/myaccount/api/v1/",
			},

			"location": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SERVICE_LOCATION", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"e2e_node": node.ResourceNode(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"e2e_node":            node.DataSourceNode(),
			"e2e_images":          image.DataSourceImages(),
			"e2e_security_groups": security_group.DataSourceSecurityGroups(),
			"e2e_ssh_keys":        ssh_key.DataSourceSshKeys(),
			"e2e_vpcs":            vpc.DataSourceVpcs(),
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
