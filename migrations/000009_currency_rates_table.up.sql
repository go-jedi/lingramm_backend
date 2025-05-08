CREATE TABLE IF NOT EXISTS currency_rates(
    id SERIAL PRIMARY KEY, -- Уникальный идентификатор.
    currency_code VARCHAR(10) NOT NULL UNIQUE, -- например, 'USD'.
    rate NUMERIC(20,6) NOT NULL, -- курс: сколько USD за 1 внутреннюю валюту.
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- Дата создания записи.
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW() -- Дата обновления записи.
);