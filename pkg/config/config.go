package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config holds all configuration for our application
type Config struct {
	ServiceName string
	Environment string
	LogLevel    string

	// Server configuration
	ServerPort int
	GrpcPort   int

	// Database configuration
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	MongoURI   string
	MongoDB    string

	// Kafka configuration
	KafkaBrokers []string `mapstructure:"KAFKA_BROKERS"`
	KafkaGroupID  string   `mapstructure:"KAFKA_GROUP_ID"`

	// JWT configuration
	JWTSecret string
	JWTExpiry int

	// Service discovery
	ServiceDiscoveryURL string
}

// LoadConfig loads configuration from environment variables and .env file
func LoadConfig(serviceName string) (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Initialize with default values
	config := &Config{
		ServiceName: serviceName,
		Environment: getEnv("ENVIRONMENT", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),

		// Server configuration
		ServerPort: getEnvAsInt("SERVICE_PORT", 8080),
		GrpcPort:   getEnvAsInt("GRPC_PORT", 9090),

		// Database configuration
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnvAsInt("DB_PORT", 5432),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "ecommerce"),
		MongoURI:   getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:    getEnv("MONGO_DB", "ecommerce"),

		// Kafka configuration
		KafkaBrokers: strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ","),
		KafkaGroupID: getEnv("KAFKA_GROUP_ID", "my-group"),

		// JWT configuration
		JWTSecret: getEnv("JWT_SECRET", "your_jwt_secret_key"),
		JWTExpiry: getEnvAsInt("JWT_EXPIRY", 24*60*60), // 24 hours in seconds

		// Service discovery
		ServiceDiscoveryURL: getEnv("SERVICE_DISCOVERY_URL", "http://localhost:8500"),
	}

	// Also load from config file if available
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/app")
	viper.AddConfigPath("$HOME/.app")

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())

		// Override with values from config file
		if viper.GetString("environment") != "" {
			config.Environment = viper.GetString("environment")
		}
		if viper.GetString("logLevel") != "" {
			config.LogLevel = viper.GetString("logLevel")
		}
		// Add more overrides as needed
	}

	return config, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsInt gets an environment variable as an integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}
