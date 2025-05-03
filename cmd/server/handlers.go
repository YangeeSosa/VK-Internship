package main

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"testEx2/api"
	"testEx2/pkg/subpub"
)

type server struct {
	api.UnimplementedPubSubServer
	pubsub subpub.SubPub
}

func (s *server) Subscribe(req *api.SubscribeRequest, stream api.PubSub_SubscribeServer) error {
	if req.Key == "" {
		return status.Error(codes.InvalidArgument, "ключ не может быть пустым")
	}

	// Создание канала для получения сообщений
	msgChan := make(chan interface{}, 100)

	// Подписка на систему pubsub
	sub, err := s.pubsub.Subscribe(req.Key, func(msg interface{}) {
		msgChan <- msg
	})
	if err != nil {
		return status.Error(codes.Internal, "не удалось подписаться")
	}
	defer sub.Unsubscribe()

	// Обработка входящих сообщений
	for {
		select {
		case msg := <-msgChan:
			if str, ok := msg.(string); ok {
				if err := stream.Send(&api.Event{Data: str}); err != nil {
					return status.Error(codes.Internal, "не удалось отправить событие")
				}
			}
		case <-stream.Context().Done():
			return nil
		}
	}
}

func (s *server) Publish(ctx context.Context, req *api.PublishRequest) (*emptypb.Empty, error) {
	if req.Key == "" {
		return nil, status.Error(codes.InvalidArgument, "ключ не может быть пустым")
	}

	if err := s.pubsub.Publish(req.Key, req.Data); err != nil {
		return nil, status.Error(codes.Internal, "не удалось опубликовать сообщение")
	}

	return &emptypb.Empty{}, nil
}
