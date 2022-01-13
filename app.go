package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"gitlab.com/pasiol/mongoUtils"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	Router    *mux.Router
	DB        *mongo.Database
	Client    *mongo.Client
	TLSConfig *tls.Config
}

func (a *App) Initialize() {

	err := godotenv.Load()
	if err != nil {
		log.Print("Reading environment failed.")
	}

	mongoConfig := mongoUtils.MongoConfig{
		Db:       os.Getenv("DB"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("PASSWORD"),
		URI:      os.Getenv("URI"),
	}
	a.DB, a.Client, err = mongoUtils.ConnectOrFail(mongoConfig, false)
	if err != nil {
		log.Fatalf("database connection error: %s", err)
	}
	log.Printf("connected to db: %s", a.DB.Name())

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/api/healthz", a.getHealth).Methods("GET")
	a.Router.HandleFunc("/api/v1/search/{filter}", authorizeRequest(a.getSearch)).Methods("GET")
	a.Router.HandleFunc("/api/v1/latest", authorizeRequest(a.getLatest)).Methods("GET")
	a.Router.HandleFunc("/api/v1/user", authorizeRequest(a.postUser)).Methods("POST")
	a.Router.HandleFunc("/api/v1/login", a.postLogin).Methods("POST")

}

func (a *App) Run() {
	corsOptions := cors.New(cors.Options{
		AllowedHeaders:   []string{"X-Requested-With", "Content-Type", "Authorization"},
		AllowedOrigins:   []string{os.Getenv("ALLOWED_ORIGINS")},
		AllowCredentials: true,
		AllowedMethods:   []string{http.MethodGet, http.MethodOptions, http.MethodConnect, http.MethodPost},
		Debug:            true,
	})

	server := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: corsOptions.Handler(a.Router),
		TLSConfig: &tls.Config{
			Certificates: readSSLCertificates(),
		},
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("starting REST-server :%s.", os.Getenv("PORT"))
	log.Printf("Version: %s , build: %s", Version, Build)
	log.Fatal(server.ListenAndServeTLS("", ""))
}
