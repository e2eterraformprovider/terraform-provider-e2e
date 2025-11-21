package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func TestSetParamsAndHeaders(t *testing.T) {
	req := httptest.NewRequest("GET", "https://api.test.com/test/", nil)
	location := "us-east"
	projectID := "123"

	client := NewClient("test-key", "test-token", "https://api.test.com/")
	modifiedReq := client.setParamsAndHeaders(req, location, projectID)

	params := modifiedReq.URL.Query()
	if params.Get("apikey") != "test-key" {
		t.Errorf("Expected apikey test-key, got %s", params.Get("apikey"))
	}

	if params.Get("location") != location {
		t.Errorf("Expected location %s, got %s", location, params.Get("location"))
	}

	if params.Get("project_id") != projectID {
		t.Errorf("Expected project_id %s, got %s", projectID, params.Get("project_id"))
	}

	if params.Get("contact_person_id") != "null" {
		t.Errorf("Expected contact_person_id null, got %s", params.Get("contact_person_id"))
	}

	if modifiedReq.Header.Get("Authorization") != "Bearer test-token" {
		t.Errorf("Expected Authorization header Bearer test-token, got %s", modifiedReq.Header.Get("Authorization"))
	}

	if modifiedReq.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", modifiedReq.Header.Get("Content-Type"))
	}

	if modifiedReq.Header.Get("User-Agent") != "terraform-e2e" {
		t.Errorf("Expected User-Agent terraform-e2e, got %s", modifiedReq.Header.Get("User-Agent"))
	}
}

