CREATE TABLE IF NOT EXISTS achievements( -- Содержит описания достижений. Хранит информацию обо всех возможных достижениях, которые могут быть присвоены пользователям.
    id SERIAL PRIMARY KEY, -- Уникальный идентификатор.
    achievement_assets_id INTEGER NOT NULL, -- Идентификатор с помощью которого можно получить url изображения для достижения.
    award_assets_id INTEGER NOT NULL, -- Идентификатор с помощью которого можно получить url изображения для награды.
    code TEXT NOT NULL UNIQUE, -- Уникальный символьный код достижения (например, streak_7, words_50).
    name TEXT NOT NULL, -- Человеко-читаемое имя достижения, которое отображается в интерфейсе пользователя (например, "7 дней", "50 слов").
    description TEXT, -- Описание или пояснение к достижению, отображаемое в интерфейсе или всплывающей подсказке.
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- Дата создания записи.
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- Дата обновления записи.
    FOREIGN KEY (achievement_assets_id) REFERENCES achievement_assets(id),
    FOREIGN KEY (award_assets_id) REFERENCES award_assets(id)
);