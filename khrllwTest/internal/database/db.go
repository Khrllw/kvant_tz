package database

import (
	"fmt"
	"khrllwTest/internal/middleware"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"khrllwTest/internal/models"
)

// Fixme not automigrate!

// ------------------------------------------------------------
// Структуры
// ------------------------------------------------------------

// Config содержит настройки подключения к БД
type Config struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
}

// loadDBConfig загружает конфигурацию БД из переменных окружения
func loadDBConfig() Config {
	return Config{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_PORT"),
	}
}

// ------------------------------------------------------------
// Основные функции
// ------------------------------------------------------------

// SetupDatabase
// Подключение к БД
func SetupDatabase(logConfig *middleware.LoggerConfig) (*gorm.DB, error) {
	config := loadDBConfig()

	// Добавляем задержку для поднятия БД
	time.Sleep(5 * time.Second)

	// Создаем логгер для GORM
	gormLogger := middleware.DBLogger(logConfig)

	db, err := connectToDB(config, gormLogger)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %v", err)
	}
	if err := configureConnectionPool(db); err != nil {
		return nil, err
	}
	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("ошибка миграций: %v", err)
	}
	log.Println("Соединение с базой данных установлено")
	return db, nil
}

// ------------------------------------------------------------
// Вспомогательные функции
// ------------------------------------------------------------

// connectToDB
// Устанавливает соединение с БД
func connectToDB(config Config, gormLogger logger.Interface) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		config.Host, config.User, config.Password, config.DBName, config.Port,
	)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
}

// configureConnectionPool
// Настраивает пул соединений БД
func configureConnectionPool(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get DB instance: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return nil
}

// autoMigrate
// Выполняет автоматическую миграцию моделей
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Order{},
	)
}
