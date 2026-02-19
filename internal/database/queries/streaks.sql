-- name: GetStreak :one
SELECT * FROM streaks
WHERE user_id = $1
LIMIT 1;

-- name: UpsertStreak :one
INSERT INTO streaks (
    user_id,
    current_streak,
    longest_streak,
    last_session_date,
    updated_at
) VALUES (
    $1, $2, $3, $4, NOW()
) ON CONFLICT (user_id) DO UPDATE SET
    current_streak    = EXCLUDED.current_streak,
    longest_streak    = GREATEST(streaks.longest_streak, EXCLUDED.longest_streak),
    last_session_date = EXCLUDED.last_session_date,
    updated_at        = NOW()
RETURNING *;