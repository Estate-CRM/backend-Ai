package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Estate-CRM/backend-go/internal/config"
	"github.com/Estate-CRM/backend-go/internal/cronjob"
	"github.com/Estate-CRM/backend-go/internal/kafka"
	"github.com/Estate-CRM/backend-go/internal/routes"
	"github.com/robfig/cron/v3"
)

func main() {
	// Start cron job scheduler
	c := cron.New()
	// Schedule to run daily at 00:00 from Monday to Saturday
	_, err := c.AddFunc("0 0 * * 1-6", cronjob.DailyBatchJob)
	if err != nil {
		log.Fatalf("❌ Failed to add cron job: %v", err)
	}
	c.AddFunc("0 0 * * 0", func() {
		log.Println("📤 [Weekly] Sending contact data...")
		kafka.StartProducer()
	})

	c.Start()
	fmt.Println("🕒 Cron job started. Scheduled every midnight (Mon-Sat)")

	// Kafka producer loop
	go func() {
		for {
			log.Println("📤 Sending contact data...")
			kafka.StartProducer()
			time.Sleep(3 * time.Second)
		}
	}()

	// Connect to the database
	config.Connectdb()

	// Setup router
	myrouter := routes.Router()

	// Start HTTP server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: myrouter,
	}

	log.Println("🚀 Server started on :8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("❌ Error starting server: %v", err)
	}
}
