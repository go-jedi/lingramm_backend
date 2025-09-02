-- За сегодня/по дате в MSK.
CREATE INDEX IF NOT EXISTS idx_user_daily_tasks_telegram_id_occurred_at ON user_daily_tasks (telegram_id, ((occurred_at AT TIME ZONE 'Europe/Moscow')::DATE));

-- Для анти-повтора.
CREATE INDEX IF NOT EXISTS idx_user_daily_tasks_telegram_id_daily_task_id_occurred_at ON user_daily_tasks (telegram_id, daily_task_id, occurred_at);

-- Идемпотентность дня.
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_daily_tasks_telegram_id_occurred_at_unique ON user_daily_tasks (telegram_id, ((occurred_at AT TIME ZONE 'Europe/Moscow')::DATE));