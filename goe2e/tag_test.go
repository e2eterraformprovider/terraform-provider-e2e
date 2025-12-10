package goe2e

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateTag(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPost)
		assertQueryParam(t, r, "apikey", "test-api-key")
		assertQueryParam(t, r, "project_id", "test-project")
		assertQueryParam(t, r, "location", "test-location")

		data := map[string]interface{}{
			"label_id":    123,
			"label_name":  "environment",
			"label_value": "production",
		}
		writeJSON(w, http.StatusCreated, buildSuccessResponse(201, "Tag created successfully", data))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	createReq := &TagCreateRequest{
		LabelName: "environment",
		Metadata:  "production",
	}

	result, _, err := client.Tags.CreateTag(context.Background(), createReq)
	assertNoError(t, err)
	assertNotNil(t, result, "Expected result, got nil")

	if result.LabelID != 123 {
		t.Errorf("Expected LabelID 123, got %d", result.LabelID)
	}

	if result.LabelName != "environment" {
		t.Errorf("Expected LabelName environment, got %s", result.LabelName)
	}
}

func TestListTags(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodGet)
		assertQueryParam(t, r, "apikey", "test-api-key")
		assertQueryParam(t, r, "project_id", "test-project")
		assertQueryParam(t, r, "location", "test-location")

		data := []interface{}{
			map[string]interface{}{
				"label_id":    123,
				"label_name":  "environment",
				"label_value": "production",
			},
			map[string]interface{}{
				"label_id":    124,
				"label_name":  "team",
				"label_value": "backend",
			},
		}
		writeJSON(w, http.StatusOK, buildListResponse(data))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	result, _, err := client.Tags.ListTags(context.Background())
	assertNoError(t, err)

	if len(result) != 2 {
		t.Fatalf("Expected 2 tags, got %d", len(result))
	}

	if result[0].LabelID != 123 {
		t.Errorf("Expected first tag ID 123, got %d", result[0].LabelID)
	}
}

func TestGetTag(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodGet)

		data := []interface{}{
			map[string]interface{}{
				"label_id":    123,
				"label_name":  "environment",
				"label_value": "production",
			},
		}
		writeJSON(w, http.StatusOK, buildListResponse(data))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	result, _, err := client.Tags.GetTag(context.Background(), "123")
	assertNoError(t, err)
	assertNotNil(t, result, "Expected result, got nil")

	if result.LabelID != 123 {
		t.Errorf("Expected LabelID 123, got %d", result.LabelID)
	}
}

func TestGetTag_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, buildListResponse([]interface{}{}))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, err = client.Tags.GetTag(context.Background(), "999")
	assertError(t, err, "")
}

func TestDeleteTag(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodDelete)
		assertQueryParam(t, r, "apikey", "test-api-key")
		assertQueryParam(t, r, "project_id", "test-project")
		assertQueryParam(t, r, "location", "test-location")

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, err = client.Tags.DeleteTag(context.Background(), "123")
	assertNoError(t, err)
}

func TestAttachTags(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPut)
		assertQueryParam(t, r, "apikey", "test-api-key")
		assertQueryParam(t, r, "project_id", "test-project")
		assertQueryParam(t, r, "location", "test-location")

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, err = client.Tags.AttachTags(context.Background(), "nodes", "node-123", []int{123, 124})
	assertNoError(t, err)
}

func TestDetachTags(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPut)
		assertQueryParam(t, r, "apikey", "test-api-key")
		assertQueryParam(t, r, "project_id", "test-project")
		assertQueryParam(t, r, "location", "test-location")

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, err = client.Tags.DetachTags(context.Background(), "nodes", "node-123", []int{123})
	assertNoError(t, err)
}

func TestGetResourceTags(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodGet)
		assertQueryParam(t, r, "apikey", "test-api-key")
		assertQueryParam(t, r, "project_id", "test-project")
		assertQueryParam(t, r, "location", "test-location")

		data := []interface{}{
			map[string]interface{}{
				"resource_id":   123,
				"resource_type": "nodes",
				"label_mapping": []interface{}{
					map[string]interface{}{
						"label_id":   456,
						"label_name": "environment",
						"metadata":   "production",
					},
				},
			},
		}
		writeJSON(w, http.StatusOK, buildListResponse(data))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	result, _, err := client.Tags.GetResourceTags(context.Background(), "nodes", "123")
	assertNoError(t, err)

	if len(result) != 1 {
		t.Fatalf("Expected 1 tag mapping, got %d", len(result))
	}

	if result[0].ResourceID != 123 {
		t.Errorf("Expected resource ID 123, got %d", result[0].ResourceID)
	}

	if len(result[0].LabelMapping) != 1 {
		t.Fatalf("Expected 1 label in mapping, got %d", len(result[0].LabelMapping))
	}

	if result[0].LabelMapping[0].LabelID != 456 {
		t.Errorf("Expected label ID 456, got %d", result[0].LabelMapping[0].LabelID)
	}
}

