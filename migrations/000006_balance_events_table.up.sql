CREATE TABLE IF NOT EXISTS balance_events( -- События по которым начисляется бонус.
    id SERIAL PRIMARY KEY, -- Уникальный идентификатор.
    code VARCHAR(50) UNIQUE NOT NULL, -- 'daily_login', 'game_completed', 'holiday_bonus'
    description TEXT, -- Описание
    amount NUMERIC(20,2) NOT NULL, -- Сумма бонуса.
    is_active BOOLEAN NOT NULL DEFAULT TRUE, -- Активно ли событие.
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- Дата создания записи.
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW() -- Дата обновления записи.
);