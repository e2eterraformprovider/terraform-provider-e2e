package dbaas_mysql

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// expandVPCList converts a set of VPC IDs (integers) to VPC metadata for attach/detach operations
func expandVPCList(ctx context.Context, goe2eClient *goe2e.Client, vpcIDs []interface{}) ([]goe2e.VPCMetadata, error) {
	if len(vpcIDs) == 0 {
		return []goe2e.VPCMetadata{}, nil
	}

	// Convert integer IDs to strings
	var vpcIDStrings []string
	for _, id := range vpcIDs {
		vpcIDStrings = append(vpcIDStrings, strconv.Itoa(id.(int)))
	}

	// Use goe2e client to expand VPC list
	vpcMetaList, err := goe2eClient.DBaaSMySQL.ExpandVPCList(ctx, vpcIDStrings)
	if err != nil {
		return nil, err
	}

	return vpcMetaList, nil
}

// buildMySQLCreateRequest builds the create request from schema data
func buildMySQLCreateRequest(ctx context.Context, d *schema.ResourceData, goe2eClient *goe2e.Client, softwareID, templateID int) (*goe2e.MySQLClusterCreateRequest, error) {
	// Extract database configuration
	dbList := d.Get(tfconstants.AttrDatabase).([]interface{})
	if len(dbList) == 0 {
		return nil, fmt.Errorf(tfconstants.DatabaseConfigurationRequired)
	}
	dbMap := dbList[0].(map[string]interface{})

	// Build database config
	dbConfig := goe2e.DBConfig{
		User:        dbMap["user"].(string),
		Password:    dbMap["password"].(string),
		Name:        dbMap["name"].(string),
		DBaaSNumber: dbMap["dbaas_number"].(int),
	}

	// Build VPC list if provided
	var vpcList []goe2e.VPCMetadata
	if vpcSet, ok := d.GetOk(tfconstants.AttrVPCs); ok {
		vpcIDs := vpcSet.(*schema.Set).List()
		var err error
		vpcList, err = expandVPCList(ctx, goe2eClient, vpcIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to expand VPC list: %w", err)
		}
	}

	createReq := &goe2e.MySQLClusterCreateRequest{
		Name:             d.Get(tfconstants.AttrDBaaSName).(string),
		SoftwareID:       softwareID,
		TemplateID:       templateID,
		Group:            d.Get(tfconstants.AttrGroup).(string),
		Database:         dbConfig,
		PublicIPRequired: d.Get(tfconstants.AttrPublicIPRequired).(bool),
		Vpcs:             vpcList,
	}

	// Add parameter group if specified
	if pgID, ok := d.GetOk(tfconstants.AttrParameterGroupID); ok {
		createReq.ParameterGroupId = pgID.(int)
	}

	// Note: Encryption settings are handled by the API based on schema values
	// but are not part of the MySQLClusterCreateRequest struct

	return createReq, nil
}

// normalizeStatus converts API status values to user-friendly values
// For MySQL, we keep the API values as-is (SUSPENDED, RUNNING, etc.)
func normalizeStatus(apiStatus string) string {
	return apiStatus
}

// customImportStateFunc handles the import of MySQL DBaaS resources
// Format: project_id:dbaas_id
func customImportStateFunc(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), tfconstants.DBaaSImportIDSeparator)
	if len(parts) != 2 {
		return nil, fmt.Errorf(ImportIDInvalidFormatTemplate, tfconstants.DBaaSImportIDFormatDescription, d.Id())
	}

	projectID := parts[0]
	dbaasID := parts[1]

	// Set the project_id in the resource data
	if err := d.Set(tfconstants.AttrProjectID, projectID); err != nil {
		return nil, fmt.Errorf("failed to set project_id: %w", err)
	}

	// Set the resource ID to the dbaas_id
	d.SetId(dbaasID)

	return []*schema.ResourceData{d}, nil
}
