package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Category struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserId    string             `bson:"user_id" json:"user_id"`
	Name      string             `bson:"name"  json:"name"`
	Type      string             `bson:"type" json:"type"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type CreateCategory struct {
	UserId string `bson:"user_id" json:"user_id"`
	Name   string `bson:"name"  json:"name"`
	Type   string `bson:"type" json:"type"`
}

type UpdateCategory struct {
	ID        string    `bson:"-"`
	UserId    string    `bson:"user_id" json:"user_id"`
	Name      string    `bson:"name"  json:"name"`
	Type      string    `bson:"type" json:"type"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type CategoriesResponse struct {
	Categories []Category `bson:"categories" json:"categories"`
	Count      int32      `bson:"count" json:"count"`
}
