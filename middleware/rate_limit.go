package middleware

import (
	"api-pharmacy-go/response"
	"net"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type clientData struct {
	requests []time.Time
}

var (
	rateData = make(map[string]*clientData)
	mu       sync.Mutex
)

// normalizeIP chuẩn hóa IP (chuyển IPv6 -> IPv4)
func normalizeIP(ip string) string {
	if ip == "::1" || ip == "0:0:0:0:0:0:0:1" {
		return "127.0.0.1"
	}
	host, _, err := net.SplitHostPort(ip)
	if err == nil {
		return host
	}
	return ip
}

// RateLimitMiddleware giới hạn số request trong khoảng thời gian
func RateLimitMiddleware(maxRequests int, duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := normalizeIP(c.ClientIP())
		if ip == "" {
			ip = "unknown"
		}
		now := time.Now()
		mu.Lock()
		data, exists := rateData[ip]
		if !exists {
			data = &clientData{requests: []time.Time{}}
			rateData[ip] = data
		}
		// Giữ lại các request còn trong khoảng thời gian
		valid := make([]time.Time, 0, len(data.requests))
		for _, t := range data.requests {
			if now.Sub(t) < duration {
				valid = append(valid, t)
			}
		}
		data.requests = valid
		// Nếu quá giới hạn
		if len(data.requests) >= maxRequests {
			mu.Unlock()
			response.RateLimit(c, "Bạn đã gửi quá nhiều yêu cầu. Vui lòng thử lại sau.")
			c.Abort()
			return
		}

		// Ghi nhận request hiện tại
		data.requests = append(data.requests, now)
		mu.Unlock()

		c.Next()
	}
}
