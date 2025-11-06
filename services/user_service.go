package services

import (
	"api-pharmacy-go/config"
	role_dto "api-pharmacy-go/dto"
	"api-pharmacy-go/dto/common"
	"api-pharmacy-go/models"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func GetAllUsers() ([]models.MUser, error) {
	var users []models.MUser
	err := config.DB.Find(&users).Error
	return users, err
}
func GetUserById(id uint64) (models.MUser, error) {
	var user models.MUser
	err := config.DB.First(&user, id).Error
	return user, err
}
func CreateUser(dto role_dto.UserDto, token *common.TokenClaims) (models.MUser, error) {
	currentTime := time.Now()
	hashedPassword, err2 := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err2 != nil {
		return models.MUser{}, err2
	}
	newUser := models.MUser{
		Username:    dto.Username,
		Email:       dto.Email,
		EmpID:       dto.EmpID,
		FullName:    dto.FullName,
		Password:    string(hashedPassword),
		Phone:       dto.Phone,
		IsActive:    dto.IsActive,
		CreateDate:  &currentTime,
		CreateBy:    &token.Username,
		UpdatedDate: &currentTime,
		UpdatedBy:   &token.Username,
	}

	err := config.DB.Create(&newUser).Error
	return newUser, err
}
func UpdateUser(id uint64, dto role_dto.UpdateUserDto, token *common.TokenClaims) (models.MUser, error) {
	currentTime := time.Now()
	currentUser, err := GetUserById(id)
	if err != nil {
		return models.MUser{}, err
	} else {
		currentUser.UpdatedBy = &token.Username
		currentUser.UpdatedDate = &currentTime
		currentUser.Username = dto.Username
		currentUser.Email = dto.Email
		currentUser.EmpID = dto.EmpID
		currentUser.FullName = dto.FullName
		currentUser.Phone = dto.Phone
		currentUser.OrgID = dto.OrgID
		currentUser.IsActive = dto.IsActive
		err2 := config.DB.Save(&currentUser).Error
		return currentUser, err2
	}
}
func DeleteUser(id uint64) error {
	er := config.DB.Delete(&models.MUser{}, id).Error
	return er
}
