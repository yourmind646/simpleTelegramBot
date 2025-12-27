-- name: IsAdminExists :one
SELECT user_id FROM admins
WHERE user_id = $1;

-- name: CreateAdmin :exec
INSERT INTO admins (user_id)
VALUES ($1);
