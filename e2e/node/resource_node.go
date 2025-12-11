package node

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/util"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceNode() *schema.Resource {
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
			tfconstants.AttrName: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "name of the Node",
				ValidateFunc: ValidateName,
			},
			tfconstants.AttrLabel: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "the label of the Node",
				Default:     "default",
			},
			tfconstants.AttrPlan: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "the plan of the Node",
				ValidateFunc: ValidatePlanName,
			},
			tfconstants.AttrImage: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "the image id or slug of the Node (format: os-version)",
				ValidateFunc: ValidateBlank,
			},

			// ============================================
			// OPTIONAL INPUT FIELDS - CREATION
			// ============================================
			"backup": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "whether the Node has backups enabled",
				Default:     false,
			},
			"default_public_ip": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "whether the Node has default public IP enabled",
				Default:     false,
			},
			"disable_password": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "whether to disable password authentication",
				Default:     false,
			},
			"enable_bitninja": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "whether to enable BitNinja security",
				Default:     false,
			},
			"is_ipv6_availed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "whether the Node has IPv6 available",
				Default:     false,
			},
			"is_saved_image": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "whether to create the Node from a saved image",
				Default:     false,
			},
			"start_script": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "the script to be run when the Node is first created",
			},
			tfconstants.AttrReserveIP: {
				Type:          schema.TypeString,
				Optional:      true,
				Deprecated:    "Use reserve_ip_id instead. This field will be removed in v4.0",
				Description:   "id of the reserved IP to attach to the Node (DEPRECATED: use reserve_ip_id)",
				ConflictsWith: []string{tfconstants.AttrReserveIPID},
			},
			tfconstants.AttrReserveIPID: {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "id of the reserved IP to attach to the Node",
				ConflictsWith: []string{tfconstants.AttrReserveIP},
			},
			tfconstants.AttrVPCID: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "id of the VPC for the Node",
			},
			"saved_image_template_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "the template id of the saved image (required when is_saved_image is true)",
				Default:     nil,
			},
			tfconstants.AttrSecurityGroupIDs: {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "list of security group ids to attach to the Node",
				Elem: &schema.Schema{
					Type:        schema.TypeInt,
					Description: "id of the security group",
				},
			},
			"ssh_keys": {
				Type:          schema.TypeList,
				Optional:      true,
				Deprecated:    "Use ssh_key_ids instead. This field will be removed in v4.0",
				Description:   "list of SSH key labels to attach to the Node (DEPRECATED: use ssh_key_ids)",
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"ssh_key_ids"},
			},
			"ssh_key_ids": {
				Type:          schema.TypeList,
				Optional:      true,
				Description:   "list of SSH key resource IDs to attach to the Node",
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"ssh_keys"},
			},
			"block_storage_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Deprecated:  "Use e2e_volume_attachment resource instead. This field will be removed in v4.0",
				Description: "list of block storage ids to attach to the Node (DEPRECATED: use e2e_volume_attachment)",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					Description:  "id of the block storage",
					ValidateFunc: validation.All(ValidateBlank, ValidateInteger),
				},
			},
			"tag_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "list of tag IDs to attach to the Node (uses label API)",
				Elem: &schema.Schema{
					Type:        schema.TypeInt,
					Description: "id of the tag",
				},
			},

			// ============================================
			// OPTIONAL STRUCTURED BLOCKS - V3
			// ============================================
			"network_interface": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "structured network configuration for the Node",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						tfconstants.AttrVPCID: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "id of the VPC for the Node",
						},
						"assign_public_ip": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "whether to assign a public IP to the Node",
						},
						"enable_ipv6": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "whether to enable IPv6 for the Node",
						},
						tfconstants.AttrSecurityGroupIDs: {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "list of security group IDs",
							Elem: &schema.Schema{
								Type:        schema.TypeInt,
								Description: "id of the security group",
							},
						},
					},
				},
			},

			tfconstants.AttrRootDisk: {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "structured root disk configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						tfconstants.AttrSizeGB: {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "the size of the root disk in gigabytes (sent to API as 'disk' field)",
						},
						tfconstants.AttrDiskType: {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Default:     "standard",
							Description: "the type of the root disk",
						},
					},
				},
			},

			// ============================================
			// OPTIONAL INPUT FIELDS - MANAGEMENT
			// ============================================
			tfconstants.AttrPowerStatus: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     goe2econstants.NodePowerStatusOn,
				Description: "the power status of the Node (power_on to start, power_off to power off)",
				ValidateFunc: validation.StringInSlice([]string{
					goe2econstants.NodePowerStatusOff,
					goe2econstants.NodePowerStatusOn,
				}, false),
			},
			tfconstants.AttrLockNode: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "whether to lock the Node",
			},

			// ============================================
			// ACTION FIELDS
			// ============================================
			tfconstants.AttrRebootNode: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "whether to reboot the Node (Node must be in running state; may cause data loss if disk-intensive processes are running)",
			},
			"reinstall_node": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "whether to reinstall the Node (Node must be in running state; this will permanently delete all data)",
			},

			// ============================================
			// COMPUTED FIELDS - STATUS
			// ============================================
			"default_sg": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "id of the default security group",
			},
			"is_active": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "whether the Node is active",
			},
			tfconstants.AttrCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the creation date for the Node",
			},
			tfconstants.AttrStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "state of the Node instance",
			},
			"price": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the price details of the Node",
			},

			// ============================================
			// COMPUTED FIELDS - NETWORK
			// ============================================
			tfconstants.AttrPublicIPAddress: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the Nodes public ipv4 address",
			},
			tfconstants.AttrPrivateIPAddress: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the Nodes private ipv4 address",
			},
			"ipv6_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the Nodes IPv6 address (if IPv6 is enabled)",
			},

			// ============================================
			// COMPUTED FIELDS - RESOURCES
			// ============================================
			tfconstants.AttrVMID: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "id of the VM",
			},
			tfconstants.AttrMemory: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "memory of the Node in megabytes",
			},
			tfconstants.AttrDisk: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the disk information of the Node",
			},
			"is_bitninja_license_active": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "whether the Node has BitNinja license active",
			},

			// ============================================
			// OPTIONAL STRUCTURED BLOCKS - BACKUP
			// ============================================
			tfconstants.AttrBackupConfig: {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "structured backup configuration for the Node",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "whether backup is enabled",
						},
						tfconstants.AttrBackupPlanID: {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "id of the backup plan",
						},
						tfconstants.AttrBackupType: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "backup frequency type (HOURLY, DAILY, WEEKLY, MONTHLY)",
							ValidateFunc: validation.StringInSlice([]string{
								"HOURLY", "DAILY", "WEEKLY", "MONTHLY",
							}, false),
						},
						tfconstants.AttrBackupExcludePaths: {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "list of paths to exclude from backup",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						tfconstants.AttrBackupNow: {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "whether to take a backup immediately",
						},
						tfconstants.AttrCompressionType: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "compression type for backup (ZLib, GZip, None)",
							ValidateFunc: validation.StringInSlice([]string{
								"ZLib", "GZip", "None",
							}, false),
						},
						tfconstants.AttrCompressionLevel: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "compression level (Low, Medium, High)",
							ValidateFunc: validation.StringInSlice([]string{
								"Low", "Medium", "High",
							}, false),
						},
						tfconstants.AttrEncryptionEnabled: {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "whether to encrypt backups",
						},
						tfconstants.AttrEncryptionKey: {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "encryption passphrase for backups",
						},
						tfconstants.AttrHoursOfDay: {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "hours of day to run backups (0-23)",
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice(
									[]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
										"10", "11", "12", "13", "14", "15", "16", "17", "18", "19",
										"20", "21", "22", "23"}, false),
							},
						},
						tfconstants.AttrStartingMinute: {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "starting minute of the hour (0-59)",
							ValidateFunc: validation.IntBetween(0, 59),
						},
						tfconstants.AttrDBEnabled: {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "whether to enable database backup",
						},
						tfconstants.AttrDBUsername: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "database username for backup",
						},
						tfconstants.AttrDBPassword: {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "database password for backup",
						},
					},
				},
			},

			tfconstants.AttrBackupStatus: {
				Type:        schema.TypeList,
				Computed:    true,
				MaxItems:    1,
				Description: "current backup status of the Node",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						tfconstants.AttrStatus: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "current backup status",
						},
						tfconstants.AttrBackupStatusDetail: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "detailed status message",
						},
						tfconstants.AttrLastRecoveryPoint: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "timestamp of last recovery point",
						},
					},
				},
			},
		},
		CreateContext: resourceCreateNode,
		ReadContext:   resourceReadNode,
		UpdateContext: resourceUpdateNode,
		DeleteContext: resourceDeleteNode,
		Exists:        resourceExistsNode,
		Importer: &schema.ResourceImporter{
			State: CustomImportStateFunc,
		},
		CustomizeDiff: resourceNodeCustomizeDiff,
	}
}

