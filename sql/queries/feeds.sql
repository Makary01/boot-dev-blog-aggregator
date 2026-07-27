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
    ON feeds.user_id = users.id;

-- name: GetFeed :one
SELECT feeds.id, feeds.name, feeds.url, users.name AS created_by_name 
FROM feeds 
LEFT JOIN users 
    ON feeds.user_id = users.id 
WHERE feeds.url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET updated_at = $1, last_fetch_at = $2
WHERE id = $3;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds ORDER BY last_fetch_at ASC NULLS FIRST;
