CREATE OR REPLACE FUNCTION public.sync_user_stats_from_xp_events(_telegram_id TEXT) RETURNS VOID
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
    )
    INSERT INTO user_stats (telegram_id, experience_points, last_active_at)
    SELECT
        _telegram_id,
        total_delta_xp,
        last_seen_at
    FROM agg
    ON CONFLICT (telegram_id) DO UPDATE SET
        experience_points = EXCLUDED.experience_points,
        last_active_at = GREATEST(
            COALESCE(user_stats.last_active_at, '-infinity'::TIMESTAMP WITH TIME ZONE),
            COALESCE(EXCLUDED.last_active_at,  '-infinity'::TIMESTAMP WITH TIME ZONE)
        ),
        updated_at = NOW()
    WHERE user_stats.experience_points IS DISTINCT FROM EXCLUDED.experience_points
    OR user_stats.last_active_at IS DISTINCT FROM EXCLUDED.last_active_at;

    RETURN;
END;
$$;