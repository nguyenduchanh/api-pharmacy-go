package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response là struct generic cho tất cả API responses
type Response[T any] struct {
	Data       T      `json:"data"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}
type ResponseAny struct {
	Data       interface{} `json:"data"`
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
}

// Success trả về response thành công
func Success[T any](c *gin.Context, statusCode int, message string, data T) {
	c.JSON(statusCode, Response[T]{
		Data:       data,
		StatusCode: statusCode,
		Message:    message,
	})
}

// Error trả về response lỗi
func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response[interface{}]{
		Data:       nil,
		StatusCode: statusCode,
		Message:    message,
	})
}

// Các hàm helper cho status code phổ biến

func OK(c *gin.Context, message string, data interface{}) {
	Success(c, http.StatusOK, message, data)
}

func Created(c *gin.Context, message string, data interface{}) {
	Success(c, http.StatusCreated, message, data)
}

func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message)
}

func InternalServerError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, message)
}
func RateLimit(c *gin.Context, message string) {
	Error(c, http.StatusTooManyRequests, message)
}
func Conflict(c *gin.Context, message string) {
	Error(c, http.StatusConflict, message)
}
