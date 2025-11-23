package dbaas_postgress

import (
	"context"
	"strconv"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourcePostgresDBaaS() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the Postgres DBaaS instance",
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
			"pg_details": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Full pg_detail map",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the DBaaS",
			},
			"status_actions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
				Description: "Postgres version",
			},
			"parameter_group_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Attached Parameter Group ID",
			},
			"disk_size": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Size of the attached disk",
			},
			"power_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Power status of the DBaaS",
			},
		},
		ReadContext: dataSourceReadPostgres,
	}
}

func dataSourceReadPostgres(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	apiClient := cfg.Client()
	var diags diag.Diagnostics

	dbaasID := d.Get("id").(string)
	projectID := d.Get("project_id").(string)
	location := d.Get("location").(string)

	res, err := apiClient.GetPostgressDB(dbaasID, projectID, location)
	if err != nil {
		return diag.FromErr(err)
	}

	data, ok := res["data"].(map[string]interface{})
	if !ok {
		return diag.Errorf("Invalid 'data' field in response")
	}

	d.SetId(dbaasID)

	master := data["master_node"].(map[string]interface{})
	db := master["database"].(map[string]interface{})
	plan := master["plan"].(map[string]interface{})
	software := plan["software"].(map[string]interface{})

	d.Set("database_id", int(db["id"].(float64)))
	d.Set("database_name", db["database"].(string))
	d.Set("database_user", db["username"].(string))
	d.Set("status", data["status"].(string))

	if actions, ok := data["status_actions"].([]interface{}); ok {
		var actionStrings []string
		for _, a := range actions {
			actionStrings = append(actionStrings, a.(string))
		}
		d.Set("status_actions", actionStrings)
	}

	publicIP := master["public_ip_address"].(string)
	privateIP := master["private_ip_address"].(string)

	d.Set("public_ip", publicIP)
	d.Set("private_ip", privateIP)
	d.Set("is_public_ip_attached", publicIP != "")

	d.Set("disk", master["disk"].(string))
	d.Set("disk_size", master["disk"].(string))
	d.Set("plan", plan["name"].(string))
	d.Set("database_version", software["version"].(string))

	if pgDetail, ok := db["pg_detail"].(map[string]interface{}); ok {
		if paramID, exists := pgDetail["parameter_group_id"]; exists {
			d.Set("parameter_group_id", int(paramID.(float64)))
		}

		pgStringMap := make(map[string]string)
		for k, v := range pgDetail {
			pgStringMap[k] = stringify(v)
		}
		d.Set("pg_details", pgStringMap)
	}

	if powerStatus, ok := master["status"].(string); ok {
		d.Set("power_status", powerStatus)
	}

	return diags
}

func stringify(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		return strconv.Itoa(v)
	case bool:
		return strconv.FormatBool(v)
	default:
		return ""
	}
}
