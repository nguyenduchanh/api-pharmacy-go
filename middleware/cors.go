package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func CORSMiddleware() gin.HandlerFunc {
	config := cors.Config{
		// 👉 Chỉ định danh sách domain/IP được phép truy cập API
		AllowOrigins: []string{
			"http://localhost:3000", // frontend dev
			"http://localhost:8888",
			"http://127.0.0.1:3000",
		},
		// Các phương thức HTTP được phép
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		// Các header được phép client gửi lên
		AllowHeaders: []string{
			"Origin",
			"Authorization",
			"Content-Type",
			"Accept",
			"X-Requested-With",
		},
		// Các header được phép client đọc từ response
		ExposeHeaders: []string{
			"Content-Length",
			"Authorization",
		},
		// Cho phép gửi cookie hoặc header xác thực
		AllowCredentials: true,
		// Thời gian cache kết quả preflight (OPTIONS)
		MaxAge: 12 * time.Hour,
	}
	return cors.New(config)
}
