package main

import (
	"log"
	"net/http"

	"github.com/Estate-CRM/backend-go/internal/config"
	"github.com/Estate-CRM/backend-go/internal/routes"
)

func main() {
	config.Connectdb()

	myrouter := routes.Router()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: myrouter,
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	log.Println("Server started on :8080")
}
