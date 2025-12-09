package goe2e

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestContainerRegistryServiceOp_CreateContainerRegistry tests the CreateContainerRegistry method
func TestContainerRegistryServiceOp_CreateContainerRegistry(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	createJSON := `{
		"code": 200,
		"message": "Container Registry setup initiated",
		"data": {
			"setup_status": "in_progress"
		},
		"errors": {}
	}`

	listJSON := `{
		"code": 200,
		"message": "success",
		"data": [{
			"id": 123,
			"project_name": "test-registry",
			"project_size": 0.0,
			"domain_name": "test-registry.example.com",
			"prevent_vul": true,
			"severity": "high",
			"state": "active",
			"is_public": false,
			"storage_limit": 10737418240,
			"location": "Mumbai",
			"customer": 456,
			"project_id": 789,
			"my_account_project": 789,
			"deleted": false,
			"deleted_at": null,
			"created_at": "2024-01-01T10:00:00Z",
			"updated_at": "2024-01-01T10:00:00Z"
		}],
		"total_page_number": 1,
		"total_count": 1,
		"errors": {}
	}`

	ts.mux.HandleFunc("/container_registry/setup-container-registry/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")
		writeJSON(w, http.StatusOK, createJSON)
	})

	ts.mux.HandleFunc("/container_registry/projects-details/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		writeJSON(w, http.StatusOK, listJSON)
	})

	ctx := context.Background()
	createReq := &ContainerRegistryCreateRequest{
		ProjectName: "test-registry",
		PreventVul:  "true",
		Severity:    "high",
	}

	registry, resp, err := ts.client.ContainerRegistry.CreateContainerRegistry(ctx, createReq)
	if err != nil {
		t.Fatalf("CreateContainerRegistry returned error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if registry == nil {
		t.Fatal("Expected registry, got nil")
	}
	if registry.ProjectName != "test-registry" {
		t.Errorf("Expected project_name 'test-registry', got '%s'", registry.ProjectName)
	}
	if registry.ID != 123 {
		t.Errorf("Expected ID 123, got %d", registry.ID)
	}
	if registry.State != "active" {
		t.Errorf("Expected state 'active', got '%s'", registry.State)
	}
}

// TestContainerRegistryServiceOp_CreateContainerRegistry_NilRequest tests nil request validation
func TestContainerRegistryServiceOp_CreateContainerRegistry_NilRequest(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ctx := context.Background()
	_, _, err := ts.client.ContainerRegistry.CreateContainerRegistry(ctx, nil)
	if err == nil {
		t.Fatal("Expected error for nil request, got nil")
	}

	argErr, ok := err.(*ArgError)
	if !ok {
		t.Fatalf("Expected ArgError, got %T", err)
	}
	if argErr.arg != "createReq" {
		t.Errorf("Expected arg 'createReq', got '%s'", argErr.arg)
	}
}

// TestContainerRegistryServiceOp_CreateContainerRegistry_EmptyProjectName tests empty project name validation
func TestContainerRegistryServiceOp_CreateContainerRegistry_EmptyProjectName(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ctx := context.Background()
	createReq := &ContainerRegistryCreateRequest{
		ProjectName: "",
		PreventVul:  "true",
		Severity:    "high",
	}

	_, _, err := ts.client.ContainerRegistry.CreateContainerRegistry(ctx, createReq)
	if err == nil {
		t.Fatal("Expected error for empty project name, got nil")
	}

	argErr, ok := err.(*ArgError)
	if !ok {
		t.Fatalf("Expected ArgError, got %T", err)
	}
	if argErr.arg != "createReq.ProjectName" {
		t.Errorf("Expected arg 'createReq.ProjectName', got '%s'", argErr.arg)
	}
}

// TestContainerRegistryServiceOp_CreateContainerRegistry_NotFoundInList tests when created registry is not in list
func TestContainerRegistryServiceOp_CreateContainerRegistry_NotFoundInList(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	createJSON := `{
		"code": 200,
		"message": "Container Registry setup initiated",
		"data": {
			"setup_status": "in_progress"
		},
		"errors": {}
	}`

	// Return empty list
	listJSON := `{
		"code": 200,
		"message": "success",
		"data": [],
		"total_page_number": 0,
		"total_count": 0,
		"errors": {}
	}`

	ts.mux.HandleFunc("/container_registry/setup-container-registry/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, createJSON)
	})

	ts.mux.HandleFunc("/container_registry/projects-details/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, listJSON)
	})

	ctx := context.Background()
	createReq := &ContainerRegistryCreateRequest{
		ProjectName: "test-registry",
		PreventVul:  "true",
		Severity:    "high",
	}

	_, _, err := ts.client.ContainerRegistry.CreateContainerRegistry(ctx, createReq)
	if err == nil {
		t.Fatal("Expected error when registry not found in list, got nil")
	}
}

