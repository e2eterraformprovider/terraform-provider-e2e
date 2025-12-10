package goe2e

import (
	"context"
	"fmt"
	"net/http"
)

const (
	cdnDistributionsPath = "cdn/distributions"
)

// CDNService is an interface for interacting with the CDN endpoints
// of the E2E Networks API.
type CDNService interface {
	CreateCDNDistribution(context.Context, *CreateCDNDistributionRequest) (*CDNDistribution, *Response, error)
	GetCDNDistributions(context.Context, string) ([]*CDNDistribution, *Response, error)
	UpdateCDNDistribution(context.Context, string, *UpdateCDNDistributionRequest) (*CDNDistribution, *Response, error)
	DeleteCDNDistribution(context.Context, string) (*Response, error)
}

// CDNServiceOp handles communication with CDN related methods of the
// E2E Networks API.
type CDNServiceOp struct {
	client *Client
}

var _ CDNService = &CDNServiceOp{}

// CDNDistribution represents a CDN distribution
type CDNDistribution struct {
	DomainID      string            `json:"domain_id"`
	DomainName    string            `json:"domain_name"`
	E2EDomainName string            `json:"e2e_domain_name"`
	Source        string            `json:"source"`
	IsEnabled     bool              `json:"is_enabled"`
	State         string            `json:"state"`
	OriginDetails *CDNOriginDetails `json:"origin_details,omitempty"`
	CacheDetails  *CDNCacheDetails  `json:"cache_details,omitempty"`
	DomainDetails *CDNDomainDetails `json:"domain_details,omitempty"`
	CreatedAt     string            `json:"created_at,omitempty"`
	UpdatedAt     string            `json:"updated_at,omitempty"`
}

// CDNOriginDetails represents CDN origin configuration
type CDNOriginDetails struct {
	Path                   string `json:"path,omitempty"`
	SSLProtocol            string `json:"ssl_protocol,omitempty"`
	ProtocolPolicy         string `json:"protocol_policy,omitempty"`
	OriginReadTimeout      int    `json:"origin_read_timeout,omitempty"`
	OriginKeepaliveTimeout int    `json:"origin_keepalive_timeout,omitempty"`
}

// CDNCacheDetails represents CDN cache configuration
type CDNCacheDetails struct {
	ViewerProtocolPolicy string   `json:"viewer_protocol_policy,omitempty"`
	AllowedHTTPMethods   []string `json:"allowed_http_methods,omitempty"`
	DefaultTTL           int      `json:"default_ttl,omitempty"`
	MinTTL               int      `json:"min_ttl,omitempty"`
	MaxTTL               int      `json:"max_ttl,omitempty"`
}

// CDNDomainDetails represents CDN domain configuration
type CDNDomainDetails struct {
	HTTPVersions []string `json:"http_versions,omitempty"`
	RootObject   string   `json:"root_object,omitempty"`
	IPv6Enabled  bool     `json:"ipv6_enabled,omitempty"`
}

// CreateCDNDistributionRequest represents a request to create a CDN distribution
type CreateCDNDistributionRequest struct {
	DomainName    string            `json:"domain_name"`
	Source        string            `json:"source"`
	IsEnabled     *bool             `json:"is_enabled,omitempty"`
	OriginDetails *CDNOriginDetails `json:"origin_details,omitempty"`
	CacheDetails  *CDNCacheDetails  `json:"cache_details,omitempty"`
	DomainDetails *CDNDomainDetails `json:"domain_details,omitempty"`
}

// UpdateCDNDistributionRequest represents a request to update a CDN distribution
type UpdateCDNDistributionRequest struct {
	IsEnabled *bool `json:"is_enabled"`
}

// Response wrappers for API calls
type cdnDistributionRoot struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    CDNDistribution `json:"data"`
}

type cdnDistributionsRoot struct {
	Code    int                `json:"code"`
	Message string             `json:"message"`
	Data    []*CDNDistribution `json:"data"`
}

// CreateCDNDistribution creates a new CDN distribution.
func (s *CDNServiceOp) CreateCDNDistribution(ctx context.Context, createReq *CreateCDNDistributionRequest) (*CDNDistribution, *Response, error) {
	if createReq == nil {
		return nil, nil, NewArgError("createReq", "cannot be nil")
	}
	if createReq.DomainName == "" {
		return nil, nil, NewArgError("domain_name", "cannot be empty")
	}
	if createReq.Source == "" {
		return nil, nil, NewArgError("source", "cannot be empty")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, cdnDistributionsPath, createReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for CDN distribution (%s): %w", createReq.DomainName, err)
	}

	root := new(cdnDistributionRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to create CDN distribution (%s): %w", createReq.DomainName, err)
	}

	return &root.Data, resp, nil
}

// GetCDNDistributions retrieves CDN distributions. If domainID is provided, retrieves a specific distribution.
// If domainID is empty, retrieves all distributions.
func (s *CDNServiceOp) GetCDNDistributions(ctx context.Context, domainID string) ([]*CDNDistribution, *Response, error) {
	path := cdnDistributionsPath

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for CDN distributions: %w", err)
	}

	// Add domain_id as query parameter if provided
	if domainID != "" {
		q := req.URL.Query()
		q.Add("domain_id", domainID)
		req.URL.RawQuery = q.Encode()
	}

	root := new(cdnDistributionsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		// Return nil distributions for 404 (not found)
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, resp, nil
		}
		if domainID != "" {
			return nil, resp, fmt.Errorf("failed to retrieve CDN distribution (ID: %s): %w", domainID, err)
		}
		return nil, resp, fmt.Errorf("failed to retrieve CDN distributions: %w", err)
	}

	return root.Data, resp, nil
}

// UpdateCDNDistribution updates a CDN distribution (primarily for enabling/disabling).
func (s *CDNServiceOp) UpdateCDNDistribution(ctx context.Context, domainID string, updateReq *UpdateCDNDistributionRequest) (*CDNDistribution, *Response, error) {
	if domainID == "" {
		return nil, nil, NewArgError("domainID", "cannot be empty")
	}
	if updateReq == nil {
		return nil, nil, NewArgError("updateReq", "cannot be nil")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPut, cdnDistributionsPath, updateReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for updating CDN distribution (ID: %s): %w", domainID, err)
	}

	// Add domain_id as query parameter
	q := req.URL.Query()
	q.Add("domain_id", domainID)
	req.URL.RawQuery = q.Encode()

	root := new(cdnDistributionRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to update CDN distribution (ID: %s): %w", domainID, err)
	}

	return &root.Data, resp, nil
}

// DeleteCDNDistribution deletes a CDN distribution.
func (s *CDNServiceOp) DeleteCDNDistribution(ctx context.Context, domainID string) (*Response, error) {
	if domainID == "" {
		return nil, NewArgError("domainID", "cannot be empty")
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, cdnDistributionsPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for deleting CDN distribution (ID: %s): %w", domainID, err)
	}

	// Add domain_id as query parameter
	q := req.URL.Query()
	q.Add("domain_id", domainID)
	req.URL.RawQuery = q.Encode()

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to delete CDN distribution (ID: %s): %w", domainID, err)
	}
	return resp, nil
}
