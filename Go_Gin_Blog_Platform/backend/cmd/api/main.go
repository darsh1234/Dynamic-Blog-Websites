package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/darshvaidya/dynamic-blog-websites/go-gin-blog-platform/backend/internal/auth"
	"github.com/darshvaidya/dynamic-blog-websites/go-gin-blog-platform/backend/internal/config"
	"github.com/darshvaidya/dynamic-blog-websites/go-gin-blog-platform/backend/internal/db"
	"github.com/darshvaidya/dynamic-blog-websites/go-gin-blog-platform/backend/internal/email"
	"github.com/darshvaidya/dynamic-blog-websites/go-gin-blog-platform/backend/internal/logging"
	"github.com/darshvaidya/dynamic-blog-websites/go-gin-blog-platform/backend/internal/repository"
	"github.com/darshvaidya/dynamic-blog-websites/go-gin-blog-platform/backend/internal/service"
	httptransport "github.com/darshvaidya/dynamic-blog-websites/go-gin-blog-platform/backend/internal/transport/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	logger := logging.NewJSONLogger(cfg.AppEnv)
	logger.Info("starting api", "env", cfg.AppEnv, "port", cfg.Port, "variant", cfg.AppVariant)

	store, err := db.New(cfg.DatabaseURL, logger)
	if err != nil {
		panic(fmt.Errorf("failed to initialize database: %w", err))
	}
	defer func() {
		if err := store.Close(); err != nil {
			logger.Error("database close failed", "error", err)
		}
	}()

	tokenManager := auth.NewTokenManager(
		cfg.JWTAccessSecret,
		cfg.JWTRefreshSecret,
		time.Duration(cfg.JWTAccessTTLMinutes)*time.Minute,
		time.Duration(cfg.JWTRefreshTTLHours)*time.Hour,
	)

	emailSender := resolveEmailSender(cfg, logger)
	userRepo := repository.NewUserRepository(store.Gorm())
	postRepo := repository.NewPostRepository(store.Gorm())
	refreshRepo := repository.NewRefreshTokenRepository(store.Gorm())
	passwordResetRepo := repository.NewPasswordResetTokenRepository(store.Gorm())

	authService := service.NewAuthService(
		logger,
		userRepo,
		refreshRepo,
		passwordResetRepo,
		tokenManager,
		emailSender,
		time.Duration(cfg.PasswordResetTTLMinutes)*time.Minute,
		cfg.FrontendBaseURL,
	)
	authHandler := httptransport.NewAuthHandler(authService)
	postService := service.NewPostService(postRepo)
	postHandler := httptransport.NewPostHandler(postService)
	adminService := service.NewAdminService(userRepo)
	adminHandler := httptransport.NewAdminHandler(adminService)

	router := httptransport.NewRouter(logger, httptransport.RouterDependencies{
		HealthChecker:       store,
		HealthCheckTimeout:  time.Duration(cfg.RequestTimeoutS) * time.Second,
		AuthHandler:         authHandler,
		PostHandler:         postHandler,
		AdminHandler:        adminHandler,
		AccessTokenVerifier: tokenManager,
	})
	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	logger.Info("api listening", "addr", server.Addr)
	shutdownGracefully(server, logger)
}

func shutdownGracefully(server *http.Server, logger *slog.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutdown signal received")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}

	logger.Info("server stopped")
}

func resolveEmailSender(cfg config.Config, logger *slog.Logger) email.Sender {
	if cfg.EmailProvider == "ses" {
		return email.NewSESSender(email.SESConfig{
			Region:  cfg.AWSRegion,
			From:    cfg.EmailFrom,
			FromARN: cfg.AWSSESFromARN,
		})
	}

	return email.NewStubSender(logger)
}
