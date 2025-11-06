package services

import (
	"api-pharmacy-go/config"
	"api-pharmacy-go/dto"
	"api-pharmacy-go/dto/common"
	"api-pharmacy-go/models"
	"time"
)

func GetAllEmp() ([]models.MEmp, error) {
	var emp []models.MEmp
	err := config.DB.Find(&emp).Error
	return emp, err
}
func GetEmpById(id uint64) (models.MEmp, error) {
	var emp models.MEmp
	err := config.DB.First(&emp, id).Error
	return emp, err
}
func GetEmpByPhone(phone string) (models.MEmp, error) {
	var emp models.MEmp
	err := config.DB.First(&emp, phone).Error
	return emp, err
}
func CreateEmp(dto dto.EmpDto, token *common.TokenClaims) (models.MEmp, error) {
	currentTime := time.Now()
	role := models.MEmp{
		Department:  dto.Department,
		EmpAddress:  dto.EmpAddress,
		EmpEmail:    dto.EmpEmail,
		EmpName:     dto.EmpName,
		EmpPhone:    dto.EmpPhone,
		OrgID:       token.OrgId,
		CreateDate:  &currentTime,
		CreateBy:    &token.Username,
		UpdatedDate: &currentTime,
		UpdatedBy:   &token.Username,
	}
	err := config.DB.Create(&role).Error
	return role, err
}
func UpdateEmp(id uint64, dto dto.UpdateEmpDto, token *common.TokenClaims) (models.MEmp, error) {
	currentTime := time.Now()
	currentEmp, err := GetEmpById(id)
	if err != nil {
		return models.MEmp{}, err
	} else {
		currentEmp.UpdatedBy = &token.Username
		currentEmp.UpdatedDate = &currentTime
		currentEmp.Department = dto.Department
		currentEmp.EmpAddress = dto.EmpAddress
		currentEmp.EmpEmail = dto.EmpEmail
		currentEmp.EmpName = dto.EmpName
		currentEmp.EmpPhone = dto.EmpPhone
		currentEmp.OrgID = dto.OrgID
		err := config.DB.Save(&currentEmp).Error
		return currentEmp, err
	}
}
func DeleteEmp(id uint64) error {
	er := config.DB.Delete(&models.MEmp{}, id).Error
	return er
}
