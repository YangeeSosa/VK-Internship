package server

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
	Server *grpc.Server
	SubPub subpub.SubPub
	Lis    net.Listener
}

func NewGRPCServer(cfg *config.Config, ps subpub.SubPub) (*GRPCServer, error) {
	grpcServer := grpc.NewServer()
	api.RegisterPubSubServer(grpcServer, &Server{pubsub: ps})

	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		return nil, err
	}

	return &GRPCServer{
		Server: grpcServer,
		SubPub: ps,
		Lis:    lis,
	}, nil
}

func (s *GRPCServer) Start() error {
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down gRPC server...")
		s.Server.GracefulStop()

		if err := s.SubPub.Close(context.Background()); err != nil {
			log.Printf("Error closing pubsub: %v", err)
		}
	}()

	log.Printf("Starting gRPC server on %s", s.Lis.Addr().String())
	return s.Server.Serve(s.Lis)
}

func (s *GRPCServer) Stop() {
	s.Server.GracefulStop()
}
