package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func TestAddParamsAndHeadersFunc(t *testing.T) {
	req := httptest.NewRequest("GET", "https://api.test.com/test/", nil)
	apiKey := "test-key"
	authToken := "test-token"
	projectID := "123"
	location := "us-east"

	modifiedReq := AddParamsAndHeaders(req, apiKey, authToken, projectID, location)

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

func TestNewSfs(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code":    200,
		"message": "SFS created successfully",
		"data": map[string]interface{}{
			"id":   "sfs-123",
			"name": "test-sfs",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/efs/create/" {
			t.Errorf("Expected path /efs/create/, got %s", r.URL.Path)
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
		var sfsCreate models.SfsCreate
		json.Unmarshal(body, &sfsCreate)

		if sfsCreate.Name == "" {
			t.Error("Expected Name in request body")
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	sfsCreate := &models.SfsCreate{
		Name: "test-sfs",
	}

	result, err := client.NewSfs(sfsCreate, "123", "us-east")

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

func TestNewSfsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid request"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)

	sfsCreate := &models.SfsCreate{
		Name: "test-sfs",
	}

	result, err := client.NewSfs(sfsCreate, "123", "us-east")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestGetSfs(t *testing.T) {
	mockResponse := map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{
			"id":    "sfs-123",
			"name":  "test-sfs",
			"state": "ACTIVE",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/efs/sfs-123/" {
			t.Errorf("Expected path /efs/sfs-123/, got %s", r.URL.Path)
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

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	result, err := client.GetSfs("sfs-123", "123", "us-east")

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

func TestGetSfsNon200Status(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "not found"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	result, err := client.GetSfs("sfs-404", "123", "us-east")

	if err == nil {
		t.Fatal("Expected error for non-200 status, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestDeleteSFs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/efs/delete/sfs-123/" {
			t.Errorf("Expected path /efs/delete/sfs-123/, got %s", r.URL.Path)
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

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	err := client.DeleteSFs("sfs-123", "123", "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestDeleteSFsNon200Status(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "server error"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	err := client.DeleteSFs("sfs-123", "123", "us-east")

	if err == nil {
		t.Fatal("Expected error for non-200 status, got nil")
	}
}

func TestGetSfss(t *testing.T) {
	mockResponse := models.ResponseSfss{
		Code:    200,
		Message: "Success",
		Data: []models.SfssRead{
			{
				ID:   1,
				Name: "sfs-one",
			},
			{
				ID:   2,
				Name: "sfs-two",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/efs/" {
			t.Errorf("Expected path /efs/, got %s", r.URL.Path)
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

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	result, err := client.GetSfss("us-east", "123")

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
		t.Errorf("Expected 2 SFS instances, got %d", len(result.Data))
	}
}

func TestGetSfssNon200Status(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "unauthorized"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	result, err := client.GetSfss("us-east", "123")

	if err == nil {
		t.Fatal("Expected error for non-200 status, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestNewSfsWithHeaders(t *testing.T) {
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

	sfsCreate := &models.SfsCreate{
		Name: "test-sfs",
	}

	_, err := client.NewSfs(sfsCreate, "123", "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestGetSfsWithHeaders(t *testing.T) {
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
			"data": map[string]interface{}{},
		})
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	_, err := client.GetSfs("sfs-123", "123", "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestDeleteSFsWithHeaders(t *testing.T) {
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
	err := client.DeleteSFs("sfs-123", "123", "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestGetSfssWithHeaders(t *testing.T) {
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
		json.NewEncoder(w).Encode(models.ResponseSfss{
			Code:    200,
			Message: "Success",
			Data:    []models.SfssRead{},
		})
	}))
	defer server.Close()

	client := NewClient("test-key", "test-token", server.URL)
	_, err := client.GetSfss("us-east", "123")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}
