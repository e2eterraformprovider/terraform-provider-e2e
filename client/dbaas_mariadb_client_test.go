package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func TestCreateMariaDB(t *testing.T) {
	mockResponse := models.DBResponse{
		Code:    200,
		Message: "success",
		Data: models.DB{
			ID:     123,
			Name:   "test-mariadb",
			Status: "active",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/" {
			t.Errorf("Expected path /rds/cluster/, got %s", r.URL.Path)
		}

		// Verify query parameters
		query := r.URL.Query()
		if query.Get("apikey") == "" {
			t.Error("Expected apikey parameter")
		}
		if query.Get("location") == "" {
			t.Error("Expected location parameter")
		}
		if query.Get("project_id") == "" {
			t.Error("Expected project_id parameter")
		}

		// Read and verify request body
		body, _ := ioutil.ReadAll(r.Body)
		var req models.MariaDBCreateRequest
		json.Unmarshal(body, &req)

		if req.Name == "" {
			t.Error("Expected Name in request body")
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	req := &models.MariaDBCreateRequest{
		Name:             "test-mariadb",
		SoftwareID:       1,
		TemplateID:       100,
		PublicIPRequired: true,
		Group:            "default",
		Database: models.DBConfig{
			User:     "testuser",
			Password: "testpass",
			Name:     "testdb",
		},
	}

	result, err := client.CreateMariaDB(req, "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.ID != 123 {
		t.Errorf("Expected ID 123, got %d", result.ID)
	}

	if result.Name != "test-mariadb" {
		t.Errorf("Expected Name test-mariadb, got %s", result.Name)
	}
}

func TestCreateMariaDB_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	req := &models.MariaDBCreateRequest{
		Name:       "test-mariadb",
		SoftwareID: 1,
		TemplateID: 100,
	}

	result, err := client.CreateMariaDB(req, "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestReadMariaDB(t *testing.T) {
	mockResponse := models.DBResponse{
		Code:    200,
		Message: "success",
		Data: models.DB{
			ID:     123,
			Name:   "test-mariadb",
			Status: "active",
			Software: models.Software{
				Name:    "MariaDB",
				Version: "10.5",
				Engine:  "mariadb",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/123/" {
			t.Errorf("Expected path //rds/cluster/123/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.ReadMariaDB("123", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.ID != 123 {
		t.Errorf("Expected ID 123, got %d", result.ID)
	}

	if result.Name != "test-mariadb" {
		t.Errorf("Expected Name test-mariadb, got %s", result.Name)
	}
}

func TestReadMariaDB_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.ReadMariaDB("999", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestMariaDBExists(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		expectedExists bool
		expectError    bool
	}{
		{
			name:           "Exists - 200 OK",
			statusCode:     http.StatusOK,
			expectedExists: true,
			expectError:    false,
		},
		{
			name:           "Not Found - 404",
			statusCode:     http.StatusNotFound,
			expectedExists: false,
			expectError:    false,
		},
		{
			name:           "Error - 500",
			statusCode:     http.StatusInternalServerError,
			expectedExists: false,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Expected GET request, got %s", r.Method)
				}

				w.WriteHeader(tt.statusCode)
				if tt.statusCode == http.StatusOK {
					json.NewEncoder(w).Encode(models.DBResponse{Code: 200})
				} else {
					w.Write([]byte("Error"))
				}
			}))
			defer server.Close()

			client := NewClient("test-key", "test-token", server.URL)

			exists, err := client.MariaDBExists("123", "test-project", "test-location")

			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}

			if exists != tt.expectedExists {
				t.Errorf("Expected exists: %v, got: %v", tt.expectedExists, exists)
			}
		})
	}
}

