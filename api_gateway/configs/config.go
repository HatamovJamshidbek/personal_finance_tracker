package configs

import (
	"fmt"
	"os"

	"github.com/spf13/cast"

	"github.com/joho/godotenv"
)

type Config struct {
	ServiceName string
	Environment string
	LoggerLevel string
	HTTPort     string

	BudgetServiceGrpcHost string
	BudgetServiceGrpcPort string

	UserServiceGrpcHost string
	UserServiceGrpcPort string
}

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env not found", err)
	}

	cfg := Config{}

	cfg.ServiceName = cast.ToString(getOrReturnDefault("SERVICE_NAME", "api_gateway"))
	// cfg.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "dev"))
	cfg.LoggerLevel = cast.ToString(getOrReturnDefault("LOGGER_LEVEL", "debug"))
	// HTTP PORT 8080
	cfg.HTTPort = cast.ToString(getOrReturnDefault("HTTP_PORT", ":8080"))
	// BUDGET SERVICE 8081
	cfg.BudgetServiceGrpcHost = cast.ToString(getOrReturnDefault("BUDGET_SERVICE_GRPC_HOST", "localhost"))
	cfg.BudgetServiceGrpcPort = cast.ToString(getOrReturnDefault("BUDGET_SERVICE_GRPC_PORT", ":8082"))
	// USER SERVICE 8082
	cfg.UserServiceGrpcHost = cast.ToString(getOrReturnDefault("USER_SERVICE_GRPC_HOST", "localhost"))
	cfg.UserServiceGrpcPort = cast.ToString(getOrReturnDefault("USER_SERVICE_GRPC_PORT", ":8081"))

	return cfg
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	value := os.Getenv(key)
	if value != "" {
		return value
	}

	return defaultValue
}
