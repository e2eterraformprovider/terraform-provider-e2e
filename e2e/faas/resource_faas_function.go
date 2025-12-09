package faas

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	e2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/util"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Default values for FaaS function configuration
const (
	defaultMemorySize  = 256
	defaultTimeout     = 30
	defaultMinReplicas = 1
	defaultMaxReplicas = 5
)

// Validation limits for FaaS function configuration
const (
	minMemorySize = 128
	minTimeout    = 1
	maxTimeout    = 900
)

// Async operation timeouts
const (
	functionCreateTimeout = 10 * time.Minute
	functionUpdateTimeout = 10 * time.Minute
	functionDeleteTimeout = 2 * time.Minute
	functionPollInterval  = 10 * time.Second
)

func ResourceFaasFunction() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Common fields - use constants and helpers
			e2econstants.AttrRegion:    config.RegionSchema(),
			e2econstants.AttrLocation:  config.LocationSchema(),
			e2econstants.AttrProjectID: config.ProjectIDSchemaResource(),

			// FaaS-specific fields
			e2econstants.AttrName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "name of the FaaS function",
			},
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "the namespace for the FaaS function",
			},
			"runtime": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "the runtime for the function (e.g., python-3.11-fastapi, node-18, go-1.21)",
			},
			"code_inline": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"code_file"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.TrimSpace(old) == strings.TrimSpace(new)
				},
				Description: "the inline code for the function",
			},
			"code_file": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"code_inline"},
				Description:   "path to local zip file containing function code",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "description of the FaaS function",
			},
			"memory_mb": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      defaultMemorySize,
				ValidateFunc: validation.IntAtLeast(minMemorySize),
				Description:  "the memory allocated to the function in megabytes",
			},
			"timeout_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      defaultTimeout,
				ValidateFunc: validation.IntBetween(minTimeout, maxTimeout),
				Description:  "the execution timeout for the function in seconds",
			},
			"min_replicas": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      defaultMinReplicas,
				ValidateFunc: validation.IntAtLeast(0),
				Description:  "the minimum number of replicas",
			},
			"max_replicas": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      defaultMaxReplicas,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "the maximum number of replicas",
			},
			"environment_variables": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "environment variables for the function",
			},
			e2econstants.AttrTags: {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "tags for the FaaS function",
			},
			"endpoint_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the endpoint URL of the function",
			},
			e2econstants.AttrStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "state of the FaaS function instance",
			},
			e2econstants.AttrCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the creation date for the FaaS function",
			},
			e2econstants.AttrUpdatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the last update date for the FaaS function",
			},
		},

		CreateContext: resourceCreateFaasFunction,
		ReadContext:   resourceReadFaasFunction,
		UpdateContext: resourceUpdateFaasFunction,
		DeleteContext: resourceDeleteFaasFunction,
		Exists:        resourceExistsFaasFunction,

		CustomizeDiff: customizeDiffFaasFunction,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// customizeDiffFaasFunction handles validation
func customizeDiffFaasFunction(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
	// Validate replica range
	minReplicas := diff.Get("min_replicas").(int)
	maxReplicas := diff.Get("max_replicas").(int)

	if minReplicas > maxReplicas {
		return fmt.Errorf("min_replicas (%d) cannot be greater than max_replicas (%d)", minReplicas, maxReplicas)
	}

	// Require at least one code source
	_, hasCodeInline := diff.GetOk("code_inline")
	_, hasCodeFile := diff.GetOk("code_file")

	if !hasCodeInline && !hasCodeFile {
		return fmt.Errorf("one of code_inline or code_file must be specified")
	}

	if hasCodeInline && hasCodeFile {
		return fmt.Errorf("code_inline and code_file are mutually exclusive")
	}

	return nil
}

