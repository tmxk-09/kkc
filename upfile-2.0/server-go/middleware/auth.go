package middleware

import (
	"strings"

	"upfile-server/util"

	"github.com/gin-gonic/gin"
)

// Auth JWT 认证中间件，等价 Java AuthInterceptor
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			util.Unauthorized(c)
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := util.ParseToken(tokenStr)
		if err != nil {
			util.Unauthorized(c)
			c.Abort()
			return
		}

		c.Set("userId", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

// GetUserID 从 gin Context 取 userId
func GetUserID(c *gin.Context) int {
	v, _ := c.Get("userId")
	if id, ok := v.(int); ok {
		return id
	}
	return 0
}

// GetUsername 从 gin Context 取 username
func GetUsername(c *gin.Context) string {
	v, _ := c.Get("username")
	if name, ok := v.(string); ok {
		return name
	}
	return ""
}

// IsAdmin 判断当前请求是否为管理员
func IsAdmin(c *gin.Context) bool {
	return GetUsername(c) == "admin"
}
