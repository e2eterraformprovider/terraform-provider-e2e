package client

import (
	"net/http"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func TestSetParamsAndHeaders(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	req, _ := http.NewRequest("GET", ts.server.URL+"/test/", nil)
	location := "us-east"
	projectID := "123"

	modifiedReq := ts.client.setParamsAndHeaders(req, location, projectID)

	params := modifiedReq.URL.Query()
	if params.Get("apikey") != "test-api-key" {
		t.Errorf("Expected apikey test-api-key, got %s", params.Get("apikey"))
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

	if modifiedReq.Header.Get("Authorization") != "Bearer test-auth-token" {
		t.Errorf("Expected Authorization header Bearer test-auth-token, got %s", modifiedReq.Header.Get("Authorization"))
	}

	if modifiedReq.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", modifiedReq.Header.Get("Content-Type"))
	}

	if modifiedReq.Header.Get("User-Agent") != "terraform-e2e" {
		t.Errorf("Expected User-Agent terraform-e2e, got %s", modifiedReq.Header.Get("User-Agent"))
	}
}

func TestCreateBucket(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/storage/buckets/test-bucket/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testURLPath(t, r, "/storage/buckets/test-bucket/")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "location", "us-east")
		testQueryParam(t, r, "project_id", "123")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Bucket created successfully",
			"data": {
				"bucket_name": "test-bucket"
			}
		}`)
	})

	bucketPayload := &models.ObjectStorePayload{
		BucketName: "test-bucket",
		Region:     "us-east",
		ProjectID:  123,
	}

	result, err := ts.client.CreateBucket(bucketPayload)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result["message"] != "Bucket created successfully" {
		t.Errorf("Expected message Bucket created successfully, got %s", result["message"])
	}
}

func TestCreateBucketError(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/storage/buckets/invalid-bucket/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusBadRequest, "invalid bucket name")
	})

	bucketPayload := &models.ObjectStorePayload{
		BucketName: "invalid-bucket",
		Region:     "us-east",
		ProjectID:  123,
	}

	result, err := ts.client.CreateBucket(bucketPayload)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestGetBuckets(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/storage/buckets/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/storage/buckets/")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "location", "us-east")
		testQueryParam(t, r, "project_id", "123")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Success",
			"data": [
				{
					"name": "bucket-1",
					"status": "ACTIVE"
				},
				{
					"name": "bucket-2",
					"status": "ACTIVE"
				}
			]
		}`)
	})

	result, err := ts.client.GetBuckets("us-east", "123")

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
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/storage/buckets/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusInternalServerError, "server error")
	})

	result, err := ts.client.GetBuckets("us-east", "123")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestGetBucket(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/storage/buckets/test-bucket/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/storage/buckets/test-bucket/")
		testQueryParam(t, r, "apikey", "test-api-key")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"data": {
				"name": "test-bucket",
				"region": "us-east"
			}
		}`)
	})

	result, err := ts.client.GetBucket("test-bucket", "us-east", "123")

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
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/storage/buckets/nonexistent-bucket/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotFound, "bucket not found")
	})

	result, err := ts.client.GetBucket("nonexistent-bucket", "us-east", "123")

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
			ts := setup()
			defer ts.teardown()

			ts.mux.HandleFunc("/storage/bucket_versioning/test-bucket/", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, http.MethodPut)
				testURLPath(t, r, "/storage/bucket_versioning/test-bucket/")

				writeJSON(w, http.StatusOK, `{
					"code": 200,
					"message": "Versioning updated successfully"
				}`)
			})

			result, err := ts.client.SetBucketVersioning("test-bucket", "us-east", "123", tt.action)

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
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/storage/bucket_versioning/test-bucket/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusBadRequest, "invalid versioning state")
	})

	result, err := ts.client.SetBucketVersioning("test-bucket", "us-east", "123", "Invalid")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestDeleteBucket(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/storage/buckets/test-bucket/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		testURLPath(t, r, "/storage/buckets/test-bucket/")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "location", "us-east")
		testQueryParam(t, r, "project_id", "123")

		w.WriteHeader(http.StatusOK)
	})

	err := ts.client.DeleteBucket("test-bucket", "us-east", "123")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestDeleteBucketError(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/storage/buckets/test-bucket/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusConflict, "bucket not empty")
	})

	err := ts.client.DeleteBucket("test-bucket", "us-east", "123")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestGetBucketsNon200Status(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/storage/buckets/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusForbidden, "forbidden")
	})

	result, err := ts.client.GetBuckets("us-east", "123")

	if err == nil {
		t.Fatal("Expected error for non-200 status, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestGetBucketNon200Status(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/storage/buckets/test-bucket/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusUnauthorized, "unauthorized")
	})

	result, err := ts.client.GetBucket("test-bucket", "us-east", "123")

	if err == nil {
		t.Fatal("Expected error for non-200 status, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestDeleteBucketNon200Status(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/storage/buckets/nonexistent-bucket/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotFound, "bucket not found")
	})

	err := ts.client.DeleteBucket("nonexistent-bucket", "us-east", "123")

	if err == nil {
		t.Fatal("Expected error for non-200 status, got nil")
	}
}
