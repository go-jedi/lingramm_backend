CREATE OR REPLACE FUNCTION public.subscription_create(_telegram_id TEXT) RETURNS subscriptions
    SECURITY DEFINER
    LANGUAGE plpgsql
AS
$$
DECLARE
    _s subscriptions; -- subscription.
    _sat TIMESTAMP; -- subscribed_at.
    _exp TIMESTAMP; -- expires_at.
BEGIN
    _sat = NOW();
    _exp = _sat + INTERVAL '1 month';

    -- create subscription.
    UPDATE subscriptions SET
        subscribed_at = _sat,
        expires_at = _exp,
        is_active = TRUE,
        updated_at = NOW()
    WHERE telegram_id = _telegram_id
    RETURNING * INTO _s;

    -- create subscription history.
    INSERT INTO subscription_history(
        telegram_id,
        action_time,
        expires_at
    ) VALUES(
        _telegram_id,
        _sat,
        _exp
    );

    RETURN _s;
END;
$$;