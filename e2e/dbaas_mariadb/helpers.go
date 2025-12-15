package dbaas_mariadb

import (
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
)

// normalizeStatus converts API status values to user-friendly values.
// For MariaDB, SUSPENDED is converted to STOPPED, all other values pass through unchanged.
func normalizeStatus(apiStatus string) string {
	if apiStatus == goe2econstants.DBaaSStatusSuspended {
		return goe2econstants.DBaaSStatusStopped
	}
	return apiStatus
}
