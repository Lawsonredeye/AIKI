-- name: CreateJob :one
INSERT INTO jobs (
    user_id,
    title,
    company_name,
    notes,
    link,
    location,
    platform,
    date_applied,
    status
) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetJobByID :one
SELECT * FROM jobs
WHERE id = $1
LIMIT 1;

-- name: GetJobs :many
SELECT * FROM jobs
WHERE user_id = $1;

-- name: UpdateJobByID :exec
UPDATE jobs
SET
    title = COALESCE(sqlc.narg('title')::text, title),
    company_name = COALESCE(sqlc.narg('company_name')::text, company_name),
    notes = COALESCE(sqlc.narg('notes')::text, notes),
    link = COALESCE(sqlc.narg('link')::text, link),
    location = COALESCE(sqlc.narg('location')::text, location),
    platform = COALESCE(sqlc.narg('platform')::text, platform),
    date_applied = COALESCE(sqlc.narg('date_applied')::timestamp, date_applied),
    status = COALESCE(sqlc.narg('status')::text, status),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteJobByID :exec
DELETE FROM jobs
WHERE id = $1;

