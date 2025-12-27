-- name: IsUserExists :one
SELECT user_id FROM users
WHERE user_id = $1;

-- name: CreateUser :exec
INSERT INTO users (user_id, username, fullname, register_date)
VALUES ($1, $2, $3, now())
ON CONFLICT (user_id) DO UPDATE
SET username = EXCLUDED.username,
    fullname = EXCLUDED.fullname;

-- name: GetUser :one
SELECT * FROM users
WHERE user_id = $1;

-- name: GetUsersCount :one
SELECT COUNT(*) FROM users;