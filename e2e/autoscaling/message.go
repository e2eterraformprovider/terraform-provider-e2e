package autoscaling

// ResourceName is used in autoscaling-specific error messages.
const ResourceName = "Scaler Group"

// CustomizeDiff / schema validation messages.
// Keep these stable because acceptance tests and user experience depend on them.
const (
	// Required V2/V3 "pair" messages
	ErrEitherVMImageOrImageRequired           = "either 'vm_image_name' or 'image' must be specified"
	ErrEitherMinNodesOrMinSizeRequired        = "either 'min_nodes' or 'min_size' must be specified"
	ErrEitherMaxNodesOrMaxSizeRequired        = "either 'max_nodes' or 'max_size' must be specified"
	ErrEitherDesiredOrDesiredCapacityRequired = "either 'desired' or 'desired_capacity' must be specified"

	// State requirement messages (provider-side invariants)
	ErrSecurityGroupUpdatesRequireRunningFmt = "security group updates require scaler group to be in 'Running' state, current: %s"
	ErrVPCUpdatesRequireStoppedFmt           = "VPC updates require scaler group to be in 'Stopped' state, current: %s"
	ErrPublicIPUpdatesRequireStoppedFmt      = "public IP updates require scaler group to be in 'Stopped' state, current: %s"
	ErrPublicIPUpdatesRequireVPC             = "public IP updates require at least one VPC to be attached"

	// Scheduled action conditional validation (schema can't express "required if")
	ErrScheduledActionTargetCapacityRequiredFmt = "%s[%d]: %s must be set (> 0) when %s is %q"
	ErrScheduledActionAdjustmentRequiredFmt     = "%s[%d]: %s must be set (non-zero) when %s is %q"

	// Deprecation warnings (CustomizeDiff)
	WarnDeprecatedV2FieldsHeader = "Deprecated V2 autoscaling fields detected. Please migrate to V3 field names/blocks to avoid breaking changes in v4.0."
	WarnDeprecatedV2FieldsFooter = "Migration notes: (1) `status` is lowercase in V3 (`running`/`stopped`), (2) prefer structured blocks: `scaling_policy`, `scheduled_action`, `vpc_config` / `network_config`."
)
