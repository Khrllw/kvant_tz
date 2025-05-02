project/                                                            \
├── cmd/                                                            \
│   └── main.go          # Точка входа в приложение                 \
├── internal/                                                       \
│   ├── handlers/        # Обработчики HTTP-запросов                \
│   ├── models/          # Модели базы данных (GORM)                \
│   ├── repository/      # Работа с базой данных (GORM)             \
│   ├── services/        # Бизнес-логика                            \
│   ├── middleware/      # Middleware (JWT, логирование)            \
│   └── utils/           # Вспомогательные функции                  \
├── migrations/          # SQL-миграции (если используются)         \
├── go.mod               # Файл зависимостей                        \
├── .env                 # Переменные окружения                     \
├── Dockerfile           # Dockerfile для запуска приложения        \
└── docker-compose.yml   # Docker Compose (PostgreSQL + приложение) \
