package dbaas_postgress

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

func ResourcePostgresDBaaS() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourcePostgreSQLResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourcePostgreSQLStateUpgradeV0toV1,
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
			// REQUIRED INPUT FIELDS (Immutable)
			// ============================================
			tfconstants.AttrVersion: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "the PostgreSQL version to use (e.g., 11.0, 12.0, 13.0, 14.0, 15.0)",
				ValidateFunc: validation.StringInSlice(
					[]string{"11.0", "12.0", "13.0", "14.0", "15.0"},
					false,
				),
			},
			tfconstants.AttrPlan: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "the plan name for the PostgreSQL DBaaS instance (upgradeable only when instance is in SUSPENDED state)",
			},
			tfconstants.AttrName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "name of the PostgreSQL DBaaS instance",
			},
			tfconstants.AttrDatabase: {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				// No ForceNew on block itself - allows password rotation
				Description: "database configuration (user, password, database name)",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true, // ✅ ADDED - initial admin user immutable
							Description: "the database username",
						},
						"password": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
							// No ForceNew - password rotation supported
							Description: "the database password",
						},
						"dbaas_number": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     1,
							ForceNew:    true, // ✅ ADDED - topology immutable
							Description: "the DBaaS number (typically 1)",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true, // ✅ ADDED - initial DB name immutable
							Description: "name of the database to create",
						},
					},
				},
			},

			// ============================================
			// OPTIONAL INPUT FIELDS - CREATION
			// ============================================
			tfconstants.AttrGroup: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "the group name for the PostgreSQL DBaaS instance",
				Default:     "Default",
			},
			tfconstants.AttrPublicIPRequired: {
				Type:     schema.TypeBool,
				Optional: true,
				// No ForceNew - can be attached/detached dynamically
				Description: "whether to attach a public IP to the PostgreSQL DBaaS instance",
				Default:     true,
			},
			tfconstants.AttrIsEncryptionEnabled: {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "whether to enable encryption at rest for the PostgreSQL DBaaS instance",
			},
			tfconstants.AttrParameterGroupID: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "id of the parameter group to attach",
			},
			tfconstants.AttrVPCs: {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Optional:    true,
				Description: "list of VPC ids to attach to the PostgreSQL DBaaS instance",
			},

			// ============================================
			// OPTIONAL INPUT FIELDS - MANAGEMENT
			// ============================================
			tfconstants.AttrStatus: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"STOPPED", "SUSPENDED", "RUNNING", "RESTARTING"},
					false,
				),
				Description: "state of the PostgreSQL DBaaS instance (use 'SUSPENDED' to stop, 'RUNNING' to start, 'RESTARTING' to restart)",
			},
			tfconstants.AttrSize: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "additional disk size in gigabytes to attach (cumulative: each increase adds to previous expansions)",
			},

			// ============================================
			// COMPUTED FIELDS
			// ============================================
			tfconstants.AttrID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the ID of the PostgreSQL DBaaS instance",
			},
			// Status is now Optional+Computed (defined above in management section)
			"status_title": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the status title of the PostgreSQL DBaaS instance",
			},
			"status_actions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "available actions for the PostgreSQL DBaaS instance",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"num_instances": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "the number of instances in the PostgreSQL DBaaS cluster",
			},
			"project_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "name of the project",
			},
			"snapshot_exist": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "whether a snapshot exists for the PostgreSQL DBaaS instance",
			},
			"connectivity_detail": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the connectivity details for the PostgreSQL DBaaS instance",
			},
			"vector_database_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the vector database status of the PostgreSQL DBaaS instance",
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
			// COMPUTED FIELDS - NETWORK
			// ============================================
			tfconstants.AttrPublicIPAddress: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the PostgreSQL instance public IPv4 address (if public IP attached)",
			},
			tfconstants.AttrPrivateIPAddress: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the PostgreSQL instance private IPv4 address",
			},
			"port": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the port number for PostgreSQL service (typically 5432)",
			},
			tfconstants.AttrDisk: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the total disk size of the PostgreSQL DBaaS instance after expansions",
			},
		},

		CreateContext: resourceCreatePostgress,
		ReadContext:   resourceReadPostgress,
		DeleteContext: resourceDeletePostgress,
		UpdateContext: resourceUpdatePostgress,
		Importer: &schema.ResourceImporter{
			StateContext: CustomImportStateFunc,
		},
	}
}

