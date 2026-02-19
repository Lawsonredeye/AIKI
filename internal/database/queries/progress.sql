-- name: UpsertDailyProgress :exec
INSERT INTO daily_progress (
    user_id,
    date,
    total_focus_seconds,
    sessions_completed
) VALUES (
    $1, $2, $3, $4
) ON CONFLICT (user_id, date) DO UPDATE SET
    total_focus_seconds = daily_progress.total_focus_seconds + EXCLUDED.total_focus_seconds,
    sessions_completed  = daily_progress.sessions_completed  + EXCLUDED.sessions_completed;

-- name: GetProgressSummary :one
SELECT
    COALESCE(SUM(total_focus_seconds), 0)::bigint AS total_focus_seconds,
    COALESCE(SUM(sessions_completed), 0)::bigint  AS sessions_completed,
    COUNT(DISTINCT date)::bigint                  AS days_active
FROM daily_progress
WHERE user_id = $1
  AND date BETWEEN $2 AND $3;