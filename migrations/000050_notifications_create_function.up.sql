CREATE OR REPLACE FUNCTION public.notifications_create(_src JSONB) RETURNS JSONB
    SECURITY DEFINER
    LANGUAGE plpgsql
AS
$$
DECLARE
    _response JSONB;
BEGIN
    WITH payload AS (
        SELECT
            s.type,
            s.telegram_id,
            JSONB_BUILD_OBJECT(
                'title', COALESCE(s.message->>'title', ''),
                'text', COALESCE(s.message->>'text', '')
            ) AS message
        FROM JSONB_TO_RECORDSET(_src) AS s(
            type TEXT,
            telegram_id TEXT,
            message JSONB
        )
        WHERE s.type IS NOT NULL
        AND s.telegram_id IS NOT NULL
        AND s.message IS NOT NULL
        AND s.type = ANY (SELECT UNNEST(ENUM_RANGE(NULL::notifications_type))::text)
    ),
    ins AS (
        INSERT INTO notifications(
            type,
            telegram_id,
            message
        )
        SELECT p.type::notifications_type, p.telegram_id, p.message
        FROM payload p
        INNER JOIN users u ON p.telegram_id = u.telegram_id
        RETURNING id, message, type, telegram_id, status, created_at, sent_at
    )
    SELECT COALESCE(
        JSONB_AGG(
            JSONB_BUILD_OBJECT(
                'id', i.id,
                'message',    i.message,
                'type',       i.type,
                'telegram_id',i.telegram_id,
                'status',     i.status,
                'created_at', i.created_at,
                'sent_at',    i.sent_at
            )
            ORDER BY i.id
        ),
        '[]'::JSONB
    )
    INTO _response
    FROM ins i;

    RETURN _response;
END;
$$;