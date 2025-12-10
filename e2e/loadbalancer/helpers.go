package loadbalancer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	e2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// LoadBalancerAPIResponse represents the actual nested API response structure
type LoadBalancerAPIResponse struct {
	Code    int                 `json:"code"`
	Message string              `json:"message"`
	Data    LoadBalancerAPIData `json:"data"`
}

// LoadBalancerAPIData represents the nested data structure in the API response
type LoadBalancerAPIData struct {
	Name              string                  `json:"name"`
	NodeDetail        LoadBalancerNodeDetail  `json:"node_detail"`
	ApplianceInstance []LoadBalancerAppliance `json:"appliance_instance"`
	LBStatus          LoadBalancerStatus      `json:"lb_status"`
}

// LoadBalancerNodeDetail represents node detail information
type LoadBalancerNodeDetail struct {
	RAM       string  `json:"ram"`
	Disk      string  `json:"disk"`
	VCPU      float64 `json:"vcpu"`
	PlanName  string  `json:"plan_name"`
	PrivateIP string  `json:"private_ip"`
	PublicIP  string  `json:"public_ip"`
}

// LoadBalancerAppliance represents appliance instance information
type LoadBalancerAppliance struct {
	Context LoadBalancerContext `json:"context"`
}

// LoadBalancerContext represents context information within appliance instance
type LoadBalancerContext struct {
	LBMode         string `json:"lb_mode"`
	HostTargetIPv6 string `json:"host_target_ipv6,omitempty"`
}

// LoadBalancerStatus represents the load balancer status information
type LoadBalancerStatus struct {
	Status      string              `json:"status"`
	DataMonitor LoadBalancerMonitor `json:"data_monitor"`
}

// LoadBalancerMonitor represents data monitor information
type LoadBalancerMonitor struct {
	Status bool `json:"status"`
}

func GetLbPort(mode string) string {
	if mode == "HTTP" {
		return "80"
	}
	return "443"
}

// ExpandBackendsWithGoe2e expands backends using goe2e client
func ExpandBackendsWithGoe2e(ctx context.Context, config []interface{}, goe2eClient *goe2e.Client, project_id string, region string) ([]goe2e.LBBackend, error) {
	backends := make([]goe2e.LBBackend, 0, len(config))

	for _, backend := range config {
		detail := backend.(map[string]interface{})

		servers, err := ExpandServersWithGoe2e(ctx, detail["servers"], goe2eClient, project_id, region)
		if err != nil {
			return nil, err
		}
		r := goe2e.LBBackend{
			Balance:        detail["balance"].(string),
			CheckboxEnable: detail["checkbox_enable"].(bool),
			DomainName:     detail["domain_name"].(string),
			CheckURL:       detail["check_url"].(string),
			Servers:        servers,
			HTTPCheck:      detail["http_check"].(bool),
			Name:           detail["name"].(string),
			ScalerID:       detail["scaler_id"].(string),
			ScalerPort:     detail["scaler_port"].(string),
		}
		backends = append(backends, r)
	}
	return backends, nil
}

func ExpandAclList(config []interface{}) ([]goe2e.LBACLList, error) {
	aclList := make([]goe2e.LBACLList, 0, len(config))

	for _, acl := range config {
		aclRule := acl.(map[string]interface{})

		r := goe2e.LBACLList{
			ACLName:         aclRule["acl_name"].(string),
			ACLCondition:    aclRule["acl_condition"].(string),
			ACLMatchingPath: aclRule["acl_matching_path"].(string),
		}

		aclList = append(aclList, r)
	}
	return aclList, nil
}

func ExpandAclMap(config []interface{}) ([]goe2e.LBACLMap, error) {
	aclMap := make([]goe2e.LBACLMap, 0, len(config))

	for _, backendlist := range config {
		aclMapData := backendlist.(map[string]interface{})

		r := goe2e.LBACLMap{
			ACLName:           aclMapData["acl_name"].(string),
			ACLConditionState: true,
			ACLBackend:        aclMapData["acl_backend"].(string),
		}

		aclMap = append(aclMap, r)
	}
	return aclMap, nil
}

func ExpandEnableEosLogger(config []interface{}) (goe2e.LBEOSDetail, error) {
	eosDetail := goe2e.LBEOSDetail{}

	for _, eosBucketInfo := range config {
		detail := eosBucketInfo.(map[string]interface{})
		eosDetail.ApplianceID = detail["appliance_id"].(int)
		eosDetail.AccessKey = detail["access_key"].(string)
		eosDetail.SecretKey = detail["secret_key"].(string)
		eosDetail.Bucket = detail["bucket"].(string)
	}
	return eosDetail, nil
}

