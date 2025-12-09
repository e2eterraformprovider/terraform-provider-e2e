package goe2e

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoadBalancerService_CreateLoadBalancer(t *testing.T) {
	tests := []struct {
		name        string
		lbName      string
		planName    string
		lbMode      string
		statusCode  int
		expectError bool
		setup       func(*LoadBalancerCreateRequest)
	}{
		{
			name:       "valid create request",
			lbName:     "test-lb",
			planName:   "E2E-LB-2",
			lbMode:     "HTTP",
			statusCode: http.StatusOK,
		},
		{
			name:        "nil request",
			expectError: true,
		},
		{
			name:        "empty LB name",
			lbName:      "",
			planName:    "E2E-LB-2",
			lbMode:      "HTTP",
			expectError: true,
		},
		{
			name:        "empty plan name",
			lbName:      "test-lb",
			planName:    "",
			lbMode:      "HTTP",
			expectError: true,
		},
		{
			name:        "empty LB mode",
			lbName:      "test-lb",
			planName:    "E2E-LB-2",
			lbMode:      "",
			expectError: true,
		},
		{
			name:        "API error response",
			lbName:      "test-lb",
			planName:    "E2E-LB-2",
			lbMode:      "HTTP",
			statusCode:  http.StatusInternalServerError,
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertHTTPMethod(t, r, http.MethodPost)
				if !lbContains(r.URL.Path, "appliances/load-balancers") {
					t.Errorf("expected path to contain 'appliances/load-balancers', got %s", r.URL.Path)
				}
				w.WriteHeader(tt.statusCode)
				if tt.statusCode == http.StatusOK {
					resp := loadBalancerRoot{
						Code:    200,
						Message: "Success",
						Data: LoadBalancer{
							ID:       "lb-123",
							Name:     tt.lbName,
							PlanName: tt.planName,
							LBMode:   tt.lbMode,
						},
					}
					json.NewEncoder(w).Encode(resp)
				} else {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"code":    tt.statusCode,
						"message": "Error",
					})
				}
			}))
			defer server.Close()
			client, err := NewClient("key", "token", "proj-123", "Mumbai", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
			assertNoError(t, err)
			if tt.expectError && tt.lbName == "" && tt.planName != "" && tt.lbMode != "" {
				_, _, err := client.LoadBalancer.CreateLoadBalancer(context.Background(), &LoadBalancerCreateRequest{
					LBName:   tt.lbName,
					PlanName: tt.planName,
					LBMode:   tt.lbMode,
				})
				assertError(t, err, "")
				return
			}
			if tt.expectError && tt.lbName != "" && tt.planName == "" && tt.lbMode != "" {
				_, _, err := client.LoadBalancer.CreateLoadBalancer(context.Background(), &LoadBalancerCreateRequest{
					LBName:   tt.lbName,
					PlanName: tt.planName,
					LBMode:   tt.lbMode,
				})
				assertError(t, err, "")
				return
			}
			if tt.expectError && tt.lbName != "" && tt.planName != "" && tt.lbMode == "" {
				_, _, err := client.LoadBalancer.CreateLoadBalancer(context.Background(), &LoadBalancerCreateRequest{
					LBName:   tt.lbName,
					PlanName: tt.planName,
					LBMode:   tt.lbMode,
				})
				assertError(t, err, "")
				return
			}
			if tt.expectError && tt.statusCode != http.StatusOK {
				_, _, err := client.LoadBalancer.CreateLoadBalancer(context.Background(), &LoadBalancerCreateRequest{
					LBName:   tt.lbName,
					PlanName: tt.planName,
					LBMode:   tt.lbMode,
				})
				assertError(t, err, "")
				return
			}
			if !tt.expectError {
				lb, resp, err := client.LoadBalancer.CreateLoadBalancer(context.Background(), &LoadBalancerCreateRequest{
					LBName:   tt.lbName,
					PlanName: tt.planName,
					LBMode:   tt.lbMode,
				})
				assertNoError(t, err)
				assertNotNil(t, lb, "load balancer")
				if lb.Name != tt.lbName {
					t.Errorf("expected name %s, got %s", tt.lbName, lb.Name)
				}
				assertStatus(t, resp, http.StatusOK)
			}
		})
	}
}
func TestLoadBalancerService_GetLoadBalancer(t *testing.T) {
	tests := []struct {
		name        string
		lbID        string
		statusCode  int
		expectError bool
		expectNil   bool
	}{
		{
			name:       "valid get request",
			lbID:       "lb-123",
			statusCode: http.StatusOK,
		},
		{
			name:        "empty ID",
			lbID:        "",
			expectError: true,
		},
		{
			name:       "not found",
			lbID:       "lb-123",
			statusCode: http.StatusNotFound,
			expectNil:  true,
		},
		{
			name:        "server error",
			lbID:        "lb-123",
			statusCode:  http.StatusInternalServerError,
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertHTTPMethod(t, r, http.MethodGet)
				if tt.statusCode == http.StatusOK {
					resp := loadBalancerRoot{
						Code:    200,
						Message: "Success",
						Data: LoadBalancer{
							ID:       tt.lbID,
							Name:     "test-lb",
							PlanName: "E2E-LB-2",
						},
					}
					json.NewEncoder(w).Encode(resp)
				} else {
					w.WriteHeader(tt.statusCode)
					json.NewEncoder(w).Encode(map[string]interface{}{
						"code":    tt.statusCode,
						"message": "Error",
					})
				}
			}))
			defer server.Close()
			client, err := NewClient("key", "token", "proj-123", "Mumbai", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
			assertNoError(t, err)
			lb, resp, err := client.LoadBalancer.GetLoadBalancer(context.Background(), tt.lbID)
			if tt.expectError {
				assertError(t, err, "")
				return
			}
			if tt.expectNil && lb != nil {
				t.Error("expected nil load balancer for not found")
			}
			if !tt.expectNil && lb == nil {
				t.Fatal("expected non-nil load balancer")
			}
			assertNotNil(t, resp, "response")
		})
	}
}
func TestLoadBalancerService_UpdateLoadBalancer(t *testing.T) {
	tests := []struct {
		name        string
		lbID        string
		nilReq      bool
		statusCode  int
		expectError bool
	}{
		{
			name:       "valid update request",
			lbID:       "lb-123",
			statusCode: http.StatusOK,
		},
		{
			name:        "nil request",
			lbID:        "lb-123",
			nilReq:      true,
			expectError: true,
		},
		{
			name:        "API error",
			lbID:        "lb-123",
			statusCode:  http.StatusBadRequest,
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertHTTPMethod(t, r, http.MethodPut)
				if tt.statusCode == http.StatusOK {
					resp := loadBalancerRoot{
						Code:    200,
						Message: "Success",
						Data: LoadBalancer{
							Name:     "updated-lb",
							PlanName: "E2E-LB-3",
						},
					}
					json.NewEncoder(w).Encode(resp)
				} else {
					w.WriteHeader(tt.statusCode)
					json.NewEncoder(w).Encode(map[string]interface{}{
						"code":    tt.statusCode,
						"message": "Error",
					})
				}
			}))
			defer server.Close()
			client, err := NewClient("key", "token", "proj-123", "Mumbai", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
			assertNoError(t, err)
			if tt.nilReq {
				_, _, err := client.LoadBalancer.UpdateLoadBalancer(context.Background(), tt.lbID, nil)
				assertError(t, err, "")
				return
			}
			lb, _, err := client.LoadBalancer.UpdateLoadBalancer(context.Background(), tt.lbID, &LoadBalancerUpdateRequest{
				LBName: "updated-lb",
			})
			if tt.expectError {
				assertError(t, err, "")
				return
			}
			assertNoError(t, err)
			assertNotNil(t, lb, "load balancer")
		})
	}
}
func TestLoadBalancerService_DeleteLoadBalancer(t *testing.T) {
	tests := []struct {
		name        string
		lbID        string
		statusCode  int
		expectError bool
	}{
		{
			name:       "valid delete request",
			lbID:       "lb-123",
			statusCode: http.StatusNoContent,
		},
		{
			name:        "not found",
			lbID:        "lb-123",
			statusCode:  http.StatusNotFound,
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertHTTPMethod(t, r, http.MethodDelete)
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()
			client, err := NewClient("key", "token", "proj-123", "Mumbai", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
			assertNoError(t, err)
			resp, err := client.LoadBalancer.DeleteLoadBalancer(context.Background(), tt.lbID)
			if tt.expectError {
				assertError(t, err, "")
				return
			}
			assertNoError(t, err)
			assertNotNil(t, resp, "response")
		})
	}
}
func TestLoadBalancerService_UpdateLoadBalancerAction(t *testing.T) {
	tests := []struct {
		name        string
		lbID        string
		actionType  string
		statusCode  int
		expectError bool
		nilReq      bool
	}{
		{
			name:       "power action",
			lbID:       "lb-123",
			actionType: "power_on",
			statusCode: http.StatusOK,
		},
		{
			name:       "rename action",
			lbID:       "lb-123",
			actionType: "rename",
			statusCode: http.StatusOK,
		},
		{
			name:       "upgrade plan action",
			lbID:       "lb-123",
			actionType: "upgrade_plan",
			statusCode: http.StatusOK,
		},
		{
			name:        "nil request",
			lbID:        "lb-123",
			actionType:  "power_on",
			nilReq:      true,
			expectError: true,
		},
		{
			name:        "empty action type",
			lbID:        "lb-123",
			actionType:  "",
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertHTTPMethod(t, r, http.MethodPut)
				if !lbContains(r.URL.Path, "actions") {
					t.Errorf("expected path to contain 'actions', got %s", r.URL.Path)
				}
				if tt.statusCode == http.StatusOK {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"code":    tt.statusCode,
						"message": "Success",
					})
				} else {
					w.WriteHeader(tt.statusCode)
					json.NewEncoder(w).Encode(map[string]interface{}{
						"code":    tt.statusCode,
						"message": "Error",
					})
				}
			}))
			defer server.Close()
			client, err := NewClient("key", "token", "proj-123", "Mumbai", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
			assertNoError(t, err)
			if tt.nilReq {
				_, err := client.LoadBalancer.UpdateLoadBalancerAction(context.Background(), tt.lbID, nil)
				assertError(t, err, "")
				return
			}
			resp, err := client.LoadBalancer.UpdateLoadBalancerAction(context.Background(), tt.lbID, &LoadBalancerActionRequest{
				Type: tt.actionType,
				Name: "test-lb",
			})
			if tt.expectError {
				assertError(t, err, "")
				return
			}
			assertNoError(t, err)
			assertNotNil(t, resp, "response")
		})
	}
}
func TestLoadBalancerService_UpdateIPv6(t *testing.T) {
	tests := []struct {
		name        string
		lbID        string
		action      string
		statusCode  int
		expectError bool
		nilReq      bool
	}{
		{
			name:       "attach IPv6",
			lbID:       "lb-123",
			action:     "attach",
			statusCode: http.StatusOK,
		},
		{
			name:       "detach IPv6",
			lbID:       "lb-123",
			action:     "detach",
			statusCode: http.StatusOK,
		},
		{
			name:        "nil request",
			lbID:        "lb-123",
			action:      "attach",
			nilReq:      true,
			expectError: true,
		},
		{
			name:        "empty action",
			lbID:        "lb-123",
			action:      "",
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertHTTPMethod(t, r, http.MethodPut)
				if !lbContains(r.URL.Path, "ipv6") {
					t.Errorf("expected path to contain 'ipv6', got %s", r.URL.Path)
				}
				if tt.statusCode == http.StatusOK {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"code":    tt.statusCode,
						"message": "Success",
					})
				} else {
					w.WriteHeader(tt.statusCode)
					json.NewEncoder(w).Encode(map[string]interface{}{
						"code":    tt.statusCode,
						"message": "Error",
					})
				}
			}))
			defer server.Close()
			client, err := NewClient("key", "token", "proj-123", "Mumbai", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
			assertNoError(t, err)
			if tt.nilReq {
				_, err := client.LoadBalancer.UpdateIPv6(context.Background(), tt.lbID, nil)
				assertError(t, err, "")
				return
			}
			resp, err := client.LoadBalancer.UpdateIPv6(context.Background(), tt.lbID, &IPv6ActionRequest{
				Action:     tt.action,
				DetachIPv6: "::1",
			})
			if tt.expectError {
				assertError(t, err, "")
				return
			}
			assertNoError(t, err)
			assertNotNil(t, resp, "response")
		})
	}
}

