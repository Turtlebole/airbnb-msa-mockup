package main

import (
	"accommodation-service/handlers"
	"accommodation-service/repositories"
	"accommodation-service/services"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize repositories
	accommodationRepo := repositories.NewAccommodationRepository()

	// Initialize services
	accommodationService := services.NewAccommodationService(accommodationRepo)

	// Initialize handlers
	accommodationHandler := handlers.NewAccommodationHandler(accommodationService)

	// Setup routes
	r := mux.NewRouter()
	r.HandleFunc("/accommodation", accommodationHandler.CreateAccommodation).Methods("POST")
	r.HandleFunc("/accommodation/{id}", accommodationHandler.GetAccommodation).Methods("GET")
	r.HandleFunc("/accommodation", accommodationHandler.GetAllAccommodations).Methods("GET")
	r.HandleFunc("/accommodation/{id}", accommodationHandler.UpdateAccommodation).Methods("PUT")

	// Start server
	port := 8080
	log.Printf("Server is running on port %d...\n", port)
	log.Fatal(http.ListenAndServe(":8080", r))
}
