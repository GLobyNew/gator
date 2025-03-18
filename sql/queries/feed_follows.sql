-- name: CreateFeedFollow :one
WITH new_row AS (
    INSERT INTO feed_follows(id, created_at, updated_at, user_id, feed_id)
    VALUES(
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
SELECT
    nr.*,
    feeds.name as feed_name,
    users.name as user_name
FROM
    new_row nr
INNER JOIN users ON nr.user_id = users.id
INNER JOIN feeds ON nr.feed_id = feeds.id;

-- name: GetFeedFollowsForUser :many 
SELECT 
    feed_follows.id, 
    feed_follows.user_id, 
    users.name AS user_name, 
    feed_follows.feed_id, 
    feeds.name AS feed_name,
    feeds.url AS feed_url
FROM 
    feed_follows
INNER JOIN 
    users ON feed_follows.user_id = users.id
INNER JOIN 
    feeds ON feed_follows.feed_id = feeds.id
WHERE 
    feed_follows.user_id = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows
USING users, feeds
WHERE users.name = $1
AND feeds.url = $2;