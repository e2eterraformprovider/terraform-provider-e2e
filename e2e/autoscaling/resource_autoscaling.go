package autoscaling

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/node"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	// Autoscaling schema keys (provider-side, not API contract).
	// Keep these file-local to avoid polluting `e2e/constants/attrs.go` with autoscaling-only nested keys.
	attrScheduledAction     = "scheduled_action"
	attrScheduledActionType = "action_type"
	attrScheduledAdjustment = "adjustment"
	attrScheduledTargetCap  = "target_capacity"
)

func addAutoscalingV2DeprecationWarnings(diff *schema.ResourceDiff) {
	// Collect all deprecated V2 fields/blocks used in the config and emit a single, actionable warning.
	// We intentionally emit one warning to reduce noise and make migration steps obvious.

	type mapping struct {
		old string
		new string
	}

	var used []mapping
	addIfUsed := func(oldKey, newKey string) {
		if v, ok := diff.GetOk(oldKey); ok {
			// Consider empty string / empty list as "not used".
			switch vv := v.(type) {
			case string:
				if vv == "" {
					return
				}
			case []interface{}:
				if len(vv) == 0 {
					return
				}
			}
			used = append(used, mapping{old: oldKey, new: newKey})
		}
	}

	addIfUsed("vm_image_name", "image")
	addIfUsed(tfconstants.AttrMinNodes, "min_size")
	addIfUsed(tfconstants.AttrMaxNodes, "max_size")
	addIfUsed(tfconstants.AttrDesired, "desired_capacity")
	addIfUsed("provision_status", "status")
	addIfUsed(tfconstants.AttrIsEncryptionEnabled, "enable_encryption")
	addIfUsed(tfconstants.AttrPublicIPRequired, "assign_public_ip")
	addIfUsed("vpc", "vpc_config (or network_config.vpc_names)")
	addIfUsed("policy", "scaling_policy")
	addIfUsed("scheduled_policy", attrScheduledAction)

	if len(used) == 0 {
		return
	}

	var b strings.Builder
	b.WriteString(WarnDeprecatedV2FieldsHeader)
	b.WriteString("\n\nFields to migrate:\n")
	for _, m := range used {
		b.WriteString(fmt.Sprintf("- %s -> %s\n", m.old, m.new))
	}
	b.WriteString("\n")
	b.WriteString(WarnDeprecatedV2FieldsFooter)

	// terraform-plugin-sdk in this repo does not expose ResourceDiff.AddWarning.
	// We still emit a deprecation warning at plan time via logs to keep the signal
	// close to configuration evaluation (CustomizeDiff).
	log.Printf("[WARN] %s", b.String())
}

func autoscalingV2DeprecationWarningDiagnostic(d *schema.ResourceData) *diag.Diagnostic {
	// Apply-time warning for deprecated V2 field usage.
	// Rationale: SDKv2 does not provide plan warnings from CustomizeDiff, but diag warnings are
	// user-visible during apply (and often during refresh-driven plans).

	type mapping struct {
		old string
		new string
	}

	var used []mapping
	addIfUsed := func(oldKey, newKey string) {
		if v, ok := d.GetOk(oldKey); ok {
			switch vv := v.(type) {
			case string:
				if vv == "" {
					return
				}
			case []interface{}:
				if len(vv) == 0 {
					return
				}
			}
			used = append(used, mapping{old: oldKey, new: newKey})
		}
	}

	addIfUsed("vm_image_name", "image")
	addIfUsed(tfconstants.AttrMinNodes, "min_size")
	addIfUsed(tfconstants.AttrMaxNodes, "max_size")
	addIfUsed(tfconstants.AttrDesired, "desired_capacity")
	addIfUsed("provision_status", "status")
	addIfUsed(tfconstants.AttrIsEncryptionEnabled, "enable_encryption")
	addIfUsed(tfconstants.AttrPublicIPRequired, "assign_public_ip")
	addIfUsed("vpc", "vpc_config (or network_config.vpc_names)")
	addIfUsed("policy", "scaling_policy")
	addIfUsed("scheduled_policy", attrScheduledAction)

	if len(used) == 0 {
		return nil
	}

	var b strings.Builder
	b.WriteString(WarnDeprecatedV2FieldsHeader)
	b.WriteString("\n\nFields to migrate:\n")
	for _, m := range used {
		b.WriteString(fmt.Sprintf("- %s -> %s\n", m.old, m.new))
	}
	b.WriteString("\n")
	b.WriteString(WarnDeprecatedV2FieldsFooter)

	return &diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Deprecated V2 autoscaling fields",
		Detail:   b.String(),
	}
}

func autoscalingScalerGroupStatusIn(status string, allowed []string) bool {
	for _, v := range allowed {
		if status == v {
			return true
		}
	}
	return false
}

