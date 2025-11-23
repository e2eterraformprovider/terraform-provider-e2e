package dbaas_mariadb

import (
	"context"
	"strconv"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceMariaDB() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReadMariaDB,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the MariaDB cluster",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project ID associated with the MariaDB cluster",
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Location/region of the MariaDB cluster (e.g. 'Delhi')",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the MariaDB cluster",
			},
			"database_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Database ID inside the MariaDB cluster",
			},
			"database_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the MariaDB database",
			},
			"database_user": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Username for the MariaDB database",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "High-level status of the MariaDB cluster",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP address of the master node (if attached)",
			},
			"private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Private IP address of the master node",
			},
			"is_public_ip_attached": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether a public IP is currently attached",
			},
			"disk": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Disk size of the master node (e.g. '400 GB')",
			},
			"plan": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the plan associated with the master node",
			},
			"software_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "MariaDB software version",
			},
			"parameter_group_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Parameter group ID attached to the DB",
			},
			"power_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VM status of the master node (e.g., 'Running', 'Stopped')",
			},
		},
	}
}

func dataSourceReadMariaDB(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	apiClient := cfg.Client()
	var diags diag.Diagnostics

	// Fetch basic identifiers
	clusterID := d.Get("id").(string)
	projectID := d.Get("project_id").(string)
	location := d.Get("location").(string)

	// Call API
	maria, err := apiClient.ReadMariaDB(clusterID, projectID, location)
	if err != nil {
		return diag.FromErr(err)
	}

	// Extract nested structures
	master := maria.MasterNode
	db := master.Database
	plan := master.Plan
	software := plan.Software

	// Set fields
	d.SetId(strconv.Itoa(maria.ID))
	d.Set("name", maria.Name)
	d.Set("database_id", db.ID)
	d.Set("database_name", db.Database)
	d.Set("database_user", db.Username)
	d.Set("status", maria.Status)
	d.Set("public_ip", master.PublicIPAddress)
	d.Set("private_ip", master.PrivateIPAddress)
	d.Set("is_public_ip_attached", master.PublicIPAddress != "")
	d.Set("disk", master.Disk)
	d.Set("plan", plan.Name)
	d.Set("software_version", software.Version)
	d.Set("power_status", master.Status)

	if db.PGDetail.ID != 0 {
		d.Set("parameter_group_id", db.PGDetail.ID)
	}

	return diags
}
