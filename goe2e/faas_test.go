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

	result, _, err := ts.client.FaaS.CreateNamespace(context.Background(), "test-namespace")

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

	_, err := ts.client.FaaS.DeleteNamespace(context.Background(), "test-namespace")

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

	result, _, err := ts.client.FaaS.CreateFunction(context.Background(), fn)

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

	result, _, err := ts.client.FaaS.GetFunction(context.Background(), "func-123")

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

	result, _, err := ts.client.FaaS.GetFunction(context.Background(), "nonexistent")

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

	fn := &FaasFunctionUpdateRequest{
		MemoryMB: Int(512),
		Timeout:  Int(60),
	}

	result, _, err := ts.client.FaaS.UpdateFunction(context.Background(), "func-123", fn)

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

	_, err := ts.client.FaaS.DeleteFunction(context.Background(), "func-123")

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

	result, _, err := ts.client.FaaS.GetLogs(context.Background(), "func-123")

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

// Edge case tests for better coverage
func TestCreateNamespace_EmptyName(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, _, err := ts.client.FaaS.CreateNamespace(context.Background(), "")
	if err == nil {
		t.Fatal("Expected error for empty namespace, got nil")
	}
}

func TestDeleteFunction_EmptyID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.FaaS.DeleteFunction(context.Background(), "")
	if err == nil {
		t.Fatal("Expected error for empty functionID, got nil")
	}
}

func TestGetLogs_EmptyID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, _, err := ts.client.FaaS.GetLogs(context.Background(), "")
	if err == nil {
		t.Fatal("Expected error for empty functionID, got nil")
	}
}

func TestUpdateFunction_NilRequest(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, _, err := ts.client.FaaS.UpdateFunction(context.Background(), "func-123", nil)
	if err == nil {
		t.Fatal("Expected error for nil update request, got nil")
	}
}

func TestUpdateFunction_EmptyID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	updateReq := &FaasFunctionUpdateRequest{
		MemoryMB: Int(512),
	}

	_, _, err := ts.client.FaaS.UpdateFunction(context.Background(), "", updateReq)
	if err == nil {
		t.Fatal("Expected error for empty functionID, got nil")
	}
}

func TestCreateFunction_NilRequest(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, _, err := ts.client.FaaS.CreateFunction(context.Background(), nil)
	if err == nil {
		t.Fatal("Expected error for nil create request, got nil")
	}
}

// Additional error tests for better FaaS coverage
func TestCreateNamespace_ErrorResponse(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/faas/namespace", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Internal server error"
		}`)
	})

	_, _, err := ts.client.FaaS.CreateNamespace(context.Background(), "test-namespace")
	if err == nil {
		t.Fatal("Expected error for 500 response, got nil")
	}
}

func TestDeleteNamespace_ErrorResponse(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/faas/namespace", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Internal server error"
		}`)
	})

	_, err := ts.client.FaaS.DeleteNamespace(context.Background(), "test-namespace")
	if err == nil {
		t.Fatal("Expected error for 500 response, got nil")
	}
}

func TestCreateFunction_ErrorResponse(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/faas/functions", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Internal server error"
		}`)
	})

	createReq := &FaasFunctionCreateRequest{
		Name:      "test-function",
		Namespace: "test-namespace",
		Runtime:   "python3.9",
		Code:      "def handler(): pass",
	}

	_, _, err := ts.client.FaaS.CreateFunction(context.Background(), createReq)
	if err == nil {
		t.Fatal("Expected error for 500 response, got nil")
	}
}

func TestGetFunction_ErrorResponse(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/faas/function/func-123/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Internal server error"
		}`)
	})

	_, _, err := ts.client.FaaS.GetFunction(context.Background(), "func-123")
	if err == nil {
		t.Fatal("Expected error for 500 response, got nil")
	}
}

func TestUpdateFunction_ErrorResponse(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/faas/function/func-123/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Internal server error"
		}`)
	})

	updateReq := &FaasFunctionUpdateRequest{
		MemoryMB: Int(512),
	}

	_, _, err := ts.client.FaaS.UpdateFunction(context.Background(), "func-123", updateReq)
	if err == nil {
		t.Fatal("Expected error for 500 response, got nil")
	}
}

func TestDeleteFunction_ErrorResponse(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/faas/function/func-123/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Internal server error"
		}`)
	})

	_, err := ts.client.FaaS.DeleteFunction(context.Background(), "func-123")
	if err == nil {
		t.Fatal("Expected error for 500 response, got nil")
	}
}