func ResourceScalerGroup() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceAutoscalingResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: ResourceAutoscalingStateUpgradeV0toV1,
				Version: 0,
			},
		},
		Schema: map[string]*schema.Schema{
			// Common fields - use constants and helpers
			tfconstants.AttrRegion:    config.RegionSchema(),
			tfconstants.AttrLocation:  config.LocationSchema(),
			tfconstants.AttrProjectID: config.ProjectIDSchemaResource(),

			// Autoscaling-specific fields
			tfconstants.AttrName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "name of the Scaler Group",
			},
			tfconstants.AttrPlan: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "the plan of the Scaler Group",
			},
			"plan_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the internal id of the plan derived from plan",
			},
			"sku_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the SKU id (same as plan_id) used for the Scaler Group",
			},
			"slug_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the slug representation of the plan",
			},

			"vm_image_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Deprecated:    "Use 'image' field instead. This field will be removed in v4.0.",
				ConflictsWith: []string{"image"},
				Description:   "the VM image name for the Scaler Group (deprecated, use 'image' instead). Either 'vm_image_name' or 'image' must be specified.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					trim := func(name string) string {
						if idx := strings.Index(name, "_"); idx != -1 {
							return name[:idx]
						}
						return name
					}
					trimmedOld := trim(old)
					trimmedNew := trim(new)

					log.Printf("[DEBUG] DiffSuppressFunc: old=%s → %s, new=%s → %s", old, trimmedOld, new, trimmedNew)

					return trimmedOld == trimmedNew
				},
			},
			"image": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"vm_image_name"},
				Description:   "the VM image name for the Scaler Group (V3 field name, preferred over 'vm_image_name'). Either 'vm_image_name' or 'image' must be specified.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					trim := func(name string) string {
						if idx := strings.Index(name, "_"); idx != -1 {
							return name[:idx]
						}
						return name
					}
					trimmedOld := trim(old)
					trimmedNew := trim(new)

					log.Printf("[DEBUG] DiffSuppressFunc: old=%s → %s, new=%s → %s", old, trimmedOld, new, trimmedNew)

					return trimmedOld == trimmedNew
				},
			},
			"vm_image_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "id of the VM image",
			},
			"vm_template_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "id of the VM template",
			},
			"running": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "current number of running nodes in the Scaler Group (V2 field)",
			},
			"running_node_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "current number of running nodes in the Scaler Group (V3 field, alias for 'running')",
			},
			"nodes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "list of nodes in the Scaler Group",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "node ID",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "node name",
						},
						"ip": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "list of IP addresses",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"public_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "public IP address",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "node status",
						},
						"cpu_usage": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CPU usage percentage",
						},
					},
				},
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "resource tags (state-only in V3.0, API support pending)",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"my_account_sg_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "id of the Security Group to attach to the Scaler Group (if not provided, a default will be fetched from the API)",
			},
			tfconstants.AttrSecurityGroupIDs: {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "list of Security Group ids currently attached to the Scaler Group",
			},

			tfconstants.AttrIsEncryptionEnabled: {
				Type:          schema.TypeBool,
				Optional:      true,
				Default:       false,
				ForceNew:      true,
				Deprecated:    "Use 'enable_encryption' field instead. This field will be removed in v4.0.",
				ConflictsWith: []string{"enable_encryption"},
				Description:   "whether to enable encryption for the Scaler Group (deprecated, use 'enable_encryption' instead)",
			},
			"enable_encryption": {
				Type:          schema.TypeBool,
				Optional:      true,
				Default:       false,
				ForceNew:      true,
				ConflictsWith: []string{tfconstants.AttrIsEncryptionEnabled},
				Description:   "whether to enable encryption for the Scaler Group (V3 field name, preferred over 'is_encryption_enabled')",
			},
			tfconstants.AttrEncryptionPassphrase: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				ForceNew:    true,
				Sensitive:   true,
				Description: "passphrase for encryption (if enabled)",
			},
			tfconstants.AttrPublicIPRequired: {
				Type:          schema.TypeBool,
				Optional:      true,
				Default:       true,
				Deprecated:    "Use 'assign_public_ip' field instead. This field will be removed in v4.0.",
				ConflictsWith: []string{"assign_public_ip"},
				Description:   "whether to assign a public IP to nodes (deprecated, use 'assign_public_ip' instead). Can only be updated when the Scaler Group is stopped and a VPC is attached.",
			},
			"assign_public_ip": {
				Type:          schema.TypeBool,
				Optional:      true,
				Default:       true,
				ConflictsWith: []string{tfconstants.AttrPublicIPRequired},
				Description:   "whether to assign a public IP to nodes (V3 field name, preferred over 'is_public_ip_required'). Can only be updated when the Scaler Group is stopped and a VPC is attached.",
			},

			"provision_status": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Use 'status' field instead. This field will be removed in v4.0.",
				ConflictsWith: []string{"status"},
				ValidateFunc:  validation.StringInSlice(tfconstants.AutoscalingScalerGroupProvisionStatusAllowed, false),
				Description:   "the provision status of the Scaler Group (deprecated, use 'status' instead). Set to 'Stopped' to stop, or 'Running' to start.",
			},
			"status": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"provision_status"},
				ValidateFunc:  validation.StringInSlice(tfconstants.AutoscalingScalerGroupStatusAllowed, false),
				Description:   "the status of the Scaler Group (V3 field name, preferred over 'provision_status'). Set to 'stopped' to stop, or 'running' to start.",
			},

			tfconstants.AttrMinNodes: {
				Type:          schema.TypeInt,
				Optional:      true,
				Deprecated:    "Use 'min_size' field instead. This field will be removed in v4.0.",
				ConflictsWith: []string{"min_size"},
				Description:   "the minimum number of nodes in the Scaler Group (deprecated, use 'min_size' instead)",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 1 {
						errs = append(errs, fmt.Errorf("%q must be at least 1, got: %d", key, v))
					}
					return
				},
			},
			"min_size": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{tfconstants.AttrMinNodes},
				Description:   "the minimum number of nodes in the Scaler Group (V3 field name, preferred over 'min_nodes')",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 1 {
						errs = append(errs, fmt.Errorf("%q must be at least 1, got: %d", key, v))
					}
					return
				},
			},

			tfconstants.AttrMaxNodes: {
				Type:          schema.TypeInt,
				Optional:      true,
				Deprecated:    "Use 'max_size' field instead. This field will be removed in v4.0.",
				ConflictsWith: []string{"max_size"},
				Description:   "the maximum number of nodes in the Scaler Group (deprecated, use 'max_size' instead)",
			},
			"max_size": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{tfconstants.AttrMaxNodes},
				Description:   "the maximum number of nodes in the Scaler Group (V3 field name, preferred over 'max_nodes')",
			},
			tfconstants.AttrDesired: {
				Type:          schema.TypeInt,
				Optional:      true,
				Deprecated:    "Use 'desired_capacity' field instead. This field will be removed in v4.0.",
				ConflictsWith: []string{"desired_capacity"},
				Description:   "the desired number of nodes in the Scaler Group (deprecated, use 'desired_capacity' instead)",
			},
			"desired_capacity": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{tfconstants.AttrDesired},
				Description:   "the desired number of nodes in the Scaler Group (V3 field name, preferred over 'desired')",
			},
			"policy_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "the policy type for the Scaler Group",
			},
			"vpc": {
				Type:          schema.TypeList,
				Optional:      true,
				Deprecated:    "Use 'vpc_config' block instead. This field will be removed in v4.0.",
				ConflictsWith: []string{"vpc_config"},
				Description:   "list of VPCs attached to the Scaler Group (deprecated, use 'vpc_config' block instead)",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"network_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ipv4_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnets": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"subnet_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cidr": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"used_ips": {
										Type:     schema.TypeInt,
										Computed: true,
									},

									"total_ips": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"vpc_config": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"vpc"},
				Description:   "list of VPCs to attach to the Scaler Group (V3 structured block, preferred over 'vpc')",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "name of the VPC to attach",
						},
						"network_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "network ID of the VPC (computed)",
						},
						"ipv4_cidr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IPv4 CIDR block of the VPC (computed)",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "state of the VPC (computed)",
						},
						"subnets": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "list of subnets in the VPC (computed)",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "subnet ID",
									},
									"subnet_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "name of the subnet",
									},
									"cidr": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CIDR block of the subnet",
									},
									"used_ips": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "number of used IPs in the subnet",
									},
									"total_ips": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "total number of IPs in the subnet",
									},
								},
							},
						},
					},
				},
			},
			"network_config": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "network configuration block for consolidated network settings",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"assign_public_ip": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "whether to assign a public IP to nodes",
						},
						"vpc_names": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "list of VPC names to attach",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"security_groups": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "list of security group IDs to attach",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},

			"policy": {
				Type:          schema.TypeList,
				Optional:      true,
				Deprecated:    "Use 'scaling_policy' block instead. This field will be removed in v4.0.",
				ConflictsWith: []string{"scaling_policy"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type":           {Type: schema.TypeString, Required: true},
						"adjust":         {Type: schema.TypeInt, Required: true},
						"parameter":      {Type: schema.TypeString, Required: true},
						"operator":       {Type: schema.TypeString, Required: true},
						"value":          {Type: schema.TypeString, Required: true},
						"period_number":  {Type: schema.TypeString, Required: true},
						"period_seconds": {Type: schema.TypeString, Required: true},
						"cooldown":       {Type: schema.TypeString, Required: true},
					},
				},
				Description: "list of elastic scaling policies (deprecated, use 'scaling_policy' block instead)",
			},
			"scaling_policy": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"policy"},
				Description:   "list of elastic scaling policies (V3 structured block, preferred over 'policy')",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"scale_up", "scale_down"}, false),
							Description:  "type of scaling action: 'scale_up' or 'scale_down'",
						},
						"adjustment": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "number of nodes to adjust by",
						},
						"metric": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"cpu_utilization", "memory_utilization"}, false),
							Description:  "metric to monitor: 'cpu_utilization' or 'memory_utilization'",
						},
						"operator": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{">", "<", ">=", "<=", "=="}, false),
							Description:  "comparison operator",
						},
						"threshold": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "threshold value for the metric",
						},
						"evaluation_periods": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "number of evaluation periods",
						},
						"period_seconds": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "period length in seconds",
						},
						"cooldown_seconds": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "cooldown period in seconds after scaling action",
						},
					},
				},
			},
			"scheduled_policy": {
				Type:          schema.TypeList,
				Optional:      true,
				Deprecated:    "Use 'scheduled_action' block instead. This field will be removed in v4.0.",
				ConflictsWith: []string{"scheduled_action"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type":       {Type: schema.TypeString, Required: true},
						"adjust":     {Type: schema.TypeString, Required: true},
						"recurrence": {Type: schema.TypeString, Required: true},
					},
				},
				Description: "list of scheduled scaling policies (deprecated, use 'scheduled_action' block instead)",
			},
			"scheduled_action": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"scheduled_policy"},
				Description:   "list of scheduled scaling actions (V3 structured block, preferred over 'scheduled_policy')",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "name of the scheduled action",
						},
						attrScheduledActionType: {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								goe2econstants.AutoscalingScheduledActionTypeScaleUp,
								goe2econstants.AutoscalingScheduledActionTypeScaleDown,
								goe2econstants.AutoscalingScheduledActionTypeSetCapacity,
							}, false),
							Description: "type of action: 'scale_up', 'scale_down', or 'set_capacity'",
						},
						attrScheduledAdjustment: {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "number of nodes to adjust by (required for 'scale_up' and 'scale_down')",
						},
						attrScheduledTargetCap: {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "target node count (required for 'set_capacity')",
						},
						"recurrence": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "cron expression defining when the action should run",
						},
					},
				},
			},
		},
		CreateContext: resourceCreateScalerGroup,
		ReadContext:   resourceReadScalerGroup,
		DeleteContext: resourceDeleteScalerGroup,
		UpdateContext: resourceUpdateScalerGroup,
		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
			addAutoscalingV2DeprecationWarnings(diff)

			// Validate that at least one of image field pair is set
			hasVMImageName := diff.Get("vm_image_name") != nil && diff.Get("vm_image_name").(string) != ""
			hasImage := diff.Get("image") != nil && diff.Get("image").(string) != ""
			if !hasVMImageName && !hasImage {
				return errors.New(ErrEitherVMImageOrImageRequired)
			}

			// Validate that at least one of each V2/V3 field pair is set
			hasMinNodes := diff.Get(tfconstants.AttrMinNodes) != nil && diff.Get(tfconstants.AttrMinNodes).(int) > 0
			hasMinSize := diff.Get("min_size") != nil && diff.Get("min_size").(int) > 0
			if !hasMinNodes && !hasMinSize {
				return errors.New(ErrEitherMinNodesOrMinSizeRequired)
			}

			hasMaxNodes := diff.Get(tfconstants.AttrMaxNodes) != nil && diff.Get(tfconstants.AttrMaxNodes).(int) > 0
			hasMaxSize := diff.Get("max_size") != nil && diff.Get("max_size").(int) > 0
			if !hasMaxNodes && !hasMaxSize {
				return errors.New(ErrEitherMaxNodesOrMaxSizeRequired)
			}

			hasDesired := diff.Get(tfconstants.AttrDesired) != nil && diff.Get(tfconstants.AttrDesired).(int) > 0
			hasDesiredCapacity := diff.Get("desired_capacity") != nil && diff.Get("desired_capacity").(int) > 0
			if !hasDesired && !hasDesiredCapacity {
				return errors.New(ErrEitherDesiredOrDesiredCapacityRequired)
			}

			// Get min size (handle both V2 and V3 fields)
			var min int
			if hasMinSize {
				min = diff.Get("min_size").(int)
			} else if hasMinNodes {
				min = diff.Get(tfconstants.AttrMinNodes).(int)
			}

			// Get desired capacity (handle both V2 and V3 fields)
			var desired int
			if hasDesiredCapacity {
				desired = diff.Get("desired_capacity").(int)
			} else if hasDesired {
				desired = diff.Get(tfconstants.AttrDesired).(int)
			}

			// Get max size (handle both V2 and V3 fields)
			var max int
			if hasMaxSize {
				max = diff.Get("max_size").(int)
			} else if hasMaxNodes {
				max = diff.Get(tfconstants.AttrMaxNodes).(int)
			}

			// Validate min <= desired <= max
			if min > 0 && desired > 0 && min > desired {
				return fmt.Errorf("min_size/min_nodes (%d) cannot be greater than desired_capacity/desired (%d)", min, desired)
			}

			if desired > 0 && max > 0 && desired > max {
				return fmt.Errorf("desired_capacity/desired (%d) cannot be greater than max_size/max_nodes (%d)", desired, max)
			}

			// Validate state requirements for security group updates
			if diff.HasChange(tfconstants.AttrSecurityGroupIDs) {
				status := getStatusFromDiff(diff)
				if status != "" && !autoscalingScalerGroupStatusIn(status, tfconstants.AutoscalingScalerGroupRunningStates) {
					return fmt.Errorf(ErrSecurityGroupUpdatesRequireRunningFmt, status)
				}
			}

			// Validate state requirements for VPC updates
			if diff.HasChange("vpc") || diff.HasChange("vpc_config") {
				status := getStatusFromDiff(diff)
				if status != "" && !autoscalingScalerGroupStatusIn(status, tfconstants.AutoscalingScalerGroupStoppedStates) {
					return fmt.Errorf(ErrVPCUpdatesRequireStoppedFmt, status)
				}
			}

			// Validate state + VPC requirements for public IP updates
			if diff.HasChange("assign_public_ip") || diff.HasChange(tfconstants.AttrPublicIPRequired) {
				status := getStatusFromDiff(diff)
				if status != "" && !autoscalingScalerGroupStatusIn(status, tfconstants.AutoscalingScalerGroupStoppedStates) {
					return fmt.Errorf(ErrPublicIPUpdatesRequireStoppedFmt, status)
				}

				// Check if VPC is attached
				vpcs := getVPCsFromDiff(diff)
				if len(vpcs) == 0 {
					return errors.New(ErrPublicIPUpdatesRequireVPC)
				}
			}

			// Validate network_config conflicts with individual fields
			networkConfig := expandNetworkConfigFromDiff(diff)
			if networkConfig != nil {
				// Check for assign_public_ip conflict
				if networkConfig.AssignPublicIP {
					if diff.Get("assign_public_ip") != nil && diff.Get("assign_public_ip") != "" {
						if diff.Get("assign_public_ip").(bool) != networkConfig.AssignPublicIP {
							return fmt.Errorf("cannot set both 'network_config.assign_public_ip' and 'assign_public_ip' fields")
						}
					}
					if diff.Get(tfconstants.AttrPublicIPRequired) != nil && diff.Get(tfconstants.AttrPublicIPRequired) != "" {
						if diff.Get(tfconstants.AttrPublicIPRequired).(bool) != networkConfig.AssignPublicIP {
							return fmt.Errorf("cannot set both 'network_config.assign_public_ip' and 'is_public_ip_required' fields")
						}
					}
				}

				// Check for VPC conflicts
				if len(networkConfig.VPCNames) > 0 {
					if diff.Get("vpc") != nil {
						return fmt.Errorf("cannot set both 'network_config.vpc_names' and 'vpc' fields")
					}
					if diff.Get("vpc_config") != nil {
						return fmt.Errorf("cannot set both 'network_config.vpc_names' and 'vpc_config' fields")
					}
				}

				// Check for security group conflicts
				if len(networkConfig.SecurityGroups) > 0 {
					if diff.Get(tfconstants.AttrSecurityGroupIDs) != nil {
						return fmt.Errorf("cannot set both 'network_config.security_groups' and 'security_group_ids' fields")
					}
				}

				// Validate state requirements for network_config changes
				if diff.HasChange("network_config") {
					// Check if any network_config field changed
					oldRaw, newRaw := diff.GetChange("network_config")
					oldConfig := expandNetworkConfigFromRaw(oldRaw)
					newConfig := expandNetworkConfigFromRaw(newRaw)

					// Check public IP change
					if oldConfig == nil || newConfig == nil || oldConfig.AssignPublicIP != newConfig.AssignPublicIP {
						status := getStatusFromDiff(diff)
						if status != "" && !autoscalingScalerGroupStatusIn(status, tfconstants.AutoscalingScalerGroupStoppedStates) {
							return fmt.Errorf("network_config.assign_public_ip updates require scaler group to be in 'Stopped' state, current: %s", status)
						}
						vpcs := getVPCsFromDiff(diff)
						if len(vpcs) == 0 && (newConfig == nil || len(newConfig.VPCNames) == 0) {
							return fmt.Errorf("network_config.assign_public_ip updates require at least one VPC to be attached")
						}
					}

					// Check VPC change
					if oldConfig == nil || newConfig == nil || !stringSlicesEqual(oldConfig.VPCNames, newConfig.VPCNames) {
						status := getStatusFromDiff(diff)
						if status != "" && !autoscalingScalerGroupStatusIn(status, tfconstants.AutoscalingScalerGroupStoppedStates) {
							return fmt.Errorf("network_config.vpc_names updates require scaler group to be in 'Stopped' state, current: %s", status)
						}
					}

					// Check security group change
					if oldConfig == nil || newConfig == nil || !intSlicesEqual(oldConfig.SecurityGroups, newConfig.SecurityGroups) {
						status := getStatusFromDiff(diff)
						if status != "" && !autoscalingScalerGroupStatusIn(status, tfconstants.AutoscalingScalerGroupRunningStates) {
							return fmt.Errorf("network_config.security_groups updates require scaler group to be in 'Running' state, current: %s", status)
						}
					}
				}
			}

			// Validate scheduled_action conditional requirements.
			// Schema can't express "required if" based on action_type, so enforce here for determinism.
			if v, ok := diff.GetOk(attrScheduledAction); ok {
				for i, raw := range v.([]interface{}) {
					action, ok := raw.(map[string]interface{})
					if !ok {
						continue
					}
					actionType, _ := action[attrScheduledActionType].(string)
					switch actionType {
					case goe2econstants.AutoscalingScheduledActionTypeSetCapacity:
						if tc, ok := action[attrScheduledTargetCap].(int); !ok || tc <= 0 {
							return fmt.Errorf(ErrScheduledActionTargetCapacityRequiredFmt, attrScheduledAction, i, attrScheduledTargetCap, attrScheduledActionType, goe2econstants.AutoscalingScheduledActionTypeSetCapacity)
						}
					case goe2econstants.AutoscalingScheduledActionTypeScaleUp, goe2econstants.AutoscalingScheduledActionTypeScaleDown:
						if adj, ok := action[attrScheduledAdjustment].(int); !ok || adj == 0 {
							return fmt.Errorf(ErrScheduledActionAdjustmentRequiredFmt, attrScheduledAction, i, attrScheduledAdjustment, attrScheduledActionType, actionType)
						}
					}
				}
			}

			return nil
		},

		Importer: &schema.ResourceImporter{
			State: node.CustomImportStateFunc,
		},
	}
}

func resourceCreateScalerGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[INFO] Starting CreateScalerGroup operation")

	cfg := m.(*config.Config)
	var diags diag.Diagnostics
	if warn := autoscalingV2DeprecationWarningDiagnostic(d); warn != nil {
		diags = append(diags, *warn)
	}

	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	projectID, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	// Create GoE2E client for this project/region
	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return append(diags, diag.FromErr(fmt.Errorf("failed to create GoE2E client: %w", err))...)
	}

	// Get image name (handle V2/V3 fields)
	imageName, err := GetImageName(d)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	// Get saved image details using GoE2E client
	savedImage, err := getSavedImageByName(ctx, goe2eClient, imageName)
	if err != nil {
		return append(diags, diag.FromErr(fmt.Errorf("failed to fetch saved image details for '%s': %w", imageName, err))...)
	}

	if err := d.Set("vm_image_id", savedImage.ImageID); err != nil {
		return append(diags, diag.FromErr(fmt.Errorf("failed to set vm_image_id: %v", err))...)
	}
	if err := d.Set("vm_template_id", savedImage.TemplateID); err != nil {
		return append(diags, diag.FromErr(fmt.Errorf("failed to set vm_template_id: %v", err))...)
	}

	log.Printf("[DEBUG] Image Details → ID: %s, TemplateID: %d, Distro: %s", savedImage.ImageID, savedImage.TemplateID, savedImage.Distro)

	// Check if network_config block is present
	networkConfig := expandNetworkConfig(d)

	// Get security group ID (network_config takes precedence)
	var sgID int
	var securityGroupIDs []int
	if networkConfig != nil && len(networkConfig.SecurityGroups) > 0 {
		// Use security groups from network_config
		securityGroupIDs = networkConfig.SecurityGroups
		sgID = securityGroupIDs[0] // Use first one for MyAccountSGID (API requirement)
		log.Printf("[INFO] Using Security Group IDs from network_config: %v", securityGroupIDs)
	} else if v, ok := d.GetOk("my_account_sg_id"); ok {
		// Use individual field
		sgID = v.(int)
		securityGroupIDs = []int{sgID}
		log.Printf("[INFO] Using user-provided Security Group ID: %d", sgID)
	} else {
		// Use default
		sgID, err = getDefaultSecurityGroupID(ctx, goe2eClient)
		if err != nil {
			return append(diags, diag.FromErr(fmt.Errorf("failed to fetch default security group ID: %w", err))...)
		}
		securityGroupIDs = []int{sgID}
		log.Printf("[INFO] Using default Security Group ID from API: %d", sgID)
		if err := d.Set("my_account_sg_id", sgID); err != nil {
			return append(diags, diag.FromErr(fmt.Errorf("failed to set my_account_sg_id: %v", err))...)
		}
	}

	if err := d.Set(tfconstants.AttrSecurityGroupIDs, securityGroupIDs); err != nil {
		return append(diags, diag.FromErr(fmt.Errorf("failed to set security_group_ids: %v", err))...)
	}

	// Expand create request (handles V2/V3 fields and network_config)
	req, err := expandCreateScalerGroupRequestV3(ctx, d, cfg, goe2eClient, projectID, region, sgID, savedImage, networkConfig)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	requestJSON := MarshalScalerGroupCreateRequestForLog(req)
	log.Printf("[DEBUG] CreateScalerGroup Request JSON:\n%s", requestJSON)

	// Create scaler group using GoE2E client
	scalerGroup, _, err := goe2eClient.Autoscaling.CreateScalerGroup(ctx, req)
	if err != nil {
		return append(diags, diag.FromErr(fmt.Errorf("failed to create scaler group: %w", err))...)
	}

	log.Printf("[INFO] ScalerGroup created with ID: %s", scalerGroup.ID)
	d.SetId(scalerGroup.ID)

	return append(diags, resourceReadScalerGroup(ctx, d, m)...)
}

func resourceReadScalerGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] Reading ScalerGroup ID: %s", d.Id())

	cfg := m.(*config.Config)

	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	projectID, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Create GoE2E client for this project/region
	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create GoE2E client: %w", err))
	}

	id := d.Id()

	// Get scaler group using GoE2E client
	group, _, err := goe2eClient.Autoscaling.GetScalerGroup(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read scaler group: %w", err))
	}

	if group == nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Retrieved ScalerGroup: %+v", group)

	// Set image fields (handle V2/V3)
	apiVMImageName := group.VMImageName
	// Set both V2 and V3 fields if they exist in state, otherwise set based on what's in state
	if _, ok := d.GetOk("vm_image_name"); ok {
		stateVMImageName := d.Get("vm_image_name").(string)
		if !strings.HasPrefix(stateVMImageName, apiVMImageName) {
			log.Printf("[INFO] Updating vm_image_name to: %s", apiVMImageName)
			if err := d.Set("vm_image_name", apiVMImageName); err != nil {
				return diag.FromErr(fmt.Errorf("failed to set vm_image_name: %v", err))
			}
		}
	}
	if _, ok := d.GetOk("image"); ok {
		if err := d.Set("image", apiVMImageName); err != nil {
			return diag.FromErr(fmt.Errorf("failed to set image: %v", err))
		}
	}

	// Set basic fields
	if err := d.Set("name", group.Name); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set name: %v", err))
	}
	if err := d.Set(tfconstants.AttrPlan, group.PlanName); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set plan: %v", err))
	}
	if err := d.Set("plan_id", strconv.Itoa(group.PlanID)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set plan_id: %v", err))
	}
	if err := d.Set("sku_id", strconv.Itoa(group.PlanID)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set sku_id: %v", err))
	}
	if err := d.Set("policy_type", group.PolicyType); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set policy_type: %v", err))
	}
	if err := d.Set("vm_image_id", strconv.Itoa(group.ImageID)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set vm_image_id: %v", err))
	}

	// Set size fields (handle V2/V3)
	if err := d.Set(tfconstants.AttrDesired, group.Desired); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set desired: %v", err))
	}
	if err := d.Set("desired_capacity", group.Desired); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set desired_capacity: %v", err))
	}
	if err := d.Set(tfconstants.AttrMinNodes, group.MinNodes); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set min_nodes: %v", err))
	}
	if err := d.Set("min_size", group.MinNodes); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set min_size: %v", err))
	}
	if err := d.Set(tfconstants.AttrMaxNodes, group.MaxNodes); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set max_nodes: %v", err))
	}
	if err := d.Set("max_size", group.MaxNodes); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set max_size: %v", err))
	}

	// Set running node count (V2 and V3)
	if err := d.Set("running", group.Running); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set running: %v", err))
	}
	if err := d.Set("running_node_count", group.Running); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set running_node_count: %v", err))
	}

	// Normalize and set status (handle V2/V3)
	normalizedStatus := NormalizeStatus(group.ProvisionStatus)
	if err := d.Set("provision_status", normalizedStatus); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set provision_status: %v", err))
	}
	// Also set V3 status field (convert to lowercase for V3)
	v3Status := strings.ToLower(normalizedStatus)
	if err := d.Set("status", v3Status); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set status: %v", err))
	}

	// Fetch attached VPCs using GoE2E client and flatten using helper
	var vpcNames []string
	attachedVPCs, _, err := goe2eClient.Autoscaling.GetAttachedVPCsForScalerGroup(ctx, id)
	if err != nil {
		log.Printf("[WARN] Failed to fetch attached VPCs from scaler group: %v", err)
	} else {
		// Extract VPC names for network_config
		for _, vpc := range attachedVPCs {
			vpcNames = append(vpcNames, vpc.Name)
		}

		fetchedVPCs, err := flattenVPCConfig(ctx, attachedVPCs, goe2eClient)
		if err != nil {
			log.Printf("[WARN] Failed to flatten VPC config: %v", err)
		} else {
			// Set VPC fields (both V2 and V3)
			if err := d.Set("vpc", fetchedVPCs); err != nil {
				return diag.FromErr(fmt.Errorf("failed to set vpc details: %v", err))
			}
			if err := d.Set("vpc_config", fetchedVPCs); err != nil {
				return diag.FromErr(fmt.Errorf("failed to set vpc_config details: %v", err))
			}
		}
	}

	// Set slug_name if template ID is available
	if templateID, ok := d.Get("vm_template_id").(int); ok && templateID > 0 {
		_, slugName, err := getPlanDetailsFromPlanName(ctx, goe2eClient, templateID, group.PlanName)
		if err == nil {
			d.Set("slug_name", slugName)
		} else {
			log.Printf("[WARN] Failed to recompute slug_name: %v", err)
		}
	}

	// Get public IP status using GoE2E client
	var assignPublicIP bool
	ipStatus, _, err := goe2eClient.Autoscaling.GetPublicIPStatus(ctx, id)
	if err != nil {
		log.Printf("[WARN] Failed to fetch public IP status: %v", err)
		// Fallback to state value if API call failed
		if v, ok := d.GetOk("assign_public_ip"); ok {
			assignPublicIP = v.(bool)
		} else if v, ok := d.GetOk(tfconstants.AttrPublicIPRequired); ok {
			assignPublicIP = v.(bool)
		}
	} else {
		assignPublicIP = ipStatus.IsPublicIPRequired
		// Set both V2 and V3 fields
		if err := d.Set(tfconstants.AttrPublicIPRequired, ipStatus.IsPublicIPRequired); err != nil {
			return diag.FromErr(fmt.Errorf("failed to set is_public_ip_required: %v", err))
		}
		if err := d.Set("assign_public_ip", ipStatus.IsPublicIPRequired); err != nil {
			return diag.FromErr(fmt.Errorf("failed to set assign_public_ip: %v", err))
		}
	}

	// Get security group IDs
	var securityGroupIDs []int
	if sgIDsRaw, ok := d.GetOk(tfconstants.AttrSecurityGroupIDs); ok {
		sgIDsList := sgIDsRaw.([]interface{})
		securityGroupIDs = make([]int, len(sgIDsList))
		for i, v := range sgIDsList {
			securityGroupIDs[i] = v.(int)
		}
	}

	// Set node details
	if err := flattenNodes(d, group.Nodes); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set nodes: %v", err))
	}

	// Flatten and set policies (both V2 and V3)
	v2Policies, v3Policies := flattenScalingPolicy(group)
	if err := d.Set("policy", v2Policies); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set policy: %v", err))
	}
	if err := d.Set("scaling_policy", v3Policies); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set scaling_policy: %v", err))
	}

	// Flatten and set scheduled policies (both V2 and V3)
	v2Scheduled, v3Scheduled := flattenScheduledAction(group)
	if err := d.Set("scheduled_policy", v2Scheduled); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set scheduled_policy: %v", err))
	}
	if err := d.Set(attrScheduledAction, v3Scheduled); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set scheduled_action: %v", err))
	}

	// Flatten and set network_config block
	networkConfig := flattenNetworkConfig(assignPublicIP, vpcNames, securityGroupIDs)
	if len(networkConfig) > 0 {
		if err := d.Set("network_config", networkConfig); err != nil {
			return diag.FromErr(fmt.Errorf("failed to set network_config: %v", err))
		}
	}

	log.Printf("[INFO] ScalerGroup ID %s state synced successfully", id)
	return nil
}

func resourceDeleteScalerGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] Deleting ScalerGroup ID: %s", d.Id())

	cfg := m.(*config.Config)

	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	projectID, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Create GoE2E client for this project/region
	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create GoE2E client: %w", err))
	}

	id := d.Id()

	// Delete scaler group using GoE2E client
	_, err = goe2eClient.Autoscaling.DeleteScalerGroup(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete scaler group: %w", err))
	}

	d.SetId("")
	log.Println("[INFO] ScalerGroup deleted successfully")
	return nil
}

// expandCreateScalerGroupRequestV3 creates a GoE2E ScalerGroupCreateRequest from schema data
// Handles both V2 and V3 field names, and new structured blocks including network_config
func expandCreateScalerGroupRequestV3(ctx context.Context, d *schema.ResourceData, cfg *config.Config, client *goe2e.Client, projectID, region string, sgID int, savedImage *goe2e.SavedImage, networkConfig *NetworkConfig) (*goe2e.ScalerGroupCreateRequest, error) {
	planName := d.Get(tfconstants.AttrPlan).(string)
	imageName, _ := GetImageName(d)

	// Get plan details (temporary: using old client for this until GoE2E has equivalent)
	planID, slugName, err := getPlanDetailsFromPlanName(ctx, client, savedImage.TemplateID, planName)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch plan details: %w", err)
	}

	// Get size fields (handle V2/V3)
	minSize := getMinSize(d)
	maxSize := getMaxSize(d)
	desired := getDesiredCapacity(d)

	// Get encryption flag (handle V2/V3)
	enableEncryption := getEnableEncryption(d)
	encryptionPassphrase := d.Get(tfconstants.AttrEncryptionPassphrase).(string)

	// Get public IP flag (network_config takes precedence)
	var assignPublicIP bool
	if networkConfig != nil {
		assignPublicIP = networkConfig.AssignPublicIP
	} else {
		assignPublicIP = getAssignPublicIP(d)
	}

	// Expand policies using helper
	elasticPolicies := expandScalingPolicy(d)

	// Expand scheduled policies using helper
	schedPolicies := expandScheduledAction(d)

	// Expand VPCs (network_config takes precedence)
	var vpcDetails []goe2e.VPCDetail
	if networkConfig != nil && len(networkConfig.VPCNames) > 0 {
		// Use VPC names from network_config
		for _, vpcName := range networkConfig.VPCNames {
			vpcDetail, _, err := client.Vpcs.GetVPCByName(ctx, vpcName)
			if err != nil {
				return nil, fmt.Errorf("failed to get VPC details for %s: %w", vpcName, err)
			}
			vpcDetails = append(vpcDetails, goe2e.VPCDetail{
				Name:      vpcDetail.Name,
				NetworkID: vpcDetail.NetworkID,
				IPv4CIDR:  vpcDetail.IPv4CIDR,
				State:     vpcDetail.State,
			})
		}
	} else {
		// Use existing VPC config logic
		vpcDetails, err = expandVPCConfig(ctx, d, client)
		if err != nil {
			return nil, err
		}
	}

	policyType := d.Get("policy_type").(string)
	if len(elasticPolicies) > 0 && policyType == "" {
		policyType = "elastic"
	}

	return &goe2e.ScalerGroupCreateRequest{
		Name:                 d.Get("name").(string),
		PlanName:             planName,
		PlanID:               planID,
		SKUID:                planID,
		SlugName:             slugName,
		VMImageID:            savedImage.ImageID,
		VMImageName:          imageName,
		VMTemplateID:         savedImage.TemplateID,
		MyAccountSGID:        sgID,
		IsEncryptionEnabled:  enableEncryption,
		EncryptionPassphrase: encryptionPassphrase,
		IsPublicIPRequired:   assignPublicIP,
		MinNodes:             strconv.Itoa(minSize),
		MaxNodes:             strconv.Itoa(maxSize),
		Desired:              strconv.Itoa(desired),
		PolicyType:           policyType,
		Policy:               elasticPolicies,
		ScheduledPolicy:      schedPolicies,
		VPC:                  vpcDetails,
	}, nil
}

func resourceUpdateScalerGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	var diags diag.Diagnostics
	if warn := autoscalingV2DeprecationWarningDiagnostic(d); warn != nil {
		diags = append(diags, *warn)
	}

	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	projectID, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	// Create GoE2E client for this project/region
	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return append(diags, diag.FromErr(fmt.Errorf("failed to create GoE2E client: %w", err))...)
	}

	id := d.Id()

	// Handle status changes (handle both V2 and V3 fields)
	if d.HasChange("provision_status") || d.HasChange("status") {
		var newStatus string
		if d.HasChange("status") {
			newStatus = d.Get("status").(string)
			// Convert V3 lowercase to V2 format for API
			if newStatus == goe2econstants.AutoscalingScalerGroupStatusRunningLower {
				newStatus = goe2econstants.AutoscalingScalerGroupStatusRunning
			} else if newStatus == goe2econstants.AutoscalingScalerGroupStatusStoppedLower {
				newStatus = goe2econstants.AutoscalingScalerGroupStatusStopped
			}
		} else {
			newStatus = d.Get("provision_status").(string)
		}

		oldStatus := getStatus(d)
		log.Printf("[INFO] Changing status from %s → %s", oldStatus, newStatus)

		_, err := goe2eClient.Autoscaling.UpdateScalerGroupStatus(ctx, id, newStatus)
		if err != nil {
			return append(diags, diag.FromErr(fmt.Errorf("failed to update status to %s: %w", newStatus, err))...)
		}

		return append(diags, resourceReadScalerGroup(ctx, d, m)...)
	}

	// Get size fields (handle V2/V3)
	minNodes := getMinSize(d)
	maxNodes := getMaxSize(d)
	desired := getDesiredCapacity(d)

	if desired < minNodes || desired > maxNodes {
		return append(diags, diag.Errorf("desired node count (%d) must be between min_size/min_nodes (%d) and max_size/max_nodes (%d)", desired, minNodes, maxNodes)...)
	}

	// If only desired changed, call separate API
	hasDesiredChange := d.HasChange(tfconstants.AttrDesired) || d.HasChange("desired_capacity")
	hasOtherChanges := d.HasChange(tfconstants.AttrMinNodes) || d.HasChange("min_size") ||
		d.HasChange(tfconstants.AttrMaxNodes) || d.HasChange("max_size") ||
		d.HasChange("policy_type") || d.HasChange("policy") || d.HasChange("scaling_policy") ||
		d.HasChange("scheduled_policy") || d.HasChange(attrScheduledAction) ||
		d.HasChange("network_config") // network_config changes handled separately above

	if hasDesiredChange && !hasOtherChanges {
		log.Printf("[INFO] Only desired node count changed; using separate API.")
		_, err := goe2eClient.Autoscaling.UpdateDesiredNodeCount(ctx, id, desired)
		if err != nil {
			return append(diags, diag.FromErr(fmt.Errorf("failed to update desired node count: %w", err))...)
		}
		return append(diags, resourceReadScalerGroup(ctx, d, m)...)
	}

	// Handle network_config block changes
	if d.HasChange("network_config") {
		oldRaw, newRaw := d.GetChange("network_config")
		oldConfig := expandNetworkConfigFromRaw(oldRaw)
		newConfig := expandNetworkConfigFromRaw(newRaw)

		// Handle public IP changes
		if oldConfig == nil || newConfig == nil || oldConfig.AssignPublicIP != newConfig.AssignPublicIP {
			var newVal bool
			if newConfig != nil {
				newVal = newConfig.AssignPublicIP
			} else {
				// Fallback to individual field if network_config removed
				newVal = getAssignPublicIP(d)
			}

			log.Printf("[INFO] assign_public_ip changed to %v (from network_config)", newVal)

			group, _, err := goe2eClient.Autoscaling.GetScalerGroup(ctx, id)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to fetch scaler group status: %w", err))
			}
			if group == nil {
				return diag.Errorf("scaler group not found")
			}
			normalizedStatus := NormalizeStatus(group.ProvisionStatus)
			if !autoscalingScalerGroupStatusIn(normalizedStatus, tfconstants.AutoscalingScalerGroupStoppedStates) {
				return diag.Errorf("ScalerGroup must be in 'Stopped' state to attach/detach public IP. Current: %s", group.ProvisionStatus)
			}

			// Check if VPC is attached
			vpcsRaw, ok := d.GetOk("vpc")
			if !ok {
				vpcsRaw, ok = d.GetOk("vpc_config")
			}
			if !ok || len(vpcsRaw.([]interface{})) == 0 {
				return diag.Errorf("At least one VPC must be attached to attach/detach public IP")
			}

			if newVal {
				log.Printf("[INFO] Triggering Public IP ATTACH (from network_config)")
				_, _, err := goe2eClient.Autoscaling.AttachPublicIPToScalerGroup(ctx, id)
				if err != nil {
					return diag.FromErr(fmt.Errorf("failed to attach public IP: %w", err))
				}
			} else {
				log.Printf("[INFO] Triggering Public IP DETACH (from network_config)")
				_, _, err := goe2eClient.Autoscaling.DetachPublicIPFromScalerGroup(ctx, id)
				if err != nil {
					return diag.FromErr(fmt.Errorf("failed to detach public IP: %w", err))
				}
			}
		}

		// Handle VPC changes
		if oldConfig == nil || newConfig == nil || !stringSlicesEqual(oldConfig.VPCNames, newConfig.VPCNames) {
			var newVPCNames []string
			if newConfig != nil {
				newVPCNames = newConfig.VPCNames
			} else {
				// Fallback to individual fields if network_config removed
				if v, ok := d.GetOk("vpc_config"); ok {
					for _, vRaw := range v.([]interface{}) {
						vMap := vRaw.(map[string]interface{})
						newVPCNames = append(newVPCNames, vMap["name"].(string))
					}
				} else if v, ok := d.GetOk("vpc"); ok {
					for _, vRaw := range v.([]interface{}) {
						vMap := vRaw.(map[string]interface{})
						newVPCNames = append(newVPCNames, vMap["name"].(string))
					}
				}
			}

			group, _, err := goe2eClient.Autoscaling.GetScalerGroup(ctx, id)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to fetch scaler group status for update: %w", err))
			}
			if group == nil {
				return diag.Errorf("scaler group not found")
			}
			normalizedStatus := NormalizeStatus(group.ProvisionStatus)
			if !autoscalingScalerGroupStatusIn(normalizedStatus, tfconstants.AutoscalingScalerGroupStoppedStates) {
				return diag.Errorf("VPCs can only be attached or detached when the scaler group is in 'Stopped' state. Current state: %q", group.ProvisionStatus)
			}

			// Get old VPC names
			var oldVPCNames []string
			if oldConfig != nil {
				oldVPCNames = oldConfig.VPCNames
			} else {
				oldVPCNames = extractVpcNames(d.Get("vpc").([]interface{}))
				if len(oldVPCNames) == 0 {
					oldVPCNames = extractVpcNames(d.Get("vpc_config").([]interface{}))
				}
			}

			toAttach := difference(newVPCNames, oldVPCNames)
			toDetach := difference(oldVPCNames, newVPCNames)

			// Attach VPCs
			for _, vpcName := range toAttach {
				vpcDetail, _, err := goe2eClient.Vpcs.GetVPCByName(ctx, vpcName)
				if err != nil {
					return diag.FromErr(fmt.Errorf("failed to get VPC details for name %q: %w", vpcName, err))
				}

				vpcDetailList := []goe2e.VPCDetail{
					{
						Name:      vpcDetail.Name,
						NetworkID: vpcDetail.NetworkID,
						IPv4CIDR:  vpcDetail.IPv4CIDR,
						State:     vpcDetail.State,
					},
				}

				attachReq := &goe2e.VPCAttachRequest{VPC: vpcDetailList}
				_, err = goe2eClient.Autoscaling.AttachVPCToScalerGroup(ctx, id, attachReq)
				if err != nil {
					return diag.FromErr(fmt.Errorf("failed to attach VPC %q: %w", vpcName, err))
				}
			}

			// Detach VPCs
			for _, vpcName := range toDetach {
				vpcDetail, _, err := goe2eClient.Vpcs.GetVPCByName(ctx, vpcName)
				if err != nil {
					return diag.FromErr(fmt.Errorf("failed to get VPC ID for name %q: %w", vpcName, err))
				}
				vpcID := strconv.Itoa(vpcDetail.NetworkID)
				_, err = goe2eClient.Autoscaling.DetachVPCFromScalerGroup(ctx, id, vpcID)
				if err != nil {
					return diag.FromErr(fmt.Errorf("failed to detach VPC %q: %w", vpcName, err))
				}
			}
		}

		// Handle security group changes
		if oldConfig == nil || newConfig == nil || !intSlicesEqual(oldConfig.SecurityGroups, newConfig.SecurityGroups) {
			var newSGIDs []int
			if newConfig != nil {
				newSGIDs = newConfig.SecurityGroups
			} else {
				// Fallback to individual field if network_config removed
				if v, ok := d.GetOk(tfconstants.AttrSecurityGroupIDs); ok {
					sgIDsList := v.([]interface{})
					newSGIDs = make([]int, len(sgIDsList))
					for i, v := range sgIDsList {
						newSGIDs[i] = v.(int)
					}
				}
			}

			group, _, err := goe2eClient.Autoscaling.GetScalerGroup(ctx, id)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to fetch scaler group status: %w", err))
			}
			if group == nil {
				return diag.Errorf("scaler group not found")
			}
			normalizedStatus := NormalizeStatus(group.ProvisionStatus)
			if !autoscalingScalerGroupStatusIn(normalizedStatus, tfconstants.AutoscalingScalerGroupRunningStates) {
				return diag.Errorf("Scaler group must be in 'Running' state to update security groups. Current: %s", group.ProvisionStatus)
			}

			if len(newSGIDs) == 0 {
				return diag.Errorf("At least one security group must be attached to the scaler group")
			}

			// Get old security group IDs
			var oldSGIDs []int
			if oldConfig != nil {
				oldSGIDs = oldConfig.SecurityGroups
			} else {
				oldRaw, _ := d.GetChange(tfconstants.AttrSecurityGroupIDs)
				oldList := expandIntList(oldRaw.([]interface{}))
				oldSGIDs = oldList
			}

			oldStr := intSliceToStringSlice(oldSGIDs)
			newStr := intSliceToStringSlice(newSGIDs)

			toAttach := difference(newStr, oldStr)
			toDetach := difference(oldStr, newStr)

			for _, sgIDStr := range toAttach {
				sgID, _ := strconv.Atoi(sgIDStr)
				log.Printf("[INFO] Attaching Security Group ID %d (from network_config)", sgID)
				_, err := goe2eClient.Autoscaling.AttachSecurityGroupToScalerGroup(ctx, id, sgID)
				if err != nil {
					return diag.FromErr(fmt.Errorf("failed to attach SG %d: %w", sgID, err))
				}
			}

			for _, sgIDStr := range toDetach {
				sgID, _ := strconv.Atoi(sgIDStr)
				log.Printf("[INFO] Detaching Security Group ID %d (from network_config)", sgID)
				_, err := goe2eClient.Autoscaling.DetachSecurityGroupFromScalerGroup(ctx, id, sgID)
				if err != nil {
					return diag.FromErr(fmt.Errorf("failed to detach SG %d: %w", sgID, err))
				}
			}
		}
	}

	// Handle individual security_group_ids changes (only if network_config didn't change)
	if d.HasChange(tfconstants.AttrSecurityGroupIDs) && !d.HasChange("network_config") {
		log.Printf("[INFO] Detected change in security_group_ids for Scaler Group %s", id)

		group, _, err := goe2eClient.Autoscaling.GetScalerGroup(ctx, id)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to fetch scaler group status: %w", err))
		}
		if group == nil {
			return diag.Errorf("scaler group not found")
		}
		normalizedStatus := NormalizeStatus(group.ProvisionStatus)
		if !autoscalingScalerGroupStatusIn(normalizedStatus, tfconstants.AutoscalingScalerGroupRunningStates) {
			return diag.Errorf("Scaler group must be in 'Running' state to update security groups. Current: %s", group.ProvisionStatus)
		}

		oldRaw, newRaw := d.GetChange(tfconstants.AttrSecurityGroupIDs)
		oldList := expandIntList(oldRaw.([]interface{}))
		newList := expandIntList(newRaw.([]interface{}))

		if len(newList) == 0 {
			return diag.Errorf("At least one security group must be attached to the scaler group")
		}

		oldStr := intSliceToStringSlice(oldList)
		newStr := intSliceToStringSlice(newList)

		toAttach := difference(newStr, oldStr)
		toDetach := difference(oldStr, newStr)

		for _, sgIDStr := range toAttach {
			sgID, _ := strconv.Atoi(sgIDStr)
			log.Printf("[INFO] Attaching Security Group ID %d", sgID)
			_, err := goe2eClient.Autoscaling.AttachSecurityGroupToScalerGroup(ctx, id, sgID)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to attach SG %d: %w", sgID, err))
			}
		}

		for _, sgIDStr := range toDetach {
			sgID, _ := strconv.Atoi(sgIDStr)
			log.Printf("[INFO] Detaching Security Group ID %d", sgID)
			_, err := goe2eClient.Autoscaling.DetachSecurityGroupFromScalerGroup(ctx, id, sgID)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to detach SG %d: %w", sgID, err))
			}
		}
	}

	// Handle individual VPC changes (only if network_config didn't change)
	if (d.HasChange("vpc") || d.HasChange("vpc_config")) && !d.HasChange("network_config") {
		group, _, err := goe2eClient.Autoscaling.GetScalerGroup(ctx, id)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to fetch scaler group status for update: %w", err))
		}
		if group == nil {
			return diag.Errorf("scaler group not found")
		}
		normalizedStatus := NormalizeStatus(group.ProvisionStatus)
		if !autoscalingScalerGroupStatusIn(normalizedStatus, tfconstants.AutoscalingScalerGroupStoppedStates) {
			return diag.Errorf("VPCs can only be attached or detached when the scaler group is in 'Stopped' state. Current state: %q", group.ProvisionStatus)
		}

		// Get old and new VPC lists (handle both V2 and V3)
		var oldRaw, newRaw interface{}
		if d.HasChange("vpc") {
			oldRaw, newRaw = d.GetChange("vpc")
		} else {
			oldRaw, newRaw = d.GetChange("vpc_config")
		}

		oldList := extractVpcNames(oldRaw.([]interface{}))
		newList := extractVpcNames(newRaw.([]interface{}))

		toAttach := difference(newList, oldList)
		toDetach := difference(oldList, newList)

		// Attach VPCs
		for _, vpcName := range toAttach {
			// Get full VPC details using GoE2E client
			vpcDetail, _, err := goe2eClient.Vpcs.GetVPCByName(ctx, vpcName)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to get VPC details for name %q: %w", vpcName, err))
			}

			// Convert to GoE2E VPCDetail format
			vpcDetailList := []goe2e.VPCDetail{
				{
					Name:      vpcDetail.Name,
					NetworkID: vpcDetail.NetworkID,
					IPv4CIDR:  vpcDetail.IPv4CIDR,
					State:     vpcDetail.State,
				},
			}

			attachReq := &goe2e.VPCAttachRequest{VPC: vpcDetailList}
			_, err = goe2eClient.Autoscaling.AttachVPCToScalerGroup(ctx, id, attachReq)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to attach VPC %q: %w", vpcName, err))
			}
		}

		// Detach VPCs
		for _, vpcName := range toDetach {
			// Get VPC details to find network ID using GoE2E client
			vpcDetail, _, err := goe2eClient.Vpcs.GetVPCByName(ctx, vpcName)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to get VPC ID for name %q: %w", vpcName, err))
			}
			vpcID := strconv.Itoa(vpcDetail.NetworkID)
			_, err = goe2eClient.Autoscaling.DetachVPCFromScalerGroup(ctx, id, vpcID)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to detach VPC %q: %w", vpcName, err))
			}
		}

		// Refresh VPC state using GoE2E client
		vpcNames := extractVpcNames(newRaw.([]interface{}))
		vpcStateList := []map[string]interface{}{}

		for _, vpcName := range vpcNames {
			vpcDetail, _, err := goe2eClient.Vpcs.GetVPCByName(ctx, vpcName)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to refresh VPC details for %q: %w", vpcName, err))
			}

			subnetList := []map[string]interface{}{}
			for _, sn := range vpcDetail.Subnets {
				subnetList = append(subnetList, map[string]interface{}{
					"id":          sn.ID,
					"subnet_name": sn.SubnetName,
					"cidr":        sn.CIDR,
					"used_ips":    sn.UsedIPs,
					"total_ips":   sn.TotalIPs,
				})
			}

			vpcStateList = append(vpcStateList, map[string]interface{}{
				"name":       vpcDetail.Name,
				"network_id": vpcDetail.NetworkID,
				"ipv4_cidr":  vpcDetail.IPv4CIDR,
				"state":      vpcDetail.State,
				"subnets":    subnetList,
			})
		}

		// Set both V2 and V3 fields
		if err := d.Set("vpc", vpcStateList); err != nil {
			return diag.FromErr(fmt.Errorf("failed to set vpc state: %w", err))
		}
		if err := d.Set("vpc_config", vpcStateList); err != nil {
			return diag.FromErr(fmt.Errorf("failed to set vpc_config state: %w", err))
		}
	}

	// Handle individual public IP changes (only if network_config didn't change)
	if (d.HasChange(tfconstants.AttrPublicIPRequired) || d.HasChange("assign_public_ip")) && !d.HasChange("network_config") {
		var newVal bool
		if d.HasChange("assign_public_ip") {
			newVal = d.Get("assign_public_ip").(bool)
		} else {
			newVal = d.Get(tfconstants.AttrPublicIPRequired).(bool)
		}

		log.Printf("[INFO] assign_public_ip/is_public_ip_required changed to %v", newVal)

		group, _, err := goe2eClient.Autoscaling.GetScalerGroup(ctx, id)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to fetch scaler group status: %w", err))
		}
		if group == nil {
			return diag.Errorf("scaler group not found")
		}
		normalizedStatus := NormalizeStatus(group.ProvisionStatus)
		if !autoscalingScalerGroupStatusIn(normalizedStatus, tfconstants.AutoscalingScalerGroupStoppedStates) {
			return diag.Errorf("ScalerGroup must be in 'Stopped' state to attach/detach public IP. Current: %s", group.ProvisionStatus)
		}

		// Check if VPC is attached
		vpcsRaw, ok := d.GetOk("vpc")
		if !ok {
			vpcsRaw, ok = d.GetOk("vpc_config")
		}
		if !ok || len(vpcsRaw.([]interface{})) == 0 {
			return diag.Errorf("At least one VPC must be attached to attach/detach public IP")
		}

		if newVal {
			log.Printf("[INFO] Triggering Public IP ATTACH")
			_, _, err := goe2eClient.Autoscaling.AttachPublicIPToScalerGroup(ctx, id)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to attach public IP: %w", err))
			}
		} else {
			log.Printf("[INFO] Triggering Public IP DETACH")
			_, _, err := goe2eClient.Autoscaling.DetachPublicIPFromScalerGroup(ctx, id)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to detach public IP: %w", err))
			}
		}
	}

	// Check if there are other changes (size, policies)
	hasSizeChange := d.HasChange(tfconstants.AttrMinNodes) || d.HasChange("min_size") ||
		d.HasChange(tfconstants.AttrMaxNodes) || d.HasChange("max_size")
	hasPolicyChange := d.HasChange("policy_type") || d.HasChange("policy") || d.HasChange("scaling_policy") ||
		d.HasChange("scheduled_policy") || d.HasChange(attrScheduledAction)

	if !hasSizeChange && !hasPolicyChange {
		log.Println("[INFO] No relevant changes detected, skipping update.")
		return nil
	}

	// Expand policies (handle V2 and V3)
	var elasticPolicies []goe2e.ElasticPolicy
	if v, ok := d.GetOk("scaling_policy"); ok {
		// V3 format
		for _, p := range v.([]interface{}) {
			pMap := p.(map[string]interface{})
			policyType := "upscale"
			if pMap["type"].(string) == "scale_down" {
				policyType = "downscale"
			}
			// Map V3 metric names to API parameter names
			metric := pMap["metric"].(string)
			parameter := metric
			if metric == "cpu_utilization" {
				parameter = "cpu"
			} else if metric == "memory_utilization" {
				parameter = "memory"
			}

			elasticPolicies = append(elasticPolicies, goe2e.ElasticPolicy{
				Type:          policyType,
				Adjust:        pMap["adjustment"].(int),
				Parameter:     parameter,
				Operator:      pMap["operator"].(string),
				Value:         pMap["threshold"].(string),
				PeriodNumber:  strconv.Itoa(pMap["evaluation_periods"].(int)),
				PeriodSeconds: strconv.Itoa(pMap["period_seconds"].(int)),
				Cooldown:      strconv.Itoa(pMap["cooldown_seconds"].(int)),
			})
		}
	} else if v, ok := d.GetOk("policy"); ok {
		// V2 format
		for _, p := range v.([]interface{}) {
			pMap := p.(map[string]interface{})
			elasticPolicies = append(elasticPolicies, goe2e.ElasticPolicy{
				Type:          pMap["type"].(string),
				Adjust:        pMap["adjust"].(int),
				Parameter:     pMap["parameter"].(string),
				Operator:      pMap["operator"].(string),
				Value:         pMap["value"].(string),
				PeriodNumber:  pMap["period_number"].(string),
				PeriodSeconds: pMap["period_seconds"].(string),
				Cooldown:      pMap["cooldown"].(string),
			})
		}
	}

	var schedPolicies []goe2e.ScheduledPolicy
	if v, ok := d.GetOk(attrScheduledAction); ok {
		// V3 format
		for _, s := range v.([]interface{}) {
			sMap := s.(map[string]interface{})
			actionType := sMap[attrScheduledActionType].(string)
			var adjust string
			if actionType == goe2econstants.AutoscalingScheduledActionTypeSetCapacity {
				adjust = strconv.Itoa(sMap[attrScheduledTargetCap].(int))
			} else {
				adjust = strconv.Itoa(sMap[attrScheduledAdjustment].(int))
			}
			schedPolicies = append(schedPolicies, goe2e.ScheduledPolicy{
				Type:       actionType,
				Adjust:     adjust,
				Recurrence: sMap["recurrence"].(string),
			})
		}
	} else if v, ok := d.GetOk("scheduled_policy"); ok {
		// V2 format
		for _, s := range v.([]interface{}) {
			sMap := s.(map[string]interface{})
			schedPolicies = append(schedPolicies, goe2e.ScheduledPolicy{
				Type:       sMap["type"].(string),
				Adjust:     sMap["adjust"].(string),
				Recurrence: sMap["recurrence"].(string),
			})
		}
	}

	var policyType string
	if len(elasticPolicies) > 0 {
		policyType = d.Get("policy_type").(string)
		if policyType == "" {
			policyType = "elastic"
		}
	}

	planID := d.Get("plan_id").(string)
	req := &goe2e.ScalerGroupUpdateRequest{
		Name:            d.Get("name").(string),
		PlanID:          planID,
		MinNodes:        minNodes,
		MaxNodes:        maxNodes,
		PolicyType:      policyType,
		Policy:          elasticPolicies,
		ScheduledPolicy: schedPolicies,
	}

	log.Printf("[INFO] Updating ScalerGroup %s with new configuration...", id)
	_, err = goe2eClient.Autoscaling.UpdateScalerGroup(ctx, id, req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update scaler group: %w", err))
	}

	return resourceReadScalerGroup(ctx, d, m)
}

