# Focus Session Flow

## Overview

The session lifecycle is split between the **client** and the **server**. The server stores session state and snapshots of elapsed time — the actual timer countdown runs on the mobile app. This avoids the need for websockets or polling.

---

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/sessions` | Start a new session |
| GET | `/api/v1/sessions/active` | Get current active or paused session |
| PATCH | `/api/v1/sessions/:id/pause` | Pause a session |
| PATCH | `/api/v1/sessions/:id/resume` | Resume a paused session |
| PATCH | `/api/v1/sessions/:id/end` | End a session (completed or abandoned) |
| GET | `/api/v1/sessions` | Get session history |

---

## Lifecycle

```
Start → [timer runs on client] → Pause ↔ Resume → End (completed / abandoned)
                                                          ↓
                                               streak + progress + badges update
```

---

## Step by Step

### 1. Start a Session — `POST /api/v1/sessions`

The user picks a duration and taps "Lock-in". The client sends the planned duration in seconds.

**Request:**
```json
{
  "duration_seconds": 1500
}
```

**Duration options from the timer picker:**
| Label | Seconds |
|-------|---------|
| 25 min | 1500 |
| 45 min | 2700 |
| 60 min | 3600 |
| Custom | 60 – 86400 |

**Response:**
```json
{
  "success": true,
  "message": "focus session started",
  "data": {
    "id": 1,
    "user_id": 1,
    "duration_seconds": 1500,
    "elapsed_seconds": 0,
    "status": "active",
    "started_at": "2026-02-18T08:00:00Z",
    "ended_at": null
  }
}
```

> A user can only have one `active` or `paused` session at a time. Starting a second session while one is running returns a `409 Conflict` error.

---

### 2. Timer Runs on the Client

Once the session is started, the app counts down locally using the `duration_seconds` value. The server is not involved until the user pauses, resumes, or ends the session.

---

### 3. Pause — `PATCH /api/v1/sessions/:id/pause`

User taps "Pause session". The client reports how many seconds have elapsed so far.

**Request:**
```json
{
  "elapsed_seconds": 750,
  "status": "paused"
}
```

The session status becomes `"paused"` and `elapsed_seconds` is saved. This means if the user closes the app, their progress is not lost.

---

### 4. Resume — `PATCH /api/v1/sessions/:id/resume`

User taps "Resume". No request body is needed — the API flips the status back to `"active"`. The client resumes its local timer from the saved `elapsed_seconds`.

**Response:**
```json
{
  "data": {
    "id": 1,
    "elapsed_seconds": 750,
    "status": "active"
  }
}
```

---

### 5. End — `PATCH /api/v1/sessions/:id/end`

Ends the session as either **completed** or **abandoned**.

**Completed** — user finished the full session:
```json
{
  "elapsed_seconds": 1500,
  "status": "completed"
}
```

**Abandoned** — user quit early:
```json
{
  "elapsed_seconds": 400,
  "status": "abandoned"
}
```

> Only `"completed"` triggers post-session side effects. Abandoned sessions are recorded but do not update streaks, progress, or badges.

---

### 6. Post-Session Side Effects (Completed Only)

When a session is ended with `status: "completed"`, three things happen automatically in the background:

**a) Daily Progress Update**
Today's `daily_progress` row is upserted with the session's focus time and session count.

**b) Streak Recalculation**

| Scenario | Result |
|----------|--------|
| Session completed on a new consecutive day | Streak + 1 |
| Session completed on the same day as last session | Streak unchanged |
| Gap of more than 1 day since last session | Streak resets to 1 |

**c) Badge Check**
All badge definitions are scanned. Any badge whose criteria is now met is automatically awarded. No extra API call is needed from the client.

---

### 7. Restore Session on App Reopen — `GET /api/v1/sessions/active`

When the app reopens, call this endpoint to check if there is an existing `active` or `paused` session. If found, the client restores the timer UI from `elapsed_seconds` and `duration_seconds`.

Returns `null` in the data field if no session is currently running.

---

## Session Statuses

| Status | Description |
|--------|-------------|
| `active` | Timer is running |
| `paused` | Timer is paused, progress saved |
| `completed` | Session finished successfully |
| `abandoned` | Session ended early |

---

## Notes

- The server never auto-expires or auto-completes a session — the client is always responsible for calling the end endpoint.
- `elapsed_seconds` should always be sent accurately from the client on pause and end, as this is what gets recorded in progress stats.
- The session `id` returned from the start endpoint should be stored by the client for all subsequent pause/resume/end calls.