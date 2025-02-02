-- +goose Up
CREATE TABLE feeds (
    id UUID PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    name TEXT UNIQUE NOT NULL,
    url TEXT UNIQUE NOT NULL,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE 
    -- CONSTRAINT fk_users FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +goose Down
DROP TABLE feeds;