package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// JWTSecretKey
// Возвращает секретный ключ из переменных окружения
func JWTSecretKey() []byte {
	key := os.Getenv("JWT_KEY")
	if key == "" {
		panic("JWT_KEY переменная окружения не установлена")
	}
	return []byte(key)
}

// JWTExpiration возвращает время жизни токена из переменных окружения
func JWTExpiration() time.Duration {
	exp := os.Getenv("JWT_EXPIRATION")
	if exp == "" {
		exp = "24h" // значение по умолчанию
	}

	duration, err := time.ParseDuration(exp)
	if err != nil {
		panic("Неверный формат JWT_EXPIRATION. Пример: 24h, 60m, 3600s")
	}
	return duration
}

// HashPassword
// Создает bcrypt хеш из пароля
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash
// Сравнивает пароль с его хешем
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJWTToken
// Создает JWT токен для пользователя
func GenerateJWTToken(userID uint) (string, error) {
	expirationTime := time.Now().Add(JWTExpiration())

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecretKey())
}

// ParseJWTToken
// Проверяет и парсит JWT токен
func ParseJWTToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return JWTSecretKey(), nil
	})
}

// GetUserIDFromToken
// Извлекает user_id из JWT токена
func GetUserIDFromToken(token *jwt.Token) (uint, error) {
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
