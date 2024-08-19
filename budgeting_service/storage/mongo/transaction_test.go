package mongo

//
//import (
//	"context"
//	"go.mongodb.org/mongo-driver/bson"
//	"testing"
//	"time"
//
//	pb "budgeting_service/genproto/budgeting_service"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//	"go.mongodb.org/mongo-driver/bson/primitive"
//	"go.mongodb.org/mongo-driver/mongo"
//)
//
//type MockLogger struct {
//	mock.Mock
//}
//
//func (m *MockLogger) Error(msg string, args ...interface{}) {
//	m.Called(msg, args)
//}
//
//func TestCreateTransaction(t *testing.T) {
//	ctx := context.Background()
//	mockColl := &mongo.Collection{}
//	mockLog := &MockLogger{}
//	mongoService := NewTransactionMongo(mockColl, mockLog)
//
//	// Prepare test data
//	request := &pb.TransactionRequest{
//		UserId:      "user123",
//		AccountId:   "account123",
//		CategoryId:  "category123",
//		Amount:      100.0,
//		Type:        "income",
//		Description: "Test transaction",
//		Date:        time.Now().Format(time.RFC3339),
//	}
//
//	// Mock InsertOne
//	mockColl.On("InsertOne", ctx, mock.Anything).Return(&mongo.InsertOneResult{
//		InsertedID: primitive.NewObjectID(),
//	}, nil)
//
//	// Execute method
//	response, err := mongoService.CreateTransaction(ctx, request)
//
//	// Assertions
//	assert.NoError(t, err)
//	assert.NotNil(t, response)
//	mockColl.AssertExpectations(t)
//}
//
//func TestUpdateTransaction(t *testing.T) {
//	ctx := context.Background()
//	mockColl := &mongo.Collection{}
//	mockLog := &MockLogger{}
//	mongoService := NewTransactionMongo(mockColl, mockLog)
//
//	// Prepare test data
//	request := &pb.Transaction{
//		Id:          "606c72ef4d1d3b1a5c18e4d1",
//		UserId:      "user123",
//		AccountId:   "account123",
//		CategoryId:  "category123",
//		Amount:      150.0,
//		Type:        "income",
//		Description: "Updated transaction",
//		Date:        time.Now().Format(time.RFC3339),
//	}
//
//	// Mock UpdateOne
//	mockColl.On("UpdateOne", ctx, mock.Anything, mock.Anything).Return(&mongo.UpdateResult{
//		MatchedCount: 1,
//	}, nil)
//	mockColl.On("FindOne", ctx, mock.Anything).Return(&mongo.SingleResult{}, nil).Run(func(args mock.Arguments) {
//		mockArgs := args.Get(1).(*pb.Transaction)
//		mockArgs.Id = request.Id
//		mockArgs.UserId = request.UserId
//		mockArgs.AccountId = request.AccountId
//		mockArgs.CategoryId = request.CategoryId
//		mockArgs.Amount = request.Amount
//		mockArgs.Type = request.Type
//		mockArgs.Description = request.Description
//		mockArgs.Date = request.Date
//		mockArgs.CreatedAt = time.Now().Format(time.RFC3339)
//		mockArgs.UpdatedAt = time.Now().Format(time.RFC3339)
//	})
//
//	// Execute method
//	response, err := mongoService.UpdateTransaction(ctx, request)
//
//	// Assertions
//	assert.NoError(t, err)
//	assert.NotNil(t, response)
//	mockColl.AssertExpectations(t)
//}
//
//func TestGetTransaction(t *testing.T) {
//	mockColl := new(MockCollection)
//	transactionID := primitive.NewObjectID()
//	expectedTransaction := &mongo.Transaction{
//		ID:          transactionID,
//		UserId:      "user1",
//		AccountId:   "account1",
//		CategoryId:  "category1",
//		Amount:      100.0,
//		Type:        "expense",
//		Description: "Test transaction",
//		Date:        time.Now(),
//		CreatedAt:   time.Now(),
//		UpdatedAt:   time.Now(),
//	}
//
//	mockColl.On("FindOne", mock.Anything, bson.M{"_id": transactionID}).Return(mongo.SingleResult{DecodeResult: expectedTransaction})
//
//	transactionMongo := &mongo.TransactionMongo{Coll: mockColl}
//
//	req := &budgeting_service.PrimaryKey{Id: transactionID.Hex()}
//	res, err := transactionMongo.GetTransaction(context.Background(), req)
//
//	assert.NoError(t, err)
//	assert.Equal(t, expectedTransaction.ID.Hex(), res.Id)
//	assert.Equal(t, expectedTransaction.UserId, res.UserId)
//}
//
//func TestGetAllTransaction(t *testing.T) {
//	mockColl := new(MockCollection)
//	transactionID := primitive.NewObjectID()
//	expectedTransactions := []*mongo.Transaction{
//		{
//			ID:          transactionID,
//			UserId:      "user1",
//			AccountId:   "account1",
//			CategoryId:  "category1",
//			Amount:      100.0,
//			Type:        "expense",
//			Description: "Test transaction",
//			Date:        time.Now(),
//			CreatedAt:   time.Now(),
//			UpdatedAt:   time.Now(),
//		},
//	}
//
//	cursor := mongotest.NewMockCursor()
//	for _, tx := range expectedTransactions {
//		cursor.Append(tx)
//	}
//
//	mockColl.On("Find", mock.Anything, bson.M{}).Return(cursor, nil)
//	mockColl.On("CountDocuments", mock.Anything, bson.M{}).Return(int64(len(expectedTransactions)), nil)
//
//	transactionMongo := &mongo.TransactionMongo{Coll: mockColl}
//
//	req := &budgeting_service.GetListRequest{Page: 1, Limit: 10}
//	res, err := transactionMongo.GetAllTransaction(context.Background(), req)
//
//	assert.NoError(t, err)
//	assert.Len(t, res.Transactions, len(expectedTransactions))
//}
//
//func TestDeleteTransaction(t *testing.T) {
//	mockColl := new(MockCollection)
//	transactionID := primitive.NewObjectID()
//
//	mockColl.On("DeleteOne", mock.Anything, bson.M{"_id": transactionID}).Return(&mongo.DeleteResult{DeletedCount: 1}, nil)
//
//	transactionMongo := &mongo.TransactionMongo{Coll: mockColl}
//
//	req := &budgeting_service.PrimaryKey{Id: transactionID.Hex()}
//	_, err := transactionMongo.DeleteTransaction(context.Background(), req)
//
//	assert.NoError(t, err)
//}
