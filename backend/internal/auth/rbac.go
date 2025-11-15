package auth

const (
	RoleAdmin   = "admin"
	RoleAnalyst = "analyst"
	RoleViewer  = "viewer"
)

func CanAccessResource(userRole, requiredRole string) bool {
	roleLevel := map[string]int{
		RoleViewer:  1,
		RoleAnalyst: 2,
		RoleAdmin:   3,
	}
	
	return roleLevel[userRole] >= roleLevel[requiredRole]
}

func CanWrite(role string) bool {
	return role == RoleAdmin || role == RoleAnalyst
}

func CanManageUsers(role string) bool {
	return role == RoleAdmin
}
