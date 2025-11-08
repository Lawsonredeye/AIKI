package middleware

import (
	"strings"

	"aiki/internal/domain"
	"aiki/internal/pkg/jwt"
	"aiki/internal/pkg/response"

	"github.com/labstack/echo/v4"
)

// Auth creates a JWT authentication middleware
func Auth(jwtManager *jwt.Manager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return response.Error(c, domain.ErrUnauthorized)
			}

			// Check if it's a Bearer token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return response.Error(c, domain.ErrUnauthorized)
			}

			tokenString := parts[1]

			// Validate token
			claims, err := jwtManager.ValidateToken(tokenString)
			if err != nil {
				return response.Error(c, domain.ErrInvalidToken)
			}

			// Set user information in context
			c.Set("user_id", claims.UserID)
			c.Set("user_email", claims.Email)

			return next(c)
		}
	}
}
