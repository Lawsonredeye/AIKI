package service

import (
	"aiki/internal/domain"
	"aiki/internal/pkg/response"
	"aiki/internal/repository"
	"context"
	"time"
)

//go:generate mockgen -source=home_service.go -destination=mocks/mock_home_service.go -package=mocks

type HomeService interface {
	// Home screen aggregated data
	GetHomeScreenData(ctx context.Context, userID int32) (*domain.HomeScreenData, error)

	// Focus sessions
	StartSession(ctx context.Context, userID int32, req *domain.StartSessionRequest) (*domain.FocusSession, error)
	PauseSession(ctx context.Context, userID int32, sessionID int32, elapsed int32) (*domain.FocusSession, error)
	ResumeSession(ctx context.Context, userID int32, sessionID int32) (*domain.FocusSession, error)
	EndSession(ctx context.Context, userID int32, sessionID int32, elapsed int32, completed bool) (*domain.FocusSession, error)
	GetActiveSession(ctx context.Context, userID int32) (*domain.FocusSession, error)
	GetSessionHistory(ctx context.Context, userID int32, limit, offset int32) ([]domain.FocusSession, error)

	// Streak
	GetStreak(ctx context.Context, userID int32) (*domain.Streak, error)

	// Badges
	GetUserBadges(ctx context.Context, userID int32) ([]domain.UserBadge, error)
	GetAllBadges(ctx context.Context) ([]domain.BadgeDefinition, error)

	// Progress
	GetProgressSummary(ctx context.Context, userID int32, period string) (*domain.ProgressSummary, error)
}

type homeService struct {
	homeRepo repository.HomeRepository
}

func NewHomeService(homeRepo repository.HomeRepository) HomeService {
	return &homeService{homeRepo: homeRepo}
}

// ─────────────────────────────────────────
// Home screen aggregated data
// ─────────────────────────────────────────

