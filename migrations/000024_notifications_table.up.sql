CREATE TABLE IF NOT EXISTS notifications(
    id BIGSERIAL PRIMARY KEY, -- Уникальный идентификатор.
    type notifications_type NOT NULL, -- Тип уведомления.
    telegram_id TEXT NOT NULL, -- Telegram id пользователя.
    status notifications_status NOT NULL DEFAULT 'PENDING', -- Статус: PENDING, SENT, FAILED.
    message JSONB NOT NULL, -- Тело уведомления в формате JSON.
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- Дата создания записи.
    sent_at TIMESTAMP WITH TIME ZONE, -- Дата отправки сообщения.
    FOREIGN KEY (telegram_id) REFERENCES users(telegram_id)
);