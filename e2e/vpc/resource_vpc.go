package vpc

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceVpc() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// ============================================
			// COMMON FIELDS
			// ============================================
			tfconstants.AttrRegion:    config.RegionSchema(),
			tfconstants.AttrLocation:  config.LocationSchema(),
			tfconstants.AttrProjectID: config.ProjectIDSchemaResource(),

			// ============================================
			// REQUIRED INPUT FIELDS (Immutable)
			// ============================================
			tfconstants.AttrName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "name of the VPC",
			},

			// ============================================
			// OPTIONAL INPUT FIELDS (Immutable)
			// ============================================
			"ipv4": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				ForceNew:    true,
				Description: "the IPv4 CIDR block for custom VPC (leave empty for E2E-managed VPC)",
			},
			"is_e2e_vpc": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				ForceNew:    true,
				Description: "whether this is an E2E-managed VPC (true: auto-allocated CIDR, false: requires ipv4 CIDR)",
			},

			// ============================================
			// COMPUTED FIELDS - IDENTIFIERS
			// ============================================
			"network_id": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "id of the VPC network",
			},

			// ============================================
			// COMPUTED FIELDS - STATUS
			// ============================================
			tfconstants.AttrCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the creation date for the VPC",
			},
			tfconstants.AttrStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "state of the VPC instance",
			},
			"is_active": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "whether the VPC is active",
			},

			// ============================================
			// COMPUTED FIELDS - NETWORK
			// ============================================
			"ipv4_cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the IPv4 CIDR block of the VPC",
			},
			"gateway_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the gateway IP address of the VPC",
			},
			"pool_size": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "the IP pool size of the VPC",
			},
		},
		ReadContext:   ResourceReadVpc,
		CreateContext: ResourceCreateVpc,
		DeleteContext: ResourceDeleteVpc,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				// Support both "vpc_id" and "project_id:vpc_id" formats
				parts := strings.Split(d.Id(), ":")

				if len(parts) == 2 {
					// Format: project_id:vpc_id
					d.Set("project_id", parts[0])
					d.SetId(parts[1])
				} else if len(parts) != 1 {
					return nil, fmt.Errorf("invalid import format: expected 'vpc_id' or 'project_id:vpc_id', got '%s'", d.Id())
				}
				// For single vpc_id, provider default project_id will be used

				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func ResourceReadVpc(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	cfg := m.(*config.Config)
	goe2eClient := cfg.Goe2eClient()
	vpcId := d.Id()

	// Get project_id with provider default support
	project_id, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get region with provider default support
	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Inside vpcs resource | read ")
	vpc, _, err := goe2eClient.Vpcs.GetVPC(ctx, vpcId)
	if err != nil {
		return diag.Errorf("Error retrieving VPC (ID: %s) in project (%s), region (%s): %s", vpcId, project_id, region, err)
	}

	d.Set(tfconstants.AttrName, vpc.Name)
	d.Set("network_id", vpc.ID)
	d.Set(tfconstants.AttrCreatedAt, vpc.CreatedAt)
	d.Set(tfconstants.AttrStatus, vpc.State)
	d.Set("ipv4_cidr", vpc.IPv4CIDR)
	d.Set("gateway_ip", vpc.GatewayIP)
	d.Set("is_active", vpc.IsActive)
	d.Set("pool_size", vpc.PoolSize)

	return diags
}
func ResourceCreateVpc(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	goe2eClient := cfg.Goe2eClient()

	// Get project_id with provider default support
	project_id, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get region with provider default support
	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Inside vpcs resource | create ")

	createReq := &goe2e.VpcCreateRequest{
		VpcName:  d.Get(tfconstants.AttrName).(string),
		IPv4:     d.Get("ipv4").(string),
		IsE2EVpc: d.Get("is_e2e_vpc").(bool),
	}
	vpc, _, err := goe2eClient.Vpcs.CreateVPC(ctx, createReq)
	if err != nil {
		return diag.Errorf("Error creating VPC (name: %s) in project (%s), region (%s): %s", createReq.VpcName, project_id, region, err)
	}

	log.Printf("[INFO] vpc creation | before setting fields")
	vpcID := int(vpc.ID)
	log.Printf("[INFO] vpc creation | network_id: %d", vpcID)

	d.SetId(strconv.Itoa(vpcID))

	return ResourceReadVpc(ctx, d, m)
}

func ResourceDeleteVpc(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	cfg := m.(*config.Config)
	goe2eClient := cfg.Goe2eClient()
	var diags diag.Diagnostics
	vpcId := d.Id()

	// Get project_id with provider default support
	project_id, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get region with provider default support
	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = goe2eClient.Vpcs.DeleteVPC(ctx, vpcId)
	if err != nil {
		return diag.Errorf("Error deleting VPC (ID: %s) in project (%s), region (%s): %s", vpcId, project_id, region, err)
	}

	d.SetId("")
	return diags

}
