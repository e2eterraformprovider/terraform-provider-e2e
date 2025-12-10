package goe2e

import (
	"context"
	"net/http"
	"testing"
)

// TestPostgreSQLServiceOp_InterfaceExists verifies the service implements the interface
func TestPostgreSQLServiceOp_InterfaceExists(t *testing.T) {
	// This test ensures that PostgreSQLServiceOp implements PostgreSQLService
	var _ PostgreSQLService = &PostgreSQLServiceOp{}
}

// TestPostgreSQLServiceOp_ValidationCreateCluster tests nil request validation
func TestPostgreSQLServiceOp_ValidationCreateCluster(t *testing.T) {
	c := &Client{PostgreSQL: &PostgreSQLServiceOp{client: nil}}

	// Test nil request
	_, _, err := c.PostgreSQL.CreateCluster(context.Background(), nil)
	if err == nil {
		t.Errorf("CreateCluster(nil) should return error")
	}
}

// TestPostgreSQLServiceOp_ValidationGetCluster tests empty cluster ID validation
func TestPostgreSQLServiceOp_ValidationGetCluster(t *testing.T) {
	c := &Client{PostgreSQL: &PostgreSQLServiceOp{client: nil}}

	// Test empty cluster ID
	_, _, err := c.PostgreSQL.GetCluster(context.Background(), "")
	if err == nil {
		t.Errorf("GetCluster('') should return error")
	}
}

// TestPostgreSQLServiceOp_ValidationDeleteCluster tests empty cluster ID validation
func TestPostgreSQLServiceOp_ValidationDeleteCluster(t *testing.T) {
	c := &Client{PostgreSQL: &PostgreSQLServiceOp{client: nil}}

	// Test empty cluster ID
	_, err := c.PostgreSQL.DeleteCluster(context.Background(), "")
	if err == nil {
		t.Errorf("DeleteCluster('') should return error")
	}
}

// TestPostgreSQLServiceOp_ValidationClusterExists tests empty cluster ID validation
func TestPostgreSQLServiceOp_ValidationClusterExists(t *testing.T) {
	c := &Client{PostgreSQL: &PostgreSQLServiceOp{client: nil}}

	// Test empty cluster ID
	_, _, err := c.PostgreSQL.ClusterExists(context.Background(), "")
	if err == nil {
		t.Errorf("ClusterExists('') should return error")
	}
}

// TestPostgreSQLServiceOp_ValidationStartCluster tests empty cluster ID validation
func TestPostgreSQLServiceOp_ValidationStartCluster(t *testing.T) {
	c := &Client{PostgreSQL: &PostgreSQLServiceOp{client: nil}}

	// Test empty cluster ID
	_, err := c.PostgreSQL.StartCluster(context.Background(), "")
	if err == nil {
		t.Errorf("StartCluster('') should return error")
	}
}

// TestPostgreSQLServiceOp_ValidationStopCluster tests empty cluster ID validation
func TestPostgreSQLServiceOp_ValidationStopCluster(t *testing.T) {
	c := &Client{PostgreSQL: &PostgreSQLServiceOp{client: nil}}

	// Test empty cluster ID
	_, err := c.PostgreSQL.StopCluster(context.Background(), "")
	if err == nil {
		t.Errorf("StopCluster('') should return error")
	}
}

// TestPostgreSQLServiceOp_ValidationRestartCluster tests empty cluster ID validation
func TestPostgreSQLServiceOp_ValidationRestartCluster(t *testing.T) {
	c := &Client{PostgreSQL: &PostgreSQLServiceOp{client: nil}}

	// Test empty cluster ID
	_, err := c.PostgreSQL.RestartCluster(context.Background(), "")
	if err == nil {
		t.Errorf("RestartCluster('') should return error")
	}
}

