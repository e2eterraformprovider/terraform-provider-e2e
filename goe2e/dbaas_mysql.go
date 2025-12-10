package goe2e

import (
	"context"
	"fmt"
	"net/http"
)

const (
	rdsClusterPath          = "rds/cluster" // Base path for cluster operations
	rdsClusterDetailPath    = "rds/cluster" // + /{id}/
	rdsClusterResumePath    = "rds/cluster" // + /{id}/resume
	rdsClusterShutdownPath  = "rds/cluster" // + /{id}/shutdown
	rdsClusterRestartPath   = "rds/cluster" // + /{id}/restart
	rdsClusterVPCAttachPath = "rds/cluster" // + /{id}/vpc-attach/
	rdsClusterVPCDetachPath = "rds/cluster" // + /{id}/vpc-detach/
	rdsClusterPGAddPath     = "rds/cluster" // + /{id}/parameter-group/{pg_id}/add
	rdsClusterPGDetachPath  = "rds/cluster" // + /{id}/parameter-group/{pg_id}/detach
	rdsClusterPublicIPPath  = "rds/cluster" // + /{id}/public-ip-attach/ or public-ip-detach/
	rdsClusterUpgradePath   = "rds/cluster" // + /{id}/rds-upgrade/
	rdsClusterDiskPath      = "rds/cluster" // + /{id}/disk-upgrade/
)

// DBaaSMySQLService is an interface for interacting with the MySQL DBaaS endpoints
// of the E2E Networks API.
type DBaaSMySQLService interface {
	// Cluster lifecycle operations
	CreateCluster(context.Context, *MySQLClusterCreateRequest) (*MySQLCluster, *Response, error)
	GetCluster(context.Context, string) (*MySQLCluster, *Response, error)
	DeleteCluster(context.Context, string) (*Response, error)

	// Power management operations
	StartCluster(context.Context, string) (*Response, error)
	StopCluster(context.Context, string) (*Response, error)
	RestartCluster(context.Context, string) (*Response, error)

	// VPC operations
	AttachVPC(context.Context, string, *MySQLVPCAttachRequest) (*Response, error)
	DetachVPC(context.Context, string, *MySQLVPCDetachRequest) (*Response, error)

	// Parameter group operations
	AttachParameterGroup(context.Context, string, string) (*Response, error)
	DetachParameterGroup(context.Context, string, string) (*Response, error)

	// Public IP operations
	AttachPublicIP(context.Context, string) (*Response, error)
	DetachPublicIP(context.Context, string) (*Response, error)

	// Upgrade operations
	UpgradePlan(context.Context, string, *MySQLPlanUpgradeRequest) (*Response, error)
	ExpandDisk(context.Context, string, *DiskExpansionRequest) (*Response, error)

	// Helper operations
	ExpandVPCList(context.Context, []string) ([]VPCMetadata, error)
	GetSoftwareID(context.Context, string, string) (int, error)
	GetTemplateID(context.Context, string, int) (int, error)
}

// DBaaSMySQLServiceOp handles communication with MySQL DBaaS related methods of the
// E2E Networks API.
type DBaaSMySQLServiceOp struct {
	client *Client
}

var _ DBaaSMySQLService = &DBaaSMySQLServiceOp{}

// MySQLClusterCreateRequest represents a request to create a MySQL cluster.
type MySQLClusterCreateRequest struct {
	Name             string        `json:"name"`
	Database         DBConfig      `json:"database"`
	Vpcs             []VPCMetadata `json:"vpcs"`
	SoftwareID       int           `json:"software_id"`
	TemplateID       int           `json:"template_id"`
	ParameterGroupId int           `json:"pg_id,omitempty"`
	PublicIPRequired bool          `json:"public_ip_required"`
	Group            string        `json:"group"`
}

