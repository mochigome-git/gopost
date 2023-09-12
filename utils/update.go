package utils

import (
	"encoding/json"
	"fmt"
	"post/model"
	"sync"
	"time"

	"gorm.io/gorm"
)

var mu sync.RWMutex
var stopProcessing = make(chan struct{})

func ProcessMQTTData(db *gorm.DB) {
	for {
		mu.RLock()
		jsonString := ExportedReceivedMessagesJSON
		mu.RUnlock()

		if jsonString == "" {
			fmt.Println("JSON string is empty")
			time.Sleep(time.Second)
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
			continue
		}

		for _, message := range messages {
			fieldName := message.Address
			fieldValue := message.Value

			if err := UpdateField(&existingRecord, fieldName, fieldValue); err != nil {
				fmt.Printf("Error updating field %s: %v\n", fieldName, err)
				continue
			}
		}

		if err := UpdateMQTTDataToDB(&existingRecord, db); err != nil {
			fmt.Printf("Error updating database: %v\n", err)
			continue
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
