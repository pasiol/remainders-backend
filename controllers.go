package main

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (u User) createUser(db *mongo.Database) error {

	// TODO validate user data

	hashedPassword, err := hashAndSalt(u.Password)
	if err != nil {
		return err
	}
	user := bson.D{{"_id", u.Username}, {"username", u.Username}, {"password", hashedPassword}, {"email", u.Email}, {"approved", false}, {"createAt", time.Now()}, {"updatedAt", time.Now()}}
	_, err = db.Collection("users").InsertOne(context.TODO(), user)
	if err != nil {
		if strings.Contains(err.Error(), "E11000 duplicate key error") {
			return errors.New("username already exists")
		}
		return err
	}
	return nil

}

func (u User) login(db *mongo.Database) bool {
	var user User
	filter := bson.D{{"username", u.Username}, {"approved", true}}
	queryOptions := options.FindOne()
	queryOptions.SetProjection(bson.D{{"password", 1}, {"username", 1}, {"_id", 0}})
	result := db.Collection("users").FindOne(context.TODO(), filter, queryOptions)

	if err := result.Decode(&user); err != nil {
		log.Printf("decoding user failed: %s", err)
		return false
	}
	if checkPasswordHash(u.Password, user.Password) {
		return true
	}
	return false
}
