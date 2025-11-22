package faas

import (
	"context"
	"log"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceFaasFunction() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the FaaS function",
			},
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The namespace for the FaaS function",
			},
			"runtime": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The runtime for the function (e.g., python-3.11-fastapi, node-18, go-1.21)",
			},
			"code_inline": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The inline code for the function",
			},
			"memory_mb": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     256,
				Description: "Memory allocation in MB",
			},
			"timeout_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     30,
				Description: "Function timeout in seconds",
			},
			"min_replicas": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Minimum number of replicas",
			},
			"max_replicas": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5,
				Description: "Maximum number of replicas",
			},
			"environment": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Environment variables for the function",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the project",
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The location/region for the function",
			},
			"endpoint_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The endpoint URL of the function",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the function",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp when the function was created",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp when the function was last updated",
			},
		},

		CreateContext: resourceCreateFaasFunction,
		ReadContext:   resourceReadFaasFunction,
		UpdateContext: resourceUpdateFaasFunction,
		DeleteContext: resourceDeleteFaasFunction,
		Exists:        resourceExistsFaasFunction,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCreateFaasFunction(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	client := cfg.Goe2eClient()
	var diags diag.Diagnostics

	opts := &goe2e.RequestOptions{
		ProjectID: d.Get("project_id").(string),
		Location:  d.Get("location").(string),
	}

	namespace := d.Get("namespace").(string)

	log.Printf("[INFO] FAAS FUNCTION CREATE STARTS")

	// First, try to create the namespace (it may already exist, which is fine)
	_, _, err := client.FaaS.CreateNamespace(ctx, namespace, opts)
	if err != nil {
		log.Printf("[WARN] Namespace creation returned error (may already exist): %v", err)
	}

	// Convert environment variables
	environment := make(map[string]string)
	if env, ok := d.GetOk("environment"); ok {
		for k, v := range env.(map[string]interface{}) {
			environment[k] = v.(string)
		}
	}

	// Create the function
	createReq := &goe2e.FaasFunctionCreateRequest{
		Name:        d.Get("name").(string),
		Namespace:   namespace,
		Runtime:     d.Get("runtime").(string),
		Code:        d.Get("code_inline").(string),
		MemoryMB:    d.Get("memory_mb").(int),
		Timeout:     d.Get("timeout_seconds").(int),
		MinReplicas: d.Get("min_replicas").(int),
		MaxReplicas: d.Get("max_replicas").(int),
		Environment: environment,
	}

	fn, _, err := client.FaaS.CreateFunction(ctx, createReq, opts)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] FAAS FUNCTION CREATE | RESPONSE: %+v", fn)

	// Set the ID and other attributes
	d.SetId(fn.ID)
	d.Set("name", fn.Name)
	d.Set("namespace", fn.Namespace)
	d.Set("runtime", fn.Runtime)
	d.Set("memory_mb", fn.MemoryMB)
	d.Set("timeout_seconds", fn.Timeout)
	d.Set("min_replicas", fn.MinReplicas)
	d.Set("max_replicas", fn.MaxReplicas)
	d.Set("endpoint_url", fn.EndpointURL)
	d.Set("status", fn.Status)
	d.Set("created_at", fn.CreatedAt)
	d.Set("updated_at", fn.UpdatedAt)

	if fn.Environment != nil {
		d.Set("environment", fn.Environment)
	}

	return diags
}

func resourceReadFaasFunction(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	client := cfg.Goe2eClient()
	var diags diag.Diagnostics

	functionID := d.Id()
	opts := &goe2e.RequestOptions{
		ProjectID: d.Get("project_id").(string),
		Location:  d.Get("location").(string),
	}

	log.Printf("[INFO] FAAS FUNCTION READ | ID: %s", functionID)

	fn, _, err := client.FaaS.GetFunction(ctx, functionID, opts)
	if err != nil {
		return diag.FromErr(err)
	}

	if fn == nil {
		log.Printf("[WARN] FaaS function with ID %s not found", functionID)
		d.SetId("")

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "FaaS function not found",
			Detail:   "The FaaS function may have been deleted manually.",
		})

		return diags
	}

	d.Set("name", fn.Name)
	d.Set("namespace", fn.Namespace)
	d.Set("runtime", fn.Runtime)
	d.Set("memory_mb", fn.MemoryMB)
	d.Set("timeout_seconds", fn.Timeout)
	d.Set("min_replicas", fn.MinReplicas)
	d.Set("max_replicas", fn.MaxReplicas)
	d.Set("endpoint_url", fn.EndpointURL)
	d.Set("status", fn.Status)
	d.Set("created_at", fn.CreatedAt)
	d.Set("updated_at", fn.UpdatedAt)

	if fn.Environment != nil {
		d.Set("environment", fn.Environment)
	}

	return diags
}

