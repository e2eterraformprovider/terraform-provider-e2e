package vpc

// Error message constants for VPC resource operations
// These constants are used for consistent error messaging across VPC resource operations.
const (
	// Error messages for CRUD operations
	errVPCRetrieve = "Error retrieving VPC (ID: %s) in project (%s), region (%s): %s"
	errVPCCreate   = "Error creating VPC (name: %s) in project (%s), region (%s): %s"
	errVPCDelete   = "Error deleting VPC (ID: %s) in project (%s), region (%s): %s"

	// Error messages for import operations
	errVPCImportFormat = "invalid import format: expected 'vpc_id' or 'project_id:vpc_id', got '%s'"

	// Error messages for data source operations
	errVPCListEmpty = "VPC list is empty in the response"
)
