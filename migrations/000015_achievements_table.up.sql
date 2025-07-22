CREATE TABLE IF NOT EXISTS achievements( -- хранит информацию обо всех возможных достижениях, которые могут быть присвоены пользователям.
    id SERIAL PRIMARY KEY, -- Уникальный идентификатор.
    asset_id INTEGER NOT NULL, -- Идентификатор с помощью которого можно получить url изображения для достижения.
    code TEXT UNIQUE NOT NULL, -- Уникальный символьный код достижения (например, streak_7, words_50).
    name TEXT NOT NULL, -- Человеко-читаемое имя достижения, которое отображается в интерфейсе пользователя (например, "7 дней", "50 слов").
    description TEXT, -- Описание или пояснение к достижению, отображаемое в интерфейсе или всплывающей подсказке.
    created_at TIMESTAMP NOT NULL DEFAULT NOW(), -- Дата создания записи.
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(), -- Дата обновления записи.
    FOREIGN KEY (asset_id) REFERENCES achievements_assets(id)
);