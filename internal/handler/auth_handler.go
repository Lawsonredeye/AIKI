package handler

import (
	"net/http"

	"aiki/internal/domain"
	"aiki/internal/pkg/response"
	"aiki/internal/service"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService service.AuthService
	validator   echo.Validator
}

func NewAuthHandler(authService service.AuthService, validator echo.Validator) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.RegisterRequest true "Registration details"
// @Success 201 {object} response.Response{data=domain.AuthResponse}
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	var req domain.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return response.ValidationError(c, "invalid request body")
	}

	if err := h.validator.Validate(&req); err != nil {
		return response.ValidationError(c, err.Error())
	}

	resp, err := h.authService.Register(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusCreated, "user registered successfully", resp)
}

// Login godoc
// @Summary Login
// @Description Authenticate user and return tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.LoginRequest true "Login credentials"
// @Success 200 {object} response.Response{data=domain.AuthResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req domain.LoginRequest
	if err := c.Bind(&req); err != nil {
		return response.ValidationError(c, "invalid request body")
	}

	if err := h.validator.Validate(&req); err != nil {
		return response.ValidationError(c, err.Error())
	}

	resp, err := h.authService.Login(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusOK, "login successful", resp)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get a new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} response.Response{data=domain.AuthResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var req domain.RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return response.ValidationError(c, "invalid request body")
	}

	if err := h.validator.Validate(&req); err != nil {
		return response.ValidationError(c, err.Error())
	}

	resp, err := h.authService.RefreshToken(c.Request().Context(), req.RefreshToken)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusOK, "token refreshed successfully", resp)
}

// Logout godoc
// @Summary Logout
// @Description Invalidate refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c echo.Context) error {
	var req domain.RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return response.ValidationError(c, "invalid request body")
	}

	if err := h.validator.Validate(&req); err != nil {
		return response.ValidationError(c, err.Error())
	}

	if err := h.authService.Logout(c.Request().Context(), req.RefreshToken); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, http.StatusOK, "logout successful", nil)
}
