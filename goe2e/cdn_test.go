package goe2e

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateCDNDistribution(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPost)
		testURLPath(t, r, "/cdn/distributions")
		assertQueryParam(t, r, "apikey", "test-api-key")
		assertQueryParam(t, r, "project_id", "test-project")
		assertQueryParam(t, r, "location", "test-location")

		data := map[string]interface{}{
			"domain_id":       "cdn-123",
			"domain_name":     "test.example.com",
			"e2e_domain_name": "test.cdn.e2enetworks.net",
			"source":          "origin.example.com",
			"is_enabled":      true,
			"state":           "DEPLOYING",
			"created_at":      "2024-01-01T00:00:00Z",
			"updated_at":      "2024-01-01T00:00:00Z",
		}
		writeJSON(w, http.StatusCreated, buildSuccessResponse(201, "CDN distribution created successfully", data))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	createReq := &CreateCDNDistributionRequest{
		DomainName: "test.example.com",
		Source:     "origin.example.com",
		IsEnabled:  Bool(true),
	}

	result, _, err := client.CDN.CreateCDNDistribution(context.Background(), createReq)
	assertNoError(t, err)
	assertNotNil(t, result, "Expected result, got nil")

	if result.DomainID != "cdn-123" {
		t.Errorf("Expected DomainID cdn-123, got %s", result.DomainID)
	}

	if result.DomainName != "test.example.com" {
		t.Errorf("Expected DomainName test.example.com, got %s", result.DomainName)
	}

	if result.E2EDomainName != "test.cdn.e2enetworks.net" {
		t.Errorf("Expected E2EDomainName test.cdn.e2enetworks.net, got %s", result.E2EDomainName)
	}

	if result.State != "DEPLOYING" {
		t.Errorf("Expected State DEPLOYING, got %s", result.State)
	}

	if !result.IsEnabled {
		t.Errorf("Expected IsEnabled true, got false")
	}
}

func TestCreateCDNDistribution_WithDetails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPost)
		testURLPath(t, r, "/cdn/distributions")

		data := map[string]interface{}{
			"domain_id":       "cdn-456",
			"domain_name":     "test2.example.com",
			"e2e_domain_name": "test2.cdn.e2enetworks.net",
			"source":          "origin.example.com",
			"is_enabled":      true,
			"state":           "ACTIVE",
			"origin_details": map[string]interface{}{
				"path":                     "/api",
				"ssl_protocol":             "TLSv1.2",
				"protocol_policy":          "https-only",
				"origin_read_timeout":      30,
				"origin_keepalive_timeout": 60,
			},
			"cache_details": map[string]interface{}{
				"viewer_protocol_policy": "redirect-to-https",
				"allowed_http_methods":   []string{"GET", "HEAD", "POST"},
				"default_ttl":            3600,
				"min_ttl":                0,
				"max_ttl":                86400,
			},
			"domain_details": map[string]interface{}{
				"http_versions": []string{"http/1.1", "http/2"},
				"root_object":   "index.html",
				"ipv6_enabled":  true,
			},
			"created_at": "2024-01-01T00:00:00Z",
		}
		writeJSON(w, http.StatusCreated, buildSuccessResponse(201, "CDN distribution created successfully", data))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	createReq := &CreateCDNDistributionRequest{
		DomainName: "test2.example.com",
		Source:     "origin.example.com",
		IsEnabled:  Bool(true),
		OriginDetails: &CDNOriginDetails{
			Path:                   "/api",
			SSLProtocol:            "TLSv1.2",
			ProtocolPolicy:         "https-only",
			OriginReadTimeout:      30,
			OriginKeepaliveTimeout: 60,
		},
		CacheDetails: &CDNCacheDetails{
			ViewerProtocolPolicy: "redirect-to-https",
			AllowedHTTPMethods:   []string{"GET", "HEAD", "POST"},
			DefaultTTL:           3600,
			MinTTL:               0,
			MaxTTL:               86400,
		},
		DomainDetails: &CDNDomainDetails{
			HTTPVersions: []string{"http/1.1", "http/2"},
			RootObject:   "index.html",
			IPv6Enabled:  true,
		},
	}

	result, _, err := client.CDN.CreateCDNDistribution(context.Background(), createReq)
	assertNoError(t, err)
	assertNotNil(t, result, "Expected result, got nil")

	if result.DomainID != "cdn-456" {
		t.Errorf("Expected DomainID cdn-456, got %s", result.DomainID)
	}

	if result.State != "ACTIVE" {
		t.Errorf("Expected State ACTIVE, got %s", result.State)
	}

	// Verify origin details
	assertNotNil(t, result.OriginDetails, "Expected OriginDetails, got nil")
	if result.OriginDetails.Path != "/api" {
		t.Errorf("Expected Path /api, got %s", result.OriginDetails.Path)
	}
	if result.OriginDetails.SSLProtocol != "TLSv1.2" {
		t.Errorf("Expected SSLProtocol TLSv1.2, got %s", result.OriginDetails.SSLProtocol)
	}

	// Verify cache details
	assertNotNil(t, result.CacheDetails, "Expected CacheDetails, got nil")
	if result.CacheDetails.DefaultTTL != 3600 {
		t.Errorf("Expected DefaultTTL 3600, got %d", result.CacheDetails.DefaultTTL)
	}
	if len(result.CacheDetails.AllowedHTTPMethods) != 3 {
		t.Errorf("Expected 3 AllowedHTTPMethods, got %d", len(result.CacheDetails.AllowedHTTPMethods))
	}

	// Verify domain details
	assertNotNil(t, result.DomainDetails, "Expected DomainDetails, got nil")
	if result.DomainDetails.RootObject != "index.html" {
		t.Errorf("Expected RootObject index.html, got %s", result.DomainDetails.RootObject)
	}
	if !result.DomainDetails.IPv6Enabled {
		t.Error("Expected IPv6Enabled true, got false")
	}
}

