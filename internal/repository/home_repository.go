package repository

import (
	"aiki/internal/domain"
	"aiki/internal/pkg/response"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate mockgen -source=home_repository.go -destination=mocks/mock_home_repository.go -package=mocks

type HomeRepository interface {
	// Focus Sessions
	CreateSession(ctx context.Context, userID int32, durationSeconds int32) (*domain.FocusSession, error)
	GetSessionByID(ctx context.Context, sessionID int32) (*domain.FocusSession, error)
	GetActiveSession(ctx context.Context, userID int32) (*domain.FocusSession, error)
	UpdateSession(ctx context.Context, sessionID int32, elapsedSeconds int32, status domain.SessionStatus, endedAt *time.Time) (*domain.FocusSession, error)
	GetUserSessionHistory(ctx context.Context, userID int32, limit, offset int32) ([]domain.FocusSession, error)

	// Streaks
	GetStreak(ctx context.Context, userID int32) (*domain.Streak, error)
	UpsertStreak(ctx context.Context, userID int32, currentStreak, longestStreak int32, lastDate *time.Time) (*domain.Streak, error)

	// Badges
	GetAllBadgeDefinitions(ctx context.Context) ([]domain.BadgeDefinition, error)
	GetUserBadges(ctx context.Context, userID int32) ([]domain.UserBadge, error)
	AwardBadge(ctx context.Context, userID, badgeID int32) error
	GetUserBadgeCount(ctx context.Context, userID int32) (int32, error)

	// Progress
	UpsertDailyProgress(ctx context.Context, userID int32, date time.Time, focusSeconds, sessions int32) error
	GetProgressSummary(ctx context.Context, userID int32, from, to time.Time) (*domain.ProgressSummary, error)
}

type homeRepository struct {
	db *pgxpool.Pool
}

func NewHomeRepository(dbPool *pgxpool.Pool) HomeRepository {
	return &homeRepository{db: dbPool}
}

// ─────────────────────────────────────────
// Focus Sessions
// ─────────────────────────────────────────

func (r *homeRepository) CreateSession(ctx context.Context, userID int32, durationSeconds int32) (*domain.FocusSession, error) {
	query := `
		INSERT INTO focus_sessions (user_id, duration_seconds, status)
		VALUES ($1, $2, 'active')
		RETURNING id, user_id, duration_seconds, elapsed_seconds, status, started_at, ended_at, created_at, updated_at
	`
	return r.scanSession(r.db.QueryRow(ctx, query, userID, durationSeconds))
}

func (r *homeRepository) GetSessionByID(ctx context.Context, sessionID int32) (*domain.FocusSession, error) {
	query := `
		SELECT id, user_id, duration_seconds, elapsed_seconds, status, started_at, ended_at, created_at, updated_at
		FROM focus_sessions WHERE id = $1
	`
	s, err := r.scanSession(r.db.QueryRow(ctx, query, sessionID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, response.ErrSessionNotFound
		}
		return nil, err
	}
	return s, nil
}

func (r *homeRepository) GetActiveSession(ctx context.Context, userID int32) (*domain.FocusSession, error) {
	query := `
		SELECT id, user_id, duration_seconds, elapsed_seconds, status, started_at, ended_at, created_at, updated_at
		FROM focus_sessions
		WHERE user_id = $1 AND status IN ('active', 'paused')
		ORDER BY started_at DESC
		LIMIT 1
	`
	s, err := r.scanSession(r.db.QueryRow(ctx, query, userID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // no active session is not an error
		}
		return nil, err
	}
	return s, nil
}

func (r *homeRepository) UpdateSession(ctx context.Context, sessionID int32, elapsedSeconds int32, status domain.SessionStatus, endedAt *time.Time) (*domain.FocusSession, error) {
	query := `
		UPDATE focus_sessions
		SET elapsed_seconds = $2, status = $3, ended_at = $4, updated_at = NOW()
		WHERE id = $1
		RETURNING id, user_id, duration_seconds, elapsed_seconds, status, started_at, ended_at, created_at, updated_at
	`
	s, err := r.scanSession(r.db.QueryRow(ctx, query, sessionID, elapsedSeconds, status, endedAt))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, response.ErrSessionNotFound
		}
		return nil, err
	}
	return s, nil
}

