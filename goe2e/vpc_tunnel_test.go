package goe2e

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Happy path tests
func TestListVPCTunnels(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodGet)
		assertQueryParam(t, r, "apikey", "test-api-key")
		assertQueryParam(t, r, "project_id", "test-project")
		assertQueryParam(t, r, "location", "test-location")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(buildListResponse([]interface{}{
			map[string]interface{}{
				"id":                      "tunnel-123",
				"name":                    "test-tunnel",
				"vpc_local_network_id":    "vpc-123",
				"vpc_peer_network_id":     "vpc-456",
				"is_peer_vpc_external":    false,
				"status":                  "ACTIVE",
				"local_traffic_selector":  []string{"10.0.0.0/16"},
				"remote_traffic_selector": []string{"10.1.0.0/16"},
				"local_gateway_ip":        "192.168.1.1",
				"remote_gateway_ip":       "192.168.2.1",
				"created_at":              "2024-01-01T00:00:00Z",
			},
		})))
	}))
	defer server.Close()
	client, err := NewClient("test-api-key", "test-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	result, _, _ := client.VPCTunnels.ListVPCTunnels(context.Background(), "vpc-123")
	assertNotNil(t, result, "Expected result, got nil")
	if len(result) != 1 {
		t.Fatalf("Expected 1 tunnel, got %d", len(result))
	}
	if result[0].ID != "tunnel-123" {
		t.Errorf("Expected ID tunnel-123, got %s", result[0].ID)
	}
	if result[0].Name != "test-tunnel" {
		t.Errorf("Expected Name test-tunnel, got %s", result[0].Name)
	}
	if result[0].Status != "ACTIVE" {
		t.Errorf("Expected Status ACTIVE, got %s", result[0].Status)
	}
}

func TestListVPCTunnels_Empty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(buildListResponse([]interface{}{})))
	}))
	defer server.Close()
	client, err := NewClient("test-api-key", "test-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	result, _, err := client.VPCTunnels.ListVPCTunnels(context.Background(), "vpc-123")
	assertNoError(t, err)
	assertNotNil(t, result, "Expected empty slice, not nil")
	if len(result) != 0 {
		t.Errorf("Expected 0 tunnels, got %d", len(result))
	}
}

func TestCreateVPCTunnel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(buildSuccessResponse(201, "Tunnel created successfully", map[string]interface{}{
			"id":                      "tunnel-123",
			"name":                    "test-tunnel",
			"vpc_local_network_id":    "vpc-123",
			"vpc_peer_network_id":     "vpc-456",
			"is_peer_vpc_external":    false,
			"status":                  "CREATING",
			"local_traffic_selector":  []string{"10.0.0.0/16"},
			"remote_traffic_selector": []string{"10.1.0.0/16"},
			"local_gateway_ip":        "192.168.1.1",
			"remote_gateway_ip":       "192.168.2.1",
			"created_at":              "2024-01-01T00:00:00Z",
		})))
	}))
	defer server.Close()
	client, err := NewClient("test-api-key", "test-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	createReq := &CreateVPCTunnelRequest{
		Name:              "test-tunnel",
		VPCLocalNetworkID: "vpc-123",
		VPCPeerNetworkID:  "vpc-456",
		IsPeerVPCExternal: false,
	}
	result, _, err := client.VPCTunnels.CreateVPCTunnel(context.Background(), createReq)
	assertNoError(t, err)
	if result.ID != "tunnel-123" {
		t.Errorf("Expected ID tunnel-123, got %s", result.ID)
	}
	if result.Name != "test-tunnel" {
		t.Errorf("Expected Name test-tunnel, got %s", result.Name)
	}
	if result.Status != "CREATING" {
		t.Errorf("Expected Status CREATING, got %s", result.Status)
	}
}
func TestCreateVPCTunnel_External(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(buildSuccessResponse(201, "Tunnel created successfully", map[string]interface{}{
			"id":                   "tunnel-123",
			"name":                 "test-tunnel-external",
			"vpc_local_network_id": "vpc-123",
			"vpc_peer_network_id":  "external-vpc",
			"is_peer_vpc_external": true,
			"status":               "CREATING",
			"created_at":           "2024-01-01T00:00:00Z",
		})))
	}))
	defer server.Close()
	client, err := NewClient("test-api-key", "test-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	createReq := &CreateVPCTunnelRequest{
		Name:              "test-tunnel-external",
		VPCLocalNetworkID: "vpc-123",
		VPCPeerNetworkID:  "external-vpc",
		IsPeerVPCExternal: true,
	}
	result, _, err := client.VPCTunnels.CreateVPCTunnel(context.Background(), createReq)
	assertNoError(t, err)
	if !result.IsPeerVPCExternal {
		t.Error("Expected IsPeerVPCExternal to be true")
	}
}