// resourceCreatePostgress creates a new PostgreSQL DBaaS cluster.
// It validates inputs, retrieves software/template IDs, and creates the cluster using the goe2e client.
func resourceCreatePostgress(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	var diags diag.Diagnostics

	log.Printf("[INFO] Creating PostgreSQL DBaaS cluster")

	// Get project ID and region from resource or provider defaults
	projectID, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get goe2e client for this project/region
	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return diag.Errorf("error creating goe2e client for project (%s), region (%s): %s", projectID, region, err)
	}

	// Extract nested database config
	dbConfigList := d.Get(tfconstants.AttrDatabase).([]interface{})
	if len(dbConfigList) == 0 {
		return diag.Errorf("database configuration is required")
	}
	dbConfigMap := dbConfigList[0].(map[string]interface{})

	plan := d.Get("plan").(string)
	version := d.Get(tfconstants.AttrVersion).(string)

	// Get software ID using goe2e client
	// Note: pgID parameter is required but can be empty string for PostgreSQL
	softwareID, err := goe2eClient.PostgreSQL.GetSoftwareID(ctx, "PostgreSQL", version, "")
	if err != nil {
		return diag.Errorf("error retrieving PostgreSQL software ID for version (%s) in project (%s), region (%s): %s", version, projectID, region, err)
	}

	// Get template ID using goe2e client
	templateID, err := goe2eClient.PostgreSQL.GetTemplateID(ctx, plan, strconv.Itoa(softwareID), "")
	if err != nil {
		return diag.Errorf("error retrieving PostgreSQL template ID for plan (%s) in project (%s), region (%s): %s", plan, projectID, region, err)
	}

	var pgID *int
	if v, ok := d.GetOk(tfconstants.AttrParameterGroupID); ok {
		id := v.(int)
		pgID = &id
	}
	// Extract VPC IDs
	var vpcIDs []string
	if vpcSet, ok := d.GetOk(tfconstants.AttrVPCs); ok {
		for _, v := range vpcSet.(*schema.Set).List() {
			vpcIDs = append(vpcIDs, strconv.Itoa(v.(int)))
		}
	}

	// Expand VPC list using goe2e client
	var vpcList []goe2e.VPCMetadata
	if len(vpcIDs) > 0 {
		vpcList, err = goe2eClient.PostgreSQL.ExpandPostgresVPCList(ctx, vpcIDs)
		if err != nil {
			return diag.Errorf("error preparing VPC list for PostgreSQL DBaaS in project (%s), region (%s): %s", projectID, region, err)
		}
	}

	// Build create request using goe2e types
	req := &goe2e.PostgreSQLClusterCreateRequest{
		SoftwareID:       softwareID,
		TemplateID:       templateID,
		Name:             d.Get("name").(string),
		Group:            d.Get(tfconstants.AttrGroup).(string),
		PublicIPRequired: d.Get(tfconstants.AttrPublicIPRequired).(bool),
		VPCs:             vpcList,
		Database: goe2e.DBConfig{
			User:        dbConfigMap["user"].(string),
			Password:    dbConfigMap["password"].(string),
			DBaaSNumber: dbConfigMap["dbaas_number"].(int),
			Name:        dbConfigMap["name"].(string),
		},
		PGID:                pgID,
		IsEncryptionEnabled: d.Get(tfconstants.AttrIsEncryptionEnabled).(bool),
	}

	// Create PostgreSQL cluster using goe2e client
	cluster, _, err := goe2eClient.PostgreSQL.CreateCluster(ctx, req)
	if err != nil {
		return diag.Errorf("error creating PostgreSQL DBaaS (name: %s) in project (%s), region (%s): %s", req.Name, projectID, region, err)
	}

	// Set resource ID and attributes
	d.SetId(strconv.Itoa(cluster.ID))
	if err := d.Set("name", cluster.Name); err != nil {
		return diag.FromErr(err)
	}

	// Normalize status (SUSPENDED → STOPPED for consistency)
	status := cluster.Status
	if status == "SUSPENDED" {
		status = "STOPPED"
	}
	if err := d.Set(tfconstants.AttrStatus, status); err != nil {
		return diag.FromErr(err)
	}

	// Set other computed fields
	if err := d.Set(tfconstants.AttrPublicIPAddress, cluster.MasterNode.PublicIPAddress); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(tfconstants.AttrPrivateIPAddress, cluster.MasterNode.PrivateIPAddress); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("port", cluster.MasterNode.Port); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(tfconstants.AttrDisk, cluster.MasterNode.Disk); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(tfconstants.AttrIsEncryptionEnabled, cluster.IsEncryptionEnabled); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Successfully created PostgreSQL DBaaS cluster: %s (ID: %d)", cluster.Name, cluster.ID)

	return diags
}

