package ssh_key

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	e2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// resourceSshKeyResourceV0 returns the schema for the previous (V0) state version.
// This is used during state upgrades from V0 to V1 to properly deserialize
// legacy state files that were created with the original SSH key resource schema.
// The V0 schema must match what was previously defined to ensure Terraform can
// correctly unmarshal the state data.
func resourceSshKeyResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			e2econstants.AttrRegion:    config.RegionSchema(),
			e2econstants.AttrLocation:  config.LocationSchema(),
			e2econstants.AttrProjectID: config.ProjectIDSchemaResource(),
			e2econstants.AttrLabel: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "label of the SSH Key",
			},
			e2econstants.AttrSSHKey: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "the SSH public key content",
			},
			e2econstants.AttrCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "timestamp when the SSH Key was created",
			},
		},
	}
}

// resourceSshKeyStateUpgradeV0toV1 upgrades state from V0 to V1.
// This function is called by Terraform when loading an old state file to upgrade it
// to the current schema version. It maintains backward compatibility by:
// 1. Preserving all V0 fields (label, ssh_key, location) as-is
// 2. Initializing new V3 fields (tags) with empty/default values
// 3. NOT automatically renaming fields - users control migration timing
// This approach allows gradual adoption of new field names without forcing changes.
func resourceSshKeyStateUpgradeV0toV1(
	ctx context.Context,
	rawState map[string]interface{},
	meta interface{},
) (map[string]interface{}, error) {
	// Preserve all existing V0 fields (label, ssh_key, location)
	// These continue to work for backward compatibility

	// Initialize new V3 fields with sensible defaults
	if _, exists := rawState["tags"]; !exists {
		rawState["tags"] = make(map[string]interface{})
	}

	// Extract values from deprecated fields for new preferred fields
	// Users can gradually migrate by updating their HCL
	if label, ok := rawState[e2econstants.AttrLabel]; ok && label != "" {
		// Don't auto-set 'name' - let users explicitly migrate
		log.Printf("[INFO] Upgraded SSH key state from v0 to v1: id=%s, label=%s",
			rawState["id"], label)
	}

	return rawState, nil
}

func ResourceSshKey() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceSshKeyResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceSshKeyStateUpgradeV0toV1,
				Version: 0,
			},
		},

		Schema: map[string]*schema.Schema{
			// ============================================
			// COMMON INFRASTRUCTURE FIELDS
			// ============================================
			e2econstants.AttrRegion: {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{e2econstants.AttrLocation},
				Description:   "region where the SSH key is stored",
				ValidateFunc:  validation.StringIsNotEmpty,
			},
			e2econstants.AttrLocation: {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Deprecated:    "Use 'region' instead. This field will be removed in v4.0.0.",
				ConflictsWith: []string{e2econstants.AttrRegion},
				Description:   "location where the SSH key is stored (deprecated, use region)",
			},
			e2econstants.AttrProjectID: config.ProjectIDSchemaResource(),

			// ============================================
			// REQUIRED IDENTIFICATION FIELDS (V3 Preferred)
			// ============================================
			e2econstants.AttrName: {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{e2econstants.AttrLabel},
				ExactlyOneOf:  []string{e2econstants.AttrName, e2econstants.AttrLabel},
				Description:   "name of the SSH key",
				ValidateFunc:  validation.StringIsNotEmpty,
			},
			"public_key": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Sensitive:     false, // Public keys are not sensitive
				ConflictsWith: []string{e2econstants.AttrSSHKey},
				ExactlyOneOf:  []string{"public_key", e2econstants.AttrSSHKey},
				Description:   "the public key material",
				ValidateFunc:  validation.StringIsNotEmpty,
			},

			// ============================================
			// DEPRECATED FIELDS (V2 Compatibility)
			// ============================================
			e2econstants.AttrLabel: {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Deprecated:    "Use 'name' instead. This field will be removed in v5.0.0.",
				ConflictsWith: []string{e2econstants.AttrName},
				ExactlyOneOf:  []string{e2econstants.AttrName, e2econstants.AttrLabel},
				Description:   "label of the SSH key (deprecated, use name)",
				ValidateFunc:  validation.StringIsNotEmpty,
			},
			e2econstants.AttrSSHKey: {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Sensitive:     false,
				Deprecated:    "Use 'public_key' instead. This field will be removed in v5.0.0.",
				ConflictsWith: []string{"public_key"},
				ExactlyOneOf:  []string{"public_key", e2econstants.AttrSSHKey},
				Description:   "the SSH public key content (deprecated, use public_key)",
				ValidateFunc:  validation.StringIsNotEmpty,
			},

			// ============================================
			// V3 OPTIONAL FIELDS
			// ============================================
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "map of tags to assign to the resource (state-only, API support pending)",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			// ============================================
			// COMPUTED FIELDS
			// ============================================
			e2econstants.AttrCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "timestamp when the SSH key was created",
			},
			e2econstants.AttrProjectName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "name of the project associated with the SSH key",
			},
		},

		CreateContext: resourceCreateSshKey,
		ReadContext:   resourceReadSshKey,
		UpdateContext: resourceUpdateSshKey,
		DeleteContext: resourceDeleteSshKey,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSshKeyImport,
		},
	}
}

