package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"post/utils"

	//"github.com/joho/godotenv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db             *gorm.DB
	broker         string // broker stores the MQTT broker's hostname.
	mqttport       string // mqttport stores the MQTT broker's port number.
	topic          string // topic stores the topic of the MQTT broker.
	mqttsStr       string // mqtts true or false
	caCertFile     string // CA Certificate location path
	clientCertFile string // Client Certificate location path
	clientKeyFile  string // Client Key location path
)

func main() {
	configureApp()

	stopProcessing := make(chan struct{})
	clientDone := make(chan struct{})
	receivedMessagesJSONChan := make(chan string) // Create a channel for received JSON data

	go utils.Client(broker, mqttport, topic, mqttsStr, caCertFile, clientCertFile, clientKeyFile, receivedMessagesJSONChan, clientDone)

	go func() {
		defer close(stopProcessing)

		for {
			select {
			case <-stopProcessing:
				return
			default:
				utils.ProcessMQTTData(db, receivedMessagesJSONChan, stopProcessing) // Pass the channels
			}
		}
	}()

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	<-sigCh

	// Signal to stop processing
	close(stopProcessing)

	// Wait for client to finish
	<-clientDone
}

func configureApp() {
	//if err := godotenv.Load(); err != nil {
	//	log.Fatalf("Error loading .env file: %v", err)
	//}
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	broker = os.Getenv("MQTT_SUB_HOST")
	mqttsStr = os.Getenv("MQTTS_ON")
	mqttport = os.Getenv("MQTT_SUB_PORT")
	topic = os.Getenv("MQTT_SUB_TOPIC")
	caCertFile = os.Getenv("MQTT_CA_CERTIFICATE")
	clientCertFile = os.Getenv("MQTT_CLIENT_CERTIFICATE")
	clientKeyFile = os.Getenv("MQTT_PRIVATE_KEY")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", host, user, password, dbname, port)
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get DB instance: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
}
