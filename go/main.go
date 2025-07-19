package main

import (
	"log"
	"time"

	"github.com/Estate-CRM/backend-go/internal/kafka"
)

func main() {
	//topic := "contacts-topic"

	go func() {
		for {
			log.Println("📤 Sending contact data...")
			kafka.StartProducer()
			time.Sleep(3 * time.Second)
		}
	}()

	select {} // keep the app alive
}
