package models

// ------------------------ USER ------------------------
// Определение структур данных пользователя и их отношений к БД

// ------------------------------------------------------------
// Структуры пользователя
// ------------------------------------------------------------

// Fixme надо добавить examples

// User
// Данные о пользователе
// Fixme возможно можно убрать поле column - gorm автоматически именует поля
type User struct {
	ID           uint    `gorm:"primaryKey" json:"id"`
	Name         string  `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Email        string  `gorm:"column:email;type:varchar(255);uniqueIndex;not null" json:"email"`
	Age          int     `gorm:"column:age;not null" json:"age"`
	PasswordHash string  `gorm:"column:password_hash;type:varchar(255);not null" json:"-"`
	Orders       []Order `gorm:"foreignKey:UserID" json:"-"`
}

// CreateUserRequest (DTO)
// Структура данных для создания пользователя
// Fixme надо подумать над полем password
type CreateUserRequest struct {
	Name     string `json:"name"     binding:"required,max=255"`
	Email    string `json:"email"    binding:"required,email,max=255"`
	Age      int    `json:"age"      binding:"required,gte=0"`
	Password string `json:"password" binding:"required,min=8"`
}

type UpdateUserRequest struct {
	Name  string `json:"name"     binding:"required,max=255"`
	Email string `json:"email"    binding:"required,email,max=255"`
	Age   int    `json:"age"      binding:"required,gte=0"`
}

// UserResponse (DTO)
// Структура данных для ответа
type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// UsersListResponse
// представляет ответ со списком пользователей и метаданными пагинации
type UsersListResponse struct {
	Page  int            `json:"page" example:"1"`
	Limit int            `json:"limit" example:"10"`
	Total int64          `json:"total" example:"100"`
	Users []UserResponse `json:"users"`
}
