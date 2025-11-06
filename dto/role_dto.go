package dto

type RoleDto struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}
type UpdateRoleDto struct {
	Id          uint64  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}
