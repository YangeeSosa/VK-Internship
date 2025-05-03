package main

import (
	"log"

	"testEx2/config"
	"testEx2/pkg/subpub"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Printf("Не удалось загрузить конфигурацию, используем значения по умолчанию: %v", err)
		cfg = &config.Config{GRPCPort: ":50051"}
	}

	// Создаем новый pubsub
	ps, err := subpub.NewSubPub()
	if err != nil {
		log.Fatalf("Не удалось создать pubsub: %v", err)
	}

	// Создаем и запускаем gRPC сервер
	grpcServer, err := NewGRPCServer(cfg, ps)
	if err != nil {
		log.Fatalf("Не удалось создать gRPC сервер: %v", err)
	}

	if err := grpcServer.Start(); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
