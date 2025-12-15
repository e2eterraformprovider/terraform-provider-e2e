package dbaas_mysql

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

func ResourceMySql() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceMySQLResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: ResourceMySQLStateUpgradeV0toV1,
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
			tfconstants.AttrVersion: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "the MySQL version to use (e.g., 5.6, 5.7, 8.0)",
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
						"dbaas_number": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     tfconstants.DBaaSDefaultDBaaSNumber,
							ForceNew:    true,
							Description: "the DBaaS number (typically 1)",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "name of the database to create",
						},
					},
				},
			},
			tfconstants.AttrPlan: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "the plan name for the MySQL DBaaS instance",
			},

			// ============================================
			// OPTIONAL INPUT FIELDS - IMMUTABLE
			// ============================================
			tfconstants.AttrDBaaSName: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "name of the MySQL DBaaS instance",
			},
			tfconstants.AttrDBLocation: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     tfconstants.DBaaSDefaultDBLocation,
				ForceNew:    true,
				Description: "the location of the MySQL DBaaS instance",
			},
			tfconstants.AttrIsEncryptionEnabled: {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "whether encryption is enabled for the MySQL DBaaS instance",
			},
			tfconstants.AttrEncryptionPassphrase: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "encryption passphrase (required if encryption enabled)",
			},
			tfconstants.AttrGroup: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     tfconstants.DBaaSDefaultGroupName,
				ForceNew:    true,
				Description: "the group name for the instance",
			},

			// ============================================
			// OPTIONAL INPUT FIELDS - MUTABLE
			// ============================================
			tfconstants.AttrVPCs: {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Optional:    true,
				Description: "list of VPC ids to attach to the MySQL DBaaS instance",
			},
			tfconstants.AttrParameterGroupID: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "id of the parameter group to attach",
			},
			tfconstants.AttrPublicIPRequired: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     tfconstants.DBaaSDefaultPublicIPRequired,
				Description: "whether to attach a public IP to the MySQL DBaaS instance",
			},
			tfconstants.AttrSize: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "additional disk size in gigabytes to attach (cumulative across updates)",
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
						goe2econstants.DBaaSStatusSuspended,
						goe2econstants.DBaaSStatusRunning,
						goe2econstants.DBaaSStatusRestarting,
						tfconstants.DBaaSPowerActionStart,
						tfconstants.DBaaSPowerActionStop,
						tfconstants.DBaaSPowerActionRestart,
					},
					false,
				),
				Description: "state of the MySQL DBaaS instance (use 'SUSPENDED' to stop, 'RUNNING' to start, 'RESTARTING' to restart)",
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
			// COMPUTED FIELDS - RESOURCES
			// ============================================
			tfconstants.AttrDisk: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the total disk size of the MySQL DBaaS instance after expansions",
			},

			// ============================================
			// COMPUTED FIELDS - NETWORK
			// ============================================
			tfconstants.AttrPublicIPAddress: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the MySQL instance public IPv4 address (if public IP attached)",
			},
			tfconstants.AttrPrivateIPAddress: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the MySQL instance private IPv4 address",
			},
			tfconstants.AttrPort: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the port number for MySQL service (typically 3306)",
			},
		},

		CreateContext: resourceCreateMySqlDB,
		ReadContext:   resourceReadMySqlDB,
		UpdateContext: resourceUpdateMySqlDB,
		DeleteContext: resourceDeleteMySqlDB,

		Importer: &schema.ResourceImporter{
			State: customImportStateFunc,
		},
	}
}

