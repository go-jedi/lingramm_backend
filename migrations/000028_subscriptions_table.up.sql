CREATE TABLE IF NOT EXISTS subscriptions(
    id BIGSERIAL PRIMARY KEY, -- Уникальный идентификатор.
    telegram_id TEXT NOT NULL UNIQUE, -- Telegram id пользователя.
    subscribed_at TIMESTAMP, -- Дата начала подписки.
    expires_at TIMESTAMP, -- Дата окончания подписки.
    is_active BOOLEAN DEFAULT FALSE, -- Флаг, указывающий, активна ли подписка.
    created_at TIMESTAMP NOT NULL DEFAULT NOW(), -- Дата создания записи.
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(), -- Дата обновления записи.
    FOREIGN KEY (telegram_id) REFERENCES users(telegram_id)
);