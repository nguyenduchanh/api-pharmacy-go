package models

import (
	"api-pharmacy-go/common"
	"time"
)

type MUser struct {
	ID          uint64         `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Email       *string        `gorm:"column:email;type:varchar(100)" json:"email"`
	EmpID       *uint64        `gorm:"column:emp_id" json:"emp_id"`
	FullName    *string        `gorm:"column:full_name;type:varchar(100)" json:"full_name"`
	IsActive    common.BitBool `gorm:"column:is_active;type:bit(1)" json:"is_active"`
	LastLogin   *time.Time     `gorm:"column:last_login;type:datetime(6)" json:"last_login"`
	OrgID       *uint64        `gorm:"column:org_id" json:"org_id"`
	Password    string         `gorm:"column:password;type:varchar(100);not null" json:"password"`
	Phone       *string        `gorm:"column:phone;type:varchar(15)" json:"phone"`
	Username    string         `gorm:"column:username;type:varchar(50);not null" json:"username"`
	CreateDate  *time.Time     `gorm:"column:create_date;type:datetime(6)" json:"create_date"`
	CreateBy    *string        `gorm:"column:create_by;type:varchar(200)" json:"create_by"`
	UpdatedDate *time.Time     `gorm:"column:updated_date;type:datetime(6)" json:"updated_date"`
	UpdatedBy   *string        `gorm:"column:updated_by;type:varchar(200)" json:"updated_by"`
	UserRoles   []UserRole     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user_roles"`
}

func (MUser) TableName() string {
	return "m_users"
}
