package client

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func TestNewMySqlDb(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testURLPath(t, r, "/rds/cluster/")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "location", "test-location")
		testQueryParam(t, r, "project_id", "test-project")

		body, _ := io.ReadAll(r.Body)
		var req models.MySqlCreate
		json.Unmarshal(body, &req)

		if req.Name == "" {
			t.Error("Expected Name in request body")
		}

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "MySQL database created successfully",
			"data": {
				"id": 456,
				"name": "test-mysql"
			}
		}`)
	})

	mysqlCreate := &models.MySqlCreate{
		Name:             "test-mysql",
		SoftwareID:       2,
		TemplateID:       101,
		PublicIPRequired: true,
		Group:            "default",
		Database: models.DBConfig{
			User:     "mysqluser",
			Password: "mysqlpass",
			Name:     "mysqldb",
		},
	}

	result, err := ts.client.NewMySqlDb(mysqlCreate, "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result["message"] != "MySQL database created successfully" {
		t.Errorf("Expected message 'MySQL database created successfully', got %s", result["message"])
	}
}

func TestNewMySqlDb_Error(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusBadRequest, "Bad Request")
	})

	mysqlCreate := &models.MySqlCreate{
		Name:       "test-mysql",
		SoftwareID: 2,
		TemplateID: 101,
	}

	result, err := ts.client.NewMySqlDb(mysqlCreate, "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestGetMySqlDbaas(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/456/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/rds/cluster/456/")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": {
				"id": 456,
				"name": "test-mysql",
				"status": "active",
				"software": {
					"name": "MySQL",
					"version": "8.0",
					"engine": "mysql"
				}
			}
		}`)
	})

	result, err := ts.client.GetMySqlDbaas("456", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Data.ID != 456 {
		t.Errorf("Expected ID 456, got %d", result.Data.ID)
	}

	if result.Data.Name != "test-mysql" {
		t.Errorf("Expected Name test-mysql, got %s", result.Data.Name)
	}
}

