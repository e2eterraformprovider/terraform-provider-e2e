package blockstorage

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/node"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Constants for magic numbers
const (
	IOPS_PER_GB         = 15
	TB_TO_GB_MULTIPLIER = 1000
)

// Valid block storage sizes
var validBlockStorageSizes = []float64{250, 500, 1000, 2000, 4000, 8000, 16000, 24000}

const validBlockStorageSizesString = "250, 500, 1000, 2000, 4000, 8000, 16000, 24000"

const (
	blockStorageSizeValidationErrorTemplate = "%q must be one of [%s] GB, got: %.0f"
	blockStorageTagsStateDescription        = "map of tags to assign to the resource (state-only, API support pending)"
	blockStorageTagsStateSetWarningTemplate = "failed to set %s: %w"

	volumeVMIDKey   = "vm_id"
	volumeVMNameKey = "vm_name"

	resizeTolerance = 1e-6
)

func isBlockStorageSizeUpgrade(prevSize, currSize float64) bool {
	return currSize > prevSize+resizeTolerance
}

// validateBlockStorageSize validates that the size is one of the allowed values
func validateBlockStorageSize(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(float64)
	for _, validSize := range validBlockStorageSizes {
		if value == validSize {
			return
		}
	}
	errors = append(errors, fmt.Errorf(blockStorageSizeValidationErrorTemplate, k, validBlockStorageSizesString, value))
	return
}

func ResourceBlockStorage() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceBlockStorageResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: ResourceBlockStorageStateUpgradeV0toV1,
				Version: 0,
			},
		},

		Schema: map[string]*schema.Schema{
			// ============================================
			// COMMON INFRASTRUCTURE FIELDS
			// ============================================
			tfconstants.AttrRegion:    config.RegionSchema(),
			tfconstants.AttrLocation:  config.LocationSchema(),
			tfconstants.AttrProjectID: config.ProjectIDSchemaResource(),

			// ============================================
			// REQUIRED FIELDS - IMMUTABLE
			// ============================================
			tfconstants.AttrName: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "name of the block storage volume",
				ValidateFunc: node.ValidateName,
			},

			// ============================================
			// REQUIRED FIELDS - UPGRADABLE
			// ============================================
			tfconstants.AttrSize: {
				Type:         schema.TypeFloat,
				Required:     true,
				Description:  "size of the block storage in GB (upgradable: can increase, cannot decrease). Valid sizes: " + validBlockStorageSizesString,
				ValidateFunc: validateBlockStorageSize,
			},

			// ============================================
			// V3 OPTIONAL FIELDS
			// ============================================
			tfconstants.AttrTags: {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: blockStorageTagsStateDescription,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			// ============================================
			// COMPUTED FIELDS - PERFORMANCE
			// ============================================
			tfconstants.AttrIOPS: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IOPS (Input/Output Operations Per Second) - calculated as size * 15",
			},

			// ============================================
			// COMPUTED FIELDS - STATUS
			// ============================================
			tfconstants.AttrStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "current status of the block storage (available, attached, creating, error, etc.)",
			},

			// ============================================
			// COMPUTED FIELDS - ATTACHMENT
			// ============================================
			tfconstants.AttrVMID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the node to which the block storage is attached (null if detached)",
			},
			tfconstants.AttrVMName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "name of the node to which the block storage is attached (null if detached)",
			},
		},

		CreateContext: resourceCreateBlockStorage,
		ReadContext:   resourceReadBlockStorage,
		UpdateContext: resourceUpdateBlockStorage,
		DeleteContext: resourceDeleteBlockStorage,
		Exists:        resourceExistsBlockStorage,
		Importer: &schema.ResourceImporter{
			StateContext: customImportBlockStorage,
		},
	}
}

