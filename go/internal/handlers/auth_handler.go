package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Estate-CRM/backend-go/internal/config"
	"github.com/Estate-CRM/backend-go/internal/middlewares"
	"github.com/Estate-CRM/backend-go/internal/model"
	"github.com/Estate-CRM/backend-go/internal/pkg"
)

type AuthHandler struct {

}

func (auth *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	//use model

	var user model.User
	result := config.DB.First(&user, "email = ? ", req.Email)
	if result.Error != nil {
		http.Error(w, "Invalid email", http.StatusUnauthorized)
		return
	}
	hashPass := pkg.PasswordUtils{}
	isValid, err := hashPass.VerifyPassword(user.Password, req.Password)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if !isValid {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	role := user.Role
	accessToken, err := middlewares.SignAccessToken(user.Email, role)
	if err != nil {
		http.Error(w, "Failed to create access token", http.StatusInternalServerError)
		return
	}
	refreshToken, err := middlewares.SignRefreshToken(user.Email, role)
	if err != nil {
		http.Error(w, "Failed to create refresh token", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message":      "Login successful",
		"user_id":      fmt.Sprint(user.ID),
		"role":         role,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})

}

func (auth *AuthHandler) RegisterClient(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Phone     string `json:"phone_number"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var existing model.User
	if result := config.DB.First(&existing, "email = ?", req.Email); result.RowsAffected > 0 {
		http.Error(w, "Email already in use", http.StatusConflict)
		return
	}
	hashPass := pkg.PasswordUtils{}
	hashedPassword, err := hashPass.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	req.Password = hashedPassword
	user := model.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Email:     req.Email,
		Password:  req.Password,
		Role:      "client",
	}
	if result := config.DB.Create(&user); result.Error != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}
	client := model.Client{
		UserID: user.ID,
	}
	if err := config.DB.Create(&client).Error; err != nil {
		http.Error(w, "Could not create client", http.StatusInternalServerError)
		return
	}
	accessToken, err := middlewares.SignAccessToken(user.Email, user.Role)
	if err != nil {
		http.Error(w, "Failed to create access token", http.StatusInternalServerError)
		return
	}
	refreshToken, err := middlewares.SignRefreshToken(user.Email, user.Role)
	if err != nil {
		http.Error(w, "Failed to create refresh token", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":      "User registered successfully",
		"user_id":      fmt.Sprint(user.ID),
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})

}

func (auth *AuthHandler) RegisterAgent(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form with 10MB max memory
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Read text fields
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	phone := r.FormValue("phone_number")
	email := r.FormValue("email")
	password := r.FormValue("password")
	hashPass := pkg.PasswordUtils{}
	hashedPassword, err := hashPass.HashPassword(password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	password = hashedPassword
	user := model.User{
		FirstName: firstName,
		LastName:  lastName,
		Phone:     phone,
		Email:     email,
		Password:  password,
		Role:      "agent",
	}
	result := config.DB.Create(&user)
	if result.Error != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Create agent model
	agent := model.Agent{
		UserID:             user.ID,
		NationalID:         r.FormValue("national_id"),
		CommercialRegister: "",
		Verified:           false, // Default to false, can be updated later
	}
	result2 := config.DB.Create(&agent)
	if result2.Error != nil {
		http.Error(w, "Failed to create agent", http.StatusInternalServerError)
		return
	}
	// Handle file uploads
	nationalIDFile, nationalIDHeader, err := r.FormFile("national_id")
	if err != nil {
		http.Error(w, "National ID image is required", http.StatusBadRequest)
		return
	}
	defer nationalIDFile.Close()

	commercialRegFile, commercialRegHeader, err := r.FormFile("commercial_register")
	if err != nil {
		http.Error(w, "Commercial Register PDF is required", http.StatusBadRequest)
		return
	}
	defer commercialRegFile.Close()
	accessToken, err := middlewares.SignAccessToken(user.Email, user.Role)
	if err != nil {
		http.Error(w, "Failed to create access token", http.StatusInternalServerError)
		return
	}
	refreshToken, err := middlewares.SignRefreshToken(user.Email, user.Role)
	if err != nil {
		http.Error(w, "Failed to create refresh token", http.StatusInternalServerError)
		return
	}
	// Respond with extracted info (without saving)
	response := map[string]interface{}{
		"message": "Received agent data successfully",
		"user": map[string]string{
			"first_name":   firstName,
			"last_name":    lastName,
			"phone":        phone,
			"email":        email,
			"password":     password,
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		},
		"files": map[string]interface{}{
			"national_id": map[string]interface{}{
				"filename": nationalIDHeader.Filename,
				"size":     nationalIDHeader.Size,
				"mime":     nationalIDHeader.Header.Get("Content-Type"),
			},
			"commercial_register": map[string]interface{}{
				"filename": commercialRegHeader.Filename,
				"size":     commercialRegHeader.Size,
				"mime":     commercialRegHeader.Header.Get("Content-Type"),
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (auth *AuthHandler) Testdata(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	// Split the "Bearer <token>"
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		http.Error(w, "Invalid Authorization format", http.StatusUnauthorized)
		return
	}

	tokenStr := parts[1]

	claims, err := middlewares.VerifyAccessToken(tokenStr)
	if err != nil {
		http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Print extracted claims
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Token is valid!\nEmail: %s\nRole: %s\n", claims.Email, claims.Role)
}
