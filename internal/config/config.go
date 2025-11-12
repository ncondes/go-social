package config

import (
	"fmt"

	"github.com/ncondes/go/social/internal/env"
)

type Config struct {
	Addr string
}

func Load() *Config {
	return &Config{
		Addr: fmt.Sprintf(":%s", env.GetString("PORT", "8080")),
	}
}