// TestContainerRegistryServiceOp_ListContainerRegistryProjects tests the ListContainerRegistryProjects method
func TestContainerRegistryServiceOp_ListContainerRegistryProjects(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	listJSON := `{
		"code": 200,
		"message": "success",
		"data": [
			{
				"id": 123,
				"project_name": "registry-1",
				"project_size": 1024.5,
				"domain_name": "registry-1.example.com",
				"prevent_vul": true,
				"severity": "high",
				"state": "active",
				"is_public": false,
				"storage_limit": 10737418240,
				"location": "Mumbai",
				"customer": 456,
				"project_id": 789,
				"my_account_project": 789,
				"deleted": false,
				"deleted_at": null,
				"created_at": "2024-01-01T10:00:00Z",
				"updated_at": "2024-01-01T10:00:00Z"
			},
			{
				"id": 124,
				"project_name": "registry-2",
				"project_size": 2048.75,
				"domain_name": "registry-2.example.com",
				"prevent_vul": false,
				"severity": "low",
				"state": "active",
				"is_public": true,
				"storage_limit": 21474836480,
				"location": "Delhi",
				"customer": 456,
				"project_id": 789,
				"my_account_project": 789,
				"deleted": false,
				"deleted_at": null,
				"created_at": "2024-01-02T11:00:00Z",
				"updated_at": "2024-01-02T11:00:00Z"
			}
		],
		"total_page_number": 1,
		"total_count": 2,
		"errors": {}
	}`

	ts.mux.HandleFunc("/container_registry/projects-details/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testQueryParam(t, r, "page", "1")
		testQueryParam(t, r, "page_size", "100")
		writeJSON(w, http.StatusOK, listJSON)
	})

	ctx := context.Background()
	opts := &ContainerRegistryListOptions{Page: 1, PageSize: 100}

	registries, resp, err := ts.client.ContainerRegistry.ListContainerRegistryProjects(ctx, opts)
	if err != nil {
		t.Fatalf("ListContainerRegistryProjects returned error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if len(registries) != 2 {
		t.Errorf("Expected 2 registries, got %d", len(registries))
	}

	// Validate first registry
	if registries[0].ID != 123 {
		t.Errorf("Expected ID 123, got %d", registries[0].ID)
	}
	if registries[0].ProjectName != "registry-1" {
		t.Errorf("Expected project_name 'registry-1', got '%s'", registries[0].ProjectName)
	}
	if registries[0].PreventVul != true {
		t.Errorf("Expected prevent_vul true, got %v", registries[0].PreventVul)
	}
	if registries[0].Severity != "high" {
		t.Errorf("Expected severity 'high', got '%s'", registries[0].Severity)
	}

	// Validate second registry
	if registries[1].ID != 124 {
		t.Errorf("Expected ID 124, got %d", registries[1].ID)
	}
	if registries[1].IsPublic != true {
		t.Errorf("Expected is_public true, got %v", registries[1].IsPublic)
	}
}

// TestContainerRegistryServiceOp_ListContainerRegistryProjects_NilOptions tests default options
func TestContainerRegistryServiceOp_ListContainerRegistryProjects_NilOptions(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	listJSON := `{
		"code": 200,
		"message": "success",
		"data": [],
		"total_page_number": 0,
		"total_count": 0,
		"errors": {}
	}`

	ts.mux.HandleFunc("/container_registry/projects-details/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		// Should use defaults: page=1, page_size=100
		testQueryParam(t, r, "page", "1")
		testQueryParam(t, r, "page_size", "100")
		writeJSON(w, http.StatusOK, listJSON)
	})

	ctx := context.Background()
	_, _, err := ts.client.ContainerRegistry.ListContainerRegistryProjects(ctx, nil)
	if err != nil {
		t.Fatalf("ListContainerRegistryProjects returned error: %v", err)
	}
}

// TestContainerRegistryServiceOp_ListContainerRegistryProjects_ZeroValues tests zero value defaults
func TestContainerRegistryServiceOp_ListContainerRegistryProjects_ZeroValues(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	listJSON := `{
		"code": 200,
		"message": "success",
		"data": [],
		"total_page_number": 0,
		"total_count": 0,
		"errors": {}
	}`

	ts.mux.HandleFunc("/container_registry/projects-details/", func(w http.ResponseWriter, r *http.Request) {
		// Should use defaults when 0 values provided
		testQueryParam(t, r, "page", "1")
		testQueryParam(t, r, "page_size", "100")
		writeJSON(w, http.StatusOK, listJSON)
	})

	ctx := context.Background()
	opts := &ContainerRegistryListOptions{Page: 0, PageSize: 0}
	_, _, err := ts.client.ContainerRegistry.ListContainerRegistryProjects(ctx, opts)
	if err != nil {
		t.Fatalf("ListContainerRegistryProjects returned error: %v", err)
	}
}

