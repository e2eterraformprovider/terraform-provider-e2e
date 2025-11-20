package models

import (
	"encoding/json"
	"testing"
)

func TestVpcMarshalling(t *testing.T) {
	vpc := Vpc{
		Created_at: "2024-01-01T00:00:00Z",
		State:      "ACTIVE",
		Name:       "test-vpc",
		Ipv4_cidr:  "10.0.0.0/24",
		Network_id: 123,
		Gateway_ip: "10.0.0.1",
		Pool_size:  254,
		Is_active:  true,
	}

	// Test marshaling
	jsonData, err := json.Marshal(vpc)
	if err != nil {
		t.Fatalf("Failed to marshal VPC: %v", err)
	}

	// Test unmarshaling
	var unmarshaledVpc Vpc
	err = json.Unmarshal(jsonData, &unmarshaledVpc)
	if err != nil {
		t.Fatalf("Failed to unmarshal VPC: %v", err)
	}

	// Verify fields
	if unmarshaledVpc.Name != vpc.Name {
		t.Errorf("Expected Name %s, got %s", vpc.Name, unmarshaledVpc.Name)
	}
	if unmarshaledVpc.Ipv4_cidr != vpc.Ipv4_cidr {
		t.Errorf("Expected Ipv4_cidr %s, got %s", vpc.Ipv4_cidr, unmarshaledVpc.Ipv4_cidr)
	}
	if unmarshaledVpc.Network_id != vpc.Network_id {
		t.Errorf("Expected Network_id %f, got %f", vpc.Network_id, unmarshaledVpc.Network_id)
	}
	if unmarshaledVpc.State != vpc.State {
		t.Errorf("Expected State %s, got %s", vpc.State, unmarshaledVpc.State)
	}
	if unmarshaledVpc.Is_active != vpc.Is_active {
		t.Errorf("Expected Is_active %v, got %v", vpc.Is_active, unmarshaledVpc.Is_active)
	}
}

func TestVpcResponseMarshalling(t *testing.T) {
	vpcResponse := VpcResponse{
		Code:    200,
		Message: "success",
		Data: Vpc{
			Name:       "test-vpc",
			Ipv4_cidr:  "10.0.0.0/24",
			Network_id: 123,
			State:      "ACTIVE",
		},
		Error: []interface{}{},
	}

	// Test marshaling
	jsonData, err := json.Marshal(vpcResponse)
	if err != nil {
		t.Fatalf("Failed to marshal VpcResponse: %v", err)
	}

	// Test unmarshaling
	var unmarshaledResponse VpcResponse
	err = json.Unmarshal(jsonData, &unmarshaledResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal VpcResponse: %v", err)
	}

	// Verify fields
	if unmarshaledResponse.Code != vpcResponse.Code {
		t.Errorf("Expected Code %d, got %d", vpcResponse.Code, unmarshaledResponse.Code)
	}
	if unmarshaledResponse.Message != vpcResponse.Message {
		t.Errorf("Expected Message %s, got %s", vpcResponse.Message, unmarshaledResponse.Message)
	}
	if unmarshaledResponse.Data.Name != vpcResponse.Data.Name {
		t.Errorf("Expected VPC Name %s, got %s", vpcResponse.Data.Name, unmarshaledResponse.Data.Name)
	}
}

func TestVpcsResponseMarshalling(t *testing.T) {
	vpcsResponse := VpcsResponse{
		Code:    200,
		Message: "success",
		Data: []Vpc{
			{
				Name:       "test-vpc-1",
				Ipv4_cidr:  "10.0.0.0/24",
				Network_id: 123,
				State:      "ACTIVE",
			},
			{
				Name:       "test-vpc-2",
				Ipv4_cidr:  "10.0.1.0/24",
				Network_id: 124,
				State:      "ACTIVE",
			},
		},
		Error: []interface{}{},
	}

	// Test marshaling
	jsonData, err := json.Marshal(vpcsResponse)
	if err != nil {
		t.Fatalf("Failed to marshal VpcsResponse: %v", err)
	}

	// Test unmarshaling
	var unmarshaledResponse VpcsResponse
	err = json.Unmarshal(jsonData, &unmarshaledResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal VpcsResponse: %v", err)
	}

	// Verify fields
	if unmarshaledResponse.Code != vpcsResponse.Code {
		t.Errorf("Expected Code %d, got %d", vpcsResponse.Code, unmarshaledResponse.Code)
	}
	if len(unmarshaledResponse.Data) != len(vpcsResponse.Data) {
		t.Errorf("Expected %d VPCs, got %d", len(vpcsResponse.Data), len(unmarshaledResponse.Data))
	}
	if len(unmarshaledResponse.Data) > 0 {
		if unmarshaledResponse.Data[0].Name != vpcsResponse.Data[0].Name {
			t.Errorf("Expected first VPC name %s, got %s", vpcsResponse.Data[0].Name, unmarshaledResponse.Data[0].Name)
		}
	}
}

func TestVpcCreateMarshalling(t *testing.T) {
	vpcCreate := VpcCreate{
		IPv4:     "10.0.0.0/24",
		IsE2EVpc: true,
		VpcName:  "test-vpc",
	}

	// Test marshaling
	jsonData, err := json.Marshal(vpcCreate)
	if err != nil {
		t.Fatalf("Failed to marshal VpcCreate: %v", err)
	}

	// Test unmarshaling
	var unmarshaledVpcCreate VpcCreate
	err = json.Unmarshal(jsonData, &unmarshaledVpcCreate)
	if err != nil {
		t.Fatalf("Failed to unmarshal VpcCreate: %v", err)
	}

	// Verify fields
	if unmarshaledVpcCreate.VpcName != vpcCreate.VpcName {
		t.Errorf("Expected VpcName %s, got %s", vpcCreate.VpcName, unmarshaledVpcCreate.VpcName)
	}
	if unmarshaledVpcCreate.IPv4 != vpcCreate.IPv4 {
		t.Errorf("Expected IPv4 %s, got %s", vpcCreate.IPv4, unmarshaledVpcCreate.IPv4)
	}
	if unmarshaledVpcCreate.IsE2EVpc != vpcCreate.IsE2EVpc {
		t.Errorf("Expected IsE2EVpc %v, got %v", vpcCreate.IsE2EVpc, unmarshaledVpcCreate.IsE2EVpc)
	}
}

func TestVpcJSONTags(t *testing.T) {
	jsonStr := `{
		"created_at": "2024-01-01T00:00:00Z",
		"state": "ACTIVE",
		"name": "test-vpc",
		"ipv4_cidr": "10.0.0.0/24",
		"network_id": 123,
		"gateway_ip": "10.0.0.1",
		"pool_size": 254,
		"is_active": true
	}`

	var vpc Vpc
	err := json.Unmarshal([]byte(jsonStr), &vpc)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if vpc.Name != "test-vpc" {
		t.Errorf("Expected Name test-vpc, got %s", vpc.Name)
	}
	if vpc.Ipv4_cidr != "10.0.0.0/24" {
		t.Errorf("Expected Ipv4_cidr 10.0.0.0/24, got %s", vpc.Ipv4_cidr)
	}
	if vpc.Network_id != 123 {
		t.Errorf("Expected Network_id 123, got %f", vpc.Network_id)
	}
}
