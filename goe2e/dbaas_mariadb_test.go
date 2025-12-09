package goe2e

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

// =====================================================
// CRUD Operations Tests
// =====================================================

func TestMariaDBCreate(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "MariaDB cluster created successfully",
			"data": {
				"id": 123,
				"name": "test-mariadb",
				"status": "creating",
				"status_title": "Creating",
				"num_instances": 1,
				"software": {
					"name": "MariaDB",
					"version": "10.6",
					"engine": "mariadb"
				},
				"master_node": {
					"node_name": "master-1",
					"instance_id": 456,
					"cluster_id": 123,
					"state": "active"
				},
				"isEncryptionEnabled": true
			}
		}`)
	})

	createReq := &MariaDBCreateRequest{
		Name:       "test-mariadb",
		SoftwareID: 10,
		TemplateID: 20,
		Database: DBConfig{
			User:     "admin",
			Password: "password",
			Name:     "testdb",
		},
		IsEncryptionEnabled: true,
	}

	result, resp, err := ts.client.MariaDB.CreateMariaDB(context.Background(), createReq)
	assertNoError(t, err)

	assertNotNil(t, result, "Expected result")

	if result.ID != 123 {
		t.Errorf("Expected ID 123, got %d", result.ID)
	}

	if result.Name != "test-mariadb" {
		t.Errorf("Expected Name 'test-mariadb', got %s", result.Name)
	}

	if result.Status != "creating" {
		t.Errorf("Expected Status 'creating', got %s", result.Status)
	}

	if !result.IsEncryptionEnabled {
		t.Error("Expected IsEncryptionEnabled true, got false")
	}

	assertStatus(t, resp, http.StatusOK)
}

func TestCreateMariaDB_NilRequest(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, _, err := ts.client.MariaDB.CreateMariaDB(context.Background(), nil)
	assertError(t, err, "")

	assertErrorType(t, err, &ArgError{})
}

func TestCreateMariaDB_EmptyName(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	createReq := &MariaDBCreateRequest{
		Name:       "",
		SoftwareID: 10,
		TemplateID: 20,
	}

	_, _, err := ts.client.MariaDB.CreateMariaDB(context.Background(), createReq)
	assertError(t, err, "")

	assertErrorType(t, err, &ArgError{})
}

func TestCreateMariaDB_ZeroSoftwareID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	createReq := &MariaDBCreateRequest{
		Name:       "test-mariadb",
		SoftwareID: 0,
		TemplateID: 20,
	}

	_, _, err := ts.client.MariaDB.CreateMariaDB(context.Background(), createReq)
	assertError(t, err, "")

	assertErrorType(t, err, &ArgError{})
}

func TestCreateMariaDB_ZeroTemplateID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	createReq := &MariaDBCreateRequest{
		Name:       "test-mariadb",
		SoftwareID: 10,
		TemplateID: 0,
	}

	_, _, err := ts.client.MariaDB.CreateMariaDB(context.Background(), createReq)
	assertError(t, err, "")

	assertErrorType(t, err, &ArgError{})
}

func TestCreateMariaDB_Error_Response(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Internal server error",
			"error": "Database creation failed"
		}`)
	})

	createReq := &MariaDBCreateRequest{
		Name:       "test-mariadb",
		SoftwareID: 10,
		TemplateID: 20,
	}

	_, resp, err := ts.client.MariaDB.CreateMariaDB(context.Background(), createReq)
	assertError(t, err, "")

	assertStatus(t, resp, http.StatusInternalServerError)
}

func TestMariaDBGet(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": {
				"id": 123,
				"name": "test-mariadb",
				"status": "active",
				"status_title": "Active",
				"num_instances": 1,
				"software": {
					"name": "MariaDB",
					"version": "10.6",
					"engine": "mariadb"
				},
				"master_node": {
					"node_name": "master-1",
					"instance_id": 456
				}
			}
		}`)
	})

	result, resp, err := ts.client.MariaDB.GetMariaDB(context.Background(), "123")
	assertNoError(t, err)

	assertNotNil(t, result, "Expected result")

	if result.ID != 123 {
		t.Errorf("Expected ID 123, got %d", result.ID)
	}

	if result.Name != "test-mariadb" {
		t.Errorf("Expected Name 'test-mariadb', got %s", result.Name)
	}

	if result.Status != "active" {
		t.Errorf("Expected Status 'active', got %s", result.Status)
	}

	assertStatus(t, resp, http.StatusOK)
}

func TestGetMariaDB_NotFound(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/999/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		writeJSON(w, http.StatusNotFound, `{
			"code": 404,
			"message": "Cluster not found"
		}`)
	})

	result, resp, err := ts.client.MariaDB.GetMariaDB(context.Background(), "999")
	if err != nil {
		t.Fatalf("GetMariaDB returned unexpected error: %v", err)
	}

	assertNil(t, result, "Expected nil result")

	assertStatus(t, resp, http.StatusNotFound)
}

func TestGetMariaDB_Empty_ID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, _, err := ts.client.MariaDB.GetMariaDB(context.Background(), "")
	assertError(t, err, "")

	assertErrorType(t, err, &ArgError{})
}

func TestGetMariaDB_Error_Response(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Internal server error"
		}`)
	})

	_, resp, err := ts.client.MariaDB.GetMariaDB(context.Background(), "123")
	assertError(t, err, "")

	assertStatus(t, resp, http.StatusInternalServerError)
}

