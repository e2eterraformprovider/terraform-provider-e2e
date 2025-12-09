package goe2e

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ============================================================================
// Happy Path Tests
// TestReserveIP_ListReserveIPs_Success tests successful list of reserve IPs
func TestReserveIP_ListReserveIPs_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/reserve_ip/" {
			t.Errorf("Expected path /reserve_ip/, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
			"code": 200,
			"message": "success",
			"data": [
				{
					"reserve_id": "1",
					"appliance_type": "node",
					"ip_address": "192.168.1.1",
					"reserved_type": "FloatingIP",
					"status": "Available",
					"vm_id": 0,
					"vm_name": "",
					"bought_at": "2025-01-01T00:00:00Z",
					"urn": "e2e:reserve_ip:us-east-1:192.168.1.1"
				}
			]
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	ips, resp, err := client.ReserveIP.ListReserveIPs(context.Background())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	if len(ips) != 1 {
		t.Errorf("Expected 1 IP, got %d", len(ips))
	}
	if ips[0].IPAddress != "192.168.1.1" {
		t.Errorf("Expected IP 192.168.1.1, got %s", ips[0].IPAddress)
	}
}

// TestReserveIP_ListReserveIPs_Empty tests empty list response
func TestReserveIP_ListReserveIPs_Empty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
			"code": 200,
			"message": "success",
			"data": []
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	ips, _, err := client.ReserveIP.ListReserveIPs(context.Background())
	assertNoError(t, err)
	if len(ips) != 0 {
		t.Errorf("Expected 0 IPs, got %d", len(ips))
	}
}

// TestReserveIP_GetReserveIP_Success tests successful get of single reserve IP
func TestReserveIP_GetReserveIP_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/reserve_ip/192.168.1.1/" {
			t.Errorf("Expected path /reserve_ip/192.168.1.1/, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
			"code": 200,
			"message": "success",
			"data": {
				"reserve_id": "1",
				"appliance_type": "node",
				"ip_address": "192.168.1.1",
				"reserved_type": "PublicIP",
				"status": "Attached",
				"vm_id": 123,
				"vm_name": "test-node",
				"bought_at": "2025-01-01T00:00:00Z",
				"urn": "e2e:reserve_ip:us-east-1:192.168.1.1"
			}
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	ip, resp, err := client.ReserveIP.GetReserveIP(context.Background(), "192.168.1.1")
	assertNoError(t, err)
	if ip == nil {
		t.Fatal("Expected IP, got nil")
	}
	if ip.IPAddress != "192.168.1.1" {
		t.Errorf("Expected IP 192.168.1.1, got %s", ip.IPAddress)
	}
	if ip.Status != "Attached" {
		t.Errorf("Expected status Attached, got %s", ip.Status)
	}
	assertNotNil(t, resp, "Expected non-nil response")
}

// TestReserveIP_DeleteReserveIP_Success tests successful delete of reserve IP
func TestReserveIP_DeleteReserveIP_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
			"code": 200,
			"message": "success"
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	resp, err := client.ReserveIP.DeleteReserveIP(context.Background(), "192.168.1.1")
	assertNoError(t, err)
	assertNotNil(t, resp, "Expected non-nil response")
}

// Input Validation Tests
// TestReserveIP_GetReserveIP_EmptyIPAddress tests validation of empty IP
func TestReserveIP_GetReserveIP_EmptyIPAddress(t *testing.T) {
	client, _ := NewClient("key", "token", "proj", "region")
	_, _, err := client.ReserveIP.GetReserveIP(context.Background(), "")
	if err == nil {
		t.Error("Expected error for empty IP address, got nil")
	}
	// Verify it's an ArgError
	if _, ok := err.(*ArgError); !ok {
		t.Errorf("Expected ArgError, got %T: %v", err, err)
	}
}

