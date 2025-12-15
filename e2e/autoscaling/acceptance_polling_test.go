package autoscaling_test

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/acceptance"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/autoscaling"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Acceptance polling helpers.
//
// Motivation: autoscaling operations are async/eventually consistent. Terraform may return success
// before the API reflects the final state (status transitions, VPC attach/detach, public IP).
// These helpers reduce flakes by polling the API with conservative intervals/timeouts.

func waitUntil(ctx context.Context, interval time.Duration, fn func() (bool, error)) error {
	t := time.NewTicker(interval)
	defer t.Stop()

	for {
		ok, err := fn()
		if err != nil {
			return err
		}
		if ok {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
		}
	}
}

func stringSetEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	ac := append([]string(nil), a...)
	bc := append([]string(nil), b...)
	sort.Strings(ac)
	sort.Strings(bc)
	for i := range ac {
		if ac[i] != bc[i] {
			return false
		}
	}
	return true
}

func intSetEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	ac := append([]int(nil), a...)
	bc := append([]int(nil), b...)
	sort.Ints(ac)
	sort.Ints(bc)
	for i := range ac {
		if ac[i] != bc[i] {
			return false
		}
	}
	return true
}

func withDefaultTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout <= 0 {
		timeout = tfconstants.StateChangeTimeoutDefault
	}
	return context.WithTimeout(context.Background(), timeout)
}

func pollInterval() time.Duration {
	// Keep this conservative for acceptance tests to avoid rate limiting.
	if tfconstants.StateChangePollInterval <= 0 {
		return 5 * time.Second
	}
	return tfconstants.StateChangePollInterval
}

func waitForScalerGroupStatusNormalized(
	projectID string,
	location string,
	scalerGroupID string,
	wantNormalizedStatus string,
	timeout time.Duration,
) error {
	ctx, cancel := withDefaultTimeout(timeout)
	defer cancel()

	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	goe2eClient, err := cfg.Goe2eClientForProject(projectID, location)
	if err != nil {
		return fmt.Errorf("failed to create GoE2E client: %w", err)
	}

	return waitUntil(ctx, pollInterval(), func() (bool, error) {
		group, _, err := goe2eClient.Autoscaling.GetScalerGroup(ctx, scalerGroupID)
		if err != nil {
			return false, err
		}
		if group == nil {
			return false, fmt.Errorf("scaler group %s not found during polling", scalerGroupID)
		}
		got := autoscaling.NormalizeStatus(group.ProvisionStatus)
		return got == wantNormalizedStatus, nil
	})
}

func waitForAttachedVPCNames(
	projectID string,
	location string,
	scalerGroupID string,
	wantVPCNames []string,
	timeout time.Duration,
) error {
	ctx, cancel := withDefaultTimeout(timeout)
	defer cancel()

	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	goe2eClient, err := cfg.Goe2eClientForProject(projectID, location)
	if err != nil {
		return fmt.Errorf("failed to create GoE2E client: %w", err)
	}

	return waitUntil(ctx, pollInterval(), func() (bool, error) {
		vpcs, _, err := goe2eClient.Autoscaling.GetAttachedVPCsForScalerGroup(ctx, scalerGroupID)
		if err != nil {
			return false, err
		}
		var got []string
		for _, v := range vpcs {
			got = append(got, v.Name)
		}
		return stringSetEqual(got, wantVPCNames), nil
	})
}

func waitForPublicIPRequired(
	projectID string,
	location string,
	scalerGroupID string,
	want bool,
	timeout time.Duration,
) error {
	ctx, cancel := withDefaultTimeout(timeout)
	defer cancel()

	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	goe2eClient, err := cfg.Goe2eClientForProject(projectID, location)
	if err != nil {
		return fmt.Errorf("failed to create GoE2E client: %w", err)
	}

	return waitUntil(ctx, pollInterval(), func() (bool, error) {
		ipStatus, _, err := goe2eClient.Autoscaling.GetPublicIPStatus(ctx, scalerGroupID)
		if err != nil {
			return false, err
		}
		return ipStatus != nil && ipStatus.IsPublicIPRequired == want, nil
	})
}

func waitForSecurityGroupIDsInState(
	s *terraform.State,
	resourceName string,
	want []int,
	timeout time.Duration,
) error {
	// The Autoscaling API client currently does not expose "list attached SG IDs" for a scaler group.
	// The provider also does not fetch SG IDs from the scaler group read endpoint. As a result, the
	// most reliable signal we can assert in acceptance tests is the Terraform state itself.
	ctx, cancel := withDefaultTimeout(timeout)
	defer cancel()

	return waitUntil(ctx, pollInterval(), func() (bool, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return false, fmt.Errorf("resource not found in state: %s", resourceName)
		}
		attrs := rs.Primary.Attributes
		countStr := attrs[tfconstants.AttrSecurityGroupIDs+".#"]
		if countStr == "" {
			return false, nil
		}
		count, err := strconv.Atoi(countStr)
		if err != nil {
			return false, fmt.Errorf("invalid %s.# value: %q", tfconstants.AttrSecurityGroupIDs, countStr)
		}
		got := make([]int, 0, count)
		for i := 0; i < count; i++ {
			v := attrs[fmt.Sprintf("%s.%d", tfconstants.AttrSecurityGroupIDs, i)]
			if v == "" {
				continue
			}
			id, err := strconv.Atoi(v)
			if err != nil {
				return false, fmt.Errorf("invalid %s.%d value: %q", tfconstants.AttrSecurityGroupIDs, i, v)
			}
			got = append(got, id)
		}
		return intSetEqual(got, want), nil
	})
}
