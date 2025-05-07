package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"os"
	"time"
)

// ------------------------------------------------------------
// Конфигурация
// ------------------------------------------------------------

// JWTConfig содержит конфигурацию для JWT аутентификации
type JWTConfig struct {
	JWTKey        string
	JWTExpiration time.Duration
}

// NewJWTConfig создает конфигурацию аутентификации из переменных окружения
func NewJWTConfig() (*JWTConfig, error) {
	key := os.Getenv("JWT_KEY")
	if key == "" {
		return nil, errors.New("JWT_KEY переменная окружения не установлена")
	}

	exp := os.Getenv("JWT_EXPIRATION")
	if exp == "" {
		exp = "24h" // значение по умолчанию
	}

	duration, err := time.ParseDuration(exp)
	if err != nil {
		return nil, errors.New("неверный формат JWT_EXPIRATION. Пример: 24h, 60m, 3600s")
	}

	return &JWTConfig{
		JWTKey:        key,
		JWTExpiration: duration,
	}, nil
}

// ------------------------------------------------------------
// Интерфейс
// ------------------------------------------------------------

// TokenManager определяет контракт для работы с JWT токенами
type TokenManager interface {
	Generate(userID uint) (string, error)
	Parse(tokenString string) (*jwt.Token, error)
	ExtractUserID(token *jwt.Token) (uint, error)
}

// ------------------------------------------------------------
// Реализация
// ------------------------------------------------------------

// JWTManager реализует TokenManager с использованием библиотеки jwt
type jwtManager struct {
	config *JWTConfig
}

// ------------------------------------------------------------
// Конструктор
// ------------------------------------------------------------

// NewTokenManager создает новый TokenManager с JWT
func NewTokenManager(config *JWTConfig) TokenManager {
	return &jwtManager{config: config}
}

// ------------------------------------------------------------
// Методы реализации
// ------------------------------------------------------------

// Generate создает JWT токен для пользователя
func (m *jwtManager) Generate(userID uint) (string, error) {
	expirationTime := time.Now().Add(m.config.JWTExpiration)

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.config.JWTKey))
}

// Parse проверяет и парсит JWT токен
func (m *jwtManager) Parse(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		log.Println(tokenString)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(m.config.JWTKey), nil
	})
}

// ExtractUserID извлекает user_id из JWT токена
func (m *jwtManager) ExtractUserID(token *jwt.Token) (uint, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, jwt.ErrTokenInvalidClaims
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, jwt.ErrTokenInvalidId
	}

	return uint(userID), nil
}