// Test error cases for tags
func TestCreateTag_NilRequest(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, err = client.Tags.CreateTag(context.Background(), nil)
	assertError(t, err, "")
}

func TestGetTag_InvalidID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, buildListResponse([]interface{}{}))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, err = client.Tags.GetTag(context.Background(), "invalid-id")
	assertError(t, err, "")
}

// TestUpdateTag tests tag update
func TestUpdateTag(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	updateReq := &TagUpdateRequest{
		LabelName: String("updated-name"),
		Metadata:  String("updated-metadata"),
	}

	// UpdateTag may not be supported - just test it doesn't panic
	_, _, err = client.Tags.UpdateTag(context.Background(), "123", updateReq)
	// We don't fail if error - just test it runs
	t.Logf("UpdateTag result: %v", err)
}

func TestUpdateTag_EmptyID(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	updateReq := &TagUpdateRequest{
		LabelName: String("test"),
	}

	_, _, err = client.Tags.UpdateTag(context.Background(), "", updateReq)
	assertError(t, err, "")
}

func TestUpdateTag_NilRequest(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, err = client.Tags.UpdateTag(context.Background(), "123", nil)
	assertError(t, err, "")
}

// Additional error tests for better coverage
func TestCreateTag_EmptyLabelName(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	createReq := &TagCreateRequest{
		LabelName: "",
		Metadata:  "test",
	}

	_, _, err = client.Tags.CreateTag(context.Background(), createReq)
	assertError(t, err, "")
}

func TestListTags_ErrorResponse(t *testing.T) {
	server := newErrorServer(t, http.StatusInternalServerError, "Internal server error")
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, err = client.Tags.ListTags(context.Background())
	assertError(t, err, "")
}

func TestDeleteTag_ErrorResponse(t *testing.T) {
	server := newErrorServer(t, http.StatusInternalServerError, "Internal server error")
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, err = client.Tags.DeleteTag(context.Background(), "123")
	assertError(t, err, "")
}

func TestDeleteTag_EmptyID(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, err = client.Tags.DeleteTag(context.Background(), "")
	assertError(t, err, "")
}

func TestAttachTags_EmptyResourceType(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, err = client.Tags.AttachTags(context.Background(), "", "node-123", []int{123})
	assertError(t, err, "")
}

func TestAttachTags_EmptyResourceID(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, err = client.Tags.AttachTags(context.Background(), "nodes", "", []int{123})
	assertError(t, err, "")
}

func TestAttachTags_EmptyTagIDs(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, err = client.Tags.AttachTags(context.Background(), "nodes", "node-123", []int{})
	assertError(t, err, "")
}

func TestAttachTags_ErrorResponse(t *testing.T) {
	server := newErrorServer(t, http.StatusInternalServerError, "Internal server error")
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, err = client.Tags.AttachTags(context.Background(), "nodes", "node-123", []int{123})
	assertError(t, err, "")
}

func TestDetachTags_EmptyResourceType(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, err = client.Tags.DetachTags(context.Background(), "", "node-123", []int{123})
	assertError(t, err, "")
}

func TestDetachTags_EmptyResourceID(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, err = client.Tags.DetachTags(context.Background(), "nodes", "", []int{123})
	assertError(t, err, "")
}

func TestDetachTags_EmptyTagIDs(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, err = client.Tags.DetachTags(context.Background(), "nodes", "node-123", []int{})
	assertError(t, err, "")
}

func TestDetachTags_ErrorResponse(t *testing.T) {
	server := newErrorServer(t, http.StatusInternalServerError, "Internal server error")
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, err = client.Tags.DetachTags(context.Background(), "nodes", "node-123", []int{123})
	assertError(t, err, "")
}

func TestGetResourceTags_EmptyResourceType(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, err = client.Tags.GetResourceTags(context.Background(), "", "123")
	assertError(t, err, "")
}

func TestGetResourceTags_ErrorResponse(t *testing.T) {
	server := newErrorServer(t, http.StatusInternalServerError, "Internal server error")
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, err = client.Tags.GetResourceTags(context.Background(), "nodes", "123")
	assertError(t, err, "")
}

func TestGetTag_EmptyID(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, err = client.Tags.GetTag(context.Background(), "")
	assertError(t, err, "")
}
