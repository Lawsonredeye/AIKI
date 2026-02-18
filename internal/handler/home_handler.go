package handler

import (
	"aiki/internal/domain"
	"aiki/internal/pkg/response"
	"aiki/internal/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type HomeHandler struct {
	homeService service.HomeService
	validator   echo.Validator
}

func NewHomeHandler(homeService service.HomeService, validator echo.Validator) *HomeHandler {
	return &HomeHandler{
		homeService: homeService,
		validator:   validator,
	}
}

// ─────────────────────────────────────────
// Home Screen
// ─────────────────────────────────────────

// GetHomeScreen godoc
// @Summary Get home screen data
// @Description Returns aggregated home screen data: streak, active session, weekly progress, recent badges
// @Tags home
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=domain.HomeScreenData}
// @Failure 401 {object} response.Response
// @Router /home [get]
func (h *HomeHandler) GetHomeScreen(c echo.Context) error {
	userID, ok := c.Get("user_id").(int32)
	if !ok {
		return response.Error(c, domain.ErrUnauthorized)
	}

	data, err := h.homeService.GetHomeScreenData(c.Request().Context(), userID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusOK, "home screen data retrieved successfully", data)
}

// ─────────────────────────────────────────
// Focus Sessions
// ─────────────────────────────────────────

// StartSession godoc
// @Summary Start a focus session
// @Description Starts a new lock-in focus session for the authenticated user
// @Tags sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.StartSessionRequest true "Session details"
// @Success 201 {object} response.Response{data=domain.FocusSession}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 409 {object} response.Response "Session already active"
// @Router /sessions [post]
func (h *HomeHandler) StartSession(c echo.Context) error {
	userID, ok := c.Get("user_id").(int32)
	if !ok {
		return response.Error(c, domain.ErrUnauthorized)
	}

	var req domain.StartSessionRequest
	if err := c.Bind(&req); err != nil {
		return response.ValidationError(c, "invalid request body")
	}
	if err := h.validator.Validate(&req); err != nil {
		return response.ValidationError(c, err.Error())
	}

	session, err := h.homeService.StartSession(c.Request().Context(), userID, &req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusCreated, "focus session started", session)
}

// GetActiveSession godoc
// @Summary Get active session
// @Description Returns the current active or paused session, if any
// @Tags sessions
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=domain.FocusSession}
// @Failure 401 {object} response.Response
// @Router /sessions/active [get]
func (h *HomeHandler) GetActiveSession(c echo.Context) error {
	userID, ok := c.Get("user_id").(int32)
	if !ok {
		return response.Error(c, domain.ErrUnauthorized)
	}

	session, err := h.homeService.GetActiveSession(c.Request().Context(), userID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusOK, "active session retrieved", session)
}

// PauseSession godoc
// @Summary Pause a focus session
// @Description Pauses an active focus session
// @Tags sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Session ID"
// @Param request body domain.UpdateSessionRequest true "Elapsed seconds"
// @Success 200 {object} response.Response{data=domain.FocusSession}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /sessions/{id}/pause [patch]
func (h *HomeHandler) PauseSession(c echo.Context) error {
	userID, ok := c.Get("user_id").(int32)
	if !ok {
		return response.Error(c, domain.ErrUnauthorized)
	}

	sessionID, err := parseIDParam(c, "id")
	if err != nil {
		return response.ValidationError(c, "invalid session id")
	}

	var req domain.UpdateSessionRequest
	if err := c.Bind(&req); err != nil {
		return response.ValidationError(c, "invalid request body")
	}

	session, err := h.homeService.PauseSession(c.Request().Context(), userID, sessionID, req.ElapsedSeconds)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusOK, "session paused", session)
}

// ResumeSession godoc
// @Summary Resume a paused session
// @Tags sessions
// @Produce json
// @Security BearerAuth
// @Param id path int true "Session ID"
// @Success 200 {object} response.Response{data=domain.FocusSession}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /sessions/{id}/resume [patch]
func (h *HomeHandler) ResumeSession(c echo.Context) error {
	userID, ok := c.Get("user_id").(int32)
	if !ok {
		return response.Error(c, domain.ErrUnauthorized)
	}

	sessionID, err := parseIDParam(c, "id")
	if err != nil {
		return response.ValidationError(c, "invalid session id")
	}

	session, err := h.homeService.ResumeSession(c.Request().Context(), userID, sessionID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusOK, "session resumed", session)
}

