-- +goose Up
-- Создаем таблицу users
CREATE TABLE IF NOT EXISTS users
(
    id            SERIAL PRIMARY KEY,
    name          VARCHAR(255)        NOT NULL,
    email         VARCHAR(255) UNIQUE NOT NULL,
    age           INT                 NOT NULL,
    password_hash VARCHAR(255)        NOT NULL
);

-- Создаем таблицу orders
CREATE TABLE IF NOT EXISTS orders
(
    id         SERIAL PRIMARY KEY,
    user_id    INT            NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    product    VARCHAR(255)   NOT NULL,
    quantity   INT            NOT NULL,
    price      DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Индекс для ускорения поиска пользователей по email
CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);

-- Индекс для ускорения поиска заказов по пользователю
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders (user_id);