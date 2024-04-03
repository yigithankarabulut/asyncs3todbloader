package main

import (
	"github.com/yigithankarabulut/asyncs3todbloader/job/app"
	"github.com/yigithankarabulut/asyncs3todbloader/job/config"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	// load configuration.
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	doneChan := make(chan struct{})
	shutdownChan := make(chan os.Signal, 2)
	signal.Notify(shutdownChan, os.Interrupt, os.Kill)
	go func() {
		if err := app.New(
			app.WithConfig(cfg),
			app.WithLogLevel("INFO"),
			app.WithDoneChan(doneChan),
		); err != nil {
			log.Fatalf("failed to create app: %v", err)
		}
	}()
	// graceful shutdown. listen for shutdown signal.
	select {
	case <-shutdownChan:
		log.Println("Shutdown signal received. Shutting down...")
		log.Println("Waiting 5 seconds for goroutines to complete the job...")
		time.Sleep(5 * time.Second)
	case <-doneChan:
		log.Println("All goroutines completed the job.")
	}
	log.Println("Exiting...")
}
