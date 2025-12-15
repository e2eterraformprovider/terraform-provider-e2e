package node

import (
	"fmt"
	"testing"

	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
)

func TestFlattenNodes(t *testing.T) {
	tests := []struct {
		name           string
		nodes          []goe2e.Node
		expectedLength int
		validateFunc   func([]interface{}) error
	}{
		{
			name:           "nil input - returns empty slice",
			nodes:          nil,
			expectedLength: 0,
			validateFunc: func(result []interface{}) error {
				if len(result) != 0 {
					return fmt.Errorf("expected empty slice for nil input, got length %d", len(result))
				}
				return nil
			},
		},
		{
			name:           "empty slice - returns empty slice",
			nodes:          []goe2e.Node{},
			expectedLength: 0,
			validateFunc: func(result []interface{}) error {
				if len(result) != 0 {
					return fmt.Errorf("expected empty slice, got length %d", len(result))
				}
				return nil
			},
		},
		{
			name: "single node - all fields present",
			nodes: []goe2e.Node{
				{
					ID:               "123",
					Name:             "test-node-1",
					IsLocked:         false,
					PrivateIPAddress: "10.0.0.1",
					PublicIPAddress:  "203.0.113.1",
					Status:           "Running",
				},
			},
			expectedLength: 1,
			validateFunc: func(result []interface{}) error {
				if len(result) != 1 {
					return fmt.Errorf("expected 1 node, got %d", len(result))
				}

				nodeMap := result[0].(map[string]interface{})

				// Check ID conversion (string to float64)
				if id, ok := nodeMap[tfconstants.AttrID].(float64); !ok || id != 123.0 {
					return fmt.Errorf("ID = %v (type %T), want 123.0 (float64)", nodeMap[tfconstants.AttrID], nodeMap[tfconstants.AttrID])
				}

				if name, ok := nodeMap[tfconstants.AttrName].(string); !ok || name != "test-node-1" {
					return fmt.Errorf("Name = %v, want test-node-1", name)
				}

				if isLocked, ok := nodeMap[tfconstants.AttrIsLocked].(bool); !ok || isLocked != false {
					return fmt.Errorf("IsLocked = %v, want false", isLocked)
				}

				if privateIP, ok := nodeMap[tfconstants.AttrPrivateIPAddress].(string); !ok || privateIP != "10.0.0.1" {
					return fmt.Errorf("PrivateIPAddress = %v, want 10.0.0.1", privateIP)
				}

				if publicIP, ok := nodeMap[tfconstants.AttrPublicIPAddress].(string); !ok || publicIP != "203.0.113.1" {
					return fmt.Errorf("PublicIPAddress = %v, want 203.0.113.1", publicIP)
				}

				if status, ok := nodeMap[tfconstants.AttrStatus].(string); !ok || status != "Running" {
					return fmt.Errorf("Status = %v, want Running", status)
				}

				// Check rescue_mode_status is always empty string
				if rescueMode, ok := nodeMap["rescue_mode_status"].(string); !ok || rescueMode != "" {
					return fmt.Errorf("rescue_mode_status = %v, want empty string", rescueMode)
				}

				return nil
			},
		},
		{
			name: "multiple nodes",
			nodes: []goe2e.Node{
				{
					ID:               "123",
					Name:             "test-node-1",
					IsLocked:         false,
					PrivateIPAddress: "10.0.0.1",
					PublicIPAddress:  "203.0.113.1",
					Status:           "Running",
				},
				{
					ID:               "456",
					Name:             "test-node-2",
					IsLocked:         true,
					PrivateIPAddress: "10.0.0.2",
					PublicIPAddress:  "203.0.113.2",
					Status:           "Stopped",
				},
			},
			expectedLength: 2,
			validateFunc: func(result []interface{}) error {
				if len(result) != 2 {
					return fmt.Errorf("expected 2 nodes, got %d", len(result))
				}

				// Validate first node
				node1 := result[0].(map[string]interface{})
				if name, _ := node1[tfconstants.AttrName].(string); name != "test-node-1" {
					return fmt.Errorf("first node name = %v, want test-node-1", name)
				}

				// Validate second node
				node2 := result[1].(map[string]interface{})
				if name, _ := node2[tfconstants.AttrName].(string); name != "test-node-2" {
					return fmt.Errorf("second node name = %v, want test-node-2", name)
				}

				return nil
			},
		},
		{
			name: "node with invalid ID - converts to 0.0",
			nodes: []goe2e.Node{
				{
					ID:               "invalid-id",
					Name:             "test-node",
					IsLocked:         false,
					PrivateIPAddress: "10.0.0.1",
					PublicIPAddress:  "203.0.113.1",
					Status:           "Running",
				},
			},
			expectedLength: 1,
			validateFunc: func(result []interface{}) error {
				nodeMap := result[0].(map[string]interface{})
				if id, ok := nodeMap[tfconstants.AttrID].(float64); !ok || id != 0.0 {
					return fmt.Errorf("ID = %v, want 0.0 for invalid ID", id)
				}
				return nil
			},
		},
		{
			name: "node with empty fields",
			nodes: []goe2e.Node{
				{
					ID:               "789",
					Name:             "",
					IsLocked:         false,
					PrivateIPAddress: "",
					PublicIPAddress:  "",
					Status:           "",
				},
			},
			expectedLength: 1,
			validateFunc: func(result []interface{}) error {
				nodeMap := result[0].(map[string]interface{})

				// All fields should be present even if empty
				if _, ok := nodeMap[tfconstants.AttrName]; !ok {
					return fmt.Errorf("Name field missing")
				}
				if _, ok := nodeMap[tfconstants.AttrPrivateIPAddress]; !ok {
					return fmt.Errorf("PrivateIPAddress field missing")
				}
				if _, ok := nodeMap[tfconstants.AttrPublicIPAddress]; !ok {
					return fmt.Errorf("PublicIPAddress field missing")
				}
				if _, ok := nodeMap[tfconstants.AttrStatus]; !ok {
					return fmt.Errorf("Status field missing")
				}
				if _, ok := nodeMap["rescue_mode_status"]; !ok {
					return fmt.Errorf("rescue_mode_status field missing")
				}

				return nil
			},
		},
		{
			name: "node with large ID number",
			nodes: []goe2e.Node{
				{
					ID:               "999999999",
					Name:             "test-node",
					IsLocked:         false,
					PrivateIPAddress: "10.0.0.1",
					PublicIPAddress:  "203.0.113.1",
					Status:           "Running",
				},
			},
			expectedLength: 1,
			validateFunc: func(result []interface{}) error {
				nodeMap := result[0].(map[string]interface{})
				if id, ok := nodeMap[tfconstants.AttrID].(float64); !ok || id != 999999999.0 {
					return fmt.Errorf("ID = %v, want 999999999.0", id)
				}
				return nil
			},
		},
		{
			name: "node with locked status",
			nodes: []goe2e.Node{
				{
					ID:               "123",
					Name:             "locked-node",
					IsLocked:         true,
					PrivateIPAddress: "10.0.0.1",
					PublicIPAddress:  "203.0.113.1",
					Status:           "Running",
				},
			},
			expectedLength: 1,
			validateFunc: func(result []interface{}) error {
				nodeMap := result[0].(map[string]interface{})
				if isLocked, ok := nodeMap[tfconstants.AttrIsLocked].(bool); !ok || isLocked != true {
					return fmt.Errorf("IsLocked = %v, want true", isLocked)
				}
				return nil
			},
		},
		{
			name: "node with decimal ID string - parses as float",
			nodes: []goe2e.Node{
				{
					ID:               "123.45",
					Name:             "test-node",
					IsLocked:         false,
					PrivateIPAddress: "10.0.0.1",
					PublicIPAddress:  "203.0.113.1",
					Status:           "Running",
				},
			},
			expectedLength: 1,
			validateFunc: func(result []interface{}) error {
				nodeMap := result[0].(map[string]interface{})
				// ParseFloat successfully parses "123.45" as 123.45
				if id, ok := nodeMap[tfconstants.AttrID].(float64); !ok || id != 123.45 {
					return fmt.Errorf("ID = %v, want 123.45 for decimal ID string", id)
				}
				return nil
			},
		},
		{
			name: "node with negative ID string - parses as negative float",
			nodes: []goe2e.Node{
				{
					ID:               "-123",
					Name:             "test-node",
					IsLocked:         false,
					PrivateIPAddress: "10.0.0.1",
					PublicIPAddress:  "203.0.113.1",
					Status:           "Running",
				},
			},
			expectedLength: 1,
			validateFunc: func(result []interface{}) error {
				nodeMap := result[0].(map[string]interface{})
				// ParseFloat will parse "-123" as -123.0
				if id, ok := nodeMap[tfconstants.AttrID].(float64); !ok || id != -123.0 {
					return fmt.Errorf("ID = %v, want -123.0", id)
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := flattenNodes(tt.nodes)

			if len(result) != tt.expectedLength {
				t.Errorf("flattenNodes() length = %v, want %v", len(result), tt.expectedLength)
				return
			}

			if tt.validateFunc != nil {
				if err := tt.validateFunc(result); err != nil {
					t.Errorf("validation failed: %v", err)
				}
			}
		})
	}
}
