package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"khrllwTest/internal/models"
	"khrllwTest/internal/repository"
)

type OrderHandler struct {
	orderRepo repository.OrderRepository
	userRepo  repository.UserRepository
}

func NewOrderHandler(orderRepo repository.OrderRepository, userRepo repository.UserRepository) *OrderHandler {
	return &OrderHandler{
		orderRepo: orderRepo,
		userRepo:  userRepo,
	}
}

// CreateOrder создает новый заказ для пользователя
// @Summary Создать заказ
// @Description Создает новый заказ для указанного пользователя
// @Tags orders
// @Accept json
// @Produce json
// @Param user_id path int true "ID пользователя"
// @Param input body models.CreateOrderRequest true "Данные заказа"
// @Security ApiKeyAuth
// @Success 201 {object} models.OrderResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{user_id}/orders [post]
func (handler *OrderHandler) CreateOrder(context *gin.Context) {
	// Получаем ID пользователя из URL
	userID, err := strconv.Atoi(context.Param("user_id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID пользователя"})
		return
	}

	// Проверяем существование пользователя
	if _, err := handler.userRepo.FindByID(uint(userID)); err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	// Парсим тело запроса
	var req models.CreateOrderRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Создаем заказ
	order := models.Order{
		UserID:   uint(userID),
		Product:  req.Product,
		Quantity: req.Quantity,
		Price:    req.Price,
	}

	if err := handler.orderRepo.Create(&order); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании заказа"})
		return
	}

	// Возвращаем созданный заказ
	context.JSON(http.StatusCreated, models.OrderResponse{
		ID:        order.ID,
		UserID:    order.UserID,
		Product:   order.Product,
		Quantity:  order.Quantity,
		Price:     order.Price,
		CreatedAt: order.CreatedAt,
	})
}

// GetUserOrders возвращает список заказов пользователя
// @Summary Получить заказы пользователя
// @Description Возвращает список всех заказов указанного пользователя
// @Tags orders
// @Produce json
// @Param user_id path int true "ID пользователя"
// @Security ApiKeyAuth
// @Success 200 {array} models.OrderResponse
// @Failure 404 {object} map[string]string
// @Router /users/{user_id}/orders [get]
func (handler *OrderHandler) GetUserOrders(context *gin.Context) {
	// Получаем ID пользователя из URL
	userID, err := strconv.Atoi(context.Param("user_id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID пользователя"})
		return
	}

	// Проверяем существование пользователя
	if _, err := handler.userRepo.FindByID(uint(userID)); err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	// Получаем заказы пользователя
	orders, err := handler.orderRepo.FindByUserID(uint(userID))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении заказов"})
		return
	}

	// Формируем ответ
	var response []models.OrderResponse
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

	context.JSON(http.StatusOK, response)
}
