CREATE TYPE notifications_status AS ENUM ('PENDING', 'SENT', 'FAILED');
CREATE TYPE notifications_type AS ENUM ('achievement', 'internal_currency', 'level', 'mini_game');