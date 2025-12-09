package goe2e

import (
	"context"
	"net/http"
	"testing"
)

// ============================================
// SUCCESS PATH TESTS
// ============================================

func TestCreateMySQLCluster(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testURLPath(t, r, "/rds/cluster/")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		writeJSON(w, http.StatusCreated, `{
			"code": 201,
			"message": "Cluster created successfully",
			"data": {
				"id": 12345
			}
		}`)
	})

	// Mock GetCluster call that follows CreateCluster
	ts.mux.HandleFunc("/rds/cluster/12345/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": {
				"id": 12345,
				"name": "test-mysql-cluster",
				"status": "CREATING",
				"software": {
					"name": "MySQL",
					"version": "8.0",
					"engine": "mysql"
				},
				"master_node": {
					"node_name": "master-node",
					"instance_id": 1,
					"cluster_id": 12345,
					"database": {
						"id": 1,
						"username": "admin",
						"database": "testdb",
						"pg_detail": {
							"pg_id": 0
						}
					},
					"disk": "100GB"
				},
				"isEncryptionEnabled": false
			}
		}`)
	})

	createReq := &MySQLClusterCreateRequest{
		Name:       "test-mysql-cluster",
		SoftwareID: 1,
		TemplateID: 10,
		Database: DBConfig{
			User:        "admin",
			Password:    "password123",
			Name:        "testdb",
			DBaaSNumber: 1,
		},
		PublicIPRequired: true,
		Group:            "default",
	}

	result, _, err := ts.client.DBaaSMySQL.CreateCluster(context.Background(), createReq)
	assertNoError(t, err)
	assertNotNil(t, result, "Expected result")

	if result.ID != 12345 {
		t.Errorf("Expected ID 12345, got %d", result.ID)
	}
	if result.Name != "test-mysql-cluster" {
		t.Errorf("Expected Name test-mysql-cluster, got %s", result.Name)
	}
}

func TestGetMySQLCluster(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/12345/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/rds/cluster/12345/")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": {
				"id": 12345,
				"name": "test-mysql-cluster",
				"status": "ACTIVE",
				"software": {
					"name": "MySQL",
					"version": "8.0",
					"engine": "mysql"
				},
				"master_node": {
					"node_name": "master-node",
					"instance_id": 1,
					"cluster_id": 12345,
					"database": {
						"id": 1,
						"username": "admin",
						"database": "testdb",
						"pg_detail": {
							"pg_id": 0
						}
					},
					"disk": "100GB"
				},
				"isEncryptionEnabled": false
			}
		}`)
	})

	result, _, err := ts.client.DBaaSMySQL.GetCluster(context.Background(), "12345")
	assertNoError(t, err)
	assertNotNil(t, result, "Expected result")

	if result.ID != 12345 {
		t.Errorf("Expected ID 12345, got %d", result.ID)
	}
	if result.Status != "ACTIVE" {
		t.Errorf("Expected Status ACTIVE, got %s", result.Status)
	}
}

func TestGetMySQLCluster_NotFound(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/nonexistent/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		writeJSON(w, http.StatusNotFound, `{
			"code": 404,
			"message": "Cluster not found"
		}`)
	})

	result, _, err := ts.client.DBaaSMySQL.GetCluster(context.Background(), "nonexistent")
	assertNoError(t, err)
	assertNil(t, result, "Expected nil result for 404")
}

func TestDeleteMySQLCluster(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/12345/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		testURLPath(t, r, "/rds/cluster/12345/")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Cluster deleted successfully"
		}`)
	})

	_, err := ts.client.DBaaSMySQL.DeleteCluster(context.Background(), "12345")
	assertNoError(t, err)
}

func TestStartMySQLCluster(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/12345/resume", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/12345/resume")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Cluster started successfully"
		}`)
	})

	_, err := ts.client.DBaaSMySQL.StartCluster(context.Background(), "12345")
	assertNoError(t, err)
}

func TestStopMySQLCluster(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/12345/shutdown", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/12345/shutdown")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Cluster stopped successfully"
		}`)
	})

	_, err := ts.client.DBaaSMySQL.StopCluster(context.Background(), "12345")
	assertNoError(t, err)
}

func TestRestartMySQLCluster(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/12345/restart", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/12345/restart")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Cluster restarted successfully"
		}`)
	})

	_, err := ts.client.DBaaSMySQL.RestartCluster(context.Background(), "12345")
	assertNoError(t, err)
}

