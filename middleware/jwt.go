package middleware

import (
	"api-pharmacy-go/response"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var JwtKey = []byte("c2VjcmV0S2V5Rm9ySldUMTYyQXVndXN0MjAyMCFAZGV2QHNlY3JldCprZXk=")

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.Unauthorized(c, "Thiếu token trong header Authorization")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})
		if err != nil || !token.Valid {
			response.Unauthorized(c, "Token không hợp lệ")
			c.Abort()
			return
		}

		c.Next()
	}
}
func DecodeTokenFromHeader(c *gin.Context) (*TokenClaims, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, errors.New("thiếu token trong header Authorization")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse token
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("phương thức ký không hợp lệ: %v", t.Header["alg"])
		}
		return JwtKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("token không hợp lệ: %v", err)
	}
	if !token.Valid {
		return nil, errors.New("token đã hết hạn hoặc không hợp lệ")
	}

	// Lấy claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("không thể đọc claims từ token")
	}

	// Lấy thông tin từ claims
	userID, _ := claims["user_id"].(float64)
	username, _ := claims["username"].(string)
	orgId, _ := claims["org_id"].(float64)
	expUnix, _ := claims["exp"].(float64)

	var roles []string
	if rawRoles, ok := claims["roles"].([]interface{}); ok {
		for _, r := range rawRoles {
			if roleStr, ok := r.(string); ok {
				roles = append(roles, roleStr)
			}
		}
	}

	return &TokenClaims{
		UserID:   int(userID),
		Username: username,
		OrgId:    int(orgId),
		Roles:    roles,
		Exp:      time.Unix(int64(expUnix), 0),
	}, nil
}

type TokenClaims struct {
	UserID   int
	Username string
	OrgId    int
	Roles    []string
	Exp      time.Time
}
