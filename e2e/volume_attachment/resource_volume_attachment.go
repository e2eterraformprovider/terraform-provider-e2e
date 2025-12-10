package volume_attachment

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	e2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceVolumeAttachment() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// ============================================
			// COMMON FIELDS
			// ============================================
			e2econstants.AttrRegion:    config.RegionSchema(),
			e2econstants.AttrLocation:  config.LocationSchema(),
			e2econstants.AttrProjectID: config.ProjectIDSchemaResource(),

			// ============================================
			// REQUIRED FIELDS (ForceNew - attachment is immutable)
			// ============================================
			e2econstants.AttrNodeID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "id of the Node to which the volume will be attached",
			},

			e2econstants.AttrVolumeID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "id of the Block Storage volume to attach",
			},

			// ============================================
			// COMPUTED FIELDS
			// ============================================
			e2econstants.AttrVMID: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "VM id of the Node",
			},

			"device_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "device name of the attached volume",
			},

			e2econstants.AttrStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "attachment status",
			},
		},

		CreateContext: resourceVolumeAttachmentCreate,
		ReadContext:   resourceVolumeAttachmentRead,
		DeleteContext: resourceVolumeAttachmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceVolumeAttachmentImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceVolumeAttachmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)

	// Get region with provider default support
	region, err := config.GetRegionOrLocation(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get projectID with provider default support
	projectID, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return diag.Errorf("Error creating goe2e client: %s", err)
	}

	nodeID := d.Get(e2econstants.AttrNodeID).(string)
	volumeID := d.Get(e2econstants.AttrVolumeID).(string)

	log.Printf("[INFO] Attaching volume (ID: %s) to node (ID: %s) in project (%s), region (%s)", volumeID, nodeID, projectID, region)

	// Get node details to retrieve VM ID and validate
	node, _, err := goe2eClient.Nodes.GetNode(ctx, nodeID)
	if err != nil {
		return diag.Errorf("Error fetching node (ID: %s) in project (%s), region (%s): %s", nodeID, projectID, region, err)
	}

	if node == nil {
		return diag.Errorf("Error: node (ID: %s) not found", nodeID)
	}

	vmID := node.VMID

	// Check if node plan supports block storage
	if len(node.Plan) >= 2 && node.Plan[0:2] == e2econstants.PREFIX_C2_NODE {
		return diag.Errorf("Cannot attach volume to node (ID: %s): C2 plan nodes do not support block storage attachment", nodeID)
	}

	// Prepare goe2e request
	attachReq := &goe2e.VolumeAttachmentRequest{
		NodeID:   nodeID,
		VolumeID: volumeID,
	}

	// Attach the volume using goe2e client
	_, _, err = goe2eClient.VolumeAttachment.AttachVolume(ctx, attachReq)
	if err != nil {
		return diag.Errorf("Error attaching volume (ID: %s) to node (ID: %s) in project (%s), region (%s): %s", volumeID, nodeID, projectID, region, err)
	}

	// Set the resource ID as a composite of node_id and volume_id
	d.SetId(fmt.Sprintf("%s/%s", nodeID, volumeID))
	d.Set(e2econstants.AttrVMID, vmID)
	d.Set(e2econstants.AttrProjectID, projectID)
	d.Set(e2econstants.AttrRegion, region)

	// Wait for attachment to complete
	if err := waitForVolumeAttachment(ctx, goe2eClient, nodeID, volumeID, projectID, region, true); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Successfully attached volume (ID: %s) to node (ID: %s)", volumeID, nodeID)

	return resourceVolumeAttachmentRead(ctx, d, m)
}

func resourceVolumeAttachmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	var diags diag.Diagnostics

	// Get region with provider default support
	region, err := config.GetRegionOrLocation(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get projectID with provider default support
	projectID, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return diag.Errorf("Error creating goe2e client: %s", err)
	}

	// Parse composite ID
	nodeID, volumeID, err := parseVolumeAttachmentID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Reading volume attachment: node (ID: %s), volume (ID: %s) in project (%s), region (%s)", nodeID, volumeID, projectID, region)

	// Verify the volume exists and get its details
	volume, _, err := goe2eClient.BlockStorage.GetBlockStorage(ctx, volumeID)
	if err != nil {
		log.Printf("[WARN] Volume (ID: %s) not found, removing from state", volumeID)
		d.SetId("")
		return diags
	}

	// Check if volume is attached to the node
	var volumeVMID string
	var volumeVMName string
	if vmDetail, ok := volume.VMDetail["vm_id"]; ok {
		if vmIDStr, ok := vmDetail.(string); ok {
			volumeVMID = vmIDStr
		} else if vmIDFloat, ok := vmDetail.(float64); ok {
			volumeVMID = strconv.Itoa(int(vmIDFloat))
		}
	}
	if vmDetail, ok := volume.VMDetail["vm_name"]; ok {
		if vmNameStr, ok := vmDetail.(string); ok {
			volumeVMName = vmNameStr
		}
	}

	// Get node details to verify attachment
	node, _, err := goe2eClient.Nodes.GetNode(ctx, nodeID)
	if err != nil {
		log.Printf("[WARN] Node (ID: %s) not found, removing attachment from state", nodeID)
		d.SetId("")
		return diags
	}
	if node == nil {
		log.Printf("[WARN] Node (ID: %s) not found, removing attachment from state", nodeID)
		d.SetId("")
		return diags
	}

	vmID := node.VMID

	// Verify the volume is attached to this specific node
	if volumeVMID == "" || volumeVMID != strconv.Itoa(vmID) {
		log.Printf("[WARN] Volume (ID: %s) is not attached to node (ID: %s), removing from state", volumeID, nodeID)
		d.SetId("")
		return diags
	}

	// Update state
	d.Set(e2econstants.AttrNodeID, nodeID)
	d.Set(e2econstants.AttrVolumeID, volumeID)
	d.Set(e2econstants.AttrVMID, vmID)
	d.Set(e2econstants.AttrStatus, volume.Status)
	d.Set(e2econstants.AttrProjectID, projectID)
	d.Set(e2econstants.AttrRegion, region)

	// Device name may not be available in API, but we can infer from vm_name
	if volumeVMName != "" {
		d.Set("device_name", volumeVMName)
	}

	log.Printf("[INFO] Successfully read volume attachment: node (ID: %s), volume (ID: %s)", nodeID, volumeID)

	return diags
}

func resourceVolumeAttachmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	var diags diag.Diagnostics

	// Get region with provider default support
	region, err := config.GetRegionOrLocation(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get projectID with provider default support
	projectID, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return diag.Errorf("Error creating goe2e client: %s", err)
	}

	// Parse composite ID
	nodeID, volumeID, err := parseVolumeAttachmentID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Detaching volume (ID: %s) from node (ID: %s) in project (%s), region (%s)", volumeID, nodeID, projectID, region)

	// Get node details to verify it still exists
	_, _, err = goe2eClient.Nodes.GetNode(ctx, nodeID)
	if err != nil {
		// If node doesn't exist, the attachment is already gone
		log.Printf("[WARN] Node (ID: %s) not found, considering volume detached", nodeID)
		return diags
	}

	// Prepare goe2e request
	detachReq := &goe2e.VolumeDetachmentRequest{
		NodeID:   nodeID,
		VolumeID: volumeID,
	}

	// Detach the volume using goe2e client
	_, err = goe2eClient.VolumeAttachment.DetachVolume(ctx, detachReq)
	if err != nil {
		return diag.Errorf("Error detaching volume (ID: %s) from node (ID: %s) in project (%s), region (%s): %s", volumeID, nodeID, projectID, region, err)
	}

	// Wait for detachment to complete
	if err := waitForVolumeAttachment(ctx, goe2eClient, nodeID, volumeID, projectID, region, false); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Successfully detached volume (ID: %s) from node (ID: %s)", volumeID, nodeID)

	d.SetId("")
	return diags
}

func resourceVolumeAttachmentImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	// Import ID format: node_id/volume_id or project_id/region/node_id/volume_id
	parts := strings.Split(d.Id(), "/")

	var nodeID, volumeID, projectID, region string

	switch len(parts) {
	case 2:
		// Simple format: node_id/volume_id
		nodeID = parts[0]
		volumeID = parts[1]
	case 4:
		// Full format: project_id/region/node_id/volume_id
		projectID = parts[0]
		region = parts[1]
		nodeID = parts[2]
		volumeID = parts[3]

		// Set project_id and region in state
		d.Set(e2econstants.AttrProjectID, projectID)
		d.Set(e2econstants.AttrRegion, region)
	default:
		return nil, fmt.Errorf("invalid import ID format: expected 'node_id/volume_id' or 'project_id/region/node_id/volume_id', got: %s", d.Id())
	}

	// Set the composite ID
	d.SetId(fmt.Sprintf("%s/%s", nodeID, volumeID))
	d.Set(e2econstants.AttrNodeID, nodeID)
	d.Set(e2econstants.AttrVolumeID, volumeID)

	// Call Read to populate the rest of the state
	diags := resourceVolumeAttachmentRead(ctx, d, m)
	if diags.HasError() {
		return nil, fmt.Errorf("error reading volume attachment during import: %v", diags)
	}

	return []*schema.ResourceData{d}, nil
}

// Helper function to parse the composite ID
func parseVolumeAttachmentID(id string) (nodeID, volumeID string, err error) {
	parts := strings.Split(id, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid volume attachment ID format: %s (expected: node_id/volume_id)", id)
	}
	return parts[0], parts[1], nil
}

// Helper function to wait for volume attachment/detachment
func waitForVolumeAttachment(ctx context.Context, goe2eClient *goe2e.Client, nodeID, volumeID, projectID, region string, shouldBeAttached bool) error {

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	timeout := time.After(3 * time.Minute)

	action := "attach"
	if !shouldBeAttached {
		action = "detach"
	}

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled while waiting for volume %s", action)
		case <-timeout:
			return fmt.Errorf("timeout waiting for volume (ID: %s) to %s to/from node (ID: %s)", volumeID, action, nodeID)
		case <-ticker.C:
			volume, _, err := goe2eClient.BlockStorage.GetBlockStorage(ctx, volumeID)
			if err != nil {
				log.Printf("[DEBUG] Error checking volume status: %s", err)
				continue
			}

			var volumeVMID string
			if vmDetail, ok := volume.VMDetail["vm_id"]; ok {
				if vmIDStr, ok := vmDetail.(string); ok {
					volumeVMID = vmIDStr
				} else if vmIDFloat, ok := vmDetail.(float64); ok {
					volumeVMID = strconv.Itoa(int(vmIDFloat))
				}
			}

			if shouldBeAttached {
				// Waiting for attachment - check if vm_id is set
				if volumeVMID != "" && volumeVMID != "null" {
					// Verify it's attached to the correct node
					node, _, err := goe2eClient.Nodes.GetNode(ctx, nodeID)
					if err != nil {
						log.Printf("[DEBUG] Error checking node status: %s", err)
						continue
					}

					if node != nil && volumeVMID == strconv.Itoa(node.VMID) {
						log.Printf("[INFO] Volume (ID: %s) successfully attached to node (ID: %s)", volumeID, nodeID)
						return nil
					}
				}
			} else {
				// Waiting for detachment - check if vm_id is cleared
				if volumeVMID == "" || volumeVMID == "null" {
					log.Printf("[INFO] Volume (ID: %s) successfully detached from node (ID: %s)", volumeID, nodeID)
					return nil
				}
			}
		}
	}
}