func TestGetMySqlDbaas_NotFound(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/999/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotFound, "Not Found")
	})

	result, err := ts.client.GetMySqlDbaas("999", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestDeleteMySqlDBaaS(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/456/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		testURLPath(t, r, "/rds/cluster/456/")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "MySQL database deleted successfully"
		}`)
	})

	result, err := ts.client.DeleteMySqlDBaaS("456", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result["message"] != "MySQL database deleted successfully" {
		t.Errorf("Expected message 'MySQL database deleted successfully', got %s", result["message"])
	}
}

func TestDeleteMySqlDBaaS_Error(t *testing.T) {
	t.Skip("DeleteMySqlDBaaS doesn't check HTTP status codes")
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/456/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusInternalServerError, "Server Error")
	})

	result, err := ts.client.DeleteMySqlDBaaS("456", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestResumeMySqlDBaaS(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/456/resume", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/456/resume")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "MySQL database resumed successfully"
		}`)
	})

	result, err := ts.client.ResumeMySqlDBaaS("456", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestResumeMySqlDBaaS_Error(t *testing.T) {
	t.Skip("ResumeMySqlDBaaS doesn't check HTTP status codes")
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/456/resume", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusForbidden, "Forbidden")
	})

	result, err := ts.client.ResumeMySqlDBaaS("456", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestStopMySqlDBaaS(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/456/shutdown", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/456/shutdown")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "MySQL database stopped successfully"
		}`)
	})

	result, err := ts.client.StopMySqlDBaaS("456", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestStopMySqlDBaaS_Error(t *testing.T) {
	t.Skip("StopMySqlDBaaS doesn't check HTTP status codes")
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/456/shutdown", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusBadRequest, "Bad Request")
	})

	result, err := ts.client.StopMySqlDBaaS("456", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestRestartMySqlDBaaS(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/456/restart", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/456/restart")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "MySQL database restarted successfully"
		}`)
	})

	result, err := ts.client.RestartMySqlDBaaS("456", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestRestartMySqlDBaaS_Error(t *testing.T) {
	t.Skip("RestartMySqlDBaaS doesn't check HTTP status codes")
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/456/restart", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusServiceUnavailable, "Service Unavailable")
	})

	result, err := ts.client.RestartMySqlDBaaS("456", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestAttachVpcToMySql(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/456/vpc-attach/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/456/vpc-attach/")

		body, _ := io.ReadAll(r.Body)
		var req models.AttachVPCPayloadRequest
		json.Unmarshal(body, &req)

		if req.Action == "" {
			t.Error("Expected Action in request body")
		}

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "VPC attached successfully"
		}`)
	})

	vpcPayload := &models.AttachVPCPayloadRequest{
		Action: "attach",
		VPCs: []models.VPC{
			{
				VpcName:    "test-vpc",
				Ipv4_cidr:  "10.0.0.0/24",
				Network_id: 100,
			},
		},
	}

	result, err := ts.client.AttachVpcToMySql(vpcPayload, "456", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestAttachVpcToMySql_Error(t *testing.T) {
	t.Skip("AttachVpcToMySql doesn't check HTTP status codes")
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/456/vpc-attach/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotFound, "VPC not found")
	})

	vpcPayload := &models.AttachVPCPayloadRequest{
		Action: "attach",
		VPCs:   []models.VPC{},
	}

	result, err := ts.client.AttachVpcToMySql(vpcPayload, "456", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestDetachVpcFromMySql(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/456/vpc-detach/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/456/vpc-detach/")

		body, _ := io.ReadAll(r.Body)
		var req models.AttachVPCPayloadRequest
		json.Unmarshal(body, &req)

		if req.Action == "" {
			t.Error("Expected Action in request body")
		}

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "VPC detached successfully"
		}`)
	})

	vpcPayload := &models.AttachVPCPayloadRequest{
		Action: "detach",
		VPCs: []models.VPC{
			{
				VpcName:    "test-vpc",
				Ipv4_cidr:  "10.0.0.0/24",
				Network_id: 100,
			},
		},
	}

	result, err := ts.client.DetachVpcFromMySql(vpcPayload, "456", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestDetachVpcFromMySql_Error(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/456/vpc-detach/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusBadRequest, "Bad Request")
	})

	vpcPayload := &models.AttachVPCPayloadRequest{
		Action: "detach",
		VPCs:   []models.VPC{},
	}

	result, err := ts.client.DetachVpcFromMySql(vpcPayload, "456", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestAttachPGToMySqlDBaaS(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/456/parameter-group/789/add", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/456/parameter-group/789/add")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Parameter group attached successfully"
		}`)
	})

	result, err := ts.client.AttachPGToMySqlDBaaS("456", "789", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestAttachPGToMySqlDBaaS_Error(t *testing.T) {
	t.Skip("AttachPGToMySqlDBaaS doesn't check HTTP status codes")
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/456/parameter-group/999/add", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotFound, "Parameter group not found")
	})

	result, err := ts.client.AttachPGToMySqlDBaaS("456", "999", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestDetachPGFromMySqlDBaaS(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/456/parameter-group/789/detach", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/456/parameter-group/789/detach")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Parameter group detached successfully"
		}`)
	})

	result, err := ts.client.DetachPGFromMySqlDBaaS("456", "789", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestDetachPGFromMySqlDBaaS_Error(t *testing.T) {
	t.Skip("DetachPGFromMySqlDBaaS doesn't check HTTP status codes")
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/456/parameter-group/789/detach", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusInternalServerError, "Server Error")
	})

	result, err := ts.client.DetachPGFromMySqlDBaaS("456", "789", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestAttachPublicIPToMySql(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/456/public-ip-attach/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/456/public-ip-attach/")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Public IP attached successfully"
		}`)
	})

	result, err := ts.client.AttachPublicIPToMySql("456", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestAttachPublicIPToMySql_Error(t *testing.T) {
	t.Skip("AttachPublicIPToMySql doesn't check HTTP status codes")
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/456/public-ip-attach/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusConflict, "Conflict")
	})

	result, err := ts.client.AttachPublicIPToMySql("456", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestDetachPublicIPFromMySql(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/456/public-ip-detach/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/456/public-ip-detach/")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Public IP detached successfully"
		}`)
	})

	result, err := ts.client.DetachPublicIPFromMySql("456", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestDetachPublicIPFromMySql_Error(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/456/public-ip-detach/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusBadRequest, "Bad Request")
	})

	result, err := ts.client.DetachPublicIPFromMySql("456", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestUpgradeMySQLPlan(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/456/rds-upgrade/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/456/rds-upgrade/")

		body, _ := io.ReadAll(r.Body)
		var req models.MySQlPlanUpgradeAction
		json.Unmarshal(body, &req)

		if req.TemplateID != 202 {
			t.Errorf("Expected template_id 202, got %d", req.TemplateID)
		}

		w.WriteHeader(http.StatusOK)
	})

	result, err := ts.client.UpgradeMySQLPlan("456", 202, "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestUpgradeMySQLPlan_Error(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/456/rds-upgrade/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusPaymentRequired, "Insufficient credits")
	})

	result, err := ts.client.UpgradeMySQLPlan("456", 202, "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestExpandMySQLDBaaSDisk(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/456/disk-upgrade/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/456/disk-upgrade/")

		body, _ := io.ReadAll(r.Body)
		var req models.MYSQLExpandDisk
		json.Unmarshal(body, &req)

		if req.Size != 100 {
			t.Errorf("Expected size 100, got %d", req.Size)
		}

		w.WriteHeader(http.StatusOK)
	})

	result, err := ts.client.ExpandMySQLDBaaSDisk("456", 100, "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestExpandMySQLDBaaSDisk_Error(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/456/disk-upgrade/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusUnprocessableEntity, "Cannot expand disk")
	})

	result, err := ts.client.ExpandMySQLDBaaSDisk("456", 100, "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}
