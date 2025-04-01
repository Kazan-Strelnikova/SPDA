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
	createEnrollment "github.com/Kazan-Strelnikova/SPDA/server/internal/http/enrollments/create"
	deleteEnrollment "github.com/Kazan-Strelnikova/SPDA/server/internal/http/enrollments/delete"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/http/events/create"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/http/events/delete"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/http/events/get"
	getall "github.com/Kazan-Strelnikova/SPDA/server/internal/http/events/getAll"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/http/events/put"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/http/middleware/auth"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/http/ping"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/http/users/login"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/http/users/register"
	tokenlogin "github.com/Kazan-Strelnikova/SPDA/server/internal/http/users/token_login"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/log"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/service"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/storage/postgres"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/storage/redis"
	"github.com/gin-gonic/gin"
	"go.elastic.co/apm/module/apmgin/v2"
	"go.elastic.co/apm/v2"
)

func main() {

	cfg := config.LoadConfig()

	log := log.SetupLogger(cfg.Env, cfg.LogHost, cfg.LogPort)

	log.Info("connecting to database")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	// log.Debug("connection info", slog.String("db url", dbURL))
	storage := postgres.New(context.Background(), dbURL)
	defer storage.Close()

	log.Debug(cfg.CacheAddr)
	cache, err := redis.New(cfg.CacheAddr)
	if err != nil {
		log.Error("error connecting to cache. proceeding without it", slog.String("err", err.Error()))
		cache = nil
	}

	log.Info("database connection established")

	service := service.New(log, storage, storage, cfg.JWTSecret, cfg.SMTPConfig, cache)

	router := gin.Default()

	tracer, err := apm.NewTracerOptions(apm.TracerOptions{
		ServiceName:    "Event-planner",
		ServiceVersion: "1.0",
	})

	if err == nil {
		router.Use(apmgin.Middleware(router, apmgin.WithTracer(tracer)))
	} else {
		log.Warn("tracer initialization error", slog.String("err", err.Error()))
	}

	router.GET("/ping", ping.New(log))

	router.POST("/users/signup", register.New(log, service, cfg.RWTimeout))
	router.POST("/users/signin", login.New(log, service, cfg.RWTimeout))
	router.GET("/users/signin/cookie", tokenlogin.New(log, service, cfg.RWTimeout))

	//TODO:
	// add login middleware and compare emails
	router.POST("/events", create.New(log, service, cfg.RWTimeout))
	router.GET("/events", getall.New(log, service, cfg.RWTimeout))
	router.GET("/events/:event_id", get.New(log, service, cfg.RWTimeout))
	router.PUT("/events/:event_id", put.New(log, service, cfg.RWTimeout))
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

	go service.Monitor(done)
	<-done
	close(done)
	log.Info("stopping the server")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.RWTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server forced to shutdown", slog.String("err", err.Error()))
	}

	log.Info("server exited properly")
}