func (r *homeRepository) GetUserSessionHistory(ctx context.Context, userID int32, limit, offset int32) ([]domain.FocusSession, error) {
	query := `
		SELECT id, user_id, duration_seconds, elapsed_seconds, status, started_at, ended_at, created_at, updated_at
		FROM focus_sessions
		WHERE user_id = $1
		ORDER BY started_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []domain.FocusSession
	for rows.Next() {
		var s domain.FocusSession
		if err := rows.Scan(
			&s.ID, &s.UserID, &s.DurationSeconds, &s.ElapsedSeconds,
			&s.Status, &s.StartedAt, &s.EndedAt, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}
	return sessions, nil
}

func (r *homeRepository) scanSession(row pgx.Row) (*domain.FocusSession, error) {
	var s domain.FocusSession
	err := row.Scan(
		&s.ID, &s.UserID, &s.DurationSeconds, &s.ElapsedSeconds,
		&s.Status, &s.StartedAt, &s.EndedAt, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// ─────────────────────────────────────────
// Streaks
// ─────────────────────────────────────────

func (r *homeRepository) GetStreak(ctx context.Context, userID int32) (*domain.Streak, error) {
	query := `
		SELECT user_id, current_streak, longest_streak, last_session_date, updated_at
		FROM streaks WHERE user_id = $1
	`
	var st domain.Streak
	var lastDate *time.Time
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&st.UserID, &st.CurrentStreak, &st.LongestStreak, &lastDate, &st.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Return a default zero streak rather than an error
			return &domain.Streak{UserID: userID}, nil
		}
		return nil, err
	}
	st.LastSessionDate = lastDate
	return &st, nil
}

func (r *homeRepository) UpsertStreak(ctx context.Context, userID int32, currentStreak, longestStreak int32, lastDate *time.Time) (*domain.Streak, error) {
	query := `
		INSERT INTO streaks (user_id, current_streak, longest_streak, last_session_date, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (user_id) DO UPDATE SET
			current_streak   = EXCLUDED.current_streak,
			longest_streak   = GREATEST(streaks.longest_streak, EXCLUDED.longest_streak),
			last_session_date = EXCLUDED.last_session_date,
			updated_at       = NOW()
		RETURNING user_id, current_streak, longest_streak, last_session_date, updated_at
	`
	var st domain.Streak
	var ld *time.Time
	err := r.db.QueryRow(ctx, query, userID, currentStreak, longestStreak, lastDate).Scan(
		&st.UserID, &st.CurrentStreak, &st.LongestStreak, &ld, &st.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	st.LastSessionDate = ld
	return &st, nil
}

// ─────────────────────────────────────────
// Badges
// ─────────────────────────────────────────

func (r *homeRepository) GetAllBadgeDefinitions(ctx context.Context) ([]domain.BadgeDefinition, error) {
	query := `SELECT id, name, description, icon_key, criteria_type, criteria_value, created_at FROM badge_definitions ORDER BY id`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var badges []domain.BadgeDefinition
	for rows.Next() {
		var b domain.BadgeDefinition
		if err := rows.Scan(&b.ID, &b.Name, &b.Description, &b.IconKey, &b.CriteriaType, &b.CriteriaValue, &b.CreatedAt); err != nil {
			return nil, err
		}
		badges = append(badges, b)
	}
	return badges, nil
}

func (r *homeRepository) GetUserBadges(ctx context.Context, userID int32) ([]domain.UserBadge, error) {
	query := `
		SELECT bd.id, bd.name, bd.description, bd.icon_key, ub.earned_at
		FROM user_badges ub
		JOIN badge_definitions bd ON bd.id = ub.badge_id
		WHERE ub.user_id = $1
		ORDER BY ub.earned_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var badges []domain.UserBadge
	for rows.Next() {
		var b domain.UserBadge
		if err := rows.Scan(&b.BadgeID, &b.Name, &b.Description, &b.IconKey, &b.EarnedAt); err != nil {
			return nil, err
		}
		badges = append(badges, b)
	}
	return badges, nil
}

func (r *homeRepository) AwardBadge(ctx context.Context, userID, badgeID int32) error {
	query := `
		INSERT INTO user_badges (user_id, badge_id) VALUES ($1, $2)
		ON CONFLICT (user_id, badge_id) DO NOTHING
	`
	_, err := r.db.Exec(ctx, query, userID, badgeID)
	return err
}

func (r *homeRepository) GetUserBadgeCount(ctx context.Context, userID int32) (int32, error) {
	var count int32
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM user_badges WHERE user_id = $1`, userID).Scan(&count)
	return count, err
}

// ─────────────────────────────────────────
// Progress
// ─────────────────────────────────────────

func (r *homeRepository) UpsertDailyProgress(ctx context.Context, userID int32, date time.Time, focusSeconds, sessions int32) error {
	query := `
		INSERT INTO daily_progress (user_id, date, total_focus_seconds, sessions_completed)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, date) DO UPDATE SET
			total_focus_seconds = daily_progress.total_focus_seconds + EXCLUDED.total_focus_seconds,
			sessions_completed  = daily_progress.sessions_completed  + EXCLUDED.sessions_completed
	`
	_, err := r.db.Exec(ctx, query, userID, date, focusSeconds, sessions)
	if err != nil {
		fmt.Println("Error upserting daily progress:", err)
	}
	return err
}

func (r *homeRepository) GetProgressSummary(ctx context.Context, userID int32, from, to time.Time) (*domain.ProgressSummary, error) {
	query := `
		SELECT
			COALESCE(SUM(total_focus_seconds), 0) AS total_focus_seconds,
			COALESCE(SUM(sessions_completed), 0)  AS sessions_completed,
			COUNT(DISTINCT date)                   AS days_active
		FROM daily_progress
		WHERE user_id = $1 AND date BETWEEN $2 AND $3
	`
	var p domain.ProgressSummary
	err := r.db.QueryRow(ctx, query, userID, from, to).Scan(
		&p.TotalFocusSeconds, &p.SessionsCompleted, &p.DaysActive,
	)
	if err != nil {
		return nil, err
	}
	p.TotalFocusHours = float64(p.TotalFocusSeconds) / 3600.0
	return &p, nil
}
