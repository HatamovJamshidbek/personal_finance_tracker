package service

import (
	pb "budgeting_service/genproto/budgeting_service"
	"budgeting_service/pkg/logger"
	"budgeting_service/storage"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AccountService struct {
	storage storage.IStorage
	log     logger.ILogger
	pb.UnimplementedAccountServiceServer
}

func NewAccountService(storage storage.IStorage, log logger.ILogger) *AccountService {
	return &AccountService{
		storage: storage,
		log:     log,
	}
}

func (service *AccountService) CreateAccount(ctx context.Context, in *pb.AccountRequest) (*pb.Account, error) {
	return service.storage.Accounts().CreateAccount(ctx, in)

}
func (service *AccountService) UpdateAccount(ctx context.Context, in *pb.Account) (*pb.Account, error) {
	return service.storage.Accounts().UpdateAccount(ctx, in)
}
func (service *AccountService) GetAccount(ctx context.Context, in *pb.PrimaryKey) (*pb.Account, error) {
	return service.storage.Accounts().GetAccount(ctx, in)
}
func (service *AccountService) GetListAccounts(ctx context.Context, in *pb.GetListRequest) (*pb.Accounts, error) {
	return service.storage.Accounts().GetAllAccount(ctx, in)
}
func (service *AccountService) DeleteAccount(ctx context.Context, in *pb.PrimaryKey) (*emptypb.Empty, error) {
	return service.storage.Accounts().DeleteAccount(ctx, in)
}
