-- Быстрее фильтр активных задач.
CREATE INDEX IF NOT EXISTS idx_daily_tasks_is_active ON daily_tasks (is_active);