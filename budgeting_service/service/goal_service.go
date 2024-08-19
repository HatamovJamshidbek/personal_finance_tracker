package service

import (
	pb "budgeting_service/genproto/budgeting_service"
	"budgeting_service/pkg/logger"
	"budgeting_service/storage"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GoalService struct {
	storage storage.IStorage
	log     logger.ILogger
	pb.UnimplementedGoalServiceServer
}

func NewGoalService(storage storage.IStorage, log logger.ILogger) *GoalService {
	return &GoalService{
		storage: storage,
		log:     log,
	}
}

func (s *GoalService) GetGoalReportProgress(ctx context.Context, in *pb.GoalProgressRequest) (*pb.GoalProgressResponse, error) {
	return s.storage.Goals().GenerateGoalProgressReport(ctx, in)

}
func (s *GoalService) CreateGoal(ctx context.Context, in *pb.GoalRequest) (*pb.Goal, error) {
	return s.storage.Goals().CreateGoal(ctx, in)

}
func (service *GoalService) UpdateGoal(ctx context.Context, in *pb.Goal) (*pb.Goal, error) {
	return service.storage.Goals().UpdateGoal(ctx, in)
}
func (service *GoalService) GetGoal(ctx context.Context, in *pb.PrimaryKey) (*pb.Goal, error) {
	return service.storage.Goals().GetGoal(ctx, in)
}
func (service *GoalService) GetListGoals(ctx context.Context, in *pb.GetListRequest) (*pb.Goals, error) {
	return service.storage.Goals().GetAllGoal(ctx, in)
}
func (service *GoalService) DeleteGoal(ctx context.Context, in *pb.PrimaryKey) (*emptypb.Empty, error) {
	return service.storage.Goals().DeleteGoal(ctx, in)
}