func resourceCreateFaasFunction(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	client := cfg.Goe2eClient()
	var diags diag.Diagnostics

	namespace := d.Get("namespace").(string)

	log.Printf("[INFO] FAAS FUNCTION CREATE STARTS")

	// First, try to create the namespace (it may already exist, which is fine)
	_, _, err := client.FaaS.CreateNamespace(ctx, namespace)
	if err != nil {
		log.Printf("[WARN] Namespace creation returned error (may already exist): %v", err)
	}

	// Get code from either code_inline or code_file
	var code string
	if codeInline, ok := d.GetOk("code_inline"); ok {
		code = codeInline.(string)
	} else if codeFile, ok := d.GetOk("code_file"); ok {
		// TODO: In Phase 3, implement file upload logic
		// For now, treat code_file as inline code (path to be read)
		code = codeFile.(string)
		log.Printf("[WARN] code_file support for actual file upload is not yet implemented. Treating as inline code.")
	}

	// Convert environment variables
	environment := make(map[string]string)
	if env, ok := d.GetOk("environment_variables"); ok {
		for k, v := range env.(map[string]interface{}) {
			environment[k] = v.(string)
		}
	}

	// Create the function
	createReq := &goe2e.FaasFunctionCreateRequest{
		Name:        d.Get("name").(string),
		Namespace:   namespace,
		Runtime:     d.Get("runtime").(string),
		Code:        code,
		MemoryMB:    d.Get("memory_mb").(int),
		Timeout:     d.Get("timeout_seconds").(int),
		MinReplicas: d.Get("min_replicas").(int),
		MaxReplicas: d.Get("max_replicas").(int),
		Environment: environment,
	}

	fn, _, err := client.FaaS.CreateFunction(ctx, createReq)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] FAAS FUNCTION CREATE | RESPONSE: %+v", fn)

	// Set the ID immediately so it can be used if waiting fails
	d.SetId(fn.ID)

	// Wait for function to be ready
	if err := util.WaitForFunctionReady(ctx, client, fn.ID, functionCreateTimeout); err != nil {
		return diag.FromErr(fmt.Errorf("function created but failed to become ready: %w", err))
	}

	// Refresh function data after waiting
	fn, _, err = client.FaaS.GetFunction(ctx, fn.ID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading function after creation: %w", err))
	}
	if fn == nil {
		return diag.FromErr(fmt.Errorf("function not found after creation"))
	}

	// Set the attributes
	if err := d.Set("name", fn.Name); err != nil {
		return diag.FromErr(fmt.Errorf("error setting name: %w", err))
	}
	if err := d.Set("namespace", fn.Namespace); err != nil {
		return diag.FromErr(fmt.Errorf("error setting namespace: %w", err))
	}
	if err := d.Set("runtime", fn.Runtime); err != nil {
		return diag.FromErr(fmt.Errorf("error setting runtime: %w", err))
	}
	if err := d.Set("memory_mb", fn.MemoryMB); err != nil {
		return diag.FromErr(fmt.Errorf("error setting memory_mb: %w", err))
	}
	if err := d.Set("timeout_seconds", fn.Timeout); err != nil {
		return diag.FromErr(fmt.Errorf("error setting timeout_seconds: %w", err))
	}
	if err := d.Set("min_replicas", fn.MinReplicas); err != nil {
		return diag.FromErr(fmt.Errorf("error setting min_replicas: %w", err))
	}
	if err := d.Set("max_replicas", fn.MaxReplicas); err != nil {
		return diag.FromErr(fmt.Errorf("error setting max_replicas: %w", err))
	}
	if err := d.Set("endpoint_url", fn.EndpointURL); err != nil {
		return diag.FromErr(fmt.Errorf("error setting endpoint_url: %w", err))
	}
	if err := d.Set(e2econstants.AttrStatus, fn.Status); err != nil {
		return diag.FromErr(fmt.Errorf("error setting status: %w", err))
	}
	if err := d.Set("created_at", fn.CreatedAt); err != nil {
		return diag.FromErr(fmt.Errorf("error setting created_at: %w", err))
	}
	if err := d.Set(e2econstants.AttrUpdatedAt, fn.UpdatedAt); err != nil {
		return diag.FromErr(fmt.Errorf("error setting updated_at: %w", err))
	}

	if fn.Environment != nil {
		if err := d.Set("environment_variables", fn.Environment); err != nil {
			return diag.FromErr(fmt.Errorf("error setting environment_variables: %w", err))
		}
	}

	return diags
}

func resourceReadFaasFunction(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	client := cfg.Goe2eClient()
	var diags diag.Diagnostics

	functionID := d.Id()

	log.Printf("[INFO] FAAS FUNCTION READ | ID: %s", functionID)

	fn, _, err := client.FaaS.GetFunction(ctx, functionID)
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

	if err := d.Set("name", fn.Name); err != nil {
		return diag.FromErr(fmt.Errorf("error setting name: %w", err))
	}
	if err := d.Set("namespace", fn.Namespace); err != nil {
		return diag.FromErr(fmt.Errorf("error setting namespace: %w", err))
	}
	if err := d.Set("runtime", fn.Runtime); err != nil {
		return diag.FromErr(fmt.Errorf("error setting runtime: %w", err))
	}
	if err := d.Set("memory_mb", fn.MemoryMB); err != nil {
		return diag.FromErr(fmt.Errorf("error setting memory_mb: %w", err))
	}
	if err := d.Set("timeout_seconds", fn.Timeout); err != nil {
		return diag.FromErr(fmt.Errorf("error setting timeout_seconds: %w", err))
	}
	if err := d.Set("min_replicas", fn.MinReplicas); err != nil {
		return diag.FromErr(fmt.Errorf("error setting min_replicas: %w", err))
	}
	if err := d.Set("max_replicas", fn.MaxReplicas); err != nil {
		return diag.FromErr(fmt.Errorf("error setting max_replicas: %w", err))
	}
	if err := d.Set("endpoint_url", fn.EndpointURL); err != nil {
		return diag.FromErr(fmt.Errorf("error setting endpoint_url: %w", err))
	}
	if err := d.Set(e2econstants.AttrStatus, fn.Status); err != nil {
		return diag.FromErr(fmt.Errorf("error setting status: %w", err))
	}
	if err := d.Set("created_at", fn.CreatedAt); err != nil {
		return diag.FromErr(fmt.Errorf("error setting created_at: %w", err))
	}
	if err := d.Set(e2econstants.AttrUpdatedAt, fn.UpdatedAt); err != nil {
		return diag.FromErr(fmt.Errorf("error setting updated_at: %w", err))
	}

	if fn.Environment != nil {
		if err := d.Set("environment_variables", fn.Environment); err != nil {
			return diag.FromErr(fmt.Errorf("error setting environment_variables: %w", err))
		}
	}

	return diags
}

