--Шаблон дял прав
CREATE ROLE app_role NOSUPERUSER NOCREATEDB NOCREATEROLE NOINHERIT;


GRANT USAGE ON SCHEMA public TO app_role;
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE
    users, offers, complexes, profiles
    TO app_role;

-- Если есть sequence (например, для ID)
GRANT USAGE ON SEQUENCE users_id_seq, offers_id_seq, complexes_id_seq TO app_role;


CREATE USER app_user WITH PASSWORD 'strong_password_123!';  -- ← замени на генерируемый или из .env
GRANT app_role TO app_user;
ALTER ROLE app_role SET statement_timeout = '10s';   -- 10 секунд
ALTER ROLE app_role SET lock_timeout = '2s';