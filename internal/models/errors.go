package models

import "errors"

// ---------------------------------- ОШИБКИ API ----------------------------------

var (
	// ---------------------------- Общие ошибки -------------------------

	ErrInvalidRequestFormat = errors.New("Неверный формат запроса. ")
	ErrInternalServerError  = errors.New("Внутренняя ошибка сервера. ")
	ErrDatabaseError        = errors.New("Ошибка базы данных. ")

	// ------------------------ Ошибки авторизации -----------------------

	ErrEmailPasswordRequired = errors.New("Email и пароль обязательны. ")
	ErrInvalidCredentials    = errors.New("Неверные учетные данные. ")
	ErrInvalidToken          = errors.New("Неверный токен авторизации. ")
	ErrTokenRequired         = errors.New("Токен авторизации отсутствует. ")
	ErrTokenGenerationFailed = errors.New("Ошибка создания токена авторизации. ")
	ErrPasswordHashFailed    = errors.New("Ошибка хеширования пароля. ")

	ErrInvalidTokenClaims = errors.New("некорректное содержимое токена")

	// ---------------------- Ошибки пользователей -----------------------

	ErrUserNotFound        = errors.New("Пользователь не найден. ")
	ErrInvalidUserID       = errors.New("Неверный ID пользователя. ")
	ErrEmailAlreadyExists  = errors.New("Пользователь с данным email уже существует. ")
	ErrInvalidUserName     = errors.New("Некорректное имя пользователя. ")
	ErrInvalidUserEmail    = errors.New("Некорректные email пользователя. ")
	ErrInvalidUserPassword = errors.New("Некорректные пароль пользователя. ")
	ErrInvalidUserAge      = errors.New("Некорректный возраст пользователя. ")

	ErrInvalidPagination   = errors.New("Некорректные параметры пагинации. ")
	ErrInvalidFilterParams = errors.New("Некорректные параметры фильтрации. ")

	// ------------------------- Репозитории -----------------------------

	ErrRecordNotFound = errors.New("Запись не найдена. ")

	// ------------------------- Ошибки заказов -------------------------

	ErrInvalidPrice    = errors.New("Некорректная цена. ")
	ErrInvalidQuantity = errors.New("Некорректное количество. ")
	ErrProductRequired = errors.New("Некорректное название продукта. ")
)
