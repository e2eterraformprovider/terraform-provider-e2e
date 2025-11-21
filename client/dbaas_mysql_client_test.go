package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func TestNewMySqlDb(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code":    float64(200),
		"message": "MySQL database created successfully",
		"data": map[string]interface{}{
			"id":   float64(456),
			"name": "test-mysql",
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
		var req models.MySqlCreate
		json.Unmarshal(body, &req)

		if req.Name == "" {
			t.Error("Expected Name in request body")
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	mysqlCreate := &models.MySqlCreate{
		Name:             "test-mysql",
		SoftwareID:       2,
		TemplateID:       101,
		PublicIPRequired: true,
		Group:            "default",
		Database: models.DBConfig{
			User:     "mysqluser",
			Password: "mysqlpass",
			Name:     "mysqldb",
		},
	}

	result, err := client.NewMySqlDb(mysqlCreate, "test-project", "test-location")

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

func TestNewMySqlDb_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	mysqlCreate := &models.MySqlCreate{
		Name:       "test-mysql",
		SoftwareID: 2,
		TemplateID: 101,
	}

	result, err := client.NewMySqlDb(mysqlCreate, "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestGetMySqlDbaas(t *testing.T) {
	mockResponse := models.DBResponse{
		Code:    200,
		Message: "success",
		Data: models.DB{
			ID:     456,
			Name:   "test-mysql",
			Status: "active",
			Software: models.Software{
				Name:    "MySQL",
				Version: "8.0",
				Engine:  "mysql",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/456/" {
			t.Errorf("Expected path /rds/cluster/456/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.GetMySqlDbaas("456", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Data.ID != 456 {
		t.Errorf("Expected ID 456, got %d", result.Data.ID)
	}

	if result.Data.Name != "test-mysql" {
		t.Errorf("Expected Name test-mysql, got %s", result.Data.Name)
	}
}

func TestGetMySqlDbaas_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.GetMySqlDbaas("999", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestDeleteMySqlDBaaS(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code":    float64(200),
		"message": "MySQL database deleted successfully",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/456/" {
			t.Errorf("Expected path /rds/cluster/456/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.DeleteMySqlDBaaS("456", "test-project", "test-location")

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

func TestDeleteMySqlDBaaS_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server Error"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.DeleteMySqlDBaaS("456", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestResumeMySqlDBaaS(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code":    float64(200),
		"message": "MySQL database resumed successfully",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/456/resume" {
			t.Errorf("Expected path /rds/cluster/456/resume, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.ResumeMySqlDBaaS("456", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestResumeMySqlDBaaS_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Forbidden"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.ResumeMySqlDBaaS("456", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestStopMySqlDBaaS(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code":    float64(200),
		"message": "MySQL database stopped successfully",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/456/shutdown" {
			t.Errorf("Expected path /rds/cluster/456/shutdown, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.StopMySqlDBaaS("456", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestStopMySqlDBaaS_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.StopMySqlDBaaS("456", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestRestartMySqlDBaaS(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code":    float64(200),
		"message": "MySQL database restarted successfully",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/456/restart" {
			t.Errorf("Expected path /rds/cluster/456/restart, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.RestartMySqlDBaaS("456", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestRestartMySqlDBaaS_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Service Unavailable"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.RestartMySqlDBaaS("456", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestAttachVpcToMySql(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code":    float64(200),
		"message": "VPC attached successfully",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/456/vpc-attach/" {
			t.Errorf("Expected path /rds/cluster/456/vpc-attach/, got %s", r.URL.Path)
		}

		// Verify request body
		body, _ := ioutil.ReadAll(r.Body)
		var req models.AttachVPCPayloadRequest
		json.Unmarshal(body, &req)

		if req.Action == "" {
			t.Error("Expected Action in request body")
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	vpcPayload := &models.AttachVPCPayloadRequest{
		Action: "attach",
		VPCs: []models.VPC{
			{
				VpcName:    "test-vpc",
				Ipv4_cidr:  "10.0.0.0/24",
				Network_id: 100,
			},
		},
	}

	result, err := client.AttachVpcToMySql(vpcPayload, "456", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestAttachVpcToMySql_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("VPC not found"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	vpcPayload := &models.AttachVPCPayloadRequest{
		Action: "attach",
		VPCs:   []models.VPC{},
	}

	result, err := client.AttachVpcToMySql(vpcPayload, "456", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestDetachVpcFromMySql(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code":    float64(200),
		"message": "VPC detached successfully",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/456/vpc-detach/" {
			t.Errorf("Expected path /rds/cluster/456/vpc-detach/, got %s", r.URL.Path)
		}

		// Verify request body
		body, _ := ioutil.ReadAll(r.Body)
		var req models.AttachVPCPayloadRequest
		json.Unmarshal(body, &req)

		if req.Action == "" {
			t.Error("Expected Action in request body")
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	vpcPayload := &models.AttachVPCPayloadRequest{
		Action: "detach",
		VPCs: []models.VPC{
			{
				VpcName:    "test-vpc",
				Ipv4_cidr:  "10.0.0.0/24",
				Network_id: 100,
			},
		},
	}

	result, err := client.DetachVpcFromMySql(vpcPayload, "456", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestDetachVpcFromMySql_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	vpcPayload := &models.AttachVPCPayloadRequest{
		Action: "detach",
		VPCs:   []models.VPC{},
	}

	result, err := client.DetachVpcFromMySql(vpcPayload, "456", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestAttachPGToMySqlDBaaS(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code":    float64(200),
		"message": "Parameter group attached successfully",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/456/parameter-group/789/add" {
			t.Errorf("Expected path /rds/cluster/456/parameter-group/789/add, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.AttachPGToMySqlDBaaS("456", "789", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestAttachPGToMySqlDBaaS_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Parameter group not found"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.AttachPGToMySqlDBaaS("456", "999", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestDetachPGFromMySqlDBaaS(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code":    float64(200),
		"message": "Parameter group detached successfully",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/456/parameter-group/789/detach" {
			t.Errorf("Expected path /rds/cluster/456/parameter-group/789/detach, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.DetachPGFromMySqlDBaaS("456", "789", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestDetachPGFromMySqlDBaaS_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server Error"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.DetachPGFromMySqlDBaaS("456", "789", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestAttachPublicIPToMySql(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code":    float64(200),
		"message": "Public IP attached successfully",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/456/public-ip-attach/" {
			t.Errorf("Expected path /rds/cluster/456/public-ip-attach/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.AttachPublicIPToMySql("456", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestAttachPublicIPToMySql_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Conflict"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.AttachPublicIPToMySql("456", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestDetachPublicIPFromMySql(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code":    float64(200),
		"message": "Public IP detached successfully",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/456/public-ip-detach/" {
			t.Errorf("Expected path /rds/cluster/456/public-ip-detach/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.DetachPublicIPFromMySql("456", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestDetachPublicIPFromMySql_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.DetachPublicIPFromMySql("456", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestUpgradeMySQLPlan(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/456/rds-upgrade/" {
			t.Errorf("Expected path /rds/cluster/456/rds-upgrade/, got %s", r.URL.Path)
		}

		// Verify request body
		body, _ := ioutil.ReadAll(r.Body)
		var req models.MySQlPlanUpgradeAction
		json.Unmarshal(body, &req)

		if req.TemplateID != 202 {
			t.Errorf("Expected template_id 202, got %d", req.TemplateID)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.UpgradeMySQLPlan("456", 202, "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestUpgradeMySQLPlan_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusPaymentRequired)
		w.Write([]byte("Insufficient credits"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.UpgradeMySQLPlan("456", 202, "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestExpandMySQLDBaaSDisk(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/cluster/456/disk-upgrade/" {
			t.Errorf("Expected path /rds/cluster/456/disk-upgrade/, got %s", r.URL.Path)
		}

		// Verify request body
		body, _ := ioutil.ReadAll(r.Body)
		var req models.MYSQLExpandDisk
		json.Unmarshal(body, &req)

		if req.Size != 100 {
			t.Errorf("Expected size 100, got %d", req.Size)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.ExpandMySQLDBaaSDisk("456", 100, "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestExpandMySQLDBaaSDisk_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("Cannot expand disk"))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.ExpandMySQLDBaaSDisk("456", 100, "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}
