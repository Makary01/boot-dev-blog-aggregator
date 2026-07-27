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
INNER JOIN feed_follows ON feed_follows.feed_id = posts.feed_id
WHERE feed_follows.user_id = $1
ORDER BY posts.created_at DESC 
LIMIT $2;
