package controllers

import (
	"api-pharmacy-go/common"
	user_dto "api-pharmacy-go/dto"
	"api-pharmacy-go/middleware"
	"api-pharmacy-go/response"
	"api-pharmacy-go/services"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"strings"
	"time"
)

func GetUsers(c *gin.Context) {
	users, err := services.GetAllUsers()
	if err != nil {
		response.InternalServerError(c, "Lỗi truy vấn dữ liệu: "+err.Error())
		return
	}
	response.OK(c, "Lấy danh sách người dùng thành công", users)
}

func GetUser(c *gin.Context) {
	id, ok := common.ParseUintParam(c, "id")
	if !ok {
		return
	}
	user, err := services.GetUserById(id)
	user.Password = ""
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "Không tìm thấy người dùng")
		} else {
			response.InternalServerError(c, "Lỗi khi truy vấn cơ sở dữ liệu: "+err.Error())
		}
	}
	response.OK(c, "Lấy thông tin người dùng thành công", user)
}

func CreateUser(c *gin.Context) {
	userInfo, _ := middleware.DecodeTokenFromHeader(c)
	var input user_dto.UserDto
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "Dữ liệu không hợp lệ: "+err.Error())
		return
	}
	newUser, er := services.CreateUser(input, userInfo)
	if er != nil {
		response.BadRequest(c, "Không thể tạo mới người dùng: "+er.Error())
		return
	}
	response.Created(c, "Tạo người dùng thành công", newUser)
}

func UpdateUser(c *gin.Context) {
	userInfo, _ := middleware.DecodeTokenFromHeader(c)
	var input user_dto.UpdateUserDto
	id, ok := common.ParseUintParam(c, "id")
	if !ok {
		return
	}
	role, err := services.GetUserById(id)
	if err != nil {
		response.NotFound(c, "Không tìm thấy người dùng")
		return
	}
	input.CreateBy = role.CreateBy
	input.CreateDate = role.CreateDate
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "Dữ liệu không hợp lệ: "+err.Error())
		return
	}
	newUser, er := services.UpdateUser(role.ID, input, userInfo)
	if er != nil {
		response.BadRequest(c, "Không thể cập nhật người dùng: "+er.Error())
		return
	}
	response.OK(c, "cập nhật thông tin người dùng thành công", newUser)
}

func DeleteUser(c *gin.Context) {
	id, ok := common.ParseUintParam(c, "id")
	if !ok {
		return
	}
	if _, err := services.GetUserById(id); err != nil {
		response.NotFound(c, "Không tìm thấy người dùng")
		return
	}
	er := services.DeleteUser(id)
	if er != nil {
		response.InternalServerError(c, "Lỗi khi truy vấn cơ sở dữ liệu: "+er.Error())
	}
	response.OK(c, "Xóa người dùng thành công", nil)
}

// ImportUsers godoc
// @Summary Import danh sách người dùng từ file Excel
// @Description Cho phép upload file Excel (.xlsx hoặc .xls) để import người dùng hàng loạt
// @Tags Users
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "File Excel (.xlsx hoặc .xls)"
// @Success 200 {object} response.Response[string] "Import thành công"
// @Failure 400 {object} response.Response[string] "File không hợp lệ hoặc lỗi đọc dữ liệu"
// @Failure 500 {object} response.Response[string] "Lỗi hệ thống"
// @Router /api/users/import [post]
func ImportUsers(c *gin.Context) {
	userInfo, _ := middleware.DecodeTokenFromHeader(c)
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "Thiếu file upload")
		return
	}

	if !common.IsAllowedExcelFile(file) {
		response.BadRequest(c, "Chỉ được upload file Excel (.xlsx hoặc .xls)")
		return
	}

	// Lưu file tạm ra ổ đĩa
	tempPath := fmt.Sprintf("./temp_%d_%s", time.Now().UnixNano(), file.Filename)
	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		response.InternalServerError(c, "Không thể lưu file tạm: "+err.Error())
		return
	}
	defer common.RemoveFileSafe(tempPath)

	// Mở file Excel
	f, err := excelize.OpenFile(tempPath)
	if err != nil {
		response.InternalServerError(c, "Không thể đọc file Excel: "+err.Error())
		return
	}

	// Lấy sheet đầu tiên trong file
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		response.BadRequest(c, "File Excel không có sheet nào")
		return
	}
	firstSheet := sheets[0]
	fmt.Println("📄 Đang đọc sheet:", firstSheet)

	rows, err := f.GetRows(firstSheet)
	if err != nil {
		response.BadRequest(c, fmt.Sprintf("Không thể đọc dữ liệu từ sheet '%s': %v", firstSheet, err))
		return
	}

	createdCount := 0
	for i, row := range rows {
		if i == 0 || len(row) < 2 {
			continue // bỏ tiêu đề hoặc dòng rỗng
		}

		username := strings.TrimSpace(row[1])
		if username == "" {
			continue
		}
		password := "123456"
		userDto := user_dto.UserDto{
			Username: username,
			Password: password,
			OrgID:    &userInfo.OrgId,
			IsActive: common.BitBool(true),
		}

		_, err := services.CreateUser(userDto, userInfo)
		if err == nil {
			createdCount++
		}
	}
	response.OK(c, fmt.Sprintf("Import thành công %d người dùng từ sheet '%s'", createdCount, firstSheet), nil)
}

// ExportUsers godoc
// @Summary Xuất danh sách người dùng ra file Excel
// @Description API cho phép tải danh sách người dùng dạng Excel. Yêu cầu đăng nhập (JWT Token).
// @Tags Users
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Security BearerAuth
// @Success 200 {file} file "File Excel chứa danh sách người dùng"
// @Failure 401 {object} response.Response[string] "Không có hoặc token không hợp lệ"
// @Failure 500 {object} response.Response[string] "Lỗi hệ thống khi export dữ liệu"
// @Router /api/users/export [post]
func ExportUsers(c *gin.Context) {
	users, err := services.GetAllUsers()
	if err != nil {
		response.InternalServerError(c, "Không thể lấy dữ liệu để xuất Excel")
		return
	}
	f := excelize.NewFile()
	sheet := "Users"
	f.NewSheet(sheet)

	headers := []string{"ID", "Username", "Full Name", "Email", "Phone", "Is Active", "OrgID", "Last Login"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, u := range users {
		row := i + 2
		f.SetCellValue(sheet, "A"+string(rune(row+'0')), u.ID)
		f.SetCellValue(sheet, "B"+string(rune(row+'0')), u.Username)
		f.SetCellValue(sheet, "C"+string(rune(row+'0')), common.DerefString(u.FullName))
		f.SetCellValue(sheet, "D"+string(rune(row+'0')), common.DerefString(u.Email))
		f.SetCellValue(sheet, "E"+string(rune(row+'0')), common.DerefString(u.Phone))
		f.SetCellValue(sheet, "F"+string(rune(row+'0')), u.IsActive)
		f.SetCellValue(sheet, "G"+string(rune(row+'0')), common.DerefUint64(u.OrgID))
		f.SetCellValue(sheet, "H"+string(rune(row+'0')), u.LastLogin)
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=users.xlsx")
	c.Header("Content-Transfer-Encoding", "binary")
	_ = f.Write(c.Writer)
}
