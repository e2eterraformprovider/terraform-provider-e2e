package volume_attachment

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceVolumeAttachment() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// ============================================
			// COMMON FIELDS
			// ============================================
			tfconstants.AttrRegion:    config.RegionSchema(),
			tfconstants.AttrLocation:  config.LocationSchema(),
			tfconstants.AttrProjectID: config.ProjectIDSchemaResource(),

			// ============================================
			// REQUIRED FIELDS (ForceNew - attachment is immutable)
			// ============================================
			tfconstants.AttrNodeID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "id of the Node to which the volume will be attached",
			},

			tfconstants.AttrVolumeID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "id of the Block Storage volume to attach",
			},

			// ============================================
			// COMPUTED FIELDS
			// ============================================
			tfconstants.AttrVMID: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "VM id of the Node",
			},

			tfconstants.AttrDeviceName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "device name of the attached volume",
			},

			tfconstants.AttrStatus: {
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
			Create: schema.DefaultTimeout(tfconstants.StateChangeTimeoutShort),
			Delete: schema.DefaultTimeout(tfconstants.StateChangeTimeoutShort),
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
		return diag.Errorf(tfconstants.ErrorCreatingGoe2eClient, err)
	}

	nodeID := d.Get(tfconstants.AttrNodeID).(string)
	volumeID := d.Get(tfconstants.AttrVolumeID).(string)

	log.Printf(LogAttachTemplate, volumeID, nodeID, projectID, region)

	// Get node details to retrieve VM ID and validate
	node, _, err := goe2eClient.Nodes.GetNode(ctx, nodeID)
	if err != nil {
		return diag.Errorf(tfconstants.ResourceOperationByIDErrorTemplate, tfconstants.OperationRetrieving, "Node", nodeID, projectID, region, err)
	}

	if node == nil {
		return diag.Errorf(ErrorNodeNotFoundTemplate, nodeID)
	}

	vmID := node.VMID

	// Check if node plan supports block storage
	if len(node.Plan) >= 2 && node.Plan[0:2] == tfconstants.PREFIX_C2_NODE {
		return diag.Errorf(ErrorC2PlanNoBlockStorageAttachmentFormat, nodeID)
	}

	// Prepare goe2e request
	attachReq := &goe2e.VolumeAttachmentRequest{
		NodeID:   nodeID,
		VolumeID: volumeID,
	}

	attachmentID := fmt.Sprintf("%s%s%s", nodeID, tfconstants.VolumeAttachmentImportDelimiter, volumeID)

	// Attach the volume using goe2e client
	_, _, err = goe2eClient.VolumeAttachment.AttachVolume(ctx, attachReq)
	if err != nil {
		return diag.Errorf(tfconstants.ResourceOperationErrorTemplate, goe2econstants.BlockStorageActionAttach, ResourceName, attachmentID, projectID, region, err)
	}

	// Set the resource ID as a composite of node_id and volume_id
	d.SetId(attachmentID)
	d.Set(tfconstants.AttrVMID, vmID)
	d.Set(tfconstants.AttrProjectID, projectID)
	d.Set(tfconstants.AttrRegion, region)

	// Wait for attachment to complete
	if err := waitForVolumeAttachment(ctx, goe2eClient, nodeID, volumeID, projectID, region, true); err != nil {
		return diag.FromErr(err)
	}

	log.Printf(LogAttachedTemplate, volumeID, nodeID)

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
		return diag.Errorf(tfconstants.ErrorCreatingGoe2eClient, err)
	}

	// Parse composite ID
	nodeID, volumeID, err := parseVolumeAttachmentID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf(LogReadTemplate, nodeID, volumeID, projectID, region)

	// Verify the volume exists and get its details
	volume, _, err := goe2eClient.BlockStorage.GetBlockStorage(ctx, volumeID)
	if err != nil {
		log.Printf(LogVolumeNotFound, volumeID)
		d.SetId("")
		return diags
	}

	// Check if volume is attached to the node
	var volumeVMID string
	var volumeVMName string
	if vmDetail, ok := volume.VMDetail[tfconstants.VolumeAttachmentVMDetailKeyVMID]; ok {
		if vmIDStr, ok := vmDetail.(string); ok {
			volumeVMID = vmIDStr
		} else if vmIDFloat, ok := vmDetail.(float64); ok {
			volumeVMID = strconv.Itoa(int(vmIDFloat))
		}
	}
	if vmDetail, ok := volume.VMDetail[tfconstants.VolumeAttachmentVMDetailKeyVMName]; ok {
		if vmNameStr, ok := vmDetail.(string); ok {
			volumeVMName = vmNameStr
		}
	}

	// Get node details to verify attachment
	node, _, err := goe2eClient.Nodes.GetNode(ctx, nodeID)
	if err != nil {
		log.Printf(LogNodeNotFound, nodeID)
		d.SetId("")
		return diags
	}
	if node == nil {
		log.Printf(LogNodeNotFound, nodeID)
		d.SetId("")
		return diags
	}

	vmID := node.VMID

	// Verify the volume is attached to this specific node
	if volumeVMID == "" || volumeVMID != strconv.Itoa(vmID) {
		log.Printf(LogNotAttached, volumeID, nodeID)
		d.SetId("")
		return diags
	}

	// Update state
	d.Set(tfconstants.AttrNodeID, nodeID)
	d.Set(tfconstants.AttrVolumeID, volumeID)
	d.Set(tfconstants.AttrVMID, vmID)
	d.Set(tfconstants.AttrStatus, volume.Status)
	d.Set(tfconstants.AttrProjectID, projectID)
	d.Set(tfconstants.AttrRegion, region)

	// Device name may not be available in API, but we can infer from vm_name
	if volumeVMName != "" {
		d.Set(tfconstants.AttrDeviceName, volumeVMName)
	}

	log.Printf(LogReadSuccess, nodeID, volumeID)

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
		return diag.Errorf(tfconstants.ErrorCreatingGoe2eClient, err)
	}

	// Parse composite ID
	nodeID, volumeID, err := parseVolumeAttachmentID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf(LogDetachTemplate, volumeID, nodeID, projectID, region)

	// Get node details to verify it still exists
	_, _, err = goe2eClient.Nodes.GetNode(ctx, nodeID)
	if err != nil {
		// If node doesn't exist, the attachment is already gone
		log.Printf(LogNodeMissingDetach, nodeID)
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
		return diag.Errorf(tfconstants.ResourceOperationErrorTemplate, goe2econstants.BlockStorageActionDetach, ResourceName, d.Id(), projectID, region, err)
	}

	// Wait for detachment to complete
	if err := waitForVolumeAttachment(ctx, goe2eClient, nodeID, volumeID, projectID, region, false); err != nil {
		return diag.FromErr(err)
	}

	log.Printf(LogDetachedTemplate, volumeID, nodeID)

	d.SetId("")
	return diags
}

func resourceVolumeAttachmentImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	// Import ID format: node_id/volume_id or project_id/region/node_id/volume_id
	parts := strings.Split(d.Id(), tfconstants.VolumeAttachmentImportDelimiter)

	var nodeID, volumeID, projectID, region string

	switch len(parts) {
	case ImportIDPartsShortCount:
		// Simple format: node_id/volume_id
		nodeID = parts[0]
		volumeID = parts[1]
	case ImportIDPartsFullCount:
		// Full format: project_id/region/node_id/volume_id
		projectID = parts[0]
		region = parts[1]
		nodeID = parts[2]
		volumeID = parts[3]

		// Set project_id and region in state
		d.Set(tfconstants.AttrProjectID, projectID)
		d.Set(tfconstants.AttrRegion, region)
	default:
		return nil, fmt.Errorf(tfconstants.ImportIDInvalidFormatTemplate, d.Id(), ImportIDFormatShortDescription+" or "+ImportIDFormatFullDescription)
	}

	// Set the composite ID
	d.SetId(fmt.Sprintf("%s%s%s", nodeID, tfconstants.VolumeAttachmentImportDelimiter, volumeID))
	d.Set(tfconstants.AttrNodeID, nodeID)
	d.Set(tfconstants.AttrVolumeID, volumeID)

	// Call Read to populate the rest of the state
	diags := resourceVolumeAttachmentRead(ctx, d, m)
	if diags.HasError() {
		return nil, fmt.Errorf(ErrorImportReadDuringImportTemplate, diags)
	}

	return []*schema.ResourceData{d}, nil
}

// Helper function to parse the composite ID
func parseVolumeAttachmentID(id string) (nodeID, volumeID string, err error) {
	parts := strings.Split(id, tfconstants.VolumeAttachmentImportDelimiter)
	if len(parts) != ImportIDPartsShortCount {
		return "", "", fmt.Errorf(ErrorParseIDTemplate, id, ImportIDFormatShortDescription)
	}
	return parts[0], parts[1], nil
}

// Helper function to wait for volume attachment/detachment
func waitForVolumeAttachment(ctx context.Context, goe2eClient *goe2e.Client, nodeID, volumeID, projectID, region string, shouldBeAttached bool) error {

	ticker := time.NewTicker(tfconstants.VolumeAttachmentPollInterval)
	defer ticker.Stop()

	timeout := time.After(tfconstants.VolumeAttachmentWaitTimeout)

	action := goe2econstants.BlockStorageActionAttach
	if !shouldBeAttached {
		action = goe2econstants.BlockStorageActionDetach
	}

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf(ErrorWaitContextCancelledTemplate, action)
		case <-timeout:
			return fmt.Errorf(ErrorWaitTimeoutTemplate, volumeID, action, nodeID)
		case <-ticker.C:
			volume, _, err := goe2eClient.BlockStorage.GetBlockStorage(ctx, volumeID)
			if err != nil {
				log.Printf(LogDebugVolumeCheck, err)
				continue
			}

			var volumeVMID string
			if vmDetail, ok := volume.VMDetail[tfconstants.VolumeAttachmentVMDetailKeyVMID]; ok {
				if vmIDStr, ok := vmDetail.(string); ok {
					volumeVMID = vmIDStr
				} else if vmIDFloat, ok := vmDetail.(float64); ok {
					volumeVMID = strconv.Itoa(int(vmIDFloat))
				}
			}

			if shouldBeAttached {
				// Waiting for attachment - check if vm_id is set
				if volumeVMID != "" && volumeVMID != tfconstants.VolumeAttachmentVMIDNullValue {
					// Verify it's attached to the correct node
					node, _, err := goe2eClient.Nodes.GetNode(ctx, nodeID)
					if err != nil {
						log.Printf(LogDebugNodeCheck, err)
						continue
					}

					if node != nil && volumeVMID == strconv.Itoa(node.VMID) {
						log.Printf(LogWaitAttached, volumeID, nodeID)
						return nil
					}
				}
			} else {
				// Waiting for detachment - check if vm_id is cleared
				if volumeVMID == "" || volumeVMID == tfconstants.VolumeAttachmentVMIDNullValue {
					log.Printf(LogWaitDetached, volumeID, nodeID)
					return nil
				}
			}
		}
	}
}
