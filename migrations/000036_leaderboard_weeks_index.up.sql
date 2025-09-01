-- Уникальный индекс.
-- Зачем: гарантирует, что у одного пользователя только одна строка на каждую неделю.
-- Когда помогает: при апдейте агрегата через upsert:
CREATE UNIQUE INDEX IF NOT EXISTS idx_leaderboard_weeks_week_start_telegram_id_week_user_uniq ON leaderboard_weeks (week_start, telegram_id);

-- Топ за неделю.
-- Зачем: хранит строки отсортированными по неделе и по убыванию XP, чтобы быстро отдавать топ.
-- Когда помогает: классический запрос «топ-100 за текущую неделю»:
CREATE INDEX IF NOT EXISTS idx_leaderboard_weeks_week_start_xp_week_top ON leaderboard_weeks (week_start, xp DESC) INCLUDE (telegram_id);

-- История пользователя.
-- Зачем: быстро отдаёт «мои недели» в обратном порядке по дате.
-- Когда помогает: лента/график для одного пользователя:
CREATE INDEX IF NOT EXISTS idx_leaderboard_weeks_telegram_id_week_start_user_history ON leaderboard_weeks (telegram_id, week_start DESC);