CREATE TABLE IF NOT EXISTS balance_transactions( -- История операций по балансу.
    id SERIAL PRIMARY KEY, -- Уникальный идентификатор.
    event_id BIGINT NOT NULL, -- Идентификатор события.
    telegram_id TEXT NOT NULL, -- Telegram id пользователя.
    amount NUMERIC(20, 2) NOT NULL, -- Сумма.
    transaction_type VARCHAR(50) NOT NULL, -- Например: 'purchase', 'bonus', 'admin_credit', 'refund'.
    source_id BIGINT, -- ID источника (например, order_id).
    description TEXT, -- Описание.
    balance_after NUMERIC(20, 2), -- Баланс пользователя после транзакции (для аудита)
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- Дата создания записи.
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- Дата обновления записи.
    FOREIGN KEY (event_id) REFERENCES balance_events(id) ON DELETE SET NULL,
    FOREIGN KEY (telegram_id) REFERENCES users(telegram_id) ON DELETE CASCADE
);