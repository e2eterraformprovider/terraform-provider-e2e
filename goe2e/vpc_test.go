package goe2e

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestCreateVPC tests creating a new VPC
func TestCreateVPC(t *testing.T) {
	server := newSuccessServer(t, map[string]interface{}{
		"network_id": 100,
		"name":       "test-vpc",
		"state":      "Active",
		"created_at": "2025-12-04T10:00:00Z",
		"ipv4_cidr":  "10.0.0.0/24",
		"gateway_ip": "10.0.0.1",
		"pool_size":  254,
		"is_active":  true,
	})
	defer server.Close()
	client, _ := NewClient("test-key", "test-token", "proj-123", "Chennai", SetBaseURL(server.URL+"/"))
	ctx := context.Background()
	createReq := &VpcCreateRequest{
		VpcName:  "test-vpc",
		IsE2EVpc: true,
		IPv4:     "",
	}
	vpc, resp, err := client.Vpcs.CreateVPC(ctx, createReq)
	assertNoError(t, err)
	assertNotNil(t, vpc, "Expected VPC to be returned")
	if vpc.Name != "test-vpc" {
		t.Errorf("Expected VPC name 'test-vpc', got %s", vpc.Name)
	}
	assertStatus(t, resp, http.StatusOK)
}

// TestCreateVPCMissingName tests that CreateVPC fails without a name
func TestCreateVPCMissingName(t *testing.T) {
	client, _ := NewClient("test-key", "test-token", "proj-123", "Chennai")
	ctx := context.Background()
	createReq := &VpcCreateRequest{
		VpcName:  "",
		IsE2EVpc: true,
		IPv4:     "",
	}
	_, _, err := client.Vpcs.CreateVPC(ctx, createReq)
	assertError(t, err, "")
}

// TestCreateVPCNilRequest tests that CreateVPC fails with nil request
func TestCreateVPCNilRequest(t *testing.T) {
	client, _ := NewClient("test-key", "test-token", "proj-123", "Chennai")
	ctx := context.Background()
	_, _, err := client.Vpcs.CreateVPC(ctx, nil)
	assertError(t, err, "")
}

// TestGetVPC tests retrieving a VPC
func TestGetVPC(t *testing.T) {
	server := newSuccessServer(t, map[string]interface{}{
		"network_id": 100,
		"name":       "test-vpc",
		"state":      "Active",
		"created_at": "2025-12-04T10:00:00Z",
		"ipv4_cidr":  "10.0.0.0/24",
		"gateway_ip": "10.0.0.1",
		"pool_size":  254,
		"is_active":  true,
	})
	defer server.Close()
	client, _ := NewClient("test-key", "test-token", "proj-123", "Chennai", SetBaseURL(server.URL+"/"))
	ctx := context.Background()
	vpc, resp, err := client.Vpcs.GetVPC(ctx, "100")
	assertNoError(t, err)
	assertNotNil(t, vpc, "Expected VPC to be returned")
	assertStatus(t, resp, http.StatusOK)
}

// TestGetVPCMissingID tests that GetVPC fails with empty VPC ID
func TestGetVPCMissingID(t *testing.T) {
	client, _ := NewClient("test-key", "test-token", "proj-123", "Chennai")
	ctx := context.Background()
	_, _, err := client.Vpcs.GetVPC(ctx, "")
	assertError(t, err, "")
}

// TestListVPCs tests listing all VPCs
func TestListVPCs(t *testing.T) {
	server := newSuccessServer(t, []interface{}{
		map[string]interface{}{
			"network_id": 100,
			"name":       "vpc-1",
			"state":      "Active",
			"created_at": "2025-12-04T10:00:00Z",
			"ipv4_cidr":  "10.0.0.0/24",
			"gateway_ip": "10.0.0.1",
			"pool_size":  254,
			"is_active":  true,
		},
		map[string]interface{}{
			"network_id": 101,
			"name":       "vpc-2",
			"state":      "Active",
			"created_at": "2025-12-04T11:00:00Z",
			"ipv4_cidr":  "10.1.0.0/24",
			"gateway_ip": "10.1.0.1",
			"pool_size":  254,
			"is_active":  true,
		},
	})
	defer server.Close()
	client, _ := NewClient("test-key", "test-token", "proj-123", "Chennai", SetBaseURL(server.URL+"/"))
	ctx := context.Background()
	vpcs, resp, err := client.Vpcs.ListVPCs(ctx)
	assertNoError(t, err)
	assertStatus(t, resp, http.StatusOK)
	if len(vpcs) != 2 {
		t.Errorf("Expected 2 VPCs, got %d", len(vpcs))
	}
}

// TestListVPCsEmpty tests listing when no VPCs exist
func TestListVPCsEmpty(t *testing.T) {
	server := newSuccessServer(t, []interface{}{})
	defer server.Close()
	client, _ := NewClient("test-key", "test-token", "proj-123", "Chennai", SetBaseURL(server.URL+"/"))
	ctx := context.Background()
	vpcs, resp, err := client.Vpcs.ListVPCs(ctx)
	assertNoError(t, err)
	assertStatus(t, resp, http.StatusOK)
	if len(vpcs) != 0 {
		t.Errorf("Expected 0 VPCs, got %d", len(vpcs))
	}
}

