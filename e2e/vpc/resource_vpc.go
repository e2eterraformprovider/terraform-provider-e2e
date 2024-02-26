package vpc

import (
	"context"
	// "fmt"
	"log"

	// "regexp"

	"strconv"
	//"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"

	// "github.com/hashicorp/terraform-plugin-log"
	// "github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResouceVpc() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Region should specified",
			},
			"vpc_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the project. It should be unique",
			},
			"network_size": {
				Type:     schema.TypeFloat,
				Optional: true,
				Default:  512,
			},
			"network_id": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "The id of network",
			},
			"pool_size": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv4_cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gateway_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},

		ReadContext:   ResourceReadVpc,
		CreateContext: ResourceCreateVpc,
		UpdateContext: ResourceUpdateVpc,
		DeleteContext: ResourceDeleteVpc,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func ResourceReadVpc(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	apiClient := m.(*client.Client)
	log.Printf("[INFO] Inside vpcs  resourcsource | read ")
	Response, err := apiClient.GetVpc(d.Id(), d.Get("project_id").(string), d.Get("region").(string))
	if err != nil {
		return diag.Errorf("error finding vpcs ")
	}

	data := Response.Data
	d.Set("created_at", data.Created_at)
	d.Set("state", data.State)
	d.Set("ipv4_cidr", data.Ipv4_cidr)
	d.Set("gateway_ip", data.Gateway_ip)
	d.Set("is_active", data.Is_active)
	d.Set("pool_size", data.Pool_size)

	return diags
}
func ResourceCreateVpc(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	apiClient := m.(*client.Client)
	log.Printf("[INFO] Inside vpcs  resource | create ")

	newvpc := models.VpcCreate{
		VpcName:     d.Get("vpc_name").(string),
		NetworkSize: d.Get("network_size").(float64),
	}
	resvpc, err := apiClient.CreateVpc(d.Get("region").(string), &newvpc, d.Get("project_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if _, codeok := resvpc["code"]; !codeok {
		return diag.Errorf(resvpc["message"].(string))
	}

	data := resvpc["data"].(map[string]interface{})
	log.Printf("[INFO] vpc creation | before setting fields")

	var vpcID int

	if networkID, ok := data["network_id"].(float64); ok {
		vpcID = int(networkID)
		log.Printf("[INFO] vpc creation | network_id: %d", vpcID)
	} else {
		log.Printf("[ERROR] vpc creation | unable to extract network_id from data")
	}

	d.SetId(strconv.Itoa(vpcID))

	return diags
}

func ResourceUpdateVpc(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func ResourceDeleteVpc(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	apiClient := m.(*client.Client)
	var diags diag.Diagnostics
	vpcId := d.Id()

	_, err := apiClient.DeleteVpc(vpcId, d.Get("project_id").(string), d.Get("region").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diags

}
