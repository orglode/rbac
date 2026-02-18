package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"rbac/conf"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取并验证Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			respondUnauthorized(c, "Authorization header is required")
			return
		}

		// 2. 提取token
		tokenString, err := extractBearerToken(authHeader)
		if err != nil {
			respondUnauthorized(c, err.Error())
			return
		}

		// 3. 解析JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return conf.ConfGlobal.Server.JwtToken, nil
		})

		// 4. 处理JWT错误
		if err != nil || !token.Valid {
			handleJWTError(c, err)
			return
		}

		// 5. 提取并验证user_id
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			respondUnauthorized(c, "Invalid token claims")
			return
		}

		uid, err := extractInt64UID(claims)
		if err != nil {
			respondUnauthorized(c, err.Error())
			return
		}
		if uid <= 0 {
			respondUnauthorized(c, "Please login")
			return
		}

		// 6. 将int64类型的user_id存入context
		c.Set("uid", uid)
		c.Next()
	}
}

// 从Authorization头提取Bearer token
func extractBearerToken(authHeader string) (string, error) {
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", fmt.Errorf("authorization header format must be 'Bearer {token}'")
	}
	return strings.TrimPrefix(authHeader, bearerPrefix), nil
}

// 从claims中提取int64类型的uid
func extractInt64UID(claims jwt.MapClaims) (int64, error) {
	userID, exists := claims["uid"]
	if !exists {
		return 0, fmt.Errorf("uid not found in token")
	}

	// JSON数字默认解码为float64
	floatID, ok := userID.(float64)
	if !ok {
		return 0, fmt.Errorf("uid must be a number")
	}

	// 转换为int64
	intID := int64(floatID)
	if floatID != float64(intID) {
		return 0, fmt.Errorf("uid must be an integer")
	}

	return intID, nil
}

// 统一处理未授权响应
func respondUnauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"code":     http.StatusUnauthorized,
		"msg":      message,
		"trace_id": c.GetString("trace_id"),
	})
	c.Abort()
}

// 处理JWT特定错误
func handleJWTError(c *gin.Context, err error) {
	var validationError *jwt.ValidationError
	if errors.As(err, &validationError) {
		switch {
		case validationError.Errors&jwt.ValidationErrorMalformed != 0:
			respondUnauthorized(c, "Token is malformed")
		case validationError.Errors&jwt.ValidationErrorExpired != 0:
			respondUnauthorized(c, "Token is expired")
		case validationError.Errors&jwt.ValidationErrorNotValidYet != 0:
			respondUnauthorized(c, "Token is not active yet")
		default:
			respondUnauthorized(c, "Token validation failed")
		}
	} else {
		respondUnauthorized(c, "Invalid token")
	}
}

type CustomClaims struct {
	Uid int64 `json:"uid"`
	jwt.RegisteredClaims
}

func GenerateToken(uid int64) (string, error) {
	// 创建Claims

	claims := CustomClaims{
		Uid: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), //有效时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                     //签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                     //生效时间
			Issuer:    "root",                                             //签发人
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(conf.ConfGlobal.Server.JwtToken)
	if err != nil {
		return "", err
	}

	return ss, nil
}
