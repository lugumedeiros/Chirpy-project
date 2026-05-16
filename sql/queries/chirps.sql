-- name: CreateChirp :one
INSERT INTO chirps(created_at, updated_at, body, user_id) VALUES (
    NOW(), NOW(), $1, $2
) RETURNING *;

-- name: DeleteChirpsByUserId :exec
DELETE FROM chirps WHERE user_id = $1;

-- name: DeleteChirp :exec
DELETE FROM chirps WHERE id = $1;

-- name: DeleteAllChirps :exec
DELETE FROM chirps;