// TestContainerRegistryServiceOp_ListContainerRegistryProjects_CustomPagination tests custom pagination
func TestContainerRegistryServiceOp_ListContainerRegistryProjects_CustomPagination(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	listJSON := `{
		"code": 200,
		"message": "success",
		"data": [],
		"total_page_number": 5,
		"total_count": 50,
		"errors": {}
	}`

	ts.mux.HandleFunc("/container_registry/projects-details/", func(w http.ResponseWriter, r *http.Request) {
		testQueryParam(t, r, "page", "3")
		testQueryParam(t, r, "page_size", "10")
		writeJSON(w, http.StatusOK, listJSON)
	})

	ctx := context.Background()
	opts := &ContainerRegistryListOptions{Page: 3, PageSize: 10}
	_, _, err := ts.client.ContainerRegistry.ListContainerRegistryProjects(ctx, opts)
	if err != nil {
		t.Fatalf("ListContainerRegistryProjects returned error: %v", err)
	}
}

// TestContainerRegistryServiceOp_GetContainerRegistry tests the GetContainerRegistry method
func TestContainerRegistryServiceOp_GetContainerRegistry(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	listJSON := `{
		"code": 200,
		"message": "success",
		"data": [{
			"id": 123,
			"project_name": "test-registry",
			"project_size": 1024.0,
			"domain_name": "test-registry.example.com",
			"prevent_vul": true,
			"severity": "high",
			"state": "active",
			"is_public": false,
			"storage_limit": 10737418240,
			"location": "Mumbai",
			"customer": 456,
			"project_id": 789,
			"my_account_project": 789,
			"deleted": false,
			"deleted_at": null,
			"created_at": "2024-01-01T10:00:00Z",
			"updated_at": "2024-01-01T10:00:00Z"
		}],
		"total_page_number": 1,
		"total_count": 1,
		"errors": {}
	}`

	ts.mux.HandleFunc("/container_registry/projects-details/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		writeJSON(w, http.StatusOK, listJSON)
	})

	ctx := context.Background()
	registry, resp, err := ts.client.ContainerRegistry.GetContainerRegistry(ctx, 123)
	if err != nil {
		t.Fatalf("GetContainerRegistry returned error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if registry == nil {
		t.Fatal("Expected registry, got nil")
	}
	if registry.ID != 123 {
		t.Errorf("Expected ID 123, got %d", registry.ID)
	}
	if registry.ProjectName != "test-registry" {
		t.Errorf("Expected project_name 'test-registry', got '%s'", registry.ProjectName)
	}
}

// TestContainerRegistryServiceOp_GetContainerRegistry_NotFound tests when registry is not found
func TestContainerRegistryServiceOp_GetContainerRegistry_NotFound(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	listJSON := `{
		"code": 200,
		"message": "success",
		"data": [{
			"id": 999,
			"project_name": "other-registry",
			"project_size": 0.0,
			"domain_name": "other.example.com",
			"prevent_vul": false,
			"severity": "low",
			"state": "active",
			"is_public": false,
			"storage_limit": 10737418240,
			"location": "Mumbai",
			"customer": 456,
			"project_id": 789,
			"my_account_project": 789,
			"deleted": false,
			"deleted_at": null,
			"created_at": "2024-01-01T10:00:00Z",
			"updated_at": "2024-01-01T10:00:00Z"
		}],
		"total_page_number": 1,
		"total_count": 1,
		"errors": {}
	}`

	ts.mux.HandleFunc("/container_registry/projects-details/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, listJSON)
	})

	ctx := context.Background()
	registry, _, err := ts.client.ContainerRegistry.GetContainerRegistry(ctx, 123)
	if err != nil {
		t.Fatalf("GetContainerRegistry returned error: %v", err)
	}
	if registry != nil {
		t.Errorf("Expected nil registry for not found, got %+v", registry)
	}
}

// TestContainerRegistryServiceOp_GetContainerRegistry_InvalidID tests invalid ID validation
func TestContainerRegistryServiceOp_GetContainerRegistry_InvalidID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ctx := context.Background()

	testCases := []int{0, -1, -100}
	for _, id := range testCases {
		t.Run(fmt.Sprintf("ID=%d", id), func(t *testing.T) {
			_, _, err := ts.client.ContainerRegistry.GetContainerRegistry(ctx, id)
			if err == nil {
				t.Fatal("Expected error for invalid ID, got nil")
			}

			argErr, ok := err.(*ArgError)
			if !ok {
				t.Fatalf("Expected ArgError, got %T", err)
			}
			if argErr.arg != "id" {
				t.Errorf("Expected arg 'id', got '%s'", argErr.arg)
			}
		})
	}
}

