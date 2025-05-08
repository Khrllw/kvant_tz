package models

import "time"

// -------------------------- ORDER --------------------------
// Определение структур данных заказа и их отношений к БД

// ------------------------------------------------------------
// Структуры заказов
// ------------------------------------------------------------

// Order
// Данные о заказе
type Order struct {
	// Уникальный идентификатор заказа
	ID uint `gorm:"primaryKey" json:"id"`

	// Идентификатор пользователя, который сделал заказ
	UserID uint `gorm:"not null" json:"user_id"`

	// Название продукта, заказанного пользователем
	Product string `gorm:"size:255;not null" json:"product"`

	// Количество заказанных единиц товара
	Quantity int `gorm:"not null" json:"quantity"`

	// Цена товара
	Price float64 `gorm:"type:decimal(10,2);not null" json:"price"`

	// Дата и время создания заказа
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// ------------------------------------------------------------
// Request/Response
// ------------------------------------------------------------

// CreateOrderRequest (DTO)
// Структура данных для создания заказа
// @Description Структура для запроса на создание нового заказа
// @Schema example: {"product": "Laptop", "quantity": 2, "price": 1500.50}
type CreateOrderRequest struct {
	// Название продукта, заказанного пользователем
	Product string `json:"product"  binding:"required,max=255" example:"Laptop"`

	// Количество заказанных единиц товара
	Quantity int `json:"quantity" binding:"required,gte=1" example:"2"`

	// Цена товара
	Price float64 `json:"price"    binding:"required,gte=0" example:"1500.50"`
}

// OrderResponse (DTO)
// Структура данных для ответа
// @Description Структура для ответа, содержащая информацию о заказе
// @Schema example: {"id": 1, "user_id": 123, "product": "Laptop", "quantity": 2, "price": 1500.50, "created_at": "2025-05-07T12:34:56Z"}
type OrderResponse struct {
	// Уникальный идентификатор заказа
	ID uint `json:"id" example:"1"`

	// Идентификатор пользователя, который сделал заказ
	UserID uint `json:"user_id" example:"123"`

	// Название продукта, заказанного пользователем
	Product string `json:"product" example:"Laptop"`

	// Количество заказанных единиц товара
	Quantity int `json:"quantity" example:"2"`

	// Цена товара
	Price float64 `json:"price" example:"1500.50"`

	// Дата и время создания заказа
	CreatedAt time.Time `json:"created_at" example:"2025-05-07T12:34:56Z"`
}
