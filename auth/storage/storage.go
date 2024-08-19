package storage

import (
	"auth_service/api/models"
	pb "auth_service/genproto/auth_service"
	"context"
)

type IStorage interface {
	Close()
	Users() IUserStorage
}

type IUserStorage interface {
	RegisterUser(ctx context.Context, request models.Register) (*models.RegisterResponse, error)
	Login(ctx context.Context, request *models.LoginRequest) (*models.LoginResponse, error)
	ResetPassword(ctx context.Context, request *models.LoginRequest) (string, error)
	UpdateUserProfile(ctx context.Context, request *pb.User) (*pb.User, error)
	GetUserProfile(ctx context.Context, request *pb.PrimaryKeyUser) (*pb.User, error)
}
