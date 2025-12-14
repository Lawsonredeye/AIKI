package router

import (
	"aiki/internal/handler"
	"aiki/internal/middleware"
	"aiki/internal/pkg/jwt"

	"github.com/labstack/echo/v4"
)

// Setup configures all routes for the application
func Setup(
	e *echo.Echo,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	jwtManager *jwt.Manager,
) {
	// API v1 group
	api := e.Group("/api/v1")

	// Health check (no auth required)
	api.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status":  "ok",
			"service": "aiki-api",
		})
	})

	// Auth routes (no auth required)
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", authHandler.Logout)
		auth.GET("/linkedin/login", authHandler.LinkedInLogin) // New LinkedIn login route
		// auth.GET("/linkedin/callback", authHandler.LinkedInCallback) // TODO: add New LinkedIn callback route
		auth.POST("/forgot-password", authHandler.ForgottenPassword)
		auth.POST("/forgot-password/validate", authHandler.ValidateForgottenPasswordOTP)
		auth.POST("/reset-password", authHandler.ResetPassword)
	}

	// User routes (auth required)
	users := api.Group("/users")
	users.Use(middleware.Auth(jwtManager))
	{
		users.GET("/me", userHandler.GetMe)
		users.PUT("/me", userHandler.UpdateMe)
		users.POST("/profile", userHandler.CreateProfile)
		users.PATCH("/profile", userHandler.UpdateProfile)
		users.POST("/upload/cv", userHandler.UploadCV)
	}
}
