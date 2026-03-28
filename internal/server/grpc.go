package server

import (
	"context"

	httpstreamv1 "github.com/OpenProjectX/http-stream/api/httpstream/v1"
	"github.com/OpenProjectX/http-stream/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	httpstreamv1.UnimplementedStreamServiceServer
	streamer *service.Streamer
}

func New(streamer *service.Streamer) *GRPCServer {
	return &GRPCServer{streamer: streamer}
}

func (s *GRPCServer) Transfer(ctx context.Context, req *httpstreamv1.TransferRequest) (*httpstreamv1.TransferResponse, error) {
	resp, err := s.streamer.Transfer(ctx, req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return resp, nil
}
