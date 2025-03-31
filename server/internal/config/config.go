package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	CacheAddr 	string
	ServerPort  string
	Env         string
	RWTimeout   time.Duration
	IdleTimeout time.Duration
	LogHost     string
	LogPort     string
	SMTPConfig  SMTPConfig
	JWTSecret   string
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
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

	SMPTPort, err := strconv.Atoi(getEnv("SMPT_PORT", "587"))
	if err != nil {
		log.Fatal("smtp port must be a number")
	}

	SMTPPassword := os.Getenv("SMTP_PASSWORD")
	if SMTPPassword == "" {
		log.Fatal("no SMTP password provided")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("missing jwt secret")
	}

	config := &Config{
		DBHost:      getEnv("POSTGRES_HOST", "localhost"),
		DBPort:      getEnv("POSTGRES_PORT", "5432"),
		DBUser:      getEnv("POSTGRES_USER", "user"),
		DBPassword:  getEnv("POSTGRES_PASSWORD", "secret"),
		DBName:      getEnv("POSTGRES_DB", "mydb"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		CacheAddr:   getEnv("CACHE_ADDR", "localhost:6379"),
		RWTimeout:   rwTimeout,
		IdleTimeout: idleTimeout,
		Env:         getEnv("ENV", "local"),
		LogHost:     getEnv("LOGSTASH_HOST", "logstash"),
		LogPort:     getEnv("LOGSTASH_PORT", "5044"),
		JWTSecret:   jwtSecret,
		SMTPConfig: SMTPConfig{
			Host:     getEnv("SMTP_HOST", "smtp.google.com"),
			Port:     SMPTPort,
			Username: getEnv("SMTP_USERNAME", "zvukovat@gmail.com"),
			Password: SMTPPassword,
		},
	}
	return config
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
