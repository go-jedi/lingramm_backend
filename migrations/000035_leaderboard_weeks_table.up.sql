CREATE TABLE IF NOT EXISTS leaderboard_weeks( -- Недельный агрегат (Лидер борд).
    id BIGSERIAL PRIMARY KEY, -- Уникальный идентификатор.
    telegram_id TEXT NOT NULL, -- Telegram id пользователя. Кто получил/потерял XP.
    xp BIGINT NOT NULL DEFAULT 0, -- Сумма XP за неделю.
    week_start DATE NOT NULL, -- Понедельник недели.
    FOREIGN KEY (telegram_id) REFERENCES users(telegram_id)
);