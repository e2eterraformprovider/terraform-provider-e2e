package goe2e

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

const (
	postgresqlClusterPath          = "rds/cluster" // Base path for cluster operations
	postgresqlClusterDetailPath    = "rds/cluster" // + /{id}/
	postgresqlClusterResumePath    = "rds/cluster" // + /{id}/resume
	postgresqlClusterShutdownPath  = "rds/cluster" // + /{id}/shutdown
	postgresqlClusterRestartPath   = "rds/cluster" // + /{id}/restart
	postgresqlClusterVPCAttachPath = "rds/cluster" // + /{id}/vpc-attach/
	postgresqlClusterVPCDetachPath = "rds/cluster" // + /{id}/vpc-detach/
	postgresqlClusterPGAddPath     = "rds/cluster" // + /{id}/parameter-group/{pg_id}/add
	postgresqlClusterPGDetachPath  = "rds/cluster" // + /{id}/parameter-group/{pg_id}/detach
	postgresqlClusterPublicIPPath  = "rds/cluster" // + /{id}/public-ip-attach/ or public-ip-detach/
	postgresqlClusterUpgradePath   = "rds/cluster" // + /{id}/rds-upgrade/
	postgresqlClusterDiskPath      = "rds/cluster" // + /{id}/disk-upgrade/
	postgresqlPlansPath            = "rds/plans"   // For fetching software and template IDs
)

// PostgreSQLService is an interface for interacting with the PostgreSQL DBaaS endpoints
// of the E2E Networks API.
type PostgreSQLService interface {
	// Cluster lifecycle operations
	CreateCluster(context.Context, *PostgreSQLClusterCreateRequest) (*PostgreSQLCluster, *Response, error)
	GetCluster(context.Context, string) (*PostgreSQLCluster, *Response, error)
	DeleteCluster(context.Context, string) (*Response, error)
	ClusterExists(context.Context, string) (bool, *Response, error)

	// Power management operations
	StartCluster(context.Context, string) (*Response, error)
	StopCluster(context.Context, string) (*Response, error)
	RestartCluster(context.Context, string) (*Response, error)

	// VPC operations
	AttachVPC(context.Context, string, *PostgreSQLVPCAttachRequest) (*Response, error)
	DetachVPC(context.Context, string, *PostgreSQLVPCAttachRequest) (*Response, error)

	// Parameter group operations
	AttachParameterGroup(context.Context, string, string) (*Response, error)
	DetachParameterGroup(context.Context, string, string) (*Response, error)

	// Public IP operations
	AttachPublicIP(context.Context, string) (*Response, error)
	DetachPublicIP(context.Context, string) (*Response, error)

	// Upgrade operations
	UpgradePlan(context.Context, string, *PostgreSQLPlanUpgradeRequest) (*Response, error)
	ExpandDisk(context.Context, string, *DiskExpansionRequest) (*Response, error)

	// Helper operations
	GetSoftwareID(context.Context, string, string, string) (int, error)
	GetTemplateID(context.Context, string, string, string) (int, error)
	ExpandPostgresVPCList(context.Context, []string) ([]VPCMetadata, error)
}

// PostgreSQLServiceOp handles communication with PostgreSQL DBaaS related methods of the
// E2E Networks API.
type PostgreSQLServiceOp struct {
	client *Client
}

var _ PostgreSQLService = &PostgreSQLServiceOp{}

// PostgreSQLClusterCreateRequest represents a request to create a PostgreSQL cluster.
type PostgreSQLClusterCreateRequest struct {
	Name                string        `json:"name"`
	SoftwareID          int           `json:"software_id"`
	TemplateID          int           `json:"template_id"`
	PublicIPRequired    bool          `json:"public_ip_required"`
	Group               string        `json:"group"`
	VPCs                []VPCMetadata `json:"vpcs"`
	Database            DBConfig      `json:"database"`
	PGID                *int          `json:"pg_id,omitempty"`
	IsEncryptionEnabled bool          `json:"isEncryptionEnabled"`
}

