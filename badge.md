# Badge System

## Overview

Badges are achievements awarded automatically when a user hits certain milestones. They are never manually triggered — the system checks and awards them every time a focus session is completed.

---

## Badge Catalog

There are 8 badges seeded into the `badge_definitions` table:

| Badge | Description | Criteria Type | Criteria Value |
|-------|-------------|---------------|----------------|
| First Lock-in | Complete your first focus session | sessions | 1 |
| 3-Day Streak | Maintain a 3-day streak | streak | 3 |
| 7-Day Streak | Maintain a 7-day streak | streak | 7 |
| 30-Day Streak | Maintain a 30-day streak | streak | 30 |
| Focus Rookie | Accumulate 1 hour of total focus time | focus_time | 3600 |
| Focus Pro | Accumulate 10 hours of total focus time | focus_time | 36000 |
| Job Hunter | Track 5 job applications | jobs | 5 |
| Consistency | Complete 10 focus sessions | sessions | 10 |

---

## How Badges Are Awarded

Badges are checked and awarded automatically at the end of every **completed** session. The process runs in this order:

1. Fetch all badge definitions from `badge_definitions`
2. Fetch the user's already-earned badges and skip any already awarded
3. Pull the user's all-time aggregate stats from `daily_progress`
4. Loop through every badge definition and check if the criteria is met
5. Award any badge whose threshold is now reached by inserting into `user_badges`

---

## Criteria Types

| Type | What It Checks |
|------|---------------|
| `streak` | User's `current_streak` in the `streaks` table |
| `sessions` | Total `sessions_completed` across all time in `daily_progress` |
| `focus_time` | Total `total_focus_seconds` across all time in `daily_progress` |
| `jobs` | Total job applications tracked (future implementation) |

---

## Key Behaviours

**Badges are permanent.** Once earned, a badge can never be lost — even if a streak resets, streak badges already earned are kept.

**Badges are awarded once.** The `user_badges` table has a `UNIQUE(user_id, badge_id)` constraint and inserts use `ON CONFLICT DO NOTHING`, so duplicates are impossible.

**Multiple badges can be awarded in one session.** If a user hits several thresholds at once (e.g. their 10th session also completes their 1 hour of focus time), all applicable badges are awarded together.

**Badge checks are non-fatal.** If the badge check fails for any reason, the session end still succeeds and the response is still returned to the client.

---

## Database Tables

### `badge_definitions`
Stores the static catalog of all available badges. Seeded once at migration time.

| Column | Type | Description |
|--------|------|-------------|
| `id` | int | Primary key |
| `name` | varchar | Badge name |
| `description` | text | Human-readable description |
| `icon_key` | varchar | Frontend icon identifier |
| `criteria_type` | varchar | `streak`, `sessions`, `focus_time`, or `jobs` |
| `criteria_value` | int | Threshold value to unlock |

### `user_badges`
Records which badges each user has earned.

| Column | Type | Description |
|--------|------|-------------|
| `user_id` | int | References `users.id` |
| `badge_id` | int | References `badge_definitions.id` |
| `earned_at` | timestamp | When the badge was awarded |

---

## API

### Get All Badge Definitions
Returns the full catalog of available badges regardless of whether the user has earned them.
```bash
GET /api/v1/badges
Authorization: Bearer <token>
```

Response:
```json
{
  "success": true,
  "message": "all badges retrieved",
  "data": [
    {
      "id": 1,
      "name": "First Lock-in",
      "description": "Complete your first focus session",
      "icon_key": "first_session",
      "criteria_type": "sessions",
      "criteria_value": 1
    }
  ]
}
```

### Get My Earned Badges
Returns only the badges the authenticated user has earned, ordered by most recently earned.
```bash
GET /api/v1/badges/me
Authorization: Bearer <token>
```

Response:
```json
{
  "success": true,
  "message": "badges retrieved",
  "data": [
    {
      "badge_id": 1,
      "name": "First Lock-in",
      "description": "Complete your first focus session",
      "icon_key": "first_session",
      "earned_at": "2026-02-18T10:00:00Z"
    }
  ]
}
```

### Home Screen (recent badges)
The `GET /api/v1/home` endpoint returns the 6 most recently earned badges and a total badge count as part of the aggregated home screen payload — no separate call needed.

```json
{
  "recent_badges": [...],
  "total_badges": 3
}
```

---

## Notes

- The `Job Hunter` badge criteria type is `jobs` but job count is not yet wired into the badge check — it is planned for a future update.
- Badge icon keys (e.g. `first_session`, `streak_3`) are identifiers for the frontend to map to actual icons or illustrations in the app.