// MySQLCluster represents a MySQL DBaaS cluster.
type MySQLCluster struct {
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

// MySQLVPCAttachRequest represents a request to attach VPCs to a MySQL cluster.
type MySQLVPCAttachRequest struct {
	Action string        `json:"action"`
	VPCs   []VPCMetadata `json:"vpcs"`
}

// MySQLVPCDetachRequest represents a request to detach VPCs from a MySQL cluster.
type MySQLVPCDetachRequest struct {
	Action string        `json:"action"`
	VPCs   []VPCMetadata `json:"vpcs"`
}

// MySQLPlanUpgradeRequest represents a request to upgrade a MySQL cluster plan.
type MySQLPlanUpgradeRequest struct {
	TemplateID int `json:"template_id"`
}

// DiskExpansionRequest represents a request to expand a MySQL cluster disk.
type DiskExpansionRequest struct {
	Size int `json:"size"`
}

// Response wrappers for API calls
type mysqlClusterRoot struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type mysqlClusterDetailRoot struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    MySQLCluster `json:"data"`
}

// CreateCluster creates a new MySQL DBaaS cluster.
func (s *DBaaSMySQLServiceOp) CreateCluster(ctx context.Context, createReq *MySQLClusterCreateRequest) (*MySQLCluster, *Response, error) {
	if createReq == nil {
		return nil, nil, NewArgError("createReq", "cannot be nil")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, rdsClusterPath+"/", createReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for MySQL cluster (%s): %w", createReq.Name, err)
	}

	root := new(mysqlClusterRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to create MySQL cluster (%s): %w", createReq.Name, err)
	}

	// The create response returns the cluster data in a generic map,
	// we need to fetch the full cluster details
	if idVal, ok := root.Data["id"].(float64); ok {
		clusterID := fmt.Sprintf("%.0f", idVal)
		return s.GetCluster(ctx, clusterID)
	}

	return nil, resp, fmt.Errorf("failed to extract cluster ID from create response")
}

// GetCluster retrieves a MySQL DBaaS cluster by ID.
func (s *DBaaSMySQLServiceOp) GetCluster(ctx context.Context, clusterID string) (*MySQLCluster, *Response, error) {
	if clusterID == "" {
		return nil, nil, NewArgError("clusterID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", rdsClusterDetailPath, clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for MySQL cluster (ID: %s): %w", clusterID, err)
	}

	root := new(mysqlClusterDetailRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		// Return nil cluster for 404 (not found)
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, resp, nil
		}
		return nil, resp, fmt.Errorf("failed to retrieve MySQL cluster (ID: %s): %w", clusterID, err)
	}

	return &root.Data, resp, nil
}

// DeleteCluster deletes a MySQL DBaaS cluster.
func (s *DBaaSMySQLServiceOp) DeleteCluster(ctx context.Context, clusterID string) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", rdsClusterDetailPath, clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for deleting MySQL cluster (ID: %s): %w", clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to delete MySQL cluster (ID: %s): %w", clusterID, err)
	}
	return resp, nil
}

// StartCluster starts/resumes a MySQL DBaaS cluster.
func (s *DBaaSMySQLServiceOp) StartCluster(ctx context.Context, clusterID string) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/resume", rdsClusterResumePath, clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for starting MySQL cluster (ID: %s): %w", clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to start MySQL cluster (ID: %s): %w", clusterID, err)
	}
	return resp, nil
}

// StopCluster stops/shuts down a MySQL DBaaS cluster.
func (s *DBaaSMySQLServiceOp) StopCluster(ctx context.Context, clusterID string) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/shutdown", rdsClusterShutdownPath, clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for stopping MySQL cluster (ID: %s): %w", clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to stop MySQL cluster (ID: %s): %w", clusterID, err)
	}
	return resp, nil
}

// RestartCluster restarts a MySQL DBaaS cluster.
func (s *DBaaSMySQLServiceOp) RestartCluster(ctx context.Context, clusterID string) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/restart", rdsClusterRestartPath, clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for restarting MySQL cluster (ID: %s): %w", clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to restart MySQL cluster (ID: %s): %w", clusterID, err)
	}
	return resp, nil
}

