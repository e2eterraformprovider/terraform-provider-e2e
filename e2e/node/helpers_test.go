package node

import (
	"reflect"
	"testing"

	e2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
)

func TestExpandNetworkInterface(t *testing.T) {
	tests := []struct {
		name                     string
		input                    []interface{}
		expectedVPCID            string
		expectedAssignPublicIP   bool
		expectedEnableIPv6       bool
		expectedSecurityGroupIDs []int
	}{
		{
			name:                     "empty input",
			input:                    []interface{}{},
			expectedVPCID:            "",
			expectedAssignPublicIP:   false,
			expectedEnableIPv6:       false,
			expectedSecurityGroupIDs: nil,
		},
		{
			name: "full network interface",
			input: []interface{}{
				map[string]interface{}{
					e2econstants.AttrVPCID:            "vpc-123",
					e2econstants.AttrAssignPublicIP:   true,
					e2econstants.AttrEnableIPv6:       true,
					e2econstants.AttrSecurityGroupIDs: []interface{}{100, 200, 300},
				},
			},
			expectedVPCID:            "vpc-123",
			expectedAssignPublicIP:   true,
			expectedEnableIPv6:       true,
			expectedSecurityGroupIDs: []int{100, 200, 300},
		},
		{
			name: "partial network interface - vpc only",
			input: []interface{}{
				map[string]interface{}{
					e2econstants.AttrVPCID: "vpc-456",
				},
			},
			expectedVPCID:            "vpc-456",
			expectedAssignPublicIP:   false,
			expectedEnableIPv6:       false,
			expectedSecurityGroupIDs: nil,
		},
		{
			name: "partial network interface - no vpc",
			input: []interface{}{
				map[string]interface{}{
					e2econstants.AttrAssignPublicIP:   true,
					e2econstants.AttrSecurityGroupIDs: []interface{}{999},
				},
			},
			expectedVPCID:            "",
			expectedAssignPublicIP:   true,
			expectedEnableIPv6:       false,
			expectedSecurityGroupIDs: []int{999},
		},
		{
			name: "network interface with empty security groups",
			input: []interface{}{
				map[string]interface{}{
					e2econstants.AttrVPCID:            "vpc-789",
					e2econstants.AttrSecurityGroupIDs: []interface{}{},
				},
			},
			expectedVPCID:            "vpc-789",
			expectedAssignPublicIP:   false,
			expectedEnableIPv6:       false,
			expectedSecurityGroupIDs: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vpcID, assignPublicIP, enableIPv6, securityGroupIDs := expandNetworkInterface(tt.input)

			if vpcID != tt.expectedVPCID {
				t.Errorf("vpcID = %v, want %v", vpcID, tt.expectedVPCID)
			}
			if assignPublicIP != tt.expectedAssignPublicIP {
				t.Errorf("assignPublicIP = %v, want %v", assignPublicIP, tt.expectedAssignPublicIP)
			}
			if enableIPv6 != tt.expectedEnableIPv6 {
				t.Errorf("enableIPv6 = %v, want %v", enableIPv6, tt.expectedEnableIPv6)
			}
			if !reflect.DeepEqual(securityGroupIDs, tt.expectedSecurityGroupIDs) {
				t.Errorf("securityGroupIDs = %v, want %v", securityGroupIDs, tt.expectedSecurityGroupIDs)
			}
		})
	}
}

func TestFlattenNetworkInterface(t *testing.T) {
	tests := []struct {
		name             string
		vpcID            string
		publicIP         string
		ipv6Address      string
		securityGroupIDs []int
		expectedNIExists bool
		expectedVPCID    string
		expectedPublicIP bool
		expectedIPv6     bool
		expectedSGCount  int
	}{
		{
			name:             "full network interface with IPs",
			vpcID:            "vpc-123",
			publicIP:         "203.0.113.1",
			ipv6Address:      "2001:db8::1",
			securityGroupIDs: []int{100, 200},
			expectedNIExists: true,
			expectedVPCID:    "vpc-123",
			expectedPublicIP: true,
			expectedIPv6:     true,
			expectedSGCount:  2,
		},
		{
			name:             "no public IP",
			vpcID:            "vpc-456",
			publicIP:         "",
			ipv6Address:      "",
			securityGroupIDs: []int{300},
			expectedNIExists: true,
			expectedVPCID:    "vpc-456",
			expectedPublicIP: false,
			expectedIPv6:     false,
			expectedSGCount:  1,
		},
		{
			name:             "no VPC",
			vpcID:            "",
			publicIP:         "203.0.113.2",
			ipv6Address:      "",
			securityGroupIDs: nil,
			expectedNIExists: true,
			expectedVPCID:    "",
			expectedPublicIP: true,
			expectedIPv6:     false,
			expectedSGCount:  0,
		},
		{
			name:             "IPv6 only",
			vpcID:            "vpc-789",
			publicIP:         "",
			ipv6Address:      "2001:db8::2",
			securityGroupIDs: []int{},
			expectedNIExists: true,
			expectedVPCID:    "vpc-789",
			expectedPublicIP: false,
			expectedIPv6:     true,
			expectedSGCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := flattenNetworkInterface(tt.vpcID, tt.publicIP, tt.ipv6Address, tt.securityGroupIDs)

			if len(result) == 0 && tt.expectedNIExists {
				t.Errorf("Expected network interface to exist, but got empty slice")
				return
			}

			if len(result) > 0 {
				ni := result[0].(map[string]interface{})

				// Check VPC ID
				vpcID, hasVPC := ni[e2econstants.AttrVPCID]
				if tt.expectedVPCID != "" {
					if !hasVPC || vpcID != tt.expectedVPCID {
						t.Errorf("vpcID = %v, want %v", vpcID, tt.expectedVPCID)
					}
				}

				// Check assign_public_ip
				assignPublicIP := ni[e2econstants.AttrAssignPublicIP].(bool)
				if assignPublicIP != tt.expectedPublicIP {
					t.Errorf("assign_public_ip = %v, want %v", assignPublicIP, tt.expectedPublicIP)
				}

				// Check enable_ipv6
				enableIPv6 := ni[e2econstants.AttrEnableIPv6].(bool)
				if enableIPv6 != tt.expectedIPv6 {
					t.Errorf("enable_ipv6 = %v, want %v", enableIPv6, tt.expectedIPv6)
				}

				// Check security groups
				if tt.expectedSGCount > 0 {
					sgList, hasSG := ni[e2econstants.AttrSecurityGroupIDs]
					if !hasSG {
						t.Errorf("Expected security_group_ids to exist")
					} else {
						sgSlice := sgList.([]interface{})
						if len(sgSlice) != tt.expectedSGCount {
							t.Errorf("security_group_ids count = %v, want %v", len(sgSlice), tt.expectedSGCount)
						}
					}
				}
			}
		})
	}
}

