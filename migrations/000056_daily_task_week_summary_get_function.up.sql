CREATE OR REPLACE FUNCTION public.daily_task_week_summary_get(
    _telegram_id TEXT
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

    WITH bounds AS (
        SELECT (
            DATE_TRUNC('week', (NOW() AT TIME ZONE 'Europe/Moscow'))::DATE
        ) AS week_start
    ),
    days AS (
        SELECT generate_series(b.week_start, b.week_start + 6, INTERVAL '1 day')::DATE AS d
        FROM bounds b
    ),
    per_day AS (
        SELECT
            d.d,
            EXISTS (
                SELECT 1
                FROM user_daily_tasks udt
                WHERE udt.telegram_id = _telegram_id
                AND (
                    (udt.occurred_at AT TIME ZONE 'Europe/Moscow')::DATE
                ) = d.d
                AND udt.is_completed = TRUE
            ) AS is_completed
        FROM days d
    )
    SELECT
        JSONB_AGG(
            JSONB_BUILD_OBJECT(
                'date', TO_CHAR(pd.d, 'YYYY-MM-DD"T"00:00:00"Z"'),
                'is_completed', pd.is_completed
            )
            ORDER BY pd.d
        )
    INTO _response
    FROM per_day pd;

    RETURN _response;
END;
$$;