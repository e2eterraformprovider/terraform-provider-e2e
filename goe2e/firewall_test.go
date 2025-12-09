package goe2e

import (
	"context"
	"net/http"
	"testing"
)

// ========================================================================
// Happy Path Tests - Successful Operations
// ========================================================================

func TestCreateFirewall(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/fortigate/create/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testURLPath(t, r, "/fortigate/create/")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		writeJSON(w, http.StatusCreated, `{
			"code": 201,
			"message": "Firewall created successfully",
			"data": {
				"id": "fw-123",
				"vm_id": 12345,
				"name": "test-firewall",
				"label": "Test Firewall",
				"plan": "fortigate-small",
				"status": "active",
				"public_ip_address": "203.0.113.1",
				"private_ip_address": "10.0.1.10",
				"memory": "4096",
				"disk": "50",
				"vcpus": "2",
				"created_at": "2024-01-01T00:00:00Z",
				"is_active": true,
				"is_locked": false,
				"is_fortigate_vm": true,
				"vpc_id": "vpc-789",
				"cn_id": "cn-456"
			}
		}`)
	})

	createReq := &FirewallCreateRequest{
		Name:  "test-firewall",
		Label: "Test Firewall",
		Plan:  "fortigate-small",
		Image: "fortinet-7.2",
		VPCID: "vpc-789",
		CNID:  "cn-456",
	}

	result, _, err := ts.client.Firewall.CreateFirewall(context.Background(), createReq)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.ID != "fw-123" {
		t.Errorf("Expected ID fw-123, got %s", result.ID)
	}

	if result.Name != "test-firewall" {
		t.Errorf("Expected Name test-firewall, got %s", result.Name)
	}

	if result.VMID != 12345 {
		t.Errorf("Expected VMID 12345, got %d", result.VMID)
	}

	if !result.IsFortigateVM {
		t.Error("Expected IsFortigateVM to be true")
	}

	if result.VPCID != "vpc-789" {
		t.Errorf("Expected VPCID vpc-789, got %s", result.VPCID)
	}
}

func TestGetFirewalls(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/fortigate/list/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/fortigate/list/")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": [
				{
					"id": "fw-123",
					"vm_id": 12345,
					"name": "firewall-1",
					"status": "active",
					"is_fortigate_vm": true,
					"vpc_id": "vpc-789"
				},
				{
					"id": "fw-456",
					"vm_id": 67890,
					"name": "firewall-2",
					"status": "active",
					"is_fortigate_vm": true,
					"vpc_id": "vpc-999"
				}
			]
		}`)
	})

	result, _, err := ts.client.Firewall.GetFirewalls(context.Background())

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 firewalls, got %d", len(result))
	}

	if result[0].ID != "fw-123" {
		t.Errorf("Expected first firewall ID fw-123, got %s", result[0].ID)
	}

	if result[1].ID != "fw-456" {
		t.Errorf("Expected second firewall ID fw-456, got %s", result[1].ID)
	}
}

func TestGetFirewall(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/fortigate/fw-123/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/fortigate/fw-123/")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": {
				"id": "fw-123",
				"vm_id": 12345,
				"name": "test-firewall",
				"label": "Test Firewall",
				"plan": "fortigate-small",
				"status": "active",
				"public_ip_address": "203.0.113.1",
				"private_ip_address": "10.0.1.10",
				"memory": "4096",
				"disk": "50",
				"vcpus": "2",
				"created_at": "2024-01-01T00:00:00Z",
				"is_active": true,
				"is_locked": false,
				"is_fortigate_vm": true,
				"vpc_id": "vpc-789"
			}
		}`)
	})

	result, _, err := ts.client.Firewall.GetFirewall(context.Background(), "fw-123")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.ID != "fw-123" {
		t.Errorf("Expected ID fw-123, got %s", result.ID)
	}

	if result.Name != "test-firewall" {
		t.Errorf("Expected Name test-firewall, got %s", result.Name)
	}
}

