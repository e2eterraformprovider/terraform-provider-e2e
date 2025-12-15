package dbaas_postgress

// ResourceName is the human-readable resource name used in shared error templates.
const ResourceName = "PostgreSQL DBaaS"

const (
	ClusterIDRequiredForUpdate    = "cluster ID is required for update"
	ClusterIDRequiredForDeletion  = "cluster ID is required for deletion"
	ClusterNotFoundTemplate       = "PostgreSQL cluster %s not found"
	ImportIDInvalidFormatTemplate = "invalid ID format: expected %s, got %s"

	// Operation-specific error messages
	ErrorRetrievingSoftwareIDTemplate       = "error retrieving PostgreSQL software ID for version (%s) in project (%s), region (%s): %s"
	ErrorRetrievingTemplateIDTemplate       = "error retrieving PostgreSQL template ID for plan (%s) in project (%s), region (%s): %s"
	ErrorPreparingVPCListTemplate           = "error preparing VPC list for PostgreSQL DBaaS (ID: %s) in project (%s), region (%s): %s"
	ErrorStoppingTemplate                   = "error stopping PostgreSQL DBaaS (ID: %s) in project (%s), region (%s): %s"
	ErrorStartingTemplate                   = "error starting PostgreSQL DBaaS (ID: %s) in project (%s), region (%s): %s"
	ErrorRestartingTemplate                 = "error restarting PostgreSQL DBaaS (ID: %s) in project (%s), region (%s): %s"
	ErrorDetachingVPCTemplate               = "error detaching VPC from PostgreSQL DBaaS (ID: %s) in project (%s), region (%s): %s"
	ErrorAttachingVPCTemplate               = "error attaching VPC to PostgreSQL DBaaS (ID: %s) in project (%s), region (%s): %s"
	ErrorAttachingPublicIPTemplate          = "error attaching public IP to PostgreSQL DBaaS (ID: %s) in project (%s), region (%s): %s"
	ErrorDetachingPublicIPTemplate          = "error detaching public IP from PostgreSQL DBaaS (ID: %s) in project (%s), region (%s): %s"
	ErrorDetachingParameterGroupTemplate    = "error detaching parameter group (ID: %d) from PostgreSQL DBaaS (ID: %s) in project (%s), region (%s): %s"
	ErrorAttachingParameterGroupTemplate    = "error attaching parameter group (ID: %d) to PostgreSQL DBaaS (ID: %s) in project (%s), region (%s): %s"
	ErrorUpgradingPlanTemplate              = "error upgrading PostgreSQL DBaaS (ID: %s) plan from (%s) to (%s) in project (%s), region (%s): %s"
	ErrorExpandingDiskTemplate              = "error expanding PostgreSQL DBaaS (ID: %s) disk by %d GB in project (%s), region (%s): %s"
	ErrorRetrievingSoftwareIDForUpgrade     = "error retrieving PostgreSQL software ID for version (%s) while upgrading plan for DBaaS (ID: %s) in project (%s), region (%s): %s"
	ErrorRetrievingTemplateIDForUpgrade     = "error retrieving PostgreSQL template ID for plan (%s) while upgrading DBaaS (ID: %s) in project (%s), region (%s): %s"
	ErrorUnsupportedStatusTemplate          = "error updating PostgreSQL DBaaS (ID: %s): unsupported status value: %s. Must be one of: STOPPED, SUSPENDED, RUNNING, RESTARTING"
	ErrorCannotUpgradePlanTemplate          = "cannot upgrade plan for PostgreSQL DBaaS (ID: %s): database must be in SUSPENDED state (current state: %s) in project (%s), region (%s). Please stop the instance first"
	ErrorCannotExpandDiskTemplate           = "cannot expand disk for PostgreSQL DBaaS (ID: %s): database must be in SUSPENDED state (current state: %s) in project (%s), region (%s)"
	ErrorCannotPerformPowerOpsTemplate      = "cannot perform power operations on PostgreSQL DBaaS (ID: %s): database is in CREATING state in project (%s), region (%s). Please wait for database creation to complete"
	ErrorCannotUpdatePublicIPTemplate       = "cannot update public IP for PostgreSQL DBaaS (ID: %s): database is in CREATING state in project (%s), region (%s). Please wait for database creation to complete"
	ErrorCannotUpdateVPCListTemplate        = "cannot update VPC list for PostgreSQL DBaaS (ID: %s): database is in CREATING state in project (%s), region (%s). Please wait for database creation to complete"
	ErrorCannotUpdateParameterGroupTemplate = "cannot update parameter group for PostgreSQL DBaaS (ID: %s): database is in CREATING state in project (%s), region (%s). Please wait for database creation to complete"
	ErrorExpandingDiskInvalidTypeTemplate   = "error expanding disk for PostgreSQL DBaaS (ID: %s): size must be an integer, got %T"
	ErrorExpandingDiskInvalidSizeTemplate   = "error expanding disk for PostgreSQL DBaaS (ID: %s): size must be greater than previous size (%d GB). Got: %d GB"
	ErrorProjectIDImmutable                 = "Cannot update project_id: this field is immutable after PostgreSQL DBaaS creation"
)