func TestExpandRootVolume(t *testing.T) {
	tests := []struct {
		name               string
		input              []interface{}
		expectedSizeGB     int
		expectedVolumeType string
	}{
		{
			name:               "empty input",
			input:              []interface{}{},
			expectedSizeGB:     0,
			expectedVolumeType: "standard",
		},
		{
			name: "full root volume",
			input: []interface{}{
				map[string]interface{}{
					e2econstants.AttrSizeGB:   100,
					e2econstants.AttrDiskType: "ssd",
				},
			},
			expectedSizeGB:     100,
			expectedVolumeType: "ssd",
		},
		{
			name: "size only",
			input: []interface{}{
				map[string]interface{}{
					e2econstants.AttrSizeGB: 250,
				},
			},
			expectedSizeGB:     250,
			expectedVolumeType: "standard",
		},
		{
			name: "type only",
			input: []interface{}{
				map[string]interface{}{
					e2econstants.AttrDiskType: "nvme",
				},
			},
			expectedSizeGB:     0,
			expectedVolumeType: "nvme",
		},
		{
			name: "empty values",
			input: []interface{}{
				map[string]interface{}{},
			},
			expectedSizeGB:     0,
			expectedVolumeType: "standard",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sizeGB, volumeType := expandRootDisk(tt.input)

			if sizeGB != tt.expectedSizeGB {
				t.Errorf("sizeGB = %v, want %v", sizeGB, tt.expectedSizeGB)
			}
			if volumeType != tt.expectedVolumeType {
				t.Errorf("volumeType = %v, want %v", volumeType, tt.expectedVolumeType)
			}
		})
	}
}

func TestFlattenRootVolume(t *testing.T) {
	tests := []struct {
		name               string
		diskInfo           string
		volumeType         string
		expectedRVExists   bool
		expectedVolumeType string
	}{
		{
			name:               "with disk info",
			diskInfo:           "100 GB",
			volumeType:         "ssd",
			expectedRVExists:   true,
			expectedVolumeType: "ssd",
		},
		{
			name:               "empty disk info",
			diskInfo:           "",
			volumeType:         "standard",
			expectedRVExists:   true,
			expectedVolumeType: "standard",
		},
		{
			name:               "different formats",
			diskInfo:           "250GB SSD",
			volumeType:         "nvme",
			expectedRVExists:   true,
			expectedVolumeType: "nvme",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := flattenRootDisk(tt.diskInfo, tt.volumeType)

			if len(result) == 0 && tt.expectedRVExists {
				t.Errorf("Expected root_volume to exist, but got empty slice")
				return
			}

			if len(result) > 0 {
				rv := result[0].(map[string]interface{})

				// Check volume type
				volumeType := rv[e2econstants.AttrDiskType].(string)
				if volumeType != tt.expectedVolumeType {
					t.Errorf("volume_type = %v, want %v", volumeType, tt.expectedVolumeType)
				}

				// Check size_gb exists (even if 0)
				if _, hasSizeGB := rv[e2econstants.AttrSizeGB]; !hasSizeGB {
					t.Errorf("Expected size_gb to exist in root_volume")
				}
			}
		})
	}
}

func TestGetStartScripts(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedLength int
		expectedValue  string
	}{
		{
			name:           "empty script",
			input:          "",
			expectedLength: 0,
			expectedValue:  "",
		},
		{
			name:           "simple script",
			input:          "#!/bin/bash\necho 'hello'",
			expectedLength: 1,
			expectedValue:  "#!/bin/bash\necho 'hello'",
		},
		{
			name:           "multiline script",
			input:          "#!/bin/bash\napt-get update\napt-get install -y nginx",
			expectedLength: 1,
			expectedValue:  "#!/bin/bash\napt-get update\napt-get install -y nginx",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetStartScripts(tt.input)

			if len(result) != tt.expectedLength {
				t.Errorf("length = %v, want %v", len(result), tt.expectedLength)
			}

			if tt.expectedLength > 0 && result[0].(string) != tt.expectedValue {
				t.Errorf("value = %v, want %v", result[0], tt.expectedValue)
			}
		})
	}
}
