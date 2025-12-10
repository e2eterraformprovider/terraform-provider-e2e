package goe2e

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMasterPlans(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodGet)
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": [
				{
					"plan": "standard-k8s",
					"k8s_version": "1.28.0",
					"specs": {
						"id": "spec-123",
						"sku_name": "standard-sku"
					}
				},
				{
					"plan": "premium-k8s",
					"k8s_version": "1.29.0",
					"specs": {
						"id": "spec-456",
						"sku_name": "premium-sku"
					}
				}
			]
		}`)
	}))
	defer server.Close()
	client, err := NewClient("test-key", "test-token", "test-project", "test-location", SetBaseURL(server.URL+"/"))
	assertNoError(t, err)
	plans, _, _ := client.Kubernetes.GetMasterPlans(context.Background())
	if len(plans) != 2 {
		t.Errorf("Expected 2 plans, got %d", len(plans))
	}
	if plans[0].Plan != "standard-k8s" {
		t.Errorf("Expected Plan standard-k8s, got %s", plans[0].Plan)
	}
	if plans[0].K8sVersion != "1.28.0" {
		t.Errorf("Expected K8sVersion 1.28.0, got %s", plans[0].K8sVersion)
	}
	if plans[0].Specs.ID != "spec-123" {
		t.Errorf("Expected Specs.ID spec-123, got %s", plans[0].Specs.ID)
	}
}

func TestGetWorkerPlans(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodGet)
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": [
				{
					"plan": "worker-standard",
					"specs": {
						"id": "worker-spec-123",
						"sku_name": "worker-standard-sku"
					}
				}
			]
		}`)
	}))
	defer server.Close()
	client, err := NewClient("test-key", "test-token", "test-project", "test-location", SetBaseURL(server.URL+"/"))
	assertNoError(t, err)
	plans, _, err := client.Kubernetes.GetWorkerPlans(context.Background())
	assertNoError(t, err)
	if len(plans) != 1 {
		t.Errorf("Expected 1 plan, got %d", len(plans))
	}
	if plans[0].Plan != "worker-standard" {
		t.Errorf("Expected Plan worker-standard, got %s", plans[0].Plan)
	}
	if plans[0].Specs.SKUName != "worker-standard-sku" {
		t.Errorf("Expected SKUName worker-standard-sku, got %s", plans[0].Specs.SKUName)
	}
}

func TestCreateKubernetesCluster(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPost)
		testURLPath(t, r, "/kubernetes/")
		assertQueryParam(t, r, "apikey", "test-api-key")
		assertQueryParam(t, r, "project_id", "test-project")
		assertQueryParam(t, r, "location", "test-location")
		writeJSON(w, http.StatusCreated, buildSuccessResponse(201, "Cluster created successfully", map[string]interface{}{
			"DOCUMENT": map[string]interface{}{
				"ID": "cluster-123",
			},
		}))
	}))
	defer server.Close()
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	createReq := &KubernetesClusterCreateRequest{
		Name:     "test-cluster",
		Version:  "1.28.0",
		VPCID:    "vpc-123",
		SKUID:    "sku-456",
		SlugName: "standard-k8s",
		NodePools: []NodePool{
			{
				Name:       "worker-pool",
				SlugName:   "worker-standard",
				SKUID:      "worker-sku-123",
				SpecsName:  "worker-standard-sku",
				WorkerNode: 2,
			},
		},
	}
	cluster, _, err := client.Kubernetes.Create(context.Background(), createReq)
	assertNoError(t, err)
	assertNotNil(t, cluster, "Expected cluster to be returned")
	if cluster.ID != "cluster-123" {
		t.Errorf("Expected ID cluster-123, got %s", cluster.ID)
	}
}

