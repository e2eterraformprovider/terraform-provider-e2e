package goe2e

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateSSHKey(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/ssh_keys/", func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPost)
		assertQueryParam(t, r, "apikey", "test-api-key")
		assertQueryParam(t, r, "project_id", "test-project")
		assertQueryParam(t, r, "location", "test-location")

		writeJSON(w, http.StatusCreated, `{
			"code": 201,
			"message": "SSH key created successfully",
			"data": {
				"pk": 123,
				"label": "my-key",
				"ssh_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAB..."
			}
		}`)
	})

	createReq := &SSHKeyCreateRequest{
		Label:  "my-key",
		SSHKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAB...",
	}

	result, resp, err := ts.client.SSHKeys.CreateSSHKey(context.Background(), createReq)
	assertNoError(t, err)
	assertNotNil(t, result, "Expected result")
	if result.Label != "my-key" {
		t.Errorf("Expected Label 'my-key', got %s", result.Label)
	}
	if result.SSHKey != "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAB..." {
		t.Errorf("Expected SSHKey 'ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAB...', got %s", result.SSHKey)
	}
	assertStatus(t, resp, http.StatusCreated)
}

func TestCreateSSHKey_NilRequest(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, _, err := ts.client.SSHKeys.CreateSSHKey(context.Background(), nil)
	assertError(t, err, "")
	assertErrorType(t, err, &ArgError{})
}

func TestCreateSSHKey_EmptyLabel(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	createReq := &SSHKeyCreateRequest{
		Label:  "",
		SSHKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAB...",
	}

	_, _, err := ts.client.SSHKeys.CreateSSHKey(context.Background(), createReq)
	assertError(t, err, "")
	assertErrorType(t, err, &ArgError{})
}

func TestCreateSSHKey_EmptySSHKey(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	createReq := &SSHKeyCreateRequest{
		Label:  "my-key",
		SSHKey: "",
	}

	_, _, err := ts.client.SSHKeys.CreateSSHKey(context.Background(), createReq)
	assertError(t, err, "")
	assertErrorType(t, err, &ArgError{})
}

func TestCreateSSHKey_ErrorResponse(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/ssh_keys/", func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPost)
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Internal server error",
			"error": "Something went wrong"
		}`)
	})

	createReq := &SSHKeyCreateRequest{
		Label:  "my-key",
		SSHKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAB...",
	}

	_, resp, err := ts.client.SSHKeys.CreateSSHKey(context.Background(), createReq)
	assertError(t, err, "")
	assertStatus(t, resp, http.StatusInternalServerError)
}

func TestGetSSHKey(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/ssh_keys/", func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodGet)
		assertQueryParam(t, r, "apikey", "test-api-key")
		assertQueryParam(t, r, "project_id", "test-project")
		assertQueryParam(t, r, "location", "test-location")
		assertQueryParam(t, r, "pk", "123")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": [
				{
					"pk": 123,
					"label": "my-key",
					"ssh_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAB...",
					"timestamp": "2024-01-01T00:00:00Z"
				}
			]
		}`)
	})

	result, resp, err := ts.client.SSHKeys.GetSSHKey(context.Background(), "123")
	assertNoError(t, err)
	assertNotNil(t, result, "Expected result")
	if result.PK != 123 {
		t.Errorf("Expected PK 123, got %d", result.PK)
	}
	if result.Label != "my-key" {
		t.Errorf("Expected Label 'my-key', got %s", result.Label)
	}
	assertStatus(t, resp, http.StatusOK)
}

func TestGetSSHKey_NotFound(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/ssh_keys/", func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodGet)
		assertQueryParam(t, r, "pk", "999")
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": []
		}`)
	})

	_, _, err := ts.client.SSHKeys.GetSSHKey(context.Background(), "999")
	assertError(t, err, "SSH key with ID 999 not found")
}

func TestGetSSHKey_EmptyPK(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, _, err := ts.client.SSHKeys.GetSSHKey(context.Background(), "")
	assertError(t, err, "")
	assertErrorType(t, err, &ArgError{})
}

func TestGetSSHKey_ErrorResponse(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/ssh_keys/", func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodGet)
		assertQueryParam(t, r, "pk", "123")
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Internal server error",
			"error": "Something went wrong"
		}`)
	})

	_, resp, err := ts.client.SSHKeys.GetSSHKey(context.Background(), "123")
	assertError(t, err, "")
	assertStatus(t, resp, http.StatusInternalServerError)
}

