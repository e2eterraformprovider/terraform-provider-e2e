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

func TestNewBlockStorage(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code":    200,
		"message": "Block storage created successfully",
		"data": map[string]interface{}{
			"id":   "bs-123",
			"name": "test-block-storage",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/block_storage/" {
			t.Errorf("Expected path /block_storage/, got %s", r.URL.Path)
		}

		// Verify request body
		body, _ := ioutil.ReadAll(r.Body)
		var blockStorageCreate models.BlockStorageCreate
		json.Unmarshal(body, &blockStorageCreate)

		if blockStorageCreate.Size == 0 {
			t.Error("Expected Size in request body")
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	blockStorageCreate := &models.BlockStorageCreate{
		Size: 100,
		Name: "test-block-storage",
	}

	result, err := client.NewBlockStorage(blockStorageCreate, 123, "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result["message"] != mockResponse["message"] {
		t.Errorf("Expected message %s, got %s", mockResponse["message"], result["message"])
	}
}

func TestGetBlockStorage(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code":    200,
		"message": "success",
		"data": map[string]interface{}{
			"id":   "bs-123",
			"name": "test-block-storage",
			"size": float64(100),
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/block_storage/bs-123/" {
			t.Errorf("Expected path /block_storage/bs-123/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.GetBlockStorage("bs-123", 123, "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result["message"] != mockResponse["message"] {
		t.Errorf("Expected message %s, got %s", mockResponse["message"], result["message"])
	}
}

func TestGetBlockStorage_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    404,
			"message": "Not found",
		})
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.GetBlockStorage("bs-nonexistent", 123, "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestDeleteBlockStorage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/block_storage/bs-123/" {
			t.Errorf("Expected path /block_storage/bs-123/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    200,
			"message": "deleted",
		})
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.DeleteBlockStorage("bs-123", 123, "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestUpdateBlockStorage(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code":    200,
		"message": "Block storage updated successfully",
		"data": map[string]interface{}{
			"id":   "bs-123",
			"size": float64(200),
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/block_storage/bs-123/vm/upgrade/" {
			t.Errorf("Expected path /block_storage/bs-123/vm/upgrade/, got %s", r.URL.Path)
		}

		// Verify request body
		body, _ := ioutil.ReadAll(r.Body)
		var blockStorageUpgrade models.BlockStorageUpgrade
		json.Unmarshal(body, &blockStorageUpgrade)

		if blockStorageUpgrade.Size == 0 {
			t.Error("Expected Size in request body")
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	blockStorageUpgrade := &models.BlockStorageUpgrade{
		Size: 200,
	}

	result, err := client.UpdateBlockStorage(blockStorageUpgrade, "bs-123", 123, "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result["message"] != mockResponse["message"] {
		t.Errorf("Expected message %s, got %s", mockResponse["message"], result["message"])
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
			mockResponse := map[string]interface{}{
				"code":    200,
				"message": "success",
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "PUT" {
					t.Errorf("Expected PUT request, got %s", r.Method)
				}

				if r.URL.Path != tt.expectedPath {
					t.Errorf("Expected path %s, got %s", tt.expectedPath, r.URL.Path)
				}

				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(mockResponse)
			}))
			defer server.Close()

			client := NewClient("test-key", "test-token", server.URL)

			blockStorageAttach := &models.BlockStorageAttach{
				VM_ID: 456,
			}

			result, err := client.AttachOrDetachBlockStorage(blockStorageAttach, tt.action, "bs-123", 123, "test-location")

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
	mockResponse := map[string]interface{}{
		"code":    200,
		"message": "success",
		"data": []interface{}{
			map[string]interface{}{
				"id":   "plan-1",
				"name": "Basic",
				"size": float64(100),
			},
			map[string]interface{}{
				"id":   "plan-2",
				"name": "Standard",
				"size": float64(500),
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/block_storage/plans/" {
			t.Errorf("Expected path /block_storage/plans/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.GetBlockStoragePlans(123, "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result["message"] != mockResponse["message"] {
		t.Errorf("Expected message %s, got %s", mockResponse["message"], result["message"])
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
			// Create a mock response
			resp := &http.Response{
				StatusCode: tt.statusCode,
				Body:       ioutil.NopCloser(bytes.NewBufferString(tt.responseBody)),
			}

			err := CheckResponseStatusForBlock(resp)

			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}
