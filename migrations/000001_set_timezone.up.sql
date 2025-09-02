ALTER DATABASE postgres RESET timezone;
ALTER DATABASE lingramm_db
SET timezone = 'Europe/Moscow';
-- roles.
ALTER ROLE admin RESET timezone;