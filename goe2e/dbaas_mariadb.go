package goe2e

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

const (
	mariadbBasePath       = "rds/cluster"
	mariadbShutdownSuffix = "shutdown"
	mariadbResumeSuffix   = "resume"
	mariadbRestartSuffix  = "restart"
	mariadbVPCAttachPath  = "vpc-attach"
	mariadbVPCDetachPath  = "vpc-detach"
	mariadbPublicIPAttach = "public-ip-attach"
	mariadbPublicIPDetach = "public-ip-detach"
	mariadbParamGroupAdd  = "parameter-group"
	mariadbPlanUpgrade    = "rds-upgrade"
	mariadbDiskUpgrade    = "disk-upgrade"
)

// MariaDBService is an interface for interacting with the MariaDB/DBaaS endpoints
// of the E2E Networks API.
type MariaDBService interface {
	// CRUD operations
	CreateMariaDB(context.Context, *MariaDBCreateRequest) (*MariaDB, *Response, error)
	GetMariaDB(context.Context, string) (*MariaDB, *Response, error)
	DeleteMariaDB(context.Context, string) (*Response, error)
	MariaDBExists(context.Context, string) (bool, *Response, error)

	// Lifecycle operations
	ShutdownMariaDB(context.Context, string) (*Response, error)
	ResumeMariaDB(context.Context, string) (*Response, error)
	RestartMariaDB(context.Context, string) (*Response, error)

	// VPC operations
	AttachVPC(context.Context, string, []string) (*Response, error)
	DetachVPC(context.Context, string, []string) (*Response, error)

	// Public IP operations
	AttachPublicIP(context.Context, string) (*Response, error)
	DetachPublicIP(context.Context, string) (*Response, error)

	// Parameter group operations
	AttachParameterGroup(context.Context, string, int) (*Response, error)
	DetachParameterGroup(context.Context, string, int) (*Response, error)

	// Upgrade operations
	UpgradePlan(context.Context, string, int) (*Response, error)
	ExpandDisk(context.Context, string, int) (*Response, error)

	// Helper operations
	ExpandVPCList(context.Context, []string) ([]VPCMetadata, error)
	GetSoftwareID(context.Context, string, string) (int, error)
	GetTemplateID(context.Context, string, int) (int, error)
}

// MariaDBServiceOp handles communication with MariaDB related methods of the
// E2E Networks API.
type MariaDBServiceOp struct {
	client *Client
}

var _ MariaDBService = &MariaDBServiceOp{}

// MariaDB represents a MariaDB cluster/database instance
type MariaDB struct {
	ID                  int      `json:"id"`
	Name                string   `json:"name"`
	Status              string   `json:"status"`
	StatusTitle         string   `json:"status_title"`
	StatusActions       []string `json:"status_actions"`
	NumInstances        int      `json:"num_instances"`
	Software            Software `json:"software"`
	MasterNode          DBNode   `json:"master_node"`
	ConnectivityDetail  string   `json:"connectivity_detail"`
	VectorDBStatus      string   `json:"vector_database_status"`
	ProjectName         string   `json:"project_name"`
	SnapshotExist       bool     `json:"snapshot_exist"`
	ZookeeperInstances  int      `json:"zookeeper_instances"`
	SlaveInstances      int      `json:"slave_instances"`
	IsEncryptionEnabled bool     `json:"isEncryptionEnabled"`
}

// MariaDBCreateRequest represents a request to create a MariaDB cluster
type MariaDBCreateRequest struct {
	Name                 string        `json:"name"`
	SoftwareID           int           `json:"software_id"`
	TemplateID           int           `json:"template_id"`
	PublicIPRequired     bool          `json:"public_ip_required"`
	Group                string        `json:"group"`
	VPCs                 []VPCMetadata `json:"vpcs,omitempty"`
	Database             DBConfig      `json:"database"`
	PGID                 int           `json:"pg_id,omitempty"`
	IsEncryptionEnabled  bool          `json:"isEncryptionEnabled"`
	EncryptionPassphrase string        `json:"encryption_passphrase,omitempty"`
}

// AttachDetachVPCRequest represents a VPC attach/detach request
type AttachDetachVPCRequest struct {
	Action string        `json:"action"`
	VPCs   []VPCMetadata `json:"vpcs"`
}

// ParameterGroupRequest represents a parameter group request
type ParameterGroupRequest struct {
	Action string `json:"action"`
}

// UpgradePlanRequest represents a plan upgrade request
type UpgradePlanRequest struct {
	TemplateID int `json:"template_id"`
}

// DiskUpgradeRequest represents a disk upgrade request
type DiskUpgradeRequest struct {
	Size int `json:"size"`
}

