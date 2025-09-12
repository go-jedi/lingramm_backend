CREATE TABLE IF NOT EXISTS user_studied_languages( -- язык, который выбрал пользователь для изучения.
    id BIGSERIAL PRIMARY KEY, -- Уникальный идентификатор.
    telegram_id TEXT NOT NULL UNIQUE, -- Telegram id пользователя.
    studied_language_id BIGINT NOT NULL, -- Идентификатор языка из таблицы studied_languages.
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- Дата создания записи.
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- Дата обновления записи.
    FOREIGN KEY (telegram_id) REFERENCES users(telegram_id),
    FOREIGN KEY (studied_language_id) REFERENCES studied_languages(id)
);