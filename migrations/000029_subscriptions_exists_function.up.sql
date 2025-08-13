CREATE OR REPLACE FUNCTION public.subscription_exists(_telegram_id TEXT) RETURNS BOOLEAN
    SECURITY DEFINER
    LANGUAGE plpgsql
AS
$$
DECLARE
    _s subscriptions;
BEGIN
    SELECT *
    FROM subscriptions
    WHERE telegram_id = _telegram_id
    INTO _s;

    IF _s.is_active AND NOW() >= _s.expires_at THEN
        UPDATE subscriptions SET
            subscribed_at = NULL,
            expires_at = NULL,
            is_active = FALSE
        WHERE telegram_id = _telegram_id
        RETURNING * INTO _s;
    END IF;

    RETURN _s.is_active;
END;
$$;