package dbaas_mariadb

import (
	"context"

	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
)

// MockMariaDBService is a mock implementation of goe2e.MariaDBService for testing
type MockMariaDBService struct {
	CreateMariaDBFunc        func(context.Context, *goe2e.MariaDBCreateRequest) (*goe2e.MariaDB, *goe2e.Response, error)
	GetMariaDBFunc           func(context.Context, string) (*goe2e.MariaDB, *goe2e.Response, error)
	DeleteMariaDBFunc        func(context.Context, string) (*goe2e.Response, error)
	MariaDBExistsFunc        func(context.Context, string) (bool, *goe2e.Response, error)
	ShutdownMariaDBFunc      func(context.Context, string) (*goe2e.Response, error)
	ResumeMariaDBFunc        func(context.Context, string) (*goe2e.Response, error)
	RestartMariaDBFunc       func(context.Context, string) (*goe2e.Response, error)
	AttachVPCFunc            func(context.Context, string, []string) (*goe2e.Response, error)
	DetachVPCFunc            func(context.Context, string, []string) (*goe2e.Response, error)
	AttachPublicIPFunc       func(context.Context, string) (*goe2e.Response, error)
	DetachPublicIPFunc       func(context.Context, string) (*goe2e.Response, error)
	AttachParameterGroupFunc func(context.Context, string, int) (*goe2e.Response, error)
	DetachParameterGroupFunc func(context.Context, string, int) (*goe2e.Response, error)
	UpgradePlanFunc          func(context.Context, string, int) (*goe2e.Response, error)
	ExpandDiskFunc           func(context.Context, string, int) (*goe2e.Response, error)
	ExpandVPCListFunc        func(context.Context, []string) ([]goe2e.VPCMetadata, error)
	GetSoftwareIDFunc        func(context.Context, string, string) (int, error)
	GetTemplateIDFunc        func(context.Context, string, int) (int, error)
}

// CreateMariaDB implements goe2e.MariaDBService
func (m *MockMariaDBService) CreateMariaDB(ctx context.Context, req *goe2e.MariaDBCreateRequest) (*goe2e.MariaDB, *goe2e.Response, error) {
	if m.CreateMariaDBFunc != nil {
		return m.CreateMariaDBFunc(ctx, req)
	}
	return nil, nil, nil
}

// GetMariaDB implements goe2e.MariaDBService
func (m *MockMariaDBService) GetMariaDB(ctx context.Context, id string) (*goe2e.MariaDB, *goe2e.Response, error) {
	if m.GetMariaDBFunc != nil {
		return m.GetMariaDBFunc(ctx, id)
	}
	return nil, nil, nil
}

// DeleteMariaDB implements goe2e.MariaDBService
func (m *MockMariaDBService) DeleteMariaDB(ctx context.Context, id string) (*goe2e.Response, error) {
	if m.DeleteMariaDBFunc != nil {
		return m.DeleteMariaDBFunc(ctx, id)
	}
	return nil, nil
}

// MariaDBExists implements goe2e.MariaDBService
func (m *MockMariaDBService) MariaDBExists(ctx context.Context, id string) (bool, *goe2e.Response, error) {
	if m.MariaDBExistsFunc != nil {
		return m.MariaDBExistsFunc(ctx, id)
	}
	return false, nil, nil
}

// ShutdownMariaDB implements goe2e.MariaDBService
func (m *MockMariaDBService) ShutdownMariaDB(ctx context.Context, id string) (*goe2e.Response, error) {
	if m.ShutdownMariaDBFunc != nil {
		return m.ShutdownMariaDBFunc(ctx, id)
	}
	return nil, nil
}

// ResumeMariaDB implements goe2e.MariaDBService
func (m *MockMariaDBService) ResumeMariaDB(ctx context.Context, id string) (*goe2e.Response, error) {
	if m.ResumeMariaDBFunc != nil {
		return m.ResumeMariaDBFunc(ctx, id)
	}
	return nil, nil
}

