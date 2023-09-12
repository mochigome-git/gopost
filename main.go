package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"post/utils"

	//"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// broker stores the MQTT broker's hostname.
var broker string

// port stores the MQTT broker's port number.
var mqttport string

// topic of the MQTT broker
var topic string

func main() {
	clientDone := make(chan struct{})
	stopProcessing := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		utils.Client(broker, mqttport, topic)
	}()

	go func() {
		defer close(clientDone)

		for {
			select {
			case <-stopProcessing:
				return
			default:
				utils.ProcessMQTTData(db)
			}
		}
	}()

	time.Sleep(3 * time.Second)

	close(stopProcessing)

	wg.Wait()
}

func init() {
	/*if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}*/
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	broker = os.Getenv("MQTT_HOST")
	mqttport = os.Getenv("MQTT_PORT")
	topic = os.Getenv("MQTT_TOPIC")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", host, user, password, dbname, port)
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("Failed to get DB instance: " + err.Error())
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
}
