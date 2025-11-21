package client

import (
	"net/http"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func TestCreateContainerRegistry(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/container_registry/setup-container-registry/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testURLPath(t, r, "/container_registry/setup-container-registry/")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Container registry created successfully",
			"data": {
				"setup_status": "success"
			}
		}`)
	})

	req := &models.CreateContainerRegistryRequest{
		ProjectName: "test-registry",
	}

	result, err := ts.client.CreateContainerRegistry(req, "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.SetupStatus != "success" {
		t.Errorf("Expected SetupStatus success, got %s", result.SetupStatus)
	}
}

func TestGetContainerRegistryProjects(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/container_registry/projects-details/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/container_registry/projects-details/")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": [
				{
					"project_id": 123,
					"project_name": "registry-1"
				},
				{
					"project_id": 456,
					"project_name": "registry-2"
				}
			]
		}`)
	})

	result, err := ts.client.GetContainerRegistryProjects("test-project", "test-location")

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
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/container_registry/setup-container-registry/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		testURLPath(t, r, "/container_registry/setup-container-registry/")
		testQueryParam(t, r, "cr_project_id", "123")
		testQueryParam(t, r, "user_id", "user-456")
		testQueryParam(t, r, "project_name", "test-registry")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Deleted successfully",
			"data": {
				"status": "deleted"
			}
		}`)
	})

	err := ts.client.DeleteContainerRegistry("123", "test-registry", "user-456", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestUpdateContainerRegistry(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/container_registry/setup-container-registry/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/container_registry/setup-container-registry/")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "updated"
		}`)
	})

	err := ts.client.UpdateContainerRegistry("test-registry", "true", "high", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}
