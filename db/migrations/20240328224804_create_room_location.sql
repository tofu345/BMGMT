-- migrate:up
CREATE TABLE location (
    id bigserial PRIMARY KEY,
    name varchar(32) NOT NULL,
    address varchar(100) NOT NULL
);

CREATE TABLE room (
    id bigserial PRIMARY KEY,
    name varchar(32) NOT NULL,
    location_id bigint REFERENCES location (id) ON DELETE CASCADE,
    tenant_id bigint REFERENCES users (id) ON DELETE SET NULL,
    tenancy_end_date timestamp
);

-- migrate:down
DROP TABLE location, room CASCADE;

