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
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
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
				Upgrade: ResourceMariaDBStateUpgradeV0toV1,
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
			tfconstants.AttrSoftwareName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "the software name (e.g., MariaDB)",
			},
			tfconstants.AttrSoftwareVersion: {
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
						tfconstants.AttrDatabaseBlockUser: {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "the database username",
						},
						tfconstants.AttrDatabaseBlockPassword: {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "the database password",
						},
						tfconstants.AttrDatabaseBlockName: {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "name of the database to create",
						},
						tfconstants.AttrDatabaseBlockDBaaSNumber: {
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
				Default:     tfconstants.DBaaSDefaultPublicIPEnabled,
				Description: "whether to attach a public IP during creation or update",
			},
			tfconstants.AttrParameterGroupID: {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     tfconstants.DBaaSDefaultParameterGroupID,
				Description: "id of the parameter group to attach (use 0 to skip)",
			},

			// ============================================
			// OPTIONAL IMMUTABLE FIELDS - SECURITY
			// ============================================
			tfconstants.AttrIsEncryptionEnabled: {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     tfconstants.DBaaSDefaultIsEncryptionEnabled,
				Description: "whether to enable encryption at rest for the MariaDB cluster",
			},
			tfconstants.AttrEncryptionPassphrase: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     tfconstants.DBaaSDefaultEncryptionPassphrase,
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
					[]string{
						goe2econstants.DBaaSStatusStopped,
						goe2econstants.DBaaSStatusRunning,
						goe2econstants.DBaaSStatusRestarting,
					},
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
			tfconstants.AttrTags: {
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
			tfconstants.AttrSoftwareID: {
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
			tfconstants.AttrPublicIPAttached: {
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
			tfconstants.AttrPort: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the port number on which the MariaDB service is accessible",
			},

			// ============================================
			// COMPUTED FIELDS - RESOURCES
			// ============================================
			tfconstants.AttrTotalDiskSize: {
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

	softwareName := d.Get(tfconstants.AttrSoftwareName).(string)
	softwareVersion := d.Get(tfconstants.AttrSoftwareVersion).(string)
	planName := d.Get(tfconstants.AttrPlan).(string)

	// Get software ID using goe2e client
	softwareID, err := goe2eClient.MariaDB.GetSoftwareID(ctx, softwareName, softwareVersion)
	if err != nil {
		return diag.Errorf(ErrorRetrievingSoftwareIDTemplate, softwareName, softwareVersion, projectID, region, err)
	}

	// Get template ID using goe2e client
	templateID, err := goe2eClient.MariaDB.GetTemplateID(ctx, planName, softwareID)
	if err != nil {
		return diag.Errorf(ErrorRetrievingTemplateIDTemplate, softwareName, planName, projectID, region, err)
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
			return diag.Errorf(ErrorPreparingVPCListTemplate, projectID, region, err)
		}
	}

	publicIPEnabled := d.Get(tfconstants.AttrPublicIPEnabled).(bool)

	parameterGroupID := tfconstants.DBaaSDefaultParameterGroupID
	if v, ok := d.GetOk(tfconstants.AttrParameterGroupID); ok {
		parameterGroupID = v.(int)
	}

	isEncryptionEnabled := tfconstants.DBaaSDefaultIsEncryptionEnabled
	if v, ok := d.GetOk(tfconstants.AttrIsEncryptionEnabled); ok {
		isEncryptionEnabled = v.(bool)
	}

	encryptionPassphrase := tfconstants.DBaaSDefaultEncryptionPassphrase
	if v, ok := d.GetOk(tfconstants.AttrEncryptionPassphrase); ok {
		encryptionPassphrase = v.(string)
	}

	// Build create request using goe2e types
	req := &goe2e.MariaDBCreateRequest{
		Name:                 d.Get(tfconstants.AttrName).(string),
		SoftwareID:           softwareID,
		TemplateID:           templateID,
		PublicIPRequired:     publicIPEnabled,
		Group:                d.Get(tfconstants.AttrGroup).(string),
		VPCs:                 vpcList,
		PGID:                 parameterGroupID,
		IsEncryptionEnabled:  isEncryptionEnabled,
		EncryptionPassphrase: encryptionPassphrase,
		Database: goe2e.DBConfig{
			User:        dbConfigMap[tfconstants.AttrDatabaseBlockUser].(string),
			Password:    dbConfigMap[tfconstants.AttrDatabaseBlockPassword].(string),
			Name:        dbConfigMap[tfconstants.AttrDatabaseBlockName].(string),
			DBaaSNumber: dbConfigMap[tfconstants.AttrDatabaseBlockDBaaSNumber].(int),
		},
	}

	// Create MariaDB cluster using goe2e client
	mariaDB, _, err := goe2eClient.MariaDB.CreateMariaDB(ctx, req)
	if err != nil {
		return diag.Errorf(tfconstants.ResourceOperationErrorTemplate, tfconstants.OperationCreating, ResourceName, req.Name, projectID, region, err)
	}

	// Set resource ID and attributes
	d.SetId(fmt.Sprintf("%d", mariaDB.ID))
	if err := d.Set(tfconstants.AttrName, mariaDB.Name); err != nil {
		return diag.FromErr(err)
	}

	// Normalize status (SUSPENDED → STOPPED)
	status := normalizeStatus(mariaDB.Status)
	if err := d.Set(tfconstants.AttrStatus, status); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(tfconstants.AttrPublicIPAddress, mariaDB.MasterNode.PublicIPAddress); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(tfconstants.AttrPrivateIPAddress, mariaDB.MasterNode.PrivateIPAddress); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(tfconstants.AttrPort, mariaDB.MasterNode.Port); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(tfconstants.AttrSoftwareID, softwareID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(tfconstants.AttrTemplateID, templateID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(tfconstants.AttrPublicIPAttached, mariaDB.MasterNode.PublicIPAddress != ""); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Successfully created %s: %s (ID: %d)", ResourceName, mariaDB.Name, mariaDB.ID)

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
		return diag.Errorf(tfconstants.ResourceOperationByIDErrorTemplate, tfconstants.OperationRetrieving, ResourceName, id, "", "", err)
	}

	// Check if resource was deleted
	if mariaDB == nil {
		log.Printf("[WARN] %s %s not found, removing from state", ResourceName, id)
		d.SetId("")
		return diags
	}

	// Set basic attributes
	if err := d.Set(tfconstants.AttrName, mariaDB.Name); err != nil {
		return diag.FromErr(err)
	}

	// Normalize status (SUSPENDED → STOPPED)
	status := normalizeStatus(mariaDB.Status)
	if err := d.Set(tfconstants.AttrStatus, status); err != nil {
		return diag.FromErr(err)
	}

	// Set software information
	if err := d.Set(tfconstants.AttrSoftwareName, mariaDB.Software.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(tfconstants.AttrSoftwareVersion, mariaDB.Software.Version); err != nil {
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
	if err := d.Set(tfconstants.AttrPort, mariaDB.MasterNode.Port); err != nil {
		return diag.FromErr(err)
	}

	// Set computed fields
	if err := d.Set(tfconstants.AttrPublicIPAttached, mariaDB.MasterNode.PublicIPAddress != ""); err != nil {
		return diag.FromErr(err)
	}

	// Parse disk size from string to int
	diskSize := 0
	if mariaDB.MasterNode.Disk != "" {
		if size, err := strconv.Atoi(mariaDB.MasterNode.Disk); err == nil {
			diskSize = size
		}
	}
	if err := d.Set(tfconstants.AttrTotalDiskSize, diskSize); err != nil {
		return diag.FromErr(err)
	}

	// Set encryption status
	if err := d.Set(tfconstants.AttrIsEncryptionEnabled, mariaDB.IsEncryptionEnabled); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Successfully read %s: %s (ID: %s)", ResourceName, mariaDB.Name, id)

	return diags
}

func resourceDeleteMariaDB(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	goe2eClient := cfg.Goe2eClient()
	var diags diag.Diagnostics

	id := d.Id()

	log.Printf("[INFO] Deleting %s: %s", ResourceName, id)

	// Delete MariaDB cluster using goe2e client
	_, err := goe2eClient.MariaDB.DeleteMariaDB(ctx, id)
	if err != nil {
		// Check if already deleted (404 or "not found" error)
		if strings.Contains(err.Error(), goe2econstants.NotFoundSubstring) || strings.Contains(err.Error(), goe2econstants.NotFoundCode) {
			log.Printf("[WARN] %s %s already deleted", ResourceName, id)
			d.SetId("")
			return diags
		}
		return diag.Errorf(tfconstants.ResourceOperationByIDErrorTemplate, tfconstants.OperationDeleting, ResourceName, id, "", "", err)
	}

	log.Printf("[INFO] Successfully deleted %s: %s", ResourceName, id)
	d.SetId("")
	return diags
}

func resourceUpdateMariaDB(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	goe2eClient := cfg.Goe2eClient()
	id := d.Id()

	// Handle status changes (power management)
	if d.HasChange(tfconstants.AttrStatus) {
		newStatus := d.Get(tfconstants.AttrStatus).(string)
		log.Printf("[INFO] Status change detected for %s %s: %s", ResourceName, id, newStatus)

		// Normalize status (SUSPENDED → STOPPED) for consistent handling
		normalizedStatus := normalizeStatus(strings.ToUpper(newStatus))
		switch normalizedStatus {
		case goe2econstants.DBaaSStatusStopped:
			if _, err := goe2eClient.MariaDB.ShutdownMariaDB(ctx, id); err != nil {
				// Rollback disk_size and plan changes on failure
				if d.HasChange(tfconstants.AttrDiskSize) {
					d.Set(tfconstants.AttrDiskSize, 0)
				}
				if d.HasChange(tfconstants.AttrPlan) {
					oldPlan, _ := d.GetChange(tfconstants.AttrPlan)
					d.Set(tfconstants.AttrPlan, oldPlan.(string))
				}
				return diag.Errorf(ErrorStoppingTemplate, id, err)
			}
			log.Printf("[INFO] Successfully stopped %s %s", ResourceName, id)
		case goe2econstants.DBaaSStatusRunning:
			if _, err := goe2eClient.MariaDB.ResumeMariaDB(ctx, id); err != nil {
				return diag.Errorf(ErrorStartingTemplate, id, err)
			}
			log.Printf("[INFO] Successfully started %s %s", ResourceName, id)
		case goe2econstants.DBaaSStatusRestarting:
			if _, err := goe2eClient.MariaDB.RestartMariaDB(ctx, id); err != nil {
				return diag.Errorf(ErrorRestartingTemplate, id, err)
			}
			log.Printf("[INFO] Successfully restarted %s %s", ResourceName, id)
		default:
			return diag.Errorf(ErrorUnsupportedStatusTemplate, id, newStatus)
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
			log.Printf("[INFO] Detaching VPCs from %s %s: %v", ResourceName, id, toDetach)
			if _, err := goe2eClient.MariaDB.DetachVPC(ctx, id, toDetach); err != nil {
				return diag.Errorf(ErrorDetachingVPCTemplate, id, err)
			}
		}

		// Attach VPCs
		if len(toAttach) > 0 {
			log.Printf("[INFO] Attaching VPCs to %s %s: %v", ResourceName, id, toAttach)
			if _, err := goe2eClient.MariaDB.AttachVPC(ctx, id, toAttach); err != nil {
				return diag.Errorf(ErrorAttachingVPCTemplate, id, err)
			}
		}
	}

	// Handle public IP changes
	if d.HasChange(tfconstants.AttrPublicIPEnabled) {
		newVal := d.Get(tfconstants.AttrPublicIPEnabled).(bool)
		log.Printf("[INFO] Public IP change detected for %s %s: %v", ResourceName, id, newVal)

		if newVal {
			if _, err := goe2eClient.MariaDB.AttachPublicIP(ctx, id); err != nil {
				return diag.Errorf(ErrorAttachingPublicIPTemplate, id, err)
			}
			log.Printf("[INFO] Successfully attached public IP to %s %s", ResourceName, id)
		} else {
			if _, err := goe2eClient.MariaDB.DetachPublicIP(ctx, id); err != nil {
				return diag.Errorf(ErrorDetachingPublicIPTemplate, id, err)
			}
			log.Printf("[INFO] Successfully detached public IP from %s %s", ResourceName, id)
		}
	}

	// Handle parameter group changes
	if d.HasChange(tfconstants.AttrParameterGroupID) {
		oldRaw, newRaw := d.GetChange(tfconstants.AttrParameterGroupID)
		oldPGID := oldRaw.(int)
		newPGID := newRaw.(int)

		log.Printf("[INFO] Parameter group change detected for %s %s: %d -> %d", ResourceName, id, oldPGID, newPGID)

		switch {
		case oldPGID != 0 && newPGID == 0:
			// Detach parameter group
			if _, err := goe2eClient.MariaDB.DetachParameterGroup(ctx, id, oldPGID); err != nil {
				return diag.Errorf(ErrorDetachingParameterGroupTemplate, oldPGID, id, err)
			}
			log.Printf("[INFO] Successfully detached parameter group %d from %s %s", oldPGID, ResourceName, id)
		case newPGID != 0 && newPGID != oldPGID:
			// Attach new parameter group
			if _, err := goe2eClient.MariaDB.AttachParameterGroup(ctx, id, newPGID); err != nil {
				return diag.Errorf(ErrorAttachingParameterGroupTemplate, newPGID, id, err)
			}
			log.Printf("[INFO] Successfully attached parameter group %d to %s %s", newPGID, ResourceName, id)
		}
	}

	// Handle plan upgrade
	if d.HasChange(tfconstants.AttrPlan) {
		oldPlan, newPlan := d.GetChange(tfconstants.AttrPlan)
		log.Printf("[INFO] Plan change detected for %s %s: %s -> %s", ResourceName, id, oldPlan.(string), newPlan.(string))

		// Verify cluster is stopped
		status := d.Get(tfconstants.AttrStatus).(string)
		if normalizeStatus(strings.ToUpper(status)) != goe2econstants.DBaaSStatusStopped {
			d.Set(tfconstants.AttrPlan, oldPlan.(string))
			return diag.Errorf(ErrorCannotUpgradePlanTemplate, id, status)
		}

		// Get software ID and template ID
		softwareName := d.Get(tfconstants.AttrSoftwareName).(string)
		softwareVersion := d.Get(tfconstants.AttrSoftwareVersion).(string)

		softwareID, err := goe2eClient.MariaDB.GetSoftwareID(ctx, softwareName, softwareVersion)
		if err != nil {
			d.Set(tfconstants.AttrPlan, oldPlan.(string))
			return diag.Errorf(ErrorRetrievingSoftwareIDForUpgrade, softwareName, softwareVersion, id, err)
		}

		templateID, err := goe2eClient.MariaDB.GetTemplateID(ctx, newPlan.(string), softwareID)
		if err != nil {
			d.Set(tfconstants.AttrPlan, oldPlan.(string))
			return diag.Errorf(ErrorRetrievingTemplateIDForUpgrade, softwareName, newPlan.(string), id, err)
		}

		// Upgrade the plan
		if _, err := goe2eClient.MariaDB.UpgradePlan(ctx, id, templateID); err != nil {
			d.Set(tfconstants.AttrPlan, oldPlan.(string))
			return diag.Errorf(ErrorUpgradingPlanTemplate, id, oldPlan.(string), newPlan.(string), err)
		}

		log.Printf("[INFO] Successfully upgraded %s %s to plan %s (template_id=%d)", ResourceName, id, newPlan, templateID)

		// Update template ID in state
		if err := d.Set(tfconstants.AttrTemplateID, templateID); err != nil {
			return diag.FromErr(err)
		}
	}

	// Handle disk expansion
	if d.HasChange(tfconstants.AttrDiskSize) {
		additionalSize := d.Get(tfconstants.AttrDiskSize).(int)

		if additionalSize > 0 {
			log.Printf("[INFO] Disk expansion requested for %s %s: +%d GB", ResourceName, id, additionalSize)

			// Verify cluster is stopped
			status := d.Get(tfconstants.AttrStatus).(string)
			if normalizeStatus(strings.ToUpper(status)) != goe2econstants.DBaaSStatusStopped {
				d.Set(tfconstants.AttrDiskSize, 0)
				return diag.Errorf(ErrorCannotExpandDiskTemplate, id, status)
			}

			// Expand the disk
			if _, err := goe2eClient.MariaDB.ExpandDisk(ctx, id, additionalSize); err != nil {
				d.Set(tfconstants.AttrDiskSize, 0)
				return diag.Errorf(ErrorExpandingDiskTemplate, id, additionalSize, err)
			}

			log.Printf("[INFO] Successfully expanded disk by %d GB for %s %s", additionalSize, ResourceName, id)

			// Reset disk_size to 0 after expansion
			if err := d.Set(tfconstants.AttrDiskSize, 0); err != nil {
				return diag.FromErr(err)
			}
		} else {
			log.Printf("[DEBUG] disk_size is 0 for %s %s, skipping expansion", ResourceName, id)
		}
	}

	// Handle tags (state-only, no API call)
	if d.HasChange(tfconstants.AttrTags) {
		log.Printf("[INFO] %s %s: tags updated (state-only, not sent to API)", ResourceName, id)
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
			tfconstants.AttrSoftwareName: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			tfconstants.AttrSoftwareVersion: {
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
						tfconstants.AttrDatabaseBlockUser: {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						tfconstants.AttrDatabaseBlockPassword: {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						tfconstants.AttrDatabaseBlockName: {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						tfconstants.AttrDatabaseBlockDBaaSNumber: {
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
				Default:  tfconstants.DBaaSDefaultPublicIPEnabled,
			},
			tfconstants.AttrParameterGroupID: {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  tfconstants.DBaaSDefaultParameterGroupID,
			},

			// ============================================
			// OPTIONAL IMMUTABLE FIELDS - SECURITY
			// ============================================
			tfconstants.AttrIsEncryptionEnabled: {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  tfconstants.DBaaSDefaultIsEncryptionEnabled,
			},
			tfconstants.AttrEncryptionPassphrase: {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Default:   tfconstants.DBaaSDefaultEncryptionPassphrase,
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
			tfconstants.AttrPublicIPAttached: {
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
			tfconstants.AttrPort: {
				Type:     schema.TypeString,
				Computed: true,
			},

			// ============================================
			// COMPUTED FIELDS - RESOURCES
			// ============================================
			tfconstants.AttrTotalDiskSize: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			// Note: V0 schema does not have "tags" field
		},
	}
}

// ResourceMariaDBStateUpgradeV0toV1 upgrades the state from V0 to V1
// V1 adds the "tags" field for state-only tag management
func ResourceMariaDBStateUpgradeV0toV1(
	ctx context.Context,
	rawState map[string]interface{},
	meta interface{},
) (map[string]interface{}, error) {
	// Add new V3 fields with defaults
	if _, exists := rawState["tags"]; !exists {
		rawState["tags"] = make(map[string]interface{})
	}

	// Preserve all existing fields - no modifications needed
	log.Printf("[INFO] Upgraded %s resource state from v0 to v1: %s", ResourceName, rawState["id"])

	return rawState, nil
}