func ValidateName(v interface{}, k string) (ws []string, es []error) {

	var errs []error
	var warns []string
	value, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected name to be string"))
		return warns, errs
	}
	if len(value) == 0 {
		errs = append(errs, fmt.Errorf("name cannot be empty"))
		return warns, errs
	}
	validNameRegexp := regexp.MustCompile(`^[a-zA-Z0-9-_]{1,50}$`)
	if !validNameRegexp.Match([]byte(value)) {
		errs = append(errs, fmt.Errorf("the name field cannot be blank, must not contain whitespace or special characters, and must be between 1 and 50 characters in length. Got %s", value))
		return warns, errs
	}
	return warns, errs
}

func resourceNodeCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
	// Emit warnings for deprecated fields
	if _, ok := d.GetOk(tfconstants.AttrSSHKeys); ok {
		log.Printf("[WARN] The ssh_keys field is deprecated and will be removed in v4.0. Please use ssh_key_ids instead.")
	}

	if _, ok := d.GetOk(tfconstants.AttrReserveIP); ok {
		log.Printf("[WARN] The reserve_ip field is deprecated and will be removed in v4.0. Please use reserve_ip_id instead.")
	}

	if blockStorageIDs, ok := d.GetOk("block_storage_ids"); ok {
		if len(blockStorageIDs.([]interface{})) > 0 {
			log.Printf("[WARN] The block_storage_ids field is deprecated and will be removed in v4.0. Please use the e2e_volume_attachment resource instead.")
		}
	}

	// Validate C2 plan restrictions at plan time
	if plan, ok := d.GetOk(tfconstants.AttrPlan); ok {
		planStr := plan.(string)
		if len(planStr) >= 2 && planStr[0:2] == tfconstants.PREFIX_C2_NODE {
			// Check if trying to attach block storage to C2 node
			if blockStorageIDs, ok := d.GetOk("block_storage_ids"); ok {
				if len(blockStorageIDs.([]interface{})) > 0 {
					return fmt.Errorf("cannot attach block storage to C2 plan nodes: C2 plans do not support block storage attachment")
				}
			}
		}
	}

	return nil
}

