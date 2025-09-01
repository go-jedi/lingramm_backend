CREATE OR REPLACE FUNCTION public.back_fill_missing_level_history(_telegram_id TEXT) RETURNS JSONB
    SECURITY DEFINER
    LANGUAGE plpgsql
AS
$$
DECLARE
    _old_level BIGINT;
    _new_level BIGINT;
    _top_level BIGINT := 0;
    _max_cum_xp BIGINT := 0;
BEGIN
    IF _telegram_id IS NULL THEN
        RAISE EXCEPTION 'telegram_id IS NULL';
    END IF;

    -- Максимальный накопленный XP.
    WITH ordered AS (
        SELECT
            e.id,
            e.occurred_at,
            e.delta_xp,
            SUM(e.delta_xp) OVER (
                ORDER BY e.occurred_at, e.id
                ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
            ) AS cum_xp
        FROM xp_events e
        WHERE e.telegram_id = _telegram_id
    )
    SELECT COALESCE(MAX(cum_xp), 0)
    INTO _max_cum_xp
    FROM ordered;

    -- Верхний достигнутый уровень.
    SELECT l.level_number
    INTO _top_level
    FROM levels l
    WHERE l.required_experience <= _max_cum_xp
    ORDER BY l.required_experience DESC
    LIMIT 1;

    WITH ordered AS (
        -- События пользователя в строгом порядке + накопительный XP.
        SELECT
            e.id,
            e.occurred_at,
            e.delta_xp,
            SUM(e.delta_xp) OVER (
                ORDER BY e.occurred_at, e.id
                ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
            ) AS cum_xp
        FROM xp_events e
        WHERE e.telegram_id = _telegram_id
    ),
    first_hits AS (
        -- Первая строка пересечения порога для каждого уровня.
        SELECT
            l.level_number,
            o.id AS event_id,
            o.occurred_at AS reached_at,
            o.cum_xp AS xp_at_event,
            ROW_NUMBER() OVER (
                PARTITION BY l.level_number
                ORDER BY o.occurred_at, o.id
            ) AS rn
        FROM levels l
        INNER JOIN ordered o ON l.required_experience <= o.cum_xp
        WHERE l.required_experience > 0
    ),
    hits AS (
        -- 4) Оставляем только первую подходящую строку для каждого уровня.
        SELECT
            level_number,
            event_id,
            reached_at,
            xp_at_event
        FROM first_hits
        WHERE rn = 1
    ),
    missing AS (
        -- Только те уровни, которых ещё нет в истории.
        SELECT
            h.level_number,
            h.event_id,
            h.reached_at,
            h.xp_at_event
        FROM hits h
        LEFT JOIN user_level_history ulh ON ulh.telegram_id  = _telegram_id
        AND h.level_number = ulh.level_number
        WHERE ulh.level_number IS NULL
    )
    INSERT INTO user_level_history(
        telegram_id,
        level_number,
        xp_event_id,
        xp_at_reach,
        reached_at
    )
    SELECT
        _telegram_id,
        m.level_number,
        m.event_id,
        CASE
            WHEN m.level_number = _top_level THEN m.xp_at_event
            ELSE l.required_experience
        END,
        m.reached_at
    FROM missing m
    INNER JOIN levels l ON m.level_number = l.level_number
    ON CONFLICT (
        telegram_id,
        level_number
    ) DO NOTHING;

    SELECT us.level
    INTO _old_level
    FROM user_stats us
    WHERE us.telegram_id = _telegram_id
    FOR UPDATE;

    _new_level := GREATEST(
        COALESCE(_old_level, 0),
        COALESCE(_top_level, 0)
    );

    IF _old_level IS DISTINCT FROM _new_level THEN
        UPDATE user_stats SET
            level = _new_level,
            updated_at = NOW()
        WHERE telegram_id = _telegram_id;
    END IF;

    RETURN JSONB_BUILD_OBJECT(
        'is_level_up',
            CASE
                WHEN _old_level IS NULL THEN FALSE
                ELSE (_new_level > _old_level)
            END,
        'old_level', _old_level,
        'new_level', _new_level
    );
END;
$$;