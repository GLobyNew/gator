-- name: CreatePost :one
WITH NP AS (
    INSERT INTO posts(id, created_at, updated_at, published_at, title, url, description, feed_id)
    VALUES(
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8
    )
    RETURNING *
)
SELECT 
    NP.*,
    feeds.name as feed_name
FROM 
    NP
INNER JOIN feeds ON NP.feed_id = feeds.id;


-- name: GetPostsByUser :many
SELECT 
    posts.id,
    posts.title,
    posts.url,
    posts.description,
    posts.published_at,
    feeds.id AS feed_id,
    feeds.name AS feed_name,
    users.id AS user_id,
    users.name AS user_name
FROM 
    posts
INNER JOIN
    feeds ON posts.feed_id = feeds.id
INNER JOIN 
    users ON feeds.user_id = users.id
WHERE 
    users.id = $1
ORDER BY posts.published_at
LIMIT $2;