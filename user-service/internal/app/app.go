package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10/translations/en"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/simonaditia/nyayurin/user-service/config"
	"github.com/simonaditia/nyayurin/user-service/internal/adapter/handler"
	"github.com/simonaditia/nyayurin/user-service/internal/adapter/repository"
	"github.com/simonaditia/nyayurin/user-service/internal/core/service"
	"github.com/simonaditia/nyayurin/user-service/utils/validator"
)

func RunServer() {
	cfg := config.NewConfig()
	db, err := cfg.ConnectionPostgres()
	if err != nil {
		log.Fatalf("[RunServer-1] Failed to connect to database: %v", err)
		return
	}

	userRepo := repository.NewUserRepository(db.DB)
	jwtService := service.NewJwtService(cfg)
	userService := service.NewUserService(userRepo, cfg, jwtService)

	e := echo.New()
	e.Use(middleware.CORS())
	// e.Group("/api/v1")

	customValidator := validator.NewValidator()
	en.RegisterDefaultTranslations(customValidator.Validator, customValidator.Translator)
	e.Validator = customValidator

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	handler.NewUserHandler(e, userService, cfg)

	go func() {
		if cfg.App.AppPort == "" {
			cfg.App.AppPort = os.Getenv("APP_PORT")
		}

		err := e.Start(":" + cfg.App.AppPort)
		if err != nil {
			log.Fatalf("[RunServer-2] Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)
	<-quit

	log.Print("[RunServer-3] Shutting down server of 5 seconds...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	e.Shutdown(ctx)
}
