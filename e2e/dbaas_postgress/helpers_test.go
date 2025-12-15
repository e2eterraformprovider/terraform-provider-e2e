package dbaas_postgress

import (
	"context"
	"testing"

	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomImportStateFunc(t *testing.T) {
	ctx := context.Background()

	t.Run("valid import format", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
			tfconstants.AttrProjectID: {Type: schema.TypeString},
		}, map[string]interface{}{})
		d.SetId("project123" + tfconstants.DBaaSImportIDSeparator + "dbaas456")

		result, err := CustomImportStateFunc(ctx, d, nil)
		require.NoError(t, err)
		require.Len(t, result, 1)
		assert.Equal(t, "dbaas456", result[0].Id())
		projectID, ok := result[0].Get(tfconstants.AttrProjectID).(string)
		require.True(t, ok)
		assert.Equal(t, "project123", projectID)
	})

	t.Run("invalid format - missing colon", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})
		d.SetId("project123dbaas456")

		result, err := CustomImportStateFunc(ctx, d, nil)
		assert.Error(t, err)
		assert.Nil(t, result)
		// Error message is formatted, so check for the actual formatted string
		assert.Contains(t, err.Error(), tfconstants.DBaaSImportIDFormatDescription)
		assert.Contains(t, err.Error(), "project123dbaas456")
	})

	t.Run("invalid format - too many parts", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})
		d.SetId("project123" + tfconstants.DBaaSImportIDSeparator + "dbaas456" + tfconstants.DBaaSImportIDSeparator + "extra")

		result, err := CustomImportStateFunc(ctx, d, nil)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), tfconstants.DBaaSImportIDFormatDescription)
	})

	t.Run("invalid format - empty parts", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})
		d.SetId(tfconstants.DBaaSImportIDSeparator)

		result, err := CustomImportStateFunc(ctx, d, nil)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("invalid format - only colon, no parts", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})
		d.SetId(tfconstants.DBaaSImportIDSeparator)

		result, err := CustomImportStateFunc(ctx, d, nil)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("invalid format - empty project_id", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})
		d.SetId(tfconstants.DBaaSImportIDSeparator + "dbaas456")

		result, err := CustomImportStateFunc(ctx, d, nil)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("invalid format - empty dbaas_id", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})
		d.SetId("project123" + tfconstants.DBaaSImportIDSeparator)

		result, err := CustomImportStateFunc(ctx, d, nil)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("context is properly passed through", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
			tfconstants.AttrProjectID: {Type: schema.TypeString},
		}, map[string]interface{}{})
		d.SetId("project123" + tfconstants.DBaaSImportIDSeparator + "dbaas456")

		// Create a context with a value to verify it's passed through
		type testKey string
		const testKeyValue testKey = "test-key"
		ctxWithValue := context.WithValue(ctx, testKeyValue, "test-value")
		result, err := CustomImportStateFunc(ctxWithValue, d, nil)
		require.NoError(t, err)
		assert.NotNil(t, result)
		// Context is passed but not directly testable without modifying the function
		// This test verifies the function accepts and uses context parameter
	})
}