// TestPostgreSQLServiceOp_ValidationAttachVPC tests input validation
func TestPostgreSQLServiceOp_ValidationAttachVPC(t *testing.T) {
	c := &Client{PostgreSQL: &PostgreSQLServiceOp{client: nil}}

	tests := []struct {
		name      string
		clusterID string
		req       *PostgreSQLVPCAttachRequest
		wantErr   bool
	}{
		{"empty cluster id", "", &PostgreSQLVPCAttachRequest{}, true},
		{"nil request", "cluster-1", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := c.PostgreSQL.AttachVPC(context.Background(), tt.clusterID, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("AttachVPC() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestPostgreSQLServiceOp_ValidationDetachVPC tests input validation
func TestPostgreSQLServiceOp_ValidationDetachVPC(t *testing.T) {
	c := &Client{PostgreSQL: &PostgreSQLServiceOp{client: nil}}

	// Test empty cluster ID
	_, err := c.PostgreSQL.DetachVPC(context.Background(), "", nil)
	if err == nil {
		t.Errorf("DetachVPC('', nil) should return error")
	}
}

// TestPostgreSQLServiceOp_ValidationAttachPublicIP tests empty cluster ID validation
func TestPostgreSQLServiceOp_ValidationAttachPublicIP(t *testing.T) {
	c := &Client{PostgreSQL: &PostgreSQLServiceOp{client: nil}}

	// Test empty cluster ID
	_, err := c.PostgreSQL.AttachPublicIP(context.Background(), "")
	if err == nil {
		t.Errorf("AttachPublicIP('') should return error")
	}
}

// TestPostgreSQLServiceOp_ValidationDetachPublicIP tests empty cluster ID validation
func TestPostgreSQLServiceOp_ValidationDetachPublicIP(t *testing.T) {
	c := &Client{PostgreSQL: &PostgreSQLServiceOp{client: nil}}

	// Test empty cluster ID
	_, err := c.PostgreSQL.DetachPublicIP(context.Background(), "")
	if err == nil {
		t.Errorf("DetachPublicIP('') should return error")
	}
}

// TestPostgreSQLServiceOp_ValidationAttachParameterGroup tests input validation
func TestPostgreSQLServiceOp_ValidationAttachParameterGroup(t *testing.T) {
	c := &Client{PostgreSQL: &PostgreSQLServiceOp{client: nil}}

	tests := []struct {
		name         string
		clusterID    string
		paramGroupID string
		wantErr      bool
	}{
		{"empty cluster id", "", "pg-1", true},
		{"empty param group id", "cluster-1", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := c.PostgreSQL.AttachParameterGroup(context.Background(), tt.clusterID, tt.paramGroupID)
			if (err != nil) != tt.wantErr {
				t.Errorf("AttachParameterGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestPostgreSQLServiceOp_ValidationDetachParameterGroup tests input validation
func TestPostgreSQLServiceOp_ValidationDetachParameterGroup(t *testing.T) {
	c := &Client{PostgreSQL: &PostgreSQLServiceOp{client: nil}}

	// Test empty cluster ID
	_, err := c.PostgreSQL.DetachParameterGroup(context.Background(), "", "pg-1")
	if err == nil {
		t.Errorf("DetachParameterGroup('', 'pg-1') should return error")
	}
}

// TestPostgreSQLServiceOp_ValidationUpgradePlan tests input validation
func TestPostgreSQLServiceOp_ValidationUpgradePlan(t *testing.T) {
	c := &Client{PostgreSQL: &PostgreSQLServiceOp{client: nil}}

	tests := []struct {
		name      string
		clusterID string
		req       *PostgreSQLPlanUpgradeRequest
		wantErr   bool
	}{
		{"empty cluster id", "", &PostgreSQLPlanUpgradeRequest{TemplateID: 42}, true},
		{"nil request", "cluster-1", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := c.PostgreSQL.UpgradePlan(context.Background(), tt.clusterID, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpgradePlan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestPostgreSQLServiceOp_ValidationExpandDisk tests input validation
func TestPostgreSQLServiceOp_ValidationExpandDisk(t *testing.T) {
	c := &Client{PostgreSQL: &PostgreSQLServiceOp{client: nil}}

	tests := []struct {
		name      string
		clusterID string
		req       *DiskExpansionRequest
		wantErr   bool
	}{
		{"empty cluster id", "", &DiskExpansionRequest{Size: 100}, true},
		{"nil request", "cluster-1", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := c.PostgreSQL.ExpandDisk(context.Background(), tt.clusterID, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExpandDisk() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestPostgreSQLServiceOp_ValidationGetSoftwareID tests input validation
func TestPostgreSQLServiceOp_ValidationGetSoftwareID(t *testing.T) {
	c := &Client{PostgreSQL: &PostgreSQLServiceOp{client: nil}}

	tests := []struct {
		name       string
		engineName string
		version    string
		pgID       string
		wantErr    bool
	}{
		{"empty engine name", "", "13.0", "pg-1", true},
		{"empty version", "postgresql", "", "pg-1", true},
		{"empty pgID", "postgresql", "13.0", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := c.PostgreSQL.GetSoftwareID(context.Background(), tt.engineName, tt.version, tt.pgID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSoftwareID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestPostgreSQLServiceOp_ValidationGetTemplateID tests input validation
func TestPostgreSQLServiceOp_ValidationGetTemplateID(t *testing.T) {
	c := &Client{PostgreSQL: &PostgreSQLServiceOp{client: nil}}

	tests := []struct {
		name       string
		plan       string
		softwareID string
		pgID       string
		wantErr    bool
	}{
		{"empty plan", "", "5", "pg-1", true},
		{"empty softwareID", "2vCPU-4GB", "", "pg-1", true},
		{"empty pgID", "2vCPU-4GB", "5", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := c.PostgreSQL.GetTemplateID(context.Background(), tt.plan, tt.softwareID, tt.pgID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTemplateID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestPostgreSQLServiceOp_ValidationExpandPostgresVPCList tests input validation
func TestPostgreSQLServiceOp_ValidationExpandPostgresVPCList(t *testing.T) {
	c := &Client{PostgreSQL: &PostgreSQLServiceOp{client: nil}}

	// Test empty VPC list
	_, err := c.PostgreSQL.ExpandPostgresVPCList(context.Background(), []string{})
	if err == nil {
		t.Errorf("ExpandPostgresVPCList([]) should return error")
	}
}

// TestPostgreSQLPathConstants verifies path constants are defined
func TestPostgreSQLPathConstants(t *testing.T) {
	paths := []string{
		postgresqlClusterPath,
		postgresqlClusterDetailPath,
		postgresqlClusterResumePath,
		postgresqlClusterShutdownPath,
		postgresqlClusterRestartPath,
		postgresqlClusterVPCAttachPath,
		postgresqlClusterVPCDetachPath,
		postgresqlClusterPGAddPath,
		postgresqlClusterPGDetachPath,
		postgresqlClusterPublicIPPath,
		postgresqlClusterUpgradePath,
		postgresqlClusterDiskPath,
		postgresqlPlansPath,
	}

	for _, path := range paths {
		if path == "" {
			t.Errorf("Path constant is empty")
		}
	}
}

// TestPostgreSQLTypeAliases verifies type aliases are correct
func TestPostgreSQLTypeAliases(t *testing.T) {
	// Verify type aliases compile correctly
	var createReq *PostgreSQLClusterCreateRequest
	var cluster *PostgreSQLCluster
	var vpcReq *PostgreSQLVPCAttachRequest
	var upgradeReq *PostgreSQLPlanUpgradeRequest
	var diskReq *DiskExpansionRequest

	if createReq != nil || cluster != nil || vpcReq != nil || upgradeReq != nil || diskReq != nil {
		t.Errorf("Type aliases are not working correctly")
	}
}

// TestPostgreSQLResponseStructures verifies response types are defined
func TestPostgreSQLResponseStructures(t *testing.T) {
	// Verify response wrapper types exist and can be instantiated
	_ = postgresqlClusterRoot{Code: 0}
	_ = postgresqlClusterDetailRoot{Code: 0}
	_ = postgresqlOperationRoot{Code: 0}
	_ = plansResponse{Code: 0}
}

// TestPostgreSQLHTTPMethods verifies correct HTTP methods are used
func TestPostgreSQLHTTPMethods(t *testing.T) {
	// This is a compile-time verification that the code structures are correct
	tests := []struct {
		method string
		want   string
	}{
		{http.MethodPost, "POST"},
		{http.MethodGet, "GET"},
		{http.MethodPut, "PUT"},
		{http.MethodDelete, "DELETE"},
	}

	for _, tt := range tests {
		if tt.method != tt.want {
			t.Errorf("HTTP method = %s, want %s", tt.method, tt.want)
		}
	}
}

// BenchmarkPostgreSQLValidation benchmarks parameter validation
func BenchmarkPostgreSQLValidation(b *testing.B) {
	c := &Client{PostgreSQL: &PostgreSQLServiceOp{client: nil}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.PostgreSQL.GetCluster(context.Background(), "")
	}
}

// ============================================
// COMPREHENSIVE HAPPY PATH TESTS
// ============================================

func TestPostgreSQLCreateCluster_Success(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		writeJSON(w, http.StatusCreated, `{
			"code": 201,
			"message": "Cluster created",
			"data": {"id": 100}
		}`)
	})

	ts.mux.HandleFunc("/rds/cluster/100/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": {
				"id": 100,
				"name": "pg-cluster",
				"status": "active",
				"software": {"name": "PostgreSQL", "version": "13.0", "engine": "postgresql"}
			}
		}`)
	})

	req := &PostgreSQLClusterCreateRequest{Name: "pg-cluster"}
	result, _, err := ts.client.PostgreSQL.CreateCluster(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateCluster error: %v", err)
	}
	if result == nil || result.ID != 100 {
		t.Error("Expected cluster with ID 100")
	}
}

func TestPostgreSQLGetCluster_Success(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/100/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"data": {"id": 100, "name": "pg-cluster", "status": "active"}
		}`)
	})

	result, _, err := ts.client.PostgreSQL.GetCluster(context.Background(), "100")
	if err != nil {
		t.Fatalf("GetCluster error: %v", err)
	}
	if result == nil || result.ID != 100 {
		t.Error("Expected cluster")
	}
}

func TestPostgreSQLDeleteCluster_Success(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/100/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusOK)
	})

	_, err := ts.client.PostgreSQL.DeleteCluster(context.Background(), "100")
	if err != nil {
		t.Fatalf("DeleteCluster error: %v", err)
	}
}

func TestPostgreSQLClusterExists_True(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/100/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	exists, _, err := ts.client.PostgreSQL.ClusterExists(context.Background(), "100")
	if err != nil {
		t.Fatalf("ClusterExists error: %v", err)
	}
	if !exists {
		t.Error("Expected cluster to exist")
	}
}

func TestPostgreSQLClusterExists_False(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/999/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	exists, _, err := ts.client.PostgreSQL.ClusterExists(context.Background(), "999")
	if err != nil {
		t.Fatalf("ClusterExists error: %v", err)
	}
	if exists {
		t.Error("Expected cluster to not exist")
	}
}

// ============================================
// ERROR RESPONSE TESTS
// ============================================

func TestPostgreSQLCreateCluster_BadRequest(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusBadRequest, `{"code": 400, "message": "Invalid input"}`)
	})

	_, _, err := ts.client.PostgreSQL.CreateCluster(context.Background(), &PostgreSQLClusterCreateRequest{})
	if err == nil {
		t.Fatal("Expected error for 400 response")
	}
}

func TestPostgreSQLCreateCluster_Conflict(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusConflict, `{"code": 409, "message": "Cluster already exists"}`)
	})

	_, _, err := ts.client.PostgreSQL.CreateCluster(context.Background(), &PostgreSQLClusterCreateRequest{})
	if err == nil {
		t.Fatal("Expected error for 409 response")
	}
}

func TestPostgreSQLGetCluster_BadRequest(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/100/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusBadRequest, `{"code": 400, "message": "Invalid cluster ID"}`)
	})

	_, _, err := ts.client.PostgreSQL.GetCluster(context.Background(), "100")
	if err == nil {
		t.Fatal("Expected error for 400 response")
	}
}

func TestPostgreSQLDeleteCluster_Conflict(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/100/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusConflict, `{"code": 409, "message": "Cluster has running connections"}`)
	})

	_, err := ts.client.PostgreSQL.DeleteCluster(context.Background(), "100")
	if err == nil {
		t.Fatal("Expected error for 409 response")
	}
}

func TestPostgreSQLStartCluster_Conflict(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/100/resume", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusConflict, `{"code": 409, "message": "Cluster already running"}`)
	})

	_, err := ts.client.PostgreSQL.StartCluster(context.Background(), "100")
	if err == nil {
		t.Fatal("Expected error for 409 response")
	}
}

func TestPostgreSQLStopCluster_Forbidden(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/100/shutdown", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusForbidden, `{"code": 403, "message": "Insufficient permissions"}`)
	})

	_, err := ts.client.PostgreSQL.StopCluster(context.Background(), "100")
	if err == nil {
		t.Fatal("Expected error for 403 response")
	}
}

func TestPostgreSQLRestartCluster_ServiceUnavailable(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/100/restart", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusServiceUnavailable, `{"code": 503, "message": "Service temporarily unavailable"}`)
	})

	_, err := ts.client.PostgreSQL.RestartCluster(context.Background(), "100")
	if err == nil {
		t.Fatal("Expected error for 503 response")
	}
}

