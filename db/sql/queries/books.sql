-- name: CreateShelf :one
INSERT INTO shelves (userid, title, description)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetShelvesByUser :many
SELECT * FROM SHELVES
WHERE userid = $1
ORDER BY created_at DESC;

-- name: AddBookToShelf :exec
INSERT INTO shelf_books (shelf_id, isbn13, google_id, title, thumbnail_url)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (shelf_id, google_id) DO NOTHING;

-- name: GetBooksInShelf :many 
SELECT * FROM shelf_books
WHERE shelf_id = $1
ORDER BY added_at DESC;

-- name: UpsertBookLog :one
INSERT INTO book_logs (userid, isbn13, google_id, title, thumbnail_url, status, rating, review)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (userid, google_id)
DO UPDATE SET 
    status = EXCLUDED.status,
    rating = EXCLUDED.rating,
    review = EXCLUDED.review,
    updated_at = NOW()
RETURNING *;

-- name: GetUserLogs :many
SELECT * FROM book_logs
WHERE userid = $1
ORDER BY updated_at DESC;