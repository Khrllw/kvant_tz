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
// @Tags Orders
// @Summary Создать новый заказ
// @Description Создает новый заказ для пользователя
// @Accept json
// @Produce json
// @Param user_id path int true "ID пользователя"
// @Param order body models.CreateOrderRequest true "Данные заказа"
// @Success 201 {object} models.OrderResponse
// @Failure 400 {object} models.ErrorLoginResponse "Неверный формат запроса/некорректные данные"
// @Failure 500 {object} models.ErrorLoginResponse "Внутренняя ошибка сервера"
// @Router /users/{user_id}/orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID, err := h.parseUserID(c)
	if err != nil {
		h.sendErrorResponse(c, http.StatusBadRequest, models.ErrInvalidUserID)
		return
	}

	var req models.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendErrorResponse(c, http.StatusBadRequest, models.ErrInvalidRequestFormat)
		return
	}

	order, err := h.orderService.CreateOrder(userID, &req)
	if err != nil {
		if errors.Is(err, models.ErrDatabaseError) {
			h.sendErrorResponse(c, http.StatusInternalServerError, models.ErrInternalServerError)
			return
		}
		h.sendErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	h.sendOrderResponse(c, http.StatusCreated, order)
}

// GetUserOrders обрабатывает запрос на получение заказов пользователя
// @Tags Orders
// @Summary Получить все заказы пользователя
// @Description Возвращает все заказы для конкретного пользователя по его ID
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {array} models.OrderResponse
// @Failure 400 {object} models.ErrorLoginResponse "Неверный формат запроса/некорректные данные"
// @Failure 500 {object} models.ErrorLoginResponse "Внутренняя ошибка сервера"
// @Router /users/{user_id}/orders [get]
func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	userID, err := h.parseUserID(c)
	if err != nil {
		h.sendErrorResponse(c, http.StatusBadRequest, models.ErrInvalidUserID)
		return
	}

	orders, err := h.orderService.GetUserOrders(userID)
	if err != nil {
		if errors.Is(err, models.ErrDatabaseError) {
			h.sendErrorResponse(c, http.StatusInternalServerError, models.ErrInternalServerError)
			return
		}
		h.sendErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, h.mapToResponse(orders))
}

// ------------------------------------------------------------
// Вспомогательные методы
// ------------------------------------------------------------

// parseUserID парсит ID пользователя из URL
func (h *OrderHandler) parseUserID(c *gin.Context) (uint, error) {
	id, err := strconv.Atoi(c.Param("user_id"))
	return uint(id), err
}

// mapToResponse преобразует заказы в формат ответа
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

// sendOrderResponse отправляет ответ с заказом
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

// sendErrorResponse отправляет ответ с ошибкой
func (h *OrderHandler) sendErrorResponse(c *gin.Context, status int, err error) {
	c.JSON(status, models.ErrorLoginResponse{
		Error: err.Error(),
	})
}
