package goe2e

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestListPersistentVolumes tests the happy path for listing persistent volumes
func TestListPersistentVolumes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/persistent_volume/cluster-123/")
		assertQueryParam(t, r, "apikey", "test-api-key")
		assertQueryParam(t, r, "project_id", "test-project")
		assertQueryParam(t, r, "location", "test-location")

		writeJSON(w, http.StatusOK, buildListResponse([]interface{}{
			map[string]interface{}{
				"id":         "pv-123",
				"name":       "test-pv-1",
				"pv_size":    100,
				"cluster_id": "cluster-123",
				"status":     "active",
				"is_dynamic": false,
				"created_at": "2024-01-01T00:00:00Z",
			},
			map[string]interface{}{
				"id":         "pv-456",
				"name":       "test-pv-2",
				"pv_size":    200,
				"cluster_id": "cluster-123",
				"status":     "provisioning",
				"is_dynamic": true,
				"created_at": "2024-01-02T00:00:00Z",
			},
		}))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	pvs, _, err := client.Kubernetes.ListPersistentVolumes(context.Background(), "cluster-123")
	assertNoError(t, err)
	if len(pvs) != 2 {
		t.Errorf("Expected 2 persistent volumes, got %d", len(pvs))
	}

	// Verify first PV
	if pvs[0].ID != "pv-123" {
		t.Errorf("Expected ID pv-123, got %s", pvs[0].ID)
	}
	if pvs[0].Name != "test-pv-1" {
		t.Errorf("Expected Name test-pv-1, got %s", pvs[0].Name)
	}
	if pvs[0].PVSize != 100 {
		t.Errorf("Expected PVSize 100, got %d", pvs[0].PVSize)
	}
	if pvs[0].IsDynamic != false {
		t.Errorf("Expected IsDynamic false, got %v", pvs[0].IsDynamic)
	}

	// Verify second PV
	if pvs[1].ID != "pv-456" {
		t.Errorf("Expected ID pv-456, got %s", pvs[1].ID)
	}
	if pvs[1].IsDynamic != true {
		t.Errorf("Expected IsDynamic true, got %v", pvs[1].IsDynamic)
	}
}

// TestListPersistentVolumes_EmptyClusterID tests validation for empty cluster ID
func TestListPersistentVolumes_EmptyClusterID(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	assertNoError(t, err)

	_, _, err = client.Kubernetes.ListPersistentVolumes(context.Background(), "")
	assertError(t, err, "")
	assertErrorType(t, err, &ArgError{})
}

// TestListPersistentVolumes_EmptyList tests handling of empty PV list
func TestListPersistentVolumes_EmptyList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodGet)
		writeJSON(w, http.StatusOK, buildListResponse([]interface{}{}))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	pvs, _, err := client.Kubernetes.ListPersistentVolumes(context.Background(), "cluster-123")
	assertNoError(t, err)
	if len(pvs) != 0 {
		t.Errorf("Expected 0 persistent volumes, got %d", len(pvs))
	}
}

// TestListPersistentVolumes_IntAndFloat64Types tests handling of both int and float64 for pv_size
func TestListPersistentVolumes_IntAndFloat64Types(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodGet)
		// Mix float64 and int types
		writeJSON(w, http.StatusOK, buildListResponse([]interface{}{
			map[string]interface{}{
				"id":         "pv-123",
				"name":       "test-pv-float",
				"pv_size":    100.0,
				"cluster_id": "cluster-123",
				"status":     "active",
				"is_dynamic": false,
				"created_at": "2024-01-01T00:00:00Z",
			},
		}))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	pvs, _, err := client.Kubernetes.ListPersistentVolumes(context.Background(), "cluster-123")
	assertNoError(t, err)
	if len(pvs) != 1 {
		t.Fatalf("Expected 1 persistent volume, got %d", len(pvs))
	}
	if pvs[0].PVSize != 100 {
		t.Errorf("Expected PVSize 100, got %d", pvs[0].PVSize)
	}
}

// TestListPersistentVolumes_APIError tests handling of API errors
func TestListPersistentVolumes_APIError(t *testing.T) {
	server := newErrorServer(t, http.StatusInternalServerError, "Internal server error")
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}), noRetryOpt())
	assertNoError(t, err)

	_, _, err = client.Kubernetes.ListPersistentVolumes(context.Background(), "cluster-123")
	assertError(t, err, "")
}

