package dto

import (
	"api-pharmacy-go/common"
	"time"
)

type UserDto struct {
	Username  string         `json:"username" binding:"required"`
	Password  string         `json:"password" binding:"required"`
	Email     *string        `json:"email"`
	EmpID     *uint64        `json:"emp_id"`
	FullName  *string        `json:"full_name"`
	IsActive  common.BitBool `json:"is_active"`
	Phone     *string        `json:"phone"`
	OrgID     *uint64        `json:"org_id"`
	UserRoles []struct {
		RoleID uint64 `json:"role_id"`
	} `json:"user_roles,omitempty"`
}
type UpdateUserDto struct {
	Username   string         `json:"username"`
	Email      *string        `json:"email"`
	EmpID      *uint64        `json:"emp_id,"`
	FullName   *string        `json:"full_name,"`
	IsActive   common.BitBool `json:"is_active"`
	OrgID      *uint64        `json:"org_id"`
	Phone      *string        `json:"phone"`
	CreateDate *time.Time     `json:"create_date"`
	CreateBy   *string        `json:"create_by"`
	UserRoles  []struct {
		RoleID uint64 `json:"role_id"`
	} `json:"user_roles,omitempty"`
}
