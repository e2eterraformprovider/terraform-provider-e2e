package faas

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/util"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceFaasFunction() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Common fields - use constants and helpers
			tfconstants.AttrRegion:    config.RegionSchema(),
			tfconstants.AttrLocation:  config.LocationSchema(),
			tfconstants.AttrProjectID: config.ProjectIDSchemaResource(),

			// FaaS-specific fields
			tfconstants.AttrName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "name of the FaaS function",
			},
			tfconstants.AttrNamespace: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "the namespace for the FaaS function",
			},
			tfconstants.AttrRuntime: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "the runtime for the function (e.g., python-3.11-fastapi, node-18, go-1.21)",
			},
			tfconstants.AttrCodeInline: {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{tfconstants.AttrCodeFile},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.TrimSpace(old) == strings.TrimSpace(new)
				},
				Description: "the inline code for the function",
			},
			tfconstants.AttrCodeFile: {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{tfconstants.AttrCodeInline},
				Description:   "path to local zip file containing function code",
			},
			tfconstants.AttrDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "description of the FaaS function",
			},
			tfconstants.AttrMemoryMB: {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      goe2econstants.FaaSDefaultMemoryMB,
				ValidateFunc: validation.IntAtLeast(tfconstants.FaaSMinMemoryMB),
				Description:  "the memory allocated to the function in megabytes",
			},
			tfconstants.AttrTimeoutSeconds: {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      goe2econstants.FaaSDefaultTimeoutSeconds,
				ValidateFunc: validation.IntBetween(tfconstants.FaaSMinTimeoutSeconds, tfconstants.FaaSMaxTimeoutSeconds),
				Description:  "the execution timeout for the function in seconds",
			},
			tfconstants.AttrMinReplicas: {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      goe2econstants.FaaSDefaultMinReplicas,
				ValidateFunc: validation.IntAtLeast(0),
				Description:  "the minimum number of replicas",
			},
			tfconstants.AttrMaxReplicas: {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      goe2econstants.FaaSDefaultMaxReplicas,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "the maximum number of replicas",
			},
			tfconstants.AttrEnvironmentVariables: {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "environment variables for the function",
			},
			tfconstants.AttrTags: {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "tags for the FaaS function",
			},
			tfconstants.AttrEndpointURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the endpoint URL of the function",
			},
			tfconstants.AttrStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "state of the FaaS function instance",
			},
			tfconstants.AttrCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the creation date for the FaaS function",
			},
			tfconstants.AttrUpdatedAt: {
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
	minReplicas := diff.Get(tfconstants.AttrMinReplicas).(int)
	maxReplicas := diff.Get(tfconstants.AttrMaxReplicas).(int)

	if minReplicas > maxReplicas {
		return fmt.Errorf(ErrMinReplicasGreaterThanMaxFmt, minReplicas, maxReplicas)
	}

	// Require at least one code source
	_, hasCodeInline := diff.GetOk(tfconstants.AttrCodeInline)
	_, hasCodeFile := diff.GetOk(tfconstants.AttrCodeFile)

	if !hasCodeInline && !hasCodeFile {
		return errors.New(ErrCodeInlineOrFileRequired)
	}

	if hasCodeInline && hasCodeFile {
		return errors.New(ErrCodeInlineAndFileExclusive)
	}

	return nil
}

