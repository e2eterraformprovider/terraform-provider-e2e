package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func TestCreateMariaDB(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testURLPath(t, r, "/rds/cluster/")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "location", "test-location")
		testQueryParam(t, r, "project_id", "test-project")

		body, _ := ioutil.ReadAll(r.Body)
		var req models.MariaDBCreateRequest
		_ = json.Unmarshal(body, &req)

		if req.Name == "" {
			t.Error("Expected Name in request body")
		}

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": {
				"id": 123,
				"name": "test-mariadb",
				"status": "active"
			}
		}`)
	})

	req := &models.MariaDBCreateRequest{
		Name:             "test-mariadb",
		SoftwareID:       1,
		TemplateID:       100,
		PublicIPRequired: true,
		Group:            "default",
		Database: models.DBConfig{
			User:     "testuser",
			Password: "testpass",
			Name:     "testdb",
		},
	}

	result, err := ts.client.CreateMariaDB(req, "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.ID != 123 {
		t.Errorf("Expected ID 123, got %d", result.ID)
	}

	if result.Name != "test-mariadb" {
		t.Errorf("Expected Name test-mariadb, got %s", result.Name)
	}
}

func TestCreateMariaDB_Error(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusInternalServerError, "Internal Server Error")
	})

	req := &models.MariaDBCreateRequest{
		Name:       "test-mariadb",
		SoftwareID: 1,
		TemplateID: 100,
	}

	result, err := ts.client.CreateMariaDB(req, "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestReadMariaDB(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/rds/cluster/123/")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": {
				"id": 123,
				"name": "test-mariadb",
				"status": "active",
				"software": {
					"name": "MariaDB",
					"version": "10.5",
					"engine": "mariadb"
				}
			}
		}`)
	})

	result, err := ts.client.ReadMariaDB("123", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.ID != 123 {
		t.Errorf("Expected ID 123, got %d", result.ID)
	}

	if result.Name != "test-mariadb" {
		t.Errorf("Expected Name test-mariadb, got %s", result.Name)
	}
}

func TestReadMariaDB_NotFound(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/999/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotFound, "Not Found")
	})

	result, err := ts.client.ReadMariaDB("999", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestMariaDBExists(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		expectedExists bool
		expectError    bool
	}{
		{
			name:           "Exists - 200 OK",
			statusCode:     http.StatusOK,
			expectedExists: true,
			expectError:    false,
		},
		{
			name:           "Not Found - 404",
			statusCode:     http.StatusNotFound,
			expectedExists: false,
			expectError:    false,
		},
		{
			name:           "Error - 500",
			statusCode:     http.StatusInternalServerError,
			expectedExists: false,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := setup()
			defer ts.teardown()

			ts.mux.HandleFunc("/rds/cluster/123/", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, http.MethodGet)

				w.WriteHeader(tt.statusCode)
				if tt.statusCode == http.StatusOK {
					writeJSON(w, tt.statusCode, `{"code": 200}`)
				} else {
					_, _ = w.Write([]byte("Error"))
				}
			})

			exists, err := ts.client.MariaDBExists("123", "test-project", "test-location")

			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}

			if exists != tt.expectedExists {
				t.Errorf("Expected exists: %v, got: %v", tt.expectedExists, exists)
			}
		})
	}
}

