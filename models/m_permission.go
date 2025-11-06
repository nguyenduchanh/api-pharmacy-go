package models

import (
	"api-pharmacy-go/dto/enum"
	"time"
)

type MPermission struct {
	ID          uint64          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ApiEndPoint string          `gorm:"column:api_end_point;type:varchar(100)" json:"api_end_point"`
	Description *string         `gorm:"column:description;type:varchar(100)" json:"description"`
	MenuId      uint64          `gorm:"column:menu_id" json:"menu_id"`
	Method      enum.HttpMethod `gorm:"column:method;type:varchar(10)" json:"method"`
	OrgID       uint64          `gorm:"column:org_id" json:"org_id"`
	CreateDate  *time.Time      `gorm:"column:create_date;type:datetime(6)" json:"create_date"`
	CreateBy    *string         `gorm:"column:create_by;type:varchar(200)" json:"create_by"`
	UpdatedDate *time.Time      `gorm:"column:updated_date;type:datetime(6)" json:"updated_date"`
	UpdatedBy   *string         `gorm:"column:updated_by;type:varchar(200)" json:"updated_by"`
}
