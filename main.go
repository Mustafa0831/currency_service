package main

import (
	"currency_service/handlers"
	"currency_service/utils"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	config, err := utils.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	db, err := utils.ConnectDB(config.DBConnection)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	handlers.InitHandlers(db)

	router := mux.NewRouter()
	router.HandleFunc("/currency/save/{date}", handlers.SaveCurrency).Methods("GET")
	router.HandleFunc("/currency/{date}/{code}", handlers.GetCurrency).Methods("GET")

	log.Printf("Server listening on port %s", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, router))
}
