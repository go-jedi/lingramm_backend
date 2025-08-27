CREATE TABLE IF NOT EXISTS users_blacklist(
    id BIGSERIAL PRIMARY KEY, -- Уникальный идентификатор.
    telegram_id TEXT NOT NULL UNIQUE, -- Telegram id пользователя.
    ban_timestamp TIMESTAMP NOT NULL DEFAULT NOW(), -- Время, когда пользователь был забанен.
    ban_reason TEXT NOT NULL, -- Причина бана.
    banned_by_telegram_id TEXT NOT NULL, -- Кто забанил пользователя.
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- Дата создания записи.
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- Дата обновления записи.
    FOREIGN KEY (telegram_id) REFERENCES users(telegram_id)
);