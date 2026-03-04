package common

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// ServiceConfig holds configuration common to all TeslaPay microservices.
type ServiceConfig struct {
	ServiceName string
	Environment string
	Port        int
	LogLevel    string

	// Database
	DatabaseURL string

	// Redis
	RedisURL      string
	RedisPassword string

	// Kafka
	KafkaBrokers []string

	// JWT
	JWTPublicKeyPath  string
	JWTPrivateKeyPath string
	JWTIssuer         string
	AccessTokenTTL    time.Duration
	RefreshTokenTTL   time.Duration

	// Observability
	TraceEndpoint string
	TraceSampling float64
}

// LoadBaseConfig reads common environment variables shared by all services.
func LoadBaseConfig(serviceName string) (*ServiceConfig, error) {
	cfg := &ServiceConfig{
		ServiceName: serviceName,
		Environment: getEnvOrDefault("ENVIRONMENT", "development"),
		LogLevel:    getEnvOrDefault("LOG_LEVEL", "info"),
		JWTIssuer:   getEnvOrDefault("JWT_ISSUER", "teslapay.eu"),
	}

	port, err := strconv.Atoi(getEnvOrDefault("PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid PORT: %w", err)
	}
	cfg.Port = port

	cfg.DatabaseURL = os.Getenv("DATABASE_URL")
	cfg.RedisURL = getEnvOrDefault("REDIS_URL", "localhost:6379")
	cfg.RedisPassword = os.Getenv("REDIS_PASSWORD")

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers != "" {
		cfg.KafkaBrokers = splitAndTrim(kafkaBrokers, ",")
	} else {
		cfg.KafkaBrokers = []string{"localhost:9092"}
	}

	cfg.JWTPublicKeyPath = os.Getenv("JWT_PUBLIC_KEY_PATH")
	cfg.JWTPrivateKeyPath = os.Getenv("JWT_PRIVATE_KEY_PATH")

	accessTTL, err := time.ParseDuration(getEnvOrDefault("ACCESS_TOKEN_TTL", "15m"))
	if err != nil {
		return nil, fmt.Errorf("invalid ACCESS_TOKEN_TTL: %w", err)
	}
	cfg.AccessTokenTTL = accessTTL

	refreshTTL, err := time.ParseDuration(getEnvOrDefault("REFRESH_TOKEN_TTL", "720h"))
	if err != nil {
		return nil, fmt.Errorf("invalid REFRESH_TOKEN_TTL: %w", err)
	}
	cfg.RefreshTokenTTL = refreshTTL

	cfg.TraceEndpoint = os.Getenv("TRACE_ENDPOINT")
	sampling, _ := strconv.ParseFloat(getEnvOrDefault("TRACE_SAMPLING", "0.1"), 64)
	cfg.TraceSampling = sampling

	return cfg, nil
}

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func splitAndTrim(s, sep string) []string {
	var result []string
	for _, part := range splitString(s, sep) {
		trimmed := trimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func splitString(s, sep string) []string {
	var parts []string
	start := 0
	for i := 0; i <= len(s)-len(sep); i++ {
		if s[i:i+len(sep)] == sep {
			parts = append(parts, s[start:i])
			start = i + len(sep)
		}
	}
	parts = append(parts, s[start:])
	return parts
}

func trimSpace(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n') {
		end--
	}
	return s[start:end]
}
