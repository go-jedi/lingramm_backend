CREATE TABLE IF NOT EXISTS user_balances( -- храним актуальный баланс пользователя.
    id BIGSERIAL PRIMARY KEY, -- Уникальный идентификатор.
    telegram_id TEXT NOT NULL UNIQUE, -- Telegram id пользователя.
    balance NUMERIC(20, 2) NOT NULL DEFAULT 0, -- Баланс пользователя.
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- Дата создания записи.
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- Дата обновления записи.
    FOREIGN KEY (telegram_id) REFERENCES users(telegram_id)
);