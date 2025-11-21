package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func TestGetKubernetesMasterPlans(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/kubernetes/plans", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/kubernetes/plans")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "location", "us-east")
		testQueryParam(t, r, "project_id", "123")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"data": [
				{
					"name": "master-plan-1",
					"cpu": 4,
					"memory": 8192
				}
			]
		}`)
	})

	result, err := ts.client.GetKubernetesMasterPlans(123, "us-east")

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
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/kubernetes/worker-plans/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/kubernetes/worker-plans/")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"data": [
				{
					"name": "worker-plan-1",
					"cpu": 2,
					"memory": 4096
				}
			]
		}`)
	})

	result, err := ts.client.GetKubernetesWorkerPlans(123, "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestNewKubernetesService(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/kubernetes/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testURLPath(t, r, "/kubernetes/")

		writeJSON(w, http.StatusCreated, `{
			"code": 201,
			"message": "Kubernetes cluster created successfully",
			"data": {
				"cluster_id": "k8s-123"
			}
		}`)
	})

	kubernetesCreate := &models.KubernetesCreate{
		Name: "test-cluster",
		NodePools: []models.NodePool{
			{
				Name:       "pool-1",
				WorkerNode: 2,
			},
		},
	}

	result, err := ts.client.NewKubernetesService(kubernetesCreate, 123, "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestGetKubernetesServiceInfo(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/kubernetes/k8s-123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/kubernetes/k8s-123")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"data": {
				"cluster_id": "k8s-123",
				"cluster_name": "test-cluster",
				"state": "ACTIVE"
			}
		}`)
	})

	result, err := ts.client.GetKubernetesServiceInfo("k8s-123", "us-east", 123)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestDeleteKubernetesService(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/kubernetes/k8s-123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		testURLPath(t, r, "/kubernetes/k8s-123")

		w.WriteHeader(http.StatusNoContent)
	})

	err := ts.client.DeleteKubernetesService("k8s-123", "us-east", 123)

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
						"name":            "pool-1",
						"worker_node":     float64(2),
						"elasticity_dict": []interface{}{},
						"policy_type":     "Manual",
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
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/kubernetes/node-pool-services/k8s-123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/kubernetes/node-pool-services/k8s-123")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"data": [
				{
					"name": "pool-1",
					"worker_node": 2
				}
			]
		}`)
	})

	result, err := ts.client.GetKubernetesNodePools("k8s-123", 123, "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestUpdateNodePoolCardinality(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/kubernetes/cluster-update/123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/kubernetes/cluster-update/123")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Node pool updated"
		}`)
	})

	resize := &models.NodePoolResize{
		NodePoolSize: 3,
	}

	result, err := ts.client.UpdateNodePoolCardinality(resize, 123, 456, "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestUpdateNodePoolCardinalityNoContent(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/kubernetes/cluster-update/123", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	resize := &models.NodePoolResize{
		NodePoolSize: 3,
	}

	result, err := ts.client.UpdateNodePoolCardinality(resize, 123, 456, "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result != nil {
		t.Fatal("Expected nil result for 204 response, got non-nil")
	}
}

func TestDeleteNodePool(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/kubernetes/delete-node-pool-service/123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		testURLPath(t, r, "/kubernetes/delete-node-pool-service/123")

		w.WriteHeader(http.StatusNoContent)
	})

	result, err := ts.client.DeleteNodePool(123, 456, "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result != nil {
		t.Fatal("Expected nil result for 204 response")
	}
}

func TestAddNodePool(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/kubernetes/add-node-pools/k8s-123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testURLPath(t, r, "/kubernetes/add-node-pools/k8s-123")

		writeJSON(w, http.StatusCreated, `{
			"code": 201,
			"message": "Node pool added"
		}`)
	})

	nodePoolAdd := &models.NodePoolAdd{
		NodePools: []models.NodePool{
			{
				Name:       "new-pool",
				WorkerNode: 2,
			},
		},
	}

	result, err := ts.client.AddNodePool(nodePoolAdd, "k8s-123", 123, "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestUpdateNodePoolDetails(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/kubernetes/update-node-pool/123/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/kubernetes/update-node-pool/123/")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Node pool details updated"
		}`)
	})

	nodePoolUpdate := &models.NodePoolUpdate{
		PlanID: "updated-plan",
	}

	result, err := ts.client.UpdateNodePoolDetails(nodePoolUpdate, 123, 456, "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestCheckNodePoolStatus(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/kubernetes/node-pool-services/k8s-123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/kubernetes/node-pool-services/k8s-123")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"data": [
				{
					"name": "pool-1",
					"status": "ACTIVE"
				}
			]
		}`)
	})

	result, err := ts.client.CheckNodePoolStatus("k8s-123", 123, "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestGetKubernetesMasterPlansError(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/kubernetes/plans", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusInternalServerError, "server error")
	})

	result, err := ts.client.GetKubernetesMasterPlans(123, "us-east")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestNewKubernetesServiceError(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/kubernetes/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusBadRequest, "invalid request")
	})

	kubernetesCreate := &models.KubernetesCreate{
		Name: "test-cluster",
	}

	result, err := ts.client.NewKubernetesService(kubernetesCreate, 123, "us-east")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}
