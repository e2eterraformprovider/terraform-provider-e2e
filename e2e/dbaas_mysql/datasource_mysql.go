package dbaas_mysql

import (
	"context"
	"fmt"
	"strconv"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceMySQLDBaaS() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the MySQL DBaaS instance",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project ID",
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Deployment location",
			},
			"database_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Database ID",
			},
			"database_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the database",
			},
			"database_user": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Database user",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the DBaaS",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP address",
			},
			"private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Private IP address",
			},
			"is_public_ip_attached": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether a public IP is attached",
			},
			"disk": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Disk size",
			},
			"plan": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Plan name",
			},
			"database_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "MySQL version",
			},
			"parameter_group_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Attached Parameter Group ID",
			},
			"power_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Power status of the DBaaS",
			},
			"is_encryption_enabled": {
				Type:		schema.TypeBool,
				Computed:	true,
				Description: "Whether encryption is enabled for the database",
			},
		},
		ReadContext: dataSourceReadMySQL,
	}
}

func dataSourceReadMySQL(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics

	dbaasID := d.Get("id").(string)
	projectID := d.Get("project_id").(string)
	location := d.Get("location").(string)

	res, err := apiClient.GetMySqlDbaas(dbaasID, projectID, location)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while fteching dbaas instance details: %s", err))
	}

	mysql := res.Data
	master := mysql.MasterNode
	db := master.Database
	plan := master.Plan
	software := plan.Software

	d.SetId(strconv.Itoa(mysql.ID))
	d.Set("database_id", db.ID)
	d.Set("database_name", db.Database)
	d.Set("database_user", db.Username)
	d.Set("status", mysql.Status)
	d.Set("public_ip", master.PublicIPAddress)
	d.Set("private_ip", master.PrivateIPAddress)
	d.Set("is_public_ip_attached", master.PublicIPAddress != "")
	d.Set("disk", master.Disk)
	d.Set("plan", plan.Name)
	d.Set("database_version", software.Version)
	d.Set("parameter_group_id", db.PGDetail.ID)
	d.Set("power_status", master.Status)
	d.Set("is_encryption_enabled", mysql.IsEncryptionEnabled)

	return diags
}
