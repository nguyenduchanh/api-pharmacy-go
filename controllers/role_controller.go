package controllers

import (
	"api-pharmacy-go/common"
	role_dto "api-pharmacy-go/dto"
	"api-pharmacy-go/middleware"
	"api-pharmacy-go/response"
	"api-pharmacy-go/services"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetRoles(c *gin.Context) {
	roles, err := services.GetAllRoles()
	if err != nil {
		response.InternalServerError(c, "Lỗi truy vấn dữ liệu: "+err.Error())
		return
	}
	response.OK(c, "Lấy dữ liệu nhóm quyền thành công", roles)
}
func GetRole(c *gin.Context) {
	id, ok := common.ParseUintParam(c, "id")
	if !ok {
		return
	}
	role, err := services.GetRoleById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "Không tìm thấy nhóm quyền")
		} else {
			response.InternalServerError(c, "Lỗi khi truy vấn cơ sở dữ liệu: "+err.Error())
		}
	}
	response.OK(c, "Lấy thông tin nhóm quyền thành công", role)
}
func CreateRole(c *gin.Context) {
	userInfo, _ := middleware.DecodeTokenFromHeader(c)
	var input role_dto.RoleDto
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "Dữ liệu không hợp lệ: "+err.Error())
		return
	}
	newRole, er := services.CreateRole(input, userInfo)
	if er != nil {
		response.BadRequest(c, "Không thể tạo mới nhóm quyền: "+er.Error())
		return
	}
	response.Created(c, "Tạo nhóm quyền thành công", newRole)
}

func UpdateRole(c *gin.Context) {
	userInfo, _ := middleware.DecodeTokenFromHeader(c)
	var input role_dto.UpdateRoleDto
	id, ok := common.ParseUintParam(c, "id")
	if !ok {
		return
	}
	role, err := services.GetRoleById(id)
	if err != nil {
		response.NotFound(c, "Không tìm thấy nhóm quyền")
		return
	}
	input.CreateBy = role.CreateBy
	input.CreateDate = role.CreateDate
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "Dữ liệu không hợp lệ: "+err.Error())
		return
	}
	newRole, er := services.UpdateRole(role.ID, input, userInfo)
	if er != nil {
		response.BadRequest(c, "Không thể cập nhật nhóm quyền: "+er.Error())
		return
	}
	response.OK(c, "Cập nhật nhóm quyền thành công", newRole)
}

func DeleteRole(c *gin.Context) {
	id, ok := common.ParseUintParam(c, "id")
	if !ok {
		return
	}
	if _, err := services.GetRoleById(id); err != nil {
		response.NotFound(c, "Không tìm thấy nhóm quyền")
		return
	}
	er := services.DeleteRole(id)
	if er != nil {
		response.InternalServerError(c, "Lỗi khi truy vấn cơ sở dữ liệu: "+er.Error())
	}
	response.OK(c, "Xóa nhóm quyền thành công", nil)
}
