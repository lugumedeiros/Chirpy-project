-- name: CreateChirp :one
INSERT INTO chirps(created_at, updated_at, body, user_id) VALUES (
    NOW(), NOW(), $1, $2
) RETURNING *;

-- name: DeleteChirpsByUserId :exec
DELETE FROM chirps WHERE user_id = $1;

-- name: DeleteChirp :exec
DELETE FROM chirps WHERE id = $1;

-- name: DeleteChirp2Step :exec
DELETE FROM chirps WHERE id = $1 AND user_id = $2;

-- name: DeleteAllChirps :exec
DELETE FROM chirps;

-- name: GetAllChirps :many
SELECT * FROM chirps;

-- name: GetChirp :one
SELECT * FROM chirps WHERE id = $1;

-- name: GetAllChirpsByUserId :many
SELECT * FROM chirps WHERE user_id = $1;