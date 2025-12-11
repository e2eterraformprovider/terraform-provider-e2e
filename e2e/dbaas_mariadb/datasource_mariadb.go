package dbaas_mariadb

import (
	"context"
	"log"
	"strconv"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceMariaDB() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReadMariaDB,
		Schema: map[string]*schema.Schema{
			// Common fields - use constants and helpers
			tfconstants.AttrRegion:    config.RegionSchema(),
			tfconstants.AttrLocation:  config.LocationSchema(),
			tfconstants.AttrProjectID: config.ProjectIDSchemaComputed(),

			// MariaDB-specific fields
			tfconstants.AttrID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "id of the MariaDB DBaaS instance",
			},
			tfconstants.AttrName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "name of the MariaDB DBaaS instance",
			},
			tfconstants.AttrDatabaseID: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "id of the database",
			},
			tfconstants.AttrDatabaseName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "name of the database",
			},
			tfconstants.AttrDatabaseUser: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the database username",
			},
			tfconstants.AttrStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "state of the MariaDB DBaaS instance",
			},
			tfconstants.AttrPublicIPAddress: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the MariaDB DBaaS instances public ipv4 address",
			},
			tfconstants.AttrPrivateIPAddress: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the MariaDB DBaaS instances private ipv4 address",
			},
			"is_public_ip_attached": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "whether a public IP is currently attached",
			},
			"disk": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the disk size of the MariaDB DBaaS instance",
			},
			tfconstants.AttrPlan: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the plan name of the MariaDB DBaaS instance",
			},
			"software_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the MariaDB software version",
			},
			tfconstants.AttrParameterGroupID: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "id of the attached parameter group",
			},
			tfconstants.AttrPowerStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the power status of the MariaDB DBaaS instance (e.g., Running, Stopped)",
			},
		},
	}
}

func dataSourceReadMariaDB(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	goe2eClient := cfg.Goe2eClient()
	var diags diag.Diagnostics

	// Fetch basic identifiers
	clusterID := d.Get("id").(string)

	log.Printf("[DEBUG] Reading MariaDB datasource for cluster ID: %s", clusterID)

	// Get MariaDB cluster using goe2e client
	maria, _, err := goe2eClient.MariaDB.GetMariaDB(ctx, clusterID)
	if err != nil {
		return diag.Errorf("error retrieving MariaDB DBaaS (ID: %s): %s", clusterID, err)
	}

	// Check if resource was deleted
	if maria == nil {
		return diag.Errorf("MariaDB cluster %s not found", clusterID)
	}

	// Extract nested structures
	master := maria.MasterNode
	db := master.Database
	plan := master.Plan
	software := plan.Software

	// Set resource ID
	d.SetId(strconv.Itoa(maria.ID))

	// Set basic fields
	if err := d.Set(tfconstants.AttrName, maria.Name); err != nil {
		return diag.FromErr(err)
	}

	// Set database fields
	if err := d.Set(tfconstants.AttrDatabaseID, db.ID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(tfconstants.AttrDatabaseName, db.Database); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(tfconstants.AttrDatabaseUser, db.Username); err != nil {
		return diag.FromErr(err)
	}

	// Normalize status (SUSPENDED → STOPPED)
	status := maria.Status
	if status == "SUSPENDED" {
		status = "STOPPED"
	}
	if err := d.Set(tfconstants.AttrStatus, status); err != nil {
		return diag.FromErr(err)
	}

	// Set network fields
	if err := d.Set(tfconstants.AttrPublicIPAddress, master.PublicIPAddress); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(tfconstants.AttrPrivateIPAddress, master.PrivateIPAddress); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_public_ip_attached", master.PublicIPAddress != ""); err != nil {
		return diag.FromErr(err)
	}

	// Set disk field
	if err := d.Set("disk", master.Disk); err != nil {
		return diag.FromErr(err)
	}

	// Set plan field
	if err := d.Set(tfconstants.AttrPlan, plan.Name); err != nil {
		return diag.FromErr(err)
	}

	// Set software version
	if err := d.Set("software_version", software.Version); err != nil {
		return diag.FromErr(err)
	}

	// Set power status
	if err := d.Set(tfconstants.AttrPowerStatus, master.Status); err != nil {
		return diag.FromErr(err)
	}

	// Set parameter group ID if present
	if db.PGDetail.ID != 0 {
		if err := d.Set(tfconstants.AttrParameterGroupID, db.PGDetail.ID); err != nil {
			return diag.FromErr(err)
		}
	}

	log.Printf("[DEBUG] Successfully read MariaDB datasource: %s (ID: %s)", maria.Name, clusterID)

	return diags
}
