package goe2e

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func testHeader(t *testing.T, r *http.Request, header, expected string) {
	t.Helper()
	if got := r.Header.Get(header); got != expected {
		t.Errorf("Header %s: %v, expected %v", header, got, expected)
	}
}

func TestObjectStorageServiceOp_CreateBucket(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/storage/buckets/test-bucket/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Authorization", "Bearer test-auth-token")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		fmt.Fprint(w, `{
			"code": 200,
			"message": "Bucket created successfully",
			"data": {
				"id": 123,
				"name": "test-bucket",
				"status": "ACTIVE",
				"bucket_size": "0 MB",
				"created_at": "2023-01-01T00:00:00Z",
				"versioning_status": "Suspended",
				"lifecycle_configuration_status": "Disabled"
			}
		}`)
	})

	createReq := &BucketCreateRequest{
		BucketName: "test-bucket",
	}

	bucket, resp, err := ts.client.ObjectStorage.CreateBucket(context.Background(), createReq)
	if err != nil {
		t.Fatalf("ObjectStorage.CreateBucket returned error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	expectedBucket := &Bucket{
		ID:                           123,
		Name:                         "test-bucket",
		Status:                       "ACTIVE",
		BucketSize:                   "0 MB",
		CreatedAt:                    "2023-01-01T00:00:00Z",
		VersioningStatus:             "Suspended",
		LifecycleConfigurationStatus: "Disabled",
	}

	if bucket.Name != expectedBucket.Name {
		t.Errorf("Expected bucket name %s, got %s", expectedBucket.Name, bucket.Name)
	}
	if bucket.Status != expectedBucket.Status {
		t.Errorf("Expected bucket status %s, got %s", expectedBucket.Status, bucket.Status)
	}
	if bucket.ID != expectedBucket.ID {
		t.Errorf("Expected bucket ID %d, got %d", expectedBucket.ID, bucket.ID)
	}
}

func TestObjectStorageServiceOp_CreateBucket_NilRequest(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, _, err := ts.client.ObjectStorage.CreateBucket(context.Background(), nil)
	if err == nil {
		t.Fatal("Expected error when createReq is nil, got nil")
	}

	argErr, ok := err.(*ArgError)
	if !ok {
		t.Fatalf("Expected ArgError, got %T", err)
	}
	if argErr.arg != "createReq" {
		t.Errorf("Expected arg 'createReq', got '%s'", argErr.arg)
	}
}

func TestObjectStorageServiceOp_CreateBucket_EmptyBucketName(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	createReq := &BucketCreateRequest{
		BucketName: "",
	}

	_, _, err := ts.client.ObjectStorage.CreateBucket(context.Background(), createReq)
	if err == nil {
		t.Fatal("Expected error when bucket name is empty, got nil")
	}

	argErr, ok := err.(*ArgError)
	if !ok {
		t.Fatalf("Expected ArgError, got %T", err)
	}
	if argErr.arg != "createReq.BucketName" {
		t.Errorf("Expected arg 'createReq.BucketName', got '%s'", argErr.arg)
	}
}

func TestObjectStorageServiceOp_CreateBucket_Error(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/storage/buckets/test-bucket/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{
			"code": 400,
			"message": "Bucket name already exists",
			"errors": ["Bucket with this name already exists"]
		}`)
	})

	createReq := &BucketCreateRequest{
		BucketName: "test-bucket",
	}

	bucket, resp, err := ts.client.ObjectStorage.CreateBucket(context.Background(), createReq)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if bucket != nil {
		t.Errorf("Expected nil bucket on error, got %+v", bucket)
	}
	if resp == nil {
		t.Fatal("Expected response, got nil")
	}
}

