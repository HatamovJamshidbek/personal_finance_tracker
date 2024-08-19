package mongo

import (
	pb "budgeting_service/genproto/budgeting_service"
	"budgeting_service/pkg/logger"
	"context"
	_ "database/sql"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	_ "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	_ "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"time"
)

type TransactionMongo struct {
	Coll *mongo.Collection
	log  logger.ILogger
}

func NewTransactionMongoStore(db *mongo.Collection, lg logger.ILogger) *TransactionMongo {
	return &TransactionMongo{
		Coll: db,
		log:  lg,
	}
}

type Transaction struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`        // ObjectID for MongoDB
	UserId      string             `bson:"user_id" json:"user_id"`         // User ID as a string
	AccountId   string             `bson:"account_id" json:"account_id"`   // Account ID as a string
	CategoryId  string             `bson:"category_id" json:"category_id"` // Category ID as a string
	Amount      float32            `bson:"amount" json:"amount"`           // Amount as a float32
	Type        string             `bson:"type" json:"type"`               // Type (spending/income) as a string
	Description string             `bson:"description" json:"description"` // Description as a string
	Date        string             `bson:"date" json:"date"`               // Date as a time.Time
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`   // Created timestamp as a time.Time
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`   // Updated timestamp as a time.Time
}

func (mongodb *TransactionMongo) CreateTransaction(ctx context.Context, request *pb.TransactionRequest) (*pb.Transaction, error) {
	var (
		err         error
		currentTime = time.Now()
	)

	transaction := &Transaction{
		UserId:      request.GetUserId(),
		AccountId:   request.GetAccountId(),
		CategoryId:  request.GetCategoryId(),
		Amount:      request.GetAmount(),
		Type:        request.GetType(),
		Description: request.GetDescription(),
		Date:        request.Date,
		CreatedAt:   currentTime,
		UpdatedAt:   currentTime,
	}
	fmt.Println("request", request.Date)
	fmt.Println("transaction", transaction.Date)

	result, err := mongodb.Coll.InsertOne(ctx, transaction)
	if err != nil {
		mongodb.log.Error("error inserting transaction", logger.Error(err))
		return nil, err
	}

	transaction.ID = result.InsertedID.(primitive.ObjectID)

	response := &pb.Transaction{
		Id:          transaction.ID.Hex(),
		UserId:      transaction.UserId,
		AccountId:   transaction.AccountId,
		CategoryId:  transaction.CategoryId,
		Amount:      transaction.Amount,
		Type:        transaction.Type,
		Description: transaction.Description,
		Date:        request.Date,
		CreatedAt:   transaction.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   transaction.UpdatedAt.Format(time.RFC3339),
	}

	return response, nil
}

func (mongodb *TransactionMongo) UpdateTransaction(ctx context.Context, request *pb.Transaction) (*pb.Transaction, error) {
	var (
		params      = bson.M{}
		currentTime = time.Now()
	)

	if request.GetUserId() != "" {
		params["user_id"] = request.GetUserId()
	}
	if request.GetAccountId() != "" {
		params["account_id"] = request.GetAccountId()
	}
	if request.GetCategoryId() != "" {
		params["category_id"] = request.GetCategoryId()
	}
	if request.GetAmount() > 0 {
		params["amount"] = request.GetAmount()
	}
	if request.GetType() != "" {
		params["type"] = request.GetType()
	}
	if request.GetDescription() != "" {
		params["description"] = request.GetDescription()
	}
	if request.GetDate() != "" {
		params["date"] = request.GetDate()
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
	transaction := pb.Transaction{}

	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("this budget is not found than one")
	}

	err = mongodb.Coll.FindOne(ctx, filter).Decode(&transaction)
	if err != nil {
		return nil, err
	}
	fmt.Println("+++++++++++++++++++++++++++------------------------", transaction)
	response := pb.Transaction{
		UserId:      transaction.UserId,
		AccountId:   transaction.AccountId,
		CategoryId:  transaction.CategoryId,
		Amount:      transaction.Amount,
		Type:        transaction.Type,
		Description: transaction.Description,
		Date:        transaction.Date,
		CreatedAt:   transaction.CreatedAt,
		UpdatedAt:   transaction.UpdatedAt,
	}

	return &response, nil
}

func (mongodb *TransactionMongo) GetTransaction(ctx context.Context, request *pb.PrimaryKey) (*pb.Transaction, error) {
	var (
		err         error
		transaction Transaction
	)
	id, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, err
	}
	err = mongodb.Coll.FindOne(ctx, bson.M{"_id": id}).Decode(&transaction)
	if err != nil {
		return nil, err
	}
	response := pb.Transaction{
		Id:          transaction.ID.String(),
		UserId:      transaction.UserId,
		CategoryId:  transaction.CategoryId,
		AccountId:   transaction.AccountId,
		Amount:      transaction.Amount,
		Type:        transaction.Type,
		Description: transaction.Description,
		Date:        transaction.Description,
		CreatedAt:   transaction.CreatedAt.String(),
		UpdatedAt:   transaction.UpdatedAt.String(),
	}
	return &response, nil
}
func (mongodb *TransactionMongo) GetAllTransaction(ctx context.Context, request *pb.GetListRequest) (*pb.Transactions, error) {
	var (
		transactions = []*pb.Transaction{}
		offset       = (request.GetPage() - 1) * request.GetLimit()
		count        int64
		err          error
		filter       = bson.M{}
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
		var example Transaction
		if err := cursor.Decode(&example); err != nil {
			log.Printf("Error decoding transaction: %v", err)
			continue
		}
		transaction := &pb.Transaction{
			Id:          example.ID.String(),
			UserId:      example.UserId,
			AccountId:   example.AccountId,
			CategoryId:  example.CategoryId,
			Amount:      example.Amount,
			Type:        example.Type,
			Description: example.Description,
			Date:        example.Date,
			CreatedAt:   example.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   example.UpdatedAt.Format(time.RFC3339),
		}
		transactions = append(transactions, transaction)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &pb.Transactions{
		Transactions: transactions,
		Count:        int32(count),
	}, nil
}
func (mongodb *TransactionMongo) DeleteTransaction(ctx context.Context, request *pb.PrimaryKey) (*emptypb.Empty, error) {
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

type Money struct {
	TotalAmount float32 `bson:"total_amount"`
	Time        string  `bson:"time" json:"time"`
	CategoryId  string  `bson:"category_id" json:"category_id"`
}

func (mongodb *TransactionMongo) GetUserTotalSpend(ctx context.Context, request *pb.GetUserMoneyRequest) (*pb.GetUserMoneyResponse, error) {
	var (
		response pb.GetUserMoneyResponse
	)

	pipeline := mongo.Pipeline{
		{
			{"$match", bson.D{
				{"user_id", request.GetUserId()},
				{"type", "expense"},
				{"date", bson.D{
					{"$gte", request.GetStartTime()},
					{"$lte", request.GetEndTime()},
				}},
			}},
		},
		{
			{"$group", bson.D{
				{"_id", "$category_id"},
				{"total_income", bson.D{{"$sum", "$amount"}}},
				{"first_date", bson.D{{"$min", "$date"}}},
			}},
		},
		{
			{"$project", bson.D{
				{"category_id", "$_id"},
				{"total_income", 1},
				{"first_date", 1},
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
		TotalIncome float32   `bson:"total_income"`
		FirstDate   time.Time `bson:"first_date"`
	}

	if err = cursor.All(ctx, &results); err != nil {
		mongodb.log.Error("Failed to decode aggregation results", logger.Error(err))
		return nil, err
	}

	for _, result := range results {
		response.CategoryId = result.CategoryId
		response.TotalAmount += result.TotalIncome
		response.Time = result.FirstDate.Format(time.RFC3339)
	}

	return &response, nil
}
func (mongodb *TransactionMongo) GetUserTotalIncome(ctx context.Context, request *pb.GetUserMoneyRequest) (*pb.GetUserMoneyResponse, error) {
	var (
		response pb.GetUserMoneyResponse
	)

	pipeline := mongo.Pipeline{
		{
			{"$match", bson.D{
				{"user_id", request.GetUserId()},
				{"type", "income"},
				{"date", bson.D{
					{"$gte", request.GetStartTime()},
					{"$lte", request.GetEndTime()},
				}},
			}},
		},
		{
			{"$group", bson.D{
				{"_id", "$category_id"},
				{"total_income", bson.D{{"$sum", "$amount"}}},
				{"first_date", bson.D{{"$min", "$date"}}},
			}},
		},
		{
			{"$project", bson.D{
				{"category_id", "$_id"},
				{"total_income", 1},
				{"first_date", 1},
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
		TotalIncome float32   `bson:"total_income"`
		FirstDate   time.Time `bson:"first_date"`
	}

	if err = cursor.All(ctx, &results); err != nil {
		mongodb.log.Error("Failed to decode aggregation results", logger.Error(err))
		return nil, err
	}

	for _, result := range results {
		response.CategoryId = result.CategoryId
		response.TotalAmount += result.TotalIncome
		response.Time = result.FirstDate.Format(time.RFC3339)
	}

	return &response, nil
}
