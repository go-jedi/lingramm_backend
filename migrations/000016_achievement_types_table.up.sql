CREATE TABLE IF NOT EXISTS achievement_types( -- справочник типов событий для достижений.
    id BIGSERIAL PRIMARY KEY, -- Уникальный идентификатор.
    name TEXT NOT NULL UNIQUE, -- Уникальное название события.
    description TEXT, -- Описание.
    streak_days_need BIGINT, -- Сколько дней подряд пользователь должен заходить в приложение.
    daily_task_streak_days_need BIGINT, -- Сколько дней подряд пользователь должен выполнять ежедневные задания.
    words_learned_need BIGINT, -- Сколько нужно выучить слов.
    tasks_completed_need BIGINT, -- Сколько заданий нужно выполнить.
    lessons_finished_need BIGINT, -- Сколько нужно пройти уроков.
    words_translate_need BIGINT, -- Сколько нужно перевести слов.
    dialog_completed_need BIGINT, -- Сколько нужно пройти диалогов.
    experience_points_need BIGINT, -- Сколько нужно опыта.
    level_need BIGINT, -- Какой уровень должен быть у пользователя.
    is_active BOOLEAN NOT NULL DEFAULT TRUE, -- Активно ли событие.
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- Дата создания записи.
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW() -- Дата обновления записи.
);

INSERT INTO achievement_types(
    name,
    description,
    streak_days_need
) VALUES(
    'streak_days_30',
    'Событие по ежедневному заходу в приложение пользователем на протяжении 30 дней подряд',
    30
);

INSERT INTO achievement_types(
    name,
    description,
    words_learned_need
) VALUES(
    'words_learned_50',
    'Событие по изучению пользователем 50 слов',
    50
);

INSERT INTO achievement_types(
    name,
    description,
    dialog_completed_need
) VALUES(
    'dialog_completed_2',
    'Событие по изучению 2 диалогов',
    2
);

INSERT INTO achievement_types(
    name,
    description,
    words_learned_need,
    dialog_completed_need
) VALUES(
    'dialog_completed_4',
    'Событие по изучению 4 диалогов',
    2,
    4
);