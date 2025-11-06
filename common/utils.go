package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func ParseUintParam(c *gin.Context, paramName string) (uint64, bool) {
	idParam := c.Param(paramName)
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Tham số " + paramName + " không hợp lệ",
		})
		return 0, false
	}
	return id, true
}

// Hàm kiểm tra loại file an toàn
func IsAllowedExcelFile(file *multipart.FileHeader) bool {
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
func DerefString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func DerefUint64(u *uint64) uint64 {
	if u != nil {
		return *u
	}
	return 0
}

func DerefTime(t *time.Time) string {
	if t != nil {
		return t.Format("2006-01-02 15:04:05")
	}
	return ""
}

// Xóa file an toàn
func RemoveFileSafe(path string) {
	if err := os.Remove(path); err != nil {
		log.Printf("Không thể xóa file tạm: %v", err)
	}
}
