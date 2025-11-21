package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func TestCreateContainerRegistry(t *testing.T) {
	mockResponse := models.CreateContainerRegistryResponse{
		Code:    200,
		Message: "Container registry created successfully",
		Data: models.CreateContainerRegistryData{
			SetupStatus: "success",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/container_registry/setup-container-registry/" {
			t.Errorf("Expected path /container_registry/setup-container-registry/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	req := &models.CreateContainerRegistryRequest{
		ProjectName: "test-registry",
	}

	result, err := client.CreateContainerRegistry(req, "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.SetupStatus != mockResponse.Data.SetupStatus {
		t.Errorf("Expected SetupStatus %s, got %s", mockResponse.Data.SetupStatus, result.SetupStatus)
	}
}

func TestGetContainerRegistryProjects(t *testing.T) {
	mockResponse := models.GetContainerRegistryProjectsResponse{
		Code:    200,
		Message: "success",
		Data: []models.ContainerRegistryProject{
			{
				ProjectID:   123,
				ProjectName: "registry-1",
			},
			{
				ProjectID:   456,
				ProjectName: "registry-2",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/container_registry/projects-details/" {
			t.Errorf("Expected path /container_registry/projects-details/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.GetContainerRegistryProjects("test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 projects, got %d", len(result))
	}

	if result[0].ProjectName != "registry-1" {
		t.Errorf("Expected ProjectName registry-1, got %s", result[0].ProjectName)
	}
}

func TestDeleteContainerRegistry(t *testing.T) {
	mockResponse := models.DeleteContainerRegistryResponse{
		Code:    200,
		Message: "Deleted successfully",
		Data: models.DeleteContainerRegistryData{
			Status: "deleted",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/container_registry/setup-container-registry/" {
			t.Errorf("Expected path /container_registry/setup-container-registry/, got %s", r.URL.Path)
		}

		query := r.URL.Query()
		if query.Get("cr_project_id") != "123" {
			t.Errorf("Expected cr_project_id 123, got %s", query.Get("cr_project_id"))
		}

		if query.Get("user_id") != "user-456" {
			t.Errorf("Expected user_id user-456, got %s", query.Get("user_id"))
		}

		if query.Get("project_name") != "test-registry" {
			t.Errorf("Expected project_name test-registry, got %s", query.Get("project_name"))
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.DeleteContainerRegistry("123", "test-registry", "user-456", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestUpdateContainerRegistry(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/container_registry/setup-container-registry/" {
			t.Errorf("Expected path /container_registry/setup-container-registry/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    200,
			"message": "updated",
		})
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.UpdateContainerRegistry("test-registry", "true", "high", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}
