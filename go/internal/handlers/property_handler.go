package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Estate-CRM/backend-go/internal/config"
	"github.com/Estate-CRM/backend-go/internal/model"
	"github.com/go-chi/chi"
)

func CreateProperty(w http.ResponseWriter, r *http.Request) {
	var property model.Property

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&property)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result := config.DB.Create(&property)
	if result.Error != nil {
		http.Error(w, "Failed to create property", http.StatusInternalServerError)
		log.Printf("DB error: %v", result.Error)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(property)
}

func DeleteProperty(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid property ID", http.StatusBadRequest)
		return
	}

	result := config.DB.Delete(&model.Property{}, id)
	if result.Error != nil {
		http.Error(w, "Failed to delete property", http.StatusInternalServerError)
		log.Printf("DB error: %v", result.Error)
		return
	}

	if result.RowsAffected == 0 {
		http.Error(w, "Property not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Property deleted successfully"))
}

func GetProperties(w http.ResponseWriter, r *http.Request) {
	var property []model.Property

	result := config.DB.Find(&property)
	if result.Error != nil {
		http.Error(w, "Failed to get properties", http.StatusInternalServerError)
		log.Printf("DB error: %v", result.Error)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(property)
}
