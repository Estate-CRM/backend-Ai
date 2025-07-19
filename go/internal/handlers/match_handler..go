package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Estate-CRM/backend-go/internal/config"
	"github.com/Estate-CRM/backend-go/internal/model"
	"github.com/Estate-CRM/backend-go/internal/pkg"
)

type ContractRequest struct {
	ContactID  int `json:"contact_id"`
	ClientID   int `json:"client_id"`
	PropertyID int `json:"property_id"`
}

type MatchHandler struct{}

func (mh *MatchHandler) HandleGenerateContract(w http.ResponseWriter, r *http.Request) {
	var req ContractRequest

	// Parse body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Generating contract for contact #%d (client #%d) and property #%d", req.ContactID, req.ClientID, req.PropertyID)

	var match model.Match
	var property model.Property
	result := config.DB.First(&match, "contact_id = ? AND client_id = ? AND property_id = ?", req.ContactID, req.ClientID, req.PropertyID)
	if result.Error != nil {
		http.Error(w, "Match not found", http.StatusNotFound)
		log.Printf("Match not found: %v", result.Error)
		return
	}
	result = config.DB.First(&property, req.PropertyID)
	if result.Error != nil {
		http.Error(w, "Property not found", http.StatusNotFound)
		log.Printf("Property not found: %v", result.Error)
		return
	}
	contractURL, err := pkg.GenerateContractPDF(req.ContactID, req.ClientID, req.PropertyID, property)
	if err != nil {
		http.Error(w, "Failed to generate contract", http.StatusInternalServerError)
		log.Printf("Error generating contract PDF: %v", err)
		return
	}
	response := map[string]interface{}{
		"message":      "Contract generated successfully",
		"contact_id":   req.ContactID,
		"client_id":    req.ClientID,
		"property_id":  req.PropertyID,
		"contract_url": contractURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func CreateMatch(w http.ResponseWriter, r *http.Request) {

	return
}
