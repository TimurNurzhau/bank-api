package config

import (
	"os"
)

type Config struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBSSLMode     string
	JWTSecret     string
	SMTPHost      string
	SMTPPort      string
	SMTPUser      string
	SMTPPass      string
	PGPPublicKey  string // Добавлено
	PGPPrivateKey string // Добавлено
	HMACSecret    string
	ServerPort    string
	LogLevel      string
}

func Load() *Config {
	return &Config{
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "bankuser"),
		DBPassword:    getEnv("DB_PASSWORD", "bankuser_pass_2024"),
		DBName:        getEnv("DB_NAME", "bankapi"),
		DBSSLMode:     getEnv("DB_SSLMODE", "disable"),
		JWTSecret:     getEnv("JWT_SECRET", "default-secret"),
		SMTPHost:      getEnv("SMTP_HOST", "localhost"),
		SMTPPort:      getEnv("SMTP_PORT", "587"),
		SMTPUser:      getEnv("SMTP_USER", ""),
		SMTPPass:      getEnv("SMTP_PASS", ""),
		PGPPublicKey:  getEnv("PGP_PUBLIC_KEY", ""),   // Добавлено
		PGPPrivateKey: getEnv("PGP_PRIVATE_KEY", ""),  // Добавлено
		HMACSecret:    getEnv("HMAC_SECRET", "hmac-secret"),
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		LogLevel:      getEnv("LOG_LEVEL", "debug"),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}