func TestGetKubernetesCluster(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testURLPath(t, r, "/kubernetes/cluster-123")
		writeJSON(w, http.StatusOK, buildListResponse([]interface{}{
			map[string]interface{}{
				"service_id":   123,
				"service_name": "test-cluster",
				"state":        "Running",
				"version":      "1.28.0",
				"created_at":   "2024-01-01T00:00:00Z",
			},
		}))
	}))
	defer server.Close()
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	cluster, _, err := client.Kubernetes.Get(context.Background(), "cluster-123")
	assertNoError(t, err)
	if cluster.ServiceID != "123" {
		t.Errorf("Expected ServiceID 123, got %s", cluster.ServiceID)
	}
	if cluster.ServiceName != "test-cluster" {
		t.Errorf("Expected ServiceName test-cluster, got %s", cluster.ServiceName)
	}
	if cluster.State != "Running" {
		t.Errorf("Expected State Running, got %s", cluster.State)
	}
}
func TestGetKubernetesCluster_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	cluster, _, err := client.Kubernetes.Get(context.Background(), "nonexistent")
	if cluster != nil {
		t.Errorf("Expected nil cluster for 404, got: %v", cluster)
	}
	assertError(t, err, "")
}

func TestDeleteKubernetesCluster(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	_, err = client.Kubernetes.Delete(context.Background(), "cluster-123")
	assertNoError(t, err)
}

func TestGetNodePools(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testURLPath(t, r, "/kubernetes/node-pool-services/cluster-123")
		writeJSON(w, http.StatusOK, buildListResponse([]interface{}{
			map[string]interface{}{
				"service_id":   456,
				"service_name": "worker-pool-1",
				"cardinality":  3,
			},
			map[string]interface{}{
				"service_id":   789,
				"service_name": "worker-pool-2",
				"state":        "Provisioning",
				"cardinality":  2,
			},
		}))
	}))
	defer server.Close()
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	nodePools, _, err := client.Kubernetes.GetNodePools(context.Background(), "cluster-123")
	assertNoError(t, err)
	if len(nodePools) != 2 {
		t.Errorf("Expected 2 node pools, got %d", len(nodePools))
	}
	if nodePools[0].ServiceID != 456 {
		t.Errorf("Expected ServiceID 456, got %v", nodePools[0].ServiceID)
	}
	if nodePools[0].ServiceName != "worker-pool-1" {
		t.Errorf("Expected ServiceName worker-pool-1, got %s", nodePools[0].ServiceName)
	}
	if nodePools[0].Cardinality != 3 {
		t.Errorf("Expected Cardinality 3, got %d", nodePools[0].Cardinality)
	}
}

func TestCheckNodePoolStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testURLPath(t, r, "/kubernetes/node-pool-services/cluster-123")
		writeJSON(w, http.StatusOK, buildListResponse([]interface{}{
			map[string]interface{}{
				"service_id": 456,
			},
		}))
	}))
	defer server.Close()
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	nodePools, _, err := client.Kubernetes.CheckNodePoolStatus(context.Background(), "cluster-123")
	assertNoError(t, err)
	if len(nodePools) != 1 {
		t.Errorf("Expected 1 node pool, got %d", len(nodePools))
	}
}

func TestAddNodePool(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPost)
		testURLPath(t, r, "/kubernetes/add-node-pools/cluster-123")
		writeJSON(w, http.StatusCreated, buildSuccessResponse(201, "Node pool added successfully", map[string]interface{}{}))
	}))
	defer server.Close()
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	addReq := &NodePoolAddRequest{
		NodePools: []NodePool{
			{
				Name: "new-worker-pool",
			},
		},
	}
	result, _, err := client.Kubernetes.AddNodePool(context.Background(), "cluster-123", addReq)
	assertNoError(t, err)
	assertNotNil(t, result, "Expected result to be returned")
	if result.Message != "Node pool added successfully" {
		t.Errorf("Expected Message 'Node pool added successfully', got %s", result.Message)
	}
}

func TestUpdateNodePoolCardinality(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/kubernetes/cluster-update/456")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	resizeReq := &NodePoolResizeRequest{
		NodePoolSize: 5,
	}
	_, err = client.Kubernetes.UpdateNodePoolCardinality(context.Background(), "456", resizeReq)
	assertNoError(t, err)
}

