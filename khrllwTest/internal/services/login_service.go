package service

import (
	"errors"
	"khrllwTest/internal/models"
	"khrllwTest/internal/repository"
	"khrllwTest/internal/utils"
)

// ------------------------------------------------------------
// Реализация
// ------------------------------------------------------------

// LoginService отвечает за бизнес-логику авторизации
type LoginService struct {
	userRepo     repository.UserRepository
	tokenManager utils.TokenManager
	passHasher   utils.PasswordHasher
}

// ------------------------------------------------------------
// Конструктор
// ------------------------------------------------------------

// NewLoginService создает новый экземпляр LoginService
func NewLoginService(
	userRepo repository.UserRepository,
	tokenManager utils.TokenManager,
	passHasher utils.PasswordHasher,
) *LoginService {

	return &LoginService{
		userRepo:     userRepo,
		tokenManager: tokenManager,
		passHasher:   passHasher,
	}
}

// ------------------------------------------------------------
// Методы реализации
// ------------------------------------------------------------

// Login выполняет авторизацию пользователя и возвращает JWT токен
func (s *LoginService) Login(email, password string) (string, error) {
	if email == "" || password == "" {
		return "", models.ErrEmailPasswordRequired
	}

	user, err := s.authenticateUser(email, password)
	if err != nil {
		return "", err
	}

	token, err := s.tokenManager.Generate(user.ID)
	if err != nil {
		return "", models.ErrTokenGenerationFailed
	}

	return token, nil
}

// authenticateUser проверяет учетные данные пользователя
func (s *LoginService) authenticateUser(email, password string) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			return nil, models.ErrInvalidCredentials
		}
		return nil, models.ErrDatabaseError
	}

	if !s.passHasher.Check(password, user.PasswordHash) {
		return nil, models.ErrInvalidCredentials
	}

	return user, nil
}