func resourceCreateNode(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	var diags diag.Diagnostics

	// Get region with provider default support
	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get projectID with provider default support
	projectID, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Handle SSH keys - support both old (ssh_keys) and new (ssh_key_ids) fields
	var sshKeysForAPI []interface{}
	var originalSSHKeys interface{} // Store original for state restore

	if sshKeyIDs, ok := d.GetOk(tfconstants.AttrSSHKeyIDs); ok {
		// New way: using ssh_key_ids (direct resource IDs)
		// IMPORTANT: API expects actual SSH key content strings, not IDs!
		originalSSHKeys = sshKeyIDs
		log.Printf("[INFO] Using ssh_key_ids: %v", sshKeyIDs)

		// Convert resource IDs to actual SSH key content
		sshKeyContents, Err := convertIDsToSshKeyContent(m, sshKeyIDs.([]interface{}), projectID, region)
		if Err != nil {
			return Err
		}
		sshKeysForAPI = sshKeyContents
		log.Printf("[INFO] Converted %d ssh_key_ids to key content for API", len(sshKeysForAPI))
	} else if sshKeyLabels, ok := d.GetOk(tfconstants.AttrSSHKeys); ok {
		// Old way: using ssh_keys (labels) - DEPRECATED
		log.Printf("[WARN] ssh_keys field is deprecated. Use ssh_key_ids instead. This field will be removed in v4.0")
		originalSSHKeys = sshKeyLabels

		new_SSH_keys, Err := convertLabelToSshKey(m, sshKeyLabels.([]interface{}), projectID, region)
		if Err != nil {
			return Err
		}
		d.Set(tfconstants.AttrSSHKeys, new_SSH_keys)
		sshKeysForAPI = new_SSH_keys
	}

	// Handle block storage - DEPRECATED field
	blockStorageIDs := d.Get("block_storage_ids").([]interface{})
	if len(blockStorageIDs) > 0 {
		log.Printf("[WARN] block_storage_ids field is deprecated. Use e2e_volume_attachment resource instead. This field will be removed in v4.0")
	}

	if len(blockStorageIDs) > 1 {
		return diag.Errorf("Only one block storage volume can be attached during node creation")
	}
	image_id := 0
	if len(blockStorageIDs) == 1 {
		// C2 plan check is now in CustomizeDiff for better UX
		image_id_string := blockStorageIDs[0].(string)

		image_id_temp, err := convertStringToInt(image_id_string)
		if err != nil {
			return diag.Errorf("Error converting block storage ID to integer: %s", err)
		}
		image_id = image_id_temp
		Error := checkBlockStorage(m, image_id_string, projectID, region)
		if Error != nil {
			return Error
		}
	}

	log.Printf("[INFO] NODE CREATE STARTS ")
	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return diag.Errorf("Error creating goe2e client: %s", err)
	}

	securityGroups, _, err := goe2eClient.Nodes.GetSecurityGroupList(ctx)
	log.Printf("[INFO] GET Security groups | RESPONSE BODY | %+v", securityGroups)
	if err != nil {
		log.Printf("[ERROR] Error getting Security Group List inside Node Create. Error : %s", err)
		return diag.Errorf("Error retrieving security groups for projectID (%s) and region (%s): %s. Please verify that projectID and region are correct", projectID, region, err)
	}
	defaultSG := getDefaultSGFromList(securityGroups)
	d.Set("default_sg", defaultSG)

	security_group := defaultSG
	if securityGroupsList, ok := d.GetOk(tfconstants.AttrSecurityGroupIDs); ok {
		if securityGroupsList != nil {
			if securityGroups, ok := securityGroupsList.([]interface{}); ok && len(securityGroups) > 0 {
				security_group = securityGroups[0].(int)
				if len(securityGroups) > 1 {
					log.Printf("Can only attach a single security group while node creation. Only the first Security Group will be attached")
					d.Set(tfconstants.AttrSecurityGroupIDs, []int{security_group})
				}
			}
		}
	}

	// Handle reserve_ip - support both old and new fields
	reserveIP := ""
	if reserveIPID, ok := d.GetOk(tfconstants.AttrReserveIPID); ok {
		// New way: using reserve_ip_id
		reserveIP = reserveIPID.(string)
	} else if oldReserveIP, ok := d.GetOk(tfconstants.AttrReserveIP); ok {
		// Old way: using reserve_ip - DEPRECATED
		log.Printf("[WARN] reserve_ip field is deprecated. Use reserve_ip_id instead. This field will be removed in v4.0")
		reserveIP = oldReserveIP.(string)
	}

	// Parse network_interface block if provided (V3 structured config)
	// This can override individual fields like vpc_id, security_group_ids, etc.
	var niVPCID string
	var niAssignPublicIP, niEnableIPv6 bool
	var niSecurityGroupIDs []int
	if niList, ok := d.GetOk(tfconstants.AttrNetworkInterface); ok {
		niVPCID, niAssignPublicIP, niEnableIPv6, niSecurityGroupIDs = expandNetworkInterface(niList.([]interface{}))
		log.Printf("[INFO] Using network_interface block: vpc_id=%s, assign_public_ip=%v, enable_ipv6=%v, security_groups=%v",
			niVPCID, niAssignPublicIP, niEnableIPv6, niSecurityGroupIDs)

		// Override individual fields if network_interface is provided
		// niVPCID will be used below in node struct
		// niEnableIPv6 will override is_ipv6_availed below
		// Security groups handled separately below
	}

	// Parse root_disk block if provided (V3 structured config - maps to API 'disk' field)
	var rootDiskSize int
	if rdList, ok := d.GetOk(tfconstants.AttrRootDisk); ok {
		sizeGB, diskType := expandRootDisk(rdList.([]interface{}))
		if sizeGB > 0 {
			rootDiskSize = sizeGB
			log.Printf("[INFO] Using root_disk block: size_gb=%d, disk_type=%s", sizeGB, diskType)
		}
	}

	// Use network_interface values if provided, otherwise use individual fields
	vpcID := d.Get(tfconstants.AttrVPCID).(string)
	if niVPCID != "" {
		vpcID = niVPCID
	}

	enableIPv6 := d.Get("is_ipv6_availed").(bool)
	if niEnableIPv6 {
		enableIPv6 = niEnableIPv6
	}

	// Convert SSH keys from []interface{} to []string
	sshKeysStr := make([]string, 0, len(sshKeysForAPI))
	for _, key := range sshKeysForAPI {
		if keyStr, ok := key.(string); ok {
			sshKeysStr = append(sshKeysStr, keyStr)
		}
	}

	// Convert start_scripts to single string (goe2e expects string, not array)
	startScript := ""
	if startScripts := GetStartScripts(d.Get("start_script").(string)); len(startScripts) > 0 {
		if scriptStr, ok := startScripts[0].(string); ok {
			startScript = scriptStr
		}
	}

	// Convert models.NodeCreate to goe2e.NodeCreateRequest
	createReq := &goe2e.NodeCreateRequest{
		Name:                 d.Get("name").(string),
		Label:                d.Get(tfconstants.AttrLabel).(string),
		Plan:                 d.Get("plan").(string),
		Image:                d.Get(tfconstants.AttrImage).(string),
		Backup:               d.Get("backup").(bool),
		DefaultPublicIP:      d.Get("default_public_ip").(bool),
		DisablePassword:      d.Get("disable_password").(bool),
		EnableBitNinja:       d.Get("enable_bitninja").(bool),
		EnableIPv6:           enableIPv6,
		IsSavedImage:         d.Get("is_saved_image").(bool),
		ReserveIP:            reserveIP,
		VPCID:                vpcID,
		SavedImageTemplateID: d.Get("saved_image_template_id").(int),
		SecurityGroupID:      security_group,
		SSHKeys:              sshKeysStr,
		StartScript:          startScript,
		ImageID:              image_id,
		Disk:                 rootDiskSize,
	}

	if vpcID != "" {
		vpc, _, err := goe2eClient.Vpcs.GetVPC(ctx, vpcID)
		if err != nil {
			return diag.Errorf("Error retrieving VPC (ID: %s) in project (%s), region (%s): %s", vpcID, projectID, region, err)
		}
		if vpc.State != "Active" {
			return diag.Errorf("Cannot create node: VPC (ID: %s) is in %s state in project (%s), region (%s) (expected: Active)", vpcID, vpc.State, projectID, region)
		}
	}

	createdNode, resp, err := goe2eClient.Nodes.CreateNode(ctx, createReq)
	if err != nil {
		return diag.Errorf("Error creating node (name: %s) in project (%s), region (%s): %s", createReq.Name, projectID, region, err)
	}

	log.Printf("[INFO] NODE CREATE | RESPONSE | %+v", resp)

	// Note: goe2e client doesn't expose is_credit_sufficient in the response struct
	// The API will return an error if credits are insufficient, so we rely on error handling above

	// Get the created node to populate all fields
	// The CreateNode response might not have all fields, so we fetch it
	if createdNode != nil && createdNode.ID != "" {
		d.SetId(createdNode.ID)
		// Fetch full node details
		fullNode, _, err := goe2eClient.Nodes.GetNode(ctx, createdNode.ID)
		if err == nil && fullNode != nil {
			d.Set("is_active", fullNode.IsActive)
			d.Set(tfconstants.AttrCreatedAt, fullNode.CreatedAt)
			d.Set(tfconstants.AttrMemory, fullNode.Memory)
			d.Set(tfconstants.AttrStatus, fullNode.Status)
			d.Set(tfconstants.AttrDisk, fullNode.Disk)
			d.Set("price", fullNode.Price)
			if fullNode.VMID > 0 {
				d.Set(tfconstants.AttrVMID, fullNode.VMID)
			}
		}
	} else {
		// Fallback: try to extract ID from response if available
		// This is a workaround if CreateNode doesn't return the ID
		return diag.Errorf("Error creating node: did not receive node ID in response")
	}

	// Restore original SSH keys to state (either ssh_keys or ssh_key_ids)
	if originalSSHKeys != nil {
		if _, ok := d.GetOk(tfconstants.AttrSSHKeyIDs); ok {
			d.Set(tfconstants.AttrSSHKeyIDs, originalSSHKeys)
		} else {
			d.Set(tfconstants.AttrSSHKeys, originalSSHKeys)
		}
	}
	return diags
}