func (s *homeService) GetHomeScreenData(ctx context.Context, userID int32) (*domain.HomeScreenData, error) {
	streak, err := s.homeRepo.GetStreak(ctx, userID)
	if err != nil {
		return nil, err
	}

	activeSession, err := s.homeRepo.GetActiveSession(ctx, userID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	weekStart := now.AddDate(0, 0, -int(now.Weekday()))
	weeklyProgress, err := s.homeRepo.GetProgressSummary(ctx, userID, weekStart, now)
	if err != nil {
		return nil, err
	}
	weeklyProgress.Period = "weekly"

	badges, err := s.homeRepo.GetUserBadges(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Show only the 6 most recent badges on home screen
	recentBadges := badges
	if len(recentBadges) > 6 {
		recentBadges = recentBadges[:6]
	}

	totalBadges, err := s.homeRepo.GetUserBadgeCount(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &domain.HomeScreenData{
		Streak:         streak,
		ActiveSession:  activeSession,
		WeeklyProgress: *weeklyProgress,
		RecentBadges:   recentBadges,
		TotalBadges:    totalBadges,
	}, nil
}

// ─────────────────────────────────────────
// Focus Sessions
// ─────────────────────────────────────────

func (s *homeService) StartSession(ctx context.Context, userID int32, req *domain.StartSessionRequest) (*domain.FocusSession, error) {
	// Check for existing active/paused session
	existing, err := s.homeRepo.GetActiveSession(ctx, userID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, response.ErrSessionAlreadyActive
	}

	return s.homeRepo.CreateSession(ctx, userID, req.DurationSeconds)
}

func (s *homeService) PauseSession(ctx context.Context, userID int32, sessionID int32, elapsed int32) (*domain.FocusSession, error) {
	session, err := s.homeRepo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session.UserID != userID {
		return nil, domain.ErrUnauthorized
	}
	if session.Status != domain.SessionStatusActive {
		return nil, response.ErrInvalidSessionStatus
	}
	return s.homeRepo.UpdateSession(ctx, sessionID, elapsed, domain.SessionStatusPaused, nil)
}

func (s *homeService) ResumeSession(ctx context.Context, userID int32, sessionID int32) (*domain.FocusSession, error) {
	session, err := s.homeRepo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session.UserID != userID {
		return nil, domain.ErrUnauthorized
	}
	if session.Status != domain.SessionStatusPaused {
		return nil, response.ErrInvalidSessionStatus
	}
	return s.homeRepo.UpdateSession(ctx, sessionID, session.ElapsedSeconds, domain.SessionStatusActive, nil)
}

func (s *homeService) EndSession(ctx context.Context, userID int32, sessionID int32, elapsed int32, completed bool) (*domain.FocusSession, error) {
	session, err := s.homeRepo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session.UserID != userID {
		return nil, domain.ErrUnauthorized
	}

	status := domain.SessionStatusAbandoned
	if completed {
		status = domain.SessionStatusCompleted
	}

	now := time.Now()
	updatedSession, err := s.homeRepo.UpdateSession(ctx, sessionID, elapsed, status, &now)
	if err != nil {
		return nil, err
	}

	// Post-session side effects (streak + progress + badges) only on completion
	if completed {
		s.handleSessionCompleted(ctx, userID, elapsed, now)
	}

	return updatedSession, nil
}

// handleSessionCompleted updates streak, daily progress, and checks for new badges.
// Errors here are non-fatal — we log them but don't fail the request.
func (s *homeService) handleSessionCompleted(ctx context.Context, userID int32, focusSeconds int32, completedAt time.Time) {
	// Update daily progress
	today := time.Date(completedAt.Year(), completedAt.Month(), completedAt.Day(), 0, 0, 0, 0, completedAt.Location())
	_ = s.homeRepo.UpsertDailyProgress(ctx, userID, today, focusSeconds, 1)

	// Update streak
	streak, err := s.homeRepo.GetStreak(ctx, userID)
	if err == nil {
		newStreak := s.calculateNewStreak(streak, today)
		_, _ = s.homeRepo.UpsertStreak(ctx, userID, newStreak, max32(newStreak, streak.LongestStreak), &today)
		streak.CurrentStreak = newStreak // for badge check below
	}

	// Check and award badges
	_ = s.checkAndAwardBadges(ctx, userID, streak)
}

func (s *homeService) calculateNewStreak(streak *domain.Streak, today time.Time) int32 {
	if streak.LastSessionDate == nil {
		return 1
	}

	lastDate := *streak.LastSessionDate
	lastDate = time.Date(lastDate.Year(), lastDate.Month(), lastDate.Day(), 0, 0, 0, 0, lastDate.Location())
	diff := int(today.Sub(lastDate).Hours() / 24)

	switch diff {
	case 0:
		// Same day — streak unchanged
		return streak.CurrentStreak
	case 1:
		// Consecutive day — increment
		return streak.CurrentStreak + 1
	default:
		// Streak broken
		return 1
	}
}

func (s *homeService) checkAndAwardBadges(ctx context.Context, userID int32, streak *domain.Streak) error {
	definitions, err := s.homeRepo.GetAllBadgeDefinitions(ctx)
	if err != nil {
		return err
	}

	// Fetch aggregate stats for other criteria types
	now := time.Now()
	allTime, _ := s.homeRepo.GetProgressSummary(ctx, userID, time.Time{}, now)

	earnedBadges, _ := s.homeRepo.GetUserBadges(ctx, userID)
	earnedMap := make(map[int32]bool, len(earnedBadges))
	for _, b := range earnedBadges {
		earnedMap[b.BadgeID] = true
	}

	for _, def := range definitions {
		if earnedMap[def.ID] {
			continue
		}

		var earned bool
		switch def.CriteriaType {
		case "streak":
			if streak != nil && streak.CurrentStreak >= def.CriteriaValue {
				earned = true
			}
		case "sessions":
			if allTime != nil && int32(allTime.SessionsCompleted) >= def.CriteriaValue {
				earned = true
			}
		case "focus_time":
			if allTime != nil && int32(allTime.TotalFocusSeconds) >= def.CriteriaValue {
				earned = true
			}
		}

		if earned {
			_ = s.homeRepo.AwardBadge(ctx, userID, def.ID)
		}
	}

	return nil
}

func (s *homeService) GetActiveSession(ctx context.Context, userID int32) (*domain.FocusSession, error) {
	return s.homeRepo.GetActiveSession(ctx, userID)
}

func (s *homeService) GetSessionHistory(ctx context.Context, userID int32, limit, offset int32) ([]domain.FocusSession, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return s.homeRepo.GetUserSessionHistory(ctx, userID, limit, offset)
}

// ─────────────────────────────────────────
// Streaks
// ─────────────────────────────────────────

func (s *homeService) GetStreak(ctx context.Context, userID int32) (*domain.Streak, error) {
	return s.homeRepo.GetStreak(ctx, userID)
}

// ─────────────────────────────────────────
// Badges
// ─────────────────────────────────────────

func (s *homeService) GetUserBadges(ctx context.Context, userID int32) ([]domain.UserBadge, error) {
	return s.homeRepo.GetUserBadges(ctx, userID)
}

func (s *homeService) GetAllBadges(ctx context.Context) ([]domain.BadgeDefinition, error) {
	return s.homeRepo.GetAllBadgeDefinitions(ctx)
}

// ─────────────────────────────────────────
// Progress
// ─────────────────────────────────────────

func (s *homeService) GetProgressSummary(ctx context.Context, userID int32, period string) (*domain.ProgressSummary, error) {
	now := time.Now()
	var from time.Time

	switch period {
	case "monthly":
		from = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	case "yearly":
		from = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	default: // weekly
		period = "weekly"
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		from = now.AddDate(0, 0, -(weekday - 1))
		from = time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, from.Location())
	}

	summary, err := s.homeRepo.GetProgressSummary(ctx, userID, from, now)
	if err != nil {
		return nil, err
	}
	summary.Period = period
	return summary, nil
}

// ─────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
