-- name: CreateUser :one
INSERT INTO users(created_at, upgraded_at, email, hashed_password) VALUES (
    NOW(), NOW(), $1, $2
) RETURNING *;

-- name: ResetUsers :exec
DELETE FROM users;

-- name: GetUser :one
SELECT * FROM users WHERE email = $1;