package client

import (
	"net/http"
	"testing"
)

func TestGetSoftwareId(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/plans/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/rds/plans/")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": {
				"database_engines": [
					{
						"id": 1,
						"name": "mysql",
						"version": "8.0"
					},
					{
						"id": 2,
						"name": "postgres",
						"version": "14.0"
					}
				]
			}
		}`)
	})

	result, err := ts.client.GetSoftwareId("test-project", "test-location", "mysql", "8.0")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result != 1 {
		t.Errorf("Expected EngineID 1, got %d", result)
	}
}

func TestGetSoftwareId_NotFound(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/plans/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": {
				"database_engines": [
					{
						"id": 1,
						"name": "mysql",
						"version": "8.0"
					}
				]
			}
		}`)
	})

	result, err := ts.client.GetSoftwareId("test-project", "test-location", "postgres", "15.0")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != -1 {
		t.Errorf("Expected -1, got %d", result)
	}
}

func TestGetTemplateId(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/plans/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/rds/plans/")
		testQueryParam(t, r, "software_id", "1")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": {
				"template_plans": [
					{
						"template_id": 100,
						"name": "small"
					},
					{
						"template_id": 101,
						"name": "medium"
					}
				]
			}
		}`)
	})

	result, err := ts.client.GetTemplateId("test-project", "test-location", "small", "1")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result != 100 {
		t.Errorf("Expected PlanTemplateID 100, got %d", result)
	}
}

func TestGetTemplateId_NotFound(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/rds/plans/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": {
				"template_plans": [
					{
						"template_id": 100,
						"name": "small"
					}
				]
			}
		}`)
	})

	result, err := ts.client.GetTemplateId("test-project", "test-location", "large", "1")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != -1 {
		t.Errorf("Expected -1, got %d", result)
	}
}

func TestExpandMariaDBVpcList(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/vpc/100/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": {
				"name": "test-vpc",
				"network_id": 100,
				"ipv4_cidr": "10.0.0.0/24",
				"state": "Active"
			}
		}`)
	})

	vpcIDs := []string{"100"}
	result, err := ts.client.ExpandMariaDBVpcList(vpcIDs, "test-project", "test-location")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 VPC, got %d", len(result))
	}

	if result[0].VPCName != "test-vpc" {
		t.Errorf("Expected VPCName test-vpc, got %s", result[0].VPCName)
	}
}

func TestExpandMariaDBVpcList_InactiveVPC(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/vpc/100/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": {
				"name": "test-vpc",
				"network_id": 100,
				"ipv4_cidr": "10.0.0.0/24",
				"state": "Inactive"
			}
		}`)
	})

	vpcIDs := []string{"100"}
	result, err := ts.client.ExpandMariaDBVpcList(vpcIDs, "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error for inactive VPC, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}
