package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/golang-jwt/jwt/v4"

	"golang.org/x/crypto/bcrypt"
)

func readSSLCertificates() []tls.Certificate {

	cert, err := tls.LoadX509KeyPair(os.Getenv("SSL_PUBLIC"), os.Getenv("SSL_PRIVATE"))
	if err != nil {
		log.Fatalf("reading ssl certificates failed: %s", err)
	}
	return []tls.Certificate{cert}
}

func getTSLConfig(cert string) *tls.Config {

	pool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(cert)
	if err != nil {
		log.Fatalf("reading certificate file failed: %s", err)
	}
	pool.AppendCertsFromPEM(ca)

	tslConfig := &tls.Config{
		RootCAs:            pool,
		InsecureSkipVerify: true,
	}

	return tslConfig
}

func hashAndSalt(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getToken(name string) (string, error) {
	signingKey := []byte(os.Getenv("JWT_SECRET"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"username": name})
	tokenString, err := token.SignedString(signingKey)
	return tokenString, err
}

func verifyToken(tokenString string) (jwt.Claims, error) {
	signingKey := []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	if err != nil {
		return nil, err
	}
	return token.Claims, err
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
	if os.Getenv("APP_DEBUG") == "true" {
		return true
	}
	return false
}

func generateMongoURI(m MongoConfig, connectDirect bool) (string, error) {
	if m.URI == "" || m.Db == "" {
		return "", errors.New("malformed mongo config struct")
	}
	var credentials, parameters string
	if m.User != "" && m.Password != "" {
		credentials = fmt.Sprintf("%s:%s@", m.User, m.Password)
	}
	if connectDirect {
		parameters = "?connect=direct"
		return fmt.Sprintf("mongodb://%s%s/%s%s", credentials, m.URI, m.Db, parameters), nil
	}
	parameters = "?retryWrites=true&w=majority"
	return fmt.Sprintf("mongodb+srv://%s%s/%s%s", credentials, m.URI, m.Db, parameters), nil
}

func ConnectOrFail(m MongoConfig, connectDirect bool) (a *mongo.Database, b *mongo.Client, err error) {
	mongoURI, err := generateMongoURI(m, connectDirect)
	if err != nil {
		return &mongo.Database{}, &mongo.Client{}, err
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return &mongo.Database{}, &mongo.Client{}, err
	}
	err = client.Connect(context.TODO())
	var DB = client.Database(m.Db)
	if err != nil {
		return &mongo.Database{}, &mongo.Client{}, err
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return &mongo.Database{}, &mongo.Client{}, err
	}
	return DB, client, nil
}
