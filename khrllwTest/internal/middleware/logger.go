package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/pborman/uuid"
	"time"
)

// ----------------- ЛОГИРОВАНИЕ -----------------

// RequestIDLogger задание уникального номера для каждого запроса
func RequestIDLogger() gin.HandlerFunc {
	return func(context *gin.Context) {
		requestID := uuid.New()
		context.Set("request_id", requestID)
		context.Next()
	}
}

// ConsoleLogger логирование в консоль
func ConsoleLogger() gin.HandlerFunc {
	return func(context *gin.Context) {
		requestID, _ := context.Get("request_id")
		// Действия перед обработкой запроса
		println(
			"\033[1;97;42m[", context.Request.Method, "]\033[0m\033[97m\t",
			time.Now().Format("2006/01/02 - 15:04:05\t"),
			context.Request.URL.Path,
			"\n\033[1;97;42m[", context.Request.Method, "]\033[0m\033[97m\t",
			requestID.(string))

		// Передаем управление следующему обработчику
		context.Next()

		// Действия после обработки запроса
		if len(context.Errors) > 0 || context.Writer.Status() >= 400 {
			println("\033[1;97;41m[ ERR ]\u001B[0m\u001B[97m\t",
				"Code: ", context.Writer.Status(), " | ",
				requestID.(string))
		} else {
			println("\033[1;97;44m[ ✓✓✓ ]\u001B[0m\u001B[97m\t", requestID.(string))
		}
	}
}
