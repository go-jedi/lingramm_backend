CREATE TABLE IF NOT EXISTS events( -- события доступные в приложении.
    id BIGSERIAL PRIMARY KEY, -- Уникальный идентификатор.
    name VARCHAR(50) UNIQUE NOT NULL, -- 'daily_login', 'game_reward', 'holiday_bonus'.
    description TEXT NOT NULL, -- Описание.
    amount NUMERIC(20,2) NOT NULL DEFAULT 0, -- Сумма бонуса.
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- Дата создания записи.
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW() -- Дата обновления записи.
);