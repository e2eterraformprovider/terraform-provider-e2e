package sfs

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
)

const (
	sfsCreateTimeout   = 10 * time.Minute
	sfsDeleteTimeout   = 5 * time.Minute
	sfsPollingInterval = 10 * time.Second
)

// normalizeSfsState converts API status to normalized state
func normalizeSfsState(status string) string {
	switch strings.ToLower(status) {
	case strings.ToLower(goe2econstants.SFSStatusCreating):
		return goe2econstants.SFSStateCreating
	case strings.ToLower(goe2econstants.SFSStatusActive):
		return goe2econstants.SFSStateActive
	case strings.ToLower(goe2econstants.SFSStatusDeleting):
		return goe2econstants.SFSStateDeleting
	case strings.ToLower(goe2econstants.SFSStatusDeleted):
		return goe2econstants.SFSStateDeleted
	case strings.ToLower(goe2econstants.SFSStatusError):
		return goe2econstants.SFSStateError
	default:
		return strings.ToLower(status)
	}
}

// waitForSfsStatus polls the SFS status until it reaches the desired state or times out
func waitForSfsStatus(ctx context.Context, client *goe2e.Client, sfsID string, desiredStatus string, timeout time.Duration) error {
	// Helper function to check status
	checkStatus := func() (bool, error) {
		if client == nil || client.Sfs == nil {
			return false, fmt.Errorf(goe2econstants.ClientOrServiceNil)
		}
		sfs, _, err := client.Sfs.GetSfs(ctx, sfsID)
		if err != nil {
			// Check if it's a 404 (not found)
			if strings.Contains(err.Error(), goe2econstants.NotFoundSubstring) || strings.Contains(err.Error(), goe2econstants.NotFoundCode) {
				if desiredStatus == goe2econstants.SFSDesiredStatusDeleted || desiredStatus == goe2econstants.SFSDesiredStatus404 {
					return true, nil
				}
				return false, fmt.Errorf(goe2econstants.SFSNotFound, sfsID)
			}
			// Log transient errors but continue polling
			log.Printf("[WARN] Error polling SFS %s status: %s", sfsID, err)
			return false, nil
		}

		if sfs == nil {
			if desiredStatus == goe2econstants.SFSDesiredStatusDeleted || desiredStatus == goe2econstants.SFSDesiredStatus404 {
				return true, nil
			}
			return false, fmt.Errorf(goe2econstants.SFSNotFound, sfsID)
		}

		currentStatus := sfs.Status
		normalizedStatus := normalizeSfsState(currentStatus)

		// Check if we've reached desired status
		if normalizedStatus == desiredStatus || currentStatus == desiredStatus {
			log.Printf("[DEBUG] SFS %s reached desired status: %s", sfsID, desiredStatus)
			return true, nil
		}

		// Check for error status
		if normalizedStatus == goe2econstants.SFSStateError {
			return false, fmt.Errorf(goe2econstants.SFSEnteredErrorState, sfsID, desiredStatus)
		}

		log.Printf("[DEBUG] Waiting for SFS %s status: current=%s, desired=%s", sfsID, currentStatus, desiredStatus)
		return false, nil
	}

	// Check context before starting
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Check immediately before starting the polling loop
	done, err := checkStatus()
	if done {
		return err
	}
	if err != nil {
		return err
	}

	ticker := time.NewTicker(sfsPollingInterval)
	defer ticker.Stop()

	timeoutChan := time.After(timeout)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-timeoutChan:
			return fmt.Errorf(goe2econstants.SFSTimeoutWaitingForStatus, sfsID, desiredStatus)

		case <-ticker.C:
			done, err := checkStatus()
			if done {
				return err
			}
			if err != nil {
				return err
			}
		}
	}
}

// waitForSfsActive waits for SFS to become Active after creation
func waitForSfsActive(ctx context.Context, client *goe2e.Client, sfsID string) error {
	return waitForSfsStatus(ctx, client, sfsID, goe2econstants.SFSDesiredStatusActive, sfsCreateTimeout)
}

// parseSfsImportID parses the import ID string
// Supports two formats:
// 1. Simple: <sfs_id>
// 2. Full: <project_id>/<region>/<sfs_id>
func parseSfsImportID(id string) (projectID, region, sfsID string, err error) {
	parts := strings.Split(id, "/")

	switch len(parts) {
	case 1:
		// Simple format: just SFS ID
		sfsID = parts[0]
		return "", "", sfsID, nil

	case 3:
		// Full format: project_id/region/sfs_id
		projectID = parts[0]
		region = parts[1]
		sfsID = parts[2]
		if projectID == "" || region == "" || sfsID == "" {
			return "", "", "", fmt.Errorf(ImportIDInvalidFormat, id)
		}
		return projectID, region, sfsID, nil

	default:
		return "", "", "", fmt.Errorf(ImportIDInvalidFormat, id)
	}
}

// getEffectiveFieldValue returns the effective value from conflicting V2/V3 field pairs
// Prefers V3 fields over V2 fields if both are set
func getEffectiveSizeGB(d interface{}, attrSizeGB, attrDiskSize string, defaultValue int) int {
	type getter interface {
		Get(string) interface{}
		GetOk(string) (interface{}, bool)
	}

	if g, ok := d.(getter); ok {
		// Prefer V3 field
		if val, ok := g.GetOk(attrSizeGB); ok && val != 0 {
			if v, ok := val.(int); ok && v > 0 {
				return v
			}
		}
		// Fall back to V2 field
		if val, ok := g.GetOk(attrDiskSize); ok && val != 0 {
			if v, ok := val.(int); ok && v > 0 {
				return v
			}
		}
	}

	return defaultValue
}

// getEffectiveIOPS returns the effective IOPS value from conflicting V2/V3 field pairs
func getEffectiveIOPS(d interface{}, attrIOPS, attrDiskIOPS string, defaultValue int) int {
	type getter interface {
		Get(string) interface{}
		GetOk(string) (interface{}, bool)
	}

	if g, ok := d.(getter); ok {
		// Prefer V3 field
		if val, ok := g.GetOk(attrIOPS); ok && val != 0 {
			if v, ok := val.(int); ok && v > 0 {
				return v
			}
		}
		// Fall back to V2 field
		if val, ok := g.GetOk(attrDiskIOPS); ok && val != 0 {
			if v, ok := val.(int); ok && v > 0 {
				return v
			}
		}
	}

	return defaultValue
}

// getEffectiveEncryptionEnabled returns the effective encryption_enabled value from conflicting V2/V3 field pairs
func getEffectiveEncryptionEnabled(d interface{}, attrEncryptionEnabled, attrIsEncryptionEnabled string) bool {
	type getter interface {
		Get(string) interface{}
		GetOk(string) (interface{}, bool)
	}

	if g, ok := d.(getter); ok {
		// Prefer V3 field
		if val, ok := g.GetOk(attrEncryptionEnabled); ok {
			if v, ok := val.(bool); ok {
				return v
			}
		}
		// Fall back to V2 field
		if val, ok := g.GetOk(attrIsEncryptionEnabled); ok {
			if v, ok := val.(bool); ok {
				return v
			}
		}
	}

	return false
}

// logDeprecationWarning logs a deprecation warning for old field usage
func logDeprecationWarning(oldField, newField string) {
	log.Printf("[WARN] The '%s' field is deprecated and will be removed in v4.0. Please use '%s' instead.", oldField, newField)
}
