package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort     string
	DatabaseURL string
}

func Load() *Config {
	_ = godotenv.Load() // ignore error if no .env

	cfg := &Config{
		AppPort:     getEnv("PORT", "8080"),
		DatabaseURL: mustGetEnv("DATABASE_URL"),
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("%s is required", key)
	}
	return v
}