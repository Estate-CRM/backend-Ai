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

func CreateContact(w http.ResponseWriter, r *http.Request) {
	var contact model.Contact

	

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&contact)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result := config.DB.Create(&contact)
	if result.Error != nil {
		http.Error(w, "Failed to create contact", http.StatusInternalServerError)
		log.Printf("DB error: %v", result.Error)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(contact)
}


func DeleteContact(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid contact ID", http.StatusBadRequest)
		return
	}

	result := config.DB.Delete(&model.Contact{}, id)
	if result.Error != nil {
		http.Error(w, "Failed to delete contact", http.StatusInternalServerError)
		log.Printf("DB error: %v", result.Error)
		return
	}

	if result.RowsAffected == 0 {
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Contact deleted successfully"))
}

func GetContacts(w http.ResponseWriter, r *http.Request){
	var contact []model.Contact
	
	result:=config.DB.Find(&contact)
	if result.Error != nil {
		http.Error(w, "Failed to get contacts", http.StatusInternalServerError)
		log.Printf("DB error: %v", result.Error)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contact)
}