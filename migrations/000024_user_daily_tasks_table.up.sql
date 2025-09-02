CREATE TABLE IF NOT EXISTS user_daily_tasks( -- Содержит ежедневные задания, которые пользователь выполняет/выполнил.
    id BIGSERIAL PRIMARY KEY, -- Уникальный идентификатор.
    daily_task_id BIGINT NOT NULL, -- Идентификатор ежедневного задания.
    telegram_id TEXT NOT NULL, -- Telegram id пользователя.
    words_learned BIGINT NOT NULL DEFAULT 0, -- Сколько слов выучено.
    tasks_completed BIGINT NOT NULL DEFAULT 0, -- Сколько заданий выполнено.
    lessons_finished BIGINT NOT NULL DEFAULT 0, -- Пройдено уроков.
    words_translate BIGINT NOT NULL DEFAULT 0, -- Переведено новых слов.
    dialog_completed BIGINT NOT NULL DEFAULT 0, -- Пройдено диалогов.
    experience_points BIGINT NOT NULL DEFAULT 0, -- Шкала опыта.
    is_completed BOOLEAN NOT NULL DEFAULT FALSE, -- Было ли выполнено ежедневное задание.
    occurred_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(), -- Когда событие произошло.
    week_start DATE GENERATED ALWAYS AS (date_trunc('week', occurred_at AT TIME ZONE 'Europe/Moscow')::DATE) STORED,
    FOREIGN KEY (daily_task_id) REFERENCES daily_tasks(id)
);