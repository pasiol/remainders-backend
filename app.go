package main

import (
	"crypto/tls"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"os"
	"time"
)

type App struct {
	API       *echo.Echo
	Db        *mongo.Database
	Client    *mongo.Client
	TLSConfig *tls.Config
	Debug     bool
}

func (a *App) Initialize() {
	a.API = echo.New()

	a.API.Validator = &CustomValidator{validator: validator.New()}
	a.API.Use(middleware.Logger())
	a.API.Use(middleware.Recover())
	a.API.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 10 * time.Second,
	}))
	err := godotenv.Load()
	if err != nil {
		log.Print("Reading environment failed.")
	}
	a.Debug = GetDebug()
	origins := SplitOrigins()
	a.API.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: origins,
		AllowMethods: []string{http.MethodGet},
	}))
	if a.Debug {
		a.API.Logger.Printf(fmt.Sprintf("CORS: %v", origins))
	}

	a.Db, a.Client, err = a.getDbConnection()
	if err != nil {
		a.API.Logger.Fatal("initializing db connection failed: %s", err)
	}
	a.API.Logger.Printf("database connection succeed db: %s", a.Db.Name())
	a.API.GET("/healthz", a.getHealthz)
	a.API.POST("/login", a.postLogin)

	authorizedEndpoints := a.API.Group("/api/v1")
	config := middleware.JWTConfig{
		Claims:     &jwtCustomClaims{},
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}
	authorizedEndpoints.Use(middleware.JWTWithConfig(config))

	route := authorizedEndpoints.GET("/latest", a.getLatest)
	route.Name = "get-latest"
	route = authorizedEndpoints.GET("/search/:filter", a.getSearch)
	route.Name = "get-search"
}

func (a *App) Run() {
	a.API.Logger.Fatal(a.API.StartTLS(fmt.Sprintf(":%s", os.Getenv("APP_PORT")), os.Getenv("SSL_CERT"), os.Getenv("SSL_KEY")))
}