// getKeyName returns the SSH key name, preferring the new V3 'name' field
// over the deprecated V2 'label' field. This enables transparent backward
// compatibility - existing configs using 'label' continue to work while
// new configs can use the preferred 'name' field.
func getKeyName(d *schema.ResourceData) string {
	if name, ok := d.GetOk(e2econstants.AttrName); ok {
		return name.(string)
	}
	return d.Get(e2econstants.AttrLabel).(string)
}

// getPublicKey returns the public key material, preferring the new V3 'public_key' field
// over the deprecated V2 'ssh_key' field. This enables transparent backward
// compatibility - existing configs using 'ssh_key' continue to work while
// new configs can use the preferred AWS-aligned 'public_key' field.
func getPublicKey(d *schema.ResourceData) string {
	if pk, ok := d.GetOk("public_key"); ok {
		return pk.(string)
	}
	return d.Get(e2econstants.AttrSSHKey).(string)
}

// setKeyName sets both the new V3 'name' field and the deprecated V2 'label' field.
// This ensures Terraform state contains both values, allowing users to reference
// either field in their configuration without errors. The provider logs deprecation
// warnings separately when deprecated fields are used.
func setKeyName(d *schema.ResourceData, name string) error {
	// Set the preferred V3 field
	if err := d.Set(e2econstants.AttrName, name); err != nil {
		return err
	}
	// Also set deprecated V2 field for backward compatibility
	return d.Set(e2econstants.AttrLabel, name)
}

// setPublicKey sets both the new V3 'public_key' field and the deprecated V2 'ssh_key' field.
// This ensures Terraform state contains both values, allowing users to reference
// either field in their configuration without errors. The provider logs deprecation
// warnings separately when deprecated fields are used.
func setPublicKey(d *schema.ResourceData, publicKey string) error {
	// Set the preferred V3 field
	if err := d.Set("public_key", publicKey); err != nil {
		return err
	}
	// Also set deprecated V2 field for backward compatibility
	return d.Set(e2econstants.AttrSSHKey, publicKey)
}

