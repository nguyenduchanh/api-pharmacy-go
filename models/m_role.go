package models

import "time"

type MRole struct {
	ID          uint64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Description *string    `gorm:"column:description;type:varchar(200)" json:"description"`
	Name        string     `gorm:"column:name;type:varchar(50);not null" json:"name"`
	CreateDate  *time.Time `gorm:"column:create_date;type:datetime(6)" json:"create_date"`
	CreateBy    *string    `gorm:"column:create_by;type:varchar(200)" json:"create_by"`
	UpdatedDate *time.Time `gorm:"column:updated_date;type:datetime(6)" json:"updated_date"`
	UpdatedBy   *string    `gorm:"column:updated_by;type:varchar(200)" json:"updated_by"`
	UserRoles   []UserRole `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE" json:"user_roles"`
}

func (MRole) TableName() string {
	return "m_roles"
}
