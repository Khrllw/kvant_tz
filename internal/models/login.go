package models

// -------------------------- LOGIN --------------------------
// Определение структур данных для аутентификации

// LoginRequest представляет структуру данных для запроса аутентификации
// @Description Структура данных для аутентификации пользователя через email и пароль
// @Schema example: {"email": "user@example.com", "password": "securepassword123"}
type LoginRequest struct {
	// Email пользователя для аутентификации
	Email string `json:"email" binding:"required,email" example:"user@example.com"`

	// Пароль пользователя
	Password string `json:"password" binding:"required,min=8" example:"securepassword123"`
}

// LoginResponse представляет структуру ответа с токеном
// @Description Структура, которая возвращает токен для аутентифицированного пользователя
// @Schema example: {"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}
type LoginResponse struct {
	// Токен аутентифицированного пользователя
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// ErrorLoginResponse представляет структуру для возвращаемых ошибок
// @Description Структура, которая содержит сообщение об ошибке
// @Schema example: {"error": "Invalid credentials"}
type ErrorLoginResponse struct {
	// Сообщение об ошибке
	Error string `json:"error"`
}
