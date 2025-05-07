package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"khrllwTest/internal/models"
	"khrllwTest/internal/utils"
)

// ------------------------------------------------------------
// Структуры
// ------------------------------------------------------------

type Authorization struct {
	tokenManager utils.TokenManager
}

// ------------------------------------------------------------
// Конструктор
// ------------------------------------------------------------

func NewAuthorization(tokenManager utils.TokenManager) *Authorization {
	return &Authorization{tokenManager: tokenManager}
}

// ------------------------------------------------------------
// Основные методы
// ------------------------------------------------------------

// Middleware проверяет JWT токен и добавляет user_id в контекст
func (m *Authorization) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Извлекаем токен из заголовка
		tokenString := m.extractToken(c)
		if tokenString == "" {
			m.abortWithError(c, http.StatusUnauthorized, models.ErrTokenRequired)
			return
		}

		// Парсим и валидируем токен
		token, err := m.tokenManager.Parse(tokenString)
		log.Println(token, err)
		if err != nil || !token.Valid {
			m.abortWithError(c, http.StatusUnauthorized, models.ErrInvalidToken)
			return
		}

		// Извлекаем user_id из токена
		userID, err := m.tokenManager.ExtractUserID(token)
		if err != nil {
			m.abortWithError(c, http.StatusUnauthorized, models.ErrInvalidTokenClaims)
			return
		}

		// Добавляем user_id в контекст
		c.Set("user_id", userID)
		c.Next()
	}
}

// extractToken извлекает токен из заголовка Authorization
func (m *Authorization) extractToken(c *gin.Context) string {
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

// abortWithError отправляет ошибку и прерывает выполнение
func (m *Authorization) abortWithError(c *gin.Context, status int, err error) {
	c.AbortWithStatusJSON(status, models.ErrorResponse{
		Error: err.Error(),
	})
}
