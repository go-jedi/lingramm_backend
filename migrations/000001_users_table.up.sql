CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY, -- Уникальный идентификатор.
    uuid TEXT NOT NULL UNIQUE, -- UUID пользователя.
    telegram_id TEXT NOT NULL UNIQUE, -- Telegram id пользователя.
    username VARCHAR(255) UNIQUE, -- Username пользователя.
    first_name VARCHAR(255), -- Имя пользователя.
    last_name VARCHAR(255), -- Фамилия пользователя.
    created_at TIMESTAMP NOT NULL DEFAULT NOW(), -- Дата создания записи.
    updated_at TIMESTAMP NOT NULL DEFAULT NOW() -- Дата обновления записи.
);