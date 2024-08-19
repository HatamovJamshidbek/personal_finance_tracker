package main

import (
	"budgeting_service/configs"
	"budgeting_service/grpc"
	"budgeting_service/pkg/logger"
	"budgeting_service/rabbitMq"
	"budgeting_service/storage/mongo"
	"budgeting_service/storage/redis"
	"context"
	"github.com/gin-gonic/gin"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := configs.Load()

	// Set logger level based on the environment
	var loggerLevel string
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
	redcon := redis.RedisConnection()
	reds := redis.NewRedisRepo(redcon)

	mongoStore, err := mongo.NewStore(context.Background(), cfg, log, reds)
	if err != nil {
		log.Error("error while connecting to mongo", logger.Error(err))
		return
	}
	time.Sleep(5 * time.Second)
	rabbitMq.SetUpConsumers(log, mongoStore)
	time.Sleep(5 * time.Second)
	grpcServer := grpc.SetUpServer(mongoStore, log)

	lis, err := net.Listen("tcp", cfg.BudgetServiceGrpcHost+cfg.BudgetServiceGrpcPort)
	if err != nil {
		log.Error("error while listening grpc host port", logger.Error(err))
		return
	}

	log.Info("Service is running...", logger.Any("grpc port", cfg.BudgetServiceGrpcPort))

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Error("error while serving grpc", logger.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down service...")
	grpcServer.GracefulStop()
}
