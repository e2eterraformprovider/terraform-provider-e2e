package loadbalancer

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

// Import format constants - File-local constants for import validation
const (
	importFormatError = "invalid import format, expected: <lb_id> or <project_id>/<region>/<lb_id>"
)

func ResourceLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Schema:        ResourceLoadBalancerSchema(),
		CreateContext: resourceCreateLoadBalancer,
		ReadContext:   resourceReadLoadBalancer,
		UpdateContext: resourceUpdateLoadBalancer,
		DeleteContext: resourceDeleteLoadBalancer,
		Exists:        resourceExistsLoadBalancer,
		CustomizeDiff: resourceLoadBalancerCustomizeDiff,
		Importer: &schema.ResourceImporter{
			State: customImportStateLoadBalancer,
		},
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceLoadBalancerResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: ResourceLoadBalancerStateUpgradeV0toV1,
				Version: 0,
			},
		},
	}
}

func customImportStateLoadBalancer(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	// Support both formats: <lb_id> or <project_id>/<region>/<lb_id>
	if len(parts) == 1 {
		// Simple format: just LB ID (uses provider defaults for project/region)
		d.SetId(parts[0])
		return []*schema.ResourceData{d}, nil
	} else if len(parts) == 3 {
		// Full format: project_id/region/lb_id
		d.Set(tfconstants.AttrProjectID, parts[0])
		d.Set(tfconstants.AttrRegion, parts[1])
		d.SetId(parts[2])
		return []*schema.ResourceData{d}, nil
	}

	return nil, errors.New(importFormatError)
}

// resourceLoadBalancerCustomizeDiff handles custom diff logic
// Validates field conflicts and emits deprecation warnings
func resourceLoadBalancerCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
	// Emit deprecation warning if location is used
	if _, ok := d.GetOk(tfconstants.AttrLocation); ok {
		log.Printf("[WARN] Parameter 'location' is deprecated and will be removed in v4.0. Please use 'region' instead")
	}

	// Validate that region and location are not both set (handled by ConflictsWith, but double-check)
	if _, hasRegion := d.GetOk(tfconstants.AttrRegion); hasRegion {
		if _, hasLocation := d.GetOk(tfconstants.AttrLocation); hasLocation {
			return fmt.Errorf("cannot set both 'region' and 'location' parameters")
		}
	}

	// Ensure at least one of name or lb_name is provided
	_, hasName := d.GetOk(tfconstants.AttrName)
	_, hasLbName := d.GetOk(tfconstants.AttrLbName)
	if !hasName && !hasLbName {
		return fmt.Errorf("either 'name' or 'lb_name' must be provided")
	}

	// Emit deprecation warning if lb_name is used
	if hasLbName {
		log.Printf("[WARN] Parameter 'lb_name' is deprecated and will be removed in v4.0. Please use 'name' instead")
	}

	// Emit deprecation warning if lb_reserve_ip is used
	if _, ok := d.GetOk("lb_reserve_ip"); ok {
		log.Printf("[WARN] Parameter 'lb_reserve_ip' is deprecated and will be removed in v4.0. Please use 'floating_ip_id' instead")
	}

	return nil
}

func ResourceLoadBalancerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// ============================================
		// COMMON FIELDS
		// ============================================
		tfconstants.AttrRegion: config.RegionSchema(),
		tfconstants.AttrLocation: func() *schema.Schema {
			s := config.LocationSchema()
			s.Deprecated = "The 'location' field is deprecated. Use 'region' instead. This field will be removed in v4.0."
			return s
		}(),
		tfconstants.AttrProjectID: config.ProjectIDSchemaResource(),

		// ============================================
		// REQUIRED INPUT FIELDS (Immutable)
		// ============================================
		tfconstants.AttrPlan: {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "the plan name of the Load Balancer",
			ValidateFunc: validation.StringInSlice([]string{
				goe2econstants.LBPlanE2ELB2,
				goe2econstants.LBPlanE2ELB3,
				goe2econstants.LBPlanE2ELB4,
				goe2econstants.LBPlanE2ELB5,
			}, false),
		},
		tfconstants.AttrName: {
			Type:          schema.TypeString,
			Optional:      true,
			Description:   "name of the Load Balancer (letters, digits, underscores, and hyphens are allowed). This is the recommended field name.",
			ValidateFunc:  node.ValidateName,
			ConflictsWith: []string{tfconstants.AttrLbName},
		},
		tfconstants.AttrLbName: {
			Type:          schema.TypeString,
			Optional:      true,
			Deprecated:    "The 'lb_name' field is deprecated. Use 'name' instead. This field will be removed in v4.0.",
			Description:   "name of the Load Balancer (letters, digits, underscores, and hyphens are allowed). Deprecated: use 'name' instead.",
			ValidateFunc:  node.ValidateName,
			ConflictsWith: []string{tfconstants.AttrName},
		},
		tfconstants.AttrLbMode: {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "the mode of the Load Balancer (HTTP, HTTPS, or Both)",
			ValidateFunc: validation.StringInSlice([]string{
				goe2econstants.LBModeHTTP,
				goe2econstants.LBModeHTTPS,
				goe2econstants.LBModeBoth,
			}, false),
		},

		// ============================================
		// OPTIONAL INPUT FIELDS - CREATION (Immutable)
		// ============================================
		tfconstants.AttrLbType: {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     goe2econstants.LBTypeExternal,
			ForceNew:    true,
			Description: "the type of Load Balancer (Internal or External)",
		},
		"node_list_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     goe2econstants.LBNodeListTypeStatic,
			ForceNew:    true,
			Description: "the node list type (S for static nodes, D for dynamic autoscaling)",
			ValidateFunc: validation.StringInSlice([]string{
				goe2econstants.LBNodeListTypeStatic,
				goe2econstants.LBNodeListTypeDynamic,
			}, false),
		},
		"floating_ip_id": {
			Type:          schema.TypeString,
			Optional:      true,
			Description:   "id of the floating IP to attach to the Load Balancer. This is the recommended field name.",
			ConflictsWith: []string{"lb_reserve_ip"},
		},
		"lb_reserve_ip": {
			Type:          schema.TypeString,
			Optional:      true,
			Default:       "",
			Deprecated:    "The 'lb_reserve_ip' field is deprecated. Use 'floating_ip_id' instead. This field will be removed in v4.0.",
			Description:   "id of the reserved IP to attach to the Load Balancer. Deprecated: use 'floating_ip_id' instead.",
			ConflictsWith: []string{"floating_ip_id"},
		},
		"enable_bitninja": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "whether to enable BitNinja security for the Load Balancer",
		},

		// ============================================
		// OPTIONAL INPUT FIELDS - CONFIGURATION (Mutable)
		// ============================================
		"checkbox_enable": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "checkbox configuration option",
		},
		"ssl_certificate_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "id of the SSL certificate (required if lb_mode is HTTPS or Both)",
		},
		"ssl_context": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "SSL context configuration for the Load Balancer",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"redirect_to_https": {
						Type:        schema.TypeBool,
						Optional:    true,
						Default:     false,
						Description: "whether to redirect HTTP to HTTPS (required if Load Balancer is set to Both)",
					},
				},
			},
		},
		tfconstants.AttrBackends: {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			Description: "list of backend details to attach to the Load Balancer",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "name of the backend",
					},
					"scaler_id": {
						Type:        schema.TypeString,
						Optional:    true,
						Default:     "",
						Description: "id of the scaler group to attach (if using autoscaling)",
					},
					"scaler_port": {
						Type:        schema.TypeString,
						Optional:    true,
						Default:     "",
						Description: "port of the scaler group to attach (if using autoscaling)",
					},
					"balance": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "the load balancing algorithm (source, roundrobin, or leastconn)",
						ValidateFunc: validation.StringInSlice([]string{
							goe2econstants.LBBalanceSource,
							goe2econstants.LBBalanceRoundRobin,
							goe2econstants.LBBalanceLeastConn,
						}, false),
					},
					"checkbox_enable": {
						Type:        schema.TypeBool,
						Optional:    true,
						Default:     false,
						Description: "whether to enable healthcheck",
					},
					"domain_name": {
						Type:        schema.TypeString,
						Optional:    true,
						Default:     goe2econstants.LBDefaultDomainName,
						Description: "the domain name for healthcheck",
					},
					"check_url": {
						Type:        schema.TypeString,
						Optional:    true,
						Default:     goe2econstants.LBDefaultCheckURL,
						Description: "the endpoint URL for healthcheck",
					},
					"servers": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "list of servers to attach to the backend",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"node_id": {
									Type:          schema.TypeString,
									Optional:      true,
									Description:   "id of the Node to attach. This is the recommended field name.",
									ConflictsWith: []string{"id"},
								},
								"id": {
									Type:          schema.TypeString,
									Optional:      true,
									Deprecated:    "The 'id' field in backend servers is deprecated. Use 'node_id' instead. This field will be removed in v4.0.",
									Description:   "id of the Node to attach. Deprecated: use 'node_id' instead.",
									ConflictsWith: []string{"node_id"},
								},
								"port": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "port number of the Node",
								},
							},
						},
					},
					"http_check": {
						Type:        schema.TypeBool,
						Optional:    true,
						Default:     false,
						Description: "whether HTTP health check is enabled",
					},
				},
			},
		},
		"acl_list": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "list of ACL rules to apply",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"acl_name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "name of the ACL rule",
					},
					"acl_condition": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "the condition in which the ACL rule will match",
					},
					"acl_matching_path": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "the path in which this rule will work",
					},
				},
			},
		},

		"acl_map": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "list of ACL routing rules to route requests according to ACL rules",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"acl_name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "name of the ACL rule",
					},
					"acl_condition_state": {
						Type:        schema.TypeBool,
						Optional:    true,
						Default:     true,
						Description: "whether the ACL condition state is enabled",
					},
					"acl_backend": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "name of the backend server",
					},
				},
			},
		},
		"vpc_list": {
			Type:        schema.TypeSet,
			Elem:        &schema.Schema{Type: schema.TypeInt},
			Optional:    true,
			Description: "list of VPC ids to attach to the Load Balancer",
		},
		"enable_eos_logger": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "configuration to enable EOS bucket logging for the Load Balancer",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"appliance_id": {
						Type:        schema.TypeInt,
						Optional:    true,
						Default:     0,
						Description: "id of the appliance",
					},
					"access_key": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "the access key of the Object Store bucket",
					},
					"secret_key": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "the secret key of the Object Store bucket",
					},
					"bucket": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "the bucket name of the Object Store",
					},
				},
			},
		},
		"tcp_backend": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "list of TCP backend configurations to attach",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"backend_name": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "name of the TCP backend",
						ValidateFunc: node.ValidateName,
					},
					"port": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "port number for the TCP backend (ports 8080, 10050, 9101, 80, or 443 are not allowed)",
						ValidateFunc: validation.StringNotInSlice(tfconstants.LBTCPDisallowedPorts, false),
					},
					"balance": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "the load balancing algorithm (source, roundrobin, or leastconn)",
						ValidateFunc: validation.StringInSlice([]string{
							goe2econstants.LBBalanceSource,
							goe2econstants.LBBalanceRoundRobin,
							goe2econstants.LBBalanceLeastConn,
						}, false),
					},
					"servers": {
						Type:        schema.TypeList,
						Required:    true,
						Description: "list of servers to attach to the TCP backend",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"node_id": {
									Type:          schema.TypeString,
									Optional:      true,
									Description:   "id of the Node to attach. This is the recommended field name.",
									ValidateFunc:  node.ValidateName,
									ConflictsWith: []string{"id"},
								},
								"id": {
									Type:          schema.TypeString,
									Optional:      true,
									Deprecated:    "The 'id' field in TCP backend servers is deprecated. Use 'node_id' instead. This field will be removed in v4.0.",
									Description:   "id of the Node to attach. Deprecated: use 'node_id' instead.",
									ValidateFunc:  node.ValidateName,
									ConflictsWith: []string{"node_id"},
								},
								"port": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "port number of the Node",
								},
							},
						},
					},
				},
			},
		},
		// ============================================
		// OPTIONAL INPUT FIELDS - MANAGEMENT (Mutable)
		// ============================================
		tfconstants.AttrPowerStatus: {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     goe2econstants.NodePowerStatusOn,
			Description: "the power state of the Load Balancer (power_on to start, power_off to power off)",
		},
		"is_ipv6_attached": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "whether IPv6 is attached to the Load Balancer",
		},
		"default_backend": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "the default backend name for the Load Balancer",
		},
		tfconstants.AttrTags: {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "tags to apply to the Load Balancer (state-only until API support is added)",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},

		// ============================================
		// COMPUTED FIELDS - STATUS
		// ============================================
		tfconstants.AttrStatus: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "state of the Load Balancer instance (API value: Creating, Running, Powered off, etc.)",
		},
		tfconstants.AttrState: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "normalized state of the Load Balancer (creating, running, stopped, etc.)",
		},

		// ============================================
		// COMPUTED FIELDS - NETWORK
		// ============================================
		tfconstants.AttrPublicIPAddress: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "the Load Balancer's public IPv4 address. This is the recommended field name.",
		},
		tfconstants.AttrPrivateIPAddress: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "the Load Balancer's private IPv4 address. This is the recommended field name.",
		},
		"public_ip": {
			Type:        schema.TypeString,
			Computed:    true,
			Deprecated:  "The 'public_ip' field is deprecated. Use 'public_ip_address' instead. This field will be removed in v4.0.",
			Description: "the Load Balancer's public IPv4 address. Deprecated: use 'public_ip_address' instead.",
		},
		"private_ip": {
			Type:        schema.TypeString,
			Computed:    true,
			Deprecated:  "The 'private_ip' field is deprecated. Use 'private_ip_address' instead. This field will be removed in v4.0.",
			Description: "the Load Balancer's private IPv4 address. Deprecated: use 'private_ip_address' instead.",
		},
		"host_target_ipv6": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "the IPv6 address allocated to the Load Balancer",
		},

		// ============================================
		// COMPUTED FIELDS - RESOURCES
		// ============================================
		tfconstants.AttrRAM: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "the RAM allocated to the Load Balancer",
		},
		tfconstants.AttrDisk: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "the disk storage allocated to the Load Balancer",
		},
		tfconstants.AttrVCPU: {
			Type:        schema.TypeFloat,
			Computed:    true,
			Description: "the number of virtual CPUs allocated to the Load Balancer",
		},
	}
}

