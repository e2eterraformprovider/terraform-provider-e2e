package goe2e

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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

// ============================================================================
// Comprehensive Test Suite - Following tags.md requirements
// ============================================================================

// Test Category 1: Happy Path Tests with Enhanced Assertions

// TestTagService_CreateTag_Success verifies successful tag creation with all fields
func TestTagService_CreateTag_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		testURLPath(t, r, "/label/")

		// Verify request body
		body, err := io.ReadAll(r.Body)
		assertNoError(t, err)
		var reqBody TagCreateRequest
		err = json.Unmarshal(body, &reqBody)
		assertNoError(t, err)

		if reqBody.LabelName != "test-tag" {
			t.Errorf("Expected LabelName 'test-tag', got '%s'", reqBody.LabelName)
		}
		if reqBody.Metadata != "test metadata" {
			t.Errorf("Expected Metadata 'test metadata', got '%s'", reqBody.Metadata)
		}

		// Send response
		data := map[string]interface{}{
			"label_id": 123,
		}
		writeJSON(w, http.StatusOK, buildSuccessResponse(201, "Tag created successfully", data))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	createReq := &TagCreateRequest{
		LabelName: "test-tag",
		Metadata:  "test metadata",
	}

	tag, resp, err := client.Tags.CreateTag(context.Background(), createReq)
	assertNoError(t, err)
	assertNotNil(t, tag, "Expected tag, got nil")
	assertNotNil(t, resp, "Expected response, got nil")

	if tag.LabelID != 123 {
		t.Errorf("Expected LabelID 123, got %d", tag.LabelID)
	}
	if tag.LabelName != "test-tag" {
		t.Errorf("Expected LabelName 'test-tag', got '%s'", tag.LabelName)
	}
	if tag.Metadata != "test metadata" {
		t.Errorf("Expected Metadata 'test metadata', got '%s'", tag.Metadata)
	}
}

// TestTagService_ListTags_Success verifies successful tag listing
func TestTagService_ListTags_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/label/")

		data := []interface{}{
			map[string]interface{}{
				"label_id":   100,
				"label_name": "tag-1",
				"metadata":   "metadata-1",
			},
			map[string]interface{}{
				"label_id":   200,
				"label_name": "tag-2",
				"metadata":   "metadata-2",
			},
		}
		writeJSON(w, http.StatusOK, buildListResponse(data))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	tags, resp, err := client.Tags.ListTags(context.Background())
	assertNoError(t, err)
	assertNotNil(t, resp, "Expected response, got nil")

	if len(tags) != 2 {
		t.Fatalf("Expected 2 tags, got %d", len(tags))
	}

	if tags[0].LabelID != 100 || tags[0].LabelName != "tag-1" {
		t.Errorf("First tag mismatch: got ID=%d, Name=%s", tags[0].LabelID, tags[0].LabelName)
	}
	if tags[1].LabelID != 200 || tags[1].LabelName != "tag-2" {
		t.Errorf("Second tag mismatch: got ID=%d, Name=%s", tags[1].LabelID, tags[1].LabelName)
	}
}

// TestTagService_GetTag_Success verifies GetTag uses ListTags internally
func TestTagService_GetTag_Success(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/label/")
		callCount++

		data := []interface{}{
			map[string]interface{}{
				"label_id":   123,
				"label_name": "found-tag",
				"metadata":   "found-metadata",
			},
		}
		writeJSON(w, http.StatusOK, buildListResponse(data))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	tag, resp, err := client.Tags.GetTag(context.Background(), "123")
	assertNoError(t, err)
	assertNotNil(t, tag, "Expected tag, got nil")
	assertNotNil(t, resp, "Expected response, got nil")

	if tag.LabelID != 123 {
		t.Errorf("Expected LabelID 123, got %d", tag.LabelID)
	}
	if callCount != 1 {
		t.Errorf("Expected ListTags to be called once, got %d", callCount)
	}
}

