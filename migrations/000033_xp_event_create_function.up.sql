CREATE OR REPLACE FUNCTION public.xp_event_create(
    _telegram_id TEXT,
    _event_type TEXT,
    _delta_xp INTEGER
) RETURNS VOID
    SECURITY DEFINER
    LANGUAGE plpgsql
AS
$$
BEGIN
    IF _telegram_id IS NULL THEN
        RAISE EXCEPTION 'telegram_id IS NULL';
    END IF;
    IF _event_type IS NULL THEN
        RAISE EXCEPTION 'event_type IS NULL';
    END IF;
    IF _delta_xp IS NULL THEN
        RAISE EXCEPTION 'delta_xp IS NULL';
    END IF;

    INSERT INTO xp_events(
        event_type_id,
        telegram_id,
        delta_xp
    )
    SELECT
        et.id,
        _telegram_id,
        _delta_xp
    FROM event_types et
    WHERE et.name = _event_type
    AND _delta_xp IS NOT NULL
    AND _delta_xp <> 0
    AND EXISTS (
        SELECT 1
        FROM users u
        WHERE u.telegram_id = _telegram_id
    );

    RETURN;
END;
$$;