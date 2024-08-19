package mongo

//
//import (
//	pb "budgeting_service/genproto/budgeting_service"
//	"budgeting_service/pkg/logger"
//	"context"
//	"github.com/stretchr/testify/assert"
//	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
//	"testing"
//)
//
//func TestCreateCategory(t *testing.T) {
//	mt := mtest.New(t, mtest.NewOptions().ClientOptions(mongoOptions))
//	defer mt.Close()
//
//	categoryColl := mt.Coll("categories")
//	catMongo := NewCategoryMongo(categoryColl, logger.NewLogger())
//
//	req := &pb.CategoryRequest{
//		UserId: "user1",
//		Name:   "Food",
//		Type:   "Expense",
//	}
//
//	createdCategory, err := catMongo.CreateCategory(context.Background(), req)
//	assert.NoError(t, err)
//	assert.NotNil(t, createdCategory)
//	assert.Equal(t, req.GetUserId(), createdCategory.GetUserId())
//}
//
//func TestUpdateCategory(t *testing.T) {
//	mt := mtest.New(t, mtest.NewOptions().ClientOptions(mongoOptions))
//	defer mt.Close()
//
//	categoryColl := mt.Coll("categories")
//	catMongo := NewCategoryMongo(categoryColl, logger.NewLogger())
//
//	// Create a category to update
//	createReq := &pb.CategoryRequest{
//		UserId: "user2",
//		Name:   "Travel",
//		Type:   "Expense",
//	}
//	category, _ := catMongo.CreateCategory(context.Background(), createReq)
//
//	// Update the category
//	updateReq := &pb.Category{
//		Id:     category.Id,
//		UserId: "user2",
//		Name:   "Travel Updated",
//		Type:   "Income",
//	}
//
//	updatedCategory, err := catMongo.UpdateCategory(context.Background(), updateReq)
//	assert.NoError(t, err)
//	assert.Equal(t, "Travel Updated", updatedCategory.GetName())
//	assert.Equal(t, "Income", updatedCategory.GetType())
//}
//
//func TestGetCategory(t *testing.T) {
//	mt := mtest.New(t, mtest.NewOptions().ClientOptions(mongoOptions))
//	defer mt.Close()
//
//	categoryColl := mt.Coll("categories")
//	catMongo := NewCategoryMongo(categoryColl, logger.NewLogger())
//
//	// Create a category to retrieve
//	createReq := &pb.CategoryRequest{
//		UserId: "user3",
//		Name:   "Books",
//		Type:   "Expense",
//	}
//	category, _ := catMongo.CreateCategory(context.Background(), createReq)
//
//	// Retrieve the category
//	getReq := &pb.PrimaryKey{Id: category.Id}
//	retrievedCategory, err := catMongo.GetCategory(context.Background(), getReq)
//	assert.NoError(t, err)
//	assert.Equal(t, category.Id, retrievedCategory.Id)
//}
//
//func TestGetAllCategory(t *testing.T) {
//	mt := mtest.New(t, mtest.NewOptions().ClientOptions(mongoOptions))
//	defer mt.Close()
//
//	categoryColl := mt.Coll("categories")
//	catMongo := NewCategoryMongo(categoryColl, logger.NewLogger())
//
//	// Create categories
//	catMongo.CreateCategory(context.Background(), &pb.CategoryRequest{UserId: "user4", Name: "Sports", Type: "Expense"})
//	catMongo.CreateCategory(context.Background(), &pb.CategoryRequest{UserId: "user4", Name: "Music", Type: "Income"})
//
//	// Retrieve all categories
//	getReq := &pb.GetListRequest{Page: 1, Limit: 10}
//	categories, err := catMongo.GetAllCategory(context.Background(), getReq)
//	assert.NoError(t, err)
//	assert.Greater(t, len(categories.GetCategories()), 0)
//}
//
//func TestDeleteCategory(t *testing.T) {
//	mt := mtest.New(t, mtest.NewOptions().ClientOptions(mongoOptions))
//	defer mt.Close()
//
//	categoryColl := mt.Coll("categories")
//	catMongo := NewCategoryMongo(categoryColl, logger.NewLogger())
//
//	// Create a category to delete
//	createReq := &pb.CategoryRequest{
//		UserId: "user5",
//		Name:   "Entertainment",
//		Type:   "Expense",
//	}
//	category, _ := catMongo.CreateCategory(context.Background(), createReq)
//
//	// Delete the category
//	deleteReq := &pb.PrimaryKey{Id: category.Id}
//	_, err := catMongo.DeleteCategory(context.Background(), deleteReq)
//	assert.NoError(t, err)
//
//	// Verify deletion
//	_, err = catMongo.GetCategory(context.Background(), deleteReq)
//	assert.Error(t, err)
//}
//
//func TestCategoryMongo_GetUserBudgetSummary(t *testing.T) {
//	mt := mtest.New(t, mtest.NewOptions().ClientOptions(mongoOptions))
//	defer mt.Close()
//
//	categoryColl := mt.Coll("categories")
//	catMongo := NewCategoryMongo(categoryColl, logger.NewLogger())
//
//	// Create categories
//	catMongo.CreateCategory(context.Background(), &pb.CategoryRequest{UserId: "user6", Name: "Travel", Type: "Expense"})
//	catMongo.CreateCategory(context.Background(), &pb.CategoryRequest{UserId: "user6", Name: "Food", Type: "Expense"})
//
//	// Get user budget summary
//	getReq := &pb.PrimaryKey{Id: "user6"}
//	summary, err := catMongo.GetUserBudgetSummary(context.Background(), getReq)
//	assert.NoError(t, err)
//	assert.Greater(t, len(summary.Results), 0)
//}
