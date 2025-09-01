CREATE OR REPLACE FUNCTION public.leaderboard_weeks_process_batch(
    _worker_name TEXT, -- имя воркера (например, 'main').
    _batch_size INTEGER, -- целевой размер batch (~100–300 мс на вызов).
    _statement_timeout_ms INTEGER DEFAULT NULL, -- локальный statement_timeout (мс).
    _lock_timeout_ms INTEGER DEFAULT NULL -- локальный lock_timeout (мс).
) RETURNS JSONB
    SECURITY DEFINER
    LANGUAGE plpgsql
AS
$$
DECLARE
    _last_event_id BIGINT;
    _current_max_id BIGINT;
    _batch_count INTEGER := 0;
    _new_event_count INTEGER := 0;
    _groups_count INTEGER := 0;
    _total_add_xp BIGINT  := 0;
    _batch_max_id BIGINT;
    _eff_batch_size INTEGER;
    _response JSONB;
BEGIN
    -- JIT на коротких батчах обычно мешает latency.
    PERFORM SET_CONFIG('jit', 'off', TRUE);

    -- Локальные таймауты на время этого вызова
    IF _statement_timeout_ms IS NOT NULL THEN
        PERFORM set_config('statement_timeout', _statement_timeout_ms || 'ms', TRUE);
    END IF;
    IF _lock_timeout_ms IS NOT NULL THEN
        PERFORM set_config('lock_timeout', _lock_timeout_ms || 'ms', TRUE);
    END IF;

    -- клампим batch_size (LIMIT не любит 0/отрицательные).
    _eff_batch_size := GREATEST(COALESCE(_batch_size, 0), 1);

    INSERT INTO leaderboard_weeks_worker_state(
        name,
        last_event_id
    )
    VALUES (
        _worker_name,
        0
    )
    ON CONFLICT (name) DO NOTHING;

    -- lock строку чекпоинта (единственный активный worker с этим name).
    SELECT
        last_event_id
    INTO _last_event_id
    FROM leaderboard_weeks_worker_state
    WHERE name = _worker_name
    FOR UPDATE;

    -- фиксируем "потолок" (верхнюю границу) на момент старта итерации.
    SELECT
        COALESCE(MAX(id), 0)
    INTO _current_max_id
    FROM xp_events;

    -- если нечего обрабатывать — возвращаем JSON сразу.
    IF _current_max_id <= _last_event_id THEN
        _response := JSONB_BUILD_OBJECT(
            'processed', FALSE,
            'from_id', _last_event_id,
            'to_id', _current_max_id,
            'batch_count', 0,
            'new_event_count', 0,
            'groups_count', 0,
            'applied_xp', 0,
            'new_last_event_id', _last_event_id
        );
        RETURN _response;
    END IF;

    -- основной CTE-поток (одна транзакция, один план).
    WITH batch AS MATERIALIZED (
        SELECT
            id
        FROM xp_events
        WHERE id > _last_event_id
        AND id <= _current_max_id
        ORDER BY id
        LIMIT _eff_batch_size
    ),
    applied AS MATERIALIZED (
        INSERT INTO leaderboard_weeks_applied_events(
            event_id
        )
        SELECT
            id
        FROM batch
        ON CONFLICT (event_id) DO NOTHING
        RETURNING event_id
    ),
    delta AS MATERIALIZED (
        SELECT
            xpe.week_start,
            xpe.telegram_id,
            SUM(xpe.delta_xp)::BIGINT AS add_xp
        FROM xp_events xpe
        INNER JOIN applied a ON xpe.id = a.event_id
        GROUP BY xpe.week_start, xpe.telegram_id
        HAVING SUM(xpe.delta_xp) <> 0
    ),
    upsert AS MATERIALIZED (
        INSERT INTO leaderboard_weeks(
            week_start,
            telegram_id,
            xp
        )
        SELECT week_start, telegram_id, add_xp
        FROM delta
        ON CONFLICT (week_start, telegram_id)
        DO UPDATE SET xp = leaderboard_weeks.xp + EXCLUDED.xp
        RETURNING 1
    ),
    stats AS (
        SELECT
            (
                SELECT
                    COUNT(*)
                FROM batch
            )::INTEGER AS batch_cnt,
            (
                SELECT
                    COUNT(*)
                FROM applied
            )::INTEGER AS new_ev_cnt,
            (
                SELECT COUNT(*)
                FROM delta
            )::INTEGER AS groups_cnt,
            COALESCE((
                SELECT
                    SUM(add_xp)
                FROM delta
            ), 0)::BIGINT AS total_add_xp,
            COALESCE((
                SELECT
                    MAX(id)
                FROM batch
            ), _last_event_id)::BIGINT AS batch_max
    )
    UPDATE leaderboard_weeks_worker_state ws SET
        last_event_id = stats.batch_max,
        updated_at = NOW()
    FROM stats
    WHERE ws.name = _worker_name
    RETURNING
        stats.batch_cnt,
        stats.new_ev_cnt,
        stats.groups_cnt,
        stats.total_add_xp,
        stats.batch_max
    INTO
        _batch_count,
        _new_event_count,
        _groups_count,
        _total_add_xp,
        _batch_max_id;

    _response := JSONB_BUILD_OBJECT(
            'processed', (_batch_count > 0),
            'from_id', _last_event_id,
            'to_id', _current_max_id,
            'batch_count', _batch_count,
            'new_event_count', _new_event_count,
            'groups_count', _groups_count,
            'applied_xp', _total_add_xp,
            'new_last_event_id', _batch_max_id
    );

    RETURN _response;
END;
$$;