func resourceReadNode(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	cfg := m.(*config.Config)
	var diags diag.Diagnostics
	copy_ssh_keys := d.Get("ssh_keys")
	log.Printf("[info] inside node Resource read")
	nodeId := d.Id()

	// Get region with provider default support
	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get projectID with provider default support
	projectID, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return diag.Errorf("Error creating goe2e client: %s", err)
	}

	node, _, err := goe2eClient.Nodes.GetNode(ctx, nodeId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
			return diags
		} else {
			return diag.Errorf("Error retrieving node (ID: %s) in project (%s), region (%s): %s", nodeId, projectID, region, err)
		}
	}
	if node == nil {
		d.SetId("")
		return diags
	}
	log.Printf("[info] node Resource read | before setting data")
	log.Printf("[info] node Resource read | data = %+v", node)
	d.Set("name", node.Name)
	d.Set(tfconstants.AttrLabel, node.Label)
	d.Set(tfconstants.AttrPlan, node.Plan)
	d.Set(tfconstants.AttrCreatedAt, node.CreatedAt)
	d.Set(tfconstants.AttrMemory, node.Memory)
	d.Set(tfconstants.AttrStatus, node.Status)
	d.Set(tfconstants.AttrDisk, node.Disk)
	d.Set("price", node.Price)
	d.Set(tfconstants.AttrLockNode, node.IsLocked)
	d.Set(tfconstants.AttrPublicIPAddress, node.PublicIPAddress)
	d.Set(tfconstants.AttrPrivateIPAddress, node.PrivateIPAddress)

	// Set IPv6 address if available
	if node.IPv6Address != "" {
		d.Set(tfconstants.AttrIPv6Address, node.IPv6Address)
	}

	d.Set("is_bitninja_license_active", node.BitNinjaLicenseActive)

	// Preserve SSH keys in state (don't overwrite with API response)
	// The API returns the actual keys, but we want to keep the original reference (labels or IDs)
	if _, ok := d.GetOk(tfconstants.AttrSSHKeyIDs); ok {
		// Using ssh_key_ids - keep as-is
	} else {
		// Using ssh_keys - keep original labels
		d.Set(tfconstants.AttrSSHKeys, copy_ssh_keys)
	}

	// Convert VMID from string to int if needed
	if node.VMID > 0 {
		d.Set(tfconstants.AttrVMID, node.VMID)
	}

	log.Printf("[info] node Resource read | after setting data")
	if node.Status == goe2econstants.NodeStatusRunning || node.Status == goe2econstants.NodeStatusCreating {
		d.Set(tfconstants.AttrPowerStatus, goe2econstants.NodePowerStatusOn)
	}
	if node.Status == goe2econstants.NodeStatusPoweredOff {
		d.Set(tfconstants.AttrPowerStatus, goe2econstants.NodePowerStatusOff)
	}
	securityGroups, _, err := goe2eClient.Nodes.GetSecurityGroupList(ctx)
	if err != nil {
		log.Printf("[ERROR] Error getting Security Group List inside Node Read. Error : %s", err)
		return diag.Errorf("Error retrieving security groups for projectID (%s) and region (%s): %s. Please verify these values are correct", projectID, region, err)
	}
	defaultSG := getDefaultSGFromList(securityGroups)
	d.Set("default_sg", defaultSG)

	return diags

}

