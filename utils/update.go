package utils

import (
	"encoding/json"
	"fmt"
	"post/model"
	"strconv"
	"sync"
	"time"

	"gorm.io/gorm"
)

var mu sync.RWMutex
var stopProcessing = make(chan struct{})
var lastProcessedTime time.Time

func ProcessMQTTData(db *gorm.DB) {
	for {
		mu.RLock()
		jsonString := ExportedReceivedMessagesJSON
		mu.RUnlock()

		if jsonString == "" {
			fmt.Println("JSON string is empty")
			time.Sleep(time.Second)
			continue
		}

		var messages []model.Message

		if err := json.Unmarshal([]byte(jsonString), &messages); err != nil {
			fmt.Printf("Error unmarshaling JSON: %v\n", err)
			time.Sleep(time.Second)
			continue
		}

		var existingRecord model.Post
		if err := FindRecordByID(1, &existingRecord, db); err != nil {
			fmt.Printf("Error finding record by ID: %v\n", err)
			time.Sleep(time.Second)
			continue
		}

		// Collect data for 35 seconds
		startTime := time.Now()
		collectedData := make(map[string][]float64) // Map to store data for each fieldName as float64

		for time.Since(startTime).Seconds() < 35 {
			for _, message := range messages {
				fieldName := message.Address

				// Check if the Value is a float64
				fieldValue, ok := message.Value.(float64)
				if !ok {
					// Attempt to convert to float64
					if floatValue, err := strconv.ParseFloat(fmt.Sprintf("%v", message.Value), 64); err == nil {
						fieldValue = floatValue
					} else {
						fmt.Printf("Error: message.Value is not a float64: %v\n", message.Value)
						continue
					}
				}

				// Append the fieldValue to the map for the corresponding fieldName
				collectedData[fieldName] = append(collectedData[fieldName], fieldValue)
			}
			time.Sleep(time.Second)
		}

		// Calculate the mean for each fieldName
		for fieldName, values := range collectedData {
			if len(values) == 0 {
				continue
			}

			var sum float64
			for _, value := range values {
				sum += value
			}
			mean := sum / float64(len(values))

			// Call UpdateField with the calculated mean
			if err := UpdateField(&existingRecord, fieldName, mean); err != nil {
				fmt.Printf("Error updating field %s: %v\n", fieldName, err)
				continue
			}
		}

		// Clear the collected data
		collectedData = make(map[string][]float64)

		// Update the database after 35 seconds
		if time.Since(startTime).Seconds() >= 35 {
			if err := UpdateMQTTDataToDB(&existingRecord, db); err != nil {
				fmt.Printf("Error updating database: %v\n", err)
			}
			lastProcessedTime = time.Now()
		}

		select {
		case <-stopProcessing:
			return
		default:
			continue
		}
	}
}

// To stop the goroutine, you can close the stopProcessing channel:
func StopProcessing() {
	close(stopProcessing)
}
