package mongo

import (
	pb "budgeting_service/genproto/budgeting_service"
	"budgeting_service/pkg/logger"
	"budgeting_service/storage"
	"budgeting_service/storage/redis"
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
	"time"
)

type Account struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" bson:"_id"`
	UserId    string             `bson:"user_id" json:"user_id"`
	Name      string             `bson:"name" json:"name"`
	Type      string             `bson:"type" json:"type"`
	Balance   float32            `bson:"balance" json:"balance"`
	Currency  string             `bson:"currency" json:"currency"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type AccountMongo struct {
	Coll *mongo.Collection
	log  logger.ILogger
	cash *redis.RedisRepo
}

func NewAccountMongoStore(collection *mongo.Collection, log logger.ILogger, redisCash *redis.RedisRepo) storage.IAccountStorage {
	return &AccountMongo{
		Coll: collection,
		log:  log,
		cash: redisCash,
	}
}

func (mongodb *AccountMongo) CreateAccount(ctx context.Context, request *pb.AccountRequest) (*pb.Account, error) {
	var (
		err         error
		currentTime = time.Now()
		example     Account
	)

	account := bson.D{
		{Key: "name", Value: request.Name},
		{Key: "type", Value: request.Type},
		{Key: "balance", Value: request.Balance},
		{Key: "currency", Value: request.Currency},
		{Key: "created_at", Value: currentTime},
	}

	mongodb.log.Info(fmt.Sprintf("Inserting document: %+v", account))

	result, err := mongodb.Coll.InsertOne(ctx, account)
	if err != nil {
		mongodb.log.Error(fmt.Sprintf("InsertOne failed: %v", err))
		return nil, err
	}

	mongodb.log.Info(fmt.Sprintf("InsertOne result: %+v", result))

	err = mongodb.Coll.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&example)
	if err != nil {
		mongodb.log.Error(fmt.Sprintf("FindOne failed: %v", err))
		return nil, err
	}

	response := pb.Account{
		Id:        example.ID.String(),
		UserId:    example.UserId,
		Name:      example.Name,
		Type:      example.Type,
		Balance:   example.Balance,
		Currency:  example.Currency,
		CreatedAt: example.CreatedAt.String(),
		UpdatedAt: example.UpdatedAt.String(),
	}
	cash := &pb.Account{
		Id:      example.ID.String(),
		Balance: example.Balance,
	}

	expiration := 24 * time.Hour
	err = mongodb.cash.AddBalanceInCache(ctx, cash, expiration)
	if err != nil {
		fmt.Println("cash is not created in for account")
	}

	return &response, nil
}

func (mongodb *AccountMongo) UpdateAccount(ctx context.Context, request *pb.Account) (*pb.Account, error) {
	var (
		params      = bson.M{}
		example     Account
		currentTime = time.Now()
	)

	if request.GetUserId() != "" {
		params["user_id"] = request.GetUserId()
	}
	if request.GetName() != "" {
		params["name"] = request.GetName()
	}
	if request.GetType() != "" {
		params["type"] = request.GetType()
	}
	if request.GetCurrency() != "" {
		params["currency"] = request.GetCurrency()
	}
	if request.GetBalance() > -1 {
		params["balance"] = request.GetBalance()
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
		return nil, fmt.Errorf("no document found with the given id")
	}

	err = mongodb.Coll.FindOne(ctx, filter).Decode(&example)
	if err != nil {
		return nil, err
	}
	response := pb.Account{
		Id:        example.ID.String(),
		UserId:    example.UserId,
		Name:      example.Name,
		Type:      example.Type,
		Balance:   example.Balance,
		Currency:  example.Currency,
		CreatedAt: example.CreatedAt.String(),
		UpdatedAt: example.UpdatedAt.String(),
	}

	return &response, nil
}

func (mongodb *AccountMongo) GetAccount(ctx context.Context, request *pb.PrimaryKey) (*pb.Account, error) {
	var (
		err     error
		example pb.Account
	)
	id, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, err
	}
	err = mongodb.Coll.FindOne(ctx, bson.M{"_id": id}).Decode(&example)
	if err != nil {
		return nil, err
	}
	response := pb.Account{
		UserId:    example.UserId,
		Name:      example.Name,
		Type:      example.Type,
		Balance:   example.Balance,
		Currency:  example.Currency,
		CreatedAt: example.CreatedAt,
		UpdatedAt: example.UpdatedAt,
	}
	return &response, nil
}
func (mongodb *AccountMongo) GetAllAccount(ctx context.Context, request *pb.GetListRequest) (*pb.Accounts, error) {
	var (
		accounts = []*pb.Account{}
		offset   = (request.GetPage() - 1) * request.GetLimit()
		count    int64
		err      error
		filter   = bson.M{}
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
		example := Account{}
		if err = cursor.Decode(&example); err != nil {
			return nil, err
		}
		account := pb.Account{
			UserId:    example.UserId,
			Name:      example.Name,
			Type:      example.Type,
			Balance:   example.Balance,
			Currency:  example.Currency,
			CreatedAt: example.CreatedAt.String(),
			UpdatedAt: example.UpdatedAt.String(),
		}

		accounts = append(accounts, &account)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &pb.Accounts{
		Accounts: accounts,
		Count:    int32(count),
	}, nil
}
func (mongodb *AccountMongo) DeleteAccount(ctx context.Context, request *pb.PrimaryKey) (*emptypb.Empty, error) {
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
