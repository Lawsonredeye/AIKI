-- name: GetAllBadgeDefinitions :many
SELECT * FROM badge_definitions
ORDER BY id;

-- name: GetUserBadges :many
SELECT
    bd.id,
    bd.name,
    bd.description,
    bd.icon_key,
    bd.criteria_type,
    bd.criteria_value,
    ub.earned_at
FROM user_badges ub
JOIN badge_definitions bd ON bd.id = ub.badge_id
WHERE ub.user_id = $1
ORDER BY ub.earned_at DESC;

-- name: AwardBadge :exec
INSERT INTO user_badges (user_id, badge_id)
VALUES ($1, $2)
ON CONFLICT (user_id, badge_id) DO NOTHING;

-- name: GetUserBadgeCount :one
SELECT COUNT(*) FROM user_badges
WHERE user_id = $1;