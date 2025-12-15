package security_group

// ============================================================================
// Security Group Resource Name Constant
// ============================================================================

// ResourceName is the resource type name for Security Group, used in error messages
const ResourceName = "Security Group"

// ============================================================================
// Security Group API Value Constants
// ============================================================================
// These constants represent valid API values for security group rules

// Rule type constants - API values
const (
	RuleTypeInbound  = "Inbound"
	RuleTypeOutbound = "Outbound"
)

// Network type constants - API values
const (
	NetworkTypeMyNetwork = "myNetwork"
	NetworkTypeManual    = "manual"
	NetworkTypeAny       = "any"
)

// Protocol constants - API values
const (
	ProtocolAll       = "All"
	ProtocolAllTCP    = "All_TCP"
	ProtocolAllUDP    = "All_UDP"
	ProtocolICMP      = "ICMP"
	ProtocolCustomTCP = "Custom_TCP"
	ProtocolCustomUDP = "Custom_UDP"
)

// ============================================================================
// Security Group Default Configuration Constants
// ============================================================================
// These are Terraform-specific default values, not API values

// Default values for security group rules
const (
	DefaultMyNetworkSize = 512
	DefaultNetworkCIDR   = "--"
	DefaultDescription   = ""
)

// ============================================================================
// Security Group-Specific Error Messages
// ============================================================================
// NOTE: For standard CRUD operations, use the generic templates from tfconstants:
// - tfconstants.ResourceOperationErrorTemplate (for create, with name)
// - tfconstants.ResourceOperationByIDErrorTemplate (for read/update/delete, with ID)
// - tfconstants.ResourceDataSourceListErrorTemplate (for data source listing)
// These are only defined here for security group-specific error scenarios

// ============================================================================

// ImportIDInvalidFormat is Security Group-specific import format error
const ImportIDInvalidFormat = "invalid ID format for import: expected 'project_id/region/sg_name', got: %s"

// DeleteDefaultSecurityGroup is the error message when trying to delete a default security group
const DeleteDefaultSecurityGroup = "Cannot delete default security group (ID: %s, name: %s). Unset is_default or default field before deleting."
