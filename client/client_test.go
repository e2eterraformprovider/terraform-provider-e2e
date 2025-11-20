package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func TestNewClient(t *testing.T) {
	apiKey := "test-api-key"
	authToken := "test-auth-token"
	apiEndpoint := "https://api.test.com/"

	client := NewClient(apiKey, authToken, apiEndpoint)

	if client.Api_key != apiKey {
		t.Errorf("Expected Api_key %s, got %s", apiKey, client.Api_key)
	}

	if client.Auth_token != authToken {
		t.Errorf("Expected Auth_token %s, got %s", authToken, client.Auth_token)
	}

	if client.Api_endpoint != apiEndpoint {
		t.Errorf("Expected Api_endpoint %s, got %s", apiEndpoint, client.Api_endpoint)
	}

	if client.HttpClient == nil {
		t.Error("Expected HttpClient to be initialized, got nil")
	}
}

func TestAddParamsAndHeaders(t *testing.T) {
	req := httptest.NewRequest("GET", "https://api.test.com/nodes/", nil)
	apiKey := "test-key"
	authToken := "test-token"
	projectID := "123"
	location := "us-east"

	modifiedReq := addParamsAndHeaders(req, apiKey, authToken, projectID, location)

	// Check query parameters
	params := modifiedReq.URL.Query()
	if params.Get("apikey") != apiKey {
		t.Errorf("Expected apikey %s, got %s", apiKey, params.Get("apikey"))
	}

	if params.Get("project_id") != projectID {
		t.Errorf("Expected project_id %s, got %s", projectID, params.Get("project_id"))
	}

	if params.Get("location") != location {
		t.Errorf("Expected location %s, got %s", location, params.Get("location"))
	}

	// Check headers
	if modifiedReq.Header.Get("Authorization") != "Bearer "+authToken {
		t.Errorf("Expected Authorization header Bearer %s, got %s", authToken, modifiedReq.Header.Get("Authorization"))
	}

	if modifiedReq.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", modifiedReq.Header.Get("Content-Type"))
	}

	if modifiedReq.Header.Get("User-Agent") != "terraform-e2e" {
		t.Errorf("Expected User-Agent terraform-e2e, got %s", modifiedReq.Header.Get("User-Agent"))
	}
}

func TestSetBasicHeaders(t *testing.T) {
	req := httptest.NewRequest("GET", "https://api.test.com/nodes/", nil)
	authToken := "test-token"

	SetBasicHeaders(authToken, req)

	if req.Header.Get("Authorization") != "Bearer "+authToken {
		t.Errorf("Expected Authorization header Bearer %s, got %s", authToken, req.Header.Get("Authorization"))
	}

	if req.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", req.Header.Get("Content-Type"))
	}

	if req.Header.Get("User-Agent") != "terraform-e2e" {
		t.Errorf("Expected User-Agent terraform-e2e, got %s", req.Header.Get("User-Agent"))
	}
}

func TestCheckResponseStatus(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		responseBody   string
		expectError    bool
		errorContains  string
	}{
		{
			name:        "Success - 200 OK",
			statusCode:  http.StatusOK,
			responseBody: `{"status": "success"}`,
			expectError: false,
		},
		{
			name:          "Error - 400 Bad Request",
			statusCode:    http.StatusBadRequest,
			responseBody:  `{"error": "invalid request"}`,
			expectError:   true,
			errorContains: "400",
		},
		{
			name:          "Error - 401 Unauthorized",
			statusCode:    http.StatusUnauthorized,
			responseBody:  `{"error": "unauthorized"}`,
			expectError:   true,
			errorContains: "401",
		},
		{
			name:          "Error - 404 Not Found",
			statusCode:    http.StatusNotFound,
			responseBody:  `{"error": "not found"}`,
			expectError:   true,
			errorContains: "404",
		},
		{
			name:          "Error - 500 Internal Server Error",
			statusCode:    http.StatusInternalServerError,
			responseBody:  `{"error": "server error"}`,
			expectError:   true,
			errorContains: "500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock response
			resp := &http.Response{
				StatusCode: tt.statusCode,
				Body:       ioutil.NopCloser(bytes.NewBufferString(tt.responseBody)),
			}

			err := CheckResponseStatus(resp)

			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}

			if tt.expectError && err != nil {
				if tt.errorContains != "" && !bytes.Contains([]byte(err.Error()), []byte(tt.errorContains)) {
					t.Errorf("Expected error to contain %s, got: %s", tt.errorContains, err.Error())
				}
			}
		})
	}
}