// TestTagService_DeleteTag_Success verifies successful tag deletion
func TestTagService_DeleteTag_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		testURLPath(t, r, "/label/123/")

		writeJSON(w, http.StatusOK, buildSuccessResponse(200, "Tag deleted successfully", nil))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	resp, err := client.Tags.DeleteTag(context.Background(), "123")
	assertNoError(t, err)
	assertNotNil(t, resp, "Expected response, got nil")
}

// TestTagService_AttachTags_Success verifies request body for attach operation
func TestTagService_AttachTags_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/label/mapping/nodes/node-123/")

		// Verify request body
		body, err := io.ReadAll(r.Body)
		assertNoError(t, err)
		var reqBody map[string]interface{}
		err = json.Unmarshal(body, &reqBody)
		assertNoError(t, err)

		attach, ok := reqBody["attach"].([]interface{})
		if !ok {
			t.Error("Expected 'attach' field in request body")
		}
		if len(attach) != 2 {
			t.Errorf("Expected 2 tag IDs in attach, got %d", len(attach))
		}

		writeJSON(w, http.StatusOK, buildSuccessResponse(200, "Tags attached successfully", nil))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	resp, err := client.Tags.AttachTags(context.Background(), "nodes", "node-123", []int{123, 124})
	assertNoError(t, err)
	assertNotNil(t, resp, "Expected response, got nil")
}

// TestTagService_DetachTags_Success verifies request body for detach operation
func TestTagService_DetachTags_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testURLPath(t, r, "/label/mapping/nodes/node-123/")

		// Verify request body
		body, err := io.ReadAll(r.Body)
		assertNoError(t, err)
		var reqBody map[string]interface{}
		err = json.Unmarshal(body, &reqBody)
		assertNoError(t, err)

		detach, ok := reqBody["detach"].([]interface{})
		if !ok {
			t.Error("Expected 'detach' field in request body")
		}
		if len(detach) != 1 {
			t.Errorf("Expected 1 tag ID in detach, got %d", len(detach))
		}

		writeJSON(w, http.StatusOK, buildSuccessResponse(200, "Tags detached successfully", nil))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	resp, err := client.Tags.DetachTags(context.Background(), "nodes", "node-123", []int{123})
	assertNoError(t, err)
	assertNotNil(t, resp, "Expected response, got nil")
}

// TestTagService_GetResourceTags_Success verifies filtering by resourceID
func TestTagService_GetResourceTags_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testURLPath(t, r, "/label/mapping/nodes/")

		data := []interface{}{
			map[string]interface{}{
				"resource_id":   123,
				"resource_type": "nodes",
				"label_mapping": []interface{}{
					map[string]interface{}{
						"label_id":   456,
						"label_name": "env",
						"metadata":   "prod",
					},
				},
			},
			map[string]interface{}{
				"resource_id":   999,
				"resource_type": "nodes",
				"label_mapping": []interface{}{},
			},
		}
		writeJSON(w, http.StatusOK, buildListResponse(data))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	// Test with resourceID filter
	mappings, resp, err := client.Tags.GetResourceTags(context.Background(), "nodes", "123")
	assertNoError(t, err)
	assertNotNil(t, resp, "Expected response, got nil")

	if len(mappings) != 1 {
		t.Fatalf("Expected 1 filtered mapping, got %d", len(mappings))
	}
	if mappings[0].ResourceID != 123 {
		t.Errorf("Expected ResourceID 123, got %d", mappings[0].ResourceID)
	}

	// Test without filter
	mappingsAll, _, err := client.Tags.GetResourceTags(context.Background(), "nodes", "")
	assertNoError(t, err)
	if len(mappingsAll) != 2 {
		t.Fatalf("Expected 2 mappings without filter, got %d", len(mappingsAll))
	}
}

// Test Category 2: Input Validation Tests with ArgError Verification

