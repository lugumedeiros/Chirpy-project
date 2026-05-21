-- name: CreateUser :one
INSERT INTO users(created_at, upgraded_at, email, hashed_password) VALUES (
    NOW(), NOW(), $1, $2
) RETURNING *;

-- name: ResetUsers :exec
DELETE FROM users;

-- name: GetUser :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: UpdateUser :exec
UPDATE users SET upgraded_at = $2, email = $3, hashed_password = $4 WHERE id = $1;

-- name: UpdateRed :exec
UPDATE users SET is_chirpy_red = $2 WHERE id = $1;