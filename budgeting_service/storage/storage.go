package storage

import (
	pb "budgeting_service/genproto/budgeting_service"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

type IStorage interface {
	Accounts() IAccountStorage
	Budgets() IBudgetStorage
	Categories() ICategoryStorage
	Goals() IGoalStorage
	Transactions() ITransactionStorage
}

type IAccountStorage interface {
	CreateAccount(ctx context.Context, request *pb.AccountRequest) (*pb.Account, error)
	UpdateAccount(ctx context.Context, request *pb.Account) (*pb.Account, error)
	GetAccount(ctx context.Context, request *pb.PrimaryKey) (*pb.Account, error)
	GetAllAccount(ctx context.Context, request *pb.GetListRequest) (*pb.Accounts, error)
	DeleteAccount(ctx context.Context, request *pb.PrimaryKey) (*emptypb.Empty, error)
}
type IBudgetStorage interface {
	CreateBudget(ctx context.Context, request *pb.BudgetRequest) (*pb.Budget, error)
	UpdateBudget(ctx context.Context, request *pb.Budget) (*pb.Budget, error)
	GetBudget(ctx context.Context, request *pb.PrimaryKey) (*pb.Budget, error)
	GetAllBudget(ctx context.Context, request *pb.GetListRequest) (*pb.Budgets, error)
	DeleteBudget(ctx context.Context, request *pb.PrimaryKey) (*emptypb.Empty, error)
	GetUserBudgetSummary(ctx context.Context, request *pb.PrimaryKey) (*pb.GetUserBudgetResponse, error)
}
type ICategoryStorage interface {
	CreateCategory(ctx context.Context, request *pb.CategoryRequest) (*pb.Category, error)
	UpdateCategory(ctx context.Context, request *pb.Category) (*pb.Category, error)
	GetCategory(ctx context.Context, request *pb.PrimaryKey) (*pb.Category, error)
	GetAllCategory(ctx context.Context, request *pb.GetListRequest) (*pb.Categories, error)
	DeleteCategory(ctx context.Context, request *pb.PrimaryKey) (*emptypb.Empty, error)
}
type IGoalStorage interface {
	CreateGoal(ctx context.Context, request *pb.GoalRequest) (*pb.Goal, error)
	UpdateGoal(ctx context.Context, request *pb.Goal) (*pb.Goal, error)
	GetGoal(ctx context.Context, request *pb.PrimaryKey) (*pb.Goal, error)
	GetAllGoal(ctx context.Context, request *pb.GetListRequest) (*pb.Goals, error)
	DeleteGoal(ctx context.Context, request *pb.PrimaryKey) (*emptypb.Empty, error)
	GenerateGoalProgressReport(ctx context.Context, request *pb.GoalProgressRequest) (*pb.GoalProgressResponse, error)
}
type ITransactionStorage interface {
	CreateTransaction(ctx context.Context, request *pb.TransactionRequest) (*pb.Transaction, error)
	UpdateTransaction(ctx context.Context, request *pb.Transaction) (*pb.Transaction, error)
	GetTransaction(ctx context.Context, request *pb.PrimaryKey) (*pb.Transaction, error)
	GetAllTransaction(ctx context.Context, request *pb.GetListRequest) (*pb.Transactions, error)
	DeleteTransaction(ctx context.Context, request *pb.PrimaryKey) (*emptypb.Empty, error)
	GetUserTotalIncome(ctx context.Context, request *pb.GetUserMoneyRequest) (*pb.GetUserMoneyResponse, error)
	GetUserTotalSpend(ctx context.Context, request *pb.GetUserMoneyRequest) (*pb.GetUserMoneyResponse, error)
}
