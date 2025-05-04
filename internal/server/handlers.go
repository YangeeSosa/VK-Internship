package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"testEx2/api"
	"testEx2/pkg/subpub"
)

type Server struct {
	api.UnimplementedPubSubServer
	pubsub subpub.SubPub
}

func (s *Server) Subscribe(req *api.SubscribeRequest, stream api.PubSub_SubscribeServer) error {
	if req.Key == "" {
		return status.Error(codes.InvalidArgument, "ключ не может быть пустым")
	}

	msgChan := make(chan interface{}, 100)

	sub, err := s.pubsub.Subscribe(req.Key, func(msg interface{}) {
		msgChan <- msg
	})
	if err != nil {
		return status.Error(codes.Internal, "не удалось подписаться")
	}
	defer sub.Unsubscribe()

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

func (s *Server) Publish(ctx context.Context, req *api.PublishRequest) (*emptypb.Empty, error) {
	if req.Key == "" {
		return nil, status.Error(codes.InvalidArgument, "ключ не может быть пустым")
	}

	if err := s.pubsub.Publish(req.Key, req.Data); err != nil {
		return nil, status.Error(codes.Internal, "не удалось опубликовать сообщение")
	}

	return &emptypb.Empty{}, nil
}
