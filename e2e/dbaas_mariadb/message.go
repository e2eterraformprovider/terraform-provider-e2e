package dbaas_mariadb

// ResourceName is the human-readable resource name used in shared error templates.
const ResourceName = "MariaDB DBaaS"

const (
	ClusterNotFoundTemplate = "MariaDB cluster %s not found"

	// Operation-specific error messages
	ErrorRetrievingSoftwareIDTemplate    = "error retrieving %s software ID for version (%s) in project (%s), region (%s): %s"
	ErrorRetrievingTemplateIDTemplate    = "error retrieving %s template ID for plan (%s) in project (%s), region (%s): %s"
	ErrorPreparingVPCListTemplate        = "error preparing VPC list for MariaDB DBaaS in project (%s), region (%s): %s"
	ErrorStoppingTemplate                = "error stopping MariaDB DBaaS (ID: %s): %s"
	ErrorStartingTemplate                = "error starting MariaDB DBaaS (ID: %s): %s"
	ErrorRestartingTemplate              = "error restarting MariaDB DBaaS (ID: %s): %s"
	ErrorDetachingVPCTemplate            = "error detaching VPC from MariaDB DBaaS (ID: %s): %s"
	ErrorAttachingVPCTemplate            = "error attaching VPC to MariaDB DBaaS (ID: %s): %s"
	ErrorAttachingPublicIPTemplate       = "error attaching public IP to MariaDB DBaaS (ID: %s): %s"
	ErrorDetachingPublicIPTemplate       = "error detaching public IP from MariaDB DBaaS (ID: %s): %s"
	ErrorDetachingParameterGroupTemplate = "error detaching parameter group (ID: %d) from MariaDB DBaaS (ID: %s): %s"
	ErrorAttachingParameterGroupTemplate = "error attaching parameter group (ID: %d) to MariaDB DBaaS (ID: %s): %s"
	ErrorUpgradingPlanTemplate           = "error upgrading MariaDB DBaaS (ID: %s) plan from (%s) to (%s): %s"
	ErrorExpandingDiskTemplate           = "error expanding MariaDB DBaaS (ID: %s) disk by %d GB: %s"
	ErrorRetrievingSoftwareIDForUpgrade  = "error retrieving %s software ID for version (%s) while upgrading plan for DBaaS (ID: %s): %s"
	ErrorRetrievingTemplateIDForUpgrade  = "error retrieving %s template ID for plan (%s) while upgrading DBaaS (ID: %s): %s"
	ErrorUnsupportedStatusTemplate       = "error updating MariaDB DBaaS (ID: %s): unsupported status value: %s. Must be one of: STOPPED, RUNNING, RESTARTING"
	ErrorCannotUpgradePlanTemplate       = "cannot upgrade plan for MariaDB DBaaS (ID: %s): database must be in STOPPED state (current state: %s). Please stop the instance first"
	ErrorCannotExpandDiskTemplate        = "cannot expand disk for MariaDB DBaaS (ID: %s): database must be in STOPPED state (current state: %s)"
)
