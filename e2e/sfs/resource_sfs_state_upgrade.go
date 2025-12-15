package sfs

import (
	"context"
	"log"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceSfsResourceV0 returns the V0 schema (legacy)
func resourceSfsResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Common fields
			tfconstants.AttrRegion:    config.RegionSchema(),
			tfconstants.AttrLocation:  config.LocationSchema(),
			tfconstants.AttrProjectID: config.ProjectIDSchemaResource(),

			// Core identity and configuration
			tfconstants.AttrName: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			tfconstants.AttrPlan: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			tfconstants.AttrVPCID: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			tfconstants.AttrDiskSize: {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			tfconstants.AttrDiskIOPS: {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			tfconstants.AttrStatus: {
				Type:     schema.TypeString,
				Computed: true,
			},
			tfconstants.AttrEncryptionPassphrase: {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				ForceNew:  true,
				Default:   "",
			},
			tfconstants.AttrIsEncryptionEnabled: {
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
	if diskSize, ok := rawState[tfconstants.FieldMigrationKeyDiskSize]; ok && diskSize != nil {
		log.Printf("[DEBUG] Migrating disk_size to size_gb: %v", diskSize)
		rawState[tfconstants.FieldMigrationKeySizeGB] = diskSize
	}

	// Copy disk_iops -> iops if disk_iops exists
	if diskIops, ok := rawState[tfconstants.FieldMigrationKeyDiskIOPS]; ok && diskIops != nil {
		log.Printf("[DEBUG] Migrating disk_iops to iops: %v", diskIops)
		rawState[tfconstants.FieldMigrationKeyIOPS] = diskIops
	}

	// Copy is_encryption_enabled -> encryption_enabled if is_encryption_enabled exists
	if isEncryption, ok := rawState[tfconstants.FieldMigrationKeyIsEncryptionEnabled]; ok && isEncryption != nil {
		log.Printf("[DEBUG] Migrating is_encryption_enabled to encryption_enabled: %v", isEncryption)
		rawState[tfconstants.FieldMigrationKeyEncryptionEnabled] = isEncryption
	}

	// Initialize new computed fields with reasonable defaults
	// state will be populated on next read from API
	if _, ok := rawState[tfconstants.AttrState]; !ok {
		rawState[tfconstants.AttrState] = ""
	}

	// mount_endpoint is alias for private_endpoint - copy if private_endpoint exists
	if privateEndpoint, ok := rawState[tfconstants.AttrPrivateEndpoint]; ok && privateEndpoint != nil {
		rawState[tfconstants.AttrMountEndpoint] = privateEndpoint
	} else if _, ok := rawState[tfconstants.AttrMountEndpoint]; !ok {
		rawState[tfconstants.AttrMountEndpoint] = ""
	}

	// Initialize tags if not present
	if _, ok := rawState[tfconstants.AttrTags]; !ok {
		rawState[tfconstants.AttrTags] = map[string]interface{}{}
	}

	// is_backup_enabled will be populated on next read from API
	if _, ok := rawState[tfconstants.AttrIsBackupEnabled]; !ok {
		rawState[tfconstants.AttrIsBackupEnabled] = false
	}

	// created_at will be populated on next read from API
	if _, ok := rawState[tfconstants.AttrCreatedAt]; !ok {
		rawState[tfconstants.AttrCreatedAt] = ""
	}

	// private_endpoint will be populated on next read from API
	if _, ok := rawState[tfconstants.AttrPrivateEndpoint]; !ok {
		rawState[tfconstants.AttrPrivateEndpoint] = ""
	}

	log.Printf("[INFO] SFS state upgrade from V0 to V1 completed successfully")
	return rawState, nil
}
