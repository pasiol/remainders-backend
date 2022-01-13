package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type remainder struct {
	To        string    `bson:"to" json:"to"`
	Title     string    `bson:"title" json:"title"`
	Message   string    `bson:"message" json:"message"`
	Type      string    `bson:"type" json:"type"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type User struct {
	UserId    primitive.ObjectID `bson:"_id" json:"userId"`
	Username  string             `bson:"username" json:"username"`
	Password  string             `bson:"password" json:"password"`
	Email     string             `bson:"email" json:"email"`
	Approved  bool               `bson:"approved" json:"approved"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type Exception struct {
	Message string `json:"message"`
}

type Payload struct {
	Message string `json:"message"`
}
