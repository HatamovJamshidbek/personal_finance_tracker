package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Goal struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserId        string             `bson:"user_id" json:"user_id"`
	Name          string             `bson:"name" json:"name"`
	Type          string             `bson:"type" json:"type"`
	TargetAmount  float32            `bson:"target_amount" json:"target_amount"`
	CurrentAmount float32            `bson:"current_amount" json:"current_amount"`
	Deadline      string             `json:"deadline" bson:"deadline"`
	Status        string             `bson:"status" json:"status"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}

type CreateGoal struct {
	UserId        string  `bson:"user_id" json:"user_id"`
	Name          string  `bson:"name" json:"name"`
	Type          string  `bson:"type" json:"type"`
	TargetAmount  float32 `bson:"target_amount" json:"target_amount"`
	CurrentAmount float32 `bson:"current_amount" json:"current_amount"`
	Deadline      string  `json:"deadline" bson:"deadline" json:"deadline"`
	Status        string  `bson:"status" json:"status" json:"status"`
}

type UpdateGoal struct {
	ID            string  `bson:"-"`
	UserId        string  `bson:"user_id" json:"user_id"`
	Name          string  `bson:"name" json:"name"`
	Type          string  `bson:"type" json:"type"`
	TargetAmount  float32 `bson:"target_amount" json:"target_amount"`
	CurrentAmount float32 `bson:"current_amount" json:"current_amount"`
	Deadline      string  `json:"deadline" bson:"deadline" json:"deadline"`
	Status        string  `bson:"status" json:"status" json:"status"`
}

type GoalsResponse struct {
	Goals []Goal `bson:"goals" json:"goals"`
	Count int32  `bson:"count" json:"count"`
}
