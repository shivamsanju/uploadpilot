package enf

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupManager(t *testing.T) *AccessManager {
	// Override the file adapter path
	os.Create("policy.csv")

	manager, err := NewAccessManager()
	if err != nil {
		t.Fatalf("Failed to create access manager: %v", err)
	}
	manager.enforcer.ClearPolicy()
	return manager
}

func TestBasic(t *testing.T) {
	t.Run("Basic Access Control", func(t *testing.T) {
		testCases := []struct {
			name         string
			setupFunc    func(*AccessManager)
			userID       string
			tenantID     string
			workspaceID  string
			roleToCheck  string
			expectAccess bool
		}{
			{
				name: "Admin has access to specific workspace",
				setupFunc: func(am *AccessManager) {
					_ = am.AddUserAccess("user1", "tenant1", "workspace1", "admin")
				},
				userID:       "user1",
				tenantID:     "tenant1",
				workspaceID:  "workspace1",
				roleToCheck:  "admin",
				expectAccess: true,
			},
			{
				name: "Admin inherits reader access",
				setupFunc: func(am *AccessManager) {
					_ = am.AddUserAccess("user1", "tenant1", "workspace1", "admin")
				},
				userID:       "user1",
				tenantID:     "tenant1",
				workspaceID:  "workspace1",
				roleToCheck:  "reader",
				expectAccess: true,
			},
			{
				name: "Admin inherits end-user access",
				setupFunc: func(am *AccessManager) {
					_ = am.AddUserAccess("user1", "tenant1", "workspace1", "admin")
				},
				userID:       "user1",
				tenantID:     "tenant1",
				workspaceID:  "workspace1",
				roleToCheck:  "end-user",
				expectAccess: true,
			},
			{
				name: "Reader has access to specific workspace",
				setupFunc: func(am *AccessManager) {
					_ = am.AddUserAccess("user2", "tenant1", "workspace1", "reader")
				},
				userID:       "user2",
				tenantID:     "tenant1",
				workspaceID:  "workspace1",
				roleToCheck:  "reader",
				expectAccess: true,
			},
			{
				name: "Reader inherits end-user access",
				setupFunc: func(am *AccessManager) {
					_ = am.AddUserAccess("user2", "tenant1", "workspace1", "reader")
				},
				userID:       "user2",
				tenantID:     "tenant1",
				workspaceID:  "workspace1",
				roleToCheck:  "end-user",
				expectAccess: true,
			},
			{
				name: "Reader does not have admin access",
				setupFunc: func(am *AccessManager) {
					_ = am.AddUserAccess("user2", "tenant1", "workspace1", "reader")
				},
				userID:       "user2",
				tenantID:     "tenant1",
				workspaceID:  "workspace1",
				roleToCheck:  "admin",
				expectAccess: false,
			},
			{
				name: "End-user has only end-user access",
				setupFunc: func(am *AccessManager) {
					_ = am.AddUserAccess("user3", "tenant1", "workspace1", "end-user")
				},
				userID:       "user3",
				tenantID:     "tenant1",
				workspaceID:  "workspace1",
				roleToCheck:  "end-user",
				expectAccess: true,
			},
			{
				name: "End-user does not have reader access",
				setupFunc: func(am *AccessManager) {
					_ = am.AddUserAccess("user3", "tenant1", "workspace1", "end-user")
				},
				userID:       "user3",
				tenantID:     "tenant1",
				workspaceID:  "workspace1",
				roleToCheck:  "reader",
				expectAccess: false,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				am := setupManager(t)
				tc.setupFunc(am)

				result := am.CheckAccess(tc.userID, tc.tenantID, tc.workspaceID, tc.roleToCheck)
				assert.Equal(t, tc.expectAccess, result, "Access check result should match expected")
			})
		}
	})
}

