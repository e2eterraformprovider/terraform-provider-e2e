package dbaas_mariadb

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceMariaDB() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceMariaDBResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceMariaDBStateUpgradeV0toV1,
				Version: 0,
			},
		},

		Schema: map[string]*schema.Schema{
			// ============================================
			// COMMON FIELDS
			// ============================================
			tfconstants.AttrRegion:    config.RegionSchema(),
			tfconstants.AttrLocation:  config.LocationSchema(),
			tfconstants.AttrProjectID: config.ProjectIDSchemaResource(),

			// ============================================
			// REQUIRED IMMUTABLE FIELDS
			// ============================================
			tfconstants.AttrName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "name of the MariaDB DBaaS instance",
			},
			"software_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "the software name (e.g., MariaDB)",
			},
			"software_version": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "the software version (e.g., 10.6)",
			},
			tfconstants.AttrGroup: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "the group name for the MariaDB DBaaS instance (e.g., Default)",
			},
			tfconstants.AttrDatabase: {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "database configuration (user, password, database name)",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "the database username",
						},
						"password": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "the database password",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "name of the database to create",
						},
						"dbaas_number": {
							Type:        schema.TypeInt,
							Required:    true,
							ForceNew:    true,
							Description: "the DBaaS number (typically 1 for single instance)",
						},
					},
				},
			},

			// ============================================
			// OPTIONAL MUTABLE FIELDS - CONFIGURATION
			// ============================================
			tfconstants.AttrPlan: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "the plan name specifying CPU/memory (e.g., DBS.16GB)",
			},
			tfconstants.AttrVPCs: {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "list of VPC ids to associate",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			tfconstants.AttrPublicIPEnabled: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "whether to attach a public IP during creation or update",
			},
			tfconstants.AttrParameterGroupID: {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "id of the parameter group to attach (use 0 to skip)",
			},

			// ============================================
			// OPTIONAL IMMUTABLE FIELDS - SECURITY
			// ============================================
			tfconstants.AttrIsEncryptionEnabled: {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "whether to enable encryption at rest for the MariaDB cluster",
			},
			tfconstants.AttrEncryptionPassphrase: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "",
				Sensitive:   true,
				Description: "passphrase for encryption (leave empty if encryption is not enabled)",
			},

			// ============================================
			// POWER MANAGEMENT
			// ============================================
			tfconstants.AttrStatus: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"STOPPED", "RUNNING", "RESTARTING"},
					false,
				),
				Description: "the operational status of the MariaDB DBaaS instance (STOPPED, RUNNING, or RESTARTING)",
			},

			// ============================================
			// DISK MANAGEMENT
			// ============================================
			tfconstants.AttrDiskSize: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "additional disk size in gigabytes to expand during update",
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
			// COMPUTED FIELDS - IDENTIFIERS
			// ============================================
			"software_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "id of the software",
			},
			tfconstants.AttrTemplateID: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "id of the template",
			},

			// ============================================
			// COMPUTED FIELDS - STATUS
			// ============================================
			"public_ip_attached": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "whether a public IP is currently attached",
			},

			// ============================================
			// COMPUTED FIELDS - NETWORK
			// ============================================
			tfconstants.AttrPublicIPAddress: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the MariaDB DBaaS instances public ipv4 address",
			},
			tfconstants.AttrPrivateIPAddress: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the MariaDB DBaaS instances private ipv4 address",
			},
			"port": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the port number on which the MariaDB service is accessible",
			},

			// ============================================
			// COMPUTED FIELDS - RESOURCES
			// ============================================
			"total_disk_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "the total disk size in gigabytes after expansion",
			},
		},

		CreateContext: resourceCreateMariaDB,
		ReadContext:   resourceReadMariaDB,
		UpdateContext: resourceUpdateMariaDB,
		DeleteContext: resourceDeleteMariaDB,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCreateMariaDB(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	goe2eClient := cfg.Goe2eClient()
	var diags diag.Diagnostics

	projectID, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	softwareName := d.Get("software_name").(string)
	softwareVersion := d.Get("software_version").(string)
	planName := d.Get(tfconstants.AttrPlan).(string)

	// Get software ID using goe2e client
	softwareID, err := goe2eClient.MariaDB.GetSoftwareID(ctx, softwareName, softwareVersion)
	if err != nil {
		return diag.Errorf("error retrieving %s software ID for version (%s) in project (%s), region (%s): %s", softwareName, softwareVersion, projectID, region, err)
	}

	// Get template ID using goe2e client
	templateID, err := goe2eClient.MariaDB.GetTemplateID(ctx, planName, softwareID)
	if err != nil {
		return diag.Errorf("error retrieving %s template ID for plan (%s) in project (%s), region (%s): %s", softwareName, planName, projectID, region, err)
	}

	// Extract database configuration
	dbConfigList := d.Get(tfconstants.AttrDatabase).([]interface{})
	dbConfigMap := dbConfigList[0].(map[string]interface{})

	// Extract VPC IDs
	var vpcIDs []string
	for _, v := range d.Get(tfconstants.AttrVPCs).([]interface{}) {
		vpcIDs = append(vpcIDs, v.(string))
	}

	// Expand VPC list using goe2e client
	var vpcList []goe2e.VPCMetadata
	if len(vpcIDs) > 0 {
		vpcList, err = goe2eClient.MariaDB.ExpandVPCList(ctx, vpcIDs)
		if err != nil {
			return diag.Errorf("error preparing VPC list for MariaDB DBaaS in project (%s), region (%s): %s", projectID, region, err)
		}
	}

	publicIPEnabled := d.Get(tfconstants.AttrPublicIPEnabled).(bool)

	parameterGroupID := 0
	if v, ok := d.GetOk(tfconstants.AttrParameterGroupID); ok {
		parameterGroupID = v.(int)
	}

	isEncryptionEnabled := false
	if v, ok := d.GetOk(tfconstants.AttrIsEncryptionEnabled); ok {
		isEncryptionEnabled = v.(bool)
	}

	encryptionPassphrase := ""
	if v, ok := d.GetOk(tfconstants.AttrEncryptionPassphrase); ok {
		encryptionPassphrase = v.(string)
	}

	// Build create request using goe2e types
	req := &goe2e.MariaDBCreateRequest{
		Name:                 d.Get("name").(string),
		SoftwareID:           softwareID,
		TemplateID:           templateID,
		PublicIPRequired:     publicIPEnabled,
		Group:                d.Get(tfconstants.AttrGroup).(string),
		VPCs:                 vpcList,
		PGID:                 parameterGroupID,
		IsEncryptionEnabled:  isEncryptionEnabled,
		EncryptionPassphrase: encryptionPassphrase,
		Database: goe2e.DBConfig{
			User:        dbConfigMap["user"].(string),
			Password:    dbConfigMap["password"].(string),
			Name:        dbConfigMap["name"].(string),
			DBaaSNumber: dbConfigMap["dbaas_number"].(int),
		},
	}

	// Create MariaDB cluster using goe2e client
	mariaDB, _, err := goe2eClient.MariaDB.CreateMariaDB(ctx, req)
	if err != nil {
		return diag.Errorf("error creating MariaDB DBaaS (name: %s) in project (%s), region (%s): %s", req.Name, projectID, region, err)
	}

	// Set resource ID and attributes
	d.SetId(fmt.Sprintf("%d", mariaDB.ID))
	if err := d.Set("name", mariaDB.Name); err != nil {
		return diag.FromErr(err)
	}

	// Normalize status (SUSPENDED → STOPPED)
	status := mariaDB.Status
	if status == "SUSPENDED" {
		status = "STOPPED"
	}
	if err := d.Set(tfconstants.AttrStatus, status); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(tfconstants.AttrPublicIPAddress, mariaDB.MasterNode.PublicIPAddress); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(tfconstants.AttrPrivateIPAddress, mariaDB.MasterNode.PrivateIPAddress); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("port", mariaDB.MasterNode.Port); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("software_id", softwareID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(tfconstants.AttrTemplateID, templateID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("public_ip_attached", mariaDB.MasterNode.PublicIPAddress != ""); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Successfully created MariaDB cluster: %s (ID: %d)", mariaDB.Name, mariaDB.ID)

	return diags
}

func resourceReadMariaDB(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	goe2eClient := cfg.Goe2eClient()
	var diags diag.Diagnostics

	id := d.Id()

	// Get MariaDB cluster using goe2e client
	mariaDB, _, err := goe2eClient.MariaDB.GetMariaDB(ctx, id)
	if err != nil {
		return diag.Errorf("error retrieving MariaDB DBaaS (ID: %s): %s", id, err)
	}

	// Check if resource was deleted
	if mariaDB == nil {
		log.Printf("[WARN] MariaDB cluster %s not found, removing from state", id)
		d.SetId("")
		return diags
	}

	// Set basic attributes
	if err := d.Set("name", mariaDB.Name); err != nil {
		return diag.FromErr(err)
	}

	// Normalize status (SUSPENDED → STOPPED)
	status := mariaDB.Status
	if status == "SUSPENDED" {
		status = "STOPPED"
	}
	if err := d.Set(tfconstants.AttrStatus, status); err != nil {
		return diag.FromErr(err)
	}

	// Set software information
	if err := d.Set("software_name", mariaDB.Software.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("software_version", mariaDB.Software.Version); err != nil {
		return diag.FromErr(err)
	}

	// Set plan information
	if err := d.Set(tfconstants.AttrPlan, mariaDB.MasterNode.Plan.Name); err != nil {
		return diag.FromErr(err)
	}

	// Set network information
	if err := d.Set(tfconstants.AttrPublicIPAddress, mariaDB.MasterNode.PublicIPAddress); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(tfconstants.AttrPrivateIPAddress, mariaDB.MasterNode.PrivateIPAddress); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("port", mariaDB.MasterNode.Port); err != nil {
		return diag.FromErr(err)
	}

	// Set computed fields
	if err := d.Set("public_ip_attached", mariaDB.MasterNode.PublicIPAddress != ""); err != nil {
		return diag.FromErr(err)
	}

	// Parse disk size from string to int
	diskSize := 0
	if mariaDB.MasterNode.Disk != "" {
		if size, err := strconv.Atoi(mariaDB.MasterNode.Disk); err == nil {
			diskSize = size
		}
	}
	if err := d.Set("total_disk_size", diskSize); err != nil {
		return diag.FromErr(err)
	}

	// Set encryption status
	if err := d.Set(tfconstants.AttrIsEncryptionEnabled, mariaDB.IsEncryptionEnabled); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Successfully read MariaDB cluster: %s (ID: %s)", mariaDB.Name, id)

	return diags
}

func resourceDeleteMariaDB(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	goe2eClient := cfg.Goe2eClient()
	var diags diag.Diagnostics

	id := d.Id()

	log.Printf("[INFO] Deleting MariaDB cluster: %s", id)

	// Delete MariaDB cluster using goe2e client
	_, err := goe2eClient.MariaDB.DeleteMariaDB(ctx, id)
	if err != nil {
		// Check if already deleted (404 or "not found" error)
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "404") {
			log.Printf("[WARN] MariaDB cluster %s already deleted", id)
			d.SetId("")
			return diags
		}
		return diag.Errorf("error deleting MariaDB DBaaS (ID: %s): %s", id, err)
	}

	log.Printf("[INFO] Successfully deleted MariaDB cluster: %s", id)
	d.SetId("")
	return diags
}

func resourceUpdateMariaDB(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	goe2eClient := cfg.Goe2eClient()
	id := d.Id()

	// Handle status changes (power management)
	if d.HasChange("status") {
		newStatus := d.Get("status").(string)
		log.Printf("[INFO] Status change detected for MariaDB cluster %s: %s", id, newStatus)

		switch strings.ToUpper(newStatus) {
		case "STOPPED":
			if _, err := goe2eClient.MariaDB.ShutdownMariaDB(ctx, id); err != nil {
				// Rollback disk_size and plan changes on failure
				if d.HasChange(tfconstants.AttrDiskSize) {
					d.Set(tfconstants.AttrDiskSize, 0)
				}
				if d.HasChange(tfconstants.AttrPlan) {
					oldPlan, _ := d.GetChange(tfconstants.AttrPlan)
					d.Set(tfconstants.AttrPlan, oldPlan.(string))
				}
				return diag.Errorf("error stopping MariaDB DBaaS (ID: %s): %s", id, err)
			}
			log.Printf("[INFO] Successfully stopped MariaDB cluster %s", id)
		case "RUNNING":
			if _, err := goe2eClient.MariaDB.ResumeMariaDB(ctx, id); err != nil {
				return diag.Errorf("error starting MariaDB DBaaS (ID: %s): %s", id, err)
			}
			log.Printf("[INFO] Successfully started MariaDB cluster %s", id)
		case "RESTARTING":
			if _, err := goe2eClient.MariaDB.RestartMariaDB(ctx, id); err != nil {
				return diag.Errorf("error restarting MariaDB DBaaS (ID: %s): %s", id, err)
			}
			log.Printf("[INFO] Successfully restarted MariaDB cluster %s", id)
		default:
			return diag.Errorf("error updating MariaDB DBaaS (ID: %s): unsupported status value: %s. Must be one of: STOPPED, RUNNING, RESTARTING", id, newStatus)
		}
	}

	// Handle VPC changes
	if d.HasChange(tfconstants.AttrVPCs) {
		oldRaw, newRaw := d.GetChange(tfconstants.AttrVPCs)
		oldVPCSet := expandStringSet(oldRaw.([]interface{}))
		newVPCSet := expandStringSet(newRaw.([]interface{}))

		var toDetach, toAttach []string

		// Find VPCs to detach
		for vpc := range oldVPCSet {
			if _, exists := newVPCSet[vpc]; !exists {
				toDetach = append(toDetach, vpc)
			}
		}

		// Find VPCs to attach
		for vpc := range newVPCSet {
			if _, exists := oldVPCSet[vpc]; !exists {
				toAttach = append(toAttach, vpc)
			}
		}

		// Detach VPCs
		if len(toDetach) > 0 {
			log.Printf("[INFO] Detaching VPCs from MariaDB cluster %s: %v", id, toDetach)
			if _, err := goe2eClient.MariaDB.DetachVPC(ctx, id, toDetach); err != nil {
				return diag.Errorf("error detaching VPC from MariaDB DBaaS (ID: %s): %s", id, err)
			}
		}

		// Attach VPCs
		if len(toAttach) > 0 {
			log.Printf("[INFO] Attaching VPCs to MariaDB cluster %s: %v", id, toAttach)
			if _, err := goe2eClient.MariaDB.AttachVPC(ctx, id, toAttach); err != nil {
				return diag.Errorf("error attaching VPC to MariaDB DBaaS (ID: %s): %s", id, err)
			}
		}
	}

	// Handle public IP changes
	if d.HasChange(tfconstants.AttrPublicIPEnabled) {
		newVal := d.Get(tfconstants.AttrPublicIPEnabled).(bool)
		log.Printf("[INFO] Public IP change detected for MariaDB cluster %s: %v", id, newVal)

		if newVal {
			if _, err := goe2eClient.MariaDB.AttachPublicIP(ctx, id); err != nil {
				return diag.Errorf("error attaching public IP to MariaDB DBaaS (ID: %s): %s", id, err)
			}
			log.Printf("[INFO] Successfully attached public IP to MariaDB cluster %s", id)
		} else {
			if _, err := goe2eClient.MariaDB.DetachPublicIP(ctx, id); err != nil {
				return diag.Errorf("error detaching public IP from MariaDB DBaaS (ID: %s): %s", id, err)
			}
			log.Printf("[INFO] Successfully detached public IP from MariaDB cluster %s", id)
		}
	}

	// Handle parameter group changes
	if d.HasChange(tfconstants.AttrParameterGroupID) {
		oldRaw, newRaw := d.GetChange(tfconstants.AttrParameterGroupID)
		oldPGID := oldRaw.(int)
		newPGID := newRaw.(int)

		log.Printf("[INFO] Parameter group change detected for MariaDB cluster %s: %d -> %d", id, oldPGID, newPGID)

		switch {
		case oldPGID != 0 && newPGID == 0:
			// Detach parameter group
			if _, err := goe2eClient.MariaDB.DetachParameterGroup(ctx, id, oldPGID); err != nil {
				return diag.Errorf("error detaching parameter group (ID: %d) from MariaDB DBaaS (ID: %s): %s", oldPGID, id, err)
			}
			log.Printf("[INFO] Successfully detached parameter group %d from MariaDB cluster %s", oldPGID, id)
		case newPGID != 0 && newPGID != oldPGID:
			// Attach new parameter group
			if _, err := goe2eClient.MariaDB.AttachParameterGroup(ctx, id, newPGID); err != nil {
				return diag.Errorf("error attaching parameter group (ID: %d) to MariaDB DBaaS (ID: %s): %s", newPGID, id, err)
			}
			log.Printf("[INFO] Successfully attached parameter group %d to MariaDB cluster %s", newPGID, id)
		}
	}

	// Handle plan upgrade
	if d.HasChange(tfconstants.AttrPlan) {
		oldPlan, newPlan := d.GetChange(tfconstants.AttrPlan)
		log.Printf("[INFO] Plan change detected for MariaDB cluster %s: %s -> %s", id, oldPlan.(string), newPlan.(string))

		// Verify cluster is stopped
		status := d.Get("status").(string)
		if strings.ToUpper(status) != "STOPPED" {
			d.Set(tfconstants.AttrPlan, oldPlan.(string))
			return diag.Errorf("cannot upgrade plan for MariaDB DBaaS (ID: %s): database must be in STOPPED state (current state: %s). Please stop the instance first", id, status)
		}

		// Get software ID and template ID
		softwareName := d.Get("software_name").(string)
		softwareVersion := d.Get("software_version").(string)

		softwareID, err := goe2eClient.MariaDB.GetSoftwareID(ctx, softwareName, softwareVersion)
		if err != nil {
			d.Set(tfconstants.AttrPlan, oldPlan.(string))
			return diag.Errorf("error retrieving %s software ID for version (%s) while upgrading plan for DBaaS (ID: %s): %s", softwareName, softwareVersion, id, err)
		}

		templateID, err := goe2eClient.MariaDB.GetTemplateID(ctx, newPlan.(string), softwareID)
		if err != nil {
			d.Set(tfconstants.AttrPlan, oldPlan.(string))
			return diag.Errorf("error retrieving %s template ID for plan (%s) while upgrading DBaaS (ID: %s): %s", softwareName, newPlan.(string), id, err)
		}

		// Upgrade the plan
		if _, err := goe2eClient.MariaDB.UpgradePlan(ctx, id, templateID); err != nil {
			d.Set(tfconstants.AttrPlan, oldPlan.(string))
			return diag.Errorf("error upgrading MariaDB DBaaS (ID: %s) plan from (%s) to (%s): %s", id, oldPlan.(string), newPlan.(string), err)
		}

		log.Printf("[INFO] Successfully upgraded MariaDB cluster %s to plan %s (template_id=%d)", id, newPlan, templateID)

		// Update template ID in state
		if err := d.Set(tfconstants.AttrTemplateID, templateID); err != nil {
			return diag.FromErr(err)
		}
	}

	// Handle disk expansion
	if d.HasChange(tfconstants.AttrDiskSize) {
		additionalSize := d.Get(tfconstants.AttrDiskSize).(int)

		if additionalSize > 0 {
			log.Printf("[INFO] Disk expansion requested for MariaDB cluster %s: +%d GB", id, additionalSize)

			// Verify cluster is stopped
			status := d.Get("status").(string)
			if strings.ToUpper(status) != "STOPPED" {
				d.Set(tfconstants.AttrDiskSize, 0)
				return diag.Errorf("cannot expand disk for MariaDB DBaaS (ID: %s): database must be in STOPPED state (current state: %s)", id, status)
			}

			// Expand the disk
			if _, err := goe2eClient.MariaDB.ExpandDisk(ctx, id, additionalSize); err != nil {
				d.Set(tfconstants.AttrDiskSize, 0)
				return diag.Errorf("error expanding MariaDB DBaaS (ID: %s) disk by %d GB: %s", id, additionalSize, err)
			}

			log.Printf("[INFO] Successfully expanded disk by %d GB for MariaDB cluster %s", additionalSize, id)

			// Reset disk_size to 0 after expansion
			if err := d.Set(tfconstants.AttrDiskSize, 0); err != nil {
				return diag.FromErr(err)
			}
		} else {
			log.Printf("[DEBUG] disk_size is 0 for MariaDB cluster %s, skipping expansion", id)
		}
	}

	// Handle tags (state-only, no API call)
	if d.HasChange("tags") {
		log.Printf("[INFO] MariaDB cluster %s: tags updated (state-only, not sent to API)", id)
	}

	// Refresh state by reading the resource
	return resourceReadMariaDB(ctx, d, m)
}

func expandStringSet(list []interface{}) map[string]struct{} {
	result := make(map[string]struct{})
	for _, v := range list {
		result[v.(string)] = struct{}{}
	}
	return result
}

// ============================================
// STATE MIGRATION: V0 → V1
// ============================================

// resourceMariaDBResourceV0 returns the V0 schema definition (before tags were added)
func resourceMariaDBResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// ============================================
			// COMMON FIELDS
			// ============================================
			tfconstants.AttrRegion:    config.RegionSchema(),
			tfconstants.AttrLocation:  config.LocationSchema(),
			tfconstants.AttrProjectID: config.ProjectIDSchemaResource(),

			// ============================================
			// REQUIRED IMMUTABLE FIELDS
			// ============================================
			tfconstants.AttrName: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"software_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"software_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			tfconstants.AttrGroup: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			tfconstants.AttrDatabase: {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"password": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"dbaas_number": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			// ============================================
			// OPTIONAL MUTABLE FIELDS - CONFIGURATION
			// ============================================
			tfconstants.AttrPlan: {
				Type:     schema.TypeString,
				Required: true,
			},
			tfconstants.AttrVPCs: {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			tfconstants.AttrPublicIPEnabled: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			tfconstants.AttrParameterGroupID: {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			// ============================================
			// OPTIONAL IMMUTABLE FIELDS - SECURITY
			// ============================================
			tfconstants.AttrIsEncryptionEnabled: {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			tfconstants.AttrEncryptionPassphrase: {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Default:   "",
				Sensitive: true,
			},

			// ============================================
			// POWER MANAGEMENT
			// ============================================
			tfconstants.AttrStatus: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			// ============================================
			// DISK MANAGEMENT
			// ============================================
			tfconstants.AttrDiskSize: {
				Type:     schema.TypeInt,
				Optional: true,
			},

			// ============================================
			// COMPUTED FIELDS - IDENTIFIERS
			// ============================================
			"software_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			tfconstants.AttrTemplateID: {
				Type:     schema.TypeInt,
				Computed: true,
			},

			// ============================================
			// COMPUTED FIELDS - STATUS
			// ============================================
			"public_ip_attached": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			// ============================================
			// COMPUTED FIELDS - NETWORK
			// ============================================
			tfconstants.AttrPublicIPAddress: {
				Type:     schema.TypeString,
				Computed: true,
			},
			tfconstants.AttrPrivateIPAddress: {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// ============================================
			// COMPUTED FIELDS - RESOURCES
			// ============================================
			"total_disk_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			// Note: V0 schema does not have "tags" field
		},
	}
}

// resourceMariaDBStateUpgradeV0toV1 upgrades the state from V0 to V1
// V1 adds the "tags" field for state-only tag management
func resourceMariaDBStateUpgradeV0toV1(
	ctx context.Context,
	rawState map[string]interface{},
	meta interface{},
) (map[string]interface{}, error) {
	// Add new V3 fields with defaults
	if _, exists := rawState["tags"]; !exists {
		rawState["tags"] = make(map[string]interface{})
	}

	// Preserve all existing fields - no modifications needed
	log.Printf("[INFO] Upgraded MariaDB resource state from v0 to v1: %s", rawState["id"])

	return rawState, nil
}
