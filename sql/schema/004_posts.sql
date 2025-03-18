-- +goose Up
CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    published_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL,
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;