package reserve_ip

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceReserveIP() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// ============================================
			// COMMON FIELDS (Input Only)
			// ============================================
			tfconstants.AttrRegion:    config.RegionSchema(),
			tfconstants.AttrLocation:  config.LocationSchema(),
			tfconstants.AttrProjectID: config.ProjectIDSchemaResource(),

			// ============================================
			// COMPUTED FIELDS - NETWORK INFORMATION
			// ============================================
			tfconstants.AttrIPAddress: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IPv4 address of the Reserved IP",
			},
			tfconstants.AttrReservedType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of Reserved IP (deprecated, use 'type' instead)",
				Deprecated:  "Use 'type' instead",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of Reserved IP (FloatingIP, PublicIP, AddonIP)",
			},

			// ============================================
			// COMPUTED FIELDS - STATUS
			// ============================================
			tfconstants.AttrStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current status of the Reserved IP (Available, Attached)",
			},

			// ============================================
			// COMPUTED FIELDS - IDENTIFIER
			// ============================================
			tfconstants.AttrReserveID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Numeric ID of the Reserved IP from E2E API",
			},
			tfconstants.AttrURN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique Resource Name (URN) in E2E format: e2e:reserve_ip:<region>:<ip_address>",
			},

			// ============================================
			// COMPUTED FIELDS - ATTACHMENT INFORMATION
			// ============================================
			tfconstants.AttrVMID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the VM to which the Reserved IP is attached (if any)",
			},
			tfconstants.AttrVMName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the VM to which the Reserved IP is attached (if any)",
			},
			tfconstants.AttrApplianceType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of appliance",
			},
			tfconstants.AttrFloatingIPNodes: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of nodes attached to this floating IP (only populated when type is FloatingIP)",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Node ID",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node name",
						},
						"vm_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "VM ID",
						},
						"ip_address_public": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public IP address of the node",
						},
						"ip_address_private": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Private IP address of the node",
						},
						"status_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status name of the node",
						},
						"security_group_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Security group status of the node",
						},
					},
				},
			},

			// ============================================
			// COMPUTED FIELDS - METADATA
			// ============================================
			tfconstants.AttrCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the Reserved IP was created/purchased",
			},
			tfconstants.AttrProjectName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the project",
			},
		},
		CreateContext: resourceCreateReserveIP,
		ReadContext:   resourceReadReserveIP,
		DeleteContext: resourceDeleteReserveIP,
		// UpdateContext removed - reserved IPs are immutable
		Importer: &schema.ResourceImporter{
			StateContext: resourceReserveIPImport,
		},
	}
}

func resourceReserveIPImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	// Parse import ID: project_id/region/ip_address
	projectID, region, ipAddress, err := parseReserveIPImportID(d.Id())
	if err != nil {
		return nil, err
	}

	cfg := m.(*config.Config)
	d.Set(tfconstants.AttrProjectID, projectID)
	d.Set(tfconstants.AttrRegion, region)
	d.Set(tfconstants.AttrLocation, region) // Also set location for backwards compatibility

	// Fetch the reserved IP to verify it exists using goe2e client
	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return nil, fmt.Errorf("error creating goe2e client during import: %w", err)
	}

	rips, _, err := goe2eClient.ReserveIP.ListReserveIPs(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving reserved IPs during import: %w", err)
	}

	// Find the reserved IP by IP address
	var found bool
	for _, item := range rips {
		if item.IPAddress == ipAddress {
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("reserved IP with IP address %s not found", ipAddress)
	}

	// Use IP address as ID (DigitalOcean pattern)
	d.SetId(ipAddress)

	return []*schema.ResourceData{d}, nil
}

func resourceCreateReserveIP(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	var diags diag.Diagnostics

	log.Printf("[INFO] RESERVED IP CREATE STARTS")

	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	projectID, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Use goe2e client for creation
	goe2eClient := cfg.Goe2eClient()

	rip, _, err := goe2eClient.ReserveIP.CreateReserveIP(ctx)
	if err != nil {
		return diag.Errorf("Error creating reserved IP in project (%s), region (%s): %s", projectID, region, err)
	}

	log.Printf("[INFO] ReserveIP CREATE | RESPONSE BODY | %+v", rip)

	// Validate response
	if rip == nil || rip.IPAddress == "" {
		return diag.Errorf("Error creating reserved IP: IP address not found in response")
	}

	// Use IP address as ID (DigitalOcean pattern)
	d.SetId(rip.IPAddress)

	// Generate URN
	urn := generateReserveIPURN(region, rip.IPAddress)

	// Set all computed fields from the response
	d.Set(tfconstants.AttrIPAddress, rip.IPAddress)
	d.Set(tfconstants.AttrStatus, rip.Status)
	d.Set(tfconstants.AttrCreatedAt, rip.BoughtAt)
	d.Set(tfconstants.AttrVMID, strconv.Itoa(rip.VMID))
	d.Set(tfconstants.AttrVMName, rip.VMName)
	d.Set(tfconstants.AttrReserveID, rip.ReserveID)
	d.Set(tfconstants.AttrApplianceType, rip.ApplianceType)
	d.Set(tfconstants.AttrReservedType, rip.ReservedType) // Keep for backwards compatibility
	d.Set("type", rip.ReservedType)                       // V3 field
	d.Set(tfconstants.AttrURN, urn)                       // V3 field
	d.Set(tfconstants.AttrProjectName, rip.ProjectName)

	// Set floating_ip_attached_nodes if type is FloatingIP
	if rip.ReservedType == "FloatingIP" && len(rip.FloatingIPAttachedNodes) > 0 {
		d.Set(tfconstants.AttrFloatingIPNodes, flattenFloatingIPAttachedNodes(rip.FloatingIPAttachedNodes))
	}

	return diags
}