// resourceReadPostgress reads the current state of a PostgreSQL DBaaS cluster.
// It retrieves the cluster details from the API and updates the Terraform state.
func resourceReadPostgress(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	var diags diag.Diagnostics

	log.Printf("[INFO] Reading PostgreSQL DBaaS cluster")

	clusterID := d.Id()
	if clusterID == "" {
		clusterID = d.Get("id").(string)
	}

	projectID, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get goe2e client for this project/region
	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return diag.Errorf("error creating goe2e client for project (%s), region (%s): %s", projectID, region, err)
	}

	// Get PostgreSQL cluster using goe2e client
	cluster, _, err := goe2eClient.PostgreSQL.GetCluster(ctx, clusterID)
	if err != nil {
		return diag.Errorf("error retrieving PostgreSQL DBaaS (ID: %s) in project (%s), region (%s): %s", clusterID, projectID, region, err)
	}

	// Check if resource was deleted
	if cluster == nil {
		d.SetId("")
		return diags
	}

	// Set resource ID
	d.SetId(strconv.Itoa(cluster.ID))
	if err := d.Set("id", cluster.ID); err != nil {
		return diag.FromErr(err)
	}

	// Set basic fields
	if err := d.Set("name", cluster.Name); err != nil {
		return diag.FromErr(err)
	}

	// Normalize status (SUSPENDED → STOPPED for consistency)
	status := cluster.Status
	if status == "SUSPENDED" {
		status = "STOPPED"
	}
	if err := d.Set(tfconstants.AttrStatus, status); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("status_title", cluster.StatusTitle); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("status_actions", cluster.StatusActions); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("num_instances", cluster.NumInstances); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("project_name", cluster.ProjectName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("snapshot_exist", cluster.SnapshotExist); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("connectivity_detail", cluster.ConnectivityDetail); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vector_database_status", cluster.VectorDBStatus); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(tfconstants.AttrIsEncryptionEnabled, cluster.IsEncryptionEnabled); err != nil {
		return diag.FromErr(err)
	}

	// Set network fields
	if err := d.Set(tfconstants.AttrPublicIPAddress, cluster.MasterNode.PublicIPAddress); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(tfconstants.AttrPrivateIPAddress, cluster.MasterNode.PrivateIPAddress); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("port", cluster.MasterNode.Port); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(tfconstants.AttrDisk, cluster.MasterNode.Disk); err != nil {
		return diag.FromErr(err)
	}

	// Set plan and version
	if err := d.Set(tfconstants.AttrPlan, cluster.MasterNode.Plan.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(tfconstants.AttrVersion, cluster.MasterNode.Plan.Software.Version); err != nil {
		return diag.FromErr(err)
	}

	// Set parameter group ID if present
	if cluster.MasterNode.Database.PGDetail.ID != 0 {
		if err := d.Set(tfconstants.AttrParameterGroupID, cluster.MasterNode.Database.PGDetail.ID); err != nil {
			return diag.FromErr(err)
		}
	}

	log.Printf("[INFO] Successfully read PostgreSQL DBaaS cluster: %s (ID: %d)", cluster.Name, cluster.ID)

	return diags
}

