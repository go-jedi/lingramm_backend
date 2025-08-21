CREATE TABLE IF NOT EXISTS leaderboard_weeks( -- Недельный агрегат (Лидер борд).
    id BIGSERIAL PRIMARY KEY, -- Уникальный идентификатор.
    telegram_id TEXT NOT NULL, -- Telegram id пользователя. Кто получил/потерял XP.
    xp BIGINT NOT NULL DEFAULT 0, -- Сумма XP за неделю.
    week_start DATE NOT NULL, -- Понедельник недели.
    FOREIGN KEY (telegram_id) REFERENCES users(telegram_id),
    CONSTRAINT leaderboard_weeks_week_start_telegram_id_uniq UNIQUE (week_start, telegram_id),
    CONSTRAINT fk_leaderboard_weeks_telegram_id FOREIGN KEY (telegram_id) REFERENCES public.users(telegram_id)
);