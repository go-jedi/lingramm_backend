CREATE OR REPLACE FUNCTION public.xp_event_create(_telegram_id TEXT, _src JSONB) RETURNS VOID
    SECURITY DEFINER
    LANGUAGE plpgsql
AS
$$
BEGIN
    IF _telegram_id IS NULL THEN
        RAISE EXCEPTION 'telegram_id IS NULL';
    END IF;

    INSERT INTO xp_events(
        telegram_id,
        delta_xp,
        reason
    )
    SELECT _telegram_id, s.delta_xp, s.reason
    FROM JSONB_TO_RECORDSET(_src) AS s(
        delta_xp INTEGER,
        reason VARCHAR(50)
    )
    INNER JOIN users u ON _telegram_id = u.telegram_id
    WHERE _telegram_id IS NOT NULL
    AND s.delta_xp IS NOT NULL
    AND s.delta_xp <> 0
    AND s.reason IS NOT NULL;

    RETURN;
END;
$$;