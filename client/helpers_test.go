package client

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestRemoveExtraKeysLoadBalancer(t *testing.T) {
	tests := []struct {
		name        string
		input       map[string]interface{}
		expected    map[string]interface{}
		expectError bool
	}{
		{
			name: "Remove enable_eos_logger with empty access_key",
			input: map[string]interface{}{
				"name": "test-lb",
				"enable_eos_logger": map[string]interface{}{
					"access_key": "",
					"secret_key": "secret",
				},
			},
			expected: map[string]interface{}{
				"name": "test-lb",
			},
			expectError: false,
		},
		{
			name: "Keep enable_eos_logger with non-empty access_key",
			input: map[string]interface{}{
				"name": "test-lb",
				"enable_eos_logger": map[string]interface{}{
					"access_key": "my-access-key",
					"secret_key": "secret",
				},
			},
			expected: map[string]interface{}{
				"name": "test-lb",
				"enable_eos_logger": map[string]interface{}{
					"access_key": "my-access-key",
					"secret_key": "secret",
				},
			},
			expectError: false,
		},
		{
			name: "Handle missing enable_eos_logger",
			input: map[string]interface{}{
				"name": "test-lb",
				"port": 80,
			},
			expected: map[string]interface{}{
				"name": "test-lb",
				"port": 80,
			},
			expectError: false,
		},
		{
			name: "Handle enable_eos_logger with nil access_key",
			input: map[string]interface{}{
				"name": "test-lb",
				"enable_eos_logger": map[string]interface{}{
					"secret_key": "secret",
				},
			},
			expected: map[string]interface{}{
				"name": "test-lb",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert input to JSON buffer
			inputJSON, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("Failed to marshal input: %v", err)
			}
			buf := bytes.NewBuffer(inputJSON)

			// Call the function
			result, err := RemoveExtraKeysLoadBalancer(buf)

			// Check error expectation
			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}

			if !tt.expectError {
				// Unmarshal result
				var resultMap map[string]interface{}
				err = json.Unmarshal(result.Bytes(), &resultMap)
				if err != nil {
					t.Fatalf("Failed to unmarshal result: %v", err)
				}

				// Compare results
				expectedJSON, _ := json.Marshal(tt.expected)
				resultJSON, _ := json.Marshal(resultMap)

				if string(expectedJSON) != string(resultJSON) {
					t.Errorf("Expected: %s, got: %s", string(expectedJSON), string(resultJSON))
				}
			}
		})
	}
}

func TestGenerateSSHKeyMap(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		expected []map[string]interface{}
	}{
		{
			name:  "Empty input",
			input: []interface{}{},
			expected: []map[string]interface{}{},
		},
		{
			name: "Single SSH key",
			input: []interface{}{
				"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC...",
			},
			expected: []map[string]interface{}{
				{
					"label":   "ssh-key-1",
					"ssh_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC...",
				},
			},
		},
		{
			name: "Multiple SSH keys",
			input: []interface{}{
				"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC1...",
				"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC2...",
				"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC3...",
			},
			expected: []map[string]interface{}{
				{
					"label":   "ssh-key-1",
					"ssh_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC1...",
				},
				{
					"label":   "ssh-key-2",
					"ssh_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC2...",
				},
				{
					"label":   "ssh-key-3",
					"ssh_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC3...",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateSSHKeyMap(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected length %d, got %d", len(tt.expected), len(result))
				return
			}

			for i := range result {
				if result[i]["label"] != tt.expected[i]["label"] {
					t.Errorf("Expected label %s, got %s", tt.expected[i]["label"], result[i]["label"])
				}
				if result[i]["ssh_key"] != tt.expected[i]["ssh_key"] {
					t.Errorf("Expected ssh_key %s, got %s", tt.expected[i]["ssh_key"], result[i]["ssh_key"])
				}
			}
		})
	}
}
