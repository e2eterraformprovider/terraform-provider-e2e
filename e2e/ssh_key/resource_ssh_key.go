package ssh_key

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/e2eterraformprovider/terraform-provider-e2e/models"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceSshKey() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{

			"label": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The label(name) of the ssh key",
				ForceNew:    true,
			},

			"ssh_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "your ssh key",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the project associated with the ssh key",
			},
			"timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp of the SSH Key",
			},
		},

		CreateContext: resourceCreateSshKey,
		ReadContext:   resourceReadSshKey,
		UpdateContext: resourceUpdateSshKey,
		DeleteContext: resourceDeleteSshKey,
		Exists:        resourceExistsSshKey,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateSshKey(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics

	log.Printf("[INFO] SSH KEY ADD STARTS ")
	ssh_key := models.AddSshKey{
		Label:  d.Get("label").(string),
		SshKey: d.Get("ssh_key").(string),
	}

	project_id := d.Get("project_id").(string)
	res, err := apiClient.AddSshKey(ssh_key, project_id)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] SSH_KEY CREATE | RESPONSE BODY | %+v", res)
	if code, codeok := res["code"].(float64); !codeok || int(code) < 200 || int(code) >= 300 {
		// log.Printf("res['code'] = %v, res['code'].(int) = %T", res["code"], res["code"].(int))
		log.Printf("code = %v, type = %T", code, code)
		return diag.Errorf("%v | %v", res["message"].(string), res["errors"])
	}

	log.Printf("[INFO] Ssh key creation | res = %+v, type = %T", res, res)
	data := res["data"].(map[string]interface{})

	ssh_key_id := strconv.FormatFloat(data["pk"].(float64), 'f', 0, 64)
	log.Printf("[INFO] SSH KEY CREATE | SSH KEY ID | %+v, Type = %T", ssh_key_id, ssh_key_id)
	log.Printf("[INFO] SSH KEY CREATE | label | %+v, Type = %T", data["label"], data["label"])
	log.Printf("[INFO] SSH KEY CREATE | ssh_key | %+v, Type = %T", data["ssh_key"], data["ssh_key"])
	log.Printf("[INFO] SSH KEY CREATE | timestamp | %+v, Type = %T", data["timestamp"], data["timestamp"])
	d.SetId(ssh_key_id)

	d.Set("label", data["label"].(string))
	d.Set("ssh_key", data["ssh_key"].(string))
	d.Set("timestamp", data["timestamp"].(string))
	return diags
}

func resourceReadSshKey(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	apiClient := m.(*client.Client)
	var diags diag.Diagnostics
	log.Printf("[info] inside SSH key Resource read")
	label := d.Get("label").(string)
	project_id := d.Get("project_id").(string)
	res, err := apiClient.GetSshKey(label, project_id)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[info] SSH Key Resource read | res = %+v, type = %T", res, res)
	data := res["data"].(map[string]interface{})

	ssh_key_id := strconv.FormatFloat(data["pk"].(float64), 'f', 0, 64)
	d.SetId(ssh_key_id)

	d.Set("label", data["label"].(string))
	d.Set("ssh_key", data["ssh_key"].(string))
	d.Set("timestamp", data["timestamp"].(string))

	log.Printf("[info] SSH Key Resource read | after setting data")
	return diags

}

func resourceDeleteSshKey(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics
	ssh_key_id := d.Id()
	project_id := d.Get("project_id").(string)
	location := d.Get("location").(string)

	err := apiClient.DeleteSshKey(ssh_key_id, project_id, location)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
func resourceUpdateSshKey(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// apiClient := m.(*client.Client)
	var diags diag.Diagnostics

	return diags

}

func resourceExistsSshKey(d *schema.ResourceData, m interface{}) (bool, error) {
	apiClient := m.(*client.Client)

	ssh_key_id := d.Id()
	project_id := d.Get("project_id").(string)
	_, err := apiClient.GetSshKey(ssh_key_id, project_id)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}
