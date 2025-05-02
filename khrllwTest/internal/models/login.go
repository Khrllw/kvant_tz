package models

// LoginRequest представляет структуру данных для запроса аутентификации
// @Description Запрос на авторизацию пользователя
type LoginRequest struct {
	// Email пользователя (должен быть валидным email)
	Email string `json:"email" binding:"required,email" example:"user@example.com"`

	// Пароль пользователя (минимум 8 символов)
	Password string `json:"password" binding:"required,min=8" example:"securepassword123"`
}

// LoginResponse представляет структуру ответа с JWT токеном
// @Description Ответ с JWT токеном при успешной авторизации
type LoginResponse struct {
	// JWT токен для аутентификации
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}