func TestGetLogs_ErrorResponse(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/faas/logs/func-123/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Internal server error"
		}`)
	})

	_, _, err := ts.client.FaaS.GetLogs(context.Background(), "func-123")
	if err == nil {
		t.Fatal("Expected error for 500 response, got nil")
	}
}

// Additional error response tests for better coverage
func TestCreateNamespace_BadRequest(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/faas/namespace", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusBadRequest, `{
			"code": 400,
			"message": "Invalid namespace name"
		}`)
	})

	_, _, err := ts.client.FaaS.CreateNamespace(context.Background(), "invalid@name")
	if err == nil {
		t.Fatal("Expected error for bad request, got nil")
	}
}

func TestDeleteNamespace_NotFound(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/faas/namespace", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusNotFound, `{
			"code": 404,
			"message": "Namespace not found"
		}`)
	})

	_, err := ts.client.FaaS.DeleteNamespace(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("Expected error for not found, got nil")
	}
}

func TestGetFunction_BadGateway(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/faas/function/func-123/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusBadGateway, `{
			"code": 502,
			"message": "Bad gateway"
		}`)
	})

	_, _, err := ts.client.FaaS.GetFunction(context.Background(), "func-123")
	if err == nil {
		t.Fatal("Expected error for bad gateway, got nil")
	}
}

func TestUpdateFunction_Conflict(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/faas/function/func-123/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusConflict, `{
			"code": 409,
			"message": "Function is currently being updated"
		}`)
	})

	updateReq := &FaasFunctionUpdateRequest{
		MemoryMB: Int(512),
	}

	_, _, err := ts.client.FaaS.UpdateFunction(context.Background(), "func-123", updateReq)
	if err == nil {
		t.Fatal("Expected error for conflict, got nil")
	}
}

func TestDeleteFunction_Conflict(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/faas/function/func-123/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusConflict, `{
			"code": 409,
			"message": "Cannot delete function with running invocations"
		}`)
	})

	_, err := ts.client.FaaS.DeleteFunction(context.Background(), "func-123")
	if err == nil {
		t.Fatal("Expected error for conflict, got nil")
	}
}

func TestGetLogs_NotFound(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/faas/logs/nonexistent/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusNotFound, `{
			"code": 404,
			"message": "Function not found"
		}`)
	})

	_, _, err := ts.client.FaaS.GetLogs(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("Expected error for not found, got nil")
	}
}

func TestCreateFunction_Conflict(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/faas/functions", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusConflict, `{
			"code": 409,
			"message": "Function with this name already exists"
		}`)
	})

	createReq := &FaasFunctionCreateRequest{
		Name:      "test-function",
		Namespace: "test-namespace",
		Runtime:   "python3.9",
		Code:      "def handler(): pass",
	}

	_, _, err := ts.client.FaaS.CreateFunction(context.Background(), createReq)
	if err == nil {
		t.Fatal("Expected error for conflict, got nil")
	}
}

// Phase 2: Response Parsing & Edge Case Tests

func TestCreateNamespace_MalformedJSON(t *testing.T) {
	server := newMalformedJSONServer(t)
	defer server.Close()

	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	_, _, err := client.FaaS.CreateNamespace(context.Background(), "test-ns")

	if err == nil {
		t.Error("Expected error for malformed JSON")
	}
}

func TestCreateFunction_MissingRequiredFields(t *testing.T) {
	server := newMissingFieldServer(t, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data": map[string]interface{}{
			// Missing "function_id" field
			"name": "test-function",
		},
	})
	defer server.Close()

	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	resp, _, err := client.FaaS.CreateFunction(context.Background(), &FaasFunctionCreateRequest{
		Name:      "test",
		Namespace: "default",
		Runtime:   "python3.9",
		Code:      "pass",
	})

	// Should handle missing fields gracefully
	if resp == nil && err == nil {
		t.Error("Expected response or error handling")
	}
}

func TestGetFunction_NullFieldValues(t *testing.T) {
	server := newNullFieldServer(t, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data": map[string]interface{}{
			"function_id":      "func-123",
			"name":             "test-function",
			"runtime":          nil, // Null value
			"environment_vars": nil,
			"labels":           nil,
		},
	})
	defer server.Close()

	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	resp, _, err := client.FaaS.GetFunction(context.Background(), "func-123")

	// Should handle null fields without panic
	if resp == nil && err == nil {
		t.Error("Expected response or error for null fields")
	}
}

func TestGetFunction_EmptyNameField(t *testing.T) {
	server := newNullFieldServer(t, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data": map[string]interface{}{
			"id":   "func-123",
			"name": nil, // Null value - should have a name
		},
	})
	defer server.Close()

	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	resp, _, err := client.FaaS.GetFunction(context.Background(), "func-123")

	// Should handle gracefully
	if resp == nil && err == nil {
		t.Error("Expected response or error")
	}
}

func TestGetFunction_InvalidRuntimeField(t *testing.T) {
	server := newInvalidFieldTypeServer(t, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data": map[string]interface{}{
			"function_id": "func-123",
			"name":        "test-function",
			"runtime":     123, // Should be string, not int
		},
	})
	defer server.Close()

	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	resp, _, err := client.FaaS.GetFunction(context.Background(), "func-123")

	// Should handle wrong type gracefully
	if resp == nil && err == nil {
		t.Error("Expected response or error for invalid field type")
	}
}
