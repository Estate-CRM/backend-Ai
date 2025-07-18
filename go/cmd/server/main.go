package main

import (
	"log"
	"net/http"

	"github.com/Estate-CRM/backend-go/internal/config"
	"github.com/Estate-CRM/backend-go/internal/kafka"

	// assumes ExportContactsToCSV is here
	// assumes SendCSVChunksToKafka is here

	// assumes GetAllContacts is here
	"github.com/Estate-CRM/backend-go/internal/routes"
)

func main() {
	config.Connectdb()

	//go scheduleWeeklyExport()

	pageSize := 100
	topic := "contacts-topic"
	kafka.ProduceContactsPaginated(pageSize,topic)
	
	myrouter := routes.Router()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: myrouter,
	}

	log.Println("Server started on :8080")

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}