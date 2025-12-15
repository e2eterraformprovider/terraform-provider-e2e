package ssh_key

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// Test: dataSourceReadSshKey
// ============================================================================

func TestDataSourceReadSshKey_Success(t *testing.T) {
	mockService := &mockSSHKeyService{
		getSSHKeyByLabelFunc: func(ctx context.Context, label string) (*goe2e.SSHKey, *goe2e.Response, error) {
			assert.Equal(t, "test-key", label)
			return &goe2e.SSHKey{
				PK:        12345,
				Label:     "test-key",
				SSHKey:    "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...",
				Timestamp: "2024-01-01T00:00:00Z",
			}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createMockConfig(t, mockService, "test-project", "us-east-1")
	resource := DataSourceSshKey()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrLabel: "test-key",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.False(t, diags.HasError(), "Read should succeed")
	assert.Equal(t, "12345", d.Id())
	assert.Equal(t, "test-key", d.Get(tfconstants.AttrName))
	assert.Equal(t, "test-key", d.Get(tfconstants.AttrLabel))
	assert.Equal(t, "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...", d.Get(tfconstants.AttrPublicKey))
	assert.Equal(t, "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...", d.Get(tfconstants.AttrSSHKey))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get(tfconstants.AttrCreatedAt))
}

func TestDataSourceReadSshKey_NotFound(t *testing.T) {
	mockService := &mockSSHKeyService{
		getSSHKeyByLabelFunc: func(ctx context.Context, label string) (*goe2e.SSHKey, *goe2e.Response, error) {
			return nil, nil, fmt.Errorf("SSH key with label %s not found", label)
		},
	}

	cfg := createMockConfig(t, mockService, "test-project", "us-east-1")
	resource := DataSourceSshKey()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrLabel: "non-existent-key",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Read should fail when key not found")
	errorMsg := diags[0].Summary + " " + diags[0].Detail
	assert.Contains(t, errorMsg, "failed to find SSH key")
}

func TestDataSourceReadSshKey_EmptyLabel(t *testing.T) {
	cfg := createMockConfig(t, &mockSSHKeyService{}, "test-project", "us-east-1")
	resource := DataSourceSshKey()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrLabel: "",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Read should fail with empty label")
	errorMsg := diags[0].Summary + " " + diags[0].Detail
	assert.Contains(t, errorMsg, "required")
}

func TestDataSourceReadSshKey_APIError(t *testing.T) {
	mockService := &mockSSHKeyService{
		getSSHKeyByLabelFunc: func(ctx context.Context, label string) (*goe2e.SSHKey, *goe2e.Response, error) {
			return nil, nil, errors.New("API error: failed to fetch SSH key")
		},
	}

	cfg := createMockConfig(t, mockService, "test-project", "us-east-1")
	resource := DataSourceSshKey()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrLabel: "test-key",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)

	require.True(t, diags.HasError(), "Read should fail on API error")
	errorMsg := diags[0].Summary + " " + diags[0].Detail
	assert.Contains(t, errorMsg, "failed to find SSH key")
}
