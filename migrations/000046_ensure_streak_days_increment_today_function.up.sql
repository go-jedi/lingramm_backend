CREATE OR REPLACE FUNCTION public.ensure_streak_days_increment_today(_telegram_id TEXT) RETURNS VOID
    SECURITY DEFINER
    LANGUAGE plpgsql
AS
$$
BEGIN
    IF _telegram_id IS NULL THEN
        RAISE EXCEPTION 'telegram_id IS NULL';
    END IF;

    WITH params AS (
        -- получаем параметры для будущего использования.
        SELECT
            _telegram_id::TEXT AS telegram_id,
            CURRENT_DATE AS today,
            NOW() AS ts
    )
    -- если вчера был учтен, то +1; если был разрыв по дате, то 1;
    -- если уже был учтен сегодня, то без изменений.
    UPDATE user_stats us SET
        streak_days =
            CASE
                WHEN us.last_streak_day = p.today - 1 THEN
                    us.streak_days + 1
                ELSE 1
            END,
        last_streak_day = p.today,
        last_active_at = p.ts,
        updated_at = now()
    FROM params p
    WHERE us.telegram_id = p.telegram_id
      -- обновляем только если наступил новый день ИЛИ
      -- это первый реальный учёт (streak=0, day уже = today из-за DEFAULT).
    AND (
        p.today > us.last_streak_day OR (
            p.today = us.last_streak_day AND us.streak_days = 0
        )
    );
END;
$$;