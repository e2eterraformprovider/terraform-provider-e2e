package goe2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

const (
	// Library version
	libraryVersion = "1.0.0"

	// Default API endpoint
	defaultBaseURL = "https://api.e2enetworks.com/myaccount/api/v1/"

	// User agent
	userAgent = "goe2e/" + libraryVersion

	// Media type
	mediaType = "application/json"

	// Retry configuration defaults
	defaultRetryMax     = 4
	defaultRetryWaitMax = 30 // seconds
	defaultRetryWaitMin = 1  // seconds
)

// Client manages communication with the E2E Networks API
type Client struct {
	// HTTP client used to communicate with the API
	client *http.Client

	// Base URL for API requests
	BaseURL *url.URL

	// User agent used when communicating with the API
	UserAgent string

	// E2E authentication credentials
	apiKey    string
	authToken string

	// Required parameters for all API calls
	projectID string
	region    string

	// Optional parameters for ai cloud API calls. They may get nuked later.
	workspaceID string
	teamID      string

	// Retry configuration
	RetryConfig RetryConfig

	// Custom headers to add to each request
	headers map[string]string

	// Services used for talking to different parts of the E2E API
	FaaS              FaasService
	Tags              TagService
	VolumeAttachment  VolumeAttachmentService
	Nodes             NodeService
	SSHKeys           SSHKeyService
	Vpcs              VpcService
	VPCTunnels        VPCTunnelService
	SecurityGroups    SecurityGroupService
	Autoscaling       AutoscalingService
	DBaaSMySQL        DBaaSMySQLService
	MariaDB           MariaDBService
	PostgreSQL        PostgreSQLService
	BlockStorage      BlockStorageService
	ContainerRegistry ContainerRegistryService
	LoadBalancer      LoadBalancerService
	ObjectStorage     ObjectStorageService
	Kubernetes        KubernetesService
	Sfs               SfsService
	ReserveIP         ReserveIPService
	Images            ImageService
	CDN               CDNService
	Backup            BackupService
	Firewall          FirewallService
}

// RetryConfig holds retry-specific configuration
type RetryConfig struct {
	RetryMax     int
	RetryWaitMin *time.Duration
	RetryWaitMax *time.Duration
}

// NewClient returns a new E2E API client with required parameters.
// This is the primary constructor for the client.
//
// All four parameters are required for every API call:
//   - apiKey: Your E2E API key
//   - authToken: Your E2E authentication token
//   - projectID: The project ID for API calls
//   - region: The region for API calls (e.g., "Mumbai", "Delhi", "Chennai")
//
// Example:
//
//	client, err := goe2e.NewClient("api-key", "auth-token", "project-123", "Mumbai")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	function, _, err := client.FaaS.GetFunction(ctx, "func-123")
func NewClient(apiKey, authToken, projectID, region string, opts ...ClientOpt) (*Client, error) {
	// Validate required parameters
	if apiKey == "" {
		return nil, NewArgError("apiKey", "cannot be empty")
	}
	if authToken == "" {
		return nil, NewArgError("authToken", "cannot be empty")
	}
	if projectID == "" {
		return nil, NewArgError("projectID", "cannot be empty")
	}
	if region == "" {
		return nil, NewArgError("region", "cannot be empty")
	}

	// Set up default retry configuration
	retryConfig := RetryConfig{
		RetryMax:     defaultRetryMax,
		RetryWaitMin: PtrTo(time.Duration(defaultRetryWaitMin) * time.Second),
		RetryWaitMax: PtrTo(time.Duration(defaultRetryWaitMax) * time.Second),
	}

	// Build options in order: required params, user options, then default retry config
	// This allows user to override default retry configuration
	fullOpts := []ClientOpt{
		setAPIKey(apiKey),
		setAuthToken(authToken),
		setProjectID(projectID),
		setRegion(region),
	}

	// Add user options
	fullOpts = append(fullOpts, opts...)

	// Add default retry config (allows user to override if they provided their own)
	fullOpts = append(fullOpts, WithRetryAndBackoffs(retryConfig))

	opts = fullOpts

	return New(nil, opts...)
}

