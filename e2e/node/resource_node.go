package node

import (
	// "context"

	"fmt"
	"log"
	"math"
	"regexp"

	"context"
	"strconv"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/e2eterraformprovider/terraform-provider-e2e/models"

	// "github.com/hashicorp/terraform-plugin-log"
	// "github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	BaseUrl = "https://api.e2enetworks.com/myaccount/api/v1/nodes"

//BaseUrl="https://api.e2enetworks."
)

func ResourceNode() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the resource, also acts as it's unique ID",
				ForceNew:     true,
				ValidateFunc: validateName,
			},
			"label": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the group",
			},
			"plan": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Plan",
			},
			"backup": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Tells you the state of your backups",
				Default:     false,
			},
			"image": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the image you have selected",
				Default:     "CentOS-7.5-Distro",
			},
			"default_public_ip": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Tells us the state of default public ip",
				Default:     false,
			},
			"disable_password": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "",
				Default:     false,
			},
			"enable_bitninja": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "",
				Default:     false,
			},
			"is_ipv6_availed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "",
				Default:     false,
			},
			"is_saved_image": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "",
				Default:     false,
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
				Default:     "ncr",
			},
			"reserve_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
				Default:     "",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
				Default:     "",
			},
			"ngc_container_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "",
				Default:     nil,
			},
			"saved_image_template_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "",
				Default:     nil,
			},
			"security_group_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "",
				Default:     150,
			},
			"ssh_keys": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			// "ssh_keys": {
			// 	Type:     schema.TypeList,
			// 	Optional: true,
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"label": {
			// 				Type:        schema.TypeString,
			// 				Computed:    true,
			// 				Description: "label of you ssh key",
			// 			},
			// 			"ssh_key": {
			// 				Type:        schema.TypeString,
			// 				Computed:    true,
			// 				Description: "your ssh key",
			// 			},
			// 		},
			// 	},
			// 	Description: "you list of ssh keys (can add more than one)",
			// 	Default:     nil,
			// },
			"is_active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"memory": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"price": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_monitored": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_bitninja_license_active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"power_status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "power_on",
			},
			"lock_node": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  "false",
			},
			"reboot_node": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  "false",
			},
			"reinstall_node": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  "false",
			},
		},

		CreateContext: resourceCreateNode,
		ReadContext:   resourceReadNode,
		UpdateContext: resourceUpdateNode,
		DeleteContext: resourceDeleteNode,
		Exists:        resourceExistsNode,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func validateName(v interface{}, k string) (ws []string, es []error) {

	var errs []error
	var warns []string
	value, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected name to be string"))
		return warns, errs
	}
	whiteSpace := regexp.MustCompile(`\s+`)
	if whiteSpace.Match([]byte(value)) {
		errs = append(errs, fmt.Errorf("name cannot contain whitespace. Got %s", value))
		return warns, errs
	}
	return warns, errs
}

func resourceCreateNode(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics

	log.Printf("[INFO] inside create ")
	node := models.Node{
		Name:              d.Get("name").(string),
		Label:             d.Get("label").(string),
		Plan:              d.Get("plan").(string),
		Backup:            d.Get("backup").(bool),
		Image:             d.Get("image").(string),
		Default_public_ip: d.Get("default_public_ip").(bool),
		Disable_password:  d.Get("disable_password").(bool),
		Enable_bitninja:   d.Get("enable_bitninja").(bool),
		Is_ipv6_availed:   d.Get("is_ipv6_availed").(bool),
		Is_saved_image:    d.Get("is_saved_image").(bool),
		Region:            d.Get("region").(string),
		Reserve_ip:        d.Get("reserve_ip").(string),
		Vpc_id:            d.Get("vpc_id").(string),
		Security_group_id: d.Get("security_group_id").(int),
		SSH_keys:          d.Get("ssh_keys").([]interface{}),
	}

	resnode, err := apiClient.NewNode(&node)
	if err != nil {
		return diag.FromErr(err)
	}

	if resnode["code"].(float64) != 200 {
		error := resnode["errors"].(string)
		return diag.Errorf(error)

	}

	data := resnode["data"].(map[string]interface{})
	nodeId := data["id"].(float64)
	nodeId = math.Round(nodeId)
	fmt.Println(data)
	d.SetId(strconv.Itoa(int(math.Round(nodeId))))
	d.Set("is_active", data["is_active"].(bool))
	d.Set("created_at", data["created_at"].(string))
	d.Set("memory", data["memory"].(string))
	d.Set("status", data["status"].(string))
	d.Set("disk", data["disk"].(string))
	d.Set("price", data["price"].(string))

	return diags
}

