package controllers

import (
	"api-pharmacy-go/config"
	"api-pharmacy-go/middleware"
	"api-pharmacy-go/models"
	"api-pharmacy-go/response"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Register godoc
// @Summary Đăng ký người dùng mới
// @Description Đăng ký người dùng mới với thông tin được cung cấp
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/auth/register [post]
func Register(c *gin.Context) {
	var input models.MUser
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
	input.Password = string(hash)
	input.IsActive = true

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tạo user-dto thành công"})
}

// Login godoc
// @Summary Đăng nhập
// @Description Đăng nhập với email và mật khẩu
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/auth/login [post]
func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": "Thiếu thông tin username hoặc password"})
		response.BadRequest(c, "Thiếu thông tin username hoặc password")
		return
	}

	// 🔹 Sửa lại: Query đơn giản trước
	var user models.MUser

	// 🔹 THÊM DEBUG: In ra query và kiểm tra
	//fmt.Printf("Trying to find user-dto with username: %s\n", input.Username)

	// Thực hiện query
	result := config.DB.Where("username = ?", input.Username).First(&user)

	if result.Error != nil {
		fmt.Printf("Database error: %v\n", result.Error) // 🔹 IN LỖI CHI TIẾT

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			//c.JSON(http.StatusUnauthorized, gin.H{"error": "Sai username hoặc password"})
			response.Unauthorized(c, "Sai username hoặc password")
		} else {
			// 🔹 PHÂN BIỆT CÁC LOẠI LỖI
			if strings.Contains(result.Error.Error(), "SQL syntax") {
				//c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi cấu hình database"})
				response.InternalServerError(c, "Lỗi cấu hình database")
			} else if strings.Contains(result.Error.Error(), "connection") {
				//c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi kết nối database"})
				response.InternalServerError(c, "Lỗi kết nối database")
			} else {
				//c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi hệ thống", "details": result.Error.Error()})
				response.InternalServerError(c, "Lỗi hệ thống")
			}
		}
		return
	}

	// 🔹 DEBUG: In ra thông tin user-dto tìm thấy
	//fmt.Printf("Found user-dto: ID=%d, Username=%s, Active=%t\n", user.ID, user.Username, user.IsActive)

	// 🔹 Kiểm tra trạng thái tài khoản (nếu có field Active)
	//isActive := user-dto.IsActive
	//if !(isActive) {
	//	c.JSON(http.StatusUnauthorized, gin.H{"error": "Tài khoản đã bị khóa"})
	//	return
	//}

	// 🔹 So sánh mật khẩu
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		//c.JSON(http.StatusUnauthorized, gin.H{"error": "Sai username hoặc password"})
		response.BadRequest(c, "Sai username hoặc password")
		return
	}

	// 🔹 Query lấy roles
	var roles []string
	err := config.DB.Table("m_roles r").
		Select("r.name").
		Joins("JOIN user_roles ur ON r.id = ur.role_id").
		Where("ur.user_id = ?", user.ID).
		Pluck("r.name", &roles).Error

	if err != nil {
		fmt.Printf("Error fetching roles: %v\n", err)
		roles = []string{}
	}

	// 🔹 Tạo JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"roles":    roles,
		"org_id":   user.OrgID,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	})

	tokenString, err := token.SignedString(middleware.JwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo token"})
		return
	}

	// 🔹 Cập nhật last_login (nếu có field LastLogin)
	config.DB.Model(&user).Update("last_login", time.Now())

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user-dto": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"fullname": user.FullName,
		},
		"roles":      roles,
		"expires_in": 24 * 60 * 60,
	})
}