// TestDeleteVPC tests deleting a VPC
func TestDeleteVPC(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(buildSuccessResponse(1, "OK", nil)))
	}))
	defer server.Close()
	client, _ := NewClient("test-key", "test-token", "proj-123", "Chennai", SetBaseURL(server.URL+"/"))
	ctx := context.Background()
	resp, err := client.Vpcs.DeleteVPC(ctx, "100")
	assertNoError(t, err)
	assertNotNil(t, resp, "Expected response to be returned")
}

// TestDeleteVPCMissingID tests that DeleteVPC fails with empty VPC ID
func TestDeleteVPCMissingID(t *testing.T) {
	client, _ := NewClient("test-key", "test-token", "proj-123", "Chennai")
	ctx := context.Background()
	_, err := client.Vpcs.DeleteVPC(ctx, "")
	assertError(t, err, "")
}

// TestVPCWithCustomIPv4 tests creating a VPC with custom IPv4 CIDR
func TestVPCWithCustomIPv4(t *testing.T) {
	server := newSuccessServer(t, map[string]interface{}{
		"network_id": 102,
		"name":       "custom-vpc",
		"state":      "Active",
		"created_at": "2025-12-04T10:00:00Z",
		"ipv4_cidr":  "192.168.0.0/24",
		"gateway_ip": "192.168.0.1",
		"pool_size":  254,
		"is_active":  true,
	})
	defer server.Close()
	client, _ := NewClient("test-key", "test-token", "proj-123", "Chennai", SetBaseURL(server.URL+"/"))
	ctx := context.Background()
	createReq := &VpcCreateRequest{
		VpcName:  "custom-vpc",
		IsE2EVpc: false,
		IPv4:     "192.168.0.0/24",
	}
	vpc, resp, err := client.Vpcs.CreateVPC(ctx, createReq)
	assertNoError(t, err)
	assertNotNil(t, vpc, "Expected VPC to be returned")
	if vpc.Name != "custom-vpc" {
		t.Errorf("Expected VPC name 'custom-vpc', got %s", vpc.Name)
	}
	assertStatus(t, resp, http.StatusOK)
}

// TestVPCQueryParameters tests that query parameters are correctly added
func TestVPCQueryParameters(t *testing.T) {
	server := newSuccessServer(t, []interface{}{})
	defer server.Close()
	client, _ := NewClient("test-key", "test-token", "proj-123", "Chennai", SetBaseURL(server.URL+"/"))
	ctx := context.Background()
	_, _, err := client.Vpcs.ListVPCs(ctx)
	assertNoError(t, err)
}

// Error response tests - CreateVPC
func TestCreateVPC_ErrorResponse(t *testing.T) {
	server := newErrorServer(t, http.StatusInternalServerError, "Internal server error")
	defer server.Close()
	client, _ := NewClient("test-key", "test-token", "proj-123", "Chennai", SetBaseURL(server.URL+"/"), noRetryOpt())
	ctx := context.Background()
	createReq := &VpcCreateRequest{
		VpcName:  "test-vpc",
		IsE2EVpc: true,
		IPv4:     "",
	}
	_, _, err := client.Vpcs.CreateVPC(ctx, createReq)
	assertError(t, err, "")
}

func TestCreateVPC_BadRequest(t *testing.T) {
	server := newErrorServer(t, http.StatusBadRequest, "Invalid CIDR block")
	defer server.Close()
	client, _ := NewClient("test-key", "test-token", "proj-123", "Chennai", SetBaseURL(server.URL+"/"), noRetryOpt())
	ctx := context.Background()
	createReq := &VpcCreateRequest{
		VpcName:  "test-vpc",
		IsE2EVpc: false,
		IPv4:     "invalid-cidr",
	}
	_, _, err := client.Vpcs.CreateVPC(ctx, createReq)
	assertError(t, err, "")
}

func TestCreateVPC_Conflict(t *testing.T) {
	server := newErrorServer(t, http.StatusConflict, "VPC with this name already exists")
	defer server.Close()
	client, _ := NewClient("test-key", "test-token", "proj-123", "Chennai", SetBaseURL(server.URL+"/"), noRetryOpt())
	ctx := context.Background()
	createReq := &VpcCreateRequest{
		VpcName:  "existing-vpc",
		IsE2EVpc: true,
		IPv4:     "",
	}
	_, _, err := client.Vpcs.CreateVPC(ctx, createReq)
	assertError(t, err, "")
}

// Error response tests - GetVPC
func TestGetVPC_ErrorResponse(t *testing.T) {
	server := newErrorServer(t, http.StatusInternalServerError, "Internal server error")
	defer server.Close()
	client, _ := NewClient("test-key", "test-token", "proj-123", "Chennai", SetBaseURL(server.URL+"/"), noRetryOpt())
	ctx := context.Background()
	_, _, err := client.Vpcs.GetVPC(ctx, "100")
	assertError(t, err, "")
}
