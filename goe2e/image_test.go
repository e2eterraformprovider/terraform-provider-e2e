package goe2e

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetOSCategories(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/images/os-category")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "test-project")
		testQueryParam(t, r, "location", "test-location")
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Success",
			"data": {
				"category_list": [
					{
						"OS": "CentOS",
						"version": [
							{
								"number_of_domains": null,
								"os": "CentOS",
								"version": "7",
								"sub_category": "CentOS",
								"software_version": ""
							},
							{
								"version": "Stream"
							}
						],
						"category": ["Linux Virtual Node"]
					}
				]
			},
			"errors": {}
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	ctx := context.Background()
	categories, _, err := client.Images.GetOSCategories(ctx)
	assertNoError(t, err)
	assertNotNil(t, categories, "Expected categories")
	if categories == nil {
		t.Fatal("categories is nil, cannot proceed with assertions")
	}
	if len(categories.CategoryList) != 1 {
		t.Errorf("Expected 1 category, got %d", len(categories.CategoryList))
	}
	if categories.CategoryList[0].OS != "CentOS" {
		t.Errorf("Expected OS 'CentOS', got %s", categories.CategoryList[0].OS)
	}
	if len(categories.CategoryList[0].Version) != 2 {
		t.Errorf("Expected 2 versions, got %d", len(categories.CategoryList[0].Version))
	}
}

func TestGetImagePlans(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testURLPath(t, r, "/images")
		testQueryParam(t, r, "category", "CentOS")
		testQueryParam(t, r, "os", "CentOS")
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Success",
			"data": [
				{
					"name": "C3.8GB",
					"plan": "C3-4vCPU-8RAM-100DISK-C3.96GB - CentOS-Stream-Delhi",
					"image": "CentOS-Stream-Distro",
					"os": {
						"name": "CentOS",
						"version": "Stream",
						"image": "CentOS-Stream-Distro",
						"category": "CentOS"
					},
					"location": "Delhi",
					"specs": {
						"id": "1088",
						"sku_name": "C3.8GB",
						"ram": "8.00",
						"cpu": 4,
						"disk_space": 100,
						"price_per_month": 2263,
						"price_per_hour": 3.1,
						"series": "C3",
						"minimum_billing_amount": 0,
						"committed_sku": [],
						"family": "CPU Intensive 3rd Generation"
					},
					"cpu_type": "vCPU",
					"gpu_card_details": {},
					"node_description": "",
					"installed_application_version": {},
					"can_support_bitninja": {
						"show_bitninja": true,
						"bitninja_cost": 760
					},
					"bitninja_discount_percentage": 0,
					"available_inventory_status": true,
					"currency": "INR",
					"is_blockstorage_attachable": true
				}
			],
			"errors": {}
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	ctx := context.Background()
	plans, _, err := client.Images.GetImagePlans(ctx, &ImagePlansRequest{
		Category: "CentOS",
		OS:       "CentOS",
	})
	assertNoError(t, err)
	assertNotNil(t, plans, "Expected plans")
	if len(plans) != 1 {
		t.Errorf("Expected 1 plan, got %d", len(plans))
	}
	if plans[0].Name != "C3.8GB" {
		t.Errorf("Expected plan name 'C3.8GB', got %s", plans[0].Name)
	}
	if plans[0].Specs.CPU != 4 {
		t.Errorf("Expected CPU 4, got %d", plans[0].Specs.CPU)
	}
}

func TestGetImagePlansNilRequest(t *testing.T) {
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	ctx := context.Background()
	_, _, err := client.Images.GetImagePlans(ctx, nil)
	assertError(t, err, "")
}

func TestImportImage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPost)
		testURLPath(t, r, "/images/import-image")
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Success",
			"data": {},
			"errors": {}
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	ctx := context.Background()
	result, _, err := client.Images.ImportImage(ctx, &ImportImageRequest{
		ImageName: "TestImage",
		PublicURL: "https://bucket.objectstore.e2enetworks.net/test.img",
		Location:  "Delhi",
		OS:        "CentOS",
	})
	assertNoError(t, err)
	assertNotNil(t, result, "Expected result")
}
func TestImportImageMissingImageName(t *testing.T) {
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	ctx := context.Background()
	_, _, err := client.Images.ImportImage(ctx, &ImportImageRequest{
		ImageName: "",
	})
	assertError(t, err, "")
}

func TestImportImageMissingPublicURL(t *testing.T) {
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	ctx := context.Background()
	_, _, err := client.Images.ImportImage(ctx, &ImportImageRequest{
		PublicURL: "",
	})
	assertError(t, err, "")
}

func TestImportImageMissingLocation(t *testing.T) {
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	ctx := context.Background()
	_, _, err := client.Images.ImportImage(ctx, &ImportImageRequest{
		Location: "",
	})
	assertError(t, err, "")
}

