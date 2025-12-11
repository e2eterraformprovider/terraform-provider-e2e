package constants

// Common attribute names used across resources and data sources
// Following HashiCorp's terraform-provider-aws pattern for consistency
const (
	// Core identity and location attributes
	AttrRegion      = "region"
	AttrLocation    = "location"
	AttrProjectID   = "project_id"
	AttrProjectName = "project_name"
	AttrName        = "name"
	AttrID          = "id"

	// Common resource attributes
	AttrPlan = "plan"

	// Common metadata and state attributes
	AttrCreatedAt  = "created_at"
	AttrUpdatedAt  = "updated_at"
	AttrStatus     = "status"
	AttrLabel      = "label"
	AttrGroup      = "group"
	AttrTimestamp  = "timestamp"
	AttrSSHKey     = "ssh_key"
	AttrSSHKeys    = "ssh_keys"
	AttrSSHKeyIDs  = "ssh_key_ids"
	AttrPublicKey  = "public_key"   // V3+ preferred field name for SSH public key material
	AttrPK         = "pk"           // Primary key identifier (especially for SSH keys)
	AttrSSHKeyList = "ssh_key_list" // List of SSH keys in data source responses
	AttrTags       = "tags"
	AttrTagIDs     = "tag_ids"

	// Network and connectivity attributes
	AttrVPCID             = "vpc_id"
	AttrVPCs              = "vpcs"
	AttrPublicIPAddress   = "public_ip_address"
	AttrPrivateIPAddress  = "private_ip_address"
	AttrIPv6Address       = "ipv6_address"
	AttrPublicIPRequired  = "public_ip_required"
	AttrPublicIPEnabled   = "public_ip_enabled"
	AttrDefaultPublicIP   = "default_public_ip"
	AttrIPAddress         = "ip_address"
	AttrSecurityGroupIDs  = "security_group_ids"
	AttrNetworkInterface  = "network_interface"
	AttrAssignPublicIP    = "assign_public_ip"
	AttrEnableIPv6        = "enable_ipv6"
	AttrReserveIP         = "reserve_ip"
	AttrReserveIPID       = "reserve_ip_id"
	AttrReserveID         = "reserve_id"                 // Numeric ID from API
	AttrApplianceType     = "appliance_type"             // Infrastructure type
	AttrReservedType      = "reserved_type"              // Deprecated in V3
	AttrURN               = "urn"                        // Unique Resource Name
	AttrFloatingIPNodes   = "floating_ip_attached_nodes" // Attached nodes list
	AttrBoughtAt          = "bought_at"                  // When IP was purchased
	AttrCNID              = "cn_id"
	AttrLocalVPCNetworkID = "local_vpc_network_id"
	AttrPeerVPCNetworkID  = "peer_vpc_network_id"
	AttrIsPeerVPCExternal = "is_peer_vpc_external"
	AttrLocalTS           = "local_traffic_selector"
	AttrRemoteTS          = "remote_traffic_selector"
	AttrLocalGatewayIP    = "local_gateway_ip"
	AttrRemoteGatewayIP   = "remote_gateway_ip"

	// Storage attributes
	AttrSize     = "size"
	AttrDisk     = "disk"
	AttrDiskSize = "disk_size"
	AttrIOPS     = "iops"
	AttrVolumeID = "volume_id"
	AttrRootDisk = "root_disk"
	AttrSizeGB   = "size_gb"
	AttrDiskType = "disk_type"

	// Compute and instance attributes
	AttrVersion    = "version"
	AttrImage      = "image"
	AttrTemplateID = "template_id"
	AttrNodeID     = "node_id"
	AttrMemory     = "memory"
	AttrRAM        = "ram"
	AttrVCPU       = "vcpu"
	AttrVMID       = "vm_id"
	AttrVMName     = "vm_name"

	// Power and state management attributes
	AttrPowerStatus = "power_status"
	AttrLockNode    = "lock_node"
	AttrIsLocked    = "is_locked"
	AttrRebootNode  = "reboot_node"

	// Security attributes
	AttrIsEncryptionEnabled    = "is_encryption_enabled"
	AttrEncryptionPassphrase   = "encryption_passphrase"
	AttrSeverity               = "severity"
	AttrPreventVulnerabilities = "prevent_vulnerabilities"
	AttrEnableBitninja         = "enable_bitninja"
	AttrDisablePassword        = "disable_password"

	// Database attributes
	AttrDatabase         = "database"
	AttrDatabaseID       = "database_id"
	AttrDatabaseName     = "database_name"
	AttrDatabaseUser     = "database_user"
	AttrParameterGroupID = "parameter_group_id"

	// Scaling attributes
	AttrMinNodes = "min_nodes"
	AttrMaxNodes = "max_nodes"
	AttrMinVMs   = "min_vms"
	AttrMaxVMs   = "max_vms"
	AttrDesired  = "desired"

	// Load balancer attributes
	AttrBackends = "backends"
	AttrLbName   = "lb_name"
	AttrLbMode   = "lb_mode"
	AttrLbType   = "lb_type"
	AttrPlanName = "plan_name"

	// Kubernetes attributes
	AttrNodePools = "node_pools"
	AttrClusterID = "cluster_id"
	AttrPVSize    = "pv_size"
	AttrIsDynamic = "is_dynamic"

	// Object storage attributes
	AttrVersioningStatus             = "versioning_status"
	AttrLifecycleConfigurationStatus = "lifecycle_configuration_status"

	// CDN attributes
	AttrDomainName    = "domain_name"
	AttrDomainID      = "domain_id"
	AttrE2EDomainName = "e2e_domain_name"
	AttrSource        = "source"
	AttrIsEnabled     = "is_enabled"
	AttrState         = "state"
	AttrOriginDetails = "origin_details"
	AttrCacheDetails  = "cache_details"
	AttrDomainDetails = "domain_details"

	// CDN Origin Details attributes
	AttrPath                   = "path"
	AttrSSLProtocol            = "ssl_protocol"
	AttrProtocolPolicy         = "protocol_policy"
	AttrOriginReadTimeout      = "origin_read_timeout"
	AttrOriginKeepaliveTimeout = "origin_keepalive_timeout"

	// CDN Cache Details attributes
	AttrViewerProtocolPolicy = "viewer_protocol_policy"
	AttrAllowedHTTPMethods   = "allowed_http_methods"
	AttrDefaultTTL           = "default_ttl"
	AttrMinTTL               = "min_ttl"
	AttrMaxTTL               = "max_ttl"

	// CDN Domain Details attributes
	AttrHTTPVersions = "http_versions"
	AttrRootObject   = "root_object"
	AttrIPv6Enabled  = "ipv6_enabled"

	// Backup attributes
	AttrBackupConfig       = "backup_config"
	AttrBackupStatus       = "backup_status"
	AttrBackupPlanID       = "plan_id"
	AttrBackupType         = "type"
	AttrBackupExcludePaths = "exclude_paths"
	AttrBackupNow          = "backup_now"
	AttrCompressionType    = "compression_type"
	AttrCompressionLevel   = "compression_level"
	AttrEncryptionEnabled  = "encryption_enabled"
	AttrEncryptionKey      = "encryption_passphrase"
	AttrHoursOfDay         = "hours_of_day"
	AttrStartingMinute     = "starting_minute"
	AttrDBEnabled          = "db_enabled"
	AttrDBUsername         = "db_username"
	AttrDBPassword         = "db_password"
	AttrBackupStatusDetail = "detail"
	AttrLastRecoveryPoint  = "last_recovery_point"
)
