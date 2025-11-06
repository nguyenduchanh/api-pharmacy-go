package services

import (
	"api-pharmacy-go/config"
	role_dto "api-pharmacy-go/dto"
	"api-pharmacy-go/dto/common"
	"api-pharmacy-go/models"
	"time"
)

func GetAllRoles() ([]models.MRole, error) {
	var roles []models.MRole
	err := config.DB.Find(&roles).Error
	return roles, err
}
func GetRoleById(id uint64) (models.MRole, error) {
	var role models.MRole
	err := config.DB.First(&role, id).Error
	return role, err
}
func CreateRole(dto role_dto.RoleDto, token *common.TokenClaims) (models.MRole, error) {
	currentTime := time.Now()
	role := models.MRole{
		Name:        dto.Name,
		Description: dto.Description,
		CreateDate:  &currentTime,
		CreateBy:    &token.Username,
		UpdatedDate: &currentTime,
		UpdatedBy:   &token.Username,
	}
	err := config.DB.Create(&role).Error
	return role, err
}
func UpdateRole(id uint64, dto role_dto.UpdateRoleDto, token *common.TokenClaims) (models.MRole, error) {
	currentTime := time.Now()
	currentRole, err := GetRoleById(id)
	if err != nil {
		return models.MRole{}, err
	} else {
		currentRole.UpdatedBy = &token.Username
		currentRole.UpdatedDate = &currentTime
		currentRole.Name = dto.Name
		currentRole.Description = dto.Description
		err := config.DB.Save(&currentRole).Error
		return currentRole, err
	}
}
func DeleteRole(id uint64) error {
	er := config.DB.Delete(&models.MEmp{}, id).Error
	return er
}
