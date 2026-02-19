package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort         string
	PostgresHost    string
	PostgresPort    string
	PostgresUser    string
	PostgresPass    string
	PostgresDB      string
	RedisAddr       string
	RedisPassword   string
	RedisDB         string
	KafkaBrokers    string
	KafkaTopic      string
	MongoURI        string
	MongoDB         string
	MongoCollection string
	JWTSecret       string
	APIKey          string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		AppPort:         getEnv("APP_PORT", "8080"),
		PostgresHost:    getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:    getEnv("POSTGRES_PORT", "5432"),
		PostgresUser:    getEnv("POSTGRES_USER", "cars"),
		PostgresPass:    getEnv("POSTGRES_PASSWORD", "cars"),
		PostgresDB:      getEnv("POSTGRES_DB", "cars"),
		RedisAddr:       getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:   getEnv("REDIS_PASSWORD", ""),
		RedisDB:         getEnv("REDIS_DB", "0"),
		KafkaBrokers:    getEnv("KAFKA_BROKERS", "localhost:9092"),
		KafkaTopic:      getEnv("KAFKA_TOPIC", "car-api-logs"),
		MongoURI:        getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:         getEnv("MONGO_DB", "cars_logs"),
		MongoCollection: getEnv("MONGO_COLLECTION", "request_logs"),
		JWTSecret:       getEnv("JWT_SECRET", "super-secret-change-me"),
		APIKey:          getEnv("API_KEY", "my-api-key-12345"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
