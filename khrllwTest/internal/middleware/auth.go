package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"khrllwTest/internal/utils"
)

// AuthMiddleware
// проверяет JWT токен и добавляет user_id в контекст
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Пропускаем публичные маршруты
		if c.Request.URL.Path == "/auth/login" {
			c.Next()
			return
		}

		// Извлекаем токен из заголовка
		tokenString := extractToken(c)
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization token required",
			})
			return
		}

		// Парсим и валидируем токен
		token, err := utils.ParseJWTToken(tokenString)
		if err != nil || !token.Valid {
			log.Printf("Invalid token error: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			return
		}

		// Извлекаем user_id из токена
		userID, err := utils.GetUserIDFromToken(token)
		if err != nil {
			log.Printf("Failed to get user ID from token: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token claims",
			})
			return
		}

		// Добавляем user_id в контекст
		c.Set("user_id", userID)
		c.Next()
	}
}

// extractToken
// извлекает токен из заголовка Authorization
func extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	// Формат: Bearer <token>
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

// AdminMiddleware
// (дополнительно) - проверка ролей пользователя
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Реализация проверки ролей при необходимости
		c.Next()
	}
}
