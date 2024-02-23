package blockstorage

import (
	"context"
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"

	"strings"
	// "time"
	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceBlockStorage() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the block storage, also acts as its unique ID",
				ForceNew:     true,
				ValidateFunc: validateName,
			},
			"size": {
				Type:     schema.TypeFloat,
				Required: true,
				// ForceNew:     true,
				Description:  "Size of the block storage in GB",
				ValidateFunc: validateSize,
			},
			"iops": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "IOPS of the block storage",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the project. It should be unique",
			},
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Location of the block storage",
				ValidateFunc: validation.StringInSlice([]string{
					"Delhi",
					"Mumbai",
				}, false),
				Default: "Delhi",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the node",
			},
		},

		CreateContext: resourceCreateBlockStorage,
		ReadContext:   resourceReadBlockStorage,
		UpdateContext: resourceUpdateBlockStorage,
		DeleteContext: resourceDeleteBlockStorage,
		Exists:        resourceExistsBlockStorage,
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

func resourceCreateBlockStorage(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics

	log.Printf("[INFO] BLOCK STORAGE CREATE STARTS ")
	blockStorage := models.BlockStorageCreate{
		Name: d.Get("name").(string),
		Size: d.Get("size").(float64),
		// IOPS: d.Get("iops").(int), //I think because we are not taking this from the end user as input
	}

	iops := calculateIOPS(blockStorage.Size)
	blockStorage.IOPS = iops

	resBlockStorage, err := apiClient.NewBlockStorage(&blockStorage, d.Get("project_id").(int), d.Get("location").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] BLOCK STORAGE CREATE | RESPONSE BODY | %+v", resBlockStorage)
	if _, codeok := resBlockStorage["code"]; !codeok {
		return diag.Errorf(resBlockStorage["message"].(string))
	}

	data := resBlockStorage["data"].(map[string]interface{})
	if data["is_credit_sufficient"] == false {
		return diag.Errorf(resBlockStorage["message"].(string))
	}
	log.Printf("[INFO] Block Storage creation | before setting fields")
	blockStorageIDFloat, ok := data["id"].(float64)
	if !ok {
		return diag.Errorf("Block ID is not a valid float64 in the response %v", data["id"])
	}

	blockStorageID := int(math.Round(blockStorageIDFloat))
	d.SetId(strconv.Itoa(blockStorageID))
	d.Set("iops", iops)
	return diags
}

func resourceReadBlockStorage(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics

	log.Printf("[INFO] BLOCK STORAGE READ STARTS")
	blockStorageID := d.Id()

	blockStorage, err := apiClient.GetBlockStorage(blockStorageID, d.Get("project_id").(int), d.Get("location").(string))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
			return diags
		}
		return diag.Errorf("error finding Block Storage with ID %s: %s", blockStorageID, err.Error())
	}

	log.Printf("[INFO] BLOCK STORAGE READ | BEFORE SETTING DATA")
	data := blockStorage["data"].(map[string]interface{})
	template := data["template"].(map[string]interface{})
	resSize := convertIntoGB(data["size"].(float64))
	d.Set("name", data["name"].(string))
	d.Set("size", resSize)
	d.Set("status", data["status"].(string))
	d.Set("iops", template["TOTAL_IOPS_SEC"].(string))

	log.Printf("[INFO] BLOCK STORAGE READ | AFTER SETTING DATA")

	return diags
}

func resourceUpdateBlockStorage(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	apiClient := m.(*client.Client)

	blockStorageID := d.Id()

	_, err := apiClient.GetBlockStorage(blockStorageID, d.Get("project_id").(int), d.Get("location").(string))
	if err != nil {
		return diag.Errorf("error finding Block Storage with ID %s", blockStorageID)

	}
	return diag.Errorf("you cannot update parameters after block storage creation. kindly destroy it and then create a new block storage.")

}

func resourceDeleteBlockStorage(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics
	blockStorageID := d.Id()
	node_status := d.Get("status").(string)
	if node_status == "Saving" || node_status == "Creating" {
		return diag.Errorf("Node in %s state", node_status)
	}
	err := apiClient.DeleteBlockStorage(blockStorageID, d.Get("project_id").(int), d.Get("location").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

func resourceExistsBlockStorage(d *schema.ResourceData, m interface{}) (bool, error) {
	apiClient := m.(*client.Client)

	blockStorageID := d.Id()
	_, err := apiClient.GetBlockStorage(blockStorageID, d.Get("project_id").(int), d.Get("location").(string))

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func calculateIOPS(size float64) int {
	switch size {
	case 250:
		return 5000
	case 500:
		return 10000
	case 1000:
		return 20000
	case 2000:
		return 40000
	case 4000:
		return 80000
	case 8000:
		return 120000
	case 16000:
		return 240000
	case 24000:
		return 360000
	default:
		return 0
	}
}

func convertIntoGB(bsSizeRes float64) float64 {
	return bsSizeRes / 1024
}

func validateSize(i interface{}, key string) (ws []string, es []error) {
	bsSize, ok := i.(float64)
	if !ok {
		es = append(es, fmt.Errorf("expected a float64"))
		return
	}

	validSizes := []float64{250, 500, 1000, 2000, 4000, 8000} // These values are the size options available on MyAccount
	valid := false
	for _, size := range validSizes {
		if bsSize == size {
			valid = true
			break
		}
	}

	if !valid {
		es = append(es, fmt.Errorf("size must be one of %v", validSizes))
	}

	return
}
