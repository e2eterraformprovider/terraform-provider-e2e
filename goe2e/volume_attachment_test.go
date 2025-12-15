package goe2e

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
)

func TestAttachVolume(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPost)
		assertQueryParam(t, r, "apikey", "test-api-key")
		assertQueryParam(t, r, "project_id", "test-project")
		assertQueryParam(t, r, "location", "test-location")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(buildSuccessResponse(200, "Volume attached successfully", map[string]interface{}{
			"node_id":   "node-123",
			"volume_id": "vol-456",
			"device":    "/dev/vdb",
			"status":    goe2econstants.VolumeAttachmentStatusAttached,
		})))
	}))
	defer server.Close()
	client, err := NewClient("test-api-key", "test-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	attachReq := &VolumeAttachmentRequest{
		NodeID:   "node-123",
		VolumeID: "vol-456",
	}
	result, _, err := client.VolumeAttachment.AttachVolume(context.Background(), attachReq)
	assertNoError(t, err)
	assertNotNil(t, result, "Expected result, got nil")
	if result.NodeID != "node-123" {
		t.Errorf("Expected NodeID node-123, got %s", result.NodeID)
	}
	if result.VolumeID != "vol-456" {
		t.Errorf("Expected VolumeID vol-456, got %s", result.VolumeID)
	}
	if result.Status != goe2econstants.VolumeAttachmentStatusAttached {
		t.Errorf("Expected Status %s, got %s", goe2econstants.VolumeAttachmentStatusAttached, result.Status)
	}
}

func TestDetachVolume(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertHTTPMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(buildSuccessResponse(200, "Volume detached successfully", map[string]interface{}{
			"node_id":   "node-123",
			"volume_id": "vol-456",
			"status":    goe2econstants.VolumeAttachmentStatusDetached,
		})))
	}))
	defer server.Close()
	client, err := NewClient("test-api-key", "test-token", "test-project", "test-location", SetBaseURL(server.URL), WithRetryAndBackoffs(RetryConfig{RetryMax: 0}))
	assertNoError(t, err)
	detachReq := &VolumeDetachmentRequest{
		NodeID:   "node-123",
		VolumeID: "vol-456",
	}
	_, err = client.VolumeAttachment.DetachVolume(context.Background(), detachReq)
	assertNoError(t, err)
}
