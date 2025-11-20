package env

import (
	"os"
	"strconv"
	"time"
)

func GetString(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}

	return value
}

func GetInt(key string, fallback int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}

	valueAsInt, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return valueAsInt
}

func GetDuration(key string, fallback time.Duration) time.Duration {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}

	valueAsDuration, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}

	return valueAsDuration
}