// resourceBlockStorageResourceV0 returns the V0 schema for state migration
func resourceBlockStorageResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			tfconstants.AttrRegion:    config.RegionSchema(),
			tfconstants.AttrLocation:  config.LocationSchema(),
			tfconstants.AttrProjectID: config.ProjectIDSchemaResource(),
			tfconstants.AttrName: {
				Type:     schema.TypeString,
				Required: true,
			},
			tfconstants.AttrSize: {
				Type:     schema.TypeFloat,
				Required: true,
			},
			tfconstants.AttrIOPS: {
				Type:     schema.TypeString,
				Computed: true,
			},
			tfconstants.AttrStatus: {
				Type:     schema.TypeString,
				Computed: true,
			},
			tfconstants.AttrVMID: {
				Type:     schema.TypeString,
				Computed: true,
			},
			tfconstants.AttrVMName: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// ResourceBlockStorageStateUpgradeV0toV1 upgrades state from v0 to v1
// Exported for testing purposes
func ResourceBlockStorageStateUpgradeV0toV1(
	ctx context.Context,
	rawState map[string]interface{},
	meta interface{},
) (map[string]interface{}, error) {
	// Add new V3 fields with defaults
	if _, exists := rawState[tfconstants.AttrTags]; !exists {
		rawState[tfconstants.AttrTags] = make(map[string]interface{})
	}

	// Preserve all existing V2 fields
	// No automatic renames or transformations

	log.Printf("[INFO] Upgraded block storage state from v0 to v1: %s", rawState[tfconstants.AttrID])
	return rawState, nil
}

func resourceCreateBlockStorage(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	var diags diag.Diagnostics

	// Get region with provider default support
	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get project_id with provider default support
	projectIDStr, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Create goe2e client with projectID and region
	goe2eClient, err := cfg.Goe2eClientForProject(projectIDStr, region)
	if err != nil {
		return diag.Errorf(tfconstants.ErrorCreatingGoe2eClient, err)
	}

	err = validateSize(ctx, d, goe2eClient)
	if err != nil {
		return diag.Errorf(ErrorValidateSizeTemplate, projectIDStr, region, err)
	}

	log.Printf("[INFO] BLOCK STORAGE CREATE STARTS")
	name := d.Get(tfconstants.AttrName).(string)
	size := d.Get(tfconstants.AttrSize).(float64)

	iops := CalculateIOPS(size)
	createReq := &goe2e.BlockStorageCreateRequest{
		Name: name,
		Size: size,
		IOPS: iops,
	}

	blockStorage, _, err := goe2eClient.BlockStorage.CreateBlockStorage(ctx, createReq)
	if err != nil {
		return diag.Errorf(tfconstants.ResourceOperationErrorTemplate, tfconstants.OperationCreating, ResourceName, name, projectIDStr, region, err)
	}

	if blockStorage == nil {
		return diag.Errorf(ErrorCreateNilResponseTemplate, name)
	}

	log.Printf("[INFO] Block Storage creation | before setting fields")
	d.SetId(strconv.Itoa(blockStorage.BlockID))
	d.Set(tfconstants.AttrIOPS, iops)

	// Initialize empty tags map for state-only tag support
	if _, ok := d.GetOk(tfconstants.AttrTags); !ok {
		if err := d.Set(tfconstants.AttrTags, make(map[string]string)); err != nil {
			return diag.FromErr(fmt.Errorf(blockStorageTagsStateSetWarningTemplate, tfconstants.AttrTags, err))
		}
	}

	return diags
}

