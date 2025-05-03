CREATE TABLE IF NOT EXISTS user_profiles(
    id SERIAL PRIMARY KEY, -- Уникальный идентификатор.
    uuid TEXT NOT NULL UNIQUE, -- UUID пользователя.
    telegram_id TEXT NOT NULL UNIQUE, -- Telegram id пользователя.
    experience_points BIGINT NOT NULL DEFAULT 0, -- Шкала опыта.
    created_at TIMESTAMP NOT NULL DEFAULT NOW(), -- Дата создания записи.
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(), -- Дата обновления записи.
    FOREIGN KEY (uuid) REFERENCES users(uuid) ON DELETE CASCADE,
    FOREIGN KEY (telegram_id) REFERENCES users(telegram_id) ON DELETE CASCADE
);