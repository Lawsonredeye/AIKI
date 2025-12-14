-- name: CreateUserProfile :one
INSERT INTO user_profile (
    user_id,
    full_name,
    current_job,
    experience_level
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetUserProfileByUserID :one
SELECT * FROM user_profile
WHERE user_id = $1
LIMIT 1;

-- name: UpdateUserProfile :one
UPDATE user_profile
SET
    full_name = COALESCE(sqlc.narg('full_name'), full_name),
    current_job = COALESCE(sqlc.narg('current_job'), current_job),
    experience_level = COALESCE(sqlc.narg('experience_level'), experience_level),
    updated_at = NOW()
WHERE user_id = $1
RETURNING *;


-- name: UploadUserCV :one
UPDATE user_profile
SET
    cv = $2,
    updated_at = NOW()
WHERE user_id = $1
RETURNING *;