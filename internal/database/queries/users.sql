-- name: CreateUser :one
INSERT INTO users (
    first_name,
    last_name,
    email,
    phone_number,
    password_hash,
    linkedin_id
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 AND is_active = TRUE
LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 AND is_active = TRUE
LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET
    first_name = COALESCE(sqlc.narg('first_name'), first_name),
    last_name = COALESCE(sqlc.narg('last_name'), last_name),
    phone_number = COALESCE(sqlc.narg('phone_number'), phone_number),
    linkedin_id = COALESCE(sqlc.narg('linkedin_id'), linkedin_id)
WHERE id = $1 AND is_active = TRUE
RETURNING *;

-- name: DeactivateUser :exec
UPDATE users
SET is_active = FALSE
WHERE id = $1;

-- name: CheckEmailExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = $1) AS exists;