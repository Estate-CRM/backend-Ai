package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Estate-CRM/backend-go/internal/config"
	"github.com/Estate-CRM/backend-go/internal/model"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func fetchContactsPage(limit, offset int) ([]model.Contact, error) {
	var page []model.Contact
	result := config.DB.Limit(limit).Offset(offset).Find(&page)
	if result.Error != nil {
		return nil, result.Error
	}
	return page, nil
}

func ProduceContactsPaginated(pageSize int, topic string) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		log.Fatalf("❌ Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("❌ Delivery failed: %v\n", ev.TopicPartition.Error)
				} else {
					log.Printf("✅ Message delivered to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	offset := 0
	pageNumber := 1

	for {
		contacts, err := fetchContactsPage(pageSize, offset)
		if err != nil {
			log.Printf("❌ DB error at offset %d: %v", offset, err)
			break
		}

		if len(contacts) == 0 {
			log.Println("✅ All contacts processed.")
			break
		}

		data, err := json.Marshal(contacts)
		if err != nil {
			log.Printf("❌ JSON marshal failed at page %d: %v", pageNumber, err)
			break
		}
	
		fmt.Printf("Value (JSON): %s\n", string(data))
	
		err = producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: int32(kafka.PartitionAny)},
			Value:          data,
			Key:            []byte(fmt.Sprintf("page-%d", pageNumber)),
		}, nil)

		if err != nil {
			log.Printf("❌ Failed to send page %d: %v", pageNumber, err)
			break
		}

		log.Printf("🚀 Sent page %d with %d contacts", pageNumber, len(contacts))

		offset += pageSize
		pageNumber++
		time.Sleep(300 * time.Millisecond) // Optional: throttle sending
	}

	producer.Flush(5000)
}
