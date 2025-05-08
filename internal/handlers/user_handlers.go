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
// @Tags Users
// @Summary Создать нового пользователя
// @Description Создает нового пользователя в системе
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "Данные пользователя"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} models.ErrorLoginResponse "Неверный формат запроса/некорректные данные"
// @Failure 500 {object} models.ErrorLoginResponse "Внутренняя ошибка сервера"
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendErrorResponse(c, http.StatusBadRequest, models.ErrInvalidRequestFormat)
		return
	}

	user, err := h.userService.CreateUser(&req)
	if err != nil {
		if errors.Is(err, models.ErrPasswordHashFailed) || errors.Is(err, models.ErrDatabaseError) {
			h.sendErrorResponse(c, http.StatusInternalServerError, err)
			return
		}
		h.sendErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	h.sendUserResponse(c, http.StatusCreated, user)
}

// GetUsers обрабатывает запрос на получение списка пользователей
// @Tags Users
// @Summary Получить список пользователей
// @Description Возвращает список пользователей с пагинацией и фильтрацией по возрасту
// @Accept json
// @Produce json
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(10)
// @Param min_age query int false "Min Age"
// @Param max_age query int false "Max Age"
// @Success 200 {object} models.UsersListResponse
// @Failure 400 {object} models.ErrorLoginResponse "Неверный формат запроса/некорректные данные"
// @Failure 500 {object} models.ErrorLoginResponse "Внутренняя ошибка сервера"
// @Router /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	page, limit, minAge, maxAge, err := h.parseQueryParams(c)
	if err != nil {
		h.sendErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	users, total, err := h.userService.GetUsers(page, limit, minAge, maxAge)
	if err != nil {
		h.sendErrorResponse(c, http.StatusInternalServerError, err)
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
// @Tags Users
// @Summary Получить пользователя по ID
// @Description Возвращает данные пользователя по его ID
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.ErrorLoginResponse "Неверный формат запроса/некорректные данные"
// @Failure 404 {object} models.ErrorLoginResponse "Пользователь не найден"
// @Failure 500 {object} models.ErrorLoginResponse "Внутренняя ошибка сервера"
// @Router /users/{user_id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID, err := h.parseUserID(c)
	if err != nil {
		h.sendErrorResponse(c, http.StatusBadRequest, models.ErrInvalidUserID)
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			h.sendErrorResponse(c, http.StatusNotFound, err)
			return
		}
		h.sendErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	h.sendUserResponse(c, http.StatusOK, user)
}

// UpdateUser обрабатывает запрос на обновление пользователя
// @Tags Users
// @Summary Обновить данные пользователя
// @Description Обновляет информацию о пользователе
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param user body models.UpdateUserRequest true "User Data"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.ErrorLoginResponse "Неверный формат запроса/некорректные данные"
// @Failure 404 {object} models.ErrorLoginResponse "Пользователь не найден"
// @Failure 500 {object} models.ErrorLoginResponse "Внутренняя ошибка сервера"
// @Router /users/{user_id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID, err := h.parseUserID(c)
	if err != nil {
		h.sendErrorResponse(c, http.StatusBadRequest, models.ErrInvalidUserID)
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendErrorResponse(c, http.StatusBadRequest, models.ErrInvalidRequestFormat)
		return
	}

	user, err := h.userService.UpdateUser(userID, &req)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			h.sendErrorResponse(c, http.StatusNotFound, err)
		}
		if errors.Is(err, models.ErrDatabaseError) {
			h.sendErrorResponse(c, http.StatusInternalServerError, err)
		}
		h.sendErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	h.sendUserResponse(c, http.StatusOK, user)
}

// DeleteUser обрабатывает запрос на удаление пользователя
// @Tags Users
// @Summary Удалить пользователя
// @Description Удаляет пользователя по его ID
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} models.ErrorLoginResponse "Неверный формат запроса/некорректные данные"
// @Failure 404 {object} models.ErrorLoginResponse "Пользователь не найден"
// @Failure 500 {object} models.ErrorLoginResponse "Внутренняя ошибка сервера"
// @Router /users/{user_id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, err := h.parseUserID(c)
	if err != nil {
		h.sendErrorResponse(c, http.StatusBadRequest, models.ErrInvalidUserID)
		return
	}
	if err := h.userService.DeleteUser(userID); err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			h.sendErrorResponse(c, http.StatusNotFound, err)
		}
		h.sendErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ------------------------------------------------------------
// Вспомогательные методы
// ------------------------------------------------------------

// parseQueryParams парсит параметры запроса
// @Description Парсит параметры запроса для пагинации и фильтрации пользователей
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

// parseUserID парсит ID пользователя
// @Description Парсит user_id из URL-параметра
func (h *UserHandler) parseUserID(c *gin.Context) (uint, error) {
	id, err := strconv.Atoi(c.Param("user_id"))
	return uint(id), err
}

// mapToResponse преобразует список пользователей в формат ответа
func (h *UserHandler) mapToResponse(users []models.User) []models.UserResponse {
	response := make([]models.UserResponse, len(users))

	for i, user := range users {
		response[i] = models.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Age:   user.Age,
		}
	}

	return response
}

// sendUserResponse отправляет успешный ответ с данными пользователя
func (h *UserHandler) sendUserResponse(c *gin.Context, status int, user *models.User) {
	c.JSON(status, models.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	})
}

// sendErrorResponse отправляет ответ с ошибкой
func (h *UserHandler) sendErrorResponse(c *gin.Context, status int, err error) {
	c.JSON(status, models.ErrorLoginResponse{
		Error: err.Error(),
	})
}