// TestCreatePersistentVolume tests the happy path for creating a persistent volume
func TestCreatePersistentVolume(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPost)
		testURLPath(t, r, "/persistent_volume/cluster-123/")
		assertQueryParam(t, r, "apikey", "test-api-key")
		assertQueryParam(t, r, "project_id", "test-project")
		assertQueryParam(t, r, "location", "test-location")

		writeJSON(w, http.StatusCreated, buildSuccessResponse(201, "Persistent volume created successfully", map[string]interface{}{
			"id":         "pv-new-123",
			"name":       "my-pv",
			"pv_size":    150,
			"cluster_id": "cluster-123",
			"status":     "provisioning",
			"is_dynamic": false,
			"created_at": "2024-01-03T00:00:00Z",
		}))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	createReq := &CreatePersistentVolumeRequest{
		Name:      "my-pv",
		PVSize:    150,
		IsDynamic: false,
	}

	pv, _, err := client.Kubernetes.CreatePersistentVolume(context.Background(), "cluster-123", createReq)
	assertNoError(t, err)
	assertNotNil(t, pv, "Expected persistent volume to be returned")
	if pv.ID != "pv-new-123" {
		t.Errorf("Expected ID pv-new-123, got %s", pv.ID)
	}
	if pv.Name != "my-pv" {
		t.Errorf("Expected Name my-pv, got %s", pv.Name)
	}
	if pv.PVSize != 150 {
		t.Errorf("Expected PVSize 150, got %d", pv.PVSize)
	}
	if pv.ClusterID != "cluster-123" {
		t.Errorf("Expected ClusterID cluster-123, got %s", pv.ClusterID)
	}
	if pv.Status != "provisioning" {
		t.Errorf("Expected Status provisioning, got %s", pv.Status)
	}
	if pv.IsDynamic != false {
		t.Errorf("Expected IsDynamic false, got %v", pv.IsDynamic)
	}
}

// TestCreatePersistentVolume_DynamicVolume tests creating a dynamic volume
func TestCreatePersistentVolume_DynamicVolume(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPost)
		writeJSON(w, http.StatusCreated, buildSuccessResponse(201, "Persistent volume created successfully", map[string]interface{}{
			"id":         "pv-dynamic-123",
			"name":       "my-dynamic-pv",
			"pv_size":    200,
			"cluster_id": "cluster-123",
			"status":     "active",
			"is_dynamic": true,
			"created_at": "2024-01-03T00:00:00Z",
		}))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	createReq := &CreatePersistentVolumeRequest{
		Name:      "my-dynamic-pv",
		PVSize:    200,
		IsDynamic: true,
	}

	pv, _, err := client.Kubernetes.CreatePersistentVolume(context.Background(), "cluster-123", createReq)
	assertNoError(t, err)
	if pv.IsDynamic != true {
		t.Errorf("Expected IsDynamic true, got %v", pv.IsDynamic)
	}
}

// TestCreatePersistentVolume_EmptyClusterID tests validation for empty cluster ID
func TestCreatePersistentVolume_EmptyClusterID(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	assertNoError(t, err)

	createReq := &CreatePersistentVolumeRequest{
		Name:   "my-pv",
		PVSize: 100,
	}

	_, _, err = client.Kubernetes.CreatePersistentVolume(context.Background(), "", createReq)
	assertError(t, err, "")
	assertErrorType(t, err, &ArgError{})
}

// TestCreatePersistentVolume_NilRequest tests validation for nil request
func TestCreatePersistentVolume_NilRequest(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	assertNoError(t, err)

	_, _, err = client.Kubernetes.CreatePersistentVolume(context.Background(), "cluster-123", nil)
	assertError(t, err, "")
	assertErrorType(t, err, &ArgError{})
}

// TestCreatePersistentVolume_EmptyName tests validation for empty name
func TestCreatePersistentVolume_EmptyName(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	assertNoError(t, err)

	createReq := &CreatePersistentVolumeRequest{
		Name:   "",
		PVSize: 100,
	}

	_, _, err = client.Kubernetes.CreatePersistentVolume(context.Background(), "cluster-123", createReq)
	assertError(t, err, "")
	assertErrorType(t, err, &ArgError{})
}