func extractVpcNames(vpcs []interface{}) []string {
	var names []string
	for _, raw := range vpcs {
		if m, ok := raw.(map[string]interface{}); ok {
			if name, ok := m["name"].(string); ok {
				names = append(names, name)
			}
		}
	}
	return names
}

func difference(a, b []string) []string {
	mb := map[string]bool{}
	for _, x := range b {
		mb[x] = true
	}
	var diff []string
	for _, x := range a {
		if !mb[x] {
			diff = append(diff, x)
		}
	}
	return diff
}

func expandIntList(raw []interface{}) []int {
	result := make([]int, len(raw))
	for i, v := range raw {
		result[i] = v.(int)
	}
	return result
}

func intSliceToStringSlice(in []int) []string {
	result := make([]string, len(in))
	for i, v := range in {
		result[i] = strconv.Itoa(v)
	}
	return result
}

// getStatusFromDiff retrieves the status from diff, handling both V2 and V3 field names
func getStatusFromDiff(diff *schema.ResourceDiff) string {
	if v, ok := diff.GetOk("status"); ok {
		return v.(string)
	}
	if v, ok := diff.GetOk("provision_status"); ok {
		return v.(string)
	}
	return ""
}

// getVPCsFromDiff retrieves VPCs from diff, handling both V2 and V3 field names
func getVPCsFromDiff(diff *schema.ResourceDiff) []interface{} {
	if v, ok := diff.GetOk("vpc_config"); ok {
		return v.([]interface{})
	}
	if v, ok := diff.GetOk("vpc"); ok {
		return v.([]interface{})
	}
	// Also check network_config
	if networkConfig := expandNetworkConfigFromDiff(diff); networkConfig != nil && len(networkConfig.VPCNames) > 0 {
		// Return a placeholder to indicate VPCs are present
		return []interface{}{map[string]interface{}{"name": networkConfig.VPCNames[0]}}
	}
	return []interface{}{}
}

