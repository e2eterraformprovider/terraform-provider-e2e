package client

import (
	"net/http"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func TestCreateScalerGroup(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/scaler/scalegroups", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "Authorization", "Bearer test-auth-token")
		testQueryParam(t, r, "apikey", "test-api-key")

		writeJSON(w, http.StatusCreated, `{
			"code": 201,
			"message": "Scaler Group created successfully",
			"data": {
				"id": "sg-123",
				"name": "test-scaler-group"
			}
		}`)
	})

	req := &models.CreateScalerGroupRequest{
		Name:    "test-scaler-group",
		Desired: "2",
	}

	result, err := ts.client.CreateScalerGroup(req, "test-project", "test-location")

	if err != nil {
		t.Fatalf("CreateScalerGroup returned error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.ID != "sg-123" {
		t.Errorf("Expected ID sg-123, got %s", result.ID)
	}

	if result.Name != "test-scaler-group" {
		t.Errorf("Expected Name test-scaler-group, got %s", result.Name)
	}
}

func TestGetScalerGroup(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/scaler/scalegroups/sg-123/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": {
				"name": "test-scaler-group",
				"desired": 3,
				"running": 3
			}
		}`)
	})

	result, err := ts.client.GetScalerGroup("sg-123", "test-project", "test-location")

	if err != nil {
		t.Fatalf("GetScalerGroup returned error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Name != "test-scaler-group" {
		t.Errorf("Name = %s, expected test-scaler-group", result.Name)
	}

	if result.Desired != 3 {
		t.Errorf("Desired = %d, expected 3", result.Desired)
	}
}

func TestDeleteScalerGroup(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/scaler/scalegroups/sg-123/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Scaler Group deleted successfully"
		}`)
	})

	err := ts.client.DeleteScalerGroup("sg-123", "test-project", "test-location")

	if err != nil {
		t.Fatalf("DeleteScalerGroup returned error: %v", err)
	}
}

func TestGetSavedImageByName(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/images/saved-images/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": [
				{
					"image_id": "img-123",
					"name": "test-image",
					"template_id": 101,
					"distro": "ubuntu"
				},
				{
					"image_id": "img-456",
					"name": "another-image",
					"template_id": 102,
					"distro": "centos"
				}
			]
		}`)
	})

	result, err := ts.client.GetSavedImageByName("test-image", "test-project", "test-location")

	if err != nil {
		t.Fatalf("GetSavedImageByName returned error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.ImageID != "img-123" {
		t.Errorf("ImageID = %s, expected img-123", result.ImageID)
	}

	if result.Name != "test-image" {
		t.Errorf("Name = %s, expected test-image", result.Name)
	}
}

func TestGetSavedImageByName_NotFound(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/images/saved-images/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": []
		}`)
	})

	result, err := ts.client.GetSavedImageByName("nonexistent-image", "test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error for non-existent image, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %+v", result)
	}

	testErrorContains(t, err, "no saved image found")
}

func TestGetDefaultSecurityGroupID(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/security_group/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": [
				{
					"id": 1,
					"is_default": true
				},
				{
					"id": 2,
					"is_default": false
				}
			]
		}`)
	})

	result, err := ts.client.GetDefaultSecurityGroupID("test-project", "test-location")

	if err != nil {
		t.Fatalf("GetDefaultSecurityGroupID returned error: %v", err)
	}

	if result != 1 {
		t.Errorf("Expected ID 1, got %d", result)
	}
}

func TestGetDefaultSecurityGroupID_NotFound(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/security_group/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": [
				{
					"id": 1,
					"is_default": false
				}
			]
		}`)
	})

	result, err := ts.client.GetDefaultSecurityGroupID("test-project", "test-location")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != 0 {
		t.Errorf("Expected 0, got %d", result)
	}
}

func TestGetPlanDetailsFromPlanName(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/images/upgradeimage/101/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"data": {
				"data": [
					{
						"name": "small",
						"plan": "plan-small",
						"specs": {
							"id": "plan-id-123"
						}
					}
				]
			}
		}`)
	})

	planID, slugName, err := ts.client.GetPlanDetailsFromPlanName(101, "small", "test-project", "test-location")

	if err != nil {
		t.Fatalf("GetPlanDetailsFromPlanName returned error: %v", err)
	}

	if planID != "plan-id-123" {
		t.Errorf("Expected planID plan-id-123, got %s", planID)
	}

	if slugName != "plan-small" {
		t.Errorf("Expected slugName plan-small, got %s", slugName)
	}
}

func TestUpdateScalerGroup(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/scaler/scalegroups/update/sg-123/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "updated"
		}`)
	})

	req := &models.UpdateScalerGroupRequest{
		Name: "updated-name",
	}

	err := ts.client.UpdateScalerGroup("sg-123", req, "test-project", "test-location")

	if err != nil {
		t.Fatalf("UpdateScalerGroup returned error: %v", err)
	}
}

func TestUpdateDesiredNodeCount(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/scaler/scalegroups/123/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)

		w.WriteHeader(http.StatusNoContent)
	})

	err := ts.client.UpdateDesiredNodeCount(123, 5, "test-project", "test-location")

	if err != nil {
		t.Fatalf("UpdateDesiredNodeCount returned error: %v", err)
	}
}

