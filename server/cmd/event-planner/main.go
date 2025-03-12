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
	"github.com/Kazan-Strelnikova/SPDA/server/internal/http/events/delete"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/http/events/get"
	getall "github.com/Kazan-Strelnikova/SPDA/server/internal/http/events/getAll"
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
	router.DELETE("/events/:event_id", func(c *gin.Context) {
		log.Info("attempt of cookie login through middleware")

		token, err := c.Cookie("token")
		if err != nil {
			log.Error("Missing authentication token", "error", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), cfg.RWTimeout)
		defer cancel()

		usr, err := service.LoginByToken(ctx, token)
		if err != nil {
			log.Error("Invalid authentication token", "error", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		log.Info("cookie login successful", slog.String("email", usr.Email))

		c.Set("email", usr.Email)

		c.Next()
	}, delete.New(log, service, cfg.RWTimeout))

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