func TestGetFirewall_NotFound(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/fortigate/nonexistent/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	result, _, err := ts.client.Firewall.GetFirewall(context.Background(), "nonexistent")

	if err != nil {
		t.Fatalf("Expected no error for 404, got: %v", err)
	}

	if result != nil {
		t.Errorf("Expected nil result for 404, got: %v", result)
	}
}

func TestUpdateFirewall(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/fortigate/fw-123/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		testURLPath(t, r, "/fortigate/fw-123/")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Firewall updated successfully",
			"data": {
				"id": "fw-123",
				"vm_id": 12345,
				"name": "updated-firewall",
				"label": "Updated Firewall",
				"plan": "fortigate-small",
				"status": "active",
				"is_fortigate_vm": true,
				"vpc_id": "vpc-789"
			}
		}`)
	})

	updateReq := &FirewallUpdateRequest{
		Name:  String("updated-firewall"),
		Label: String("Updated Firewall"),
	}

	result, _, err := ts.client.Firewall.UpdateFirewall(context.Background(), "fw-123", updateReq)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Name != "updated-firewall" {
		t.Errorf("Expected Name updated-firewall, got %s", result.Name)
	}

	if result.Label != "Updated Firewall" {
		t.Errorf("Expected Label 'Updated Firewall', got %s", result.Label)
	}
}

func TestDeleteFirewall(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/fortigate/fw-123/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		testURLPath(t, r, "/fortigate/fw-123/")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		w.WriteHeader(http.StatusOK)
	})

	_, err := ts.client.Firewall.DeleteFirewall(context.Background(), "fw-123")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

// ========================================================================
// Input Validation Tests - Empty/Nil Parameters
// ========================================================================

func TestCreateFirewall_NilRequest(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, _, err := ts.client.Firewall.CreateFirewall(context.Background(), nil)
	if err == nil {
		t.Fatal("Expected error for nil create request, got nil")
	}
}

func TestCreateFirewall_EmptyName(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	createReq := &FirewallCreateRequest{
		Name:  "",
		Plan:  "fortigate-small",
		Image: "fortinet-7.2",
		VPCID: "vpc-789",
	}

	_, _, err := ts.client.Firewall.CreateFirewall(context.Background(), createReq)
	if err == nil {
		t.Fatal("Expected error for empty name, got nil")
	}
}

func TestCreateFirewall_EmptyPlan(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	createReq := &FirewallCreateRequest{
		Name:  "test-firewall",
		Plan:  "",
		Image: "fortinet-7.2",
		VPCID: "vpc-789",
	}

	_, _, err := ts.client.Firewall.CreateFirewall(context.Background(), createReq)
	if err == nil {
		t.Fatal("Expected error for empty plan, got nil")
	}
}

func TestCreateFirewall_EmptyImage(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	createReq := &FirewallCreateRequest{
		Name:  "test-firewall",
		Plan:  "fortigate-small",
		Image: "",
		VPCID: "vpc-789",
	}

	_, _, err := ts.client.Firewall.CreateFirewall(context.Background(), createReq)
	if err == nil {
		t.Fatal("Expected error for empty image, got nil")
	}
}

func TestCreateFirewall_EmptyVPCID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	createReq := &FirewallCreateRequest{
		Name:  "test-firewall",
		Plan:  "fortigate-small",
		Image: "fortinet-7.2",
		VPCID: "",
	}

	_, _, err := ts.client.Firewall.CreateFirewall(context.Background(), createReq)
	if err == nil {
		t.Fatal("Expected error for empty vpc_id, got nil")
	}
}

func TestGetFirewall_EmptyID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, _, err := ts.client.Firewall.GetFirewall(context.Background(), "")
	if err == nil {
		t.Fatal("Expected error for empty firewallID, got nil")
	}
}

func TestUpdateFirewall_EmptyID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	updateReq := &FirewallUpdateRequest{
		Name: String("updated-firewall"),
	}

	_, _, err := ts.client.Firewall.UpdateFirewall(context.Background(), "", updateReq)
	if err == nil {
		t.Fatal("Expected error for empty firewallID, got nil")
	}
}

func TestUpdateFirewall_NilRequest(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, _, err := ts.client.Firewall.UpdateFirewall(context.Background(), "fw-123", nil)
	if err == nil {
		t.Fatal("Expected error for nil update request, got nil")
	}
}

func TestDeleteFirewall_EmptyID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.Firewall.DeleteFirewall(context.Background(), "")
	if err == nil {
		t.Fatal("Expected error for empty firewallID, got nil")
	}
}

// ========================================================================
// Error Handling Tests - Various HTTP Status Codes
// ========================================================================

func TestCreateFirewall_BadRequest(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/fortigate/create/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusBadRequest, `{
			"code": 400,
			"message": "Invalid firewall configuration"
		}`)
	})

	createReq := &FirewallCreateRequest{
		Name:  "test-firewall",
		Plan:  "fortigate-small",
		Image: "fortinet-7.2",
		VPCID: "vpc-789",
	}

	_, _, err := ts.client.Firewall.CreateFirewall(context.Background(), createReq)
	if err == nil {
		t.Fatal("Expected error for bad request, got nil")
	}
}

func TestCreateFirewall_Unauthorized(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/fortigate/create/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusUnauthorized, `{
			"code": 401,
			"message": "Unauthorized"
		}`)
	})

	createReq := &FirewallCreateRequest{
		Name:  "test-firewall",
		Plan:  "fortigate-small",
		Image: "fortinet-7.2",
		VPCID: "vpc-789",
	}

	_, _, err := ts.client.Firewall.CreateFirewall(context.Background(), createReq)
	if err == nil {
		t.Fatal("Expected error for unauthorized, got nil")
	}
}

func TestCreateFirewall_Forbidden(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/fortigate/create/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusForbidden, `{
			"code": 403,
			"message": "Forbidden"
		}`)
	})

	createReq := &FirewallCreateRequest{
		Name:  "test-firewall",
		Plan:  "fortigate-small",
		Image: "fortinet-7.2",
		VPCID: "vpc-789",
	}

	_, _, err := ts.client.Firewall.CreateFirewall(context.Background(), createReq)
	if err == nil {
		t.Fatal("Expected error for forbidden, got nil")
	}
}

func TestCreateFirewall_Conflict(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/fortigate/create/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusConflict, `{
			"code": 409,
			"message": "Firewall with this name already exists"
		}`)
	})

	createReq := &FirewallCreateRequest{
		Name:  "test-firewall",
		Plan:  "fortigate-small",
		Image: "fortinet-7.2",
		VPCID: "vpc-789",
	}

	_, _, err := ts.client.Firewall.CreateFirewall(context.Background(), createReq)
	if err == nil {
		t.Fatal("Expected error for conflict, got nil")
	}
}

func TestCreateFirewall_InternalServerError(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/fortigate/create/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Internal server error"
		}`)
	})

	createReq := &FirewallCreateRequest{
		Name:  "test-firewall",
		Plan:  "fortigate-small",
		Image: "fortinet-7.2",
		VPCID: "vpc-789",
	}

	_, _, err := ts.client.Firewall.CreateFirewall(context.Background(), createReq)
	if err == nil {
		t.Fatal("Expected error for internal server error, got nil")
	}
}

