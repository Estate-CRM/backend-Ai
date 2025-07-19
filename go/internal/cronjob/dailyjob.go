package cronjob

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

// This job simulates a daily batch upload of property data to the AI backend.
// The total set of properties (e.g. 12,000) is divided across 6 days of the week.
// Each day, a specific slice of properties is processed based on how many days have passed since the last full CSV upload.
// This approach ensures balanced daily processing and avoids overloading the system.

const totalProperties = 12000
const numDays = 6
const batchSize = totalProperties / numDays
const maxConcurrentRequests = 10

// Simulate CSV upload time (set to last Sunday)
var lastCSVUpload = time.Now().AddDate(0, 0, -1) // Yesterday

// DailyBatchJob processes a daily batch of properties
func DailyBatchJob() {
	log.Println("🔁 Daily batch job started...")

	// Calculate how many days since the last upload
	daysSinceUpload := int(time.Since(lastCSVUpload).Hours() / 24)
	if daysSinceUpload >= numDays {
		log.Println("✅ All daily batches for this week are completed.")
		return
	}

	start := daysSinceUpload * batchSize
	end := start + batchSize
	log.Printf("📦 Processing properties from %d to %d...\n", start, end)

	// Simulate the batch of properties with IDs
	properties := mockProperties(start, end)

	// Use concurrency with semaphore
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrentRequests)

	for _, prop := range properties {
		wg.Add(1)
		sem <- struct{}{}

		go func(p Property) {
			defer wg.Done()
			defer func() { <-sem }()
			processProperty(p)
		}(prop)
	}

	wg.Wait()
	log.Println("✅ Daily batch job finished.")

}

// Property represents a mocked property
type Property struct {
	ID    int
	Title string
}

// mockProperties generates fake property data
func mockProperties(start, end int) []Property {
	props := make([]Property, 0, end-start)
	for i := start; i < end; i++ {
		props = append(props, Property{
			ID:    i,
			Title: fmt.Sprintf("Property #%d", i),
		})
	}
	return props
}

// processProperty simulates sending property to AI backend
func processProperty(p Property) {
	// Simulate delay
	time.Sleep(time.Duration(rand.Intn(200)+100) * time.Millisecond)
	log.Printf("📩 Sent %s to AI backend.\n", p.Title)
}