func TestMySQLAttachVPC(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/12345/vpc-attach/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/12345/vpc-attach/")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "VPC attached successfully"
		}`)
	})

	attachReq := &MySQLVPCAttachRequest{
		Action: "attach",
		VPCs: []VPCMetadata{
			{
				VPCName:   "test-vpc",
				IPv4CIDR:  "10.0.0.0/24",
				NetworkID: "123",
			},
		},
	}

	_, err := ts.client.DBaaSMySQL.AttachVPC(context.Background(), "12345", attachReq)

	assertNoError(t, err)
}

func TestMySQLDetachVPC(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/12345/vpc-detach/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/12345/vpc-detach/")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "VPC detached successfully"
		}`)
	})

	detachReq := &MySQLVPCDetachRequest{
		Action: "detach",
		VPCs: []VPCMetadata{
			{
				NetworkID: "123",
			},
		},
	}

	_, err := ts.client.DBaaSMySQL.DetachVPC(context.Background(), "12345", detachReq)

	assertNoError(t, err)
}

func TestAttachParameterGroup(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/12345/parameter-group/456/add", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/12345/parameter-group/456/add")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Parameter group attached successfully"
		}`)
	})

	_, err := ts.client.DBaaSMySQL.AttachParameterGroup(context.Background(), "12345", "456")

	assertNoError(t, err)
}

func TestDetachParameterGroup(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/12345/parameter-group/456/detach", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/12345/parameter-group/456/detach")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Parameter group detached successfully"
		}`)
	})

	_, err := ts.client.DBaaSMySQL.DetachParameterGroup(context.Background(), "12345", "456")

	assertNoError(t, err)
}

func TestAttachPublicIP(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/12345/public-ip-attach/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/12345/public-ip-attach/")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Public IP attached successfully"
		}`)
	})

	_, err := ts.client.DBaaSMySQL.AttachPublicIP(context.Background(), "12345")

	assertNoError(t, err)
}

func TestDetachPublicIP(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/12345/public-ip-detach/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/12345/public-ip-detach/")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Public IP detached successfully"
		}`)
	})

	_, err := ts.client.DBaaSMySQL.DetachPublicIP(context.Background(), "12345")

	assertNoError(t, err)
}

func TestMySQLUpgradePlan(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/12345/rds-upgrade/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/12345/rds-upgrade/")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Plan upgraded successfully"
		}`)
	})

	upgradeReq := &MySQLPlanUpgradeRequest{
		TemplateID: 20,
	}

	_, err := ts.client.DBaaSMySQL.UpgradePlan(context.Background(), "12345", upgradeReq)

	assertNoError(t, err)
}

func TestExpandDisk(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/12345/disk-upgrade/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/rds/cluster/12345/disk-upgrade/")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Disk expanded successfully"
		}`)
	})

	expandReq := &DiskExpansionRequest{
		Size: 50,
	}

	_, err := ts.client.DBaaSMySQL.ExpandDisk(context.Background(), "12345", expandReq)

	assertNoError(t, err)
}

// ============================================
// EDGE CASE / VALIDATION TESTS
// ============================================

func TestCreateMySQLCluster_NilRequest(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, _, err := ts.client.DBaaSMySQL.CreateCluster(context.Background(), nil)
	assertError(t, err, "")
}

func TestGetMySQLCluster_EmptyID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, _, err := ts.client.DBaaSMySQL.GetCluster(context.Background(), "")
	assertError(t, err, "")
}

func TestDeleteMySQLCluster_EmptyID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.DBaaSMySQL.DeleteCluster(context.Background(), "")
	assertError(t, err, "")
}

func TestStartMySQLCluster_EmptyID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.DBaaSMySQL.StartCluster(context.Background(), "")
	assertError(t, err, "")
}

func TestStopMySQLCluster_EmptyID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.DBaaSMySQL.StopCluster(context.Background(), "")
	assertError(t, err, "")
}

func TestRestartMySQLCluster_EmptyID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.DBaaSMySQL.RestartCluster(context.Background(), "")
	assertError(t, err, "")
}

func TestMySQLAttachVPC_EmptyID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	attachReq := &MySQLVPCAttachRequest{
		Action: "attach",
		VPCs:   []VPCMetadata{},
	}

	_, err := ts.client.DBaaSMySQL.AttachVPC(context.Background(), "", attachReq)
	assertError(t, err, "")
}