func resourceCreateSshKey(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	goe2eClient := cfg.Goe2eClient()
	var diags diag.Diagnostics

	log.Printf("[INFO] Creating SSH key")

	// Get region from configuration or provider defaults
	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Extract key name using helper (prefers V3 'name', falls back to V2 'label')
	keyName := getKeyName(d)
	// Extract public key using helper (prefers V3 'public_key', falls back to V2 'ssh_key')
	publicKey := getPublicKey(d)

	log.Printf("[DEBUG] Creating SSH key: name=%s, region=%s", keyName, region)

	// Create the SSH key via goe2e client (which includes retry/backoff support)
	createReq := &goe2e.SSHKeyCreateRequest{
		Label:  keyName,
		SSHKey: publicKey,
	}

	sshKey, _, err := goe2eClient.SSHKeys.CreateSSHKey(ctx, createReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create SSH key: %w", err))
	}

	// Store the API's primary key (PK) as the Terraform resource ID
	d.SetId(strconv.Itoa(sshKey.PK))

	// Populate state with all returned fields, using setters for backward compatibility
	// This ensures both V2 and V3 field names are present in state
	if err := setKeyName(d, sshKey.Label); err != nil {
		return diag.FromErr(err)
	}
	if err := setPublicKey(d, sshKey.SSHKey); err != nil {
		return diag.FromErr(err)
	}

	// Set computed fields
	if err := d.Set(e2econstants.AttrCreatedAt, sshKey.Timestamp); err != nil {
		return diag.FromErr(err)
	}

	// Set region in both V3 preferred field and V2 deprecated field
	if err := d.Set(e2econstants.AttrRegion, region); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(e2econstants.AttrLocation, region); err != nil {
		return diag.FromErr(err)
	}

	// Initialize empty tags map for state-only tag support
	// Tags are not persisted to the API in V3.0 but are stored in Terraform state
	if _, ok := d.GetOk("tags"); !ok {
		if err := d.Set("tags", make(map[string]interface{})); err != nil {
			return diag.FromErr(err)
		}
	}

	log.Printf("[INFO] SSH key created: id=%s", d.Id())
	return diags
}

func resourceReadSshKey(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	goe2eClient := cfg.Goe2eClient()
	var diags diag.Diagnostics

	pk := d.Id()
	if pk == "" {
		return diag.Errorf("SSH key ID is empty, invalid state")
	}

	log.Printf("[DEBUG] Reading SSH key: id=%s", pk)

	// Fetch SSH key using goe2e client
	sshKey, _, err := goe2eClient.SSHKeys.GetSSHKey(ctx, pk)
	if err != nil {
		// Check if it's a "not found" error
		if strings.Contains(err.Error(), "not found") {
			log.Printf("[WARN] SSH key with ID %s not found, removing from state", pk)
			d.SetId("")
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "SSH key not found",
				Detail:   "The SSH key may have been deleted outside of Terraform.",
			}}
		}
		return diag.FromErr(fmt.Errorf("failed to read SSH key (ID: %s): %w", pk, err))
	}

	if sshKey == nil {
		log.Printf("[WARN] SSH key with ID %s not found", pk)
		d.SetId("")
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "SSH key not found",
			Detail:   "The SSH key may have been deleted outside of Terraform.",
		}}
	}

	// Set all fields (both preferred and deprecated for backward compatibility)
	if err := setKeyName(d, sshKey.Label); err != nil {
		return diag.FromErr(err)
	}
	if err := setPublicKey(d, sshKey.SSHKey); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(e2econstants.AttrCreatedAt, sshKey.Timestamp); err != nil {
		return diag.FromErr(err)
	}

	// Set region (store the value from state or config)
	region, _ := cfg.GetRegionOrDefault(d)
	if err := d.Set(e2econstants.AttrRegion, region); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(e2econstants.AttrLocation, region); err != nil {
		return diag.FromErr(err)
	}

	// Tags are state-only, preserve existing tags
	if tags, ok := d.GetOk("tags"); ok {
		if err := d.Set("tags", tags); err != nil {
			return diag.FromErr(err)
		}
	}

	log.Printf("[DEBUG] SSH key read successfully: id=%s", pk)
	return diags
}