func TestCreateFirewall_BadGateway(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/fortigate/create/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusBadGateway, `{
			"code": 502,
			"message": "Bad gateway"
		}`)
	})

	createReq := &FirewallCreateRequest{
		Name:  "test-firewall",
		Plan:  "fortigate-small",
		Image: "fortinet-7.2",
		VPCID: "vpc-789",
	}

	_, _, err := ts.client.Firewall.CreateFirewall(context.Background(), createReq)
	if err == nil {
		t.Fatal("Expected error for bad gateway, got nil")
	}
}

func TestCreateFirewall_ServiceUnavailable(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/fortigate/create/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusServiceUnavailable, `{
			"code": 503,
			"message": "Service unavailable"
		}`)
	})

	createReq := &FirewallCreateRequest{
		Name:  "test-firewall",
		Plan:  "fortigate-small",
		Image: "fortinet-7.2",
		VPCID: "vpc-789",
	}

	_, _, err := ts.client.Firewall.CreateFirewall(context.Background(), createReq)
	if err == nil {
		t.Fatal("Expected error for service unavailable, got nil")
	}
}

func TestCreateFirewall_GatewayTimeout(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/fortigate/create/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusGatewayTimeout, `{
			"code": 504,
			"message": "Gateway timeout"
		}`)
	})

	createReq := &FirewallCreateRequest{
		Name:  "test-firewall",
		Plan:  "fortigate-small",
		Image: "fortinet-7.2",
		VPCID: "vpc-789",
	}

	_, _, err := ts.client.Firewall.CreateFirewall(context.Background(), createReq)
	if err == nil {
		t.Fatal("Expected error for gateway timeout, got nil")
	}
}

