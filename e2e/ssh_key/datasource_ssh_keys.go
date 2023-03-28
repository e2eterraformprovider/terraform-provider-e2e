package ssh_key

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

func DataSourceSshKeys() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ssh_key_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pk": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"label": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ssh_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"timestamp": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
		ReadContext: dataSourceReadSshKeys,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func dataSourceReadSshKeys(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	apiClient := m.(*client.Client)
	log.Printf("[INFO] Inside sshkeys data source ")
	Response, err := apiClient.GetSshKeys()
	if err != nil {
		return diag.Errorf("error finding ssh keys")
	}
	d.Set("ssh_key_list", flattenSshKeys(&Response.Data))
	d.SetId("ssh_key_list")

	return diags
}

func flattenSshKeys(sshKeyList *[]models.SshKey) []interface{} {

	if sshKeyList != nil {
		ois := make([]interface{}, len(*sshKeyList), len(*sshKeyList))

		for i, sshKey := range *sshKeyList {
			oi := make(map[string]interface{})
			oi["label"] = sshKey.Label
			oi["ssh_key"] = sshKey.Ssh_key
			oi["pk"] = sshKey.Pk
			oi["timestamp"] = sshKey.Timestamp
			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}
