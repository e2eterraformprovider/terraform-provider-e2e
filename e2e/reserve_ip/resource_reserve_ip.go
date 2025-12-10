package reserve_ip

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	e2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
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
			e2econstants.AttrRegion:    config.RegionSchema(),
			e2econstants.AttrLocation:  config.LocationSchema(),
			e2econstants.AttrProjectID: config.ProjectIDSchemaResource(),

			// ============================================
			// COMPUTED FIELDS - NETWORK INFORMATION
			// ============================================
			e2econstants.AttrIPAddress: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IPv4 address of the Reserved IP",
			},
			e2econstants.AttrReservedType: {
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
			e2econstants.AttrStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current status of the Reserved IP (Available, Attached)",
			},

			// ============================================
			// COMPUTED FIELDS - IDENTIFIER
			// ============================================
			e2econstants.AttrReserveID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Numeric ID of the Reserved IP from E2E API",
			},
			e2econstants.AttrURN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique Resource Name (URN) in E2E format: e2e:reserve_ip:<region>:<ip_address>",
			},

			// ============================================
			// COMPUTED FIELDS - ATTACHMENT INFORMATION
			// ============================================
			e2econstants.AttrVMID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the VM to which the Reserved IP is attached (if any)",
			},
			e2econstants.AttrVMName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the VM to which the Reserved IP is attached (if any)",
			},
			e2econstants.AttrApplianceType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of appliance",
			},
			e2econstants.AttrFloatingIPNodes: {
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
			e2econstants.AttrCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the Reserved IP was created/purchased",
			},
			e2econstants.AttrProjectName: {
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
	d.Set(e2econstants.AttrProjectID, projectID)
	d.Set(e2econstants.AttrRegion, region)
	d.Set(e2econstants.AttrLocation, region) // Also set location for backwards compatibility

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
	d.Set(e2econstants.AttrIPAddress, rip.IPAddress)
	d.Set(e2econstants.AttrStatus, rip.Status)
	d.Set(e2econstants.AttrCreatedAt, rip.BoughtAt)
	d.Set(e2econstants.AttrVMID, strconv.Itoa(rip.VMID))
	d.Set(e2econstants.AttrVMName, rip.VMName)
	d.Set(e2econstants.AttrReserveID, rip.ReserveID)
	d.Set(e2econstants.AttrApplianceType, rip.ApplianceType)
	d.Set(e2econstants.AttrReservedType, rip.ReservedType) // Keep for backwards compatibility
	d.Set("type", rip.ReservedType)                        // V3 field
	d.Set(e2econstants.AttrURN, urn)                       // V3 field
	d.Set(e2econstants.AttrProjectName, rip.ProjectName)

	// Set floating_ip_attached_nodes if type is FloatingIP
	if rip.ReservedType == "FloatingIP" && len(rip.FloatingIPAttachedNodes) > 0 {
		d.Set(e2econstants.AttrFloatingIPNodes, flattenFloatingIPAttachedNodes(rip.FloatingIPAttachedNodes))
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
	d.Set(e2econstants.AttrIPAddress, data.IPAddress)
	d.Set(e2econstants.AttrStatus, data.Status)
	d.Set(e2econstants.AttrCreatedAt, data.BoughtAt)
	d.Set(e2econstants.AttrVMID, strconv.Itoa(data.VMID))
	d.Set(e2econstants.AttrVMName, data.VMName)
	d.Set(e2econstants.AttrApplianceType, data.ApplianceType)
	d.Set(e2econstants.AttrReservedType, data.ReservedType) // Keep for backwards compatibility
	d.Set("type", data.ReservedType)                        // V3 field
	d.Set(e2econstants.AttrURN, urn)                        // V3 field
	d.Set(e2econstants.AttrReserveID, data.ReserveID)
	d.Set(e2econstants.AttrProjectName, data.ProjectName)

	// Set floating_ip_attached_nodes if type is FloatingIP
	if data.ReservedType == "FloatingIP" && len(data.FloatingIPAttachedNodes) > 0 {
		d.Set(e2econstants.AttrFloatingIPNodes, flattenFloatingIPAttachedNodes(data.FloatingIPAttachedNodes))
	} else {
		d.Set(e2econstants.AttrFloatingIPNodes, []map[string]interface{}{})
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
