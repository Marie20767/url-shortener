-- name: GetUserByID :one
SELECT id, name, email FROM users WHERE id = $1;

-- name: GetUserByName :one
SELECT id, name, email FROM users WHERE name = $1;