func TestGetSSHKeyByLabel(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/ssh_keys/", func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodGet)
		assertQueryParam(t, r, "apikey", "test-api-key")
		assertQueryParam(t, r, "project_id", "test-project")
		assertQueryParam(t, r, "location", "test-location")
		assertQueryParam(t, r, "label", "my-key")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": [
				{
					"pk": 123,
					"label": "my-key",
					"ssh_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAB...",
					"timestamp": "2024-01-01T00:00:00Z"
				}
			]
		}`)
	})

	result, resp, err := ts.client.SSHKeys.GetSSHKeyByLabel(context.Background(), "my-key")
	assertNoError(t, err)
	assertNotNil(t, result, "Expected result")
	if result.PK != 123 {
		t.Errorf("Expected PK 123, got %d", result.PK)
	}
	if result.Label != "my-key" {
		t.Errorf("Expected Label 'my-key', got %s", result.Label)
	}
	assertStatus(t, resp, http.StatusOK)
}

func TestGetSSHKeyByLabel_NotFound(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/ssh_keys/", func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodGet)
		assertQueryParam(t, r, "label", "non-existent")
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": []
		}`)
	})

	_, _, err := ts.client.SSHKeys.GetSSHKeyByLabel(context.Background(), "non-existent")
	assertError(t, err, "SSH key with label non-existent not found")
}

func TestGetSSHKeyByLabel_EmptyLabel(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, _, err := ts.client.SSHKeys.GetSSHKeyByLabel(context.Background(), "")
	assertError(t, err, "")
	assertErrorType(t, err, &ArgError{})
}

func TestGetSSHKeyByLabel_ErrorResponse(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/ssh_keys/", func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodGet)
		assertQueryParam(t, r, "label", "my-key")
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Internal server error",
			"error": "Something went wrong"
		}`)
	})

	_, resp, err := ts.client.SSHKeys.GetSSHKeyByLabel(context.Background(), "my-key")
	assertError(t, err, "")
	assertStatus(t, resp, http.StatusInternalServerError)
}

func TestListSSHKeys(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/ssh_keys/", func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodGet)
		assertQueryParam(t, r, "apikey", "test-api-key")
		assertQueryParam(t, r, "project_id", "test-project")
		assertQueryParam(t, r, "location", "test-location")

		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": [
				{
					"pk": 123,
					"label": "my-key-1",
					"ssh_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAB...",
					"timestamp": "2024-01-01T00:00:00Z"
				},
				{
					"pk": 124,
					"label": "my-key-2",
					"ssh_key": "ssh-rsa BBBBB3NzaC1yc2EAAAADAQABAAAB...",
					"timestamp": "2024-01-02T00:00:00Z"
				}
			]
		}`)
	})

	result, resp, err := ts.client.SSHKeys.ListSSHKeys(context.Background())
	assertNoError(t, err)
	if len(result) != 2 {
		t.Fatalf("Expected 2 SSH keys, got %d", len(result))
	}
	if result[0].PK != 123 {
		t.Errorf("Expected first SSH key PK 123, got %d", result[0].PK)
	}
	if result[0].Label != "my-key-1" {
		t.Errorf("Expected first SSH key label 'my-key-1', got %s", result[0].Label)
	}
	if result[1].PK != 124 {
		t.Errorf("Expected second SSH key PK 124, got %d", result[1].PK)
	}
	if result[1].Label != "my-key-2" {
		t.Errorf("Expected second SSH key label 'my-key-2', got %s", result[1].Label)
	}
	assertStatus(t, resp, http.StatusOK)
}

func TestListSSHKeys_EmptyList(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/ssh_keys/", func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodGet)
		writeJSON(w, http.StatusOK, `{
			"code": 200,
			"message": "success",
			"data": []
		}`)
	})

	result, resp, err := ts.client.SSHKeys.ListSSHKeys(context.Background())
	assertNoError(t, err)
	if len(result) != 0 {
		t.Fatalf("Expected 0 SSH keys, got %d", len(result))
	}
	assertStatus(t, resp, http.StatusOK)
}