// resourceUpdateSshKey handles tag updates (state-only in V3.0)
// SSH key itself is immutable, so only tags can be updated
func resourceUpdateSshKey(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	pk := d.Id()
	log.Printf("[DEBUG] Updating SSH key: id=%s", pk)

	// Only tags are updateable (state-only, not sent to API)
	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		log.Printf("[DEBUG] SSH key tags changed. Old: %v, New: %v", oldTags, newTags)
		// Tags are stored in state only, no API call needed
		if err := d.Set("tags", newTags); err != nil {
			return diag.FromErr(err)
		}
	}

	// ForceNew fields should trigger recreation (not update)
	// If other fields change, Terraform will handle the recreation

	log.Printf("[INFO] SSH key updated: id=%s", pk)
	return diags
}

func resourceDeleteSshKey(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	goe2eClient := cfg.Goe2eClient()
	var diags diag.Diagnostics

	pk := d.Id()
	if pk == "" {
		return diag.Errorf("SSH key ID is empty, cannot delete")
	}

	log.Printf("[DEBUG] Deleting SSH key: id=%s", pk)

	_, err := goe2eClient.SSHKeys.DeleteSSHKey(ctx, pk)
	if err != nil {
		// Check if key not found (treat as success since we want it gone)
		if strings.Contains(err.Error(), "not found") {
			log.Printf("[WARN] SSH key not found during delete (already deleted), treating as success")
			d.SetId("")
			return diags
		}
		return diag.FromErr(fmt.Errorf("failed to delete SSH key (ID: %s): %w", pk, err))
	}

	d.SetId("")
	log.Printf("[INFO] SSH key deleted: id=%s", pk)
	return diags
}

func resourceSshKeyImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	cfg := m.(*config.Config)
	goe2eClient := cfg.Goe2eClient()

	// Expected format: project_id:ssh_key_id or project_id:region:ssh_key_id
	parts := strings.Split(d.Id(), ":")

	var projectID, region, sshKeyID string

	if len(parts) == 2 {
		projectID = parts[0]
		sshKeyID = parts[1]
		// Get default region from config
		region = cfg.DefaultRegion
		if region == "" {
			return nil, fmt.Errorf("region must be specified in import ID or provider default_region must be set. Use format: project_id:region:ssh_key_id")
		}
	} else if len(parts) == 3 {
		projectID = parts[0]
		region = parts[1]
		sshKeyID = parts[2]
	} else {
		return nil, fmt.Errorf("invalid import ID format. Expected: project_id:ssh_key_id or project_id:region:ssh_key_id")
	}

	// Fetch SSH key from API to populate all fields
	sshKey, _, err := goe2eClient.SSHKeys.GetSSHKey(ctx, sshKeyID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch SSH key for import (ID: %s): %w", sshKeyID, err)
	}

	if sshKey == nil {
		return nil, fmt.Errorf("SSH key not found during import (ID: %s)", sshKeyID)
	}

	// Set resource ID
	d.SetId(sshKeyID)

	// Set all fields with V3 preferred names
	if err := setKeyName(d, sshKey.Label); err != nil {
		return nil, err
	}
	if err := setPublicKey(d, sshKey.SSHKey); err != nil {
		return nil, err
	}

	// Set infrastructure fields
	if err := d.Set(e2econstants.AttrProjectID, projectID); err != nil {
		return nil, err
	}
	if err := d.Set(e2econstants.AttrRegion, region); err != nil {
		return nil, err
	}
	if err := d.Set(e2econstants.AttrLocation, region); err != nil {
		return nil, err
	}

	// Set computed fields
	if err := d.Set(e2econstants.AttrCreatedAt, sshKey.Timestamp); err != nil {
		return nil, err
	}

	// Initialize empty tags
	if err := d.Set("tags", make(map[string]interface{})); err != nil {
		return nil, err
	}

	log.Printf("[INFO] SSH key imported: id=%s, name=%s", sshKeyID, sshKey.Label)
	return []*schema.ResourceData{d}, nil
}