// TestReserveIP_DeleteReserveIP_EmptyIPAddress tests validation of empty IP for delete
func TestReserveIP_DeleteReserveIP_EmptyIPAddress(t *testing.T) {
	client, _ := NewClient("key", "token", "proj", "region")
	_, err := client.ReserveIP.DeleteReserveIP(context.Background(), "")
	if err == nil {
		t.Error("Expected error for empty IP address, got nil")
	}
}

// Error Handling Tests
// TestReserveIP_Get_NotFound tests 404 response
func TestReserveIP_Get_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{
			"code": 404,
			"message": "not found"
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	ip, resp, err := client.ReserveIP.GetReserveIP(context.Background(), "192.168.1.1")
	if err == nil {
		t.Errorf("Expected error for 404, got nil")
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}
	if ip != nil {
		t.Error("Expected nil IP for 404, got non-nil")
	}
}

// TestReserveIP_List_BadRequest tests 400 response
func TestReserveIP_List_BadRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{
			"code": 400,
			"message": "bad request"
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	_, _, err := client.ReserveIP.ListReserveIPs(context.Background())
	if err == nil {
		t.Error("Expected error for 400 status, got nil")
	}
}

// TestReserveIP_Get_Unauthorized tests 401 response
func TestReserveIP_Get_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{
			"code": 401,
			"message": "unauthorized"
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	_, _, err := client.ReserveIP.GetReserveIP(context.Background(), "192.168.1.1")
	if err == nil {
		t.Error("Expected error for 401 status, got nil")
	}
}