// CreateLoadBalancerObjectWithGoe2e creates a goe2e LoadBalancerCreateRequest from schema data
func CreateLoadBalancerObjectWithGoe2e(ctx context.Context, goe2eClient *goe2e.Client, d *schema.ResourceData, region string, projectID string) (*goe2e.LoadBalancerCreateRequest, diag.Diagnostics) {
	log.Printf("[INFO] LOAD BALANCER OBJECT CREATION STARTS (goe2e)")

	// Handle name vs lb_name (prefer name, fallback to lb_name)
	var lbName string
	if name, ok := d.GetOk(tfconstants.AttrName); ok {
		lbName = name.(string)
	} else if lbNameVal, ok := d.GetOk(tfconstants.AttrLbName); ok {
		lbName = lbNameVal.(string)
	}

	// Handle floating_ip_id vs lb_reserve_ip (prefer floating_ip_id, fallback to lb_reserve_ip)
	var lbReserveIP string
	if floatingIPID, ok := d.GetOk("floating_ip_id"); ok {
		lbReserveIP = floatingIPID.(string)
	} else if lbReserveIPVal, ok := d.GetOk("lb_reserve_ip"); ok {
		lbReserveIP = lbReserveIPVal.(string)
	}

	loadBalancerObj := &goe2e.LoadBalancerCreateRequest{
		PlanName:         d.Get(tfconstants.AttrPlan).(string),
		LBName:           lbName,
		LBType:           d.Get(tfconstants.AttrLbType).(string),
		LBMode:           d.Get(tfconstants.AttrLbMode).(string),
		LBPort:           GetLbPort(d.Get(tfconstants.AttrLbMode).(string)),
		NodeListType:     d.Get("node_list_type").(string),
		CheckBoxEnable:   d.Get("checkbox_enable").(string),
		LBReserveIP:      lbReserveIP,
		SSLCertificateID: d.Get("ssl_certificate_id").(string),
		EnableBitNinja:   d.Get("enable_bitninja").(bool),
		IsIPv6Attached:   d.Get("is_ipv6_attached").(bool),
		DefaultBackend:   d.Get("default_backend").(string),
		Location:         region,
	}

	enableEosLogger, ok := d.GetOk("enable_eos_logger")
	if ok {
		eosDetail, err := ExpandEnableEosLogger(enableEosLogger.([]interface{}))
		if err != nil {
			return nil, diag.FromErr(err)
		}
		loadBalancerObj.EnableEOSLogger = eosDetail
	}

	aclList, ok := d.GetOk("acl_list")
	if ok {
		aclListDetail, err := ExpandAclList(aclList.([]interface{}))
		if err != nil {
			return nil, diag.FromErr(err)
		}
		loadBalancerObj.ACLList = aclListDetail
	} else {
		loadBalancerObj.ACLList = make([]goe2e.LBACLList, 0)
	}

	aclMap, ok := d.GetOk("acl_map")
	if ok {
		aclMapDetail, err := ExpandAclMap(aclMap.([]interface{}))
		if err != nil {
			return nil, diag.FromErr(err)
		}
		loadBalancerObj.ACLMap = aclMapDetail
	} else {
		loadBalancerObj.ACLMap = make([]goe2e.LBACLMap, 0)
	}

	tcpBackend, ok := d.GetOk("tcp_backend")
	if ok {
		tcpBackendDetail, err := ExpandTcpBackendWithGoe2e(ctx, tcpBackend.([]interface{}), goe2eClient, projectID, region)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		loadBalancerObj.TCPBackend = tcpBackendDetail
	} else {
		loadBalancerObj.TCPBackend = make([]goe2e.LBTCPBackend, 0)
	}

	backends, ok := d.GetOk(tfconstants.AttrBackends)
	if ok {
		backendDetail, err := ExpandBackendsWithGoe2e(ctx, backends.([]interface{}), goe2eClient, projectID, region)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		loadBalancerObj.Backends = backendDetail
	} else {
		loadBalancerObj.Backends = make([]goe2e.LBBackend, 0)
	}

	vpcList, ok := d.GetOk("vpc_list")
	if ok {
		vpcListDetail, err := ExpandVpcListWithGoe2e(ctx, d, vpcList.(*schema.Set).List(), goe2eClient)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		loadBalancerObj.VPCList = vpcListDetail
	} else {
		loadBalancerObj.VPCList = make([]goe2e.LBVPCDetail, 0)
	}

	sslContext, ok := d.GetOk("ssl_context")
	if ok {
		sslContextList := sslContext.([]interface{})
		detail := sslContextList[0].(map[string]interface{})
		loadBalancerObj.SSLContext = detail
	} else {
		loadBalancerObj.SSLContext = map[string]interface{}{"redirect_to_https": false}
	}

	return loadBalancerObj, nil
}

