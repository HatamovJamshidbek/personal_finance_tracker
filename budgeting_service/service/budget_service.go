package service

import (
	pb "budgeting_service/genproto/budgeting_service"
	"budgeting_service/pkg/logger"
	"budgeting_service/storage"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

type BudgetService struct {
	storage storage.IStorage
	log     logger.ILogger
	pb.UnimplementedBudgetServiceServer
}

func NewBudgetService(storage storage.IStorage, log logger.ILogger) *BudgetService {
	return &BudgetService{
		storage: storage,
		log:     log,
	}
}

func (service *BudgetService) CreateBudget(ctx context.Context, in *pb.BudgetRequest) (*pb.Budget, error) {
	return service.storage.Budgets().CreateBudget(ctx, in)

}
func (service *BudgetService) UpdateBudget(ctx context.Context, in *pb.Budget) (*pb.Budget, error) {
	return service.storage.Budgets().UpdateBudget(ctx, in)
}
func (service *BudgetService) GetBudget(ctx context.Context, in *pb.PrimaryKey) (*pb.Budget, error) {
	return service.storage.Budgets().GetBudget(ctx, in)
}
func (service *BudgetService) GetListBudgets(ctx context.Context, in *pb.GetListRequest) (*pb.Budgets, error) {
	return service.storage.Budgets().GetAllBudget(ctx, in)
}
func (service *BudgetService) DeleteBudget(ctx context.Context, in *pb.PrimaryKey) (*emptypb.Empty, error) {
	return service.storage.Budgets().DeleteBudget(ctx, in)
}

func (service *BudgetService) GetBudgetSummary(ctx context.Context, in *pb.PrimaryKey) (*pb.GetUserBudgetResponse, error) {
	return service.storage.Budgets().GetUserBudgetSummary(ctx, in)
}
