-- name: CreateUser :one
INSERT INTO users(created_at, upgraded_at, email) VALUES (
    NOW(), NOW(), $1
) RETURNING *;