func TestCreateBucket(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code":    200,
		"message": "Bucket created successfully",
		"data": map[string]interface{}{
			"bucket_name": "test-bucket",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/storage/buckets/test-bucket/" {
			t.Errorf("Expected path /storage/buckets/test-bucket/, got %s", r.URL.Path)
		}

		query := r.URL.Query()
		if query.Get("apikey") == "" {
			t.Error("Expected apikey parameter")
		}
		if query.Get("location") == "" {
			t.Error("Expected location parameter")
		}
		if query.Get("project_id") == "" {
			t.Error("Expected project_id parameter")
		}

		body, _ := ioutil.ReadAll(r.Body)
		var payload models.ObjectStorePayload
		json.Unmarshal(body, &payload)

		if payload.BucketName == "" {
			t.Error("Expected BucketName in request body")
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	bucketPayload := &models.ObjectStorePayload{
		BucketName: "test-bucket",
		Region:     "us-east",
		ProjectID:  123,
	}

	result, err := client.CreateBucket(bucketPayload)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result["message"] != mockResponse["message"] {
		t.Errorf("Expected message %s, got %s", mockResponse["message"], result["message"])
	}
}

func TestCreateBucketError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid bucket name"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	bucketPayload := &models.ObjectStorePayload{
		BucketName: "invalid-bucket",
		Region:     "us-east",
		ProjectID:  123,
	}

	result, err := client.CreateBucket(bucketPayload)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestGetBuckets(t *testing.T) {
	mockResponse := models.ResponseBuckets{
		Code:    200,
		Message: "Success",
		Data: []models.ObjectStore{
			{
				Name:   "bucket-1",
				Status: "ACTIVE",
			},
			{
				Name:   "bucket-2",
				Status: "ACTIVE",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/storage/buckets/" {
			t.Errorf("Expected path /storage/buckets/, got %s", r.URL.Path)
		}

		query := r.URL.Query()
		if query.Get("apikey") == "" {
			t.Error("Expected apikey parameter")
		}
		if query.Get("location") == "" {
			t.Error("Expected location parameter")
		}
		if query.Get("project_id") == "" {
			t.Error("Expected project_id parameter")
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	result, err := client.GetBuckets("us-east", "123")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Code != 200 {
		t.Errorf("Expected code 200, got %d", result.Code)
	}

	if len(result.Data) != 2 {
		t.Errorf("Expected 2 buckets, got %d", len(result.Data))
	}
}

func TestGetBucketsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "server error"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	result, err := client.GetBuckets("us-east", "123")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestGetBucket(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{
			"name":   "test-bucket",
			"region": "us-east",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/storage/buckets/test-bucket/" {
			t.Errorf("Expected path /storage/buckets/test-bucket/, got %s", r.URL.Path)
		}

		query := r.URL.Query()
		if query.Get("apikey") == "" {
			t.Error("Expected apikey parameter")
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	result, err := client.GetBucket("test-bucket", "us-east", "123")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	code := result["code"].(float64)
	if code != 200 {
		t.Errorf("Expected code 200, got %v", code)
	}
}

func TestGetBucketError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "bucket not found"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	result, err := client.GetBucket("nonexistent-bucket", "us-east", "123")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestSetBucketVersioning(t *testing.T) {
	tests := []struct {
		name   string
		action string
	}{
		{
			name:   "Enable versioning",
			action: "Enabled",
		},
		{
			name:   "Suspend versioning",
			action: "Suspended",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockResponse := map[string]interface{}{
				"code":    200,
				"message": "Versioning updated successfully",
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "PUT" {
					t.Errorf("Expected PUT request, got %s", r.Method)
				}

				if r.URL.Path != "/storage/bucket_versioning/test-bucket/" {
					t.Errorf("Expected path /storage/bucket_versioning/test-bucket/, got %s", r.URL.Path)
				}

				body, _ := ioutil.ReadAll(r.Body)
				var data map[string]string
				json.Unmarshal(body, &data)

				if data["bucket_name"] != "test-bucket" {
					t.Errorf("Expected bucket_name test-bucket, got %s", data["bucket_name"])
				}

				if data["new_versioning_state"] != tt.action {
					t.Errorf("Expected new_versioning_state %s, got %s", tt.action, data["new_versioning_state"])
				}

				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(mockResponse)
			}))
			defer server.Close()

			client := NewClient("test-key", "test-token", server.URL)
			result, err := client.SetBucketVersioning("test-bucket", "us-east", "123", tt.action)

			if err != nil {
				t.Fatalf("Expected no error, got: %v", err)
			}

			if result == nil {
				t.Fatal("Expected result, got nil")
			}
		})
	}
}

func TestSetBucketVersioningError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid versioning state"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	result, err := client.SetBucketVersioning("test-bucket", "us-east", "123", "Invalid")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestDeleteBucket(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/storage/buckets/test-bucket/" {
			t.Errorf("Expected path /storage/buckets/test-bucket/, got %s", r.URL.Path)
		}

		query := r.URL.Query()
		if query.Get("apikey") == "" {
			t.Error("Expected apikey parameter")
		}
		if query.Get("location") == "" {
			t.Error("Expected location parameter")
		}
		if query.Get("project_id") == "" {
			t.Error("Expected project_id parameter")
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	err := client.DeleteBucket("test-bucket", "us-east", "123")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestDeleteBucketError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(`{"error": "bucket not empty"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	err := client.DeleteBucket("test-bucket", "us-east", "123")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestCreateBucketWithHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			t.Error("Expected Authorization header")
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		if r.Header.Get("User-Agent") != "terraform-e2e" {
			t.Errorf("Expected User-Agent terraform-e2e, got %s", r.Header.Get("User-Agent"))
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    200,
			"message": "Success",
		})
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	bucketPayload := &models.ObjectStorePayload{
		BucketName: "test-bucket",
		Region:     "us-east",
		ProjectID:  123,
	}

	_, err := client.CreateBucket(bucketPayload)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestGetBucketsNon200Status(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"error": "forbidden"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	result, err := client.GetBuckets("us-east", "123")

	if err == nil {
		t.Fatal("Expected error for non-200 status, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestGetBucketNon200Status(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "unauthorized"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	result, err := client.GetBucket("test-bucket", "us-east", "123")

	if err == nil {
		t.Fatal("Expected error for non-200 status, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestDeleteBucketNon200Status(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "bucket not found"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	err := client.DeleteBucket("nonexistent-bucket", "us-east", "123")

	if err == nil {
		t.Fatal("Expected error for non-200 status, got nil")
	}
}
