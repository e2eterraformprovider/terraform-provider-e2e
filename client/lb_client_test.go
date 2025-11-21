package client

import (
	"net/http"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func TestAddParamsAndHeader(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	req, _ := http.NewRequest("GET", ts.server.URL+"/test/", nil)
	location := "us-east"
	projectID := "123"

	modifiedReq, err := ts.client.AddParamsAndHeader(req, location, projectID)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	params := modifiedReq.URL.Query()
	if params.Get("apikey") != "test-api-key" {
		t.Errorf("Expected apikey test-api-key, got %s", params.Get("apikey"))
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

	if modifiedReq.Header.Get("Authorization") != "Bearer test-auth-token" {
		t.Errorf("Expected Authorization header Bearer test-auth-token, got %s", modifiedReq.Header.Get("Authorization"))
	}

	if modifiedReq.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", modifiedReq.Header.Get("Content-Type"))
	}

	if modifiedReq.Header.Get("User-Agent") != "terraform-e2e" {
		t.Errorf("Expected User-Agent terraform-e2e, got %s", modifiedReq.Header.Get("User-Agent"))
	}
}

func TestNewLoadBalancer(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/appliances/load-balancers/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testURLPath(t, r, "/appliances/load-balancers/")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "123")
		testQueryParam(t, r, "location", "us-east")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Load balancer created successfully",
			"data": {
				"id": "lb-123",
				"name": "test-lb"
			}
		}`)
	})

	lbCreate := &models.LoadBalancerCreate{
		LbName:   "test-lb",
		Location: "us-east",
	}

	result, err := ts.client.NewLoadBalancer(lbCreate, "123")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result["message"] != "Load balancer created successfully" {
		t.Errorf("Expected message Load balancer created successfully, got %s", result["message"])
	}
}

func TestGetLoadBalancerInfo(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/appliances/lb-123/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/appliances/lb-123/")
		testQueryParam(t, r, "apikey", "test-api-key")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"data": {
				"id": "lb-123",
				"name": "test-lb",
				"state": "ACTIVE"
			}
		}`)
	})

	result, err := ts.client.GetLoadBalancerInfo("lb-123", "us-east", "123")

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
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/appliances/lb-404/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotFound, "not found")
	})

	result, err := ts.client.GetLoadBalancerInfo("lb-404", "us-east", "123")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestDeleteLoadBalancer(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/appliances/lb-123/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		testURLPath(t, r, "/appliances/lb-123/")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "location", "us-east")
		testQueryParam(t, r, "project_id", "123")

		w.WriteHeader(http.StatusOK)
	})

	err := ts.client.DeleteLoadBalancer("lb-123", "us-east", "123")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestDeleteLoadBalancerError(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/appliances/lb-123/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusInternalServerError, "server error")
	})

	err := ts.client.DeleteLoadBalancer("lb-123", "us-east", "123")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestUpdateLoadBalancerAction(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/appliances/load-balancers/lb-123/actions/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/appliances/load-balancers/lb-123/actions/")

		w.WriteHeader(http.StatusOK)
	})

	data := map[string]interface{}{
		"action": "start",
	}

	err := ts.client.UpdateLoadBalancerAction(data, "lb-123", "us-east", "123")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestUpdateLoadBalancerActionError(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/appliances/load-balancers/lb-123/actions/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusBadRequest, "invalid action")
	})

	data := map[string]interface{}{
		"action": "invalid",
	}

	err := ts.client.UpdateLoadBalancerAction(data, "lb-123", "us-east", "123")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestIPV6LoadBalancerAction(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/appliances/load-balancers/lb-123/ipv6/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/appliances/load-balancers/lb-123/ipv6/")

		w.WriteHeader(http.StatusOK)
	})

	data := map[string]interface{}{
		"enable_ipv6": true,
	}

	err := ts.client.IPV6LoadBalancerAction(data, "lb-123", "us-east", "123")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestIPV6LoadBalancerActionError(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/appliances/load-balancers/lb-123/ipv6/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusInternalServerError, "server error")
	})

	data := map[string]interface{}{
		"enable_ipv6": true,
	}

	err := ts.client.IPV6LoadBalancerAction(data, "lb-123", "us-east", "123")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestLoadBalancerBackendUpdate(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/appliances/load-balancers/lb-123/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/appliances/load-balancers/lb-123/")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Load balancer updated successfully"
		}`)
	})

	lbCreate := &models.LoadBalancerCreate{
		LbName:   "updated-lb",
		Location: "us-east",
	}

	result, err := ts.client.LoadBalancerBackendUpdate(lbCreate, "lb-123", "us-east", "123")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result["message"] != "Load balancer updated successfully" {
		t.Errorf("Expected message Load balancer updated successfully, got %s", result["message"])
	}
}

func TestLoadBalancerBackendUpdateError(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/appliances/load-balancers/lb-123/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusBadRequest, "invalid request")
	})

	lbCreate := &models.LoadBalancerCreate{
		LbName:   "updated-lb",
		Location: "us-east",
	}

	result, err := ts.client.LoadBalancerBackendUpdate(lbCreate, "lb-123", "us-east", "123")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestNewLoadBalancerError(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/appliances/load-balancers/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusBadRequest, "invalid request")
	})

	lbCreate := &models.LoadBalancerCreate{
		LbName:   "test-lb",
		Location: "us-east",
	}

	result, err := ts.client.NewLoadBalancer(lbCreate, "123")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}
