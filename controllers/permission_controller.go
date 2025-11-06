package controllers

import (
	"api-pharmacy-go/common"
	"api-pharmacy-go/dto"
	"api-pharmacy-go/middleware"
	"api-pharmacy-go/response"
	"api-pharmacy-go/services"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetPermissions(c *gin.Context) {
	permissions, err := services.GetAllPermission()
	if err != nil {
		response.InternalServerError(c, "Lỗi truy vấn dữ liệu: "+err.Error())
		return
	}
	response.OK(c, "Lấy dữ liệu quyền thành công", permissions)
}
func GetPermission(c *gin.Context) {
	id, ok := common.ParseUintParam(c, "id")
	if !ok {
		return
	}
	per, err := services.GetPermissionById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "Không tìm thấy quyền")
		} else {
			response.InternalServerError(c, "Lỗi khi truy vấn cơ sở dữ liệu: "+err.Error())
		}
	}
	response.OK(c, "Lấy thông tin quyền thành công", per)
}
func CreatePermission(c *gin.Context) {
	userInfo, _ := middleware.DecodeTokenFromHeader(c)
	var input dto.PermissionDto
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "Dữ liệu không hợp lệ: "+err.Error())
		return
	}
	newPer, er := services.CreatePermission(input, userInfo)
	if er != nil {
		response.BadRequest(c, "Không thể tạo mới quyền: "+er.Error())
		return
	}
	response.Created(c, "Tạo quyền thành công", newPer)
}

func UpdatePermission(c *gin.Context) {
	userInfo, _ := middleware.DecodeTokenFromHeader(c)
	var input dto.UpdatePermissionDto
	id, ok := common.ParseUintParam(c, "id")
	if !ok {
		return
	}
	per, err := services.GetPermissionById(id)
	if err != nil {
		response.NotFound(c, "Không tìm thấy quyền")
		return
	}
	input.CreateBy = per.CreateBy
	input.CreateDate = per.CreateDate
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "Dữ liệu không hợp lệ: "+err.Error())
		return
	}
	perRes, er := services.UpdatePermission(per.ID, input, userInfo)
	if er != nil {
		response.BadRequest(c, "Không thể cập nhật nhân viên: "+er.Error())
		return
	}
	response.OK(c, "Cập nhật nhân viên thành công", perRes)
}
func DeletePermission(c *gin.Context) {
	id, ok := common.ParseUintParam(c, "id")
	if !ok {
		return
	}
	_, er := services.GetPermissionById(id)
	if er != nil {
		response.NotFound(c, "Không tìm thấy nhân viên")
		return
	}

	er2 := services.DeleteEmp(id)
	if er2 != nil {
		response.InternalServerError(c, "Lỗi khi truy vấn cơ sở dữ liệu: "+er2.Error())
	}
	response.OK(c, "Xóa nhân viên thành công", nil)
}
