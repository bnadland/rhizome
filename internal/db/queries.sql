-- name: GetPageBySlug :one
SELECT * FROM pages
WHERE slug = $1 LIMIT 1;

-- name: SavePage :exec
INSERT INTO pages (slug, title, content)
VALUES ($1, $2, $3)
ON CONFLICT (slug)
DO UPDATE SET title = $2, content = $3
RETURNING *;