func TestGetFirewalls_ErrorResponse(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/fortigate/list/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Internal server error"
		}`)
	})

	_, _, err := ts.client.Firewall.GetFirewalls(context.Background())
	if err == nil {
		t.Fatal("Expected error for 500 response, got nil")
	}
}

func TestGetFirewall_ErrorResponse(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/fortigate/fw-123/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Internal server error"
		}`)
	})

	_, _, err := ts.client.Firewall.GetFirewall(context.Background(), "fw-123")
	if err == nil {
		t.Fatal("Expected error for 500 response, got nil")
	}
}

func TestUpdateFirewall_NotFound(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/fortigate/nonexistent/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusNotFound, `{
			"code": 404,
			"message": "Firewall not found"
		}`)
	})

	updateReq := &FirewallUpdateRequest{
		Name: String("updated-firewall"),
	}

	_, _, err := ts.client.Firewall.UpdateFirewall(context.Background(), "nonexistent", updateReq)
	if err == nil {
		t.Fatal("Expected error for not found, got nil")
	}
}

func TestUpdateFirewall_Conflict(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/fortigate/fw-123/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusConflict, `{
			"code": 409,
			"message": "Firewall is currently being updated"
		}`)
	})

	updateReq := &FirewallUpdateRequest{
		Name: String("updated-firewall"),
	}

	_, _, err := ts.client.Firewall.UpdateFirewall(context.Background(), "fw-123", updateReq)
	if err == nil {
		t.Fatal("Expected error for conflict, got nil")
	}
}

func TestDeleteFirewall_NotFound(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/fortigate/nonexistent/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusNotFound, `{
			"code": 404,
			"message": "Firewall not found"
		}`)
	})

	_, err := ts.client.Firewall.DeleteFirewall(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("Expected error for not found, got nil")
	}
}

func TestDeleteFirewall_Conflict(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/fortigate/fw-123/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusConflict, `{
			"code": 409,
			"message": "Cannot delete firewall with active connections"
		}`)
	})

	_, err := ts.client.Firewall.DeleteFirewall(context.Background(), "fw-123")
	if err == nil {
		t.Fatal("Expected error for conflict, got nil")
	}
}

// ========================================================================
// Response Parsing Tests - Malformed JSON, Missing Fields, Null Values
// ========================================================================

func TestCreateFirewall_MalformedJSON(t *testing.T) {
	server := newMalformedJSONServer(t)
	defer server.Close()

	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())

	createReq := &FirewallCreateRequest{
		Name:  "test-firewall",
		Plan:  "fortigate-small",
		Image: "fortinet-7.2",
		VPCID: "vpc-789",
	}

	_, _, err := client.Firewall.CreateFirewall(context.Background(), createReq)

	if err == nil {
		t.Error("Expected error for malformed JSON")
	}
}

