package postgres

import (
	// "context"

	//"time"
	// "github.com/e2eterraformprovider/terraform-provider-e2e/client"
	// "github.com/e2eterraformprovider/terraform-provider-e2e/constants"

	// "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/security_group"
	// "github.com/e2eterraformprovider/terraform-provider-e2e/models"

	// "github.com/hashicorp/terraform-plugin-log"
	// "github.com/hashicorp/terraform-plugin-log/tflog"

	// "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"context"
	"log"
	"strconv"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	// "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourcePostgresDBaaS() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "location where DBAAS is deployed",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project ID to which the DBaaS instance belongs",
			},
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status_title": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status_actions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"num_instances": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"software": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name":    {Type: schema.TypeString, Computed: true},
						"version": {Type: schema.TypeString, Computed: true},
						"engine":  {Type: schema.TypeString, Computed: true},
					},
				},
			},
			"master_node": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_name":          {Type: schema.TypeString, Computed: true},
						"instance_id":        {Type: schema.TypeInt, Computed: true},
						"cluster_id":         {Type: schema.TypeInt, Computed: true},
						"node_id":            {Type: schema.TypeInt, Computed: true},
						"vm_id":              {Type: schema.TypeInt, Computed: true},
						"port":               {Type: schema.TypeString, Computed: true},
						"public_ip_address":  {Type: schema.TypeString, Computed: true},
						"private_ip_address": {Type: schema.TypeString, Computed: true},
						"zabbix_host_id":     {Type: schema.TypeInt, Computed: true},
						"status":             {Type: schema.TypeString, Computed: true},
						"db_status":          {Type: schema.TypeString, Computed: true},
						"created_at":         {Type: schema.TypeString, Computed: true},
						"ssl":                {Type: schema.TypeBool, Computed: true},
						"domain":             {Type: schema.TypeString, Computed: true},
						"public_port":        {Type: schema.TypeString, Computed: true},
						"committed_info":     {Type: schema.TypeString, Computed: true},
						"database": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id":       {Type: schema.TypeInt, Computed: true},
									"username": {Type: schema.TypeString, Computed: true},
									"database": {Type: schema.TypeString, Computed: true},
								},
							},
						},
						"plan": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name":                       {Type: schema.TypeString, Computed: true},
									"price":                      {Type: schema.TypeString, Computed: true},
									"template_id":                {Type: schema.TypeInt, Computed: true},
									"ram":                        {Type: schema.TypeString, Computed: true},
									"cpu":                        {Type: schema.TypeString, Computed: true},
									"disk":                       {Type: schema.TypeString, Computed: true},
									"currency":                   {Type: schema.TypeString, Computed: true},
									"price_per_hour":             {Type: schema.TypeFloat, Computed: true},
									"price_per_month":            {Type: schema.TypeFloat, Computed: true},
									"available_inventory_status": {Type: schema.TypeBool, Computed: true},
									"software": {
										Type:     schema.TypeList,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name":    {Type: schema.TypeString, Computed: true},
												"version": {Type: schema.TypeString, Computed: true},
												"engine":  {Type: schema.TypeString, Computed: true},
											},
										},
									},
									"committed_sku": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"committed_sku_id":       {Type: schema.TypeInt, Computed: true},
												"committed_sku_name":     {Type: schema.TypeString, Computed: true},
												"committed_node_message": {Type: schema.TypeString, Computed: true},
												"committed_sku_price":    {Type: schema.TypeFloat, Computed: true},
												"committed_upto_date":    {Type: schema.TypeString, Computed: true},
												"committed_days":         {Type: schema.TypeInt, Computed: true},
											},
										},
									},
								},
							},
						},
						"committed_details": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"committed_sku_id":       {Type: schema.TypeInt, Computed: true},
									"committed_sku_name":     {Type: schema.TypeString, Computed: true},
									"committed_node_message": {Type: schema.TypeString, Computed: true},
									"committed_sku_price":    {Type: schema.TypeFloat, Computed: true},
									"committed_upto_date":    {Type: schema.TypeString, Computed: true},
									"committed_days":         {Type: schema.TypeInt, Computed: true},
								},
							},
						},
					},
				},
			},
			"project_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snapshot_exist": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"connectivity_detail": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vector_database_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_encryption_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},

		CreateContext: resourceCreatePostgress,
		// ReadContext:   resourceReadReserveIP,
		// DeleteContext: resourceDeleteReserveIP,
		// UpdateContext: resourceUpdateReserveIP,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func convertToString(data map[string]interface{}, key string) string {
	if data[key] != nil {
		return data[key].(string)
	}
	return ""
}
func resourceCreatePostgress(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics

	log.Printf("[INFO] DBAAS_POSTGRESS CREATE STARTS")

	res, err := apiClient.CreateDbaasPostgress(d.Get("project_id").(string), d.Get("location").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] DBAAS POSTGRESS CREATE | RESPONSE BODY | ------ > %+v", res)

	if _, codeok := res["code"]; !codeok || res["is_limit_available"] == false {
		return diag.Errorf(res["message"].(string))
	}

	// Extract top-level `data`
	dataArr, ok := res["data"].([]interface{})
	if !ok || len(dataArr) == 0 {
		return diag.Errorf("Invalid or empty data field in response")
	}
	data := dataArr[0].(map[string]interface{})

	// Set top-level fields
	id := int(data["id"].(float64))
	d.SetId(strconv.Itoa(id))
	d.Set("id", id)
	d.Set("name", data["name"])
	d.Set("status", data["status"])
	d.Set("status_title", data["status_title"])
	d.Set("status_actions", data["status_actions"])
	d.Set("num_instances", int(data["num_instances"].(float64)))
	d.Set("project_name", data["project_name"])
	d.Set("snapshot_exist", data["snapshot_exist"])
	d.Set("connectivity_detail", data["connectivity_detail"])
	d.Set("vector_database_status", data["vector_database_status"])
	d.Set("is_encryption_enabled", data["is_encryption_enabled"])

	// Set `software`
	if sw, ok := data["software"].(map[string]interface{}); ok {
		d.Set("software", []interface{}{map[string]interface{}{
			"name":    sw["name"],
			"version": sw["version"],
			"engine":  sw["engine"],
		}})
	}

	// Set `master_node`
	if masterNodeRaw, ok := data["master_node"].(map[string]interface{}); ok {
		masterNode := map[string]interface{}{}

		// Direct fields
		for _, key := range []string{
			"node_name", "instance_id", "cluster_id", "node_id", "vm_id", "port", "public_ip_address", "private_ip_address",
			"zabbix_host_id", "status", "db_status", "created_at", "ssl", "domain", "public_port", "committed_info",
		} {
			masterNode[key] = masterNodeRaw[key]
		}

		// master_node.database
		if dbArr, ok := masterNodeRaw["database"].([]interface{}); ok && len(dbArr) > 0 {
			db := dbArr[0].(map[string]interface{})
			masterNode["database"] = []interface{}{map[string]interface{}{
				"id":       db["id"],
				"username": db["username"],
				"database": db["database"],
			}}
		}

		// master_node.plan
		if planRaw, ok := masterNodeRaw["plan"].(map[string]interface{}); ok {
			plan := map[string]interface{}{}
			for _, key := range []string{
				"name", "price", "template_id", "ram", "cpu", "disk", "currency", "price_per_hour", "price_per_month", "available_inventory_status",
			} {
				plan[key] = planRaw[key]
			}

			// plan.software
			if psw, ok := planRaw["software"].(map[string]interface{}); ok {
				plan["software"] = []interface{}{map[string]interface{}{
					"name":    psw["name"],
					"version": psw["version"],
					"engine":  psw["engine"],
				}}
			}

			// plan.committed_sku
			if skuList, ok := planRaw["committed_sku"].([]interface{}); ok {
				var skus []interface{}
				for _, skuRaw := range skuList {
					sku := skuRaw.(map[string]interface{})
					skus = append(skus, map[string]interface{}{
						"committed_sku_id":       sku["committed_sku_id"],
						"committed_sku_name":     sku["committed_sku_name"],
						"committed_node_message": sku["committed_node_message"],
						"committed_sku_price":    sku["committed_sku_price"],
						"committed_upto_date":    sku["committed_upto_date"],
						"committed_days":         sku["committed_days"],
					})
				}
				plan["committed_sku"] = skus
			}

			masterNode["plan"] = []interface{}{plan}
		}

		// master_node.committed_details
		if cdetails, ok := masterNodeRaw["committed_details"].([]interface{}); ok {
			var cdList []interface{}
			for _, cdRaw := range cdetails {
				cd := cdRaw.(map[string]interface{})
				cdList = append(cdList, map[string]interface{}{
					"committed_sku_id":       cd["committed_sku_id"],
					"committed_sku_name":     cd["committed_sku_name"],
					"committed_node_message": cd["committed_node_message"],
					"committed_sku_price":    cd["committed_sku_price"],
					"committed_upto_date":    cd["committed_upto_date"],
					"committed_days":         cd["committed_days"],
				})
			}
			masterNode["committed_details"] = cdList
		}

		d.Set("master_node", []interface{}{masterNode})
	}

	return diags
}
