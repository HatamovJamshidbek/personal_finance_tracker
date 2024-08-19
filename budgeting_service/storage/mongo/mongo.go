package mongo

import (
	"budgeting_service/configs"
	"budgeting_service/pkg/logger"
	"budgeting_service/storage"
	"budgeting_service/storage/redis"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	client    *mongo.Client
	db        *mongo.Database
	cfg       configs.Config
	log       logger.ILogger
	redisCash *redis.RedisRepo
}

// Accounts returns a new instance of AccountMongo with necessary dependencies
func (s *Store) Accounts() storage.IAccountStorage {
	return NewAccountMongoStore(s.db.Collection("account"), s.log, s.redisCash)
}

func (s *Store) Categories() storage.ICategoryStorage {
	return NewCategoryMongoStore(s.db.Collection("category"), s.log)
}

func (s *Store) Budgets() storage.IBudgetStorage {
	return NewBudgetMongoStore(s.db.Collection("budget"), s.log)
}

func (s *Store) Goals() storage.IGoalStorage {
	return NewGoalMongoStore(s.db.Collection("goal"), s.log)
}

func (s *Store) Transactions() storage.ITransactionStorage {
	return NewTransactionMongoStore(s.db.Collection("transaction"), s.log)
}

func NewStore(ctx context.Context, cfg configs.Config, log logger.ILogger, redisCash *redis.RedisRepo) (*Store, error) {
	uri := fmt.Sprintf(
		`mongodb://%s:%s@%s:%s/%s?authSource=admin&authMechanism=SCRAM-SHA-256`,
		cfg.MongoUser,
		cfg.MongoPassword,
		cfg.MongoHost,
		cfg.MongoPort,
		cfg.MongoDB,
	)

	clientOptions := options.Client().ApplyURI(uri).
		SetMaxPoolSize(100).
		SetConnectTimeout(10 * time.Second)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Error("Failed to connect to MongoDB", logger.Error(err))
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Error("Failed to ping MongoDB", logger.Error(err))
		return nil, err
	}

	db := client.Database(cfg.MongoDB)

	return &Store{
		client:    client,
		db:        db,
		cfg:       cfg,
		log:       log,
		redisCash: redisCash,
	}, nil
}

// Close disconnects the MongoDB client
func (s *Store) Close() {
	if err := s.client.Disconnect(context.Background()); err != nil {
		s.log.Error("Error disconnecting from MongoDB", logger.Error(err))
	}
}
