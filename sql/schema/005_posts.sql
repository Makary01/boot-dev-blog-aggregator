-- +goose Up
CREATE TABLE posts (
    id UUID PRIMARY KEY,
    feed_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    published_at TIMESTAMP,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    description TEXT,
    UNIQUE(url),
    CONSTRAINT fd_feed_id 
        FOREIGN KEY (feed_id) 
        REFERENCES feeds(id)
);

-- +goose Down
DROP TABLE posts;
