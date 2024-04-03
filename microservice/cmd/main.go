package main

import (
	"github.com/yigithankarabulut/asyncs3todbloader/microservice/apiserver"
	"github.com/yigithankarabulut/asyncs3todbloader/microservice/config"
	"log"
)

func main() {
	cnf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	if err := apiserver.New(
		apiserver.WithConfig(cnf),
		apiserver.WithServerEnv("development"),
	); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