// AttachVPC attaches VPCs to a MySQL DBaaS cluster.
func (s *DBaaSMySQLServiceOp) AttachVPC(ctx context.Context, clusterID string, attachReq *MySQLVPCAttachRequest) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}
	if attachReq == nil {
		return nil, NewArgError("attachReq", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%s/vpc-attach/", rdsClusterVPCAttachPath, clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, attachReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for attaching VPC to MySQL cluster (ID: %s): %w", clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to attach VPC to MySQL cluster (ID: %s): %w", clusterID, err)
	}
	return resp, nil
}

// DetachVPC detaches VPCs from a MySQL DBaaS cluster.
func (s *DBaaSMySQLServiceOp) DetachVPC(ctx context.Context, clusterID string, detachReq *MySQLVPCDetachRequest) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}
	if detachReq == nil {
		return nil, NewArgError("detachReq", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%s/vpc-detach/", rdsClusterVPCDetachPath, clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, detachReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for detaching VPC from MySQL cluster (ID: %s): %w", clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to detach VPC from MySQL cluster (ID: %s): %w", clusterID, err)
	}
	return resp, nil
}

// AttachParameterGroup attaches a parameter group to a MySQL DBaaS cluster.
func (s *DBaaSMySQLServiceOp) AttachParameterGroup(ctx context.Context, clusterID string, parameterGroupID string) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}
	if parameterGroupID == "" {
		return nil, NewArgError("parameterGroupID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/parameter-group/%s/add", rdsClusterPGAddPath, clusterID, parameterGroupID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for attaching parameter group (%s) to MySQL cluster (ID: %s): %w", parameterGroupID, clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to attach parameter group (%s) to MySQL cluster (ID: %s): %w", parameterGroupID, clusterID, err)
	}
	return resp, nil
}

// DetachParameterGroup detaches a parameter group from a MySQL DBaaS cluster.
func (s *DBaaSMySQLServiceOp) DetachParameterGroup(ctx context.Context, clusterID string, parameterGroupID string) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}
	if parameterGroupID == "" {
		return nil, NewArgError("parameterGroupID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/parameter-group/%s/detach", rdsClusterPGDetachPath, clusterID, parameterGroupID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for detaching parameter group (%s) from MySQL cluster (ID: %s): %w", parameterGroupID, clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to detach parameter group (%s) from MySQL cluster (ID: %s): %w", parameterGroupID, clusterID, err)
	}
	return resp, nil
}

// AttachPublicIP attaches a public IP to a MySQL DBaaS cluster.
func (s *DBaaSMySQLServiceOp) AttachPublicIP(ctx context.Context, clusterID string) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/public-ip-attach/", rdsClusterPublicIPPath, clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for attaching public IP to MySQL cluster (ID: %s): %w", clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to attach public IP to MySQL cluster (ID: %s): %w", clusterID, err)
	}
	return resp, nil
}

// DetachPublicIP detaches a public IP from a MySQL DBaaS cluster.
func (s *DBaaSMySQLServiceOp) DetachPublicIP(ctx context.Context, clusterID string) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/public-ip-detach/", rdsClusterPublicIPPath, clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for detaching public IP from MySQL cluster (ID: %s): %w", clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to detach public IP from MySQL cluster (ID: %s): %w", clusterID, err)
	}
	return resp, nil
}

// UpgradePlan upgrades the plan of a MySQL DBaaS cluster.
func (s *DBaaSMySQLServiceOp) UpgradePlan(ctx context.Context, clusterID string, upgradeReq *MySQLPlanUpgradeRequest) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}
	if upgradeReq == nil {
		return nil, NewArgError("upgradeReq", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%s/rds-upgrade/", rdsClusterUpgradePath, clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, upgradeReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for upgrading plan of MySQL cluster (ID: %s): %w", clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to upgrade plan of MySQL cluster (ID: %s): %w", clusterID, err)
	}
	return resp, nil
}

