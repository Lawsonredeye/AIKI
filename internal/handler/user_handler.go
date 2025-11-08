package handler

import (
	"net/http"

	"aiki/internal/domain"
	"aiki/internal/pkg/response"
	"aiki/internal/service"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService service.UserService
	validator   echo.Validator
}

func NewUserHandler(userService service.UserService, validator echo.Validator) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   validator,
	}
}

// GetMe godoc
// @Summary Get current user
// @Description Get the currently authenticated user's profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=domain.User}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /users/me [get]
func (h *UserHandler) GetMe(c echo.Context) error {
	// Get user ID from JWT claims (set by auth middleware)
	userID, ok := c.Get("user_id").(int32)
	if !ok {
		return response.Error(c, domain.ErrUnauthorized)
	}

	user, err := h.userService.GetByID(c.Request().Context(), userID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusOK, "user retrieved successfully", user)
}

// UpdateMe godoc
// @Summary Update current user
// @Description Update the currently authenticated user's profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.UpdateUserRequest true "Update user details"
// @Success 200 {object} response.Response{data=domain.User}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /users/me [put]
func (h *UserHandler) UpdateMe(c echo.Context) error {
	// Get user ID from JWT claims
	userID, ok := c.Get("user_id").(int32)
	if !ok {
		return response.Error(c, domain.ErrUnauthorized)
	}

	var req domain.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return response.ValidationError(c, "invalid request body")
	}

	if err := h.validator.Validate(&req); err != nil {
		return response.ValidationError(c, err.Error())
	}

	user, err := h.userService.Update(c.Request().Context(), userID, &req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusOK, "user updated successfully", user)
}
