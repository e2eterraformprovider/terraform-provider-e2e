package goe2e

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestCreateNode tests creating a new node
func TestCreateNode(t *testing.T) {
	data := map[string]interface{}{
		"id":     12345,
		"vm_id":  54321,
		"name":   "test-node",
		"plan":   "c1.small",
		"status": "Creating",
	}
	server := newSuccessServer(t, data)
	defer server.Close()
	client, err := NewClient("test-key", "test-token", "proj-123", "Chennai", SetBaseURL(server.URL+"/"))
	assertNoError(t, err)
	ctx := context.Background()
	createReq := &NodeCreateRequest{
		Name:  "test-node",
		Plan:  "c1.small",
		Image: "ubuntu-20.04",
	}
	node, resp, err := client.Nodes.CreateNode(ctx, createReq)
	assertNoError(t, err)
	assertNotNil(t, node, "Expected node to be returned")
	assertStatus(t, resp, http.StatusOK)
}

// TestCreateNodeMissingName tests that CreateNode fails without a name
func TestCreateNodeMissingName(t *testing.T) {
	client, err := NewClient("test-key", "test-token", "proj-123", "Chennai")
	assertNoError(t, err)
	ctx := context.Background()
	createReq := &NodeCreateRequest{
		Name:  "",
		Plan:  "c1.small",
		Image: "ubuntu-20.04",
	}
	_, _, err = client.Nodes.CreateNode(ctx, createReq)
	assertError(t, err, "")
}

// TestGetNode tests retrieving a node
func TestGetNode(t *testing.T) {
	data := map[string]interface{}{
		"id":     12345,
		"vm_id":  54321,
		"name":   "test-node",
		"plan":   "c1.small",
		"status": "Running",
	}
	server := newSuccessServer(t, data)
	defer server.Close()
	client, err := NewClient("test-key", "test-token", "proj-123", "Chennai", SetBaseURL(server.URL+"/"))
	assertNoError(t, err)
	ctx := context.Background()
	node, resp, err := client.Nodes.GetNode(ctx, "12345")
	assertNoError(t, err)
	assertNotNil(t, node, "Expected node to be returned")
	assertStatus(t, resp, http.StatusOK)
}

// TestGetNodeMissingID tests that GetNode fails with empty node ID
func TestGetNodeMissingID(t *testing.T) {
	client, err := NewClient("test-key", "test-token", "proj-123", "Chennai")
	assertNoError(t, err)
	ctx := context.Background()
	_, _, err = client.Nodes.GetNode(ctx, "")
	assertError(t, err, "")
}

// TestUpdateNode tests updating a node
func TestUpdateNode(t *testing.T) {
	data := map[string]interface{}{
		"id":   12345,
		"name": "updated-node",
	}
	server := newSuccessServer(t, data)
	defer server.Close()
	client, err := NewClient("test-key", "test-token", "proj-123", "Chennai", SetBaseURL(server.URL+"/"))
	assertNoError(t, err)
	ctx := context.Background()
	updateReq := &NodeUpdateRequest{
		Name: "updated-node",
	}
	node, resp, err := client.Nodes.UpdateNode(ctx, "12345", updateReq)
	assertNoError(t, err)
	assertNotNil(t, node, "Expected node to be returned")
	assertStatus(t, resp, http.StatusOK)
}

// TestDeleteNode tests deleting a node
func TestDeleteNode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodDelete)
		writeJSON(w, http.StatusOK, buildSuccessResponse(1, "OK", nil))
	}))
	defer server.Close()
	client, err := NewClient("test-key", "test-token", "proj-123", "Chennai", SetBaseURL(server.URL+"/"))
	assertNoError(t, err)
	ctx := context.Background()
	resp, err := client.Nodes.DeleteNode(ctx, "12345")
	assertNoError(t, err)
	assertNotNil(t, resp, "Expected response to be returned")
}

// TestPowerOn tests powering on a node
func TestPowerOn(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPut)
		writeJSON(w, http.StatusOK, buildSuccessResponse(200, "OK", nil))
	}))
	defer server.Close()
	client, err := NewClient("test-key", "test-token", "proj-123", "Chennai", SetBaseURL(server.URL+"/"))
	assertNoError(t, err)
	ctx := context.Background()
	resp, err := client.Nodes.PowerOn(ctx, "12345")
	assertNoError(t, err)
	assertNotNil(t, resp, "Expected response to be returned")
}

// TestPowerOff tests powering off a node
func TestPowerOff(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPut)
		writeJSON(w, http.StatusOK, buildSuccessResponse(200, "OK", nil))
	}))
	defer server.Close()
	client, err := NewClient("test-key", "test-token", "proj-123", "Chennai", SetBaseURL(server.URL+"/"))
	assertNoError(t, err)
	ctx := context.Background()
	resp, err := client.Nodes.PowerOff(ctx, "12345")
	assertNoError(t, err)
	assertNotNil(t, resp, "Expected response to be returned")
}