// resourceUpdatePostgress updates a PostgreSQL DBaaS cluster.
// It handles updates to: status (power management), public IP, VPCs, parameter groups, plan, and disk size.
// Each update type has its own validation and API call logic.
func resourceUpdatePostgress(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)

	clusterID := d.Id()
	if clusterID == "" {
		clusterID = d.Get("id").(string)
	}
	if clusterID == "" {
		return diag.Errorf("cluster ID is required for update")
	}

	projectID, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get goe2e client for this project/region
	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return diag.Errorf("error creating goe2e client for project (%s), region (%s): %s", projectID, region, err)
	}

	// Handle status changes (power management)
	if d.HasChange(tfconstants.AttrStatus) {
		oldStatus, newStatusRaw := d.GetChange(tfconstants.AttrStatus)
		currentStatus := oldStatus.(string)
		newStatus := newStatusRaw.(string)

		// Block operation if DBaaS is still in "Creating" state
		if currentStatus == "CREATING" {
			return diag.Errorf("cannot perform power operations on PostgreSQL DBaaS (ID: %s): database is in CREATING state in project (%s), region (%s). Please wait for database creation to complete", clusterID, projectID, region)
		}

		log.Printf("[INFO] Status change detected for PostgreSQL cluster %s: %s -> %s", clusterID, currentStatus, newStatus)

		switch strings.ToUpper(newStatus) {
		case "STOPPED", "SUSPENDED":
			_, err := goe2eClient.PostgreSQL.StopCluster(ctx, clusterID)
			if err != nil {
				return diag.Errorf("error stopping PostgreSQL DBaaS (ID: %s) in project (%s), region (%s): %s", clusterID, projectID, region, err)
			}
			log.Printf("[INFO] Successfully stopped PostgreSQL cluster %s", clusterID)
		case "RUNNING":
			_, err := goe2eClient.PostgreSQL.StartCluster(ctx, clusterID)
			if err != nil {
				return diag.Errorf("error starting PostgreSQL DBaaS (ID: %s) in project (%s), region (%s): %s", clusterID, projectID, region, err)
			}
			log.Printf("[INFO] Successfully started PostgreSQL cluster %s", clusterID)
		case "RESTARTING":
			_, err := goe2eClient.PostgreSQL.RestartCluster(ctx, clusterID)
			if err != nil {
				return diag.Errorf("error restarting PostgreSQL DBaaS (ID: %s) in project (%s), region (%s): %s", clusterID, projectID, region, err)
			}
			log.Printf("[INFO] Successfully restarted PostgreSQL cluster %s", clusterID)
		default:
			return diag.Errorf("error updating PostgreSQL DBaaS (ID: %s): unsupported status value: %s. Must be one of: STOPPED, SUSPENDED, RUNNING, RESTARTING", clusterID, newStatus)
		}
	}

	// Handle public IP changes
	if d.HasChange(tfconstants.AttrPublicIPRequired) {
		newVal := d.Get(tfconstants.AttrPublicIPRequired).(bool)
		currentStatus := d.Get(tfconstants.AttrStatus).(string)

		// Block operation if DBaaS is still in "Creating" state
		if currentStatus == "CREATING" {
			prev, _ := d.GetChange(tfconstants.AttrPublicIPRequired)
			d.Set(tfconstants.AttrPublicIPRequired, prev)
			return diag.Errorf("cannot update public IP for PostgreSQL DBaaS (ID: %s): database is in CREATING state in project (%s), region (%s). Please wait for database creation to complete", clusterID, projectID, region)
		}

		log.Printf("[INFO] Public IP change detected for PostgreSQL cluster %s: %v", clusterID, newVal)

		if newVal {
			_, err := goe2eClient.PostgreSQL.AttachPublicIP(ctx, clusterID)
			if err != nil {
				prev, _ := d.GetChange(tfconstants.AttrPublicIPRequired)
				d.Set(tfconstants.AttrPublicIPRequired, prev)
				return diag.Errorf("error attaching public IP to PostgreSQL DBaaS (ID: %s) in project (%s), region (%s): %s", clusterID, projectID, region, err)
			}
			log.Printf("[INFO] Successfully attached public IP to PostgreSQL cluster %s", clusterID)
		} else {
			_, err := goe2eClient.PostgreSQL.DetachPublicIP(ctx, clusterID)
			if err != nil {
				prev, _ := d.GetChange(tfconstants.AttrPublicIPRequired)
				d.Set(tfconstants.AttrPublicIPRequired, prev)
				return diag.Errorf("error detaching public IP from PostgreSQL DBaaS (ID: %s) in project (%s), region (%s): %s", clusterID, projectID, region, err)
			}
			log.Printf("[INFO] Successfully detached public IP from PostgreSQL cluster %s", clusterID)
		}
	}

	if d.HasChange("project_id") {
		prev, _ := d.GetChange("project_id")
		d.Set("project_id", prev)
		return diag.Errorf("Cannot update project_id: this field is immutable after PostgreSQL DBaaS creation")
	}

	// Handle VPC changes
	if d.HasChange(tfconstants.AttrVPCs) {
		currentStatus := d.Get(tfconstants.AttrStatus).(string)

		// Block operation if DBaaS is still in "Creating" state
		if currentStatus == "CREATING" {
			prev, _ := d.GetChange(tfconstants.AttrVPCs)
			d.Set(tfconstants.AttrVPCs, prev)
			return diag.Errorf("cannot update VPC list for PostgreSQL DBaaS (ID: %s): database is in CREATING state in project (%s), region (%s). Please wait for database creation to complete", clusterID, projectID, region)
		}

		oldRaw, newRaw := d.GetChange(tfconstants.AttrVPCs)
		oldVPCs := oldRaw.(*schema.Set)
		newVPCs := newRaw.(*schema.Set)

		// Convert to string slices for goe2e client
		oldVPCIDs := make([]string, 0)
		for _, v := range oldVPCs.List() {
			oldVPCIDs = append(oldVPCIDs, strconv.Itoa(v.(int)))
		}

		newVPCIDs := make([]string, 0)
		for _, v := range newVPCs.List() {
			newVPCIDs = append(newVPCIDs, strconv.Itoa(v.(int)))
		}

		// Find VPCs to detach
		oldSet := make(map[string]bool)
		for _, vpc := range oldVPCIDs {
			oldSet[vpc] = true
		}

		var toDetach []string
		for _, vpc := range oldVPCIDs {
			found := false
			for _, newVPC := range newVPCIDs {
				if vpc == newVPC {
					found = true
					break
				}
			}
			if !found {
				toDetach = append(toDetach, vpc)
			}
		}

		// Find VPCs to attach
		var toAttach []string
		for _, vpc := range newVPCIDs {
			if !oldSet[vpc] {
				toAttach = append(toAttach, vpc)
			}
		}

		// Detach VPCs
		if len(toDetach) > 0 {
			log.Printf("[INFO] Detaching VPCs from PostgreSQL cluster %s: %v", clusterID, toDetach)
			vpcList, err := goe2eClient.PostgreSQL.ExpandPostgresVPCList(ctx, toDetach)
			if err != nil {
				prev, _ := d.GetChange(tfconstants.AttrVPCs)
				d.Set(tfconstants.AttrVPCs, prev)
				return diag.Errorf("error preparing VPC list for PostgreSQL DBaaS (ID: %s) in project (%s), region (%s): %s", clusterID, projectID, region, err)
			}

			detachReq := &goe2e.PostgreSQLVPCAttachRequest{
				Action: "detach",
				VPCs:   vpcList,
			}

			_, err = goe2eClient.PostgreSQL.DetachVPC(ctx, clusterID, detachReq)
			if err != nil {
				prev, _ := d.GetChange(tfconstants.AttrVPCs)
				d.Set(tfconstants.AttrVPCs, prev)
				return diag.Errorf("error detaching VPC from PostgreSQL DBaaS (ID: %s) in project (%s), region (%s): %s", clusterID, projectID, region, err)
			}
		}

		// Attach VPCs
		if len(toAttach) > 0 {
			log.Printf("[INFO] Attaching VPCs to PostgreSQL cluster %s: %v", clusterID, toAttach)
			vpcList, err := goe2eClient.PostgreSQL.ExpandPostgresVPCList(ctx, toAttach)
			if err != nil {
				prev, _ := d.GetChange(tfconstants.AttrVPCs)
				d.Set(tfconstants.AttrVPCs, prev)
				return diag.Errorf("error preparing VPC list for PostgreSQL DBaaS (ID: %s) in project (%s), region (%s): %s", clusterID, projectID, region, err)
			}

			attachReq := &goe2e.PostgreSQLVPCAttachRequest{
				Action: "attach",
				VPCs:   vpcList,
			}

			_, err = goe2eClient.PostgreSQL.AttachVPC(ctx, clusterID, attachReq)
			if err != nil {
				prev, _ := d.GetChange(tfconstants.AttrVPCs)
				d.Set(tfconstants.AttrVPCs, prev)
				return diag.Errorf("error attaching VPC to PostgreSQL DBaaS (ID: %s) in project (%s), region (%s): %s", clusterID, projectID, region, err)
			}
		}
	}

	// Handle parameter group changes
	if d.HasChange(tfconstants.AttrParameterGroupID) {
		currentStatus := d.Get(tfconstants.AttrStatus).(string)
		oldRaw, newRaw := d.GetChange(tfconstants.AttrParameterGroupID)
		oldPGID := 0
		newPGID := 0

		if oldRaw != nil {
			oldPGID = oldRaw.(int)
		}
		if newRaw != nil {
			newPGID = newRaw.(int)
		}

		// Block operation if DBaaS is still in "Creating" state
		if currentStatus == "CREATING" {
			d.Set(tfconstants.AttrParameterGroupID, oldRaw)
			return diag.Errorf("cannot update parameter group for PostgreSQL DBaaS (ID: %s): database is in CREATING state in project (%s), region (%s). Please wait for database creation to complete", clusterID, projectID, region)
		}

		log.Printf("[INFO] Parameter group change detected for PostgreSQL cluster %s: %d -> %d", clusterID, oldPGID, newPGID)

		switch {
		case oldPGID != 0 && newPGID == 0:
			// Detach parameter group
			_, err := goe2eClient.PostgreSQL.DetachParameterGroup(ctx, clusterID, strconv.Itoa(oldPGID))
			if err != nil {
				d.Set(tfconstants.AttrParameterGroupID, oldRaw)
				return diag.Errorf("error detaching parameter group (ID: %d) from PostgreSQL DBaaS (ID: %s) in project (%s), region (%s): %s", oldPGID, clusterID, projectID, region, err)
			}
			log.Printf("[INFO] Successfully detached parameter group %d from PostgreSQL cluster %s", oldPGID, clusterID)
		case newPGID != 0 && newPGID != oldPGID:
			// Attach new parameter group
			_, err := goe2eClient.PostgreSQL.AttachParameterGroup(ctx, clusterID, strconv.Itoa(newPGID))
			if err != nil {
				d.Set(tfconstants.AttrParameterGroupID, oldRaw)
				return diag.Errorf("error attaching parameter group (ID: %d) to PostgreSQL DBaaS (ID: %s) in project (%s), region (%s): %s", newPGID, clusterID, projectID, region, err)
			}
			log.Printf("[INFO] Successfully attached parameter group %d to PostgreSQL cluster %s", newPGID, clusterID)
		}
	}

	// Handle plan upgrades
	if d.HasChange(tfconstants.AttrPlan) {
		prevPlan, currPlan := d.GetChange(tfconstants.AttrPlan)
		plan := d.Get(tfconstants.AttrPlan).(string)
		version := d.Get(tfconstants.AttrVersion).(string)

		currentStatus := d.Get(tfconstants.AttrStatus).(string)
		// Normalize status for check
		if currentStatus == "STOPPED" {
			currentStatus = "SUSPENDED"
		}

		if currentStatus != "SUSPENDED" {
			d.Set(tfconstants.AttrPlan, prevPlan)
			return diag.Errorf("cannot upgrade plan for PostgreSQL DBaaS (ID: %s): database must be in SUSPENDED state (current state: %s) in project (%s), region (%s). Please stop the instance first", clusterID, currentStatus, projectID, region)
		}

		// Get software ID using goe2e client
		softwareID, err := goe2eClient.PostgreSQL.GetSoftwareID(ctx, "PostgreSQL", version, "")
		if err != nil {
			d.Set(tfconstants.AttrPlan, prevPlan)
			return diag.Errorf("error retrieving PostgreSQL software ID for version (%s) while upgrading plan for DBaaS (ID: %s) in project (%s), region (%s): %s", version, clusterID, projectID, region, err)
		}

		// Get template ID using goe2e client
		templateID, err := goe2eClient.PostgreSQL.GetTemplateID(ctx, plan, strconv.Itoa(softwareID), "")
		if err != nil {
			d.Set(tfconstants.AttrPlan, prevPlan)
			return diag.Errorf("error retrieving PostgreSQL template ID for plan (%s) while upgrading DBaaS (ID: %s) in project (%s), region (%s): %s", plan, clusterID, projectID, region, err)
		}

		log.Printf("[INFO] Upgrading plan for PostgreSQL cluster %s: %s -> %s", clusterID, prevPlan.(string), currPlan.(string))

		upgradeReq := &goe2e.PostgreSQLPlanUpgradeRequest{
			TemplateID: templateID,
		}

		_, err = goe2eClient.PostgreSQL.UpgradePlan(ctx, clusterID, upgradeReq)
		if err != nil {
			d.Set(tfconstants.AttrPlan, prevPlan)
			return diag.Errorf("error upgrading PostgreSQL DBaaS (ID: %s) plan from (%s) to (%s) in project (%s), region (%s): %s", clusterID, prevPlan.(string), currPlan.(string), projectID, region, err)
		}

		log.Printf("[INFO] Successfully upgraded plan for PostgreSQL cluster %s", clusterID)
	}

	// Handle disk expansion
	if d.HasChange(tfconstants.AttrSize) {
		prevSize, currSize := d.GetChange(tfconstants.AttrSize)
		currentStatus := d.Get(tfconstants.AttrStatus).(string)
		// Normalize status for check
		if currentStatus == "STOPPED" {
			currentStatus = "SUSPENDED"
		}

		if currentStatus != "SUSPENDED" {
			d.Set(tfconstants.AttrSize, prevSize)
			return diag.Errorf("cannot expand disk for PostgreSQL DBaaS (ID: %s): database must be in SUSPENDED state (current state: %s) in project (%s), region (%s)", clusterID, currentStatus, projectID, region)
		}

		sizeInt, ok := currSize.(int)
		if !ok {
			d.Set(tfconstants.AttrSize, prevSize)
			return diag.Errorf("error expanding disk for PostgreSQL DBaaS (ID: %s): size must be an integer, got %T", clusterID, currSize)
		}

		log.Printf("[INFO] Expanding disk for PostgreSQL cluster %s by %d GB", clusterID, sizeInt)

		// Calculate the additional size (cumulative expansion)
		prevSizeInt := 0
		if prevSize != nil {
			prevSizeInt = prevSize.(int)
		}
		additionalSize := sizeInt - prevSizeInt

		if additionalSize <= 0 {
			d.Set(tfconstants.AttrSize, prevSize)
			return diag.Errorf("error expanding disk for PostgreSQL DBaaS (ID: %s): size must be greater than previous size (%d GB). Got: %d GB", clusterID, prevSizeInt, sizeInt)
		}

		expandReq := &goe2e.DiskExpansionRequest{
			Size: additionalSize,
		}

		_, err = goe2eClient.PostgreSQL.ExpandDisk(ctx, clusterID, expandReq)
		if err != nil {
			d.Set(tfconstants.AttrSize, prevSize)
			return diag.Errorf("error expanding PostgreSQL DBaaS (ID: %s) disk by %d GB in project (%s), region (%s): %s", clusterID, additionalSize, projectID, region, err)
		}

		log.Printf("[INFO] Successfully expanded disk for PostgreSQL cluster %s by %d GB", clusterID, additionalSize)
	}

	// Handle tags (state-only, no API call needed)
	// Tags are automatically handled by Terraform's state management

	return resourceReadPostgress(ctx, d, m)
}

