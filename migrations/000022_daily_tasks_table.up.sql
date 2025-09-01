CREATE TABLE IF NOT EXISTS daily_tasks( -- Содержит описание ежедневных заданий. Хранит информацию обо всех возможных ежедневных заданий, которые будут выполнять пользователи.
    id BIGSERIAL PRIMARY KEY, -- Уникальный идентификатор.
    words_learned_need BIGINT NOT NULL DEFAULT 0 , -- Сколько нужно выучить слов.
    tasks_completed_need BIGINT NOT NULL DEFAULT 0, -- Сколько заданий нужно выполнить.
    lessons_finished_need BIGINT NOT NULL DEFAULT 0, -- Сколько нужно пройти уроков.
    words_translate_need BIGINT NOT NULL DEFAULT 0, -- Сколько нужно перевести слов.
    dialog_completed_need BIGINT NOT NULL DEFAULT 0, -- Сколько нужно пройти диалогов.
    experience_points_need BIGINT NOT NULL DEFAULT 0, -- Сколько нужно опыта.
    is_active BOOLEAN NOT NULL DEFAULT TRUE, -- Активно ли событие.
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- Дата создания записи.
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW() -- Дата обновления записи.
);

INSERT INTO daily_tasks(
    words_translate_need,
    dialog_completed_need
) VALUES(
    6,
    1
);

INSERT INTO daily_tasks(
    dialog_completed_need,
    experience_points_need
) VALUES(
    2,
    50
);