func TestUpdateScalerGroupStatus(t *testing.T) {
	tests := []struct {
		name         string
		status       string
		expectedPath string
		wantError    bool
	}{
		{
			name:         "Stop scaler group",
			status:       "Stopped",
			expectedPath: "/scaler/scalegroups/123/stop/",
			wantError:    false,
		},
		{
			name:         "Start scaler group",
			status:       "Running",
			expectedPath: "/scaler/scalegroups/123/start/",
			wantError:    false,
		},
		{
			name:         "Invalid status",
			status:       "Invalid",
			expectedPath: "",
			wantError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := setup()
			defer ts.teardown()

			if !tt.wantError {
				ts.mux.HandleFunc(tt.expectedPath, func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, http.MethodPut)
					writeJSON(w, http.StatusOK, `{"code": 200, "message": "status updated"}`)
				})
			}

			err := ts.client.UpdateScalerGroupStatus(123, tt.status, "test-project", "test-location")

			if (err != nil) != tt.wantError {
				t.Errorf("UpdateScalerGroupStatus() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestGetVpcDetailsByName(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/vpc/list/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		writeJSON(w, http.StatusOK, `{
			"data": [
				{
					"name": "test-vpc",
					"network_id": 100,
					"ipv4_cidr": "10.0.0.0/24"
				},
				{
					"name": "another-vpc",
					"network_id": 101,
					"ipv4_cidr": "10.0.1.0/24"
				}
			]
		}`)
	})

	result, err := ts.client.GetVpcDetailsByName("test-project", "test-location", "test-vpc")

	if err != nil {
		t.Fatalf("GetVpcDetailsByName returned error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Name != "test-vpc" {
		t.Errorf("Expected Name test-vpc, got %s", result.Name)
	}

	if result.NetworkID != 100 {
		t.Errorf("Expected NetworkID 100, got %d", result.NetworkID)
	}
}

func TestAttachVPCToScalerGroup(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/scaler/scalegroups/sg-123/vpc/action/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)

		w.WriteHeader(http.StatusOK)
	})

	vpcs := []models.VPCDetail{
		{
			Name:      "test-vpc",
			NetworkID: 100,
			IPv4CIDR:  "10.0.0.0/24",
		},
	}

	err := ts.client.AttachVPCToScalerGroup("sg-123", vpcs, "test-project", "test-location")

	if err != nil {
		t.Fatalf("AttachVPCToScalerGroup returned error: %v", err)
	}
}

func TestDetachVPCFromScalerGroup(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/scaler/scalegroups/sg-123/vpc/action/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		testQueryParam(t, r, "vpc_id", "vpc-456")

		w.WriteHeader(http.StatusOK)
	})

	err := ts.client.DetachVPCFromScalerGroup("sg-123", "vpc-456", "test-project", "test-location")

	if err != nil {
		t.Fatalf("DetachVPCFromScalerGroup returned error: %v", err)
	}
}

func TestGetPublicIPStatus(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/scaler/scalegroups/sg-123/public_ip/action/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": {
				"is_public_ip_required": true
			}
		}`)
	})

	result, err := ts.client.GetPublicIPStatus("sg-123", "test-project", "test-location")

	if err != nil {
		t.Fatalf("GetPublicIPStatus returned error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if !result.IsPublicIPRequired {
		t.Error("Expected IsPublicIPRequired to be true, got false")
	}
}

func TestAttachPublicIP(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/scaler/scalegroups/sg-123/public_ip/action/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Public IP attached",
			"data": "192.168.1.1"
		}`)
	})

	result, err := ts.client.AttachPublicIP("sg-123", "test-project", "test-location")

	if err != nil {
		t.Fatalf("AttachPublicIP returned error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Data != "192.168.1.1" {
		t.Errorf("Expected Data 192.168.1.1, got %s", result.Data)
	}
}

func TestDetachPublicIP(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/scaler/scalegroups/sg-123/public_ip/action/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Public IP detached",
			"data": ""
		}`)
	})

	result, err := ts.client.DetachPublicIP("sg-123", "test-project", "test-location")

	if err != nil {
		t.Fatalf("DetachPublicIP returned error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}
}

func TestGetAttachedVPCsForScalerGroup(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/scaler/scalegroups/sg-123/vpc/action/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		writeJSON(w, http.StatusOK, `{
			"data": [
				{
					"name": "vpc-1",
					"network_id": 100
				},
				{
					"name": "vpc-2",
					"network_id": 101
				}
			]
		}`)
	})

	result, err := ts.client.GetAttachedVPCsForScalerGroup("sg-123", "test-project", "test-location")

	if err != nil {
		t.Fatalf("GetAttachedVPCsForScalerGroup returned error: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 VPCs, got %d", len(result))
	}

	if result[0].Name != "vpc-1" {
		t.Errorf("Expected Name vpc-1, got %s", result[0].Name)
	}
}

func TestDetachSecurityGroupFromScalergroup(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/scaler/scalegroups/security_groups/sg-123/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		testQueryParam(t, r, "security_group_id", "456")

		w.WriteHeader(http.StatusOK)
	})

	err := ts.client.DetachSecurityGroupFromScalergroup("sg-123", 456, "test-project", "test-location")

	if err != nil {
		t.Fatalf("DetachSecurityGroupFromScalergroup returned error: %v", err)
	}
}

func TestAddSecurityGroupToScalergroup(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/scaler/scalegroups/security_groups/sg-123/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)

		w.WriteHeader(http.StatusOK)
	})

	err := ts.client.AddSecurityGroupToScalergroup("sg-123", 456, "test-project", "test-location")

	if err != nil {
		t.Fatalf("AddSecurityGroupToScalergroup returned error: %v", err)
	}
}
