package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func TestCreateFaasNamespace(t *testing.T) {
	mockResponse := models.FaasNamespaceResponse{
		Code:    201,
		Message: "Namespace created successfully",
		Data: models.FaasNamespaceData{
			Name: "test-namespace",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/faas/namespace" {
			t.Errorf("Expected path /faas/namespace, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.CreateFaasNamespace("test-namespace", "test-project", "test-location")

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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/faas/namespace" {
			t.Errorf("Expected path /faas/namespace, got %s", r.URL.Path)
		}

		query := r.URL.Query()
		if query.Get("namespace") != "test-namespace" {
			t.Errorf("Expected namespace test-namespace, got %s", query.Get("namespace"))
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.DeleteFaasNamespace("test-namespace", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestCreateFaasFunction(t *testing.T) {
	mockResponse := models.FaasFunctionResponse{
		Code:    201,
		Message: "Function created successfully",
		Data: models.FaasFunction{
			ID:     "func-123",
			Name:   "test-function",
			Status: "active",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/faas/functions" {
			t.Errorf("Expected path /faas/functions, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	fn := &models.FaasFunctionCreate{
		Name:      "test-function",
		Namespace: "test-namespace",
		Runtime:   "python3.9",
	}

	result, err := client.CreateFaasFunction(fn, "test-project", "test-location")

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
	mockResponse := models.FaasFunctionResponse{
		Code:    200,
		Message: "success",
		Data: models.FaasFunction{
			ID:     "func-123",
			Name:   "test-function",
			Status: "active",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/faas/function/func-123/" {
			t.Errorf("Expected path /faas/function/func-123/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.GetFaasFunction("func-123", "test-project", "test-location")

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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.GetFaasFunction("nonexistent", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error for 404, got: %v", err)
	}

	if result != nil {
		t.Errorf("Expected nil result for 404, got: %v", result)
	}
}

func TestUpdateFaasFunction(t *testing.T) {
	mockResponse := models.FaasFunctionResponse{
		Code:    200,
		Message: "Function updated successfully",
		Data: models.FaasFunction{
			ID:     "func-123",
			Name:   "updated-function",
			Status: "active",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/faas/function/func-123/" {
			t.Errorf("Expected path /faas/function/func-123/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	fn := &models.FaasFunctionUpdate{}

	result, err := client.UpdateFaasFunction("func-123", fn, "test-project", "test-location")

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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/faas/function/func-123/" {
			t.Errorf("Expected path /faas/function/func-123/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.DeleteFaasFunction("func-123", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestGetFaasLogs(t *testing.T) {
	mockResponse := models.FaasLogsResponse{
		Code:    200,
		Message: "success",
		Data: []models.FaasLog{
			{Message: "Log line 1"},
			{Message: "Log line 2"},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/faas/logs/func-123/" {
			t.Errorf("Expected path /faas/logs/func-123/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.GetFaasLogs("func-123", "test-project", "test-location")

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