func TestUpdateNodePoolDetails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/kubernetes/update-node-pool/456/")
		writeJSON(w, http.StatusOK, buildSuccessResponse(200, "Node pool updated successfully", map[string]interface{}{}))
	}))
	defer server.Close()
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	updateReq := &NodePoolUpdateRequest{
		MinVms:      2,
		Cardinality: 3,
		MaxVms:      5,
		PlanID:      "plan-123",
	}
	result, _, err := client.Kubernetes.UpdateNodePoolDetails(context.Background(), "456", updateReq)
	assertNoError(t, err)
	if result.Message != "Node pool updated successfully" {
		t.Errorf("Expected Message 'Node pool updated successfully', got %s", result.Message)
	}
}

func TestDeleteNodePool(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodDelete)
		testURLPath(t, r, "/kubernetes/delete-node-pool-service/456")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	_, err = client.Kubernetes.DeleteNodePool(context.Background(), "456")
	assertNoError(t, err)
}

// Edge case tests for better coverage
func TestCreateKubernetesCluster_NilRequest(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	assertNoError(t, err)
	_, _, err = client.Kubernetes.Create(context.Background(), nil)
	assertError(t, err, "")
}

func TestCreateKubernetesCluster_EmptyName(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	assertNoError(t, err)
	createReq := &KubernetesClusterCreateRequest{
		Name:  "",
		VPCID: "vpc-123",
	}
	_, _, err = client.Kubernetes.Create(context.Background(), createReq)
	assertError(t, err, "")
}

func TestCreateKubernetesCluster_EmptyVPCID(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	assertNoError(t, err)
	createReq := &KubernetesClusterCreateRequest{
		Name:  "test-cluster",
		VPCID: "",
	}
	_, _, err = client.Kubernetes.Create(context.Background(), createReq)
	assertError(t, err, "")
}

func TestGetKubernetesCluster_EmptyID(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	assertNoError(t, err)
	_, _, err = client.Kubernetes.Get(context.Background(), "")
	assertError(t, err, "")
}

func TestDeleteKubernetesCluster_EmptyID(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	assertNoError(t, err)
	_, err = client.Kubernetes.Delete(context.Background(), "")
	assertError(t, err, "")
}

func TestGetNodePools_EmptyClusterID(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	assertNoError(t, err)
	_, _, err = client.Kubernetes.GetNodePools(context.Background(), "")
	assertError(t, err, "")
}

func TestAddNodePool_EmptyClusterID(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	assertNoError(t, err)
	addReq := &NodePoolAddRequest{
		NodePools: []NodePool{},
	}
	_, _, err = client.Kubernetes.AddNodePool(context.Background(), "", addReq)
	assertError(t, err, "")
}

func TestAddNodePool_NilRequest(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	assertNoError(t, err)
	_, _, err = client.Kubernetes.AddNodePool(context.Background(), "cluster-123", nil)
	assertError(t, err, "")
}

func TestUpdateNodePoolCardinality_EmptyServiceID(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	assertNoError(t, err)
	resizeReq := &NodePoolResizeRequest{
		NodePoolSize: 5,
	}
	_, err = client.Kubernetes.UpdateNodePoolCardinality(context.Background(), "", resizeReq)
	assertError(t, err, "")
}

func TestUpdateNodePoolCardinality_NilRequest(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	assertNoError(t, err)
	_, err = client.Kubernetes.UpdateNodePoolCardinality(context.Background(), "456", nil)
	assertError(t, err, "")
}

func TestUpdateNodePoolDetails_EmptyServiceID(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	assertNoError(t, err)
	updateReq := &NodePoolUpdateRequest{
		MinVms: 2,
	}
	_, _, err = client.Kubernetes.UpdateNodePoolDetails(context.Background(), "", updateReq)
	assertError(t, err, "")
}

func TestUpdateNodePoolDetails_NilRequest(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	assertNoError(t, err)
	_, _, err = client.Kubernetes.UpdateNodePoolDetails(context.Background(), "456", nil)
	assertError(t, err, "")
}

func TestDeleteNodePool_EmptyServiceID(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	assertNoError(t, err)
	_, err = client.Kubernetes.DeleteNodePool(context.Background(), "")
	assertError(t, err, "")
}