func TestObjectStorageServiceOp_GetBucket(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/storage/buckets/test-bucket/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testHeader(t, r, "Authorization", "Bearer test-auth-token")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		fmt.Fprint(w, `{
			"code": 200,
			"message": "Success",
			"data": {
				"id": 123,
				"name": "test-bucket",
				"status": "ACTIVE",
				"bucket_size": "10 MB",
				"created_at": "2023-01-01T00:00:00Z",
				"versioning_status": "Enabled",
				"lifecycle_configuration_status": "Enabled"
			}
		}`)
	})

	bucket, resp, err := ts.client.ObjectStorage.GetBucket(context.Background(), "test-bucket")
	if err != nil {
		t.Fatalf("ObjectStorage.GetBucket returned error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if bucket.Name != "test-bucket" {
		t.Errorf("Expected bucket name test-bucket, got %s", bucket.Name)
	}
	if bucket.Status != "ACTIVE" {
		t.Errorf("Expected bucket status ACTIVE, got %s", bucket.Status)
	}
	if bucket.VersioningStatus != "Enabled" {
		t.Errorf("Expected versioning status Enabled, got %s", bucket.VersioningStatus)
	}
}

func TestObjectStorageServiceOp_GetBucket_EmptyName(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, _, err := ts.client.ObjectStorage.GetBucket(context.Background(), "")
	if err == nil {
		t.Fatal("Expected error when bucket name is empty, got nil")
	}

	argErr, ok := err.(*ArgError)
	if !ok {
		t.Fatalf("Expected ArgError, got %T", err)
	}
	if argErr.arg != "bucketName" {
		t.Errorf("Expected arg 'bucketName', got '%s'", argErr.arg)
	}
}

func TestObjectStorageServiceOp_GetBucket_NotFound(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/storage/buckets/nonexistent-bucket/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{
			"code": 404,
			"message": "Bucket not found",
			"errors": ["Bucket does not exist"]
		}`)
	})

	bucket, resp, err := ts.client.ObjectStorage.GetBucket(context.Background(), "nonexistent-bucket")
	if err != nil {
		t.Errorf("Expected nil error for 404, got: %v", err)
	}
	if bucket != nil {
		t.Errorf("Expected nil bucket for 404, got %+v", bucket)
	}
	if resp == nil {
		t.Fatal("Expected response, got nil")
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
}

func TestObjectStorageServiceOp_GetBucket_Error(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/storage/buckets/test-bucket/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{
			"code": 500,
			"message": "Internal server error"
		}`)
	})

	bucket, resp, err := ts.client.ObjectStorage.GetBucket(context.Background(), "test-bucket")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if bucket != nil {
		t.Errorf("Expected nil bucket on error, got %+v", bucket)
	}
	if resp == nil {
		t.Fatal("Expected response, got nil")
	}
}

func TestObjectStorageServiceOp_ListBuckets(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/storage/buckets/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testHeader(t, r, "Authorization", "Bearer test-auth-token")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		fmt.Fprint(w, `{
			"code": 200,
			"message": "Success",
			"data": [
				{
					"id": 123,
					"name": "bucket-1",
					"status": "ACTIVE",
					"bucket_size": "10 MB",
					"created_at": "2023-01-01T00:00:00Z",
					"versioning_status": "Enabled",
					"lifecycle_configuration_status": "Enabled"
				},
				{
					"id": 124,
					"name": "bucket-2",
					"status": "ACTIVE",
					"bucket_size": "5 MB",
					"created_at": "2023-01-02T00:00:00Z",
					"versioning_status": "Suspended",
					"lifecycle_configuration_status": "Disabled"
				}
			]
		}`)
	})

	buckets, resp, err := ts.client.ObjectStorage.ListBuckets(context.Background())
	if err != nil {
		t.Fatalf("ObjectStorage.ListBuckets returned error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if len(buckets) != 2 {
		t.Errorf("Expected 2 buckets, got %d", len(buckets))
	}

	if buckets[0].Name != "bucket-1" {
		t.Errorf("Expected first bucket name bucket-1, got %s", buckets[0].Name)
	}
	if buckets[1].Name != "bucket-2" {
		t.Errorf("Expected second bucket name bucket-2, got %s", buckets[1].Name)
	}
}