// PostgreSQLCluster represents a PostgreSQL DBaaS cluster.
type PostgreSQLCluster struct {
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

// PostgreSQLVPCAttachRequest represents a request to attach/detach VPCs to a PostgreSQL cluster.
type PostgreSQLVPCAttachRequest struct {
	Action string        `json:"action"`
	VPCs   []VPCMetadata `json:"vpcs"`
}

// PostgreSQLPlanUpgradeRequest represents a request to upgrade a PostgreSQL cluster plan.
type PostgreSQLPlanUpgradeRequest struct {
	TemplateID int `json:"template_id"`
}

// Response wrappers for API calls
type postgresqlClusterRoot struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type postgresqlClusterDetailRoot struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    PostgreSQLCluster `json:"data"`
}

type postgresqlOperationRoot struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type plansResponse struct {
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

// CreateCluster creates a new PostgreSQL DBaaS cluster.
func (s *PostgreSQLServiceOp) CreateCluster(ctx context.Context, createReq *PostgreSQLClusterCreateRequest) (*PostgreSQLCluster, *Response, error) {
	if createReq == nil {
		return nil, nil, NewArgError("createReq", "cannot be nil")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, postgresqlClusterPath+"/", createReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for PostgreSQL cluster (%s): %w", createReq.Name, err)
	}

	root := new(postgresqlClusterRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to create PostgreSQL cluster (%s): %w", createReq.Name, err)
	}

	// The create response returns the cluster data in a generic map,
	// we need to fetch the full cluster details
	if idVal, ok := root.Data["id"].(float64); ok {
		clusterID := fmt.Sprintf("%.0f", idVal)
		return s.GetCluster(ctx, clusterID)
	}

	return nil, resp, fmt.Errorf("failed to extract cluster ID from create response")
}

// GetCluster retrieves a PostgreSQL DBaaS cluster by ID.
func (s *PostgreSQLServiceOp) GetCluster(ctx context.Context, clusterID string) (*PostgreSQLCluster, *Response, error) {
	if clusterID == "" {
		return nil, nil, NewArgError("clusterID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", postgresqlClusterDetailPath, clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request for PostgreSQL cluster (ID: %s): %w", clusterID, err)
	}

	root := new(postgresqlClusterDetailRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		// Return nil cluster for 404 (not found)
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, resp, nil
		}
		return nil, resp, fmt.Errorf("failed to retrieve PostgreSQL cluster (ID: %s): %w", clusterID, err)
	}

	return &root.Data, resp, nil
}

// DeleteCluster deletes a PostgreSQL DBaaS cluster.
func (s *PostgreSQLServiceOp) DeleteCluster(ctx context.Context, clusterID string) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", postgresqlClusterDetailPath, clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for deleting PostgreSQL cluster (ID: %s): %w", clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to delete PostgreSQL cluster (ID: %s): %w", clusterID, err)
	}
	return resp, nil
}

// ClusterExists checks if a PostgreSQL DBaaS cluster exists.
func (s *PostgreSQLServiceOp) ClusterExists(ctx context.Context, clusterID string) (bool, *Response, error) {
	if clusterID == "" {
		return false, nil, NewArgError("clusterID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/", postgresqlClusterDetailPath, clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return false, nil, fmt.Errorf("failed to create request for PostgreSQL cluster existence check (ID: %s): %w", clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		// 404 means doesn't exist
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return false, resp, nil
		}
		return false, resp, fmt.Errorf("failed to check PostgreSQL cluster existence (ID: %s): %w", clusterID, err)
	}

	return true, resp, nil
}

// StartCluster starts/resumes a PostgreSQL DBaaS cluster.
func (s *PostgreSQLServiceOp) StartCluster(ctx context.Context, clusterID string) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/resume", postgresqlClusterResumePath, clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for starting PostgreSQL cluster (ID: %s): %w", clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to start PostgreSQL cluster (ID: %s): %w", clusterID, err)
	}
	return resp, nil
}

// StopCluster stops/shuts down a PostgreSQL DBaaS cluster.
func (s *PostgreSQLServiceOp) StopCluster(ctx context.Context, clusterID string) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/shutdown", postgresqlClusterShutdownPath, clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for stopping PostgreSQL cluster (ID: %s): %w", clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to stop PostgreSQL cluster (ID: %s): %w", clusterID, err)
	}
	return resp, nil
}

