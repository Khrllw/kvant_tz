package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"khrllwTest/internal/models"
	"khrllwTest/internal/services"
)

// ------------------------------------------------------------
// Структуры
// ------------------------------------------------------------

// UserHandler обрабатывает HTTP-запросы для работы с пользователями
type UserHandler struct {
	userService *service.UserService
}

// ------------------------------------------------------------
// Конструктор
// ------------------------------------------------------------

// NewUserHandler создает новый экземпляр UserHandler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// ------------------------------------------------------------
// Методы обработки запросов
// ------------------------------------------------------------

// CreateUser обрабатывает запрос на создание пользователя
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendErrorResponse(c, http.StatusBadRequest, models.ErrInvalidRequestFormat.Error())
		return
	}

	user, err := h.userService.CreateUser(&req)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	h.sendUserResponse(c, http.StatusCreated, user)
}

// GetUsers обрабатывает запрос на получение списка пользователей
func (h *UserHandler) GetUsers(c *gin.Context) {
	page, limit, minAge, maxAge, err := h.parseQueryParams(c)
	if err != nil {
		h.sendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	users, total, err := h.userService.GetUsers(page, limit, minAge, maxAge)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	response := models.UsersListResponse{
		Page:  page,
		Limit: limit,
		Total: total,
		Users: h.mapToResponse(users),
	}

	c.JSON(http.StatusOK, response)
}

// GetUserByID обрабатывает запрос на получение пользователя по ID
func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID, err := h.parseUserID(c)
	if err != nil {
		h.sendErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	h.sendUserResponse(c, http.StatusOK, user)
}

// UpdateUser обрабатывает запрос на обновление пользователя
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID, err := h.parseUserID(c)
	if err != nil {
		h.sendErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendErrorResponse(c, http.StatusBadRequest, models.ErrInvalidRequestFormat.Error())
		return
	}

	user, err := h.userService.UpdateUser(userID, &req)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	h.sendUserResponse(c, http.StatusOK, user)
}

// DeleteUser обрабатывает запрос на удаление пользователя
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, err := h.parseUserID(c)
	if err != nil {
		h.sendErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := h.userService.DeleteUser(userID); err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ------------------------------------------------------------
// Вспомогательные методы
// ------------------------------------------------------------

func (h *UserHandler) parseQueryParams(c *gin.Context) (page, limit, minAge, maxAge int, err error) {
	page, err = strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		return 0, 0, 0, 0, models.ErrInvalidPagination
	}

	limit, err = strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		return 0, 0, 0, 0, models.ErrInvalidPagination
	}

	if minAgeStr := c.Query("min_age"); minAgeStr != "" {
		minAge, err = strconv.Atoi(minAgeStr)
		if err != nil {
			return 0, 0, 0, 0, models.ErrInvalidFilterParams
		}
	}

	if maxAgeStr := c.Query("max_age"); maxAgeStr != "" {
		maxAge, err = strconv.Atoi(maxAgeStr)
		if err != nil {
			return 0, 0, 0, 0, models.ErrInvalidFilterParams
		}
	}

	return page, limit, minAge, maxAge, nil
}

func (h *UserHandler) parseUserID(c *gin.Context) (uint, error) {
	id, err := strconv.Atoi(c.Param("user_id"))
	return uint(id), err
}

func (h *UserHandler) handleServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrUserNotFound):
		h.sendErrorResponse(c, http.StatusNotFound, "User not found")
	case errors.Is(err, models.ErrEmailAlreadyExists):
		h.sendErrorResponse(c, http.StatusBadRequest, "Email already exists")
	case errors.Is(err, models.ErrInvalidPagination):
		h.sendErrorResponse(c, http.StatusBadRequest, "Invalid pagination parameters")
	case errors.Is(err, models.ErrInvalidFilterParams):
		h.sendErrorResponse(c, http.StatusBadRequest, "Invalid filter parameters")
	default:
		h.sendErrorResponse(c, http.StatusInternalServerError, "Internal server error")
	}
}

func (h *UserHandler) mapToResponse(users []models.User) []models.UserResponse {
	response := make([]models.UserResponse, len(users))

	for i, user := range users {
		// Direct assignment avoids extra allocation from composite literal
		response[i] = models.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Age:   user.Age,
		}
	}

	return response
}

func (h *UserHandler) sendUserResponse(c *gin.Context, status int, user *models.User) {
	c.JSON(status, models.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	})
}

func (h *UserHandler) sendErrorResponse(c *gin.Context, status int, message string) {
	c.JSON(status, models.ErrorResponse{
		Error: message,
	})
}
