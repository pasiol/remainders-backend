package main

import (
	"context"
	"errors"
	"net/url"
	"os"
	"regexp"
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

func connectOrFail(uri string) (*mongo.Database, *mongo.Client, error) {

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, nil, err
	}
	err = client.Connect(context.Background())
	if err != nil {
		return nil, nil, err
	}
	db, err := getDbName(uri)
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

func getDbName(uri string) (string, error) {

	beginning, err := regexp.Compile(`(mongodb(\+srv|)://(\S*):(\S*)@(\S*)/)`)
	if err != nil {
		return "", err
	}

	match := beginning.FindAllString(uri, -1)
	if len(match) > 0 {

		dbName := strings.Replace(uri, match[0], "", 1)
		ending, err := regexp.Compile(`\?(\S*)`)
		if err != nil {
			return "", err
		}
		match = ending.FindAllString(dbName, -1)
		if len(match) > 0 {
			dbName = strings.Replace(string(dbName), match[0], "", 1)
		}

		return string(dbName), nil
	}
	return "", errors.New("cannot extract db name")
}
