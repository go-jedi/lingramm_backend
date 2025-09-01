CREATE TABLE IF NOT EXISTS studied_languages( -- список языков, которые пользователь может изучать.
    id BIGSERIAL PRIMARY KEY, -- Уникальный идентификатор.
    name TEXT NOT NULL, -- 'English', 'Spain'.
    description TEXT NOT NULL, -- Описание.
    lang VARCHAR(2) NOT NULL UNIQUE, -- Язык перевода (ru, en, и т.д.).
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- Дата создания записи.
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW() -- Дата обновления записи.
);

INSERT INTO studied_languages(
    name,
    description,
    lang
) VALUES(
    'English',
    'English course (A1–B2): grammar, high-frequency vocabulary, pronunciation, reading, listening, and writing, with plenty of dialogue practice for travel, study, and work.',
    'en'
);

INSERT INTO studied_languages(
    name,
    description,
    lang
) VALUES(
    'Spain',
    'Curso de español (A1–B2): gramática, vocabulario de uso frecuente, pronunciación, lectura, escucha y escritura, con mucha práctica en diálogos para viajar y comunicarte.',
    'es'
);