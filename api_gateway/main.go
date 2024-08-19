package main

import (
	"api_gateway/api"
	"api_gateway/api/handler"
	"api_gateway/configs"
	"api_gateway/grpc/client"
	"api_gateway/pkg/logger"
	"api_gateway/rabbitMq"
	"fmt"
	"github.com/casbin/casbin/v2"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := configs.Load()

	loggerLevel := logger.LevelDebug

	switch cfg.Environment {
	case configs.DebugMode:
		loggerLevel = logger.LevelDebug
		gin.SetMode(gin.DebugMode)
	case configs.TestMode:
		loggerLevel = logger.LevelDebug
		gin.SetMode(gin.TestMode)
	default:
		loggerLevel = logger.LevelInfo
		gin.SetMode(gin.ReleaseMode)
	}

	log := logger.NewLogger(cfg.ServiceName, loggerLevel)
	defer logger.Cleanup(log)

	services, err := client.NewGrpcClients(cfg)
	if err != nil {
		log.Error("Error while initializing grpc clients", logger.Error(err))
		return
	}
	time.Sleep(10 * time.Second)
	writer, err := rabbitMq.NewRabbitMqProducerInt("amqp://user:1111@rabbitmq:5672")
	if err != nil {
		log.Error("this error is newRabbitMq", logger.Error(err))
		return
	}

	h := handler.New(cfg, services, log, writer)

	enforcer, err := casbin.NewEnforcer("casbin/model.conf", "casbin/policy.csv")
	if err != nil {
		fmt.Println("++++++++++", err)

	}
	r := api.New(h, enforcer)

	log.Info("Server is running ...", logger.Any("port", cfg.HTTPort))
	if err = r.Run(cfg.HTTPort); err != nil {
		log.Error("Error while running server", logger.Error(err))
	}
}
