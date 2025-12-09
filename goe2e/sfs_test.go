package goe2e

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateSfs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPost)
		assertQueryParam(t, r, "apikey", "test-api-key")
		assertQueryParam(t, r, "project_id", "test-project")
		assertQueryParam(t, r, "location", "test-location")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(buildSuccessResponse(201, "SFS created successfully", map[string]interface{}{
			"efs_id":               12345,
			"is_credit_sufficient": true,
		})))
	}))
	defer server.Close()
	client, err := NewClient("test-api-key", "test-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	createReq := &SfsCreateRequest{
		Name:                "test-sfs",
		Plan:                "basic",
		VPCID:               "vpc-123",
		DiskSize:            100,
		DiskIOPS:            3000,
		IsEncryptionEnabled: false,
	}
	result, _, err := client.Sfs.CreateSfs(context.Background(), createReq)
	assertNoError(t, err)
	assertNotNil(t, result, "Expected result, got nil")
	if result.ID != "12345" {
		t.Errorf("Expected ID 12345, got %s", result.ID)
	}
	if result.Name != "test-sfs" {
		t.Errorf("Expected Name test-sfs, got %s", result.Name)
	}
}

func TestCreateSfs_WithEncryption(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(buildSuccessResponse(201, "SFS created successfully", map[string]interface{}{
			"efs_id":               67890,
			"is_credit_sufficient": true,
		})))
	}))
	defer server.Close()
	client, err := NewClient("test-api-key", "test-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	createReq := &SfsCreateRequest{
		Name:                 "encrypted-sfs",
		Plan:                 "premium",
		VPCID:                "vpc-456",
		DiskSize:             200,
		DiskIOPS:             5000,
		IsEncryptionEnabled:  true,
		EncryptionPassphrase: "secret-passphrase",
	}
	result, _, err := client.Sfs.CreateSfs(context.Background(), createReq)
	assertNoError(t, err)
	if result.ID != "67890" {
		t.Errorf("Expected ID 67890, got %s", result.ID)
	}
}

func TestCreateSfs_StringID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(buildSuccessResponse(201, "SFS created successfully", map[string]interface{}{
			"efs_id":               "sfs-12345",
			"is_credit_sufficient": true,
		})))
	}))
	defer server.Close()
	client, err := NewClient("test-api-key", "test-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	createReq := &SfsCreateRequest{
		Name:     "test-sfs",
		Plan:     "basic",
		VPCID:    "vpc-123",
		DiskSize: 100,
		DiskIOPS: 3000,
	}
	result, _, err := client.Sfs.CreateSfs(context.Background(), createReq)
	assertNoError(t, err)
	if result.ID != "sfs-12345" {
		t.Errorf("Expected ID sfs-12345, got %s", result.ID)
	}
}
func TestGetSfs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(buildSuccessResponse(200, "success", map[string]interface{}{
			"name":                "test-sfs",
			"status":              "active",
			"vpc_id":              "vpc-123",
			"efs_disk_size":       "100GB",
			"disk_iops":           3000,
			"plan_name":           "basic",
			"isEncryptionEnabled": false,
			"private_endpoint":    "10.0.0.1",
			"is_backup_enabled":   true,
		})))
	}))
	defer server.Close()
	client, err := NewClient("test-api-key", "test-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	result, _, err := client.Sfs.GetSfs(context.Background(), "12345")
	assertNoError(t, err)
	if result.Status != "active" {
		t.Errorf("Expected Status active, got %s", result.Status)
	}
	if result.VPCID != "vpc-123" {
		t.Errorf("Expected VPCID vpc-123, got %s", result.VPCID)
	}
	if result.DiskSize != "100GB" {
		t.Errorf("Expected DiskSize 100GB, got %s", result.DiskSize)
	}
	if result.DiskIOPS != 3000 {
		t.Errorf("Expected DiskIOPS 3000, got %d", result.DiskIOPS)
	}
	if result.PlanName != "basic" {
		t.Errorf("Expected PlanName basic, got %s", result.PlanName)
	}
	if result.IsEncryptionEnabled {
		t.Error("Expected IsEncryptionEnabled false, got true")
	}
	if result.PrivateIPAddress != "10.0.0.1" {
		t.Errorf("Expected PrivateIPAddress 10.0.0.1, got %s", result.PrivateIPAddress)
	}
	if !result.IsBackupEnabled {
		t.Error("Expected IsBackupEnabled true, got false")
	}
}
func TestGetSfs_WithFloatIOPS(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(buildSuccessResponse(200, "success", map[string]interface{}{
			"name":          "test-sfs",
			"status":        "active",
			"vpc_id":        "vpc-123",
			"efs_disk_size": "100GB",
			"disk_iops":     3000.0,
			"plan_name":     "basic",
		})))
	}))
	defer server.Close()
	client, err := NewClient("test-api-key", "test-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	result, _, err := client.Sfs.GetSfs(context.Background(), "12345")
	assertNoError(t, err)
	if result.DiskIOPS != 3000 {
		t.Errorf("Expected DiskIOPS 3000, got %d", result.DiskIOPS)
	}
}

func TestGetSfs_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{
			"code": 404,
			"message": "not found"
		}`)
	}))
	defer server.Close()
	client, err := NewClient("test-api-key", "test-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	result, _, err := client.Sfs.GetSfs(context.Background(), "nonexistent")
	if result != nil {
		t.Errorf("Expected nil result for 404, got: %v", result)
	}
	assertError(t, err, "")
}

func TestDeleteSfs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
			"code": 200,
			"message": "success"
		}`)
	}))
	defer server.Close()
	client, err := NewClient("test-api-key", "test-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	_, err = client.Sfs.DeleteSfs(context.Background(), "12345")
	assertNoError(t, err)
}

