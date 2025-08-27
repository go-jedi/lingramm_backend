CREATE TABLE IF NOT EXISTS achievement_types( -- справочник типов событий для достижений.
    id BIGSERIAL PRIMARY KEY, -- Уникальный идентификатор.
    name TEXT NOT NULL UNIQUE, -- Уникальное название события.
    description TEXT, -- Описание.
    is_active BOOLEAN NOT NULL DEFAULT TRUE, -- Активно ли событие.
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- Дата создания записи.
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW() -- Дата обновления записи.
);

INSERT INTO achievement_types(
    name,
    description
) VALUES(
    'streak_days_30',
    'Событие по ежедневному заходу в приложение пользователем на протяжении 30 дней подряд'
);

INSERT INTO achievement_types(
    name,
    description
) VALUES(
    'words_50',
    'Событие по изучению пользователем 50 слов'
);