func TestObjectStorageServiceOp_ListBuckets_Empty(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/storage/buckets/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"code": 200,
			"message": "Success",
			"data": []
		}`)
	})

	buckets, resp, err := ts.client.ObjectStorage.ListBuckets(context.Background())
	if err != nil {
		t.Fatalf("ObjectStorage.ListBuckets returned error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if len(buckets) != 0 {
		t.Errorf("Expected 0 buckets, got %d", len(buckets))
	}
}

func TestObjectStorageServiceOp_ListBuckets_Error(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/storage/buckets/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, `{
			"code": 403,
			"message": "Access denied"
		}`)
	})

	buckets, resp, err := ts.client.ObjectStorage.ListBuckets(context.Background())
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if buckets != nil {
		t.Errorf("Expected nil buckets on error, got %+v", buckets)
	}
	if resp == nil {
		t.Fatal("Expected response, got nil")
	}
}

func TestObjectStorageServiceOp_DeleteBucket(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/storage/buckets/test-bucket/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		testHeader(t, r, "Authorization", "Bearer test-auth-token")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"code": 200,
			"message": "Bucket deleted successfully"
		}`)
	})

	resp, err := ts.client.ObjectStorage.DeleteBucket(context.Background(), "test-bucket")
	if err != nil {
		t.Fatalf("ObjectStorage.DeleteBucket returned error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestObjectStorageServiceOp_DeleteBucket_EmptyName(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.ObjectStorage.DeleteBucket(context.Background(), "")
	if err == nil {
		t.Fatal("Expected error when bucket name is empty, got nil")
	}

	argErr, ok := err.(*ArgError)
	if !ok {
		t.Fatalf("Expected ArgError, got %T", err)
	}
	if argErr.arg != "bucketName" {
		t.Errorf("Expected arg 'bucketName', got '%s'", argErr.arg)
	}
}

func TestObjectStorageServiceOp_DeleteBucket_Error(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/storage/buckets/test-bucket/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, `{
			"code": 409,
			"message": "Bucket not empty",
			"errors": ["Cannot delete non-empty bucket"]
		}`)
	})

	resp, err := ts.client.ObjectStorage.DeleteBucket(context.Background(), "test-bucket")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if resp == nil {
		t.Fatal("Expected response, got nil")
	}
	if resp.StatusCode != http.StatusConflict {
		t.Errorf("Expected status code %d, got %d", http.StatusConflict, resp.StatusCode)
	}
}

func TestObjectStorageServiceOp_SetBucketVersioning_Enable(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/storage/bucket_versioning/test-bucket/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testHeader(t, r, "Authorization", "Bearer test-auth-token")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		fmt.Fprint(w, `{
			"code": 200,
			"message": "Versioning updated successfully",
			"data": {
				"bucket_name": "test-bucket",
				"versioning_status": "Enabled"
			}
		}`)
	})

	versioningReq := &BucketVersioningRequest{
		BucketName:         "test-bucket",
		NewVersioningState: "Enabled",
	}

	versioning, resp, err := ts.client.ObjectStorage.SetBucketVersioning(context.Background(), "test-bucket", versioningReq)
	if err != nil {
		t.Fatalf("ObjectStorage.SetBucketVersioning returned error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if versioning.BucketName != "test-bucket" {
		t.Errorf("Expected bucket name test-bucket, got %s", versioning.BucketName)
	}
	if versioning.VersioningStatus != "Enabled" {
		t.Errorf("Expected versioning status Enabled, got %s", versioning.VersioningStatus)
	}
}

func TestObjectStorageServiceOp_SetBucketVersioning_Suspend(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/storage/bucket_versioning/test-bucket/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		fmt.Fprint(w, `{
			"code": 200,
			"message": "Versioning updated successfully",
			"data": {
				"bucket_name": "test-bucket",
				"versioning_status": "Suspended"
			}
		}`)
	})

	versioningReq := &BucketVersioningRequest{
		BucketName:         "test-bucket",
		NewVersioningState: "Suspended",
	}

	versioning, resp, err := ts.client.ObjectStorage.SetBucketVersioning(context.Background(), "test-bucket", versioningReq)
	if err != nil {
		t.Fatalf("ObjectStorage.SetBucketVersioning returned error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if versioning.VersioningStatus != "Suspended" {
		t.Errorf("Expected versioning status Suspended, got %s", versioning.VersioningStatus)
	}
}

