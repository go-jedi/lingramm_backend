CREATE TABLE IF NOT EXISTS event_types( -- справочник типов событий.
    id BIGSERIAL PRIMARY KEY, -- Уникальный идентификатор.
    name TEXT NOT NULL UNIQUE, -- Уникальное название события.
    description TEXT, -- Описание.
    xp INTEGER NOT NULL DEFAULT 0, -- XP за событие.
    amount NUMERIC(20,2) NOT NULL DEFAULT 0.00, -- Сумма бонуса.
    notification_message TEXT, -- Текст уведомления клиенту, если нужно отправлять уведомления.
    is_send_notification BOOLEAN NOT NULL DEFAULT FALSE, -- нужно ли по выполнения события отправлять уведомление.
    is_active BOOLEAN NOT NULL DEFAULT TRUE, -- Активно ли событие.
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- Дата создания записи.
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW() -- Дата обновления записи.
);

INSERT INTO event_types(
    name,
    description,
    xp
) VALUES(
    'daily_login',
    'Событие по ежедневному заходу в приложение пользователем',
    5
);

INSERT INTO event_types(
    name,
    description,
    xp,
    amount,
    notification_message,
    is_send_notification
) VALUES(
    'mini_game_reward',
    'Событие по выполнению мини игры пользователем',
    20,
    5.00,
    'Мини игра успешно выполнена',
    TRUE
);