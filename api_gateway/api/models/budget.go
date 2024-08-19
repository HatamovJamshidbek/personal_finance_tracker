package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	_ "go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

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

type CreateBudget struct {
	UserId     string  `bson:"user_id" json:"user_id"`
	CategoryId string  `bson:"category_id" json:"category_id"`
	Amount     float32 `bson:"amount" json:"amount"`
	Period     string  `bson:"period" json:"period"`
	StartTime  string  `bson:"start_time" json:"start_time"`
	EndTime    string  `bson:"end_time" json:"end_time"`
}

type UpdateBudget struct {
	ID         string  `bson:"-"`
	UserId     string  `bson:"user_id" json:"user_id"`
	CategoryId string  `bson:"category_id" json:"category_id"`
	Amount     float32 `bson:"amount" json:"amount"`
	Period     string  `bson:"period" json:"period"`
	StartTime  string  `bson:"start_time" json:"start_time"`
	EndTime    string  `bson:"end_time" json:"end_time"`
}

type BudgetsResponse struct {
	Budgets []Budget `bson:"budgets" json:"budgets"`
	Count   int32    `bson:"count" json:"count"`
}

type BudgetSummaryItem struct {
	CategoryId  string  `json:"category_id" example:"category1" description:"ID of the budget category"`
	TotalAmount float32 `json:"total_amount" example:"1234.56" description:"Total amount for the category"`
	StartTime   string  `json:"start_time" example:"2024-01-01T00:00:00Z" description:"Start time of the budget period"`
	EndTime     string  `json:"end_time" example:"2024-01-31T23:59:59Z" description:"End time of the budget period"`
	Period      string  `json:"period" example:"monthly" description:"Period of the budget (e.g., daily, weekly, monthly)"`
}

type GetUserBudgetResponse struct {
	Results []*BudgetSummaryItem `json:"results" description:"List of budget summary items"`
}
