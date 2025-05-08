package repository

import (
	"errors"
	"gorm.io/gorm"
	"khrllwTest/internal/models"
)

// ------------------------------------------------------------
// Интерфейсы
// ------------------------------------------------------------

// OrderRepository определяет контракт для работы с заказами
type OrderRepository interface {

	// Create
	// Создание нового заказа
	Create(order *models.Order) error

	// FindByUserID
	// Поиск всех заказов пользователя по ID
	FindByUserID(userID uint) ([]models.Order, error)

	/* Update
	// Обновление данных заказа
	Update(order *models.Order) error

	// Delete
	// Удаление заказа по ID
	Delete(id uint) error

	*/
}

// ------------------------------------------------------------
// Конструктор
// ------------------------------------------------------------

// NewOrderRepository создает новый экземпляр OrderRepository
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &OrderRepositoryImpl{db: db}
}

// ------------------------------------------------------------
// Реализация
// ------------------------------------------------------------

// OrderRepositoryImpl - реализация для GORM
type OrderRepositoryImpl struct {
	db *gorm.DB // Экземпляр подключения к БД
}

// ------------------------------------------------------------
// Методы OrderRepositoryImpl
// ------------------------------------------------------------

func (r *OrderRepositoryImpl) Create(order *models.Order) error {
	// INSERT INTO orders (...) VALUES (...)
	return r.db.Create(order).Error
}

func (r *OrderRepositoryImpl) FindByUserID(userID uint) ([]models.Order, error) {
	var orders []models.Order
	// SELECT * FROM orders WHERE user_id = ?
	err := r.db.Where("user_id = ?", userID).Find(&orders).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return orders, nil
	}
	return orders, err
}

/*

func (r *OrderRepositoryImpl) Update(order *models.Order) error {
	// UPDATE orders SET ... WHERE id = ?
	return r.db.Save(order).Error
}

func (r *OrderRepositoryImpl) Delete(id uint) error {
	// DELETE FROM orders WHERE id = ?
	return r.db.Delete(&models.Order{}, id).Error
}

*/
