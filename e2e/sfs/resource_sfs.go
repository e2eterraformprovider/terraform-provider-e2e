package sfs

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/node"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceSfs() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			// Common fields
			tfconstants.AttrRegion:    config.RegionSchema(),
			tfconstants.AttrLocation:  config.LocationSchema(),
			tfconstants.AttrProjectID: config.ProjectIDSchemaResource(),

			// Resource identity and configuration
			tfconstants.AttrName: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "name of the SFS",
				ValidateFunc: validateName,
			},
			tfconstants.AttrPlan: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "the plan of the SFS",
			},
			tfconstants.AttrVPCID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "id of the VPC for the SFS",
			},

			// V3 Storage fields (preferred)
			tfconstants.AttrSizeGB: {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{tfconstants.AttrDiskSize},
				Description:   "the size of the SFS volume in gigabytes",
			},
			tfconstants.AttrIOPS: {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"disk_iops"},
				Description:   "the IOPS value of the SFS",
			},

			// V2 Storage fields (deprecated)
			tfconstants.AttrDiskSize: {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				Deprecated:    "Use size_gb instead. This field will be removed in v4.0.",
				ConflictsWith: []string{tfconstants.AttrSizeGB},
				Description:   "DEPRECATED: Use size_gb instead. The size of the disk in gigabytes.",
			},
			"disk_iops": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				Deprecated:    "Use iops instead. This field will be removed in v4.0.",
				ConflictsWith: []string{tfconstants.AttrIOPS},
				Description:   "DEPRECATED: Use iops instead. The IOPS of the disk.",
			},

			// V3 Encryption field (preferred)
			tfconstants.AttrEncryptionEnabled: {
				Type:          schema.TypeBool,
				Optional:      true,
				ForceNew:      true,
				Default:       false,
				ConflictsWith: []string{tfconstants.AttrIsEncryptionEnabled},
				Description:   "whether to enable encryption for the SFS",
			},

			// V2 Encryption field (deprecated)
			tfconstants.AttrIsEncryptionEnabled: {
				Type:          schema.TypeBool,
				Optional:      true,
				ForceNew:      true,
				Default:       false,
				Deprecated:    "Use encryption_enabled instead. This field will be removed in v4.0.",
				ConflictsWith: []string{tfconstants.AttrEncryptionEnabled},
				Description:   "DEPRECATED: Use encryption_enabled instead. Whether encryption is enabled.",
			},

			// Encryption passphrase (shared by both versions)
			tfconstants.AttrEncryptionPassphrase: {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				ForceNew:    true,
				Default:     "",
				Description: "passphrase for encryption, if encryption is enabled. This field is optional and should only be set if encryption_enabled or is_encryption_enabled is true.",
			},

			// Tags (state-only until API support)
			tfconstants.AttrTags: {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "map of tags for the SFS instance",
			},

			// Computed fields - Status
			tfconstants.AttrStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the API status of the SFS instance (Creating, Active, Error, Deleting)",
			},
			tfconstants.AttrState: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the normalized state of the SFS instance (creating, active, error, deleting)",
			},

			// Computed fields - Networking
			"private_endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the NFS mount endpoint for the SFS",
			},
			"mount_endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "alias for private_endpoint - the NFS mount endpoint for the SFS",
			},

			// Computed fields - Storage and backup
			"is_backup_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "whether backups are enabled for the SFS",
			},

			// Computed fields - Metadata
			tfconstants.AttrCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the creation timestamp of the SFS",
			},
		},
		CreateContext: resourceCreateSfs,
		ReadContext:   resourceReadSfs,
		DeleteContext: resourceDeleteSfs,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceSfsResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceSfsStateUpgradeV0toV1,
				Version: 0,
			},
		},
		Importer: &schema.ResourceImporter{
			State: node.CustomImportStateFunc,
		},
	}
}

func validateName(v interface{}, k string) (ws []string, es []error) {

	var errs []error
	var warns []string
	value, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected name to be string"))
		return warns, errs
	}
	whiteSpace := regexp.MustCompile(`\s+`)
	if whiteSpace.Match([]byte(value)) {
		errs = append(errs, fmt.Errorf("name cannot contain whitespace. Got %s", value))
		return warns, errs
	}
	return warns, errs
}

