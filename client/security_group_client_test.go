package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func TestGetSecurityGroupList(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code": 200,
		"data": []interface{}{
			map[string]interface{}{
				"id":   "sg-1",
				"name": "default",
			},
			map[string]interface{}{
				"id":   "sg-2",
				"name": "custom",
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

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	result, err := client.GetSecurityGroupList("123", "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	code := result["code"].(float64)
	if code != 200 {
		t.Errorf("Expected code 200, got %v", code)
	}
}

func TestGetSecurityGroupListError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "server error"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	result, err := client.GetSecurityGroupList("123", "us-east")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestGetSecurityGroup(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code": 200,
		"data": []interface{}{
			map[string]interface{}{
				"id":   "sg-1",
				"name": "test-sg",
			},
			map[string]interface{}{
				"id":   "sg-2",
				"name": "other-sg",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	result, err := client.GetSecurityGroup("test-sg", "123", "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result["name"] != "test-sg" {
		t.Errorf("Expected name test-sg, got %v", result["name"])
	}
}

func TestGetSecurityGroupNotFound(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code": 200,
		"data": []interface{}{
			map[string]interface{}{
				"id":   "sg-1",
				"name": "other-sg",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	result, err := client.GetSecurityGroup("nonexistent-sg", "123", "us-east")

	if err == nil {
		t.Fatal("Expected error for not found security group, got nil")
	}

	if result != nil {
		t.Error("Expected nil result when security group not found")
	}
}

func TestCreateSecurityGroups(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/security_group/" {
			t.Errorf("Expected path /security_group/, got %s", r.URL.Path)
		}

		body, _ := ioutil.ReadAll(r.Body)
		var payload models.SecurityGroupCreateRequest
		json.Unmarshal(body, &payload)

		if payload.Name == "" {
			t.Error("Expected Name in request body")
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    200,
			"message": "Security group created",
		})
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	payload := models.SecurityGroupCreateRequest{
		Name:        "test-sg",
		Description: "Test security group",
	}

	err := client.CreateSecurityGroups(payload, "123", "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestCreateSecurityGroupsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid request"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	payload := models.SecurityGroupCreateRequest{
		Name: "test-sg",
	}

	err := client.CreateSecurityGroups(payload, "123", "us-east")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestUpdateSecurityGroups(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/security_group/sg-123/" {
			t.Errorf("Expected path /security_group/sg-123/, got %s", r.URL.Path)
		}

		body, _ := ioutil.ReadAll(r.Body)
		var payload models.SecurityGroupUpdateRequest
		json.Unmarshal(body, &payload)

		if payload.Name == "" {
			t.Error("Expected Name in request body")
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    200,
			"message": "Security group updated",
		})
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	payload := models.SecurityGroupUpdateRequest{
		Name:        "updated-sg",
		Description: "Updated security group",
	}

	err := client.UpdateSecurityGroups(payload, "sg-123", "123", "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestUpdateSecurityGroupsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "not found"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	payload := models.SecurityGroupUpdateRequest{
		Name: "updated-sg",
	}

	err := client.UpdateSecurityGroups(payload, "sg-404", "123", "us-east")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestMakeDefaultSecurityGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/security_group/sg-123/mark-default/" {
			t.Errorf("Expected path /security_group/sg-123/mark-default/, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    200,
			"message": "Security group marked as default",
		})
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	err := client.MakeDefaultSecurityGroup("sg-123", "123", "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestMakeDefaultSecurityGroupError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "cannot mark as default"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	err := client.MakeDefaultSecurityGroup("sg-123", "123", "us-east")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestDetachSecurityGroup(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code":    200,
		"message": "Security group detached",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/security_group/123/detach/" {
			t.Errorf("Expected path /security_group/123/detach/, got %s", r.URL.Path)
		}

		body, _ := ioutil.ReadAll(r.Body)
		var payload models.UpdateSecurityGroups
		json.Unmarshal(body, &payload)

		if len(payload.SecurityGroupList) == 0 {
			t.Error("Expected SecurityGroupList in request body")
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	payload := &models.UpdateSecurityGroups{
		SecurityGroupList: []int{1, 2},
	}

	result, err := client.DetachSecurityGroup(payload, 123, "456", "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestDetachSecurityGroupError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid request"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	payload := &models.UpdateSecurityGroups{
		SecurityGroupList: []int{1},
	}

	result, err := client.DetachSecurityGroup(payload, 123, "456", "us-east")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestAttachSecurityGroup(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code":    200,
		"message": "Security group attached",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/security_group/123/attach/" {
			t.Errorf("Expected path /security_group/123/attach/, got %s", r.URL.Path)
		}

		body, _ := ioutil.ReadAll(r.Body)
		var payload models.UpdateSecurityGroups
		json.Unmarshal(body, &payload)

		if len(payload.SecurityGroupList) == 0 {
			t.Error("Expected SecurityGroupList in request body")
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	payload := &models.UpdateSecurityGroups{
		SecurityGroupList: []int{1, 2},
	}

	result, err := client.AttachSecurityGroup(payload, 123, "456", "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestAttachSecurityGroupError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid request"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	payload := &models.UpdateSecurityGroups{
		SecurityGroupList: []int{1},
	}

	result, err := client.AttachSecurityGroup(payload, 123, "456", "us-east")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestDeleteSecurityGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/security_group/sg-123/" {
			t.Errorf("Expected path /security_group/sg-123/, got %s", r.URL.Path)
		}

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

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    200,
			"message": "Security group deleted",
		})
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	err := client.DeleteSecurityGroup("sg-123", "123", "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestDeleteSecurityGroupError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(`{"error": "security group in use"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	err := client.DeleteSecurityGroup("sg-123", "123", "us-east")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestGetSecurityGroupInvalidData(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code": 200,
		"data": "invalid-not-an-array",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	result, err := client.GetSecurityGroup("test-sg", "123", "us-east")

	if err == nil {
		t.Fatal("Expected error for invalid data format, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestCreateSecurityGroupsWithHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			t.Error("Expected Authorization header")
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		if r.Header.Get("User-Agent") != "terraform-e2e" {
			t.Errorf("Expected User-Agent terraform-e2e, got %s", r.Header.Get("User-Agent"))
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code": 200,
		})
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	payload := models.SecurityGroupCreateRequest{
		Name: "test-sg",
	}

	err := client.CreateSecurityGroups(payload, "123", "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}
