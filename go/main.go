package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Estate-CRM/backend-go/internal/cronjob"
	"github.com/Estate-CRM/backend-go/internal/kafka"
	"github.com/robfig/cron/v3"
)

func main() {
	//topic := "contacts-topic"
	c := cron.New()
	c.AddFunc("@every 1m", cronjob.DailyBatchJob)
	c.Start()
	fmt.Println("🕒 Cron job started, running every minute...")
	go func() {
		for {
			log.Println("📤 Sending contact data...")
			kafka.StartProducer()
			time.Sleep(3 * time.Second)
		}
	}()

	select {} // keep the app alive
}
