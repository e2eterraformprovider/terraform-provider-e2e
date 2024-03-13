package ssh_key

import (
	"context"
	"log"

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
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Location(region) where the ssh key is to be created",
			},
			"pk": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique Id for the SSH Key",
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
	res, err := apiClient.AddSshKey(ssh_key, project_id, d.Get("location").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] SSH_KEY CREATE | RESPONSE BODY | %+v", res)
	if _, codeok := res["code"]; !codeok {
		return diag.Errorf(res["message"].(string))
	}

	data := res["data"].(map[string]interface{})

	log.Printf("[INFO] Ssh key creation | before setting fields")

	d.Set("label", data["label"].(string))
	d.Set("ssh_key", data["ssh_key"].(string))
	d.Set("pk", data["pk"].(string))
	d.Set("timestamp", data["timestamp"].(string))
	d.Set("project_id", data["project_id"].(string))
	return diags
}

func resourceReadSshKey(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	apiClient := m.(*client.Client)
	var diags diag.Diagnostics
	log.Printf("[info] inside SSH key Resource read")
	sshKeyId := d.Get("pk").(string)
	project_id := d.Get("project_id").(string)
	location := d.Get("location").(string)
	res, err := apiClient.ReadSshKey(sshKeyId, project_id, location)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[info] node Resource read | before setting data")
	data := res["data"].(map[string]interface{})
	d.Set("label", data["label"].(string))
	d.Set("ssh_key", data["ssh_key"].(string))
	d.Set("pk", data["pk"].(string))
	d.Set("timestamp", data["timestamp"].(string))
	d.Set("project_id", data["project_id"].(string))

	log.Printf("[info] SSH Key Resource read | after setting data")
	return diags

}

func resourceDeleteSshKey(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics
	ssh_key_id := d.Get("pk").(string)
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

// func resourceExistsSshKey(d *schema.ResourceData, m interface{}) (bool, error) {
// 	apiClient := m.(*client.Client)

// 	nodeId := d.Id()
// 	project_id := d.Get("project_id").(string)
// 	_, err := apiClient.GetNode(nodeId, project_id)

// 	if err != nil {
// 		if strings.Contains(err.Error(), "not found") {
// 			return false, nil
// 		} else {
// 			return false, err
// 		}
// 	}
// 	return true, nil
// }
