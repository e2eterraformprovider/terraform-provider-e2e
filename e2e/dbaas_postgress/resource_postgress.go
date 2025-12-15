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
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
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
				Upgrade: ResourcePostgreSQLStateUpgradeV0toV1,
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
					[]string{
						goe2econstants.PostgreSQLVersion11,
						goe2econstants.PostgreSQLVersion12,
						goe2econstants.PostgreSQLVersion13,
						goe2econstants.PostgreSQLVersion14,
						goe2econstants.PostgreSQLVersion15,
					},
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
						tfconstants.AttrDatabaseBlockUser: {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true, // ✅ ADDED - initial admin user immutable
							Description: "the database username",
						},
						tfconstants.AttrDatabaseBlockPassword: {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
							// No ForceNew - password rotation supported
							Description: "the database password",
						},
						tfconstants.AttrDatabaseBlockDBaaSNumber: {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     tfconstants.DBaaSDefaultDBaaSNumber,
							ForceNew:    true, // ✅ ADDED - topology immutable
							Description: "the DBaaS number (typically 1)",
						},
						tfconstants.AttrDatabaseBlockName: {
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
				Default:     tfconstants.DBaaSDefaultGroupName,
			},
			tfconstants.AttrPublicIPRequired: {
				Type:     schema.TypeBool,
				Optional: true,
				// No ForceNew - can be attached/detached dynamically
				Description: "whether to attach a public IP to the PostgreSQL DBaaS instance",
				Default:     tfconstants.DBaaSDefaultPublicIPRequired,
			},
			tfconstants.AttrIsEncryptionEnabled: {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     tfconstants.DBaaSDefaultIsEncryptionEnabled,
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
					[]string{
						goe2econstants.DBaaSStatusStopped,
						goe2econstants.DBaaSStatusSuspended,
						goe2econstants.DBaaSStatusRunning,
						goe2econstants.DBaaSStatusRestarting,
					},
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
			tfconstants.AttrStatusTitle: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the status title of the PostgreSQL DBaaS instance",
			},
			tfconstants.AttrStatusActions: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "available actions for the PostgreSQL DBaaS instance",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			tfconstants.AttrNumInstances: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "the number of instances in the PostgreSQL DBaaS cluster",
			},
			tfconstants.AttrProjectName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "name of the project",
			},
			tfconstants.AttrSnapshotExist: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "whether a snapshot exists for the PostgreSQL DBaaS instance",
			},
			tfconstants.AttrConnectivityDetail: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the connectivity details for the PostgreSQL DBaaS instance",
			},
			tfconstants.AttrVectorDatabaseStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the vector database status of the PostgreSQL DBaaS instance",
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
			tfconstants.AttrPort: {
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

	log.Printf("[INFO] Creating %s", ResourceName)

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
		return diag.Errorf(tfconstants.DatabaseConfigurationRequired)
	}
	dbConfigMap := dbConfigList[0].(map[string]interface{})

	plan := d.Get(tfconstants.AttrPlan).(string)
	version := d.Get(tfconstants.AttrVersion).(string)

	// Get software ID using goe2e client
	// Note: pgID parameter is required but can be empty string for PostgreSQL
	softwareID, err := goe2eClient.PostgreSQL.GetSoftwareID(ctx, goe2econstants.DBaaSSoftwarePostgreSQL, version, "")
	if err != nil {
		return diag.Errorf(ErrorRetrievingSoftwareIDTemplate, version, projectID, region, err)
	}

	// Get template ID using goe2e client
	templateID, err := goe2eClient.PostgreSQL.GetTemplateID(ctx, plan, strconv.Itoa(softwareID), "")
	if err != nil {
		return diag.Errorf(ErrorRetrievingTemplateIDTemplate, plan, projectID, region, err)
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
	clusterName := d.Get(tfconstants.AttrName).(string)
	if len(vpcIDs) > 0 {
		vpcList, err = goe2eClient.PostgreSQL.ExpandPostgresVPCList(ctx, vpcIDs)
		if err != nil {
			return diag.Errorf(ErrorPreparingVPCListTemplate, clusterName, projectID, region, err)
		}
	}

	// Build create request using goe2e types
	req := &goe2e.PostgreSQLClusterCreateRequest{
		SoftwareID:       softwareID,
		TemplateID:       templateID,
		Name:             clusterName,
		Group:            d.Get(tfconstants.AttrGroup).(string),
		PublicIPRequired: d.Get(tfconstants.AttrPublicIPRequired).(bool),
		VPCs:             vpcList,
		Database: goe2e.DBConfig{
			User:        dbConfigMap[tfconstants.AttrDatabaseBlockUser].(string),
			Password:    dbConfigMap[tfconstants.AttrDatabaseBlockPassword].(string),
			DBaaSNumber: dbConfigMap[tfconstants.AttrDatabaseBlockDBaaSNumber].(int),
			Name:        dbConfigMap[tfconstants.AttrDatabaseBlockName].(string),
		},
		PGID:                pgID,
		IsEncryptionEnabled: d.Get(tfconstants.AttrIsEncryptionEnabled).(bool),
	}

	// Create PostgreSQL cluster using goe2e client
	cluster, _, err := goe2eClient.PostgreSQL.CreateCluster(ctx, req)
	if err != nil {
		return diag.Errorf(tfconstants.ResourceOperationErrorTemplate, tfconstants.OperationCreating, ResourceName, req.Name, projectID, region, err)
	}

	// Set resource ID and attributes
	d.SetId(strconv.Itoa(cluster.ID))
	if err := d.Set(tfconstants.AttrName, cluster.Name); err != nil {
		return diag.FromErr(err)
	}

	// Normalize status (SUSPENDED → STOPPED for consistency)
	status := cluster.Status
	if status == goe2econstants.DBaaSStatusSuspended {
		status = goe2econstants.DBaaSStatusStopped
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
	if err := d.Set(tfconstants.AttrPort, cluster.MasterNode.Port); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(tfconstants.AttrDisk, cluster.MasterNode.Disk); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(tfconstants.AttrIsEncryptionEnabled, cluster.IsEncryptionEnabled); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Successfully created %s: %s (ID: %d)", ResourceName, cluster.Name, cluster.ID)

	return diags
}

// resourceReadPostgress reads the current state of a PostgreSQL DBaaS cluster.
// It retrieves the cluster details from the API and updates the Terraform state.
func resourceReadPostgress(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	var diags diag.Diagnostics

	log.Printf("[INFO] Reading %s", ResourceName)

	clusterID := d.Id()
	if clusterID == "" {
		clusterID = d.Get(tfconstants.AttrID).(string)
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
		return diag.Errorf(tfconstants.ErrorCreatingGoe2eClient, err)
	}

	// Get PostgreSQL cluster using goe2e client
	cluster, _, err := goe2eClient.PostgreSQL.GetCluster(ctx, clusterID)
	if err != nil {
		return diag.Errorf(tfconstants.ResourceOperationByIDErrorTemplate, tfconstants.OperationRetrieving, ResourceName, clusterID, projectID, region, err)
	}

	// Check if resource was deleted
	if cluster == nil {
		d.SetId("")
		return diags
	}

	// Set resource ID
	d.SetId(strconv.Itoa(cluster.ID))
	if err := d.Set(tfconstants.AttrID, cluster.ID); err != nil {
		return diag.FromErr(err)
	}

	// Set basic fields
	if err := d.Set(tfconstants.AttrName, cluster.Name); err != nil {
		return diag.FromErr(err)
	}

	// Normalize status (SUSPENDED → STOPPED for consistency)
	status := cluster.Status
	if status == goe2econstants.DBaaSStatusSuspended {
		status = goe2econstants.DBaaSStatusStopped
	}
	if err := d.Set(tfconstants.AttrStatus, status); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(tfconstants.AttrStatusTitle, cluster.StatusTitle); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(tfconstants.AttrStatusActions, cluster.StatusActions); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(tfconstants.AttrNumInstances, cluster.NumInstances); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(tfconstants.AttrProjectName, cluster.ProjectName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(tfconstants.AttrSnapshotExist, cluster.SnapshotExist); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(tfconstants.AttrConnectivityDetail, cluster.ConnectivityDetail); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(tfconstants.AttrVectorDatabaseStatus, cluster.VectorDBStatus); err != nil {
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
	if err := d.Set(tfconstants.AttrPort, cluster.MasterNode.Port); err != nil {
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

	log.Printf("[INFO] Successfully read %s: %s (ID: %d)", ResourceName, cluster.Name, cluster.ID)

	return diags
}

// resourceUpdatePostgress updates a PostgreSQL DBaaS cluster.
// It handles updates to: status (power management), public IP, VPCs, parameter groups, plan, and disk size.
// Each update type has its own validation and API call logic.
func resourceUpdatePostgress(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)

	clusterID := d.Id()
	if clusterID == "" {
		clusterID = d.Get(tfconstants.AttrID).(string)
	}
	if clusterID == "" {
		return diag.Errorf(ClusterIDRequiredForUpdate)
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
		if currentStatus == goe2econstants.DBaaSStatusCreating {
			return diag.Errorf(ErrorCannotPerformPowerOpsTemplate, clusterID, projectID, region)
		}

		log.Printf("[INFO] Status change detected for %s %s: %s -> %s", ResourceName, clusterID, currentStatus, newStatus)

		switch strings.ToUpper(newStatus) {
		case goe2econstants.DBaaSStatusStopped, goe2econstants.DBaaSStatusSuspended:
			_, err := goe2eClient.PostgreSQL.StopCluster(ctx, clusterID)
			if err != nil {
				return diag.Errorf(ErrorStoppingTemplate, clusterID, projectID, region, err)
			}
			log.Printf("[INFO] Successfully stopped %s %s", ResourceName, clusterID)
		case goe2econstants.DBaaSStatusRunning:
			_, err := goe2eClient.PostgreSQL.StartCluster(ctx, clusterID)
			if err != nil {
				return diag.Errorf(ErrorStartingTemplate, clusterID, projectID, region, err)
			}
			log.Printf("[INFO] Successfully started %s %s", ResourceName, clusterID)
		case goe2econstants.DBaaSStatusRestarting:
			_, err := goe2eClient.PostgreSQL.RestartCluster(ctx, clusterID)
			if err != nil {
				return diag.Errorf(ErrorRestartingTemplate, clusterID, projectID, region, err)
			}
			log.Printf("[INFO] Successfully restarted %s %s", ResourceName, clusterID)
		default:
			return diag.Errorf(ErrorUnsupportedStatusTemplate, clusterID, newStatus)
		}
	}

	// Handle public IP changes
	if d.HasChange(tfconstants.AttrPublicIPRequired) {
		newVal := d.Get(tfconstants.AttrPublicIPRequired).(bool)
		currentStatus := d.Get(tfconstants.AttrStatus).(string)

		// Block operation if DBaaS is still in "Creating" state
		if currentStatus == goe2econstants.DBaaSStatusCreating {
			prev, _ := d.GetChange(tfconstants.AttrPublicIPRequired)
			d.Set(tfconstants.AttrPublicIPRequired, prev)
			return diag.Errorf(ErrorCannotUpdatePublicIPTemplate, clusterID, projectID, region)
		}

		log.Printf("[INFO] Public IP change detected for %s %s: %v", ResourceName, clusterID, newVal)

		if newVal {
			_, err := goe2eClient.PostgreSQL.AttachPublicIP(ctx, clusterID)
			if err != nil {
				prev, _ := d.GetChange(tfconstants.AttrPublicIPRequired)
				d.Set(tfconstants.AttrPublicIPRequired, prev)
				return diag.Errorf(ErrorAttachingPublicIPTemplate, clusterID, projectID, region, err)
			}
			log.Printf("[INFO] Successfully attached public IP to %s %s", ResourceName, clusterID)
		} else {
			_, err := goe2eClient.PostgreSQL.DetachPublicIP(ctx, clusterID)
			if err != nil {
				prev, _ := d.GetChange(tfconstants.AttrPublicIPRequired)
				d.Set(tfconstants.AttrPublicIPRequired, prev)
				return diag.Errorf(ErrorDetachingPublicIPTemplate, clusterID, projectID, region, err)
			}
			log.Printf("[INFO] Successfully detached public IP from %s %s", ResourceName, clusterID)
		}
	}

	if d.HasChange(tfconstants.AttrProjectID) {
		prev, _ := d.GetChange(tfconstants.AttrProjectID)
		d.Set(tfconstants.AttrProjectID, prev)
		return diag.Errorf(ErrorProjectIDImmutable)
	}

	// Handle VPC changes
	if d.HasChange(tfconstants.AttrVPCs) {
		currentStatus := d.Get(tfconstants.AttrStatus).(string)

		// Block operation if DBaaS is still in "Creating" state
		if currentStatus == goe2econstants.DBaaSStatusCreating {
			prev, _ := d.GetChange(tfconstants.AttrVPCs)
			d.Set(tfconstants.AttrVPCs, prev)
			return diag.Errorf(ErrorCannotUpdateVPCListTemplate, clusterID, projectID, region)
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
			log.Printf("[INFO] Detaching VPCs from %s %s: %v", ResourceName, clusterID, toDetach)
			vpcList, err := goe2eClient.PostgreSQL.ExpandPostgresVPCList(ctx, toDetach)
			if err != nil {
				prev, _ := d.GetChange(tfconstants.AttrVPCs)
				d.Set(tfconstants.AttrVPCs, prev)
				return diag.Errorf(ErrorPreparingVPCListTemplate, clusterID, projectID, region, err)
			}

			detachReq := &goe2e.PostgreSQLVPCAttachRequest{
				Action: goe2econstants.ActionDetach,
				VPCs:   vpcList,
			}

			_, err = goe2eClient.PostgreSQL.DetachVPC(ctx, clusterID, detachReq)
			if err != nil {
				prev, _ := d.GetChange(tfconstants.AttrVPCs)
				d.Set(tfconstants.AttrVPCs, prev)
				return diag.Errorf(ErrorDetachingVPCTemplate, clusterID, projectID, region, err)
			}
		}

		// Attach VPCs
		if len(toAttach) > 0 {
			log.Printf("[INFO] Attaching VPCs to %s %s: %v", ResourceName, clusterID, toAttach)
			vpcList, err := goe2eClient.PostgreSQL.ExpandPostgresVPCList(ctx, toAttach)
			if err != nil {
				prev, _ := d.GetChange(tfconstants.AttrVPCs)
				d.Set(tfconstants.AttrVPCs, prev)
				return diag.Errorf(ErrorPreparingVPCListTemplate, clusterID, projectID, region, err)
			}

			attachReq := &goe2e.PostgreSQLVPCAttachRequest{
				Action: goe2econstants.ActionAttach,
				VPCs:   vpcList,
			}

			_, err = goe2eClient.PostgreSQL.AttachVPC(ctx, clusterID, attachReq)
			if err != nil {
				prev, _ := d.GetChange(tfconstants.AttrVPCs)
				d.Set(tfconstants.AttrVPCs, prev)
				return diag.Errorf(ErrorAttachingVPCTemplate, clusterID, projectID, region, err)
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
		if currentStatus == goe2econstants.DBaaSStatusCreating {
			d.Set(tfconstants.AttrParameterGroupID, oldRaw)
			return diag.Errorf(ErrorCannotUpdateParameterGroupTemplate, clusterID, projectID, region)
		}

		log.Printf("[INFO] Parameter group change detected for %s %s: %d -> %d", ResourceName, clusterID, oldPGID, newPGID)

		switch {
		case oldPGID != 0 && newPGID == 0:
			// Detach parameter group
			_, err := goe2eClient.PostgreSQL.DetachParameterGroup(ctx, clusterID, strconv.Itoa(oldPGID))
			if err != nil {
				d.Set(tfconstants.AttrParameterGroupID, oldRaw)
				return diag.Errorf(ErrorDetachingParameterGroupTemplate, oldPGID, clusterID, projectID, region, err)
			}
			log.Printf("[INFO] Successfully detached parameter group %d from %s %s", oldPGID, ResourceName, clusterID)
		case newPGID != 0 && newPGID != oldPGID:
			// Attach new parameter group
			_, err := goe2eClient.PostgreSQL.AttachParameterGroup(ctx, clusterID, strconv.Itoa(newPGID))
			if err != nil {
				d.Set(tfconstants.AttrParameterGroupID, oldRaw)
				return diag.Errorf(ErrorAttachingParameterGroupTemplate, newPGID, clusterID, projectID, region, err)
			}
			log.Printf("[INFO] Successfully attached parameter group %d to %s %s", newPGID, ResourceName, clusterID)
		}
	}

	// Handle plan upgrades
	if d.HasChange(tfconstants.AttrPlan) {
		prevPlan, currPlan := d.GetChange(tfconstants.AttrPlan)
		plan := d.Get(tfconstants.AttrPlan).(string)
		version := d.Get(tfconstants.AttrVersion).(string)

		currentStatus := d.Get(tfconstants.AttrStatus).(string)
		// Normalize status for check
		if currentStatus == goe2econstants.DBaaSStatusStopped {
			currentStatus = goe2econstants.DBaaSStatusSuspended
		}

		if currentStatus != goe2econstants.DBaaSStatusSuspended {
			d.Set(tfconstants.AttrPlan, prevPlan)
			return diag.Errorf(ErrorCannotUpgradePlanTemplate, clusterID, currentStatus, projectID, region)
		}

		// Get software ID using goe2e client
		softwareID, err := goe2eClient.PostgreSQL.GetSoftwareID(ctx, goe2econstants.DBaaSSoftwarePostgreSQL, version, "")
		if err != nil {
			d.Set(tfconstants.AttrPlan, prevPlan)
			return diag.Errorf(ErrorRetrievingSoftwareIDForUpgrade, version, clusterID, projectID, region, err)
		}

		// Get template ID using goe2e client
		templateID, err := goe2eClient.PostgreSQL.GetTemplateID(ctx, plan, strconv.Itoa(softwareID), "")
		if err != nil {
			d.Set(tfconstants.AttrPlan, prevPlan)
			return diag.Errorf(ErrorRetrievingTemplateIDForUpgrade, plan, clusterID, projectID, region, err)
		}

		log.Printf("[INFO] Upgrading plan for %s %s: %s -> %s", ResourceName, clusterID, prevPlan.(string), currPlan.(string))

		upgradeReq := &goe2e.PostgreSQLPlanUpgradeRequest{
			TemplateID: templateID,
		}

		_, err = goe2eClient.PostgreSQL.UpgradePlan(ctx, clusterID, upgradeReq)
		if err != nil {
			d.Set(tfconstants.AttrPlan, prevPlan)
			return diag.Errorf(ErrorUpgradingPlanTemplate, clusterID, prevPlan.(string), currPlan.(string), projectID, region, err)
		}

		log.Printf("[INFO] Successfully upgraded plan for %s %s", ResourceName, clusterID)
	}

	// Handle disk expansion
	if d.HasChange(tfconstants.AttrSize) {
		prevSize, currSize := d.GetChange(tfconstants.AttrSize)
		currentStatus := d.Get(tfconstants.AttrStatus).(string)
		// Normalize status for check
		if currentStatus == goe2econstants.DBaaSStatusStopped {
			currentStatus = goe2econstants.DBaaSStatusSuspended
		}

		if currentStatus != goe2econstants.DBaaSStatusSuspended {
			d.Set(tfconstants.AttrSize, prevSize)
			return diag.Errorf(ErrorCannotExpandDiskTemplate, clusterID, currentStatus, projectID, region)
		}

		sizeInt, ok := currSize.(int)
		if !ok {
			d.Set(tfconstants.AttrSize, prevSize)
			return diag.Errorf(ErrorExpandingDiskInvalidTypeTemplate, clusterID, currSize)
		}

		log.Printf("[INFO] Expanding disk for %s %s by %d GB", ResourceName, clusterID, sizeInt)

		// Calculate the additional size (cumulative expansion)
		prevSizeInt := 0
		if prevSize != nil {
			prevSizeInt = prevSize.(int)
		}
		additionalSize := sizeInt - prevSizeInt

		if additionalSize <= 0 {
			d.Set(tfconstants.AttrSize, prevSize)
			return diag.Errorf(ErrorExpandingDiskInvalidSizeTemplate, clusterID, prevSizeInt, sizeInt)
		}

		expandReq := &goe2e.DiskExpansionRequest{
			Size: additionalSize,
		}

		_, err = goe2eClient.PostgreSQL.ExpandDisk(ctx, clusterID, expandReq)
		if err != nil {
			d.Set(tfconstants.AttrSize, prevSize)
			return diag.Errorf(ErrorExpandingDiskTemplate, clusterID, additionalSize, projectID, region, err)
		}

		log.Printf("[INFO] Successfully expanded disk for %s %s by %d GB", ResourceName, clusterID, additionalSize)
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
		clusterID = d.Get(tfconstants.AttrID).(string)
	}
	if clusterID == "" {
		return diag.Errorf(ClusterIDRequiredForDeletion)
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
		return diag.Errorf(tfconstants.ErrorCreatingGoe2eClient, err)
	}

	// Check current status before deletion
	status := d.Get(tfconstants.AttrStatus).(string)
	if status == goe2econstants.DBaaSStatusCreating {
		return diag.Errorf(tfconstants.ResourceDeleteStateErrorTemplate, ResourceName, clusterID, ResourceName, goe2econstants.DBaaSStatusCreating, projectID, region, "Please wait for database creation to complete")
	}

	// Delete PostgreSQL cluster using goe2e client
	_, err = goe2eClient.PostgreSQL.DeleteCluster(ctx, clusterID)
	if err != nil {
		// Check if resource was already deleted (404)
		exists, _, checkErr := goe2eClient.PostgreSQL.ClusterExists(ctx, clusterID)
		if checkErr == nil && !exists {
			log.Printf("[WARN] %s %s was already deleted", ResourceName, clusterID)
			d.SetId("")
			return diags
		}
		return diag.Errorf(tfconstants.ResourceOperationByIDErrorTemplate, tfconstants.OperationDeleting, ResourceName, clusterID, projectID, region, err)
	}

	d.SetId("")
	log.Printf("[INFO] Successfully deleted %s: %s", ResourceName, clusterID)

	return diags
}

// CustomImportStateFunc implements the custom import function for PostgreSQL resources.
// Format: project_id:dbaas_id
func CustomImportStateFunc(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
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
						tfconstants.AttrDatabaseBlockUser: {
							Type:     schema.TypeString,
							Required: true,
							// V0: No ForceNew on individual fields
						},
						tfconstants.AttrDatabaseBlockPassword: {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
							ForceNew:  true, // V0: ForceNew on password
						},
						tfconstants.AttrDatabaseBlockDBaaSNumber: {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  tfconstants.DBaaSDefaultDBaaSNumber,
							// V0: No ForceNew
						},
						tfconstants.AttrDatabaseBlockName: {
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
				Default:  tfconstants.DBaaSDefaultGroupName,
				// V0: No ForceNew (this was the bug we're fixing)
			},
			tfconstants.AttrPublicIPRequired: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  tfconstants.DBaaSDefaultPublicIPRequired,
				ForceNew: true, // V0: ForceNew on public_ip_required
			},
			tfconstants.AttrIsEncryptionEnabled: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  tfconstants.DBaaSDefaultIsEncryptionEnabled,
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

// ResourcePostgreSQLStateUpgradeV0toV1 upgrades the state from V0 to V1
// V1 adds the "tags" field and fixes ForceNew semantics
func ResourcePostgreSQLStateUpgradeV0toV1(
	ctx context.Context,
	rawState map[string]interface{},
	meta interface{},
) (map[string]interface{}, error) {
	// Add new V3 fields with defaults
	if _, exists := rawState["tags"]; !exists {
		rawState["tags"] = make(map[string]interface{})
	}

	// Migrate vpc_list to vpcs if present
	if vpcList, exists := rawState[tfconstants.FieldMigrationKeyVPCList]; exists && vpcList != nil {
		rawState[tfconstants.AttrVPCs] = vpcList
		// Keep vpc_list for backwards compatibility during migration
	}

	// Migrate detach_public_ip to public_ip_required if present
	if detachPublicIP, exists := rawState[tfconstants.FieldMigrationKeyDetachPublicIP]; exists {
		// If detach_public_ip is true, then public_ip_required should be false
		if detach, ok := detachPublicIP.(bool); ok {
			rawState[tfconstants.AttrPublicIPRequired] = !detach
		}
	}

	// Migrate power_status to status if present
	if powerStatus, exists := rawState[tfconstants.FieldMigrationKeyPowerStatus]; exists {
		// Map old power_status values to new status values
		status := ""
		if powerStatusStr, ok := powerStatus.(string); ok && powerStatusStr != "" {
			switch powerStatusStr {
			case tfconstants.DBaaSPowerActionStart:
				status = goe2econstants.DBaaSStatusRunning
			case tfconstants.DBaaSPowerActionStop:
				status = goe2econstants.DBaaSStatusSuspended
			case tfconstants.DBaaSPowerActionRestart:
				status = goe2econstants.DBaaSStatusRestarting
			default:
				status = powerStatusStr
			}
			if status != "" {
				rawState[tfconstants.AttrStatus] = status
			}
		}
	}

	// Preserve all existing fields - no data loss
	log.Printf("[INFO] Upgraded %s resource state from v0 to v1: %s", ResourceName, rawState["id"])

	return rawState, nil
}
