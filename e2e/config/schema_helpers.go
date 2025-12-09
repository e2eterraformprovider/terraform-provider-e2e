package config

import (
	"sync"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// RegionSchema returns the region field schema
// Used by all resources and data sources that need region parameter
//
// Following HashiCorp's terraform-provider-aws pattern using sync.OnceValue
// to ensure singleton schema instances for better memory efficiency.
var RegionSchema = sync.OnceValue(func() *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeString,
		Optional:      true,
		ConflictsWith: []string{constants.AttrLocation},
		Description:   "Region for the resource. If not specified, uses provider default_region or E2E_REGION environment variable.",
	}
})

// LocationSchema returns the deprecated location field schema
// Kept for backwards compatibility, will be removed in v3.0.0
//
// This schema conflicts with 'region' to ensure only one is specified.
var LocationSchema = sync.OnceValue(func() *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeString,
		Optional:      true,
		Deprecated:    "Use 'region' instead. This parameter will be removed in version 3.0.0",
		ConflictsWith: []string{constants.AttrRegion},
		Description:   "DEPRECATED: Use 'region' instead",
	}
})

// ProjectIDSchemaResource returns the project_id field schema for resources
// Includes ForceNew since project_id changes require resource recreation
//
// Used by resources where project_id is a creation-time parameter.
var ProjectIDSchemaResource = sync.OnceValue(func() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "ID of the project as a string (e.g., \"12345\"). If not specified, uses provider default_project_id or E2E_PROJECT_ID environment variable.",
	}
})

// ProjectIDSchemaComputed returns a computed project_id schema
// Used by data sources where project_id can be computed from context
//
// This variant does NOT have ForceNew since data sources don't support it.
// It allows project_id to be both specified and computed.
var ProjectIDSchemaComputed = sync.OnceValue(func() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "ID of the project as a string (e.g., \"12345\"). If not specified, uses provider default_project_id or E2E_PROJECT_ID environment variable.",
	}
})

// ParseImportID parses import ID formats like "id" or "project/region/id"
// Returns a slice of parts split by '/'
//
// This is a helper function used by resource Import functions to handle
// both simple and complex import ID formats.
//
// Usage:
//
//	parts := config.ParseImportID(d.Id())
//	if len(parts) == 3 {
//	    // Format: project_id/region/resource_id
//	    projectID, region, resourceID := parts[0], parts[1], parts[2]
//	} else if len(parts) == 1 {
//	    // Format: resource_id (use provider defaults for project/region)
//	    resourceID := parts[0]
//	}
func ParseImportID(id string) []string {
	var parts []string
	start := 0
	for i := 0; i <= len(id); i++ {
		if i == len(id) || id[i] == '/' {
			if i > start {
				parts = append(parts, id[start:i])
			}
			start = i + 1
		}
	}
	return parts
}

// CDNOriginDetailsSchema returns the schema for CDN origin details block
// Used for configuring CDN distribution origin settings
var CDNOriginDetailsSchema = sync.OnceValue(func() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			constants.AttrPath: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Path to origin resource (e.g., '/api')",
			},
			constants.AttrSSLProtocol: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SSL/TLS protocol version (e.g., 'TLSv1.2')",
			},
			constants.AttrProtocolPolicy: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Protocol policy for origin connections (e.g., 'http-only', 'https-only', 'match-viewer')",
			},
			constants.AttrOriginReadTimeout: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Timeout for reading response from origin in seconds",
			},
			constants.AttrOriginKeepaliveTimeout: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Timeout for keeping idle connections to origin in seconds",
			},
		},
	}
})

// CDNCacheDetailsSchema returns the schema for CDN cache details block
// Used for configuring CDN distribution caching behavior
var CDNCacheDetailsSchema = sync.OnceValue(func() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			constants.AttrViewerProtocolPolicy: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Viewer protocol policy (e.g., 'allow-all', 'https-only', 'redirect-to-https')",
			},
			constants.AttrAllowedHTTPMethods: {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of allowed HTTP methods (e.g., ['GET', 'HEAD', 'POST', 'PUT', 'DELETE'])",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			constants.AttrDefaultTTL: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Default time-to-live for cached objects in seconds",
			},
			constants.AttrMinTTL: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Minimum time-to-live for cached objects in seconds",
			},
			constants.AttrMaxTTL: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Maximum time-to-live for cached objects in seconds",
			},
		},
	}
})

// CDNDomainDetailsSchema returns the schema for CDN domain details block
// Used for configuring CDN distribution domain settings
var CDNDomainDetailsSchema = sync.OnceValue(func() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			constants.AttrHTTPVersions: {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Supported HTTP versions (e.g., ['http/1.1', 'http/2'])",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			constants.AttrRootObject: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Default root object (e.g., 'index.html')",
			},
			constants.AttrIPv6Enabled: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable IPv6 for the distribution",
			},
		},
	}
})
