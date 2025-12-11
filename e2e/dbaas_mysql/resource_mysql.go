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
				Upgrade: resourceMySQLStateUpgradeV0toV1,
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
							Default:     1,
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
			"dbaas_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "name of the MySQL DBaaS instance",
			},
			"db_location": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Delhi",
				ForceNew:    true,
				Description: "the location of the MySQL DBaaS instance",
			},
			"is_encryption_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "whether encryption is enabled for the MySQL DBaaS instance",
			},
			"encryption_passphrase": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "encryption passphrase (required if encryption enabled)",
			},
			tfconstants.AttrGroup: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Default",
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
				Default:     true,
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
					[]string{"STOPPED", "SUSPENDED", "RUNNING", "RESTARTING", "start", "stop", "restart"},
					false,
				),
				Description: "state of the MySQL DBaaS instance (use 'SUSPENDED' to stop, 'RUNNING' to start, 'RESTARTING' to restart)",
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
			"port": {
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

	softwareName := "MySQL"
	softwareVersion := d.Get(tfconstants.AttrVersion).(string)
	planName := d.Get(tfconstants.AttrPlan).(string)

	// Get software ID using goe2e client
	softwareID, err := goe2eClient.DBaaSMySQL.GetSoftwareID(ctx, softwareName, softwareVersion)
	if err != nil {
		return diag.Errorf("error retrieving %s software ID for version (%s) in project (%s), region (%s): %s", softwareName, softwareVersion, projectID, region, err)
	}

	// Get template ID using goe2e client
	templateID, err := goe2eClient.DBaaSMySQL.GetTemplateID(ctx, planName, softwareID)
	if err != nil {
		return diag.Errorf("error retrieving %s template ID for plan (%s) in project (%s), region (%s): %s", softwareName, planName, projectID, region, err)
	}

	// Build create request using helper
	createReq, err := buildMySQLCreateRequest(ctx, d, goe2eClient, softwareID, templateID)
	if err != nil {
		return diag.FromErr(err)
	}

	// Create MySQL cluster using goe2e client
	mysql, _, err := goe2eClient.DBaaSMySQL.CreateCluster(ctx, createReq)
	if err != nil {
		return diag.Errorf("error creating MySQL DBaaS (name: %s) in project (%s), region (%s): %s", createReq.Name, projectID, region, err)
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
	if err := d.Set("port", mysql.MasterNode.Port); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Successfully created MySQL cluster: %s (ID: %d)", mysql.Name, mysql.ID)

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
		return diag.Errorf("error retrieving MySQL DBaaS (ID: %s): %s", id, err)
	}

	// Check if resource was deleted
	if mysql == nil {
		log.Printf("[WARN] MySQL cluster %s not found, removing from state", id)
		d.SetId("")
		return diags
	}

	// Set status
	if err := d.Set(tfconstants.AttrStatus, normalizeStatus(mysql.Status)); err != nil {
		return diag.FromErr(err)
	}

	// Set encryption status
	if err := d.Set("is_encryption_enabled", mysql.IsEncryptionEnabled); err != nil {
		return diag.FromErr(err)
	}

	// Set parameter group ID (with nil check for PGDetail)
	pgID := 0
	if mysql.MasterNode.Database.PGDetail.ID != 0 {
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
	if err := d.Set("port", mysql.MasterNode.Port); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Successfully read MySQL cluster: %s (ID: %s)", mysql.Name, id)

	return diags
}

func resourceUpdateMySqlDB(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	goe2eClient := cfg.Goe2eClient()
	id := d.Id()

	// Handle status changes (power management)
	if d.HasChange("status") {
		newStatus := d.Get("status").(string)
		log.Printf("[INFO] Status change detected for MySQL cluster %s: %s", id, newStatus)

		switch strings.ToUpper(newStatus) {
		case "SUSPENDED", "STOPPED":
			if _, err := goe2eClient.DBaaSMySQL.StopCluster(ctx, id); err != nil {
				return diag.Errorf("error stopping MySQL DBaaS (ID: %s): %s", id, err)
			}
			log.Printf("[INFO] Successfully stopped MySQL cluster %s", id)
		case "RUNNING", "START":
			if _, err := goe2eClient.DBaaSMySQL.StartCluster(ctx, id); err != nil {
				return diag.Errorf("error starting MySQL DBaaS (ID: %s): %s", id, err)
			}
			log.Printf("[INFO] Successfully started MySQL cluster %s", id)
		case "RESTARTING", "RESTART":
			if _, err := goe2eClient.DBaaSMySQL.RestartCluster(ctx, id); err != nil {
				return diag.Errorf("error restarting MySQL DBaaS (ID: %s): %s", id, err)
			}
			log.Printf("[INFO] Successfully restarted MySQL cluster %s", id)
		default:
			return diag.Errorf("error updating MySQL DBaaS (ID: %s): unsupported status value: %s. Must be one of: SUSPENDED, STOPPED, RUNNING, START, RESTARTING, RESTART", id, newStatus)
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
			log.Printf("[INFO] Attaching VPCs to MySQL cluster %s: %v", id, added)
			vpcDetails, err := expandVPCList(ctx, goe2eClient, added)
			if err != nil {
				d.Set(tfconstants.AttrVPCs, prevSet)
				return diag.Errorf("error preparing VPC list for MySQL DBaaS (ID: %s): %s", id, err)
			}
			attachReq := &goe2e.MySQLVPCAttachRequest{
				Action: "attach",
				VPCs:   vpcDetails,
			}
			if _, err := goe2eClient.DBaaSMySQL.AttachVPC(ctx, id, attachReq); err != nil {
				d.Set(tfconstants.AttrVPCs, prevSet)
				return diag.Errorf("error attaching VPC to MySQL DBaaS (ID: %s): %s", id, err)
			}
		}

		// Detach removed VPCs
		if len(removed) > 0 {
			log.Printf("[INFO] Detaching VPCs from MySQL cluster %s: %v", id, removed)
			vpcDetails, err := expandVPCList(ctx, goe2eClient, removed)
			if err != nil {
				d.Set(tfconstants.AttrVPCs, prevSet)
				return diag.Errorf("error preparing VPC list for MySQL DBaaS (ID: %s): %s", id, err)
			}
			detachReq := &goe2e.MySQLVPCDetachRequest{
				Action: "detach",
				VPCs:   vpcDetails,
			}
			if _, err := goe2eClient.DBaaSMySQL.DetachVPC(ctx, id, detachReq); err != nil {
				d.Set(tfconstants.AttrVPCs, prevSet)
				return diag.Errorf("error detaching VPC from MySQL DBaaS (ID: %s): %s", id, err)
			}
		}
	}

	// Handle parameter group changes
	if d.HasChange(tfconstants.AttrParameterGroupID) {
		oldRaw, newRaw := d.GetChange(tfconstants.AttrParameterGroupID)
		oldPGID := oldRaw.(int)
		newPGID := newRaw.(int)

		log.Printf("[INFO] Parameter group change detected for MySQL cluster %s: %d -> %d", id, oldPGID, newPGID)

		switch {
		case oldPGID != 0 && newPGID == 0:
			// Detach parameter group
			if _, err := goe2eClient.DBaaSMySQL.DetachParameterGroup(ctx, id, strconv.Itoa(oldPGID)); err != nil {
				return diag.Errorf("error detaching parameter group (ID: %d) from MySQL DBaaS (ID: %s): %s", oldPGID, id, err)
			}
			log.Printf("[INFO] Successfully detached parameter group %d from MySQL cluster %s", oldPGID, id)
		case newPGID != 0 && newPGID != oldPGID:
			// Attach new parameter group
			if _, err := goe2eClient.DBaaSMySQL.AttachParameterGroup(ctx, id, strconv.Itoa(newPGID)); err != nil {
				return diag.Errorf("error attaching parameter group (ID: %d) to MySQL DBaaS (ID: %s): %s", newPGID, id, err)
			}
			log.Printf("[INFO] Successfully attached parameter group %d to MySQL cluster %s", newPGID, id)
		}
	}

	// Handle public IP changes
	if d.HasChange(tfconstants.AttrPublicIPRequired) {
		newVal := d.Get(tfconstants.AttrPublicIPRequired).(bool)
		log.Printf("[INFO] Public IP change detected for MySQL cluster %s: %v", id, newVal)

		if newVal {
			if _, err := goe2eClient.DBaaSMySQL.AttachPublicIP(ctx, id); err != nil {
				return diag.Errorf("error attaching public IP to MySQL DBaaS (ID: %s): %s", id, err)
			}
			log.Printf("[INFO] Successfully attached public IP to MySQL cluster %s", id)
		} else {
			if _, err := goe2eClient.DBaaSMySQL.DetachPublicIP(ctx, id); err != nil {
				return diag.Errorf("error detaching public IP from MySQL DBaaS (ID: %s): %s", id, err)
			}
			log.Printf("[INFO] Successfully detached public IP from MySQL cluster %s", id)
		}
	}

	// Handle plan upgrade
	if d.HasChange("plan") {
		oldPlan, newPlan := d.GetChange("plan")
		log.Printf("[INFO] Plan change detected for MySQL cluster %s: %s -> %s", id, oldPlan.(string), newPlan.(string))

		// Verify cluster is suspended
		status := d.Get("status").(string)
		if strings.ToUpper(status) != "SUSPENDED" && strings.ToUpper(status) != "STOPPED" {
			d.Set("plan", oldPlan.(string))
			return diag.Errorf("cannot upgrade plan for MySQL DBaaS (ID: %s): database must be in SUSPENDED/STOPPED state (current state: %s). Please stop the instance first", id, status)
		}

		// Get version for software ID lookup
		version := d.Get(tfconstants.AttrVersion).(string)

		// Get software ID and template ID
		softwareID, err := goe2eClient.DBaaSMySQL.GetSoftwareID(ctx, "MySQL", version)
		if err != nil {
			d.Set("plan", oldPlan.(string))
			return diag.Errorf("error retrieving MySQL software ID for version (%s) while upgrading plan for DBaaS (ID: %s): %s", version, id, err)
		}

		templateID, err := goe2eClient.DBaaSMySQL.GetTemplateID(ctx, newPlan.(string), softwareID)
		if err != nil {
			d.Set("plan", oldPlan.(string))
			return diag.Errorf("error retrieving MySQL template ID for plan (%s) while upgrading DBaaS (ID: %s): %s", newPlan.(string), id, err)
		}

		// Upgrade the plan
		upgradeReq := &goe2e.MySQLPlanUpgradeRequest{
			TemplateID: templateID,
		}
		if _, err := goe2eClient.DBaaSMySQL.UpgradePlan(ctx, id, upgradeReq); err != nil {
			d.Set("plan", oldPlan.(string))
			return diag.Errorf("error upgrading MySQL DBaaS (ID: %s) plan from (%s) to (%s): %s", id, oldPlan.(string), newPlan.(string), err)
		}

		log.Printf("[INFO] Successfully upgraded MySQL cluster %s to plan %s (template_id=%d)", id, newPlan, templateID)
	}

	// Handle disk expansion
	if d.HasChange(tfconstants.AttrSize) {
		additionalSize := d.Get(tfconstants.AttrSize).(int)

		if additionalSize > 0 {
			log.Printf("[INFO] Disk expansion requested for MySQL cluster %s: +%d GB", id, additionalSize)

			// Expand the disk
			expandReq := &goe2e.DiskExpansionRequest{
				Size: additionalSize,
			}
			if _, err := goe2eClient.DBaaSMySQL.ExpandDisk(ctx, id, expandReq); err != nil {
				d.Set(tfconstants.AttrSize, 0)
				return diag.Errorf("error expanding MySQL DBaaS (ID: %s) disk by %d GB: %s", id, additionalSize, err)
			}

			log.Printf("[INFO] Successfully expanded disk by %d GB for MySQL cluster %s", additionalSize, id)

			// Reset size to 0 after expansion (cumulative behavior)
			if err := d.Set(tfconstants.AttrSize, 0); err != nil {
				return diag.FromErr(err)
			}
		} else {
			log.Printf("[DEBUG] size is 0 for MySQL cluster %s, skipping expansion", id)
		}
	}

	// Handle tags (state-only, no API call)
	if d.HasChange("tags") {
		log.Printf("[INFO] MySQL cluster %s: tags updated (state-only, not sent to API)", id)
	}

	// Refresh state by reading the resource
	return resourceReadMySqlDB(ctx, d, m)
}

func resourceDeleteMySqlDB(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	goe2eClient := cfg.Goe2eClient()
	var diags diag.Diagnostics

	id := d.Id()

	log.Printf("[INFO] Deleting MySQL cluster: %s", id)

	// Delete MySQL cluster using goe2e client
	_, err := goe2eClient.DBaaSMySQL.DeleteCluster(ctx, id)
	if err != nil {
		// Check if already deleted (404 or "not found" error)
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "404") {
			log.Printf("[WARN] MySQL cluster %s already deleted", id)
			d.SetId("")
			return diags
		}
		return diag.Errorf("error deleting MySQL DBaaS (ID: %s): %s", id, err)
	}

	log.Printf("[INFO] Successfully deleted MySQL cluster: %s", id)
	d.SetId("")
	return diags
}

func CustomImportStateFunc(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid ID format: expected project_id:dbaas_id")
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
							Default:  1,
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
			"dbaas_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"db_location": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Delhi",
				ForceNew: true,
			},
			"is_encryption_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"encryption_passphrase": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			tfconstants.AttrGroup: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Default",
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
				Default:  true,
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
					[]string{"STOPPED", "SUSPENDED", "RUNNING", "RESTARTING", "start", "stop", "restart"},
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
			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// Note: V0 schema does not have "tags" field
		},
	}
}

// resourceMySQLStateUpgradeV0toV1 upgrades the state from V0 to V1
// V1 adds the "tags" field for state-only tag management
func resourceMySQLStateUpgradeV0toV1(
	ctx context.Context,
	rawState map[string]interface{},
	meta interface{},
) (map[string]interface{}, error) {
	// Add new V1 fields with defaults
	if _, exists := rawState["tags"]; !exists {
		rawState["tags"] = make(map[string]interface{})
	}

	// Preserve all existing fields - no modifications needed
	log.Printf("[INFO] Upgraded MySQL resource state from v0 to v1: %s", rawState["id"])

	return rawState, nil
}
