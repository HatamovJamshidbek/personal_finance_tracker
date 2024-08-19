package service

import (
	"auth_service/api/models"
	pb "auth_service/genproto/auth_service"
	"auth_service/pkg/logger"
	"auth_service/storage"
	"context"
)

type UserService struct {
	log     logger.ILogger
	storage storage.IStorage
	pb.UnimplementedUserServiceServer
}

func NewUserService(storage storage.IStorage, log logger.ILogger) *UserService {
	return &UserService{
		storage: storage,
		log:     log,
	}
}

func (service *UserService) RegisterUser(ctx context.Context, in models.Register) (*models.RegisterResponse, error) {
	return service.storage.Users().RegisterUser(ctx, in)
}
func (service *UserService) Login(ctx context.Context, in *models.LoginRequest) (*models.LoginResponse, error) {
	return service.storage.Users().Login(ctx, in)
}
func (service *UserService) ResetPassword(ctx context.Context, in *models.LoginRequest) (string, error) {
	return service.storage.Users().ResetPassword(ctx, in)
}

func (service *UserService) UpdateUserProfile(ctx context.Context, in *pb.User) (*pb.User, error) {
	return service.storage.Users().UpdateUserProfile(ctx, in)

}
func (service *UserService) GetUserProfile(ctx context.Context, in *pb.PrimaryKeyUser) (*pb.User, error) {
	return service.storage.Users().GetUserProfile(ctx, in)

}
