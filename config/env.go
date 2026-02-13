package config

import (
	"os"
	"strings"
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetCORSOrigins() []string {
	origins := GetEnv("CORS_ORIGINS", "http://localhost:3000,http://localhost:3002")
	return strings.Split(origins, ",")
}

func GetPort() string {
	return ":" + GetEnv("PORT", "8080")
}

func GetDatabaseURL() string {
	return GetEnv("DATABASE_URL", "")
}