// EndSession godoc
// @Summary End a focus session
// @Description Ends a session as completed or abandoned
// @Tags sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Session ID"
// @Param request body domain.UpdateSessionRequest true "Elapsed seconds and status"
// @Success 200 {object} response.Response{data=domain.FocusSession}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /sessions/{id}/end [patch]
func (h *HomeHandler) EndSession(c echo.Context) error {
	userID, ok := c.Get("user_id").(int32)
	if !ok {
		return response.Error(c, domain.ErrUnauthorized)
	}

	sessionID, err := parseIDParam(c, "id")
	if err != nil {
		return response.ValidationError(c, "invalid session id")
	}

	var req domain.UpdateSessionRequest
	if err := c.Bind(&req); err != nil {
		return response.ValidationError(c, "invalid request body")
	}
	if err := h.validator.Validate(&req); err != nil {
		return response.ValidationError(c, err.Error())
	}

	completed := req.Status == domain.SessionStatusCompleted
	session, err := h.homeService.EndSession(c.Request().Context(), userID, sessionID, req.ElapsedSeconds, completed)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusOK, "session ended", session)
}

// GetSessionHistory godoc
// @Summary Get session history
// @Tags sessions
// @Produce json
// @Security BearerAuth
// @Param limit  query int false "Limit (default 20)"
// @Param offset query int false "Offset (default 0)"
// @Success 200 {object} response.Response{data=[]domain.FocusSession}
// @Failure 401 {object} response.Response
// @Router /sessions [get]
func (h *HomeHandler) GetSessionHistory(c echo.Context) error {
	userID, ok := c.Get("user_id").(int32)
	if !ok {
		return response.Error(c, domain.ErrUnauthorized)
	}

	limit := int32(20)
	offset := int32(0)
	if l := c.QueryParam("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil {
			limit = int32(v)
		}
	}
	if o := c.QueryParam("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil {
			offset = int32(v)
		}
	}

	sessions, err := h.homeService.GetSessionHistory(c.Request().Context(), userID, limit, offset)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusOK, "session history retrieved", sessions)
}

// ─────────────────────────────────────────
// Streaks
// ─────────────────────────────────────────

// GetStreak godoc
// @Summary Get user streak
// @Tags streaks
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=domain.Streak}
// @Failure 401 {object} response.Response
// @Router /streaks [get]
func (h *HomeHandler) GetStreak(c echo.Context) error {
	userID, ok := c.Get("user_id").(int32)
	if !ok {
		return response.Error(c, domain.ErrUnauthorized)
	}

	streak, err := h.homeService.GetStreak(c.Request().Context(), userID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusOK, "streak retrieved", streak)
}

// ─────────────────────────────────────────
// Badges
// ─────────────────────────────────────────

// GetUserBadges godoc
// @Summary Get user's earned badges
// @Tags badges
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]domain.UserBadge}
// @Failure 401 {object} response.Response
// @Router /badges/me [get]
func (h *HomeHandler) GetUserBadges(c echo.Context) error {
	userID, ok := c.Get("user_id").(int32)
	if !ok {
		return response.Error(c, domain.ErrUnauthorized)
	}

	badges, err := h.homeService.GetUserBadges(c.Request().Context(), userID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusOK, "badges retrieved", badges)
}

// GetAllBadges godoc
// @Summary Get all badge definitions
// @Tags badges
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]domain.BadgeDefinition}
// @Router /badges [get]
func (h *HomeHandler) GetAllBadges(c echo.Context) error {
	badges, err := h.homeService.GetAllBadges(c.Request().Context())
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusOK, "all badges retrieved", badges)
}

// ─────────────────────────────────────────
// Progress
// ─────────────────────────────────────────

// GetProgress godoc
// @Summary Get progress summary
// @Description Returns focus time stats for the given period (weekly | monthly | yearly)
// @Tags progress
// @Produce json
// @Security BearerAuth
// @Param period query string false "Period: weekly (default), monthly, yearly"
// @Success 200 {object} response.Response{data=domain.ProgressSummary}
// @Failure 401 {object} response.Response
// @Router /progress [get]
func (h *HomeHandler) GetProgress(c echo.Context) error {
	userID, ok := c.Get("user_id").(int32)
	if !ok {
		return response.Error(c, domain.ErrUnauthorized)
	}

	period := c.QueryParam("period")
	if period == "" {
		period = "weekly"
	}

	summary, err := h.homeService.GetProgressSummary(c.Request().Context(), userID, period)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusOK, "progress retrieved", summary)
}

// ─────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────

func parseIDParam(c echo.Context, param string) (int32, error) {
	v, err := strconv.Atoi(c.Param(param))
	if err != nil {
		return 0, err
	}
	return int32(v), nil
}
