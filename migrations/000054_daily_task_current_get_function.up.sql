CREATE OR REPLACE FUNCTION public.daily_task_current_get(_telegram_id TEXT) RETURNS JSONB
    SECURITY DEFINER
    LANGUAGE plpgsql
AS
$$
DECLARE
    _today_msk DATE := (NOW() AT TIME ZONE 'Europe/Moscow')::DATE;
    _response JSONB;
BEGIN
    IF _telegram_id IS NULL THEN
        RAISE EXCEPTION 'telegram_id IS NULL';
    END IF;

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
        'progress_percent', JSONB_STRIP_NULLS(
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