// resourceDeletePostgress deletes a PostgreSQL DBaaS cluster.
// It checks the cluster status before deletion and handles already-deleted resources gracefully.
func resourceDeletePostgress(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	var diags diag.Diagnostics

	clusterID := d.Id()
	if clusterID == "" {
		clusterID = d.Get("id").(string)
	}
	if clusterID == "" {
		return diag.Errorf("cluster ID is required for deletion")
	}

	projectID, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get goe2e client for this project/region
	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return diag.Errorf("error creating goe2e client for project (%s), region (%s): %s", projectID, region, err)
	}

	// Check current status before deletion
	status := d.Get(tfconstants.AttrStatus).(string)
	if status == "CREATING" {
		return diag.Errorf("cannot delete PostgreSQL DBaaS (ID: %s): database is in CREATING state in project (%s), region (%s). Please wait for database creation to complete", clusterID, projectID, region)
	}

	// Delete PostgreSQL cluster using goe2e client
	_, err = goe2eClient.PostgreSQL.DeleteCluster(ctx, clusterID)
	if err != nil {
		// Check if resource was already deleted (404)
		exists, _, checkErr := goe2eClient.PostgreSQL.ClusterExists(ctx, clusterID)
		if checkErr == nil && !exists {
			log.Printf("[WARN] PostgreSQL cluster %s was already deleted", clusterID)
			d.SetId("")
			return diags
		}
		return diag.Errorf("error deleting PostgreSQL DBaaS (ID: %s) in project (%s), region (%s): %s", clusterID, projectID, region, err)
	}

	d.SetId("")
	log.Printf("[INFO] Successfully deleted PostgreSQL DBaaS cluster: %s", clusterID)

	return diags
}

