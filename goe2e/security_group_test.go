package goe2e

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestCreateSecurityGroup_Success tests successful creation of a security group
func TestCreateSecurityGroup_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPost)
		testURLPath(t, r, "/security_group/")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")
		writeJSON(w, http.StatusCreated, `{
			"code": 201,
			"message": "Security group created",
			"data": {
				"id": "sg-123",
				"name": "test-sg",
				"description": "Test security group",
				"is_default": false,
				"rules": []
			}
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	ctx := context.Background()
	req := &SecurityGroupCreateRequest{
		Name:        "test-sg",
		Description: "Test security group",
		Rules:       []Rule{},
		Default:     false,
	}
	sg, _, err := client.SecurityGroups.CreateSecurityGroup(ctx, req)
	assertNoError(t, err)
	assertNotNil(t, sg, "Expected security group")
	if sg.Name != "test-sg" {
		t.Errorf("Expected Name test-sg, got %s", sg.Name)
	}
	if sg.ID != "sg-123" {
		t.Errorf("Expected ID sg-123, got %s", sg.ID)
	}
}

// TestCreateSecurityGroup_NilRequest tests CreateSecurityGroup with nil request
func TestCreateSecurityGroup_NilRequest(t *testing.T) {
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	ctx := context.Background()
	_, _, err := client.SecurityGroups.CreateSecurityGroup(ctx, nil)
	assertError(t, err, "")
}

// TestCreateSecurityGroup_EmptyName tests CreateSecurityGroup with empty name
func TestCreateSecurityGroup_EmptyName(t *testing.T) {
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	ctx := context.Background()
	req := &SecurityGroupCreateRequest{
		Name: "",
	}
	_, _, err := client.SecurityGroups.CreateSecurityGroup(ctx, req)
	assertError(t, err, "")
}

// TestCreateSecurityGroup_ServerError tests CreateSecurityGroup with server error
func TestCreateSecurityGroup_ServerError(t *testing.T) {
	server := newErrorServer(t, http.StatusInternalServerError, "Internal error")
	defer server.Close()
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	ctx := context.Background()
	req := &SecurityGroupCreateRequest{
		Name: "test-sg",
	}
	_, _, err := client.SecurityGroups.CreateSecurityGroup(ctx, req)
	assertError(t, err, "")
}

// TestGetSecurityGroup_Success tests successful retrieval of a security group
func TestGetSecurityGroup_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/security_group/sg-123/")
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Security group retrieved",
			"data": {
				"id": "sg-123",
				"name": "test-sg",
				"description": "Test security group"
			}
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	ctx := context.Background()
	sg, _, err := client.SecurityGroups.GetSecurityGroup(ctx, "sg-123")
	assertNoError(t, err)
	assertNotNil(t, sg, "Expected security group to be returned")
}

// TestGetSecurityGroup_EmptyID tests GetSecurityGroup with empty ID
func TestGetSecurityGroup_EmptyID(t *testing.T) {
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	ctx := context.Background()
	_, _, err := client.SecurityGroups.GetSecurityGroup(ctx, "")
	assertError(t, err, "")
}

// TestGetSecurityGroup_NotFound tests GetSecurityGroup with 404 response
func TestGetSecurityGroup_NotFound(t *testing.T) {
	server := newErrorServer(t, http.StatusNotFound, "Not found")
	defer server.Close()
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	ctx := context.Background()
	sg, resp, err := client.SecurityGroups.GetSecurityGroup(ctx, "sg-notfound")
	if sg != nil {
		t.Errorf("Expected nil security group for 404, got %v", sg)
	}
	assertStatus(t, resp, http.StatusNotFound)
	assertError(t, err, "")
}

// TestGetSecurityGroup_OtherError tests GetSecurityGroup with non-404 error
func TestGetSecurityGroup_OtherError(t *testing.T) {
	server := newErrorServer(t, http.StatusForbidden, "Forbidden")
	defer server.Close()
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	ctx := context.Background()
	sg, _, err := client.SecurityGroups.GetSecurityGroup(ctx, "sg-err")
	if sg != nil {
		t.Errorf("Expected nil security group, got %v", sg)
	}
	assertError(t, err, "")
}

// TestUpdateSecurityGroup_Success tests successful update of a security group
func TestUpdateSecurityGroup_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPut)
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Security group updated",
			"data": {
				"id": "sg-123",
				"name": "updated-sg",
				"description": "Updated security group"
			}
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	ctx := context.Background()
	updateReq := &SecurityGroupUpdateRequest{
		Name:        "updated-sg",
		Description: "Updated security group",
	}
	sg, _, err := client.SecurityGroups.UpdateSecurityGroup(ctx, "sg-123", updateReq)
	assertNoError(t, err)
	if sg.Name != "updated-sg" {
		t.Errorf("Expected Name updated-sg, got %s", sg.Name)
	}
}

// TestUpdateSecurityGroup_EmptyID tests UpdateSecurityGroup with empty ID
func TestUpdateSecurityGroup_EmptyID(t *testing.T) {
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	ctx := context.Background()
	updateReq := &SecurityGroupUpdateRequest{
		Name: "updated-sg",
	}
	_, _, err := client.SecurityGroups.UpdateSecurityGroup(ctx, "", updateReq)
	assertError(t, err, "")
}

// TestUpdateSecurityGroup_NilRequest tests UpdateSecurityGroup with nil request
func TestUpdateSecurityGroup_NilRequest(t *testing.T) {
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	ctx := context.Background()
	_, _, err := client.SecurityGroups.UpdateSecurityGroup(ctx, "sg-123", nil)
	assertError(t, err, "")
}

// TestUpdateSecurityGroup_ServerError tests UpdateSecurityGroup with server error
func TestUpdateSecurityGroup_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusBadRequest, `{"code": 400, "message": "Bad request"}`)
	}))
	defer server.Close()
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	ctx := context.Background()
	updateReq := &SecurityGroupUpdateRequest{
		Name: "updated-sg",
	}
	_, _, err := client.SecurityGroups.UpdateSecurityGroup(ctx, "sg-123", updateReq)
	assertError(t, err, "")
}

// TestDeleteSecurityGroup_Success tests successful deletion of a security group
func TestDeleteSecurityGroup_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodDelete)
		writeJSON(w, http.StatusOK, `{"code": 200, "message": "Success"}`)
	}))
	defer server.Close()
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	ctx := context.Background()
	_, err := client.SecurityGroups.DeleteSecurityGroup(ctx, "sg-123")
	assertNoError(t, err)
}

// TestDeleteSecurityGroup_EmptyID tests DeleteSecurityGroup with empty ID
func TestDeleteSecurityGroup_EmptyID(t *testing.T) {
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	ctx := context.Background()
	_, err := client.SecurityGroups.DeleteSecurityGroup(ctx, "")
	assertError(t, err, "")
}

// TestDeleteSecurityGroup_ServerError tests DeleteSecurityGroup with server error
func TestDeleteSecurityGroup_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, `{"code": 500, "message": "Internal error"}`)
	}))
	defer server.Close()
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	ctx := context.Background()
	_, err := client.SecurityGroups.DeleteSecurityGroup(ctx, "sg-123")
	assertError(t, err, "")
}

// TestDeleteSecurityGroup_Forbidden tests DeleteSecurityGroup with forbidden error
func TestDeleteSecurityGroup_Forbidden(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusForbidden, `{"code": 403, "message": "Cannot delete default security group"}`)
	}))
	defer server.Close()
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	ctx := context.Background()
	_, err := client.SecurityGroups.DeleteSecurityGroup(ctx, "sg-protected")
	assertError(t, err, "")
}

// TestGetSecurityGroupList_Success tests successful retrieval of security group list
func TestGetSecurityGroupList_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Security groups retrieved",
			"data": [
				{
					"id": "sg-1",
					"name": "sg-name-1",
					"description": "First SG",
					"is_default": true,
					"rules": []
				},
				{
					"id": "sg-2",
					"name": "sg-name-2",
					"description": "Second SG",
					"is_default": false,
					"rules": []
				}
			]
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	ctx := context.Background()
	sgs, _, err := client.SecurityGroups.GetSecurityGroupList(ctx)
	assertNoError(t, err)
	assertNotNil(t, sgs, "Expected security group list")
	if len(sgs) != 2 {
		t.Errorf("Expected 2 security groups, got %d", len(sgs))
	}
	if sgs[0].Name != "sg-name-1" {
		t.Errorf("Expected first SG name sg-name-1, got %s", sgs[0].Name)
	}
	if sgs[1].Name != "sg-name-2" {
		t.Errorf("Expected second SG name sg-name-2, got %s", sgs[1].Name)
	}
}

// TestGetSecurityGroupList_Empty tests retrieval of empty security group list
func TestGetSecurityGroupList_Empty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Success",
			"data": []
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	ctx := context.Background()
	sgs, _, err := client.SecurityGroups.GetSecurityGroupList(ctx)
	assertNoError(t, err)
	assertNotNil(t, sgs, "Expected empty list")
	if len(sgs) != 0 {
		t.Errorf("Expected 0 security groups, got %d", len(sgs))
	}
}

// TestGetSecurityGroupList_ServerError tests GetSecurityGroupList with server error
func TestGetSecurityGroupList_ServerError(t *testing.T) {
	server := newErrorServer(t, http.StatusInternalServerError, "Internal server error")
	defer server.Close()
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	ctx := context.Background()
	_, _, err := client.SecurityGroups.GetSecurityGroupList(ctx)
	assertError(t, err, "")
}

// TestMakeDefaultSecurityGroup_Success tests successful marking of default security group
func TestMakeDefaultSecurityGroup_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/security_group/sg-123/mark-default/")
		writeJSON(w, http.StatusOK, `{"code": 200, "message": "Success"}`)
	}))
	defer server.Close()
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	ctx := context.Background()
	_, err := client.SecurityGroups.MakeDefaultSecurityGroup(ctx, "sg-123")
	assertNoError(t, err)
}

// TestMakeDefaultSecurityGroup_EmptyID tests MakeDefaultSecurityGroup with empty ID
func TestMakeDefaultSecurityGroup_EmptyID(t *testing.T) {
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	ctx := context.Background()
	_, err := client.SecurityGroups.MakeDefaultSecurityGroup(ctx, "")
	assertError(t, err, "")
}

// TestMakeDefaultSecurityGroup_ServerError tests MakeDefaultSecurityGroup with server error
func TestMakeDefaultSecurityGroup_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		writeJSON(w, http.StatusNotFound, `{"code": 404, "message": "Not found"}`)
	}))
	defer server.Close()
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	ctx := context.Background()
	_, err := client.SecurityGroups.MakeDefaultSecurityGroup(ctx, "sg-123")
	assertError(t, err, "")
}

// TestAttachSecurityGroup_Success tests successful attachment of security group to node
func TestAttachSecurityGroup_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPost)
		testURLPath(t, r, "/security_group/123/attach/")
		writeJSON(w, http.StatusOK, `{"code": 200, "message": "Success"}`)
	}))
	defer server.Close()
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	ctx := context.Background()
	attachReq := &SecurityGroupAttachRequest{
		SecurityGroupIDs: []int{1, 2},
	}
	_, err := client.SecurityGroups.AttachSecurityGroup(ctx, 123, attachReq)
	assertNoError(t, err)
}

// TestAttachSecurityGroup_InvalidVMID tests AttachSecurityGroup with invalid VM ID
func TestAttachSecurityGroup_InvalidVMID(t *testing.T) {
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	ctx := context.Background()
	attachReq := &SecurityGroupAttachRequest{
		SecurityGroupIDs: []int{1},
	}
	_, err := client.SecurityGroups.AttachSecurityGroup(ctx, 0, attachReq)
	assertError(t, err, "")
}

// TestAttachSecurityGroup_NegativeVMID tests AttachSecurityGroup with negative VM ID
func TestAttachSecurityGroup_NegativeVMID(t *testing.T) {
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	ctx := context.Background()
	attachReq := &SecurityGroupAttachRequest{
		SecurityGroupIDs: []int{1},
	}
	_, err := client.SecurityGroups.AttachSecurityGroup(ctx, -5, attachReq)
	assertError(t, err, "")
}

// TestSecurityGroupAttachNilRequest tests AttachSecurityGroup with nil request
func TestSecurityGroupAttachNilRequest(t *testing.T) {
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	ctx := context.Background()
	_, err := client.SecurityGroups.AttachSecurityGroup(ctx, 123, nil)
	assertError(t, err, "")
}

// TestAttachSecurityGroup_ServerError tests AttachSecurityGroup with server error
func TestAttachSecurityGroup_ServerError(t *testing.T) {
	server := newErrorServer(t, http.StatusInternalServerError, "Internal server error")
	defer server.Close()
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	ctx := context.Background()
	attachReq := &SecurityGroupAttachRequest{
		SecurityGroupIDs: []int{1},
	}
	_, err := client.SecurityGroups.AttachSecurityGroup(ctx, 123, attachReq)
	assertError(t, err, "")
}

// TestDetachSecurityGroup_Success tests successful detachment of security group from node
func TestDetachSecurityGroup_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPost)
		testURLPath(t, r, "/security_group/456/detach/")
		writeJSON(w, http.StatusOK, `{"code": 200, "message": "Success"}`)
	}))
	defer server.Close()
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	ctx := context.Background()
	detachReq := &SecurityGroupAttachRequest{
		SecurityGroupIDs: []int{1},
	}
	_, err := client.SecurityGroups.DetachSecurityGroup(ctx, 456, detachReq)
	assertNoError(t, err)
}

// TestDetachSecurityGroup_InvalidVMID tests DetachSecurityGroup with invalid VM ID
func TestDetachSecurityGroup_InvalidVMID(t *testing.T) {
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	ctx := context.Background()
	detachReq := &SecurityGroupAttachRequest{
		SecurityGroupIDs: []int{1},
	}
	_, err := client.SecurityGroups.DetachSecurityGroup(ctx, 0, detachReq)
	assertError(t, err, "")
}

// TestDetachSecurityGroup_NegativeVMID tests DetachSecurityGroup with negative VM ID
func TestDetachSecurityGroup_NegativeVMID(t *testing.T) {
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	ctx := context.Background()
	detachReq := &SecurityGroupAttachRequest{
		SecurityGroupIDs: []int{1},
	}
	_, err := client.SecurityGroups.DetachSecurityGroup(ctx, -10, detachReq)
	assertError(t, err, "")
}

// TestSecurityGroupDetachNilRequest tests DetachSecurityGroup with nil request
func TestSecurityGroupDetachNilRequest(t *testing.T) {
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	ctx := context.Background()
	_, err := client.SecurityGroups.DetachSecurityGroup(ctx, 456, nil)
	assertError(t, err, "")
}

// TestDetachSecurityGroup_ServerError tests DetachSecurityGroup with server error
func TestDetachSecurityGroup_ServerError(t *testing.T) {
	server := newErrorServer(t, http.StatusInternalServerError, "Internal server error")
	defer server.Close()
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	ctx := context.Background()
	detachReq := &SecurityGroupAttachRequest{
		SecurityGroupIDs: []int{1},
	}
	_, err := client.SecurityGroups.DetachSecurityGroup(ctx, 456, detachReq)
	assertError(t, err, "")
}

// TestSecurityGroupService_ContextCanceled tests that operations handle canceled context
func TestSecurityGroupService_ContextCanceled(t *testing.T) {
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, _, err := client.SecurityGroups.CreateSecurityGroup(ctx, &SecurityGroupCreateRequest{
		Name: "test-sg",
	})
	if err != context.Canceled && err != nil && !strings.Contains(err.Error(), "context canceled") {
		t.Errorf("Expected context canceled error, got %v", err)
	}
}

// TestCreateSecurityGroup_WithoutDescription tests creation without description
func TestCreateSecurityGroup_WithoutDescription(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPost)
		writeJSON(w, http.StatusCreated, `{
			"code": 201,
			"message": "Security group created",
			"data": {
				"id": "sg-789",
				"name": "minimal-sg",
				"description": ""
			}
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	ctx := context.Background()
	req := &SecurityGroupCreateRequest{
		Name: "minimal-sg",
	}
	sg, _, err := client.SecurityGroups.CreateSecurityGroup(ctx, req)
	assertNoError(t, err)
	if sg.Name != "minimal-sg" {
		t.Errorf("Expected Name minimal-sg, got %s", sg.Name)
	}
}

// TestUpdateSecurityGroup_WithoutRules tests update without rules
func TestUpdateSecurityGroup_WithoutRules(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPut)
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Updated",
			"data": {
				"id": "sg-999",
				"name": "updated",
				"description": "no rules"
			}
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	ctx := context.Background()
	updateReq := &SecurityGroupUpdateRequest{
		Name:        "updated",
		Description: "no rules",
	}
	sg, _, err := client.SecurityGroups.UpdateSecurityGroup(ctx, "sg-999", updateReq)
	assertNoError(t, err)
	if sg.Name != "updated" {
		t.Errorf("Expected Name updated, got %s", sg.Name)
	}
}

