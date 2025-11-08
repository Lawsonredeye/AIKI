package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"aiki/internal/config"
	"aiki/internal/database"
	"aiki/internal/handler"
	"aiki/internal/middleware"
	"aiki/internal/pkg/jwt"
	"aiki/internal/pkg/validator"
	"aiki/internal/repository"
	"aiki/internal/router"
	"aiki/internal/service"

	"github.com/labstack/echo/v4"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	db, err := database.NewPostgresPool(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("✓ Database connection established")

	// Initialize Redis client
	redis, err := database.NewRedisClient(&cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redis.Close()
	log.Println("✓ Redis connection established")

	// Initialize JWT manager
	jwtManager := jwt.NewManager(
		cfg.JWT.Secret,
		cfg.JWT.AccessExpiry,
		cfg.JWT.RefreshExpiry,
	)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, jwtManager)
	userService := service.NewUserService(userRepo)

	// Initialize Echo
	e := echo.New()
	e.HideBanner = true
	e.Validator = validator.New()

	// Global middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService, e.Validator)
	userHandler := handler.NewUserHandler(userService, e.Validator)

	// Setup routes
	router.Setup(e, authHandler, userHandler, jwtManager)

	// Start server
	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	go func() {
		log.Printf("✓ Server starting on %s", serverAddr)
		if err := e.Start(serverAddr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
