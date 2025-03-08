package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost		string
	DBPort		string
	DBUser		string
	DBPassword	string
	DBName		string
	ServerPort	string
	Env			string
	RWTimeout	time.Duration
	IdleTimeout	time.Duration
}

func LoadConfig() *Config {
    if err := godotenv.Load(); err != nil {
        log.Println("Warning: No .env file found, using system environment variables")
    }

    rwTimeout, err := time.ParseDuration(getEnv("RW_TIMEOUT", "5s"))
    if err != nil {
        log.Fatal("error parsing duration")
    }

    idleTimeout, err := time.ParseDuration(getEnv("IDLE_TIMEOUT", "30s"))
    if err != nil {
        log.Fatal("error parsing duration")
    }

    config := &Config{
        DBHost:      getEnv("POSTGRES_HOST", "localhost"),
        DBPort:      getEnv("POSTGRES_PORT", "5432"),
        DBUser:      getEnv("POSTGRES_USER", "user"),
        DBPassword:  getEnv("POSTGRES_PASSWORD", "secret"),
        DBName:      getEnv("POSTGRES_DB", "mydb"),
        ServerPort:  getEnv("SERVER_PORT", "8080"),
        RWTimeout:   rwTimeout,
        IdleTimeout: idleTimeout,
		Env: 		 getEnv("ENV", "local"),
    }
    return config
}

func getEnv(key, fallback string) string {
    if value:= os.Getenv(key); value != "" {
        return value
    }
    return fallback
}