func resourceCreateLoadBalancer(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	// Use goe2e client
	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return diag.Errorf("Error creating goe2e client: %s", err)
	}

	loadBalancerObj, diags := CreateLoadBalancerObjectWithGoe2e(ctx, goe2eClient, d, region, projectID)
	if diags != nil {
		return diags
	}

	// Get lbName for error messages
	var lbName string
	if name, ok := d.GetOk(tfconstants.AttrName); ok {
		lbName = name.(string)
	} else if lbNameVal, ok := d.GetOk(tfconstants.AttrLbName); ok {
		lbName = lbNameVal.(string)
	}

	lb, _, err := goe2eClient.LoadBalancer.CreateLoadBalancer(ctx, loadBalancerObj)
	if err != nil {
		return diag.Errorf("Error creating load balancer (name: %s) in project (%s), region (%s): %s", lbName, projectID, region, err)
	}
	log.Printf("[INFO] LOAD BALANCER CREATE | RESPONSE | %+v", lb)

	if !lb.IsCreditSufficient {
		return diag.Errorf("Cannot create load balancer (name: %s) in project (%s), region (%s): insufficient credits. Please add credits to your account", lbName, projectID, region)
	}
	log.Printf("[INFO] load balancer creation | before setting fields")

	// Set ID (goe2e client returns ID as string)
	d.SetId(lb.ID)

	// Wait for load balancer to reach Running status
	log.Printf("[INFO] Waiting for load balancer %s to reach Running status", lb.ID)
	if err := waitForLoadBalancerStatus(ctx, goe2eClient, lb.ID, goe2econstants.LBStateRunning, int(tfconstants.LBCreateTimeout.Minutes())); err != nil {
		return diag.Errorf("Error waiting for load balancer to become ready: %s", err)
	}

	// Set both new and deprecated computed fields for backwards compatibility
	if lb.PublicIPAddress != "" {
		d.Set(tfconstants.AttrPublicIPAddress, lb.PublicIPAddress)
		d.Set("public_ip", lb.PublicIPAddress) // Deprecated field
	}

	// Store tags in state if provided
	if tags, ok := d.GetOk(tfconstants.AttrTags); ok {
		if err := d.Set(tfconstants.AttrTags, tags); err != nil {
			return diag.FromErr(fmt.Errorf("error setting tags: %w", err))
		}
	}

	return diags
}