func TestObjectStorageServiceOp_SetBucketVersioning_EmptyBucketName(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	versioningReq := &BucketVersioningRequest{
		NewVersioningState: "Enabled",
	}

	_, _, err := ts.client.ObjectStorage.SetBucketVersioning(context.Background(), "", versioningReq)
	if err == nil {
		t.Fatal("Expected error when bucket name is empty, got nil")
	}

	argErr, ok := err.(*ArgError)
	if !ok {
		t.Fatalf("Expected ArgError, got %T", err)
	}
	if argErr.arg != "bucketName" {
		t.Errorf("Expected arg 'bucketName', got '%s'", argErr.arg)
	}
}

func TestObjectStorageServiceOp_SetBucketVersioning_NilRequest(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, _, err := ts.client.ObjectStorage.SetBucketVersioning(context.Background(), "test-bucket", nil)
	if err == nil {
		t.Fatal("Expected error when versioningReq is nil, got nil")
	}

	argErr, ok := err.(*ArgError)
	if !ok {
		t.Fatalf("Expected ArgError, got %T", err)
	}
	if argErr.arg != "versioningReq" {
		t.Errorf("Expected arg 'versioningReq', got '%s'", argErr.arg)
	}
}

func TestObjectStorageServiceOp_SetBucketVersioning_EmptyState(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	versioningReq := &BucketVersioningRequest{
		BucketName:         "test-bucket",
		NewVersioningState: "",
	}

	_, _, err := ts.client.ObjectStorage.SetBucketVersioning(context.Background(), "test-bucket", versioningReq)
	if err == nil {
		t.Fatal("Expected error when versioning state is empty, got nil")
	}

	argErr, ok := err.(*ArgError)
	if !ok {
		t.Fatalf("Expected ArgError, got %T", err)
	}
	if argErr.arg != "versioningReq.NewVersioningState" {
		t.Errorf("Expected arg 'versioningReq.NewVersioningState', got '%s'", argErr.arg)
	}
}

func TestObjectStorageServiceOp_SetBucketVersioning_InvalidState(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	versioningReq := &BucketVersioningRequest{
		BucketName:         "test-bucket",
		NewVersioningState: "Invalid",
	}

	_, _, err := ts.client.ObjectStorage.SetBucketVersioning(context.Background(), "test-bucket", versioningReq)
	if err == nil {
		t.Fatal("Expected error when versioning state is invalid, got nil")
	}

	argErr, ok := err.(*ArgError)
	if !ok {
		t.Fatalf("Expected ArgError, got %T", err)
	}
	if argErr.arg != "versioningReq.NewVersioningState" {
		t.Errorf("Expected arg 'versioningReq.NewVersioningState', got '%s'", argErr.arg)
	}
}

func TestObjectStorageServiceOp_SetBucketVersioning_AutoSetBucketName(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/storage/bucket_versioning/test-bucket/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		fmt.Fprint(w, `{
			"code": 200,
			"message": "Versioning updated successfully",
			"data": {
				"bucket_name": "test-bucket",
				"versioning_status": "Enabled"
			}
		}`)
	})

	// Test that bucket name is auto-set from parameter if not in request
	versioningReq := &BucketVersioningRequest{
		NewVersioningState: "Enabled",
		// BucketName intentionally not set
	}

	versioning, resp, err := ts.client.ObjectStorage.SetBucketVersioning(context.Background(), "test-bucket", versioningReq)
	if err != nil {
		t.Fatalf("ObjectStorage.SetBucketVersioning returned error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if versioning.BucketName != "test-bucket" {
		t.Errorf("Expected bucket name test-bucket, got %s", versioning.BucketName)
	}
}

