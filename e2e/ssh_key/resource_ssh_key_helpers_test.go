package ssh_key

import (
	"testing"

	e2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestGetKeyName tests the getKeyName helper function
func TestGetKeyName(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]interface{}
		expected string
	}{
		{
			name: "prefers_name_over_label",
			data: map[string]interface{}{
				e2econstants.AttrName:  "my-key",
				e2econstants.AttrLabel: "old-label",
			},
			expected: "my-key",
		},
		{
			name: "falls_back_to_label_when_name_missing",
			data: map[string]interface{}{
				e2econstants.AttrLabel: "my-label",
			},
			expected: "my-label",
		},
		{
			name: "returns_name_when_only_name_set",
			data: map[string]interface{}{
				e2econstants.AttrName: "preferred-name",
			},
			expected: "preferred-name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a ResourceData with test values
			d := schema.TestResourceDataRaw(t, ResourceSshKey().Schema, tt.data)

			result := getKeyName(d)
			if result != tt.expected {
				t.Errorf("getKeyName() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// TestGetPublicKey tests the getPublicKey helper function
func TestGetPublicKey(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]interface{}
		expected string
	}{
		{
			name: "prefers_public_key_over_ssh_key",
			data: map[string]interface{}{
				e2econstants.AttrPublicKey: "ssh-rsa AAAAB3NzaC1yc2EAAAA...",
				e2econstants.AttrSSHKey:    "ssh-rsa AAAAB3NzaC1yc2EAAAA_old...",
			},
			expected: "ssh-rsa AAAAB3NzaC1yc2EAAAA...",
		},
		{
			name: "falls_back_to_ssh_key_when_public_key_missing",
			data: map[string]interface{}{
				e2econstants.AttrSSHKey: "ssh-rsa AAAAB3NzaC1yc2EAAAA...",
			},
			expected: "ssh-rsa AAAAB3NzaC1yc2EAAAA...",
		},
		{
			name: "returns_public_key_when_only_public_key_set",
			data: map[string]interface{}{
				e2econstants.AttrPublicKey: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAA...",
			},
			expected: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAA...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, ResourceSshKey().Schema, tt.data)

			result := getPublicKey(d)
			if result != tt.expected {
				t.Errorf("getPublicKey() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// TestSetKeyName tests the setKeyName helper function
func TestSetKeyName(t *testing.T) {
	d := schema.TestResourceDataRaw(t, ResourceSshKey().Schema, map[string]interface{}{})

	keyName := "test-key"
	if err := setKeyName(d, keyName); err != nil {
		t.Fatalf("setKeyName() error = %v", err)
	}

	// Verify both fields are set
	if name, ok := d.GetOk(e2econstants.AttrName); !ok || name != keyName {
		t.Errorf("name field not set correctly: got %v", name)
	}

	if label, ok := d.GetOk(e2econstants.AttrLabel); !ok || label != keyName {
		t.Errorf("label field not set correctly: got %v", label)
	}
}

// TestSetPublicKey tests the setPublicKey helper function
func TestSetPublicKey(t *testing.T) {
	d := schema.TestResourceDataRaw(t, ResourceSshKey().Schema, map[string]interface{}{})

	pubKey := "ssh-rsa AAAAB3NzaC1yc2EAAAA..."
	if err := setPublicKey(d, pubKey); err != nil {
		t.Fatalf("setPublicKey() error = %v", err)
	}

	// Verify both fields are set
	if pk, ok := d.GetOk(e2econstants.AttrPublicKey); !ok || pk != pubKey {
		t.Errorf("public_key field not set correctly: got %v", pk)
	}

	if sshKey, ok := d.GetOk(e2econstants.AttrSSHKey); !ok || sshKey != pubKey {
		t.Errorf("ssh_key field not set correctly: got %v", sshKey)
	}
}