// TestReboot tests rebooting a node
func TestReboot(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPut)
		writeJSON(w, http.StatusOK, buildSuccessResponse(200, "OK", nil))
	}))
	defer server.Close()
	client, err := NewClient("test-key", "test-token", "proj-123", "Chennai", SetBaseURL(server.URL+"/"))
	assertNoError(t, err)
	ctx := context.Background()
	resp, err := client.Nodes.Reboot(ctx, "12345")
	assertNoError(t, err)
	assertNotNil(t, resp, "Expected response to be returned")
}

// TestReinstall tests reinstalling a node
func TestReinstall(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPut)
		writeJSON(w, http.StatusOK, buildSuccessResponse(200, "OK", nil))
	}))
	defer server.Close()
	client, err := NewClient("test-key", "test-token", "proj-123", "Chennai", SetBaseURL(server.URL+"/"))
	assertNoError(t, err)
	ctx := context.Background()
	reinstallReq := &NodeReinstallRequest{
		Image: "ubuntu-20.04",
	}
	resp, err := client.Nodes.Reinstall(ctx, "12345", reinstallReq)
	assertNoError(t, err)
	assertNotNil(t, resp, "Expected response to be returned")
}

// TestLockNode tests locking a node
func TestLockNode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPut)
		writeJSON(w, http.StatusOK, buildSuccessResponse(200, "OK", nil))
	}))
	defer server.Close()
	client, err := NewClient("test-key", "test-token", "proj-123", "Chennai", SetBaseURL(server.URL+"/"))
	assertNoError(t, err)
	ctx := context.Background()
	resp, err := client.Nodes.LockNode(ctx, "12345")
	assertNoError(t, err)
	assertNotNil(t, resp, "Expected response to be returned")
}

// TestUnlockNode tests unlocking a node
func TestUnlockNode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPut)
		writeJSON(w, http.StatusOK, buildSuccessResponse(200, "OK", nil))
	}))
	defer server.Close()
	client, err := NewClient("test-key", "test-token", "proj-123", "Chennai", SetBaseURL(server.URL+"/"))
	assertNoError(t, err)
	ctx := context.Background()
	resp, err := client.Nodes.UnlockNode(ctx, "12345")
	assertNoError(t, err)
	assertNotNil(t, resp, "Expected response to be returned")
}

// TestAttachSecurityGroup tests attaching a security group
func TestAttachSecurityGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPost)
		writeJSON(w, http.StatusOK, buildSuccessResponse(200, "OK", nil))
	}))
	defer server.Close()
	client, err := NewClient("test-key", "test-token", "proj-123", "Chennai", SetBaseURL(server.URL+"/"))
	assertNoError(t, err)
	ctx := context.Background()
	sgReq := &SecurityGroupRequest{
		SecurityGroupList: []int{1, 2, 3},
	}
	resp, err := client.Nodes.AttachSecurityGroup(ctx, "12345", sgReq)
	assertNoError(t, err)
	assertNotNil(t, resp, "Expected response to be returned")
}

// TestDetachSecurityGroup tests detaching a security group
func TestDetachSecurityGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPost)
		writeJSON(w, http.StatusOK, buildSuccessResponse(200, "OK", nil))
	}))
	defer server.Close()
	client, err := NewClient("test-key", "test-token", "proj-123", "Chennai", SetBaseURL(server.URL+"/"))
	assertNoError(t, err)
	ctx := context.Background()
	sgReq := &SecurityGroupRequest{
		SecurityGroupList: []int{1},
	}
	resp, err := client.Nodes.DetachSecurityGroup(ctx, "12345", sgReq)
	assertNoError(t, err)
	assertNotNil(t, resp, "Expected response to be returned")
}

// TestGetSecurityGroupList tests getting security group list
func TestGetSecurityGroupList(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"id": 1, "name": "default", "is_default": true},
		map[string]interface{}{"id": 2, "name": "web", "is_default": false},
	}
	server := newSuccessServer(t, data)
	defer server.Close()
	client, err := NewClient("test-key", "test-token", "proj-123", "Chennai", SetBaseURL(server.URL+"/"))
	assertNoError(t, err)
	ctx := context.Background()
	sgs, resp, err := client.Nodes.GetSecurityGroupList(ctx)
	assertNoError(t, err)
	assertNotNil(t, sgs, "Expected security group list to be returned")
	assertStatus(t, resp, http.StatusOK)
}

// TestAttachVPC tests attaching a VPC
func TestNodeAttachVPC(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPost)
		writeJSON(w, http.StatusOK, buildSuccessResponse(200, "OK", nil))
	}))
	defer server.Close()
	client, err := NewClient("test-key", "test-token", "proj-123", "Chennai", SetBaseURL(server.URL+"/"))
	assertNoError(t, err)
	ctx := context.Background()
	vpcReq := &NodeVPCAttachRequest{
		VPCID: "vpc-123",
	}
	resp, err := client.Nodes.AttachVPC(ctx, "12345", vpcReq)
	assertNoError(t, err)
	assertNotNil(t, resp, "Expected response to be returned")
}