func resourceReadBlockStorage(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	var diags diag.Diagnostics

	log.Printf("[INFO] BLOCK STORAGE READ STARTS")
	blockStorageID := d.Id()

	// Get region with provider default support
	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get project_id with provider default support
	projectIDStr, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Create goe2e client with projectID and region
	goe2eClient, err := cfg.Goe2eClientForProject(projectIDStr, region)
	if err != nil {
		return diag.Errorf(tfconstants.ErrorCreatingGoe2eClient, err)
	}

	blockStorage, resp, err := goe2eClient.BlockStorage.GetBlockStorage(ctx, blockStorageID)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			log.Printf(LogNotFoundRemoveFromStateTemplate, blockStorageID)
			d.SetId("")
			return diags
		}
		if strings.Contains(err.Error(), goe2econstants.NotFoundSubstring) {
			d.SetId("")
			return diags
		}
		return diag.Errorf(tfconstants.ResourceOperationByIDErrorTemplate, tfconstants.OperationRetrieving, ResourceName, blockStorageID, projectIDStr, region, err)
	}

	if blockStorage == nil {
		log.Printf(LogNotFoundRemoveFromStateTemplate, blockStorageID)
		d.SetId("")
		return diags
	}

	log.Printf("[INFO] BLOCK STORAGE READ | BEFORE SETTING DATA")
	d.Set(tfconstants.AttrName, blockStorage.Name)
	d.Set(tfconstants.AttrStatus, blockStorage.Status)
	d.Set(tfconstants.AttrIOPS, blockStorage.Template.TotalIOPSSec)

	// Handle VM attachment details
	if blockStorage.VMDetail != nil {
		if vmID, ok := blockStorage.VMDetail[volumeVMIDKey]; ok {
			vmIDFloat, ok := vmID.(float64)
			if ok {
				d.Set(tfconstants.AttrVMID, strconv.Itoa(int(vmIDFloat)))
			}
		}
		if vmName, ok := blockStorage.VMDetail[volumeVMNameKey]; ok {
			if vmNameStr, ok := vmName.(string); ok {
				d.Set(tfconstants.AttrVMName, vmNameStr)
			}
		}
	} else {
		d.Set(tfconstants.AttrVMID, nil)
		d.Set(tfconstants.AttrVMName, nil)
	}

	// Tags are state-only, preserve existing tags
	if tags, ok := d.GetOk(tfconstants.AttrTags); ok {
		if err := d.Set(tfconstants.AttrTags, tags); err != nil {
			return diag.FromErr(err)
		}
	}

	log.Printf("[INFO] BLOCK STORAGE READ | AFTER SETTING DATA")

	return diags
}

