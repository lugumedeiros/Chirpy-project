-- name: AddRefreshToken :one
INSERT INTO refresh_tokens(token, user_id, expires_at, revoked_at, created_at, updated_at)
VALUES($1, $2, $3, $4, NOW(), NOW()) RETURNING *;

-- name: GetRefreshTokenToken :one
SELECT * FROM refresh_tokens WHERE token = $1;

-- name: GetRefreshTokenUser :one
SELECT * FROM refresh_tokens WHERE user_id = $1;

-- name: ResetRefreshToken :exec
DELETE FROM refresh_tokens;

-- name: RevokeToken :exec
UPDATE refresh_tokens SET revoked_at = $2 WHERE token = $1;