func resourceCreateFaasFunction(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	client := cfg.Goe2eClient()
	var diags diag.Diagnostics

	namespace := d.Get(tfconstants.AttrNamespace).(string)

	log.Printf("[INFO] FAAS FUNCTION CREATE STARTS")

	// First, try to create the namespace (it may already exist, which is fine)
	_, _, err := client.FaaS.CreateNamespace(ctx, namespace)
	if err != nil {
		log.Printf(LogNamespaceCreateWarningFmt, err)
	}

	// Get code from either code_inline or code_file
	var code string
	if codeInline, ok := d.GetOk(tfconstants.AttrCodeInline); ok {
		code = codeInline.(string)
	} else if codeFile, ok := d.GetOk(tfconstants.AttrCodeFile); ok {
		// TODO: In Phase 3, implement file upload logic
		// For now, treat code_file as inline code (path to be read)
		code = codeFile.(string)
		log.Print(LogCodeFileNotImplemented)
	}

	// Convert environment variables
	environment := make(map[string]string)
	if env, ok := d.GetOk(tfconstants.AttrEnvironmentVariables); ok {
		for k, v := range env.(map[string]interface{}) {
			environment[k] = v.(string)
		}
	}

	// Create the function
	createReq := &goe2e.FaasFunctionCreateRequest{
		Name:        d.Get(tfconstants.AttrName).(string),
		Namespace:   namespace,
		Runtime:     d.Get(tfconstants.AttrRuntime).(string),
		Code:        code,
		MemoryMB:    d.Get(tfconstants.AttrMemoryMB).(int),
		Timeout:     d.Get(tfconstants.AttrTimeoutSeconds).(int),
		MinReplicas: d.Get(tfconstants.AttrMinReplicas).(int),
		MaxReplicas: d.Get(tfconstants.AttrMaxReplicas).(int),
		Environment: environment,
	}

	fn, _, err := client.FaaS.CreateFunction(ctx, createReq)
	if err != nil {
		return diag.FromErr(err)
	}

	// SECURITY: Do not log full response to avoid exposing sensitive data (code field)
	log.Printf("[INFO] FAAS FUNCTION CREATE | ID: %s, Name: %s, Status: %s", fn.ID, fn.Name, fn.Status)

	// Set the ID immediately so it can be used if waiting fails
	d.SetId(fn.ID)

	// Wait for function to be ready
	if err := util.WaitForFunctionReady(ctx, client, fn.ID, tfconstants.FaaSCreateTimeout); err != nil {
		return diag.FromErr(fmt.Errorf("function created but failed to become ready: %w", err))
	}

	// Refresh function data after waiting
	fn, _, err = client.FaaS.GetFunction(ctx, fn.ID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading function after creation: %w", err))
	}
	if fn == nil {
		return diag.FromErr(errors.New(ErrFunctionNotFoundAfterCreate))
	}

	// Set the attributes
	if err := d.Set(tfconstants.AttrName, fn.Name); err != nil {
		return diag.FromErr(fmt.Errorf(tfconstants.ErrorSettingStateFormat(tfconstants.AttrName), err))
	}
	if err := d.Set(tfconstants.AttrNamespace, fn.Namespace); err != nil {
		return diag.FromErr(fmt.Errorf(tfconstants.ErrorSettingStateFormat(tfconstants.AttrNamespace), err))
	}
	if err := d.Set(tfconstants.AttrRuntime, fn.Runtime); err != nil {
		return diag.FromErr(fmt.Errorf(tfconstants.ErrorSettingStateFormat(tfconstants.AttrRuntime), err))
	}
	if err := d.Set(tfconstants.AttrMemoryMB, fn.MemoryMB); err != nil {
		return diag.FromErr(fmt.Errorf(tfconstants.ErrorSettingStateFormat(tfconstants.AttrMemoryMB), err))
	}
	if err := d.Set(tfconstants.AttrTimeoutSeconds, fn.Timeout); err != nil {
		return diag.FromErr(fmt.Errorf(tfconstants.ErrorSettingStateFormat(tfconstants.AttrTimeoutSeconds), err))
	}
	if err := d.Set(tfconstants.AttrMinReplicas, fn.MinReplicas); err != nil {
		return diag.FromErr(fmt.Errorf(tfconstants.ErrorSettingStateFormat(tfconstants.AttrMinReplicas), err))
	}
	if err := d.Set(tfconstants.AttrMaxReplicas, fn.MaxReplicas); err != nil {
		return diag.FromErr(fmt.Errorf(tfconstants.ErrorSettingStateFormat(tfconstants.AttrMaxReplicas), err))
	}
	if err := d.Set(tfconstants.AttrEndpointURL, fn.EndpointURL); err != nil {
		return diag.FromErr(fmt.Errorf(tfconstants.ErrorSettingStateFormat(tfconstants.AttrEndpointURL), err))
	}
	if err := d.Set(tfconstants.AttrStatus, fn.Status); err != nil {
		return diag.FromErr(fmt.Errorf(tfconstants.ErrorSettingStateFormat(tfconstants.AttrStatus), err))
	}
	if err := d.Set(tfconstants.AttrCreatedAt, fn.CreatedAt); err != nil {
		return diag.FromErr(fmt.Errorf(tfconstants.ErrorSettingStateFormat(tfconstants.AttrCreatedAt), err))
	}
	if err := d.Set(tfconstants.AttrUpdatedAt, fn.UpdatedAt); err != nil {
		return diag.FromErr(fmt.Errorf(tfconstants.ErrorSettingStateFormat(tfconstants.AttrUpdatedAt), err))
	}

	if fn.Environment != nil {
		if err := d.Set(tfconstants.AttrEnvironmentVariables, fn.Environment); err != nil {
			return diag.FromErr(fmt.Errorf(tfconstants.ErrorSettingStateFormat(tfconstants.AttrEnvironmentVariables), err))
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
		log.Printf(LogFunctionNotFoundWarningFmt, functionID)
		d.SetId("")

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "FaaS function not found",
			Detail:   "The FaaS function may have been deleted manually.",
		})

		return diags
	}

	if err := d.Set(tfconstants.AttrName, fn.Name); err != nil {
		return diag.FromErr(fmt.Errorf(tfconstants.ErrorSettingStateFormat(tfconstants.AttrName), err))
	}
	if err := d.Set(tfconstants.AttrNamespace, fn.Namespace); err != nil {
		return diag.FromErr(fmt.Errorf(tfconstants.ErrorSettingStateFormat(tfconstants.AttrNamespace), err))
	}
	if err := d.Set(tfconstants.AttrRuntime, fn.Runtime); err != nil {
		return diag.FromErr(fmt.Errorf(tfconstants.ErrorSettingStateFormat(tfconstants.AttrRuntime), err))
	}
	if err := d.Set(tfconstants.AttrMemoryMB, fn.MemoryMB); err != nil {
		return diag.FromErr(fmt.Errorf(tfconstants.ErrorSettingStateFormat(tfconstants.AttrMemoryMB), err))
	}
	if err := d.Set(tfconstants.AttrTimeoutSeconds, fn.Timeout); err != nil {
		return diag.FromErr(fmt.Errorf(tfconstants.ErrorSettingStateFormat(tfconstants.AttrTimeoutSeconds), err))
	}
	if err := d.Set(tfconstants.AttrMinReplicas, fn.MinReplicas); err != nil {
		return diag.FromErr(fmt.Errorf(tfconstants.ErrorSettingStateFormat(tfconstants.AttrMinReplicas), err))
	}
	if err := d.Set(tfconstants.AttrMaxReplicas, fn.MaxReplicas); err != nil {
		return diag.FromErr(fmt.Errorf(tfconstants.ErrorSettingStateFormat(tfconstants.AttrMaxReplicas), err))
	}
	if err := d.Set(tfconstants.AttrEndpointURL, fn.EndpointURL); err != nil {
		return diag.FromErr(fmt.Errorf(tfconstants.ErrorSettingStateFormat(tfconstants.AttrEndpointURL), err))
	}
	if err := d.Set(tfconstants.AttrStatus, fn.Status); err != nil {
		return diag.FromErr(fmt.Errorf(tfconstants.ErrorSettingStateFormat(tfconstants.AttrStatus), err))
	}
	if err := d.Set(tfconstants.AttrCreatedAt, fn.CreatedAt); err != nil {
		return diag.FromErr(fmt.Errorf(tfconstants.ErrorSettingStateFormat(tfconstants.AttrCreatedAt), err))
	}
	if err := d.Set(tfconstants.AttrUpdatedAt, fn.UpdatedAt); err != nil {
		return diag.FromErr(fmt.Errorf(tfconstants.ErrorSettingStateFormat(tfconstants.AttrUpdatedAt), err))
	}

	if fn.Environment != nil {
		if err := d.Set(tfconstants.AttrEnvironmentVariables, fn.Environment); err != nil {
			return diag.FromErr(fmt.Errorf(tfconstants.ErrorSettingStateFormat(tfconstants.AttrEnvironmentVariables), err))
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

	if d.HasChange(tfconstants.AttrCodeInline) {
		updateReq.Code = goe2e.String(d.Get(tfconstants.AttrCodeInline).(string))
		hasChanges = true
	}

	if d.HasChange(tfconstants.AttrCodeFile) {
		// TODO: In Phase 3, implement file upload logic
		updateReq.Code = goe2e.String(d.Get(tfconstants.AttrCodeFile).(string))
		hasChanges = true
		log.Print(LogCodeFileNotImplemented)
	}

	if d.HasChange(tfconstants.AttrMemoryMB) {
		updateReq.MemoryMB = goe2e.Int(d.Get(tfconstants.AttrMemoryMB).(int))
		hasChanges = true
	}

	if d.HasChange(tfconstants.AttrTimeoutSeconds) {
		updateReq.Timeout = goe2e.Int(d.Get(tfconstants.AttrTimeoutSeconds).(int))
		hasChanges = true
	}

	if d.HasChange(tfconstants.AttrMinReplicas) {
		updateReq.MinReplicas = goe2e.Int(d.Get(tfconstants.AttrMinReplicas).(int))
		hasChanges = true
	}

	if d.HasChange(tfconstants.AttrMaxReplicas) {
		updateReq.MaxReplicas = goe2e.Int(d.Get(tfconstants.AttrMaxReplicas).(int))
		hasChanges = true
	}

	if d.HasChange(tfconstants.AttrEnvironmentVariables) {
		environment := make(map[string]string)
		if env, ok := d.GetOk(tfconstants.AttrEnvironmentVariables); ok {
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

	// SECURITY: Do not log full response to avoid exposing sensitive data (code field)
	log.Printf("[INFO] FAAS FUNCTION UPDATE | ID: %s, Name: %s, Status: %s", fn.ID, fn.Name, fn.Status)

	// Wait for function to be ready after update
	if err := util.WaitForFunctionReady(ctx, client, functionID, tfconstants.FaaSUpdateTimeout); err != nil {
		return diag.FromErr(fmt.Errorf("function updated but failed to become ready: %w", err))
	}

	// Refresh function data after waiting
	fn, _, err = client.FaaS.GetFunction(ctx, functionID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading function after update: %w", err))
	}
	if fn == nil {
		return diag.FromErr(errors.New(ErrFunctionNotFoundAfterUpdate))
	}

	// Update the state with the response
	if err := d.Set(tfconstants.AttrMemoryMB, fn.MemoryMB); err != nil {
		return diag.FromErr(fmt.Errorf(tfconstants.ErrorSettingStateFormat(tfconstants.AttrMemoryMB), err))
	}
	if err := d.Set(tfconstants.AttrTimeoutSeconds, fn.Timeout); err != nil {
		return diag.FromErr(fmt.Errorf(tfconstants.ErrorSettingStateFormat(tfconstants.AttrTimeoutSeconds), err))
	}
	if err := d.Set(tfconstants.AttrMinReplicas, fn.MinReplicas); err != nil {
		return diag.FromErr(fmt.Errorf(tfconstants.ErrorSettingStateFormat(tfconstants.AttrMinReplicas), err))
	}
	if err := d.Set(tfconstants.AttrMaxReplicas, fn.MaxReplicas); err != nil {
		return diag.FromErr(fmt.Errorf(tfconstants.ErrorSettingStateFormat(tfconstants.AttrMaxReplicas), err))
	}
	if err := d.Set(tfconstants.AttrEndpointURL, fn.EndpointURL); err != nil {
		return diag.FromErr(fmt.Errorf(tfconstants.ErrorSettingStateFormat(tfconstants.AttrEndpointURL), err))
	}
	if err := d.Set(tfconstants.AttrStatus, fn.Status); err != nil {
		return diag.FromErr(fmt.Errorf(tfconstants.ErrorSettingStateFormat(tfconstants.AttrStatus), err))
	}
	if err := d.Set(tfconstants.AttrUpdatedAt, fn.UpdatedAt); err != nil {
		return diag.FromErr(fmt.Errorf(tfconstants.ErrorSettingStateFormat(tfconstants.AttrUpdatedAt), err))
	}

	if fn.Environment != nil {
		if err := d.Set(tfconstants.AttrEnvironmentVariables, fn.Environment); err != nil {
			return diag.FromErr(fmt.Errorf(tfconstants.ErrorSettingStateFormat(tfconstants.AttrEnvironmentVariables), err))
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
