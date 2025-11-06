package dto

type UserDto struct {
	Username  string  `json:"username" binding:"required"`
	Password  string  `json:"password" binding:"required"`
	Email     *string `json:"email,omitempty"`
	EmpID     *uint64 `json:"emp_id,omitempty"`
	FullName  *string `json:"full_name,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
	OrgID     *uint64 `json:"org_id,omitempty"`
	Phone     *string `json:"phone,omitempty"`
	UserRoles []struct {
		RoleID uint64 `json:"role_id"`
	} `json:"user_roles,omitempty"`
}
