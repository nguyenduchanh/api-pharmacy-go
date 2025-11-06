package controllers

import (
	"api-pharmacy-go/common"
	"api-pharmacy-go/config"
	user_dto "api-pharmacy-go/dto"
	"api-pharmacy-go/models"
	"api-pharmacy-go/response"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// users godoc
// @Summary Đăng ký người dùng mới
// @Description Đăng ký người dùng mới với thông tin được cung cấp
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/users [get]
func GetUsers(c *gin.Context) {
	var users []models.MUser
	if err := config.DB.Preload("UserRoles").Find(&users).Error; err != nil {
		response.InternalServerError(c, "Lỗi khi lấy danh sách người dùng")
		return
	}

	response.OK(c, "Lấy danh sách người dùng thành công", users)
}

// users godoc
// @Summary Đăng ký người dùng mới
// @Description Đăng ký người dùng mới với thông tin được cung cấp
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/user-dto [get]
func GetUser(c *gin.Context) {
	idParam := c.Param("id")
	// Chuyển param sang uint (hoặc int) — tránh lỗi khi id không hợp lệ
	id, err := strconv.ParseUint(idParam, 10, 64)
	println(idParam)
	println(id)
	if err != nil {
		response.BadRequest(c, "ID không hợp lệ")
		return
	}

	var user models.MUser
	if err := config.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "Không tìm thấy người dùng")
		} else {
			response.InternalServerError(c, "Lỗi khi truy vấn cơ sở dữ liệu: "+err.Error())
		}
		return
	}

	response.OK(c, "Lấy thông tin người dùng thành công", user)
}

func CreateUser(c *gin.Context) {
	var input user_dto.UserDto
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "Dữ liệu không hợp lệ: "+err.Error())
		return
	}
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		response.InternalServerError(c, "Không thể mã hóa mật khẩu")
		return
	}
	currentTime := time.Now()
	// Tạo user-dto object
	user := models.MUser{
		Username:  input.Username,
		Password:  string(hashedPassword),
		LastLogin: &currentTime,
		Email:     input.Email,
		EmpID:     input.EmpID,
		FullName:  input.FullName,
		OrgID:     input.OrgID,
		Phone:     input.Phone,
	}
	// Xử lý IsActive
	if input.IsActive != nil {
		user.IsActive = common.BitBool(*input.IsActive)
	} else {
		user.IsActive = common.BitBool(true) // Mặc định là active
	}
	// Bắt đầu transaction
	tx := config.DB.Begin()
	// Tạo user-dto
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), "Duplicate entry") {
			if strings.Contains(err.Error(), "username") {
				response.BadRequest(c, "Username đã tồn tại")
			} else {
				response.BadRequest(c, "Dữ liệu trùng lặp")
			}
			return
		}
		response.InternalServerError(c, "Không thể tạo người dùng: "+err.Error())
		return
	}
	// Xử lý UserRoles nếu có
	if len(input.UserRoles) > 0 {
		var userRoles []models.UserRole
		for _, role := range input.UserRoles {
			userRoles = append(userRoles, models.UserRole{
				UserID: user.ID,
				RoleID: role.RoleID,
			})
		}

		if err := tx.Create(&userRoles).Error; err != nil {
			tx.Rollback()
			response.BadRequest(c, "Không thể gán roles: "+err.Error())
			return
		}
	}
	// Commit transaction
	tx.Commit()
	// Preload user-dto với roles để trả về đầy đủ
	if err := config.DB.Preload("UserRoles").First(&user, user.ID).Error; err != nil {
		response.InternalServerError(c, "Lỗi khi tải thông tin user-dto")
		return
	}
	// Không trả về password
	user.Password = ""

	response.Created(c, "Tạo người dùng thành công", user)
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.MUser
	if err := config.DB.First(&user, id).Error; err != nil {
		//c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy user-dto"})
		response.BadRequest(c, "Dữ liệu không hợp lệ")
		return
	}
	c.ShouldBindJSON(&user)
	config.DB.Save(&user)
	//c.JSON(http.StatusOK, user-dto)
	response.OK(c, "cập nhật thông tin người dùng thành công", user)
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	config.DB.Delete(&models.MUser{}, id)
	//c.JSON(http.StatusOK, gin.H{"message": "Xóa user-dto thành công"})
	response.OK(c, "Xóa user-dto thành công", nil)
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
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "Thiếu file upload")
		return
	}

	if !isAllowedExcelFile(file) {
		response.BadRequest(c, "Chỉ được upload file Excel (.xlsx hoặc .xls)")
		return
	}

	// Lưu file tạm ra ổ đĩa
	tempPath := fmt.Sprintf("./temp_%d_%s", time.Now().UnixNano(), file.Filename)
	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		response.InternalServerError(c, "Không thể lưu file tạm: "+err.Error())
		return
	}
	defer removeFileSafe(tempPath)

	// Mở file Excel
	f, err := excelize.OpenFile(tempPath)
	if err != nil {
		response.InternalServerError(c, "Không thể đọc file Excel: "+err.Error())
		return
	}

	// ✅ Lấy sheet đầu tiên trong file
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

		user := models.MUser{
			Username: username,
			Password: password,
			IsActive: true,
		}

		if err := config.DB.Create(&user).Error; err == nil {
			createdCount++
		}
	}

	response.OK(c, fmt.Sprintf("Import thành công %d người dùng từ sheet '%s'", createdCount, firstSheet), nil)
}

