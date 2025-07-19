package kafka

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Estate-CRM/backend-go/internal/model"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func StartProducer() {
	// Create Kafka producer
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",
	})
	if err != nil {
		log.Fatalf("❌ Failed to create producer: %s", err)
	}
	defer p.Close()

	// Handle delivery reports
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("❌ Delivery failed: %v\n", ev.TopicPartition.Error)
				} else {
					log.Printf("✅ Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	topic := "contacts-topic"

	for {
		// Simulate contact data
		contacts := []model.Contact{
			{ID: 1, ClientID: 101, Latitude: 36.75, Longitude: 3.06},
			{ID: 2, ClientID: 102, Latitude: 35.69, Longitude: -0.63},
		}

		// Convert to JSON
		valueBytes, err := json.Marshal(contacts)
		if err != nil {
			log.Printf("❌ Failed to marshal contacts: %v\n", err)
			continue
		}

		// Send message
		err = p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &topic,
				Partition: int32(kafka.PartitionAny),
			},
			Value: valueBytes,
			Key:   []byte("page-1"),
		}, nil)

		if err != nil {
			log.Printf("❌ Failed to send message: %v\n", err)
		}

		time.Sleep(3 * time.Second)
	}
}
