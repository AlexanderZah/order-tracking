package config

import (
	"log"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort   string `env:"HTTP_PORT" env-default:"8083"`
	DBURL      string `env:"DB_URL" env-default:"postgres://user:password@postgres:5432/orderdb?sslmode=disable"`
	KafkaAddr  string `env:"KAFKA_ADDR" env-default:"kafka:9092"`
	ETAService string `env:"ETA_SERVICE" env-default:"eta-service:50051"`
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default vars")
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