// TestMakeDefaultSecurityGroup_WithResponse tests MakeDefaultSecurityGroup with body response
func TestMakeDefaultSecurityGroup_WithResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPut)
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Marked as default"
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	ctx := context.Background()
	_, err := client.SecurityGroups.MakeDefaultSecurityGroup(ctx, "sg-123")
	assertNoError(t, err)
}

// TestAttachSecurityGroup_EmptySecurityGroupList tests AttachSecurityGroup with empty list
func TestAttachSecurityGroup_EmptySecurityGroupList(t *testing.T) {
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	ctx := context.Background()
	attachReq := &SecurityGroupAttachRequest{
		SecurityGroupIDs: []int{},
	}
	_, err := client.SecurityGroups.AttachSecurityGroup(ctx, 789, attachReq)
	assertError(t, err, "")
}

// TestDetachSecurityGroup_EmptySecurityGroupList tests DetachSecurityGroup with empty list
func TestDetachSecurityGroup_EmptySecurityGroupList(t *testing.T) {
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	ctx := context.Background()
	detachReq := &SecurityGroupAttachRequest{
		SecurityGroupIDs: []int{},
	}
	_, err := client.SecurityGroups.DetachSecurityGroup(ctx, 999, detachReq)
	assertError(t, err, "")
}

