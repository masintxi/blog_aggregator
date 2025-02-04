-- +goose Up
CREATE TABLE posts (
    id UUID PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    title TEXT,
    url TEXT UNIQUE NOT NULL,
    description TEXT,
    published_at TIMESTAMPTZ,
    feed_id UUID  NOT NULL REFERENCES feeds (id) ON DELETE CASCADE
);

-- +goose Down
DELETE TABLE posts;