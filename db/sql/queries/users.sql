-- name: CreateUserOrUpdate :one
INSERT INTO users (email, google_id, name, avatar_url)
VALUES ($1, $2, $3, $4)
ON CONFLICT (email)
