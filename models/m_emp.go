package models

import "time"

type MEmp struct {
	ID          uint64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Department  string     `gorm:"column:email;type:varchar(100)" json:"department"`
	EmpAddress  *string    `gorm:"column:emp_address;type:varchar(100)" json:"emp_address"`
	EmpEmail    *string    `gorm:"column:emp_email;type:varchar(100)" json:"emp_email"`
	EmpName     *string    `gorm:"column:emp_name;type:varchar(100)" json:"emp_name"`
	EmpPhone    *string    `gorm:"column:emp_phone;type:varchar(100)" json:"emp_phone"`
	OrgID       uint64     `gorm:"column:org_id" json:"org_id"`
	CreateDate  *time.Time `gorm:"column:create_date;type:datetime(6)" json:"create_date"`
	CreateBy    *string    `gorm:"column:create_by;type:varchar(200)" json:"create_by"`
	UpdatedDate *time.Time `gorm:"column:updated_date;type:datetime(6)" json:"updated_date"`
	UpdatedBy   *string    `gorm:"column:updated_by;type:varchar(200)" json:"updated_by"`
}

func (MEmp) TableName() string {
	return "m_employees"
}
