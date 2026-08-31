package models

type ScaleGroupNode struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	IP       []string `json:"ip"`
	PublicIP string   `json:"public_ip"`
	Status   string   `json:"status"`
	RealCPU  string   `json:"real_cpu"`
}

type ElasticPolicy struct {
	Type          string `json:"type"`
	Adjust        int    `json:"adjust"`
	Parameter     string `json:"parameter"`
	Operator      string `json:"operator"`
	Value         string `json:"value"`
	PeriodNumber  string `json:"period_number"`
	PeriodSeconds string `json:"period"`
	Cooldown      string `json:"cooldown"`
}

type ScheduledPolicy struct {
	Type       string `json:"type"`
	Adjust     string `json:"adjust"`
	Recurrence string `json:"recurrence"`
}

type CreateScalerGroupRequest struct {
	Name                 string            `json:"name"`
	PlanID               string            `json:"plan_id"`
	PlanName             string            `json:"plan_name"`
	SlugName             string            `json:"slug_name"`
	SKUID                string            `json:"sku_id"`
	VMImageID            string            `json:"vm_image_id"`
	VMImageName          string            `json:"vm_image_name"`
	VMTemplateID         int               `json:"vm_template_id"`
	MyAccountSGID        int               `json:"my_account_sg_id"`
	IsEncryptionEnabled  bool              `json:"isEncryptionEnabled"`
	EncryptionPassphrase string            `json:"encryption_passphrase,omitempty"`
	IsPublicIPRequired   bool              `json:"is_public_ip_required"`
	MinNodes             string            `json:"min_nodes"`
	MaxNodes             string            `json:"max_nodes"`
	Desired              string            `json:"desired"`
	PolicyType           string            `json:"policy_type,omitempty"`
	Policy               []ElasticPolicy   `json:"policy,omitempty"`
	ScheduledPolicy      []ScheduledPolicy `json:"scheduled_policy,omitempty"`
	VPC                  []VPCDetail       `json:"vpc,omitempty"`
}

type VPCDetail struct {
	Name      string         `json:"name,omitempty"`
	NetworkID int            `json:"network_id"`
	IPv4CIDR  string         `json:"ipv4_cidr,omitempty"`
	State     string         `json:"state,omitempty"`
	Subnets   []SubnetDetail `json:"subnets,omitempty"`
}

type SubnetDetail struct {
	ID         int    `json:"id"`
	SubnetName string `json:"subnet_name"`
	CIDR       string `json:"cidr"`
	UsedIPs    int    `json:"usedIPs"`
	TotalIPs   int    `json:"totalIPs"`
}

type CreateScalerGroupResponse struct {
	Code    int                      `json:"code"`
	Message string                   `json:"message"`
	Errors  map[string]string        `json:"errors"`
	Data    ScalerGroupCreateDetails `json:"data"`
}

type ScalerGroupCreateDetails struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	VMImageName     string `json:"vm_image_name"`
	ProvisionStatus string `json:"provision_status"`
	Running         int    `json:"running"`
	Desired         int    `json:"desired"`
	Tags            string `json:"tags"`
	MinNodes        int    `json:"min_nodes"`
	MaxNodes        int    `json:"max_nodes"`
	CustomerID      int    `json:"customer_id"`
	PlanID          int    `json:"plan_id"`
	ImageID         int    `json:"image_id"`
	VPCNames        string `json:"vpc_names"`
	ElasticPolicy   struct {
		Type                 string `json:"type"`
		Policy               string `json:"policy"`
		PolicyMeasure        string `json:"policy_measure"`
		PolicyOp             string `json:"policy_op"`
		UpscalePolicyValue   int    `json:"upscale_policy_value"`
		DownscalePolicyValue int    `json:"downscale_policy_value"`
		WaitForPeriod        int    `json:"wait_for_period"`
		WaitPeriod           int    `json:"wait_period"`
		Cooldown             int    `json:"cooldown"`
	} `json:"elastic_policy"`
	ScheduledPolicy struct {
		Type                string `json:"type"`
		ScheduledPolicyIP   string `json:"scheduled_policy_ip"`
		UpscaleRecurrence   string `json:"upscale_recurrence"`
		UpscaleAdjust       int    `json:"upscale_adjust"`
		DownscaleRecurrence string `json:"downscale_recurrence"`
		DownscaleAdjust     int    `json:"downscale_adjust"`
	} `json:"scheduled_policy"`
}

type GetScalerGroupResponse struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Errors  map[string]string    `json:"errors"`
	Data    ScalerGroupGetDetail `json:"data"`
}