func resourceReadLoadBalancer(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	var diags diag.Diagnostics

	log.Printf("=============INSIDE RESOURCE READ LOAD BALANCER==========================")
	lbId := d.Id()

	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	projectID, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Use goe2e client with typed nested response
	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return diag.Errorf("Error creating goe2e client: %s", err)
	}

	apiResponse, err := getLoadBalancerWithNestedResponse(ctx, goe2eClient, lbId)
	if err != nil {
		return diag.Errorf("Error retrieving load balancer (ID: %s) in project (%s), region (%s): %s", lbId, projectID, region, err)
	}

	log.Printf("[INFO] LOADBALANCER READ | BEFORE SETTING DATA")

	// Set both new and deprecated computed fields for backwards compatibility
	if apiResponse.Data.NodeDetail.PrivateIP != "" {
		d.Set(tfconstants.AttrPrivateIPAddress, apiResponse.Data.NodeDetail.PrivateIP)
		d.Set("private_ip", apiResponse.Data.NodeDetail.PrivateIP) // Deprecated field
	}
	if apiResponse.Data.NodeDetail.PublicIP != "" {
		d.Set(tfconstants.AttrPublicIPAddress, apiResponse.Data.NodeDetail.PublicIP)
		d.Set("public_ip", apiResponse.Data.NodeDetail.PublicIP) // Deprecated field
	}
	d.Set(tfconstants.AttrRAM, apiResponse.Data.NodeDetail.RAM)
	d.Set(tfconstants.AttrDisk, apiResponse.Data.NodeDetail.Disk)
	d.Set(tfconstants.AttrVCPU, apiResponse.Data.NodeDetail.VCPU)

	// Set both name and lb_name for backwards compatibility
	if apiResponse.Data.Name != "" {
		d.Set(tfconstants.AttrName, apiResponse.Data.Name)
		d.Set(tfconstants.AttrLbName, apiResponse.Data.Name) // Deprecated field
	}

	d.Set(tfconstants.AttrPlan, apiResponse.Data.NodeDetail.PlanName)

	// Extract lb_mode and host_target_ipv6 from appliance_instance[0].context
	if len(apiResponse.Data.ApplianceInstance) > 0 {
		context := apiResponse.Data.ApplianceInstance[0].Context
		d.Set(tfconstants.AttrLbMode, context.LBMode)

		if d.Get("is_ipv6_attached").(bool) {
			if context.HostTargetIPv6 != "" {
				d.Set("host_target_ipv6", context.HostTargetIPv6)
			} else {
				d.Set("is_ipv6_attached", false)
			}
		}
	}

	// Set status using SetLoadBalancerStatus helper (needs map for now)
	lbStatusMap := map[string]interface{}{
		"status": apiResponse.Data.LBStatus.Status,
		"data_monitor": map[string]interface{}{
			"status": apiResponse.Data.LBStatus.DataMonitor.Status,
		},
	}
	err = SetLoadBalancerStatus(d, lbStatusMap)
	if err != nil {
		return diag.Errorf("Error setting load balancer status for ID (%s) in project (%s), region (%s): %s", lbId, projectID, region, err)
	}

	// Set normalized state field
	if status, ok := d.GetOk("status"); ok {
		statusStr := status.(string)
		d.Set(tfconstants.AttrState, normalizeLoadBalancerState(statusStr))
	}

	if d.Get("status").(string) == goe2econstants.LBStatusPoweredOff {
		d.Set(tfconstants.AttrPowerStatus, goe2econstants.NodePowerStatusOff)
	} else {
		d.Set(tfconstants.AttrPowerStatus, goe2econstants.NodePowerStatusOn)
	}

	// Preserve tags in state (state-only until API support)
	if tags, ok := d.GetOk(tfconstants.AttrTags); ok {
		d.Set(tfconstants.AttrTags, tags)
	}

	return diags
}

