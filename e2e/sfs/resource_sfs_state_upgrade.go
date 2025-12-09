package sfs

import (
	"context"
	"log"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	e2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceSfsResourceV0 returns the V0 schema (legacy)
func resourceSfsResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Common fields
			e2econstants.AttrRegion:    config.RegionSchema(),
			e2econstants.AttrLocation:  config.LocationSchema(),
			e2econstants.AttrProjectID: config.ProjectIDSchemaResource(),

			// Core identity and configuration
			e2econstants.AttrName: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			e2econstants.AttrPlan: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			e2econstants.AttrVPCID: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			e2econstants.AttrDiskSize: {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"disk_iops": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			e2econstants.AttrStatus: {
				Type:     schema.TypeString,
				Computed: true,
			},
			e2econstants.AttrEncryptionPassphrase: {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				ForceNew:  true,
				Default:   "",
			},
			e2econstants.AttrIsEncryptionEnabled: {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
		},
	}
}

// resourceSfsStateUpgradeV0toV1 handles upgrade from schema version 0 to 1
func resourceSfsStateUpgradeV0toV1(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	log.Printf("[INFO] Upgrading SFS state from V0 to V1")

	// Handle nil state
	if rawState == nil {
		rawState = make(map[string]interface{})
	}

	// Copy disk_size -> size_gb if disk_size exists
	if diskSize, ok := rawState[e2econstants.AttrDiskSize]; ok && diskSize != nil {
		log.Printf("[DEBUG] Migrating disk_size to size_gb: %v", diskSize)
		rawState[e2econstants.AttrSizeGB] = diskSize
	}

	// Copy disk_iops -> iops if disk_iops exists
	if diskIops, ok := rawState["disk_iops"]; ok && diskIops != nil {
		log.Printf("[DEBUG] Migrating disk_iops to iops: %v", diskIops)
		rawState[e2econstants.AttrIOPS] = diskIops
	}

	// Copy is_encryption_enabled -> encryption_enabled if is_encryption_enabled exists
	if isEncryption, ok := rawState[e2econstants.AttrIsEncryptionEnabled]; ok && isEncryption != nil {
		log.Printf("[DEBUG] Migrating is_encryption_enabled to encryption_enabled: %v", isEncryption)
		rawState[e2econstants.AttrEncryptionEnabled] = isEncryption
	}

	// Initialize new computed fields with reasonable defaults
	// state will be populated on next read from API
	if _, ok := rawState[e2econstants.AttrState]; !ok {
		rawState[e2econstants.AttrState] = ""
	}

	// mount_endpoint is alias for private_endpoint - copy if private_endpoint exists
	if privateEndpoint, ok := rawState["private_endpoint"]; ok && privateEndpoint != nil {
		rawState["mount_endpoint"] = privateEndpoint
	} else if _, ok := rawState["mount_endpoint"]; !ok {
		rawState["mount_endpoint"] = ""
	}

	// Initialize tags if not present
	if _, ok := rawState[e2econstants.AttrTags]; !ok {
		rawState[e2econstants.AttrTags] = map[string]interface{}{}
	}

	// is_backup_enabled will be populated on next read from API
	if _, ok := rawState["is_backup_enabled"]; !ok {
		rawState["is_backup_enabled"] = false
	}

	// created_at will be populated on next read from API
	if _, ok := rawState[e2econstants.AttrCreatedAt]; !ok {
		rawState[e2econstants.AttrCreatedAt] = ""
	}

	// private_endpoint will be populated on next read from API
	if _, ok := rawState["private_endpoint"]; !ok {
		rawState["private_endpoint"] = ""
	}

	log.Printf("[INFO] SFS state upgrade from V0 to V1 completed successfully")
	return rawState, nil
}
