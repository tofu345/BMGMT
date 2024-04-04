-- migrate:up
DROP TYPE IF EXISTS issue_type;
CREATE TYPE issue_type AS ENUM (
    'electrical',
    'plumbing'
);

CREATE TABLE location_issues (
    id bigserial PRIMARY KEY,
    room_id bigint NOT NULL REFERENCES room (id) ON DELETE CASCADE,
    issue_type issue_type NOT NULL,
    info text NOT NULL,
    resolved bool NOT NULL
);

-- migrate:down
DROP TABLE location_issues;
