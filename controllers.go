package main

import (
	"context"
	"log"
	"os"
	"time"

	mongodb "gitlab.com/pasiol/mongoUtils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (a *App) getDbConnection() (*mongo.Database, *mongo.Client, error) {
	var err error
	var db *mongo.Database
	var client *mongo.Client
	for i := 1; i <= 10; i++ {
		m := mongodb.MongoConfig{
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("PASSWORD"),
			Db:       os.Getenv("APP_DB"),
			URI:      os.Getenv("APP_DB_URI"),
		}
		db, client, err = mongodb.ConnectOrFail(m, false)
		if err == nil {
			break
		}
		a.API.Logger.Printf("connecting to database failed, iteration: %d, err: %s", i, err)
		time.Sleep(10 * time.Second)
	}
	return db, client, err
}

func searchRecipients(searchPhrase string, db *mongo.Database) ([]Remainder, error) {
	var remainders []Remainder

	log.Printf("Trying to find recipients filter: %s", searchPhrase)
	queryOptions := options.Find()
	queryOptions.SetSort(bson.D{{"updated_at", -1}})
	queryOptions.SetLimit(200)

	cursor, err := db.Collection("sended").Find(context.TODO(), bson.D{{"to", bson.D{{"$regex", searchPhrase}, {"$options", "im"}}}}, queryOptions)
	if err != nil {
		return []Remainder{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Printf("closing cursor failed: %s", err)
		}
	}(cursor, context.TODO())
	for cursor.Next(context.TODO()) {
		var currentRemainder Remainder
		if err = cursor.Decode(&currentRemainder); err != nil {
			log.Printf("decoding remainder failed: err")
		}

		remainders = append(remainders, currentRemainder)
	}
	return remainders, nil
}

func find(db *mongo.Database) ([]Remainder, error) {
	var remainders []Remainder

	queryOptions := options.Find()
	queryOptions.SetSort(bson.D{{"updated_at", -1}})
	queryOptions.SetLimit(25)

	cursor, err := db.Collection("sended").Find(context.TODO(), bson.D{{}}, queryOptions)
	if err != nil {
		return []Remainder{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Printf("closing cursor failed: %s", err)
		}
	}(cursor, context.TODO())
	for cursor.Next(context.TODO()) {
		var currentRemainder Remainder
		if err = cursor.Decode(&currentRemainder); err != nil {
			log.Printf("decoding remainder failed: %s", err.Error())
		}
		remainders = append(remainders, currentRemainder)
	}
	return remainders, nil
}

func (u *User) Login(db *mongo.Database) bool {
	var user User
	filter := bson.D{{"username", u.Username}, {"approved", true}}
	queryOptions := options.FindOne()
	queryOptions.SetProjection(bson.D{{"password", 1}, {"username", 1}, {"_id", 0}})
	result := db.Collection("users").FindOne(context.Background(), filter, queryOptions)

	if err := result.Decode(&user); err != nil {
		return false
	}
	if checkPasswordHash(u.Password, user.Password) {
		return true
	}
	return false
}
