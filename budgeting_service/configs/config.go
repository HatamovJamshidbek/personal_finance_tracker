package configs

import (
	"fmt"
	"os"

	"github.com/spf13/cast"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoHost     string
	MongoPort     string
	MongoUser     string
	MongoPassword string
	MongoDB       string

	ServiceName string
	Environment string
	LoggerLevel string

	BudgetServiceGrpcHost string
	BudgetServiceGrpcPort string
	EmailPassword         string
}

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env not found", err)
	}

	cfg := Config{}

	cfg.MongoHost = cast.ToString(getOrReturnDefault("MONGO_HOST", "mongo"))
	cfg.MongoPort = cast.ToString(getOrReturnDefault("MONGO_PORT", "27017"))
	cfg.MongoUser = cast.ToString(getOrReturnDefault("MONGO_USER", "mongosh"))
	cfg.MongoPassword = cast.ToString(getOrReturnDefault("MONGO_PASSWORD", "1111"))
	cfg.MongoDB = cast.ToString(getOrReturnDefault("MONGO_DB", "budget_service"))

	cfg.ServiceName = cast.ToString(getOrReturnDefault("SERVICE_NAME", "budget"))
	// cfg.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "dev"))
	cfg.LoggerLevel = cast.ToString(getOrReturnDefault("LOGGER_LEVEL", "debug"))

	cfg.BudgetServiceGrpcHost = cast.ToString(getOrReturnDefault("BUDGET_SERVICE_GRPC_HOST", "budget_docker"))
	cfg.BudgetServiceGrpcPort = cast.ToString(getOrReturnDefault("BUDGET_SERVICE_GRPC_PORT", ":8082"))

	return cfg
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	value := os.Getenv(key)
	if value != "" {
		return value
	}

	return defaultValue
}