// RestartCluster restarts a PostgreSQL DBaaS cluster.
func (s *PostgreSQLServiceOp) RestartCluster(ctx context.Context, clusterID string) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/restart", postgresqlClusterRestartPath, clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for restarting PostgreSQL cluster (ID: %s): %w", clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to restart PostgreSQL cluster (ID: %s): %w", clusterID, err)
	}
	return resp, nil
}

// AttachVPC attaches VPCs to a PostgreSQL DBaaS cluster.
func (s *PostgreSQLServiceOp) AttachVPC(ctx context.Context, clusterID string, attachReq *PostgreSQLVPCAttachRequest) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}
	if attachReq == nil {
		return nil, NewArgError("attachReq", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%s/vpc-attach/", postgresqlClusterVPCAttachPath, clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, attachReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for attaching VPC to PostgreSQL cluster (ID: %s): %w", clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to attach VPC to PostgreSQL cluster (ID: %s): %w", clusterID, err)
	}
	return resp, nil
}

// DetachVPC detaches VPCs from a PostgreSQL DBaaS cluster.
func (s *PostgreSQLServiceOp) DetachVPC(ctx context.Context, clusterID string, detachReq *PostgreSQLVPCAttachRequest) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}
	if detachReq == nil {
		return nil, NewArgError("detachReq", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%s/vpc-detach/", postgresqlClusterVPCDetachPath, clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, detachReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for detaching VPC from PostgreSQL cluster (ID: %s): %w", clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to detach VPC from PostgreSQL cluster (ID: %s): %w", clusterID, err)
	}
	return resp, nil
}

// AttachPublicIP attaches a public IP to a PostgreSQL DBaaS cluster.
func (s *PostgreSQLServiceOp) AttachPublicIP(ctx context.Context, clusterID string) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/public-ip-attach/", postgresqlClusterPublicIPPath, clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for attaching public IP to PostgreSQL cluster (ID: %s): %w", clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to attach public IP to PostgreSQL cluster (ID: %s): %w", clusterID, err)
	}
	return resp, nil
}

// DetachPublicIP detaches a public IP from a PostgreSQL DBaaS cluster.
func (s *PostgreSQLServiceOp) DetachPublicIP(ctx context.Context, clusterID string) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/public-ip-detach/", postgresqlClusterPublicIPPath, clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for detaching public IP from PostgreSQL cluster (ID: %s): %w", clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to detach public IP from PostgreSQL cluster (ID: %s): %w", clusterID, err)
	}
	return resp, nil
}

// AttachParameterGroup attaches a parameter group to a PostgreSQL DBaaS cluster.
func (s *PostgreSQLServiceOp) AttachParameterGroup(ctx context.Context, clusterID string, parameterGroupID string) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}
	if parameterGroupID == "" {
		return nil, NewArgError("parameterGroupID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/parameter-group/%s/add", postgresqlClusterPGAddPath, clusterID, parameterGroupID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for attaching parameter group (%s) to PostgreSQL cluster (ID: %s): %w", parameterGroupID, clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to attach parameter group (%s) to PostgreSQL cluster (ID: %s): %w", parameterGroupID, clusterID, err)
	}
	return resp, nil
}

// DetachParameterGroup detaches a parameter group from a PostgreSQL DBaaS cluster.
func (s *PostgreSQLServiceOp) DetachParameterGroup(ctx context.Context, clusterID string, parameterGroupID string) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}
	if parameterGroupID == "" {
		return nil, NewArgError("parameterGroupID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/parameter-group/%s/detach", postgresqlClusterPGDetachPath, clusterID, parameterGroupID)

	// Send empty JSON object as body
	req, err := s.client.NewRequest(ctx, http.MethodPut, path, map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("failed to create request for detaching parameter group (%s) from PostgreSQL cluster (ID: %s): %w", parameterGroupID, clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to detach parameter group (%s) from PostgreSQL cluster (ID: %s): %w", parameterGroupID, clusterID, err)
	}
	return resp, nil
}

// UpgradePlan upgrades the plan of a PostgreSQL DBaaS cluster.
func (s *PostgreSQLServiceOp) UpgradePlan(ctx context.Context, clusterID string, upgradeReq *PostgreSQLPlanUpgradeRequest) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}
	if upgradeReq == nil {
		return nil, NewArgError("upgradeReq", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%s/rds-upgrade/", postgresqlClusterUpgradePath, clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, upgradeReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for upgrading plan of PostgreSQL cluster (ID: %s): %w", clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to upgrade plan of PostgreSQL cluster (ID: %s): %w", clusterID, err)
	}
	return resp, nil
}

