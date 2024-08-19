package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Account struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserId    string             `bson:"user_id" json:"user_id"`
	Name      string             `bson:"name" json:"name"`
	Type      string             `bson:"type" json:"type"`
	Balance   float32            `bson:"balance" json:"balance"`
	Currency  string             `bson:"currency" json:"currency"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type CreateAccount struct {
	UserId   string  `bson:"user_id" json:"user_id"`
	Name     string  `bson:"name" json:"name"`
	Type     string  `bson:"type" json:"type"`
	Balance  float32 `bson:"balance" json:"balance"`
	Currency string  `bson:"currency" json:"currency"`
}

type UpdateAccount struct {
	ID       string  `bson:"-"`
	UserId   string  `bson:"user_id" json:"user_id"`
	Name     string  `bson:"name" json:"name"`
	Type     string  `bson:"type" json:"type"`
	Balance  float32 `bson:"balance" json:"balance"`
	Currency string  `bson:"currency" json:"currency"`
}

type AccountsResponse struct {
	Accounts []Account `bson:"accounts" json:"accounts" `
	Count    int32     `bson:"count" json:"count"`
}
