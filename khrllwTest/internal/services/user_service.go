package service

import (
	"errors"
	"khrllwTest/internal/models"
	"khrllwTest/internal/repository"
	"khrllwTest/internal/utils"
)

// ------------------------------------------------------------
// Структуры
// ------------------------------------------------------------

// UserService реализует бизнес-логику работы с пользователями
type UserService struct {
	userRepo   repository.UserRepository
	passHasher utils.PasswordHasher
}

// ------------------------------------------------------------
// Конструктор
// ------------------------------------------------------------

// NewUserService создает новый экземпляр UserService
func NewUserService(
	userRepo repository.UserRepository,
	passHasher utils.PasswordHasher,
) *UserService {
	return &UserService{
		userRepo:   userRepo,
		passHasher: passHasher,
	}
}

// ------------------------------------------------------------
// Основные методы
// ------------------------------------------------------------

// CreateUser создает нового пользователя
func (s *UserService) CreateUser(req *models.CreateUserRequest) (*models.User, error) {
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	hashedPassword, err := s.passHasher.Hash(req.Password)
	if err != nil {
		return nil, models.ErrPasswordHashFailed
	}

	user := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		Age:          req.Age,
		PasswordHash: hashedPassword,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, models.ErrDatabaseError
	}

	return user, nil
}

// GetUsers возвращает список пользователей с пагинацией и фильтрацией
func (s *UserService) GetUsers(page, limit, minAge, maxAge int) ([]models.User, int64, error) {
	users, total, err := s.userRepo.GetAll(
		(page-1)*limit,
		limit,
		minAge,
		maxAge,
	)
	if err != nil {
		return nil, 0, models.ErrDatabaseError
	}

	return users, total, nil
}

// GetUserByID возвращает пользователя по ID
func (s *UserService) GetUserByID(userID uint) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			return nil, models.ErrUserNotFound
		}
		return nil, models.ErrDatabaseError
	}
	return user, nil
}

// UpdateUser обновляет данные пользователя
func (s *UserService) UpdateUser(userID uint, req *models.UpdateUserRequest) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			return nil, models.ErrUserNotFound
		}
		return nil, models.ErrDatabaseError
	}

	if err := s.updateUserFields(user, req); err != nil {
		return nil, err
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, models.ErrDatabaseError
	}

	return user, nil
}

// DeleteUser удаляет пользователя по ID
func (s *UserService) DeleteUser(userID uint) error {
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			return models.ErrUserNotFound
		}
		return models.ErrDatabaseError
	}
	if err := s.userRepo.Delete(userID); err != nil {
		return models.ErrDatabaseError
	}
	return nil
}

// ------------------------------------------------------------
// Вспомогательные методы
// ------------------------------------------------------------

// validateCreateRequest проверяет данные запроса на создание пользователя
func (s *UserService) validateCreateRequest(req *models.CreateUserRequest) error {
	if req.Name == "" {
		return models.ErrInvalidUserName
	}
	if req.Email == "" {
		return models.ErrInvalidUserEmail
	}
	if req.Password == "" {
		return models.ErrInvalidUserPassword
	}
	if req.Age <= 0 || req.Age > 150 {
		return models.ErrInvalidUserAge
	}
	if _, err := s.userRepo.FindByEmail(req.Email); err == nil {
		return models.ErrEmailAlreadyExists
	}
	return nil
}

// updateUserFields обновляет поля пользователя
func (s *UserService) updateUserFields(user *models.User, req *models.UpdateUserRequest) error {
	if req.Name == "" {
		return models.ErrInvalidUserName
	}
	if req.Email == "" {
		return models.ErrInvalidUserEmail
	}
	if req.Age <= 0 || req.Age > 150 {
		return models.ErrInvalidUserAge
	}
	user.Name = req.Name
	user.Email = req.Email
	user.Age = req.Age
	return nil
}
