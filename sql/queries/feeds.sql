-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeedsByUser :one
SELECT * FROM feeds INNER JOIN users ON feeds.user_id = users.id WHERE users.name = $1;

-- name: ResetFeedsDB :exec
TRUNCATE TABLE feeds;

-- name: DeleteFeedByURL :exec
DELETE FROM feeds WHERE url = $1;

-- name: ListFeeds :many
SELECT * FROM feeds ORDER BY name;

-- name: GetFeedByURL :one
SELECT * FROM feeds WHERE url = $1;

-- name: GetFeedByID :one
SELECT * FROM feeds WHERE id = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = NOW(), updated_at = NOW()
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY last_fetched_at NULLS FIRST
LIMIT 1;