// Additional error tests for better coverage
func TestGetMasterPlans_ErrorResponse(t *testing.T) {
	server := newErrorServer(t, http.StatusInternalServerError, "Internal server error")
	defer server.Close()
	client, err := NewClient("test-key", "test-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	_, _, err = client.Kubernetes.GetMasterPlans(context.Background())
	assertError(t, err, "")
}

func TestGetWorkerPlans_ErrorResponse(t *testing.T) {
	server := newErrorServer(t, http.StatusInternalServerError, "Internal server error")
	defer server.Close()
	client, err := NewClient("test-key", "test-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	_, _, err = client.Kubernetes.GetWorkerPlans(context.Background())
	assertError(t, err, "")
}

func TestCreateKubernetesCluster_ErrorResponse(t *testing.T) {
	server := newErrorServer(t, http.StatusBadRequest, "Bad request")
	defer server.Close()
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	createReq := &KubernetesClusterCreateRequest{
		Name:  "test-cluster",
		VPCID: "vpc-123",
	}
	_, _, err = client.Kubernetes.Create(context.Background(), createReq)
	assertError(t, err, "")
}

func TestGetKubernetesCluster_ErrorResponse(t *testing.T) {
	server := newErrorServer(t, http.StatusInternalServerError, "Internal server error")
	defer server.Close()
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	_, _, err = client.Kubernetes.Get(context.Background(), "cluster-123")
	assertError(t, err, "")
}

func TestGetKubernetesCluster_EmptyData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, buildListResponse([]interface{}{}))
	}))
	defer server.Close()
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	cluster, _, err := client.Kubernetes.Get(context.Background(), "cluster-123")
	if cluster != nil {
		t.Errorf("Expected nil cluster for empty data, got: %v", cluster)
	}
	assertError(t, err, "")
}

func TestDeleteKubernetesCluster_ErrorResponse(t *testing.T) {
	server := newErrorServer(t, http.StatusInternalServerError, "Internal server error")
	defer server.Close()
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	_, err = client.Kubernetes.Delete(context.Background(), "cluster-123")
	assertError(t, err, "")
}

func TestGetNodePools_ErrorResponse(t *testing.T) {
	server := newErrorServer(t, http.StatusInternalServerError, "Internal server error")
	defer server.Close()
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	_, _, err = client.Kubernetes.GetNodePools(context.Background(), "cluster-123")
	assertError(t, err, "")
}

func TestAddNodePool_ErrorResponse(t *testing.T) {
	server := newErrorServer(t, http.StatusBadRequest, "Bad request")
	defer server.Close()
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	addReq := &NodePoolAddRequest{
		NodePools: []NodePool{
			{
				Name: "new-pool",
			},
		},
	}
	_, _, err = client.Kubernetes.AddNodePool(context.Background(), "cluster-123", addReq)
	assertError(t, err, "")
}

func TestUpdateNodePoolCardinality_ErrorResponse(t *testing.T) {
	server := newErrorServer(t, http.StatusBadRequest, "Bad request")
	defer server.Close()
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	resizeReq := &NodePoolResizeRequest{
		NodePoolSize: 5,
	}
	_, err = client.Kubernetes.UpdateNodePoolCardinality(context.Background(), "456", resizeReq)
	assertError(t, err, "")
}

func TestUpdateNodePoolDetails_ErrorResponse(t *testing.T) {
	server := newErrorServer(t, http.StatusBadRequest, "Bad request")
	defer server.Close()
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	updateReq := &NodePoolUpdateRequest{
		MinVms: 2,
	}
	_, _, err = client.Kubernetes.UpdateNodePoolDetails(context.Background(), "456", updateReq)
	assertError(t, err, "")
}

func TestDeleteNodePool_ErrorResponse(t *testing.T) {
	server := newErrorServer(t, http.StatusInternalServerError, "Internal server error")
	defer server.Close()
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	_, err = client.Kubernetes.DeleteNodePool(context.Background(), "456")
	assertError(t, err, "")
}