func TestDeleteMariaDB(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/123/" {
			t.Errorf("Expected path //rds/cluster/123/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.DeleteMariaDB("123", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestDeleteMariaDB_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.DeleteMariaDB("123", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestShutdownMariaDB(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/123/shutdown" {
			t.Errorf("Expected path //rds/cluster/123/shutdown, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.ShutdownMariaDB("123", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestShutdownMariaDB_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.ShutdownMariaDB("123", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestResumeMariaDB(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/123/resume" {
			t.Errorf("Expected path //rds/cluster/123/resume, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.ResumeMariaDB("123", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestResumeMariaDB_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Forbidden"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.ResumeMariaDB("123", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestRestartMariaDB(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/123/restart" {
			t.Errorf("Expected path //rds/cluster/123/restart, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.RestartMariaDB("123", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestRestartMariaDB_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Service Unavailable"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.RestartMariaDB("123", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestAttachVPCToMariaDB(t *testing.T) {
	// Mock VPC response for ExpandMariaDBVpcList
	mockVpcResponse := models.VpcResponse{
		Code:    200,
		Message: "success",
		Data: models.Vpc{
			Name:       "test-vpc",
			Network_id: 100,
			Ipv4_cidr:  "10.0.0.0/24",
			State:      "Active",
		},
	}

	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++

		// First call is to expand VPC list (GET /vpc/100/)
		if r.Method == "GET" && callCount == 1 {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(mockVpcResponse)
			return
		}

		// Second call is to attach VPC (PUT /rds/cluster/123/vpc-attach/)
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request for attach, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/123/vpc-attach/" {
			t.Errorf("Expected path //rds/cluster/123/vpc-attach/, got %s", r.URL.Path)
		}

		// Verify request body
		body, _ := ioutil.ReadAll(r.Body)
		var req models.AttachDetachVPCRequest
		json.Unmarshal(body, &req)

		if req.Action != "attach" {
			t.Errorf("Expected action 'attach', got %s", req.Action)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.AttachVPCToMariaDB("123", "test-project", "test-location", []string{"100"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestAttachVPCToMariaDB_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return error for VPC expansion
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("VPC not found"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.AttachVPCToMariaDB("123", "test-project", "test-location", []string{"999"})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestDetachVPCFromMariaDB(t *testing.T) {
	// Mock VPC response for ExpandMariaDBVpcList
	mockVpcResponse := models.VpcResponse{
		Code:    200,
		Message: "success",
		Data: models.Vpc{
			Name:       "test-vpc",
			Network_id: 100,
			Ipv4_cidr:  "10.0.0.0/24",
			State:      "Active",
		},
	}

	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++

		// First call is to expand VPC list
		if r.Method == "GET" && callCount == 1 {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(mockVpcResponse)
			return
		}

		// Second call is to detach VPC
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request for detach, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/123/vpc-detach/" {
			t.Errorf("Expected path //rds/cluster/123/vpc-detach/, got %s", r.URL.Path)
		}

		// Verify request body
		body, _ := ioutil.ReadAll(r.Body)
		var req models.AttachDetachVPCRequest
		json.Unmarshal(body, &req)

		if req.Action != "detach" {
			t.Errorf("Expected action 'detach', got %s", req.Action)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.DetachVPCFromMariaDB("123", "test-project", "test-location", []string{"100"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestDetachVPCFromMariaDB_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return error for VPC expansion
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("VPC not found"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.DetachVPCFromMariaDB("123", "test-project", "test-location", []string{"999"})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestAttachPublicIPToMariaDB(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/123/public-ip-attach/" {
			t.Errorf("Expected path //rds/cluster/123/public-ip-attach/, got %s", r.URL.Path)
		}

		// Verify request body
		body, _ := ioutil.ReadAll(r.Body)
		var payload map[string]string
		json.Unmarshal(body, &payload)

		if payload["action"] != "attach" {
			t.Errorf("Expected action 'attach', got %s", payload["action"])
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.AttachPublicIPToMariaDB("123", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestAttachPublicIPToMariaDB_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Conflict"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.AttachPublicIPToMariaDB("123", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestDetachPublicIPFromMariaDB(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/123/public-ip-detach/" {
			t.Errorf("Expected path //rds/cluster/123/public-ip-detach/, got %s", r.URL.Path)
		}

		// Verify request body
		body, _ := ioutil.ReadAll(r.Body)
		var payload map[string]string
		json.Unmarshal(body, &payload)

		if payload["action"] != "detach" {
			t.Errorf("Expected action 'detach', got %s", payload["action"])
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.DetachPublicIPFromMariaDB("123", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestDetachPublicIPFromMariaDB_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.DetachPublicIPFromMariaDB("123", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestAttachParameterGroupToMariaDB(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/123/parameter-group/456/add" {
			t.Errorf("Expected path //rds/cluster/123/parameter-group/456/add, got %s", r.URL.Path)
		}

		// Verify request body
		body, _ := ioutil.ReadAll(r.Body)
		var req models.ParameterGroupRequest
		json.Unmarshal(body, &req)

		if req.Action != "add" {
			t.Errorf("Expected action 'add', got %s", req.Action)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.AttachParameterGroupToMariaDB("123", 456, "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestAttachParameterGroupToMariaDB_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Parameter group not found"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.AttachParameterGroupToMariaDB("123", 999, "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestDetachParameterGroupFromMariaDB(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/123/parameter-group/456/detach" {
			t.Errorf("Expected path //rds/cluster/123/parameter-group/456/detach, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.DetachParameterGroupFromMariaDB("123", 456, "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestDetachParameterGroupFromMariaDB_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server Error"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.DetachParameterGroupFromMariaDB("123", 456, "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestUpgradeMariaDBPlan(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/123/rds-upgrade/" {
			t.Errorf("Expected path //rds/cluster/123/rds-upgrade/, got %s", r.URL.Path)
		}

		// Verify request body
		body, _ := ioutil.ReadAll(r.Body)
		var payload map[string]interface{}
		json.Unmarshal(body, &payload)

		templateID := int(payload["template_id"].(float64))
		if templateID != 200 {
			t.Errorf("Expected template_id 200, got %d", templateID)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.UpgradeMariaDBPlan("123", "test-project", "test-location", 200)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestUpgradeMariaDBPlan_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusPaymentRequired)
		w.Write([]byte("Insufficient credits"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.UpgradeMariaDBPlan("123", "test-project", "test-location", 200)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestExpandMariaDBDisk(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/123/disk-upgrade/" {
			t.Errorf("Expected path //rds/cluster/123/disk-upgrade/, got %s", r.URL.Path)
		}

		// Verify request body
		body, _ := ioutil.ReadAll(r.Body)
		var req models.DiskUpgradeRequest
		json.Unmarshal(body, &req)

		if req.Size != 50 {
			t.Errorf("Expected size 50, got %d", req.Size)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.ExpandMariaDBDisk("123", "test-project", "test-location", 50)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestExpandMariaDBDisk_ZeroSize(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Expected no HTTP request for zero size")
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.ExpandMariaDBDisk("123", "test-project", "test-location", 0)

	if err != nil {
		t.Fatalf("Expected no error for zero size (should skip), got: %v", err)
	}
}

func TestExpandMariaDBDisk_NegativeSize(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Expected no HTTP request for negative size")
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.ExpandMariaDBDisk("123", "test-project", "test-location", -10)

	if err != nil {
		t.Fatalf("Expected no error for negative size (should skip), got: %v", err)
	}
}

func TestExpandMariaDBDisk_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("Cannot expand disk"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.ExpandMariaDBDisk("123", "test-project", "test-location", 50)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}
