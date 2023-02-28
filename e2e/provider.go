package e2e

import (
	"awakeninggit.e2enetworks.net/cloud/terraform-provider-e2e/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SERVICE_LOCATION", ""),
			},
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
		},
		ResourcesMap: map[string]*schema.Resource{
			"e2e_node": resourceNode(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	location := d.Get("location").(string)
	api_key := d.Get("api_key").(string)
	auth_token := d.Get("auth_token").(string)
	return client.NewClient(location, api_key, auth_token), nil
}
