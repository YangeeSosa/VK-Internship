package main

import (
	"log"

	"testEx2/config"
	grpcserver "testEx2/internal/server"
	"testEx2/pkg/subpub"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Printf("Failed to load configuration, using default values: %v", err)
		cfg = &config.Config{GRPCPort: ":50051"}
	}

	ps, err := subpub.NewSubPub()
	if err != nil {
		log.Fatalf("Failed to create pubsub: %v", err)
	}

	grpcServer, err := grpcserver.NewGRPCServer(cfg, ps)
	if err != nil {
		log.Fatalf("Failed to create gRPC server: %v", err)
	}

	if err := grpcServer.Start(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
