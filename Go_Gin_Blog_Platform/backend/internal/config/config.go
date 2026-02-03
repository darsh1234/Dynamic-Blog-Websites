package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	AppEnv                  string
	Port                    string
	DatabaseURL             string
	JWTAccessSecret         string
	JWTRefreshSecret        string
	JWTAccessTTLMinutes     int
	JWTRefreshTTLHours      int
	PasswordResetTTLMinutes int
	EmailProvider           string
	EmailFrom               string
	AWSRegion               string
	AWSSESFromARN           string
	CORSOrigins             []string
	AppVariant              string
	FrontendBaseURL         string
	RequestTimeoutS         int
}

func Load() (Config, error) {
	cfg := Config{
		AppEnv:                  getEnv("APP_ENV", "local"),
		Port:                    getEnv("PORT", "8080"),
		DatabaseURL:             getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/blog_platform?sslmode=disable"),
		JWTAccessSecret:         getEnv("JWT_ACCESS_SECRET", "local-access-secret-change-me"),
		JWTRefreshSecret:        getEnv("JWT_REFRESH_SECRET", "local-refresh-secret-change-me"),
		JWTAccessTTLMinutes:     getEnvInt("JWT_ACCESS_TTL_MINUTES", 15),
		JWTRefreshTTLHours:      getEnvInt("JWT_REFRESH_TTL_HOURS", 168),
		PasswordResetTTLMinutes: getEnvInt("PASSWORD_RESET_TTL_MINUTES", 30),
		EmailProvider:           getEnv("EMAIL_PROVIDER", "stub"),
		EmailFrom:               getEnv("EMAIL_FROM", "no-reply@localhost"),
		AWSRegion:               getEnv("AWS_REGION", "us-east-1"),
		AWSSESFromARN:           getEnv("AWS_SES_FROM_ARN", ""),
		CORSOrigins:             splitCSV(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173")),
		AppVariant:              getEnv("APP_VARIANT", "blog_a"),
		FrontendBaseURL:         getEnv("FRONTEND_BASE_URL", "http://localhost:5173"),
		RequestTimeoutS:         getEnvInt("REQUEST_TIMEOUT_SECONDS", 10),
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) Validate() error {
	if c.Port == "" {
		return fmt.Errorf("PORT cannot be empty")
	}

	if c.JWTAccessSecret == "" || c.JWTRefreshSecret == "" {
		return fmt.Errorf("JWT_ACCESS_SECRET and JWT_REFRESH_SECRET are required")
	}

	if c.EmailProvider != "stub" && c.EmailProvider != "ses" {
		return fmt.Errorf("EMAIL_PROVIDER must be 'stub' or 'ses'")
	}

	if c.RequestTimeoutS <= 0 {
		return fmt.Errorf("REQUEST_TIMEOUT_SECONDS must be > 0")
	}

	if c.JWTAccessTTLMinutes <= 0 || c.JWTRefreshTTLHours <= 0 {
		return fmt.Errorf("JWT_ACCESS_TTL_MINUTES and JWT_REFRESH_TTL_HOURS must be > 0")
	}

	if c.PasswordResetTTLMinutes <= 0 {
		return fmt.Errorf("PASSWORD_RESET_TTL_MINUTES must be > 0")
	}

	return nil
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvInt(key string, fallback int) int {
	raw := getEnv(key, "")
	if raw == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return parsed
}

func splitCSV(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		value := strings.TrimSpace(p)
		if value == "" {
			continue
		}
		out = append(out, value)
	}
	return out
}
