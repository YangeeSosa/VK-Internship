package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"testEx2/api"
	"testEx2/config"
	"testEx2/pkg/subpub"
)

type GRPCServer struct {
	server *grpc.Server
	pubsub subpub.SubPub
	lis    net.Listener
}

func NewGRPCServer(cfg *config.Config, ps subpub.SubPub) (*GRPCServer, error) {
	grpcServer := grpc.NewServer()
	api.RegisterPubSubServer(grpcServer, &server{pubsub: ps})

	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		return nil, err
	}

	return &GRPCServer{
		server: grpcServer,
		pubsub: ps,
		lis:    lis,
	}, nil
}

func (s *GRPCServer) Start() error {
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Завершаем работу gRPC сервера...")
		s.server.GracefulStop()

		if err := s.pubsub.Close(context.Background()); err != nil {
			log.Printf("Ошибка при закрытии pubsub: %v", err)
		}
	}()

	log.Printf("Запускаем gRPC сервер на %s", s.lis.Addr().String())
	return s.server.Serve(s.lis)
}

func (s *GRPCServer) Stop() {
	s.server.GracefulStop()
}
