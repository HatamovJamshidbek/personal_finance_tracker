package mongo

//
//import (
//	pb "budgeting_service/genproto/budgeting_service"
//	"context"
//	"fmt"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//	"go.mongodb.org/mongo-driver/mongo"
//	"go.mongodb.org/mongo-driver/mongo/options"
//	"time"
//)
//
//type MockLogger struct {
//	mock.Mock
//}
//
//func (m *MockLogger) Error(msg string, err error) {
//	m.Called(msg, err)
//}
//
//func setupTest() (*BudgetMongo, *mongo.Collection, func()) {
//	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
//	client, err := mongo.Connect(context.Background(), clientOptions)
//	if err != nil {
//		panic(err)
//	}
//
//	db := client.Database("testdb")
//	coll := db.Collection("budgets")
//
//	return NewBudgetMongo(coll, &MockLogger{}), coll, func() {
//		_ = client.Disconnect(context.Background())
//	}
//}
//func TestUpdateBudget(t *testing.T) {
//	budgetMongo, coll, teardown := setupTest()
//	defer teardown()
//
//	// Create a budget first
//	request := &pb.BudgetRequest{
//		UserId:     "user1",
//		CategoryId: "cat1",
//		Amount:     100.0,
//		Period:     "monthly",
//		StartTime:  time.Now().Format(time.RFC3339),
//		EndTime:    time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339),
//	}
//	_, err := budgetMongo.CreateBudget(context.Background(), request)
//	assert.NoError(t, err)
//
//	// Update the budget
//	budgetID := "inserted_budget_id_here" // Replace with the actual ID from CreateBudget
//	updateRequest := &pb.Budget{
//		Id:        budgetID,
//		Amount:    150.0,
//		StartTime: time.Now().Add(-10 * 24 * time.Hour).Format(time.RFC3339),
//		EndTime:   time.Now().Add(40 * 24 * time.Hour).Format(time.RFC3339),
//	}
//	resp, err := budgetMongo.UpdateBudget(context.Background(), updateRequest)
//	assert.NoError(t, err)
//	assert.NotNil(t, resp)
//	assert.Equal(t, updateRequest.GetAmount(), resp.GetAmount())
//}
//func TestGetBudget(t *testing.T) {
//	budgetMongo, coll, teardown := setupTest()
//	defer teardown()
//
//	// Create a budget first
//	request := &pb.BudgetRequest{
//		UserId:     "user1",
//		CategoryId: "cat1",
//		Amount:     100.0,
//		Period:     "monthly",
//		StartTime:  time.Now().Format(time.RFC3339),
//		EndTime:    time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339),
//	}
//	createResp, err := budgetMongo.CreateBudget(context.Background(), request)
//	assert.NoError(t, err)
//
//	// Get the budget
//	getRequest := &pb.PrimaryKey{Id: createResp.GetId()}
//	resp, err := budgetMongo.GetBudget(context.Background(), getRequest)
//	assert.NoError(t, err)
//	assert.NotNil(t, resp)
//	assert.Equal(t, createResp.GetId(), resp.GetId())
//}
//
//func TestDeleteBudget(t *testing.T) {
//	budgetMongo, coll, teardown := setupTest()
//	defer teardown()
//
//	// Create a budget first
//	request := &pb.BudgetRequest{
//		UserId:     "user1",
//		CategoryId: "cat1",
//		Amount:     100.0,
//		Period:     "monthly",
//		StartTime:  time.Now().Format(time.RFC3339),
//		EndTime:    time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339),
//	}
//	createResp, err := budgetMongo.CreateBudget(context.Background(), request)
//	assert.NoError(t, err)
//
//	// Delete the budget
//	deleteRequest := &pb.PrimaryKey{Id: createResp.GetId()}
//	_, err = budgetMongo.DeleteBudget(context.Background(), deleteRequest)
//	assert.NoError(t, err)
//
//	// Verify deletion
//	resp, err := budgetMongo.GetBudget(context.Background(), deleteRequest)
//	assert.Error(t, err)
//	assert.Nil(t, resp)
//}
//
//func TestGetUserBudgetSummary(t *testing.T) {
//	budgetMongo, coll, teardown := setupTest()
//	defer teardown()
//
//	// Create some budgets
//	for i := 0; i < 3; i++ {
//		_, err := budgetMongo.CreateBudget(context.Background(), &pb.BudgetRequest{
//			UserId:     "user1",
//			CategoryId: fmt.Sprintf("cat%d", i),
//			Amount:     float32(100 + i*10),
//			Period:     "monthly",
//			StartTime:  time.Now().Format(time.RFC3339),
//			EndTime:    time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339),
//		})
//		assert.NoError(t, err)
//	}
//
//	// Get user budget summary
//	summaryRequest := &pb.PrimaryKey{Id: "user1"}
//	resp, err := budgetMongo.GetUserBudgetSummary(context.Background(), summaryRequest)
//	assert.NoError(t, err)
//	assert.NotNil(t, resp)
//	assert.Len(t, resp.GetResults(), 3)
//}