func TestImportImageMissingOS(t *testing.T) {
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	ctx := context.Background()
	_, _, err := client.Images.ImportImage(ctx, &ImportImageRequest{
		OS: "",
	})
	assertError(t, err, "")
}

func TestGetWindowsImagePermission(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testURLPath(t, r, "/images/window-image-permissions")
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Success",
			"data": {
				"is_windows_allowed": false
			},
			"errors": {}
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	ctx := context.Background()
	permission, _, err := client.Images.GetWindowsImagePermission(ctx)
	assertNoError(t, err)
	assertNotNil(t, permission, "Expected permission")
	if permission.IsWindowsAllowed != false {
		t.Errorf("Expected IsWindowsAllowed false, got %v", permission.IsWindowsAllowed)
	}
}

func TestGetSavedImages(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testURLPath(t, r, "/images/saved-images")
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Success",
			"data": [
				{
					"template_id": 28894,
					"is_windows": false,
					"vm_info": [],
					"image_type": "E2E Image",
					"os_distribution": "DEL-TestingnnooDE_1726058019",
					"name": "TestingnnooDE_1726058019_1726058019",
					"image_id": "24440",
					"distro": "CentOS-Stream",
					"sku_type": "C3VPS",
					"image_state": "Ready",
					"running_vms": "0",
					"cloning_ops": "0",
					"image_size": "95.368 GB",
					"creation_time": "11-09-2024 18:03:37"
				}
			]
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	ctx := context.Background()
	images, _, err := client.Images.GetSavedImages(ctx)
	assertNoError(t, err)
	assertNotNil(t, images, "Expected images")
	if len(images) != 1 {
		t.Errorf("Expected 1 image, got %d", len(images))
	}
	if images[0].TemplateID != 28894 {
		t.Errorf("Expected template_id 28894, got %d", images[0].TemplateID)
	}
	if images[0].Name != "TestingnnooDE_1726058019_1726058019" {
		t.Errorf("Expected name 'TestingnnooDE_1726058019_1726058019', got %s", images[0].Name)
	}
	if images[0].ImageState != "Ready" {
		t.Errorf("Expected image_state 'Ready', got %s", images[0].ImageState)
	}
}

func TestRenameImage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/images/24440")
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Success",
			"data": {
				"status": true,
				"message": "Image name changed successfully"
			}
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	ctx := context.Background()
	result, _, err := client.Images.RenameImage(ctx, "24440", &RenameImageRequest{
		ActionType: "rename",
		Name:       "NewImageName",
	})
	assertNoError(t, err)
	if result.Status != true {
		t.Errorf("Expected status true, got %v", result.Status)
	}
	if result.Message != "Image name changed successfully" {
		t.Errorf("Expected message 'Image name changed successfully', got %s", result.Message)
	}
}

func TestRenameImageEmptyImageID(t *testing.T) {
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	ctx := context.Background()
	_, _, err := client.Images.RenameImage(ctx, "", &RenameImageRequest{
		Name:       "NewName",
		ActionType: "rename",
	})
	assertError(t, err, "")
}

func TestRenameImageNilRequest(t *testing.T) {
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	ctx := context.Background()
	_, _, err := client.Images.RenameImage(ctx, "24440", nil)
	assertError(t, err, "")
}

func TestRenameImageEmptyActionType(t *testing.T) {
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	ctx := context.Background()
	_, _, err := client.Images.RenameImage(ctx, "24440", &RenameImageRequest{
		ActionType: "",
		Name:       "NewName",
	})
	assertError(t, err, "")
}

func TestRenameImageEmptyName(t *testing.T) {
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	ctx := context.Background()
	_, _, err := client.Images.RenameImage(ctx, "24440", &RenameImageRequest{
		ActionType: "rename",
		Name:       "",
	})
	assertError(t, err, "")
}

func TestDeleteImage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodDelete)
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Success",
			"data": {
				"message": "Image deleted successfully"
			}
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	ctx := context.Background()
	result, _, err := client.Images.DeleteImage(ctx, "24440")
	assertNoError(t, err)
	if result.Message != "Image deleted successfully" {
		t.Errorf("Expected message 'Image deleted successfully', got %s", result.Message)
	}
}

func TestDeleteImageEmptyImageID(t *testing.T) {
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	ctx := context.Background()
	_, _, err := client.Images.DeleteImage(ctx, "")
	assertError(t, err, "")
}

func TestImportImageNilRequest(t *testing.T) {
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location")
	ctx := context.Background()
	_, _, err := client.Images.ImportImage(ctx, nil)
	assertError(t, err, "")
}

