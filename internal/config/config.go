package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Database   DatabaseConfig
	TimescaleDB TimescaleConfig
	JWT        JWTConfig
	Server     ServerConfig
	MQTT       MQTTConfig
	Redis      RedisConfig
	S3         S3Config
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type TimescaleConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

type ServerConfig struct {
	Host string
	Port int
	Mode string
}

type MQTTConfig struct {
	Broker       string
	ClientID     string
	Username     string
	Password     string
	TopicSensor  string
	TopicControl string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type S3Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	Region    string
}

func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	config := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "swiflet_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		TimescaleDB: TimescaleConfig{
			Host:     getEnv("TIMESCALE_HOST", "localhost"),
			Port:     getEnvAsInt("TIMESCALE_PORT", 5432),
			User:     getEnv("TIMESCALE_USER", "postgres"),
			Password: getEnv("TIMESCALE_PASSWORD", "password"),
			DBName:   getEnv("TIMESCALE_DB", "swiflet_timeseries"),
			SSLMode:  getEnv("TIMESCALE_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "default-secret-change-this"),
			Expiry: getEnvAsDuration("JWT_EXPIRY", 24*time.Hour),
		},
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnvAsInt("SERVER_PORT", 8080),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		MQTT: MQTTConfig{
			Broker:       getEnv("MQTT_BROKER", "tcp://localhost:1883"),
			ClientID:     getEnv("MQTT_CLIENT_ID", "swiflet-backend"),
			Username:     getEnv("MQTT_USERNAME", ""),
			Password:     getEnv("MQTT_PASSWORD", ""),
			TopicSensor:  getEnv("MQTT_TOPIC_SENSOR", "sensors/+/data"),
			TopicControl: getEnv("MQTT_TOPIC_CONTROL", "control/+/command"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		S3: S3Config{
			Endpoint:  getEnv("S3_ENDPOINT", ""),
			AccessKey: getEnv("S3_ACCESS_KEY", ""),
			SecretKey: getEnv("S3_SECRET_KEY", ""),
			Bucket:    getEnv("S3_BUCKET", "swiftlead-storage"),
			Region:    getEnv("S3_REGION", "us-east-1"),
		},
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}