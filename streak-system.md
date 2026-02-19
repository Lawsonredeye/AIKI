# Streak System

## Overview

A streak tracks how many consecutive days a user has completed at least one focus session. It is calculated automatically every time a session is ended as `completed` — no extra API call is needed.

---

## How It's Calculated

The core logic compares the **date of the last completed session** to **today's date**:

| Scenario | Result |
|----------|--------|
| First session ever | Streak starts at 1 |
| Session completed on a new consecutive day | Streak + 1 |
| Session completed on the same day as last session | Streak unchanged |
| Gap of more than 1 day since last session | Streak resets to 1 |

---

## Example Walkthrough

| Day | Action | Streak |
|-----|--------|--------|
| Monday | Complete session | 1 |
| Tuesday | Complete session | 2 |
| Tuesday | Complete another session | 2 (same day, no change) |
| Wednesday | Complete session | 3 |
| Thursday | No session | — |
| Friday | Complete session | 1 (reset, missed Thursday) |

---

## What Triggers a Streak Update

Only ending a session with `status: "completed"` triggers a streak recalculation.

| Action | Updates Streak? |
|--------|----------------|
| Start session | No |
| Pause session | No |
| Resume session | No |
| End session (`abandoned`) | No |
| End session (`completed`) | **Yes** |

---

## What's Stored

The `streaks` table holds three values per user:

| Field | Description |
|-------|-------------|
| `current_streak` | Number of consecutive days with at least one completed session |
| `longest_streak` | All-time highest streak ever achieved |
| `last_session_date` | Date of the last completed session, used for the next calculation |

The `longest_streak` is handled at the database level using `GREATEST(streaks.longest_streak, EXCLUDED.longest_streak)` — so it never decreases even if the current streak resets.

---

## API

### Get Streak
```bash
GET /api/v1/streaks
Authorization: Bearer <token>
```

Response:
```json
{
  "success": true,
  "message": "streak retrieved",
  "data": {
    "user_id": 1,
    "current_streak": 3,
    "longest_streak": 7,
    "last_session_date": "2026-02-18T00:00:00Z",
    "updated_at": "2026-02-18T10:00:00Z"
  }
}
```

### Get Home Screen (includes streak)
```bash
GET /api/v1/home
Authorization: Bearer <token>
```

The home screen endpoint returns the streak as part of the aggregated response so the frontend doesn't need a separate call.

---

## Important Notes

- A user can complete multiple sessions in one day — only the first one on a new day increments the streak. Subsequent sessions that day leave it unchanged.
- The streak is stored per user in the `streaks` table. If no record exists yet (new user), `GetStreak` returns a default zero streak rather than an error.
- Time zone is based on the server's local time when the session is completed. If your users are in different time zones, this is something to consider for a future improvement.