func resourceCreateSfs(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	// Create goe2e client with specific projectID and region
	client, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return diag.Errorf("Error creating goe2e client: %s", err)
	}

	log.Printf("[INFO] SFS CREATE STARTS")

	name := d.Get(tfconstants.AttrName).(string)
	plan := d.Get(tfconstants.AttrPlan).(string)
	vpcID := d.Get(tfconstants.AttrVPCID).(string)

	// Get size: prefer V3 field size_gb over V2 field disk_size
	sizeGB := getEffectiveSizeGB(d, tfconstants.AttrSizeGB, tfconstants.AttrDiskSize, 0)
	if sizeGB == 0 {
		return diag.Errorf("Error creating SFS (name: %s): size_gb or disk_size must be specified", name)
	}

	// Log deprecation warning if old field is used
	if _, ok := d.GetOk(tfconstants.AttrDiskSize); ok && !d.HasChanges(tfconstants.AttrDiskSize) {
		if _, ok2 := d.GetOk(tfconstants.AttrSizeGB); !ok2 {
			logDeprecationWarning(tfconstants.AttrDiskSize, tfconstants.AttrSizeGB)
		}
	}

	// Get IOPS: prefer V3 field iops over V2 field disk_iops
	iops := getEffectiveIOPS(d, tfconstants.AttrIOPS, "disk_iops", 0)
	if iops == 0 {
		return diag.Errorf("Error creating SFS (name: %s): iops or disk_iops must be specified", name)
	}

	// Log deprecation warning if old field is used
	if _, ok := d.GetOk("disk_iops"); ok && !d.HasChanges("disk_iops") {
		if _, ok2 := d.GetOk(tfconstants.AttrIOPS); !ok2 {
			logDeprecationWarning("disk_iops", tfconstants.AttrIOPS)
		}
	}

	// Get encryption: prefer V3 field encryption_enabled over V2 field is_encryption_enabled
	isEncrypted := getEffectiveEncryptionEnabled(d, tfconstants.AttrEncryptionEnabled, tfconstants.AttrIsEncryptionEnabled)

	// Log deprecation warning if old encryption field is used
	if _, ok := d.GetOk(tfconstants.AttrIsEncryptionEnabled); ok && !d.HasChanges(tfconstants.AttrIsEncryptionEnabled) {
		if _, ok2 := d.GetOk(tfconstants.AttrEncryptionEnabled); !ok2 {
			logDeprecationWarning(tfconstants.AttrIsEncryptionEnabled, tfconstants.AttrEncryptionEnabled)
		}
	}

	createReq := &goe2e.SfsCreateRequest{
		Name:                name,
		Plan:                plan,
		VPCID:               vpcID,
		DiskSize:            sizeGB,
		DiskIOPS:            iops,
		IsEncryptionEnabled: isEncrypted,
	}

	if pass, ok := d.GetOk(tfconstants.AttrEncryptionPassphrase); ok {
		createReq.EncryptionPassphrase = pass.(string)
	}

	sfs, _, err := client.Sfs.CreateSfs(ctx, createReq)
	if err != nil {
		return diag.Errorf("Error creating SFS (name: %s) in project (%s), region (%s): %s", createReq.Name, projectID, region, err)
	}

	log.Printf("[INFO] SFS CREATE | RESPONSE: %+v", sfs)

	if sfs == nil || sfs.ID == "" {
		return diag.Errorf("Error creating SFS (name: %s) in project (%s), region (%s): unable to retrieve valid 'efs_id' from API response", createReq.Name, projectID, region)
	}

	d.SetId(sfs.ID)

	// Store tags in state (state-only until API supports it)
	if tags, ok := d.GetOk(tfconstants.AttrTags); ok {
		if err := d.Set(tfconstants.AttrTags, tags); err != nil {
			return diag.FromErr(fmt.Errorf("error setting tags: %w", err))
		}
	}

	// Poll for SFS to become Active
	if err := waitForSfsActive(ctx, client, sfs.ID); err != nil {
		return diag.Errorf("Error waiting for SFS (ID: %s) to become active in project (%s), region (%s): %s", sfs.ID, projectID, region, err)
	}

	return diags
}

