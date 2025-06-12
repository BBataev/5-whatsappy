package config

import "os"

type Config struct {
	ListenAddr string

	PostgreHost string
	PostgrePort string
	PostgreUser string
	PostgrePass string
	PostgreName string

	HTTPTimeout string
	JWToken     string
}

func Load() *Config {
	return &Config{
		ListenAddr:  getEnv("PORT", ""),
		PostgreHost: getEnv("PG_HOST", ""),
		PostgrePort: getEnv("PG_PORT", ""),
		PostgreUser: getEnv("PG_USER", ""),
		PostgrePass: getEnv("PG_PASS", ""),
		PostgreName: getEnv("PG_NAME", ""),
		HTTPTimeout: getEnv("HTTP_TIMEOUT", ""),
		JWToken:     getEnv("JWT_TOKEN", ""),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
