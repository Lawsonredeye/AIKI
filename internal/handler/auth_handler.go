package handler

import (
	"aiki/internal/pkg/otp_token"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"aiki/internal/config"
	"aiki/internal/domain"
	"aiki/internal/pkg/response"
	"aiki/internal/service"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type AuthHandler struct {
	authService service.AuthService
	validator   echo.Validator
	config      config.Config
	redis       *redis.Client
}

const (
	ExpireInMinute = 5 * time.Minute
	ExpireInHours  = 1 * time.Hour
)

// NewAuthHandler Update the constructor to accept config
func NewAuthHandler(authService service.AuthService, validator echo.Validator, redis *redis.Client) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator,
		redis:       redis,
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

// LinkedInLogin godoc
// @Summary Initiate LinkedIn login
// @Description Redirects user to LinkedIn for authentication
// @Tags auth
// @Produce html
// @Success 302 "Redirects to LinkedIn"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /auth/linkedin/login [get]
func (h *AuthHandler) LinkedInLogin(c echo.Context) error {
	// Generate a random state to protect against CSRF attacks
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return response.Error(c, domain.ErrInternalServer)
	}
	state := base64.URLEncoding.EncodeToString(b)

	// Set the state in a cookie for later validation
	cookie := &http.Cookie{
		Name:     "linkedin_oauth_state",
		Value:    state,
		Path:     "/",
		Expires:  time.Now().Add(5 * time.Minute), // State valid for 5 minutes
		HttpOnly: true,
		Secure:   h.config.Server.Env == "production", // Only secure in production
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)

	// Construct the LinkedIn authorization URL
	// Note: LinkedIn's OAuth 2.0 endpoints typically start with https://www.linkedin.com/oauth/v2/
	// You need to register your redirect URI with LinkedIn.
	// For this example, let's assume the redirect URI is /auth/linkedin/callback on your server.
	linkedInAuthURL := "https://www.linkedin.com/oauth/v2/authorization"
	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", h.config.LinkedIn.ClientID)
	params.Add("redirect_uri", fmt.Sprintf(h.config.LinkedIn.ClientCallbackUrl, h.config.Server.Port))
	params.Add("state", state)
	params.Add("scope", "r_liteprofile r_emailaddress") // Request basic profile and email

	fullURL := linkedInAuthURL + "?" + params.Encode()

	return c.Redirect(http.StatusFound, fullURL)
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

func (h *AuthHandler) ForgottenPassword(c echo.Context) error {
	var req domain.ForgotPasswordRequest
	if err := c.Bind(&req); err != nil {
		return response.ValidationError(c, "invalid request body")
	}

	if err := h.validator.Validate(&req); err != nil {
		return response.ValidationError(c, err.Error())
	}
	if _, err := h.authService.ForgottenPassword(c.Request().Context(), &req); err != nil {
		return response.Error(c, err)
	}
	// using the session_id to keep track of users movement and the otp sent
	c.Logger().Info("Generating session id for tracking")
	sessionId := uuid.New()
	data := make(map[string]string)
	data["session_id"] = sessionId.String()
	ctx := context.Background()
	c.Logger().Info("Generating otp token")
	otp, err := otp_token.OTPGenerator()
	if err != nil {
		c.Logger().Errorf("failed to generate otp token: %v", err)
		return response.Error(c, domain.ErrInternalServer)
	}
	value := map[string]string{
		"otp":        otp,
		"email":      req.Email,
		"created_at": time.Now().Format(time.RFC3339),
	}
	bValue, err := json.Marshal(&value)
	if err != nil {
		c.Logger().Errorf("failed to marshal data: %v", err)
	}
	if os.Getenv("environment") != "production" {
		data["otp"] = otp
	}
	c.Logger().Info("Storing session id for tracking")
	key := sessionKey(sessionId.String())
	if err := h.redis.SetEx(ctx, key, bValue, ExpireInMinute).Err(); err != nil {
		c.Logger().Errorf("failed to store session id for tracking: %v", err)
	}
	return response.Success(c, http.StatusOK, "if the email exists, a password reset link has been sent", data)
}

func (h *AuthHandler) ValidateForgottenPasswordOTP(c echo.Context) error {
	var req domain.ValidateForgotPasswordOTP
	var err error
	if err = c.Bind(&req); err != nil {
		return response.ValidationError(c, "invalid request body")
	}
	if err := h.validator.Validate(&req); err != nil {
		return response.ValidationError(c, err.Error())
	}
	// get otp from redis and compare user input
	ctx := context.Background()
	key := sessionKey(req.SessionId)
	resp := h.redis.Get(ctx, key)
	if resp.Err() != nil {
		c.Logger().Errorf("Failed to find session id for current user: %v", resp)
		return response.Error(c, errors.New("invalid OTP token"))
	}
	c.Logger().Infof("Getting session id for tracking user session: %v", resp.Val())
	jsonData, err := resp.Bytes()
	if err != nil {
		c.Logger().Errorf("Failed to get json data: %v", err)
		return response.Error(c, errors.New("something went wrong please try again"))
	}

	var value map[string]string
	if err := json.Unmarshal(jsonData, &value); err != nil {
		c.Logger().Errorf("failed to unmarshal session data: %v", err)
		return response.Error(c, errors.New("something went wrong please try again"))
	}
	fmt.Printf("Validating session id for tracking user session: %v", resp.Val())
	fmt.Printf("redis cache value: %v, session from request: %v\n", value["session_id"], req.SessionId)
	if value["otp"] != req.Otp {
		c.Logger().Error("Invalid otp for tracking user session")
		return response.Error(c, errors.New("invalid OTP token"))
	}
	c.Logger().Info("user validation successful")
	value["isValid"] = "true"
	jBytes, err := json.Marshal(value)
	if err != nil {
		c.Logger().Errorf("failed to marshal data: %v", err)
		return response.Error(c, errors.New("something went wrong please try again"))
	}
	if err := h.redis.Set(ctx, key, jBytes, ExpireInMinute).Err(); err != nil {
		c.Logger().Errorf("failed to store session id for tracking: %v", err)
		return response.Error(c, errors.New("something went wrong please try again"))
	}
	return response.Success(c, http.StatusOK, "user validation successful", value)
}

func (h *AuthHandler) ResetPassword(c echo.Context) error {
	var req domain.ResetPasswordRequest
	if err := c.Bind(&req); err != nil {
		return response.ValidationError(c, "invalid request body")
	}

	if err := h.validator.Validate(&req); err != nil {
		return response.ValidationError(c, err.Error())
	}
	ctx := context.Background()
	key := sessionKey(req.SessionId)
	resp := h.redis.Get(ctx, key)
	if resp.Err() != nil {
		c.Logger().Errorf("Failed to find session id for current user: %v", resp)
		return response.Error(c, errors.New("invalid operation step"))
	}
	jsonData, err := resp.Bytes()
	if err != nil {
		c.Logger().Errorf("Failed to get json data: %v", err)
		return response.Error(c, errors.New("something went wrong please try again"))
	}
	value := make(map[string]string)
	if err := json.Unmarshal(jsonData, &value); err != nil {
		c.Logger().Errorf("failed to unmarshal session data: %v", err)
		return response.Error(c, errors.New("something went wrong please try again"))
	}
	if value["isValid"] != "true" {
		c.Logger().Error("Invalid user operational step")
		return response.Error(c, errors.New("unauthorized request access"))
	}

	if err := h.authService.ResetPassword(c.Request().Context(), &req); err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, http.StatusOK, "password has been reset successfully", nil)
}

func sessionKey(key string) string {
	return fmt.Sprintf("forgotten-password-%s", key)
}