func TestGetVPCTunnel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(buildSuccessResponse(200, "success", map[string]interface{}{
			"id":                      "tunnel-123",
			"name":                    "test-tunnel",
			"vpc_local_network_id":    "vpc-123",
			"vpc_peer_network_id":     "vpc-456",
			"is_peer_vpc_external":    false,
			"status":                  "ACTIVE",
			"local_traffic_selector":  []string{"10.0.0.0/16"},
			"remote_traffic_selector": []string{"10.1.0.0/16"},
			"local_gateway_ip":        "192.168.1.1",
			"remote_gateway_ip":       "192.168.2.1",
			"created_at":              "2024-01-01T00:00:00Z",
		})))
	}))
	defer server.Close()
	client, err := NewClient("test-api-key", "test-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	result, _, err := client.VPCTunnels.GetVPCTunnel(context.Background(), "tunnel-123")
	assertNoError(t, err)
	if result.Status != "ACTIVE" {
		t.Errorf("Expected Status ACTIVE, got %s", result.Status)
	}
}

func TestGetVPCTunnel_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{
			"code": 404,
			"message": "not found"
		}`)
	}))
	defer server.Close()
	client, err := NewClient("test-api-key", "test-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	result, _, err := client.VPCTunnels.GetVPCTunnel(context.Background(), "nonexistent")
	if result != nil {
		t.Errorf("Expected nil result for 404, got: %v", result)
	}
	assertError(t, err, "")
}

func TestPauseVPCTunnel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPut)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(buildSuccessResponse(200, "Tunnel paused successfully", map[string]interface{}{})))
	}))
	defer server.Close()
	client, err := NewClient("test-api-key", "test-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	_, err = client.VPCTunnels.PauseVPCTunnel(context.Background(), "tunnel-123")
	assertNoError(t, err)
}

func TestRestartVPCTunnel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPut)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(buildSuccessResponse(200, "Tunnel restarted successfully", map[string]interface{}{})))
	}))
	defer server.Close()
	client, err := NewClient("test-api-key", "test-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	_, err = client.VPCTunnels.RestartVPCTunnel(context.Background(), "tunnel-123")
	assertNoError(t, err)
}

func TestDeleteVPCTunnel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodDelete)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(buildSuccessResponse(200, "Tunnel deleted successfully", map[string]interface{}{})))
	}))
	defer server.Close()
	client, err := NewClient("test-api-key", "test-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	_, err = client.VPCTunnels.DeleteVPCTunnel(context.Background(), "tunnel-123")
	assertNoError(t, err)
}

// Validation tests
func TestListVPCTunnels_EmptyVPCNetworkID(t *testing.T) {
	client, err := NewClient("test-api-key", "test-token", "test-project", "test-location", WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	_, _, err = client.VPCTunnels.ListVPCTunnels(context.Background(), "")
	assertError(t, err, "")
}

func TestCreateVPCTunnel_NilRequest(t *testing.T) {
	client, err := NewClient("test-api-key", "test-token", "test-project", "test-location", WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	_, _, err = client.VPCTunnels.CreateVPCTunnel(context.Background(), nil)
	assertError(t, err, "")
}

func TestCreateVPCTunnel_EmptyName(t *testing.T) {
	client, err := NewClient("test-api-key", "test-token", "test-project", "test-location", WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	createReq := &CreateVPCTunnelRequest{
		Name:              "",
		VPCLocalNetworkID: "vpc-123",
		VPCPeerNetworkID:  "vpc-456",
		IsPeerVPCExternal: false,
	}
	_, _, err = client.VPCTunnels.CreateVPCTunnel(context.Background(), createReq)
	assertError(t, err, "")
}
