package reserve_ip

import (
	"fmt"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
)

// ============================================================================
// File-local constants for Reserve IP URN and Import formats
// ============================================================================

const (
	// reserveIPURNFormat is the format string for generating Reserve IP URNs
	// Format: e2e:reserve_ip:<region>:<ip_address>
	reserveIPURNFormat = "e2e:reserve_ip:%s:%s"

	// reserveIPImportFormatError is the error message template for invalid import IDs
	reserveIPImportFormatError = "invalid import ID format, expected: project_id/region/ip_address, got: %s"
)

// generateReserveIPURN generates a URN for a reserved IP in the format:
// e2e:reserve_ip:<region>:<ip_address>
func generateReserveIPURN(region, ipAddress string) string {
	return fmt.Sprintf(reserveIPURNFormat, region, ipAddress)
}

// parseReserveIPImportID parses an import ID string in the format:
// project_id/region/ip_address
// Returns projectID, region, ipAddress, and error
func parseReserveIPImportID(id string) (projectID, region, ipAddress string, err error) {
	parts := strings.Split(id, "/")
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf(reserveIPImportFormatError, id)
	}
	return parts[0], parts[1], parts[2], nil
}

// flattenFloatingIPAttachedNodes converts a slice of FloatingIPAttachedNode to
// a slice of maps suitable for Terraform state
func flattenFloatingIPAttachedNodes(nodes []goe2e.FloatingIPAttachedNode) []map[string]interface{} {
	if nodes == nil {
		return []map[string]interface{}{}
	}

	result := make([]map[string]interface{}, len(nodes))
	for i, node := range nodes {
		result[i] = map[string]interface{}{
			"id":                    node.ID,
			"name":                  node.Name,
			"vm_id":                 node.VMID,
			"ip_address_public":     node.IPAddressPublic,
			"ip_address_private":    node.IPAddressPrivate,
			"status_name":           node.StatusName,
			"security_group_status": node.SecurityGroupStatus,
		}
	}
	return result
}
