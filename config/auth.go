package config

import (
    "time"
    "github.com/golang-jwt/jwt/v4"
)

// JWTSecret 是用来生成 token 的密钥
var JWTSecret = []byte("your_secret_key") // 将此密钥替换为安全的随机字符串

// GenerateToken 生成 JWT token
func GenerateToken(userID uint64) (string, error) {
    // 定义 token 的声明信息
    claims := jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(time.Hour * 1).Unix(), // 1 小时后过期
    }

    // 生成 token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(JWTSecret)
}

// ParseToken 解析并验证 token
func ParseToken(tokenStr string) (uint64, error) {
    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        // 确保采用的签名算法是 HS256
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
        }
        return JWTSecret, nil
    })

    if err != nil {
        return 0, err
    }

    // 获取用户 ID
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        userID := uint64(claims["user_id"].(float64)) // jwt数字默认float64，所以先断言然后转换
        return userID, nil
    }

    return 0, jwt.NewValidationError("invalid token", jwt.ValidationErrorSignatureInvalid)
}