func resourceReadReserveIP(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	var diags diag.Diagnostics

	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	projectID, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Use ID (IP address) for lookup (DigitalOcean pattern)
	ipAddress := d.Id()

	// Use goe2e client for retrieval
	goe2eClient := cfg.Goe2eClient()

	rips, _, err := goe2eClient.ReserveIP.ListReserveIPs(ctx)
	if err != nil {
		return diag.Errorf("Error retrieving reserved IP list in project (%s), region (%s): %s", projectID, region, err)
	}

	log.Printf("[INFO] ReserveIP READ | RESPONSE | %+v", rips)

	// Find the reserved IP by IP address
	var data *goe2e.ReserveIP
	for i := range rips {
		if rips[i].IPAddress == ipAddress {
			data = &rips[i]
			break
		}
	}

	log.Printf("[INFO] FILTER DATA | %+v", data)

	if data == nil || data.IPAddress == "" {
		log.Printf("[WARN] Reserved IP (address: %s) not found in project (%s), region (%s)", ipAddress, projectID, region)
		d.SetId("")
		return diags
	}

	// Generate URN
	urn := generateReserveIPURN(region, data.IPAddress)

	// Set all computed fields
	log.Printf("[INFO] ReserveIP READ | SETTING DATA %+v", data)
	d.Set(tfconstants.AttrIPAddress, data.IPAddress)
	d.Set(tfconstants.AttrStatus, data.Status)
	d.Set(tfconstants.AttrCreatedAt, data.BoughtAt)
	d.Set(tfconstants.AttrVMID, strconv.Itoa(data.VMID))
	d.Set(tfconstants.AttrVMName, data.VMName)
	d.Set(tfconstants.AttrApplianceType, data.ApplianceType)
	d.Set(tfconstants.AttrReservedType, data.ReservedType) // Keep for backwards compatibility
	d.Set("type", data.ReservedType)                       // V3 field
	d.Set(tfconstants.AttrURN, urn)                        // V3 field
	d.Set(tfconstants.AttrReserveID, data.ReserveID)
	d.Set(tfconstants.AttrProjectName, data.ProjectName)

	// Set floating_ip_attached_nodes if type is FloatingIP
	if data.ReservedType == "FloatingIP" && len(data.FloatingIPAttachedNodes) > 0 {
		d.Set(tfconstants.AttrFloatingIPNodes, flattenFloatingIPAttachedNodes(data.FloatingIPAttachedNodes))
	} else {
		d.Set(tfconstants.AttrFloatingIPNodes, []map[string]interface{}{})
	}

	return diags
}

// resourceUpdateReserveIP removed - reserved IPs are immutable and cannot be updated

func resourceDeleteReserveIP(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	var diags diag.Diagnostics

	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	projectID, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Use ID (IP address) for deletion (DigitalOcean pattern)
	ipAddress := d.Id()

	// Use goe2e client for deletion
	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return diag.Errorf("Error creating goe2e client: %s", err)
	}

	// Check if IP is attached and log warning
	rips, _, err := goe2eClient.ReserveIP.ListReserveIPs(ctx)
	if err == nil {
		for _, rip := range rips {
			if rip.IPAddress == ipAddress {
				if rip.Status == "Attached" || len(rip.FloatingIPAttachedNodes) > 0 {
					log.Printf("[WARN] Reserved IP (%s) is currently attached. The API will handle detachment automatically.", ipAddress)
				}
				break
			}
		}
	}

	_, err = goe2eClient.ReserveIP.DeleteReserveIP(ctx, ipAddress)
	if err != nil {
		return diag.Errorf("Error deleting reserved IP (address: %s) in project (%s), region (%s): %s", ipAddress, projectID, region, err)
	}

	d.SetId("")
	return diags
}
