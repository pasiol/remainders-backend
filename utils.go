package main

import (
	"context"
	"errors"
	"net/url"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"golang.org/x/crypto/bcrypt"
)

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func SplitOrigins() ([]string, error) {
	s, exists := os.LookupEnv("ALLOWED_ORIGINS")
	if !exists {
		return []string{}, errors.New("ALLOWED_ORIGINS variable missing")
	}
	origins := strings.Split(s, ",")

	for _, origin := range origins {
		uri, err := url.ParseRequestURI(origin)
		if err == nil && (uri.Scheme != "https" && uri.Scheme != "http") {
			return []string{}, errors.New("malformed uri")
		}
	}
	return origins, nil
}

func GetDebug() bool {
	return os.Getenv("APP_DEBUG") == "true"
}

func connectOrFail(uri string, db string) (*mongo.Database, *mongo.Client, error) {

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, nil, err
	}
	err = client.Connect(context.Background())
	if err != nil {
		return nil, nil, err
	}
	var DB = client.Database(db)
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, nil, err
	}

	return DB, client, nil
}
