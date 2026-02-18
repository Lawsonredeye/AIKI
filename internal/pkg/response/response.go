package response

import (
	"errors"
	"net/http"

	"aiki/internal/domain"

	"github.com/labstack/echo/v4"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Success sends a successful response
func Success(c echo.Context, statusCode int, message string, data interface{}) error {
	return c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error sends an error response
func Error(c echo.Context, err error) error {
	statusCode := domain.GetHTTPStatus(err)

	return c.JSON(statusCode, Response{
		Success: false,
		Error:   err.Error(),
	})
}

// ValidationError sends a validation error response
func ValidationError(c echo.Context, message string) error {
	return c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Error:   message,
	})
}

var (
	ErrSessionNotFound      = errors.New("focus session not found")
	ErrSessionAlreadyActive = errors.New("you already have an active focus session")
	ErrSessionNotActive     = errors.New("focus session is not active")
	ErrInvalidSessionStatus = errors.New("invalid session status transition")
	ErrStreakNotFound       = errors.New("streak record not found")
)
