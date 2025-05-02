package repository

import (
	"gorm.io/gorm"
	"khrllwTest/internal/models"
)

// OrderRepositoryImpl - реализация для GORM
type OrderRepositoryImpl struct {
	db *gorm.DB // Экземпляр подключения к БД
}

// OrderRepository - определяет контракт для работы с заказами
type OrderRepository interface {

	// Create
	//Создание нового заказа
	Create(order *models.Order) error

	// FindByID
	//Поиск заказа по ID
	FindByID(id uint) (*models.Order, error)

	// FindByUserID
	// Поиск всех заказов пользователя
	FindByUserID(userID uint) ([]models.Order, error)

	// Update
	//Обновление данных заказа
	Update(order *models.Order) error

	// Delete
	//Удаление заказа по ID
	Delete(id uint) error
}

// NewOrderRepository создает новый экземпляр OrderRepository
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &OrderRepositoryImpl{db: db}
}

func (r *OrderRepositoryImpl) Create(order *models.Order) error {
	// INSERT INTO orders (...) VALUES (...)
	return r.db.Create(order).Error
}

func (r *OrderRepositoryImpl) FindByID(id uint) (*models.Order, error) {
	var order models.Order
	err := r.db.First(&order, id).Error
	return &order, err
}

func (r *OrderRepositoryImpl) FindByUserID(userID uint) ([]models.Order, error) {
	var orders []models.Order
	// SELECT * FROM orders WHERE user_id = ?
	err := r.db.Where("user_id = ?", userID).Find(&orders).Error
	return orders, err
}

func (r *OrderRepositoryImpl) Update(order *models.Order) error {
	// UPDATE orders SET ... WHERE id = ?
	return r.db.Save(order).Error
}

func (r *OrderRepositoryImpl) Delete(id uint) error {
	// DELETE FROM orders WHERE id = ?
	return r.db.Delete(&models.Order{}, id).Error
}
