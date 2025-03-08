package main

import (
	"context"
	"fmt"

	"github.com/Kazan-Strelnikova/SPDA/server/internal/config"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/http/ping"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/log"
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

	router := gin.Default()
	// TODO:
	// router.SetTrustedProxies()

	router.GET("/ping", ping.New(log))

	log.Info("starting up server")
	router.Run(":" + cfg.ServerPort)
}