func resourceUpdateLoadBalancer(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)

	lbId := d.Id()

	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	projectID, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	lb_status := d.Get("status").(string)

	// Use goe2e client
	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return diag.Errorf("Error creating goe2e client: %s", err)
	}

	// Get current load balancer info for plan upgrade check and status checks
	apiResponse, err := getLoadBalancerWithNestedResponse(ctx, goe2eClient, lbId)
	if err != nil {
		return diag.Errorf("Error retrieving load balancer (ID: %s) in project (%s), region (%s): %s", lbId, projectID, region, err)
	}

	if d.HasChange(tfconstants.AttrPowerStatus) {
		disablePowerStatusList := []string{goe2econstants.LBStatusCreating, goe2econstants.LBStatusDeploying, goe2econstants.LBStatusUpgrading}

		if CheckStatus(disablePowerStatusList, lb_status) {
			return diag.Errorf("Cannot change power status for load balancer (ID: %s): load balancer is in %s state in project (%s), region (%s). Power status can only be changed when load balancer is in %s or %s state", lbId, lb_status, projectID, region, goe2econstants.LBStatusRunning, goe2econstants.LBStatusPoweredOff)
		}

		actionReq := &goe2e.LoadBalancerActionRequest{
			Type: d.Get(tfconstants.AttrPowerStatus).(string),
		}
		_, err = goe2eClient.LoadBalancer.UpdateLoadBalancerAction(ctx, lbId, actionReq)
		if err != nil {
			return diag.Errorf("Error updating power status for load balancer (ID: %s) in project (%s), region (%s): %s", lbId, projectID, region, err)
		}
		// Wait for power action to complete
		var targetStatus string
		if d.Get(tfconstants.AttrPowerStatus).(string) == goe2econstants.NodePowerStatusOn {
			targetStatus = goe2econstants.LBStateRunning
		} else {
			targetStatus = goe2econstants.LBStateStopped
		}
		if err := waitForLoadBalancerStatus(ctx, goe2eClient, lbId, targetStatus, int(tfconstants.LBPowerActionTimeout.Minutes())); err != nil {
			return diag.Errorf("Error waiting for load balancer power action to complete: %s", err)
		}
		return resourceReadLoadBalancer(ctx, d, m)
	}

	if d.HasChange(tfconstants.AttrPlan) {
		currentPlanName := apiResponse.Data.NodeDetail.PlanName
		newPlanName := d.Get(tfconstants.AttrPlan).(string)
		currentNum, err := extractPlanNumber(currentPlanName)
		if err != nil {
			return diag.FromErr(err)
		}
		newNum, err := extractPlanNumber(newPlanName)
		if err != nil {
			return diag.FromErr(err)
		}
		if newNum < currentNum {
			return diag.Errorf("Cannot downgrade plan for load balancer (ID: %s) from %s to %s in project (%s), region (%s): plan downgrades are not supported. Please specify a plan equal to or higher than the current plan", lbId, currentPlanName, newPlanName, projectID, region)
		}
		actionReq := &goe2e.LoadBalancerActionRequest{
			Type:     goe2econstants.LBActionUpgradePlan,
			Name:     apiResponse.Data.Name,
			PlanName: newPlanName,
		}
		_, err = goe2eClient.LoadBalancer.UpdateLoadBalancerAction(ctx, lbId, actionReq)
		if err != nil {
			return diag.Errorf("Error upgrading plan for load balancer (ID: %s) from %s to %s in project (%s), region (%s): %s", lbId, currentPlanName, newPlanName, projectID, region, err)
		}
		// Wait for upgrade to complete
		if err := waitForLoadBalancerStatus(ctx, goe2eClient, lbId, goe2econstants.LBStateRunning, int(tfconstants.LBPlanUpgradeTimeout.Minutes())); err != nil {
			return diag.Errorf("Error waiting for load balancer plan upgrade to complete: %s", err)
		}
		return resourceReadLoadBalancer(ctx, d, m)
	}

	if lb_status == goe2econstants.LBStatusPoweredOff {
		return diag.Errorf("Cannot update load balancer (ID: %s): load balancer is in %s state in project (%s), region (%s). Load balancer must be powered on to update configuration", lbId, lb_status, projectID, region)
	}

	// Handle name or lb_name changes
	if d.HasChange(tfconstants.AttrName) || d.HasChange(tfconstants.AttrLbName) {
		var newName string
		if name, ok := d.GetOk(tfconstants.AttrName); ok {
			newName = name.(string)
		} else if lbName, ok := d.GetOk(tfconstants.AttrLbName); ok {
			newName = lbName.(string)
		}

		actionReq := &goe2e.LoadBalancerActionRequest{
			Type: goe2econstants.LBActionRename,
			Name: newName,
		}
		_, err = goe2eClient.LoadBalancer.UpdateLoadBalancerAction(ctx, lbId, actionReq)
		if err != nil {
			return diag.Errorf("Error renaming load balancer (ID: %s) in project (%s), region (%s): %s", lbId, projectID, region, err)
		}
		// Rename is immediate, no need to wait
		return resourceReadLoadBalancer(ctx, d, m)
	}

	// Handle tags updates (state-only)
	if d.HasChange(tfconstants.AttrTags) {
		// Tags are state-only, just update state
		if tags, ok := d.GetOk(tfconstants.AttrTags); ok {
			d.Set(tfconstants.AttrTags, tags)
		}
		// Continue to other updates
	}

	if d.HasChange("is_ipv6_attached") {
		ipv6_attach := d.Get("is_ipv6_attached").(bool)
		var ipv6Req *goe2e.IPv6ActionRequest
		if ipv6_attach {
			ipv6Req = &goe2e.IPv6ActionRequest{
				Action: goe2econstants.LBIPv6ActionAttach,
			}
		} else {
			// Get IPv6 address from appliance instance context
			var detachIPv6 string
			if len(apiResponse.Data.ApplianceInstance) > 0 {
				detachIPv6 = apiResponse.Data.ApplianceInstance[0].Context.HostTargetIPv6
			}
			ipv6Req = &goe2e.IPv6ActionRequest{
				Action:     goe2econstants.LBIPv6ActionDetach,
				DetachIPv6: detachIPv6,
			}
		}
		_, err = goe2eClient.LoadBalancer.UpdateIPv6(ctx, lbId, ipv6Req)
		if err != nil {
			return diag.Errorf("Error updating IPv6 configuration for load balancer (ID: %s) in project (%s), region (%s): %s", lbId, projectID, region, err)
		}
		return resourceReadLoadBalancer(ctx, d, m)
	}

	loadBalancerObj, diags := CreateLoadBalancerObjectWithGoe2e(ctx, goe2eClient, d, region, projectID)
	if diags != nil {
		return diags
	}

	// Convert to UpdateRequest
	updateReq := &goe2e.LoadBalancerUpdateRequest{
		CheckBoxEnable:   loadBalancerObj.CheckBoxEnable,
		SSLCertificateID: loadBalancerObj.SSLCertificateID,
		SSLContext:       loadBalancerObj.SSLContext,
		Backends:         loadBalancerObj.Backends,
		ACLList:          loadBalancerObj.ACLList,
		ACLMap:           loadBalancerObj.ACLMap,
		VPCList:          loadBalancerObj.VPCList,
		EnableEOSLogger:  loadBalancerObj.EnableEOSLogger,
		TCPBackend:       loadBalancerObj.TCPBackend,
		DefaultBackend:   loadBalancerObj.DefaultBackend,
	}

	res, _, err := goe2eClient.LoadBalancer.UpdateLoadBalancer(ctx, lbId, updateReq)
	if err != nil {
		return diag.Errorf("Error updating backend configuration for load balancer (ID: %s) in project (%s), region (%s): %s", lbId, projectID, region, err)
	}

	if !res.IsCreditSufficient {
		return diag.Errorf("Cannot update load balancer (ID: %s) in project (%s), region (%s): insufficient credits. Please add credits to your account", lbId, projectID, region)
	}
	return resourceReadLoadBalancer(ctx, d, m)
}

