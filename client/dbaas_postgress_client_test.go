package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func TestCreatePostgressDB(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code":    float64(200),
		"message": "PostgreSQL database created successfully",
		"data": map[string]interface{}{
			"id":   float64(789),
			"name": "test-postgres",
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
		var req models.DBCreateRequest
		json.Unmarshal(body, &req)

		if req.Name == "" {
			t.Error("Expected Name in request body")
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	dbCreate := models.DBCreateRequest{
		Name:             "test-postgres",
		SoftwareID:       3,
		TemplateID:       103,
		PublicIPRequired: true,
		Group:            "default",
		Database: models.DBConfig{
			User:     "pguser",
			Password: "pgpass",
			Name:     "pgdb",
		},
	}

	result, err := client.CreatePostgressDB(dbCreate, "test-project", "test-location")

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

func TestCreatePostgressDB_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	dbCreate := models.DBCreateRequest{
		Name:       "test-postgres",
		SoftwareID: 3,
		TemplateID: 103,
	}

	result, err := client.CreatePostgressDB(dbCreate, "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestGetPostgressDB(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code":    float64(200),
		"message": "success",
		"data": map[string]interface{}{
			"id":     float64(789),
			"name":   "test-postgres",
			"status": "active",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/789/" {
			t.Errorf("Expected path /rds/cluster/789/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.GetPostgressDB("789", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	data := result["data"].(map[string]interface{})
	if data["name"] != "test-postgres" {
		t.Errorf("Expected name test-postgres, got %s", data["name"])
	}
}

func TestGetPostgressDB_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.GetPostgressDB("999", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestDeletePostgressDB(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/789/" {
			t.Errorf("Expected path /rds/cluster/789/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.DeletePostgressDB("789", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestDeletePostgressDB_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server Error"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.DeletePostgressDB("789", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestStopPostgressDB(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/789/shutdown" {
			t.Errorf("Expected path /rds/cluster/789/shutdown, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.StopPostgressDB("789", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestStopPostgressDB_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.StopPostgressDB("789", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestStartPostgressDB(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/789/resume" {
			t.Errorf("Expected path /rds/cluster/789/resume, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.StartPostgressDB("789", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestStartPostgressDB_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Forbidden"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.StartPostgressDB("789", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestRestartPostgressDB(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/789/restart" {
			t.Errorf("Expected path /rds/cluster/789/restart, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.RestartPostgressDB("789", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestRestartPostgressDB_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Service Unavailable"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.RestartPostgressDB("789", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestAttachPublicIpPostgressDB(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/789/public-ip-attach/" {
			t.Errorf("Expected path /rds/cluster/789/public-ip-attach/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.AttachPublicIpPostgressDB("789", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestAttachPublicIpPostgressDB_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Conflict"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.AttachPublicIpPostgressDB("789", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestDetachPublicIpPostgressDB(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/789/public-ip-detach/" {
			t.Errorf("Expected path /rds/cluster/789/public-ip-detach/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.DetachPublicIpPostgressDB("789", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestDetachPublicIpPostgressDB_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.DetachPublicIpPostgressDB("789", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestAttachVPCPostgressDB(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/789/vpc-attach/" {
			t.Errorf("Expected path /rds/cluster/789/vpc-attach/, got %s", r.URL.Path)
		}

		// Verify request body
		body, _ := ioutil.ReadAll(r.Body)
		var req models.AttachVPCPayloadRequest
		json.Unmarshal(body, &req)

		if req.Action == "" {
			t.Error("Expected Action in request body")
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	vpcPayload := models.AttachVPCPayloadRequest{
		Action: "attach",
		VPCs: []models.VPC{
			{
				VpcName:    "test-vpc",
				Ipv4_cidr:  "10.0.0.0/24",
				Network_id: 100,
			},
		},
	}

	err := client.AttachVPCPostgressDB(vpcPayload, "789", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestAttachVPCPostgressDB_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("VPC not found"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	vpcPayload := models.AttachVPCPayloadRequest{
		Action: "attach",
		VPCs:   []models.VPC{},
	}

	err := client.AttachVPCPostgressDB(vpcPayload, "789", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestDetachVPCPostgressDB(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/789/vpc-detach/" {
			t.Errorf("Expected path /rds/cluster/789/vpc-detach/, got %s", r.URL.Path)
		}

		// Verify request body
		body, _ := ioutil.ReadAll(r.Body)
		var req models.AttachVPCPayloadRequest
		json.Unmarshal(body, &req)

		if req.Action == "" {
			t.Error("Expected Action in request body")
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	vpcPayload := models.AttachVPCPayloadRequest{
		Action: "detach",
		VPCs: []models.VPC{
			{
				VpcName:    "test-vpc",
				Ipv4_cidr:  "10.0.0.0/24",
				Network_id: 100,
			},
		},
	}

	err := client.DetachVPCPostgressDB(vpcPayload, "789", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestDetachVPCPostgressDB_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	vpcPayload := models.AttachVPCPayloadRequest{
		Action: "detach",
		VPCs:   []models.VPC{},
	}

	err := client.DetachVPCPostgressDB(vpcPayload, "789", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestUpgradePostgressPlan(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/789/rds-upgrade/" {
			t.Errorf("Expected path /rds/cluster/789/rds-upgrade/, got %s", r.URL.Path)
		}

		// Verify request body
		body, _ := ioutil.ReadAll(r.Body)
		var req map[string]interface{}
		json.Unmarshal(body, &req)

		templateID := int(req["template_id"].(float64))
		if templateID != 204 {
			t.Errorf("Expected template_id 204, got %d", templateID)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.UpgradePostgressPlan("789", 204, "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestUpgradePostgressPlan_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusPaymentRequired)
		w.Write([]byte("Insufficient credits"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.UpgradePostgressPlan("789", 204, "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestUpdateParameterGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/789/parameter-group/555/add" {
			t.Errorf("Expected path /rds/cluster/789/parameter-group/555/add, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.UpdateParameterGroup("789", "555", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestUpdateParameterGroup_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Parameter group not found"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.UpdateParameterGroup("789", "999", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestUpgradeDiskStorage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/789/disk-upgrade/" {
			t.Errorf("Expected path /rds/cluster/789/disk-upgrade/, got %s", r.URL.Path)
		}

		// Verify request body
		body, _ := ioutil.ReadAll(r.Body)
		var req map[string]interface{}
		json.Unmarshal(body, &req)

		size := int(req["size"].(float64))
		if size != 150 {
			t.Errorf("Expected size 150, got %d", size)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.UpgradeDiskStorage("789", 150, "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestUpgradeDiskStorage_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("Cannot upgrade disk"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.UpgradeDiskStorage("789", 150, "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}
