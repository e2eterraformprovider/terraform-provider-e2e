package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func TestAddSshKey(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code":    200,
		"message": "SSH key added successfully",
		"data": map[string]interface{}{
			"pk":    float64(123),
			"label": "test-key",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/ssh_keys/" {
			t.Errorf("Expected path /ssh_keys/, got %s", r.URL.Path)
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
		var addSshKey models.AddSshKey
		json.Unmarshal(body, &addSshKey)

		if addSshKey.Label == "" {
			t.Error("Expected Label in request body")
		}
		if addSshKey.SshKey == "" {
			t.Error("Expected SshKey in request body")
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	addSshKey := models.AddSshKey{
		Label:    "test-key",
		SshKey:   "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQ...",
		Location: "us-east",
	}

	result, err := client.AddSshKey(addSshKey, "123")

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

func TestAddSshKeyError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid SSH key"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	addSshKey := models.AddSshKey{
		Label:    "test-key",
		SshKey:   "invalid-key",
		Location: "us-east",
	}

	result, err := client.AddSshKey(addSshKey, "123")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestGetSshKey(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code": 200,
		"data": []interface{}{
			map[string]interface{}{
				"pk":      float64(123),
				"label":   "test-key",
				"ssh_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQ...",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/ssh_keys/" {
			t.Errorf("Expected path /ssh_keys/, got %s", r.URL.Path)
		}

		query := r.URL.Query()
		if query.Get("apikey") == "" {
			t.Error("Expected apikey parameter")
		}
		if query.Get("project_id") == "" {
			t.Error("Expected project_id parameter")
		}
		if query.Get("label") == "" {
			t.Error("Expected label parameter")
		}
		if query.Get("location") == "" {
			t.Error("Expected location parameter")
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	result, err := client.GetSshKey("test-key", "123", "us-east")

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

func TestGetSshKeyError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "not found"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	result, err := client.GetSshKey("nonexistent-key", "123", "us-east")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestDeleteSshKey(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		expectErr  bool
	}{
		{
			name:       "Delete with 200 status",
			statusCode: http.StatusOK,
			expectErr:  false,
		},
		{
			name:       "Delete with 204 status",
			statusCode: http.StatusNoContent,
			expectErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "DELETE" {
					t.Errorf("Expected DELETE request, got %s", r.Method)
				}

				if r.URL.Path != "/delete_ssh_key/123/" {
					t.Errorf("Expected path /delete_ssh_key/123/, got %s", r.URL.Path)
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

				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := NewClient("test-key", "test-token", server.URL)
			err := client.DeleteSshKey("123", "456", "us-east")

			if (err != nil) != tt.expectErr {
				t.Errorf("Expected error: %v, got: %v", tt.expectErr, err)
			}
		})
	}
}

func TestDeleteSshKeyError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "server error"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	err := client.DeleteSshKey("123", "456", "us-east")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestGetSshKeys(t *testing.T) {
	mockResponse := models.SshKeyResponse{
		Code:    200,
		Message: "Success",
		Data: []models.SshKey{
			{
				Pk:      123,
				Label:   "key-1",
				Ssh_key: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQ1...",
			},
			{
				Pk:      124,
				Label:   "key-2",
				Ssh_key: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQ2...",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/ssh_keys/" {
			t.Errorf("Expected path /ssh_keys/, got %s", r.URL.Path)
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
	result, err := client.GetSshKeys("us-east", "123")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Code != 200 {
		t.Errorf("Expected code 200, got %d", result.Code)
	}

	if len(result.Data) != 2 {
		t.Errorf("Expected 2 SSH keys, got %d", len(result.Data))
	}

	if result.Data[0].Label != "key-1" {
		t.Errorf("Expected first key label key-1, got %s", result.Data[0].Label)
	}
}

func TestGetSshKeysError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"code": 500, "message": "server error", "data": [], "error": []}`))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	result, err := client.GetSshKeys("us-east", "123")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Code != 500 {
		t.Errorf("Expected code 500, got %d", result.Code)
	}
}

func TestGetSshKeyByPk(t *testing.T) {
	mockResponse := struct {
		Code  int             `json:"code"`
		Data  []models.SshKey `json:"data"`
		Error []interface{}   `json:"error"`
	}{
		Code: 200,
		Data: []models.SshKey{
			{
				Pk:     123,
				Label:  "test-key",
				Ssh_key: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQ...",
			},
		},
		Error: []interface{}{},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/ssh_keys/" {
			t.Errorf("Expected path /ssh_keys/, got %s", r.URL.Path)
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
		if query.Get("pk") == "" {
			t.Error("Expected pk parameter")
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	result, err := client.GetSshKeyByPk("123", "456", "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Pk != 123 {
		t.Errorf("Expected pk 123, got %d", result.Pk)
	}

	if result.Label != "test-key" {
		t.Errorf("Expected label test-key, got %s", result.Label)
	}
}

func TestGetSshKeyByPkNotFound(t *testing.T) {
	mockResponse := struct {
		Code  int             `json:"code"`
		Data  []models.SshKey `json:"data"`
		Error []interface{}   `json:"error"`
	}{
		Code:  200,
		Data:  []models.SshKey{},
		Error: []interface{}{},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	result, err := client.GetSshKeyByPk("999", "456", "us-east")

	if err == nil {
		t.Fatal("Expected error for not found SSH key, got nil")
	}

	if result != nil {
		t.Error("Expected nil result when SSH key not found")
	}
}

func TestGetSshKeyByPk404Status(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "not found"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	result, err := client.GetSshKeyByPk("123", "456", "us-east")

	if err == nil {
		t.Fatal("Expected error for 404 status, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on 404")
	}
}

func TestGetSshKeyByPkNon200Status(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "server error"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	result, err := client.GetSshKeyByPk("123", "456", "us-east")

	if err == nil {
		t.Fatal("Expected error for non-200 status, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestGetSshKeyByPkMismatchedPk(t *testing.T) {
	mockResponse := struct {
		Code  int             `json:"code"`
		Data  []models.SshKey `json:"data"`
		Error []interface{}   `json:"error"`
	}{
		Code: 200,
		Data: []models.SshKey{
			{
				Pk:     456,
				Label:  "different-key",
				Ssh_key: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQ...",
			},
		},
		Error: []interface{}{},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	result, err := client.GetSshKeyByPk("123", "789", "us-east")

	if err == nil {
		t.Fatal("Expected error for mismatched pk, got nil")
	}

	if result != nil {
		t.Error("Expected nil result when pk doesn't match")
	}
}

func TestAddSshKeyWithHeaders(t *testing.T) {
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
			"code":    200,
			"message": "Success",
		})
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	addSshKey := models.AddSshKey{
		Label:    "test-key",
		SshKey:   "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQ...",
		Location: "us-east",
	}

	_, err := client.AddSshKey(addSshKey, "123")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestGetSshKeysWithHeaders(t *testing.T) {
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
		json.NewEncoder(w).Encode(models.SshKeyResponse{
			Code:    200,
			Message: "Success",
			Data:    []models.SshKey{},
		})
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	_, err := client.GetSshKeys("us-east", "123")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestDeleteSshKeyWithHeaders(t *testing.T) {
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
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	err := client.DeleteSshKey("123", "456", "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}
