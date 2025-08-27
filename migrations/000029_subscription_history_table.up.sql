CREATE TABLE IF NOT EXISTS subscription_history(
    id BIGSERIAL PRIMARY KEY, -- Уникальный идентификатор.
    telegram_id TEXT NOT NULL, -- Telegram id пользователя.
    action_time TIMESTAMP NOT NULL, -- Время совершения действия.
    expires_at TIMESTAMP NOT NULL, -- Дата окончания подписки.
    created_at TIMESTAMP NOT NULL DEFAULT NOW(), -- Дата создания записи.
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(), -- Дата обновления записи.
    FOREIGN KEY (telegram_id) REFERENCES users(telegram_id)
);