func TestGetFirewalls_MalformedJSON(t *testing.T) {
	server := newMalformedJSONServer(t)
	defer server.Close()

	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	_, _, err := client.Firewall.GetFirewalls(context.Background())

	if err == nil {
		t.Error("Expected error for malformed JSON")
	}
}

func TestGetFirewall_MalformedJSON(t *testing.T) {
	server := newMalformedJSONServer(t)
	defer server.Close()

	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	_, _, err := client.Firewall.GetFirewall(context.Background(), "fw-123")

	if err == nil {
		t.Error("Expected error for malformed JSON")
	}
}

func TestCreateFirewall_MissingRequiredFields(t *testing.T) {
	server := newMissingFieldServer(t, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data": map[string]interface{}{
			// Missing "id" field
			"name": "test-firewall",
		},
	})
	defer server.Close()

	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())

	createReq := &FirewallCreateRequest{
		Name:  "test-firewall",
		Plan:  "fortigate-small",
		Image: "fortinet-7.2",
		VPCID: "vpc-789",
	}

	resp, _, err := client.Firewall.CreateFirewall(context.Background(), createReq)

	// Should handle missing fields gracefully
	if resp == nil && err == nil {
		t.Error("Expected response or error handling")
	}
}

func TestGetFirewall_NullFieldValues(t *testing.T) {
	server := newNullFieldServer(t, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data": map[string]interface{}{
			"id":     "fw-123",
			"name":   "test-firewall",
			"status": nil, // Null value
			"label":  nil,
			"memory": nil,
		},
	})
	defer server.Close()

	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	resp, _, err := client.Firewall.GetFirewall(context.Background(), "fw-123")

	// Should handle null fields without panic
	if resp == nil && err == nil {
		t.Error("Expected response or error for null fields")
	}
}

func TestGetFirewall_InvalidFieldType(t *testing.T) {
	server := newInvalidFieldTypeServer(t, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data": map[string]interface{}{
			"id":              "fw-123",
			"name":            "test-firewall",
			"vm_id":           "not-an-int", // Should be int, not string
			"is_fortigate_vm": "not-a-bool", // Should be bool, not string
		},
	})
	defer server.Close()

	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	resp, _, err := client.Firewall.GetFirewall(context.Background(), "fw-123")

	// Should handle wrong type gracefully
	if resp == nil && err == nil {
		t.Error("Expected response or error for invalid field type")
	}
}

func TestGetFirewalls_EmptyList(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/fortigate/list/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": []
		}`)
	})

	result, _, err := ts.client.Firewall.GetFirewalls(context.Background())

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if len(result) != 0 {
		t.Errorf("Expected empty list, got %d items", len(result))
	}
}

func TestCreateFirewall_WithAllOptionalFields(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/fortigate/create/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusCreated, `{
			"code": 201,
			"message": "Firewall created successfully",
			"data": {
				"id": "fw-123",
				"vm_id": 12345,
				"name": "test-firewall",
				"label": "Test Firewall",
				"plan": "fortigate-small",
				"status": "active",
				"is_fortigate_vm": true,
				"vpc_id": "vpc-789",
				"cn_id": "cn-456"
			}
		}`)
	})

	createReq := &FirewallCreateRequest{
		Name:                 "test-firewall",
		Label:                "Test Firewall",
		Plan:                 "fortigate-small",
		Image:                "fortinet-7.2",
		VPCID:                "vpc-789",
		CNID:                 "cn-456",
		SSHKeys:              []string{"ssh-key-1", "ssh-key-2"},
		StartScripts:         []string{"#!/bin/bash", "echo 'Hello'"},
		Backups:              true,
		EnableBitninja:       true,
		DisablePassword:      true,
		IsSavedImage:         false,
		SavedImageTemplateID: 0,
		ReserveIP:            "203.0.113.100",
		IsIPv6Availed:        true,
		DefaultPublicIP:      true,
	}

	result, _, err := ts.client.Firewall.CreateFirewall(context.Background(), createReq)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.ID != "fw-123" {
		t.Errorf("Expected ID fw-123, got %s", result.ID)
	}
}
