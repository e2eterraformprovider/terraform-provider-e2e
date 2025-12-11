package util

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// WaitForState waits for a resource to reach a target state using StateChangeConf.
// This is the low-level utility function that all other wait functions use internally.
// It follows the same pattern as DigitalOcean's WaitForAction but for resource-based polling.
//
// Parameters:
//   - ctx: Context for cancellation
//   - refreshFunc: Function that polls the current state
//   - pending: List of states that indicate we should keep waiting
//   - target: List of states that indicate success
//   - timeout: Maximum time to wait
//   - delay: Time between polling attempts (0 = use default)
//   - minTimeout: Minimum time between polls (0 = use default)
//   - notFoundChecks: Number of "not found" responses to allow (0 = use default)
//
// Example:
//
//	err := WaitForState(ctx, refreshFunc, []string{"pending"}, []string{"ready"}, 5*time.Minute, 0, 0, 0)
func WaitForState(
	ctx context.Context,
	refreshFunc resource.StateRefreshFunc,
	pending []string,
	target []string,
	timeout time.Duration,
	delay time.Duration,
	minTimeout time.Duration,
	notFoundChecks int,
) error {
	if delay == 0 {
		delay = tfconstants.StateChangeDefaultDelay
	}
	if minTimeout == 0 {
		minTimeout = tfconstants.StateChangeRetryBackoff
	}
	if notFoundChecks == 0 {
		notFoundChecks = tfconstants.DefaultNotFoundChecks
	}

	stateConf := &resource.StateChangeConf{
		Pending:        pending,
		Target:         target,
		Refresh:        refreshFunc,
		Timeout:        timeout,
		Delay:          delay,
		MinTimeout:     minTimeout,
		NotFoundChecks: notFoundChecks,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

// WaitForFunctionReady waits for a FaaS function to reach the "Ready" status.
// This function polls the function status until it becomes "Ready" or times out.
//
// Example:
//
//	err := util.WaitForFunctionReady(ctx, client, functionID, 10*time.Minute)
func WaitForFunctionReady(ctx context.Context, client *goe2e.Client, functionID string, timeout time.Duration) error {
	log.Printf("[INFO] Waiting for FaaS function %s to be ready", functionID)

	refreshFunc := func() (interface{}, string, error) {
		log.Printf("[DEBUG] Refreshing FaaS function status for function_id=%s", functionID)

		fn, _, err := client.FaaS.GetFunction(ctx, functionID)
		if err != nil {
			return nil, "", fmt.Errorf("error checking function status: %w", err)
		}

		if fn == nil {
			return nil, "", fmt.Errorf("function not found")
		}

		log.Printf("[DEBUG] FaaS function %s status: %s", functionID, fn.Status)

		// Check for error states
		if fn.Status == goe2econstants.FaaSStatusError || fn.Status == goe2econstants.FaaSStatusFailed {
			return nil, "", fmt.Errorf("function entered error state: %s", fn.Status)
		}

		return fn, fn.Status, nil
	}

	err := WaitForState(
		ctx,
		refreshFunc,
		tfconstants.FaaSPendingStates,
		[]string{goe2econstants.FaaSStatusReady},
		timeout,
		tfconstants.StateChangePollInterval,
		tfconstants.StateChangeRetryBackoff,
		tfconstants.DefaultNotFoundChecks,
	)

	if err != nil {
		return fmt.Errorf("error waiting for function to be ready: %w", err)
	}

	log.Printf("[INFO] FaaS function %s is ready", functionID)
	return nil
}

// WaitForNodePowerState waits for a node to reach a specific power state (RUNNING or POWERED_OFF).
// This is typically used when you need to wait for a node to finish powering on or off
// before performing another operation (e.g., plan upgrade).
//
// Example:
//
//	err := util.WaitForNodePowerState(m, nodeID, projectID, region)
func WaitForNodePowerState(m interface{}, nodeID string, projectID string, region string) error {
	cfg := m.(*config.Config)
	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return fmt.Errorf("error creating goe2e client: %w", err)
	}

	ctx := context.Background()
	log.Printf("[INFO] Waiting for node %s to reach power state (RUNNING or POWERED_OFF)", nodeID)

	refreshFunc := func() (interface{}, string, error) {
		log.Printf("[DEBUG] Refreshing node status for node_id=%s", nodeID)

		node, _, err := goe2eClient.Nodes.GetNode(ctx, nodeID)
		if err != nil {
			if IsNotFoundError(err) {
				log.Printf("[DEBUG] Node %s not found (may have been deleted)", nodeID)
				return nil, tfconstants.WaitStateDeleted, nil
			}
			return nil, "", fmt.Errorf("error retrieving node: %w", err)
		}

		if node == nil {
			return nil, tfconstants.WaitStateDeleted, nil
		}

		status := node.Status
		if status == "" {
			return node, tfconstants.WaitStateUnknown, fmt.Errorf("node status is empty or not found in API response")
		}

		log.Printf("[DEBUG] Node %s current status: %s", nodeID, status)

		// Check for error states
		if status == goe2econstants.NodeStatusFailed {
			return node, status, fmt.Errorf("node entered error state: %s", status)
		}

		return node, status, nil
	}

	err = WaitForState(
		ctx,
		refreshFunc,
		append([]string{goe2econstants.NodeStatusCreating}, tfconstants.NodePowerPendingStates...),
		[]string{goe2econstants.NodeStatusRunning, goe2econstants.NodeStatusPoweredOff},
		tfconstants.StateChangeTimeoutDefault,
		tfconstants.StateChangePollInterval,
		tfconstants.StateChangeRetryBackoff,
		tfconstants.DefaultNotFoundChecks,
	)

	if err != nil {
		return fmt.Errorf("error waiting for node power state: %w", err)
	}

	log.Printf("[INFO] Node %s reached target power state", nodeID)
	return nil
}

// WaitForNodeLCMState waits for a node LCM (Life Cycle Management) state to exit specific states.
// This is used when waiting for volume attachment/detachment operations to complete.
// The function waits until the LCM state is NOT in the excluded states (HOTPLUG, HOTPLUG_PROLOG_POWEROFF, HOTPLUG_EPILOG_POWEROFF).
//
// Example:
//
//	err := util.WaitForNodeLCMState(m, nodeID, projectID, region)
func WaitForNodeLCMState(m interface{}, nodeID string, projectID string, region string) error {
	cfg := m.(*config.Config)
	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return fmt.Errorf("error creating goe2e client: %w", err)
	}

	ctx := context.Background()
	excludeStates := []string{
		goe2econstants.NodeLCMStateHotplug,
		goe2econstants.NodeLCMStateHotplugPrologPoweroff,
		goe2econstants.NodeLCMStateHotplugEpilogPoweroff,
	}

	log.Printf("[INFO] Waiting for node %s LCM state to exit states: %v", nodeID, excludeStates)

	refreshFunc := func() (interface{}, string, error) {
		log.Printf("[DEBUG] Refreshing node LCM state for node_id=%s", nodeID)

		lcmState, _, err := goe2eClient.Nodes.GetLCMState(ctx, nodeID)
		if err != nil {
			return nil, "", fmt.Errorf("error retrieving node LCM state: %w", err)
		}

		if lcmState == nil {
			return nil, "", fmt.Errorf("node LCM state is nil")
		}

		state := lcmState.LCMState
		log.Printf("[DEBUG] Node %s current LCM state: %s", nodeID, state)

		// Check if state is in exclude list (meaning we should keep waiting)
		for _, excludeState := range excludeStates {
			if state == excludeState {
				return lcmState, state, nil
			}
		}

		// State is not in exclude list, we've reached target
		return lcmState, tfconstants.NodeLCMReadyState, nil
	}

	err = WaitForState(
		ctx,
		refreshFunc,
		excludeStates, // pending states (keep waiting while in these)
		tfconstants.NodeLCMReadyTargetStates,
		tfconstants.StateChangeTimeoutDefault,
		tfconstants.StateChangePollInterval,
		tfconstants.StateChangeRetryBackoff,
		tfconstants.DefaultNotFoundChecks,
	)

	if err != nil {
		return fmt.Errorf("error waiting for node LCM state: %w", err)
	}

	log.Printf("[INFO] Node %s LCM state exited target states", nodeID)
	return nil
}

