package dbaas_mariadb

import (
	"context"
	"strconv"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
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
			"is_encryption_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether encryption is enabled for the database",
			},
		},
	}
}

func dataSourceReadMariaDB(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
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
	_ = d.Set("name", maria.Name)
	_ = d.Set("database_id", db.ID)
	_ = d.Set("database_name", db.Database)
	_ = d.Set("database_user", db.Username)
	_ = d.Set("status", maria.Status)
	_ = d.Set("public_ip", master.PublicIPAddress)
	_ = d.Set("private_ip", master.PrivateIPAddress)
	_ = d.Set("is_public_ip_attached", master.PublicIPAddress != "")
	_ = d.Set("disk", master.Disk)
	_ = d.Set("plan", plan.Name)
	_ = d.Set("software_version", software.Version)
	_ = d.Set("power_status", master.Status)
	_ = d.Set("is_encryption_enabled", maria.IsEncryptionEnabled)

	if db.PGDetail.ID != 0 {
		_ = d.Set("parameter_group_id", db.PGDetail.ID)
	}

	return diags
}
