package models

// --------------------------- USER ---------------------------
// Определение структур данных пользователя и их отношений к БД

// ------------------------------------------------------------
// Структуры пользователя
// ------------------------------------------------------------

// User
// Структура данных пользователя
type User struct {
	// Уникальный идентификатор пользователя
	ID uint `gorm:"primaryKey" json:"id"`

	// Имя пользователя
	Name string `gorm:"type:varchar(255);not null" json:"name"`

	// Email пользователя
	Email string `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`

	// Возраст пользователя
	Age int `gorm:"not null" json:"age"`

	// Хэш пароля пользователя
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"`

	// Список заказов пользователя
	Orders []Order `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"-"`
}

// ------------------------------------------------------------
// Request/Response
// ------------------------------------------------------------

// CreateUserRequest (DTO)
// Структура данных для создания пользователя
// @Description Структура для запроса на создание нового пользователя
// @Schema example: {"name": "John Doe", "email": "john@example.com", "age": 30, "password": "securepassword123"}
type CreateUserRequest struct {
	// Имя пользователя
	Name string `json:"name"     binding:"required,max=255" example:"John Doe"`

	// Email пользователя
	Email string `json:"email"    binding:"required,email,max=255" example:"john@example.com"`

	// Возраст пользователя
	Age int `json:"age"      binding:"required,gte=0" example:"30"`

	// Введенный пользователем пароль
	Password string `json:"password" binding:"required,min=8" example:"securepassword123"`
}

// UpdateUserRequest
// Структура данных для обновления информации о пользователе
// @Description Структура для запроса на обновление данных пользователя
// @Schema example: {"name": "John Doe", "email": "john@example.com", "age": 30}
type UpdateUserRequest struct {
	// Имя пользователя
	Name string `json:"name"     binding:"required,max=255" example:"John Doe"`

	// Email пользователя
	Email string `json:"email"    binding:"required,email,max=255" example:"john@example.com"`

	// Возраст пользователя
	Age int `json:"age"      binding:"required,gte=0" example:"30"`
}

// UserResponse (DTO)
// Структура данных для ответа на запросы о пользователях
// @Description Структура ответа, содержащая информацию о пользователе
// @Schema example: {"id": 1, "name": "John Doe", "email": "john@example.com", "age": 30}
type UserResponse struct {
	// Уникальный идентификатор пользователя
	ID uint `json:"id" example:"1"`

	// Имя пользователя
	Name string `json:"name" example:"John Doe"`

	// Email пользователя
	Email string `json:"email" example:"john@example.com"`

	// Возраст пользователя
	Age int `json:"age" example:"30"`
}

// UsersListResponse
// Ответ со списком пользователей и метаданными пагинации
// @Description Структура ответа с пользователями и информацией о пагинации
// @Schema example: {"page": 1, "limit": 10, "total": 100, "users": [{"id": 1, "name": "John Doe", "email": "john@example.com", "age": 30}]}
type UsersListResponse struct {
	// Номер текущей страницы в результате пагинации
	Page int `json:"page" example:"1"`

	// Количество элементов (пользователей) на одной странице
	Limit int `json:"limit" example:"10"`

	// Общее количество пользователей, соответствующих запросу (до применения пагинации)
	Total int64 `json:"total" example:"100"`

	// Список пользователей на текущей странице
	Users []UserResponse `json:"users"`
}