func TestMariaDBExists_True(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		w.WriteHeader(http.StatusOK)
	})

	exists, resp, err := ts.client.MariaDB.MariaDBExists(context.Background(), "123")
	assertNoError(t, err)

	if !exists {
		t.Error("Expected cluster to exist")
	}

	assertStatus(t, resp, http.StatusOK)
}

func TestMariaDBExists_False(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/999/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		writeJSON(w, http.StatusNotFound, `{
			"code": 404,
			"message": "Cluster not found"
		}`)
	})

	exists, resp, err := ts.client.MariaDB.MariaDBExists(context.Background(), "999")
	if err != nil {
		t.Fatalf("MariaDBExists returned unexpected error: %v", err)
	}

	if exists {
		t.Error("Expected cluster to not exist")
	}

	assertStatus(t, resp, http.StatusNotFound)
}

func TestMariaDBExists_Empty_ID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, _, err := ts.client.MariaDB.MariaDBExists(context.Background(), "")
	assertError(t, err, "")

	assertErrorType(t, err, &ArgError{})
}

func TestMariaDBExists_Error_Response(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Internal server error"
		}`)
	})

	_, resp, err := ts.client.MariaDB.MariaDBExists(context.Background(), "123")
	assertError(t, err, "")

	assertStatus(t, resp, http.StatusInternalServerError)
}

func TestMariaDBDelete(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		w.WriteHeader(http.StatusOK)
	})

	resp, err := ts.client.MariaDB.DeleteMariaDB(context.Background(), "123")
	assertNoError(t, err)

	assertStatus(t, resp, http.StatusOK)
}

func TestDeleteMariaDB_Empty_ID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.MariaDB.DeleteMariaDB(context.Background(), "")
	assertError(t, err, "")

	assertErrorType(t, err, &ArgError{})
}

func TestDeleteMariaDB_Error_Response(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)

		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Internal server error"
		}`)
	})

	resp, err := ts.client.MariaDB.DeleteMariaDB(context.Background(), "123")
	assertError(t, err, "")

	assertStatus(t, resp, http.StatusInternalServerError)
}

// =====================================================
// Lifecycle Operations Tests
// =====================================================

func TestMariaDBShutdown(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/shutdown", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		w.WriteHeader(http.StatusOK)
	})

	resp, err := ts.client.MariaDB.ShutdownMariaDB(context.Background(), "123")
	assertNoError(t, err)

	assertStatus(t, resp, http.StatusOK)
}

func TestShutdownMariaDB_Empty_ID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.MariaDB.ShutdownMariaDB(context.Background(), "")
	assertError(t, err, "")

	assertErrorType(t, err, &ArgError{})
}

func TestShutdownMariaDB_Error_Response(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/shutdown", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)

		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Failed to shutdown cluster"
		}`)
	})

	resp, err := ts.client.MariaDB.ShutdownMariaDB(context.Background(), "123")
	assertError(t, err, "")

	assertStatus(t, resp, http.StatusInternalServerError)
}

func TestMariaDBResume(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/resume", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		w.WriteHeader(http.StatusOK)
	})

	resp, err := ts.client.MariaDB.ResumeMariaDB(context.Background(), "123")
	assertNoError(t, err)

	assertStatus(t, resp, http.StatusOK)
}

func TestResumeMariaDB_Empty_ID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.MariaDB.ResumeMariaDB(context.Background(), "")
	assertError(t, err, "")

	assertErrorType(t, err, &ArgError{})
}

func TestResumeMariaDB_Error_Response(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/resume", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)

		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Failed to resume cluster"
		}`)
	})

	resp, err := ts.client.MariaDB.ResumeMariaDB(context.Background(), "123")
	assertError(t, err, "")

	assertStatus(t, resp, http.StatusInternalServerError)
}

