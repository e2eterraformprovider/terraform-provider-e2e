package sfs

import (
	// "context"

	"context"
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"

	//"time"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/e2eterraformprovider/terraform-provider-e2e/models"

	// "github.com/hashicorp/terraform-plugin-log"
	// "github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceSfs() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
         
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the resource, also acts as it's unique ID",
				ForceNew:     true,
				ValidateFunc: validateName,
			},
			"plan": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of the Plan",
			},
			"vpc_id":{
				Type:        schema.TypeInt,
				Required:    true,
				Description: "virtual private cloud id of sfs",
			},
			"disk_size":{
				Type:        schema.TypeString,
				Required:    true,
				Description: "size of disk to be created",
			},
			"vpc_name":{
				Type:        schema.TypeString,
				Required:    true
				Description:  "name of virtual private cloud",
			}
			
		
	}
}
}