func resourceDeleteLoadBalancer(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	var diags diag.Diagnostics
	lbId := d.Id()

	region, err := cfg.GetRegionOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	projectID, err := cfg.GetProjectIDOrDefault(d)
	if err != nil {
		return diag.FromErr(err)
	}

	lb_status := d.Get("status").(string)
	disableDeleteLbStatusList := []string{goe2econstants.LBStatusCreating, goe2econstants.LBStatusDeploying, goe2econstants.LBStatusUpgrading}

	if CheckStatus(disableDeleteLbStatusList, lb_status) {
		return diag.Errorf("Cannot delete load balancer (ID: %s): load balancer is in %s state in project (%s), region (%s). Load balancer can only be deleted when not in %s, %s, or %s state", lbId, lb_status, projectID, region, goe2econstants.LBStatusCreating, goe2econstants.LBStatusDeploying, goe2econstants.LBStatusUpgrading)
	}

	// Use goe2e client
	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return diag.Errorf("Error creating goe2e client: %s", err)
	}

	_, err = goe2eClient.LoadBalancer.DeleteLoadBalancer(ctx, lbId)
	if err != nil {
		return diag.Errorf("Error deleting load balancer (ID: %s) in project (%s), region (%s): %s", lbId, projectID, region, err)
	}

	// Wait for load balancer deletion to complete
	log.Printf("[INFO] Waiting for load balancer %s deletion to complete", lbId)
	if err := waitForLoadBalancerDeletion(ctx, goe2eClient, lbId, int(tfconstants.LBDeleteTimeout.Minutes())); err != nil {
		// Log warning but don't fail - deletion may still be in progress
		log.Printf("[WARN] Timeout waiting for load balancer deletion, but deletion may still be in progress: %s", err)
	}

	d.SetId("")
	return diags
}

func resourceExistsLoadBalancer(d *schema.ResourceData, m interface{}) (bool, error) {
	return true, nil
}

// extractPlanNumber extracts the numeric part from plan names like "E2E-LB-2" -> 2
func extractPlanNumber(planName string) (int, error) {
	parts := strings.Split(planName, "-")
	if len(parts) < 3 {
		return 0, fmt.Errorf("invalid plan name format: %s", planName)
	}
	return strconv.Atoi(parts[2])
}