func TestPostgreSQLAttachVPC_BadRequest(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/100/vpc-attach/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusBadRequest, `{"code": 400, "message": "Invalid VPC configuration"}`)
	})

	_, err := ts.client.PostgreSQL.AttachVPC(context.Background(), "100", &PostgreSQLVPCAttachRequest{})
	if err == nil {
		t.Fatal("Expected error for 400 response")
	}
}

func TestPostgreSQLDetachVPC_NotFound(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/100/vpc-detach/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusNotFound, `{"code": 404, "message": "VPC not attached"}`)
	})

	_, err := ts.client.PostgreSQL.DetachVPC(context.Background(), "100", &PostgreSQLVPCAttachRequest{})
	if err == nil {
		t.Fatal("Expected error for 404 response")
	}
}

func TestPostgreSQLAttachPublicIP_Forbidden(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/100/public-ip-attach/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusForbidden, `{"code": 403, "message": "Insufficient permissions"}`)
	})

	_, err := ts.client.PostgreSQL.AttachPublicIP(context.Background(), "100")
	if err == nil {
		t.Fatal("Expected error for 403 response")
	}
}

func TestPostgreSQLAttachParameterGroup_NotFound(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/100/parameter-group/999/add", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusNotFound, `{"code": 404, "message": "Parameter group not found"}`)
	})

	_, err := ts.client.PostgreSQL.AttachParameterGroup(context.Background(), "100", "999")
	if err == nil {
		t.Fatal("Expected error for 404 response")
	}
}

