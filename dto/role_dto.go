package dto

import "time"

type RoleDto struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}
type UpdateRoleDto struct {
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	CreateDate  *time.Time `json:"create_at"`
	CreateBy    *string    `json:"create_by"`
}
