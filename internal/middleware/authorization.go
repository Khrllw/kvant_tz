package middleware

import (
	"errors"
	"khrllwTest/internal/repository"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"khrllwTest/internal/models"
	"khrllwTest/internal/utils"
)

// ------------------------------------------------------------
// Структуры
// ------------------------------------------------------------

// Authorization Middleware для авторизации с использованием JWT
type Authorization struct {
	tokenManager utils.TokenManager
	userRepos    repository.UserRepository
}

// ------------------------------------------------------------
// Конструктор
// ------------------------------------------------------------

// NewAuthorization создает новый экземпляр Authorization
func NewAuthorization(tokenManager utils.TokenManager, userRepos repository.UserRepository) *Authorization {
	return &Authorization{
		tokenManager: tokenManager,
		userRepos:    userRepos,
	}
}

// ------------------------------------------------------------
// Основные методы
// ------------------------------------------------------------

// Middleware проверяет JWT токен и добавляет user_id в контекст и сверяет его с параметром маршрута
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

		// Извлекаем параметр id из маршрута и проверяем его совпадение с user_id
		paramID := c.Param("user_id")
		if paramID != "" && paramID != strconv.Itoa(int(userID)) {
			m.abortWithError(c, http.StatusUnauthorized, models.ErrInvalidTokenClaims)
			return
		}

		_, err = m.userRepos.FindByID(userID)
		if err != nil {
			if errors.Is(err, models.ErrDatabaseError) {
				m.abortWithError(c, http.StatusInternalServerError, err)
			}
			m.abortWithError(c, http.StatusUnauthorized, err)
		}

		// Добавляем user_id в контекст
		c.Set("user_id", userID)
	}
}

// ------------------------------------------------------------
// Вспомогательные методы
// ------------------------------------------------------------

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
	c.AbortWithStatusJSON(status, models.ErrorLoginResponse{
		Error: err.Error(),
	})
}