func TestMariaDBRestart(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/restart", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		w.WriteHeader(http.StatusOK)
	})

	resp, err := ts.client.MariaDB.RestartMariaDB(context.Background(), "123")
	assertNoError(t, err)

	assertStatus(t, resp, http.StatusOK)
}

func TestRestartMariaDB_Empty_ID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.MariaDB.RestartMariaDB(context.Background(), "")
	assertError(t, err, "")

	assertErrorType(t, err, &ArgError{})
}

func TestRestartMariaDB_Error_Response(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/restart", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)

		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Failed to restart cluster"
		}`)
	})

	resp, err := ts.client.MariaDB.RestartMariaDB(context.Background(), "123")
	assertError(t, err, "")

	assertStatus(t, resp, http.StatusInternalServerError)
}

// =====================================================
// VPC Operations Tests
// =====================================================

func TestMariaDBAttachVPC(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	// Mock VPC Get endpoint
	ts.mux.HandleFunc("/vpc/100/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": {
				"network_id": 100,
				"name": "test-vpc",
				"ipv4_cidr": "10.0.0.0/24",
				"state": "active"
			}
		}`)
	})

	// Mock VPC attach endpoint
	ts.mux.HandleFunc("/rds/cluster/123/vpc-attach/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)

		// Verify payload
		var payload AttachDetachVPCRequest
		json.NewDecoder(r.Body).Decode(&payload)

		if payload.Action != "attach" {
			t.Errorf("Expected action 'attach', got %s", payload.Action)
		}

		if len(payload.VPCs) != 1 {
			t.Errorf("Expected 1 VPC, got %d", len(payload.VPCs))
		}

		w.WriteHeader(http.StatusOK)
	})

	resp, err := ts.client.MariaDB.AttachVPC(context.Background(), "123", []string{"100"})
	assertNoError(t, err)

	assertStatus(t, resp, http.StatusOK)
}

func TestAttachVPC_Empty_ClusterID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.MariaDB.AttachVPC(context.Background(), "", []string{"100"})
	assertError(t, err, "")

	assertErrorType(t, err, &ArgError{})
}

func TestAttachVPC_EmptyVPCList(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.MariaDB.AttachVPC(context.Background(), "123", []string{})
	assertError(t, err, "")

	assertErrorType(t, err, &ArgError{})
}

func TestAttachVPC_VPCNotFound(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	// Mock VPC Get endpoint returning not found
	ts.mux.HandleFunc("/vpc/details/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		writeJSON(w, http.StatusNotFound, `{
			"code": 404,
			"message": "VPC not found"
		}`)
	})

	_, err := ts.client.MariaDB.AttachVPC(context.Background(), "123", []string{"999"})
	assertError(t, err, "")
}

func TestAttachVPC_Error_Response(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	// Mock VPC Get endpoint
	ts.mux.HandleFunc("/vpc/details/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"data": {
				"network_id": 100,
				"name": "test-vpc",
				"ipv4_cidr": "10.0.0.0/24"
			}
		}`)
	})

	// Mock VPC attach endpoint with error
	ts.mux.HandleFunc("/rds/cluster/123/vpc-attach/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Failed to attach VPC"
		}`)
	})

	resp, err := ts.client.MariaDB.AttachVPC(context.Background(), "123", []string{"100"})
	assertError(t, err, "")

	if resp != nil && resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, resp.StatusCode)
	}
}

func TestMariaDBDetachVPC(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	// Mock VPC Get endpoint
	ts.mux.HandleFunc("/vpc/100/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": {
				"network_id": 100,
				"name": "test-vpc",
				"ipv4_cidr": "10.0.0.0/24",
				"state": "active"
			}
		}`)
	})

	// Mock VPC detach endpoint
	ts.mux.HandleFunc("/rds/cluster/123/vpc-detach/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)

		// Verify payload
		var payload AttachDetachVPCRequest
		json.NewDecoder(r.Body).Decode(&payload)

		if payload.Action != "detach" {
			t.Errorf("Expected action 'detach', got %s", payload.Action)
		}

		w.WriteHeader(http.StatusOK)
	})

	resp, err := ts.client.MariaDB.DetachVPC(context.Background(), "123", []string{"100"})
	assertNoError(t, err)

	assertStatus(t, resp, http.StatusOK)
}

func TestDetachVPC_Empty_ClusterID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.MariaDB.DetachVPC(context.Background(), "", []string{"100"})
	assertError(t, err, "")

	assertErrorType(t, err, &ArgError{})
}