func TestDeleteMariaDB(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		testURLPath(t, r, "/rds/cluster/123/")

		w.WriteHeader(http.StatusOK)
	})

	err := ts.client.DeleteMariaDB("123", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestDeleteMariaDB_Error(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	err := ts.client.DeleteMariaDB("123", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestShutdownMariaDB(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/shutdown", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/123/shutdown")

		w.WriteHeader(http.StatusOK)
	})

	err := ts.client.ShutdownMariaDB("123", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestShutdownMariaDB_Error(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/shutdown", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusBadRequest, "Bad Request")
	})

	err := ts.client.ShutdownMariaDB("123", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestResumeMariaDB(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/resume", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/123/resume")

		w.WriteHeader(http.StatusOK)
	})

	err := ts.client.ResumeMariaDB("123", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestResumeMariaDB_Error(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/resume", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusForbidden, "Forbidden")
	})

	err := ts.client.ResumeMariaDB("123", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestRestartMariaDB(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/restart", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/123/restart")

		w.WriteHeader(http.StatusOK)
	})

	err := ts.client.RestartMariaDB("123", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestRestartMariaDB_Error(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/restart", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusServiceUnavailable, "Service Unavailable")
	})

	err := ts.client.RestartMariaDB("123", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestAttachVPCToMariaDB(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	callCount := 0
	ts.mux.HandleFunc("/vpc/100/", func(w http.ResponseWriter, r *http.Request) {
		if callCount == 0 {
			callCount++
			writeJSON(w, http.StatusOK, `{
				"code": 200,
				"message": "success",
				"data": {
					"name": "test-vpc",
					"network_id": 100,
					"ipv4_cidr": "10.0.0.0/24",
					"state": "Active"
				}
			}`)
		}
	})

	ts.mux.HandleFunc("/rds/cluster/123/vpc-attach/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/123/vpc-attach/")

		body, _ := ioutil.ReadAll(r.Body)
		var req models.AttachDetachVPCRequest
		_ = json.Unmarshal(body, &req)

		if req.Action != "attach" {
			t.Errorf("Expected action 'attach', got %s", req.Action)
		}

		w.WriteHeader(http.StatusOK)
	})

	err := ts.client.AttachVPCToMariaDB("123", "test-project", "test-location", []string{"100"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestAttachVPCToMariaDB_Error(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/vpc/999/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotFound, "VPC not found")
	})

	err := ts.client.AttachVPCToMariaDB("123", "test-project", "test-location", []string{"999"})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestDetachVPCFromMariaDB(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	callCount := 0
	ts.mux.HandleFunc("/vpc/100/", func(w http.ResponseWriter, r *http.Request) {
		if callCount == 0 {
			callCount++
			writeJSON(w, http.StatusOK, `{
				"code": 200,
				"message": "success",
				"data": {
					"name": "test-vpc",
					"network_id": 100,
					"ipv4_cidr": "10.0.0.0/24",
					"state": "Active"
				}
			}`)
		}
	})

	ts.mux.HandleFunc("/rds/cluster/123/vpc-detach/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/123/vpc-detach/")

		body, _ := ioutil.ReadAll(r.Body)
		var req models.AttachDetachVPCRequest
		_ = json.Unmarshal(body, &req)

		if req.Action != "detach" {
			t.Errorf("Expected action 'detach', got %s", req.Action)
		}

		w.WriteHeader(http.StatusOK)
	})

	err := ts.client.DetachVPCFromMariaDB("123", "test-project", "test-location", []string{"100"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestDetachVPCFromMariaDB_Error(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/vpc/999/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotFound, "VPC not found")
	})

	err := ts.client.DetachVPCFromMariaDB("123", "test-project", "test-location", []string{"999"})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestAttachPublicIPToMariaDB(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/public-ip-attach/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/123/public-ip-attach/")

		body, _ := ioutil.ReadAll(r.Body)
		var payload map[string]string
		json.Unmarshal(body, &payload)

		if payload["action"] != "attach" {
			t.Errorf("Expected action 'attach', got %s", payload["action"])
		}

		w.WriteHeader(http.StatusOK)
	})

	err := ts.client.AttachPublicIPToMariaDB("123", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestAttachPublicIPToMariaDB_Error(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/public-ip-attach/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusConflict, "Conflict")
	})

	err := ts.client.AttachPublicIPToMariaDB("123", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestDetachPublicIPFromMariaDB(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/public-ip-detach/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/123/public-ip-detach/")

		body, _ := ioutil.ReadAll(r.Body)
		var payload map[string]string
		json.Unmarshal(body, &payload)

		if payload["action"] != "detach" {
			t.Errorf("Expected action 'detach', got %s", payload["action"])
		}

		w.WriteHeader(http.StatusOK)
	})

	err := ts.client.DetachPublicIPFromMariaDB("123", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestDetachPublicIPFromMariaDB_Error(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/public-ip-detach/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusBadRequest, "Bad Request")
	})

	err := ts.client.DetachPublicIPFromMariaDB("123", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestAttachParameterGroupToMariaDB(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/parameter-group/456/add", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/123/parameter-group/456/add")

		body, _ := ioutil.ReadAll(r.Body)
		var req models.ParameterGroupRequest
		_ = json.Unmarshal(body, &req)

		if req.Action != "add" {
			t.Errorf("Expected action 'add', got %s", req.Action)
		}

		w.WriteHeader(http.StatusOK)
	})

	err := ts.client.AttachParameterGroupToMariaDB("123", 456, "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestAttachParameterGroupToMariaDB_Error(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/parameter-group/999/add", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotFound, "Parameter group not found")
	})

	err := ts.client.AttachParameterGroupToMariaDB("123", 999, "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestDetachParameterGroupFromMariaDB(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/parameter-group/456/detach", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/123/parameter-group/456/detach")

		w.WriteHeader(http.StatusOK)
	})

	err := ts.client.DetachParameterGroupFromMariaDB("123", 456, "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestDetachParameterGroupFromMariaDB_Error(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/parameter-group/456/detach", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusInternalServerError, "Server Error")
	})

	err := ts.client.DetachParameterGroupFromMariaDB("123", 456, "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestUpgradeMariaDBPlan(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/rds-upgrade/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/123/rds-upgrade/")

		body, _ := ioutil.ReadAll(r.Body)
		var payload map[string]interface{}
		json.Unmarshal(body, &payload)

		templateID := int(payload["template_id"].(float64))
		if templateID != 200 {
			t.Errorf("Expected template_id 200, got %d", templateID)
		}

		w.WriteHeader(http.StatusOK)
	})

	err := ts.client.UpgradeMariaDBPlan("123", "test-project", "test-location", 200)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestUpgradeMariaDBPlan_Error(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/rds-upgrade/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusPaymentRequired, "Insufficient credits")
	})

	err := ts.client.UpgradeMariaDBPlan("123", "test-project", "test-location", 200)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestExpandMariaDBDisk(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/disk-upgrade/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/123/disk-upgrade/")

		body, _ := ioutil.ReadAll(r.Body)
		var req models.DiskUpgradeRequest
		_ = json.Unmarshal(body, &req)

		if req.Size != 50 {
			t.Errorf("Expected size 50, got %d", req.Size)
		}

		w.WriteHeader(http.StatusOK)
	})

	err := ts.client.ExpandMariaDBDisk("123", "test-project", "test-location", 50)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestExpandMariaDBDisk_ZeroSize(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/disk-upgrade/", func(w http.ResponseWriter, r *http.Request) {
		t.Error("Expected no HTTP request for zero size")
	})

	err := ts.client.ExpandMariaDBDisk("123", "test-project", "test-location", 0)

	if err != nil {
		t.Fatalf("Expected no error for zero size (should skip), got: %v", err)
	}
}

func TestExpandMariaDBDisk_NegativeSize(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/disk-upgrade/", func(w http.ResponseWriter, r *http.Request) {
		t.Error("Expected no HTTP request for negative size")
	})

	err := ts.client.ExpandMariaDBDisk("123", "test-project", "test-location", -10)

	if err != nil {
		t.Fatalf("Expected no error for negative size (should skip), got: %v", err)
	}
}

func TestExpandMariaDBDisk_Error(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/disk-upgrade/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusUnprocessableEntity, "Cannot expand disk")
	})

	err := ts.client.ExpandMariaDBDisk("123", "test-project", "test-location", 50)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}
