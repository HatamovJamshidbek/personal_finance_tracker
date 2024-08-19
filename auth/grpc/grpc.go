package grpc

import (
	pb "auth_service/genproto/auth_service"
	"auth_service/pkg/logger"
	"auth_service/service"
	"auth_service/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func SetUpServer(storage storage.IStorage, log logger.ILogger) *grpc.Server {
	grpcServer := grpc.NewServer()

	pb.RegisterUserServiceServer(grpcServer, service.NewUserService(storage, log))

	reflection.Register(grpcServer)
	return grpcServer
}