// ExpandDisk expands the disk of a PostgreSQL DBaaS cluster.
func (s *PostgreSQLServiceOp) ExpandDisk(ctx context.Context, clusterID string, expandReq *DiskExpansionRequest) (*Response, error) {
	if clusterID == "" {
		return nil, NewArgError("clusterID", "cannot be empty")
	}
	if expandReq == nil {
		return nil, NewArgError("expandReq", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%s/disk-upgrade/", postgresqlClusterDiskPath, clusterID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, expandReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for expanding disk of PostgreSQL cluster (ID: %s): %w", clusterID, err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, fmt.Errorf("failed to expand disk of PostgreSQL cluster (ID: %s): %w", clusterID, err)
	}
	return resp, nil
}

// GetSoftwareID retrieves the software/engine ID for a given PostgreSQL version.
func (s *PostgreSQLServiceOp) GetSoftwareID(ctx context.Context, name, version string, pgID string) (int, error) {
	if name == "" {
		return -1, NewArgError("name", "cannot be empty")
	}
	if version == "" {
		return -1, NewArgError("version", "cannot be empty")
	}
	if pgID == "" {
		return -1, NewArgError("pgID", "cannot be empty")
	}

	path := postgresqlPlansPath + "/"

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return -1, fmt.Errorf("failed to create request for getting PostgreSQL software ID (name: %s, version: %s): %w", name, version, err)
	}

	root := new(plansResponse)
	_, err = s.client.Do(ctx, req, root)
	if err != nil {
		return -1, fmt.Errorf("failed to retrieve PostgreSQL software plans (name: %s, version: %s): %w", name, version, err)
	}

	for _, engine := range root.Data.DatabaseEngines {
		if engine.EngineName == name && engine.EngineVersion == version {
			return engine.EngineID, nil
		}
	}

	return -1, fmt.Errorf("matching PostgreSQL software not found (name: %s, version: %s)", name, version)
}

// GetTemplateID retrieves the template ID for a given PostgreSQL plan.
func (s *PostgreSQLServiceOp) GetTemplateID(ctx context.Context, plan, softwareID, pgID string) (int, error) {
	if plan == "" {
		return -1, NewArgError("plan", "cannot be empty")
	}
	if softwareID == "" {
		return -1, NewArgError("softwareID", "cannot be empty")
	}
	if pgID == "" {
		return -1, NewArgError("pgID", "cannot be empty")
	}

	path := postgresqlPlansPath + "/"

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return -1, fmt.Errorf("failed to create request for getting PostgreSQL template ID (plan: %s, softwareID: %s): %w", plan, softwareID, err)
	}

	q := req.URL.Query()
	q.Add("software_id", softwareID)
	req.URL.RawQuery = q.Encode()

	root := new(plansResponse)
	_, err = s.client.Do(ctx, req, root)
	if err != nil {
		return -1, fmt.Errorf("failed to retrieve PostgreSQL template plans (plan: %s, softwareID: %s): %w", plan, softwareID, err)
	}

	for _, template := range root.Data.TemplatePlans {
		if template.PlanName == plan {
			return template.PlanTemplateID, nil
		}
	}

	return -1, fmt.Errorf("matching PostgreSQL plan template not found (plan: %s)", plan)
}

// ExpandPostgresVPCList expands a list of VPC IDs into VPCMetadata by fetching VPC details.
func (s *PostgreSQLServiceOp) ExpandPostgresVPCList(ctx context.Context, vpcIDs []string) ([]VPCMetadata, error) {
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