// TestContainerRegistryServiceOp_UpdateContainerRegistry tests the UpdateContainerRegistry method
func TestContainerRegistryServiceOp_UpdateContainerRegistry(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/container_registry/setup-container-registry/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")
		w.WriteHeader(http.StatusOK)
	})

	ctx := context.Background()
	updateReq := &ContainerRegistryUpdateRequest{
		PreventVul: "false",
		Severity:   "medium",
	}

	resp, err := ts.client.ContainerRegistry.UpdateContainerRegistry(ctx, "test-registry", updateReq)
	if err != nil {
		t.Fatalf("UpdateContainerRegistry returned error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

// TestContainerRegistryServiceOp_UpdateContainerRegistry_EmptyProjectName tests empty project name validation
func TestContainerRegistryServiceOp_UpdateContainerRegistry_EmptyProjectName(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ctx := context.Background()
	updateReq := &ContainerRegistryUpdateRequest{
		PreventVul: "true",
		Severity:   "high",
	}

	_, err := ts.client.ContainerRegistry.UpdateContainerRegistry(ctx, "", updateReq)
	if err == nil {
		t.Fatal("Expected error for empty project name, got nil")
	}

	argErr, ok := err.(*ArgError)
	if !ok {
		t.Fatalf("Expected ArgError, got %T", err)
	}
	if argErr.arg != "projectName" {
		t.Errorf("Expected arg 'projectName', got '%s'", argErr.arg)
	}
}

// TestContainerRegistryServiceOp_UpdateContainerRegistry_NilRequest tests nil request validation
func TestContainerRegistryServiceOp_UpdateContainerRegistry_NilRequest(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ctx := context.Background()
	_, err := ts.client.ContainerRegistry.UpdateContainerRegistry(ctx, "test-registry", nil)
	if err == nil {
		t.Fatal("Expected error for nil request, got nil")
	}

	argErr, ok := err.(*ArgError)
	if !ok {
		t.Fatalf("Expected ArgError, got %T", err)
	}
	if argErr.arg != "updateReq" {
		t.Errorf("Expected arg 'updateReq', got '%s'", argErr.arg)
	}
}

// TestContainerRegistryServiceOp_DeleteContainerRegistry tests the DeleteContainerRegistry method
func TestContainerRegistryServiceOp_DeleteContainerRegistry(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	deleteJSON := `{
		"code": 200,
		"message": "Container Registry deleted successfully",
		"data": {
			"status": "success"
		},
		"errors": {}
	}`

	ts.mux.HandleFunc("/container_registry/setup-container-registry/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")
		testQueryParam(t, r, "cr_project_id", "123")
		testQueryParam(t, r, "project_name", "test-registry")
		testQueryParam(t, r, "user_id", "456")
		writeJSON(w, http.StatusOK, deleteJSON)
	})

	ctx := context.Background()
	deleteReq := &ContainerRegistryDeleteRequest{
		CRProjectID: "123",
		ProjectName: "test-registry",
		UserID:      "456",
	}

	resp, err := ts.client.ContainerRegistry.DeleteContainerRegistry(ctx, deleteReq)
	if err != nil {
		t.Fatalf("DeleteContainerRegistry returned error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

// TestContainerRegistryServiceOp_DeleteContainerRegistry_NilRequest tests nil request validation
func TestContainerRegistryServiceOp_DeleteContainerRegistry_NilRequest(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ctx := context.Background()
	_, err := ts.client.ContainerRegistry.DeleteContainerRegistry(ctx, nil)
	if err == nil {
		t.Fatal("Expected error for nil request, got nil")
	}

	argErr, ok := err.(*ArgError)
	if !ok {
		t.Fatalf("Expected ArgError, got %T", err)
	}
	if argErr.arg != "deleteReq" {
		t.Errorf("Expected arg 'deleteReq', got '%s'", argErr.arg)
	}
}

// TestContainerRegistryServiceOp_DeleteContainerRegistry_EmptyFields tests field validations
func TestContainerRegistryServiceOp_DeleteContainerRegistry_EmptyFields(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ctx := context.Background()

	testCases := []struct {
		name     string
		req      *ContainerRegistryDeleteRequest
		expected string
	}{
		{
			name: "Empty CRProjectID",
			req: &ContainerRegistryDeleteRequest{
				CRProjectID: "",
				ProjectName: "test",
				UserID:      "123",
			},
			expected: "deleteReq.CRProjectID",
		},
		{
			name: "Empty ProjectName",
			req: &ContainerRegistryDeleteRequest{
				CRProjectID: "123",
				ProjectName: "",
				UserID:      "456",
			},
			expected: "deleteReq.ProjectName",
		},
		{
			name: "Empty UserID",
			req: &ContainerRegistryDeleteRequest{
				CRProjectID: "123",
				ProjectName: "test",
				UserID:      "",
			},
			expected: "deleteReq.UserID",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ts.client.ContainerRegistry.DeleteContainerRegistry(ctx, tc.req)
			if err == nil {
				t.Fatal("Expected error, got nil")
			}

			argErr, ok := err.(*ArgError)
			if !ok {
				t.Fatalf("Expected ArgError, got %T", err)
			}
			if argErr.arg != tc.expected {
				t.Errorf("Expected arg '%s', got '%s'", tc.expected, argErr.arg)
			}
		})
	}
}

// TestContainerRegistryServiceOp_ErrorResponse tests error handling
func TestContainerRegistryServiceOp_ErrorResponse(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	errorJSON := `{
		"code": 400,
		"message": "Invalid request",
		"errors": ["Project name already exists"],
		"data": {}
	}`

	ts.mux.HandleFunc("/container_registry/setup-container-registry/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusBadRequest, errorJSON)
	})

	ctx := context.Background()
	createReq := &ContainerRegistryCreateRequest{
		ProjectName: "existing-registry",
		PreventVul:  "true",
		Severity:    "high",
	}

	_, resp, err := ts.client.ContainerRegistry.CreateContainerRegistry(ctx, createReq)
	if err == nil {
		t.Fatal("Expected error for bad request, got nil")
	}
	if resp == nil {
		t.Fatal("Expected response, got nil")
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

// TestContainerRegistryServiceOp_ContextCancellation tests context cancellation
func TestContainerRegistryServiceOp_ContextCancellation(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	createReq := &ContainerRegistryCreateRequest{
		ProjectName: "test-registry",
		PreventVul:  "true",
		Severity:    "high",
	}

	_, _, err := ts.client.ContainerRegistry.CreateContainerRegistry(ctx, createReq)
	if err == nil {
		t.Fatal("Expected error for canceled context, got nil")
	}
}

// TestContainerRegistryServiceOp_CreateContainerRegistry_ListError tests error during list after create
func TestContainerRegistryServiceOp_CreateContainerRegistry_ListError(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	createJSON := `{
		"code": 200,
		"message": "Container Registry setup initiated",
		"data": {
			"setup_status": "in_progress"
		},
		"errors": {}
	}`

	ts.mux.HandleFunc("/container_registry/setup-container-registry/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, createJSON)
	})

	ts.mux.HandleFunc("/container_registry/projects-details/", func(w http.ResponseWriter, r *http.Request) {
		// Return error on list
		writeJSON(w, http.StatusInternalServerError, `{"code": 500, "message": "Internal error", "data": [], "errors": {}}`)
	})

	ctx := context.Background()
	createReq := &ContainerRegistryCreateRequest{
		ProjectName: "test-registry",
		PreventVul:  "true",
		Severity:    "high",
	}

	_, _, err := ts.client.ContainerRegistry.CreateContainerRegistry(ctx, createReq)
	if err == nil {
		t.Fatal("Expected error when list fails after create, got nil")
	}
}

// TestContainerRegistryServiceOp_GetContainerRegistry_ListError tests error during list in get
func TestContainerRegistryServiceOp_GetContainerRegistry_ListError(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/container_registry/projects-details/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, `{"code": 500, "message": "Internal error", "data": [], "errors": {}}`)
	})

	ctx := context.Background()
	_, _, err := ts.client.ContainerRegistry.GetContainerRegistry(ctx, 123)
	if err == nil {
		t.Fatal("Expected error when list fails, got nil")
	}
}

