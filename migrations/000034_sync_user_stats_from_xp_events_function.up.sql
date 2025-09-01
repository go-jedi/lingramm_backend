CREATE OR REPLACE FUNCTION public.sync_user_stats_from_xp_events(
    _telegram_id TEXT,
    _src JSONB
) RETURNS VOID
    SECURITY DEFINER
    LANGUAGE plpgsql
AS
$$
BEGIN
    IF _telegram_id IS NULL THEN
        RAISE EXCEPTION 'telegram_id IS NULL';
    END IF;

    WITH agg AS (
        SELECT
            COALESCE(SUM(delta_xp), 0)::BIGINT AS total_delta_xp,
            MAX(occurred_at) AS last_seen_at
        FROM xp_events
        WHERE telegram_id = _telegram_id
    ),
    actions_data AS (
        SELECT
            COALESCE(NULLIF((_src->>'words_learned')::BIGINT, NULL), 0) AS words_learned,
            COALESCE(NULLIF((_src->>'tasks_completed')::BIGINT, NULL), 0) AS tasks_completed,
            COALESCE(NULLIF((_src->>'lessons_finished')::BIGINT, NULL), 0) AS lessons_finished,
            COALESCE(NULLIF((_src->>'words_translate')::BIGINT, NULL), 0) AS words_translate,
            COALESCE(NULLIF((_src->>'dialog_completed')::BIGINT, NULL), 0) AS dialog_completed
    )
    INSERT INTO user_stats (
        telegram_id,
        words_learned,
        tasks_completed,
        lessons_finished,
        words_translate,
        dialog_completed,
        experience_points,
        last_active_at
    )
    SELECT
        _telegram_id,
        words_learned,
        tasks_completed,
        lessons_finished,
        words_translate,
        dialog_completed,
        total_delta_xp,
        last_seen_at
    FROM agg, actions_data
    ON CONFLICT (telegram_id) DO UPDATE SET
        words_learned = user_stats.words_learned + EXCLUDED.words_learned,
        tasks_completed = user_stats.tasks_completed + EXCLUDED.tasks_completed,
        lessons_finished = user_stats.lessons_finished + EXCLUDED.lessons_finished,
        words_translate = user_stats.words_translate + EXCLUDED.words_translate,
        dialog_completed = user_stats.dialog_completed + EXCLUDED.dialog_completed,
        experience_points = EXCLUDED.experience_points,
        last_active_at = GREATEST(
            COALESCE(user_stats.last_active_at, '-infinity'::TIMESTAMP WITH TIME ZONE),
            COALESCE(EXCLUDED.last_active_at,  '-infinity'::TIMESTAMP WITH TIME ZONE)
        ),
        updated_at = NOW()
    WHERE user_stats.words_learned IS DISTINCT FROM EXCLUDED.words_learned
    OR user_stats.tasks_completed IS DISTINCT FROM EXCLUDED.tasks_completed
    OR user_stats.lessons_finished IS DISTINCT FROM EXCLUDED.lessons_finished
    OR user_stats.words_translate IS DISTINCT FROM EXCLUDED.words_translate
    OR user_stats.dialog_completed IS DISTINCT FROM EXCLUDED.dialog_completed
    OR user_stats.experience_points IS DISTINCT FROM EXCLUDED.experience_points
    OR user_stats.last_active_at IS DISTINCT FROM EXCLUDED.last_active_at;

    RETURN;
END;
$$;