CREATE OR REPLACE FUNCTION public.get_level_info(_telegram_id TEXT) RETURNS JSONB
    SECURITY DEFINER
    LANGUAGE plpgsql
AS
$$
DECLARE
    _response jsonb;
BEGIN
    IF _telegram_id IS NULL THEN
        RAISE EXCEPTION 'telegram_id IS NULL';
    END IF;

    WITH us AS (
        SELECT
            experience_points AS xp_total
        FROM user_stats
        WHERE telegram_id = _telegram_id
    ),
    lv AS (
        SELECT
            level_number,
            level_name,
            required_experience,
            LEAD(required_experience) OVER (
                ORDER BY required_experience
            ) AS next_req,
            LEAD(level_number) OVER (
                ORDER BY required_experience
            ) AS next_level
         FROM levels
    ),
    cur AS (
        SELECT lv.*
        FROM lv, us
        WHERE lv.required_experience <= us.xp_total
        ORDER BY lv.required_experience DESC
        LIMIT 1
    )
    SELECT JSONB_BUILD_OBJECT(
        'xp_total',         us.xp_total, -- текущий суммарный опыт.
        'level',            cur.level_number, -- номер текущего уровня, соответствующий самому большому порогу required_experience, который не превышает xp_total.
        'level_name',       cur.level_name, -- человеко читаемое имя текущего уровня из таблицы.
        'level_floor_xp',   cur.required_experience, -- нижняя граница текущего уровня (порог входа), то есть required_experience для поля level. Сколько XP нужно было, чтобы вообще попасть на этот уровень.
        'level_ceil_xp',    cur.next_req, -- нижняя граница следующего уровня (его порог). Если вы уже на максимальном уровне, то значение NULL.
        'next_level',       cur.next_level, -- номер следующего уровня. NULL, если текущий уровень — максимальный.
        'xp_in_level',      GREATEST(0, us.xp_total - cur.required_experience), -- сколько XP уже набрано внутри текущего уровня.
        'xp_level_size',    COALESCE(cur.next_req - cur.required_experience, 0), -- «ширина» текущего уровня в XP, то есть сколько XP нужно всего внутри уровня, чтобы перейти на следующий.
        'xp_to_next',       COALESCE(cur.next_req - us.xp_total, 0), -- сколько XP осталось до следующего уровня.
        'progress_ratio', -- доля прогресса внутри текущего уровня (для прогресс-бара), число от 0 до 1.
        CASE
            WHEN cur.next_req IS NULL OR cur.next_req = cur.required_experience THEN 1.0
            ELSE (us.xp_total - cur.required_experience)::FLOAT
                / (cur.next_req - cur.required_experience)
        END
    )
    FROM us, cur
    INTO _response;

    IF _response IS NULL THEN
        RAISE EXCEPTION 'level information data is empty';
    END IF;

    RETURN _response;
END;
$$;