func resourceReadSfs(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	var diags diag.Diagnostics

	log.Printf("[INFO] Inside SFS Resource Read")
	sfsID := d.Id()

	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	projectID, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Create goe2e client with specific projectID and region
	client, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return diag.Errorf("Error creating goe2e client: %s", err)
	}

	sfs, _, err := client.Sfs.GetSfs(ctx, sfsID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "404") {
			log.Printf("[WARN] SFS with ID %s not found", sfsID)
			d.SetId("")
			return diags
		}
		return diag.Errorf("Error retrieving SFS (ID: %s) in project (%s), region (%s): %s", sfsID, projectID, region, err)
	}

	if sfs == nil {
		log.Printf("[WARN] SFS with ID %s not found", sfsID)
		d.SetId("")
		return diags
	}

	// Set core identity fields
	if err := d.Set(tfconstants.AttrName, sfs.Name); err != nil {
		return diag.FromErr(fmt.Errorf("error setting name: %w", err))
	}
	if err := d.Set(tfconstants.AttrPlan, sfs.PlanName); err != nil {
		return diag.FromErr(fmt.Errorf("error setting plan: %w", err))
	}
	if err := d.Set(tfconstants.AttrVPCID, sfs.VPCID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting vpc_id: %w", err))
	}

	// Set size fields - prefer V3 fields but maintain V2 for backwards compatibility
	if sfs.DiskSize != "" {
		diskSizeStr := strings.TrimSpace(strings.ReplaceAll(sfs.DiskSize, "GB", ""))
		if sizeInt, err := strconv.Atoi(diskSizeStr); err == nil {
			// Set both V3 and V2 fields for backwards compatibility
			if err := d.Set(tfconstants.AttrSizeGB, sizeInt); err != nil {
				return diag.FromErr(fmt.Errorf("error setting size_gb: %w", err))
			}
			if err := d.Set(tfconstants.AttrDiskSize, sizeInt); err != nil {
				return diag.FromErr(fmt.Errorf("error setting disk_size: %w", err))
			}
		}
	}

	// Set IOPS fields - prefer V3 fields but maintain V2 for backwards compatibility
	if err := d.Set(tfconstants.AttrIOPS, sfs.DiskIOPS); err != nil {
		return diag.FromErr(fmt.Errorf("error setting iops: %w", err))
	}
	if err := d.Set("disk_iops", sfs.DiskIOPS); err != nil {
		return diag.FromErr(fmt.Errorf("error setting disk_iops: %w", err))
	}

	// Set encryption fields - prefer V3 fields but maintain V2 for backwards compatibility
	if err := d.Set(tfconstants.AttrEncryptionEnabled, sfs.IsEncryptionEnabled); err != nil {
		return diag.FromErr(fmt.Errorf("error setting encryption_enabled: %w", err))
	}
	if err := d.Set(tfconstants.AttrIsEncryptionEnabled, sfs.IsEncryptionEnabled); err != nil {
		return diag.FromErr(fmt.Errorf("error setting is_encryption_enabled: %w", err))
	}

	// Set computed fields
	if err := d.Set(tfconstants.AttrStatus, sfs.Status); err != nil {
		return diag.FromErr(fmt.Errorf("error setting status: %w", err))
	}

	// Normalize and set state field
	normalizedState := normalizeSfsState(sfs.Status)
	if err := d.Set(tfconstants.AttrState, normalizedState); err != nil {
		return diag.FromErr(fmt.Errorf("error setting state: %w", err))
	}

	// Set networking fields
	if err := d.Set("private_endpoint", sfs.PrivateIPAddress); err != nil {
		return diag.FromErr(fmt.Errorf("error setting private_endpoint: %w", err))
	}
	if err := d.Set("mount_endpoint", sfs.PrivateIPAddress); err != nil {
		return diag.FromErr(fmt.Errorf("error setting mount_endpoint: %w", err))
	}

	// Set backup field
	if err := d.Set("is_backup_enabled", sfs.IsBackupEnabled); err != nil {
		return diag.FromErr(fmt.Errorf("error setting is_backup_enabled: %w", err))
	}

	return diags
}

func resourceDeleteSfs(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	var diags diag.Diagnostics
	sfsID := d.Id()

	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	projectID, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Check if SFS is in Creating state
	status := d.Get(tfconstants.AttrStatus).(string)
	if status == "Creating" {
		return diag.Errorf("Cannot delete SFS (ID: %s): SFS is in Creating state in project (%s), region (%s). Please wait for SFS creation to complete", sfsID, projectID, region)
	}

	// Create goe2e client with specific projectID and region
	client, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return diag.Errorf("Error creating goe2e client: %s", err)
	}

	_, err = client.Sfs.DeleteSfs(ctx, sfsID)
	if err != nil {
		return diag.Errorf("Error deleting SFS (ID: %s) in project (%s), region (%s): %s", sfsID, projectID, region, err)
	}

	d.SetId("")
	return diags
}
