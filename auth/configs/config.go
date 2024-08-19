package configs

import (
	"fmt"
	"os"

	"github.com/spf13/cast"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string

	ServiceName string
	Environment string
	LoggerLevel string

	AuthServiceGrpcHost string
	AuthServiceGrpcPort string
	EmailPassword       string
}

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env not found", err)
	}

	cfg := Config{}

	cfg.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST", "postgres-db"))
	cfg.PostgresPort = cast.ToString(getOrReturnDefault("POSTGRES_PORT", "5433"))
	cfg.PostgresUser = cast.ToString(getOrReturnDefault("POSTGRES_USER", "postgres"))
	cfg.PostgresPassword = cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD", "1111"))
	cfg.PostgresDB = cast.ToString(getOrReturnDefault("POSTGRES_DB", "user"))
	cfg.EmailPassword = cast.ToString(getOrReturnDefault("EMAIL_PASSWORD", "user"))

	cfg.ServiceName = cast.ToString(getOrReturnDefault("SERVICE_NAME", "auth_service"))
	// cfg.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "dev"))
	cfg.LoggerLevel = cast.ToString(getOrReturnDefault("LOGGER_LEVEL", "debug"))

	cfg.AuthServiceGrpcHost = cast.ToString(getOrReturnDefault("AUTH_SERVICE_GRPC_HOST", "auth-services1"))
	cfg.AuthServiceGrpcPort = cast.ToString(getOrReturnDefault("AUTH_SERVICE_GRPC_PORT", ":8081"))

	return cfg
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	value := os.Getenv(key)
	if value != "" {
		return value
	}

	return defaultValue
}
