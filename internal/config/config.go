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
	Auth        AuthConfig
	Redis       RedisConfig
}

type AuthConfig struct {
	Basic BasicAuthConfig
	JWT   JWTConfig
}

type JWTConfig struct {
	Secret   string
	Issuer   string
	Audience string
	Duration time.Duration
}

type BasicAuthConfig struct {
	Username string
	Password string
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

type RedisConfig struct {
	Enabled  bool
	Addr     string
	Password string
	DB       int
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
		Auth: AuthConfig{
			Basic: BasicAuthConfig{
				Username: env.GetString("BASIC_AUTH_USERNAME", "admin"),
				Password: env.GetString("BASIC_AUTH_PASSWORD", "password"),
			},
			JWT: JWTConfig{
				Secret:   env.GetString("JWT_SECRET", ""),
				Issuer:   env.GetString("JWT_ISSUER", "social"),
				Audience: env.GetString("JWT_AUDIENCE", "social"),
				Duration: env.GetDuration("JWT_DURATION", 24*time.Hour),
			},
		},
		Redis: RedisConfig{
			Enabled:  env.GetBool("REDIS_ENABLED", false),
			Addr:     env.GetString("REDIS_ADDR", "localhost:6379"),
			Password: env.GetString("REDIS_PASSWORD", ""),
			DB:       env.GetInt("REDIS_DB", 0),
		},
	}
}