func TestWildcardWorkspace(t *testing.T) {

	// Test cases for wildcard workspace access
	t.Run("Tenant-Level (All Workspaces) Access", func(t *testing.T) {
		testCases := []struct {
			name         string
			setupFunc    func(*AccessManager)
			userID       string
			tenantID     string
			workspaceID  string
			roleToCheck  string
			expectAccess bool
		}{
			{
				name: "All workspaces access applies to specific workspace",
				setupFunc: func(am *AccessManager) {
					_ = am.AddUserAccess("user1", "tenant1", "*", "admin")
				},
				userID:       "user1",
				tenantID:     "tenant1",
				workspaceID:  "workspace1",
				roleToCheck:  "admin",
				expectAccess: true,
			},
			{
				name: "All workspaces access applies to different workspace",
				setupFunc: func(am *AccessManager) {
					_ = am.AddUserAccess("user1", "tenant1", "*", "admin")
				},
				userID:       "user1",
				tenantID:     "tenant1",
				workspaceID:  "workspace2",
				roleToCheck:  "admin",
				expectAccess: true,
			},
			{
				name: "Empty workspace treated as wildcard",
				setupFunc: func(am *AccessManager) {
					_ = am.AddUserAccess("user1", "tenant1", "", "admin")
				},
				userID:       "user1",
				tenantID:     "tenant1",
				workspaceID:  "workspace3",
				roleToCheck:  "admin",
				expectAccess: true,
			},
			{
				name: "Specific workspace access doesn't apply to other workspaces",
				setupFunc: func(am *AccessManager) {
					_ = am.AddUserAccess("user2", "tenant1", "workspace1", "reader")
				},
				userID:       "user2",
				tenantID:     "tenant1",
				workspaceID:  "workspace2",
				roleToCheck:  "reader",
				expectAccess: false,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				am := setupManager(t)
				tc.setupFunc(am)

				result := am.CheckAccess(tc.userID, tc.tenantID, tc.workspaceID, tc.roleToCheck)
				assert.Equal(t, tc.expectAccess, result, "Access check result should match expected")
			})
		}
	})
}

func TestCrossTenantAccess(t *testing.T) {

	// Test cases for cross-tenant access
	t.Run("Cross-Tenant Access Control", func(t *testing.T) {
		testCases := []struct {
			name         string
			setupFunc    func(*AccessManager)
			userID       string
			tenantID     string
			workspaceID  string
			roleToCheck  string
			expectAccess bool
		}{
			{
				name: "Access doesn't apply across tenants",
				setupFunc: func(am *AccessManager) {
					_ = am.AddUserAccess("user1", "tenant1", "*", "admin")
				},
				userID:       "user1",
				tenantID:     "tenant2", // Different tenant
				workspaceID:  "workspace1",
				roleToCheck:  "admin",
				expectAccess: false,
			},
			{
				name: "User with access to multiple tenants - correct tenant",
				setupFunc: func(am *AccessManager) {
					_ = am.AddUserAccess("user2", "tenant1", "*", "admin")
					_ = am.AddUserAccess("user2", "tenant2", "*", "reader")
				},
				userID:       "user2",
				tenantID:     "tenant2",
				workspaceID:  "workspace1",
				roleToCheck:  "reader",
				expectAccess: true,
			},
			{
				name: "User with access to multiple tenants - role correct",
				setupFunc: func(am *AccessManager) {
					_ = am.AddUserAccess("user3", "tenant1", "*", "admin")
					_ = am.AddUserAccess("user3", "tenant2", "*", "reader")
				},
				userID:       "user3",
				tenantID:     "tenant2",
				workspaceID:  "workspace1",
				roleToCheck:  "admin",
				expectAccess: false,
			},
			{
				name: "Different roles in different tenants",
				setupFunc: func(am *AccessManager) {
					_ = am.AddUserAccess("user4", "tenant1", "*", "admin")
					_ = am.AddUserAccess("user4", "tenant2", "*", "end-user")
				},
				userID:       "user4",
				tenantID:     "tenant2",
				workspaceID:  "workspace1",
				roleToCheck:  "reader",
				expectAccess: false,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				am := setupManager(t)
				tc.setupFunc(am)

				result := am.CheckAccess(tc.userID, tc.tenantID, tc.workspaceID, tc.roleToCheck)
				assert.Equal(t, tc.expectAccess, result, "Access check result should match expected")
			})
		}
	})
}

