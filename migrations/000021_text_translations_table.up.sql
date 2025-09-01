CREATE TABLE IF NOT EXISTS text_translations( -- Хранит переводы значений по lang для каждого кода (content_id). Один content_id может иметь переводы на ru, en, fr, и т.д.
    id BIGSERIAL PRIMARY KEY, -- Уникальный идентификатор.
    content_id BIGINT NOT NULL, -- Ссылка на text_contents.id.
    lang VARCHAR(2) NOT NULL, -- Язык перевода (ru, en, и т.д.).
    value TEXT NOT NULL, -- Переведённый текст.
    FOREIGN KEY (content_id) REFERENCES text_contents(id),
    CONSTRAINT uq_content_lang UNIQUE (content_id, lang)
);

INSERT INTO text_translations (content_id, lang, value) VALUES
((SELECT id FROM text_contents WHERE code = 'welcome_title'), 'ru', 'Добро пожаловать!'),
((SELECT id FROM text_contents WHERE code = 'welcome_title'), 'en', 'Welcome!'),
((SELECT id FROM text_contents WHERE code = 'welcome_subtitle'), 'ru', 'Готовы улучшить ваши знания иностранных языков?'),
((SELECT id FROM text_contents WHERE code = 'welcome_subtitle'), 'en', 'Ready to improve your language skills?');