func resourceUpdateFaasFunction(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	client := cfg.Goe2eClient()
	var diags diag.Diagnostics

	functionID := d.Id()
	opts := &goe2e.RequestOptions{
		ProjectID: d.Get("project_id").(string),
		Location:  d.Get("location").(string),
	}

	log.Printf("[INFO] FAAS FUNCTION UPDATE | ID: %s", functionID)

	updateReq := &goe2e.FaasFunctionUpdateRequest{}
	hasChanges := false

	if d.HasChange("code_inline") {
		updateReq.Code = goe2e.String(d.Get("code_inline").(string))
		hasChanges = true
	}

	if d.HasChange("memory_mb") {
		updateReq.MemoryMB = goe2e.Int(d.Get("memory_mb").(int))
		hasChanges = true
	}

	if d.HasChange("timeout_seconds") {
		updateReq.Timeout = goe2e.Int(d.Get("timeout_seconds").(int))
		hasChanges = true
	}

	if d.HasChange("min_replicas") {
		updateReq.MinReplicas = goe2e.Int(d.Get("min_replicas").(int))
		hasChanges = true
	}

	if d.HasChange("max_replicas") {
		updateReq.MaxReplicas = goe2e.Int(d.Get("max_replicas").(int))
		hasChanges = true
	}

	if d.HasChange("environment") {
		environment := make(map[string]string)
		if env, ok := d.GetOk("environment"); ok {
			for k, v := range env.(map[string]interface{}) {
				environment[k] = v.(string)
			}
		}
		updateReq.Environment = environment
		hasChanges = true
	}

	if !hasChanges {
		log.Printf("[INFO] No changes detected for FaaS function %s", functionID)
		return diags
	}

	fn, _, err := client.FaaS.UpdateFunction(ctx, functionID, updateReq, opts)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] FAAS FUNCTION UPDATE | RESPONSE: %+v", fn)

	// Update the state with the response
	d.Set("memory_mb", fn.MemoryMB)
	d.Set("timeout_seconds", fn.Timeout)
	d.Set("min_replicas", fn.MinReplicas)
	d.Set("max_replicas", fn.MaxReplicas)
	d.Set("endpoint_url", fn.EndpointURL)
	d.Set("status", fn.Status)
	d.Set("updated_at", fn.UpdatedAt)

	if fn.Environment != nil {
		d.Set("environment", fn.Environment)
	}

	return diags
}

func resourceDeleteFaasFunction(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	client := cfg.Goe2eClient()
	var diags diag.Diagnostics

	functionID := d.Id()
	opts := &goe2e.RequestOptions{
		ProjectID: d.Get("project_id").(string),
		Location:  d.Get("location").(string),
	}

	log.Printf("[INFO] FAAS FUNCTION DELETE | ID: %s", functionID)

	_, err := client.FaaS.DeleteFunction(ctx, functionID, opts)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diags
}

func resourceExistsFaasFunction(d *schema.ResourceData, m interface{}) (bool, error) {
	cfg := m.(*config.Config)
	client := cfg.Goe2eClient()

	functionID := d.Id()
	opts := &goe2e.RequestOptions{
		ProjectID: d.Get("project_id").(string),
		Location:  d.Get("location").(string),
	}

	fn, _, err := client.FaaS.GetFunction(context.Background(), functionID, opts)
	if err != nil {
		return false, err
	}

	return fn != nil, nil
}