func resourceUpdateBlockStorage(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	var diags diag.Diagnostics
	blockStorageID := d.Id()

	// Get region with provider default support
	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get project_id with provider default support
	projectIDStr, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Create goe2e client with projectID and region
	goe2eClient, err := cfg.Goe2eClientForProject(projectIDStr, region)
	if err != nil {
		return diag.Errorf(tfconstants.ErrorCreatingGoe2eClient, err)
	}

	status := d.Get(tfconstants.AttrStatus).(string)

	if status == goe2econstants.BlockStorageStatusError {
		return diag.Errorf(ErrorUpdateInErrorStateTemplate, blockStorageID, projectIDStr, region)
	}

	blockStorage, _, err := goe2eClient.BlockStorage.GetBlockStorage(ctx, blockStorageID)
	if err != nil {
		if strings.Contains(err.Error(), goe2econstants.NotFoundSubstring) {
			d.SetId("")
			return diags
		}
		return diag.Errorf(tfconstants.ResourceOperationByIDErrorTemplate, tfconstants.OperationRetrieving, ResourceName, blockStorageID, projectIDStr, region, err)
	}

	if blockStorage == nil {
		d.SetId("")
		return diags
	}

	// Region and project_id changes are handled by ForceNew: true in schema
	// No need for manual checks here

	// Handle size upgrades
	if d.HasChange(tfconstants.AttrSize) {
		prevSize, currSize := d.GetChange(tfconstants.AttrSize)
		err := validateSize(ctx, d, goe2eClient)
		if err != nil {
			d.Set(tfconstants.AttrSize, prevSize)
			return diag.Errorf("Error validating "+ResourceName+" (ID: %s) size in project (%s), region (%s): %s", blockStorageID, projectIDStr, region, err)
		}
		log.Printf("[INFO] prevSize %v, currSize %v", prevSize, currSize)

		if d.Get(tfconstants.AttrStatus) == goe2econstants.BlockStorageStatusAttached {
			if isBlockStorageSizeUpgrade(prevSize.(float64), currSize.(float64)) {
				log.Printf("[INFO] BLOCK STORAGE UPGRADE STARTS")
				currentName := d.Get(tfconstants.AttrName).(string)

				// Get VM ID from block storage
				var vmID float64
				if blockStorage.VMDetail != nil {
					if vmIDVal, ok := blockStorage.VMDetail[volumeVMIDKey]; ok {
						if vmIDFloat, ok := vmIDVal.(float64); ok {
							vmID = vmIDFloat
						}
					}
				}

				if vmID == 0 {
					d.Set(tfconstants.AttrSize, prevSize)
					return diag.Errorf(ErrorResizeVMIDMissingTemplate, blockStorageID)
				}

				upgradeReq := &goe2e.BlockStorageUpgradeRequest{
					Name: currentName,
					Size: currSize.(float64),
					VMID: vmID,
				}
				log.Printf("[INFO] BlockStorage details for update : %+v", upgradeReq)

				_, err := goe2eClient.BlockStorage.UpgradeBlockStorage(ctx, blockStorageID, upgradeReq)
				if err != nil {
					d.Set(tfconstants.AttrSize, prevSize)
					if checkErrorForSpecificMessage(err, goe2econstants.NodeLCMStateDiskResize) || checkErrorForSpecificMessage(err, goe2econstants.NodeLCMStateDiskResizePoweroff) {
						return diag.Errorf(ErrorResizeConcurrentDiskOperationTemplate, blockStorageID, currentName, vmID)
					}
					return diag.Errorf(tfconstants.ResourceOperationByIDErrorTemplate, tfconstants.OperationUpdating, ResourceName, blockStorageID, projectIDStr, region, err)
				}
				log.Printf("[INFO] BLOCK STORAGE UPGRADE completed successfully")
				// Continue to read to refresh state
			} else {
				d.Set(tfconstants.AttrSize, prevSize)
				return diag.Errorf(ErrorResizeReduceNotAllowedTemplate, blockStorageID)
			}
		} else {
			d.Set(tfconstants.AttrSize, prevSize)
			return diag.Errorf(ErrorResizeRequiresAttachmentTemplate, blockStorageID)
		}
	}

	// Handle tags updates (state-only, no API call)
	if d.HasChange(tfconstants.AttrTags) {
		oldTags, newTags := d.GetChange(tfconstants.AttrTags)
		log.Printf("[DEBUG] Block storage tags changed. Old: %v, New: %v", oldTags, newTags)
		// Tags are stored in state only, no API call needed
		if err := d.Set(tfconstants.AttrTags, newTags); err != nil {
			return diag.FromErr(err)
		}
	}

	// Name changes are now handled by ForceNew: true in schema
	// No need for manual check here
	return resourceReadBlockStorage(ctx, d, m)
}

func resourceDeleteBlockStorage(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	var diags diag.Diagnostics
	blockStorageID := d.Id()

	// Get project_id with provider default support
	projectIDStr, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get region with provider default support
	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Create goe2e client with projectID and region
	goe2eClient, err := cfg.Goe2eClientForProject(projectIDStr, region)
	if err != nil {
		return diag.Errorf(tfconstants.ErrorCreatingGoe2eClient, err)
	}

	status := d.Get(tfconstants.AttrStatus).(string)
	vmID := d.Get(tfconstants.AttrVMID).(string)

	if status == goe2econstants.BlockStorageStatusSaving || status == goe2econstants.BlockStorageStatusCreating {
		return diag.Errorf(ErrorDeleteInStateTemplate, blockStorageID, status, projectIDStr, region)
	}
	if status == goe2econstants.BlockStorageStatusAttached {
		return diag.Errorf(ErrorDeleteWhileAttachedTemplate, blockStorageID, vmID)
	}

	log.Printf("[INFO] Deleting block storage: %s (region: %s, project: %s)", blockStorageID, region, projectIDStr)

	_, err = goe2eClient.BlockStorage.DeleteBlockStorage(ctx, blockStorageID)
	if err != nil {
		// Idempotency: if already deleted, succeed
		if strings.Contains(err.Error(), goe2econstants.NotFoundSubstring) || strings.Contains(err.Error(), goe2econstants.NotFoundCode) {
			log.Printf("[WARN] Block storage %s already deleted", blockStorageID)
			d.SetId("")
			return diags
		}
		return diag.Errorf("Error deleting "+ResourceName+" (ID: %s): %s", blockStorageID, err)
	}

	log.Printf("[INFO] Successfully deleted block storage: %s", blockStorageID)
	d.SetId("")
	return diags
}

