CREATE TABLE IF NOT EXISTS leaderboard_weeks_applied_events(
    id BIGSERIAL PRIMARY KEY, -- Уникальный идентификатор.
    event_id BIGINT NOT NULL, -- это xp_events.id.
    applied_at TIMESTAMPTZ NOT NULL DEFAULT now() -- когда учли.
);

-- Делает обработку эффективно «ровно один раз». Если воркер по любой
-- причине снова возьмёт те же xp_events.id, мы вставим их в эту таблицу
-- с ON CONFLICT DO NOTHING и поймём, что они уже учтены → не прибавим XP второй раз.