func TestGetCDNDistributions_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/cdn/distributions")
		assertQueryParam(t, r, "apikey", "test-api-key")
		assertQueryParam(t, r, "project_id", "test-project")
		assertQueryParam(t, r, "location", "test-location")

		data := []interface{}{
			map[string]interface{}{
				"domain_id":       "cdn-123",
				"domain_name":     "test1.example.com",
				"e2e_domain_name": "test1.cdn.e2enetworks.net",
				"source":          "origin1.example.com",
				"is_enabled":      true,
				"state":           "ACTIVE",
			},
			map[string]interface{}{
				"domain_id":       "cdn-456",
				"domain_name":     "test2.example.com",
				"e2e_domain_name": "test2.cdn.e2enetworks.net",
				"source":          "origin2.example.com",
				"is_enabled":      false,
				"state":           "DISABLED",
			},
		}
		writeJSON(w, http.StatusOK, buildListResponse(data))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	result, _, err := client.CDN.GetCDNDistributions(context.Background(), "")
	assertNoError(t, err)
	assertNotNil(t, result, "Expected result, got nil")

	if len(result) != 2 {
		t.Fatalf("Expected 2 distributions, got %d", len(result))
	}

	if result[0].DomainID != "cdn-123" {
		t.Errorf("Expected first distribution ID cdn-123, got %s", result[0].DomainID)
	}

	if result[1].DomainID != "cdn-456" {
		t.Errorf("Expected second distribution ID cdn-456, got %s", result[1].DomainID)
	}
}

func TestGetCDNDistributions_ByID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/cdn/distributions")
		assertQueryParam(t, r, "domain_id", "cdn-123")
		assertQueryParam(t, r, "apikey", "test-api-key")
		assertQueryParam(t, r, "project_id", "test-project")
		assertQueryParam(t, r, "location", "test-location")

		data := []interface{}{
			map[string]interface{}{
				"domain_id":       "cdn-123",
				"domain_name":     "test.example.com",
				"e2e_domain_name": "test.cdn.e2enetworks.net",
				"source":          "origin.example.com",
				"is_enabled":      true,
				"state":           "ACTIVE",
				"created_at":      "2024-01-01T00:00:00Z",
			},
		}
		writeJSON(w, http.StatusOK, buildListResponse(data))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	result, _, err := client.CDN.GetCDNDistributions(context.Background(), "cdn-123")
	assertNoError(t, err)
	assertNotNil(t, result, "Expected result, got nil")

	if len(result) != 1 {
		t.Fatalf("Expected 1 distribution, got %d", len(result))
	}

	if result[0].DomainID != "cdn-123" {
		t.Errorf("Expected DomainID cdn-123, got %s", result[0].DomainID)
	}
}

func TestGetCDNDistributions_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	result, _, err := client.CDN.GetCDNDistributions(context.Background(), "nonexistent")
	assertNoError(t, err)

	if result != nil {
		t.Errorf("Expected nil result for 404, got: %v", result)
	}
}

func TestUpdateCDNDistribution(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/cdn/distributions")
		assertQueryParam(t, r, "domain_id", "cdn-123")
		assertQueryParam(t, r, "apikey", "test-api-key")
		assertQueryParam(t, r, "project_id", "test-project")
		assertQueryParam(t, r, "location", "test-location")

		data := map[string]interface{}{
			"domain_id":       "cdn-123",
			"domain_name":     "test.example.com",
			"e2e_domain_name": "test.cdn.e2enetworks.net",
			"source":          "origin.example.com",
			"is_enabled":      false,
			"state":           "DISABLED",
			"updated_at":      "2024-01-01T01:00:00Z",
		}
		writeJSON(w, http.StatusOK, buildSuccessResponse(200, "CDN distribution updated successfully", data))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	updateReq := &UpdateCDNDistributionRequest{
		IsEnabled: Bool(false),
	}

	result, _, err := client.CDN.UpdateCDNDistribution(context.Background(), "cdn-123", updateReq)
	assertNoError(t, err)
	assertNotNil(t, result, "Expected result, got nil")

	if result.DomainID != "cdn-123" {
		t.Errorf("Expected DomainID cdn-123, got %s", result.DomainID)
	}

	if result.IsEnabled {
		t.Error("Expected IsEnabled false, got true")
	}

	if result.State != "DISABLED" {
		t.Errorf("Expected State DISABLED, got %s", result.State)
	}
}

