CREATE TABLE IF NOT EXISTS user_daily_tasks(
    id BIGSERIAL PRIMARY KEY, -- Уникальный идентификатор.
    daily_task_id BIGINT NOT NULL, -- Идентификатор ежедневного задания.
    telegram_id TEXT NOT NULL, -- Telegram id пользователя.
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- Дата создания записи.
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- Дата обновления записи.
    FOREIGN KEY (daily_task_id) REFERENCES daily_tasks(id)
);