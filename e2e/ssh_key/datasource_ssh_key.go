package ssh_key

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSshKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReadSshKey,
		Schema: map[string]*schema.Schema{
			// Common fields
			tfconstants.AttrRegion:    config.RegionSchema(),
			tfconstants.AttrLocation:  config.LocationSchema(),
			tfconstants.AttrProjectID: config.ProjectIDSchemaComputed(),

			// Resource-specific fields
			tfconstants.AttrLabel: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "label (name) of the SSH key",
				ForceNew:    true,
			},
			// V3 preferred field name
			tfconstants.AttrName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "name of the SSH key",
			},
			tfconstants.AttrSSHKey: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the SSH public key content",
			},
			tfconstants.AttrPublicKey: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the public key material",
			},
			tfconstants.AttrProjectName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "name of the project associated with the SSH key",
			},
			tfconstants.AttrCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the creation date for the SSH key",
			},
		},
	}
}

// dataSourceReadSshKey reads an SSH key data source
func dataSourceReadSshKey(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	goe2eClient := cfg.Goe2eClient()
	var diags diag.Diagnostics

	label := d.Get(tfconstants.AttrLabel).(string)
	if label == "" {
		return diag.Errorf("SSH key label is required")
	}

	log.Printf("[INFO] Reading SSH key data source: label=%s", label)

	// Fetch SSH key by label using goe2e client
	sshKey, _, err := goe2eClient.SSHKeys.GetSSHKeyByLabel(ctx, label)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to find SSH key with label %s: %w", label, err))
	}

	// Check if key was found
	if sshKey == nil {
		return diag.Errorf("SSH key with label %s not found", label)
	}

	log.Printf("[DEBUG] SSH key found: pk=%d, label=%s", sshKey.PK, sshKey.Label)

	// Set the resource ID
	d.SetId(strconv.Itoa(sshKey.PK))

	// Set all fields using helper functions for backward compatibility
	if err := setKeyName(d, sshKey.Label); err != nil {
		return diag.FromErr(err)
	}
	if err := setPublicKey(d, sshKey.SSHKey); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(tfconstants.AttrCreatedAt, sshKey.Timestamp); err != nil {
		return diag.FromErr(err)
	}

	// Set project name (currently empty from goe2e API, but we compute if available)
	// TODO: Fetch project name when API supports it
	if err := d.Set(tfconstants.AttrProjectName, ""); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] SSH key data source read successfully: label=%s, pk=%d", label, sshKey.PK)
	return diags
}