func resourceExistsBlockStorage(d *schema.ResourceData, m interface{}) (bool, error) {
	ctx := context.Background()
	cfg := m.(*config.Config)

	blockStorageID := d.Id()

	// Get region with provider default support
	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return false, err
	}

	// Get project_id with provider default support
	projectIDStr, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return false, err
	}

	// Create goe2e client with projectID and region
	goe2eClient, err := cfg.Goe2eClientForProject(projectIDStr, region)
	if err != nil {
		return false, fmt.Errorf("error creating goe2e client: %s", err)
	}

	blockStorage, resp, err := goe2eClient.BlockStorage.GetBlockStorage(ctx, blockStorageID)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return false, nil
		}
		if strings.Contains(err.Error(), goe2econstants.NotFoundSubstring) {
			return false, nil
		}
		return false, err
	}

	return blockStorage != nil, nil
}

// customImportBlockStorage handles importing block storage resources
// Expected import format: project_id/region/block_storage_id
func customImportBlockStorage(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != ImportIDPartsCount {
		return nil, fmt.Errorf(tfconstants.ImportIDInvalidFormatTemplate, d.Id(), ImportIDFormatDescription)
	}

	projectID := parts[0]
	region := parts[1]
	blockStorageID := parts[2]

	// Set the individual fields with correct schema keys
	d.Set(tfconstants.AttrProjectID, projectID)
	d.Set(tfconstants.AttrRegion, region)
	d.SetId(blockStorageID)

	// Initialize empty tags
	d.Set(tfconstants.AttrTags, make(map[string]string))

	log.Printf(LogImportTemplate, blockStorageID, projectID, region)

	return []*schema.ResourceData{d}, nil
}

// CalculateIOPS calculates IOPS based on size (size * 15)
// Exported for testing purposes
func CalculateIOPS(size float64) string {
	iops := size * IOPS_PER_GB
	return strconv.Itoa(int(iops))
}

func validateSize(ctx context.Context, d *schema.ResourceData, goe2eClient *goe2e.Client) error {
	requestedSize := d.Get(tfconstants.AttrSize).(float64)

	// Try to get plans from API to validate against available sizes
	plans, _, err := goe2eClient.BlockStorage.GetBlockStoragePlans(ctx)
	if err != nil {
		// If API call fails, fallback to known valid sizes
		log.Printf("[WARN] Failed to retrieve block storage plans from API, using fallback validation: %s", err)
	} else if len(plans) > 0 {
		// Parse size from plan name (e.g., "250 GB", "1 TB")
		// Note: Actual parsing logic depends on API response format
		log.Printf("[DEBUG] Retrieved %d block storage plans from API", len(plans))
		for _, plan := range plans {
			log.Printf("[DEBUG] Available plan: %+v", plan)
		}
	}

	// Validate against known valid sizes
	for _, validSize := range validBlockStorageSizes {
		if requestedSize == validSize {
			return nil
		}
	}

	return fmt.Errorf("block storage size %.0f GB is not valid. Valid sizes are: %s GB", requestedSize, validBlockStorageSizesString)
}

func checkErrorForSpecificMessage(err error, message string) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), message)
}