// WaitForDBaaSPowerState waits for a DBaaS instance to reach SUSPENDED state.
// This is used when suspending a PostgreSQL DBaaS instance.
//
// Example:
//
//	err := util.WaitForDBaaSPowerState(m, dbaasID, projectID, region)
func WaitForDBaaSPowerState(m interface{}, dbaasID string, projectID string, region string) error {
	cfg := m.(*config.Config)
	goe2eClient, err := cfg.Goe2eClientForProject(projectID, region)
	if err != nil {
		return fmt.Errorf("error creating goe2e client for project (%s), region (%s): %w", projectID, region, err)
	}

	ctx := context.Background()
	log.Printf("[INFO] Waiting for DBaaS instance %s to reach SUSPENDED state", dbaasID)

	refreshFunc := func() (interface{}, string, error) {
		cluster, _, err := goe2eClient.PostgreSQL.GetCluster(ctx, dbaasID)
		if err != nil {
			log.Printf("[ERROR] Error fetching DBaaS Info during wait: %s", err)
			return nil, "", err
		}

		if cluster == nil {
			return nil, "", fmt.Errorf("DBaaS cluster not found")
		}

		log.Printf("[DEBUG] DBaaS instance %s current status: %s", dbaasID, cluster.Status)
		return cluster, cluster.Status, nil
	}

	err = WaitForState(
		ctx,
		refreshFunc,
		tfconstants.DBaaSSuspendPendingStates,
		tfconstants.DBaaSSuspendTargetStates,
		tfconstants.StateChangeTimeoutShort,
		tfconstants.StateChangePollInterval,
		tfconstants.StateChangeRetryBackoff,
		tfconstants.DefaultNotFoundChecks,
	)

	if err != nil {
		return fmt.Errorf("timeout: DBaaS did not reach SUSPENDED state in time: %w", err)
	}

	log.Printf("[INFO] DBaaS instance %s reached SUSPENDED state", dbaasID)
	return nil
}
