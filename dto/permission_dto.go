package dto

import (
	"api-pharmacy-go/dto/enum"
	"time"
)

type PermissionDto struct {
	ApiEndPoint string          `json:"api_end_point"`
	Description *string         `json:"description"`
	Method      enum.HttpMethod `json:"method"`
	MenuId      uint64          `json:"menu_id"`
	OrgID       uint64          `json:"org_id"`
}
type UpdatePermissionDto struct {
	ApiEndPoint string          `json:"api_end_point"`
	Description *string         `json:"description"`
	Method      enum.HttpMethod `json:"method"`
	MenuId      uint64          `json:"menu_id"`
	OrgID       uint64          `json:"org_id"`
	CreateDate  *time.Time      `json:"create_at"`
	CreateBy    *string         `json:"create_by"`
}