func TestMySQLAttachVPC_NilRequest(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.DBaaSMySQL.AttachVPC(context.Background(), "12345", nil)
	assertError(t, err, "")
}

func TestMySQLDetachVPC_EmptyID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	detachReq := &MySQLVPCDetachRequest{
		Action: "detach",
		VPCs:   []VPCMetadata{},
	}

	_, err := ts.client.DBaaSMySQL.DetachVPC(context.Background(), "", detachReq)
	assertError(t, err, "")
}

func TestMySQLDetachVPC_NilRequest(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.DBaaSMySQL.DetachVPC(context.Background(), "12345", nil)
	assertError(t, err, "")
}

func TestAttachParameterGroup_EmptyClusterID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.DBaaSMySQL.AttachParameterGroup(context.Background(), "", "456")
	assertError(t, err, "")
}

func TestAttachParameterGroup_EmptyPGID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.DBaaSMySQL.AttachParameterGroup(context.Background(), "12345", "")
	assertError(t, err, "")
}

func TestDetachParameterGroup_EmptyClusterID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.DBaaSMySQL.DetachParameterGroup(context.Background(), "", "456")
	assertError(t, err, "")
}

func TestDetachParameterGroup_EmptyPGID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.DBaaSMySQL.DetachParameterGroup(context.Background(), "12345", "")
	assertError(t, err, "")
}

func TestAttachPublicIP_EmptyID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.DBaaSMySQL.AttachPublicIP(context.Background(), "")
	assertError(t, err, "")
}

func TestDetachPublicIP_EmptyID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.DBaaSMySQL.DetachPublicIP(context.Background(), "")
	assertError(t, err, "")
}

func TestMySQLUpgradePlan_EmptyID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	upgradeReq := &MySQLPlanUpgradeRequest{
		TemplateID: 20,
	}

	_, err := ts.client.DBaaSMySQL.UpgradePlan(context.Background(), "", upgradeReq)
	assertError(t, err, "")
}

func TestMySQLUpgradePlan_NilRequest(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.DBaaSMySQL.UpgradePlan(context.Background(), "12345", nil)
	assertError(t, err, "")
}

func TestExpandDisk_EmptyID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	expandReq := &DiskExpansionRequest{
		Size: 50,
	}

	_, err := ts.client.DBaaSMySQL.ExpandDisk(context.Background(), "", expandReq)
	assertError(t, err, "")
}

func TestExpandDisk_NilRequest(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.DBaaSMySQL.ExpandDisk(context.Background(), "12345", nil)
	assertError(t, err, "")
}

// ============================================
// ERROR RESPONSE TESTS
// ============================================

func TestCreateMySQLCluster_ErrorResponse(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Internal server error"
		}`)
	})

	createReq := &MySQLClusterCreateRequest{
		Name:       "test-cluster",
		SoftwareID: 1,
		TemplateID: 10,
	}

	_, _, err := ts.client.DBaaSMySQL.CreateCluster(context.Background(), createReq)
	assertError(t, err, "")
}

func TestGetMySQLCluster_ErrorResponse(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/12345/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Internal server error"
		}`)
	})

	_, _, err := ts.client.DBaaSMySQL.GetCluster(context.Background(), "12345")
	assertError(t, err, "")
}

func TestDeleteMySQLCluster_ErrorResponse(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/12345/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Internal server error"
		}`)
	})

	_, err := ts.client.DBaaSMySQL.DeleteCluster(context.Background(), "12345")
	assertError(t, err, "")
}

func TestStartMySQLCluster_ErrorResponse(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/12345/resume", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Internal server error"
		}`)
	})

	_, err := ts.client.DBaaSMySQL.StartCluster(context.Background(), "12345")
	assertError(t, err, "")
}

func TestStopMySQLCluster_ErrorResponse(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/12345/shutdown", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Internal server error"
		}`)
	})

	_, err := ts.client.DBaaSMySQL.StopCluster(context.Background(), "12345")
	assertError(t, err, "")
}

func TestRestartMySQLCluster_ErrorResponse(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/12345/restart", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Internal server error"
		}`)
	})

	_, err := ts.client.DBaaSMySQL.RestartCluster(context.Background(), "12345")
	assertError(t, err, "")
}

