package e2e

import (
	// "context"
	"fmt"
	"math"
	"regexp"

	"strconv"
	"strings"

	"github.com/devteametwoe/terraform-provider-e2e/client"
	"github.com/devteametwoe/terraform-provider-e2e/models"

	// "github.com/hashicorp/terraform-plugin-log"
	// "github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNode() *schema.Resource {
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
			"backups": {
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
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
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
		},
		Create: resourceCreateNode,
		Read:   resourceReadNode,
		Update: resourceUpdateNode,
		Delete: resourceDeleteNode,
		// Exists: resourceExistsNode,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func validateName(v interface{}, k string) (ws []string, es []error) {
	fmt.Println("hii")
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

func resourceCreateNode(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	tfssh_Keys := d.Get("ssh_keys").(*schema.Set).List()
	ssh_keys := make([]string, len(tfssh_Keys))
	for i, tfssh_key := range tfssh_Keys {
		ssh_keys[i] = tfssh_key.(string)
	}
	node := models.Node{
		Name:                    d.Get("name").(string),
		Label:                   d.Get("label").(string),
		Plan:                    d.Get("plan").(string),
		Backups:                 d.Get("backups").(bool),
		Image:                   d.Get("image").(string),
		Default_public_ip:       d.Get("default_public_ip").(bool),
		Disable_password:        d.Get("disable_password").(bool),
		Enable_bitninja:         d.Get("enable_bitninja").(bool),
		Is_ipv6_availed:         d.Get("is_ipv6_availed").(bool),
		Is_saved_image:          d.Get("is_saved_image").(bool),
		Region:                  d.Get("region").(string),
		Reserve_ip:              d.Get("reserve_ip").(string),
		Vpc_id:                  d.Get("vpc_id").(string),
		Ngc_container_id:        d.Get("ngc_container_id").(int),
		Saved_image_template_id: d.Get("saved_image_template_id").(int),
		Security_group_id:       d.Get("security_group_id").(int),
		SSH_keys:                ssh_keys,
	}

	resnode, err := apiClient.NewNode(&node)
	data := resnode["data"].(map[string]interface{})
	nodeId := data["id"].(float64)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(int(math.Round(nodeId))))
	d.Set("is_active", data["is_active"].(bool))
	d.Set("created_at", data["created_at"].(string))
	d.Set("memory", data["memory"].(string))
	d.Set("status", data["status"].(string))
	d.Set("disk", data["disk"].(string))
	d.Set("price", data["price"].(string))
	return nil
}

func resourceReadNode(d *schema.ResourceData, m interface{}) error {

	apiClient := m.(*client.Client)

	nodeId := d.Id()
	node, err := apiClient.GetNode(nodeId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
		} else {
			return fmt.Errorf("error finding Item with ID %s", nodeId)
		}
	}
	data := node["data"].(map[string]interface{})
	//d.SetId(nodeId)

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
	//d.Set("default_public_ip", data["Default_public_ip"].(string))
	//d.Set("disable_password", data["Disable_password"].(bool))
	//d.Set("enable_bitninja", data["Enable_bitninja"].(bool))
	//d.Set("is_ipv6_availed", data["Is_ipv6_availed"].(string))
	//d.Set("is_saved_image", data["Is_saved_image"].(bool))
	//d.Set("region", data["region"].(string)))
	//d.Set("reserved_ip", data["Reserve_ip"].(string))
	//d.Set("vpc_id", data["Vpc_id"].(string))
	//d.Set("ngc_container_id", node.Ngc_container_id)
	//d.Set("saved_image_template_id", node.Saved_image_template_id)
	//d.Set("security_group_id", data["security_group_id"].(int))
	//d.Set("ssh_keys", node.SSH_keys)

	return nil

}

func resourceUpdateNode(d *schema.ResourceData, m interface{}) error {
	// apiClient := m.(*client.Client)

	// tfssh_Keys := d.Get("ssh_keys").(*schema.Set).List()
	// ssh_keys := make([]string, len(tfssh_Keys))
	// for i, tfssh_key := range tfssh_Keys {
	// 	ssh_keys[i] = tfssh_key.(string)
	// }

	// node := server.Node{
	// 	Name:                    d.Get("name").(string),
	// 	Label:                   d.Get("label").(string),
	// 	Plan:                    d.Get("plan").(string),
	// 	Backups:                 d.Get("backups").(bool),
	// 	Image:                   d.Get("image").(string),
	// 	Default_public_ip:       d.Get("default_public_ip").(bool),
	// 	Disable_password:        d.Get("disable_password").(bool),
	// 	Enable_bitninja:         d.Get("enable_bitninja").(bool),
	// 	Is_ipv6_availed:         d.Get("is_ipv6_availed").(bool),
	// 	Is_saved_image:          d.Get("is_saved_image").(bool),
	// 	Region:                  d.Get("region").(string),
	// 	Reserve_ip:              d.Get("reserve_ip").(string),
	// 	Vpc_id:                  d.Get("vpc_id").(string),
	// 	Ngc_container_id:        d.Get("ngc_container_id").(int),
	// 	Saved_image_template_id: d.Get("saved_image_template_id").(int),
	// 	Security_group_id:       d.Get("security_group_id").(int),
	// 	SSH_keys:                ssh_keys,
	// }
	// fmt.Print("update")
	// err := apiClient.UpdateNode(&node)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func resourceDeleteNode(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	nodeId := d.Id()

	err := apiClient.DeleteNode(nodeId)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
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
