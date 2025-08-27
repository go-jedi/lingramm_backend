CREATE TABLE IF NOT EXISTS user_achievements ( -- Хранит информацию о достижениях, полученных пользователями. Фиксирует какое достижение получил пользователь.
    id BIGSERIAL PRIMARY KEY, -- Уникальный идентификатор.
    telegram_id TEXT NOT NULL, -- Telegram id пользователя.
    achievement_id INTEGER NOT NULL, -- Идентификатор достижения.
    unlocked_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- Время, когда достижение было разблокировано пользователем.
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- Дата создания записи.
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- Дата обновления записи.
    FOREIGN KEY (telegram_id) REFERENCES users(telegram_id),
    FOREIGN KEY (achievement_id) REFERENCES achievements(id),
    CONSTRAINT unique_user_achievement UNIQUE(telegram_id, achievement_id) -- Запрещает повторное добавление одного и того же достижения для одного пользователя.
);