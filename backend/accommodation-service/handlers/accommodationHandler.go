// handlers/accommodation_handler.go
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yourusername/accommodation-service/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccommodationHandler struct {
	accommodationService *services.AccommodationService
}

func NewAccommodationHandler(service *services.AccommodationService) *AccommodationHandler {
	return &AccommodationHandler{
		accommodationService: service,
	}
}

func (h *AccommodationHandler) GetAccommodation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	accommodation, err := h.accommodationService.GetAccommodation(id)
	if err != nil {
		http.Error(w, "Accommodation not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accommodation)
}

func (h *AccommodationHandler) GetAllAccommodations(w http.ResponseWriter, r *http.Request) {
	accommodations := h.accommodationService.GetAllAccommodations()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accommodations)
}

func (h *AccommodationHandler) CreateAccommodation(w http.ResponseWriter, r *http.Request) {
	var accommodation *repositories.Accommodation
	err := json.NewDecoder(r.Body).Decode(&accommodation)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	createdAccommodation, err := h.accommodationService.CreateAccommodation(accommodation)
	if err != nil {
		http.Error(w, "Failed to create accommodation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdAccommodation)
}

func (h *AccommodationHandler) UpdateAccommodation(w http.ResponseWriter, r *http.Request) {
	var accommodation *repositories.Accommodation
	err := json.NewDecoder(r.Body).Decode(&accommodation)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = h.accommodationService.UpdateAccommodation(accommodation)
	if err != nil {
		http.Error(w, "Failed to update accommodation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
