package client

import (
	"net/http"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func TestAddParamsAndHeadersFunc(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://api.test.com/test/", nil)
	apiKey := "test-key"
	authToken := "test-token"
	projectID := "123"
	location := "us-east"

	modifiedReq := AddParamsAndHeaders(req, apiKey, authToken, projectID, location)

	params := modifiedReq.URL.Query()
	if params.Get("apikey") != apiKey {
		t.Errorf("Expected apikey %s, got %s", apiKey, params.Get("apikey"))
	}

	if params.Get("project_id") != projectID {
		t.Errorf("Expected project_id %s, got %s", projectID, params.Get("project_id"))
	}

	if params.Get("location") != location {
		t.Errorf("Expected location %s, got %s", location, params.Get("location"))
	}

	if modifiedReq.Header.Get("Authorization") != "Bearer "+authToken {
		t.Errorf("Expected Authorization header Bearer %s, got %s", authToken, modifiedReq.Header.Get("Authorization"))
	}

	if modifiedReq.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", modifiedReq.Header.Get("Content-Type"))
	}

	if modifiedReq.Header.Get("User-Agent") != "terraform-e2e" {
		t.Errorf("Expected User-Agent terraform-e2e, got %s", modifiedReq.Header.Get("User-Agent"))
	}
}

func TestNewSfs(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/efs/create/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testURLPath(t, r, "/efs/create/")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "SFS created successfully",
			"data": {
				"id": "sfs-123",
				"name": "test-sfs"
			}
		}`)
	})

	sfsCreate := &models.SfsCreate{
		Name: "test-sfs",
	}

	result, err := ts.client.NewSfs(sfsCreate, "123", "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result["message"] != "SFS created successfully" {
		t.Errorf("Expected message 'SFS created successfully', got %s", result["message"])
	}
}

func TestNewSfsError(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/efs/create/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusBadRequest, `{"error": "invalid request"}`)
	})

	sfsCreate := &models.SfsCreate{
		Name: "test-sfs",
	}

	result, err := ts.client.NewSfs(sfsCreate, "123", "us-east")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestGetSfs(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/efs/sfs-123/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/efs/sfs-123/")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"data": {
				"id": "sfs-123",
				"name": "test-sfs",
				"state": "ACTIVE"
			}
		}`)
	})

	result, err := ts.client.GetSfs("sfs-123", "123", "us-east")

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

func TestGetSfsNon200Status(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/efs/sfs-404/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusNotFound, `{"error": "not found"}`)
	})

	result, err := ts.client.GetSfs("sfs-404", "123", "us-east")

	if err == nil {
		t.Fatal("Expected error for non-200 status, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestDeleteSFs(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/efs/delete/sfs-123/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		testURLPath(t, r, "/efs/delete/sfs-123/")

		w.WriteHeader(http.StatusOK)
	})

	err := ts.client.DeleteSFs("sfs-123", "123", "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestDeleteSFsNon200Status(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/efs/delete/sfs-123/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, `{"error": "server error"}`)
	})

	err := ts.client.DeleteSFs("sfs-123", "123", "us-east")

	if err == nil {
		t.Fatal("Expected error for non-200 status, got nil")
	}
}

func TestGetSfss(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/efs/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/efs/")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Success",
			"data": [
				{
					"id": 1,
					"name": "sfs-one"
				},
				{
					"id": 2,
					"name": "sfs-two"
				}
			]
		}`)
	})

	result, err := ts.client.GetSfss("us-east", "123")

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
		t.Errorf("Expected 2 SFS instances, got %d", len(result.Data))
	}
}

func TestGetSfssNon200Status(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/efs/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusUnauthorized, `{"error": "unauthorized"}`)
	})

	result, err := ts.client.GetSfss("us-east", "123")

	if err == nil {
		t.Fatal("Expected error for non-200 status, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestNewSfsWithHeaders(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/efs/create/", func(w http.ResponseWriter, r *http.Request) {
		testHeader(t, r, "Authorization", "Bearer test-auth-token")
		testHeader(t, r, "Content-Type", "application/json")
		testHeader(t, r, "User-Agent", "terraform-e2e")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Success"
		}`)
	})

	sfsCreate := &models.SfsCreate{
		Name: "test-sfs",
	}

	_, err := ts.client.NewSfs(sfsCreate, "123", "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestGetSfsWithHeaders(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/efs/sfs-123/", func(w http.ResponseWriter, r *http.Request) {
		testHeader(t, r, "Authorization", "Bearer test-auth-token")
		testHeader(t, r, "Content-Type", "application/json")
		testHeader(t, r, "User-Agent", "terraform-e2e")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"data": {}
		}`)
	})

	_, err := ts.client.GetSfs("sfs-123", "123", "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestDeleteSFsWithHeaders(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/efs/delete/sfs-123/", func(w http.ResponseWriter, r *http.Request) {
		testHeader(t, r, "Authorization", "Bearer test-auth-token")
		testHeader(t, r, "Content-Type", "application/json")
		testHeader(t, r, "User-Agent", "terraform-e2e")

		w.WriteHeader(http.StatusOK)
	})

	err := ts.client.DeleteSFs("sfs-123", "123", "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestGetSfssWithHeaders(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/efs/", func(w http.ResponseWriter, r *http.Request) {
		testHeader(t, r, "Authorization", "Bearer test-auth-token")
		testHeader(t, r, "Content-Type", "application/json")
		testHeader(t, r, "User-Agent", "terraform-e2e")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Success",
			"data": []
		}`)
	})

	_, err := ts.client.GetSfss("us-east", "123")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}
