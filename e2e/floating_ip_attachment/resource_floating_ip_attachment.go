package floating_ip_attachment

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	e2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceFloatingIPAttachment() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// ============================================
			// COMMON FIELDS (Input Only)
			// ============================================
			e2econstants.AttrRegion:    config.RegionSchema(),
			e2econstants.AttrLocation:  config.LocationSchema(),
			e2econstants.AttrProjectID: config.ProjectIDSchemaResource(),

			// ============================================
			// REQUIRED FIELDS
			// ============================================
			e2econstants.AttrIPAddress: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The floating IP address to attach",
			},
			"node_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of node IDs to attach the floating IP to",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		CreateContext: resourceFloatingIPAttachmentCreate,
		ReadContext:   resourceFloatingIPAttachmentRead,
		UpdateContext: resourceFloatingIPAttachmentUpdate,
		DeleteContext: resourceFloatingIPAttachmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceFloatingIPAttachmentImport,
		},
	}
}

func resourceFloatingIPAttachmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	var diags diag.Diagnostics

	log.Printf("[INFO] FLOATING IP ATTACHMENT CREATE STARTS")

	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	projectID, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	ipAddress := d.Get(e2econstants.AttrIPAddress).(string)
	nodeIDsInterface := d.Get("node_ids").([]interface{})
	nodeIDs := make([]string, len(nodeIDsInterface))
	for i, v := range nodeIDsInterface {
		nodeIDs[i] = v.(string)
	}

	if len(nodeIDs) == 0 {
		return diag.Errorf("node_ids cannot be empty")
	}

	// Use goe2e client for attachment
	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return diag.Errorf("Error creating goe2e client: %s", err)
	}

	attachReq := &goe2e.FloatingIPAttachmentRequest{
		IPAddress: ipAddress,
		NodeIDs:   nodeIDs,
	}

	_, err = goe2eClient.ReserveIP.AttachFloatingIP(ctx, attachReq)
	if err != nil {
		return diag.Errorf("Error attaching floating IP (%s) to nodes in project (%s), region (%s): %s", ipAddress, projectID, region, err)
	}

	log.Printf("[INFO] Floating IP (%s) attached to nodes: %v", ipAddress, nodeIDs)

	// Set resource ID as ip_address (since one IP can be attached to multiple nodes)
	d.SetId(ipAddress)
	d.Set(e2econstants.AttrIPAddress, ipAddress)
	d.Set("node_ids", nodeIDs)
	d.Set(e2econstants.AttrProjectID, projectID)
	d.Set(e2econstants.AttrRegion, region)

	return diags
}

func resourceFloatingIPAttachmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	ipAddress := d.Id()

	// Use goe2e client for retrieval
	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return diag.Errorf("Error creating goe2e client: %s", err)
	}

	rips, _, err := goe2eClient.ReserveIP.ListReserveIPs(ctx)
	if err != nil {
		return diag.Errorf("Error retrieving reserved IPs in project (%s), region (%s): %s", projectID, region, err)
	}

	// Find the reserved IP by IP address
	var reserveIP *goe2e.ReserveIP
	for i := range rips {
		if rips[i].IPAddress == ipAddress {
			reserveIP = &rips[i]
			break
		}
	}

	if reserveIP == nil {
		log.Printf("[WARN] Floating IP (address: %s) not found in project (%s), region (%s)", ipAddress, projectID, region)
		d.SetId("")
		return diags
	}

	// Verify it's a FloatingIP type
	if reserveIP.ReservedType != "FloatingIP" {
		log.Printf("[WARN] Reserved IP (%s) is not a FloatingIP type, removing from state", ipAddress)
		d.SetId("")
		return diags
	}

	// Extract node IDs from attached nodes
	nodeIDs := make([]string, 0)
	if len(reserveIP.FloatingIPAttachedNodes) > 0 {
		for _, node := range reserveIP.FloatingIPAttachedNodes {
			nodeIDs = append(nodeIDs, fmt.Sprintf("%d", node.ID))
		}
	}

	// If no nodes are attached, the resource should be removed from state
	if len(nodeIDs) == 0 {
		log.Printf("[WARN] Floating IP (%s) has no attached nodes, removing from state", ipAddress)
		d.SetId("")
		return diags
	}

	// Update state
	d.Set(e2econstants.AttrIPAddress, ipAddress)
	d.Set("node_ids", nodeIDs)
	d.Set(e2econstants.AttrProjectID, projectID)
	d.Set(e2econstants.AttrRegion, region)

	return diags
}

func resourceFloatingIPAttachmentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)

	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	projectID, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	ipAddress := d.Id()

	// Use goe2e client
	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return diag.Errorf("Error creating goe2e client: %s", err)
	}

	// Get old and new node IDs
	oldNodeIDsInterface, newNodeIDsInterface := d.GetChange("node_ids")
	oldNodeIDs := make([]string, 0)
	newNodeIDs := make([]string, 0)

	for _, v := range oldNodeIDsInterface.([]interface{}) {
		oldNodeIDs = append(oldNodeIDs, v.(string))
	}
	for _, v := range newNodeIDsInterface.([]interface{}) {
		newNodeIDs = append(newNodeIDs, v.(string))
	}

	// Find nodes to detach (in old but not in new)
	nodesToDetach := make([]string, 0)
	oldNodeSet := make(map[string]bool)
	for _, id := range oldNodeIDs {
		oldNodeSet[id] = true
	}
	for _, id := range newNodeIDs {
		if !oldNodeSet[id] {
			// This is a new node, will be handled by attach
		} else {
			delete(oldNodeSet, id)
		}
	}
	for id := range oldNodeSet {
		nodesToDetach = append(nodesToDetach, id)
	}

	// Find nodes to attach (in new but not in old)
	nodesToAttach := make([]string, 0)
	newNodeSet := make(map[string]bool)
	for _, id := range newNodeIDs {
		newNodeSet[id] = true
	}
	for _, id := range oldNodeIDs {
		delete(newNodeSet, id)
	}
	for id := range newNodeSet {
		nodesToAttach = append(nodesToAttach, id)
	}

	// Validate that newNodeIDs is not empty after processing changes
	if len(newNodeIDs) == 0 {
		return diag.Errorf("node_ids cannot be empty. A floating IP attachment must have at least one node attached")
	}

	// Detach nodes that are no longer in the list
	if len(nodesToDetach) > 0 {
		detachReq := &goe2e.FloatingIPDetachmentRequest{
			IPAddress: ipAddress,
			NodeIDs:   nodesToDetach,
		}
		_, err = goe2eClient.ReserveIP.DetachFloatingIP(ctx, detachReq)
		if err != nil {
			return diag.Errorf("Error detaching floating IP (%s) from nodes in project (%s), region (%s): %s", ipAddress, projectID, region, err)
		}
		log.Printf("[INFO] Detached floating IP (%s) from nodes: %v", ipAddress, nodesToDetach)
	}

	// Attach new nodes
	if len(nodesToAttach) > 0 {
		attachReq := &goe2e.FloatingIPAttachmentRequest{
			IPAddress: ipAddress,
			NodeIDs:   nodesToAttach,
		}
		_, err = goe2eClient.ReserveIP.AttachFloatingIP(ctx, attachReq)
		if err != nil {
			return diag.Errorf("Error attaching floating IP (%s) to nodes in project (%s), region (%s): %s", ipAddress, projectID, region, err)
		}
		log.Printf("[INFO] Attached floating IP (%s) to nodes: %v", ipAddress, nodesToAttach)
	}

	// Update state
	d.Set("node_ids", newNodeIDs)

	return resourceFloatingIPAttachmentRead(ctx, d, m)
}

func resourceFloatingIPAttachmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	ipAddress := d.Id()
	nodeIDsInterface := d.Get("node_ids").([]interface{})
	nodeIDs := make([]string, len(nodeIDsInterface))
	for i, v := range nodeIDsInterface {
		nodeIDs[i] = v.(string)
	}

	// Use goe2e client for detachment
	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return diag.Errorf("Error creating goe2e client: %s", err)
	}

	detachReq := &goe2e.FloatingIPDetachmentRequest{
		IPAddress: ipAddress,
		NodeIDs:   nodeIDs,
	}

	_, err = goe2eClient.ReserveIP.DetachFloatingIP(ctx, detachReq)
	if err != nil {
		return diag.Errorf("Error detaching floating IP (%s) from nodes in project (%s), region (%s): %s", ipAddress, projectID, region, err)
	}

	log.Printf("[INFO] Floating IP (%s) detached from nodes: %v", ipAddress, nodeIDs)

	d.SetId("")
	return diags
}

func resourceFloatingIPAttachmentImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	// Import ID format: project_id/region/ip_address
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid import ID format, expected: project_id/region/ip_address")
	}

	projectID := parts[0]
	region := parts[1]
	ipAddress := parts[2]

	_ = m.(*config.Config)
	d.Set(e2econstants.AttrProjectID, projectID)
	d.Set(e2econstants.AttrRegion, region)
	d.SetId(ipAddress)

	// Call Read to populate the rest of the state
	diags := resourceFloatingIPAttachmentRead(ctx, d, m)
	if diags.HasError() {
		return nil, fmt.Errorf("error reading floating IP attachment during import: %v", diags)
	}

	return []*schema.ResourceData{d}, nil
}
