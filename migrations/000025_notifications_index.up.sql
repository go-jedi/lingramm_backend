-- Быстрый поиск уведомлений для конкретного пользователя и статуса
CREATE INDEX IF NOT EXISTS idx_notifications_telegram_status
    ON notifications (telegram_id, status);

-- Индекс по дате создания, если нужны сортировки
CREATE INDEX IF NOT EXISTS idx_notifications_created_at
    ON notifications (created_at DESC);

-- Индекс по типу уведомления, если часто фильтруются по типу
CREATE INDEX IF NOT EXISTS idx_notifications_type
    ON notifications (type);