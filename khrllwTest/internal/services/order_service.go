package service

import (
	"errors"
	"khrllwTest/internal/models"
	"khrllwTest/internal/repository"
)

// ------------------------------------------------------------
// Структуры
// ------------------------------------------------------------

// OrderService реализует бизнес-логику работы с заказами
type OrderService struct {
	orderRepo repository.OrderRepository
	userRepo  repository.UserRepository
}

// ------------------------------------------------------------
// Конструктор
// ------------------------------------------------------------

// NewOrderService создает новый экземпляр OrderService
func NewOrderService(
	orderRepo repository.OrderRepository,
	userRepo repository.UserRepository,
) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		userRepo:  userRepo,
	}
}

// ------------------------------------------------------------
// Основные методы
// ------------------------------------------------------------

// CreateOrder создает новый заказ для пользователя
func (s *OrderService) CreateOrder(userID uint, req *models.CreateOrderRequest) (*models.Order, error) {
	if err := s.validateUserExists(userID); err != nil {
		return nil, err
	}
	if err := s.validateOrderRequest(req); err != nil {
		return nil, err
	}
	order := &models.Order{
		UserID:   userID,
		Product:  req.Product,
		Quantity: req.Quantity,
		Price:    req.Price,
	}
	if err := s.orderRepo.Create(order); err != nil {
		return nil, models.ErrDatabaseError
	}
	return order, nil
}

// GetUserOrders возвращает заказы пользователя
func (s *OrderService) GetUserOrders(userID uint) ([]models.Order, error) {
	if err := s.validateUserExists(userID); err != nil {
		return nil, err
	}

	orders, err := s.orderRepo.FindByUserID(userID)
	if err != nil {
		return nil, models.ErrDatabaseError
	}

	return orders, nil
}

// ------------------------------------------------------------
// Вспомогательные методы
// ------------------------------------------------------------

// validateUserExists проверяет существование пользователя
func (s *OrderService) validateUserExists(userID uint) error {
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return err
		}
		return models.ErrDatabaseError
	}
	return nil
}

// validateOrderRequest проверяет валидность данных заказа
func (s *OrderService) validateOrderRequest(req *models.CreateOrderRequest) error {
	if req.Product == "" {
		return models.ErrProductRequired
	}
	if req.Quantity <= 0 {
		return models.ErrInvalidQuantity
	}
	if req.Price <= 0 {
		return models.ErrInvalidPrice
	}
	return nil
}