func TestDeleteCDNDistribution(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodDelete)
		testURLPath(t, r, "/cdn/distributions")
		assertQueryParam(t, r, "domain_id", "cdn-123")
		assertQueryParam(t, r, "apikey", "test-api-key")
		assertQueryParam(t, r, "project_id", "test-project")
		assertQueryParam(t, r, "location", "test-location")

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, err = client.CDN.DeleteCDNDistribution(context.Background(), "cdn-123")
	assertNoError(t, err)
}

// Edge case tests for better coverage
func TestCreateCDNDistribution_NilRequest(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, err = client.CDN.CreateCDNDistribution(context.Background(), nil)
	assertError(t, err, "")
}

func TestCreateCDNDistribution_EmptyDomainName(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	createReq := &CreateCDNDistributionRequest{
		DomainName: "",
		Source:     "origin.example.com",
	}

	_, _, err = client.CDN.CreateCDNDistribution(context.Background(), createReq)
	assertError(t, err, "")
}

func TestCreateCDNDistribution_EmptySource(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	createReq := &CreateCDNDistributionRequest{
		DomainName: "test.example.com",
		Source:     "",
	}

	_, _, err = client.CDN.CreateCDNDistribution(context.Background(), createReq)
	assertError(t, err, "")
}

func TestUpdateCDNDistribution_EmptyID(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	updateReq := &UpdateCDNDistributionRequest{
		IsEnabled: Bool(false),
	}

	_, _, err = client.CDN.UpdateCDNDistribution(context.Background(), "", updateReq)
	assertError(t, err, "")
}

func TestUpdateCDNDistribution_NilRequest(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, err = client.CDN.UpdateCDNDistribution(context.Background(), "cdn-123", nil)
	assertError(t, err, "")
}

func TestDeleteCDNDistribution_EmptyID(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, err = client.CDN.DeleteCDNDistribution(context.Background(), "")
	assertError(t, err, "")
}

// Error response tests
func TestCreateCDNDistribution_ErrorResponse(t *testing.T) {
	server := newErrorServer(t, http.StatusInternalServerError, "Internal server error")
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	createReq := &CreateCDNDistributionRequest{
		DomainName: "test.example.com",
		Source:     "origin.example.com",
	}

	_, _, err = client.CDN.CreateCDNDistribution(context.Background(), createReq)
	assertError(t, err, "")
}

func TestGetCDNDistributions_ErrorResponse(t *testing.T) {
	server := newErrorServer(t, http.StatusInternalServerError, "Internal server error")
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, err = client.CDN.GetCDNDistributions(context.Background(), "cdn-123")
	assertError(t, err, "")
}

func TestUpdateCDNDistribution_ErrorResponse(t *testing.T) {
	server := newErrorServer(t, http.StatusInternalServerError, "Internal server error")
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	updateReq := &UpdateCDNDistributionRequest{
		IsEnabled: Bool(false),
	}

	_, _, err = client.CDN.UpdateCDNDistribution(context.Background(), "cdn-123", updateReq)
	assertError(t, err, "")
}

func TestDeleteCDNDistribution_ErrorResponse(t *testing.T) {
	server := newErrorServer(t, http.StatusInternalServerError, "Internal server error")
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, err = client.CDN.DeleteCDNDistribution(context.Background(), "cdn-123")
	assertError(t, err, "")
}

func TestCreateCDNDistribution_BadRequest(t *testing.T) {
	server := newErrorServer(t, http.StatusBadRequest, "Invalid domain name format")
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	createReq := &CreateCDNDistributionRequest{
		DomainName: "invalid domain",
		Source:     "origin.example.com",
	}

	_, _, err = client.CDN.CreateCDNDistribution(context.Background(), createReq)
	assertError(t, err, "")
}

func TestGetCDNDistributions_EmptyList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodGet)
		writeJSON(w, http.StatusOK, buildListResponse([]interface{}{}))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	result, _, err := client.CDN.GetCDNDistributions(context.Background(), "")
	assertNoError(t, err)

	if result == nil {
		t.Fatal("Expected empty slice, got nil")
	}

	if len(result) != 0 {
		t.Errorf("Expected empty list, got %d distributions", len(result))
	}
}