func TestListSfss(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(buildListResponse([]interface{}{
			map[string]interface{}{
				"id":                1,
				"name":              "sfs-1",
				"efs_disk_size":     "100GB",
				"status":            "active",
				"private_endpoint":  "10.0.0.1",
				"iops":              3000,
				"is_backup_enabled": false,
				"plan_name":         "basic",
			},
			map[string]interface{}{
				"id":                2,
				"name":              "sfs-2",
				"efs_disk_size":     "200GB",
				"status":            "creating",
				"private_endpoint":  "10.0.0.2",
				"iops":              5000,
				"is_backup_enabled": true,
				"plan_name":         "premium",
			},
		})))
	}))
	defer server.Close()
	client, err := NewClient("test-api-key", "test-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	result, _, err := client.Sfs.ListSfss(context.Background())
	assertNoError(t, err)
	if len(result) != 2 {
		t.Errorf("Expected 2 SFS instances, got %d", len(result))
	}
	if result[0].Name != "sfs-1" {
		t.Errorf("Expected first SFS name sfs-1, got %s", result[0].Name)
	}
	if result[1].Name != "sfs-2" {
		t.Errorf("Expected second SFS name sfs-2, got %s", result[1].Name)
	}
}

func TestListSfss_Empty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(buildListResponse([]interface{}{})))
	}))
	defer server.Close()
	client, err := NewClient("test-api-key", "test-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	result, _, err := client.Sfs.ListSfss(context.Background())
	assertNoError(t, err)
	if len(result) != 0 {
		t.Errorf("Expected 0 SFS instances, got %d", len(result))
	}
}

// Edge case tests for better coverage
func TestCreateSfs_NilRequest(t *testing.T) {
	client, err := NewClient("test-api-key", "test-token", "test-project", "test-location", WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	_, _, err = client.Sfs.CreateSfs(context.Background(), nil)
	assertError(t, err, "")
}

func TestCreateSfs_EmptyName(t *testing.T) {
	client, err := NewClient("test-api-key", "test-token", "test-project", "test-location", WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	createReq := &SfsCreateRequest{
		Name:     "",
		Plan:     "basic",
		VPCID:    "vpc-123",
		DiskSize: 100,
		DiskIOPS: 3000,
	}
	_, _, err = client.Sfs.CreateSfs(context.Background(), createReq)
	assertError(t, err, "")
}

func TestCreateSfs_EmptyPlan(t *testing.T) {
	client, err := NewClient("test-api-key", "test-token", "test-project", "test-location", WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	createReq := &SfsCreateRequest{
		Name:     "test-sfs",
		Plan:     "",
		VPCID:    "vpc-123",
		DiskSize: 100,
		DiskIOPS: 3000,
	}
	_, _, err = client.Sfs.CreateSfs(context.Background(), createReq)
	assertError(t, err, "")
}

func TestCreateSfs_EmptyVPCID(t *testing.T) {
	client, err := NewClient("test-api-key", "test-token", "test-project", "test-location", WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	createReq := &SfsCreateRequest{
		Name:     "test-sfs",
		Plan:     "basic",
		VPCID:    "",
		DiskSize: 100,
		DiskIOPS: 3000,
	}
	_, _, err = client.Sfs.CreateSfs(context.Background(), createReq)
	assertError(t, err, "")
}

func TestCreateSfs_InvalidDiskSize(t *testing.T) {
	client, err := NewClient("test-api-key", "test-token", "test-project", "test-location", WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	createReq := &SfsCreateRequest{
		Name:     "test-sfs",
		Plan:     "basic",
		VPCID:    "vpc-123",
		DiskSize: 0,
		DiskIOPS: 3000,
	}
	_, _, err = client.Sfs.CreateSfs(context.Background(), createReq)
	assertError(t, err, "")
}

func TestCreateSfs_InvalidDiskIOPS(t *testing.T) {
	client, err := NewClient("test-api-key", "test-token", "test-project", "test-location", WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	createReq := &SfsCreateRequest{
		Name:     "test-sfs",
		Plan:     "basic",
		VPCID:    "vpc-123",
		DiskSize: 100,
		DiskIOPS: 0,
	}
	_, _, err = client.Sfs.CreateSfs(context.Background(), createReq)
	assertError(t, err, "")
}

func TestGetSfs_EmptyID(t *testing.T) {
	client, err := NewClient("test-api-key", "test-token", "test-project", "test-location", WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	_, _, err = client.Sfs.GetSfs(context.Background(), "")
	assertError(t, err, "")
}

func TestDeleteSfs_EmptyID(t *testing.T) {
	client, err := NewClient("test-api-key", "test-token", "test-project", "test-location", WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	_, err = client.Sfs.DeleteSfs(context.Background(), "")
	assertError(t, err, "")
}

// Error response tests
func TestCreateSfs_ErrorResponse(t *testing.T) {
	server := newErrorServer(t, http.StatusInternalServerError, "Internal server error")
	defer server.Close()
	client, err := NewClient("test-api-key", "test-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	createReq := &SfsCreateRequest{
		Name:     "test-sfs",
		Plan:     "basic",
		VPCID:    "vpc-123",
		DiskSize: 100,
		DiskIOPS: 3000,
	}
	_, _, err = client.Sfs.CreateSfs(context.Background(), createReq)
	assertError(t, err, "")
}
