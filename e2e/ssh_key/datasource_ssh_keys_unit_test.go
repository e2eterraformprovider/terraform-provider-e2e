package ssh_key

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// Test: dataSourceReadSshKeys
// ============================================================================

func TestDataSourceReadSshKeys_Success(t *testing.T) {
	mockService := &mockSSHKeyService{
		listSSHKeysFunc: func(ctx context.Context) ([]goe2e.SSHKey, *goe2e.Response, error) {
			return []goe2e.SSHKey{
				{
					PK:        12345,
					Label:     "test-key-1",
					SSHKey:    "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...",
					Timestamp: "2024-01-01T00:00:00Z",
				},
				{
					PK:        67890,
					Label:     "test-key-2",
					SSHKey:    "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAA...",
					Timestamp: "2024-01-02T00:00:00Z",
				},
			}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createMockConfigForProject(t, mockService, "test-project", "us-east-1")
	resource := DataSourceSshKeys()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.False(t, diags.HasError(), "Read should succeed")
	assert.Equal(t, tfconstants.AttrSSHKeyList, d.Id())

	sshKeyList := d.Get(tfconstants.AttrSSHKeyList).([]interface{})
	require.Len(t, sshKeyList, 2, "Should have 2 SSH keys")

	// Validate first key
	key1 := sshKeyList[0].(map[string]interface{})
	assert.Equal(t, 12345, key1[tfconstants.AttrPK])
	assert.Equal(t, "test-key-1", key1[tfconstants.AttrLabel])
	assert.Equal(t, "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...", key1[tfconstants.AttrSSHKey])
	assert.Equal(t, "2024-01-01T00:00:00Z", key1[tfconstants.AttrCreatedAt])

	// Validate second key
	key2 := sshKeyList[1].(map[string]interface{})
	assert.Equal(t, 67890, key2[tfconstants.AttrPK])
	assert.Equal(t, "test-key-2", key2[tfconstants.AttrLabel])
}

func TestDataSourceReadSshKeys_EmptyList(t *testing.T) {
	mockService := &mockSSHKeyService{
		listSSHKeysFunc: func(ctx context.Context) ([]goe2e.SSHKey, *goe2e.Response, error) {
			return []goe2e.SSHKey{}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createMockConfigForProject(t, mockService, "test-project", "us-east-1")
	resource := DataSourceSshKeys()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.False(t, diags.HasError(), "Read should succeed with empty list")
	assert.Equal(t, tfconstants.AttrSSHKeyList, d.Id())

	sshKeyList := d.Get(tfconstants.AttrSSHKeyList).([]interface{})
	assert.Len(t, sshKeyList, 0, "Should have empty list")
}

func TestDataSourceReadSshKeys_APIError(t *testing.T) {
	mockService := &mockSSHKeyService{
		listSSHKeysFunc: func(ctx context.Context) ([]goe2e.SSHKey, *goe2e.Response, error) {
			return nil, nil, errors.New("API error: failed to list SSH keys")
		},
	}

	cfg := createMockConfigForProject(t, mockService, "test-project", "us-east-1")
	resource := DataSourceSshKeys()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Read should fail on API error")
	// diag.Errorf sets the Summary field, not Detail
	assert.Contains(t, diags[0].Summary, "error finding ssh keys")
}

// ============================================================================
// Test: flattenSshKeys
// ============================================================================

func TestFlattenSshKeys(t *testing.T) {
	tests := []struct {
		name           string
		sshKeys        []goe2e.SSHKey
		expectedLength int
		validateFunc   func(*testing.T, []interface{})
	}{
		{
			name:           "nil input - returns empty slice",
			sshKeys:        nil,
			expectedLength: 0,
			validateFunc: func(t *testing.T, result []interface{}) {
				assert.Len(t, result, 0)
			},
		},
		{
			name:           "empty slice - returns empty slice",
			sshKeys:        []goe2e.SSHKey{},
			expectedLength: 0,
			validateFunc: func(t *testing.T, result []interface{}) {
				assert.Len(t, result, 0)
			},
		},
		{
			name: "single SSH key - all fields present",
			sshKeys: []goe2e.SSHKey{
				{
					PK:        12345,
					Label:     "test-key-1",
					SSHKey:    "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...",
					Timestamp: "2024-01-01T00:00:00Z",
				},
			},
			expectedLength: 1,
			validateFunc: func(t *testing.T, result []interface{}) {
				require.Len(t, result, 1)
				keyMap := result[0].(map[string]interface{})
				assert.Equal(t, 12345, keyMap[tfconstants.AttrPK])
				assert.Equal(t, "test-key-1", keyMap[tfconstants.AttrLabel])
				assert.Equal(t, "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...", keyMap[tfconstants.AttrSSHKey])
				assert.Equal(t, "2024-01-01T00:00:00Z", keyMap[tfconstants.AttrCreatedAt])
			},
		},
		{
			name: "multiple SSH keys",
			sshKeys: []goe2e.SSHKey{
				{
					PK:        12345,
					Label:     "test-key-1",
					SSHKey:    "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...",
					Timestamp: "2024-01-01T00:00:00Z",
				},
				{
					PK:        67890,
					Label:     "test-key-2",
					SSHKey:    "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAA...",
					Timestamp: "2024-01-02T00:00:00Z",
				},
			},
			expectedLength: 2,
			validateFunc: func(t *testing.T, result []interface{}) {
				require.Len(t, result, 2)
				key1 := result[0].(map[string]interface{})
				assert.Equal(t, "test-key-1", key1[tfconstants.AttrLabel])
				key2 := result[1].(map[string]interface{})
				assert.Equal(t, "test-key-2", key2[tfconstants.AttrLabel])
			},
		},
		{
			name: "SSH key with empty fields",
			sshKeys: []goe2e.SSHKey{
				{
					PK:        99999,
					Label:     "test-key-empty",
					SSHKey:    "",
					Timestamp: "",
				},
			},
			expectedLength: 1,
			validateFunc: func(t *testing.T, result []interface{}) {
				require.Len(t, result, 1)
				keyMap := result[0].(map[string]interface{})
				assert.Equal(t, 99999, keyMap[tfconstants.AttrPK])
				assert.Equal(t, "test-key-empty", keyMap[tfconstants.AttrLabel])
				assert.Equal(t, "", keyMap[tfconstants.AttrSSHKey])
				assert.Equal(t, "", keyMap[tfconstants.AttrCreatedAt])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := flattenSshKeys(tt.sshKeys)

			assert.Len(t, result, tt.expectedLength)
			if tt.validateFunc != nil {
				tt.validateFunc(t, result)
			}
		})
	}
}

// Helper function for datasources that use Goe2eClientForProject
func createMockConfigForProject(t *testing.T, mockService *mockSSHKeyService, defaultProjectID, defaultRegion string) *config.Config {
	cfg, err := config.NewConfig("test-api-key-12345", "test-auth-token-12345", "https://api.e2enetworks.com/myaccount/api/v1/")
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	cfg.DefaultProjectID = defaultProjectID
	cfg.DefaultRegion = defaultRegion

	mockClient := &goe2e.Client{}
	mockClient.SSHKeys = mockService
	cfg.SetGoe2eClientForTesting(mockClient)

	return cfg
}
