package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func TestGetSoftwareId(t *testing.T) {
	mockResponse := models.PlanResponse{
		Code:    200,
		Message: "success",
		Data: models.PlanData{
			DatabaseEngines: []models.EngineDefinition{
				{
					EngineID:      1,
					EngineName:    "mysql",
					EngineVersion: "8.0",
				},
				{
					EngineID:      2,
					EngineName:    "postgres",
					EngineVersion: "14.0",
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/plans/" {
			t.Errorf("Expected path /rds/plans/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.GetSoftwareId("test-project", "test-location", "mysql", "8.0")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result != 1 {
		t.Errorf("Expected EngineID 1, got %d", result)
	}
}

func TestGetSoftwareId_NotFound(t *testing.T) {
	mockResponse := models.PlanResponse{
		Code:    200,
		Message: "success",
		Data: models.PlanData{
			DatabaseEngines: []models.EngineDefinition{
				{
					EngineID:      1,
					EngineName:    "mysql",
					EngineVersion: "8.0",
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.GetSoftwareId("test-project", "test-location", "postgres", "15.0")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != -1 {
		t.Errorf("Expected -1, got %d", result)
	}
}

func TestGetTemplateId(t *testing.T) {
	mockResponse := models.PlanResponse{
		Code:    200,
		Message: "success",
		Data: models.PlanData{
			TemplatePlans: []models.PlanTemplate{
				{
					PlanTemplateID: 100,
					PlanName:       "small",
				},
				{
					PlanTemplateID: 101,
					PlanName:       "medium",
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/rds/plans/" {
			t.Errorf("Expected path /rds/plans/, got %s", r.URL.Path)
		}

		query := r.URL.Query()
		if query.Get("software_id") != "1" {
			t.Errorf("Expected software_id 1, got %s", query.Get("software_id"))
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.GetTemplateId("test-project", "test-location", "small", "1")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result != 100 {
		t.Errorf("Expected PlanTemplateID 100, got %d", result)
	}
}

func TestGetTemplateId_NotFound(t *testing.T) {
	mockResponse := models.PlanResponse{
		Code:    200,
		Message: "success",
		Data: models.PlanData{
			TemplatePlans: []models.PlanTemplate{
				{
					PlanTemplateID: 100,
					PlanName:       "small",
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.GetTemplateId("test-project", "test-location", "large", "1")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != -1 {
		t.Errorf("Expected -1, got %d", result)
	}
}

func TestExpandMariaDBVpcList(t *testing.T) {
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

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockVpcResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	vpcIDs := []string{"100"}
	result, err := client.ExpandMariaDBVpcList(vpcIDs, "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 VPC, got %d", len(result))
	}

	if result[0].VPCName != "test-vpc" {
		t.Errorf("Expected VPCName test-vpc, got %s", result[0].VPCName)
	}
}

func TestExpandMariaDBVpcList_InactiveVPC(t *testing.T) {
	mockVpcResponse := models.VpcResponse{
		Code:    200,
		Message: "success",
		Data: models.Vpc{
			Name:       "test-vpc",
			Network_id: 100,
			Ipv4_cidr:  "10.0.0.0/24",
			State:      "Inactive",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockVpcResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	vpcIDs := []string{"100"}
	result, err := client.ExpandMariaDBVpcList(vpcIDs, "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error for inactive VPC, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}
