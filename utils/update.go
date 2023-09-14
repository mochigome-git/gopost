package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"post/model"
	"strconv"
	"sync"
	"time"

	"gorm.io/gorm"
)

var mu sync.RWMutex
var stopProcessing = make(chan struct{})
var lastProcessedTime time.Time

const (
	dataCollectionDuration = 35 * time.Second
)

func handleError(err error, message string) {
	if err != nil {
		log.Printf("Error: %s: %v\n", message, err)
	}
}

func ProcessMQTTData(db *gorm.DB) {
	for {
		mu.RLock()
		jsonString := ExportedReceivedMessagesJSON
		mu.RUnlock()

		if jsonString == "" {
			log.Println("JSON string is empty")
			time.Sleep(time.Second)
			continue
		}

		var messages []model.Message

		if err := json.Unmarshal([]byte(jsonString), &messages); err != nil {
			handleError(err, "unmarshaling JSON")
			time.Sleep(time.Second)
			continue
		}

		var existingRecord model.Post
		if err := FindRecordByID(1, &existingRecord, db); err != nil {
			handleError(err, "finding record by ID")
			time.Sleep(time.Second)
			continue
		}

		collectedDataMap := make(map[string][]float64)

		startTime := time.Now()

		for time.Since(startTime) < dataCollectionDuration {
			for _, message := range messages {
				fieldName := message.Address

				floatValue, ok := message.Value.(float64)
				if !ok {
					floatValue, err := strconv.ParseFloat(fmt.Sprintf("%v", message.Value), 64)
					if err != nil {
						log.Printf("Error: message.Value is not a float64: %v\n", message.Value)
						continue
					}
					// Update the value if parsing was successful
					message.Value = floatValue
				}

				collectedDataMap[fieldName] = append(collectedDataMap[fieldName], floatValue)
			}
			time.Sleep(time.Second)
		}

		for fieldName, values := range collectedDataMap {
			if len(values) == 0 {
				continue
			}

			var sum float64
			for _, value := range values {
				sum += value
			}
			mean := sum / float64(len(values))

			if err := UpdateField(&existingRecord, fieldName, mean); err != nil {
				handleError(err, fmt.Sprintf("updating field %s", fieldName))
				continue
			}
		}

		collectedDataMap = make(map[string][]float64)

		if time.Since(startTime) >= dataCollectionDuration {
			if err := UpdateMQTTDataToDB(&existingRecord, db); err != nil {
				handleError(err, "updating database")
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
