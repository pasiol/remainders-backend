package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

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

func (a *App) getSearch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf(" %v", vars["filter"])
	remainders, err := searchRecipients(vars["filter"], a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, remainders)
}

func (a *App) postUser(w http.ResponseWriter, r *http.Request) {
	var u User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "malformed user object")
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("closing request body failed: %s", err)
		}
	}(r.Body)

	if err := u.createUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, "new user created")
}

func (a *App) postLogin(w http.ResponseWriter, r *http.Request) {
	var u User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "malformed user object")
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("closing request body failed: %s", err)
		}
	}(r.Body)
	if u.login(a.DB) {
		token, err := getToken(u.Username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Generating token failed: " + err.Error()))
			return
		} else {
			w.Header().Set("Authorization", "Bearer "+token)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Token: " + token))
			return
		}
	}
	respondWithJSON(w, http.StatusNetworkAuthenticationRequired, "user authentication failed")
}