// Helper function to check if string contains substring
func lbContains(s, substr string) bool {
	if len(s) >= len(substr) && len(substr) > 0 {
		for i := 0; i <= len(s)-len(substr); i++ {
			if s[i:i+len(substr)] == substr {
				return true
			}
		}
	}
	return s == substr || len(substr) == 0
}

// TestLoadBalancerCreateRequest tests request struct marshaling
func TestLoadBalancerCreateRequest_Marshal(t *testing.T) {
	req := LoadBalancerCreateRequest{
		PlanName:       "E2E-LB-2",
		LBName:         "test-lb",
		LBMode:         "HTTP",
		EnableBitNinja: true,
		Location:       "Mumbai",
	}
	data, err := json.Marshal(req)
	assertNoError(t, err)
	var unmarshaled LoadBalancerCreateRequest
	err = json.Unmarshal(data, &unmarshaled)
	assertNoError(t, err)
	if unmarshaled.PlanName != req.PlanName {
		t.Errorf("expected PlanName %s, got %s", req.PlanName, unmarshaled.PlanName)
	}
	if unmarshaled.LBName != req.LBName {
		t.Errorf("expected LBName %s, got %s", req.LBName, unmarshaled.LBName)
	}
}

// TestLoadBalancer_ComplexStructure tests complex nested structures
func TestLoadBalancer_ComplexStructure(t *testing.T) {
	lb := LoadBalancer{
		ID:       "lb-123",
		Name:     "test-lb",
		PlanName: "E2E-LB-2",
		Backends: []LBBackend{
			{
				Name:    "backend-1",
				Balance: "roundrobin",
				Servers: []LBServer{
					{
						BackendName: "server-1",
						BackendIP:   "10.0.0.1",
						BackendPort: "8080",
					},
				},
			},
		},
		VPCList: []LBVPCDetail{
			{
				VPCName:   "vpc-1",
				IPv4CIDR:  "10.0.0.0/16",
				NetworkID: 1,
			},
		},
	}
	data, err := json.Marshal(lb)
	assertNoError(t, err)
	var unmarshaled LoadBalancer
	err = json.Unmarshal(data, &unmarshaled)
	assertNoError(t, err)
	if unmarshaled.ID != lb.ID {
		t.Errorf("expected ID %s, got %s", lb.ID, unmarshaled.ID)
	}
	if len(unmarshaled.Backends) != 1 {
		t.Errorf("expected 1 backend, got %d", len(unmarshaled.Backends))
	}
	if unmarshaled.Backends[0].Servers[0].BackendIP != "10.0.0.1" {
		t.Errorf("expected backend IP 10.0.0.1, got %s", unmarshaled.Backends[0].Servers[0].BackendIP)
	}
}