func TestDetachVPC_EmptyVPCList(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.MariaDB.DetachVPC(context.Background(), "123", []string{})
	assertError(t, err, "")

	assertErrorType(t, err, &ArgError{})
}

func TestDetachVPC_Error_Response(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	// Mock VPC Get endpoint
	ts.mux.HandleFunc("/vpc/details/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"data": {
				"network_id": 100,
				"name": "test-vpc",
				"ipv4_cidr": "10.0.0.0/24"
			}
		}`)
	})

	// Mock VPC detach endpoint with error
	ts.mux.HandleFunc("/rds/cluster/123/vpc-detach/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Failed to detach VPC"
		}`)
	})

	resp, err := ts.client.MariaDB.DetachVPC(context.Background(), "123", []string{"100"})
	assertError(t, err, "")

	if resp != nil && resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, resp.StatusCode)
	}
}

// =====================================================
// Public IP Operations Tests
// =====================================================

func TestMariaDBAttachPublicIP(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/public-ip-attach/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)

		// Verify payload
		var payload publicIPActionRequest
		json.NewDecoder(r.Body).Decode(&payload)

		if payload.Action != "attach" {
			t.Errorf("Expected action 'attach', got %s", payload.Action)
		}

		w.WriteHeader(http.StatusOK)
	})

	resp, err := ts.client.MariaDB.AttachPublicIP(context.Background(), "123")
	assertNoError(t, err)

	assertStatus(t, resp, http.StatusOK)
}

func TestAttachPublicIP_Empty_ID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.MariaDB.AttachPublicIP(context.Background(), "")
	assertError(t, err, "")

	assertErrorType(t, err, &ArgError{})
}

func TestAttachPublicIP_Error_Response(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/public-ip-attach/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Failed to attach public IP"
		}`)
	})

	resp, err := ts.client.MariaDB.AttachPublicIP(context.Background(), "123")
	assertError(t, err, "")

	assertStatus(t, resp, http.StatusInternalServerError)
}

func TestMariaDBDetachPublicIP(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/public-ip-detach/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)

		// Verify payload
		var payload publicIPActionRequest
		json.NewDecoder(r.Body).Decode(&payload)

		if payload.Action != "detach" {
			t.Errorf("Expected action 'detach', got %s", payload.Action)
		}

		w.WriteHeader(http.StatusOK)
	})

	resp, err := ts.client.MariaDB.DetachPublicIP(context.Background(), "123")
	assertNoError(t, err)

	assertStatus(t, resp, http.StatusOK)
}

func TestDetachPublicIP_Empty_ID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.MariaDB.DetachPublicIP(context.Background(), "")
	assertError(t, err, "")

	assertErrorType(t, err, &ArgError{})
}

func TestDetachPublicIP_Error_Response(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/public-ip-detach/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Failed to detach public IP"
		}`)
	})

	resp, err := ts.client.MariaDB.DetachPublicIP(context.Background(), "123")
	assertError(t, err, "")

	assertStatus(t, resp, http.StatusInternalServerError)
}

// =====================================================
// Parameter Group Operations Tests
// =====================================================

func TestMariaDBAttachParameterGroup(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/parameter-group/5/add", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)

		var payload ParameterGroupRequest
		json.NewDecoder(r.Body).Decode(&payload)

		if payload.Action != "add" {
			t.Errorf("Expected action 'add', got %s", payload.Action)
		}

		w.WriteHeader(http.StatusOK)
	})

	resp, err := ts.client.MariaDB.AttachParameterGroup(context.Background(), "123", 5)
	assertNoError(t, err)

	assertStatus(t, resp, http.StatusOK)
}

func TestAttachParameterGroup_Empty_ClusterID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.MariaDB.AttachParameterGroup(context.Background(), "", 5)
	assertError(t, err, "")

	assertErrorType(t, err, &ArgError{})
}

func TestAttachParameterGroup_InvalidID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.MariaDB.AttachParameterGroup(context.Background(), "123", 0)
	assertError(t, err, "")

	assertErrorType(t, err, &ArgError{})
}

func TestAttachParameterGroup_NegativeID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.MariaDB.AttachParameterGroup(context.Background(), "123", -1)
	assertError(t, err, "")

	assertErrorType(t, err, &ArgError{})
}

