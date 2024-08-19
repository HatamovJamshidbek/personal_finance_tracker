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
	"time"
)

type Category struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" bson:"_id"`
	UserId    string             `bson:"user_id" json:"user_id"`
	Name      string             `bson:"name" json:"name"`
	Type      string             `bson:"type" json:"type"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type CategoryMongo struct {
	Coll *mongo.Collection
	log  logger.ILogger
}

func NewCategoryMongoStore(db *mongo.Collection, lg logger.ILogger) *CategoryMongo {
	return &CategoryMongo{
		Coll: db,
		log:  lg,
	}
}

func (mongodb *CategoryMongo) CreateCategory(ctx context.Context, request *pb.CategoryRequest) (*pb.Category, error) {
	var (
		err         error
		example     Category
		currentTime = time.Now()
	)

	budget := bson.D{
		{Key: "user_id", Value: request.GetUserId()},
		{Key: "name", Value: request.GetName()},
		{Key: "type", Value: request.GetType()},
		{Key: "created_at", Value: currentTime},
	}
	result, err := mongodb.Coll.InsertOne(ctx, &budget)
	if err != nil {
		mongodb.log.Error("this error is insert one ", logger.Error(err))
		return nil, err
	}
	err = mongodb.Coll.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&example)
	if err != nil {
		return nil, err
	}
	response := pb.Category{
		Id:        example.ID.String(),
		UserId:    example.UserId,
		Name:      example.Name,
		Type:      example.Type,
		CreatedAt: example.CreatedAt.String(),
		UpdatedAt: example.UpdatedAt.String(),
	}
	return &response, nil

}
func (mongodb *CategoryMongo) UpdateCategory(ctx context.Context, request *pb.Category) (*pb.Category, error) {
	var (
		params      = bson.M{}
		example     Category
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

		params["updated_at"] = currentTime
	}

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
		return nil, fmt.Errorf("this budget is not found than one")
	}

	err = mongodb.Coll.FindOne(ctx, filter).Decode(&example)
	if err != nil {
		return nil, err
	}
	response := pb.Category{
		Id:        example.ID.String(),
		UserId:    example.UserId,
		Name:      example.Name,
		Type:      example.Type,
		CreatedAt: example.CreatedAt.String(),
		UpdatedAt: example.UpdatedAt.String(),
	}
	return &response, nil
}

func (mongodb *CategoryMongo) GetCategory(ctx context.Context, request *pb.PrimaryKey) (*pb.Category, error) {
	var (
		err     error
		example Category
	)
	id, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		return nil, err
	}
	err = mongodb.Coll.FindOne(ctx, bson.M{"_id": id}).Decode(&example)
	if err != nil {
		return nil, err
	}
	response := pb.Category{
		Id:        example.ID.String(),
		UserId:    example.UserId,
		Name:      example.Name,
		Type:      example.Type,
		CreatedAt: example.CreatedAt.String(),
		UpdatedAt: example.UpdatedAt.String(),
	}
	return &response, nil
}
func (mongodb *CategoryMongo) GetAllCategory(ctx context.Context, request *pb.GetListRequest) (*pb.Categories, error) {
	var (
		categories = []*pb.Category{}
		offset     = (request.GetPage() - 1) * request.GetLimit()
		count      int64
		err        error
		filter     = bson.M{}
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
		var example Category
		if err = cursor.Decode(&example); err != nil {
			return nil, err
		}
		category := pb.Category{
			Id:        example.ID.String(),
			UserId:    example.UserId,
			Name:      example.Name,
			Type:      example.Type,
			CreatedAt: example.CreatedAt.String(),
			UpdatedAt: example.UpdatedAt.String(),
		}

		categories = append(categories, &category)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &pb.Categories{
		Categories: categories,
		Count:      int32(count),
	}, nil
}
func (mongodb *CategoryMongo) DeleteCategory(ctx context.Context, request *pb.PrimaryKey) (*emptypb.Empty, error) {
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