// TestNewArgError tests error creation
func TestLoadBalancerService_ErrorHandling(t *testing.T) {
	client, _ := NewClient("key", "token", "proj-123", "Mumbai", WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	// Test CreateLoadBalancer with nil request
	_, _, err := client.LoadBalancer.CreateLoadBalancer(context.Background(), nil)
	assertError(t, err, "createReq")
	// Test GetLoadBalancer with empty ID
	_, _, err = client.LoadBalancer.GetLoadBalancer(context.Background(), "")
	assertError(t, err, "lbID")
	// Test DeleteLoadBalancer with empty ID
	_, err = client.LoadBalancer.DeleteLoadBalancer(context.Background(), "")
	assertError(t, err, "")
	// Test UpdateLoadBalancerAction with nil request
	_, err = client.LoadBalancer.UpdateLoadBalancerAction(context.Background(), "lb-123", nil)
	assertError(t, err, "")
	// Test UpdateIPv6 with nil request
	_, err = client.LoadBalancer.UpdateIPv6(context.Background(), "lb-123", nil)
	assertError(t, err, "")
}

// Additional HTTP status code error tests
func TestLoadBalancerService_CreateLoadBalancer_BadRequest(t *testing.T) {
	server := newErrorServer(t, http.StatusBadRequest, "Invalid plan name")
	defer server.Close()
	client, err := NewClient("key", "token", "proj-123", "Mumbai", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	_, _, err = client.LoadBalancer.CreateLoadBalancer(context.Background(), &LoadBalancerCreateRequest{
		LBName:   "test-lb",
		PlanName: "INVALID-PLAN",
		LBMode:   "HTTP",
	})
	assertError(t, err, "")
}
func TestLoadBalancerService_GetLoadBalancer_Forbidden(t *testing.T) {
	server := newErrorServer(t, http.StatusForbidden, "Access denied")
	defer server.Close()
	client, err := NewClient("key", "token", "proj-123", "Mumbai", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	_, _, err = client.LoadBalancer.GetLoadBalancer(context.Background(), "lb-123")
	assertError(t, err, "")
}

func TestLoadBalancerService_UpdateLoadBalancer_ServiceUnavailable(t *testing.T) {
	server := newErrorServer(t, http.StatusServiceUnavailable, "Service unavailable")
	defer server.Close()
	client, err := NewClient("key", "token", "proj-123", "Mumbai", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	_, _, err = client.LoadBalancer.UpdateLoadBalancer(context.Background(), "lb-123", &LoadBalancerUpdateRequest{
		LBName: "new-name",
	})
	assertError(t, err, "")
}

func TestLoadBalancerService_DeleteLoadBalancer_Conflict(t *testing.T) {
	server := newErrorServer(t, http.StatusConflict, "Cannot delete load balancer with attached resources")
	defer server.Close()
	client, err := NewClient("key", "token", "proj-123", "Mumbai", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	_, err = client.LoadBalancer.DeleteLoadBalancer(context.Background(), "lb-123")
	assertError(t, err, "")
}

func TestLoadBalancerService_UpdateAction_Unauthorized(t *testing.T) {
	server := newErrorServer(t, http.StatusUnauthorized, "Unauthorized")
	defer server.Close()
	client, err := NewClient("key", "token", "proj-123", "Mumbai", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	_, err = client.LoadBalancer.UpdateLoadBalancerAction(context.Background(), "lb-123", &LoadBalancerActionRequest{
		Type: "power_on",
		Name: "test",
	})
	assertError(t, err, "")
}

func TestLoadBalancerService_UpdateIPv6_GatewayTimeout(t *testing.T) {
	server := newErrorServer(t, http.StatusGatewayTimeout, "Gateway timeout")
	defer server.Close()
	client, err := NewClient("key", "token", "proj-123", "Mumbai", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	_, err = client.LoadBalancer.UpdateIPv6(context.Background(), "lb-123", &IPv6ActionRequest{
		Action: "attach",
	})
	assertError(t, err, "")
}

func TestLoadBalancerService_CreateLoadBalancer_UnprocessableEntity(t *testing.T) {
	server := newErrorServer(t, http.StatusUnprocessableEntity, "Invalid configuration")
	defer server.Close()
	client, err := NewClient("key", "token", "proj-123", "Mumbai", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	_, _, err = client.LoadBalancer.CreateLoadBalancer(context.Background(), &LoadBalancerCreateRequest{
		LBName:   "test-lb",
		PlanName: "E2E-LB-2",
		LBMode:   "INVALID_MODE",
	})
	assertError(t, err, "")
}

// Phase 2: Response Parsing & Edge Case Tests
func TestCreateLoadBalancer_MalformedJSON(t *testing.T) {
	server := newMalformedJSONServer(t)
	defer server.Close()
	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	_, _, err := client.LoadBalancer.CreateLoadBalancer(context.Background(), &LoadBalancerCreateRequest{
		LBName:   "test-lb",
		PlanName: "E2E-LB-2",
		LBMode:   "HA",
	})
	if err == nil {
		t.Error("Expected error for malformed JSON")
	}
}

func TestGetLoadBalancer_MissingRequiredFields(t *testing.T) {
	server := newMissingFieldServer(t, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data": map[string]interface{}{
			// Missing "lb_id" field
			"name": "test-lb",
		},
	})
	defer server.Close()
	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	resp, _, err := client.LoadBalancer.GetLoadBalancer(context.Background(), "lb-123")
	// Should handle missing fields gracefully
	if resp == nil && err == nil {
		t.Error("Expected response or error handling")
	}
}

func TestGetLoadBalancer_NullFieldValues(t *testing.T) {
	server := newNullFieldServer(t, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data": map[string]interface{}{
			"lb_id":          "lb-123",
			"name":           "test-lb",
			"health_check":   nil, // Null value
			"backend_config": nil,
			"listener_rules": nil,
		},
	})
	defer server.Close()
	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	resp, _, err := client.LoadBalancer.GetLoadBalancer(context.Background(), "lb-123")
	// Should handle null fields without panic
	if resp == nil && err == nil {
		t.Error("Expected response or error for null fields")
	}
}

func TestUpdateLoadBalancer_NullConfigFields(t *testing.T) {
	server := newNullFieldServer(t, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data": map[string]interface{}{
			"health_check": nil,
		},
	})
	defer server.Close()
	client, _ := NewClient("key", "token", "proj", "region", SetBaseURL(server.URL), noRetryOpt())
	resp, _, err := client.LoadBalancer.UpdateLoadBalancer(context.Background(), "lb-123", &LoadBalancerUpdateRequest{})
	// Should handle gracefully
	if resp == nil && err == nil {
		t.Error("Expected response or error")
	}
}
