package models

// LoginRequest представляет структуру данных для запроса аутентификации
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"securepassword123"`
}

// LoginResponse представляет структуру ответа с токеном
type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
