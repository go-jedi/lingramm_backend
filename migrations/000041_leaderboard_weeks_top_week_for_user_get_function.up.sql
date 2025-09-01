CREATE OR REPLACE FUNCTION public.leaderboard_weeks_top_week_for_user_get(
    _telegram_id TEXT,
    _limit INTEGER,
    _tz TEXT
) RETURNS JSONB
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
    IF _limit IS NULL THEN
        RAISE EXCEPTION 'limit IS NULL';
    END IF;
    IF _tz IS NULL THEN
        RAISE EXCEPTION 'tz IS NULL';
    END IF;

    WITH params AS (
        SELECT
            DATE_TRUNC('week', (NOW() AT TIME ZONE _tz))::DATE AS ws,
            _limit::INTEGER AS lim,
            _telegram_id::TEXT AS telegram_id
    ),
    ranked AS (
        SELECT
            DENSE_RANK() OVER (
                ORDER BY lbw.xp DESC, lbw.telegram_id
            ) AS position,
            lbw.telegram_id,
            COALESCE(
                NULLIF(u.username, ''),
                NULLIF(CONCAT_WS(' ', u.first_name, u.last_name), '')
            ) AS display_name,
            lbw.xp
        FROM leaderboard_weeks lbw
        INNER JOIN params p ON lbw.week_start = p.ws
        LEFT JOIN users u ON lbw.telegram_id = u.telegram_id
    ),
    topn AS (
        SELECT *
        FROM ranked
        ORDER BY position
        LIMIT (
            SELECT
                lim
            FROM params
        )
    ),
    me_row AS (
        SELECT r.*
        FROM ranked r
        INNER JOIN params p ON r.telegram_id = p.telegram_id
    ),
    unioned AS (
        SELECT
            t.position,
            t.telegram_id,
            t.display_name,
            t.xp,
            FALSE AS is_me,
            0 AS ord
        FROM topn t
        UNION ALL
        SELECT
            mr.position,
            mr.telegram_id,
            mr.display_name,
            mr.xp,
            TRUE AS is_me,
            1 AS ord
        FROM me_row mr
        WHERE NOT EXISTS (
            SELECT 1
            FROM topn t
            WHERE t.telegram_id = mr.telegram_id
        )
    )
    SELECT COALESCE(
        JSONB_AGG(
            JSONB_BUILD_OBJECT(
                'position', position,
                'telegram_id', telegram_id,
                'display_name', display_name,
                'xp', xp,
                'medal',
                CASE position
                    WHEN 1 THEN 'gold'
                    WHEN 2 THEN 'silver'
                    WHEN 3 THEN 'bronze'
                    END
            )
            ORDER BY position, ord
        ),
        '[]'::JSONB
    )
    INTO _response
    FROM unioned;

    RETURN _response;
END;
$$;