func TestMySQLAttachVPC_ErrorResponse(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/12345/vpc-attach/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Internal server error"
		}`)
	})

	attachReq := &MySQLVPCAttachRequest{
		Action: "attach",
		VPCs:   []VPCMetadata{},
	}

	_, err := ts.client.DBaaSMySQL.AttachVPC(context.Background(), "12345", attachReq)
	assertError(t, err, "")
}

func TestMySQLDetachVPC_ErrorResponse(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/12345/vpc-detach/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Internal server error"
		}`)
	})

	detachReq := &MySQLVPCDetachRequest{
		Action: "detach",
		VPCs:   []VPCMetadata{},
	}

	_, err := ts.client.DBaaSMySQL.DetachVPC(context.Background(), "12345", detachReq)
	assertError(t, err, "")
}

func TestAttachParameterGroup_ErrorResponse(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/12345/parameter-group/456/add", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Internal server error"
		}`)
	})

	_, err := ts.client.DBaaSMySQL.AttachParameterGroup(context.Background(), "12345", "456")
	assertError(t, err, "")
}

func TestDetachParameterGroup_ErrorResponse(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/12345/parameter-group/456/detach", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Internal server error"
		}`)
	})

	_, err := ts.client.DBaaSMySQL.DetachParameterGroup(context.Background(), "12345", "456")
	assertError(t, err, "")
}

func TestAttachPublicIP_ErrorResponse(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/12345/public-ip-attach/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Internal server error"
		}`)
	})

	_, err := ts.client.DBaaSMySQL.AttachPublicIP(context.Background(), "12345")
	assertError(t, err, "")
}

func TestDetachPublicIP_ErrorResponse(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/12345/public-ip-detach/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Internal server error"
		}`)
	})

	_, err := ts.client.DBaaSMySQL.DetachPublicIP(context.Background(), "12345")
	assertError(t, err, "")
}

func TestMySQLUpgradePlan_ErrorResponse(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/12345/rds-upgrade/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Internal server error"
		}`)
	})

	upgradeReq := &MySQLPlanUpgradeRequest{
		TemplateID: 20,
	}

	_, err := ts.client.DBaaSMySQL.UpgradePlan(context.Background(), "12345", upgradeReq)
	assertError(t, err, "")
}

func TestExpandDisk_ErrorResponse(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/12345/disk-upgrade/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Internal server error"
		}`)
	})

	expandReq := &DiskExpansionRequest{
		Size: 50,
	}

	_, err := ts.client.DBaaSMySQL.ExpandDisk(context.Background(), "12345", expandReq)
	assertError(t, err, "")
}

// ============================================
// ADDITIONAL COVERAGE TESTS
// ============================================

func TestCreateMySQLCluster_MissingIDInResponse(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusCreated, `{
			"code": 201,
			"message": "Cluster created successfully",
			"data": {
				"name": "test-cluster"
			}
		}`)
	})

	createReq := &MySQLClusterCreateRequest{
		Name:       "test-cluster",
		SoftwareID: 1,
		TemplateID: 10,
	}

	_, _, err := ts.client.DBaaSMySQL.CreateCluster(context.Background(), createReq)
	assertError(t, err, "")
}

// ============================================
// COMPREHENSIVE ERROR RESPONSE TESTS
// ============================================

func TestCreateMySQLCluster_BadRequest(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusBadRequest, `{
			"code": 400,
			"message": "Invalid request parameters"
		}`)
	})

	createReq := &MySQLClusterCreateRequest{
		Name:       "test-cluster",
		SoftwareID: 1,
		TemplateID: 10,
	}

	_, _, err := ts.client.DBaaSMySQL.CreateCluster(context.Background(), createReq)
	assertError(t, err, "")
}

func TestCreateMySQLCluster_Conflict(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusConflict, `{
			"code": 409,
			"message": "Cluster name already exists"
		}`)
	})

	createReq := &MySQLClusterCreateRequest{
		Name:       "existing-cluster",
		SoftwareID: 1,
		TemplateID: 10,
	}

	_, _, err := ts.client.DBaaSMySQL.CreateCluster(context.Background(), createReq)
	assertError(t, err, "")
}

func TestGetMySQLCluster_Conflict(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/12345/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusConflict, `{
			"code": 409,
			"message": "Cluster is in invalid state"
		}`)
	})

	_, _, err := ts.client.DBaaSMySQL.GetCluster(context.Background(), "12345")
	assertError(t, err, "")
}

