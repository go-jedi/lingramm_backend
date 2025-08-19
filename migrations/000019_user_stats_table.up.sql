CREATE TABLE IF NOT EXISTS user_stats ( -- Хранение прогресса пользователя.
    id BIGSERIAL PRIMARY KEY, -- Уникальный идентификатор.
    telegram_id TEXT NOT NULL UNIQUE, -- Telegram id пользователя.
    streak_days BIGINT NOT NULL DEFAULT 0, -- Сколько дней подряд заходил.
    words_learned BIGINT NOT NULL DEFAULT 0, -- Сколько слов выучено.
    tasks_completed BIGINT NOT NULL DEFAULT 0, -- Сколько заданий выполнено.
    lessons_finished BIGINT NOT NULL DEFAULT 0, -- Пройдено уроков.
    experience_points BIGINT NOT NULL DEFAULT 0, -- Шкала опыта.
    level BIGINT NOT NULL DEFAULT 1, -- Уровень пользователя.
    last_active_at TIMESTAMP WITH TIME ZONE, -- Последняя активность.
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- Дата создания записи.
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- Дата обновления записи.
    FOREIGN KEY (telegram_id) REFERENCES users(telegram_id)
);