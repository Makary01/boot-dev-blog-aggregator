-- name: CreateFeedFollow :one
WITH inserted AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
SELECT inserted.id, users.name AS user_name, feeds.name AS feed_name 
FROM inserted
JOIN users ON users.id = inserted.user_id
JOIN feeds ON feeds.id = inserted.feed_id;

-- name: GetFeedFollowings :many
SELECT feed_follows.id, users.name AS user_name, feeds.name AS feed_name 
FROM feed_follows
JOIN users ON users.id = feed_follows.user_id
JOIN feeds ON feeds.id = feed_follows.feed_id
WHERE users.name = $1;

-- name: DeleteFeedFollowing :exec
DELETE FROM feed_follows
USING users, feeds
WHERE feed_follows.user_id = users.id
AND feed_follows.feed_id = feeds.id
AND feeds.url = $1 
AND users.name = $2;