// Response wrappers for API calls
type mariadbRoot struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    MariaDB     `json:"data"`
	Errors  interface{} `json:"errors,omitempty"`
}

type publicIPActionRequest struct {
	Action string `json:"action"`
}

type mariadbPlansResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		DatabaseEngines []struct {
			EngineID      int    `json:"engine_id"`
			EngineName    string `json:"engine_name"`
			EngineVersion string `json:"engine_version"`
		} `json:"database_engines"`
		TemplatePlans []struct {
			PlanTemplateID int    `json:"plan_template_id"`
			PlanName       string `json:"plan_name"`
		} `json:"template_plans"`
	} `json:"data"`
}

// CreateMariaDB creates a new MariaDB cluster.
func (s *MariaDBServiceOp) CreateMariaDB(ctx context.Context, createReq *MariaDBCreateRequest) (*MariaDB, *Response, error) {
	if createReq == nil {
		return nil, nil, NewArgError("createReq", "cannot be nil")
	}
	if createReq.Name == "" {
		return nil, nil, NewArgError("createReq.Name", "cannot be empty")
	}
	if createReq.SoftwareID == 0 {
		return nil, nil, NewArgError("createReq.SoftwareID", "cannot be zero")
	}
	if createReq.TemplateID == 0 {
		return nil, nil, NewArgError("createReq.TemplateID", "cannot be zero")
	}

	path := mariadbBasePath + "/"

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for MariaDB cluster (%s): %w", createReq.Name, err)
	}

	root := new(mariadbRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to create MariaDB cluster (%s): %w", createReq.Name, err)
	}

	return &root.Data, resp, nil
}

// GetMariaDB retrieves a MariaDB cluster by ID.
func (s *MariaDBServiceOp) GetMariaDB(ctx context.Context, clusterID string) (*MariaDB, *Response, error) {
	if clusterID == "" {
		return nil, nil, NewArgError("clusterID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", mariadbBasePath, clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for MariaDB cluster (ID: %s): %w", clusterID, err)
	}

	root := new(mariadbRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		// Return nil cluster for 404 (not found)
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, resp, nil
		}
		return nil, resp, fmt.Errorf("failed to retrieve MariaDB cluster (ID: %s): %w", clusterID, err)
	}

	return &root.Data, resp, nil
}

// MariaDBExists checks if a MariaDB cluster exists.
func (s *MariaDBServiceOp) MariaDBExists(ctx context.Context, clusterID string) (bool, *Response, error) {
	if clusterID == "" {
		return false, nil, NewArgError("clusterID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", mariadbBasePath, clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return false, nil, fmt.Errorf("failed to create request for MariaDB cluster existence check (ID: %s): %w", clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		// 404 means doesn't exist
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return false, resp, nil
		}
		return false, resp, fmt.Errorf("failed to check MariaDB cluster existence (ID: %s): %w", clusterID, err)
	}

	return true, resp, nil
}

// DeleteMariaDB deletes a MariaDB cluster.
func (s *MariaDBServiceOp) DeleteMariaDB(ctx context.Context, clusterID string) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", mariadbBasePath, clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for deleting MariaDB cluster (ID: %s): %w", clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to delete MariaDB cluster (ID: %s): %w", clusterID, err)
	}

	return resp, nil
}

// ShutdownMariaDB shuts down a MariaDB cluster.
func (s *MariaDBServiceOp) ShutdownMariaDB(ctx context.Context, clusterID string) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/%s", mariadbBasePath, clusterID, mariadbShutdownSuffix)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for shutting down MariaDB cluster (ID: %s): %w", clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to shutdown MariaDB cluster (ID: %s): %w", clusterID, err)
	}

	return resp, nil
}

// ResumeMariaDB resumes a MariaDB cluster.
func (s *MariaDBServiceOp) ResumeMariaDB(ctx context.Context, clusterID string) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/%s", mariadbBasePath, clusterID, mariadbResumeSuffix)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for resuming MariaDB cluster (ID: %s): %w", clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to resume MariaDB cluster (ID: %s): %w", clusterID, err)
	}

	return resp, nil
}

// RestartMariaDB restarts a MariaDB cluster.
func (s *MariaDBServiceOp) RestartMariaDB(ctx context.Context, clusterID string) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/%s", mariadbBasePath, clusterID, mariadbRestartSuffix)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for restarting MariaDB cluster (ID: %s): %w", clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to restart MariaDB cluster (ID: %s): %w", clusterID, err)
	}

	return resp, nil
}

