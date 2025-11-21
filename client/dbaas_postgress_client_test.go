package client

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func TestCreatePostgressDB(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testURLPath(t, r, "/rds/cluster/")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "location", "test-location")
		testQueryParam(t, r, "project_id", "test-project")

		body, _ := io.ReadAll(r.Body)
		var req models.DBCreateRequest
		json.Unmarshal(body, &req)

		if req.Name == "" {
			t.Error("Expected Name in request body")
		}

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "PostgreSQL database created successfully",
			"data": {
				"id": 789,
				"name": "test-postgres"
			}
		}`)
	})

	dbCreate := models.DBCreateRequest{
		Name:             "test-postgres",
		SoftwareID:       3,
		TemplateID:       103,
		PublicIPRequired: true,
		Group:            "default",
		Database: models.DBConfig{
			User:     "pguser",
			Password: "pgpass",
			Name:     "pgdb",
		},
	}

	result, err := ts.client.CreatePostgressDB(dbCreate, "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result["message"] != "PostgreSQL database created successfully" {
		t.Errorf("Expected message 'PostgreSQL database created successfully', got %s", result["message"])
	}
}

func TestCreatePostgressDB_Error(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusBadRequest, "Bad Request")
	})

	dbCreate := models.DBCreateRequest{
		Name:       "test-postgres",
		SoftwareID: 3,
		TemplateID: 103,
	}

	result, err := ts.client.CreatePostgressDB(dbCreate, "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestGetPostgressDB(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/789/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/rds/cluster/789/")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": {
				"id": 789,
				"name": "test-postgres",
				"status": "active"
			}
		}`)
	})

	result, err := ts.client.GetPostgressDB("789", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	data := result["data"].(map[string]interface{})
	if data["name"] != "test-postgres" {
		t.Errorf("Expected name test-postgres, got %s", data["name"])
	}
}

