package autoscaling_test

import (
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/autoscaling"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	"github.com/stretchr/testify/assert"
)

func TestMarshalScalerGroupCreateRequestForLog_RedactsEncryptionPassphrase(t *testing.T) {
	req := &goe2e.ScalerGroupCreateRequest{
		Name:                 "example",
		IsEncryptionEnabled:  true,
		EncryptionPassphrase: "super-secret-passphrase",
	}

	got := autoscaling.MarshalScalerGroupCreateRequestForLog(req)
	assert.NotContains(t, string(got), "super-secret-passphrase", "log payload must not contain raw secrets")
	assert.Contains(t, string(got), "[REDACTED]", "log payload must include redacted marker when secret is present")
}

func TestMarshalScalerGroupCreateRequestForLog_NilRequest(t *testing.T) {
	got := autoscaling.MarshalScalerGroupCreateRequestForLog(nil)
	assert.Equal(t, "null", string(got))
}
