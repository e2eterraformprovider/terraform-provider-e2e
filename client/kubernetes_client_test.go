package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func TestGetKubernetesMasterPlans(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code": 200,
		"data": []interface{}{
			map[string]interface{}{
				"name":   "master-plan-1",
				"cpu":    4,
				"memory": 8192,
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/kubernetes/plans" {
			t.Errorf("Expected path /kubernetes/plans, got %s", r.URL.Path)
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
	result, err := client.GetKubernetesMasterPlans(123, "us-east")

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

func TestGetKubernetesWorkerPlans(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code": 200,
		"data": []interface{}{
			map[string]interface{}{
				"name":   "worker-plan-1",
				"cpu":    2,
				"memory": 4096,
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/kubernetes/worker-plans/" {
			t.Errorf("Expected path /kubernetes/worker-plans/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	result, err := client.GetKubernetesWorkerPlans(123, "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestNewKubernetesService(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code":    201,
		"message": "Kubernetes cluster created successfully",
		"data": map[string]interface{}{
			"cluster_id": "k8s-123",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/kubernetes/" {
			t.Errorf("Expected path /kubernetes/, got %s", r.URL.Path)
		}

		body, _ := ioutil.ReadAll(r.Body)
		var kubernetesCreate models.KubernetesCreate
		json.Unmarshal(body, &kubernetesCreate)

		if kubernetesCreate.Name == "" {
			t.Error("Expected Name in request body")
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	kubernetesCreate := &models.KubernetesCreate{
		Name: "test-cluster",
		NodePools: []models.NodePool{
			{
				Name:       "pool-1",
				WorkerNode: 2,
			},
		},
	}

	result, err := client.NewKubernetesService(kubernetesCreate, 123, "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestGetKubernetesServiceInfo(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{
			"cluster_id":   "k8s-123",
			"cluster_name": "test-cluster",
			"state":        "ACTIVE",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/kubernetes/k8s-123" {
			t.Errorf("Expected path /kubernetes/k8s-123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	result, err := client.GetKubernetesServiceInfo("k8s-123", "us-east", 123)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestDeleteKubernetesService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/kubernetes/k8s-123" {
			t.Errorf("Expected path /kubernetes/k8s-123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	err := client.DeleteKubernetesService("k8s-123", "us-east", 123)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestRemoveExtraFieldsFromKubernetes(t *testing.T) {
	tests := []struct {
		name        string
		input       map[string]interface{}
		expectError bool
	}{
		{
			name: "Remove worker_node when 0",
			input: map[string]interface{}{
				"node_pools": []interface{}{
					map[string]interface{}{
						"name":        "pool-1",
						"worker_node": float64(0),
					},
				},
			},
			expectError: false,
		},
		{
			name: "Remove elasticity_dict when worker_node >= 2",
			input: map[string]interface{}{
				"node_pools": []interface{}{
					map[string]interface{}{
						"name":             "pool-1",
						"worker_node":      float64(2),
						"elasticity_dict":  []interface{}{},
						"policy_type":      "Manual",
					},
				},
			},
			expectError: false,
		},
		{
			name: "Keep elasticity_dict when policy_type is Scheduled",
			input: map[string]interface{}{
				"node_pools": []interface{}{
					map[string]interface{}{
						"name":            "pool-1",
						"worker_node":     float64(1),
						"elasticity_dict": []interface{}{},
						"policy_type":     "Scheduled",
					},
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputJSON, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("Failed to marshal input: %v", err)
			}
			buf := bytes.NewBuffer(inputJSON)

			result, err := RemoveExtraFieldsFromKubernetes(buf)

			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}

			if !tt.expectError && result.Len() == 0 {
				t.Error("Expected non-empty result buffer")
			}
		})
	}
}

func TestGetKubernetesNodePools(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code": 200,
		"data": []interface{}{
			map[string]interface{}{
				"name":        "pool-1",
				"worker_node": 2,
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/kubernetes/node-pool-services/k8s-123" {
			t.Errorf("Expected path /kubernetes/node-pool-services/k8s-123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	result, err := client.GetKubernetesNodePools("k8s-123", 123, "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestUpdateNodePoolCardinality(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code":    200,
		"message": "Node pool updated",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/kubernetes/cluster-update/123" {
			t.Errorf("Expected path /kubernetes/cluster-update/123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	resize := &models.NodePoolResize{
		NodePoolSize: 3,
	}

	result, err := client.UpdateNodePoolCardinality(resize, 123, 456, "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestUpdateNodePoolCardinalityNoContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	resize := &models.NodePoolResize{
		NodePoolSize: 3,
	}

	result, err := client.UpdateNodePoolCardinality(resize, 123, 456, "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result != nil {
		t.Fatal("Expected nil result for 204 response, got non-nil")
	}
}

func TestDeleteNodePool(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/kubernetes/delete-node-pool-service/123" {
			t.Errorf("Expected path /kubernetes/delete-node-pool-service/123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	result, err := client.DeleteNodePool(123, 456, "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result != nil {
		t.Fatal("Expected nil result for 204 response")
	}
}

func TestAddNodePool(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code":    201,
		"message": "Node pool added",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/kubernetes/add-node-pools/k8s-123" {
			t.Errorf("Expected path /kubernetes/add-node-pools/k8s-123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	nodePoolAdd := &models.NodePoolAdd{
		NodePools: []models.NodePool{
			{
				Name:       "new-pool",
				WorkerNode: 2,
			},
		},
	}

	result, err := client.AddNodePool(nodePoolAdd, "k8s-123", 123, "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestUpdateNodePoolDetails(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code":    200,
		"message": "Node pool details updated",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/kubernetes/update-node-pool/123/" {
			t.Errorf("Expected path /kubernetes/update-node-pool/123/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	nodePoolUpdate := &models.NodePoolUpdate{
		PlanID: "updated-plan",
	}

	result, err := client.UpdateNodePoolDetails(nodePoolUpdate, 123, 456, "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestCheckNodePoolStatus(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code": 200,
		"data": []interface{}{
			map[string]interface{}{
				"name":   "pool-1",
				"status": "ACTIVE",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/kubernetes/node-pool-services/k8s-123" {
			t.Errorf("Expected path /kubernetes/node-pool-services/k8s-123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	result, err := client.CheckNodePoolStatus("k8s-123", 123, "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestGetKubernetesMasterPlansError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "server error"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	result, err := client.GetKubernetesMasterPlans(123, "us-east")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestNewKubernetesServiceError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid request"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	kubernetesCreate := &models.KubernetesCreate{
		Name: "test-cluster",
	}

	result, err := client.NewKubernetesService(kubernetesCreate, 123, "us-east")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}