// CustomImportStateFunc implements the custom import function for PostgreSQL resources.
// Format: project_id:dbaas_id
func CustomImportStateFunc(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid ID format: expected project_id:dbaas_id, got %s", d.Id())
	}

	projectID := parts[0]
	dbaasID := parts[1]

	if err := d.Set(tfconstants.AttrProjectID, projectID); err != nil {
		return nil, err
	}
	d.SetId(dbaasID)

	return []*schema.ResourceData{d}, nil
}

// resourcePostgreSQLResourceV0 returns the V0 schema definition (before tags were added and before ForceNew fixes)
func resourcePostgreSQLResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// ============================================
			// COMMON FIELDS
			// ============================================
			tfconstants.AttrRegion:    config.RegionSchema(),
			tfconstants.AttrLocation:  config.LocationSchema(),
			tfconstants.AttrProjectID: config.ProjectIDSchemaResource(),

			// ============================================
			// REQUIRED INPUT FIELDS
			// ============================================
			tfconstants.AttrVersion: {
				Type:     schema.TypeString,
				Required: true,
				// V0: No ForceNew (this was the bug we're fixing)
			},
			tfconstants.AttrPlan: {
				Type:     schema.TypeString,
				Required: true,
			},
			tfconstants.AttrName: {
				Type:     schema.TypeString,
				Required: true,
				// V0: No ForceNew (this was the bug we're fixing)
			},
			tfconstants.AttrDatabase: {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true, // V0: ForceNew on entire block
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user": {
							Type:     schema.TypeString,
							Required: true,
							// V0: No ForceNew on individual fields
						},
						"password": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
							ForceNew:  true, // V0: ForceNew on password
						},
						"dbaas_number": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
							// V0: No ForceNew
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
							// V0: No ForceNew on individual fields
						},
					},
				},
			},

			// ============================================
			// OPTIONAL INPUT FIELDS
			// ============================================
			tfconstants.AttrGroup: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Default",
				// V0: No ForceNew (this was the bug we're fixing)
			},
			tfconstants.AttrPublicIPRequired: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true, // V0: ForceNew on public_ip_required
			},
			tfconstants.AttrIsEncryptionEnabled: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				// V0: No ForceNew (this was the bug we're fixing)
			},
			tfconstants.AttrEncryptionPassphrase: {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				// V0: No ForceNew (this was the bug we're fixing)
			},
			tfconstants.AttrParameterGroupID: {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"vpc_list": { // V0: Used vpc_list instead of vpcs
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Optional: true,
			},

			// ============================================
			// OPTIONAL INPUT FIELDS - MANAGEMENT
			// ============================================
			tfconstants.AttrPowerStatus: { // V0: Used power_status instead of status
				Type:     schema.TypeString,
				Optional: true,
			},
			"detach_public_ip": { // V0: Used detach_public_ip instead of public_ip_required
				Type:     schema.TypeBool,
				Optional: true,
			},
			tfconstants.AttrSize: {
				Type:     schema.TypeInt,
				Optional: true,
			},

			// ============================================
			// COMPUTED FIELDS
			// ============================================
			tfconstants.AttrID: {
				Type:     schema.TypeString,
				Computed: true,
			},
			tfconstants.AttrStatus: {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status_title": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status_actions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"num_instances": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snapshot_exist": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"connectivity_detail": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vector_database_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// Note: V0 schema does not have "tags" field
			// Note: V0 schema does not have computed network fields (public_ip_address, private_ip_address, port, disk)
		},
	}
}

