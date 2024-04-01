-- migrate:up
CREATE TABLE loc_admin (
    id bigserial PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users (id) ON DELETE SET NULL,
    location_id bigint NOT NULL REFERENCES location (id) ON DELETE CASCADE,
    permissions text[],
    UNIQUE (user_id, location_id)
);

-- migrate:down
DROP TABLE loc_admin;