// TestReserveIP_Get_Forbidden tests 403 response
func TestReserveIP_Get_Forbidden(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, `{
			"code": 403,
			"message": "forbidden"
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	_, _, err := client.ReserveIP.GetReserveIP(context.Background(), "192.168.1.1")
	if err == nil {
		t.Error("Expected error for 403 status, got nil")
	}
}

// TestReserveIP_Delete_ServerError tests 500 response
func TestReserveIP_Delete_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{
			"code": 500,
			"message": "server error"
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	_, err := client.ReserveIP.DeleteReserveIP(context.Background(), "192.168.1.1")
	if err == nil {
		t.Error("Expected error for 500 status, got nil")
	}
}

// TestReserveIP_Delete_BadGateway tests 502 response
func TestReserveIP_Delete_BadGateway(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, `{
			"code": 502,
			"message": "bad gateway"
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	_, err := client.ReserveIP.DeleteReserveIP(context.Background(), "192.168.1.1")
	if err == nil {
		t.Error("Expected error for 502 status, got nil")
	}
}

// TestReserveIP_Delete_ServiceUnavailable tests 503 response
func TestReserveIP_Delete_ServiceUnavailable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, `{
			"code": 503,
			"message": "service unavailable"
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	_, err := client.ReserveIP.DeleteReserveIP(context.Background(), "192.168.1.1")
	if err == nil {
		t.Error("Expected error for 503 status, got nil")
	}
}

// TestReserveIP_Delete_GatewayTimeout tests 504 response
func TestReserveIP_Delete_GatewayTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusGatewayTimeout)
		fmt.Fprintf(w, `{
			"code": 504,
			"message": "gateway timeout"
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	_, err := client.ReserveIP.DeleteReserveIP(context.Background(), "192.168.1.1")
	if err == nil {
		t.Error("Expected error for 504 status, got nil")
	}
}

// Response Parsing Tests
// TestReserveIP_Parse_MalformedJSON tests handling of malformed JSON
func TestReserveIP_Parse_MalformedJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{invalid json}`)
	}))
	defer server.Close()
	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	_, _, err := client.ReserveIP.GetReserveIP(context.Background(), "192.168.1.1")
	if err == nil {
		t.Error("Expected error for malformed JSON, got nil")
	}
}

// TestReserveIP_Parse_MissingFields tests handling of missing optional fields
func TestReserveIP_Parse_MissingFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
			"code": 200,
			"message": "success",
			"data": {
				"reserve_id": "1",
				"appliance_type": "node",
				"ip_address": "192.168.1.1",
				"reserved_type": "FloatingIP",
				"status": "Available",
				"bought_at": "2025-01-01T00:00:00Z"
			}
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	ip, _, err := client.ReserveIP.GetReserveIP(context.Background(), "192.168.1.1")
	if err != nil {
		t.Errorf("Expected no error for missing optional fields, got %v", err)
	}
	if ip == nil {
		t.Fatal("Expected IP, got nil")
	}
}

// TestReserveIP_Parse_WithAttachedNodes tests parsing of attached nodes array
func TestReserveIP_Parse_WithAttachedNodes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
			"code": 200,
			"message": "success",
			"data": {
				"reserve_id": "1",
				"ip_address": "192.168.1.1",
				"attached_nodes": [
					{
						"id": 123,
						"name": "test-node-1"
					},
					{
						"id": 124,
						"name": "test-node-2"
					}
				]
			}
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	ip, _, err := client.ReserveIP.GetReserveIP(context.Background(), "192.168.1.1")
	assertNoError(t, err)
	if len(ip.AttachedNodes) != 2 {
		t.Errorf("Expected 2 attached nodes, got %d", len(ip.AttachedNodes))
	}
	if ip.AttachedNodes[0].ID != 123 {
		t.Errorf("Expected node ID 123, got %d", ip.AttachedNodes[0].ID)
	}
	if ip.AttachedNodes[0].Name != "test-node-1" {
		t.Errorf("Expected node name test-node-1, got %s", ip.AttachedNodes[0].Name)
	}
}

// CreateReserveIP Tests
// TestReserveIP_CreateReserveIP_Success tests successful creation of reserve IP
func TestReserveIP_CreateReserveIP_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
			"code": 200,
			"message": "success",
			"data": {
				"reserve_id": "42",
				"appliance_type": "node",
				"ip_address": "203.0.113.1",
				"reserved_type": "FloatingIP",
				"status": "Available",
				"vm_id": 0,
				"vm_name": "",
				"bought_at": "2025-01-15T10:30:00Z",
				"urn": "e2e:reserve_ip:us-east-1:203.0.113.1",
				"project_name": "test-project"
			}
		}`)
	}))
	defer server.Close()
	client, err := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	ip, resp, err := client.ReserveIP.CreateReserveIP(context.Background())
	assertNoError(t, err)
	assertNotNil(t, resp, "Expected non-nil response")
	assertStatus(t, resp, http.StatusOK)
	assertNotNil(t, ip, "Expected non-nil ReserveIP")
	// Verify all response fields
	if ip.IPAddress != "203.0.113.1" {
		t.Errorf("Expected IP 203.0.113.1, got %s", ip.IPAddress)
	}
	if ip.ReserveID != "42" {
		t.Errorf("Expected ReserveID 42, got %s", ip.ReserveID)
	}
	if ip.ReservedType != "FloatingIP" {
		t.Errorf("Expected ReservedType FloatingIP, got %s", ip.ReservedType)
	}
	if ip.Status != "Available" {
		t.Errorf("Expected status Available, got %s", ip.Status)
	}
	if ip.ProjectName != "test-project" {
		t.Errorf("Expected ProjectName test-project, got %s", ip.ProjectName)
	}
	if ip.ApplianceType != "node" {
		t.Errorf("Expected ApplianceType node, got %s", ip.ApplianceType)
	}
	if ip.URN != "e2e:reserve_ip:us-east-1:203.0.113.1" {
		t.Errorf("Expected URN e2e:reserve_ip:us-east-1:203.0.113.1, got %s", ip.URN)
	}
}

