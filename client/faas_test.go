package client

import (
	"net/http"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func TestCreateFaasNamespace(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/faas/namespace", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testURLPath(t, r, "/faas/namespace")

		writeJSON(w, http.StatusCreated, `{
			"code": 201,
			"message": "Namespace created successfully",
			"data": {
				"name": "test-namespace"
			}
		}`)
	})

	result, err := ts.client.CreateFaasNamespace("test-namespace", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Data.Name != "test-namespace" {
		t.Errorf("Expected Name test-namespace, got %s", result.Data.Name)
	}
}

func TestDeleteFaasNamespace(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/faas/namespace", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		testURLPath(t, r, "/faas/namespace")
		testQueryParam(t, r, "namespace", "test-namespace")

		w.WriteHeader(http.StatusOK)
	})

	err := ts.client.DeleteFaasNamespace("test-namespace", "test-project", "test-location")

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

		writeJSON(w, http.StatusCreated, `{
			"code": 201,
			"message": "Function created successfully",
			"data": {
				"id": "func-123",
				"name": "test-function",
				"status": "active"
			}
		}`)
	})

	fn := &models.FaasFunctionCreate{
		Name:      "test-function",
		Namespace: "test-namespace",
		Runtime:   "python3.9",
	}

	result, err := ts.client.CreateFaasFunction(fn, "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Data.ID != "func-123" {
		t.Errorf("Expected ID func-123, got %s", result.Data.ID)
	}
}

func TestGetFaasFunction(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/faas/function/func-123/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/faas/function/func-123/")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": {
				"id": "func-123",
				"name": "test-function",
				"status": "active"
			}
		}`)
	})

	result, err := ts.client.GetFaasFunction("func-123", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Data.ID != "func-123" {
		t.Errorf("Expected ID func-123, got %s", result.Data.ID)
	}
}

func TestGetFaasFunction_NotFound(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/faas/function/nonexistent/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	result, err := ts.client.GetFaasFunction("nonexistent", "test-project", "test-location")

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

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Function updated successfully",
			"data": {
				"id": "func-123",
				"name": "updated-function",
				"status": "active"
			}
		}`)
	})

	fn := &models.FaasFunctionUpdate{}

	result, err := ts.client.UpdateFaasFunction("func-123", fn, "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Data.Name != "updated-function" {
		t.Errorf("Expected Name updated-function, got %s", result.Data.Name)
	}
}

func TestDeleteFaasFunction(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/faas/function/func-123/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		testURLPath(t, r, "/faas/function/func-123/")

		w.WriteHeader(http.StatusOK)
	})

	err := ts.client.DeleteFaasFunction("func-123", "test-project", "test-location")

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

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": [
				{"message": "Log line 1"},
				{"message": "Log line 2"}
			]
		}`)
	})

	result, err := ts.client.GetFaasLogs("func-123", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if len(result.Data) != 2 {
		t.Errorf("Expected 2 log lines, got %d", len(result.Data))
	}
}