// AttachVPC attaches VPCs to a MariaDB cluster.
func (s *MariaDBServiceOp) AttachVPC(ctx context.Context, clusterID string, vpcIDs []string) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}
	if len(vpcIDs) == 0 {
		return nil, NewArgError("vpcIDs", "cannot be empty")
	}

	// Expand VPC metadata
	vpcMetaList, err := s.ExpandVPCList(ctx, vpcIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to expand VPC metadata: %w", err)
	}

	payload := AttachDetachVPCRequest{
		Action: "attach",
		VPCs:   vpcMetaList,
	}

	path := fmt.Sprintf("%s/%s/%s/", mariadbBasePath, clusterID, mariadbVPCAttachPath)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for attaching VPCs to MariaDB cluster (ID: %s): %w", clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to attach VPCs to MariaDB cluster (ID: %s): %w", clusterID, err)
	}

	return resp, nil
}

// DetachVPC detaches VPCs from a MariaDB cluster.
func (s *MariaDBServiceOp) DetachVPC(ctx context.Context, clusterID string, vpcIDs []string) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}
	if len(vpcIDs) == 0 {
		return nil, NewArgError("vpcIDs", "cannot be empty")
	}

	// Expand VPC metadata
	vpcMetaList, err := s.ExpandVPCList(ctx, vpcIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to expand VPC metadata: %w", err)
	}

	payload := AttachDetachVPCRequest{
		Action: "detach",
		VPCs:   vpcMetaList,
	}

	path := fmt.Sprintf("%s/%s/%s/", mariadbBasePath, clusterID, mariadbVPCDetachPath)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for detaching VPCs from MariaDB cluster (ID: %s): %w", clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to detach VPCs from MariaDB cluster (ID: %s): %w", clusterID, err)
	}

	return resp, nil
}

// AttachPublicIP attaches a public IP to a MariaDB cluster.
func (s *MariaDBServiceOp) AttachPublicIP(ctx context.Context, clusterID string) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}

	payload := publicIPActionRequest{Action: "attach"}
	path := fmt.Sprintf("%s/%s/%s/", mariadbBasePath, clusterID, mariadbPublicIPAttach)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for attaching public IP to MariaDB cluster (ID: %s): %w", clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to attach public IP to MariaDB cluster (ID: %s): %w", clusterID, err)
	}

	return resp, nil
}

// DetachPublicIP detaches a public IP from a MariaDB cluster.
func (s *MariaDBServiceOp) DetachPublicIP(ctx context.Context, clusterID string) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}

	payload := publicIPActionRequest{Action: "detach"}
	path := fmt.Sprintf("%s/%s/%s/", mariadbBasePath, clusterID, mariadbPublicIPDetach)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for detaching public IP from MariaDB cluster (ID: %s): %w", clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to detach public IP from MariaDB cluster (ID: %s): %w", clusterID, err)
	}

	return resp, nil
}

// AttachParameterGroup attaches a parameter group to a MariaDB cluster.
func (s *MariaDBServiceOp) AttachParameterGroup(ctx context.Context, clusterID string, parameterGroupID int) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}
	if parameterGroupID <= 0 {
		return nil, NewArgError("parameterGroupID", "must be greater than zero")
	}

	payload := ParameterGroupRequest{Action: "add"}
	path := fmt.Sprintf("%s/%s/%s/%d/add", mariadbBasePath, clusterID, mariadbParamGroupAdd, parameterGroupID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for attaching parameter group %d to MariaDB cluster (ID: %s): %w", parameterGroupID, clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to attach parameter group %d to MariaDB cluster (ID: %s): %w", parameterGroupID, clusterID, err)
	}

	return resp, nil
}

// DetachParameterGroup detaches a parameter group from a MariaDB cluster.
func (s *MariaDBServiceOp) DetachParameterGroup(ctx context.Context, clusterID string, parameterGroupID int) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}
	if parameterGroupID <= 0 {
		return nil, NewArgError("parameterGroupID", "must be greater than zero")
	}

	path := fmt.Sprintf("%s/%s/%s/%d/detach", mariadbBasePath, clusterID, mariadbParamGroupAdd, parameterGroupID)

	// Send empty JSON object as body
	req, err := s.client.NewRequest(ctx, http.MethodPut, path, map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("failed to create request for detaching parameter group %d from MariaDB cluster (ID: %s): %w", parameterGroupID, clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to detach parameter group %d from MariaDB cluster (ID: %s): %w", parameterGroupID, clusterID, err)
	}

	return resp, nil
}