// TestDetachVPC tests detaching a VPC
func TestNodeDetachVPC(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPost)
		writeJSON(w, http.StatusOK, buildSuccessResponse(200, "OK", nil))
	}))
	defer server.Close()
	client, err := NewClient("test-key", "test-token", "proj-123", "Chennai", SetBaseURL(server.URL+"/"))
	assertNoError(t, err)
	ctx := context.Background()
	resp, err := client.Nodes.DetachVPC(ctx, "12345")
	assertNoError(t, err)
	assertNotNil(t, resp, "Expected response to be returned")
}

// TestGetLCMState tests getting node LCM state
func TestGetLCMState(t *testing.T) {
	data := map[string]interface{}{
		"lcm_state": "RUNNING",
		"state":     "running",
	}
	server := newSuccessServer(t, data)
	defer server.Close()
	client, err := NewClient("test-key", "test-token", "proj-123", "Chennai", SetBaseURL(server.URL+"/"))
	assertNoError(t, err)
	ctx := context.Background()
	lcmState, resp, err := client.Nodes.GetLCMState(ctx, "12345")
	assertNoError(t, err)
	assertNotNil(t, lcmState, "Expected LCM state to be returned")
	assertStatus(t, resp, http.StatusOK)
}

// TestUpgradePlan tests upgrading a node's plan
func TestNodeUpgradePlan(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPut)
		writeJSON(w, http.StatusOK, buildSuccessResponse(200, "OK", nil))
	}))
	defer server.Close()
	client, err := NewClient("test-key", "test-token", "proj-123", "Chennai", SetBaseURL(server.URL+"/"))
	assertNoError(t, err)
	ctx := context.Background()
	upgradeReq := &NodePlanUpgradeRequest{
		Plan: "c1.medium",
	}
	resp, err := client.Nodes.UpgradePlan(ctx, "12345", upgradeReq)
	assertNoError(t, err)
	assertNotNil(t, resp, "Expected response to be returned")
}

// TestUpgradePlanMissingPlan tests that UpgradePlan fails without a plan
func TestUpgradePlanMissingPlan(t *testing.T) {
	client, err := NewClient("test-key", "test-token", "proj-123", "Chennai")
	assertNoError(t, err)
	ctx := context.Background()
	upgradeReq := &NodePlanUpgradeRequest{}
	_, err = client.Nodes.UpgradePlan(ctx, "12345", upgradeReq)
	assertError(t, err, "")
}

// Additional edge case tests for coverage
func TestUpdateNode_EmptyID(t *testing.T) {
	client, err := NewClient("test-key", "test-token", "proj-123", "Chennai")
	assertNoError(t, err)
	updateReq := &NodeUpdateRequest{
		Name: "new-name",
	}
	_, _, err = client.Nodes.UpdateNode(context.Background(), "", updateReq)
	assertError(t, err, "")
}

func TestUpdateNode_NilRequest(t *testing.T) {
	client, err := NewClient("test-key", "test-token", "proj-123", "Chennai")
	assertNoError(t, err)
	_, _, err = client.Nodes.UpdateNode(context.Background(), "node-123", nil)
	assertError(t, err, "")
}

func TestReinstall_NilRequest(t *testing.T) {
	client, err := NewClient("test-key", "test-token", "proj-123", "Chennai")
	assertNoError(t, err)
	_, err = client.Nodes.Reinstall(context.Background(), "node-123", nil)
	assertError(t, err, "")
}

func TestAttachSecurityGroup_NilRequest(t *testing.T) {
	client, err := NewClient("test-key", "test-token", "proj-123", "Chennai")
	assertNoError(t, err)
	_, err = client.Nodes.AttachSecurityGroup(context.Background(), "node-123", nil)
	assertError(t, err, "")
}

func TestDetachSecurityGroup_NilRequest(t *testing.T) {
	client, err := NewClient("test-key", "test-token", "proj-123", "Chennai")
	assertNoError(t, err)
	_, err = client.Nodes.DetachSecurityGroup(context.Background(), "node-123", nil)
	assertError(t, err, "")
}

func TestNodeAttachVPC_NilRequest(t *testing.T) {
	client, err := NewClient("test-key", "test-token", "proj-123", "Chennai")
	assertNoError(t, err)
	_, err = client.Nodes.AttachVPC(context.Background(), "node-123", nil)
	assertError(t, err, "")
}

// Additional error tests for better node coverage
func TestCreateNode_ErrorResponse(t *testing.T) {
	server := newErrorServer(t, http.StatusInternalServerError, "Internal server error")
	defer server.Close()
	client, err := NewClient("test-key", "test-token", "proj-123", "Chennai", SetBaseURL(server.URL+"/"), noRetryOpt())
	assertNoError(t, err)
	ctx := context.Background()
	createReq := &NodeCreateRequest{
		Name:  "test-node",
		Plan:  "c1.small",
		Image: "ubuntu-20.04",
	}
	_, _, err = client.Nodes.CreateNode(ctx, createReq)
	assertError(t, err, "")
}