func TestDeleteMySQLCluster_Conflict(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/12345/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusConflict, `{
			"code": 409,
			"message": "Cluster still has attached resources"
		}`)
	})

	_, err := ts.client.DBaaSMySQL.DeleteCluster(context.Background(), "12345")
	assertError(t, err, "")
}

func TestStartMySQLCluster_Conflict(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/12345/resume", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusConflict, `{
			"code": 409,
			"message": "Cluster already in running state"
		}`)
	})

	_, err := ts.client.DBaaSMySQL.StartCluster(context.Background(), "12345")
	assertError(t, err, "")
}

func TestMySQLAttachVPC_BadRequest(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/12345/vpc-attach/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusBadRequest, `{
			"code": 400,
			"message": "Invalid VPC CIDR"
		}`)
	})

	attachReq := &MySQLVPCAttachRequest{
		Action: "attach",
		VPCs:   []VPCMetadata{},
	}

	_, err := ts.client.DBaaSMySQL.AttachVPC(context.Background(), "12345", attachReq)
	assertError(t, err, "")
}

func TestMySQLDetachVPC_NotFound(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/12345/vpc-detach/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusNotFound, `{
			"code": 404,
			"message": "VPC not attached to cluster"
		}`)
	})

	detachReq := &MySQLVPCDetachRequest{
		Action: "detach",
		VPCs:   []VPCMetadata{},
	}

	_, err := ts.client.DBaaSMySQL.DetachVPC(context.Background(), "12345", detachReq)
	assertError(t, err, "")
}

func TestAttachParameterGroup_NotFound(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/12345/parameter-group/999/add", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusNotFound, `{
			"code": 404,
			"message": "Parameter group not found"
		}`)
	})

	_, err := ts.client.DBaaSMySQL.AttachParameterGroup(context.Background(), "12345", "999")
	assertError(t, err, "")
}

func TestMySQLUpgradePlan_Forbidden(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/12345/rds-upgrade/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusForbidden, `{
			"code": 403,
			"message": "Insufficient permissions to upgrade plan"
		}`)
	})

	upgradeReq := &MySQLPlanUpgradeRequest{
		TemplateID: 20,
	}

	_, err := ts.client.DBaaSMySQL.UpgradePlan(context.Background(), "12345", upgradeReq)
	assertError(t, err, "")
}

func TestExpandDisk_Forbidden(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/12345/disk-upgrade/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusForbidden, `{
			"code": 403,
			"message": "Insufficient permissions to expand disk"
		}`)
	})

	expandReq := &DiskExpansionRequest{
		Size: 50,
	}

	_, err := ts.client.DBaaSMySQL.ExpandDisk(context.Background(), "12345", expandReq)
	assertError(t, err, "")
}

func TestCreateMySQLCluster_NilRequest_2(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, _, err := ts.client.DBaaSMySQL.CreateCluster(context.Background(), nil)
	assertError(t, err, "")
}

// Phase 2: Response Parsing & Edge Case Tests

func TestCreateCluster_MalformedJSON(t *testing.T) {
	server := newMalformedJSONServer(t)
	defer server.Close()

	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	_, _, err := client.DBaaSMySQL.CreateCluster(context.Background(), &MySQLClusterCreateRequest{})

	if err == nil {
		t.Error("Expected error for malformed JSON")
	}
}

func TestGetCluster_MissingRequiredFields(t *testing.T) {
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
	resp, _, err := client.DBaaSMySQL.GetCluster(context.Background(), "cluster-123")

	// Should handle missing fields gracefully
	if resp == nil && err == nil {
		t.Error("Expected response or error handling")
	}
}

func TestGetCluster_NullFieldValues(t *testing.T) {
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
	resp, _, err := client.DBaaSMySQL.GetCluster(context.Background(), "cluster-123")

	// Should handle null fields without panic
	if resp == nil && err == nil {
		t.Error("Expected response or error for null fields")
	}
}

func TestGetMySQLCluster_InvalidVersionField(t *testing.T) {
	server := newInvalidFieldTypeServer(t, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data": map[string]interface{}{
			"cluster_id": "cluster-123",
			"name":       "test-cluster",
			"version":    123, // Should be string, not int
		},
	})
	defer server.Close()

	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	resp, _, err := client.DBaaSMySQL.GetCluster(context.Background(), "cluster-123")

	// Should handle wrong type gracefully
	if resp == nil && err == nil {
		t.Error("Expected response or error for invalid field type")
	}
}
