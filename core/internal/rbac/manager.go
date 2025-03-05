package rbac

import (
	"fmt"
	"log"

	"github.com/casbin/casbin/v2"
)

type AccessManager struct {
	enforcer *casbin.Enforcer
}

func NewAccessManager(adapter interface{}) (*AccessManager, error) {
	enforcer, err := casbin.NewEnforcer("internal/rbac/rbac_model.conf", adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to create enforcer: %w", err)
	}

	return &AccessManager{enforcer: enforcer}, nil
}

func (am *AccessManager) SetupPolicy() error {
	if err := am.enforcer.LoadPolicy(); err != nil {
		return fmt.Errorf("failed to load policy during setup: %w", err)
	}

	err := am.AddNewRole(Admin, Reader)
	if err != nil {
		return fmt.Errorf("failed to add admin-reader inheritance: %w", err)
	}

	err = am.AddNewRole(Reader, Uploader)
	if err != nil {
		return fmt.Errorf("failed to add reader-uploader inheritance: %w", err)
	}

	if err := am.enforcer.SavePolicy(); err != nil {
		return fmt.Errorf("failed to save policy after setup: %w", err)
	}

	return nil
}

// AddAccess grants access to a subject for a specific tenant and workspace with a specific role
func (am *AccessManager) AddAccess(sub, tenantID, workspaceID string, role AppRole) error {
	if workspaceID == "" || workspaceID == "*" {
		workspaceID = "*"
	}

	_, err := am.enforcer.AddPolicy(sub, tenantID, workspaceID, string(role))
	if err != nil {
		return fmt.Errorf("failed to add access policy: %w", err)
	}

	return am.enforcer.SavePolicy()
}

// CheckAccess verifies if a subject has access to a specific workspace with a minimum required role
func (am *AccessManager) CheckAccess(sub, tenantID, workspaceID string, requiredRole AppRole) bool {
	allowed, err := am.enforcer.Enforce(sub, tenantID, workspaceID, string(requiredRole))
	if err != nil {
		log.Printf("Error checking access: %v", err)
		return false
	}
	return allowed
}

// RemoveAccess removes access for a subject
func (am *AccessManager) RemoveAccess(sub, tenantID, workspaceID string) error {
	// Remove specific workspace access
	if workspaceID != "" && workspaceID != "*" {
		_, err := am.enforcer.RemoveFilteredPolicy(0, sub, tenantID, workspaceID)
		if err != nil {
			return fmt.Errorf("failed to remove workspace access: %w", err)
		}
	} else {
		_, err := am.enforcer.RemoveFilteredPolicy(0, sub, tenantID)
		if err != nil {
			return fmt.Errorf("failed to remove tenant access: %w", err)
		}
	}

	return am.enforcer.SavePolicy()
}

// AddNewRole adds a new role to the system with inheritance
func (am *AccessManager) AddNewRole(newRole, inheritFromRole AppRole) error {
	_, err := am.enforcer.AddGroupingPolicy(string(inheritFromRole), string(newRole))
	if err != nil {
		return fmt.Errorf("failed to add new role: %w", err)
	}

	return am.enforcer.SavePolicy()
}

// GetSubjectRoles returns all roles assigned to a subject across tenants and workspaces
func (am *AccessManager) GetSubjectRoles(sub string) ([][]string, error) {
	return am.enforcer.GetFilteredPolicy(0, sub)
}

// GetSubjectTenantAccess returns all access entries for a subject in a specific tenant
func (am *AccessManager) GetSubjectTenantAccess(sub, tenantID string) ([][]string, error) {
	return am.enforcer.GetFilteredPolicy(0, sub, tenantID)
}
