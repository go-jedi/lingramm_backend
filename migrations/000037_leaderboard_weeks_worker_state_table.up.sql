CREATE TABLE IF NOT EXISTS leaderboard_weeks_worker_state(
    id BIGSERIAL PRIMARY KEY, -- Уникальный идентификатор.
    name TEXT NOT NULL, -- имя воркера/задачи (например, 'main').
    last_event_id BIGINT NOT NULL, -- Самый большой xp_events.id, чьи эффекты уже в таблице leaderboard_weeks.
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), -- Когда чекпоинт обновили.
    CONSTRAINT ux_leaderboard_weeks_worker_state_name UNIQUE (name)
);

INSERT INTO leaderboard_weeks_worker_state(name, last_event_id) VALUES('worker_1', 0);

-- Таблица держит «до какого события из xp_events мы гарантированно всё учли в агрегате».
-- Это позволяет воркеру после рестарта продолжать с нужного места, а не с нуля.