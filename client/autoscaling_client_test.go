package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func TestCreateScalerGroup(t *testing.T) {
	mockResponse := models.CreateScalerGroupResponse{
		Code:    201,
		Message: "Scaler Group created successfully",
		Data: models.ScalerGroupCreateDetails{
			ID:   "sg-123",
			Name: "test-scaler-group",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/scaler/scalegroups" {
			t.Errorf("Expected path /scaler/scalegroups, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	req := &models.CreateScalerGroupRequest{
		Name:    "test-scaler-group",
		Desired: "2",
	}

	result, err := client.CreateScalerGroup(req, "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.ID != mockResponse.Data.ID {
		t.Errorf("Expected ID %s, got %s", mockResponse.Data.ID, result.ID)
	}

	if result.Name != mockResponse.Data.Name {
		t.Errorf("Expected Name %s, got %s", mockResponse.Data.Name, result.Name)
	}
}

func TestGetScalerGroup(t *testing.T) {
	mockResponse := models.GetScalerGroupResponse{
		Code:    200,
		Message: "success",
		Data: models.ScalerGroupGetDetail{
			Name:    "test-scaler-group",
			Desired: 3,
			Running: 3,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/scaler/scalegroups/sg-123/" {
			t.Errorf("Expected path /scaler/scalegroups/sg-123/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.GetScalerGroup("sg-123", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Name != mockResponse.Data.Name {
		t.Errorf("Expected Name %s, got %s", mockResponse.Data.Name, result.Name)
	}

	if result.Desired != mockResponse.Data.Desired {
		t.Errorf("Expected Desired %d, got %d", mockResponse.Data.Desired, result.Desired)
	}
}

func TestDeleteScalerGroup(t *testing.T) {
	mockResponse := models.DeleteScalerGroupResponse{
		Code:    200,
		Message: "Scaler Group deleted successfully",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/scaler/scalegroups/sg-123/" {
			t.Errorf("Expected path /scaler/scalegroups/sg-123/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.DeleteScalerGroup("sg-123", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestGetSavedImageByName(t *testing.T) {
	mockResponse := models.ListSavedImagesResponse{
		Code:    200,
		Message: "success",
		Data: []models.SavedImage{
			{
				ImageID:    "img-123",
				Name:       "test-image",
				TemplateID: 101,
				Distro:     "ubuntu",
			},
			{
				ImageID:    "img-456",
				Name:       "another-image",
				TemplateID: 102,
				Distro:     "centos",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/images/saved-images/" {
			t.Errorf("Expected path /images/saved-images/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.GetSavedImageByName("test-image", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.ImageID != "img-123" {
		t.Errorf("Expected ImageID img-123, got %s", result.ImageID)
	}

	if result.Name != "test-image" {
		t.Errorf("Expected Name test-image, got %s", result.Name)
	}
}

func TestGetSavedImageByName_NotFound(t *testing.T) {
	mockResponse := models.ListSavedImagesResponse{
		Code:    200,
		Message: "success",
		Data:    []models.SavedImage{},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.GetSavedImageByName("nonexistent-image", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestGetDefaultSecurityGroupID(t *testing.T) {
	mockResponse := models.GetScalerSecurityGroupsResponse{
		Code:    200,
		Message: "success",
		Data: []models.ScalerSecurityGroup{
			{
				ID:        1,
				IsDefault: true,
			},
			{
				ID:        2,
				IsDefault: false,
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/security_group/" {
			t.Errorf("Expected path /security_group/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.GetDefaultSecurityGroupID("test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result != 1 {
		t.Errorf("Expected ID 1, got %d", result)
	}
}

func TestGetDefaultSecurityGroupID_NotFound(t *testing.T) {
	mockResponse := models.GetScalerSecurityGroupsResponse{
		Code:    200,
		Message: "success",
		Data: []models.ScalerSecurityGroup{
			{
				ID:        1,
				IsDefault: false,
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.GetDefaultSecurityGroupID("test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != 0 {
		t.Errorf("Expected 0, got %d", result)
	}
}

func TestGetPlanDetailsFromPlanName(t *testing.T) {
	mockResponse := struct {
		Code int `json:"code"`
		Data struct {
			Data []struct {
				Name  string `json:"name"`
				Plan  string `json:"plan"`
				Specs struct {
					ID string `json:"id"`
				} `json:"specs"`
			} `json:"data"`
		} `json:"data"`
	}{
		Code: 200,
	}
	mockResponse.Data.Data = []struct {
		Name  string `json:"name"`
		Plan  string `json:"plan"`
		Specs struct {
			ID string `json:"id"`
		} `json:"specs"`
	}{
		{
			Name: "small",
			Plan: "plan-small",
			Specs: struct {
				ID string `json:"id"`
			}{
				ID: "plan-id-123",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/images/upgradeimage/101/" {
			t.Errorf("Expected path /images/upgradeimage/101/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	planID, slugName, err := client.GetPlanDetailsFromPlanName(101, "small", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if planID != "plan-id-123" {
		t.Errorf("Expected planID plan-id-123, got %s", planID)
	}

	if slugName != "plan-small" {
		t.Errorf("Expected slugName plan-small, got %s", slugName)
	}
}

func TestUpdateScalerGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/scaler/scalegroups/update/sg-123/" {
			t.Errorf("Expected path /scaler/scalegroups/update/sg-123/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    200,
			"message": "updated",
		})
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	req := &models.UpdateScalerGroupRequest{
		Name: "updated-name",
	}

	err := client.UpdateScalerGroup("sg-123", req, "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestUpdateDesiredNodeCount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/scaler/scalegroups/123/" {
			t.Errorf("Expected path /scaler/scalegroups/123/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.UpdateDesiredNodeCount(123, 5, "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestUpdateScalerGroupStatus(t *testing.T) {
	tests := []struct {
		name           string
		status         string
		expectedPath   string
		expectError    bool
	}{
		{
			name:         "Stop scaler group",
			status:       "Stopped",
			expectedPath: "/scaler/scalegroups/123/stop/",
			expectError:  false,
		},
		{
			name:         "Start scaler group",
			status:       "Running",
			expectedPath: "/scaler/scalegroups/123/start/",
			expectError:  false,
		},
		{
			name:         "Invalid status",
			status:       "Invalid",
			expectedPath: "",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "PUT" {
					t.Errorf("Expected PUT request, got %s", r.Method)
				}

				if tt.expectedPath != "" && r.URL.Path != tt.expectedPath {
					t.Errorf("Expected path %s, got %s", tt.expectedPath, r.URL.Path)
				}

				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"code":    200,
					"message": "status updated",
				})
			}))
			defer server.Close()

			client := NewClient("test-key", "test-token", server.URL)

			err := client.UpdateScalerGroupStatus(123, tt.status, "test-project", "test-location")

			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}

func TestGetVpcDetailsByName(t *testing.T) {
	mockResponse := struct {
		Data []models.VPCDetail `json:"data"`
	}{
		Data: []models.VPCDetail{
			{
				Name:      "test-vpc",
				NetworkID: 100,
				IPv4CIDR:  "10.0.0.0/24",
			},
			{
				Name:      "another-vpc",
				NetworkID: 101,
				IPv4CIDR:  "10.0.1.0/24",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/vpc/list/" {
			t.Errorf("Expected path /vpc/list/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.GetVpcDetailsByName("test-project", "test-location", "test-vpc")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Name != "test-vpc" {
		t.Errorf("Expected Name test-vpc, got %s", result.Name)
	}

	if result.NetworkID != 100 {
		t.Errorf("Expected NetworkID 100, got %d", result.NetworkID)
	}
}

func TestAttachVPCToScalerGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/scaler/scalegroups/sg-123/vpc/action/" {
			t.Errorf("Expected path /scaler/scalegroups/sg-123/vpc/action/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	vpcs := []models.VPCDetail{
		{
			Name:      "test-vpc",
			NetworkID: 100,
			IPv4CIDR:  "10.0.0.0/24",
		},
	}

	err := client.AttachVPCToScalerGroup("sg-123", vpcs, "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestDetachVPCFromScalerGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/scaler/scalegroups/sg-123/vpc/action/" {
			t.Errorf("Expected path /scaler/scalegroups/sg-123/vpc/action/, got %s", r.URL.Path)
		}

		query := r.URL.Query()
		if query.Get("vpc_id") != "vpc-456" {
			t.Errorf("Expected vpc_id vpc-456, got %s", query.Get("vpc_id"))
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.DetachVPCFromScalerGroup("sg-123", "vpc-456", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestGetPublicIPStatus(t *testing.T) {
	mockResponse := models.PublicIPStatusResponse{
		Code:    200,
		Message: "success",
		Data: models.PublicIPStatusData{
			IsPublicIPRequired: true,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/scaler/scalegroups/sg-123/public_ip/action/" {
			t.Errorf("Expected path /scaler/scalegroups/sg-123/public_ip/action/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.GetPublicIPStatus("sg-123", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if !result.IsPublicIPRequired {
		t.Error("Expected IsPublicIPRequired to be true, got false")
	}
}

func TestAttachPublicIP(t *testing.T) {
	mockResponse := models.PublicIPActionResponse{
		Code:    200,
		Message: "Public IP attached",
		Data:    "192.168.1.1",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/scaler/scalegroups/sg-123/public_ip/action/" {
			t.Errorf("Expected path /scaler/scalegroups/sg-123/public_ip/action/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.AttachPublicIP("sg-123", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Data != "192.168.1.1" {
		t.Errorf("Expected Data 192.168.1.1, got %s", result.Data)
	}
}

func TestDetachPublicIP(t *testing.T) {
	mockResponse := models.PublicIPActionResponse{
		Code:    200,
		Message: "Public IP detached",
		Data:    "",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/scaler/scalegroups/sg-123/public_ip/action/" {
			t.Errorf("Expected path /scaler/scalegroups/sg-123/public_ip/action/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.DetachPublicIP("sg-123", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestGetAttachedVPCsForScalerGroup(t *testing.T) {
	mockResponse := struct {
		Data []models.VPCPartial `json:"data"`
	}{
		Data: []models.VPCPartial{
			{
				Name:      "vpc-1",
				NetworkID: 100,
			},
			{
				Name:      "vpc-2",
				NetworkID: 101,
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/scaler/scalegroups/sg-123/vpc/action/" {
			t.Errorf("Expected path /scaler/scalegroups/sg-123/vpc/action/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	result, err := client.GetAttachedVPCsForScalerGroup("sg-123", "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 VPCs, got %d", len(result))
	}

	if result[0].Name != "vpc-1" {
		t.Errorf("Expected Name vpc-1, got %s", result[0].Name)
	}
}

func TestDetachSecurityGroupFromScalergroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/scaler/scalegroups/security_groups/sg-123/" {
			t.Errorf("Expected path /scaler/scalegroups/security_groups/sg-123/, got %s", r.URL.Path)
		}

		query := r.URL.Query()
		if query.Get("security_group_id") != "456" {
			t.Errorf("Expected security_group_id 456, got %s", query.Get("security_group_id"))
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.DetachSecurityGroupFromScalergroup("sg-123", 456, "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestAddSecurityGroupToScalergroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/scaler/scalegroups/security_groups/sg-123/" {
			t.Errorf("Expected path /scaler/scalegroups/security_groups/sg-123/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	err := client.AddSecurityGroupToScalergroup("sg-123", 456, "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}