// New returns a new E2E API client instance with full control.
// Use this for advanced scenarios where you need a custom HTTP client.
// For normal usage, use NewClient instead.
func New(httpClient *http.Client, opts ...ClientOpt) (*Client, error) {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 60 * time.Second}
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{
		client:    httpClient,
		BaseURL:   baseURL,
		UserAgent: userAgent,
		headers:   make(map[string]string),
		RetryConfig: RetryConfig{
			RetryMax: -1, // Sentinel: not yet configured
		},
	}

	// Apply all options
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	// Validate required fields were set via options
	if c.apiKey == "" || c.authToken == "" || c.projectID == "" || c.region == "" {
		return nil, fmt.Errorf("apiKey, authToken, projectID, and region are required (use NewClient or provide via options)")
	}

	// Initialize services
	c.FaaS = &FaasServiceOp{client: c}
	c.Tags = &TagServiceOp{client: c}
	c.VolumeAttachment = &VolumeAttachmentServiceOp{client: c}
	c.Nodes = &NodeServiceOp{client: c}
	c.SSHKeys = &SSHKeyServiceOp{client: c}
	c.Vpcs = &VpcServiceOp{client: c}
	c.VPCTunnels = &VPCTunnelServiceOp{client: c}
	c.SecurityGroups = &SecurityGroupServiceOp{client: c}
	c.Autoscaling = &AutoscalingServiceOp{client: c}
	c.DBaaSMySQL = &DBaaSMySQLServiceOp{client: c}
	c.MariaDB = &MariaDBServiceOp{client: c}
	c.PostgreSQL = &PostgreSQLServiceOp{client: c}
	c.BlockStorage = &BlockStorageServiceOp{client: c}
	c.ContainerRegistry = &ContainerRegistryServiceOp{client: c}
	c.LoadBalancer = &LoadBalancerServiceOp{client: c}
	c.Kubernetes = &KubernetesServiceOp{client: c}
	c.ObjectStorage = &ObjectStorageServiceOp{client: c}
	c.Sfs = &SfsServiceOp{client: c}
	c.ReserveIP = &ReserveIPServiceOp{client: c}
	c.Images = &ImageServiceOp{client: c}
	c.CDN = &CDNServiceOp{client: c}
	c.Backup = &BackupServiceOp{client: c}
	c.Firewall = &FirewallServiceOp{client: c}

	return c, nil
}

// ClientOpt are options for New.
type ClientOpt func(*Client) error

// Internal setters (not exported - used by NewClient)
func setAPIKey(apiKey string) ClientOpt {
	return func(c *Client) error {
		c.apiKey = apiKey
		return nil
	}
}

func setAuthToken(token string) ClientOpt {
	return func(c *Client) error {
		c.authToken = token
		return nil
	}
}

func setProjectID(projectID string) ClientOpt {
	return func(c *Client) error {
		c.projectID = projectID
		return nil
	}
}

func setRegion(region string) ClientOpt {
	return func(c *Client) error {
		c.region = region
		return nil
	}
}

// Public options (exported - users can use these for advanced scenarios)

// WithWorkspace sets the workspace ID (optional, for multi-tenant scenarios).
// When set, the workspace_id query parameter will be added to all API calls.
func WithWorkspace(workspaceID string) ClientOpt {
	return func(c *Client) error {
		c.workspaceID = workspaceID
		return nil
	}
}

// WithTeam sets the team ID (optional, for team-based scenarios).
// When set, the team_id query parameter will be added to all API calls.
func WithTeam(teamID string) ClientOpt {
	return func(c *Client) error {
		c.teamID = teamID
		return nil
	}
}

// SetBaseURL sets the base URL for API requests.
// If an invalid URL is provided, an error is returned.
func SetBaseURL(baseURL string) ClientOpt {
	return func(c *Client) error {
		u, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.BaseURL = u
		return nil
	}
}

// SetUserAgent sets a custom user agent string
func SetUserAgent(ua string) ClientOpt {
	return func(c *Client) error {
		c.UserAgent = ua
		return nil
	}
}