type ScalerGroupGetDetail struct {
	ID                      int              `json:"id"`
	Name                    string           `json:"name"`
	Running                 int              `json:"running"`
	Desired                 int              `json:"desired"`
	ProvisionStatus         string           `json:"provision_status"`
	Tags                    string           `json:"tags"`
	MinNodes                int              `json:"min_nodes"`
	MaxNodes                int              `json:"max_nodes"`
	PlanName                string           `json:"plan_name"`
	VMImageName             string           `json:"vm_image_name"`
	CustomerID              int              `json:"customer_id"`
	PlanID                  int              `json:"plan_id"`
	ImageID                 int              `json:"image_id"`
	Nodes                   []ScaleGroupNode `json:"nodes"`
	PolicyOp                string           `json:"policy_op"`
	Policy                  string           `json:"policy"`
	UpscalePolicyValue      int              `json:"upscale_policy_value"`
	DownscalePolicyValue    int              `json:"downscale_policy_value"`
	WaitForPeriod           int              `json:"wait_for_period"`
	WaitPeriod              int              `json:"wait_period"`
	Cooldown                int              `json:"cooldown"`
	PolicyMeasure           string           `json:"policy_measure"`
	PolicyUpscaleOperator   string           `json:"policy_upscale_operator"`
	PolicyDownscaleOperator string           `json:"policy_downscale_operator"`
	PolicyType              string           `json:"policy_type"`
	ParameterEvaluatedValue float64          `json:"parameter_evaluated_value"`
	ScheduledPolicyOp       string           `json:"scheduled_policy_op"`
	UpscaleRecurrence       string           `json:"upscale_recurrence"`
	UpscaleAdjust           int              `json:"upscale_adjust"`
	DownscaleRecurrence     string           `json:"downscale_recurrence"`
	DownscaleAdjust         int              `json:"downscale_adjust"`
}

type DeleteScalerGroupResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
	Data    map[string]any    `json:"data"`
}

type SavedImage struct {
	Name               string `json:"name"`
	ImageID            string `json:"image_id"`
	TemplateID         int    `json:"template_id"`
	Distro             string `json:"distro"`
	SKUType            string `json:"sku_type"`
	OSDistribution     string `json:"os_distribution"`
	NodePlansAvailable bool   `json:"node_plans_available"`
	AutoScaleTemplate  bool   `json:"auto_scale_template"`
}

type ListSavedImagesResponse struct {
	Code    int                    `json:"code"`
	Data    []SavedImage           `json:"data"`
	Message string                 `json:"message"`
	Errors  map[string]interface{} `json:"errors"`
}

type ScalerSecurityGroup struct {
	ID        int  `json:"id"`
	IsDefault bool `json:"is_default"`
}

type GetScalerSecurityGroupsResponse struct {
	Code    int                   `json:"code"`
	Data    []ScalerSecurityGroup `json:"data"`
	Message string                `json:"message"`
}

type UpdateScalerGroupRequest struct {
	Name            string            `json:"name"`
	PlanID          string            `json:"plan_id"`
	MinNodes        int               `json:"min_nodes"`
	MaxNodes        int               `json:"max_nodes"`
	PolicyType      string            `json:"policy_type,omitempty"`

	Policy          []ElasticPolicy   `json:"policy"`
	ScheduledPolicy []ScheduledPolicy `json:"scheduled_policy"`
}

type UpdateDesiredNodeCountRequest struct {
	Cardinality int `json:"cardinality"`
}

type AttachVPCRequest struct {
	VPCID string `json:"vpc_id"`
}

type AttachVPCResponse struct {
	Code    int               `json:"code"`
	Data    string            `json:"data"`
	Errors  map[string]string `json:"errors"`
	Message string            `json:"message"`
}

type DetachVPCResponse struct {
	Code    int               `json:"code"`
	Data    string            `json:"data"`
	Errors  map[string]string `json:"errors"`
	Message string            `json:"message"`
}

type PublicIPStatusResponse struct {
	Code    int                    `json:"code"`
	Data    PublicIPStatusData     `json:"data"`
	Errors  map[string]interface{} `json:"errors"`
	Message string                 `json:"message"`
}

type PublicIPStatusData struct {
	IsPublicIPRequired bool `json:"is_public_ip_required"`
}

type PublicIPActionResponse struct {
	Code    int                    `json:"code"`
	Data    string                 `json:"data"`
	Errors  map[string]interface{} `json:"errors"`
	Message string                 `json:"message"`
}

type VPCPartial struct {
	Name      string `json:"name"`
	NetworkID int    `json:"network_id"`
	IPv4CIDR  string `json:"ipv4_cidr"`
}
