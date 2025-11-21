package client

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func TestGetSecurityGroupList(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/security_group/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/security_group/")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "location", "us-east")
		testQueryParam(t, r, "project_id", "123")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"data": [
				{"id": "sg-1", "name": "default"},
				{"id": "sg-2", "name": "custom"}
			]
		}`)
	})

	result, err := ts.client.GetSecurityGroupList("123", "us-east")

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
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/security_group/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusInternalServerError, "server error")
	})

	result, err := ts.client.GetSecurityGroupList("123", "us-east")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestGetSecurityGroup(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/security_group/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"data": [
				{"id": "sg-1", "name": "test-sg"},
				{"id": "sg-2", "name": "other-sg"}
			]
		}`)
	})

	result, err := ts.client.GetSecurityGroup("test-sg", "123", "us-east")

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
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/security_group/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"data": [
				{"id": "sg-1", "name": "other-sg"}
			]
		}`)
	})

	result, err := ts.client.GetSecurityGroup("nonexistent-sg", "123", "us-east")

	if err == nil {
		t.Fatal("Expected error for not found security group, got nil")
	}

	if result != nil {
		t.Error("Expected nil result when security group not found")
	}
}

func TestCreateSecurityGroups(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/security_group/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testURLPath(t, r, "/security_group/")

		body, _ := io.ReadAll(r.Body)
		var payload models.SecurityGroupCreateRequest
		json.Unmarshal(body, &payload)

		if payload.Name == "" {
			t.Error("Expected Name in request body")
		}

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Security group created"
		}`)
	})

	payload := models.SecurityGroupCreateRequest{
		Name:        "test-sg",
		Description: "Test security group",
	}

	err := ts.client.CreateSecurityGroups(payload, "123", "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestCreateSecurityGroupsError(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/security_group/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusBadRequest, "invalid request")
	})

	payload := models.SecurityGroupCreateRequest{
		Name: "test-sg",
	}

	err := ts.client.CreateSecurityGroups(payload, "123", "us-east")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestUpdateSecurityGroups(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/security_group/sg-123/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/security_group/sg-123/")

		body, _ := io.ReadAll(r.Body)
		var payload models.SecurityGroupUpdateRequest
		json.Unmarshal(body, &payload)

		if payload.Name == "" {
			t.Error("Expected Name in request body")
		}

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Security group updated"
		}`)
	})

	payload := models.SecurityGroupUpdateRequest{
		Name:        "updated-sg",
		Description: "Updated security group",
	}

	err := ts.client.UpdateSecurityGroups(payload, "sg-123", "123", "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestUpdateSecurityGroupsError(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/security_group/sg-404/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotFound, "not found")
	})

	payload := models.SecurityGroupUpdateRequest{
		Name: "updated-sg",
	}

	err := ts.client.UpdateSecurityGroups(payload, "sg-404", "123", "us-east")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestMakeDefaultSecurityGroup(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/security_group/sg-123/mark-default/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testURLPath(t, r, "/security_group/sg-123/mark-default/")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Security group marked as default"
		}`)
	})

	err := ts.client.MakeDefaultSecurityGroup("sg-123", "123", "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestMakeDefaultSecurityGroupError(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/security_group/sg-123/mark-default/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusBadRequest, "cannot mark as default")
	})

	err := ts.client.MakeDefaultSecurityGroup("sg-123", "123", "us-east")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestDetachSecurityGroup(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/security_group/123/detach/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testURLPath(t, r, "/security_group/123/detach/")

		body, _ := io.ReadAll(r.Body)
		var payload models.UpdateSecurityGroups
		json.Unmarshal(body, &payload)

		if len(payload.SecurityGroupList) == 0 {
			t.Error("Expected SecurityGroupList in request body")
		}

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Security group detached"
		}`)
	})

	payload := &models.UpdateSecurityGroups{
		SecurityGroupList: []int{1, 2},
	}

	result, err := ts.client.DetachSecurityGroup(payload, 123, "456", "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestDetachSecurityGroupError(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/security_group/123/detach/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusBadRequest, "invalid request")
	})

	payload := &models.UpdateSecurityGroups{
		SecurityGroupList: []int{1},
	}

	result, err := ts.client.DetachSecurityGroup(payload, 123, "456", "us-east")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestAttachSecurityGroup(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/security_group/123/attach/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testURLPath(t, r, "/security_group/123/attach/")

		body, _ := io.ReadAll(r.Body)
		var payload models.UpdateSecurityGroups
		json.Unmarshal(body, &payload)

		if len(payload.SecurityGroupList) == 0 {
			t.Error("Expected SecurityGroupList in request body")
		}

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Security group attached"
		}`)
	})

	payload := &models.UpdateSecurityGroups{
		SecurityGroupList: []int{1, 2},
	}

	result, err := ts.client.AttachSecurityGroup(payload, 123, "456", "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestAttachSecurityGroupError(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/security_group/123/attach/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusBadRequest, "invalid request")
	})

	payload := &models.UpdateSecurityGroups{
		SecurityGroupList: []int{1},
	}

	result, err := ts.client.AttachSecurityGroup(payload, 123, "456", "us-east")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestDeleteSecurityGroup(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/security_group/sg-123/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		testURLPath(t, r, "/security_group/sg-123/")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "location", "us-east")
		testQueryParam(t, r, "project_id", "123")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Security group deleted"
		}`)
	})

	err := ts.client.DeleteSecurityGroup("sg-123", "123", "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestDeleteSecurityGroupError(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/security_group/sg-123/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusConflict, "security group in use")
	})

	err := ts.client.DeleteSecurityGroup("sg-123", "123", "us-east")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestGetSecurityGroupInvalidData(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/security_group/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"data": "invalid-not-an-array"
		}`)
	})

	result, err := ts.client.GetSecurityGroup("test-sg", "123", "us-east")

	if err == nil {
		t.Fatal("Expected error for invalid data format, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestCreateSecurityGroupsWithHeaders(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/security_group/", func(w http.ResponseWriter, r *http.Request) {
		testHeader(t, r, "Authorization", "Bearer test-auth-token")
		testHeader(t, r, "Content-Type", "application/json")
		testHeader(t, r, "User-Agent", "terraform-e2e")

		writeJSON(w, http.StatusOK, `{"code": 200}`)
	})

	payload := models.SecurityGroupCreateRequest{
		Name: "test-sg",
	}

	err := ts.client.CreateSecurityGroups(payload, "123", "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}
