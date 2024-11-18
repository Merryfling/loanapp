package middleware

import (
    "loanapp/config"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
)

// AuthMiddleware 用于验证 JWT token
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 从请求头中获取 token
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing or invalid"})
            c.Abort()
            return
        }

        // 移除 "Bearer " 前缀
        tokenString = strings.TrimPrefix(tokenString, "Bearer ")

        // 解析 token
        userId, err := config.ParseToken(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        // 将用户 ID 存储在上下文中，以便后续处理
        c.Set("userId", userId)
        c.Next()
    }
}
