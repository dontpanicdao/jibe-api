package main

import (
	"log"
	"net/http"
	"time"

	"github.com/dontpanicdao/jibe-api/internal/data"
	"github.com/dontpanicdao/jibe-api/internal/handlers"
	"github.com/gorilla/mux"
)

const ALPHA int = 1

func main() {
	/*
		SETUP
	*/
	data.InitConfig()
	
	data.InitStarkCuve()

	data.InitDB()
	defer data.CloseDB()

	err := data.InitTypes(ALPHA)
	if err != nil {
		log.Println("ERR Initializing Typed Cert: ", err)
	}

	r := mux.NewRouter()

	// TODO: hand out JWT tokens based on a valid signature
	a := handlers.NewAuth()

	r.Use(a.AuthMiddleware)

	/*
		GETS
	*/
	r.HandleFunc("/v1/elements", handlers.ElementsFetch).Methods("GET")
	r.HandleFunc("/v1/elements/{element_id}", handlers.ElementFetch).Methods("GET")
	r.HandleFunc("/v1/elements/{element_id}/custom", handlers.FetchCustomCert).Methods("GET")
	r.HandleFunc("/v1/elements/{element_id}/protons", handlers.ProtonsFetch).Methods("GET")
	r.HandleFunc("/v1/elements/{element_id}/protons/{proton_id}", handlers.ProtonFetch).Methods("GET")

	/*
		POSTS
	*/
	r.HandleFunc("/v1/elements", handlers.CreateElement).Methods("POST")
	r.HandleFunc("/v1/elements/{element_id}/cert", handlers.ProposeCert).Methods("POST")
	r.HandleFunc("/v1/elements/{element_id}/cert/attempt", handlers.GradeCert).Methods("POST")
	r.HandleFunc("/v1/elements/{element_id}/custom", handlers.CreateCustomCert).Methods("POST")
	r.HandleFunc("/v1/elements/{element_id}/protons", handlers.AddProton).Methods("POST")

	/*
		INIT
	*/
	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8081",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	log.Println("======================================")
	log.Println("Starting Jibe API")
	log.Println("======================================")
	log.Fatal(srv.ListenAndServe())
}
