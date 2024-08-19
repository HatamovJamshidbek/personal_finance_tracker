package mongo

import (
	pb "budgeting_service/genproto/budgeting_service"
	"budgeting_service/pkg/logger"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type Goal struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserId        string             `bson:"user_id" json:"user_id"`
	Name          string             `bson:"name" json:"name"`
	Type          string             `bson:"type" json:"type"`
	TargetAmount  float32            `bson:"target_amount" json:"target_amount"`
	CurrentAmount float32            `bson:"current_amount" json:"current_amount"`
	Deadline      string             `json:"deadline" bson:"deadline"`
	Status        string             `bson:"status" json:"status"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type GoalMongo struct {
	Coll *mongo.Collection
	log  logger.ILogger
}

func NewGoalMongoStore(db *mongo.Collection, lg logger.ILogger) *GoalMongo {
	return &GoalMongo{
		Coll: db,
		log:  lg,
	}
}

func (mongodb *GoalMongo) CreateGoal(ctx context.Context, request *pb.GoalRequest) (*pb.Goal, error) {
	var (
		err         error
		example     Goal
		currentTime = time.Now()
	)

	budget := bson.D{
		{Key: "user_id", Value: request.GetUserId()},
		{Key: "name", Value: request.GetName()},
		{Key: "target_amount", Value: request.GetTargetAmount()},
		{Key: "status", Value: request.GetStatus()},
		{Key: "current_amount", Value: request.GetCurrentAmount()},
		{Key: "deadline", Value: request.GetDeadline()},
		{Key: "created_at", Value: currentTime},
		{Key: "updated_at", Value: currentTime},
	}
	result, err := mongodb.Coll.InsertOne(ctx, &budget)
	if err != nil {
		mongodb.log.Error("insert one error", logger.Error(err))
		return nil, err
	}
	err = mongodb.Coll.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&example)
	if err != nil {
		return nil, err
	}
	response := pb.Goal{
		Id:            example.ID.String(),
		UserId:        example.UserId,
		Name:          example.Name,
		TargetAmount:  example.TargetAmount,
		CurrentAmount: example.CurrentAmount,
		Deadline:      example.Deadline,
		Status:        example.Status,
		CreatedAt:     example.CreatedAt.String(),
		UpdatedAt:     example.UpdatedAt.String(),
	}
	return &response, nil
}

func (mongodb *GoalMongo) UpdateGoal(ctx context.Context, request *pb.Goal) (*pb.Goal, error) {
	var (
		params      = bson.M{}
		example     Goal
		currentTime = time.Now()
	)

	if request.GetUserId() != "" {
		params["user_id"] = request.GetUserId()
	}
	if request.GetName() != "" {
		params["name"] = request.GetName()
	}
	if request.Deadline != "" {
		params["deadline"] = request.GetDeadline()
	}
	if request.CurrentAmount > 0 {
		params["current_amount"] = request.GetCurrentAmount()
	}
	if request.TargetAmount > 0 {
		params["target_amount"] = request.GetTargetAmount()
	}
	if request.Status != "" {
		params["status"] = request.GetStatus()
	}
	params["updated_at"] = currentTime

	id, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": params}

	result, err := mongodb.Coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("goal not found")
	}

	err = mongodb.Coll.FindOne(ctx, filter).Decode(&example)
	if err != nil {
		return nil, err
	}
	response := pb.Goal{
		Id:            example.ID.String(),
		UserId:        example.UserId,
		Name:          example.Name,
		TargetAmount:  example.TargetAmount,
		CurrentAmount: example.CurrentAmount,
		Deadline:      example.Deadline,
		Status:        example.Status,
		CreatedAt:     example.CreatedAt.String(),
		UpdatedAt:     example.UpdatedAt.String(),
	}
	return &response, nil
}

func (mongodb *GoalMongo) GetGoal(ctx context.Context, request *pb.PrimaryKey) (*pb.Goal, error) {
	var (
		err     error
		example Goal
	)
	id, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, err
	}
	err = mongodb.Coll.FindOne(ctx, bson.M{"_id": id}).Decode(&example)
	if err != nil {
		return nil, err
	}
	response := pb.Goal{
		Id:            example.ID.String(),
		UserId:        example.UserId,
		Name:          example.Name,
		TargetAmount:  example.TargetAmount,
		CurrentAmount: example.CurrentAmount,
		Deadline:      example.Deadline,
		Status:        example.Status,
		CreatedAt:     example.CreatedAt.String(),
		UpdatedAt:     example.UpdatedAt.String(),
	}
	return &response, nil
}

func (mongodb *GoalMongo) GetAllGoal(ctx context.Context, request *pb.GetListRequest) (*pb.Goals, error) {
	var (
		goals  = []*pb.Goal{}
		offset = (request.GetPage() - 1) * request.GetLimit()
		count  int64
		err    error
		filter = bson.M{}
	)

	if request.GetSearch() != "" {
		filter["name"] = bson.M{"$regex": request.GetSearch(), "$options": "i"}
	}

	count, err = mongodb.Coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	cursor, err := mongodb.Coll.Find(ctx, filter, options.Find().
		SetSkip(int64(offset)).
		SetLimit(int64(request.GetLimit())).
		SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var example Goal
		if err = cursor.Decode(&example); err != nil {
			return nil, err
		}
		goal := pb.Goal{
			Id:            example.ID.String(),
			UserId:        example.UserId,
			Name:          example.Name,
			TargetAmount:  example.TargetAmount,
			CurrentAmount: example.CurrentAmount,
			Deadline:      example.Deadline,
			Status:        example.Status,
			CreatedAt:     example.CreatedAt.String(),
			UpdatedAt:     example.UpdatedAt.String(),
		}
		goals = append(goals, &goal)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &pb.Goals{
		Goals: goals,
		Count: int32(count),
	}, nil
}

func (mongodb *GoalMongo) DeleteGoal(ctx context.Context, request *pb.PrimaryKey) (*emptypb.Empty, error) {
	var (
		err    error
		filter bson.M
	)
	id, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, err
	}
	filter = bson.M{"_id": id}
	_, err = mongodb.Coll.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
func (mongodb *GoalMongo) GenerateGoalProgressReport(ctx context.Context, request *pb.GoalProgressRequest) (*pb.GoalProgressResponse, error) {
	var (
		response pb.GoalProgressResponse
	)

	pipeline := mongo.Pipeline{
		{
			{"$match", bson.D{
				{"user_id", request.GetUserId()},
				{"deadline", bson.D{
					{"$gte", request.GetStartTime()},
					{"$lte", request.GetEndTime()},
				}},
			}},
		},
		{
			{"$group", bson.D{
				{"_id", "$status"},
				{"target_amount_sum", bson.D{{"$sum", "$target_amount"}}},
				{"current_amount_sum", bson.D{{"$sum", "$current_amount"}}},
				{"total_amount", bson.D{{"$sum", bson.D{{"$add", bson.A{"$target_amount", "$current_amount"}}}}}},
			}},
		},
		{
			{"$project", bson.D{
				{"status", "$_id"},
				{"target_amount_sum", 1},
				{"current_amount_sum", 1},
				{"total_amount", 1},
				{"_id", 0},
			}},
		},
	}

	cursor, err := mongodb.Coll.Aggregate(ctx, pipeline)
	if err != nil {
		mongodb.log.Error("Failed to execute aggregation pipeline", logger.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		Status           string  `bson:"status"`
		TargetAmountSum  float64 `bson:"target_amount_sum"`
		CurrentAmountSum float64 `bson:"current_amount_sum"`
		TotalAmount      float64 `bson:"total_amount"`
	}

	if err = cursor.All(ctx, &results); err != nil {
		mongodb.log.Error("Failed to decode aggregation results", logger.Error(err))
		return nil, err
	}

	for _, result := range results {
		progress := &pb.GoalProgressItem{
			Status:           result.Status,
			TargetAmountSum:  float32(result.TargetAmountSum),
			CurrentAmountSum: float32(result.CurrentAmountSum),
			TotalAmount:      float32(result.TotalAmount),
		}
		response.Results = append(response.Results, progress)
	}

	return &response, nil
}
