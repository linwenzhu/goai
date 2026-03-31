package middleware

import (
	"github.com/gin-gonic/gin"
	"goAi/pkg/jwt"
	"goAi/pkg/response"
	"strings"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从请求头取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Fail(c, 401, "请先登录")
			c.Abort()
			return
		}

		// 2. 检查格式是否是 "Bearer xxx"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Fail(c, 401, "token格式错误")
			c.Abort()
			return
		}

		// 3. 解析 token
		claims, err := jwt.ParseToken(parts[1])
		if err != nil {
			response.Fail(c, 401, "token无效或已过期")
			c.Abort()
			return
		}

		// 4. 把用户信息存入上下文，后续 handler 可以直接取
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}
