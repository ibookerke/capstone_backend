package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/ibookerke/capstone_backend/internal/delivery/http"
	"github.com/ibookerke/capstone_backend/internal/repository"
	"github.com/ibookerke/capstone_backend/internal/service"

	"github.com/ibookerke/capstone_backend/internal/config"
	"github.com/ibookerke/capstone_backend/internal/pkg/pgx"
	"github.com/ibookerke/capstone_backend/internal/pkg/trm/manager"
	"github.com/ibookerke/capstone_backend/internal/server"
)

// main is the entry point of the application
//
//	@title						Product API [capstone_backend]
//	@version					1.0
//	@description				This is swagger documentation for Product API
//	@termsOfService				http://swagger.io/terms/
//	@contact.name				API Support
//	@contact.url				http://www.swagger.io/support
//	@securityDefinitions.apikey	Bearer
//	@in							header
//	@name						Authorization
func main() {
	/**
	|
	| Doing some basic initialization
	|--------------------------------------------------------------------------
	*/
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	conf, err := config.Get()
	if err != nil {
		slog.Error("couldn't get config", "err", err)
		return
	}

	slogHandler := slog.Handler(slog.NewTextHandler(os.Stdout, nil))
	if !conf.Project.Debug {
		slogHandler = slog.NewJSONHandler(os.Stdout, nil)
	}
	logger := slog.New(slogHandler).With("svc", conf.Project.ServiceName)
	logger.Info("Starting the application")

	pool, err := pgx.NewPgxPool(ctx, conf.Database.DSN)
	if err != nil {
		logger.Error("couldn't create pgx pool", "err", err)
		return
	}
	defer pool.Close()

	// migrating database scheme using migrate library
	m, err := migrate.New("file://migrations", conf.Database.DSN)
	if err != nil {
		logger.Error("couldn't create migrate instance", "err", err)
		return
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Error("couldn't migrate database", "err", err)
		return
	}

	trManager := manager.Must(pgx.NewDefaultFactory(pool))
	userRepository := repository.NewUserRepository(pool, pgx.DefaultCtxGetter, trManager)
	notificationSettingsRepository := repository.NewNotificationSettingsRepository(pool, pgx.DefaultCtxGetter, trManager)

	authService := service.NewAuthService(trManager, userRepository, conf.Auth)
	userService := service.NewUserService(trManager, userRepository, notificationSettingsRepository, conf.Auth)
	notificationSettingsService := service.NewNotificationSettingsService(trManager, notificationSettingsRepository, conf.Auth)
	/**
	|
	|
	| Initializing the http server and router
	|--------------------------------------------------------------------------
	|
	*/

	router := gin.Default()

	http.NewAuthHandler(logger, authService, conf.Auth).RegisterAuthRoutes(router)
	http.NewUserHandler(logger, userService, conf.Auth).RegisterUserRoutes(router)
	http.NewNotificationSettingsHandler(logger, notificationSettingsService, conf.Auth).RegisterNotificationSettingsRoutes(router)

	httpServer := server.NewHTTPServer(logger, conf.Rest, router)
	go func() {
		if err := httpServer.Run(); err != nil {
			logger.Error("failed to start http server", "err", err)
			cancelFn()
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case v := <-quit:
		logger.Info("received exit signal", "signal", v)
		cancelFn()
	case <-ctx.Done():
		logger.Info("context done")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err = httpServer.Shutdown(ctx); err != nil {
		logger.Error("failed to shutdown server", "err", err)
	}
	logger.Info("server stopped")
}
