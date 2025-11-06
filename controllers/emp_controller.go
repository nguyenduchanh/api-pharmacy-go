package controllers

import (
	"api-pharmacy-go/config"
	dto "api-pharmacy-go/dto"
	"api-pharmacy-go/middleware"
	"api-pharmacy-go/models"
	"api-pharmacy-go/response"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

func GetEmps(c *gin.Context) {
	var emps []models.MEmp
	config.DB.Find(&emps)
	c.JSON(http.StatusOK, emps)
}
func GetEmp(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "ID không hợp lệ")
		return
	}
	var emp models.MEmp
	if err := config.DB.First(&emp, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "Không tìm thấy nhân viên")
		} else {
			response.InternalServerError(c, "Lỗi khi truy vấn cơ sở dữ liệu: "+err.Error())
		}
		return
	}

	response.OK(c, "Lấy thông tin nhân viên thành công", emp)
}
func CreateEmp(c *gin.Context) {
	userInfo, _ := middleware.DecodeTokenFromHeader(c)
	var input dto.EmpDto
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "Dữ liệu không hợp lệ: "+err.Error())
		return
	}
	var emp models.MEmp
	if err := config.DB.First(&emp, input.EmpPhone).Error; err != nil {
		response.Conflict(c, "Đã tồn tại nhân viên với số điện thoại trên")
		return
	}
	currentTime := time.Now()
	orgId := uint64(userInfo.OrgId)
	var newEmp = models.MEmp{
		Department:  input.Department,
		EmpAddress:  input.EmpAddress,
		EmpEmail:    input.EmpEmail,
		EmpName:     input.EmpName,
		EmpPhone:    input.EmpPhone,
		OrgID:       orgId,
		CreateBy:    &userInfo.Username,
		CreateDate:  &currentTime,
		UpdatedDate: &currentTime,
		UpdatedBy:   &userInfo.Username,
	}
	if err := config.DB.Create(&newEmp).Error; err != nil {
		response.BadRequest(c, "Không thể tạo mới nhân viên: "+err.Error())
		return
	}
	response.Created(c, "Tạo nhân viên thành công", newEmp)
}

func UpdateEmp(c *gin.Context) {
	userInfo, _ := middleware.DecodeTokenFromHeader(c)
	id := c.Param("id")
	var input dto.UpdateEmpDto
	var emp models.MEmp
	if err := config.DB.First(&emp, id).Error; err != nil {
		response.NotFound(c, "Không tìm thấy nhân viên")
		return
	}
	c.ShouldBindJSON(&input)
	currentTime := time.Now()
	newRole := models.MEmp{
		ID:          input.Id,
		Department:  input.Department,
		EmpAddress:  input.EmpAddress,
		EmpEmail:    input.EmpEmail,
		EmpName:     input.EmpName,
		EmpPhone:    input.EmpPhone,
		CreateBy:    emp.CreateBy,
		CreateDate:  emp.CreateDate,
		UpdatedBy:   &userInfo.Username,
		UpdatedDate: &currentTime,
	}
	config.DB.Save(&newRole)
	response.OK(c, "Cập nhật nhân viên thành công", newRole)
}
func DeleteEmp(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.First(&models.MEmp{}, id).Error; err != nil {
		response.NotFound(c, "Không tìm thấy nhân viên")
		return
	}
	if er := config.DB.Delete(&models.MEmp{}, id).Error; er != nil {
		response.InternalServerError(c, "Lỗi khi truy vấn cơ sở dữ liệu: "+er.Error())
	}
	response.OK(c, "Xóa nhân viên thành công", nil)
}
