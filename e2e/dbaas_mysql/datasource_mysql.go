package dbaas_mysql

import (
	"context"
	"fmt"
	"strconv"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	e2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceMySQLDBaaS() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReadMySQL,
		Schema: map[string]*schema.Schema{
			// Common fields - use constants and helpers
			e2econstants.AttrRegion:    config.RegionSchema(),
			e2econstants.AttrLocation:  config.LocationSchema(),
			e2econstants.AttrProjectID: config.ProjectIDSchemaComputed(),

			// MySQL-specific fields
			e2econstants.AttrID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "id of the MySQL DBaaS instance",
			},
			e2econstants.AttrDatabaseID: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "id of the database",
			},
			e2econstants.AttrDatabaseName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "name of the database",
			},
			e2econstants.AttrDatabaseUser: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the database username",
			},
			e2econstants.AttrStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "state of the MySQL DBaaS instance",
			},
			e2econstants.AttrPublicIPAddress: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the MySQL DBaaS instances public ipv4 address",
			},
			e2econstants.AttrPrivateIPAddress: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the MySQL DBaaS instances private ipv4 address",
			},
			"is_public_ip_attached": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "whether a public IP is attached to the MySQL DBaaS instance",
			},
			e2econstants.AttrDisk: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the disk size of the MySQL DBaaS instance",
			},
			e2econstants.AttrPlan: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the plan name of the MySQL DBaaS instance",
			},
			"database_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the MySQL version",
			},
			e2econstants.AttrParameterGroupID: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "id of the attached parameter group",
			},
			e2econstants.AttrPowerStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the power status of the MySQL DBaaS instance",
			},
		},
	}
}

func dataSourceReadMySQL(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	goe2eClient := cfg.Goe2eClient()
	var diags diag.Diagnostics

	dbaasID := d.Get("id").(string)

	// Get MySQL cluster using goe2e client
	mysql, _, err := goe2eClient.DBaaSMySQL.GetCluster(ctx, dbaasID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while fetching MySQL DBaaS instance details: %s", err))
	}

	// Handle case where cluster was not found
	if mysql == nil {
		return diag.Errorf("MySQL DBaaS instance (ID: %s) not found", dbaasID)
	}

	// Extract nested data
	master := mysql.MasterNode
	db := master.Database
	plan := master.Plan
	software := plan.Software

	// Set all attributes
	d.SetId(strconv.Itoa(mysql.ID))
	d.Set(e2econstants.AttrDatabaseID, db.ID)
	d.Set(e2econstants.AttrDatabaseName, db.Database)
	d.Set(e2econstants.AttrDatabaseUser, db.Username)
	d.Set(e2econstants.AttrStatus, mysql.Status)
	d.Set(e2econstants.AttrPublicIPAddress, master.PublicIPAddress)
	d.Set(e2econstants.AttrPrivateIPAddress, master.PrivateIPAddress)
	d.Set("is_public_ip_attached", master.PublicIPAddress != "")
	d.Set(e2econstants.AttrDisk, master.Disk)
	d.Set(e2econstants.AttrPlan, plan.Name)
	d.Set("database_version", software.Version)

	// Handle PGDetail safely (check if ID is set)
	if db.PGDetail.ID != 0 {
		d.Set(e2econstants.AttrParameterGroupID, db.PGDetail.ID)
	} else {
		d.Set(e2econstants.AttrParameterGroupID, 0)
	}

	d.Set(e2econstants.AttrPowerStatus, master.Status)

	return diags
}
