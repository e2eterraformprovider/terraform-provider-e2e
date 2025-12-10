package ssh_key

import (
	"context"
	"log"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	e2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSshKeys() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReadSshKeys,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			// Common fields
			e2econstants.AttrRegion:    config.RegionSchema(),
			e2econstants.AttrLocation:  config.LocationSchema(),
			e2econstants.AttrProjectID: config.ProjectIDSchemaComputed(),

			// Resource-specific fields
			"ssh_key_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "list of SSH keys which can be used to launch resources",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pk": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "id of the SSH key",
						},
						"label": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "label (name) of the SSH key",
						},
						"ssh_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the SSH key",
						},
						e2econstants.AttrCreatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the creation date for the SSH key",
						},
					},
				},
			},
		},
	}
}

func dataSourceReadSshKeys(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	cfg := m.(*config.Config)

	log.Printf("[INFO] Inside sshkeys data source ")

	// Get region with provider default support
	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get project_id with provider default support
	projectID, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return diag.Errorf("error creating goe2e client: %v", err)
	}

	sshKeys, _, err := goe2eClient.SSHKeys.ListSSHKeys(ctx)
	if err != nil {
		return diag.Errorf("error finding ssh keys: %v", err)
	}
	d.Set("ssh_key_list", flattenSshKeys(sshKeys))
	d.SetId("ssh_key_list")

	return diags
}

func flattenSshKeys(sshKeyList []goe2e.SSHKey) []interface{} {

	if len(sshKeyList) > 0 {
		ois := make([]interface{}, len(sshKeyList))

		for i, sshKey := range sshKeyList {
			oi := make(map[string]interface{})
			oi["label"] = sshKey.Label
			oi["ssh_key"] = sshKey.SSHKey
			oi["pk"] = sshKey.PK
			// Map API's "Timestamp" field to Terraform's "created_at" for consistency
			oi[e2econstants.AttrCreatedAt] = sshKey.Timestamp
			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}
