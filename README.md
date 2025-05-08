# 📦 KhrllwTest API

RESTful API на Go (Gin) с поддержкой Swagger-документации, PostgreSQL и Docker. Реализована базовая аутентификация,
управление пользователями и заказами.

---

## 📌 Оглавление

- [🌟 Особенности](#-особенности)
- [🛠 Технологии](#-технологии)
- [⚙️ Требования](#-требования)
- [📦 Установка и запуск](#-установка-и-запуск)
- [🔐 Переменные окружения](#-переменные-окружения)
- [🗂️ Структура проекта](#-структура-проекта)
- [🔄 Архитектура обработки](#-архитектура-обработки)
- [🧪 Тестирование](#-тестирование)
- [✅ Тестовое покрытие](#-тестовое-покрытие)

---

## 🌟 Особенности

- Аутентификация через JWT
- CRUD операции для пользователей
- Управление заказами
- Пагинация и фильтрация
- Swagger-документация
- Логирование запросов
- Разделение слоёв приложения (Handlers, Services, Repositories)

---

## 🛠 Технологии

- **Язык**: Go 1.20+
- **Фреймворк**: Gin
- **База данных**: PostgreSQL
- **ORM**: GORM
- **Документация**: Swagger (Swag)
- **Логирование**: Кастомный логгер с цветами
- **Контейнеризация**: Docker, Docker Compose

---

## ⚙️ Требования

Для работы с проектом необходимо установить следующие инструменты:

1. **Go (версии 1.20 и выше)**</br>
   Скачайте и установите Go с официального сайта:
   [Go 1.20+](https://go.dev/dl/)

2. **Docker и Docker Compose**</br>
   Для контейнеризации приложения используйте Docker и Docker Compose. Инструкции по установке можно найти на
   официальном сайте:
   [Docker и Compose](https://docs.docker.com/compose/install/)

3. **Swag CLI (для генерации документации в формате Swagger)**</br>
   Для автоматической генерации документации API с использованием Swagger, установите Swag CLI:

   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

   Дополнительную информацию о Swag CLI можно найти на GitHub:
   [Swaggo GitHub](https://github.com/swaggo/swag)

---

## 📦 Установка и запуск

Чтобы запустить проект локально, следуйте этим шагам:

1. **Клонируйте репозиторий**
   
   Сначала клонируйте репозиторий на вашу локальную машину:

   ```bash
   git clone https://github.com/Khrllw/kvant_tz.git
   ```

   Перейдите в каталог с проектом:

   ```bash
   cd kvant_tz
   ```

2. **Добавьте файл .env**

   Содержание файла указано [здесь](#-переменные-окружения)


3. **Сгенерируйте Swagger-документацию**
   Для генерации документации API в формате Swagger, выполните команду:

   ```bash
   swag init -g cmd/main.go
   ```

   Это создаст необходимые файлы для документации вашего API в проекте.

4. **Соберите и запустите проект**
   Используйте Docker Compose для сборки и запуска проекта:

   ```bash
   docker compose up --build
   ```

   Эта команда:

   * Соберет все контейнеры, указанные в `docker-compose.yml`
   * Запустит проект в контейнерах Docker

После выполнения этих команд приложение будет доступно в вашем браузере по адресу `http://localhost:8080` (или другому,
порту указанному в конфигурации).

---

## 🔐 Переменные окружения

Используется `.env` файл или переменные Docker:

| Переменная       | Описание                | Пример           |
|------------------|-------------------------|------------------|
| `JWT_KEY`        | Секретный ключ для JWT  | `supersecretkey` |
| `JWT_EXPIRATION` | Время жизни токена      | `24h`            |
| `DB_HOST`        | Хост базы данных        | `db`             |
| `DB_PORT`        | Порт базы данных        | `5432`           |
| `DB_USER`        | Пользователь PostgreSQL | `postgres_adm`   |
| `DB_PASSWORD`    | Пароль PostgreSQL       | `password`       |
| `DB_NAME`        | Название базы данных    | `khrllw_test`    |

---

## 🗂️ Структура проекта

```
project/
├── cmd/                   # Точка входа (main.go)
├── internal/
│   ├── handlers/          # Подключение БД
│   ├── handlers/          # HTTP обработчики
│   ├── models/            # Модели данных (GORM)
│   ├── repository/        # Работа с БД
│   ├── services/          # Бизнес-логика
│   ├── middleware/        # JWT, логирование
│   ├── utils/             # Хелперы (хеширование, токены)
│   └── migrations/        # SQL миграции
├── docs/                  # Swagger
├── .env                   # Конфигурация
├── tests                  # Тесты
├── Dockerfile             # Сборка Docker-образа
├── docker-compose.yml     # Compose-файл
├── go.mod/go.sum          # Зависимости
└── README.md              # Документация
```

---

## 🔄 Архитектура обработки

       1. Запрос от клиента
       2. Маршрутизатор Gin
       3. Промежуточное ПО (Авторизация / Логирование)
       4. Обработчики (UserHandler / OrderHandler)
       5. Сервисы (UserService / OrderService)
       6. Репозитории (UserRepository / OrderRepository)
       7. База данных PostgreSQL

---

## 🧪 Тестирование

*Добавьте инструкцию по запуску тестов, если применимо*

Чтобы запустить тесты, выполните следующие шаги:

1. Запустите Docker контейнер.
2. Перейдите в папку `tests`.
3. Запустите файлы, которые заканчиваются на `_test`, используя команду Go:

```bash
  go test ./tests/*_test.go
```

---

## ✅ Тестовое покрытие

### 🛡️ JWT

**Аутентификация**

* `Test1_ValidJWTLogin`
* `Test2_LoginWithWrongPassword`
* `Test3_LoginWithUnknownEmail`

**Защищённые маршруты**

* `Test4_RequestWithoutToken`
* `Test5_RequestWithInvalidTokenFormat`
* `Test9_MissingBearerPrefix`
* `Test12_ValidTokenAccess`

**Некорректные токены**

* `Test6_RequestWithTamperedToken`
* `Test7_ExpiredTokenSimulation`
* `Test8_UseTokenOfAnotherUser`
* `Test10_TokenReuseAfterUserDeletion`
* `Test11_MalformedJWTStructure`
* `Test13_TokenWithFutureIat`
* `Test14_TokenWithNoneAlg`
* `Test15_TokenWithInvalidSignature`
* `Test16_TokenInBodyInsteadOfHeader`
* `Test18_TokenWithExtraClaims`
* `Test19_TokenReuseMultipleTimes`

### 👤 Пользователи

**Создание**

* `TestUser1_CreateValidUser`
* `TestUser2_CreateDuplicateUser`
* `TestUser9_CreateUserWithEmptyName`
* `TestUser10_CreateUserWithInvalidAge`

**Получение**

* `TestUser3_GetUserByID`
* `TestUser6_GetUsersList`
* `TestUser7_GetUsersWithPagination`

**Обновление**

* `TestUser4_UpdateUserName`
* `TestUser8_UpdateUserInvalidEmail`
* `TestUser13_UpdateOtherUser`

**Удаление**

* `TestUser5_DeleteUser`
* `TestUser14_DeleteOtherUser`

### 📦 Заказы

**Создание**

* `TestOrder1_CreateOrder`
* `TestOrder6_CreateOrderWithoutAuth`
* `TestOrder10_CreateOrderWithInvalidPayload`

**Получение**

* `TestOrder2_ListOrders`
* `TestOrder3_GetSingleOrderByID`
* `TestOrder7_GetOtherUsersOrder`

**Обновление**

* `TestOrder4_UpdateOrder`
* `TestOrder8_UpdateOrderNotOwned`

---

