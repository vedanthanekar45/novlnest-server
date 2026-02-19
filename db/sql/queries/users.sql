-- name: CreateUserOrUpdate :one
INSERT INTO users (email, google_id, name, avatar_url)
VALUES ($1, $2, $3, $4)
ON CONFLICT (email)
DO UPDATE SET
    google_id = EXCLUDED.google_id,
    name = EXCLUDED.name,
    avatar_url = EXCLUDED.avatar_url,
    updated_at = NOW()
RETURNING *; 

-- name: GetUsersByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;