func resourceCreateMySqlDB(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	softwareName := goe2econstants.DBaaSSoftwareMySQL
	softwareVersion := d.Get(tfconstants.AttrVersion).(string)
	planName := d.Get(tfconstants.AttrPlan).(string)

	// Get software ID using goe2e client
	softwareID, err := goe2eClient.DBaaSMySQL.GetSoftwareID(ctx, softwareName, softwareVersion)
	if err != nil {
		return diag.Errorf(ErrorRetrievingSoftwareIDTemplate, softwareName, softwareVersion, projectID, region, err)
	}

	// Get template ID using goe2e client
	templateID, err := goe2eClient.DBaaSMySQL.GetTemplateID(ctx, planName, softwareID)
	if err != nil {
		return diag.Errorf(ErrorRetrievingTemplateIDTemplate, softwareName, planName, projectID, region, err)
	}

	// Build create request using helper
	createReq, err := buildMySQLCreateRequest(ctx, d, goe2eClient, softwareID, templateID)
	if err != nil {
		return diag.FromErr(err)
	}

	// Create MySQL cluster using goe2e client
	mysql, _, err := goe2eClient.DBaaSMySQL.CreateCluster(ctx, createReq)
	if err != nil {
		return diag.Errorf(tfconstants.ResourceOperationErrorTemplate, tfconstants.OperationCreating, ResourceName, createReq.Name, projectID, region, err)
	}

	// Set resource ID
	d.SetId(strconv.Itoa(mysql.ID))

	// Set computed fields from create response
	if err := d.Set(tfconstants.AttrStatus, normalizeStatus(mysql.Status)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(tfconstants.AttrPublicIPAddress, mysql.MasterNode.PublicIPAddress); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(tfconstants.AttrPrivateIPAddress, mysql.MasterNode.PrivateIPAddress); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(tfconstants.AttrPort, mysql.MasterNode.Port); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Successfully created %s: %s (ID: %d)", ResourceName, mysql.Name, mysql.ID)

	return diags
}

func resourceReadMySqlDB(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	goe2eClient := cfg.Goe2eClient()
	var diags diag.Diagnostics

	id := d.Id()

	// Get MySQL cluster using goe2e client
	mysql, _, err := goe2eClient.DBaaSMySQL.GetCluster(ctx, id)
	if err != nil {
		return diag.Errorf(tfconstants.ResourceOperationByIDErrorTemplate, tfconstants.OperationRetrieving, ResourceName, id, "", "", err)
	}

	// Check if resource was deleted
	if mysql == nil {
		log.Printf("[WARN] %s %s not found, removing from state", ResourceName, id)
		d.SetId("")
		return diags
	}

	// Set status
	if err := d.Set(tfconstants.AttrStatus, normalizeStatus(mysql.Status)); err != nil {
		return diag.FromErr(err)
	}

	// Set encryption status
	if err := d.Set(tfconstants.AttrIsEncryptionEnabled, mysql.IsEncryptionEnabled); err != nil {
		return diag.FromErr(err)
	}

	// Set parameter group ID (with nil check for PGDetail)
	pgID := tfconstants.DBaaSDefaultParameterGroupID
	if mysql.MasterNode.Database.PGDetail.ID != tfconstants.DBaaSDefaultParameterGroupID {
		pgID = mysql.MasterNode.Database.PGDetail.ID
	}
	if err := d.Set(tfconstants.AttrParameterGroupID, pgID); err != nil {
		return diag.FromErr(err)
	}

	// Set disk size
	if err := d.Set(tfconstants.AttrDisk, mysql.MasterNode.Disk); err != nil {
		return diag.FromErr(err)
	}

	// Set network information
	if err := d.Set(tfconstants.AttrPublicIPAddress, mysql.MasterNode.PublicIPAddress); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(tfconstants.AttrPrivateIPAddress, mysql.MasterNode.PrivateIPAddress); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(tfconstants.AttrPort, mysql.MasterNode.Port); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Successfully read %s: %s (ID: %s)", ResourceName, mysql.Name, id)

	return diags
}

func resourceUpdateMySqlDB(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	goe2eClient := cfg.Goe2eClient()
	id := d.Id()

	// Handle status changes (power management)
	if d.HasChange(tfconstants.AttrStatus) {
		newStatus := d.Get(tfconstants.AttrStatus).(string)
		log.Printf("[INFO] Status change detected for %s %s: %s", ResourceName, id, newStatus)

		switch strings.ToUpper(newStatus) {
		case goe2econstants.DBaaSStatusSuspended, goe2econstants.DBaaSStatusStopped:
			if _, err := goe2eClient.DBaaSMySQL.StopCluster(ctx, id); err != nil {
				return diag.Errorf(ErrorStoppingTemplate, id, err)
			}
			log.Printf("[INFO] Successfully stopped %s %s", ResourceName, id)
		case goe2econstants.DBaaSStatusRunning, strings.ToUpper(tfconstants.DBaaSPowerActionStart):
			if _, err := goe2eClient.DBaaSMySQL.StartCluster(ctx, id); err != nil {
				return diag.Errorf(ErrorStartingTemplate, id, err)
			}
			log.Printf("[INFO] Successfully started %s %s", ResourceName, id)
		case goe2econstants.DBaaSStatusRestarting, strings.ToUpper(tfconstants.DBaaSPowerActionRestart):
			if _, err := goe2eClient.DBaaSMySQL.RestartCluster(ctx, id); err != nil {
				return diag.Errorf(ErrorRestartingTemplate, id, err)
			}
			log.Printf("[INFO] Successfully restarted %s %s", ResourceName, id)
		default:
			return diag.Errorf(ErrorUnsupportedStatusTemplate, id, newStatus)
		}
	}

	// Handle VPC changes
	if d.HasChange(tfconstants.AttrVPCs) {
		prevRaw, newRaw := d.GetChange(tfconstants.AttrVPCs)
		prevSet := prevRaw.(*schema.Set)
		newSet := newRaw.(*schema.Set)

		// Calculate added and removed VPCs
		var added, removed []interface{}
		for _, v := range newSet.List() {
			if !prevSet.Contains(v) {
				added = append(added, v)
			}
		}
		for _, v := range prevSet.List() {
			if !newSet.Contains(v) {
				removed = append(removed, v)
			}
		}

		// Attach new VPCs
		if len(added) > 0 {
			log.Printf("[INFO] Attaching VPCs to %s %s: %v", ResourceName, id, added)
			vpcDetails, err := expandVPCList(ctx, goe2eClient, added)
			if err != nil {
				d.Set(tfconstants.AttrVPCs, prevSet)
				return diag.Errorf(ErrorPreparingVPCListTemplate, id, err)
			}
			attachReq := &goe2e.MySQLVPCAttachRequest{
				Action: goe2econstants.ActionAttach,
				VPCs:   vpcDetails,
			}
			if _, err := goe2eClient.DBaaSMySQL.AttachVPC(ctx, id, attachReq); err != nil {
				d.Set(tfconstants.AttrVPCs, prevSet)
				return diag.Errorf(ErrorAttachingVPCTemplate, id, err)
			}
		}

		// Detach removed VPCs
		if len(removed) > 0 {
			log.Printf("[INFO] Detaching VPCs from %s %s: %v", ResourceName, id, removed)
			vpcDetails, err := expandVPCList(ctx, goe2eClient, removed)
			if err != nil {
				d.Set(tfconstants.AttrVPCs, prevSet)
				return diag.Errorf(ErrorPreparingVPCListTemplate, id, err)
			}
			detachReq := &goe2e.MySQLVPCDetachRequest{
				Action: goe2econstants.ActionDetach,
				VPCs:   vpcDetails,
			}
			if _, err := goe2eClient.DBaaSMySQL.DetachVPC(ctx, id, detachReq); err != nil {
				d.Set(tfconstants.AttrVPCs, prevSet)
				return diag.Errorf(ErrorDetachingVPCTemplate, id, err)
			}
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
			if _, err := goe2eClient.DBaaSMySQL.DetachParameterGroup(ctx, id, strconv.Itoa(oldPGID)); err != nil {
				return diag.Errorf(ErrorDetachingParameterGroupTemplate, oldPGID, id, err)
			}
			log.Printf("[INFO] Successfully detached parameter group %d from %s %s", oldPGID, ResourceName, id)
		case newPGID != 0 && newPGID != oldPGID:
			// Attach new parameter group
			if _, err := goe2eClient.DBaaSMySQL.AttachParameterGroup(ctx, id, strconv.Itoa(newPGID)); err != nil {
				return diag.Errorf(ErrorAttachingParameterGroupTemplate, newPGID, id, err)
			}
			log.Printf("[INFO] Successfully attached parameter group %d to %s %s", newPGID, ResourceName, id)
		}
	}

	// Handle public IP changes
	if d.HasChange(tfconstants.AttrPublicIPRequired) {
		newVal := d.Get(tfconstants.AttrPublicIPRequired).(bool)
		log.Printf("[INFO] Public IP change detected for %s %s: %v", ResourceName, id, newVal)

		if newVal {
			if _, err := goe2eClient.DBaaSMySQL.AttachPublicIP(ctx, id); err != nil {
				return diag.Errorf(ErrorAttachingPublicIPTemplate, id, err)
			}
			log.Printf("[INFO] Successfully attached public IP to %s %s", ResourceName, id)
		} else {
			if _, err := goe2eClient.DBaaSMySQL.DetachPublicIP(ctx, id); err != nil {
				return diag.Errorf(ErrorDetachingPublicIPTemplate, id, err)
			}
			log.Printf("[INFO] Successfully detached public IP from %s %s", ResourceName, id)
		}
	}

	// Handle plan upgrade
	if d.HasChange(tfconstants.AttrPlan) {
		oldPlan, newPlan := d.GetChange(tfconstants.AttrPlan)
		log.Printf("[INFO] Plan change detected for %s %s: %s -> %s", ResourceName, id, oldPlan.(string), newPlan.(string))

		// Verify cluster is suspended
		status := d.Get(tfconstants.AttrStatus).(string)
		if strings.ToUpper(status) != goe2econstants.DBaaSStatusSuspended && strings.ToUpper(status) != goe2econstants.DBaaSStatusStopped {
			d.Set(tfconstants.AttrPlan, oldPlan.(string))
			return diag.Errorf(ErrorCannotUpgradePlanTemplate, id, status)
		}

		// Get version for software ID lookup
		version := d.Get(tfconstants.AttrVersion).(string)

		// Get software ID and template ID
		softwareID, err := goe2eClient.DBaaSMySQL.GetSoftwareID(ctx, goe2econstants.DBaaSSoftwareMySQL, version)
		if err != nil {
			d.Set(tfconstants.AttrPlan, oldPlan.(string))
			return diag.Errorf(ErrorRetrievingSoftwareIDForUpgrade, version, id, err)
		}

		templateID, err := goe2eClient.DBaaSMySQL.GetTemplateID(ctx, newPlan.(string), softwareID)
		if err != nil {
			d.Set(tfconstants.AttrPlan, oldPlan.(string))
			return diag.Errorf(ErrorRetrievingTemplateIDForUpgrade, newPlan.(string), id, err)
		}

		// Upgrade the plan
		upgradeReq := &goe2e.MySQLPlanUpgradeRequest{
			TemplateID: templateID,
		}
		if _, err := goe2eClient.DBaaSMySQL.UpgradePlan(ctx, id, upgradeReq); err != nil {
			d.Set(tfconstants.AttrPlan, oldPlan.(string))
			return diag.Errorf(ErrorUpgradingPlanTemplate, id, oldPlan.(string), newPlan.(string), err)
		}

		log.Printf("[INFO] Successfully upgraded %s %s to plan %s (template_id=%d)", ResourceName, id, newPlan, templateID)
	}

	// Handle disk expansion
	if d.HasChange(tfconstants.AttrSize) {
		additionalSize := d.Get(tfconstants.AttrSize).(int)

		if additionalSize > 0 {
			log.Printf("[INFO] Disk expansion requested for %s %s: +%d GB", ResourceName, id, additionalSize)

			// Expand the disk
			expandReq := &goe2e.DiskExpansionRequest{
				Size: additionalSize,
			}
			if _, err := goe2eClient.DBaaSMySQL.ExpandDisk(ctx, id, expandReq); err != nil {
				d.Set(tfconstants.AttrSize, 0)
				return diag.Errorf(ErrorExpandingDiskTemplate, id, additionalSize, err)
			}

			log.Printf("[INFO] Successfully expanded disk by %d GB for %s %s", additionalSize, ResourceName, id)

			// Reset size to 0 after expansion (cumulative behavior)
			if err := d.Set(tfconstants.AttrSize, 0); err != nil {
				return diag.FromErr(err)
			}
		} else {
			log.Printf("[DEBUG] size is 0 for %s %s, skipping expansion", ResourceName, id)
		}
	}

	// Handle tags (state-only, no API call)
	if d.HasChange(tfconstants.AttrTags) {
		log.Printf("[INFO] %s %s: tags updated (state-only, not sent to API)", ResourceName, id)
	}

	// Refresh state by reading the resource
	return resourceReadMySqlDB(ctx, d, m)
}

func resourceDeleteMySqlDB(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	goe2eClient := cfg.Goe2eClient()
	var diags diag.Diagnostics

	id := d.Id()

	log.Printf("[INFO] Deleting %s: %s", ResourceName, id)

	// Delete MySQL cluster using goe2e client
	_, err := goe2eClient.DBaaSMySQL.DeleteCluster(ctx, id)
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

func CustomImportStateFunc(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), tfconstants.DBaaSImportIDSeparator)
	if len(parts) != 2 {
		return nil, fmt.Errorf(ImportIDInvalidFormatTemplate, tfconstants.DBaaSImportIDFormatDescription, d.Id())
	}

	projectID := parts[0]
	dbaasID := parts[1]

	if err := d.Set(tfconstants.AttrProjectID, projectID); err != nil {
		return nil, err
	}
	d.SetId(dbaasID)

	return []*schema.ResourceData{d}, nil
}

