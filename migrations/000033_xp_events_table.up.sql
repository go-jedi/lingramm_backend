CREATE TABLE IF NOT EXISTS xp_events( -- журнал событий XP.
    id BIGSERIAL PRIMARY KEY, -- Уникальный идентификатор.
    event_type_id BIGINT NOT NULL, -- Идентификатор события.
    telegram_id TEXT NOT NULL, -- Telegram id пользователя. Кто получил/потерял XP.
    delta_xp INTEGER NOT NULL CHECK (delta_xp <> 0), -- +XP / -XP за событие.
    occurred_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(), -- Когда событие произошло.
    inserted_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(), -- Когда попало в БД.
    week_start DATE GENERATED ALWAYS AS (date_trunc('week', occurred_at AT TIME ZONE 'Europe/Moscow')::DATE) STORED,
    FOREIGN KEY (event_type_id) REFERENCES event_types(id),
    FOREIGN KEY (telegram_id) REFERENCES users(telegram_id)
);

-- Неизменяемая хроника всех изменений XP. Из неё можно пересчитать
-- любые агрегаты (неделя, месяц, сезон), отладить спорные ситуации
-- и делать аналитику.

-- Строка = одно событие изменения XP.

-- по week_start:
-- Начало недели (понедельник 00:00) в выбранном TZ.
-- Поменяйте 'UTC' на нужный бизнес-TZ (например, 'Europe/Moscow').