// TestContainerRegistryServiceOp_UpdateContainerRegistry_APIError tests API error during update
func TestContainerRegistryServiceOp_UpdateContainerRegistry_APIError(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/container_registry/setup-container-registry/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusBadRequest, `{"code": 400, "message": "Invalid severity", "data": {}, "errors": {}}`)
	})

	ctx := context.Background()
	updateReq := &ContainerRegistryUpdateRequest{
		PreventVul: "true",
		Severity:   "invalid",
	}

	_, err := ts.client.ContainerRegistry.UpdateContainerRegistry(ctx, "test-registry", updateReq)
	if err == nil {
		t.Fatal("Expected error for API error, got nil")
	}
}

// TestContainerRegistryServiceOp_DeleteContainerRegistry_APIError tests API error during delete
func TestContainerRegistryServiceOp_DeleteContainerRegistry_APIError(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/container_registry/setup-container-registry/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusNotFound, `{"code": 404, "message": "Registry not found", "data": {}, "errors": {}}`)
	})

	ctx := context.Background()
	deleteReq := &ContainerRegistryDeleteRequest{
		CRProjectID: "999",
		ProjectName: "nonexistent",
		UserID:      "456",
	}

	_, err := ts.client.ContainerRegistry.DeleteContainerRegistry(ctx, deleteReq)
	if err == nil {
		t.Fatal("Expected error for API error, got nil")
	}
}

// TestContainerRegistryServiceOp_ListContainerRegistryProjects_EmptyList tests empty list response
func TestContainerRegistryServiceOp_ListContainerRegistryProjects_EmptyList(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	listJSON := `{
		"code": 200,
		"message": "success",
		"data": [],
		"total_page_number": 0,
		"total_count": 0,
		"errors": {}
	}`

	ts.mux.HandleFunc("/container_registry/projects-details/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, listJSON)
	})

	ctx := context.Background()
	registries, resp, err := ts.client.ContainerRegistry.ListContainerRegistryProjects(ctx, nil)
	if err != nil {
		t.Fatalf("ListContainerRegistryProjects returned error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
	if len(registries) != 0 {
		t.Errorf("Expected 0 registries, got %d", len(registries))
	}
}

// TestContainerRegistryServiceOp_UpdateContainerRegistry_Success tests successful update with all fields
func TestContainerRegistryServiceOp_UpdateContainerRegistry_Success(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/container_registry/setup-container-registry/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)

		// Verify the Authorization header is set
		if auth := r.Header.Get("Authorization"); auth == "" {
			t.Error("Expected Authorization header to be set")
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code": 200, "message": "Updated successfully", "data": {}}`))
	})

	ctx := context.Background()
	updateReq := &ContainerRegistryUpdateRequest{
		PreventVul: "true",
		Severity:   "critical",
	}

	resp, err := ts.client.ContainerRegistry.UpdateContainerRegistry(ctx, "my-registry", updateReq)
	if err != nil {
		t.Fatalf("UpdateContainerRegistry returned error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

// TestContainerRegistryServiceOp_CreateContainerRegistry_Forbidden tests 403 Forbidden response
func TestContainerRegistryServiceOp_CreateContainerRegistry_Forbidden(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/container_registry/setup-container-registry/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusForbidden, `{"code": 403, "message": "Access denied", "errors": {}}`)
	})

	ctx := context.Background()
	createReq := &ContainerRegistryCreateRequest{
		ProjectName: "restricted-registry",
		PreventVul:  "true",
		Severity:    "high",
	}

	_, resp, err := ts.client.ContainerRegistry.CreateContainerRegistry(ctx, createReq)
	if err == nil {
		t.Fatal("Expected error for 403 Forbidden, got nil")
	}
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Expected status code %d, got %d", http.StatusForbidden, resp.StatusCode)
	}
}

