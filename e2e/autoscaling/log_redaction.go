package autoscaling

import (
	"encoding/json"
	"fmt"

	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
)

const logRedactedValue = "[REDACTED]"

// MarshalScalerGroupCreateRequestForLog returns a JSON representation of the request that is safe
// to write to logs (i.e., secrets are redacted).
//
// NOTE: This is intentionally conservative: if new sensitive fields are added to the request
// struct in the future, they must be added to this redaction step.
func MarshalScalerGroupCreateRequestForLog(req *goe2e.ScalerGroupCreateRequest) []byte {
	if req == nil {
		return []byte("null")
	}

	// Shallow copy to avoid mutating the real request.
	safe := *req

	// Redact secrets.
	if safe.EncryptionPassphrase != "" {
		safe.EncryptionPassphrase = logRedactedValue
	}

	b, err := json.MarshalIndent(&safe, "", "  ")
	if err != nil {
		return []byte(fmt.Sprintf(`{"error":"failed to marshal scaler group request for log: %s"}`, err))
	}
	return b
}
