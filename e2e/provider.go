package e2e

import (
	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	blockstorage "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/block_storage"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/image"

	// "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/loadbalancer"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/node"
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
				Default:     "https://api-stage.e2enetworks.net/myaccount/api/v1/",
				DefaultFunc: schema.EnvDefaultFunc("SERVICE_API_ENDPOINT", "https://api-stage.e2enetworks.net/myaccount/api/v1/"),
				Description: "specify the endpoint , default endpoint is https://api-stage.e2enetworks.net/myaccount/api/v1/",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"example_node":         node.ResourceNode(),
			"example_image":        image.ResourceImage(),
			"example_blockstorage": blockstorage.ResourceBlockStorage(),
			// "example_loadbalancer": loadbalancer.ResourceLoadBalancer(),
			// "example_vpc":          vpc.ResouceVpc(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"example_node":         node.DataSourceNode(),
			"example_images":       image.DataSourceImages(),
			"example_blockstorage": blockstorage.DataSourceBlockStorage(),
			//"example_security_groups": security_group.DataSourceSecurityGroups(),
			"example_ssh_keys":    ssh_key.DataSourceSshKeys(),
			"example_vpcs":        vpc.DataSourceVpcs(),
			"example_reserve_ips": reserve_ip.DataSourceReserveIps(),
			"example_nodes":       node.DataSourceNodes(),
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
