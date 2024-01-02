package loadbalancer

import (
	"context"
	"log"
	"math"
	"strconv"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/node"
	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var nameRegex string = "^[a-zA-Z0-9-_]{0,50}$"

func ResourceLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Schema:        ResouceLoadBalancerSchema(),
		CreateContext: resourceCreateLoadBalancer,
		ReadContext:   resourceReadLoadBalancer,
		UpdateContext: resourceUpdateLoadBalancer,
		DeleteContext: resourceDeleteLoadBalancer,
		Exists:        resourceExistsLoadBalancer,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func ResouceLoadBalancerSchema() map[string]*schema.Schema {
	// need to implement acl list and acl map
	return map[string]*schema.Schema{
		"plan_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "It is the plan of which load balancer is going to launch",
			ValidateFunc: validation.StringInSlice([]string{
				"E2E-LB-2",
				"E2E-LB-3",
				"E2E-LB-4",
				"E2E-LB-5",
			}, false),
		},
		"lb_name": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "It is the name of load balancer, letter,digit,underscore,hyphen are allowed",
			ValidateFunc: node.ValidateName,
		},
		"lb_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "External",
			Description: "It is used to define internal or extenal load balancer",
		},
		"lb_mode": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "In which mode load balancer is going to launch http, https, both",
			ValidateFunc: validation.StringInSlice([]string{
				"HTTP",
				"HTTPS",
				"Both",
			}, false),
		},
		"node_list_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "S",
			Description: "It is used to find out either node is static or dynamic autoscaling",
			ValidateFunc: validation.StringInSlice([]string{
				"S",
				"D",
			}, false),
		},
		"checkbox_enable": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "This checkbox is", // need description for this checkbox
		},
		"lb_reserve_ip": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "This field is for any reserve IP which is going to attach on load balancer",
		},
		"ssl_certificate_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "This field is used to set ssl sertificate if lb mode is https or both",
		},
		"ssl_context": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "This field is used to set ssl context",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"redirect_to_https": {
						Type:        schema.TypeBool,
						Optional:    true,
						Default:     false,
						Description: "If Load balancer is set to both http and https then this option need to select",
					},
				},
			},
		},
		"enable_bitninja": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Modular security tool used to enable load balancer from wide range of cyber attacks",
		},
		"backends": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "This will contain the backend details which will be attached to load balancer",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"balance": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "This will contain the type of algorithim used while load balancing",
						ValidateFunc: validation.StringInSlice([]string{
							"source",
							"roundrobin",
							"leastconn",
						}, false),
					},
					"checkbox_enable": {
						Type:        schema.TypeBool,
						Optional:    true,
						Default:     false,
						Description: "This checkbox is to enable healthcheck",
					},
					"domain_name": {
						Type:        schema.TypeString,
						Optional:    true,
						Default:     "localhost",
						Description: "domain name for healthcheck",
					},
					"check_url": {
						Type:        schema.TypeString,
						Optional:    true,
						Default:     "/",
						Description: "endpoint of healthckeck to ping",
					},
					"servers": {
						Type:        schema.TypeList,
						Required:    true,
						Description: "description of servers that are going to attach on backend",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"backend_name": {
									Type:         schema.TypeString,
									Required:     true,
									Description:  "Node name which you want to attach",
									ValidateFunc: node.ValidateName,
								},
								"backend_ip": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "Private IP of the node",
								},
								"backend_port": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "Port Number of the node",
								},
							},
						},
					},
					"http_check": {
						Type:        schema.TypeBool,
						Optional:    true,
						Default:     false,
						Description: "Check if http health check in enable",
					},
				},
			},
		},
		"acl_list": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "This will give the acl rule which you want to apply",
			Elem:        &schema.Schema{Type: schema.TypeMap},
		},
		"acl_map": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "This will give you how you want to route request according to acl rule",
			Elem:        &schema.Schema{Type: schema.TypeMap},
		},
		"vpc_list": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "detail of vpc attach to load balancer",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"network_id": {
						Type:        schema.TypeFloat,
						Required:    true,
						Description: "Network Id of vpc",
					},
					"vpc_name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "name of the VPC",
					},
					"ipv4_cidr": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "VPC ipv4 cidr",
					},
				},
			},
		},
		"enable_eos_logger": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "If you want to get the logs of loadbalancer. Please connect eos bucket",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"appliance_id": {
						Type:        schema.TypeInt,
						Required:    true,
						Default:     0,
						Description: "ID of the appliance",
					},
					"access_key": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Access key of your object storage bucket",
					},
					"secret_key": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Secret key of your object storage bucket",
					},
					"bucket": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Bucket name of your object storage bucket",
					},
				},
			},
		},
		"scaler_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "Need scalar ID if you want to attach autoscaling",
		},
		"scaler_port": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "Need scalar port if you want to attach autoscaling",
		},
		"tcp_backend": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Need Information of TCP backend If user want to attach",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"backend_name": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "Your TCP backend name",
						ValidateFunc: node.ValidateName,
					},
					"port": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Port number for your TCP backend. 8080, 10050, 9101,80 or 443 port not allowed",
						ValidateFunc: validation.StringNotInSlice([]string{
							"8080",
							"10050",
							"9101",
							"80",
							"443",
						}, false),
					},
					"balance": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "This will contain the type of algorithim used while load balancing",
						ValidateFunc: validation.StringInSlice([]string{
							"source",
							"roundrobin",
							"leastconn",
						}, false),
					},
					"servers": {
						Type:        schema.TypeList,
						Required:    true,
						Description: "description of servers that are going to attach on backend",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"backend_name": {
									Type:         schema.TypeString,
									Required:     true,
									Description:  "Node name which you want to attach",
									ValidateFunc: node.ValidateName,
								},
								"backend_ip": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "Private IP of the node",
								},
								"backend_port": {
									Type:        schema.TypeInt,
									Required:    true,
									Description: "Port Number of the node",
								},
							},
						},
					},
				},
			},
		},
		"is_ipv6_attached": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "This is used to attach IPV6 on your load balancer",
		},
		"public_ip": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Public IP of load balancer",
		},
		"private_ip": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Private IP of load balancer",
		},
		"ram": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "This is the ram allotted to your loadbalancer",
		},
		"disk": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "This is the disk storage allotted to your loadbalancer",
		},
		"vcpu": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "This is the vcpu allotted to your loadbalancer",
		},
		"location": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "Delhi",
			Description: "This is the region of your loadbalancer",
		},
		"host_target_ipv6": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Computed:    true,
			Description: "This is the ipv6 allotted to your loadbalancer",
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "This is the status of your loadbalancer",
		},
	}
}

