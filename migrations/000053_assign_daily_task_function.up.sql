CREATE OR REPLACE FUNCTION public.assign_daily_task(_telegram_id TEXT) RETURNS JSONB
    SECURITY DEFINER
    LANGUAGE plpgsql
AS
$$
DECLARE
    _today_msk DATE := (NOW() AT TIME ZONE 'Europe/Moscow')::DATE;
    _udt_today_id BIGINT;
    _picked_task_id BIGINT;
    _prev_id BIGINT;
    _prev_date DATE;
    _prev_completed BOOLEAN;
    _words_learned_need BIGINT;
    _tasks_completed_need BIGINT;
    _lessons_finished_need BIGINT;
    _words_translate_need BIGINT;
    _dialog_completed_need BIGINT;
    _experience_points_need BIGINT;
    _last_daily_task_streak_days DATE;
    _daily_task_streak_days BIGINT;
    _completed_now BOOLEAN;
    _response JSONB;
BEGIN
    IF _telegram_id IS NULL THEN
        RAISE EXCEPTION 'telegram_id IS NULL';
    END IF;

    PERFORM PG_ADVISORY_XACT_LOCK(HASHTEXT(_telegram_id));

    -- блокируем строку в таблице user_stats.
    PERFORM 1
    FROM user_stats
    WHERE telegram_id = _telegram_id
    FOR UPDATE;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'user_stats row is missing and cannot be created (no users row?) for %', _telegram_id;
    END IF;

    -- если на сегодня уже есть назначение ежедневное название, то возвращаем.
    SELECT udt.id
    INTO _udt_today_id
    FROM user_daily_tasks udt
    WHERE udt.telegram_id = _telegram_id
    AND (
        (udt.occurred_at AT TIME ZONE 'Europe/Moscow')::DATE
    ) = _today_msk
    ORDER BY udt.occurred_at DESC
    LIMIT 1;

    IF _udt_today_id IS NOT NULL THEN
        SELECT JSONB_BUILD_OBJECT(
            'id', udt.id,
            'date', TO_CHAR(_today_msk, 'YYYY-MM-DD"T"00:00:00"Z"'),
            'is_completed', udt.is_completed,
            'requirements', JSONB_STRIP_NULLS(
                JSONB_BUILD_OBJECT(
                    'words_learned_need',
                        CASE WHEN dt.words_learned_need > 0
                            THEN dt.words_learned_need
                        END,
                    'tasks_completed_need',
                        CASE WHEN dt.tasks_completed_need > 0
                            THEN dt.tasks_completed_need
                        END,
                    'lessons_finished_need',
                        CASE WHEN dt.lessons_finished_need > 0
                            THEN dt.lessons_finished_need
                        END,
                    'words_translate_need',
                        CASE WHEN dt.words_translate_need > 0
                            THEN dt.words_translate_need
                        END,
                    'dialog_completed_need',
                        CASE WHEN dt.dialog_completed_need > 0
                            THEN dt.dialog_completed_need
                        END,
                    'experience_points_need',
                        CASE WHEN dt.experience_points_need > 0
                            THEN dt.experience_points_need
                        END
                )
            ),
            'progress',
            JSONB_STRIP_NULLS(
                JSONB_BUILD_OBJECT(
                    'words_learned',
                        CASE WHEN dt.words_learned_need > 0
                            THEN udt.words_learned
                        END,
                    'tasks_completed',
                        CASE WHEN dt.tasks_completed_need > 0
                            THEN udt.tasks_completed
                        END,
                    'lessons_finished',
                        CASE WHEN dt.lessons_finished_need > 0
                            THEN udt.lessons_finished
                        END,
                    'words_translate',
                        CASE WHEN dt.words_translate_need > 0
                            THEN udt.words_translate
                        END,
                    'dialog_completed',
                        CASE WHEN dt.dialog_completed_need > 0
                            THEN udt.dialog_completed
                        END,
                    'experience_points',
                        CASE WHEN dt.experience_points_need > 0
                            THEN udt.experience_points
                        END
                )
            ),
            'progress_percent',
            JSONB_STRIP_NULLS(
                JSONB_BUILD_OBJECT(
                    'words_learned',
                        CASE WHEN dt.words_learned_need > 0
                            THEN LEAST(ROUND((udt.words_learned::NUMERIC / NULLIF(dt.words_learned_need,0)) * 100), 100)::INTEGER
                        END,
                    'tasks_completed',
                        CASE WHEN dt.tasks_completed_need > 0
                            THEN LEAST(ROUND((udt.tasks_completed::NUMERIC / NULLIF(dt.tasks_completed_need,0)) * 100), 100)::INTEGER
                        END,
                    'lessons_finished',
                        CASE WHEN dt.lessons_finished_need > 0
                            THEN LEAST(ROUND((udt.lessons_finished::NUMERIC / NULLIF(dt.lessons_finished_need,0)) * 100), 100)::INTEGER
                        END,
                    'words_translate',
                        CASE WHEN dt.words_translate_need > 0
                            THEN LEAST(ROUND((udt.words_translate::NUMERIC / NULLIF(dt.words_translate_need,0)) * 100), 100)::INTEGER
                        END,
                    'dialog_completed',
                        CASE WHEN dt.dialog_completed_need > 0
                            THEN LEAST(ROUND((udt.dialog_completed::NUMERIC / NULLIF(dt.dialog_completed_need,0)) * 100), 100)::INTEGER
                        END,
                    'experience_points',
                        CASE WHEN dt.experience_points_need > 0
                            THEN LEAST(ROUND((udt.experience_points::NUMERIC / NULLIF(dt.experience_points_need,0)) * 100), 100)::INTEGER
                        END
                )
            )
        )
        INTO _response
        FROM user_daily_tasks udt
        INNER JOIN daily_tasks dt ON udt.daily_task_id = dt.id
        WHERE udt.id = _udt_today_id;

        RETURN _response;
    END IF;

    -- вердикт по вчерашнему/последнему ежедневному заданию.
    SELECT
        udt.id,
        (
            (udt.occurred_at AT TIME ZONE 'Europe/Moscow')::DATE
        ) AS d,
        udt.is_completed,
        dt.words_learned_need,
        dt.tasks_completed_need,
        dt.lessons_finished_need,
        dt.words_translate_need,
        dt.dialog_completed_need,
        dt.experience_points_need
    INTO
        _prev_id,
        _prev_date,
        _prev_completed,
        _words_learned_need,
        _tasks_completed_need,
        _lessons_finished_need,
        _words_translate_need,
        _dialog_completed_need,
        _experience_points_need
    FROM user_daily_tasks udt
    INNER JOIN daily_tasks dt ON udt.daily_task_id = dt.id
    WHERE udt.telegram_id = _telegram_id
    AND (
        (udt.occurred_at AT TIME ZONE 'Europe/Moscow')::DATE
    ) < _today_msk
    ORDER BY udt.occurred_at DESC
    LIMIT 1
    FOR UPDATE;

    IF _prev_id IS NOT NULL THEN
        -- если is_completed ещё не выставлен корректно — вычислим по факту прогресса vs требований.
        IF _prev_completed IS DISTINCT FROM TRUE THEN
            SELECT
                (
                    (_words_learned_need = 0 OR udt.words_learned >= _words_learned_need)
                    AND (_tasks_completed_need = 0 OR udt.tasks_completed >= _tasks_completed_need)
                    AND (_lessons_finished_need = 0 OR udt.lessons_finished >= _lessons_finished_need)
                    AND (_words_translate_need = 0 OR udt.words_translate >= _words_translate_need)
                    AND (_dialog_completed_need = 0 OR udt.dialog_completed >= _dialog_completed_need)
                    AND (_experience_points_need = 0 OR udt.experience_points >= _experience_points_need)
                )
            INTO _completed_now
            FROM user_daily_tasks udt
            WHERE udt.id = _prev_id
            FOR UPDATE;

            UPDATE user_daily_tasks SET
                is_completed = _completed_now
            WHERE id = _prev_id;

            _prev_completed := _completed_now;
        END IF;

        -- обновление streak, если вчера/последний день действительно выполнен.
        IF _prev_completed THEN
            SELECT
                last_daily_task_streak_days,
                daily_task_streak_days
            INTO
                _last_daily_task_streak_days,
                _daily_task_streak_days
            FROM user_stats
            WHERE telegram_id = _telegram_id
            FOR UPDATE;

            IF _last_daily_task_streak_days = _prev_date - 1 THEN
                UPDATE user_stats SET
                    daily_task_streak_days = _daily_task_streak_days + 1,
                    last_daily_task_streak_days = _prev_date,
                    updated_at = NOW()
                WHERE telegram_id = _telegram_id;
            ELSE
                UPDATE user_stats SET
                    daily_task_streak_days = 1,
                    last_daily_task_streak_days = _prev_date,
                    updated_at = NOW()
                WHERE telegram_id = _telegram_id;
            END IF;
        END IF;
    END IF;

    -- Найти ежедневное задание на сегодня с анти-повтором 4 дня (по MSK-дате).
    WITH recent AS (
        SELECT
            DISTINCT daily_task_id
        FROM user_daily_tasks
        WHERE telegram_id = _telegram_id
        AND (
            (occurred_at AT TIME ZONE 'Europe/Moscow')::DATE
        ) >= _today_msk - 4
    ),
    candidates AS (
        SELECT
            dt.id
        FROM daily_tasks dt
        WHERE dt.is_active
        AND NOT EXISTS (
            SELECT 1 FROM recent r WHERE r.daily_task_id = dt.id
        )
    ),
    numbered AS (
        SELECT
            id,
            ROW_NUMBER() OVER (ORDER BY id) rn,
            COUNT(*) OVER() cnt
        FROM candidates
    ),
    pick AS (
        SELECT
            n.id
        FROM numbered n
        WHERE n.rn = 1 + FLOOR(RANDOM() * GREATEST(n.cnt,1))::INTEGER
        LIMIT 1
    )
    SELECT
        id
    INTO _picked_task_id
    FROM pick;

    -- если ежедневных заданий нет (все были за последние 4 дня) — разрешаем любые активные.
    IF _picked_task_id IS NULL THEN
        WITH all_active AS (
            SELECT
                id,
                ROW_NUMBER() OVER (ORDER BY id) rn,
                COUNT(*) OVER() cnt
            FROM daily_tasks
            WHERE is_active
        ),
        pick2 AS (
            SELECT
                a.id
            FROM all_active a
            WHERE a.rn = 1 + FLOOR(RANDOM() * GREATEST(a.cnt,1))::INTEGER
            LIMIT 1
        )
        SELECT
            id
        INTO _picked_task_id
        FROM pick2;
    END IF;

    INSERT INTO user_daily_tasks(
        daily_task_id,
        telegram_id,
        occurred_at
    )
    VALUES (
        _picked_task_id,
        _telegram_id,
        NOW()
    );

    SELECT JSONB_BUILD_OBJECT(
        'id', udt.id,
        'date', TO_CHAR(_today_msk, 'YYYY-MM-DD"T"00:00:00"Z"'),
        'is_completed', udt.is_completed,
        'requirements', JSONB_STRIP_NULLS(
            JSONB_BUILD_OBJECT(
                'words_learned_need',
                    CASE WHEN dt.words_learned_need > 0
                        THEN dt.words_learned_need
                    END,
                'tasks_completed_need',
                    CASE WHEN dt.tasks_completed_need > 0
                        THEN dt.tasks_completed_need
                    END,
                'lessons_finished_need',
                    CASE WHEN dt.lessons_finished_need > 0
                        THEN dt.lessons_finished_need
                    END,
                'words_translate_need',
                    CASE WHEN dt.words_translate_need > 0
                        THEN dt.words_translate_need
                    END,
                'dialog_completed_need',
                    CASE WHEN dt.dialog_completed_need > 0
                        THEN dt.dialog_completed_need
                    END,
                'experience_points_need',
                    CASE WHEN dt.experience_points_need > 0
                        THEN dt.experience_points_need
                    END
            )
        ),
        'progress',
        JSONB_STRIP_NULLS(
            JSONB_BUILD_OBJECT(
                'words_learned',
                    CASE WHEN dt.words_learned_need > 0
                        THEN udt.words_learned
                    END,
                'tasks_completed',
                    CASE WHEN dt.tasks_completed_need > 0
                        THEN udt.tasks_completed
                    END,
                'lessons_finished',
                    CASE WHEN dt.lessons_finished_need > 0
                        THEN udt.lessons_finished
                    END,
                'words_translate',
                    CASE WHEN dt.words_translate_need > 0
                        THEN udt.words_translate
                    END,
                'dialog_completed',
                    CASE WHEN dt.dialog_completed_need > 0
                        THEN udt.dialog_completed
                    END,
                'experience_points',
                    CASE WHEN dt.experience_points_need > 0
                        THEN udt.experience_points
                    END
                )
        ),
        'progress_percent',
        JSONB_STRIP_NULLS(
            JSONB_BUILD_OBJECT(
                'words_learned',
                    CASE WHEN dt.words_learned_need > 0
                        THEN LEAST(ROUND((udt.words_learned::NUMERIC / NULLIF(dt.words_learned_need,0)) * 100), 100)::INTEGER
                    END,
                'tasks_completed',
                    CASE WHEN dt.tasks_completed_need > 0
                        THEN LEAST(ROUND((udt.tasks_completed::NUMERIC / NULLIF(dt.tasks_completed_need,0)) * 100), 100)::INTEGER
                    END,
                'lessons_finished',
                    CASE WHEN dt.lessons_finished_need > 0
                        THEN LEAST(ROUND((udt.lessons_finished::NUMERIC / NULLIF(dt.lessons_finished_need,0)) * 100), 100)::INTEGER
                    END,
                'words_translate',
                    CASE WHEN dt.words_translate_need > 0
                        THEN LEAST(ROUND((udt.words_translate::NUMERIC / NULLIF(dt.words_translate_need,0)) * 100), 100)::INTEGER
                    END,
                'dialog_completed',
                    CASE WHEN dt.dialog_completed_need > 0
                        THEN LEAST(ROUND((udt.dialog_completed::NUMERIC / NULLIF(dt.dialog_completed_need,0)) * 100), 100)::INTEGER
                    END,
                'experience_points',
                    CASE WHEN dt.experience_points_need > 0
                        THEN LEAST(ROUND((udt.experience_points::NUMERIC / NULLIF(dt.experience_points_need,0)) * 100), 100)::INTEGER
                    END
            )
        )
    )
    INTO _response
    FROM user_daily_tasks udt
    INNER JOIN daily_tasks dt ON udt.daily_task_id = dt.id
    WHERE udt.telegram_id = _telegram_id
    AND (
        (udt.occurred_at AT TIME ZONE 'Europe/Moscow')::DATE
    ) = _today_msk
    ORDER BY udt.occurred_at DESC
    LIMIT 1;

    RETURN _response;
END;
$$;