// TestGetSecurityGroupList_ManyItems tests retrieval with many security groups
func TestGetSecurityGroupList_ManyItems(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Success",
			"data": [
				{"id": "sg-1", "name": "sg-1", "description": "", "is_default": false, "rules": []},
				{"id": "sg-2", "name": "sg-2", "description": "", "is_default": false, "rules": []},
				{"id": "sg-3", "name": "sg-3", "description": "", "is_default": false, "rules": []},
				{"id": "sg-4", "name": "sg-4", "description": "", "is_default": false, "rules": []},
				{"id": "sg-5", "name": "sg-5", "description": "", "is_default": true, "rules": []}
			]
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	ctx := context.Background()
	sgs, _, err := client.SecurityGroups.GetSecurityGroupList(ctx)
	assertNoError(t, err)
	if len(sgs) != 5 {
		t.Errorf("Expected 5 security groups, got %d", len(sgs))
	}
	if !sgs[4].IsDefault {
		t.Error("Expected 5th SG to be default")
	}
}

// TestCreateSecurityGroup_Default tests creation with default flag set
func TestCreateSecurityGroup_Default(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPost)
		writeJSON(w, http.StatusCreated, `{
			"code": 201,
			"message": "Security group created",
			"data": {
				"id": "sg-def",
				"name": "default-sg",
				"description": "Default SG",
				"is_default": true
			}
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	ctx := context.Background()
	req := &SecurityGroupCreateRequest{
		Name:        "default-sg",
		Description: "Default SG",
		Default:     true,
	}
	sg, _, err := client.SecurityGroups.CreateSecurityGroup(ctx, req)
	assertNoError(t, err)
	if !sg.IsDefault {
		t.Error("Expected security group to be default")
	}
}

// TestAttachSecurityGroup_LargeList tests AttachSecurityGroup with many security groups
func TestAttachSecurityGroup_LargeList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPost)
		writeJSON(w, http.StatusOK, `{"code": 200, "message": "Success"}`)
	}))
	defer server.Close()
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	ctx := context.Background()
	attachReq := &SecurityGroupAttachRequest{
		SecurityGroupIDs: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}
	_, err := client.SecurityGroups.AttachSecurityGroup(ctx, 555, attachReq)
	assertNoError(t, err)
}

// TestDetachSecurityGroup_LargeList tests DetachSecurityGroup with many security groups
func TestDetachSecurityGroup_LargeList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPost)
		writeJSON(w, http.StatusOK, `{"code": 200, "message": "Success"}`)
	}))
	defer server.Close()
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	ctx := context.Background()
	detachReq := &SecurityGroupAttachRequest{
		SecurityGroupIDs: []int{1, 2, 3, 4, 5},
	}
	_, err := client.SecurityGroups.DetachSecurityGroup(ctx, 666, detachReq)
	assertNoError(t, err)
}

// TestSecurityGroupService_MultipleRules tests security group with multiple rules
func TestSecurityGroupService_MultipleRules(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Success",
			"data": {
				"id": "sg-456",
				"name": "complex-sg",
				"description": "SG with rules",
				"rules": [
					{
						"id": 1,
						"rule_type": "ingress",
						"protocol_name": "tcp",
						"port_range": "80",
						"network": "0.0.0.0/0",
						"description": "HTTP"
					},
					{
						"id": 2,
						"rule_type": "ingress",
						"protocol_name": "tcp",
						"port_range": "443",
						"network": "0.0.0.0/0",
						"description": "HTTPS"
					}
				]
			}
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	ctx := context.Background()
	rules := []Rule{
		{
			RuleType:     "ingress",
			ProtocolName: "tcp",
			PortRange:    "80",
			Network:      "0.0.0.0/0",
			Description:  "HTTP",
		},
		{
			RuleType:     "ingress",
			ProtocolName: "tcp",
			PortRange:    "443",
			Network:      "0.0.0.0/0",
			Description:  "HTTPS",
		},
	}
	req := &SecurityGroupCreateRequest{
		Name:        "complex-sg",
		Description: "SG with rules",
		Rules:       rules,
	}
	sg, _, err := client.SecurityGroups.CreateSecurityGroup(ctx, req)
	assertNoError(t, err)
	if len(sg.Rules) != 2 {
		t.Errorf("Expected 2 rules, got %d", len(sg.Rules))
	}
}