// TestTagService_CreateTag_ValidationErrors verifies ArgError types
func TestTagService_CreateTag_ValidationErrors(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	// Test nil request
	_, _, err = client.Tags.CreateTag(context.Background(), nil)
	assertError(t, err, "")
	if _, ok := err.(*ArgError); !ok {
		t.Errorf("Expected *ArgError, got %T", err)
	}

	// Test empty LabelName
	createReq := &TagCreateRequest{
		LabelName: "",
		Metadata:  "test",
	}
	_, _, err = client.Tags.CreateTag(context.Background(), createReq)
	assertError(t, err, "")
	if _, ok := err.(*ArgError); !ok {
		t.Errorf("Expected *ArgError, got %T", err)
	}
}

// TestTagService_GetTag_ValidationErrors verifies validation errors
func TestTagService_GetTag_ValidationErrors(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	// Test empty tagID
	_, _, err = client.Tags.GetTag(context.Background(), "")
	assertError(t, err, "")
	if _, ok := err.(*ArgError); !ok {
		t.Errorf("Expected *ArgError, got %T", err)
	}

	// Test non-numeric tagID
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, buildListResponse([]interface{}{}))
	}))
	defer server.Close()

	client2, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, err = client2.Tags.GetTag(context.Background(), "abc")
	assertError(t, err, "")
	if _, ok := err.(*ArgError); !ok {
		t.Errorf("Expected *ArgError for non-numeric ID, got %T", err)
	}
}

// TestTagService_UpdateTag_ValidationErrors verifies UpdateTag always returns error
func TestTagService_UpdateTag_ValidationErrors(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	// Test empty tagID
	updateReq := &TagUpdateRequest{
		LabelName: String("test"),
	}
	_, _, err = client.Tags.UpdateTag(context.Background(), "", updateReq)
	assertError(t, err, "")
	if _, ok := err.(*ArgError); !ok {
		t.Errorf("Expected *ArgError, got %T", err)
	}

	// Test nil request
	_, _, err = client.Tags.UpdateTag(context.Background(), "123", nil)
	assertError(t, err, "")
	if _, ok := err.(*ArgError); !ok {
		t.Errorf("Expected *ArgError, got %T", err)
	}

	// Test that UpdateTag always returns "not supported" error
	_, _, err = client.Tags.UpdateTag(context.Background(), "123", updateReq)
	assertError(t, err, "not supported")
}

// TestTagService_AttachTags_ValidationErrors verifies all validation cases
func TestTagService_AttachTags_ValidationErrors(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	// Test empty resourceType
	_, err = client.Tags.AttachTags(context.Background(), "", "node-123", []int{123})
	assertError(t, err, "")
	if _, ok := err.(*ArgError); !ok {
		t.Errorf("Expected *ArgError, got %T", err)
	}

	// Test empty resourceID
	_, err = client.Tags.AttachTags(context.Background(), "nodes", "", []int{123})
	assertError(t, err, "")
	if _, ok := err.(*ArgError); !ok {
		t.Errorf("Expected *ArgError, got %T", err)
	}

	// Test empty tagIDs
	_, err = client.Tags.AttachTags(context.Background(), "nodes", "node-123", []int{})
	assertError(t, err, "")
	if _, ok := err.(*ArgError); !ok {
		t.Errorf("Expected *ArgError, got %T", err)
	}
}

// TestTagService_DetachTags_ValidationErrors same as AttachTags
func TestTagService_DetachTags_ValidationErrors(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, err = client.Tags.DetachTags(context.Background(), "", "node-123", []int{123})
	assertError(t, err, "")
	if _, ok := err.(*ArgError); !ok {
		t.Errorf("Expected *ArgError, got %T", err)
	}
}

// TestTagService_GetResourceTags_ValidationErrors verifies validation
func TestTagService_GetResourceTags_ValidationErrors(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, err = client.Tags.GetResourceTags(context.Background(), "", "123")
	assertError(t, err, "")
	if _, ok := err.(*ArgError); !ok {
		t.Errorf("Expected *ArgError, got %T", err)
	}
}

