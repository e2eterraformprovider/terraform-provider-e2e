package dbaas_mariadb_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/dbaas_mariadb"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// Data Source Read Operation Tests
// ============================================================================

func TestDataSourceReadMariaDB_Success(t *testing.T) {
	mockService := &dbaas_mariadb.MockMariaDBService{
		GetMariaDBFunc: func(ctx context.Context, id string) (*goe2e.MariaDB, *goe2e.Response, error) {
			return &goe2e.MariaDB{
				ID:     12345,
				Name:   "test-mariadb",
				Status: goe2econstants.DBaaSStatusRunning,
				MasterNode: goe2e.DBNode{
					Database: goe2e.DBCreds{
						ID:       1,
						Database: "testdb",
						Username: "testuser",
						PGDetail: goe2e.PGDetail{
							ID: 5,
						},
					},
					Plan: goe2e.Plan{
						Name: "DBS.16GB",
						Software: goe2e.Software{
							Version: "10.6",
						},
					},
					PublicIPAddress:  "1.2.3.4",
					PrivateIPAddress: "10.0.0.1",
					Disk:             "100",
					Status:           "Running",
				},
			}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createTestConfigWithMock(mockService)
	resource := dbaas_mariadb.DataSourceMariaDB()

	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrID: "12345",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)
	require.False(t, diags.HasError(), "Read should succeed")
	assert.Equal(t, "12345", d.Id())
	assert.Equal(t, "test-mariadb", d.Get(tfconstants.AttrName))
	assert.Equal(t, 1, d.Get(tfconstants.AttrDatabaseID))
	assert.Equal(t, "testdb", d.Get(tfconstants.AttrDatabaseName))
	assert.Equal(t, "testuser", d.Get(tfconstants.AttrDatabaseUser))
}

func TestDataSourceReadMariaDB_InvalidID(t *testing.T) {
	mockService := &dbaas_mariadb.MockMariaDBService{
		GetMariaDBFunc: func(ctx context.Context, id string) (*goe2e.MariaDB, *goe2e.Response, error) {
			return nil, nil, errors.New("cluster not found")
		},
	}

	cfg := createTestConfigWithMock(mockService)
	resource := dbaas_mariadb.DataSourceMariaDB()

	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrID: "99999",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)
	require.True(t, diags.HasError(), "Read should fail with invalid ID")
}

func TestDataSourceReadMariaDB_StatusNormalization(t *testing.T) {
	mockService := &dbaas_mariadb.MockMariaDBService{
		GetMariaDBFunc: func(ctx context.Context, id string) (*goe2e.MariaDB, *goe2e.Response, error) {
			return &goe2e.MariaDB{
				ID:     12345,
				Name:   "test-mariadb",
				Status: goe2econstants.DBaaSStatusSuspended, // Should be normalized to STOPPED
				MasterNode: goe2e.DBNode{
					Plan: goe2e.Plan{
						Name: "DBS.16GB",
						Software: goe2e.Software{
							Version: "10.6",
						},
					},
					Database: goe2e.DBCreds{
						ID: 1,
					},
				},
			}, &goe2e.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}

	cfg := createTestConfigWithMock(mockService)
	resource := dbaas_mariadb.DataSourceMariaDB()

	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrID: "12345",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)
	require.False(t, diags.HasError(), "Read should succeed")
	// Status should be normalized from SUSPENDED to STOPPED
	assert.Equal(t, goe2econstants.DBaaSStatusStopped, d.Get(tfconstants.AttrStatus))
}

func TestDataSourceReadMariaDB_NilResponse(t *testing.T) {
	mockService := &dbaas_mariadb.MockMariaDBService{
		GetMariaDBFunc: func(ctx context.Context, id string) (*goe2e.MariaDB, *goe2e.Response, error) {
			return nil, nil, nil
		},
	}

	cfg := createTestConfigWithMock(mockService)
	resource := dbaas_mariadb.DataSourceMariaDB()

	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrID: "12345",
	})

	diags := resource.ReadContext(context.Background(), d, cfg)
	require.True(t, diags.HasError(), "Read should fail with nil response")
}
