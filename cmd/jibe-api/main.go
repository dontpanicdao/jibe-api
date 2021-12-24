package main

import (
	"net/http"
	"os"
	"time"

	"github.com/dontpanicdao/jibe-api/internal/data"
	"github.com/dontpanicdao/jibe-api/internal/handlers"
	"github.com/gorilla/mux"
)

func main() {
	/*
		SETUP
	*/
	data.InitConfig()

	data.InitDB()
	defer data.CloseDB()

	r := mux.NewRouter()

	// TODO: hand out JWT tokens based on a valid signature
	r.Use(handlers.AuthMiddleware)

	/*
		GETS
	*/
	r.HandleFunc("/v1/subjects", handlers.SubjectsFetch).Methods("GET")
	r.HandleFunc("/v1/subjects/{subject_id}", handlers.SubjectFetch).Methods("GET")
	r.HandleFunc("/v1/subjects/{subject_id}/cert", handlers.CertFetch).Methods("GET")
	r.HandleFunc("/v1/subjects/{subject_id}/phases", handlers.PhasesFetch).Methods("GET")
	r.HandleFunc("/v1/subjects/{subject_id}/phases/{phase_id}", handlers.PhaseFetch).Methods("GET")

	/*
		POSTS
	*/
	r.HandleFunc("/v1/subjects/{subject_id}/cert/key", handlers.CertKey).Methods("POST")

	/*
		INIT
	*/
	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8081",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	log.Info("======================================")
	log.Info("Starting Jibe API")
	log.Info("======================================")
	log.Fatal(srv.ListenAndServe())
}