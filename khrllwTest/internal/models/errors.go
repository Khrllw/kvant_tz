package models

import "errors"

// Ошибки приложения
var (
	// Общие ошибки
	ErrInvalidRequestFormat  = errors.New("неверный формат запроса")
	ErrEmailPasswordRequired = errors.New("email и пароль обязательны")
	ErrInternalServerError   = errors.New("внутренняя ошибка сервера")

	// Ошибки аутентификации
	ErrInvalidCredentials  = errors.New("неверные учетные данные")
	ErrInvalidToken        = errors.New("неверный токен")
	ErrTokenRequired       = errors.New("токен отсутствует")
	ErrInvalidTokenClaims  = errors.New("ErrInvalidTokenClaims")
	ErrInvalidFilterParams = errors.New("неверные параметры фильтра")

	// Ошибки пользователей
	ErrUserNotFound       = errors.New("пользователь не найден")
	ErrEmailAlreadyExists = errors.New("email уже существует")
	ErrEmail              = errors.New("email err")
	//ErrRecordNotFound     = errors.New("record не найдена")

	ErrInvalidPagination = errors.New("неверная пагинация")
	ErrRecordNotFound    = errors.New("запись не найдена")

	// Ошибки заказов
	ErrOrderNotFound = errors.New("заказ не найден")

	ErrInvalidPrice    = errors.New("неверная цена")
	ErrInvalidQuantity = errors.New("неверное количество")
	ErrProductRequired = errors.New("неверный продукт")
	//ErrRecordNotFound = errors.New("запись не найдена")

)
