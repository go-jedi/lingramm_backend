CREATE TABLE IF NOT EXISTS text_contents( -- Хранит уникальные коды текстов, их принадлежность к странице (page) и общую описательную информацию (обычно на основном языке, например русском).
    id BIGSERIAL PRIMARY KEY, -- Уникальный идентификатор.
    code TEXT NOT NULL UNIQUE, -- Уникальный код ключа (например menu_title).
    page TEXT NOT NULL, -- Название группы/страницы (Menu, Chat).
    description TEXT -- Общий текст/описание.
);

INSERT INTO text_contents (code, page, description) VALUES
('welcome_title', 'welcome', 'Добро пожаловать!'),
('welcome_subtitle', 'welcome', 'Готовы улучшить ваши знания иностранных языков?');