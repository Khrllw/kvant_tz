package service

import (
	"errors"
	"fmt"

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
		return nil, fmt.Errorf("validation error: %w", err)
	}

	hashedPassword, err := s.passHasher.Hash(req.Password)
	if err != nil {
		return nil, fmt.Errorf("password hashing failed: %w", err)
	}

	user := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		Age:          req.Age,
		PasswordHash: hashedPassword,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}

	return user, nil
}

// GetUsers возвращает список пользователей с пагинацией
func (s *UserService) GetUsers(page, limit, minAge, maxAge int) ([]models.User, int64, error) {
	if page < 1 || limit < 1 {
		return nil, 0, models.ErrInvalidPagination
	}

	users, total, err := s.userRepo.GetAll(
		(page-1)*limit,
		limit,
		minAge,
		maxAge,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("repository error: %w", err)
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
		return nil, fmt.Errorf("repository error: %w", err)
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
		return nil, fmt.Errorf("repository error: %w", err)
	}

	s.updateUserFields(user, req)

	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}

	return user, nil
}

// DeleteUser удаляет пользователя
func (s *UserService) DeleteUser(userID uint) error {
	if err := s.userRepo.Delete(userID); err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			return models.ErrUserNotFound
		}
		return fmt.Errorf("repository error: %w", err)
	}
	return nil
}

// ------------------------------------------------------------
// Вспомогательные методы
// ------------------------------------------------------------

func (s *UserService) validateCreateRequest(req *models.CreateUserRequest) error {
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return models.ErrInvalidRequestFormat
	}

	if _, err := s.userRepo.FindByEmail(req.Email); err == nil {
		return models.ErrEmailAlreadyExists
	}

	return nil
}

func (s *UserService) updateUserFields(user *models.User, req *models.UpdateUserRequest) {
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Age > 0 {
		user.Age = req.Age
	}
}

func (s *UserService) mapToResponse(users []*models.User) []models.UserResponse {
	response := make([]models.UserResponse, 0, len(users))
	for _, user := range users {
		response = append(response, models.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Age:   user.Age,
		})
	}
	return response
}
