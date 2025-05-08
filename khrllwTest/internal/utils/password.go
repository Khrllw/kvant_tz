package utils

import "golang.org/x/crypto/bcrypt"

// ------------------------------------------------------------
// Интерфейс
// ------------------------------------------------------------

// PasswordHasher определяет контракт для хеширования паролей
type PasswordHasher interface {

	// Hash хэширует пароль
	Hash(password string) (string, error)

	// Check проверяет, соответствует ли пароль его хешу
	Check(password, hash string) bool
}

// ------------------------------------------------------------
// Реализация
// ------------------------------------------------------------

// bcryptHasher реализует PasswordHasher используя bcrypt
type bcryptHasher struct {
	// cost стоимость bcryptHasher
	cost int
}

// ------------------------------------------------------------
// Конструктор
// ------------------------------------------------------------

// NewPasswordHasher создает новый bcryptHasher
func NewPasswordHasher(cost int) PasswordHasher {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}
	return &bcryptHasher{cost: cost}
}

// ------------------------------------------------------------
// Методы реализации
// ------------------------------------------------------------

// Hash создает bcrypt хеш из пароля
func (h *bcryptHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	return string(bytes), err
}

// Check сравнивает пароль с его хешем
func (h *bcryptHasher) Check(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
