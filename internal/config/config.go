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
	RateLimit   RateLimitConfig
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

type RateLimitTier struct {
	RequestsPerWindow int
	Window            time.Duration
}

type RateLimitConfig struct {
	Enabled bool
	// Server protection
	Global RateLimitTier
	// IP-based (anonymous users)
	StrictIP   RateLimitTier
	ModerateIP RateLimitTier
	// Operation-based (authenticated users)
	ReadOps  RateLimitTier
	WriteOps RateLimitTier
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

		RateLimit: RateLimitConfig{
			Enabled: env.GetBool("RATE_LIMIT_ENABLED", true),

			Global: RateLimitTier{
				RequestsPerWindow: env.GetInt("RATE_LIMIT_GLOBAL_REQUESTS", 10000),
				Window:            env.GetDuration("RATE_LIMIT_GLOBAL_WINDOW", 1*time.Minute),
			},

			StrictIP: RateLimitTier{
				RequestsPerWindow: env.GetInt("RATE_LIMIT_STRICT_IP_REQUESTS", 5),
				Window:            env.GetDuration("RATE_LIMIT_STRICT_IP_WINDOW", 1*time.Hour),
			},

			ModerateIP: RateLimitTier{
				RequestsPerWindow: env.GetInt("RATE_LIMIT_MODERATE_IP_REQUESTS", 100),
				Window:            env.GetDuration("RATE_LIMIT_MODERATE_IP_WINDOW", 1*time.Minute),
			},

			// Read operations - GET requests (generous)
			ReadOps: RateLimitTier{
				RequestsPerWindow: env.GetInt("RATE_LIMIT_READ_REQUESTS", 300),
				Window:            env.GetDuration("RATE_LIMIT_READ_WINDOW", 1*time.Minute),
			},

			// Write operations - POST/PUT/PATCH/DELETE (moderate)
			WriteOps: RateLimitTier{
				RequestsPerWindow: env.GetInt("RATE_LIMIT_WRITE_REQUESTS", 30),
				Window:            env.GetDuration("RATE_LIMIT_WRITE_WINDOW", 1*time.Minute),
			},
		},
	}
}
