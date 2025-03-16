package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Kazan-Strelnikova/SPDA/server/internal/config"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/http/events/create"
	createEnrollment "github.com/Kazan-Strelnikova/SPDA/server/internal/http/enrollments/create"
	deleteEnrollment "github.com/Kazan-Strelnikova/SPDA/server/internal/http/enrollments/delete"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/http/events/delete"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/http/events/get"
	getall "github.com/Kazan-Strelnikova/SPDA/server/internal/http/events/getAll"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/http/middleware/auth"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/http/ping"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/http/users/login"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/http/users/register"
	tokenlogin "github.com/Kazan-Strelnikova/SPDA/server/internal/http/users/token_login"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/log"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/service"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/storage/postgres"
	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.LoadConfig()

	log := log.SetupLogger(cfg.Env)

	log.Info("connecting to database")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	log.Debug("connection info", slog.String("db url", dbURL))
	storage := postgres.New(context.Background(), dbURL)
	defer storage.Close()

	log.Info("database connection established")

	//TODO:
	//Move the secret out
	service := service.New(log, storage, storage, "filler_secret")

	router := gin.Default()
	// TODO:
	// router.SetTrustedProxies()

	router.GET("/ping", ping.New(log))

	router.POST("/users/signup", register.New(log, service, cfg.RWTimeout))
	router.POST("/users/signin", login.New(log, service, cfg.RWTimeout))
	router.GET("/users/signin/cookie", tokenlogin.New(log, service, cfg.RWTimeout))

	router.POST("/events", create.New(log, service, cfg.RWTimeout))
	router.GET("/events", getall.New(log, service, cfg.RWTimeout))
	router.GET("/events/:event_id", get.New(log, service, cfg.RWTimeout))
	router.POST("/events/:event_id/enrollment", auth.New(log, service, cfg.RWTimeout), createEnrollment.New(log, service, cfg.RWTimeout))
	router.DELETE("/events/:event_id/enrollment", auth.New(log, service, cfg.RWTimeout), deleteEnrollment.New(log, service, cfg.RWTimeout))
	router.DELETE("/events/:event_id", auth.New(log, service, cfg.RWTimeout), delete.New(log, service, cfg.RWTimeout))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  cfg.RWTimeout,
		WriteTimeout: cfg.RWTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	log.Info("starting up server")
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("listen error", slog.String("err", err.Error()))
		}
	}()

	<-done
	log.Info("stopping the server")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.RWTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server forced to shutdown", slog.String("err", err.Error()))
	}

	log.Info("server exited properly")
}
