package enf

import (
	"fmt"
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
)

// AccessManager handles all access control operations
type AccessManager struct {
	enforcer *casbin.Enforcer
}

// NewAccessManager creates a new access manager
func NewAccessManager() (*AccessManager, error) {
	// Define the RBAC model with tenants and workspaces only
	m, err := model.NewModelFromString(`
[request_definition]
r = sub, tenant, workspace, role

[policy_definition]
p = sub, tenant, workspace, role

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && r.tenant == p.tenant && (p.workspace == "*" || r.workspace == p.workspace) && (g(r.role, p.role) || r.role == p.role)
`)
	if err != nil {
		return nil, fmt.Errorf("failed to create model: %w", err)
	}

	// Create the adapter
	adapter := fileadapter.NewAdapter("policy.csv")

	// Create the enforcer
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to create enforcer: %w", err)
	}

	// Set up role inheritance (admin -> reader -> end-user)
	_, err = enforcer.AddGroupingPolicy("reader", "admin")
	if err != nil {
		return nil, fmt.Errorf("failed to add admin-reader inheritance: %w", err)
	}

	_, err = enforcer.AddGroupingPolicy("end-user", "reader")
	if err != nil {
		return nil, fmt.Errorf("failed to add reader-end-user inheritance: %w", err)
	}

	return &AccessManager{enforcer: enforcer}, nil
}

// AddUserAccess grants access to a user for a specific tenant and workspace with a specific role
func (am *AccessManager) AddUserAccess(userID, tenantID, workspaceID, role string) error {
	// If workspaceID is empty string or "*", grant access to all workspaces in the tenant
	if workspaceID == "" || workspaceID == "*" {
		workspaceID = "*"
	}

	// Add policy for the user
	_, err := am.enforcer.AddPolicy(userID, tenantID, workspaceID, role)
	if err != nil {
		return fmt.Errorf("failed to add access policy: %w", err)
	}

	return am.enforcer.SavePolicy()
}

// CheckAccess verifies if a user has access to a specific workspace with a minimum required role
func (am *AccessManager) CheckAccess(userID, tenantID, workspaceID, requiredRole string) bool {
	allowed, err := am.enforcer.Enforce(userID, tenantID, workspaceID, requiredRole)
	if err != nil {
		log.Printf("Error checking access: %v", err)
		return false
	}
	return allowed
}

// RemoveAccess removes access for a user
func (am *AccessManager) RemoveAccess(userID, tenantID, workspaceID string) error {
	// Remove specific workspace access
	if workspaceID != "" && workspaceID != "*" {
		_, err := am.enforcer.RemoveFilteredPolicy(0, userID, tenantID, workspaceID)
		if err != nil {
			return fmt.Errorf("failed to remove workspace access: %w", err)
		}
	} else {
		// Remove all workspaces for this tenant
		_, err := am.enforcer.RemoveFilteredPolicy(0, userID, tenantID)
		if err != nil {
			return fmt.Errorf("failed to remove tenant access: %w", err)
		}
	}

	return am.enforcer.SavePolicy()
}

// AddNewRole adds a new role to the system with inheritance
func (am *AccessManager) AddNewRole(newRole, inheritFromRole string) error {
	_, err := am.enforcer.AddGroupingPolicy(newRole, inheritFromRole)
	if err != nil {
		return fmt.Errorf("failed to add new role: %w", err)
	}

	return am.enforcer.SavePolicy()
}

// GetUserRoles returns all roles assigned to a user across tenants and workspaces
func (am *AccessManager) GetUserRoles(userID string) ([][]string, error) {
	return am.enforcer.GetFilteredPolicy(0, userID)
}

// GetUserTenantAccess returns all access entries for a user in a specific tenant
func (am *AccessManager) GetUserTenantAccess(userID, tenantID string) ([][]string, error) {
	return am.enforcer.GetFilteredPolicy(0, userID, tenantID)
}