func resourceUpdateNode(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	cfg := m.(*config.Config)
	nodeId := d.Id()

	// Get region with provider default support
	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get projectID with provider default support
	projectID, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return diag.Errorf("Error creating goe2e client: %s", err)
	}

	status := d.Get(tfconstants.AttrStatus).(string)
	if status == goe2econstants.NodeStatusFailed {
		rollbackChanges(d)
		return diag.Errorf("Cannot update node (ID: %s): node is in %s state in project (%s), region (%s). Please contact support at cloud-platform@e2enetworks.com", nodeId, status, projectID, region)
	}
	_, _, err = goe2eClient.Nodes.GetNode(ctx, nodeId)
	if err != nil {
		return diag.Errorf("Error retrieving node (ID: %s) in project (%s), region (%s): %s", nodeId, projectID, region, err)
	}

	if d.HasChange("start_script") {
		start_script, _ := d.GetChange("start_script")
		d.Set("region", start_script)
		return diag.Errorf("Cannot update start_script: this field is immutable after node creation")
	}

	if d.HasChange("name") {
		log.Printf("[INFO] nodeId = %v, name = %s ", d.Id(), d.Get("name").(string))
		updateReq := &goe2e.NodeUpdateRequest{
			Name: d.Get("name").(string),
		}
		_, _, err := goe2eClient.Nodes.UpdateNode(ctx, nodeId, updateReq)
		if err != nil {
			return diag.Errorf("Error renaming node (ID: %s) in project (%s), region (%s): %s", nodeId, projectID, region, err)
		}
	}

	if d.HasChange(tfconstants.AttrPowerStatus) {
		nodestatus := d.Get(tfconstants.AttrStatus).(string)
		if nodestatus == goe2econstants.NodeStatusCreating || nodestatus == goe2econstants.NodeStatusReinstalling {
			prevBlockIDArray, _ := d.GetChange("block_storage_ids")
			d.Set("block_storage_ids", prevBlockIDArray)
			return diag.Errorf("Cannot change power status for node (ID: %s): node is in %s state in project (%s), region (%s)", nodeId, nodestatus, projectID, region)
		}
		if d.Get(tfconstants.AttrLockNode).(bool) {
			return diag.Errorf("Cannot change power status for node (ID: %s): node is locked", nodeId)
		}
		log.Printf("[INFO] %s ", d.Get(tfconstants.AttrPowerStatus).(string))
		powerStatus := d.Get(tfconstants.AttrPowerStatus).(string)
		if powerStatus == goe2econstants.NodePowerStatusOn {
			_, err = goe2eClient.Nodes.PowerOn(ctx, nodeId)
		} else if powerStatus == goe2econstants.NodePowerStatusOff {
			_, err = goe2eClient.Nodes.PowerOff(ctx, nodeId)
		} else {
			return diag.Errorf("Invalid power status: %s", powerStatus)
		}
		if err != nil {
			return diag.Errorf("Error changing power status for node (ID: %s) in project (%s), region (%s): %s", nodeId, projectID, region, err)
		}
	}

	if d.HasChange(tfconstants.AttrLockNode) {
		if d.Get(tfconstants.AttrStatus).(string) == goe2econstants.NodeStatusCreating || d.Get(tfconstants.AttrStatus).(string) == goe2econstants.NodeStatusReinstalling {
			return diag.Errorf("Cannot update node (ID: %s): node is in %s state in project (%s), region (%s)", nodeId, d.Get(tfconstants.AttrStatus).(string), projectID, region)
		}
		if d.Get(tfconstants.AttrLockNode).(bool) {
			_, err := goe2eClient.Nodes.LockNode(ctx, nodeId)
			if err != nil {
				return diag.Errorf("Error locking node (ID: %s) in project (%s), region (%s): %s", nodeId, projectID, region, err)
			}
		}
		if !d.Get(tfconstants.AttrLockNode).(bool) {
			_, err := goe2eClient.Nodes.UnlockNode(ctx, nodeId)
			if err != nil {
				return diag.Errorf("Error unlocking node (ID: %s) in project (%s), region (%s): %s", nodeId, projectID, region, err)
			}
		}
	}

	if d.HasChange(tfconstants.AttrRebootNode) {

		if d.Get(tfconstants.AttrRebootNode).(bool) {
			d.Set(tfconstants.AttrRebootNode, false)
			if d.Get(tfconstants.AttrStatus).(string) == goe2econstants.NodeStatusCreating || d.Get(tfconstants.AttrStatus).(string) == goe2econstants.NodeStatusReinstalling {
				return diag.Errorf("Cannot reboot node (ID: %s): node is in %s state in project (%s), region (%s)", nodeId, d.Get(tfconstants.AttrStatus).(string), projectID, region)
			}
			if d.Get(tfconstants.AttrStatus).(string) == goe2econstants.NodeStatusPoweredOff {
				return diag.Errorf("Cannot reboot node (ID: %s): node must be powered on (current state: %s)", nodeId, d.Get(tfconstants.AttrStatus).(string))
			}
			_, err := goe2eClient.Nodes.Reboot(ctx, nodeId)
			if err != nil {
				return diag.Errorf("Error rebooting node (ID: %s) in project (%s), region (%s): %s", nodeId, projectID, region, err)
			}
		}
	}
	if d.HasChange("reinstall_node") {
		if d.Get(tfconstants.AttrStatus).(string) == goe2econstants.NodeStatusCreating {
			return diag.Errorf("Cannot reinstall node (ID: %s): node is in creating state", nodeId)
		}
		if d.Get(tfconstants.AttrStatus).(string) == goe2econstants.NodeStatusReinstalling {
			return diag.Errorf("Cannot reinstall node (ID: %s): node is already being reinstalled", nodeId)
		}
		if d.Get("reinstall_node").(bool) {
			if d.Get(tfconstants.AttrStatus).(string) == goe2econstants.NodeStatusPoweredOff {
				d.Set("reinstall_node", false)
				return diag.Errorf("Cannot reinstall node (ID: %s): node must be powered on (current state: %s)", nodeId, goe2econstants.NodeStatusPoweredOff)
			}
			if d.Get(tfconstants.AttrStatus).(string) == goe2econstants.NodeStatusReinstalling {
				d.Set("reinstall_node", false)
				return diag.Errorf("Cannot reinstall node (ID: %s): node is already being reinstalled", nodeId)
			}
			reinstallReq := &goe2e.NodeReinstallRequest{
				Image: d.Get(tfconstants.AttrImage).(string),
			}
			_, err := goe2eClient.Nodes.Reinstall(ctx, nodeId, reinstallReq)
			d.Set("reinstall_node", false)
			if err != nil {
				return diag.Errorf("Error reinstalling node (ID: %s) in project (%s), region (%s): %s", nodeId, projectID, region, err)
			}
		}
	}

	if d.HasChange("save_image") {
		if d.Get("save_image") == true {
			d.Set("save_image", false)
			if d.Get("save_image_name").(string) == "" {
				return diag.Errorf("save_image_name is required when save_image is enabled")
			}

			saveReq := &goe2e.NodeSaveImageRequest{
				ActionType: "save_images",
				Name:       d.Get("save_image_name").(string),
			}
			_, _, err := goe2eClient.Nodes.SaveImage(ctx, nodeId, saveReq)
			if err != nil {
				return diag.Errorf("Error saving image for node (ID: %s) in project (%s), region (%s): %s", nodeId, projectID, region, err)
			}
		}
	}

	if d.HasChange(tfconstants.AttrSecurityGroupIDs) {
		oldSGData, newSGData := d.GetChange(tfconstants.AttrSecurityGroupIDs)
		if d.Get(tfconstants.AttrStatus).(string) != goe2econstants.NodeStatusRunning {
			d.Set(tfconstants.AttrSecurityGroupIDs, oldSGData)
			return diag.Errorf("Cannot update security groups for node (ID: %s): node must be in running state (current state: %s)", nodeId, d.Get(tfconstants.AttrStatus).(string))
		}
		security_groups_list := d.Get(tfconstants.AttrSecurityGroupIDs).([]interface{})

		if len(security_groups_list) <= 0 {
			d.Set(tfconstants.AttrSecurityGroupIDs, oldSGData)
			return diag.Errorf("At least one security group must be attached to node (ID: %s)", nodeId)
		}
		oldSGList := oldSGData.([]interface{})
		newSGList := newSGData.([]interface{})
		sgMap := make(map[int]int)
		for _, sgID := range newSGList {
			sgMap[sgID.(int)] = 1
		}
		for _, sgID := range oldSGList {
			if count, ok := sgMap[sgID.(int)]; ok {
				sgMap[sgID.(int)] = count - 1
			} else {
				sgMap[sgID.(int)] = -1
			}
		}
		var toBeAttached []int
		for key, value := range sgMap {
			if value == -1 {
				log.Printf("----------HAVE TO DETACH THE SECURITY GROUP WITH ID %+v ------------------", key)
				sgReq := &goe2e.SecurityGroupRequest{
					SecurityGroupList: []int{key},
				}
				_, err := goe2eClient.Nodes.DetachSecurityGroup(ctx, nodeId, sgReq)
				if err != nil {
					return diag.Errorf("Error detaching security group (ID: %d) from node (ID: %s) in project (%s), region (%s): %s", key, nodeId, projectID, region, err)
				}
				continue
			}
			if value >= 1 {
				toBeAttached = append(toBeAttached, key)
			}
		}
		if len(toBeAttached) >= 1 {
			sgReq := &goe2e.SecurityGroupRequest{
				SecurityGroupList: toBeAttached,
			}
			_, err := goe2eClient.Nodes.AttachSecurityGroup(ctx, nodeId, sgReq)
			if err != nil {
				return diag.Errorf("Error attaching security groups to node (ID: %s) in project (%s), region (%s): %s", nodeId, projectID, region, err)
			}
		}
	}

	if d.HasChange(tfconstants.AttrLabel) {
		log.Printf("[INFO] nodeId = %v changed label = %s ", d.Id(), d.Get(tfconstants.AttrLabel).(string))
		updateReq := &goe2e.NodeUpdateRequest{
			Label: d.Get(tfconstants.AttrLabel).(string),
		}
		_, _, err = goe2eClient.Nodes.UpdateNode(ctx, nodeId, updateReq)
		if err != nil {
			return diag.Errorf("Error updating label for node (ID: %s) in project (%s), region (%s): %s", nodeId, projectID, region, err)
		}
	}

	if d.HasChange("ssh_keys") {
		prevSshKeys, currSshKeys := d.GetChange("ssh_keys")

		log.Printf("[INFO] nodeId = %v changed ssh_keys = %s ", d.Id(), d.Get("ssh_keys"))
		log.Printf("[INFO] type of ssh_keys data = %T", d.Get("ssh_keys"))

		new_SSH_keys, Err := convertLabelToSshKey(m, d.Get("ssh_keys").([]interface{}), projectID, d.Get("region").(string))
		if Err != nil {
			d.Set("ssh_keys", prevSshKeys)
			return Err
		}
		d.Set("ssh_keys", new_SSH_keys)
		// Use goe2e client for SSH key updates
		sshKeys := goe2e.GenerateSSHKeyMap(d.Get("ssh_keys").([]interface{}))
		sshReq := &goe2e.SSHUpdateRequest{
			Action:  "add_ssh_keys",
			SSHKeys: sshKeys,
		}
		_, err = goe2eClient.Nodes.UpdateSSH(ctx, nodeId, sshReq)
		d.Set("ssh_keys", currSshKeys)
		if err != nil {
			d.Set("ssh_keys", prevSshKeys)
			return diag.Errorf("Error updating SSH keys for node (ID: %s) in project (%s), region (%s): %s", nodeId, projectID, region, err)
		}

	}
	if d.HasChange("region") {
		prevLocation, currLocation := d.GetChange("region")
		log.Printf("[INFO] prevLocation %s, currLocation %s", prevLocation.(string), currLocation.(string))
		d.Set("region", prevLocation)
		return diag.Errorf("Cannot update region: this field is immutable after node creation")
	}
	if d.HasChange(tfconstants.AttrImage) {
		prevImage, currImage := d.GetChange(tfconstants.AttrImage)
		log.Printf("[INFO] prevImage %s, currImage %s", prevImage.(string), currImage.(string))
		d.Set(tfconstants.AttrImage, prevImage.(string))
		return diag.Errorf("Cannot update image: this field is immutable after node creation")
	}
	if d.HasChange(tfconstants.AttrPlan) {
		prevPlan, currPlan := d.GetChange(tfconstants.AttrPlan)

		if d.HasChange(tfconstants.AttrPowerStatus) {
			_ = util.WaitForNodePowerState(m, nodeId, projectID, region)
		}

		log.Printf("[INFO] prevPlan %s, currPlan %s", prevPlan.(string), currPlan.(string))

		if d.Get(tfconstants.AttrStatus).(string) != goe2econstants.NodeStatusPoweredOff {
			d.Set(tfconstants.AttrPlan, prevPlan)
			return diag.Errorf("Cannot upgrade plan for node (ID: %s): node must be powered off (current state: %s)", nodeId, d.Get(tfconstants.AttrStatus).(string))
		}
		upgradeReq := &goe2e.NodePlanUpgradeRequest{
			Plan:  d.Get(tfconstants.AttrPlan).(string),
			Image: d.Get(tfconstants.AttrImage).(string),
		}
		_, err = goe2eClient.Nodes.UpgradePlan(ctx, nodeId, upgradeReq)

		if err != nil {
			d.Set(tfconstants.AttrPlan, prevPlan)
			return diag.Errorf("Error upgrading plan for node (ID: %s) in project (%s), region (%s): %s", nodeId, projectID, region, err)
		}
	}

	if d.HasChange("block_storage_ids") {

		log.Printf("[INFO] Power_status changing is = %v", d.HasChange(tfconstants.AttrPowerStatus))
		if d.HasChange(tfconstants.AttrPowerStatus) {
			err := util.WaitForNodePowerState(m, nodeId, projectID, region)
			if err != nil {
				return diag.Errorf("Error waiting for node (ID: %s) power state change in project (%s), region (%s): %s", nodeId, projectID, region, err)
			}
		}

		prevBlockIDArray, currBlockIDArray := d.GetChange("block_storage_ids")

		if d.Get(tfconstants.AttrPlan).(string)[0:2] == tfconstants.PREFIX_C2_NODE {
			d.Set("block_storage_ids", prevBlockIDArray)
			return diag.Errorf("Cannot attach block storage to node (ID: %s): C2 plan nodes do not support block storage attachment", nodeId)
		}

		detachingIDs := UniqueArrayElements(prevBlockIDArray.([]interface{}), currBlockIDArray.([]interface{}))
		attachingIDs := UniqueArrayElements(currBlockIDArray.([]interface{}), prevBlockIDArray.([]interface{}))
		CommonIDs := prevBlockIDArray.([]interface{})
		log.Printf("[INFO] detachingIDs %+v, attachingIDs %+v, CommonIDs %+v", detachingIDs, attachingIDs, CommonIDs)
		log.Printf("[INFO] prevIDArray %v, currIDArray %v", prevBlockIDArray, currBlockIDArray)

		// Get goe2e client for volume attachment operations
		goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
		if err != nil {
			d.Set("block_storage_ids", prevBlockIDArray)
			return diag.Errorf("Error creating goe2e client: %s", err)
		}

		for i, detachingID := range detachingIDs {
			blockStorageID := detachingID.(string)
			detachReq := &goe2e.VolumeDetachmentRequest{
				NodeID:   nodeId,
				VolumeID: blockStorageID,
			}
			_, err := goe2eClient.VolumeAttachment.DetachVolume(ctx, detachReq)
			if err != nil {
				d.Set("block_storage_ids", CommonIDs)
				return diag.Errorf("Error detaching block storage (ID: %s) from node (ID: %s) in project (%s), region (%s): %s", blockStorageID, nodeId, projectID, region, err)
			}
			CommonIDs = removeArrayElement(CommonIDs, detachingID)
			// Wait for some time before detaching the next block storage
			if err := util.WaitForNodeLCMState(m, nodeId, projectID, region); err != nil {
				return diag.FromErr(err)
			}
			if i == len(detachingIDs)-1 {
				break
			}
		}
		for i, attachingID := range attachingIDs {
			blockStorageID := attachingID.(string)
			Error := checkBlockStorage(m, blockStorageID, projectID, region)
			if Error != nil {
				d.Set("block_storage_ids", CommonIDs)
				log.Printf("[ERROR] Error attaching block storage CommonIDs = %+v", CommonIDs)
				return Error
			}
			attachReq := &goe2e.VolumeAttachmentRequest{
				NodeID:   nodeId,
				VolumeID: blockStorageID,
			}
			_, _, err := goe2eClient.VolumeAttachment.AttachVolume(ctx, attachReq)
			if err != nil {
				d.Set("block_storage_ids", CommonIDs)
				log.Printf("[ERROR] Error attaching block storage CommonIDs = %+v", CommonIDs)
				return diag.Errorf("Error attaching block storage (ID: %s) to node (ID: %s) in project (%s), region (%s): %s", blockStorageID, nodeId, projectID, region, err)
			}
			CommonIDs = append(CommonIDs, attachingID)
			// Wait for some time before attaching the next block storage
			if i == len(attachingIDs)-1 {
				break
			}
			if err := util.WaitForNodeLCMState(m, nodeId, projectID, region); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	// Handle tag_ids changes (attach/detach tags via label API)
	if d.HasChange(tfconstants.AttrTagIDs) {
		oldTagIDs, newTagIDs := d.GetChange(tfconstants.AttrTagIDs)
		oldList := oldTagIDs.([]interface{})
		newList := newTagIDs.([]interface{})

		// Calculate tags to attach and detach
		var toAttach, toDetach []int
		oldMap := make(map[int]bool)
		newMap := make(map[int]bool)

		for _, id := range oldList {
			oldMap[id.(int)] = true
		}
		for _, id := range newList {
			newMap[id.(int)] = true
		}

		// Find tags to detach (in old but not in new)
		for _, id := range oldList {
			tagID := id.(int)
			if !newMap[tagID] {
				toDetach = append(toDetach, tagID)
			}
		}

		// Find tags to attach (in new but not in old)
		for _, id := range newList {
			tagID := id.(int)
			if !oldMap[tagID] {
				toAttach = append(toAttach, tagID)
			}
		}

		log.Printf("[INFO] Tag changes for node (ID: %s): attaching %v, detaching %v", nodeId, toAttach, toDetach)

		// TODO: Implement tag attach/detach API calls
		// This will be implemented when e2e_tag resource is created
		// API endpoint: PUT /label/mapping/nodes/{node_id}/
		// Payload: {"attach": [tag_ids], "detach": [tag_ids]}
		if len(toAttach) > 0 || len(toDetach) > 0 {
			log.Printf("[WARN] Tag management API not yet implemented. Tags to attach: %v, detach: %v", toAttach, toDetach)
			// When implemented, call: apiClient.UpdateNodeTags(nodeId, toAttach, toDetach, projectID, region)
		}
	}

	return resourceReadNode(ctx, d, m)

}

func resourceDeleteNode(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	var diags diag.Diagnostics
	nodeId := d.Id()

	// Get region with provider default support
	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get projectID with provider default support
	projectID, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	node_status := d.Get(tfconstants.AttrStatus).(string)
	if node_status == goe2econstants.NodeStatusSaving || node_status == goe2econstants.NodeStatusCreating {
		return diag.Errorf("Cannot delete node (ID: %s): node is in %s state in project (%s), region (%s)", nodeId, node_status, projectID, region)
	}

	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return diag.Errorf("Error creating goe2e client: %s", err)
	}

	_, err = goe2eClient.Nodes.DeleteNode(ctx, nodeId)
	if err != nil {
		return diag.Errorf("Error deleting node (ID: %s) in project (%s), region (%s): %s", nodeId, projectID, region, err)
	}
	d.SetId("")
	return diags
}