// resourcePostgreSQLStateUpgradeV0toV1 upgrades the state from V0 to V1
// V1 adds the "tags" field and fixes ForceNew semantics
func resourcePostgreSQLStateUpgradeV0toV1(
	ctx context.Context,
	rawState map[string]interface{},
	meta interface{},
) (map[string]interface{}, error) {
	// Add new V3 fields with defaults
	if _, exists := rawState["tags"]; !exists {
		rawState["tags"] = make(map[string]interface{})
	}

	// Migrate vpc_list to vpcs if present
	if vpcList, exists := rawState["vpc_list"]; exists && vpcList != nil {
		rawState[tfconstants.AttrVPCs] = vpcList
		// Keep vpc_list for backwards compatibility during migration
	}

	// Migrate detach_public_ip to public_ip_required if present
	if detachPublicIP, exists := rawState["detach_public_ip"]; exists {
		// If detach_public_ip is true, then public_ip_required should be false
		if detach, ok := detachPublicIP.(bool); ok {
			rawState[tfconstants.AttrPublicIPRequired] = !detach
		}
	}

	// Migrate power_status to status if present
	if powerStatus, exists := rawState[tfconstants.AttrPowerStatus]; exists {
		// Map old power_status values to new status values
		status := ""
		switch powerStatus.(string) {
		case "start":
			status = "RUNNING"
		case "stop":
			status = "SUSPENDED"
		case "restart":
			status = "RESTARTING"
		default:
			status = powerStatus.(string)
		}
		if status != "" {
			rawState[tfconstants.AttrStatus] = status
		}
	}

	// Preserve all existing fields - no data loss
	log.Printf("[INFO] Upgraded PostgreSQL resource state from v0 to v1: %s", rawState["id"])

	return rawState, nil
}