func TestGetImagePlansWithAllParameters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testURLPath(t, r, "/images")
		testQueryParam(t, r, "category", "CentOS")
		testQueryParam(t, r, "os", "CentOS")
		testQueryParam(t, r, "osversion", "7")
		testQueryParam(t, r, "display_category", "Linux Virtual Node")
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Success",
			"data": [
				{
					"name": "C3.8GB",
					"plan": "C3-4vCPU-8RAM-100DISK-C3.96GB - CentOS-7-Delhi",
					"image": "CentOS-7-Distro",
					"os": {
						"name": "CentOS",
						"version": "7",
						"image": "CentOS-7-Distro",
						"category": "CentOS"
					},
					"location": "Delhi",
					"specs": {
						"id": "1088",
						"sku_name": "C3.8GB",
						"ram": "8.00",
						"cpu": 4,
						"disk_space": 100,
						"price_per_month": 2263,
						"price_per_hour": 3.1,
						"series": "C3",
						"minimum_billing_amount": 0,
						"family": "CPU Intensive 3rd Generation"
					},
					"cpu_type": "vCPU"
				}
			],
			"errors": {}
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	ctx := context.Background()
	plans, _, err := client.Images.GetImagePlans(ctx, &ImagePlansRequest{
		Category:        "CentOS",
		OS:              "CentOS",
		OSVersion:       "7",
		DisplayCategory: "Linux Virtual Node",
	})
	assertNoError(t, err)
	if plans[0].OS.Version != "7" {
		t.Errorf("Expected OS version '7', got %s", plans[0].OS.Version)
	}
}

func TestGetImagePlansEmptyRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Success",
			"data": [],
			"errors": {}
		}`)
	}))
	defer server.Close()
	client, _ := NewClient("test-api-key", "test-auth-token", "test-project", "test-location", SetBaseURL(server.URL+"/"), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	ctx := context.Background()
	plans, _, err := client.Images.GetImagePlans(ctx, &ImagePlansRequest{})
	assertNoError(t, err)
	assertNotNil(t, plans, "Expected plans slice")
}

func TestImageStructMarshaling(t *testing.T) {
	planJSON := `{
		"name": "C3.8GB",
		"plan": "C3-4vCPU-8RAM-100DISK-C3.96GB - CentOS-Stream-Delhi",
		"image": "CentOS-Stream-Distro",
		"os": {
			"name": "CentOS",
			"version": "Stream",
			"image": "CentOS-Stream-Distro",
			"category": "CentOS"
		},
		"location": "Delhi",
		"specs": {
			"id": "1088",
			"sku_name": "C3.8GB",
			"ram": "8.00",
			"cpu": 4,
			"disk_space": 100,
			"price_per_month": 2263,
			"price_per_hour": 3.1,
			"series": "C3",
			"minimum_billing_amount": 0,
			"family": "CPU Intensive 3rd Generation"
		},
		"cpu_type": "vCPU",
		"gpu_card_details": {},
		"node_description": "",
		"installed_application_version": {},
		"can_support_bitninja": {
			"show_bitninja": true,
			"bitninja_cost": 760
		},
		"bitninja_discount_percentage": 0,
		"available_inventory_status": true,
		"currency": "INR",
		"is_blockstorage_attachable": true
	}`
	var plan ImagePlan
	err := json.Unmarshal([]byte(planJSON), &plan)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	if plan.Name != "C3.8GB" {
		t.Errorf("Expected name 'C3.8GB', got %s", plan.Name)
	}
	if plan.Specs.CPU != 4 {
		t.Errorf("Expected CPU 4, got %d", plan.Specs.CPU)
	}
	if plan.CanSupportBitNinja.ShowBitNinja != true {
		t.Errorf("Expected ShowBitNinja true, got %v", plan.CanSupportBitNinja.ShowBitNinja)
	}
}

func TestSavedImageStructMarshaling(t *testing.T) {
	imageJSON := `{
		"template_id": 28894,
		"is_windows": false,
		"vm_info": [],
		"image_type": "E2E Image",
		"os_distribution": "DEL-TestingnnooDE_1726058019",
		"name": "TestingnnooDE_1726058019_1726058019",
		"image_id": "24440",
		"distro": "CentOS-Stream",
		"sku_type": "C3VPS",
		"image_state": "Ready",
		"running_vms": "0",
		"cloning_ops": "0",
		"image_size": "95.368 GB",
		"creation_time": "11-09-2024 18:03:37"
	}`
	var image SavedImage
	err := json.Unmarshal([]byte(imageJSON), &image)
	assertNoError(t, err)
	if image.TemplateID != 28894 {
		t.Errorf("Expected template_id 28894, got %d", image.TemplateID)
	}
	if image.ImageState != "Ready" {
		t.Errorf("Expected image_state 'Ready', got %s", image.ImageState)
	}
	if image.ImageSize != "95.368 GB" {
		t.Errorf("Expected image_size '95.368 GB', got %s", image.ImageSize)
	}
}
