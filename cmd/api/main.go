package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/dppi/dppierp-api/internal/config"
	"github.com/dppi/dppierp-api/internal/handler"
	"github.com/dppi/dppierp-api/internal/middleware"
	"github.com/dppi/dppierp-api/internal/repository"
	"github.com/dppi/dppierp-api/internal/service"
	"github.com/dppi/dppierp-api/pkg/database"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	// Panic recovery for startup crashes
	defer func() {
		if r := recover(); r != nil {
			log.Error().Interface("panic", r).Msg("Application panicked during startup")
			// Print stack trace for systemd
			fmt.Fprintf(os.Stderr, "PANIC: %v\n", r)
			panic(r) // Re-throw to ensure exit code is still non-zero (wait, panic exits with 2 anyway)
		}
	}()

	// Setup logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().
		Str("version", Version).
		Str("build_time", BuildTime).
		Str("git_commit", GitCommit).
		Msg("Starting DPPI ERP API")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Set Gin mode
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to database
	db, err := database.NewMySQLConnection(&cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	log.Info().Msg("Connected to database successfully")

	// Initialize repositories
	fabricRepo := repository.NewFabricRepository(db)
	rackRepo := repository.NewRackRepository(db)
	userRepo := repository.NewUserRepository(db)
	masterRepo := repository.NewMasterRepository(db)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWT.Secret)

	// Initialize services
	checkpointService := service.NewCheckpointService(fabricRepo, rackRepo)
	authService := service.NewAuthService(userRepo, authMiddleware)
	masterService := service.NewMasterService(masterRepo)

	// Initialize handlers
	checkpointHandler := handler.NewCheckpointHandler(checkpointService)
	authHandler := handler.NewAuthHandler(authService)
	masterHandler := handler.NewMasterHandler(masterService)

	// Setup router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORSMiddleware(cfg.CORS.AllowedOrigins))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Auth routes (public)
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/token/refresh", authHandler.RefreshToken)
		authGroup.POST("/forgot-password/request", authHandler.ForgotPasswordRequest)
		authGroup.POST("/forgot-password/reset", authHandler.ResetPassword)
	}

	// Auth routes (protected)
	authProtected := router.Group("/auth")
	authProtected.Use(authMiddleware.Authenticate())
	{
		authProtected.GET("/me", authHandler.Me)
		authProtected.POST("/logout", authHandler.Logout)
	}

	// Profile routes (protected)
	profileGroup := router.Group("/profile")
	profileGroup.Use(authMiddleware.Authenticate())
	{
		profileGroup.GET("", authHandler.Me)
		profileGroup.POST("/change-password", authHandler.ChangePassword)
	}

	// Check Point routes (protected)
	checkpointGroup := router.Group("/check-point/v1")
	checkpointGroup.Use(authMiddleware.Authenticate())
	{
		checkpointGroup.GET("/overview", checkpointHandler.GetOverview)
		checkpointGroup.POST("/scan", checkpointHandler.ScanQR)
		checkpointGroup.POST("/move", checkpointHandler.MoveStage)
		checkpointGroup.POST("/scan-rack", checkpointHandler.ScanRack)
		checkpointGroup.POST("/relocation", checkpointHandler.Relocate)
	}

	// Master Data routes (protected)
	masterGroup := router.Group("/master")
	masterGroup.Use(authMiddleware.Authenticate())
	{
		masterGroup.GET("/blocks", masterHandler.GetBlocks)
		masterGroup.GET("/racks", masterHandler.GetRacks)
		masterGroup.GET("/relaxation-blocks", masterHandler.GetRelaxationBlocks)
		masterGroup.GET("/relaxation-racks", masterHandler.GetRelaxationRacks)
	}

	// Create server
	srv := &http.Server{
		Addr: "0.0.0.0:" + cfg.App.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Info().Str("port", cfg.App.Port).Msg("Starting server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	fmt.Println("Server exited properly")
}
