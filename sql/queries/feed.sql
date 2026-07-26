-- name: CreateFeed :one
INSERT INTO feeds (id, user_id, created_at, updated_at, name, url)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT feeds.id, feeds.name, feeds.url, users.name AS created_by_name 
FROM feeds 
LEFT JOIN users 
    ON feeds.user_id = users.id ;
