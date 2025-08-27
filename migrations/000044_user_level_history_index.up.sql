CREATE INDEX IF NOT EXISTS idx_user_level_history_telegram_id_reached_at ON user_level_history (telegram_id, reached_at DESC);
CREATE INDEX IF NOT EXISTS idx_user_level_history_xp_event_id ON user_level_history (xp_event_id);
CREATE INDEX IF NOT EXISTS idx_user_level_history_level_number ON user_level_history (level_number);