func TestAttachParameterGroup_Error_Response(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/parameter-group/5/add", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Failed to attach parameter group"
		}`)
	})

	resp, err := ts.client.MariaDB.AttachParameterGroup(context.Background(), "123", 5)
	assertError(t, err, "")

	assertStatus(t, resp, http.StatusInternalServerError)
}

func TestMariaDBDetachParameterGroup(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/parameter-group/5/detach", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)

		w.WriteHeader(http.StatusOK)
	})

	resp, err := ts.client.MariaDB.DetachParameterGroup(context.Background(), "123", 5)
	assertNoError(t, err)

	assertStatus(t, resp, http.StatusOK)
}

func TestDetachParameterGroup_Empty_ClusterID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.MariaDB.DetachParameterGroup(context.Background(), "", 5)
	assertError(t, err, "")

	assertErrorType(t, err, &ArgError{})
}

func TestDetachParameterGroup_InvalidID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.MariaDB.DetachParameterGroup(context.Background(), "123", 0)
	assertError(t, err, "")

	assertErrorType(t, err, &ArgError{})
}

func TestDetachParameterGroup_Error_Response(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/parameter-group/5/detach", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Failed to detach parameter group"
		}`)
	})

	resp, err := ts.client.MariaDB.DetachParameterGroup(context.Background(), "123", 5)
	assertError(t, err, "")

	assertStatus(t, resp, http.StatusInternalServerError)
}

// =====================================================
// Upgrade Operations Tests
// =====================================================

func TestMariaDBUpgradePlan(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/rds-upgrade/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)

		var payload UpgradePlanRequest
		json.NewDecoder(r.Body).Decode(&payload)

		if payload.TemplateID != 42 {
			t.Errorf("Expected template ID 42, got %d", payload.TemplateID)
		}

		w.WriteHeader(http.StatusOK)
	})

	resp, err := ts.client.MariaDB.UpgradePlan(context.Background(), "123", 42)
	assertNoError(t, err)

	assertStatus(t, resp, http.StatusOK)
}

func TestUpgradePlan_Empty_ClusterID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.MariaDB.UpgradePlan(context.Background(), "", 42)
	assertError(t, err, "")

	assertErrorType(t, err, &ArgError{})
}

func TestUpgradePlan_InvalidTemplateID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.MariaDB.UpgradePlan(context.Background(), "123", 0)
	assertError(t, err, "")

	assertErrorType(t, err, &ArgError{})
}

func TestUpgradePlan_NegativeTemplateID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.MariaDB.UpgradePlan(context.Background(), "123", -1)
	assertError(t, err, "")

	assertErrorType(t, err, &ArgError{})
}

func TestUpgradePlan_Error_Response(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/rds-upgrade/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Failed to upgrade plan"
		}`)
	})

	resp, err := ts.client.MariaDB.UpgradePlan(context.Background(), "123", 42)
	assertError(t, err, "")

	assertStatus(t, resp, http.StatusInternalServerError)
}

func TestMariaDBExpandDisk(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/disk-upgrade/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)

		var payload DiskUpgradeRequest
		json.NewDecoder(r.Body).Decode(&payload)

		if payload.Size != 50 {
			t.Errorf("Expected size 50, got %d", payload.Size)
		}

		w.WriteHeader(http.StatusOK)
	})

	resp, err := ts.client.MariaDB.ExpandDisk(context.Background(), "123", 50)
	assertNoError(t, err)

	assertStatus(t, resp, http.StatusOK)
}

func TestExpandDisk_Empty_ClusterID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.MariaDB.ExpandDisk(context.Background(), "", 50)
	assertError(t, err, "")

	assertErrorType(t, err, &ArgError{})
}

func TestExpandDisk_ZeroSize(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.MariaDB.ExpandDisk(context.Background(), "123", 0)
	assertError(t, err, "")

	assertErrorType(t, err, &ArgError{})
}

func TestExpandDisk_NegativeSize(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.MariaDB.ExpandDisk(context.Background(), "123", -10)
	assertError(t, err, "")

	assertErrorType(t, err, &ArgError{})
}

func TestExpandDisk_Error_Response(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/disk-upgrade/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Failed to expand disk"
		}`)
	})

	resp, err := ts.client.MariaDB.ExpandDisk(context.Background(), "123", 50)
	assertError(t, err, "")

	assertStatus(t, resp, http.StatusInternalServerError)
}

// =====================================================
// Helper Operations Tests
// =====================================================

