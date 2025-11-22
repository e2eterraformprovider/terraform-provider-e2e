package goe2e

import (
	"context"
	"net/http"
	"testing"
)

func TestCreateFaasNamespace(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/faas/namespace", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testURLPath(t, r, "/faas/namespace")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		writeJSON(w, http.StatusCreated, `{
			"code": 201,
			"message": "Namespace created successfully",
			"data": {
				"name": "test-namespace"
			}
		}`)
	})

	opts := &RequestOptions{
		ProjectID: "test-project",
		Location:  "test-location",
	}

	result, _, err := ts.client.FaaS.CreateNamespace(context.Background(), "test-namespace", opts)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Name != "test-namespace" {
		t.Errorf("Expected Name test-namespace, got %s", result.Name)
	}
}

func TestDeleteFaasNamespace(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/faas/namespace", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		testURLPath(t, r, "/faas/namespace")
		testQueryParam(t, r, "namespace", "test-namespace")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		w.WriteHeader(http.StatusOK)
	})

	opts := &RequestOptions{
		ProjectID: "test-project",
		Location:  "test-location",
	}

	_, err := ts.client.FaaS.DeleteNamespace(context.Background(), "test-namespace", opts)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestCreateFaasFunction(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/faas/functions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testURLPath(t, r, "/faas/functions")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		writeJSON(w, http.StatusCreated, `{
			"code": 201,
			"message": "Function created successfully",
			"data": {
				"id": "func-123",
				"name": "test-function",
				"namespace": "test-namespace",
				"runtime": "python3.9",
				"status": "active",
				"memory_mb": 256,
				"timeout_seconds": 30,
				"min_replicas": 1,
				"max_replicas": 5,
				"endpoint_url": "https://func.example.com",
				"created_at": "2024-01-01T00:00:00Z",
				"updated_at": "2024-01-01T00:00:00Z"
			}
		}`)
	})

	opts := &RequestOptions{
		ProjectID: "test-project",
		Location:  "test-location",
	}

	fn := &FaasFunctionCreateRequest{
		Name:        "test-function",
		Namespace:   "test-namespace",
		Runtime:     "python3.9",
		Code:        "def handler(event): return event",
		MemoryMB:    256,
		Timeout:     30,
		MinReplicas: 1,
		MaxReplicas: 5,
	}

	result, _, err := ts.client.FaaS.CreateFunction(context.Background(), fn, opts)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.ID != "func-123" {
		t.Errorf("Expected ID func-123, got %s", result.ID)
	}

	if result.Name != "test-function" {
		t.Errorf("Expected Name test-function, got %s", result.Name)
	}
}

func TestGetFaasFunction(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/faas/function/func-123/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/faas/function/func-123/")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": {
				"id": "func-123",
				"name": "test-function",
				"namespace": "test-namespace",
				"runtime": "python3.9",
				"status": "active",
				"memory_mb": 256,
				"timeout_seconds": 30,
				"min_replicas": 1,
				"max_replicas": 5,
				"endpoint_url": "https://func.example.com",
				"created_at": "2024-01-01T00:00:00Z",
				"updated_at": "2024-01-01T00:00:00Z"
			}
		}`)
	})

	opts := &RequestOptions{
		ProjectID: "test-project",
		Location:  "test-location",
	}

	result, _, err := ts.client.FaaS.GetFunction(context.Background(), "func-123", opts)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.ID != "func-123" {
		t.Errorf("Expected ID func-123, got %s", result.ID)
	}
}

func TestGetFaasFunction_NotFound(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/faas/function/nonexistent/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	opts := &RequestOptions{
		ProjectID: "test-project",
		Location:  "test-location",
	}

	result, _, err := ts.client.FaaS.GetFunction(context.Background(), "nonexistent", opts)

	if err != nil {
		t.Fatalf("Expected no error for 404, got: %v", err)
	}

	if result != nil {
		t.Errorf("Expected nil result for 404, got: %v", result)
	}
}

func TestUpdateFaasFunction(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/faas/function/func-123/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/faas/function/func-123/")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Function updated successfully",
			"data": {
				"id": "func-123",
				"name": "test-function",
				"namespace": "test-namespace",
				"runtime": "python3.9",
				"status": "active",
				"memory_mb": 512,
				"timeout_seconds": 60,
				"min_replicas": 1,
				"max_replicas": 5,
				"endpoint_url": "https://func.example.com",
				"created_at": "2024-01-01T00:00:00Z",
				"updated_at": "2024-01-01T01:00:00Z"
			}
		}`)
	})

	opts := &RequestOptions{
		ProjectID: "test-project",
		Location:  "test-location",
	}

	fn := &FaasFunctionUpdateRequest{
		MemoryMB: Int(512),
		Timeout:  Int(60),
	}

	result, _, err := ts.client.FaaS.UpdateFunction(context.Background(), "func-123", fn, opts)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.MemoryMB != 512 {
		t.Errorf("Expected MemoryMB 512, got %d", result.MemoryMB)
	}

	if result.Timeout != 60 {
		t.Errorf("Expected Timeout 60, got %d", result.Timeout)
	}
}

func TestDeleteFaasFunction(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/faas/function/func-123/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		testURLPath(t, r, "/faas/function/func-123/")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		w.WriteHeader(http.StatusOK)
	})

	opts := &RequestOptions{
		ProjectID: "test-project",
		Location:  "test-location",
	}

	_, err := ts.client.FaaS.DeleteFunction(context.Background(), "func-123", opts)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestGetFaasLogs(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/faas/logs/func-123/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/faas/logs/func-123/")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": [
				{"message": "Log line 1", "timestamp": "2024-01-01T00:00:00Z"},
				{"message": "Log line 2", "timestamp": "2024-01-01T00:00:01Z"}
			]
		}`)
	})

	opts := &RequestOptions{
		ProjectID: "test-project",
		Location:  "test-location",
	}

	result, _, err := ts.client.FaaS.GetLogs(context.Background(), "func-123", opts)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if len(result.Logs) != 2 {
		t.Errorf("Expected 2 log lines, got %d", len(result.Logs))
	}

	if result.Logs[0].Message != "Log line 1" {
		t.Errorf("Expected first log message 'Log line 1', got %s", result.Logs[0].Message)
	}
}