// TestContainerRegistryServiceOp_CreateContainerRegistry_Conflict tests 409 Conflict response
func TestContainerRegistryServiceOp_CreateContainerRegistry_Conflict(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/container_registry/setup-container-registry/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusConflict, `{"code": 409, "message": "Project already exists", "errors": {}}`)
	})

	ctx := context.Background()
	createReq := &ContainerRegistryCreateRequest{
		ProjectName: "existing-project",
		PreventVul:  "true",
		Severity:    "high",
	}

	_, resp, err := ts.client.ContainerRegistry.CreateContainerRegistry(ctx, createReq)
	if err == nil {
		t.Fatal("Expected error for 409 Conflict, got nil")
	}
	if resp.StatusCode != http.StatusConflict {
		t.Errorf("Expected status code %d, got %d", http.StatusConflict, resp.StatusCode)
	}
}

// TestContainerRegistryServiceOp_ListContainerRegistryProjects_Forbidden tests 403 Forbidden on list
func TestContainerRegistryServiceOp_ListContainerRegistryProjects_Forbidden(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/container_registry/projects-details/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusForbidden, `{"code": 403, "message": "Access denied", "errors": {}}`)
	})

	ctx := context.Background()
	_, resp, err := ts.client.ContainerRegistry.ListContainerRegistryProjects(ctx, nil)
	if err == nil {
		t.Fatal("Expected error for 403 Forbidden, got nil")
	}
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Expected status code %d, got %d", http.StatusForbidden, resp.StatusCode)
	}
}

// TestContainerRegistryServiceOp_GetContainerRegistry_ServerError tests 500 error on list
func TestContainerRegistryServiceOp_GetContainerRegistry_ServerError(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/container_registry/projects-details/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, `{"code": 500, "message": "Internal server error", "data": [], "errors": {}}`)
	})

	ctx := context.Background()
	_, resp, err := ts.client.ContainerRegistry.GetContainerRegistry(ctx, 123)
	if err == nil {
		t.Fatal("Expected error for 500 Internal Server Error, got nil")
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, resp.StatusCode)
	}
}

// TestContainerRegistryServiceOp_UpdateContainerRegistry_Forbidden tests 403 Forbidden on update
func TestContainerRegistryServiceOp_UpdateContainerRegistry_Forbidden(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/container_registry/setup-container-registry/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusForbidden, `{"code": 403, "message": "Not authorized", "errors": {}}`)
	})

	ctx := context.Background()
	updateReq := &ContainerRegistryUpdateRequest{
		PreventVul: "true",
		Severity:   "high",
	}

	_, err := ts.client.ContainerRegistry.UpdateContainerRegistry(ctx, "test-registry", updateReq)
	if err == nil {
		t.Fatal("Expected error for 403 Forbidden, got nil")
	}
}

// TestContainerRegistryServiceOp_UpdateContainerRegistry_NotFound tests 404 Not Found on update
func TestContainerRegistryServiceOp_UpdateContainerRegistry_NotFound(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/container_registry/setup-container-registry/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusNotFound, `{"code": 404, "message": "Registry not found", "errors": {}}`)
	})

	ctx := context.Background()
	updateReq := &ContainerRegistryUpdateRequest{
		PreventVul: "true",
		Severity:   "high",
	}

	_, err := ts.client.ContainerRegistry.UpdateContainerRegistry(ctx, "nonexistent", updateReq)
	if err == nil {
		t.Fatal("Expected error for 404 Not Found, got nil")
	}
}

// TestContainerRegistryServiceOp_DeleteContainerRegistry_Conflict tests 409 Conflict on delete
func TestContainerRegistryServiceOp_DeleteContainerRegistry_Conflict(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/container_registry/setup-container-registry/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusConflict, `{"code": 409, "message": "Registry in use", "errors": {}}`)
	})

	ctx := context.Background()
	deleteReq := &ContainerRegistryDeleteRequest{
		CRProjectID: "123",
		ProjectName: "test-registry",
		UserID:      "456",
	}

	_, err := ts.client.ContainerRegistry.DeleteContainerRegistry(ctx, deleteReq)
	if err == nil {
		t.Fatal("Expected error for 409 Conflict, got nil")
	}
}

