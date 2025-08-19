CREATE TABLE IF NOT EXISTS levels( -- Уровни (пороговый XP для каждого уровня).
    id BIGSERIAL PRIMARY KEY, -- Уникальный идентификатор.
    level_name VARCHAR(50) NOT NULL UNIQUE, -- Название уровня.
    level_number BIGINT NOT NULL UNIQUE, -- Числовое значение уровня.
    required_experience BIGINT NOT NULL, -- Требуемый опыт.
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- Дата создания записи.
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- Дата обновления записи.
    CONSTRAINT check_levels_required_experience_nonneg CHECK (required_experience >= 0)
);

INSERT INTO levels (level_name, level_number, required_experience) VALUES
('level 1', 1, 0),
('level 2', 2, 100),
('level 3', 3, 300),
('level 4', 4, 600),
('level 5', 5, 1000),
('level 6', 6, 1500),
('level 7', 7, 2000);