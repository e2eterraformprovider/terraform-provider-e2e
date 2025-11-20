package models

import (
	"encoding/json"
	"testing"
)

func TestBlockStorageCreateMarshalling(t *testing.T) {
	blockStorageCreate := BlockStorageCreate{
		Name: "test-volume",
		Size: 100.0,
		IOPS: "high",
	}

	// Test marshaling
	jsonData, err := json.Marshal(blockStorageCreate)
	if err != nil {
		t.Fatalf("Failed to marshal BlockStorageCreate: %v", err)
	}

	// Test unmarshaling
	var unmarshaledCreate BlockStorageCreate
	err = json.Unmarshal(jsonData, &unmarshaledCreate)
	if err != nil {
		t.Fatalf("Failed to unmarshal BlockStorageCreate: %v", err)
	}

	// Verify fields
	if unmarshaledCreate.Name != blockStorageCreate.Name {
		t.Errorf("Expected Name %s, got %s", blockStorageCreate.Name, unmarshaledCreate.Name)
	}
	if unmarshaledCreate.Size != blockStorageCreate.Size {
		t.Errorf("Expected Size %f, got %f", blockStorageCreate.Size, unmarshaledCreate.Size)
	}
	if unmarshaledCreate.IOPS != blockStorageCreate.IOPS {
		t.Errorf("Expected IOPS %s, got %s", blockStorageCreate.IOPS, unmarshaledCreate.IOPS)
	}
}

func TestBlockStorageUpgradeMarshalling(t *testing.T) {
	blockStorageUpgrade := BlockStorageUpgrade{
		Name:  "test-volume",
		Size:  200.0,
		VM_ID: 12345.0,
	}

	// Test marshaling
	jsonData, err := json.Marshal(blockStorageUpgrade)
	if err != nil {
		t.Fatalf("Failed to marshal BlockStorageUpgrade: %v", err)
	}

	// Test unmarshaling
	var unmarshaledUpgrade BlockStorageUpgrade
	err = json.Unmarshal(jsonData, &unmarshaledUpgrade)
	if err != nil {
		t.Fatalf("Failed to unmarshal BlockStorageUpgrade: %v", err)
	}

	// Verify fields
	if unmarshaledUpgrade.Name != blockStorageUpgrade.Name {
		t.Errorf("Expected Name %s, got %s", blockStorageUpgrade.Name, unmarshaledUpgrade.Name)
	}
	if unmarshaledUpgrade.Size != blockStorageUpgrade.Size {
		t.Errorf("Expected Size %f, got %f", blockStorageUpgrade.Size, unmarshaledUpgrade.Size)
	}
	if unmarshaledUpgrade.VM_ID != blockStorageUpgrade.VM_ID {
		t.Errorf("Expected VM_ID %f, got %f", blockStorageUpgrade.VM_ID, unmarshaledUpgrade.VM_ID)
	}
}

func TestBlockStorageMarshalling(t *testing.T) {
	blockStorage := BlockStorage{
		BlockID: 123,
		Name:    "test-volume",
		Size:    100,
		Status:  "AVAILABLE",
		Template: ResponseTemplate{
			DevPrefix:    "/dev/vdb",
			Driver:       "qcow2",
			TotalIOPSSec: "1000",
		},
		VMDetail:  map[string]interface{}{"vm_id": 456, "vm_name": "test-vm"},
		CreatedOn: "2024-01-01T00:00:00Z",
	}

	// Test marshaling
	jsonData, err := json.Marshal(blockStorage)
	if err != nil {
		t.Fatalf("Failed to marshal BlockStorage: %v", err)
	}

	// Test unmarshaling
	var unmarshaledStorage BlockStorage
	err = json.Unmarshal(jsonData, &unmarshaledStorage)
	if err != nil {
		t.Fatalf("Failed to unmarshal BlockStorage: %v", err)
	}

	// Verify fields
	if unmarshaledStorage.BlockID != blockStorage.BlockID {
		t.Errorf("Expected BlockID %d, got %d", blockStorage.BlockID, unmarshaledStorage.BlockID)
	}
	if unmarshaledStorage.Name != blockStorage.Name {
		t.Errorf("Expected Name %s, got %s", blockStorage.Name, unmarshaledStorage.Name)
	}
	if unmarshaledStorage.Size != blockStorage.Size {
		t.Errorf("Expected Size %d, got %d", blockStorage.Size, unmarshaledStorage.Size)
	}
	if unmarshaledStorage.Status != blockStorage.Status {
		t.Errorf("Expected Status %s, got %s", blockStorage.Status, unmarshaledStorage.Status)
	}
	if unmarshaledStorage.Template.Driver != blockStorage.Template.Driver {
		t.Errorf("Expected Driver %s, got %s", blockStorage.Template.Driver, unmarshaledStorage.Template.Driver)
	}
}

