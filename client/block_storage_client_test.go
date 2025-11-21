package client

import (
	"net/http"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func TestNewBlockStorage(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/block_storage/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testURLPath(t, r, "/block_storage/")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Block storage created successfully",
			"data": {
				"id": "bs-123",
				"name": "test-block-storage"
			}
		}`)
	})

	blockStorageCreate := &models.BlockStorageCreate{
		Size: 100,
		Name: "test-block-storage",
	}

	result, err := ts.client.NewBlockStorage(blockStorageCreate, 123, "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result["message"] != "Block storage created successfully" {
		t.Errorf("Expected message 'Block storage created successfully', got %s", result["message"])
	}
}

func TestGetBlockStorage(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/block_storage/bs-123/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/block_storage/bs-123/")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": {
				"id": "bs-123",
				"name": "test-block-storage",
				"size": 100
			}
		}`)
	})

	result, err := ts.client.GetBlockStorage("bs-123", 123, "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result["message"] != "success" {
		t.Errorf("Expected message 'success', got %s", result["message"])
	}
}

func TestGetBlockStorage_Error(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/block_storage/bs-nonexistent/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusNotFound, `{
			"code": 404,
			"message": "Not found"
		}`)
	})

	result, err := ts.client.GetBlockStorage("bs-nonexistent", 123, "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestDeleteBlockStorage(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/block_storage/bs-123/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		testURLPath(t, r, "/block_storage/bs-123/")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "deleted"
		}`)
	})

	err := ts.client.DeleteBlockStorage("bs-123", 123, "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestUpdateBlockStorage(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/block_storage/bs-123/vm/upgrade/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/block_storage/bs-123/vm/upgrade/")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Block storage updated successfully",
			"data": {
				"id": "bs-123",
				"size": 200
			}
		}`)
	})

	blockStorageUpgrade := &models.BlockStorageUpgrade{
		Size: 200,
	}

	result, err := ts.client.UpdateBlockStorage(blockStorageUpgrade, "bs-123", 123, "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result["message"] != "Block storage updated successfully" {
		t.Errorf("Expected message 'Block storage updated successfully', got %s", result["message"])
	}
}

func TestAttachOrDetachBlockStorage(t *testing.T) {
	tests := []struct {
		name         string
		action       string
		expectedPath string
	}{
		{
			name:         "Attach block storage",
			action:       "attach",
			expectedPath: "/block_storage/bs-123/vm/attach/",
		},
		{
			name:         "Detach block storage",
			action:       "detach",
			expectedPath: "/block_storage/bs-123/vm/detach/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := setup()
			defer ts.teardown()

			ts.mux.HandleFunc(tt.expectedPath, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, http.MethodPut)
				testURLPath(t, r, tt.expectedPath)

				writeJSON(w, http.StatusOK, `{
					"code": 200,
					"message": "success"
				}`)
			})

			blockStorageAttach := &models.BlockStorageAttach{
				VM_ID: 456,
			}

			result, err := ts.client.AttachOrDetachBlockStorage(blockStorageAttach, tt.action, "bs-123", 123, "test-location")

			if err != nil {
				t.Fatalf("Expected no error, got: %v", err)
			}

			if result == nil {
				t.Fatal("Expected result, got nil")
			}
		})
	}
}

func TestGetBlockStoragePlans(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/block_storage/plans/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/block_storage/plans/")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": [
				{
					"id": "plan-1",
					"name": "Basic",
					"size": 100
				},
				{
					"id": "plan-2",
					"name": "Standard",
					"size": 500
				}
			]
		}`)
	})

	result, err := ts.client.GetBlockStoragePlans(123, "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result["message"] != "success" {
		t.Errorf("Expected message 'success', got %s", result["message"])
	}
}

func TestCheckResponseStatusForBlock(t *testing.T) {
	tests := []struct {
		name        string
		statusCode  int
		responseBody string
		expectError bool
	}{
		{
			name:         "Success - 200 OK",
			statusCode:   http.StatusOK,
			responseBody: `{"code": 200, "message": "success"}`,
			expectError:  false,
		},
		{
			name:         "Error - 400 Bad Request",
			statusCode:   http.StatusBadRequest,
			responseBody: `{"errors": "invalid request"}`,
			expectError:  true,
		},
		{
			name:         "Error - 404 Not Found",
			statusCode:   http.StatusNotFound,
			responseBody: `{"errors": "not found"}`,
			expectError:  true,
		},
		{
			name:         "Error - 500 Internal Server Error",
			statusCode:   http.StatusInternalServerError,
			responseBody: `{"errors": "server error"}`,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := setup()
			defer ts.teardown()

			ts.mux.HandleFunc("/test/", func(w http.ResponseWriter, r *http.Request) {
				writeJSON(w, tt.statusCode, tt.responseBody)
			})

			req, _ := http.NewRequest("GET", ts.server.URL+"/test/", nil)
			resp, err := http.DefaultClient.Do(req)
			if err != nil && resp != nil {
				defer resp.Body.Close()
			}
			if err == nil {
				defer resp.Body.Close()
				err = CheckResponseStatusForBlock(resp)
			}

			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}
