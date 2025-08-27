CREATE TABLE IF NOT EXISTS balance_transactions( -- История операций по балансу.
    id BIGSERIAL PRIMARY KEY, -- Уникальный идентификатор.
    event_type_id BIGINT NOT NULL, -- Идентификатор события.
    telegram_id TEXT NOT NULL, -- Telegram id пользователя.
    amount NUMERIC(20, 2) NOT NULL, -- Сумма.
    description TEXT, -- Описание.
    balance_after NUMERIC(20, 2), -- Баланс пользователя после транзакции (для аудита)
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- Дата создания записи.
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- Дата обновления записи.
    FOREIGN KEY (event_type_id) REFERENCES event_types(id),
    FOREIGN KEY (telegram_id) REFERENCES users(telegram_id)
);