// TestReserveIP_CreateReserveIP_Error tests error responses (400, 500) on create
func TestReserveIP_CreateReserveIP_Error(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		responseBody   string
		expectedErrMsg string
	}{
		{
			name:       "400 Bad Request",
			statusCode: http.StatusBadRequest,
			responseBody: `{
				"code": 400,
				"message": "bad request",
				"errors": ["Invalid request parameters"]
			}`,
			expectedErrMsg: "bad request",
		},
		{
			name:       "500 Internal Server Error",
			statusCode: http.StatusInternalServerError,
			responseBody: `{
				"code": 500,
				"message": "server error",
				"errors": ["Internal server error occurred"]
			}`,
			expectedErrMsg: "server error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertHTTPMethod(t, r, http.MethodPost)
				if r.URL.Path != "/reserve_ip/" {
					t.Errorf("Expected path /reserve_ip/, got %s", r.URL.Path)
				}
				w.WriteHeader(tt.statusCode)
				fmt.Fprint(w, tt.responseBody)
			}))
			defer server.Close()
			client, err := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}
			ip, resp, err := client.ReserveIP.CreateReserveIP(context.Background())
			// Verify error is returned
			assertError(t, err, tt.expectedErrMsg)
			assertErrorType(t, err, (*ErrorResponse)(nil))
			// Verify response is not nil and has correct status code
			assertNotNil(t, resp, "Expected non-nil response on error")
			assertStatus(t, resp, tt.statusCode)
			// Verify resource is nil on error
			assertNil(t, ip, "Expected nil ReserveIP on error")
			// Verify error response details
			errResp, ok := err.(*ErrorResponse)
			if !ok {
				t.Fatalf("Expected *ErrorResponse, got %T", err)
			}
			if errResp.Code != tt.statusCode {
				t.Errorf("Expected error code %d, got %d", tt.statusCode, errResp.Code)
			}
			if errResp.Message != tt.expectedErrMsg {
				t.Errorf("Expected error message '%s', got '%s'", tt.expectedErrMsg, errResp.Message)
			}
		})
	}
}

// AttachFloatingIP Tests
// TestReserveIP_AttachFloatingIP_Success tests successful attachment
func TestReserveIP_AttachFloatingIP_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/reserve_ip/203.0.113.1/attach/" {
			t.Errorf("Expected path /reserve_ip/203.0.113.1/attach/, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
			"code": 200,
			"message": "success"
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	req := &FloatingIPAttachmentRequest{
		IPAddress: "203.0.113.1",
		NodeIDs:   []string{"node-1", "node-2"},
	}
	resp, err := client.ReserveIP.AttachFloatingIP(context.Background(), req)
	assertNoError(t, err)
	assertNotNil(t, resp, "Expected non-nil response")
}

// TestReserveIP_AttachFloatingIP_NilRequest tests validation of nil request
func TestReserveIP_AttachFloatingIP_NilRequest(t *testing.T) {
	client, _ := NewClient("key", "token", "proj", "region")
	_, err := client.ReserveIP.AttachFloatingIP(context.Background(), nil)
	if err == nil {
		t.Error("Expected error for nil request, got nil")
	}
	if _, ok := err.(*ArgError); !ok {
		t.Errorf("Expected ArgError, got %T", err)
	}
}

// TestReserveIP_AttachFloatingIP_EmptyIPAddress tests validation of empty IP address
func TestReserveIP_AttachFloatingIP_EmptyIPAddress(t *testing.T) {
	client, _ := NewClient("key", "token", "proj", "region")
	req := &FloatingIPAttachmentRequest{
		IPAddress: "",
		NodeIDs:   []string{"node-1"},
	}
	_, err := client.ReserveIP.AttachFloatingIP(context.Background(), req)
	if err == nil {
		t.Error("Expected error for empty IP address, got nil")
	}
}

// TestReserveIP_AttachFloatingIP_EmptyNodeIDs tests validation of empty node IDs
func TestReserveIP_AttachFloatingIP_EmptyNodeIDs(t *testing.T) {
	client, _ := NewClient("key", "token", "proj", "region")
	req := &FloatingIPAttachmentRequest{
		IPAddress: "203.0.113.1",
		NodeIDs:   []string{},
	}
	_, err := client.ReserveIP.AttachFloatingIP(context.Background(), req)
	if err == nil {
		t.Error("Expected error for empty node IDs, got nil")
	}
}