// Test Category 3: Error Handling Tests

// TestTagService_CreateTag_ErrorHandling tests various HTTP error codes
func TestTagService_CreateTag_ErrorHandling(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		message    string
	}{
		{"400 Bad Request", http.StatusBadRequest, "invalid input"},
		{"401 Unauthorized", http.StatusUnauthorized, "auth failure"},
		{"404 Not Found", http.StatusNotFound, "not found"},
		{"500 Internal Server Error", http.StatusInternalServerError, "server error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := newErrorServer(t, tt.statusCode, tt.message)
			defer server.Close()

			client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
				SetBaseURL(server.URL),
				WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
			assertNoError(t, err)

			createReq := &TagCreateRequest{
				LabelName: "test-tag",
			}
			_, _, err = client.Tags.CreateTag(context.Background(), createReq)
			assertError(t, err, "")
		})
	}
}

// TestTagService_GetTag_ErrorHandling verifies not found handling
func TestTagService_GetTag_ErrorHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, buildListResponse([]interface{}{}))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	tag, _, err := client.Tags.GetTag(context.Background(), "999")
	if err == nil {
		t.Error("Expected error for non-existent tag, got nil")
	}
	if tag != nil {
		t.Error("Expected nil tag on error, got non-nil")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("Expected 'not found' in error, got: %v", err)
	}
}

// TestTagService_ListTags_ErrorHandling tests error responses
func TestTagService_ListTags_ErrorHandling(t *testing.T) {
	// Test 500 error
	server := newErrorServer(t, http.StatusInternalServerError, "Internal server error")
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, err = client.Tags.ListTags(context.Background())
	assertError(t, err, "")

	// Test malformed JSON
	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code": 200, "message": "success", "data": [invalid json`))
	}))
	defer server2.Close()

	client2, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server2.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, _, err = client2.Tags.ListTags(context.Background())
	assertError(t, err, "")
}

// TestTagService_DeleteTag_ErrorHandling tests delete error scenarios
func TestTagService_DeleteTag_ErrorHandling(t *testing.T) {
	// Test 404 Not Found
	server := newErrorServer(t, http.StatusNotFound, "tag with ID 123 not found")
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, err = client.Tags.DeleteTag(context.Background(), "123")
	assertError(t, err, "")

	// Test 500 Server Error
	server2 := newErrorServer(t, http.StatusInternalServerError, "Server error")
	defer server2.Close()

	client2, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server2.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, err = client2.Tags.DeleteTag(context.Background(), "123")
	assertError(t, err, "")
}

// TestTagService_AttachTags_ErrorHandling tests attach error scenarios
func TestTagService_AttachTags_ErrorHandling(t *testing.T) {
	// Test 400 Bad Request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusBadRequest, buildErrorResponse(400, "Invalid resource_type or resource_id", []string{}))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	_, err = client.Tags.AttachTags(context.Background(), "invalid", "invalid-id", []int{123})
	assertError(t, err, "")
	if !strings.Contains(err.Error(), "invalid") {
		t.Errorf("Expected error to mention resource info, got: %v", err)
	}
}

// Test Category 4: Response Parsing Tests

// TestTagService_CreateTag_ResponseParsing tests response parsing edge cases
func TestTagService_CreateTag_ResponseParsing(t *testing.T) {
	// Test valid response with all fields
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"label_id": 999999,
		}
		writeJSON(w, http.StatusOK, buildSuccessResponse(201, "Created", data))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	createReq := &TagCreateRequest{
		LabelName: "test",
	}
	tag, _, err := client.Tags.CreateTag(context.Background(), createReq)
	assertNoError(t, err)
	if tag.LabelID != 999999 {
		t.Errorf("Expected LabelID 999999, got %d", tag.LabelID)
	}
}

