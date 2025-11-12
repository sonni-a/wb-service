package config

import (
	"log"
	"os"
)

type Config struct {
	PostgresURL  string
	KafkaBrokers string
}

func Load() *Config {
	cfg := &Config{
		PostgresURL:  getEnv("DATABASE_URL", "postgres://postgres:postgres@db:5432/demo_service"),
		KafkaBrokers: getEnv("KAFKA_BROKERS", "kafka:9092"),
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	log.Printf("ENV %s not set, using default: %s", key, defaultValue)
	return defaultValue
}