// TestCreatePersistentVolume_ZeroSize tests validation for zero size
func TestCreatePersistentVolume_ZeroSize(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	assertNoError(t, err)

	createReq := &CreatePersistentVolumeRequest{
		Name:   "my-pv",
		PVSize: 0,
	}

	_, _, err = client.Kubernetes.CreatePersistentVolume(context.Background(), "cluster-123", createReq)
	assertError(t, err, "")
	assertErrorType(t, err, &ArgError{})
}

// TestCreatePersistentVolume_NegativeSize tests validation for negative size
func TestCreatePersistentVolume_NegativeSize(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	assertNoError(t, err)

	createReq := &CreatePersistentVolumeRequest{
		Name:   "my-pv",
		PVSize: -100,
	}

	_, _, err = client.Kubernetes.CreatePersistentVolume(context.Background(), "cluster-123", createReq)
	assertError(t, err, "")
	assertErrorType(t, err, &ArgError{})
}

// TestGetPersistentVolume tests the happy path for getting a persistent volume
func TestGetPersistentVolume(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/persistent_volume/cluster-123/pv-456/")
		assertQueryParam(t, r, "apikey", "test-api-key")
		assertQueryParam(t, r, "project_id", "test-project")
		assertQueryParam(t, r, "location", "test-location")

		writeJSON(w, http.StatusOK, buildSuccessResponse(200, "success", map[string]interface{}{
			"id":         "pv-456",
			"name":       "specific-pv",
			"pv_size":    250,
			"cluster_id": "cluster-123",
			"status":     "active",
			"is_dynamic": true,
			"created_at": "2024-01-04T00:00:00Z",
		}))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	pv, _, err := client.Kubernetes.GetPersistentVolume(context.Background(), "cluster-123", "pv-456")
	assertNoError(t, err)
	assertNotNil(t, pv, "Expected persistent volume to be returned")
	if pv.ID != "pv-456" {
		t.Errorf("Expected ID pv-456, got %s", pv.ID)
	}
	if pv.Name != "specific-pv" {
		t.Errorf("Expected Name specific-pv, got %s", pv.Name)
	}
	if pv.PVSize != 250 {
		t.Errorf("Expected PVSize 250, got %d", pv.PVSize)
	}
	if pv.Status != "active" {
		t.Errorf("Expected Status active, got %s", pv.Status)
	}
}

// TestGetPersistentVolume_NotFound tests handling of 404 response
func TestGetPersistentVolume_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodGet)
		writeJSON(w, http.StatusNotFound, buildErrorResponse(404, "Persistent volume not found", nil))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	pv, resp, err := client.Kubernetes.GetPersistentVolume(context.Background(), "cluster-123", "pv-nonexistent")
	// Should return nil PV with no error for 404
	assertNoError(t, err)
	if pv != nil {
		t.Errorf("Expected nil PV for 404, got: %v", pv)
	}
	assertNotNil(t, resp, "Expected response to be returned")
	assertStatus(t, resp, http.StatusNotFound)
}

// TestGetPersistentVolume_EmptyClusterID tests validation for empty cluster ID
func TestGetPersistentVolume_EmptyClusterID(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	assertNoError(t, err)

	_, _, err = client.Kubernetes.GetPersistentVolume(context.Background(), "", "pv-123")
	assertError(t, err, "")
	assertErrorType(t, err, &ArgError{})
}

// TestGetPersistentVolume_EmptyPVID tests validation for empty PV ID
func TestGetPersistentVolume_EmptyPVID(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	assertNoError(t, err)

	_, _, err = client.Kubernetes.GetPersistentVolume(context.Background(), "cluster-123", "")
	assertError(t, err, "")
	assertErrorType(t, err, &ArgError{})
}

// TestGetPersistentVolume_ServerError tests handling of server errors
func TestGetPersistentVolume_ServerError(t *testing.T) {
	server := newErrorServer(t, http.StatusInternalServerError, "Internal server error")
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}), noRetryOpt())
	assertNoError(t, err)

	_, _, err = client.Kubernetes.GetPersistentVolume(context.Background(), "cluster-123", "pv-456")
	assertError(t, err, "")
}

