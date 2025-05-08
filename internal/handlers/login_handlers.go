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

// LoginHandler обрабатывает HTTP-запросы для входа пользователя
type LoginHandler struct {
	loginService *service.LoginService
}

// ------------------------------------------------------------
// Конструктор
// ------------------------------------------------------------

// NewLoginHandler создает новый экземпляр LoginHandler
func NewLoginHandler(authService *service.LoginService) *LoginHandler {
	return &LoginHandler{
		loginService: authService,
	}
}

// ------------------------------------------------------------
// Основные методы
// ------------------------------------------------------------

// Login godoc
// @Tags Authorization
// @Summary Авторизация пользователя
// @Description Вход в систему с email и паролем
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Данные для входа"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.ErrorLoginResponse "Неверный формат запроса"
// @Failure 401 {object} models.ErrorLoginResponse "Некорректные данные"
// @Failure 500 {object} models.ErrorLoginResponse "Внутрення ошибка сервера"
// @Router /auth/login [post]
func (h *LoginHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendErrorResponse(c, http.StatusBadRequest, models.ErrInvalidRequestFormat)
		return
	}

	token, err := h.loginService.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, models.ErrDatabaseError) || errors.Is(err, models.ErrTokenGenerationFailed) {
			h.sendErrorResponse(c, http.StatusInternalServerError, models.ErrInternalServerError)
			return
		}
		h.sendErrorResponse(c, http.StatusUnauthorized, err)
		return
	}

	h.sendSuccessResponse(c, token)
}

// Ответ с ошибкой авторизации
func (h *LoginHandler) sendErrorResponse(c *gin.Context, statusCode int, err error) {
	c.JSON(statusCode, models.ErrorLoginResponse{
		Error: err.Error(),
	})
}

// Ответ успешным входом в систему
func (h *LoginHandler) sendSuccessResponse(c *gin.Context, token string) {
	c.JSON(http.StatusOK, models.LoginResponse{
		Token: token,
	})
}
