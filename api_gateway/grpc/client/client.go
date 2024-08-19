package client

import (
	"api_gateway/configs"
	pbu "api_gateway/genproto/auth_service"
	pbb "api_gateway/genproto/budgeting_service"

	"google.golang.org/grpc"
)

type IServiceManager interface {
	// ESTATE SERVICE
	GoalService() pbb.GoalServiceClient
	TransactionService() pbb.TransactionServiceClient
	CategoryService() pbb.CategoryServiceClient
	BudgetService() pbb.BudgetServiceClient
	AccountService() pbb.AccountServiceClient

	// USER SERVICE
	UserService() pbu.UserServiceClient
}

type grpcClients struct {
	// ESTATE SERVICE
	goalService        pbb.GoalServiceClient
	transactionService pbb.TransactionServiceClient
	categoryService    pbb.CategoryServiceClient
	budgetService      pbb.BudgetServiceClient
	accountService     pbb.AccountServiceClient

	// USER SERVICE
	userService pbu.UserServiceClient
}

// USER SERVICE
func (g *grpcClients) UserService() pbu.UserServiceClient {
	return g.userService
}

// BUDGETING SERVICE
func (g *grpcClients) GoalService() pbb.GoalServiceClient {
	return g.goalService
}

func (g *grpcClients) TransactionService() pbb.TransactionServiceClient {
	return g.transactionService
}

func (g *grpcClients) CategoryService() pbb.CategoryServiceClient {
	return g.categoryService
}

func (g *grpcClients) BudgetService() pbb.BudgetServiceClient {
	return g.budgetService
}

func (g *grpcClients) AccountService() pbb.AccountServiceClient {
	return g.accountService
}

func NewGrpcClients(cfg configs.Config) (IServiceManager, error) {
	// CONNECTION WITH BUDGET SERVICE
	connBudgetService, err := grpc.Dial(
		cfg.BudgetServiceGrpcHost+cfg.BudgetServiceGrpcPort,
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	// CONNECTION WITH USER SERVICE
	connUserService, err := grpc.Dial(
		cfg.UserServiceGrpcHost+cfg.UserServiceGrpcPort,
		grpc.WithInsecure(),
	)

	if err != nil {
		return nil, err
	}

	return &grpcClients{
		// BUDGET SERVICE
		goalService:        pbb.NewGoalServiceClient(connBudgetService),
		transactionService: pbb.NewTransactionServiceClient(connBudgetService),
		categoryService:    pbb.NewCategoryServiceClient(connBudgetService),
		budgetService:      pbb.NewBudgetServiceClient(connBudgetService),
		accountService:     pbb.NewAccountServiceClient(connBudgetService),

		// USER SERVICE
		userService: pbu.NewUserServiceClient(connUserService),
	}, nil

}
