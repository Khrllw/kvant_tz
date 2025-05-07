package middleware

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm/logger"
)

const (
	colorReset         = "\033[0m"
	colorWhite         = "\033[97m"
	colorGreen         = "\033[42m"
	colorRed           = "\033[41m"
	colorCyan          = "\033[46m"
	colorBold          = "\033[1m"
	timeFormat         = "2006/01/02 - 15:04:05"
	slowQueryThreshold = 200 * time.Millisecond
)

// LoggerConfig конфигурация для логгеров
type LoggerConfig struct {
	Output   io.Writer
	Colorful bool
	LogSQL   bool
}

// defaultConfig возвращает конфигурацию по умолчанию
func defaultConfig() *LoggerConfig {
	return &LoggerConfig{
		Output:   os.Stdout,
		Colorful: true,
		LogSQL:   true,
	}
}

// ----------------- HTTP ЛОГГЕР -----------------

// RequestLogger создает middleware для логирования HTTP запросов
func RequestLogger(config *LoggerConfig) gin.HandlerFunc {
	if config == nil {
		config = defaultConfig()
	}

	return func(c *gin.Context) {
		start := time.Now()
		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		logRequestStart(config, c, requestID)
		c.Next()
		logRequestEnd(config, c, requestID, start)
	}
}

// colorizeMethod применяет цвет только к названию метода
func colorizeMethod(method, color string) string {
	return fmt.Sprintf("%s%s%s[ %s ]%s",
		colorBold, colorWhite, color,
		method,
		colorReset)
}

// formatLogMessage форматирует сообщение лога
func formatLogMessage(method, timeStr, path, requestID string, colorful bool) string {
	var coloredMethod string
	if colorful {
		coloredMethod = colorizeMethod(method, colorCyan)
	} else {
		coloredMethod = fmt.Sprintf("[ %s ]", method)
	}

	return fmt.Sprintf("%s\t%s\t%s\n%s\t%s",
		coloredMethod,
		timeStr,
		path,
		coloredMethod,
		requestID,
	)
}

// Записывает сообщение в лог с обработкой ошибок
func writeLog(output io.Writer, msg string) {
	if _, err := fmt.Fprintln(output, msg); err != nil {
		log.Printf("Failed to write log: %v", err)
	}
}

// Логирует начало HTTP запроса
func logRequestStart(config *LoggerConfig, c *gin.Context, requestID string) {
	msg := formatLogMessage(
		c.Request.Method,
		time.Now().Format(timeFormat),
		c.Request.URL.Path,
		requestID,
		config.Colorful,
	)

	writeLog(config.Output, msg)
}

// formatStatusMessage форматирует сообщение о статусе
func formatStatusMessage(method string, status int, requestID string, latency time.Duration, colorful bool) string {
	var statusLabel, color string

	/*if status >= 400 {
		statusLabel = "ERR"
		color = colorRed
	} else {*/

	statusLabel = "✓✓✓"
	color = colorGreen

	var coloredLabel string
	if colorful {
		coloredLabel = colorizeMethod(statusLabel, color)
	} else {
		coloredLabel = fmt.Sprintf("[ %s ]", statusLabel)
	}

	return fmt.Sprintf("%s\tCode: %d | %s | %v",
		coloredLabel,
		status,
		requestID,
		latency)
}

// Логирует ошибки запроса
func logRequestErrors(config *LoggerConfig, c *gin.Context, requestID string) {
	for _, err := range c.Errors {
		var errMsg string
		if config.Colorful {
			errMsg = fmt.Sprintf("%s\t%s: %v",
				colorizeMethod("ERR", colorRed),
				requestID,
				err)
		} else {
			errMsg = fmt.Sprintf("[ ERR ]\t%s: %v", requestID, err)
		}
		writeLog(config.Output, errMsg)
	}
}

// Логирует завершение HTTP запроса
func logRequestEnd(config *LoggerConfig, c *gin.Context, requestID string, start time.Time) {
	status := c.Writer.Status()
	latency := time.Since(start)

	statusMsg := formatStatusMessage(
		c.Request.Method,
		status,
		requestID,
		latency,
		config.Colorful,
	)

	writeLog(config.Output, statusMsg)
	//logRequestErrors(config, c, requestID)
}

// ----------------- DB ЛОГГЕР -----------------

// DBLogger создает настроенный логгер для GORM
func DBLogger(config *LoggerConfig) logger.Interface {
	if config == nil {
		config = defaultConfig()
	}

	return logger.New(
		log.New(config.Output, "[ DB ]  ", log.LstdFlags),
		logger.Config{
			SlowThreshold:             slowQueryThreshold,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false,
			ParameterizedQueries:      config.LogSQL,
			Colorful:                  config.Colorful,
		},
	)
}
