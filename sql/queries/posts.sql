-- name: CreatePost :exec
INSERT INTO posts (id, feed_id, created_at, updated_at, published_at, title, url, description)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;

-- name: GetPostsForUser :many
SELECT * 
FROM posts 
ORDER BY created_at DESC 
LIMIT $1;
