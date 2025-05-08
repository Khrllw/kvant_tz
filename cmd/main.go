// @title KhrllwTest API
// @version 1.0
// @description API для управления пользователями и заказами

// @host localhost:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"io"
	"khrllwTest/internal/db"
	"khrllwTest/internal/handlers"
	"khrllwTest/internal/middleware"
	"khrllwTest/internal/repository"
	service "khrllwTest/internal/services"
	"khrllwTest/internal/utils"
	"log"
	"os"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "khrllwTest/docs"
)

// loadEnv загружает переменные окружения
func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found - using system environment variables")
	}
}

// getPort получает порт из переменных окружения
func getPort() string {
	if port := os.Getenv("API_PORT"); port != "" {
		return port
	}
	return "8080"
}

// initDatabase инициализирует подключение к БД
func initDatabase(logConfig *middleware.LoggerConfig) *gorm.DB {
	db, err := db.SetupDatabase(logConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	return db
}

// setupRouter настраивает маршруты API
func setupRouter(userHandler *handlers.UserHandler,
	orderHandler *handlers.OrderHandler,
	loginHandler *handlers.LoginHandler,
	authorization *middleware.Authorization,
	logConfig *middleware.LoggerConfig) *gin.Engine {

	router := gin.Default()

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ------------------------- Обработка запросов -------------------------

	// Роут для авторизации пользователя
	router.POST("/auth/login", loginHandler.Login)

	// Роут для создания пользователя (без авторизации)
	router.POST("/users", userHandler.CreateUser)

	// Группа для работы с пользователями (требует авторизации)
	usersGroup := router.Group("/users")
	usersGroup.Use(authorization.Middleware())
	{
		usersGroup.GET("", userHandler.GetUsers)

		usersIDGroup := usersGroup.Group("/:user_id")

		// Подключение логирования для всех запросов группы
		usersIDGroup.Use(middleware.RequestLogger(logConfig))
		{
			usersIDGroup.GET("", userHandler.GetUserByID)
			usersIDGroup.PUT("", userHandler.UpdateUser)
			usersIDGroup.DELETE("", userHandler.DeleteUser)
			usersIDGroup.GET("/orders", orderHandler.GetUserOrders)
			usersIDGroup.POST("/orders", orderHandler.CreateOrder)
		}
	}

	return router
}

// startServer запускает HTTP сервер
func startServer(router *gin.Engine, port string) {
	log.Printf("Starting server on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func main() {
	// Настройка логирования
	logFile, _ := os.OpenFile("api.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	logConfig := &middleware.LoggerConfig{
		Output:   multiWriter, // Вывод в файл и консоль
		Colorful: true,        // Цветной вывод
		LogSQL:   true,        // Логировать SQL запросы
	}

	// ------------------------- INIT -------------------------
	loadEnv()
	port := getPort()

	// Инициализация БД
	db := initDatabase(logConfig)

	userRepo := repository.NewUserRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	// Инициализация обработчиков
	passHasher := utils.NewPasswordHasher(0)
	userService := service.NewUserService(userRepo, passHasher)
	userHandler := handlers.NewUserHandler(userService)

	orderService := service.NewOrderService(orderRepo, userRepo)
	orderHandler := handlers.NewOrderHandler(orderService)

	authConfig, err := utils.NewJWTConfig()
	if err != nil {
		log.Fatalf("Ошибка инициализации конфигурации аутентификации: %v", err)
	}

	tokenManager := utils.NewTokenManager(authConfig)
	authService := service.NewLoginService(userRepo, tokenManager, passHasher)
	authHandler := handlers.NewLoginHandler(authService)

	authorizationMiddleware := middleware.NewAuthorization(tokenManager, userRepo)

	// ----------------- ROUTER -----------------
	// Настройка роутера
	router := setupRouter(userHandler, orderHandler, authHandler, authorizationMiddleware, logConfig)

	// ------------------ RUN ------------------
	// Запуск сервера
	startServer(router, port)
}
