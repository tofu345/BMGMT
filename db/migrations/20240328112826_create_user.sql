-- migrate:up
CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    email varchar(32) NOT NULL UNIQUE,
    first_name varchar(32) NOT NULL,
    last_name varchar(32) NOT NULL,
    password varchar(100) NOT NULL,
    is_superuser boolean,
    created_at timestamp DEFAULT NOW() NOT NULL,
    updated_at timestamp DEFAULT NOW() NOT NULL
);

CREATE OR REPLACE FUNCTION update_updated_at_col ()
    RETURNS TRIGGER
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = (now() at time zone 'UTC');
    RETURN NEW;
END;
$$;

CREATE OR REPLACE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE PROCEDURE update_updated_at_col ();

-- migrate:down
DROP TABLE IF EXISTS users