func resourceReadNode(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	apiClient := m.(*client.Client)
	var diags diag.Diagnostics
	log.Printf("[info] inside read")
	nodeId := d.Id()

	node, err := apiClient.GetNode(nodeId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
		} else {
			return diag.Errorf("error finding Item with ID %s", nodeId)

		}
	}
	data := node["data"].(map[string]interface{})
	// if d.Get("status").(string)=="Running"{
	// 	if d.Get("lock_node").(string)=="true" && data["is_locked"].(bool)==false{

	// 	}
	// }
	d.Set("name", data["name"].(string))
	d.Set("label", data["label"].(string))
	d.Set("plan", data["plan"].(string))
	d.Set("backup", data["backup"].(bool))
	d.Set("is_active", data["is_active"].(bool))
	d.Set("created_at", data["created_at"].(string))
	d.Set("memory", data["memory"].(string))
	d.Set("status", data["status"].(string))
	d.Set("disk", data["disk"].(string))
	d.Set("price", data["price"].(string))
	d.Set("lock_node", data["is_locked"].(bool))
	d.Set("public_ip_address", data["public_ip_address"].(string))
	d.Set("private_ip_address", data["private_ip_address"].(string))
	d.Set("is_monitored", data["is_monitored"].(bool))
	d.Set("is_bitninja_license_active", data["is_bitninja_license_active"].(bool))

	if d.Get("status").(string) == "Running" {
		d.Set("power_status", "power_on")
	}
	if d.Get("status").(string) == "Powered off" {
		d.Set("power_status", "power_off")
	}

	return diags

}

func resourceUpdateNode(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	apiClient := m.(*client.Client)

	nodeId := d.Id()

	_, err := apiClient.GetNode(nodeId)
	if err != nil {

		return diag.Errorf("error finding Item with ID %s", nodeId)

	}

	if d.HasChange("power_status") {
		nodestatus := d.Get("status").(string)

		if nodestatus == "Creating" || nodestatus == "Reinstalling" {
			return diag.Errorf("Node is in %s state", d.Get("status").(string))
		}
		if d.Get("lock_node").(bool) == true {
			return diag.Errorf("cannot change the power status as the node is locked")
		}

		log.Printf("[INFO] %s ", d.Get("power_status").(string))
		//fmt.Printf(d.Get("power_status").(string))
		apiClient.UpdateNode(nodeId, d.Get("power_status").(string))

	}

	if d.HasChange("lock_node") {

		if d.Get("status").(string) == "Creating" || d.Get("status").(string) == "Reinstalling" {
			return diag.Errorf("Cannot update as the node is in %s state", d.Get("status").(string))
		}

		if d.Get("lock_node").(bool) == true {
			err := apiClient.UpdateNode(nodeId, "lock_vm")
			if err != nil {
				return diag.FromErr(err)
			}
		}
		if d.Get("lock_node").(bool) == false {
			err := apiClient.UpdateNode(nodeId, "unlock_vm")
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("reboot_node") {

		if d.Get("status").(string) == "Creating" || d.Get("status").(string) == "Reinstalling" {
			return diag.Errorf("Cannot update as the node is in %s state", d.Get("status").(string))
		}

		if d.Get("reboot_node").(bool) == true {
			if d.Get("status").(string) == "Powered off" {
				return diag.Errorf("cannot reboot as the node is powered off")
			}

			err := apiClient.UpdateNode(nodeId, "reboot")
			d.Set("reboot_node", false)
			if err != nil {
				return diag.FromErr(err)
			}
		}

	}
	if d.HasChange("reinstall_node") {

		if d.Get("status").(string) == "Creating" {
			return diag.Errorf("Node is in creating state")
		}
		if d.Get("status").(string) == "Reinstalling" {
			return diag.Errorf("Node already in Reinstalling state")
		}

		if d.Get("reinstall_node").(bool) == true {
			if d.Get("status").(string) == "Powered off" {
				d.Set("reinstall_node", false)
				return diag.Errorf("cannot reinstall as the node is powered off")
			}
			if d.Get("status").(string) == "Reinstalling" {
				d.Set("reinstall_node", false)
				return diag.Errorf("Node already in Reinstalling state")
			}
			err := apiClient.UpdateNode(nodeId, "reinstall")

			d.Set("reinstall_node", false)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	return resourceReadNode(ctx, d, m)

}

func resourceDeleteNode(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics
	nodeId := d.Id()

	err := apiClient.DeleteNode(nodeId)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

func resourceExistsNode(d *schema.ResourceData, m interface{}) (bool, error) {
	apiClient := m.(*client.Client)

	nodeId := d.Id()
	_, err := apiClient.GetNode(nodeId)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}
