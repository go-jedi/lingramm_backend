ALTER DATABASE postgres RESET timezone;
ALTER DATABASE lingvogramm_db SET timezone = 'Europe/Moscow';

-- roles.
ALTER ROLE admin RESET timezone;