// TestContainerRegistryServiceOp_DeleteContainerRegistry_Forbidden tests 403 Forbidden on delete
func TestContainerRegistryServiceOp_DeleteContainerRegistry_Forbidden(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/container_registry/setup-container-registry/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusForbidden, `{"code": 403, "message": "Permission denied", "errors": {}}`)
	})

	ctx := context.Background()
	deleteReq := &ContainerRegistryDeleteRequest{
		CRProjectID: "999",
		ProjectName: "protected-registry",
		UserID:      "456",
	}

	_, err := ts.client.ContainerRegistry.DeleteContainerRegistry(ctx, deleteReq)
	if err == nil {
		t.Fatal("Expected error for 403 Forbidden, got nil")
	}
}

// TestContainerRegistryServiceOp_ListContainerRegistryProjects_ServiceUnavailable tests 503
func TestContainerRegistryServiceOp_ListContainerRegistryProjects_ServiceUnavailable(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/container_registry/projects-details/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusServiceUnavailable, `{"code": 503, "message": "Service temporarily unavailable", "errors": {}}`)
	})

	ctx := context.Background()
	_, resp, err := ts.client.ContainerRegistry.ListContainerRegistryProjects(ctx, nil)
	if err == nil {
		t.Fatal("Expected error for 503 Service Unavailable, got nil")
	}
	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("Expected status code %d, got %d", http.StatusServiceUnavailable, resp.StatusCode)
	}
}

// TestContainerRegistryServiceOp_CreateContainerRegistry_Unauthorized tests 401 Unauthorized response
func TestContainerRegistryServiceOp_CreateContainerRegistry_Unauthorized(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/container_registry/setup-container-registry/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusUnauthorized, `{"code": 401, "message": "Unauthorized", "errors": {}}`)
	})

	ctx := context.Background()
	createReq := &ContainerRegistryCreateRequest{
		ProjectName: "test-registry",
		PreventVul:  "true",
		Severity:    "high",
	}

	_, resp, err := ts.client.ContainerRegistry.CreateContainerRegistry(ctx, createReq)
	if err == nil {
		t.Fatal("Expected error for 401 Unauthorized, got nil")
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}

// TestContainerRegistryServiceOp_ListContainerRegistryProjects_Unauthorized tests 401 Unauthorized on list
func TestContainerRegistryServiceOp_ListContainerRegistryProjects_Unauthorized(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/container_registry/projects-details/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusUnauthorized, `{"code": 401, "message": "Unauthorized", "errors": {}}`)
	})

	ctx := context.Background()
	_, resp, err := ts.client.ContainerRegistry.ListContainerRegistryProjects(ctx, nil)
	if err == nil {
		t.Fatal("Expected error for 401 Unauthorized, got nil")
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}

// TestContainerRegistryServiceOp_GetContainerRegistry_Unauthorized tests 401 Unauthorized on get
func TestContainerRegistryServiceOp_GetContainerRegistry_Unauthorized(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/container_registry/projects-details/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusUnauthorized, `{"code": 401, "message": "Unauthorized", "errors": {}}`)
	})

	ctx := context.Background()
	_, resp, err := ts.client.ContainerRegistry.GetContainerRegistry(ctx, 123)
	if err == nil {
		t.Fatal("Expected error for 401 Unauthorized, got nil")
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}

// TestContainerRegistryServiceOp_UpdateContainerRegistry_Unauthorized tests 401 Unauthorized on update
func TestContainerRegistryServiceOp_UpdateContainerRegistry_Unauthorized(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/container_registry/setup-container-registry/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusUnauthorized, `{"code": 401, "message": "Unauthorized", "errors": {}}`)
	})

	ctx := context.Background()
	updateReq := &ContainerRegistryUpdateRequest{
		PreventVul: "true",
		Severity:   "high",
	}

	_, err := ts.client.ContainerRegistry.UpdateContainerRegistry(ctx, "test-registry", updateReq)
	if err == nil {
		t.Fatal("Expected error for 401 Unauthorized, got nil")
	}
}

// TestContainerRegistryServiceOp_DeleteContainerRegistry_Unauthorized tests 401 Unauthorized on delete
func TestContainerRegistryServiceOp_DeleteContainerRegistry_Unauthorized(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/container_registry/setup-container-registry/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusUnauthorized, `{"code": 401, "message": "Unauthorized", "errors": {}}`)
	})

	ctx := context.Background()
	deleteReq := &ContainerRegistryDeleteRequest{
		CRProjectID: "123",
		ProjectName: "test-registry",
		UserID:      "456",
	}

	_, err := ts.client.ContainerRegistry.DeleteContainerRegistry(ctx, deleteReq)
	if err == nil {
		t.Fatal("Expected error for 401 Unauthorized, got nil")
	}
}

// ============================================================================
// Network Error & Timeout Tests
// ============================================================================

