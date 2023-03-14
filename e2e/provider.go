package e2e

import (
	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/node"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{

			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SERVICE_API_KEY", ""),
			},
			"auth_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SERVICE_AUTH_TOKEN", ""),
			},
			"api_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "https://api.e2enetworks.com/myaccount/api/v1/",
				DefaultFunc: schema.EnvDefaultFunc("SERVICE_API_ENDPOINT", "https://api.e2enetworks.com/myaccount/api/v1"),
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
			"e2e_node": node.DataSourceNode(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	api_key := d.Get("api_key").(string)
	auth_token := d.Get("auth_token").(string)
	location := d.Get("location").(string)
	api_endpoint := d.Get("api_endpoint").(string)
	return client.NewClient(api_key, auth_token, location, api_endpoint), nil
}
