package models

import (
	"encoding/json"
	"testing"
)

func TestSshKeyMarshalling(t *testing.T) {
	sshKey := SshKey{
		Label:     "my-key",
		Ssh_key:   "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC...",
		Pk:        123,
		Timestamp: "2024-01-01T00:00:00Z",
	}

	// Test marshaling
	jsonData, err := json.Marshal(sshKey)
	if err != nil {
		t.Fatalf("Failed to marshal SshKey: %v", err)
	}

	// Test unmarshaling
	var unmarshaledKey SshKey
	err = json.Unmarshal(jsonData, &unmarshaledKey)
	if err != nil {
		t.Fatalf("Failed to unmarshal SshKey: %v", err)
	}

	// Verify fields
	if unmarshaledKey.Label != sshKey.Label {
		t.Errorf("Expected Label %s, got %s", sshKey.Label, unmarshaledKey.Label)
	}
	if unmarshaledKey.Ssh_key != sshKey.Ssh_key {
		t.Errorf("Expected Ssh_key %s, got %s", sshKey.Ssh_key, unmarshaledKey.Ssh_key)
	}
	if unmarshaledKey.Pk != sshKey.Pk {
		t.Errorf("Expected Pk %d, got %d", sshKey.Pk, unmarshaledKey.Pk)
	}
	if unmarshaledKey.Timestamp != sshKey.Timestamp {
		t.Errorf("Expected Timestamp %s, got %s", sshKey.Timestamp, unmarshaledKey.Timestamp)
	}
}

func TestSshKeyResponseMarshalling(t *testing.T) {
	sshKeyResponse := SshKeyResponse{
		Code:    200,
		Message: "success",
		Data: []SshKey{
			{
				Label:     "key-1",
				Ssh_key:   "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC1...",
				Pk:        1,
				Timestamp: "2024-01-01T00:00:00Z",
			},
			{
				Label:     "key-2",
				Ssh_key:   "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC2...",
				Pk:        2,
				Timestamp: "2024-01-02T00:00:00Z",
			},
		},
		Error: []interface{}{},
	}

	// Test marshaling
	jsonData, err := json.Marshal(sshKeyResponse)
	if err != nil {
		t.Fatalf("Failed to marshal SshKeyResponse: %v", err)
	}

	// Test unmarshaling
	var unmarshaledResponse SshKeyResponse
	err = json.Unmarshal(jsonData, &unmarshaledResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal SshKeyResponse: %v", err)
	}

	// Verify fields
	if unmarshaledResponse.Code != sshKeyResponse.Code {
		t.Errorf("Expected Code %d, got %d", sshKeyResponse.Code, unmarshaledResponse.Code)
	}
	if len(unmarshaledResponse.Data) != len(sshKeyResponse.Data) {
		t.Errorf("Expected %d keys, got %d", len(sshKeyResponse.Data), len(unmarshaledResponse.Data))
	}
	if len(unmarshaledResponse.Data) > 0 {
		if unmarshaledResponse.Data[0].Label != sshKeyResponse.Data[0].Label {
			t.Errorf("Expected first key label %s, got %s", sshKeyResponse.Data[0].Label, unmarshaledResponse.Data[0].Label)
		}
	}
}

func TestAddSshKeyMarshalling(t *testing.T) {
	addSshKey := AddSshKey{
		Label:    "my-key",
		SshKey:   "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC...",
		Location: "us-east",
	}

	// Test marshaling
	jsonData, err := json.Marshal(addSshKey)
	if err != nil {
		t.Fatalf("Failed to marshal AddSshKey: %v", err)
	}

	// Test unmarshaling
	var unmarshaledAddSshKey AddSshKey
	err = json.Unmarshal(jsonData, &unmarshaledAddSshKey)
	if err != nil {
		t.Fatalf("Failed to unmarshal AddSshKey: %v", err)
	}

	// Verify fields
	if unmarshaledAddSshKey.Label != addSshKey.Label {
		t.Errorf("Expected Label %s, got %s", addSshKey.Label, unmarshaledAddSshKey.Label)
	}
	if unmarshaledAddSshKey.SshKey != addSshKey.SshKey {
		t.Errorf("Expected SshKey %s, got %s", addSshKey.SshKey, unmarshaledAddSshKey.SshKey)
	}
	if unmarshaledAddSshKey.Location != addSshKey.Location {
		t.Errorf("Expected Location %s, got %s", addSshKey.Location, unmarshaledAddSshKey.Location)
	}
}

func TestSshKeyJSONTags(t *testing.T) {
	jsonStr := `{
		"label": "test-key",
		"ssh_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC...",
		"pk": 456,
		"timestamp": "2024-01-15T10:30:00Z"
	}`

	var sshKey SshKey
	err := json.Unmarshal([]byte(jsonStr), &sshKey)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if sshKey.Label != "test-key" {
		t.Errorf("Expected Label test-key, got %s", sshKey.Label)
	}
	if sshKey.Pk != 456 {
		t.Errorf("Expected Pk 456, got %d", sshKey.Pk)
	}
	if sshKey.Timestamp != "2024-01-15T10:30:00Z" {
		t.Errorf("Expected Timestamp 2024-01-15T10:30:00Z, got %s", sshKey.Timestamp)
	}
}

func TestSshKeyResponseWithErrors(t *testing.T) {
	jsonStr := `{
		"code": 400,
		"message": "Bad Request",
		"data": [],
		"error": ["Invalid SSH key format"]
	}`

	var response SshKeyResponse
	err := json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if response.Code != 400 {
		t.Errorf("Expected Code 400, got %d", response.Code)
	}
	if len(response.Error) != 1 {
		t.Errorf("Expected 1 error, got %d", len(response.Error))
	}
	if len(response.Data) != 0 {
		t.Errorf("Expected empty Data, got %d items", len(response.Data))
	}
}
