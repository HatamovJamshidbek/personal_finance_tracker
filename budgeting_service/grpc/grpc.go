package grpc

import (
	pb "budgeting_service/genproto/budgeting_service"
	"budgeting_service/pkg/logger"
	"budgeting_service/service"
	"budgeting_service/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	_ "google.golang.org/grpc/reflection"
)

func SetUpServer(storage storage.IStorage, log logger.ILogger) *grpc.Server {
	grpcServer := grpc.NewServer()

	pb.RegisterAccountServiceServer(grpcServer, service.NewAccountService(storage, log))
	pb.RegisterBudgetServiceServer(grpcServer, service.NewBudgetService(storage, log))
	pb.RegisterCategoryServiceServer(grpcServer, service.NewCategoryService(storage, log))
	pb.RegisterGoalServiceServer(grpcServer, service.NewGoalService(storage, log))
	pb.RegisterTransactionServiceServer(grpcServer, service.NewTransactionService(storage, log))

	reflection.Register(grpcServer)
	return grpcServer
}
