CREATE OR REPLACE FUNCTION public.sync_user_daily_task_progress(
    _telegram_id TEXT,
    _src JSONB
) RETURNS VOID
    SECURITY DEFINER
    LANGUAGE plpgsql
AS
$$
DECLARE
    _today_msk DATE := (NOW() AT TIME ZONE 'Europe/Moscow')::DATE;
    _udt_id BIGINT;
    _words_learned_need BIGINT;
    _tasks_completed_need BIGINT;
    _lessons_finished_need BIGINT;
    _words_translate_need BIGINT;
    _dialog_completed_need BIGINT;
    _experience_points_need BIGINT;
    _done BOOLEAN;
BEGIN
    IF _telegram_id IS NULL THEN
        RAISE EXCEPTION 'telegram_id IS NULL';
    END IF;

    SELECT
        udt.id
    INTO _udt_id
    FROM user_daily_tasks udt
    WHERE udt.telegram_id = _telegram_id
    AND (
        (udt.occurred_at AT TIME ZONE 'Europe/Moscow')::DATE
    ) = _today_msk
    ORDER BY udt.occurred_at DESC
    LIMIT 1
    FOR UPDATE;

    UPDATE user_daily_tasks
    SET
        words_learned = GREATEST(0, words_learned + COALESCE((_src->>'words_learned')::BIGINT, 0)),
        tasks_completed = GREATEST(0, tasks_completed + COALESCE((_src->>'tasks_completed')::BIGINT, 0)),
        lessons_finished = GREATEST(0, lessons_finished + COALESCE((_src->>'lessons_finished')::BIGINT, 0)),
        words_translate = GREATEST(0, words_translate + COALESCE((_src->>'words_translate')::BIGINT, 0)),
        dialog_completed = GREATEST(0, dialog_completed + COALESCE((_src->>'dialog_completed')::BIGINT, 0)),
        experience_points = GREATEST(0, experience_points + COALESCE((_src->>'experience_points')::BIGINT, 0))
    WHERE id = _udt_id;

    SELECT
        dt.words_learned_need,
        dt.tasks_completed_need,
        dt.lessons_finished_need,
        dt.words_translate_need,
        dt.dialog_completed_need,
        dt.experience_points_need
    INTO
        _words_learned_need,
        _tasks_completed_need,
        _lessons_finished_need,
        _words_translate_need,
        _dialog_completed_need,
        _experience_points_need
    FROM user_daily_tasks udt
    INNER JOIN daily_tasks dt ON udt.daily_task_id = dt.id
    WHERE udt.id = _udt_id;

    SELECT
        (
            (_words_learned_need = 0 OR udt.words_learned >= _words_learned_need)
            AND (_tasks_completed_need = 0 OR udt.tasks_completed >= _tasks_completed_need)
            AND (_lessons_finished_need = 0 OR udt.lessons_finished >= _lessons_finished_need)
            AND (_words_translate_need = 0 OR udt.words_translate >= _words_translate_need)
            AND (_dialog_completed_need = 0 OR udt.dialog_completed >= _dialog_completed_need)
            AND (_experience_points_need = 0 OR udt.experience_points >= _experience_points_need)
        )
    INTO _done
    FROM user_daily_tasks udt
    WHERE udt.id = _udt_id;

    UPDATE user_daily_tasks
    SET is_completed = _done
    WHERE id = _udt_id;

    RETURN;
END;
$$;