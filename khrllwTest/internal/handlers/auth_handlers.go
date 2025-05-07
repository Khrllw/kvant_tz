package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"khrllwTest/internal/models"
	"khrllwTest/internal/services"
)

// ------------------------------------------------------------
// Структуры
// ------------------------------------------------------------

// AuthHandler обрабатывает HTTP-запросы для аутентификации
type AuthHandler struct {
	authService *service.AuthService
}

// ------------------------------------------------------------
// Конструктор
// ------------------------------------------------------------

// NewAuthHandler создает новый экземпляр AuthHandler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// ------------------------------------------------------------
// Методы реализации
// ------------------------------------------------------------

// Login обрабатывает запрос на аутентификацию
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendErrorResponse(c, http.StatusBadRequest, models.ErrInvalidRequestFormat.Error())
		return
	}

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		h.handleAuthError(c, err)
		return
	}

	h.sendSuccessResponse(c, token)
}

// handleAuthError обрабатывает ошибки аутентификации
func (h *AuthHandler) handleAuthError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrInvalidCredentials):
		h.sendErrorResponse(c, http.StatusUnauthorized, "Неверные учетные данные")
	case errors.Is(err, models.ErrEmailPasswordRequired):
		h.sendErrorResponse(c, http.StatusBadRequest, err.Error())
	default:
		h.sendErrorResponse(c, http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}
}

// sendErrorResponse отправляет ответ с ошибкой
func (h *AuthHandler) sendErrorResponse(c *gin.Context, statusCode int, errorMsg string) {
	c.JSON(statusCode, models.ErrorResponse{
		Error: errorMsg,
	})
}

// sendSuccessResponse отправляет успешный ответ
func (h *AuthHandler) sendSuccessResponse(c *gin.Context, token string) {
	c.JSON(http.StatusOK, models.LoginResponse{
		Token: token,
	})
}