// RestartMariaDB implements goe2e.MariaDBService
func (m *MockMariaDBService) RestartMariaDB(ctx context.Context, id string) (*goe2e.Response, error) {
	if m.RestartMariaDBFunc != nil {
		return m.RestartMariaDBFunc(ctx, id)
	}
	return nil, nil
}

// AttachVPC implements goe2e.MariaDBService
func (m *MockMariaDBService) AttachVPC(ctx context.Context, id string, vpcIDs []string) (*goe2e.Response, error) {
	if m.AttachVPCFunc != nil {
		return m.AttachVPCFunc(ctx, id, vpcIDs)
	}
	return nil, nil
}

// DetachVPC implements goe2e.MariaDBService
func (m *MockMariaDBService) DetachVPC(ctx context.Context, id string, vpcIDs []string) (*goe2e.Response, error) {
	if m.DetachVPCFunc != nil {
		return m.DetachVPCFunc(ctx, id, vpcIDs)
	}
	return nil, nil
}

// AttachPublicIP implements goe2e.MariaDBService
func (m *MockMariaDBService) AttachPublicIP(ctx context.Context, id string) (*goe2e.Response, error) {
	if m.AttachPublicIPFunc != nil {
		return m.AttachPublicIPFunc(ctx, id)
	}
	return nil, nil
}

// DetachPublicIP implements goe2e.MariaDBService
func (m *MockMariaDBService) DetachPublicIP(ctx context.Context, id string) (*goe2e.Response, error) {
	if m.DetachPublicIPFunc != nil {
		return m.DetachPublicIPFunc(ctx, id)
	}
	return nil, nil
}

// AttachParameterGroup implements goe2e.MariaDBService
func (m *MockMariaDBService) AttachParameterGroup(ctx context.Context, id string, pgID int) (*goe2e.Response, error) {
	if m.AttachParameterGroupFunc != nil {
		return m.AttachParameterGroupFunc(ctx, id, pgID)
	}
	return nil, nil
}

// DetachParameterGroup implements goe2e.MariaDBService
func (m *MockMariaDBService) DetachParameterGroup(ctx context.Context, id string, pgID int) (*goe2e.Response, error) {
	if m.DetachParameterGroupFunc != nil {
		return m.DetachParameterGroupFunc(ctx, id, pgID)
	}
	return nil, nil
}

// UpgradePlan implements goe2e.MariaDBService
func (m *MockMariaDBService) UpgradePlan(ctx context.Context, id string, templateID int) (*goe2e.Response, error) {
	if m.UpgradePlanFunc != nil {
		return m.UpgradePlanFunc(ctx, id, templateID)
	}
	return nil, nil
}

// ExpandDisk implements goe2e.MariaDBService
func (m *MockMariaDBService) ExpandDisk(ctx context.Context, id string, size int) (*goe2e.Response, error) {
	if m.ExpandDiskFunc != nil {
		return m.ExpandDiskFunc(ctx, id, size)
	}
	return nil, nil
}

// ExpandVPCList implements goe2e.MariaDBService
func (m *MockMariaDBService) ExpandVPCList(ctx context.Context, vpcIDs []string) ([]goe2e.VPCMetadata, error) {
	if m.ExpandVPCListFunc != nil {
		return m.ExpandVPCListFunc(ctx, vpcIDs)
	}
	return nil, nil
}

// GetSoftwareID implements goe2e.MariaDBService
func (m *MockMariaDBService) GetSoftwareID(ctx context.Context, softwareName, softwareVersion string) (int, error) {
	if m.GetSoftwareIDFunc != nil {
		return m.GetSoftwareIDFunc(ctx, softwareName, softwareVersion)
	}
	return 0, nil
}

// GetTemplateID implements goe2e.MariaDBService
func (m *MockMariaDBService) GetTemplateID(ctx context.Context, planName string, softwareID int) (int, error) {
	if m.GetTemplateIDFunc != nil {
		return m.GetTemplateIDFunc(ctx, planName, softwareID)
	}
	return 0, nil
}