// ExpandDisk expands the disk of a MySQL DBaaS cluster.
func (s *DBaaSMySQLServiceOp) ExpandDisk(ctx context.Context, clusterID string, expandReq *DiskExpansionRequest) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}
	if expandReq == nil {
		return nil, NewArgError("expandReq", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%s/disk-upgrade/", rdsClusterDiskPath, clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, expandReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for expanding disk of MySQL cluster (ID: %s): %w", clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to expand disk of MySQL cluster (ID: %s): %w", clusterID, err)
	}
	return resp, nil
}

// ExpandVPCList expands a list of VPC IDs into VPCMetadata by fetching VPC details.
// This is a helper method used by AttachVPC and DetachVPC.
func (s *DBaaSMySQLServiceOp) ExpandVPCList(ctx context.Context, vpcIDs []string) ([]VPCMetadata, error) {
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
			NetworkID: fmt.Sprintf("%.0f", vpc.ID),
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

// GetSoftwareID retrieves the software/engine ID for a given MySQL version.
func (s *DBaaSMySQLServiceOp) GetSoftwareID(ctx context.Context, name, version string) (int, error) {
	if name == "" {
		return -1, NewArgError("name", "cannot be empty")
	}
	if version == "" {
		return -1, NewArgError("version", "cannot be empty")
	}

	path := "rds/plans/"

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return -1, fmt.Errorf("failed to create request for getting MySQL software ID (name: %s, version: %s): %w", name, version, err)
	}

	root := new(mysqlPlansResponse)
	_, err = s.client.Do(ctx, req, root)
	if err != nil {
		return -1, fmt.Errorf("failed to retrieve MySQL software plans (name: %s, version: %s): %w", name, version, err)
	}

	for _, engine := range root.Data.DatabaseEngines {
		if engine.EngineName == name && engine.EngineVersion == version {
			return engine.EngineID, nil
		}
	}

	return -1, fmt.Errorf("matching MySQL software not found (name: %s, version: %s)", name, version)
}

// GetTemplateID retrieves the template ID for a given MySQL plan.
func (s *DBaaSMySQLServiceOp) GetTemplateID(ctx context.Context, plan string, softwareID int) (int, error) {
	if plan == "" {
		return -1, NewArgError("plan", "cannot be empty")
	}
	if softwareID <= 0 {
		return -1, NewArgError("softwareID", "must be greater than zero")
	}

	path := "rds/plans/"

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return -1, fmt.Errorf("failed to create request for getting MySQL template ID (plan: %s, softwareID: %d): %w", plan, softwareID, err)
	}

	q := req.URL.Query()
	q.Add("software_id", fmt.Sprintf("%d", softwareID))
	req.URL.RawQuery = q.Encode()

	root := new(mysqlPlansResponse)
	_, err = s.client.Do(ctx, req, root)
	if err != nil {
		return -1, fmt.Errorf("failed to retrieve MySQL template plans (plan: %s, softwareID: %d): %w", plan, softwareID, err)
	}

	for _, template := range root.Data.TemplatePlans {
		if template.PlanName == plan {
			return template.PlanTemplateID, nil
		}
	}

	return -1, fmt.Errorf("matching MySQL plan template not found (plan: %s)", plan)
}

// Response structures for plan queries
type mysqlPlansResponse struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Data    mysqlPlansData `json:"data"`
}

type mysqlPlansData struct {
	DatabaseEngines []mysqlDatabaseEngine `json:"database_engines"`
	TemplatePlans   []mysqlTemplatePlan   `json:"plans"`
}

type mysqlDatabaseEngine struct {
	EngineID      int    `json:"id"`
	EngineName    string `json:"name"`
	EngineVersion string `json:"version"`
}

type mysqlTemplatePlan struct {
	PlanTemplateID int    `json:"id"`
	PlanName       string `json:"name"`
}
