package config

import (
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestRegionSchema(t *testing.T) {
	s := RegionSchema()

	if s.Type != schema.TypeString {
		t.Errorf("Expected TypeString, got %v", s.Type)
	}
	if !s.Optional {
		t.Error("Expected Optional to be true")
	}
	if len(s.ConflictsWith) != 1 || s.ConflictsWith[0] != constants.AttrLocation {
		t.Errorf("Expected ConflictsWith to contain '%s', got %v", constants.AttrLocation, s.ConflictsWith)
	}
	if s.Description == "" {
		t.Error("Expected Description to be set")
	}
}

func TestLocationSchema(t *testing.T) {
	s := LocationSchema()

	if s.Type != schema.TypeString {
		t.Errorf("Expected TypeString, got %v", s.Type)
	}
	if !s.Optional {
		t.Error("Expected Optional to be true")
	}
	if s.Deprecated == "" {
		t.Error("Expected Deprecated message to be set")
	}
	if len(s.ConflictsWith) != 1 || s.ConflictsWith[0] != constants.AttrRegion {
		t.Errorf("Expected ConflictsWith to contain '%s', got %v", constants.AttrRegion, s.ConflictsWith)
	}
}

func TestProjectIDSchemaResource(t *testing.T) {
	s := ProjectIDSchemaResource()

	if s.Type != schema.TypeString {
		t.Errorf("Expected TypeString, got %v", s.Type)
	}
	if !s.Optional {
		t.Error("Expected Optional to be true")
	}
	if !s.ForceNew {
		t.Error("Expected ForceNew to be true for resources")
	}
	if s.Computed {
		t.Error("Expected Computed to be false for resources")
	}
	if s.Description == "" {
		t.Error("Expected Description to be set")
	}
}

func TestProjectIDSchemaComputed(t *testing.T) {
	s := ProjectIDSchemaComputed()

	if s.Type != schema.TypeString {
		t.Errorf("Expected TypeString, got %v", s.Type)
	}
	if !s.Optional {
		t.Error("Expected Optional to be true")
	}
	if !s.Computed {
		t.Error("Expected Computed to be true for data sources")
	}
	if s.ForceNew {
		t.Error("Expected ForceNew to be false for data sources")
	}
	if s.Description == "" {
		t.Error("Expected Description to be set")
	}
}

func TestSyncOnceValue(t *testing.T) {
	// Verify sync.OnceValue returns same pointer instance for RegionSchema
	schema1 := RegionSchema()
	schema2 := RegionSchema()

	if schema1 != schema2 {
		t.Error("Expected sync.OnceValue to return same instance for RegionSchema")
	}

	// Verify sync.OnceValue returns same pointer instance for LocationSchema
	schema3 := LocationSchema()
	schema4 := LocationSchema()

	if schema3 != schema4 {
		t.Error("Expected sync.OnceValue to return same instance for LocationSchema")
	}

	// Verify sync.OnceValue returns same pointer instance for ProjectIDSchemaResource
	schema5 := ProjectIDSchemaResource()
	schema6 := ProjectIDSchemaResource()

	if schema5 != schema6 {
		t.Error("Expected sync.OnceValue to return same instance for ProjectIDSchemaResource")
	}

	// Verify sync.OnceValue returns same pointer instance for ProjectIDSchemaComputed
	schema7 := ProjectIDSchemaComputed()
	schema8 := ProjectIDSchemaComputed()

	if schema7 != schema8 {
		t.Error("Expected sync.OnceValue to return same instance for ProjectIDSchemaComputed")
	}
}
