package utils

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthConfig содержит конфигурацию для аутентификации
type AuthConfig struct {
	JWTKey        string
	JWTExpiration time.Duration
}

// NewAuthConfig создает конфигурацию аутентификации из переменных окружения
func NewAuthConfig() (*AuthConfig, error) {
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

	return &AuthConfig{
		JWTKey:        key,
		JWTExpiration: duration,
	}, nil
}

// PasswordHasher предоставляет методы для работы с паролями
type PasswordHasher struct{}

// Hash создает bcrypt хеш из пароля
func (h *PasswordHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Check сравнивает пароль с его хешем
func (h *PasswordHasher) Check(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// JWTManager предоставляет методы для работы с JWT токенами
type JWTManager struct {
	config *AuthConfig
}

// NewJWTManager создает новый JWTManager
func NewJWTManager(config *AuthConfig) *JWTManager {
	return &JWTManager{config: config}
}

// Generate создает JWT токен для пользователя
func (m *JWTManager) Generate(userID uint) (string, error) {
	expirationTime := time.Now().Add(m.config.JWTExpiration)

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.config.JWTKey))
}

// Parse проверяет и парсит JWT токен
func (m *JWTManager) Parse(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		log.Println(tokenString)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(m.config.JWTKey), nil
	})
}

// ExtractUserID извлекает user_id из JWT токена
func (m *JWTManager) ExtractUserID(token *jwt.Token) (uint, error) {
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
