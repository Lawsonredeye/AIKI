-- name: CreateFocusSession :one
INSERT INTO focus_sessions (
    user_id,
    duration_seconds,
    status
) VALUES (
    $1, $2, 'active'
) RETURNING *;

-- name: GetFocusSessionByID :one
SELECT * FROM focus_sessions
WHERE id = $1
LIMIT 1;

-- name: GetActiveSession :one
SELECT * FROM focus_sessions
WHERE user_id = $1 AND status IN ('active', 'paused')
ORDER BY started_at DESC
LIMIT 1;

-- name: UpdateFocusSession :one
UPDATE focus_sessions
SET
    elapsed_seconds = $2,
    status = $3,
    ended_at = $4,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetUserSessionHistory :many
SELECT * FROM focus_sessions
WHERE user_id = $1
ORDER BY started_at DESC
LIMIT $2 OFFSET $3;