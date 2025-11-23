package goe2e

import (
	"bytes"
	"context"
	"encoding/json"
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

	// Retry configuration
	RetryConfig RetryConfig

	// Custom headers to add to each request
	headers map[string]string

	// Services used for talking to different parts of the E2E API
	// Phase 1: Only FaaS
	FaaS FaasService
}

// RetryConfig holds retry-specific configuration
type RetryConfig struct {
	RetryMax     int
	RetryWaitMin *time.Duration
	RetryWaitMax *time.Duration
}

// NewClient returns a new E2E API client with credentials and retry support.
// This is the primary constructor for the client.
func NewClient(apiKey, authToken string, opts ...ClientOpt) (*Client, error) {
	// Set up default retry configuration
	retryConfig := RetryConfig{
		RetryMax:     defaultRetryMax,
		RetryWaitMin: PtrTo(time.Duration(defaultRetryWaitMin) * time.Second),
		RetryWaitMax: PtrTo(time.Duration(defaultRetryWaitMax) * time.Second),
	}

	// Prepend credentials and retry config to options
	opts = append([]ClientOpt{
		SetAPIKey(apiKey),
		SetAuthToken(authToken),
		WithRetryAndBackoffs(retryConfig),
	}, opts...)

	return New(nil, opts...)
}

// New returns a new E2E API client instance.
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
	}

	// Apply all options
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	// Initialize services - Phase 1: only FaaS
	c.FaaS = &FaasServiceOp{client: c}

	return c, nil
}

// ClientOpt are options for New.
type ClientOpt func(*Client) error

// SetAPIKey sets the API key for authentication
func SetAPIKey(apiKey string) ClientOpt {
	return func(c *Client) error {
		c.apiKey = apiKey
		return nil
	}
}

// SetAuthToken sets the authentication token
func SetAuthToken(token string) ClientOpt {
	return func(c *Client) error {
		c.authToken = token
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
func WithRetryAndBackoffs(retryConfig RetryConfig) ClientOpt {
	return func(c *Client) error {
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
func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

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

// RequestOptions contains the common parameters required for E2E API requests.
// These parameters (ProjectID and Location) are required for every API call.
type RequestOptions struct {
	ProjectID string // The project ID
	Location  string // The region/location (e.g., "Chennai", "Delhi", "Mumbai")
}