func resourceExistsNode(d *schema.ResourceData, m interface{}) (bool, error) {
	cfg := m.(*config.Config)
	nodeId := d.Id()

	// Get region with provider default support
	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return false, err
	}

	// Get projectID with provider default support
	projectID, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return false, err
	}

	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return false, err
	}

	ctx := context.Background()
	_, _, err = goe2eClient.Nodes.GetNode(ctx, nodeId)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func CustomImportStateFunc(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid ID format: expected projectID/region/resource_id")
	}

	projectID := parts[0]
	region := parts[1]
	nodeID := parts[2]

	// Set the individual fields
	d.Set(tfconstants.AttrProjectID, projectID)
	d.Set("region", region)

	// Use node ID as actual Terraform resource ID
	d.SetId(nodeID)

	return []*schema.ResourceData{d}, nil
}

func convertStringToInt(str string) (int, error) {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func getDefaultSGFromList(securityGroups []goe2e.SecurityGroupInfo) int {
	for _, sg := range securityGroups {
		if sg.IsDefault {
			log.Printf("------------Default security group is: %+v -------------", sg.ID)
			return sg.ID
		}
	}
	return 0
}

func UniqueArrayElements(arr1 []interface{}, arr2 []interface{}) []interface{} {
	var res []interface{}
	for _, v := range arr1 {
		if !isContains(arr2, v) {
			res = append(res, v)
		}
	}
	return res
}

func isContains(arr []interface{}, val interface{}) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

func CommonArrayElements(arr1 []interface{}, arr2 []interface{}) []interface{} {
	var res []interface{}
	for _, v := range arr1 {
		if isContains(arr2, v) {
			res = append(res, v)
		}
	}
	return res
}

func removeArrayElement(arr []interface{}, val interface{}) []interface{} {
	var res []interface{}
	for _, v := range arr {
		if v != val {
			res = append(res, v)
		}
	}
	return res
}

func ValidatePlanName(v interface{}, k string) (ws []string, es []error) {

	var errs []error
	var warns []string
	value, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected plan to be string"))
		return warns, errs
	}
	if value == "" {
		errs = append(errs, fmt.Errorf("plan name cannot be empty"))
		return warns, errs
	}

	whiteSpace := regexp.MustCompile(`\s+`)
	if whiteSpace.Match([]byte(value)) {
		errs = append(errs, fmt.Errorf("plan cannot contain whitespace. got %s", value))
		return warns, errs
	}
	return warns, errs
}

