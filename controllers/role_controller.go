package controllers

import (
	"api-pharmacy-go/config"
	role_dto "api-pharmacy-go/dto"
	"api-pharmacy-go/middleware"
	"api-pharmacy-go/models"
	"api-pharmacy-go/response"
	"errors"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetRoles(c *gin.Context) {
	username, _ := c.Get("username")
	println(username)
	var roles []models.MRole
	config.DB.Find(&roles)
	c.JSON(http.StatusOK, roles)
}
func GetRole(c *gin.Context) {
	//c.JSON(http.StatusOK, role)
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "ID không hợp lệ")
		return
	}
	var role models.MRole
	if err := config.DB.First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "Không tìm thấy quyền")
		} else {
			response.InternalServerError(c, "Lỗi khi truy vấn cơ sở dữ liệu: "+err.Error())
		}
		return
	}

	response.OK(c, "Lấy thông tin quyền thành công", role)
}
func CreateRole(c *gin.Context) {
	userInfo, _ := middleware.DecodeTokenFromHeader(c)
	println(userInfo.Username)
	var input role_dto.RoleDto
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "Dữ liệu không hợp lệ: "+err.Error())
		return
	}
	currentTime := time.Now()
	role := models.MRole{
		Name:        input.Name,
		Description: input.Description,
		CreateDate:  &currentTime,
		CreateBy:    &userInfo.Username,
		UpdatedDate: &currentTime,
		UpdatedBy:   &userInfo.Username,
	}
	if err := config.DB.Create(&role).Error; err != nil {
		response.BadRequest(c, "Không thể tạo mới quyền: "+err.Error())
		return
	}
	response.Created(c, "Tạo quyền thành công", role)
}

func UpdateRole(c *gin.Context) {
	userInfo, _ := middleware.DecodeTokenFromHeader(c)
	id := c.Param("id")
	var input role_dto.UpdateRoleDto
	var role models.MRole
	if err := config.DB.First(&role, id).Error; err != nil {
		response.NotFound(c, "Không tìm thấy role")
		return
	}
	c.ShouldBindJSON(&input)
	currentTime := time.Now()
	newRole := models.MRole{
		ID:          input.Id,
		Name:        input.Name,
		Description: input.Description,
		CreateDate:  role.CreateDate,
		CreateBy:    role.CreateBy,
		UpdatedBy:   &userInfo.Username,
		UpdatedDate: &currentTime,
	}
	config.DB.Save(&newRole)
	response.OK(c, "Cập nhật quyền thành công", newRole)
}

func DeleteRole(c *gin.Context) {
	id := c.Param("id")
	config.DB.Delete(&models.MRole{}, id)
	response.OK(c, "Xóa role thành công", nil)
}
