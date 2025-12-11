package ssh_key

import (
	"context"
	"log"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
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
			tfconstants.AttrRegion:    config.RegionSchema(),
			tfconstants.AttrLocation:  config.LocationSchema(),
			tfconstants.AttrProjectID: config.ProjectIDSchemaComputed(),

			// Resource-specific fields
			tfconstants.AttrSSHKeyList: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "list of SSH keys which can be used to launch resources",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						tfconstants.AttrPK: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "id of the SSH key",
						},
						tfconstants.AttrLabel: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "label (name) of the SSH key",
						},
						tfconstants.AttrSSHKey: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the SSH key",
						},
						tfconstants.AttrCreatedAt: {
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
	d.Set(tfconstants.AttrSSHKeyList, flattenSshKeys(sshKeys))
	d.SetId(tfconstants.AttrSSHKeyList)

	return diags
}

func flattenSshKeys(sshKeyList []goe2e.SSHKey) []interface{} {

	if len(sshKeyList) > 0 {
		ois := make([]interface{}, len(sshKeyList))

		for i, sshKey := range sshKeyList {
			oi := make(map[string]interface{})
			oi[tfconstants.AttrLabel] = sshKey.Label
			oi[tfconstants.AttrSSHKey] = sshKey.SSHKey
			oi[tfconstants.AttrPK] = sshKey.PK
			// Map API's "Timestamp" field to Terraform's "created_at" for consistency
			oi[tfconstants.AttrCreatedAt] = sshKey.Timestamp
			ois[i] = oi
		}

		return ois
	}
	return make([]interface{}, 0)
}