func ValidateBlank(v interface{}, k string) (ws []string, es []error) {

	var errs []error
	var warns []string
	value, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected %s to be string", k))
		return warns, errs
	}
	stripped := strings.TrimSpace(value)
	if stripped == "" {
		errs = append(errs, fmt.Errorf("%s cannot be blank", k))
		return warns, errs
	}
	return warns, errs
}

func ValidateInteger(v interface{}, k string) (ws []string, es []error) {
	var errs []error
	var warns []string

	str, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected %s to be string", k))
		return warns, errs
	}
	// validate block storage id ("123" -> correct, "abc" -> incorrect, "123abc" -> incorrect)
	_, err := strconv.Atoi(str)
	if err != nil {
		errs = append(errs, fmt.Errorf("%s only contains numeric value", k))
		return warns, errs
	}
	return warns, errs
}

func rollbackChanges(d *schema.ResourceData) {
	prevImage, _ := d.GetChange(tfconstants.AttrImage)
	prevName, _ := d.GetChange(tfconstants.AttrName)
	prevPlan, _ := d.GetChange(tfconstants.AttrPlan)
	prevLocation, _ := d.GetChange(tfconstants.AttrLocation)
	prevProjectId, _ := d.GetChange(tfconstants.AttrProjectID)
	prevRegion, _ := d.GetChange(tfconstants.AttrRegion)
	prevLabel, _ := d.GetChange(tfconstants.AttrLabel)
	prevBackup, _ := d.GetChange("backup")
	prevDefaultPublicIp, _ := d.GetChange("default_public_ip")
	prevDisablePassword, _ := d.GetChange("disable_password")
	prevEnableBitninja, _ := d.GetChange("enable_bitninja")
	prevIsIpv6Availed, _ := d.GetChange("is_ipv6_availed")
	prevIsSavedImage, _ := d.GetChange("is_saved_image")
	prevReserveIp, _ := d.GetChange("reserve_ip")
	prevSavedImageTemplateId, _ := d.GetChange("saved_image_template_id")
	prevSshKey, _ := d.GetChange("ssh_keys")
	prevVpcId, _ := d.GetChange(tfconstants.AttrVPCID)
	prevBlockStorageIds, _ := d.GetChange("block_storage_ids")
	prevSecurityGroupIds, _ := d.GetChange(tfconstants.AttrSecurityGroupIDs)
	prevLockNode, _ := d.GetChange(tfconstants.AttrLockNode)
	prevPowerStatus, _ := d.GetChange(tfconstants.AttrPowerStatus)
	prevRebootNode, _ := d.GetChange(tfconstants.AttrRebootNode)
	prevReinstallNode, _ := d.GetChange("reinstall_node")

	d.Set(tfconstants.AttrImage, prevImage)
	d.Set(tfconstants.AttrName, prevName)
	d.Set(tfconstants.AttrPlan, prevPlan)
	d.Set(tfconstants.AttrLocation, prevLocation)
	d.Set(tfconstants.AttrProjectID, prevProjectId)
	d.Set(tfconstants.AttrRegion, prevRegion)
	d.Set(tfconstants.AttrLabel, prevLabel)
	d.Set("backup", prevBackup)
	d.Set("default_public_ip", prevDefaultPublicIp)
	d.Set("disable_password", prevDisablePassword)
	d.Set("enable_bitninja", prevEnableBitninja)
	d.Set("is_ipv6_availed", prevIsIpv6Availed)
	d.Set("is_saved_image", prevIsSavedImage)
	d.Set("reserve_ip", prevReserveIp)
	d.Set("saved_image_template_id", prevSavedImageTemplateId)
	d.Set("ssh_keys", prevSshKey)
	d.Set(tfconstants.AttrVPCID, prevVpcId)
	d.Set("block_storage_ids", prevBlockStorageIds)
	d.Set(tfconstants.AttrSecurityGroupIDs, prevSecurityGroupIds)

	d.Set(tfconstants.AttrLockNode, prevLockNode)
	d.Set(tfconstants.AttrPowerStatus, prevPowerStatus)
	d.Set(tfconstants.AttrRebootNode, prevRebootNode)
	d.Set("reinstall_node", prevReinstallNode)
}
