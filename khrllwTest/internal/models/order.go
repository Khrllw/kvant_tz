package models

import "time"

// -------------------------- ORDER --------------------------
// Определение структур данных заказа и их отношений к БД

// ------------------------------------------------------------
// Структуры заказов
// ------------------------------------------------------------

// Order
// Данные о заказе
// Fixme возможно можно убрать поле column - gorm автоматически именует поля
type Order struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Product   string    `gorm:"size:255;not null" json:"product"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	Price     float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// CreateOrderRequest (DTO)
// Структура данных для создания заказа
type CreateOrderRequest struct {
	//UserID   uint    `json:"user_id"  binding:"required"`
	Product  string  `json:"product"  binding:"required,max=255"`
	Quantity int     `json:"quantity" binding:"required,gte=1"`
	Price    float64 `json:"price"    binding:"required,gte=0"`
}

// OrderResponse (DTO)
// Структура данных для ответа
type OrderResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Product   string    `json:"product"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}
