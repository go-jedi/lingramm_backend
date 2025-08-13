CREATE OR REPLACE FUNCTION public.user_create(_src JSON) RETURNS users
    SECURITY DEFINER
    LANGUAGE plpgsql
AS
$$
DECLARE
    _u users;
BEGIN
    INSERT INTO users(
        telegram_id,
        username,
        first_name,
        last_name
    ) VALUES(
        _src->>'telegram_id',
        _src->>'username',
        _src->>'first_name',
        _src->>'last_name'
    )
    RETURNING * INTO _u;

    INSERT INTO user_balances(
        telegram_id
    ) VALUES(
        _src->>'telegram_id'
    );

    INSERT INTO user_stats(
        telegram_id
    ) VALUES(
        _src->>'telegram_id'
    );

    INSERT INTO subscriptions(
        telegram_id
    ) VALUES(
        _src->>'telegram_id'
    );

    RETURN _u;
END;
$$;