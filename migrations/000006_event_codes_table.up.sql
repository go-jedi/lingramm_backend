CREATE TABLE IF NOT EXISTS event_codes( -- справочник кодов событий.
    id BIGSERIAL PRIMARY KEY, -- Уникальный идентификатор.
    name TEXT NOT NULL UNIQUE, -- Уникальное название кода.
    description TEXT, -- Описание.
    is_active BOOLEAN NOT NULL DEFAULT TRUE, -- Активно ли событие.
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- Дата создания записи.
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW() -- Дата обновления записи.
);