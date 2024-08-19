package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`        // ObjectID for MongoDB
	UserId      string             `bson:"user_id" json:"user_id"`         // User ID as a string
	AccountId   string             `bson:"account_id" json:"account_id"`   // Account ID as a string
	CategoryId  string             `bson:"category_id" json:"category_id"` // Category ID as a string
	Amount      float32            `bson:"amount" json:"amount"`           // Amount as a float32
	Type        string             `bson:"type" json:"type"`               // Type (spending/income) as a string
	Description string             `bson:"description" json:"description"` // Description as a string
	Date        time.Time          `bson:"date" json:"date"`               // Date as a time.Time
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`   // Created timestamp as a time.Time
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`   // Updated timestamp as a time.Time
}

type CreateTransaction struct {
	UserId      string    `bson:"user_id" json:"user_id"`         // User ID as a string
	AccountId   string    `bson:"account_id" json:"account_id"`   // Account ID as a string
	CategoryId  string    `bson:"category_id" json:"category_id"` // Category ID as a string
	Amount      float32   `bson:"amount" json:"amount"`           // Amount as a float32
	Type        string    `bson:"type" json:"type"`               // Type (spending/income) as a string
	Description string    `bson:"description" json:"description"` // Description as a string
	Date        time.Time `bson:"date" json:"date"`               // Date as a time.Time
}

type UpdateTransaction struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`        // ObjectID for MongoDB
	UserId      string             `bson:"user_id" json:"user_id"`         // User ID as a string
	AccountId   string             `bson:"account_id" json:"account_id"`   // Account ID as a string
	CategoryId  string             `bson:"category_id" json:"category_id"` // Category ID as a string
	Amount      float32            `bson:"amount" json:"amount"`           // Amount as a float32
	Type        string             `bson:"type" json:"type"`               // Type (spending/income) as a string
	Description string             `bson:"description" json:"description"` // Description as a string
	Date        time.Time          `bson:"date" json:"date"`               // Date as a time.Time
}

type TransactionsResponse struct {
	Transactions []Transaction `bson:"transactions" json:"transactions"` // List of transactions
	Count        int32         `bson:"count" json:"count"`               // Count of transactions
}

type GetUserMoneyRequest struct {
	UserId    string    `bson:"user_id" json:"user_id"`       // User ID as a string
	StartTime time.Time `bson:"start_time" json:"start_time"` // Start time as a time.Time
	EndTime   time.Time `bson:"end_time" json:"end_time"`     // End time as a time.Time
}

type GetUserMoneyResponse struct {
	CategoryId  string    `bson:"category_id" json:"category_id"`   // Category ID as a string
	TotalAmount float32   `bson:"total_amount" json:"total_amount"` // Total amount as a float32
	Time        time.Time `bson:"time" json:"time"`                 // Time as a time.Time
}
type GetUserMoneysResponse struct {
	Moneys []GetUserMoneyResponse `bson:"moneys" json:"moneys"` // List of user money responses
}