// ExpandTcpBackendWithGoe2e expands TCP backends using goe2e client
func ExpandTcpBackendWithGoe2e(ctx context.Context, config []interface{}, goe2eClient *goe2e.Client, project_id string, region string) ([]goe2e.LBTCPBackend, error) {
	tcpBackends := make([]goe2e.LBTCPBackend, 0, len(config))

	for _, tcpBackend := range config {
		detail := tcpBackend.(map[string]interface{})

		servers, err := ExpandServersWithGoe2e(ctx, detail["servers"], goe2eClient, project_id, region)
		if err != nil {
			return nil, err
		}
		r := goe2e.LBTCPBackend{
			BackendName: detail["backend_name"].(string),
			Port:        detail["port"].(string),
			Balance:     detail["balance"].(string),
			Servers:     servers,
		}

		tcpBackends = append(tcpBackends, r)
	}
	return tcpBackends, nil
}

func SetLoadBalancerStatus(d *schema.ResourceData, status_detail interface{}) error {
	haproxyStatus := status_detail.(map[string]interface{})
	dataMonitor := haproxyStatus["data_monitor"].(map[string]interface{})
	status := haproxyStatus["status"].(string)

	if status == goe2econstants.LBStatusRunningAPI {
		if len(dataMonitor) == 0 {
			d.Set(e2econstants.AttrStatus, goe2econstants.LBStatusBackendUnavailable)
			return nil
		}
		if !dataMonitor["status"].(bool) {
			d.Set(e2econstants.AttrStatus, goe2econstants.LBStatusBackendFailure)
		} else {
			d.Set(e2econstants.AttrStatus, goe2econstants.LBStatusRunning)
		}
	} else if status == goe2econstants.LBStatusPoweredOffAPI {
		d.Set(e2econstants.AttrStatus, goe2econstants.LBStatusPoweredOff)
	} else if status == goe2econstants.LBStatusCreating {
		d.Set(e2econstants.AttrStatus, goe2econstants.LBStatusCreating)
	} else if status == goe2econstants.LBStatusDeploying {
		d.Set(e2econstants.AttrStatus, goe2econstants.LBStatusDeploying)
	} else if status == goe2econstants.LBStatusUpgradingAPI {
		d.Set(e2econstants.AttrStatus, goe2econstants.LBStatusUpgrading)
	} else {
		d.Set(e2econstants.AttrStatus, goe2econstants.LBStatusError)
	}
	return nil
}

func CheckStatus(statuslist []string, status string) bool {
	for _, s := range statuslist {
		if strings.EqualFold(s, status) {
			return true
		}
	}
	return false
}

// Updated ExpandServers to use goe2e client for node lookups
func ExpandServersWithGoe2e(ctx context.Context, server_details interface{}, goe2eClient *goe2e.Client, project_id string, region string) ([]goe2e.LBServer, error) {
	var servers []goe2e.LBServer

	for _, server := range server_details.([]interface{}) {
		server_detail := server.(map[string]interface{})

		// Handle node_id vs id (prefer node_id, fallback to id)
		var nodeID string
		if nodeIDVal, ok := server_detail["node_id"].(string); ok && nodeIDVal != "" {
			nodeID = nodeIDVal
		} else if idVal, ok := server_detail["id"].(string); ok && idVal != "" {
			nodeID = idVal
		} else {
			return nil, fmt.Errorf("either 'node_id' or 'id' must be provided for server")
		}

		// Use goe2e client to get node info
		node, _, err := goe2eClient.Nodes.GetNode(ctx, nodeID)
		if err != nil {
			return nil, err
		}

		if node.Status != goe2econstants.NodeStatusRunning {
			return nil, fmt.Errorf("Node with id %s is not in running state", nodeID)
		}

		r := goe2e.LBServer{
			BackendName: node.Name,
			BackendIP:   node.PrivateIPAddress,
			BackendPort: server_detail["port"].(string),
		}

		servers = append(servers, r)
	}
	if len(servers) == 0 {
		return make([]goe2e.LBServer, 0), nil
	}
	return servers, nil
}

