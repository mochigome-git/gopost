package utils

import (
	"encoding/json"
	"fmt"
	"post/model"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// CollectedData holds data for each fieldName
type CollectedData map[string][]float64

// ProcessMQTTData processes MQTT data and updates the database.
func ProcessMQTTData(db *gorm.DB, receivedMessagesJSONChan <-chan string, stopProcessing <-chan struct{}) {
	var collectedData CollectedData
	collectedData = make(CollectedData)

	for {
		select {
		case jsonString := <-receivedMessagesJSONChan:
			if jsonString == "" {
				fmt.Println("JSON string is empty")
				continue
			}

			var messages []model.Message
			if err := json.Unmarshal([]byte(jsonString), &messages); err != nil {
				fmt.Printf("Error unmarshaling JSON: %v\n", err)
				continue
			}

			var existingRecord model.Post
			if err := FindRecordByID(1, &existingRecord, db); err != nil {
				fmt.Printf("Error finding record by ID: %v\n", err)
				continue
			}

			startTime := time.Now()

			for {
				for _, message := range messages {
					fieldName := message.Address
					fieldValue, ok := convertToFloat64(message.Value)
					if !ok {
						fmt.Printf("Error: message.Value is not a float64: %v\n", message.Value)
						continue
					}
					collectedData[fieldName] = append(collectedData[fieldName], fieldValue)
				}
				time.Sleep(time.Second)

				if time.Since(startTime).Seconds() >= 35 {
					break
				}
			}

			for fieldName, values := range collectedData {
				if len(values) == 0 {
					continue
				}

				clearCacheAndData(collectedData)

				mean := calculateMean(values)
				if err := UpdateField(&existingRecord, fieldName, mean); err != nil {
					fmt.Printf("Error updating field %s: %v\n", fieldName, err)
					continue
				}
			}

			if err := UpdateMQTTDataToDB(&existingRecord, db); err != nil {
				fmt.Printf("Error updating database: %v\n", err)
			}
			return

		case <-stopProcessing:
			return
		}
	}
}

// convertToFloat64 converts a value to float64.
func convertToFloat64(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err == nil {
			return f, true
		}
	}
	return 0, false
}

// calculateMean calculates the mean of a slice of float64 values.
func calculateMean(values []float64) float64 {
	var sum float64
	for _, value := range values {
		sum += value
	}
	return sum / float64(len(values))
}

// Define a function to clear the cache and data.
func clearCacheAndData(collectedData CollectedData) CollectedData {
	// Create a new empty map to replace the existing one.
	return make(CollectedData)
}