func resourceCreateLoadBalancer(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics

	log.Printf("[INFO] LOAD BALANCER CREATION STARTS")
	backends, err := ExpandBackends(d.Get("backends").([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	loadBalancerObj := models.LoadBalancerCreate{
		PlanName:         d.Get("plan_name").(string),
		LbName:           d.Get("lb_name").(string),
		LbType:           d.Get("lb_type").(string),
		LbMode:           d.Get("lb_mode").(string),
		LbPort:           GetLbPort(d.Get("lb_mode").(string)),
		NodeListType:     d.Get("node_list_type").(string),
		CheckBoxEnable:   d.Get("checkbox_enable").(string),
		LbReserveIp:      d.Get("lb_reserve_ip").(string),
		SslCertificateId: d.Get("ssl_certificate_id").(string),
		EnableBitninja:   d.Get("enable_bitninja").(bool),
		Backends:         backends,
		ScalerId:         d.Get("scaler_id").(string),
		ScalerPort:       d.Get("scaler_port").(string),
		IsIpv6Attached:   d.Get("is_ipv6_attached").(bool),
	}
	log.Println("=========================== LOAD BALANCER OBJECT ============================")
	log.Println(loadBalancerObj)
	enableEosLogger, ok := d.GetOk("enable_eos_logger")
	log.Println("===========================GET EOS DETAIL================================")
	log.Println(enableEosLogger)
	log.Println(ok)
	if ok {
		eosDetail, err := ExpandEnableEosLogger(enableEosLogger.(*schema.Set).List())
		if err != nil {
			return diag.FromErr(err)
		}
		loadBalancerObj.EnableEosLogger = eosDetail
	}
	log.Println("=========================== LOAD BALANCER OBJECT ============================")
	log.Println(loadBalancerObj)
	aclList, ok := d.GetOk("acl_list")
	if ok {
		aclListDetail, err := ExpandAclList(aclList.(*schema.Set).List())
		if err != nil {
			return diag.FromErr(err)
		}
		loadBalancerObj.AclList = aclListDetail
	} else {
		loadBalancerObj.AclList = make([]models.AclListInfo, 0)
	}
	log.Println("=========================== LOAD BALANCER OBJECT ============================")
	log.Println(loadBalancerObj)
	aclMap, ok := d.GetOk("acl_map")
	if ok {
		aclMapDetail, err := ExpandAclMap(aclMap.(*schema.Set).List())
		if err != nil {
			return diag.FromErr(err)
		}
		loadBalancerObj.AclMap = aclMapDetail
	} else {
		loadBalancerObj.AclMap = make([]models.AclMapInfo, 0)
	}
	log.Println("=========================== LOAD BALANCER OBJECT ============================")
	log.Println(loadBalancerObj)
	tcpBackend, ok := d.GetOk("tcp_backend")
	if ok {
		tcpBackendDetail, err := ExpandTcpBackend(tcpBackend.(*schema.Set).List())
		if err != nil {
			return diag.FromErr(err)
		}
		loadBalancerObj.TcpBackend = tcpBackendDetail
	} else {
		loadBalancerObj.TcpBackend = make([]models.TcpBackendDetail, 0)
	}

	vpcList, ok := d.GetOk("vpc_list")
	if ok {
		vpcListDetail, err := ExpandVpcList(vpcList.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		loadBalancerObj.VpcList = vpcListDetail
	} else {
		loadBalancerObj.VpcList = make([]models.VpcDetail, 0)
	}

	sslContext, ok := d.GetOk("ssl_context")
	if ok {
		sslContextList := sslContext.([]interface{})
		detail := sslContextList[0].(map[string]interface{})
		loadBalancerObj.SslContext = detail
	} else {
		loadBalancerObj.SslContext = map[string]interface{}{"redirect_to_https": false}
	}
	log.Println("=========================== LOAD BALANCER OBJECT BEFORE CREATION ============================")
	log.Println(loadBalancerObj)
	response, err := apiClient.NewLoadBalancer(&loadBalancerObj)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] LOAD BALANCER CREATE | RESPONSE BODY | %+v", response)

	if _, codeok := response["code"]; !codeok {
		return diag.Errorf(response["message"].(string))
	}

	data := response["data"].(map[string]interface{})
	if data["is_credit_sufficient"] == false {
		return diag.Errorf("Credit is not sufficient")
	}
	log.Printf("[INFO] load balancer creation | before setting fields")

	lbId := data["id"].(float64)
	lbId = math.Round(lbId)
	d.SetId(strconv.Itoa(int(math.Round(lbId))))
	d.Set("public_ip", data["IP"].(string))
	return diags
}

func resourceReadLoadBalancer(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return make(diag.Diagnostics, 0)
}

func resourceUpdateLoadBalancer(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return make(diag.Diagnostics, 0)
}

func resourceDeleteLoadBalancer(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return make(diag.Diagnostics, 0)
}

func resourceExistsLoadBalancer(d *schema.ResourceData, m interface{}) (bool, error) {
	return true, nil
}
