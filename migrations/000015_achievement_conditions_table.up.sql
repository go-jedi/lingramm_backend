CREATE TABLE IF NOT EXISTS achievement_conditions ( -- Хранит условия получения достижений.
    id SERIAL PRIMARY KEY, -- Уникальный идентификатор условия.
    achievement_id INTEGER NOT NULL, -- Какому достижению принадлежит.
    condition_type TEXT NOT NULL UNIQUE, -- Тип условия (например, 'streak_days', 'words_learned', 'level', и т.д.).
    operator TEXT NOT NULL CHECK (operator IN ('=', '>=', '<=', '>', '<')), -- Логический оператор сравнения.
    value INTEGER NOT NULL, -- Значение, с которым сравнивается прогресс пользователя.
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- Дата создания записи.
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- Дата обновления записи.
    FOREIGN KEY (achievement_id) REFERENCES achievements(id) ON DELETE CASCADE
);