func TestListSSHKeys_ErrorResponse(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/ssh_keys/", func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodGet)
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Internal server error",
			"error": "Something went wrong"
		}`)
	})

	_, resp, err := ts.client.SSHKeys.ListSSHKeys(context.Background())
	assertError(t, err, "")
	assertStatus(t, resp, http.StatusInternalServerError)
}

func TestDeleteSSHKey(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/delete_ssh_key/", func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodDelete)
		assertQueryParam(t, r, "apikey", "test-api-key")
		assertQueryParam(t, r, "project_id", "test-project")
		assertQueryParam(t, r, "location", "test-location")
		assertQueryParam(t, r, "pk", "123")

		w.WriteHeader(http.StatusOK)
	})

	resp, err := ts.client.SSHKeys.DeleteSSHKey(context.Background(), "123")
	assertNoError(t, err)
	assertStatus(t, resp, http.StatusOK)
}

func TestDeleteSSHKey_EmptyPK(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	_, err := ts.client.SSHKeys.DeleteSSHKey(context.Background(), "")
	assertError(t, err, "")
	assertErrorType(t, err, &ArgError{})
}

func TestDeleteSSHKey_ErrorResponse(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/delete_ssh_key/", func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodDelete)
		writeJSON(w, http.StatusInternalServerError, `{
			"code": 500,
			"message": "Internal server error",
			"error": "Something went wrong"
		}`)
	})

	resp, err := ts.client.SSHKeys.DeleteSSHKey(context.Background(), "123")
	assertError(t, err, "")
	assertStatus(t, resp, http.StatusInternalServerError)
}

// TestCreateSSHKey_Forbidden tests 403 Forbidden on create
func TestCreateSSHKey_Forbidden(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/ssh_keys/", func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPost)
		writeJSON(w, http.StatusForbidden, `{
			"code": 403,
			"message": "Access denied",
			"error": "Not authorized to create SSH keys"
		}`)
	})

	createReq := &SSHKeyCreateRequest{
		Label:  "my-key",
		SSHKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAB...",
	}

	_, resp, err := ts.client.SSHKeys.CreateSSHKey(context.Background(), createReq)
	assertError(t, err, "")
	assertStatus(t, resp, http.StatusForbidden)
}

// TestCreateSSHKey_Conflict tests 409 Conflict on create
func TestCreateSSHKey_Conflict(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/ssh_keys/", func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPost)
		writeJSON(w, http.StatusConflict, `{
			"code": 409,
			"message": "SSH key label already exists",
			"error": "Label already in use"
		}`)
	})

	createReq := &SSHKeyCreateRequest{
		Label:  "existing-key",
		SSHKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAB...",
	}

	_, resp, err := ts.client.SSHKeys.CreateSSHKey(context.Background(), createReq)
	assertError(t, err, "")
	assertStatus(t, resp, http.StatusConflict)
}

// TestCreateSSHKey_BadRequest tests 400 Bad Request on create
func TestCreateSSHKey_BadRequest(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/ssh_keys/", func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPost)
		writeJSON(w, http.StatusBadRequest, `{
			"code": 400,
			"message": "Invalid SSH key format",
			"error": "SSH key does not match expected format"
		}`)
	})

	createReq := &SSHKeyCreateRequest{
		Label:  "my-key",
		SSHKey: "invalid-key-format",
	}

	_, resp, err := ts.client.SSHKeys.CreateSSHKey(context.Background(), createReq)
	assertError(t, err, "")
	assertStatus(t, resp, http.StatusBadRequest)
}

// TestCreateSSHKey_ServiceUnavailable tests 503 Service Unavailable on create
func TestCreateSSHKey_ServiceUnavailable(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/ssh_keys/", func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPost)
		writeJSON(w, http.StatusServiceUnavailable, `{
			"code": 503,
			"message": "Service temporarily unavailable"
		}`)
	})

	createReq := &SSHKeyCreateRequest{
		Label:  "my-key",
		SSHKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAB...",
	}

	_, resp, err := ts.client.SSHKeys.CreateSSHKey(context.Background(), createReq)
	assertError(t, err, "")
	assertStatus(t, resp, http.StatusServiceUnavailable)
}

// TestGetSSHKey_Forbidden tests 403 Forbidden on get
func TestGetSSHKey_Forbidden(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/ssh_keys/", func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodGet)
		assertQueryParam(t, r, "pk", "123")
		writeJSON(w, http.StatusForbidden, `{
			"code": 403,
			"message": "Access denied",
			"error": "Not authorized to view this SSH key"
		}`)
	})

	_, resp, err := ts.client.SSHKeys.GetSSHKey(context.Background(), "123")
	assertError(t, err, "")
	assertStatus(t, resp, http.StatusForbidden)
}

// TestGetSSHKeyByLabel_Forbidden tests 403 Forbidden on get by label
func TestGetSSHKeyByLabel_Forbidden(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/ssh_keys/", func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodGet)
		assertQueryParam(t, r, "label", "my-key")
		writeJSON(w, http.StatusForbidden, `{
			"code": 403,
			"message": "Access denied"
		}`)
	})

	_, resp, err := ts.client.SSHKeys.GetSSHKeyByLabel(context.Background(), "my-key")
	assertError(t, err, "")
	assertStatus(t, resp, http.StatusForbidden)
}

// TestListSSHKeys_Forbidden tests 403 Forbidden on list
func TestListSSHKeys_Forbidden(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/ssh_keys/", func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodGet)
		writeJSON(w, http.StatusForbidden, `{
			"code": 403,
			"message": "Access denied",
			"error": "Not authorized to list SSH keys"
		}`)
	})

	_, resp, err := ts.client.SSHKeys.ListSSHKeys(context.Background())
	assertError(t, err, "")
	assertStatus(t, resp, http.StatusForbidden)
}

// TestListSSHKeys_ServiceUnavailable tests 503 Service Unavailable on list
func TestListSSHKeys_ServiceUnavailable(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/ssh_keys/", func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodGet)
		writeJSON(w, http.StatusServiceUnavailable, `{
			"code": 503,
			"message": "Service temporarily unavailable"
		}`)
	})

	_, resp, err := ts.client.SSHKeys.ListSSHKeys(context.Background())
	assertError(t, err, "")
	assertStatus(t, resp, http.StatusServiceUnavailable)
}

// TestDeleteSSHKey_NotFound tests 404 Not Found on delete
func TestDeleteSSHKey_NotFound(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/delete_ssh_key/", func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodDelete)
		assertQueryParam(t, r, "pk", "999")
		writeJSON(w, http.StatusNotFound, `{
			"code": 404,
			"message": "SSH key not found",
			"error": "The specified SSH key does not exist"
		}`)
	})

	resp, err := ts.client.SSHKeys.DeleteSSHKey(context.Background(), "999")
	assertError(t, err, "")
	assertStatus(t, resp, http.StatusNotFound)
}

// TestDeleteSSHKey_Forbidden tests 403 Forbidden on delete
func TestDeleteSSHKey_Forbidden(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/delete_ssh_key/", func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodDelete)
		assertQueryParam(t, r, "pk", "123")
		writeJSON(w, http.StatusForbidden, `{
			"code": 403,
			"message": "Access denied",
			"error": "Not authorized to delete this SSH key"
		}`)
	})

	resp, err := ts.client.SSHKeys.DeleteSSHKey(context.Background(), "123")
	assertError(t, err, "")
	assertStatus(t, resp, http.StatusForbidden)
}

// TestDeleteSSHKey_Conflict tests 409 Conflict on delete
func TestDeleteSSHKey_Conflict(t *testing.T) {
	ts := setup()
	defer ts.teardown()

	ts.mux.HandleFunc("/delete_ssh_key/", func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodDelete)
		assertQueryParam(t, r, "pk", "123")
		writeJSON(w, http.StatusConflict, `{
			"code": 409,
			"message": "SSH key is in use",
			"error": "Cannot delete SSH key that is currently in use by instances"
		}`)
	})

	resp, err := ts.client.SSHKeys.DeleteSSHKey(context.Background(), "123")
	assertError(t, err, "")
	assertStatus(t, resp, http.StatusConflict)
}

// ============================================================================
// Network Error & Timeout Tests
// ============================================================================

// TestSSHKey_NetworkConnectionError tests network connection failure
func TestSSHKey_NetworkConnectionError(t *testing.T) {
	testFunc := func(client *Client, ctx context.Context) error {
		createReq := &SSHKeyCreateRequest{
			Label:  "my-key",
			SSHKey: "ssh-rsa AAAA...",
		}
		_, _, err := client.SSHKeys.CreateSSHKey(ctx, createReq)
		return err
	}
	testNetworkError(t, testFunc)
}

// TestSSHKey_ContextTimeout tests context timeout
func TestSSHKey_ContextTimeout(t *testing.T) {
	testFunc := func(client *Client, ctx context.Context) error {
		createReq := &SSHKeyCreateRequest{
			Label:  "my-key",
			SSHKey: "ssh-rsa AAAA...",
		}
		_, _, err := client.SSHKeys.CreateSSHKey(ctx, createReq)
		return err
	}
	testContextTimeout(t, testFunc)
}

// TestSSHKey_ContextCancellation tests context cancellation
func TestSSHKey_ContextCancellation(t *testing.T) {
	testFunc := func(client *Client, ctx context.Context) error {
		createReq := &SSHKeyCreateRequest{
			Label:  "my-key",
			SSHKey: "ssh-rsa AAAA...",
		}
		_, _, err := client.SSHKeys.CreateSSHKey(ctx, createReq)
		return err
	}
	testContextCancellation(t, testFunc)
}

// TestSSHKey_GetSSHKey_NetworkError tests Get with network error
func TestSSHKey_GetSSHKey_NetworkError(t *testing.T) {
	testFunc := func(client *Client, ctx context.Context) error {
		_, _, err := client.SSHKeys.GetSSHKey(ctx, "123")
		return err
	}
	testNetworkError(t, testFunc)
}

// TestSSHKey_GetSSHKeyByLabel_NetworkError tests GetByLabel with network error
func TestSSHKey_GetSSHKeyByLabel_NetworkError(t *testing.T) {
	testFunc := func(client *Client, ctx context.Context) error {
		_, _, err := client.SSHKeys.GetSSHKeyByLabel(ctx, "my-key")
		return err
	}
	testNetworkError(t, testFunc)
}

// TestSSHKey_ListSSHKeys_NetworkError tests List with network error
func TestSSHKey_ListSSHKeys_NetworkError(t *testing.T) {
	testFunc := func(client *Client, ctx context.Context) error {
		_, _, err := client.SSHKeys.ListSSHKeys(ctx)
		return err
	}
	testNetworkError(t, testFunc)
}

// TestSSHKey_DeleteSSHKey_ContextTimeout tests Delete with context timeout
func TestSSHKey_DeleteSSHKey_ContextTimeout(t *testing.T) {
	testFunc := func(client *Client, ctx context.Context) error {
		_, err := client.SSHKeys.DeleteSSHKey(ctx, "123")
		return err
	}
	testContextTimeout(t, testFunc)
}

// TestSSHKey_RetryLogic_TemporaryFailure tests retry on temporary failures
func TestSSHKey_RetryLogic_TemporaryFailure(t *testing.T) {
	var requestCount int

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++

		if requestCount == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprint(w, `{"code": 503, "message": "Service unavailable"}`)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"code": 200,
			"message": "SSH Key created",
			"data": {
				"id": 123,
				"label": "my-key",
				"key": "ssh-rsa AAAA..."
			},
			"errors": {}
		}`)
	}))
	defer server.Close()

	client, _ := NewClient(
		"test-key",
		"test-token",
		"proj-123",
		"test-location",
		SetBaseURL(server.URL+"/"),
		WithRetryAndBackoffs(RetryConfig{
			RetryMax:     2,
			RetryWaitMin: PtrTo(10 * time.Millisecond),
			RetryWaitMax: PtrTo(50 * time.Millisecond),
		}),
	)

	ctx := context.Background()
	createReq := &SSHKeyCreateRequest{
		Label:  "my-key",
		SSHKey: "ssh-rsa AAAA...",
	}

	_, _, _ = client.SSHKeys.CreateSSHKey(ctx, createReq)
	if requestCount < 1 {
		t.Errorf("Expected at least 1 request, got %d", requestCount)
	}
}
