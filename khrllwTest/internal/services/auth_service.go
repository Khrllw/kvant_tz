package service

import (
	"errors"
	"fmt"

	"khrllwTest/internal/models"
	"khrllwTest/internal/repository"
	"khrllwTest/internal/utils"
)

// ------------------------------------------------------------
// Реализация
// ------------------------------------------------------------

// AuthService отвечает за бизнес-логику аутентификации
type AuthService struct {
	userRepo     repository.UserRepository
	tokenManager utils.TokenManager
	passHasher   utils.PasswordHasher
}

// ------------------------------------------------------------
// Конструктор
// ------------------------------------------------------------

// NewAuthService создает новый экземпляр AuthService
func NewAuthService(
	userRepo repository.UserRepository,
	tokenManager utils.TokenManager,
	passHasher utils.PasswordHasher,
) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		tokenManager: tokenManager,
		passHasher:   passHasher,
	}
}

// ------------------------------------------------------------
// Методы реализации
// ------------------------------------------------------------

// Login выполняет аутентификацию пользователя и возвращает JWT токен
func (s *AuthService) Login(email, password string) (string, error) {
	if email == "" || password == "" {
		return "", models.ErrEmailPasswordRequired
	}

	user, err := s.authenticateUser(email, password)
	if err != nil {
		return "", fmt.Errorf("authentication failed: %w", err)
	}

	token, err := s.tokenManager.Generate(user.ID)
	if err != nil {
		return "", fmt.Errorf("token generation failed: %w", err)
	}

	return token, nil
}

// authenticateUser проверяет учетные данные пользователя
func (s *AuthService) authenticateUser(email, password string) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			return nil, models.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	if !s.passHasher.Check(password, user.PasswordHash) {
		return nil, models.ErrInvalidCredentials
	}

	return user, nil
}
