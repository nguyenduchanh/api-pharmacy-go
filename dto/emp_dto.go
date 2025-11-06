package dto

type EmpDto struct {
	Department string  `json:"department"`
	EmpAddress *string `json:"emp_address"`
	EmpEmail   *string `json:"emp_email"`
	EmpName    *string `json:"emp_name"`
	EmpPhone   *string `json:"emp_phone"`
	OrgID      *uint64 `json:"org_id"`
}
type UpdateEmpDto struct {
	Id         uint64  `json:"id"`
	Department string  `json:"department"`
	EmpAddress *string `json:"emp_address"`
	EmpEmail   *string `json:"emp_email"`
	EmpName    *string `json:"emp_name"`
	EmpPhone   *string `json:"emp_phone"`
	OrgID      *uint64 `json:"org_id"`
}