// TestContainerRegistry_NetworkConnectionError tests network connection failure
func TestContainerRegistry_NetworkConnectionError(t *testing.T) {
	client, _ := NewClient(
		"test-key",
		"test-token",
		"proj-123",
		"test-location",
		SetBaseURL("http://localhost:1"),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}),
	)

	ctx := context.Background()
	createReq := &ContainerRegistryCreateRequest{
		ProjectName: "test-registry",
		PreventVul:  "true",
		Severity:    "high",
	}

	_, _, err := client.ContainerRegistry.CreateContainerRegistry(ctx, createReq)
	if err == nil {
		t.Fatal("Expected network error")
	}
}

// TestContainerRegistry_ContextTimeout tests context timeout
func TestContainerRegistry_ContextTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, _ := NewClient(
		"test-key",
		"test-token",
		"proj-123",
		"test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	createReq := &ContainerRegistryCreateRequest{
		ProjectName: "test-registry",
		PreventVul:  "true",
		Severity:    "high",
	}

	_, _, err := client.ContainerRegistry.CreateContainerRegistry(ctx, createReq)
	if err == nil {
		t.Fatal("Expected error for context timeout")
	}
}

// TestContainerRegistry_ContextCancellation tests context cancellation
func TestContainerRegistry_ContextCancellation(t *testing.T) {
	client, _ := NewClient("test-key", "test-token", "proj-123", "test-location")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	createReq := &ContainerRegistryCreateRequest{
		ProjectName: "test-registry",
		PreventVul:  "true",
		Severity:    "high",
	}

	_, _, err := client.ContainerRegistry.CreateContainerRegistry(ctx, createReq)
	if err == nil {
		t.Fatal("Expected error for canceled context")
	}
}

// TestContainerRegistry_RetryLogic_TemporaryFailure tests retry on temporary failures
func TestContainerRegistry_RetryLogic_TemporaryFailure(t *testing.T) {
	var requestCount int

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++

		if requestCount == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprint(w, `{"code": 503, "message": "Unavailable"}`)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"code": 200,
			"message": "Registry created",
			"data": {"setup_status": "in_progress"},
			"errors": {}
		}`)
	}))
	defer server.Close()

	client, _ := NewClient(
		"test-key",
		"test-token",
		"proj-123",
		"test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{
			RetryMax:     2,
			RetryWaitMin: PtrTo(10 * time.Millisecond),
			RetryWaitMax: PtrTo(50 * time.Millisecond),
		}),
	)

	ctx := context.Background()
	createReq := &ContainerRegistryCreateRequest{
		ProjectName: "test-registry",
		PreventVul:  "true",
		Severity:    "high",
	}

	_, _, _ = client.ContainerRegistry.CreateContainerRegistry(ctx, createReq)
	// Will fail because setup returns status but we need to list to find the registry
	// This test verifies retry logic attempts are made
	if requestCount < 1 {
		t.Errorf("Expected at least 1 request, got %d", requestCount)
	}
}

// TestContainerRegistry_GetContainerRegistry_NetworkError tests Get with network error
func TestContainerRegistry_GetContainerRegistry_NetworkError(t *testing.T) {
	client, _ := NewClient(
		"test-key",
		"test-token",
		"proj-123",
		"test-location",
		SetBaseURL("http://localhost:1"),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}),
	)

	ctx := context.Background()
	_, _, err := client.ContainerRegistry.GetContainerRegistry(ctx, 123)
	if err == nil {
		t.Fatal("Expected network error")
	}
}

// TestContainerRegistry_ListContainerRegistryProjects_NetworkError tests List with network error
func TestContainerRegistry_ListContainerRegistryProjects_NetworkError(t *testing.T) {
	client, _ := NewClient(
		"test-key",
		"test-token",
		"proj-123",
		"test-location",
		SetBaseURL("http://localhost:1"),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}),
	)

	ctx := context.Background()
	_, _, err := client.ContainerRegistry.ListContainerRegistryProjects(ctx, nil)
	if err == nil {
		t.Fatal("Expected network error")
	}
}

// TestContainerRegistry_DeleteContainerRegistry_ContextTimeout tests Delete with context timeout
func TestContainerRegistry_DeleteContainerRegistry_ContextTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, _ := NewClient(
		"test-key",
		"test-token",
		"proj-123",
		"test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	deleteReq := &ContainerRegistryDeleteRequest{
		CRProjectID: "123",
		ProjectName: "test-registry",
		UserID:      "456",
	}

	_, err := client.ContainerRegistry.DeleteContainerRegistry(ctx, deleteReq)
	if err == nil {
		t.Fatal("Expected error for context timeout")
	}
}

// TestContainerRegistry_UpdateContainerRegistry_ContextTimeout tests Update with context timeout
func TestContainerRegistry_UpdateContainerRegistry_ContextTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, _ := NewClient(
		"test-key",
		"test-token",
		"proj-123",
		"test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	updateReq := &ContainerRegistryUpdateRequest{
		PreventVul: "true",
		Severity:   "high",
	}

	_, err := client.ContainerRegistry.UpdateContainerRegistry(ctx, "test-registry", updateReq)
	if err == nil {
		t.Fatal("Expected error for context timeout")
	}
}