func resourceUpdateFaasFunction(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	client := cfg.Goe2eClient()
	var diags diag.Diagnostics

	functionID := d.Id()

	log.Printf("[INFO] FAAS FUNCTION UPDATE | ID: %s", functionID)

	updateReq := &goe2e.FaasFunctionUpdateRequest{}
	hasChanges := false

	if d.HasChange("code_inline") {
		updateReq.Code = goe2e.String(d.Get("code_inline").(string))
		hasChanges = true
	}

	if d.HasChange("code_file") {
		// TODO: In Phase 3, implement file upload logic
		updateReq.Code = goe2e.String(d.Get("code_file").(string))
		hasChanges = true
		log.Printf("[WARN] code_file support for actual file upload is not yet implemented. Treating as inline code.")
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

	if d.HasChange("environment_variables") {
		environment := make(map[string]string)
		if env, ok := d.GetOk("environment_variables"); ok {
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

	fn, _, err := client.FaaS.UpdateFunction(ctx, functionID, updateReq)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] FAAS FUNCTION UPDATE | RESPONSE: %+v", fn)

	// Wait for function to be ready after update
	if err := util.WaitForFunctionReady(ctx, client, functionID, functionUpdateTimeout); err != nil {
		return diag.FromErr(fmt.Errorf("function updated but failed to become ready: %w", err))
	}

	// Refresh function data after waiting
	fn, _, err = client.FaaS.GetFunction(ctx, functionID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading function after update: %w", err))
	}
	if fn == nil {
		return diag.FromErr(fmt.Errorf("function not found after update"))
	}

	// Update the state with the response
	if err := d.Set("memory_mb", fn.MemoryMB); err != nil {
		return diag.FromErr(fmt.Errorf("error setting memory_mb: %w", err))
	}
	if err := d.Set("timeout_seconds", fn.Timeout); err != nil {
		return diag.FromErr(fmt.Errorf("error setting timeout_seconds: %w", err))
	}
	if err := d.Set("min_replicas", fn.MinReplicas); err != nil {
		return diag.FromErr(fmt.Errorf("error setting min_replicas: %w", err))
	}
	if err := d.Set("max_replicas", fn.MaxReplicas); err != nil {
		return diag.FromErr(fmt.Errorf("error setting max_replicas: %w", err))
	}
	if err := d.Set("endpoint_url", fn.EndpointURL); err != nil {
		return diag.FromErr(fmt.Errorf("error setting endpoint_url: %w", err))
	}
	if err := d.Set(e2econstants.AttrStatus, fn.Status); err != nil {
		return diag.FromErr(fmt.Errorf("error setting status: %w", err))
	}
	if err := d.Set(e2econstants.AttrUpdatedAt, fn.UpdatedAt); err != nil {
		return diag.FromErr(fmt.Errorf("error setting updated_at: %w", err))
	}

	if fn.Environment != nil {
		if err := d.Set("environment_variables", fn.Environment); err != nil {
			return diag.FromErr(fmt.Errorf("error setting environment_variables: %w", err))
		}
	}

	return diags
}

func resourceDeleteFaasFunction(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	client := cfg.Goe2eClient()
	var diags diag.Diagnostics

	functionID := d.Id()

	log.Printf("[INFO] FAAS FUNCTION DELETE | ID: %s", functionID)

	_, err := client.FaaS.DeleteFunction(ctx, functionID)
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

	fn, _, err := client.FaaS.GetFunction(context.Background(), functionID)
	if err != nil {
		return false, err
	}

	return fn != nil, nil
}
