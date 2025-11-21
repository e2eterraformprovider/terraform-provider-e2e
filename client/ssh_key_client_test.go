package client

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
)

func TestAddSshKey(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/ssh_keys/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testURLPath(t, r, "/ssh_keys/")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "123")
		testQueryParam(t, r, "location", "us-east")

		body, _ := io.ReadAll(r.Body)
		var addSshKey models.AddSshKey
		json.Unmarshal(body, &addSshKey)

		if addSshKey.Label == "" {
			t.Error("Expected Label in request body")
		}
		if addSshKey.SshKey == "" {
			t.Error("Expected SshKey in request body")
		}

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "SSH key added successfully",
			"data": {
				"pk": 123,
				"label": "test-key"
			}
		}`)
	})

	addSshKey := models.AddSshKey{
		Label:    "test-key",
		SshKey:   "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQ...",
		Location: "us-east",
	}

	result, err := ts.client.AddSshKey(addSshKey, "123")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result["message"] != "SSH key added successfully" {
		t.Errorf("Expected message 'SSH key added successfully', got %s", result["message"])
	}
}

func TestAddSshKeyError(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/ssh_keys/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusBadRequest, "invalid SSH key")
	})

	addSshKey := models.AddSshKey{
		Label:    "test-key",
		SshKey:   "invalid-key",
		Location: "us-east",
	}

	result, err := ts.client.AddSshKey(addSshKey, "123")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestGetSshKey(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/ssh_keys/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/ssh_keys/")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "123")
		testQueryParam(t, r, "label", "test-key")
		testQueryParam(t, r, "location", "us-east")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"data": [
				{
					"pk": 123,
					"label": "test-key",
					"ssh_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQ..."
				}
			]
		}`)
	})

	result, err := ts.client.GetSshKey("test-key", "123", "us-east")

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
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/ssh_keys/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotFound, "not found")
	})

	result, err := ts.client.GetSshKey("nonexistent-key", "123", "us-east")

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
			ts := setup()
			defer ts.teardown()

			ts.mux.HandleFunc("/delete_ssh_key/123/", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, http.MethodDelete)
				testURLPath(t, r, "/delete_ssh_key/123/")
				testQueryParam(t, r, "apikey", "test-api-key")
				testQueryParam(t, r, "project_id", "456")
				testQueryParam(t, r, "location", "us-east")

				w.WriteHeader(tt.statusCode)
			})

			err := ts.client.DeleteSshKey("123", "456", "us-east")

			if (err != nil) != tt.expectErr {
				t.Errorf("Expected error: %v, got: %v", tt.expectErr, err)
			}
		})
	}
}

func TestDeleteSshKeyError(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/delete_ssh_key/123/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusInternalServerError, "server error")
	})

	err := ts.client.DeleteSshKey("123", "456", "us-east")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestGetSshKeys(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/ssh_keys/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/ssh_keys/")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "location", "us-east")
		testQueryParam(t, r, "project_id", "123")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Success",
			"data": [
				{
					"pk": 123,
					"label": "key-1",
					"ssh_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQ1..."
				},
				{
					"pk": 124,
					"label": "key-2",
					"ssh_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQ2..."
				}
			]
		}`)
	})

	result, err := ts.client.GetSshKeys("us-east", "123")

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
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/ssh_keys/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "server error",
			"data": [],
			"error": []
		}`)
	})

	result, err := ts.client.GetSshKeys("us-east", "123")

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
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/ssh_keys/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/ssh_keys/")
		testQueryParam(t, r, "apikey", "test-api-key")
		testQueryParam(t, r, "project_id", "456")
		testQueryParam(t, r, "location", "us-east")
		testQueryParam(t, r, "pk", "123")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"data": [
				{
					"pk": 123,
					"label": "test-key",
					"ssh_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQ..."
				}
			],
			"error": []
		}`)
	})

	result, err := ts.client.GetSshKeyByPk("123", "456", "us-east")

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
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/ssh_keys/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"data": [],
			"error": []
		}`)
	})

	result, err := ts.client.GetSshKeyByPk("999", "456", "us-east")

	if err == nil {
		t.Fatal("Expected error for not found SSH key, got nil")
	}

	if result != nil {
		t.Error("Expected nil result when SSH key not found")
	}
}

func TestGetSshKeyByPk404Status(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/ssh_keys/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotFound, "not found")
	})

	result, err := ts.client.GetSshKeyByPk("123", "456", "us-east")

	if err == nil {
		t.Fatal("Expected error for 404 status, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on 404")
	}
}

func TestGetSshKeyByPkNon200Status(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/ssh_keys/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusInternalServerError, "server error")
	})

	result, err := ts.client.GetSshKeyByPk("123", "456", "us-east")

	if err == nil {
		t.Fatal("Expected error for non-200 status, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestGetSshKeyByPkMismatchedPk(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/ssh_keys/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"data": [
				{
					"pk": 456,
					"label": "different-key",
					"ssh_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQ..."
				}
			],
			"error": []
		}`)
	})

	result, err := ts.client.GetSshKeyByPk("123", "789", "us-east")

	if err == nil {
		t.Fatal("Expected error for mismatched pk, got nil")
	}

	if result != nil {
		t.Error("Expected nil result when pk doesn't match")
	}
}

func TestAddSshKeyWithHeaders(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/ssh_keys/", func(w http.ResponseWriter, r *http.Request) {
		testHeader(t, r, "Authorization", "Bearer test-auth-token")
		testHeader(t, r, "Content-Type", "application/json")
		testHeader(t, r, "User-Agent", "terraform-e2e")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Success"
		}`)
	})

	addSshKey := models.AddSshKey{
		Label:    "test-key",
		SshKey:   "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQ...",
		Location: "us-east",
	}

	_, err := ts.client.AddSshKey(addSshKey, "123")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestGetSshKeysWithHeaders(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/ssh_keys/", func(w http.ResponseWriter, r *http.Request) {
		testHeader(t, r, "Authorization", "Bearer test-auth-token")
		testHeader(t, r, "Content-Type", "application/json")
		testHeader(t, r, "User-Agent", "terraform-e2e")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "Success",
			"data": []
		}`)
	})

	_, err := ts.client.GetSshKeys("us-east", "123")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestDeleteSshKeyWithHeaders(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/delete_ssh_key/123/", func(w http.ResponseWriter, r *http.Request) {
		testHeader(t, r, "Authorization", "Bearer test-auth-token")
		testHeader(t, r, "Content-Type", "application/json")
		testHeader(t, r, "User-Agent", "terraform-e2e")

		w.WriteHeader(http.StatusOK)
	})

	err := ts.client.DeleteSshKey("123", "456", "us-east")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}
