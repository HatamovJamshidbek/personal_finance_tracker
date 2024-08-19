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

const dateTimeLayout = "2006-01-02 15:04:05"

type Budget struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserId     string             `bson:"user_id" json:"user_id"`
	CategoryId string             `bson:"category_id" json:"category_id"`
	Amount     float32            `bson:"amount" json:"amount"`
	Period     string             `bson:"period" json:"period"`
	StartTime  time.Time          `bson:"start_time" json:"start_time"`
	EndTime    time.Time          `bson:"end_time" json:"end_time"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

type BudgetMongo struct {
	Coll *mongo.Collection
	log  logger.ILogger
}

func NewBudgetMongoStore(db *mongo.Collection, lg logger.ILogger) *BudgetMongo {
	return &BudgetMongo{
		Coll: db,
		log:  lg,
	}
}

func (mongodb *BudgetMongo) CreateBudget(ctx context.Context, request *pb.BudgetRequest) (*pb.Budget, error) {
	currentTime := time.Now()
	fmt.Println("+++++++++")

	budget := bson.D{
		{Key: "user_id", Value: request.GetUserId()},
		{Key: "category_id", Value: request.GetCategoryId()},
		{Key: "amount", Value: request.GetAmount()},
		{Key: "period", Value: request.GetPeriod()},
		{Key: "start_time", Value: request.GetStartTime()},
		{Key: "end_time", Value: request.GetEndTime()},
		{Key: "created_at", Value: currentTime},
		{Key: "updated_at", Value: currentTime},
	}
	fmt.Println("+++++++++++", request)
	result, err := mongodb.Coll.InsertOne(ctx, budget)
	if err != nil {
		mongodb.log.Error("InsertOne failed", logger.Error(err))
		return nil, err
	}

	var createdBudget Budget
	err = mongodb.Coll.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&createdBudget)
	if err != nil {
		return nil, err
	}
	fmt.Println("+++++++budget", createdBudget)

	response := pb.Budget{
		Id:         createdBudget.ID.Hex(),
		UserId:     createdBudget.UserId,
		CategoryId: createdBudget.CategoryId,
		Amount:     createdBudget.Amount,
		Period:     createdBudget.Period,
		StartTime:  createdBudget.StartTime.Format(dateTimeLayout),
		EndTime:    createdBudget.EndTime.Format(dateTimeLayout),
		CreatedAt:  createdBudget.CreatedAt.Format(dateTimeLayout),
		UpdatedAt:  createdBudget.UpdatedAt.Format(dateTimeLayout),
	}
	fmt.Println("+++++++response", response)

	return &response, nil
}

func (mongodb *BudgetMongo) UpdateBudget(ctx context.Context, request *pb.Budget) (*pb.Budget, error) {
	params := bson.M{}
	if request.GetUserId() != "" {
		params["user_id"] = request.GetUserId()
	}
	if request.GetCategoryId() != "" {
		params["category_id"] = request.GetCategoryId()
	}
	if request.GetAmount() > 0 {
		params["amount"] = request.GetAmount()
	}
	if request.GetPeriod() != "" {
		params["period"] = request.GetPeriod()
	}
	if request.GetStartTime() != "" {
		params["start_time"] = request.GetStartTime()
	}
	if request.GetEndTime() != "" {
		params["end_time"] = request.GetEndTime()
	}
	params["updated_at"] = time.Now()

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
		return nil, fmt.Errorf("budget not found")
	}

	var budget Budget
	err = mongodb.Coll.FindOne(ctx, filter).Decode(&budget)
	if err != nil {
		return nil, err
	}

	response := pb.Budget{
		Id:         budget.ID.Hex(),
		UserId:     budget.UserId,
		CategoryId: budget.CategoryId,
		Amount:     budget.Amount,
		Period:     budget.Period,
		StartTime:  budget.StartTime.Format(dateTimeLayout),
		EndTime:    budget.EndTime.Format(dateTimeLayout),
		CreatedAt:  budget.CreatedAt.Format(dateTimeLayout),
		UpdatedAt:  budget.UpdatedAt.Format(dateTimeLayout),
	}

	return &response, nil
}

func (mongodb *BudgetMongo) GetBudget(ctx context.Context, request *pb.PrimaryKey) (*pb.Budget, error) {
	id, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, err
	}

	var budget Budget
	err = mongodb.Coll.FindOne(ctx, bson.M{"_id": id}).Decode(&budget)
	if err != nil {
		return nil, err
	}

	response := pb.Budget{
		Id:         budget.ID.Hex(),
		UserId:     budget.UserId,
		CategoryId: budget.CategoryId,
		Amount:     budget.Amount,
		Period:     budget.Period,
		StartTime:  budget.StartTime.Format(dateTimeLayout),
		EndTime:    budget.EndTime.Format(dateTimeLayout),
		CreatedAt:  budget.CreatedAt.Format(dateTimeLayout),
		UpdatedAt:  budget.UpdatedAt.Format(dateTimeLayout),
	}

	return &response, nil
}

func (mongodb *BudgetMongo) GetAllBudget(ctx context.Context, request *pb.GetListRequest) (*pb.Budgets, error) {
	var (
		budgets = []*pb.Budget{}
		offset  = (request.GetPage() - 1) * request.GetLimit()
		count   int64
		err     error
		filter  = bson.M{}
	)

	if request.GetSearch() != "" {
		filter["category_id"] = bson.M{"$regex": request.GetSearch(), "$options": "i"}
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
		var example Budget
		if err := cursor.Decode(&example); err != nil {
			return nil, err
		}
		budget := pb.Budget{
			Id:         example.ID.Hex(),
			UserId:     example.UserId,
			CategoryId: example.CategoryId,
			Amount:     example.Amount,
			Period:     example.Period,
			StartTime:  example.StartTime.Format(dateTimeLayout),
			EndTime:    example.EndTime.Format(dateTimeLayout),
			CreatedAt:  example.CreatedAt.Format(dateTimeLayout),
			UpdatedAt:  example.UpdatedAt.Format(dateTimeLayout),
		}
		budgets = append(budgets, &budget)
	}

	return &pb.Budgets{
		Count:   int32(count),
		Budgets: budgets,
	}, nil
}

func (mongodb *BudgetMongo) DeleteBudget(ctx context.Context, request *pb.PrimaryKey) (*emptypb.Empty, error) {
	id, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, err
	}

	result, err := mongodb.Coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return nil, err
	}

	if result.DeletedCount == 0 {
		return nil, fmt.Errorf("budget not found")
	}

	return &emptypb.Empty{}, nil
}

func (mongodb *BudgetMongo) GetUserBudgetSummary(ctx context.Context, request *pb.PrimaryKey) (*pb.GetUserBudgetResponse, error) {
	var (
		response pb.GetUserBudgetResponse
	)
	fmt.Println("+++++++++++++Request", request.GetId())

	userId := request.GetId()

	pipeline := mongo.Pipeline{
		{
			{"$match", bson.D{
				{"user_id", userId},
			}},
		},
		{
			{"$group", bson.D{
				{"_id", bson.D{
					{"category_id", "$category_id"},
					{"period", "$period"},
				}},
				{"total_amount", bson.D{{"$sum", "$amount"}}},
				{"start_time", bson.D{{"$first", "$start_time"}}},
				{"end_time", bson.D{{"$last", "$end_time"}}},
			}},
		},
		{
			{"$project", bson.D{
				{"category_id", "$_id.category_id"},
				{"period", "$_id.period"},
				{"total_amount", 1},
				{"start_time", 1},
				{"end_time", 1},
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
		CategoryId  string    `bson:"category_id"`
		TotalAmount float32   `bson:"total_amount"`
		StartTime   time.Time `bson:"start_time"`
		EndTime     time.Time `bson:"end_time"`
		Period      string    `bson:"period"`
	}

	if err = cursor.All(ctx, &results); err != nil {
		mongodb.log.Error("Failed to decode aggregation results", logger.Error(err))
		return nil, err
	}
	fmt.Println("+++++++++++++Result", results)

	for _, result := range results {
		item := &pb.BudgetSummaryItem{
			CategoryId:  result.CategoryId,
			TotalAmount: result.TotalAmount,
			StartTime:   result.StartTime.Format(time.RFC3339),
			EndTime:     result.EndTime.Format(time.RFC3339),
			Period:      result.Period,
		}
		response.Results = append(response.Results, item)
	}

	return &response, nil
}
