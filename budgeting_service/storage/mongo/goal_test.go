package mongo

//
//import (
//	"context"
//	"go.mongodb.org/mongo-driver/bson"
//	"go.mongodb.org/mongo-driver/bson/primitive"
//	"testing"
//	"time"
//
//	pb "budgeting_service/genproto/budgeting_service"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//	"go.mongodb.org/mongo-driver/mongo/mongotest"
//)
//
//// Mock Logger
//type MockLogger struct {
//	mock.Mock
//}
//
//func (m *MockLogger) Error(message string, err error) {
//	m.Called(message, err)
//}
//
//func TestCreateGoal(t *testing.T) {
//	client, coll := mongotest.New(t)
//	mockLogger := new(MockLogger)
//	mongoDB := NewGoalMongo(coll, mockLogger)
//
//	req := &pb.GoalRequest{
//		UserId:        "user123",
//		Name:          "Save for Vacation",
//		TargetAmount:  1000.0,
//		CurrentAmount: 200.0,
//		Deadline:      time.Now().AddDate(0, 0, 30).Format(time.RFC3339),
//		Status:        "in_progress",
//	}
//
//	createdGoal, err := mongoDB.CreateGoal(context.Background(), req)
//	assert.NoError(t, err)
//	assert.NotNil(t, createdGoal)
//	assert.Equal(t, req.GetUserId(), createdGoal.GetUserId())
//	assert.Equal(t, req.GetName(), createdGoal.GetName())
//	assert.Equal(t, req.GetTargetAmount(), createdGoal.GetTargetAmount())
//	assert.Equal(t, req.GetCurrentAmount(), createdGoal.GetCurrentAmount())
//	assert.Equal(t, req.GetDeadline(), createdGoal.GetDeadline())
//	assert.Equal(t, req.GetStatus(), createdGoal.GetStatus())
//}
//func TestUpdateGoal(t *testing.T) {
//	client, coll := mongotest.New(t)
//	mockLogger := new(MockLogger)
//	mongoDB := NewGoalMongo(coll, mockLogger)
//
//	goalID := primitive.NewObjectID()
//	_, _ = coll.InsertOne(context.Background(), bson.M{
//		"_id":            goalID,
//		"user_id":        "user123",
//		"name":           "Save for Vacation",
//		"target_amount":  1000.0,
//		"current_amount": 200.0,
//		"deadline":       time.Now().AddDate(0, 0, 30).Format(time.RFC3339),
//		"status":         "in_progress",
//		"created_at":     time.Now(),
//		"updated_at":     time.Now(),
//	})
//
//	req := &pb.Goal{
//		Id:            goalID.Hex(),
//		UserId:        "user123",
//		Name:          "Updated Goal Name",
//		TargetAmount:  1200.0,
//		CurrentAmount: 250.0,
//		Deadline:      time.Now().AddDate(0, 0, 60).Format(time.RFC3339),
//		Status:        "completed",
//	}
//
//	updatedGoal, err := mongoDB.UpdateGoal(context.Background(), req)
//	assert.NoError(t, err)
//	assert.NotNil(t, updatedGoal)
//	assert.Equal(t, req.GetId(), updatedGoal.GetId())
//	assert.Equal(t, req.GetName(), updatedGoal.GetName())
//	assert.Equal(t, req.GetTargetAmount(), updatedGoal.GetTargetAmount())
//	assert.Equal(t, req.GetCurrentAmount(), updatedGoal.GetCurrentAmount())
//	assert.Equal(t, req.GetDeadline(), updatedGoal.GetDeadline())
//	assert.Equal(t, req.GetStatus(), updatedGoal.GetStatus())
//}
//func TestGetGoal(t *testing.T) {
//	client, coll := mongotest.New(t)
//	mockLogger := new(MockLogger)
//	mongoDB := NewGoalMongo(coll, mockLogger)
//
//	goalID := primitive.NewObjectID()
//	_, _ = coll.InsertOne(context.Background(), bson.M{
//		"_id":            goalID,
//		"user_id":        "user123",
//		"name":           "Save for Vacation",
//		"target_amount":  1000.0,
//		"current_amount": 200.0,
//		"deadline":       time.Now().AddDate(0, 0, 30).Format(time.RFC3339),
//		"status":         "in_progress",
//		"created_at":     time.Now(),
//		"updated_at":     time.Now(),
//	})
//
//	req := &pb.PrimaryKey{
//		Id: goalID.Hex(),
//	}
//
//	goal, err := mongoDB.GetGoal(context.Background(), req)
//	assert.NoError(t, err)
//	assert.NotNil(t, goal)
//	assert.Equal(t, goalID.Hex(), goal.GetId())
//}
//func TestGetAllGoal(t *testing.T) {
//	client, coll := mongotest.New(t)
//	mockLogger := new(MockLogger)
//	mongoDB := NewGoalMongo(coll, mockLogger)
//
//	_, _ = coll.InsertMany(context.Background(), []interface{}{
//		bson.M{
//			"user_id":        "user123",
//			"name":           "Goal 1",
//			"target_amount":  1000.0,
//			"current_amount": 200.0,
//			"deadline":       time.Now().AddDate(0, 0, 30).Format(time.RFC3339),
//			"status":         "in_progress",
//			"created_at":     time.Now(),
//			"updated_at":     time.Now(),
//		},
//		bson.M{
//			"user_id":        "user123",
//			"name":           "Goal 2",
//			"target_amount":  1500.0,
//			"current_amount": 300.0,
//			"deadline":       time.Now().AddDate(0, 0, 60).Format(time.RFC3339),
//			"status":         "in_progress",
//			"created_at":     time.Now(),
//			"updated_at":     time.Now(),
//		},
//	})
//
//	req := &pb.GetListRequest{
//		Page:  1,
//		Limit: 10,
//	}
//
//	goals, err := mongoDB.GetAllGoal(context.Background(), req)
//	assert.NoError(t, err)
//	assert.NotNil(t, goals)
//	assert.Equal(t, 2, len(goals.GetGoals()))
//}
//func TestDeleteGoal(t *testing.T) {
//	client, coll := mongotest.New(t)
//	mockLogger := new(MockLogger)
//	mongoDB := NewGoalMongo(coll, mockLogger)
//
//	goalID := primitive.NewObjectID()
//	_, _ = coll.InsertOne(context.Background(), bson.M{
//		"_id":            goalID,
//		"user_id":        "user123",
//		"name":           "Goal to Delete",
//		"target_amount":  1000.0,
//		"current_amount": 200.0,
//		"deadline":       time.Now().AddDate(0, 0, 30).Format(time.RFC3339),
//		"status":         "in_progress",
//		"created_at":     time.Now(),
//		"updated_at":     time.Now(),
//	})
//
//	req := &pb.PrimaryKey{
//		Id: goalID.Hex(),
//	}
//
//	_, err := mongoDB.DeleteGoal(context.Background(), req)
//	assert.NoError(t, err)
//
//	var result bson.M
//	err = coll.FindOne(context.Background(), bson.M{"_id": goalID}).Decode(&result)
//	assert.Error(t, err)
//}
