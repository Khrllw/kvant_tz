package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"khrllwTest/internal/database"
	"khrllwTest/internal/handlers"
	"khrllwTest/internal/middleware"
	"khrllwTest/internal/repository"
	"log"
)

func main() {

	// ----------------- INIT -----------------
	PORT := ":8080"

	// Загрузка переменных окружения
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// 2. Инициализируем БД
	db := database.SetupDatabase()

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	// Инициализация обработчиков
	userHandler := handlers.NewUserHandler(userRepo)
	orderHandler := handlers.NewOrderHandler(orderRepo, userRepo)

	// ----------------- ROUTER -----------------

	router := gin.Default()

	// ----------- Обработка запросов -----------
	router.POST("/auth/login")

	// Каскадная группировка запросов
	usersGroup := router.Group("/users")
	{
		usersGroup.GET("", userHandler.GetUsers)
		usersGroup.POST("", userHandler.AddUser)

		usersIDGroup := usersGroup.Group("/:user_id")
		// Подключение логирования для всех запросов группы
		usersIDGroup.Use(middleware.RequestIDLogger())
		usersIDGroup.Use(middleware.ConsoleLogger())
		{
			usersIDGroup.GET("", userHandler.GetUserFromID)
			usersIDGroup.PUT("", userHandler.UpdateUser)
			usersIDGroup.DELETE("", userHandler.DeleteUser)
			usersIDGroup.GET("/orders", orderHandler.GetUserOrders)
			usersIDGroup.POST("/orders", orderHandler.CreateOrder)
		}
	}

	// ------------------ RUN ------------------
	// Запуск на localhost
	err := router.Run(PORT)
	if err != nil {
		return
	}
}