func TestComplexAccess(t *testing.T) {

	// Test multiple access policies and removal
	t.Run("Complex Access Patterns", func(t *testing.T) {
		testCases := []struct {
			name         string
			setupFunc    func(*AccessManager)
			checkFunc    func(*AccessManager) bool
			expectResult bool
		}{
			{
				name: "User with different roles in different workspaces",
				setupFunc: func(am *AccessManager) {
					_ = am.AddUserAccess("user1", "tenant1", "workspace1", "admin")
					_ = am.AddUserAccess("user1", "tenant1", "workspace2", "reader")
				},
				checkFunc: func(am *AccessManager) bool {
					adminAccess := am.CheckAccess("user1", "tenant1", "workspace1", "admin")
					readerAccess := am.CheckAccess("user1", "tenant1", "workspace2", "reader")
					return adminAccess && readerAccess
				},
				expectResult: true,
			},
			{
				name: "Remove specific workspace access",
				setupFunc: func(am *AccessManager) {
					_ = am.AddUserAccess("user1", "tenant1", "workspace1", "admin")
					_ = am.AddUserAccess("user1", "tenant1", "workspace2", "reader")
					_ = am.RemoveAccess("user1", "tenant1", "workspace1")
				},
				checkFunc: func(am *AccessManager) bool {
					workspace1Access := am.CheckAccess("user1", "tenant1", "workspace1", "admin")
					workspace2Access := am.CheckAccess("user1", "tenant1", "workspace2", "reader")
					return !workspace1Access && workspace2Access
				},
				expectResult: true,
			},
			{
				name: "Remove all tenant access",
				setupFunc: func(am *AccessManager) {
					_ = am.AddUserAccess("user1", "tenant1", "workspace1", "admin")
					_ = am.AddUserAccess("user1", "tenant1", "workspace2", "reader")
					_ = am.AddUserAccess("user1", "tenant2", "workspace1", "end-user")
					_ = am.RemoveAccess("user1", "tenant1", "")
				},
				checkFunc: func(am *AccessManager) bool {
					tenant1Access := am.CheckAccess("user1", "tenant1", "workspace1", "admin") ||
						am.CheckAccess("user1", "tenant1", "workspace2", "reader")
					tenant2Access := am.CheckAccess("user1", "tenant2", "workspace1", "end-user")
					return !tenant1Access && tenant2Access
				},
				expectResult: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				am := setupManager(t)
				tc.setupFunc(am)

				result := tc.checkFunc(am)
				assert.Equal(t, tc.expectResult, result, "Complex access scenario result should match expected")
			})
		}
	})
}