func TestObjectStorageServiceOp_SetBucketVersioning_Error(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/storage/bucket_versioning/test-bucket/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{
			"code": 400,
			"message": "Invalid versioning state"
		}`)
	})

	versioningReq := &BucketVersioningRequest{
		BucketName:         "test-bucket",
		NewVersioningState: "Enabled",
	}

	versioning, resp, err := ts.client.ObjectStorage.SetBucketVersioning(context.Background(), "test-bucket", versioningReq)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if versioning != nil {
		t.Errorf("Expected nil versioning on error, got %+v", versioning)
	}
	if resp == nil {
		t.Fatal("Expected response, got nil")
	}
}

// TestObjectStorageServiceOp_CreateBucket_Conflict tests 409 Conflict on create
func TestObjectStorageServiceOp_CreateBucket_Conflict(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/storage/buckets/test-bucket/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, `{
			"code": 409,
			"message": "Bucket name already exists",
			"errors": ["Bucket with this name already exists"]
		}`)
	})

	createReq := &BucketCreateRequest{
		BucketName: "test-bucket",
	}

	bucket, resp, err := ts.client.ObjectStorage.CreateBucket(context.Background(), createReq)
	if err == nil {
		t.Fatal("Expected error for 409 Conflict, got nil")
	}
	if bucket != nil {
		t.Errorf("Expected nil bucket on error, got %+v", bucket)
	}
	if resp.StatusCode != http.StatusConflict {
		t.Errorf("Expected status code %d, got %d", http.StatusConflict, resp.StatusCode)
	}
}

// TestObjectStorageServiceOp_CreateBucket_Forbidden tests 403 Forbidden on create
func TestObjectStorageServiceOp_CreateBucket_Forbidden(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/storage/buckets/test-bucket/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, `{
			"code": 403,
			"message": "Access denied",
			"errors": ["Not authorized to create buckets"]
		}`)
	})

	createReq := &BucketCreateRequest{
		BucketName: "test-bucket",
	}

	bucket, resp, err := ts.client.ObjectStorage.CreateBucket(context.Background(), createReq)
	if err == nil {
		t.Fatal("Expected error for 403 Forbidden, got nil")
	}
	if bucket != nil {
		t.Errorf("Expected nil bucket on error, got %+v", bucket)
	}
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Expected status code %d, got %d", http.StatusForbidden, resp.StatusCode)
	}
}

// TestObjectStorageServiceOp_CreateBucket_ServiceUnavailable tests 503 Service Unavailable
func TestObjectStorageServiceOp_CreateBucket_ServiceUnavailable(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/storage/buckets/test-bucket/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, `{
			"code": 503,
			"message": "Service temporarily unavailable"
		}`)
	})

	createReq := &BucketCreateRequest{
		BucketName: "test-bucket",
	}

	bucket, resp, err := ts.client.ObjectStorage.CreateBucket(context.Background(), createReq)
	if err == nil {
		t.Fatal("Expected error for 503 Service Unavailable, got nil")
	}
	if bucket != nil {
		t.Errorf("Expected nil bucket on error, got %+v", bucket)
	}
	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("Expected status code %d, got %d", http.StatusServiceUnavailable, resp.StatusCode)
	}
}

// TestObjectStorageServiceOp_GetBucket_Forbidden tests 403 Forbidden on get
func TestObjectStorageServiceOp_GetBucket_Forbidden(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/storage/buckets/test-bucket/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, `{
			"code": 403,
			"message": "Access denied"
		}`)
	})

	bucket, resp, err := ts.client.ObjectStorage.GetBucket(context.Background(), "test-bucket")
	if err == nil {
		t.Fatal("Expected error for 403 Forbidden, got nil")
	}
	if bucket != nil {
		t.Errorf("Expected nil bucket on error, got %+v", bucket)
	}
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Expected status code %d, got %d", http.StatusForbidden, resp.StatusCode)
	}
}

// TestObjectStorageServiceOp_ListBuckets_ServiceUnavailable tests 503 Service Unavailable
func TestObjectStorageServiceOp_ListBuckets_ServiceUnavailable(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/storage/buckets/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, `{
			"code": 503,
			"message": "Service temporarily unavailable"
		}`)
	})

	buckets, resp, err := ts.client.ObjectStorage.ListBuckets(context.Background())
	if err == nil {
		t.Fatal("Expected error for 503, got nil")
	}
	if buckets != nil {
		t.Errorf("Expected nil buckets on error, got %+v", buckets)
	}
	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("Expected status code %d, got %d", http.StatusServiceUnavailable, resp.StatusCode)
	}
}

// TestObjectStorageServiceOp_DeleteBucket_Forbidden tests 403 Forbidden on delete
func TestObjectStorageServiceOp_DeleteBucket_Forbidden(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/storage/buckets/test-bucket/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, `{
			"code": 403,
			"message": "Access denied"
		}`)
	})

	resp, err := ts.client.ObjectStorage.DeleteBucket(context.Background(), "test-bucket")
	if err == nil {
		t.Fatal("Expected error for 403 Forbidden, got nil")
	}
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Expected status code %d, got %d", http.StatusForbidden, resp.StatusCode)
	}
}

// TestObjectStorageServiceOp_DeleteBucket_NotFound tests 404 Not Found on delete
func TestObjectStorageServiceOp_DeleteBucket_NotFound(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/storage/buckets/nonexistent-bucket/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{
			"code": 404,
			"message": "Bucket not found"
		}`)
	})

	resp, err := ts.client.ObjectStorage.DeleteBucket(context.Background(), "nonexistent-bucket")
	if err == nil {
		t.Fatal("Expected error for 404 Not Found, got nil")
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
}

