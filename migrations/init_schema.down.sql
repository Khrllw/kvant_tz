-- Откатываем изменения в обратном порядке
DROP INDEX IF EXISTS idx_orders_user_id;
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS users;