func TestGetPostgressDB_Error(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/999/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotFound, "Not Found")
	})

	result, err := ts.client.GetPostgressDB("999", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestDeletePostgressDB(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/789/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		testURLPath(t, r, "/rds/cluster/789/")

		w.WriteHeader(http.StatusOK)
	})

	err := ts.client.DeletePostgressDB("789", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestDeletePostgressDB_Error(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/789/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusInternalServerError, "Server Error")
	})

	err := ts.client.DeletePostgressDB("789", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestStopPostgressDB(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/789/shutdown", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/789/shutdown")

		w.WriteHeader(http.StatusOK)
	})

	err := ts.client.StopPostgressDB("789", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestStopPostgressDB_Error(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/789/shutdown", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusBadRequest, "Bad Request")
	})

	err := ts.client.StopPostgressDB("789", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestStartPostgressDB(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/789/resume", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/789/resume")

		w.WriteHeader(http.StatusOK)
	})

	err := ts.client.StartPostgressDB("789", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestStartPostgressDB_Error(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/789/resume", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusForbidden, "Forbidden")
	})

	err := ts.client.StartPostgressDB("789", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestRestartPostgressDB(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/789/restart", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/789/restart")

		w.WriteHeader(http.StatusOK)
	})

	err := ts.client.RestartPostgressDB("789", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestRestartPostgressDB_Error(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/789/restart", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusServiceUnavailable, "Service Unavailable")
	})

	err := ts.client.RestartPostgressDB("789", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestAttachPublicIpPostgressDB(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/789/public-ip-attach/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/789/public-ip-attach/")

		w.WriteHeader(http.StatusOK)
	})

	err := ts.client.AttachPublicIpPostgressDB("789", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestAttachPublicIpPostgressDB_Error(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/789/public-ip-attach/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusConflict, "Conflict")
	})

	err := ts.client.AttachPublicIpPostgressDB("789", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestDetachPublicIpPostgressDB(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/789/public-ip-detach/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/789/public-ip-detach/")

		w.WriteHeader(http.StatusOK)
	})

	err := ts.client.DetachPublicIpPostgressDB("789", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestDetachPublicIpPostgressDB_Error(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/789/public-ip-detach/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusBadRequest, "Bad Request")
	})

	err := ts.client.DetachPublicIpPostgressDB("789", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestAttachVPCPostgressDB(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/789/vpc-attach/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/789/vpc-attach/")

		body, _ := io.ReadAll(r.Body)
		var req models.AttachVPCPayloadRequest
		json.Unmarshal(body, &req)

		if req.Action == "" {
			t.Error("Expected Action in request body")
		}

		w.WriteHeader(http.StatusOK)
	})

	vpcPayload := models.AttachVPCPayloadRequest{
		Action: "attach",
		VPCs: []models.VPC{
			{
				VpcName:    "test-vpc",
				Ipv4_cidr:  "10.0.0.0/24",
				Network_id: 100,
			},
		},
	}

	err := ts.client.AttachVPCPostgressDB(vpcPayload, "789", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestAttachVPCPostgressDB_Error(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/789/vpc-attach/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotFound, "VPC not found")
	})

	vpcPayload := models.AttachVPCPayloadRequest{
		Action: "attach",
		VPCs:   []models.VPC{},
	}

	err := ts.client.AttachVPCPostgressDB(vpcPayload, "789", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestDetachVPCPostgressDB(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/789/vpc-detach/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/789/vpc-detach/")

		body, _ := io.ReadAll(r.Body)
		var req models.AttachVPCPayloadRequest
		json.Unmarshal(body, &req)

		if req.Action == "" {
			t.Error("Expected Action in request body")
		}

		w.WriteHeader(http.StatusOK)
	})

	vpcPayload := models.AttachVPCPayloadRequest{
		Action: "detach",
		VPCs: []models.VPC{
			{
				VpcName:    "test-vpc",
				Ipv4_cidr:  "10.0.0.0/24",
				Network_id: 100,
			},
		},
	}

	err := ts.client.DetachVPCPostgressDB(vpcPayload, "789", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestDetachVPCPostgressDB_Error(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/789/vpc-detach/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusBadRequest, "Bad Request")
	})

	vpcPayload := models.AttachVPCPayloadRequest{
		Action: "detach",
		VPCs:   []models.VPC{},
	}

	err := ts.client.DetachVPCPostgressDB(vpcPayload, "789", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestUpgradePostgressPlan(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/789/rds-upgrade/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/789/rds-upgrade/")

		body, _ := io.ReadAll(r.Body)
		var req map[string]interface{}
		json.Unmarshal(body, &req)

		templateID := int(req["template_id"].(float64))
		if templateID != 204 {
			t.Errorf("Expected template_id 204, got %d", templateID)
		}

		w.WriteHeader(http.StatusOK)
	})

	result, err := ts.client.UpgradePostgressPlan("789", 204, "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestUpgradePostgressPlan_Error(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/789/rds-upgrade/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusPaymentRequired, "Insufficient credits")
	})

	result, err := ts.client.UpgradePostgressPlan("789", 204, "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestUpdateParameterGroup(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/789/parameter-group/555/add", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/789/parameter-group/555/add")

		w.WriteHeader(http.StatusOK)
	})

	err := ts.client.UpdateParameterGroup("789", "555", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestUpdateParameterGroup_Error(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/789/parameter-group/999/add", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotFound, "Parameter group not found")
	})

	err := ts.client.UpdateParameterGroup("789", "999", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestUpgradeDiskStorage(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/789/disk-upgrade/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/789/disk-upgrade/")

		body, _ := io.ReadAll(r.Body)
		var req map[string]interface{}
		json.Unmarshal(body, &req)

		size := int(req["size"].(float64))
		if size != 150 {
			t.Errorf("Expected size 150, got %d", size)
		}

		w.WriteHeader(http.StatusOK)
	})

	err := ts.client.UpgradeDiskStorage("789", 150, "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestUpgradeDiskStorage_Error(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/789/disk-upgrade/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusUnprocessableEntity, "Cannot upgrade disk")
	})

	err := ts.client.UpgradeDiskStorage("789", 150, "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}