func TestCheckResponseCreatedStatus(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		responseBody   string
		expectError    bool
		errorContains  string
	}{
		{
			name:        "Success - 201 Created",
			statusCode:  http.StatusCreated,
			responseBody: `{"status": "created"}`,
			expectError: false,
		},
		{
			name:          "Error - 200 OK (not created)",
			statusCode:    http.StatusOK,
			responseBody:  `{"status": "ok"}`,
			expectError:   true,
			errorContains: "201",
		},
		{
			name:          "Error - 400 Bad Request",
			statusCode:    http.StatusBadRequest,
			responseBody:  `{"error": "invalid request"}`,
			expectError:   true,
			errorContains: "201",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock response
			resp := &http.Response{
				StatusCode: tt.statusCode,
				Body:       ioutil.NopCloser(bytes.NewBufferString(tt.responseBody)),
			}

			err := CheckResponseCreatedStatus(resp)

			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}

			if tt.expectError && err != nil {
				if tt.errorContains != "" && !bytes.Contains([]byte(err.Error()), []byte(tt.errorContains)) {
					t.Errorf("Expected error to contain %s, got: %s", tt.errorContains, err.Error())
				}
			}
		})
	}
}

func TestGetVpcs(t *testing.T) {
	// Create a mock server
	mockResponse := models.VpcsResponse{
		Code:    200,
		Message: "success",
		Data: []models.Vpc{
			{
				Name:       "test-vpc-1",
				Ipv4_cidr:  "10.0.0.0/24",
				Network_id: 1,
				State:      "ACTIVE",
			},
			{
				Name:       "test-vpc-2",
				Ipv4_cidr:  "10.0.1.0/24",
				Network_id: 2,
				State:      "ACTIVE",
			},
		},
		Error: []interface{}{},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		// Verify URL path
		if r.URL.Path != "/vpc/list/" {
			t.Errorf("Expected path /vpc/list/, got %s", r.URL.Path)
		}

		// Verify query parameters
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

		// Verify headers
		if r.Header.Get("Authorization") == "" {
			t.Error("Expected Authorization header")
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Create client with mock server URL
	client := NewClient("test-key", "test-token", server.URL+"/")

	// Call GetVpcs
	result, err := client.GetVpcs("test-location", "test-project")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Code != mockResponse.Code {
		t.Errorf("Expected code %d, got %d", mockResponse.Code, result.Code)
	}

	if len(result.Data) != len(mockResponse.Data) {
		t.Errorf("Expected %d VPCs, got %d", len(mockResponse.Data), len(result.Data))
	}

	// Verify first VPC
	if len(result.Data) > 0 {
		if result.Data[0].Name != mockResponse.Data[0].Name {
			t.Errorf("Expected VPC name %s, got %s", mockResponse.Data[0].Name, result.Data[0].Name)
		}
	}
}

func TestCreateVpc(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code":    201,
		"message": "VPC created successfully",
		"data": map[string]interface{}{
			"network_id": float64(123),
			"name":       "test-vpc",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/vpc/" {
			t.Errorf("Expected path /vpc/, got %s", r.URL.Path)
		}

		// Read and verify request body
		body, _ := ioutil.ReadAll(r.Body)
		var vpcCreate models.VpcCreate
		json.Unmarshal(body, &vpcCreate)

		if vpcCreate.VpcName == "" {
			t.Error("Expected VpcName in request body")
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL+"/")

	vpcCreate := &models.VpcCreate{
		VpcName:  "test-vpc",
		IPv4:     "10.0.0.0/24",
		IsE2EVpc: true,
	}

	result, err := client.CreateVpc("test-location", vpcCreate, "test-project")

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
