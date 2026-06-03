package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort      string
	GinMode      string
	TemporalHost string
	TaskQueue    string
	DBConn       string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("💡 Note: Not found .env")
	}

	return &Config{
		AppPort:      getEnv("PORT", "8081"),
		GinMode:      getEnv("GIN_MODE", "debug"),
		TemporalHost: getEnv("TEMPORAL_HOST", "localhost:7233"),
		TaskQueue:    getEnv("TASK_QUEUE", "loan-task-queue"),
		DBConn:       getEnv("DB_CONN", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
