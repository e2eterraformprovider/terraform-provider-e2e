package config

import (
	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
)

// CombinedConfig holds both old and new clients during migration.
// Phase 1: FaaS uses NewClient (goe2e), all others use OldClient.
// Future phases will migrate services one by one until OldClient is removed.
type CombinedConfig struct {
	// OldClient is the legacy client - used by all services except FaaS
	OldClient *client.Client

	// NewClient is the goe2e client - currently used only by FaaS (Phase 1)
	NewClient *goe2e.Client
}