// TestDeletePersistentVolume tests the happy path for deleting a persistent volume
func TestDeletePersistentVolume(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodDelete)
		testURLPath(t, r, "/persistent_volume/cluster-123/pv-delete/")
		assertQueryParam(t, r, "apikey", "test-api-key")
		assertQueryParam(t, r, "project_id", "test-project")
		assertQueryParam(t, r, "location", "test-location")

		writeJSON(w, http.StatusOK, buildErrorResponse(200, "Persistent volume deleted successfully", nil))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, err = client.Kubernetes.DeletePersistentVolume(context.Background(), "cluster-123", "pv-delete")
	assertNoError(t, err)
}

// TestDeletePersistentVolume_EmptyClusterID tests validation for empty cluster ID
func TestDeletePersistentVolume_EmptyClusterID(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	assertNoError(t, err)

	_, err = client.Kubernetes.DeletePersistentVolume(context.Background(), "", "pv-123")
	assertError(t, err, "")
	assertErrorType(t, err, &ArgError{})
}

// TestDeletePersistentVolume_EmptyPVID tests validation for empty PV ID
func TestDeletePersistentVolume_EmptyPVID(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	assertNoError(t, err)

	_, err = client.Kubernetes.DeletePersistentVolume(context.Background(), "cluster-123", "")
	assertError(t, err, "")
	assertErrorType(t, err, &ArgError{})
}

// TestDeletePersistentVolume_NotFound tests handling of deleting non-existent PV
func TestDeletePersistentVolume_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodDelete)
		writeJSON(w, http.StatusNotFound, buildErrorResponse(404, "Persistent volume not found", nil))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, err = client.Kubernetes.DeletePersistentVolume(context.Background(), "cluster-123", "pv-nonexistent")
	assertError(t, err, "")
}

// TestDeletePersistentVolume_InUse tests handling of deleting in-use PV
func TestDeletePersistentVolume_InUse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodDelete)
		writeJSON(w, http.StatusConflict, buildErrorResponse(409, "Persistent volume is in use and cannot be deleted", nil))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, err = client.Kubernetes.DeletePersistentVolume(context.Background(), "cluster-123", "pv-in-use")
	assertError(t, err, "")
}

// TestPersistentVolume_MissingFieldsHandling tests handling of missing optional fields
func TestPersistentVolume_MissingFieldsHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodGet)
		// Response with minimal fields
		writeJSON(w, http.StatusOK, buildListResponse([]interface{}{
			map[string]interface{}{
				"id":   "pv-minimal",
				"name": "minimal-pv",
			},
		}))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	pvs, _, err := client.Kubernetes.ListPersistentVolumes(context.Background(), "cluster-123")
	assertNoError(t, err)
	if len(pvs) != 1 {
		t.Fatalf("Expected 1 persistent volume, got %d", len(pvs))
	}

	// Verify minimal fields are set, others are zero values
	if pvs[0].ID != "pv-minimal" {
		t.Errorf("Expected ID pv-minimal, got %s", pvs[0].ID)
	}
	if pvs[0].Name != "minimal-pv" {
		t.Errorf("Expected Name minimal-pv, got %s", pvs[0].Name)
	}
	if pvs[0].PVSize != 0 {
		t.Errorf("Expected PVSize 0 for missing field, got %d", pvs[0].PVSize)
	}
	if pvs[0].Status != "" {
		t.Errorf("Expected empty Status for missing field, got %s", pvs[0].Status)
	}
	if pvs[0].IsDynamic != false {
		t.Errorf("Expected IsDynamic false for missing field, got %v", pvs[0].IsDynamic)
	}
}

// TestPersistentVolume_LargePVSize tests handling of large PV sizes
func TestPersistentVolume_LargePVSize(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPost)
		writeJSON(w, http.StatusCreated, buildSuccessResponse(201, "success", map[string]interface{}{
			"id":         "pv-large",
			"name":       "large-pv",
			"pv_size":    10000,
			"cluster_id": "cluster-123",
			"status":     "active",
			"is_dynamic": false,
			"created_at": "2024-01-01T00:00:00Z",
		}))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	createReq := &CreatePersistentVolumeRequest{
		Name:   "large-pv",
		PVSize: 10000,
	}

	pv, _, err := client.Kubernetes.CreatePersistentVolume(context.Background(), "cluster-123", createReq)
	assertNoError(t, err)
	if pv.PVSize != 10000 {
		t.Errorf("Expected PVSize 10000, got %d", pv.PVSize)
	}
}