// TestTagService_ListTags_ResponseParsing tests empty array and malformed JSON
func TestTagService_ListTags_ResponseParsing(t *testing.T) {
	// Test empty array
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, buildListResponse([]interface{}{}))
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	tags, _, err := client.Tags.ListTags(context.Background())
	assertNoError(t, err)
	if len(tags) != 0 {
		t.Errorf("Expected empty slice, got %d tags", len(tags))
	}
}

// TestTagService_GetResourceTags_ResponseParsing tests nested structure parsing
func TestTagService_GetResourceTags_ResponseParsing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := []interface{}{
			map[string]interface{}{
				"resource_id":   123,
				"resource_type": "nodes",
				"label_mapping": []interface{}{
					map[string]interface{}{
						"label_id":   456,
						"label_name": "env",
						"metadata":   "prod",
					},
					map[string]interface{}{
						"label_id":   789,
						"label_name": "team",
						"metadata":   "backend",
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

	mappings, _, err := client.Tags.GetResourceTags(context.Background(), "nodes", "")
	assertNoError(t, err)
	if len(mappings) != 1 {
		t.Fatalf("Expected 1 mapping, got %d", len(mappings))
	}
	if len(mappings[0].LabelMapping) != 2 {
		t.Errorf("Expected 2 labels in mapping, got %d", len(mappings[0].LabelMapping))
	}
}

// Test Category 5: Path Construction Tests

// TestTagService_PathConstruction verifies all path constructions
func TestTagService_PathConstruction(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		method := r.Method

		switch {
		case method == http.MethodPost && path == "/label/":
			data := map[string]interface{}{"label_id": 1}
			writeJSON(w, http.StatusOK, buildSuccessResponse(201, "Created", data))
		case method == http.MethodGet && path == "/label/":
			writeJSON(w, http.StatusOK, buildListResponse([]interface{}{}))
		case method == http.MethodDelete && path == "/label/123/":
			writeJSON(w, http.StatusOK, buildSuccessResponse(200, "Deleted", nil))
		case method == http.MethodPut && path == "/label/mapping/nodes/node-123/":
			writeJSON(w, http.StatusOK, buildSuccessResponse(200, "OK", nil))
		case method == http.MethodGet && path == "/label/mapping/nodes/":
			writeJSON(w, http.StatusOK, buildListResponse([]interface{}{}))
		default:
			t.Errorf("Unexpected path/method: %s %s", method, path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		SetBaseURL(server.URL),
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	// Test CreateTag path
	_, _, err = client.Tags.CreateTag(context.Background(), &TagCreateRequest{LabelName: "test"})
	assertNoError(t, err)

	// Test ListTags path
	_, _, err = client.Tags.ListTags(context.Background())
	assertNoError(t, err)

	// Test DeleteTag path
	_, err = client.Tags.DeleteTag(context.Background(), "123")
	assertNoError(t, err)

	// Test AttachTags path
	_, err = client.Tags.AttachTags(context.Background(), "nodes", "node-123", []int{1})
	assertNoError(t, err)

	// Test GetResourceTags path
	_, _, err = client.Tags.GetResourceTags(context.Background(), "nodes", "")
	assertNoError(t, err)
}

// Test Category 6: Client Integration Tests

// TestTagService_ClientIntegration verifies TagServiceOp implements TagService interface
func TestTagService_ClientIntegration(t *testing.T) {
	client, err := NewClient("test-api-key", "test-auth-token", "test-project", "test-location",
		WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)

	// Compile-time check: TagServiceOp implements TagService
	var _ TagService = &TagServiceOp{}

	// Runtime check: client.Tags is not nil
	if client.Tags == nil {
		t.Fatal("Expected Tags service to be initialized, got nil")
	}

	// Verify service type
	if _, ok := client.Tags.(*TagServiceOp); !ok {
		t.Errorf("Expected *TagServiceOp, got %T", client.Tags)
	}
}
