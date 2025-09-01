CREATE OR REPLACE FUNCTION public.unlock_available_achievements(_telegram_id TEXT) RETURNS JSONB
    SECURITY DEFINER
    LANGUAGE plpgsql
AS
$$
DECLARE
    _response JSONB;
BEGIN
    IF _telegram_id IS NULL THEN
        RAISE EXCEPTION 'telegram_id IS NULL';
    END IF;

    WITH metrics AS (
        SELECT
            us.streak_days,
            us.words_learned,
            us.tasks_completed,
            us.lessons_finished,
            us.experience_points,
            us.level,
            us.daily_task_streak_days,
            us.words_translate,
            us.dialog_completed
        FROM user_stats us
        WHERE us.telegram_id = _telegram_id
    ),
    eligible AS (
        SELECT a.id, a.name
        FROM achievements a
        INNER JOIN achievement_types at ON a.achievement_type_id = at.id
        INNER JOIN metrics m ON TRUE
        WHERE at.is_active
        AND (
            at.streak_days_need IS NULL
            OR m.streak_days >= at.streak_days_need
        )
        AND (
            at.daily_task_streak_days_need IS NULL
            OR m.daily_task_streak_days >= at.daily_task_streak_days_need
        )
        AND (
            at.words_learned_need IS NULL
            OR m.words_learned >= at.words_learned_need
        )
        AND (
            at.tasks_completed_need IS NULL
            OR m.tasks_completed >= at.tasks_completed_need
        )
        AND (
            at.lessons_finished_need IS NULL
            OR m.lessons_finished >= at.lessons_finished_need
        )
        AND (
            at.words_translate_need IS NULL
            OR m.words_translate >= at.words_translate_need
        )
        AND (
            at.dialog_completed_need IS NULL
            OR m.dialog_completed >= at.dialog_completed_need
        )
        AND (
            at.experience_points_need IS NULL
            OR m.experience_points >= at.experience_points_need
        )
        AND (
            at.level_need IS NULL
            OR m.level >= at.level_need
        )
    ),
    inserted AS (
        INSERT INTO user_achievements(
            telegram_id,
            achievement_id,
            unlocked_at
        )
        SELECT _telegram_id, e.id, NOW()
        FROM eligible e
        ON CONFLICT (telegram_id, achievement_id) DO NOTHING
        RETURNING achievement_id, unlocked_at
    )
    SELECT COALESCE(
        JSONB_AGG(
            JSONB_BUILD_OBJECT(
                'achievement_id', i.achievement_id,
                'achievement_name', a.name,
                'unlocked_at', i.unlocked_at
            )
        ),
        '[]'::JSONB
    )
    INTO _response
    FROM inserted i
    INNER JOIN achievements a ON i.achievement_id = a.id;

    RETURN _response;
END;
$$;