// TestReserveIP_AttachFloatingIP_NotFound tests 404 response
func TestReserveIP_AttachFloatingIP_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{
			"code": 404,
			"message": "not found"
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	req := &FloatingIPAttachmentRequest{
		IPAddress: "203.0.113.1",
		NodeIDs:   []string{"node-1"},
	}
	_, err := client.ReserveIP.AttachFloatingIP(context.Background(), req)
	if err == nil {
		t.Error("Expected error for 404 status, got nil")
	}
}

// TestReserveIP_AttachFloatingIP_ServerError tests 500 response
func TestReserveIP_AttachFloatingIP_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{
			"code": 500,
			"message": "server error"
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	req := &FloatingIPAttachmentRequest{
		IPAddress: "203.0.113.1",
		NodeIDs:   []string{"node-1"},
	}
	_, err := client.ReserveIP.AttachFloatingIP(context.Background(), req)
	if err == nil {
		t.Error("Expected error for 500 status, got nil")
	}
}

// DetachFloatingIP Tests
// TestReserveIP_DetachFloatingIP_Success tests successful detachment
func TestReserveIP_DetachFloatingIP_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/reserve_ip/203.0.113.1/detach/" {
			t.Errorf("Expected path /reserve_ip/203.0.113.1/detach/, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
			"code": 200,
			"message": "success"
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	req := &FloatingIPDetachmentRequest{
		IPAddress: "203.0.113.1",
		NodeIDs:   []string{"node-1"},
	}
	resp, err := client.ReserveIP.DetachFloatingIP(context.Background(), req)
	assertNoError(t, err)
	assertNotNil(t, resp, "Expected non-nil response")
}

// TestReserveIP_DetachFloatingIP_NilRequest tests validation of nil request
func TestReserveIP_DetachFloatingIP_NilRequest(t *testing.T) {
	client, _ := NewClient("key", "token", "proj", "region")
	_, err := client.ReserveIP.DetachFloatingIP(context.Background(), nil)
	if err == nil {
		t.Error("Expected error for nil request, got nil")
	}
}

// TestReserveIP_DetachFloatingIP_EmptyIPAddress tests validation of empty IP address
func TestReserveIP_DetachFloatingIP_EmptyIPAddress(t *testing.T) {
	client, _ := NewClient("key", "token", "proj", "region")
	req := &FloatingIPDetachmentRequest{
		IPAddress: "",
		NodeIDs:   []string{"node-1"},
	}
	_, err := client.ReserveIP.DetachFloatingIP(context.Background(), req)
	if err == nil {
		t.Error("Expected error for empty IP address, got nil")
	}
}

// TestReserveIP_DetachFloatingIP_EmptyNodeIDs tests validation of empty node IDs
func TestReserveIP_DetachFloatingIP_EmptyNodeIDs(t *testing.T) {
	client, _ := NewClient("key", "token", "proj", "region")
	req := &FloatingIPDetachmentRequest{
		IPAddress: "203.0.113.1",
		NodeIDs:   []string{},
	}
	_, err := client.ReserveIP.DetachFloatingIP(context.Background(), req)
	if err == nil {
		t.Error("Expected error for empty node IDs, got nil")
	}
}

// TestReserveIP_DetachFloatingIP_NotFound tests 404 response
func TestReserveIP_DetachFloatingIP_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{
			"code": 404,
			"message": "not found"
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	req := &FloatingIPDetachmentRequest{
		IPAddress: "203.0.113.1",
		NodeIDs:   []string{"node-1"},
	}
	_, err := client.ReserveIP.DetachFloatingIP(context.Background(), req)
	if err == nil {
		t.Error("Expected error for 404 status, got nil")
	}
}

