package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func TestAddParamsAndHeader(t *testing.T) {
	req := httptest.NewRequest("GET", "https://api.test.com/test/", nil)
	location := "us-east"
	projectID := "123"

	client := NewClient("test-key", "test-token", "https://api.test.com/")
	modifiedReq, err := client.AddParamsAndHeader(req, location, projectID)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	params := modifiedReq.URL.Query()
	if params.Get("apikey") != "test-key" {
		t.Errorf("Expected apikey test-key, got %s", params.Get("apikey"))
	}

	if params.Get("location") != location {
		t.Errorf("Expected location %s, got %s", location, params.Get("location"))
	}

	if params.Get("project_id") != projectID {
		t.Errorf("Expected project_id %s, got %s", projectID, params.Get("project_id"))
	}

	if params.Get("contact_person_id") != "null" {
		t.Errorf("Expected contact_person_id null, got %s", params.Get("contact_person_id"))
	}

	if modifiedReq.Header.Get("Authorization") != "Bearer test-token" {
		t.Errorf("Expected Authorization header Bearer test-token, got %s", modifiedReq.Header.Get("Authorization"))
	}

	if modifiedReq.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", modifiedReq.Header.Get("Content-Type"))
	}

	if modifiedReq.Header.Get("User-Agent") != "terraform-e2e" {
		t.Errorf("Expected User-Agent terraform-e2e, got %s", modifiedReq.Header.Get("User-Agent"))
	}
}

func TestNewLoadBalancer(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code":    200,
		"message": "Load balancer created successfully",
		"data": map[string]interface{}{
			"id":   "lb-123",
			"name": "test-lb",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/appliances/load-balancers/" {
			t.Errorf("Expected path /appliances/load-balancers/, got %s", r.URL.Path)
		}

		query := r.URL.Query()
		if query.Get("apikey") == "" {
			t.Error("Expected apikey parameter")
		}
		if query.Get("project_id") == "" {
			t.Error("Expected project_id parameter")
		}
		if query.Get("location") == "" {
			t.Error("Expected location parameter")
		}

		body, _ := ioutil.ReadAll(r.Body)
		var lbCreate models.LoadBalancerCreate
		json.Unmarshal(body, &lbCreate)

		if lbCreate.LbName == "" {
			t.Error("Expected LbName in request body")
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	lbCreate := &models.LoadBalancerCreate{
		LbName:   "test-lb",
		Location: "us-east",
	}

	result, err := client.NewLoadBalancer(lbCreate, "123")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result["message"] != mockResponse["message"] {
		t.Errorf("Expected message %s, got %s", mockResponse["message"], result["message"])
	}
}

func TestGetLoadBalancerInfo(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{
			"id":    "lb-123",
			"name":  "test-lb",
			"state": "ACTIVE",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/appliances/lb-123/" {
			t.Errorf("Expected path /appliances/lb-123/, got %s", r.URL.Path)
		}

		query := r.URL.Query()
		if query.Get("apikey") == "" {
			t.Error("Expected apikey parameter")
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	result, err := client.GetLoadBalancerInfo("lb-123", "us-east", "123")

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

func TestGetLoadBalancerInfoError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "not found"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	result, err := client.GetLoadBalancerInfo("lb-404", "us-east", "123")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestDeleteLoadBalancer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/appliances/lb-123/" {
			t.Errorf("Expected path /appliances/lb-123/, got %s", r.URL.Path)
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
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	err := client.DeleteLoadBalancer("lb-123", "us-east", "123")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestDeleteLoadBalancerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "server error"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	err := client.DeleteLoadBalancer("lb-123", "us-east", "123")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestUpdateLoadBalancerAction(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/appliances/load-balancers/lb-123/actions/" {
			t.Errorf("Expected path /appliances/load-balancers/lb-123/actions/, got %s", r.URL.Path)
		}

		body, _ := ioutil.ReadAll(r.Body)
		var data map[string]interface{}
		json.Unmarshal(body, &data)

		if data["action"] == nil {
			t.Error("Expected action in request body")
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	data := map[string]interface{}{
		"action": "start",
	}

	err := client.UpdateLoadBalancerAction(data, "lb-123", "us-east", "123")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestUpdateLoadBalancerActionError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid action"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	data := map[string]interface{}{
		"action": "invalid",
	}

	err := client.UpdateLoadBalancerAction(data, "lb-123", "us-east", "123")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestIPV6LoadBalancerAction(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/appliances/load-balancers/lb-123/ipv6/" {
			t.Errorf("Expected path /appliances/load-balancers/lb-123/ipv6/, got %s", r.URL.Path)
		}

		body, _ := ioutil.ReadAll(r.Body)
		var data map[string]interface{}
		json.Unmarshal(body, &data)

		if data["enable_ipv6"] == nil {
			t.Error("Expected enable_ipv6 in request body")
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	data := map[string]interface{}{
		"enable_ipv6": true,
	}

	err := client.IPV6LoadBalancerAction(data, "lb-123", "us-east", "123")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestIPV6LoadBalancerActionError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "server error"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	data := map[string]interface{}{
		"enable_ipv6": true,
	}

	err := client.IPV6LoadBalancerAction(data, "lb-123", "us-east", "123")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestLoadBalancerBackendUpdate(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code":    200,
		"message": "Load balancer updated successfully",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		if r.URL.Path != "/appliances/load-balancers/lb-123/" {
			t.Errorf("Expected path /appliances/load-balancers/lb-123/, got %s", r.URL.Path)
		}

		body, _ := ioutil.ReadAll(r.Body)
		var lbCreate models.LoadBalancerCreate
		json.Unmarshal(body, &lbCreate)

		if lbCreate.LbName == "" {
			t.Error("Expected LbName in request body")
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	lbCreate := &models.LoadBalancerCreate{
		LbName:   "updated-lb",
		Location: "us-east",
	}

	result, err := client.LoadBalancerBackendUpdate(lbCreate, "lb-123", "us-east", "123")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result["message"] != mockResponse["message"] {
		t.Errorf("Expected message %s, got %s", mockResponse["message"], result["message"])
	}
}

func TestLoadBalancerBackendUpdateError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid request"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	lbCreate := &models.LoadBalancerCreate{
		LbName:   "updated-lb",
		Location: "us-east",
	}

	result, err := client.LoadBalancerBackendUpdate(lbCreate, "lb-123", "us-east", "123")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestNewLoadBalancerWithRemoveExtraKeys(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code":    200,
		"message": "Load balancer created successfully",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		var data map[string]interface{}
		json.Unmarshal(body, &data)

		// Verify that enable_eos_logger is removed if access_key is empty
		if eosLogger, ok := data["enable_eos_logger"].(map[string]interface{}); ok {
			if accessKey, ok := eosLogger["access_key"].(string); ok && accessKey == "" {
				t.Error("Expected enable_eos_logger to be removed when access_key is empty")
			}
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	lbCreate := &models.LoadBalancerCreate{
		LbName:   "test-lb",
		Location: "us-east",
	}

	result, err := client.NewLoadBalancer(lbCreate, "123")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestNewLoadBalancerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid request"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	lbCreate := &models.LoadBalancerCreate{
		LbName:   "test-lb",
		Location: "us-east",
	}

	result, err := client.NewLoadBalancer(lbCreate, "123")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}