// expandNetworkConfigFromDiff extracts network_config from ResourceDiff
func expandNetworkConfigFromDiff(diff *schema.ResourceDiff) *NetworkConfig {
	if v, ok := diff.GetOk("network_config"); ok {
		return expandNetworkConfigFromRaw(v)
	}
	return nil
}

// resourceAutoscalingResourceV0 returns the V0 schema for state migration
// This represents the schema before V3 changes (without V3 field names and structured blocks)
func resourceAutoscalingResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Common fields
			tfconstants.AttrRegion:    config.RegionSchema(),
			tfconstants.AttrLocation:  config.LocationSchema(),
			tfconstants.AttrProjectID: config.ProjectIDSchemaResource(),

			// Core identity
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
			"plan_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sku_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"slug_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// Image (V2 only)
			"vm_image_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vm_image_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vm_template_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			// Security Groups
			"my_account_sg_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			tfconstants.AttrSecurityGroupIDs: {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},

			// Encryption (V2 only)
			tfconstants.AttrIsEncryptionEnabled: {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			tfconstants.AttrEncryptionPassphrase: {
				Type:      schema.TypeString,
				Optional:  true,
				Default:   "",
				ForceNew:  true,
				Sensitive: true,
			},

			// Public IP (V2 only)
			tfconstants.AttrPublicIPRequired: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			// Status (V2 only)
			"provision_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(tfconstants.AutoscalingScalerGroupProvisionStatusAllowed, false),
			},

			// Size (V2 only)
			tfconstants.AttrMinNodes: {
				Type:     schema.TypeInt,
				Required: true,
			},
			tfconstants.AttrMaxNodes: {
				Type:     schema.TypeInt,
				Required: true,
			},
			tfconstants.AttrDesired: {
				Type:     schema.TypeInt,
				Required: true,
			},

			// Policy
			"policy_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type":           {Type: schema.TypeString, Required: true},
						"adjust":         {Type: schema.TypeInt, Required: true},
						"parameter":      {Type: schema.TypeString, Required: true},
						"operator":       {Type: schema.TypeString, Required: true},
						"value":          {Type: schema.TypeString, Required: true},
						"period_number":  {Type: schema.TypeString, Required: true},
						"period_seconds": {Type: schema.TypeString, Required: true},
						"cooldown":       {Type: schema.TypeString, Required: true},
					},
				},
			},
			"scheduled_policy": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type":       {Type: schema.TypeString, Required: true},
						"adjust":     {Type: schema.TypeString, Required: true},
						"recurrence": {Type: schema.TypeString, Required: true},
					},
				},
			},

			// VPC (V2 only)
			"vpc": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"network_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ipv4_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnets": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id":          {Type: schema.TypeInt, Computed: true},
									"subnet_name": {Type: schema.TypeString, Computed: true},
									"cidr":        {Type: schema.TypeString, Computed: true},
									"used_ips":    {Type: schema.TypeInt, Computed: true},
									"total_ips":   {Type: schema.TypeInt, Computed: true},
								},
							},
						},
					},
				},
			},

			// Computed fields
			"running": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

// ResourceAutoscalingStateUpgradeV0toV1 upgrades state from v0 to v1
// This function preserves all V2 fields and adds new V3 fields with defaults
// V2 fields remain functional for backwards compatibility
// Exported for testing purposes
func ResourceAutoscalingStateUpgradeV0toV1(
	ctx context.Context,
	rawState map[string]interface{},
	meta interface{},
) (map[string]interface{}, error) {
	// Add new V3 tags field with default empty map (state-only)
	if _, exists := rawState["tags"]; !exists {
		rawState["tags"] = make(map[string]interface{})
	}

	// Add computed V3 fields with defaults
	if _, exists := rawState["running_node_count"]; !exists {
		// Copy from "running" if it exists
		if running, ok := rawState["running"]; ok {
			rawState["running_node_count"] = running
		} else {
			rawState["running_node_count"] = 0
		}
	}

	if _, exists := rawState["nodes"]; !exists {
		rawState["nodes"] = []interface{}{}
	}

	// Preserve all existing V2 fields
	// No automatic renames - V2 fields remain functional
	// V3 fields will be populated on next read/refresh

	log.Printf("[INFO] Upgraded autoscaling state from v0 to v1: %s", rawState["id"])
	return rawState, nil
}
