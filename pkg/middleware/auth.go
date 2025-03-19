package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuthMiddleware is a middleware for JWT authentication
func JWTAuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供授權標頭"})
			c.Abort()
			return
		}

		// Check if the header has the Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "授權格式無效"})
			c.Abort()
			return
		}

		// Extract the token
		tokenString := parts[1]

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("非預期的簽名方法: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無效的令牌: " + err.Error()})
			c.Abort()
			return
		}

		// Check if the token is valid
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Check if the token is expired
			if exp, ok := claims["exp"].(float64); ok {
				if time.Now().Unix() > int64(exp) {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "令牌已過期"})
					c.Abort()
					return
				}
			}

			// Set the user ID in the context
			if userID, ok := claims["user_id"].(string); ok {
				c.Set("user_id", userID)
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "令牌缺少用戶 ID"})
				c.Abort()
				return
			}

			// Set the user role in the context
			if role, ok := claims["role"].(string); ok {
				c.Set("role", role)
			}

			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無效的令牌"})
			c.Abort()
			return
		}
	}
}

// GenerateJWT generates a JWT token
func GenerateJWT(userID, role, secretKey string, expiryHours int) (string, error) {
	// Create the claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * time.Duration(expiryHours)).Unix(),
		"iat":     time.Now().Unix(),
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ExtractUserID extracts the user ID from the JWT token
func ExtractUserID(tokenString, secretKey string) (string, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("非預期的簽名方法: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return "", err
	}

	// Extract the user ID from the claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["user_id"].(string); ok {
			return userID, nil
		}
		return "", errors.New("令牌缺少用戶 ID")
	}

	return "", errors.New("無效的令牌")
}
