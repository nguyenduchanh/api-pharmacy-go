package services

import (
	"api-pharmacy-go/config"
	"api-pharmacy-go/dto"
	"api-pharmacy-go/dto/common"
	"api-pharmacy-go/models"
	"time"
)

func GetAllPermission() ([]models.MPermission, error) {
	var permission []models.MPermission
	er := config.DB.Find(&permission).Error
	return permission, er
}
func GetPermissionById(id uint64) (models.MPermission, error) {
	var permission models.MPermission
	err := config.DB.First(&permission, id).Error
	return permission, err
}
func CreatePermission(dto dto.PermissionDto, token *common.TokenClaims) (models.MPermission, error) {
	currentTime := time.Now()
	permission := models.MPermission{
		ApiEndPoint: dto.ApiEndPoint,
		Description: dto.Description,
		Method:      dto.Method,
		MenuId:      dto.MenuId,
		OrgID:       token.OrgId,
		CreateDate:  &currentTime,
		CreateBy:    &token.Username,
		UpdatedDate: &currentTime,
		UpdatedBy:   &token.Username,
	}
	err := config.DB.Create(&permission).Error
	return permission, err
}
func UpdatePermission(id uint64, dto dto.UpdatePermissionDto, token *common.TokenClaims) (models.MPermission, error) {
	currentTime := time.Now()
	currentPermission, err := GetPermissionById(id)
	if err != nil {
		return models.MPermission{}, err
	} else {
		currentPermission.UpdatedBy = &token.Username
		currentPermission.UpdatedDate = &currentTime
		currentPermission.ApiEndPoint = dto.ApiEndPoint
		currentPermission.Description = dto.Description
		currentPermission.Method = dto.Method
		currentPermission.MenuId = dto.MenuId
		currentPermission.OrgID = dto.OrgID
		err2 := config.DB.Save(&currentPermission).Error
		return currentPermission, err2
	}
}
func DeletePermission(id uint64) error {
	er := config.DB.Delete(&models.MPermission{}, id).Error
	return er
}