// Updated ExpandVpcList to use goe2e client
func ExpandVpcListWithGoe2e(ctx context.Context, d *schema.ResourceData, vpc_list []interface{}, goe2eClient *goe2e.Client) ([]goe2e.LBVPCDetail, error) {
	var vpc_details []goe2e.LBVPCDetail

	for _, id := range vpc_list {
		vpcID := strconv.Itoa(id.(int))
		vpc, _, err := goe2eClient.Vpcs.GetVPC(ctx, vpcID)
		if err != nil {
			return nil, err
		}

		if vpc.State != "Active" {
			return nil, fmt.Errorf("Can not attach vpc currently, vpc is in %s state", vpc.State)
		}

		r := goe2e.LBVPCDetail{
			NetworkID: vpc.ID,
			VPCName:   vpc.Name,
			IPv4CIDR:  vpc.IPv4CIDR,
		}

		vpc_details = append(vpc_details, r)
	}
	return vpc_details, nil
}

// getLoadBalancerWithNestedResponse fetches a load balancer using goe2e client
// and returns the typed nested response structure matching the actual API response
func getLoadBalancerWithNestedResponse(ctx context.Context, goe2eClient *goe2e.Client, lbID string) (*LoadBalancerAPIResponse, error) {
	if lbID == "" {
		return nil, fmt.Errorf("load balancer ID cannot be empty")
	}

	path := fmt.Sprintf("appliances/%s/", lbID)

	req, err := goe2eClient.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for Load Balancer (ID: %s): %w", lbID, err)
	}

	// Use a bytes.Buffer to capture the raw response body
	var responseBody bytes.Buffer
	resp, err := goe2eClient.Do(ctx, req, &responseBody)
	if err != nil {
		// Check if it's a 404 (not found)
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("load balancer with ID %s not found", lbID)
		}
		return nil, fmt.Errorf("failed to retrieve Load Balancer (ID: %s): %w", lbID, err)
	}

	// Unmarshal into typed nested struct
	var apiResponse LoadBalancerAPIResponse
	if err := json.Unmarshal(responseBody.Bytes(), &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Load Balancer response (ID: %s): %w", lbID, err)
	}

	return &apiResponse, nil
}

// waitForLoadBalancerStatus polls the load balancer status until it reaches the target status
// or until timeout occurs. Returns error if timeout or context cancellation.
// targetStatus should be a normalized status value (e.g., "running", "stopped")
func waitForLoadBalancerStatus(ctx context.Context, goe2eClient *goe2e.Client, lbID string, targetStatus string, timeoutMinutes int) error {
	timeout := time.Duration(timeoutMinutes) * time.Minute
	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled while waiting for load balancer %s to reach status %s", lbID, targetStatus)
		case <-ticker.C:
			if time.Now().After(deadline) {
				return fmt.Errorf("timeout waiting for load balancer %s to reach status %s after %d minutes", lbID, targetStatus, timeoutMinutes)
			}

			apiResponse, err := getLoadBalancerWithNestedResponse(ctx, goe2eClient, lbID)
			if err != nil {
				// If 404, check if we're waiting for deletion
				if strings.Contains(err.Error(), "not found") && targetStatus == "deleted" {
					return nil // Successfully deleted
				}
				return fmt.Errorf("error checking load balancer status: %w", err)
			}

			// Get status from lb_status.status field
			currentStatus := apiResponse.Data.LBStatus.Status

			// Normalize status for comparison
			normalizedStatus := normalizeLoadBalancerState(currentStatus)

			// Check if we've reached the target status (compare normalized values)
			if normalizedStatus == targetStatus {
				return nil
			}

			// Check for error states
			if normalizedStatus == "error" {
				return fmt.Errorf("load balancer %s entered error state (current status: %s)", lbID, currentStatus)
			}
		}
	}
}

// waitForLoadBalancerDeletion polls until the load balancer is deleted (404 response)
func waitForLoadBalancerDeletion(ctx context.Context, goe2eClient *goe2e.Client, lbID string, timeoutMinutes int) error {
	timeout := time.Duration(timeoutMinutes) * time.Minute
	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled while waiting for load balancer %s deletion", lbID)
		case <-ticker.C:
			if time.Now().After(deadline) {
				return fmt.Errorf("timeout waiting for load balancer %s deletion after %d minutes", lbID, timeoutMinutes)
			}

			_, err := getLoadBalancerWithNestedResponse(ctx, goe2eClient, lbID)
			if err != nil {
				// 404 means successfully deleted
				if strings.Contains(err.Error(), "not found") {
					return nil
				}
				// Other errors might be transient, continue polling
			}
		}
	}
}
