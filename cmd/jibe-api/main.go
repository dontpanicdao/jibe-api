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
	r.HandleFunc("/v1/facts/job/{job_id}", handlers.CheckFactJobStatus).Methods("GET")
	r.HandleFunc("/v1/facts/{fact}", handlers.CheckFactStatus).Methods("GET")
	r.HandleFunc("/v1/facts/{fact}/l1_tx/{tx_hash}", handlers.ShipSharpStark).Methods("GET")
	r.HandleFunc("/v1/user/{public_key}", handlers.UserFetch).Methods("GET")

	/*
		POSTS
	*/
	r.HandleFunc("/v1/elements", handlers.CreateElement).Methods("POST")
	r.HandleFunc("/v1/elements/{element_id}/cert", handlers.ProposeCert).Methods("POST")
	r.HandleFunc("/v1/elements/{element_id}/cert/rubric", handlers.AddRubric).Methods("POST")
	r.HandleFunc("/v1/elements/{element_id}/cert/attempt", handlers.GradeCert).Methods("POST")
	r.HandleFunc("/v1/elements/{element_id}/custom", handlers.CreateCustomCert).Methods("POST")
	r.HandleFunc("/v1/elements/{element_id}/protons", handlers.AddProton).Methods("POST")


	/*
		WEBAUTHN
	*/
	r.HandleFunc("/v1/registration/init", handlers.BeginRegistration).Methods("POST")
	r.HandleFunc("/v1/registration/end", handlers.FinishRegistration).Methods("POST")
	r.HandleFunc("/v1/login/init/{public_key}", handlers.BeginLogin).Methods("GET")
	r.HandleFunc("/v1/login/end", handlers.FinishLogin).Methods("POST")

	/*
		INIT
	*/
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8081",
		WriteTimeout: 25 * time.Second,
		ReadTimeout:  25 * time.Second,
	}

	log.Println("======================================")
	log.Println("Starting Jibe API")
	log.Println("======================================")
	log.Fatal(srv.ListenAndServe())
}
