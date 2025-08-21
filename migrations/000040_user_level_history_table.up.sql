CREATE TABLE IF NOT EXISTS user_level_history( -- Когда достигнут конкретный уровень, чем подтверждено, и с каким XP.
    id BIGSERIAL PRIMARY KEY, -- Уникальный идентификатор.
    telegram_id TEXT NOT NULL, -- Telegram id пользователя.
    level_number BIGINT NOT NULL, -- Номер достигнутого уровня.
    xp_event_id BIGINT, -- Идентификатор XP-события, из-за которого произошёл апгрейд уровня.
    xp_at_reach BIGINT NOT NULL, -- Суммарный XP пользователя на момент достижения уровня (reached_at). Позволяет доказать, что порог из levels.required_experience действительно преодолён.
    reached_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- Момент, когда уровень фактически был достигнут.
    FOREIGN KEY (telegram_id) REFERENCES users(telegram_id),
    FOREIGN KEY (level_number) REFERENCES levels(level_number),
    FOREIGN KEY (xp_event_id) REFERENCES xp_events(id),
    CONSTRAINT unique_user_level_history_telegram_id_level_number UNIQUE (telegram_id, level_number),
    CONSTRAINT check_user_level_history_xp_at_reach_nonneg CHECK (xp_at_reach >= 0)
);