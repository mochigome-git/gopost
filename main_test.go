package main

import (
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"post/utils"
	"syscall"
	"testing"
	"time"
)

func TestMain(*testing.T) {
	go func() {
		http.ListenAndServe("192.168.0.126:16060", nil)
		// http://192.168.0.126:16060/debug/pprof/
	}()
	configureApp()

	stopProcessing := make(chan struct{})
	clientDone := make(chan struct{})
	receivedMessagesJSONChan := make(chan string) // Create a channel for received JSON data

	go utils.Client(broker, mqttport, topic, receivedMessagesJSONChan, clientDone)

	go func() {
		defer close(stopProcessing)

		for {
			select {
			case <-stopProcessing:
				return
			default:
				utils.ProcessMQTTData(db, receivedMessagesJSONChan, stopProcessing) // Pass the channels
			}
			time.Sleep(time.Millisecond * 100)
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