// Hàm kiểm tra loại file an toàn
func isAllowedExcelFile(file *multipart.FileHeader) bool {
	// Cho phép chỉ 2 phần mở rộng Excel
	allowedExt := []string{".xlsx", ".xls"}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	fmt.Println("File extension:", ext)

	isAllowedExt := false
	for _, e := range allowedExt {
		if ext == e {
			isAllowedExt = true
			break
		}
	}
	if !isAllowedExt {
		fmt.Println("Extension không hợp lệ")
		return false
	}

	// Kiểm tra MIME thực tế
	f, err := file.Open()
	if err != nil {
		fmt.Println("Không mở được file:", err)
		return false
	}
	defer f.Close()

	buff := make([]byte, 512)
	if _, err := f.Read(buff); err != nil {
		fmt.Println("Không đọc được header:", err)
		return false
	}

	contentType := http.DetectContentType(buff)
	fmt.Println("Content-Type:", contentType)

	// Các MIME hợp lệ cho Excel
	allowedMIMEs := []string{
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"application/vnd.ms-excel",
		"application/zip", // .xlsx thường có MIME này
	}

	for _, mime := range allowedMIMEs {
		if strings.Contains(contentType, mime) {
			return true
		}
	}

	// Nếu MIME không nằm trong danh sách thì từ chối
	fmt.Println("MIME không hợp lệ:", contentType)
	return false
}

// Xóa file an toàn
func removeFileSafe(path string) {
	if err := os.Remove(path); err != nil {
		log.Printf("Không thể xóa file tạm: %v", err)
	}
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
	var users []models.MUser
	if err := config.DB.Find(&users).Error; err != nil {
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
		f.SetCellValue(sheet, "C"+string(rune(row+'0')), derefString(u.FullName))
		f.SetCellValue(sheet, "D"+string(rune(row+'0')), derefString(u.Email))
		f.SetCellValue(sheet, "E"+string(rune(row+'0')), derefString(u.Phone))
		f.SetCellValue(sheet, "F"+string(rune(row+'0')), u.IsActive)
		f.SetCellValue(sheet, "G"+string(rune(row+'0')), derefUint64(u.OrgID))
		f.SetCellValue(sheet, "H"+string(rune(row+'0')), u.LastLogin)
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=users.xlsx")
	c.Header("Content-Transfer-Encoding", "binary")
	_ = f.Write(c.Writer)
}
func derefString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func derefUint64(u *uint64) uint64 {
	if u != nil {
		return *u
	}
	return 0
}

func derefTime(t *time.Time) string {
	if t != nil {
		return t.Format("2006-01-02 15:04:05")
	}
	return ""
}
