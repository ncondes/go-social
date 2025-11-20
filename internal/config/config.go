package config

import (
	"fmt"
	"time"

	"github.com/ncondes/go/social/internal/env"
)

type Config struct {
	Addr string
	DB   DBConfig
}

type DBConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
}

func Load() *Config {
	return &Config{
		Addr: fmt.Sprintf(":%s", env.GetString("PORT", "8080")),
		DB: DBConfig{
			Addr:         env.GetString("DB_ADDR", "postgres://postgres:password@localhost:5432/social?sslmode=disable"),
			MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 10),
			MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 5),
			MaxIdleTime:  env.GetDuration("DB_MAX_IDLE_TIME", 5*time.Minute),
		},
	}
}