func TestBlockStorageResponseMarshalling(t *testing.T) {
	blockStorageResponse := BlockStorageResponse{
		Code:    200,
		Message: "success",
		Data: []BlockStorage{
			{
				BlockID: 1,
				Name:    "volume-1",
				Size:    100,
				Status:  "AVAILABLE",
			},
			{
				BlockID: 2,
				Name:    "volume-2",
				Size:    200,
				Status:  "IN-USE",
			},
		},
		Errors: map[string]interface{}{},
	}

	// Test marshaling
	jsonData, err := json.Marshal(blockStorageResponse)
	if err != nil {
		t.Fatalf("Failed to marshal BlockStorageResponse: %v", err)
	}

	// Test unmarshaling
	var unmarshaledResponse BlockStorageResponse
	err = json.Unmarshal(jsonData, &unmarshaledResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal BlockStorageResponse: %v", err)
	}

	// Verify fields
	if unmarshaledResponse.Code != blockStorageResponse.Code {
		t.Errorf("Expected Code %d, got %d", blockStorageResponse.Code, unmarshaledResponse.Code)
	}
	if len(unmarshaledResponse.Data) != len(blockStorageResponse.Data) {
		t.Errorf("Expected %d volumes, got %d", len(blockStorageResponse.Data), len(unmarshaledResponse.Data))
	}
	if len(unmarshaledResponse.Data) > 0 {
		if unmarshaledResponse.Data[0].Name != blockStorageResponse.Data[0].Name {
			t.Errorf("Expected first volume name %s, got %s", blockStorageResponse.Data[0].Name, unmarshaledResponse.Data[0].Name)
		}
	}
}

func TestBlockStorageAttachMarshalling(t *testing.T) {
	blockStorageAttach := BlockStorageAttach{
		VM_ID: 12345,
	}

	// Test marshaling
	jsonData, err := json.Marshal(blockStorageAttach)
	if err != nil {
		t.Fatalf("Failed to marshal BlockStorageAttach: %v", err)
	}

	// Test unmarshaling
	var unmarshaledAttach BlockStorageAttach
	err = json.Unmarshal(jsonData, &unmarshaledAttach)
	if err != nil {
		t.Fatalf("Failed to unmarshal BlockStorageAttach: %v", err)
	}

	// Verify fields
	if unmarshaledAttach.VM_ID != blockStorageAttach.VM_ID {
		t.Errorf("Expected VM_ID %d, got %d", blockStorageAttach.VM_ID, unmarshaledAttach.VM_ID)
	}
}

func TestErrorResponseMarshalling(t *testing.T) {
	errorResponse := ErrorResponse{
		Errors: "Volume not found",
	}

	// Test marshaling
	jsonData, err := json.Marshal(errorResponse)
	if err != nil {
		t.Fatalf("Failed to marshal ErrorResponse: %v", err)
	}

	// Test unmarshaling
	var unmarshaledError ErrorResponse
	err = json.Unmarshal(jsonData, &unmarshaledError)
	if err != nil {
		t.Fatalf("Failed to unmarshal ErrorResponse: %v", err)
	}

	// Verify fields
	if unmarshaledError.Errors != errorResponse.Errors {
		t.Errorf("Expected Errors %s, got %s", errorResponse.Errors, unmarshaledError.Errors)
	}
}

func TestResponseTemplateMarshalling(t *testing.T) {
	responseTemplate := ResponseTemplate{
		DevPrefix:    "/dev/vdc",
		Driver:       "qcow2",
		TotalIOPSSec: "2000",
	}

	// Test marshaling
	jsonData, err := json.Marshal(responseTemplate)
	if err != nil {
		t.Fatalf("Failed to marshal ResponseTemplate: %v", err)
	}

	// Test unmarshaling
	var unmarshaledTemplate ResponseTemplate
	err = json.Unmarshal(jsonData, &unmarshaledTemplate)
	if err != nil {
		t.Fatalf("Failed to unmarshal ResponseTemplate: %v", err)
	}

	// Verify fields
	if unmarshaledTemplate.DevPrefix != responseTemplate.DevPrefix {
		t.Errorf("Expected DevPrefix %s, got %s", responseTemplate.DevPrefix, unmarshaledTemplate.DevPrefix)
	}
	if unmarshaledTemplate.Driver != responseTemplate.Driver {
		t.Errorf("Expected Driver %s, got %s", responseTemplate.Driver, unmarshaledTemplate.Driver)
	}
	if unmarshaledTemplate.TotalIOPSSec != responseTemplate.TotalIOPSSec {
		t.Errorf("Expected TotalIOPSSec %s, got %s", responseTemplate.TotalIOPSSec, unmarshaledTemplate.TotalIOPSSec)
	}
}

func TestBlockStorageJSONTags(t *testing.T) {
	jsonStr := `{
		"block_id": 789,
		"name": "json-test-volume",
		"size": 150,
		"status": "AVAILABLE",
		"template": {
			"DEV_PREFIX": "/dev/vde",
			"DRIVER": "raw",
			"TOTAL_IOPS_SEC": "1500"
		},
		"vm_detail": {},
		"created_on": "2024-02-01T12:00:00Z"
	}`

	var blockStorage BlockStorage
	err := json.Unmarshal([]byte(jsonStr), &blockStorage)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if blockStorage.BlockID != 789 {
		t.Errorf("Expected BlockID 789, got %d", blockStorage.BlockID)
	}
	if blockStorage.Name != "json-test-volume" {
		t.Errorf("Expected Name json-test-volume, got %s", blockStorage.Name)
	}
	if blockStorage.Size != 150 {
		t.Errorf("Expected Size 150, got %d", blockStorage.Size)
	}
}