// WithRetryAndBackoffs sets the retry policy with exponential backoff.
// When RetryMax is greater than 0, the client will automatically retry
// failed requests with exponential backoff.
// If retry config is already set (RetryMax != -1), this will be skipped to respect
// user configuration that explicitly sets retry behavior (including disabling).
func WithRetryAndBackoffs(retryConfig RetryConfig) ClientOpt {
	return func(c *Client) error {
		// Don't override if retry config is already explicitly set by user
		if c.RetryConfig.RetryMax != -1 {
			return nil
		}

		c.RetryConfig = retryConfig

		if retryConfig.RetryMax > 0 {
			retryableClient := retryablehttp.NewClient()
			retryableClient.RetryMax = retryConfig.RetryMax

			if retryConfig.RetryWaitMin != nil {
				retryableClient.RetryWaitMin = *retryConfig.RetryWaitMin
			}
			if retryConfig.RetryWaitMax != nil {
				retryableClient.RetryWaitMax = *retryConfig.RetryWaitMax
			}

			retryableClient.CheckRetry = retryablehttp.DefaultRetryPolicy
			retryableClient.HTTPClient = c.client
			retryableClient.ErrorHandler = retryablehttp.PassthroughErrorHandler

			c.client = retryableClient.StandardClient()
		}

		return nil
	}
}

// SetStaticRateLimit configures rate limiting (requests per second).
// TODO: Implement rate limiting similar to DigitalOcean's approach using rate.Limiter.
// For now, this is a placeholder that accepts the configuration but doesn't enforce it.
func SetStaticRateLimit(rps float64) ClientOpt {
	return func(c *Client) error {
		// TODO: Implement rate limiting
		// Example implementation would use golang.org/x/time/rate:
		// c.rateLimiter = rate.NewLimiter(rate.Limit(rps), int(rps))
		return nil
	}
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// which will be resolved to the BaseURL of the Client. Relative URLs should
// always be specified without a preceding slash. If specified, the value
// pointed to by body is JSON encoded and included as the request body.
//
// Standard query parameters (apikey, project_id, location) are automatically
// added from the client configuration.
func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	// Add standard E2E query parameters from client
	q := u.Query()
	q.Add("apikey", c.apiKey)
	q.Add("project_id", c.projectID)
	q.Add("location", c.region) // API uses 'location' parameter

	// Add optional parameters if set
	if c.workspaceID != "" {
		q.Add("workspace_id", c.workspaceID)
	}
	if c.teamID != "" {
		q.Add("team_id", c.teamID)
	}

	u.RawQuery = q.Encode()

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		if err := enc.Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	// Set standard headers
	req.Header.Set("Authorization", "Bearer "+c.authToken)
	req.Header.Set("Content-Type", mediaType)
	req.Header.Set("Accept", mediaType)
	req.Header.Set("User-Agent", c.UserAgent)

	// Add any custom headers
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	return req, nil
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer interface,
// the raw response body will be written to v, without attempting to first
// decode it.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		// If the context was canceled, return immediately
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return nil, err
	}
	defer resp.Body.Close()

	response := newResponse(resp)

	err = CheckResponse(resp)
	if err != nil {
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			decErr := json.NewDecoder(resp.Body).Decode(v)
			if decErr == io.EOF {
				decErr = nil // ignore EOF errors caused by empty response body
			}
			if decErr != nil {
				err = decErr
			}
		}
	}

	return response, err
}

// Response is an E2E Networks API response. This wraps the standard http.Response
// returned from E2E Networks.
type Response struct {
	*http.Response
}

func newResponse(r *http.Response) *Response {
	return &Response{Response: r}
}

// CheckResponse checks the API response for errors, and returns them if present.
// A response is considered an error if it has a status code outside the 200 range.
// API error responses are expected to have either no response body, or a JSON
// response body that maps to ErrorResponse.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := io.ReadAll(r.Body)
	if err == nil && data != nil {
		_ = json.Unmarshal(data, errorResponse)
	}

	// Re-populate response body for later use
	r.Body = io.NopCloser(bytes.NewBuffer(data))

	return errorResponse
}

// RequestOptions is deprecated and will be removed in v2.0.
// Use the client-level configuration instead (projectID and region are set at client creation).
//
// Deprecated: Parameters are now set at client creation time via NewClient.
type RequestOptions struct {
	ProjectID string // Deprecated: Set via NewClient projectID parameter
	Region    string // Deprecated: Set via NewClient region parameter
}
