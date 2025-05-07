package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"khrllwTest/internal/models"
	"khrllwTest/internal/services"
)

// ------------------------------------------------------------
// Структуры
// ------------------------------------------------------------

// OrderHandler обрабатывает HTTP-запросы для работы с заказами
type OrderHandler struct {
	orderService *service.OrderService
}

// ------------------------------------------------------------
// Конструктор
// ------------------------------------------------------------

// NewOrderHandler создает новый экземпляр OrderHandler
func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// ------------------------------------------------------------
// Методы обработки запросов
// ------------------------------------------------------------

// CreateOrder обрабатывает запрос на создание заказа
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID, err := h.parseUserID(c)
	if err != nil {
		h.sendErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req models.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendErrorResponse(c, http.StatusBadRequest, models.ErrInvalidRequestFormat.Error())
		return
	}

	order, err := h.orderService.CreateOrder(userID, &req)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	h.sendOrderResponse(c, http.StatusCreated, order)
}

// GetUserOrders обрабатывает запрос на получение заказов пользователя
func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	userID, err := h.parseUserID(c)
	if err != nil {
		h.sendErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	orders, err := h.orderService.GetUserOrders(userID)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, h.mapToResponse(orders))
}

// ------------------------------------------------------------
// Вспомогательные методы
// ------------------------------------------------------------

func (h *OrderHandler) parseUserID(c *gin.Context) (uint, error) {
	id, err := strconv.Atoi(c.Param("user_id"))
	return uint(id), err
}

func (h *OrderHandler) handleServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrUserNotFound):
		h.sendErrorResponse(c, http.StatusNotFound, "User not found")
	case errors.Is(err, models.ErrProductRequired):
		h.sendErrorResponse(c, http.StatusBadRequest, "Product is required")
	case errors.Is(err, models.ErrInvalidQuantity):
		h.sendErrorResponse(c, http.StatusBadRequest, "Invalid quantity")
	case errors.Is(err, models.ErrInvalidPrice):
		h.sendErrorResponse(c, http.StatusBadRequest, "Invalid price")
	default:
		h.sendErrorResponse(c, http.StatusInternalServerError, "Internal server error")
	}
}

func (h *OrderHandler) mapToResponse(orders []models.Order) []models.OrderResponse {
	response := make([]models.OrderResponse, 0, len(orders))
	for _, order := range orders {
		response = append(response, models.OrderResponse{
			ID:        order.ID,
			UserID:    order.UserID,
			Product:   order.Product,
			Quantity:  order.Quantity,
			Price:     order.Price,
			CreatedAt: order.CreatedAt,
		})
	}
	return response
}

func (h *OrderHandler) sendOrderResponse(c *gin.Context, status int, order *models.Order) {
	c.JSON(status, models.OrderResponse{
		ID:        order.ID,
		UserID:    order.UserID,
		Product:   order.Product,
		Quantity:  order.Quantity,
		Price:     order.Price,
		CreatedAt: order.CreatedAt,
	})
}

func (h *OrderHandler) sendErrorResponse(c *gin.Context, status int, message string) {
	c.JSON(status, models.ErrorResponse{
		Error: message,
	})
}
