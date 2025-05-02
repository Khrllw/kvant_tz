package handlers

import (
	"khrllwTest/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"khrllwTest/internal/models"
	"khrllwTest/internal/repository"
)

// UserHandler
// Структура для работы с запросами пользователя
type UserHandler struct{ userRepo repository.UserRepository }

// NewUserHandler
// Новая структура для работы с запросами пользователя
func NewUserHandler(userRepo repository.UserRepository) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

// ----------------------------- HANDLERS -----------------------------

// CreateUser Fixme дописать или полностью убрать хеширование
// CreateUser создание нового пользователя

func (handler *UserHandler) CreateUser(context *gin.Context) {
	var req models.CreateUserRequest

	// Парсинг входящего JSON
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка уникальности email
	if _, err := handler.userRepo.FindByEmail(req.Email); err == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "email уже существует"})
		return
	}

	// 3. Хеширование пароля
	//hashedPassword, err := utils.HashPassword(req.Password)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при хешировании пароля"})
	//	return
	//}

	// Создание структуры пользователя
	user := models.User{
		Name:         req.Name,
		Email:        req.Email,
		Age:          req.Age,
		PasswordHash: req.Password,
	}

	// Сохранение данных пользователя в БД
	if err := handler.userRepo.Create(&user); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при создании пользователя"})
		return
	}

	// Возвращение ответа без пароля
	context.JSON(http.StatusCreated, models.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	})
}

// GetUsers
// Возвращает список пользователей с пагинацией
func (handler *UserHandler) GetUsers(context *gin.Context) {
	// Получение параметров запроса
	page, _ := strconv.Atoi(context.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(context.DefaultQuery("limit", "10"))
	minAge, _ := strconv.Atoi(context.Query("min_age"))
	maxAge, _ := strconv.Atoi(context.Query("max_age"))

	// Вычисление offset для пагинации
	offset := (page - 1) * limit

	// Получение пользователей из репозитория
	users, total, err := handler.userRepo.GetAll(offset, limit, minAge, maxAge)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при получении пользователей"})
		return
	}

	// Формирование ответа
	var response models.UsersListResponse
	for _, user := range users {
		response.Users = append(response.Users, models.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Age:   user.Age,
		})
	}

	response.Page = page
	response.Limit = limit
	response.Total = total

	context.JSON(http.StatusOK, response)
}

func (handler *UserHandler) AddUser(context *gin.Context) {
	var req models.CreateUserRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка уникальности email
	if _, err := handler.userRepo.FindByEmail(req.Email); err == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	user := models.User{
		Name:         req.Name,
		Email:        req.Email,
		Age:          req.Age,
		PasswordHash: hashedPassword,
	}

	if err := handler.userRepo.Create(&user); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	context.JSON(http.StatusCreated, models.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	})
}

func (handler *UserHandler) GetUserFromID(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("user_id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := handler.userRepo.FindByID(uint(id))
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	context.JSON(http.StatusOK, models.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	})
}

func (handler *UserHandler) UpdateUser(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("user_id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req models.UpdateUserRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := handler.userRepo.FindByID(uint(id))
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Обновляем только разрешенные поля
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Age > 0 {
		user.Age = req.Age
	}

	if err := handler.userRepo.Update(user); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	context.JSON(http.StatusOK, models.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	})
}

func (handler *UserHandler) DeleteUser(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("user_id"))
	// Обработка ошибок входных данных
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	if err := handler.userRepo.Delete(uint(id)); err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	context.Status(http.StatusNoContent)
}
