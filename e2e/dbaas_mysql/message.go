package dbaas_mysql

// ResourceName is the human-readable resource name used in shared error templates.
const ResourceName = "MySQL DBaaS"

const (
	ClusterNotFoundTemplate       = "MySQL DBaaS instance (ID: %s) not found"
	ImportIDInvalidFormatTemplate = "invalid ID format: expected %s, got: %s"

	// Operation-specific error messages
	ErrorRetrievingSoftwareIDTemplate    = "error retrieving %s software ID for version (%s) in project (%s), region (%s): %s"
	ErrorRetrievingTemplateIDTemplate    = "error retrieving %s template ID for plan (%s) in project (%s), region (%s): %s"
	ErrorPreparingVPCListTemplate        = "error preparing VPC list for MySQL DBaaS (ID: %s): %s"
	ErrorStoppingTemplate                = "error stopping MySQL DBaaS (ID: %s): %s"
	ErrorStartingTemplate                = "error starting MySQL DBaaS (ID: %s): %s"
	ErrorRestartingTemplate              = "error restarting MySQL DBaaS (ID: %s): %s"
	ErrorDetachingVPCTemplate            = "error detaching VPC from MySQL DBaaS (ID: %s): %s"
	ErrorAttachingVPCTemplate            = "error attaching VPC to MySQL DBaaS (ID: %s): %s"
	ErrorAttachingPublicIPTemplate       = "error attaching public IP to MySQL DBaaS (ID: %s): %s"
	ErrorDetachingPublicIPTemplate       = "error detaching public IP from MySQL DBaaS (ID: %s): %s"
	ErrorDetachingParameterGroupTemplate = "error detaching parameter group (ID: %d) from MySQL DBaaS (ID: %s): %s"
	ErrorAttachingParameterGroupTemplate = "error attaching parameter group (ID: %d) to MySQL DBaaS (ID: %s): %s"
	ErrorUpgradingPlanTemplate           = "error upgrading MySQL DBaaS (ID: %s) plan from (%s) to (%s): %s"
	ErrorExpandingDiskTemplate           = "error expanding MySQL DBaaS (ID: %s) disk by %d GB: %s"
	ErrorRetrievingSoftwareIDForUpgrade  = "error retrieving MySQL software ID for version (%s) while upgrading plan for DBaaS (ID: %s): %s"
	ErrorRetrievingTemplateIDForUpgrade  = "error retrieving MySQL template ID for plan (%s) while upgrading DBaaS (ID: %s): %s"
	ErrorUnsupportedStatusTemplate       = "error updating MySQL DBaaS (ID: %s): unsupported status value: %s. Must be one of: SUSPENDED, STOPPED, RUNNING, START, RESTARTING, RESTART"
	ErrorCannotUpgradePlanTemplate       = "cannot upgrade plan for MySQL DBaaS (ID: %s): database must be in SUSPENDED/STOPPED state (current state: %s). Please stop the instance first"
)