// TestReserveIP_DetachFloatingIP_ServerError tests 500 response
func TestReserveIP_DetachFloatingIP_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{
			"code": 500,
			"message": "server error"
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	req := &FloatingIPDetachmentRequest{
		IPAddress: "203.0.113.1",
		NodeIDs:   []string{"node-1"},
	}
	_, err := client.ReserveIP.DetachFloatingIP(context.Background(), req)
	if err == nil {
		t.Error("Expected error for 500 status, got nil")
	}
}

// Additional Response Parsing Tests
// TestReserveIP_Parse_WithFloatingIPAttachedNodes tests parsing of FloatingIPAttachedNodes
func TestReserveIP_Parse_WithFloatingIPAttachedNodes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
			"code": 200,
			"message": "success",
			"data": {
				"reserve_id": "1",
				"ip_address": "203.0.113.1",
				"floating_ip_attached_nodes": [
					{
						"id": 123,
						"name": "test-node-1",
						"vm_id": 456,
						"ip_address_public": "198.51.100.1",
						"ip_address_private": "10.0.0.1",
						"status_name": "running",
						"security_group_status": "active"
					},
					{
						"id": 124,
						"name": "test-node-2",
						"vm_id": 457,
						"ip_address_public": "198.51.100.2",
						"ip_address_private": "10.0.0.2",
						"status_name": "running",
						"security_group_status": "active"
					}
				]
			}
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	ip, _, err := client.ReserveIP.GetReserveIP(context.Background(), "203.0.113.1")
	assertNoError(t, err)
	if len(ip.FloatingIPAttachedNodes) != 2 {
		t.Errorf("Expected 2 floating IP attached nodes, got %d", len(ip.FloatingIPAttachedNodes))
	}
	if ip.FloatingIPAttachedNodes[0].ID != 123 {
		t.Errorf("Expected node ID 123, got %d", ip.FloatingIPAttachedNodes[0].ID)
	}
	if ip.FloatingIPAttachedNodes[0].Name != "test-node-1" {
		t.Errorf("Expected node name test-node-1, got %s", ip.FloatingIPAttachedNodes[0].Name)
	}
	if ip.FloatingIPAttachedNodes[0].VMID != 456 {
		t.Errorf("Expected VMID 456, got %d", ip.FloatingIPAttachedNodes[0].VMID)
	}
	if ip.FloatingIPAttachedNodes[0].IPAddressPublic != "198.51.100.1" {
		t.Errorf("Expected public IP 198.51.100.1, got %s", ip.FloatingIPAttachedNodes[0].IPAddressPublic)
	}
	if ip.FloatingIPAttachedNodes[0].IPAddressPrivate != "10.0.0.1" {
		t.Errorf("Expected private IP 10.0.0.1, got %s", ip.FloatingIPAttachedNodes[0].IPAddressPrivate)
	}
	if ip.FloatingIPAttachedNodes[0].StatusName != "running" {
		t.Errorf("Expected status_name running, got %s", ip.FloatingIPAttachedNodes[0].StatusName)
	}
	if ip.FloatingIPAttachedNodes[0].SecurityGroupStatus != "active" {
		t.Errorf("Expected security_group_status active, got %s", ip.FloatingIPAttachedNodes[0].SecurityGroupStatus)
	}
}

// TestReserveIP_Parse_WithProjectName tests parsing of ProjectName field
func TestReserveIP_Parse_WithProjectName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
			"code": 200,
			"message": "success",
			"data": {
				"reserve_id": "1",
				"ip_address": "203.0.113.1",
				"project_name": "my-test-project"
			}
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	ip, _, err := client.ReserveIP.GetReserveIP(context.Background(), "203.0.113.1")
	assertNoError(t, err)
	if ip.ProjectName != "my-test-project" {
		t.Errorf("Expected ProjectName my-test-project, got %s", ip.ProjectName)
	}
}
