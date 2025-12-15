package dbaas_mysql

import (
	"context"

	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	"github.com/stretchr/testify/mock"
)

// MockDBaaSMySQLService is a unified mock implementation for DBaaSMySQL service
// It supports both testify mock patterns and stub modes for flexible testing
type MockDBaaSMySQLService struct {
	mock.Mock
	StubMode bool // If true, methods return "not implemented" errors without using mock.Called
}

// ExpandVPCList expands a list of VPC IDs to their metadata
func (m *MockDBaaSMySQLService) ExpandVPCList(ctx context.Context, vpcIDs []string) ([]goe2e.VPCMetadata, error) {
	args := m.Called(ctx, vpcIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]goe2e.VPCMetadata), args.Error(1)
}

// CreateCluster creates a new MySQL cluster
func (m *MockDBaaSMySQLService) CreateCluster(ctx context.Context, req *goe2e.MySQLClusterCreateRequest) (*goe2e.MySQLCluster, *goe2e.Response, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*goe2e.Response), args.Error(2)
	}
	return args.Get(0).(*goe2e.MySQLCluster), args.Get(1).(*goe2e.Response), args.Error(2)
}

// GetCluster retrieves a MySQL cluster by ID
func (m *MockDBaaSMySQLService) GetCluster(ctx context.Context, id string) (*goe2e.MySQLCluster, *goe2e.Response, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*goe2e.Response), args.Error(2)
	}
	return args.Get(0).(*goe2e.MySQLCluster), args.Get(1).(*goe2e.Response), args.Error(2)
}

// DeleteCluster deletes a MySQL cluster
func (m *MockDBaaSMySQLService) DeleteCluster(ctx context.Context, id string) (*goe2e.Response, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goe2e.Response), args.Error(1)
}

// StartCluster starts a stopped MySQL cluster
func (m *MockDBaaSMySQLService) StartCluster(ctx context.Context, id string) (*goe2e.Response, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goe2e.Response), args.Error(1)
}

// StopCluster stops a running MySQL cluster
func (m *MockDBaaSMySQLService) StopCluster(ctx context.Context, id string) (*goe2e.Response, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goe2e.Response), args.Error(1)
}

// RestartCluster restarts a MySQL cluster
func (m *MockDBaaSMySQLService) RestartCluster(ctx context.Context, id string) (*goe2e.Response, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goe2e.Response), args.Error(1)
}

// AttachVPC attaches a VPC to a MySQL cluster
func (m *MockDBaaSMySQLService) AttachVPC(ctx context.Context, id string, req *goe2e.MySQLVPCAttachRequest) (*goe2e.Response, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goe2e.Response), args.Error(1)
}

// DetachVPC detaches a VPC from a MySQL cluster
func (m *MockDBaaSMySQLService) DetachVPC(ctx context.Context, id string, req *goe2e.MySQLVPCDetachRequest) (*goe2e.Response, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goe2e.Response), args.Error(1)
}

// AttachParameterGroup attaches a parameter group to a MySQL cluster
func (m *MockDBaaSMySQLService) AttachParameterGroup(ctx context.Context, id, pgID string) (*goe2e.Response, error) {
	args := m.Called(ctx, id, pgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goe2e.Response), args.Error(1)
}

// DetachParameterGroup detaches a parameter group from a MySQL cluster
func (m *MockDBaaSMySQLService) DetachParameterGroup(ctx context.Context, id, pgID string) (*goe2e.Response, error) {
	args := m.Called(ctx, id, pgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goe2e.Response), args.Error(1)
}

// AttachPublicIP attaches a public IP to a MySQL cluster
func (m *MockDBaaSMySQLService) AttachPublicIP(ctx context.Context, id string) (*goe2e.Response, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goe2e.Response), args.Error(1)
}

// DetachPublicIP detaches a public IP from a MySQL cluster
func (m *MockDBaaSMySQLService) DetachPublicIP(ctx context.Context, id string) (*goe2e.Response, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goe2e.Response), args.Error(1)
}

// UpgradePlan upgrades a MySQL cluster's plan
func (m *MockDBaaSMySQLService) UpgradePlan(ctx context.Context, id string, req *goe2e.MySQLPlanUpgradeRequest) (*goe2e.Response, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goe2e.Response), args.Error(1)
}

// ExpandDisk expands the disk size of a MySQL cluster
func (m *MockDBaaSMySQLService) ExpandDisk(ctx context.Context, id string, req *goe2e.DiskExpansionRequest) (*goe2e.Response, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goe2e.Response), args.Error(1)
}

// GetSoftwareID retrieves the software ID for a given name and version
func (m *MockDBaaSMySQLService) GetSoftwareID(ctx context.Context, name, version string) (int, error) {
	args := m.Called(ctx, name, version)
	return args.Int(0), args.Error(1)
}

// GetTemplateID retrieves the template ID for a given plan and software ID
func (m *MockDBaaSMySQLService) GetTemplateID(ctx context.Context, plan string, softwareID int) (int, error) {
	args := m.Called(ctx, plan, softwareID)
	return args.Int(0), args.Error(1)
}