func TestMariaDBExpandVPCList(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	// Mock VPC Get endpoint for both VPCs
	ts.mux.HandleFunc("/vpc/100/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"data": {
				"network_id": 100,
				"name": "vpc-100",
				"ipv4_cidr": "10.0.0.0/24"
			}
		}`)
	})

	ts.mux.HandleFunc("/vpc/200/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"data": {
				"network_id": 200,
				"name": "vpc-200",
				"ipv4_cidr": "10.0.1.0/24"
			}
		}`)
	})

	vpcList, err := ts.client.MariaDB.ExpandVPCList(context.Background(), []string{"100", "200"})
	assertNoError(t, err)

	if len(vpcList) != 2 {
		t.Fatalf("Expected 2 VPCs, got %d", len(vpcList))
	}

	if vpcList[0].NetworkID != "100" {
		t.Errorf("Expected NetworkID '100', got %s", vpcList[0].NetworkID)
	}

	if vpcList[0].VPCName != "vpc-100" {
		t.Errorf("Expected VPCName 'vpc-100', got %s", vpcList[0].VPCName)
	}

	if vpcList[1].NetworkID != "200" {
		t.Errorf("Expected NetworkID '200', got %s", vpcList[1].NetworkID)
	}
}

func TestMariaDBExpandVPCList_EmptyList(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.MariaDB.ExpandVPCList(context.Background(), []string{})
	assertError(t, err, "")

	assertErrorType(t, err, &ArgError{})
}

func TestMariaDBExpandVPCList_VPCNotFound(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/vpc/999/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusNotFound, `{
			"code": 404,
			"message": "VPC not found"
		}`)
	})

	_, err := ts.client.MariaDB.ExpandVPCList(context.Background(), []string{"999"})
	assertError(t, err, "")
}

func TestMariaDBExpandVPCList_SkipEmptyIDs(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/vpc/100/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"data": {
				"network_id": 100,
				"name": "vpc-100",
				"ipv4_cidr": "10.0.0.0/24"
			}
		}`)
	})

	ts.mux.HandleFunc("/vpc/200/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"data": {
				"network_id": 200,
				"name": "vpc-200",
				"ipv4_cidr": "10.1.0.0/24"
			}
		}`)
	})

	vpcList, err := ts.client.MariaDB.ExpandVPCList(context.Background(), []string{"100", "200"})
	assertNoError(t, err)

	if len(vpcList) != 2 {
		t.Fatalf("Expected 2 VPCs, got %d", len(vpcList))
	}

	if vpcList[0].NetworkID != "100" {
		t.Errorf("Expected NetworkID '100', got %s", vpcList[0].NetworkID)
	}

	if vpcList[0].VPCName != "vpc-100" {
		t.Errorf("Expected VPCName 'vpc-100', got %s", vpcList[0].VPCName)
	}

	if vpcList[1].NetworkID != "200" {
		t.Errorf("Expected NetworkID '200', got %s", vpcList[1].NetworkID)
	}
}

func TestExpandVPCList_EmptyList(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.MariaDB.ExpandVPCList(context.Background(), []string{})
	assertError(t, err, "")

	assertErrorType(t, err, &ArgError{})
}

func TestExpandVPCList_VPCNotFound(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/vpc/details/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusNotFound, `{
			"code": 404,
			"message": "VPC not found"
		}`)
	})

	_, err := ts.client.MariaDB.ExpandVPCList(context.Background(), []string{"999"})
	assertError(t, err, "")
}

func TestExpandVPCList_SkipEmptyIDs(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/vpc/100/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"data": {
				"network_id": 100,
				"name": "vpc-100",
				"ipv4_cidr": "10.0.0.0/24"
			}
		}`)
	})

	vpcList, err := ts.client.MariaDB.ExpandVPCList(context.Background(), []string{"100", ""})
	assertNoError(t, err)

	if len(vpcList) != 1 {
		t.Fatalf("Expected 1 VPC (empty ID skipped), got %d", len(vpcList))
	}
}

func TestExpandVPCList_AllEmptyIDs(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.MariaDB.ExpandVPCList(context.Background(), []string{"", "", ""})
	assertError(t, err, "")
}

// =====================================================
// Edge Cases & Error Handling Tests
// =====================================================

func TestCreateMariaDB_WithVPCs(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		var createReq MariaDBCreateRequest
		json.NewDecoder(r.Body).Decode(&createReq)

		if len(createReq.VPCs) != 1 {
			t.Errorf("Expected 1 VPC, got %d", len(createReq.VPCs))
		}

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"data": {
				"id": 123,
				"name": "test-mariadb",
				"status": "creating"
			}
		}`)
	})

	createReq := &MariaDBCreateRequest{
		Name:       "test-mariadb",
		SoftwareID: 10,
		TemplateID: 20,
		VPCs: []VPCMetadata{
			{
				NetworkID: "100",
				VPCName:   "test-vpc",
				IPv4CIDR:  "10.0.0.0/24",
			},
		},
	}

	result, _, err := ts.client.MariaDB.CreateMariaDB(context.Background(), createReq)
	assertNoError(t, err)

	if result.ID != 123 {
		t.Errorf("Expected ID 123, got %d", result.ID)
	}
}

