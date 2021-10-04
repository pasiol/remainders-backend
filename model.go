package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type remainder struct {
	To        string    `bson:"to" json:"to"`
	Title     string    `bson:"title" json:"title"`
	Message   string    `bson:"message" json:"message"`
	Type      string    `bson:"type" json:"type"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type User struct {
	UserId    primitive.ObjectID `bson:"_id json":"userId"`
	Username  string             `bson:"userName" json:"userName"`
	Password  string             `bson:"password" json:"password"`
	Email     string             `bson:"email" json:"email"`
	Approved  bool               `bson:"approved" json:"approved"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

/*func searchRecipients(searchPhrase string, db *mongo.Database) ([]remainder, error) {
	var remainders []remainder

	log.Printf("Trying to find recipients filter: %s", searchPhrase)
	options := options.Find()
	options.SetSort(bson.D{{"updated_at", -1}})
	options.SetLimit(200)

	cursor, err := db.Collection("sended").Find(context.TODO(), bson.D{{"to", bson.D{{"$regex", searchPhrase}, {"$options", "im"}}}}, options)
	if err != nil {
		return []remainder{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Printf("closing cursor failed: %s", err)
		}
	}(cursor, context.TODO())
	for cursor.Next(context.TODO()) {
		var currentRemainder remainder
		if err = cursor.Decode(&currentRemainder); err != nil {
			log.Printf("decoding remainder failed: err")
		}
		transformedRemainders, err := transformRemainder(currentRemainder)
		if err != nil {
			return []remainder{}, err
		}
		remainders = append(remainders, transformedRemainders...)
	}
	return remainders, nil
}*/

func getLatest(db *mongo.Database) ([]remainder, error) {
	var remainders []remainder

	options := options.Find()
	options.SetSort(bson.D{{"updated_at", -1}})
	options.SetLimit(25)

	cursor, err := db.Collection("sended").Find(context.TODO(), bson.D{{}}, options)
	if err != nil {
		return []remainder{}, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Printf("closing cursor failed: %s", err)
		}
	}(cursor, context.TODO())
	for cursor.Next(context.TODO()) {
		var currentRemainder remainder
		if err = cursor.Decode(&currentRemainder); err != nil {
			log.Printf("decoding user failed: err")
		}
		remainders = append(remainders, currentRemainder)
	}
	return remainders, nil

}
