package service

import (
	pb "budgeting_service/genproto/budgeting_service"
	"budgeting_service/pkg/logger"
	"budgeting_service/storage"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CategoryService struct {
	storage storage.IStorage
	log     logger.ILogger
	pb.UnimplementedCategoryServiceServer
}

func NewCategoryService(storage storage.IStorage, log logger.ILogger) *CategoryService {
	return &CategoryService{
		storage: storage,
		log:     log,
	}
}

func (service *CategoryService) CreateCategory(ctx context.Context, in *pb.CategoryRequest) (*pb.Category, error) {
	return service.storage.Categories().CreateCategory(ctx, in)

}
func (service *CategoryService) UpdateCategory(ctx context.Context, in *pb.Category) (*pb.Category, error) {
	return service.storage.Categories().UpdateCategory(ctx, in)
}
func (service *CategoryService) GetCategory(ctx context.Context, in *pb.PrimaryKey) (*pb.Category, error) {
	return service.storage.Categories().GetCategory(ctx, in)
}
func (service *CategoryService) GetListCategories(ctx context.Context, in *pb.GetListRequest) (*pb.Categories, error) {
	return service.storage.Categories().GetAllCategory(ctx, in)
}
func (service *CategoryService) DeleteCategory(ctx context.Context, in *pb.PrimaryKey) (*emptypb.Empty, error) {
	return service.storage.Categories().DeleteCategory(ctx, in)
}