// ============================================
// STATE MIGRATION: V0 → V1
// ============================================

// resourceMySQLResourceV0 returns the V0 schema definition (before tags were added)
func resourceMySQLResourceV0() *schema.Resource {
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
			tfconstants.AttrVersion: {
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
						"dbaas_number": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  tfconstants.DBaaSDefaultDBaaSNumber,
							ForceNew: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			tfconstants.AttrPlan: {
				Type:     schema.TypeString,
				Required: true,
			},

			// ============================================
			// OPTIONAL INPUT FIELDS - IMMUTABLE
			// ============================================
			tfconstants.AttrDBaaSName: {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			tfconstants.AttrDBLocation: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  tfconstants.DBaaSDefaultDBLocation,
				ForceNew: true,
			},
			tfconstants.AttrIsEncryptionEnabled: {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			tfconstants.AttrEncryptionPassphrase: {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			tfconstants.AttrGroup: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  tfconstants.DBaaSDefaultGroupName,
				ForceNew: true,
			},

			// ============================================
			// OPTIONAL INPUT FIELDS - MUTABLE
			// ============================================
			tfconstants.AttrVPCs: {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Optional: true,
			},
			tfconstants.AttrParameterGroupID: {
				Type:     schema.TypeInt,
				Optional: true,
			},
			tfconstants.AttrPublicIPRequired: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  tfconstants.DBaaSDefaultPublicIPRequired,
			},
			tfconstants.AttrSize: {
				Type:     schema.TypeInt,
				Optional: true,
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
						goe2econstants.DBaaSStatusSuspended,
						goe2econstants.DBaaSStatusRunning,
						goe2econstants.DBaaSStatusRestarting,
						tfconstants.DBaaSPowerActionStart,
						tfconstants.DBaaSPowerActionStop,
						tfconstants.DBaaSPowerActionRestart,
					},
					false,
				),
			},

			// ============================================
			// COMPUTED FIELDS - RESOURCES
			// ============================================
			tfconstants.AttrDisk: {
				Type:     schema.TypeString,
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
			// Note: V0 schema does not have "tags" field
		},
	}
}

// ResourceMySQLStateUpgradeV0toV1 upgrades the state from V0 to V1
// V1 adds the "tags" field for state-only tag management
func ResourceMySQLStateUpgradeV0toV1(
	ctx context.Context,
	rawState map[string]interface{},
	meta interface{},
) (map[string]interface{}, error) {
	// Add new V1 fields with defaults and ensure tags is a map
	if tagsVal, exists := rawState[tfconstants.AttrTags]; exists {
		// If tags exists but is not a map, convert it to an empty map
		if _, ok := tagsVal.(map[string]interface{}); !ok {
			rawState[tfconstants.AttrTags] = make(map[string]interface{})
		}
	} else {
		// If tags doesn't exist, create an empty map
		rawState[tfconstants.AttrTags] = make(map[string]interface{})
	}

	// Preserve all existing fields - no modifications needed
	log.Printf("[INFO] Upgraded %s resource state from v0 to v1: %s", ResourceName, rawState[tfconstants.AttrID])

	return rawState, nil
}
