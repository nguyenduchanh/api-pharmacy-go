package models

type UserRole struct {
	UserID uint64 `gorm:"column:user_id;not null" json:"user_id"`
	RoleID uint64 `gorm:"column:role_id;not null" json:"role_id"`

	User MUser `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"user-dto"`
	Role MRole `gorm:"foreignKey:RoleID;references:ID;constraint:OnDelete:CASCADE" json:"role"`
}

func (UserRole) TableName() string {
	return "user_roles"
}