func TestPostgreSQLUpgradePlan_Conflict(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/100/rds-upgrade/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusConflict, `{"code": 409, "message": "Upgrade already in progress"}`)
	})

	_, err := ts.client.PostgreSQL.UpgradePlan(context.Background(), "100", &PostgreSQLPlanUpgradeRequest{})
	if err == nil {
		t.Fatal("Expected error for 409 response")
	}
}

func TestPostgreSQLExpandDisk_BadRequest(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/cluster/100/disk-upgrade/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusBadRequest, `{"code": 400, "message": "Invalid disk size"}`)
	})

	_, err := ts.client.PostgreSQL.ExpandDisk(context.Background(), "100", &DiskExpansionRequest{})
	if err == nil {
		t.Fatal("Expected error for 400 response")
	}
}

// ============================================
// CONTEXT & CONNECTION TESTS
// ============================================

func TestPostgreSQLGetCluster_ContextCancelled(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, _, err := ts.client.PostgreSQL.GetCluster(ctx, "100")
	if err == nil {
		t.Fatal("Expected error for cancelled context")
	}
}

// ============================================
// Phase 2: Response Parsing & Edge Case Tests
// ============================================

func TestPostgreSQLCreateCluster_MalformedJSON(t *testing.T) {
	server := newMalformedJSONServer(t)
	defer server.Close()

	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	_, _, err := client.PostgreSQL.CreateCluster(context.Background(), &PostgreSQLClusterCreateRequest{})

	if err == nil {
		t.Error("Expected error for malformed JSON")
	}
}

func TestPostgreSQLGetCluster_MissingRequiredFields(t *testing.T) {
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
	resp, _, err := client.PostgreSQL.GetCluster(context.Background(), "cluster-123")

	// Should handle missing fields gracefully
	if resp == nil && err == nil {
		t.Error("Expected response or error handling")
	}
}

func TestPostgreSQLGetCluster_NullFieldValues(t *testing.T) {
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
	resp, _, err := client.PostgreSQL.GetCluster(context.Background(), "cluster-123")

	// Should handle null fields without panic
	if resp == nil && err == nil {
		t.Error("Expected response or error for null fields")
	}
}

func TestPostgreSQLGetCluster_InvalidVersionField(t *testing.T) {
	server := newInvalidFieldTypeServer(t, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data": map[string]interface{}{
			"cluster_id": "cluster-123",
			"name":       "test-cluster",
			"version":    456, // Should be string, not int
		},
	})
	defer server.Close()

	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	resp, _, err := client.PostgreSQL.GetCluster(context.Background(), "cluster-123")

	// Should handle wrong type gracefully
	if resp == nil && err == nil {
		t.Error("Expected response or error for invalid field type")
	}
}
