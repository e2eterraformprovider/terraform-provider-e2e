package dbaas_postgress

import (
	"context"
	"log"
	"strconv"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourcePostgresDBaaS() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReadPostgres,
		Schema: map[string]*schema.Schema{
			// Common fields - use constants and helpers
			tfconstants.AttrRegion:    config.RegionSchema(),
			tfconstants.AttrLocation:  config.LocationSchema(),
			tfconstants.AttrProjectID: config.ProjectIDSchemaComputed(),

			// Postgres-specific fields
			tfconstants.AttrID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "id of the PostgreSQL DBaaS instance",
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
			"pg_details": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "full parameter group detail map",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			tfconstants.AttrStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "state of the PostgreSQL DBaaS instance",
			},
			"status_actions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "list of available status actions",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			tfconstants.AttrPublicIPAddress: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the PostgreSQL DBaaS instances public ipv4 address",
			},
			tfconstants.AttrPrivateIPAddress: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the PostgreSQL DBaaS instances private ipv4 address",
			},
			"is_public_ip_attached": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "whether a public IP is attached to the PostgreSQL DBaaS instance",
			},
			tfconstants.AttrPlan: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the plan name of the PostgreSQL DBaaS instance",
			},
			"database_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the PostgreSQL version",
			},
			tfconstants.AttrParameterGroupID: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "id of the attached parameter group",
			},
			tfconstants.AttrDiskSize: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the size of the attached disk",
			},
			tfconstants.AttrPowerStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the power status of the PostgreSQL DBaaS instance",
			},
		},
	}
}

func dataSourceReadPostgres(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	var diags diag.Diagnostics

	clusterID := d.Get("id").(string)

	log.Printf("[DEBUG] Reading PostgreSQL datasource for cluster ID: %s", clusterID)

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

	// Get goe2e client for this project/region
	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return diag.Errorf("error creating goe2e client for project (%s), region (%s): %s", projectID, region, err)
	}

	// Get PostgreSQL cluster using goe2e client
	cluster, _, err := goe2eClient.PostgreSQL.GetCluster(ctx, clusterID)
	if err != nil {
		return diag.Errorf("error retrieving PostgreSQL DBaaS (ID: %s): %s", clusterID, err)
	}

	// Check if resource was deleted
	if cluster == nil {
		return diag.Errorf("PostgreSQL cluster %s not found", clusterID)
	}

	// Set resource ID
	d.SetId(strconv.Itoa(cluster.ID))

	// Extract nested structures
	master := cluster.MasterNode
	db := master.Database
	plan := master.Plan
	software := plan.Software

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
	status := cluster.Status
	if status == "SUSPENDED" {
		status = "STOPPED"
	}
	if err := d.Set(tfconstants.AttrStatus, status); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("status_actions", cluster.StatusActions); err != nil {
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
	if err := d.Set(tfconstants.AttrDiskSize, master.Disk); err != nil {
		return diag.FromErr(err)
	}

	// Set plan field
	if err := d.Set(tfconstants.AttrPlan, plan.Name); err != nil {
		return diag.FromErr(err)
	}

	// Set software version
	if err := d.Set("database_version", software.Version); err != nil {
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

	log.Printf("[DEBUG] Successfully read PostgreSQL datasource: %s (ID: %s)", cluster.Name, clusterID)

	return diags
}