func TestCreateMariaDB_WithEncryption(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/", func(w http.ResponseWriter, r *http.Request) {
		var createReq MariaDBCreateRequest
		json.NewDecoder(r.Body).Decode(&createReq)

		if !createReq.IsEncryptionEnabled {
			t.Error("Expected encryption to be enabled")
		}

		if createReq.EncryptionPassphrase == "" {
			t.Error("Expected encryption passphrase")
		}

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"data": {
				"id": 123,
				"name": "test-mariadb",
				"isEncryptionEnabled": true
			}
		}`)
	})

	createReq := &MariaDBCreateRequest{
		Name:                 "test-mariadb",
		SoftwareID:           10,
		TemplateID:           20,
		IsEncryptionEnabled:  true,
		EncryptionPassphrase: "secret123",
	}

	result, _, err := ts.client.MariaDB.CreateMariaDB(context.Background(), createReq)
	assertNoError(t, err)

	if !result.IsEncryptionEnabled {
		t.Error("Expected encryption enabled in response")
	}
}

func TestMariaDBContextCancellation(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, _, err := ts.client.MariaDB.GetMariaDB(ctx, "123")
	assertError(t, err, "")
}

func TestMariaDBMultipleVPCAttach(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	// Mock multiple VPC Get endpoints
	ts.mux.HandleFunc("/vpc/100/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, `{"code": 200, "data": {"network_id": 100, "name": "vpc-100", "ipv4_cidr": "10.0.0.0/24"}}`)
	})
	ts.mux.HandleFunc("/vpc/200/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, `{"code": 200, "data": {"network_id": 200, "name": "vpc-200", "ipv4_cidr": "10.0.1.0/24"}}`)
	})
	ts.mux.HandleFunc("/vpc/300/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, `{"code": 200, "data": {"network_id": 300, "name": "vpc-300", "ipv4_cidr": "10.0.2.0/24"}}`)
	})

	ts.mux.HandleFunc("/rds/cluster/123/vpc-attach/", func(w http.ResponseWriter, r *http.Request) {
		var payload AttachDetachVPCRequest
		json.NewDecoder(r.Body).Decode(&payload)

		if len(payload.VPCs) != 3 {
			t.Errorf("Expected 3 VPCs, got %d", len(payload.VPCs))
		}

		w.WriteHeader(http.StatusOK)
	})

	_, err := ts.client.MariaDB.AttachVPC(context.Background(), "123", []string{"100", "200", "300"})
	assertNoError(t, err)
}

// ============================================
// ADDITIONAL COMPREHENSIVE ERROR TESTS
// ============================================

func TestMariaDBShutdown_Conflict(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/shutdown", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusConflict, `{
			"code": 409,
			"message": "Cluster already stopped"
		}`)
	})

	resp, err := ts.client.MariaDB.ShutdownMariaDB(context.Background(), "123")
	assertError(t, err, "")

	assertStatus(t, resp, http.StatusConflict)
}

func TestMariaDBResume_Conflict(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/resume", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusConflict, `{
			"code": 409,
			"message": "Cluster already running"
		}`)
	})

	resp, err := ts.client.MariaDB.ResumeMariaDB(context.Background(), "123")
	assertError(t, err, "")

	assertStatus(t, resp, http.StatusConflict)
}

func TestMariaDBRestart_Forbidden(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/restart", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusForbidden, `{
			"code": 403,
			"message": "Insufficient permissions"
		}`)
	})

	resp, err := ts.client.MariaDB.RestartMariaDB(context.Background(), "123")
	assertError(t, err, "")

	assertStatus(t, resp, http.StatusForbidden)
}

func TestMariaDBAttachVPC_BadRequest(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/vpc/100/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, `{"code": 200, "data": {"network_id": 100, "name": "vpc-100", "ipv4_cidr": "10.0.0.0/24"}}`)
	})

	ts.mux.HandleFunc("/rds/cluster/123/vpc-attach/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusBadRequest, `{
			"code": 400,
			"message": "Invalid VPC configuration"
		}`)
	})

	resp, err := ts.client.MariaDB.AttachVPC(context.Background(), "123", []string{"100"})
	assertError(t, err, "")

	assertStatus(t, resp, http.StatusBadRequest)
}

func TestMariaDBDetachVPC_Conflict(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/vpc/100/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, `{"code": 200, "data": {"network_id": 100, "name": "vpc-100", "ipv4_cidr": "10.0.0.0/24"}}`)
	})

	ts.mux.HandleFunc("/rds/cluster/123/vpc-detach/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusConflict, `{
			"code": 409,
			"message": "VPC not attached to cluster"
		}`)
	})

	resp, err := ts.client.MariaDB.DetachVPC(context.Background(), "123", []string{"100"})
	assertError(t, err, "")

	assertStatus(t, resp, http.StatusConflict)
}

func TestMariaDBAttachPublicIP_Forbidden(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/public-ip-attach/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusForbidden, `{
			"code": 403,
			"message": "Insufficient permissions"
		}`)
	})

	resp, err := ts.client.MariaDB.AttachPublicIP(context.Background(), "123")
	assertError(t, err, "")

	assertStatus(t, resp, http.StatusForbidden)
}

func TestMariaDBDetachPublicIP_NotFound(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/public-ip-detach/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusNotFound, `{
			"code": 404,
			"message": "Public IP not attached"
		}`)
	})

	resp, err := ts.client.MariaDB.DetachPublicIP(context.Background(), "123")
	assertError(t, err, "")

	assertStatus(t, resp, http.StatusNotFound)
}

func TestMariaDBAttachParameterGroup_Conflict(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/parameter-group/456/add", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusConflict, `{
			"code": 409,
			"message": "Parameter group already attached"
		}`)
	})

	resp, err := ts.client.MariaDB.AttachParameterGroup(context.Background(), "123", 456)
	assertError(t, err, "")

	assertStatus(t, resp, http.StatusConflict)
}

func TestMariaDBDetachParameterGroup_NotFound(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/123/parameter-group/456/detach", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusNotFound, `{
			"code": 404,
			"message": "Parameter group not attached"
		}`)
	})

	resp, err := ts.client.MariaDB.DetachParameterGroup(context.Background(), "123", 456)
	assertError(t, err, "")

	assertStatus(t, resp, http.StatusNotFound)
}

func TestMariaDBContextCancelled_Create(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, _, err := ts.client.MariaDB.CreateMariaDB(ctx, &MariaDBCreateRequest{})
	if err == nil {
		t.Fatal("Expected error for cancelled context")
	}
}

func TestMariaDBContextCancelled_Delete(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := ts.client.MariaDB.DeleteMariaDB(ctx, "123")
	if err == nil {
		t.Fatal("Expected error for cancelled context")
	}
}

// Phase 2: Response Parsing & Edge Case Tests

func TestMariaDBCreateCluster_MalformedJSON(t *testing.T) {
	server := newMalformedJSONServer(t)
	defer server.Close()

	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	_, _, err := client.MariaDB.CreateMariaDB(context.Background(), &MariaDBCreateRequest{
		Name:       "test-cluster",
		SoftwareID: 1,
		TemplateID: 1,
	})

	if err == nil {
		t.Error("Expected error for malformed JSON")
	}
}

func TestMariaDBGetCluster_MissingRequiredFields(t *testing.T) {
	server := newMissingFieldServer(t, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data": map[string]interface{}{
			// Missing "cluster_id" field
			"name": "test-cluster",
		},
	})
	defer server.Close()

	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	resp, _, err := client.MariaDB.GetMariaDB(context.Background(), "cluster-123")

	// Should handle missing fields gracefully
	if resp == nil && err == nil {
		t.Error("Expected response or error handling")
	}
}

func TestMariaDBGetCluster_NullFieldValues(t *testing.T) {
	server := newNullFieldServer(t, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data": map[string]interface{}{
			"cluster_id":      "cluster-123",
			"name":            "test-cluster",
			"connection_info": nil, // Null value
			"backup_config":   nil,
			"vpc_id":          nil,
		},
	})
	defer server.Close()

	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	resp, _, err := client.MariaDB.GetMariaDB(context.Background(), "cluster-123")

	// Should handle null fields without panic
	if resp == nil && err == nil {
		t.Error("Expected response or error for null fields")
	}
}

func TestMariaDBGetCluster_InvalidVersionField(t *testing.T) {
	server := newInvalidFieldTypeServer(t, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data": map[string]interface{}{
			"cluster_id": "cluster-123",
			"name":       "test-cluster",
			"version":    789, // Should be string, not int
		},
	})
	defer server.Close()

	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	resp, _, err := client.MariaDB.GetMariaDB(context.Background(), "cluster-123")

	// Should handle wrong type gracefully
	if resp == nil && err == nil {
		t.Error("Expected response or error for invalid field type")
	}
}
