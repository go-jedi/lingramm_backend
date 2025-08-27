CREATE TABLE IF NOT EXISTS achievement_conditions ( -- Хранит условия получения достижений.
    id BIGSERIAL PRIMARY KEY, -- Уникальный идентификатор условия.
    achievement_id INTEGER NOT NULL, -- Какому достижению принадлежит.
    achievement_type_id BIGINT NOT NULL, -- Идентификатор события.
    operator TEXT NOT NULL CHECK (operator IN ('=', '>=', '<=', '>', '<')), -- Логический оператор сравнения.
    value INTEGER NOT NULL, -- Значение, с которым сравнивается прогресс пользователя.
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- Дата создания записи.
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- Дата обновления записи.
    FOREIGN KEY (achievement_id) REFERENCES achievements(id),
    FOREIGN KEY (achievement_type_id) REFERENCES achievement_types(id)
);