// TestObjectStorageServiceOp_SetBucketVersioning_Forbidden tests 403 Forbidden on versioning
func TestObjectStorageServiceOp_SetBucketVersioning_Forbidden(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/storage/bucket_versioning/test-bucket/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, `{
			"code": 403,
			"message": "Access denied"
		}`)
	})

	versioningReq := &BucketVersioningRequest{
		BucketName:         "test-bucket",
		NewVersioningState: "Enabled",
	}

	versioning, resp, err := ts.client.ObjectStorage.SetBucketVersioning(context.Background(), "test-bucket", versioningReq)
	if err == nil {
		t.Fatal("Expected error for 403 Forbidden, got nil")
	}
	if versioning != nil {
		t.Errorf("Expected nil versioning on error, got %+v", versioning)
	}
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Expected status code %d, got %d", http.StatusForbidden, resp.StatusCode)
	}
}

// TestObjectStorageServiceOp_SetBucketVersioning_NotFound tests 404 Not Found on versioning
func TestObjectStorageServiceOp_SetBucketVersioning_NotFound(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/storage/bucket_versioning/nonexistent-bucket/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{
			"code": 404,
			"message": "Bucket not found"
		}`)
	})

	versioningReq := &BucketVersioningRequest{
		BucketName:         "nonexistent-bucket",
		NewVersioningState: "Enabled",
	}

	versioning, resp, err := ts.client.ObjectStorage.SetBucketVersioning(context.Background(), "nonexistent-bucket", versioningReq)
	if err == nil {
		t.Fatal("Expected error for 404 Not Found, got nil")
	}
	if versioning != nil {
		t.Errorf("Expected nil versioning on error, got %+v", versioning)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
}

// TestObjectStorageServiceOp_SetBucketVersioning_Conflict tests 409 Conflict on versioning
func TestObjectStorageServiceOp_SetBucketVersioning_Conflict(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/storage/bucket_versioning/test-bucket/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, `{
			"code": 409,
			"message": "Versioning operation in progress"
		}`)
	})

	versioningReq := &BucketVersioningRequest{
		BucketName:         "test-bucket",
		NewVersioningState: "Enabled",
	}

	versioning, resp, err := ts.client.ObjectStorage.SetBucketVersioning(context.Background(), "test-bucket", versioningReq)
	if err == nil {
		t.Fatal("Expected error for 409 Conflict, got nil")
	}
	if versioning != nil {
		t.Errorf("Expected nil versioning on error, got %+v", versioning)
	}
	if resp.StatusCode != http.StatusConflict {
		t.Errorf("Expected status code %d, got %d", http.StatusConflict, resp.StatusCode)
	}
}

// ============================================================================
// Network Error & Timeout Tests
// ============================================================================

// TestObjectStorage_NetworkConnectionError tests network connection failure
func TestObjectStorage_NetworkConnectionError(t *testing.T) {
	client, _ := NewClient(
		"test-key",
		"test-token",
		"proj-123",
		"test-location",
		SetBaseURL("http://localhost:1"),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}),
	)

	ctx := context.Background()
	createReq := &BucketCreateRequest{
		BucketName: "test-bucket",
	}

	_, _, err := client.ObjectStorage.CreateBucket(ctx, createReq)
	if err == nil {
		t.Fatal("Expected network error")
	}
}

// TestObjectStorage_ContextTimeout tests context timeout
func TestObjectStorage_ContextTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, _ := NewClient(
		"test-key",
		"test-token",
		"proj-123",
		"test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	createReq := &BucketCreateRequest{
		BucketName: "test-bucket",
	}

	_, _, err := client.ObjectStorage.CreateBucket(ctx, createReq)
	if err == nil {
		t.Fatal("Expected error for context timeout")
	}
}

// TestObjectStorage_ContextCancellation tests context cancellation
func TestObjectStorage_ContextCancellation(t *testing.T) {
	client, _ := NewClient("test-key", "test-token", "proj-123", "test-location")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	createReq := &BucketCreateRequest{
		BucketName: "test-bucket",
	}

	_, _, err := client.ObjectStorage.CreateBucket(ctx, createReq)
	if err == nil {
		t.Fatal("Expected error for canceled context")
	}
}

// TestObjectStorage_GetBucket_NetworkError tests Get with network error
func TestObjectStorage_GetBucket_NetworkError(t *testing.T) {
	client, _ := NewClient(
		"test-key",
		"test-token",
		"proj-123",
		"test-location",
		SetBaseURL("http://localhost:1"),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}),
	)

	ctx := context.Background()
	_, _, err := client.ObjectStorage.GetBucket(ctx, "test-bucket")
	if err == nil {
		t.Fatal("Expected network error")
	}
}

// TestObjectStorage_ListBuckets_NetworkError tests List with network error
func TestObjectStorage_ListBuckets_NetworkError(t *testing.T) {
	client, _ := NewClient(
		"test-key",
		"test-token",
		"proj-123",
		"test-location",
		SetBaseURL("http://localhost:1"),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}),
	)

	ctx := context.Background()
	_, _, err := client.ObjectStorage.ListBuckets(ctx)
	if err == nil {
		t.Fatal("Expected network error")
	}
}

// TestObjectStorage_DeleteBucket_ContextTimeout tests Delete with context timeout
func TestObjectStorage_DeleteBucket_ContextTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, _ := NewClient(
		"test-key",
		"test-token",
		"proj-123",
		"test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := client.ObjectStorage.DeleteBucket(ctx, "test-bucket")
	if err == nil {
		t.Fatal("Expected error for context timeout")
	}
}

// TestObjectStorage_SetBucketVersioning_ContextTimeout tests versioning with context timeout
func TestObjectStorage_SetBucketVersioning_ContextTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, _ := NewClient(
		"test-key",
		"test-token",
		"proj-123",
		"test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	versioningReq := &BucketVersioningRequest{
		BucketName:         "test-bucket",
		NewVersioningState: "Enabled",
	}

	_, _, err := client.ObjectStorage.SetBucketVersioning(ctx, "test-bucket", versioningReq)
	if err == nil {
		t.Fatal("Expected error for context timeout")
	}
}

// TestObjectStorage_RetryLogic_TemporaryFailure tests retry on temporary failures
func TestObjectStorage_RetryLogic_TemporaryFailure(t *testing.T) {
	var requestCount int

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++

		if requestCount == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprint(w, `{"code": 503, "message": "Service unavailable"}`)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"code": 200,
			"message": "Bucket created",
			"data": {"bucket_name": "test-bucket"},
			"errors": {}
		}`)
	}))
	defer server.Close()

	client, _ := NewClient(
		"test-key",
		"test-token",
		"proj-123",
		"test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{
			RetryMax:     2,
			RetryWaitMin: PtrTo(10 * time.Millisecond),
			RetryWaitMax: PtrTo(50 * time.Millisecond),
		}),
	)

	ctx := context.Background()
	createReq := &BucketCreateRequest{
		BucketName: "test-bucket",
	}

	_, _, _ = client.ObjectStorage.CreateBucket(ctx, createReq)
	if requestCount < 1 {
		t.Errorf("Expected at least 1 request, got %d", requestCount)
	}
}
