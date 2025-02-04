-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS(
    INSERT INTO feed_follows (
        id, created_at, updated_at, user_id, feed_id
    )
    VALUES (
        $1, $2, $3,
        (SELECT id FROM users WHERE users.name = $4),
        (SELECT id FROM feeds WHERE feeds.url = $5)
    ) 
    RETURNING *
)
SELECT
    inserted_feed_follow.*,
    users.name as user_name,
    feeds.name as feed_name
FROM inserted_feed_follow
INNER JOIN users ON users.id = inserted_feed_follow.user_id
INNER JOIN feeds ON feeds.id = inserted_feed_follow.feed_id;

-- name: GetFeedFollowsForUser :many
SELECT
    feed_follows.*,
    users.name as user_name,
    feeds.name as feed_name
FROM feed_follows 
INNER JOIN users ON feed_follows.user_id = users.id 
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE users.name = $1
ORDER BY feeds.name;

-- name: DeleteFollowForUser :exec
DELETE FROM feed_follows
USING feeds, users
WHERE feed_follows.feed_id = feeds.id 
AND feed_follows.user_id = users.id
AND feeds.url = $1
AND users.name = $2;

-- DELETE FROM feed_follows
-- WHERE feed_follows.feed_id = (
--     SELECT feeds.id
--     FROM feeds
--     INNER JOIN users ON feeds.user_id = users.id
--     WHERE feeds.url = $1 AND users.name = $2
-- );


