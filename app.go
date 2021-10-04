package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"os"

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
		log.Fatal("Reading environment failed.")
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
	a.Router.HandleFunc("/api/v1/search/{filter}", a.search).Methods("GET")
	a.Router.HandleFunc("/api/v1/latest", a.getLatest).Methods("GET")
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	bytes, err := w.Write(response)
	if err != nil {
		log.Printf("writing response failed: %s", err)
	}
	log.Printf("response bytes %d", bytes)
}

func (a *App) getLatest(w http.ResponseWriter, _ *http.Request) {
	remainders, err := getLatest(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, remainders)
}

func (a *App) search(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf(" %v", vars["filter"])
	remainders, err := searchRecipients(vars["filter"], a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, remainders)
}

func (a *App) Run() {
	corsOptions := cors.New(cors.Options{
		AllowedOrigins:   []string{os.Getenv("ALLOWED_ORIGINS")},
		AllowCredentials: true,
		AllowedMethods:   []string{http.MethodGet, http.MethodOptions},
		Debug:            true,
	})

	server := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: corsOptions.Handler(a.Router),
		TLSConfig: &tls.Config{
			Certificates: readSSLCertificates(),
		},
	}

	log.Print("starting REST-api server")
	log.Fatal(server.ListenAndServeTLS("", ""))
}