// resourceLoadBalancerResourceV0 returns the V0 schema for state migration
func resourceLoadBalancerResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			tfconstants.AttrRegion:    config.RegionSchema(),
			tfconstants.AttrLocation:  config.LocationSchema(),
			tfconstants.AttrProjectID: config.ProjectIDSchemaResource(),
			tfconstants.AttrPlan: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			tfconstants.AttrLbName: {
				Type:     schema.TypeString,
				Required: true,
			},
			tfconstants.AttrLbMode: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			tfconstants.AttrLbType: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "External",
				ForceNew: true,
			},
			"node_list_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "S",
				ForceNew: true,
			},
			"lb_reserve_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"enable_bitninja": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			tfconstants.AttrStatus: {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			tfconstants.AttrPublicIPAddress: {
				Type:     schema.TypeString,
				Computed: true,
			},
			tfconstants.AttrPrivateIPAddress: {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_target_ipv6": {
				Type:     schema.TypeString,
				Computed: true,
			},
			tfconstants.AttrRAM: {
				Type:     schema.TypeString,
				Computed: true,
			},
			tfconstants.AttrDisk: {
				Type:     schema.TypeString,
				Computed: true,
			},
			tfconstants.AttrVCPU: {
				Type:     schema.TypeFloat,
				Computed: true,
			},
		},
	}
}

// ResourceLoadBalancerStateUpgradeV0toV1 upgrades state from V0 to V1
// Renames deprecated fields to new fields, adds new computed fields
func ResourceLoadBalancerStateUpgradeV0toV1(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	// Rename location to region if location exists and region doesn't
	if location, ok := rawState[tfconstants.AttrLocation].(string); ok && location != "" {
		if _, hasRegion := rawState[tfconstants.AttrRegion]; !hasRegion || rawState[tfconstants.AttrRegion] == "" {
			rawState[tfconstants.AttrRegion] = location
		}
		// Keep location in state for backwards compatibility
	}

	// Rename lb_name to name if lb_name exists and name doesn't
	if lbName, ok := rawState[tfconstants.AttrLbName].(string); ok && lbName != "" {
		if _, hasName := rawState[tfconstants.AttrName]; !hasName || rawState[tfconstants.AttrName] == "" {
			rawState[tfconstants.AttrName] = lbName
		}
		// Keep lb_name in state for backwards compatibility
	}

	// Rename lb_reserve_ip to floating_ip_id if lb_reserve_ip exists and floating_ip_id doesn't
	if lbReserveIP, ok := rawState["lb_reserve_ip"].(string); ok && lbReserveIP != "" {
		if _, hasFloatingIPID := rawState["floating_ip_id"]; !hasFloatingIPID || rawState["floating_ip_id"] == "" {
			rawState["floating_ip_id"] = lbReserveIP
		}
		// Keep lb_reserve_ip in state for backwards compatibility
	}

	// Rename public_ip to public_ip_address if public_ip exists and public_ip_address doesn't
	if publicIP, ok := rawState["public_ip"].(string); ok && publicIP != "" {
		if _, hasPublicIPAddress := rawState[tfconstants.AttrPublicIPAddress]; !hasPublicIPAddress || rawState[tfconstants.AttrPublicIPAddress] == "" {
			rawState[tfconstants.AttrPublicIPAddress] = publicIP
		}
		// Keep public_ip in state for backwards compatibility
	}

	// Rename private_ip to private_ip_address if private_ip exists and private_ip_address doesn't
	if privateIP, ok := rawState["private_ip"].(string); ok && privateIP != "" {
		if _, hasPrivateIPAddress := rawState[tfconstants.AttrPrivateIPAddress]; !hasPrivateIPAddress || rawState[tfconstants.AttrPrivateIPAddress] == "" {
			rawState[tfconstants.AttrPrivateIPAddress] = privateIP
		}
		// Keep private_ip in state for backwards compatibility
	}

	// Add new computed fields with default values
	if _, ok := rawState[tfconstants.AttrState]; !ok {
		// Normalize status to state if available
		if status, ok := rawState[tfconstants.AttrStatus].(string); ok {
			rawState[tfconstants.AttrState] = normalizeLoadBalancerState(status)
		} else {
			rawState[tfconstants.AttrState] = ""
		}
	}

	if _, ok := rawState[tfconstants.AttrTags]; !ok {
		rawState[tfconstants.AttrTags] = map[string]interface{}{}
	}

	return rawState, nil
}

// normalizeLoadBalancerState normalizes the API status to a normalized state value
func normalizeLoadBalancerState(status string) string {
	switch status {
	case goe2econstants.LBStatusCreating, goe2econstants.LBStatusDeploying:
		return goe2econstants.LBStateCreating
	case goe2econstants.LBStatusRunning:
		return goe2econstants.LBStateRunning
	case goe2econstants.LBStatusPoweredOff:
		return goe2econstants.LBStateStopped
	case goe2econstants.LBStatusUpgrading:
		return goe2econstants.LBStateUpgrading
	case goe2econstants.LBStatusError, "Failed":
		return goe2econstants.LBStateError
	default:
		return strings.ToLower(status)
	}
}
