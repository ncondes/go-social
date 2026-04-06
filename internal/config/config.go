package config

import (
	"fmt"
	"time"

	"github.com/ncondes/go/social/internal/env"
)

type Config struct {
	Addr        string
	DB          DBConfig
	Env         string
	APIBaseURL  string
	MailConfig  MailConfig
	FrontendURL string
}

type DBConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
}

type MailConfig struct {
	FromEmail string
	APIKey    string
	Exp       time.Duration
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
		Env:         env.GetString("ENV", "development"),
		APIBaseURL:  env.GetString("API_BASE_URL", "localhost:8080"),
		FrontendURL: env.GetString("FRONTEND_URL", "http://localhost:5173"),
		MailConfig: MailConfig{
			FromEmail: env.GetString("MAIL_FROM_EMAIL", "noreply@example.com"),
			APIKey:    env.GetString("SENDGRID_API_KEY", ""),
			Exp:       env.GetDuration("MAIL_EXPIRATION_TIME", 24*time.Hour),
		},
	}
}
