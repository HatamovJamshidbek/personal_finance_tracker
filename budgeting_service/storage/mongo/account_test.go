package mongo_test

//
//import (
//	"context"
//	"fmt"
//	"go.mongodb.org/mongo-driver/mongo"
//	"go.mongodb.org/mongo-driver/mongo/options"
//	"testing"
//
//	"budgeting_service/genproto/budgeting_service"
//	"budgeting_service/pkg/logger"
//	"budgeting_service/storage/redis"
//	"github.com/stretchr/testify/assert"
//)
//
//func setupTestDB(t *testing.T) *mongo.Collection {
//	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
//	client, err := mongo.Connect(context.Background(), clientOptions)
//	if err != nil {
//		t.Fatalf("Failed to connect to MongoDB: %v", err)
//	}
//
//	db := client.Database("testdb")
//	return db.Collection("accounts")
//}
//
//func TestCreateAccount(t *testing.T) {
//	collection := setupTestDB(t)
//	defer collection.Drop(context.Background())
//
//	logger := logger.NewLogger()       // Mock logger if needed
//	redisCache := redis.NewRedisRepo() // Mock Redis repo if needed
//	accountMongo := mongo(collection, logger, redisCache)
//
//	request := &budgeting_service.AccountRequest{
//		Name:     "Test Account",
//		Type:     "Savings",
//		Balance:  1000.0,
//		Currency: "USD",
//	}
//
//	resp, err := accountMongo.CreateAccount(context.Background(), request)
//	assert.NoError(t, err)
//	assert.NotNil(t, resp)
//	assert.Equal(t, request.Name, resp.Name)
//	assert.Equal(t, request.Type, resp.Type)
//	assert.Equal(t, request.Balance, resp.Balance)
//	assert.Equal(t, request.Currency, resp.Currency)
//}
//func TestUpdateAccount(t *testing.T) {
//	collection := setupTestDB(t)
//	defer collection.Drop(context.Background())
//
//	logger := logger.NewLogger() // Mock logger if needed
//	redisCache := redis.NewRedisRepo() // Mock Redis repo if needed
//
//	accountMongo := mongo.NewAccountMongo(collection, logger, redisCache)
//
//	// First, create an account to update
//	createReq := &budgeting_service.AccountRequest{
//		Name:      "Test Account",
//		Type:      "Savings",
//		Balance:   1000.0,
//		Currency:  "USD",
//	}
//	createResp, err := accountMongo.CreateAccount(context.Background(), createReq)
//	assert.NoError(t, err)
//
//	// Now update the account
//	updateReq := &budgeting_service.Account{
//		Id:        createResp.Id,
//		Name:      "Updated Account",
//		Balance:   2000.0,
//		Currency:  "EUR",
//	}
//
//	updateResp, err := accountMongo.UpdateAccount(context.Background(), updateReq)
//	assert.NoError(t, err)
//	assert.NotNil(t, updateResp)
//	assert.Equal(t, updateReq.Name, updateResp.Name)
//	assert.Equal(t, updateReq.Balance, updateResp.Balance)
//	assert.Equal(t, updateReq.Currency, updateResp.Currency)
//}
//
//func TestGetAccount(t *testing.T) {
//	collection := setupTestDB(t)
//	defer collection.Drop(context.Background())
//
//	logger := logger.NewLogger() // Mock logger if needed
//	redisCache := redis.NewRedisRepo() // Mock Redis repo if needed
//
//	accountMongo := mongo.NewAccountMongo(collection, logger, redisCache)
//
//	// First, create an account to retrieve
//	createReq := &budgeting_service.AccountRequest{
//		Name:      "Test Account",
//		Type:      "Savings",
//		Balance:   1000.0,
//		Currency:  "USD",
//	}
//	createResp, err := accountMongo.CreateAccount(context.Background(), createReq)
//	assert.NoError(t, err)
//
//	// Retrieve the account
//	getReq := &budgeting_service.PrimaryKey{Id: createResp.Id}
//	getResp, err := accountMongo.GetAccount(context.Background(), getReq)
//	assert.NoError(t, err)
//	assert.NotNil(t, getResp)
//	assert.Equal(t, createResp.Id, getResp.Id)
//	assert.Equal(t, createResp.Name, getResp.Name)
//}
//
//func TestGetAllAccount(t *testing.T) {
//	collection := setupTestDB(t)
//	defer collection.Drop(context.Background())
//
//	logger := logger.NewLogger() // Mock logger if needed
//	redisCache := redis.NewRedisRepo() // Mock Redis repo if needed
//
//	accountMongo := mongo.NewAccountMongo(collection, logger, redisCache)
//
//	// Create some accounts
//	for i := 0; i < 5; i++ {
//		_, err := accountMongo.CreateAccount(context.Background(), &budgeting_service.AccountRequest{
//			Name:      fmt.Sprintf("Test Account %d", i),
//			Type:      "Savings",
//			Balance:   float32(i * 100),
//			Currency:  "USD",
//		})
//		assert.NoError(t, err)
//	}
//
//	// Retrieve all accounts
//	listReq := &budgeting_service.GetListRequest{
//		Page:  1,
//		Limit: 10,
//	}
//	listResp, err := accountMongo.GetAllAccount(context.Background(), listReq)
//	assert.NoError(t, err)
//	assert.NotNil(t, listResp)
//	assert.Equal(t, int32(5), listResp.Count)
//}
//
//func TestDeleteAccount(t *testing.T) {
//	collection := setupTestDB(t)
//	defer collection.Drop(context.Background())
//
//	logger := logger.NewLogger() // Mock logger if needed
//	redisCache := redis.NewRedisRepo() // Mock Redis repo if needed
//
//	accountMongo := mongo.NewAccountMongo(collection, logger, redisCache)
//
//	// First, create an account to delete
//	createReq := &budgeting_service.AccountRequest{
//		Name:      "Test Account",
//		Type:      "Savings",
//		Balance:   1000.0,
//		Currency:  "USD",
//	}
//	createResp, err := accountMongo.CreateAccount(context.Background(), createReq)
//	assert.NoError(t, err)
//
//	// Delete the account
//	deleteReq := &budgeting_service.PrimaryKey{Id: createResp.Id}
//	_, err = accountMongo.DeleteAccount(context.Background(), deleteReq)
//	assert.NoError(t, err)
//
//	// Verify the account was deleted
//	getResp, err := accountMongo.GetAccount(context.Background(), deleteReq)
//	assert.Error(t, err)
//	assert.Nil(t, getResp)
//}
//
