package main

import (
	"context"
	"fmt"

	"github.com/Kazan-Strelnikova/SPDA/server/internal/config"
	"github.com/Kazan-Strelnikova/SPDA/server/internal/storage/postgres"
	_ "github.com/gin-gonic/gin"
)

func main() {

	cfg := config.LoadConfig()
		
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	storage := postgres.New(context.Background(), dbURL)

	_ = storage
	// router := gin.Default()
	
}