// UpgradePlan upgrades the plan of a MariaDB cluster.
func (s *MariaDBServiceOp) UpgradePlan(ctx context.Context, clusterID string, templateID int) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}
	if templateID <= 0 {
		return nil, NewArgError("templateID", "must be greater than zero")
	}

	payload := UpgradePlanRequest{TemplateID: templateID}
	path := fmt.Sprintf("%s/%s/%s/", mariadbBasePath, clusterID, mariadbPlanUpgrade)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for upgrading plan of MariaDB cluster (ID: %s): %w", clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to upgrade plan of MariaDB cluster (ID: %s): %w", clusterID, err)
	}

	return resp, nil
}

// ExpandDisk expands the disk of a MariaDB cluster.
func (s *MariaDBServiceOp) ExpandDisk(ctx context.Context, clusterID string, additionalSize int) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}
	if additionalSize <= 0 {
		return nil, NewArgError("additionalSize", "must be greater than zero")
	}

	payload := DiskUpgradeRequest{Size: additionalSize}
	path := fmt.Sprintf("%s/%s/%s/", mariadbBasePath, clusterID, mariadbDiskUpgrade)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for expanding disk of MariaDB cluster (ID: %s): %w", clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to expand disk of MariaDB cluster (ID: %s): %w", clusterID, err)
	}

	return resp, nil
}

// ExpandVPCList expands a list of VPC IDs into VPCMetadata by fetching VPC details.
// This is a helper method used by AttachVPC and DetachVPC.
func (s *MariaDBServiceOp) ExpandVPCList(ctx context.Context, vpcIDs []string) ([]VPCMetadata, error) {
	if len(vpcIDs) == 0 {
		return nil, NewArgError("vpcIDs", "cannot be empty")
	}

	var vpcMetaList []VPCMetadata

	for _, vpcID := range vpcIDs {
		if vpcID == "" {
			continue
		}

		// Get VPC details using the Vpcs service
		vpc, _, err := s.client.Vpcs.GetVPC(ctx, vpcID)
		if err != nil {
			return nil, fmt.Errorf("failed to get VPC details for ID %s: %w", vpcID, err)
		}

		if vpc == nil {
			return nil, fmt.Errorf("VPC with ID %s not found", vpcID)
		}

		vpcMeta := VPCMetadata{
			NetworkID: strconv.FormatFloat(vpc.ID, 'f', -1, 64),
			VPCName:   vpc.Name,
			IPv4CIDR:  vpc.IPv4CIDR,
		}
		vpcMetaList = append(vpcMetaList, vpcMeta)
	}

	if len(vpcMetaList) == 0 {
		return nil, fmt.Errorf("no valid VPC IDs provided")
	}

	return vpcMetaList, nil
}

// GetSoftwareID retrieves the software/engine ID for a given MariaDB version.
func (s *MariaDBServiceOp) GetSoftwareID(ctx context.Context, name, version string) (int, error) {
	if name == "" {
		return -1, NewArgError("name", "cannot be empty")
	}
	if version == "" {
		return -1, NewArgError("version", "cannot be empty")
	}

	path := "rds/plans/"

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return -1, fmt.Errorf("failed to create request for getting MariaDB software ID (name: %s, version: %s): %w", name, version, err)
	}

	root := new(mariadbPlansResponse)
	_, err = s.client.Do(ctx, req, root)
	if err != nil {
		return -1, fmt.Errorf("failed to retrieve MariaDB software plans (name: %s, version: %s): %w", name, version, err)
	}

	for _, engine := range root.Data.DatabaseEngines {
		if engine.EngineName == name && engine.EngineVersion == version {
			return engine.EngineID, nil
		}
	}

	return -1, fmt.Errorf("matching MariaDB software not found (name: %s, version: %s)", name, version)
}

// GetTemplateID retrieves the template ID for a given MariaDB plan.
func (s *MariaDBServiceOp) GetTemplateID(ctx context.Context, plan string, softwareID int) (int, error) {
	if plan == "" {
		return -1, NewArgError("plan", "cannot be empty")
	}
	if softwareID <= 0 {
		return -1, NewArgError("softwareID", "must be greater than zero")
	}

	path := "rds/plans/"

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return -1, fmt.Errorf("failed to create request for getting MariaDB template ID (plan: %s, softwareID: %d): %w", plan, softwareID, err)
	}

	q := req.URL.Query()
	q.Add("software_id", strconv.Itoa(softwareID))
	req.URL.RawQuery = q.Encode()

	root := new(mariadbPlansResponse)
	_, err = s.client.Do(ctx, req, root)
	if err != nil {
		return -1, fmt.Errorf("failed to retrieve MariaDB template plans (plan: %s, softwareID: %d): %w", plan, softwareID, err)
	}

	for _, template := range root.Data.TemplatePlans {
		if template.PlanName == plan {
			return template.PlanTemplateID, nil
		}
	}

	return -1, fmt.Errorf("matching MariaDB plan template not found (plan: %s)", plan)
}