func TestOtherThings(t *testing.T) {
	// Test new role addition
	t.Run("Extending Role System", func(t *testing.T) {
		testCases := []struct {
			name         string
			setupFunc    func(*AccessManager)
			userID       string
			tenantID     string
			workspaceID  string
			roleToCheck  string
			expectAccess bool
		}{
			{
				name: "Add new role inheriting from reader",
				setupFunc: func(am *AccessManager) {
					_ = am.AddNewRole("reader", "manager")
					_ = am.AddUserAccess("user1", "tenant1", "workspace1", "manager")
				},
				userID:       "user1",
				tenantID:     "tenant1",
				workspaceID:  "workspace1",
				roleToCheck:  "reader",
				expectAccess: true,
			},
			{
				name: "Add new role inheriting from reader - check end-user inheritance",
				setupFunc: func(am *AccessManager) {
					_ = am.AddNewRole("reader", "manager")
					_ = am.AddUserAccess("user1", "tenant1", "workspace1", "manager")
				},
				userID:       "user1",
				tenantID:     "tenant1",
				workspaceID:  "workspace1",
				roleToCheck:  "end-user",
				expectAccess: true,
			},
			{
				name: "Add new role inheriting from reader - no admin access",
				setupFunc: func(am *AccessManager) {
					_ = am.AddNewRole("reader", "manager")
					_ = am.AddUserAccess("user1", "tenant1", "workspace1", "manager")
				},
				userID:       "user1",
				tenantID:     "tenant1",
				workspaceID:  "workspace1",
				roleToCheck:  "admin",
				expectAccess: false,
			},
			{
				name: "Chain of inheritance - super-reader -> reader -> end-user",
				setupFunc: func(am *AccessManager) {
					_ = am.AddNewRole("reader", "super-reader")
					_ = am.AddUserAccess("user1", "tenant1", "workspace1", "super-reader")
				},
				userID:       "user1",
				tenantID:     "tenant1",
				workspaceID:  "workspace1",
				roleToCheck:  "end-user",
				expectAccess: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				am := setupManager(t)
				tc.setupFunc(am)

				result := am.CheckAccess(tc.userID, tc.tenantID, tc.workspaceID, tc.roleToCheck)
				assert.Equal(t, tc.expectAccess, result, "Role inheritance check should match expected")
				os.Remove("polic.csv")
			})
		}
	})

	// Test GetUserRoles and GetUserTenantAccess
	t.Run("Query Functions", func(t *testing.T) {
		t.Run("GetUserRoles", func(t *testing.T) {
			am := setupManager(t)

			// Setup multiple access entries
			_ = am.AddUserAccess("user1", "tenant1", "workspace1", "admin")
			_ = am.AddUserAccess("user1", "tenant1", "workspace2", "reader")
			_ = am.AddUserAccess("user1", "tenant2", "workspace3", "end-user")

			roles, err := am.GetUserRoles("user1")
			assert.NoError(t, err, "GetUserRoles should not return error")
			assert.Len(t, roles, 3, "Should return 3 role entries")

			// Verify the content of the roles
			foundRoles := make(map[string]bool)
			for _, role := range roles {
				if len(role) >= 4 {
					key := role[1] + ":" + role[2] + ":" + role[3] // tenant:workspace:role
					foundRoles[key] = true
				}
			}

			assert.True(t, foundRoles["tenant1:workspace1:admin"], "Should find admin role for workspace1")
			assert.True(t, foundRoles["tenant1:workspace2:reader"], "Should find reader role for workspace2")
			assert.True(t, foundRoles["tenant2:workspace3:end-user"], "Should find end-user role for workspace3")
		})

		t.Run("GetUserTenantAccess", func(t *testing.T) {
			am := setupManager(t)

			// Setup multiple access entries
			_ = am.AddUserAccess("user1", "tenant1", "workspace1", "admin")
			_ = am.AddUserAccess("user1", "tenant1", "workspace2", "reader")
			_ = am.AddUserAccess("user1", "tenant2", "workspace3", "end-user")

			// Get access entries for tenant1
			access, err := am.GetUserTenantAccess("user1", "tenant1")
			assert.NoError(t, err, "GetUserTenantAccess should not return error")
			assert.Len(t, access, 2, "Should return 2 access entries for tenant1")

			// Verify the content
			foundWorkspaces := make(map[string]bool)
			for _, entry := range access {
				if len(entry) >= 3 {
					foundWorkspaces[entry[2]] = true // workspace
				}
			}

			assert.True(t, foundWorkspaces["workspace1"], "Should find workspace1 access")
			assert.True(t, foundWorkspaces["workspace2"], "Should find workspace2 access")
			assert.False(t, foundWorkspaces["workspace3"], "Should not find workspace3 access (different tenant)")
		})
	})

	// Test edge cases and error conditions
	t.Run("Edge Cases", func(t *testing.T) {
		testCases := []struct {
			name         string
			setupFunc    func(*AccessManager)
			checkFunc    func(*AccessManager) bool
			expectResult bool
		}{
			{
				name: "Non-existent user has no access",
				setupFunc: func(am *AccessManager) {
					// No setup needed
				},
				checkFunc: func(am *AccessManager) bool {
					return am.CheckAccess("nonexistent", "tenant1", "workspace1", "admin")
				},
				expectResult: false,
			},
			{
				name: "Access check with non-existent role",
				setupFunc: func(am *AccessManager) {
					_ = am.AddUserAccess("user1", "tenant1", "workspace1", "admin")
				},
				checkFunc: func(am *AccessManager) bool {
					return am.CheckAccess("user1", "tenant1", "workspace1", "nonexistent-role")
				},
				expectResult: false,
			},
			{
				name: "Removing non-existent access doesn't affect existing access",
				setupFunc: func(am *AccessManager) {
					_ = am.AddUserAccess("user1", "tenant1", "workspace1", "admin")
					_ = am.RemoveAccess("user1", "tenant2", "workspace1") // Non-existent access
				},
				checkFunc: func(am *AccessManager) bool {
					return am.CheckAccess("user1", "tenant1", "workspace1", "admin")
				},
				expectResult: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				am := setupManager(t)
				tc.setupFunc(am)

				result := tc.checkFunc(am)
				assert.Equal(t, tc.expectResult, result, "Edge case result should match expected")
			})
		}
	})

	// Test concurrent role assignment scenarios
	t.Run("Multiple Role Scenarios", func(t *testing.T) {
		testCases := []struct {
			name         string
			setupFunc    func(*AccessManager)
			checkFunc    func(*AccessManager) bool
			expectResult bool
		}{
			{
				name: "User with both admin and reader roles in same workspace",
				setupFunc: func(am *AccessManager) {
					_ = am.AddUserAccess("user1", "tenant1", "workspace1", "admin")
					_ = am.AddUserAccess("user1", "tenant1", "workspace1", "reader")
				},
				checkFunc: func(am *AccessManager) bool {
					// User should have both roles
					adminAccess := am.CheckAccess("user1", "tenant1", "workspace1", "admin")
					readerAccess := am.CheckAccess("user1", "tenant1", "workspace1", "reader")
					return adminAccess && readerAccess
				},
				expectResult: true,
			},
			{
				name: "Admin in one workspace, reader in another workspace",
				setupFunc: func(am *AccessManager) {
					_ = am.AddUserAccess("user1", "tenant1", "workspace1", "admin")
					_ = am.AddUserAccess("user1", "tenant1", "workspace2", "reader")
				},
				checkFunc: func(am *AccessManager) bool {
					return am.CheckAccess("user1", "tenant1", "workspace1", "reader") &&
						am.CheckAccess("user1", "tenant1", "workspace2", "admin")
				},
				expectResult: false, // Should fail because admin role has reader access too
			},
			{
				name: "Mixed specific and wildcard workspace access",
				setupFunc: func(am *AccessManager) {
					_ = am.AddUserAccess("user1", "tenant1", "*", "reader")
					_ = am.AddUserAccess("user1", "tenant1", "workspace1", "admin")
				},
				checkFunc: func(am *AccessManager) bool {
					// Should have admin access in workspace1
					ws1AdminAccess := am.CheckAccess("user1", "tenant1", "workspace1", "admin")
					// Should have reader access in any workspace including workspace2
					ws2ReaderAccess := am.CheckAccess("user1", "tenant1", "workspace2", "reader")
					// Should NOT have admin access in workspace2
					ws2AdminAccess := am.CheckAccess("user1", "tenant1", "workspace2", "admin")

					return ws1AdminAccess && ws2ReaderAccess && !ws2AdminAccess
				},
				expectResult: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				am := setupManager(t)
				tc.setupFunc(am)

				result := tc.checkFunc(am)
				assert.Equal(t, tc.expectResult, result, "Multiple role scenario result should match expected")
			})
		}
	})

	// Comprehensive Test: Cover all scenarios in a single complex setup
	t.Run("Comprehensive Access Control Test", func(t *testing.T) {
		am := setupManager(t)

		// Setup a complex access scenario
		_ = am.AddUserAccess("user1", "tenant1", "*", "admin")           // All workspaces admin in tenant1
		_ = am.AddUserAccess("user1", "tenant2", "workspace1", "reader") // Specific workspace in tenant2

		_ = am.AddUserAccess("user2", "tenant1", "workspace1", "reader")   // Specific workspace in tenant1
		_ = am.AddUserAccess("user2", "tenant1", "workspace2", "end-user") // Another workspace in tenant1

		_ = am.AddNewRole("admin", "super-admin")                             // New role inheriting from admin
		_ = am.AddUserAccess("user3", "tenant1", "workspace1", "super-admin") // User with new role

		// Tests for user1
		assert.True(t, am.CheckAccess("user1", "tenant1", "workspace1", "admin"), "user1 should have admin access in tenant1:workspace1")
		assert.True(t, am.CheckAccess("user1", "tenant1", "workspace2", "admin"), "user1 should have admin access in tenant1:workspace2")
		assert.True(t, am.CheckAccess("user1", "tenant1", "workspace3", "reader"), "user1 should have reader access in tenant1:workspace3 (inherited)")
		assert.True(t, am.CheckAccess("user1", "tenant2", "workspace1", "reader"), "user1 should have reader access in tenant2:workspace1")
		assert.False(t, am.CheckAccess("user1", "tenant2", "workspace2", "reader"), "user1 should not have reader access in tenant2:workspace2")

		// Tests for user2
		assert.True(t, am.CheckAccess("user2", "tenant1", "workspace1", "reader"), "user2 should have reader access in tenant1:workspace1")
		assert.True(t, am.CheckAccess("user2", "tenant1", "workspace1", "end-user"), "user2 should have end-user access in tenant1:workspace1 (inherited)")
		assert.True(t, am.CheckAccess("user2", "tenant1", "workspace2", "end-user"), "user2 should have end-user access in tenant1:workspace2")
		assert.False(t, am.CheckAccess("user2", "tenant1", "workspace2", "reader"), "user2 should not have reader access in tenant1:workspace2")
		assert.False(t, am.CheckAccess("user2", "tenant2", "workspace1", "end-user"), "user2 should not have any access in tenant2")

		// Tests for user3 with custom role
		assert.True(t, am.CheckAccess("user3", "tenant1", "workspace1", "super-admin"), "user3 should have super-admin access in tenant1:workspace1")
		assert.True(t, am.CheckAccess("user3", "tenant1", "workspace1", "admin"), "user3 should have admin access in tenant1:workspace1 (inherited)")
		assert.True(t, am.CheckAccess("user3", "tenant1", "workspace1", "reader"), "user3 should have reader access in tenant1:workspace1 (inherited)")
		assert.False(t, am.CheckAccess("user3", "tenant1", "workspace2", "super-admin"), "user3 should not have super-admin access in tenant1:workspace2")

		// Remove access test
		_ = am.RemoveAccess("user1", "tenant1", "") // Remove all tenant1 access
		assert.False(t, am.CheckAccess("user1", "tenant1", "workspace1", "admin"), "user1 should no longer have access in tenant1")
		assert.True(t, am.CheckAccess("user1", "tenant2", "workspace1", "reader"), "user1 should still have access in tenant2")
	})
}
